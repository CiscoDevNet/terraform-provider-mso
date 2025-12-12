package mso

import (
	"errors"
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceNDOSchemaTemplateDeploy() *schema.Resource {
	return &schema.Resource{
		Create: resourceNDOSchemaTemplateDeployExecute,
		Read:   resourceNDOSchemaTemplateDeployRead,
		Update: resourceNDOSchemaTemplateDeployExecute,
		Delete: resourceNDOSchemaTemplateDeployDelete,

		SchemaVersion: version,

		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			// Plan time validation.
			msoClient := v.(*client.Client)
			if msoClient.GetPlatform() != "nd" {
				return errors.New(`The 'mso_schema_template_deploy_ndo' resource is only supported for nd based platforms, 'platform=nd' must be configured in the provider section of your configuration.`)
			}
			return nil
		},

		Schema: (map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"template_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"template_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "application",
				ValidateFunc: validation.StringInSlice([]string{
					"application",
					"tenant",
					"l3out",
					"fabric_policy",
					"fabric_resource",
					"monitoring_tenant",
					"monitoring_access",
					"service_device",
				}, false),
			},

			"template_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"re_deploy": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"force_apply": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "always-deploy",
			},

			"site_ids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of site IDs where template is deployed. If not provided, will be retrieved from API during undeploy",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"undeploy": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set to true to undeploy the template without destroying the resource. Only supported for non-application template types",
			},

			"undeploy_on_destroy": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Undeploys the template when the Terraform resource is destroyed. Only supported for non-application template types",
			},
		}),
	}
}

func resourceNDOSchemaTemplateDeployExecute(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Template Deploy Execution", d.Id())

	undeploy := d.Get("undeploy").(bool)

	if undeploy {
		return executeTemplateUndeploy(d, m)
	}

	templateType := d.Get("template_type").(string)
	_, templateIdProvided := d.GetOk("template_id")
	templateName, templateNameProvided := d.GetOk("template_name")
	path := "api/v1/task"

	msoClient := m.(*client.Client)
	var payloadStr string

	if templateType == "application" && !templateIdProvided {
		schemaId, schemaIdProvided := d.GetOk("schema_id")

		if !schemaIdProvided || !templateNameProvided {
			return fmt.Errorf("when 'template_id' is not provided, both 'schema_id' and 'template_name' must be set for template_type %s", templateType)
		}
		schemaValidate := models.SchemValidate{SchmaId: d.Get("schema_id").(string)}
		_, err := msoClient.ReadSchemaValidate(&schemaValidate)
		if err != nil {
			return err
		}
		payloadStr = fmt.Sprintf(`{"schemaId": "%s", "templateName": "%s", "isRedeploy": %v}`, schemaId.(string), templateName.(string), d.Get("re_deploy").(bool))
	} else {
		if !templateIdProvided {
			if !templateNameProvided {
				return fmt.Errorf("when 'template_id' is not provided, 'template_name' must be set for template_type %s", templateType)
			}
			templateId, err := GetTemplateIdByNameAndType(msoClient, templateName.(string), templateType)
			if err != nil {
				return err
			}
			if err := d.Set("template_id", templateId); err != nil {
				return fmt.Errorf("error setting resolved template_id in state: %w", err)
			}
		}
		payloadStr = fmt.Sprintf(`{"templateId": "%s", "isRedeploy": %v}`, d.Get("template_id").(string), d.Get("re_deploy").(bool))
	}

	payload, err := container.ParseJSON([]byte(payloadStr))
	if err != nil {
		log.Printf("[DEBUG] Parse of JSON failed with err: %s.", err)
		return err
	}

	req, err := msoClient.MakeRestRequest("POST", path, payload, true)
	if err != nil {
		log.Printf("[DEBUG] MakeRestRequest failed with err: %s.", err)
		return err
	}
	cont, resp, err := msoClient.Do(req)
	if resp.StatusCode != 202 || err != nil {
		log.Printf("[DEBUG] Request failed with resp: %v. Err: %s.", resp, err)
		return err
	}

	taskId, ok := cont.S("id").Data().(string)
	if !ok || taskId == "" {
		log.Printf("[DEBUG] Task ID not found or is invalid. Data was: %v", cont.S("id").Data())
		return fmt.Errorf("task ID not found or is invalid")
	}

	schemaId, ok := cont.S("reqDetails", "schemaId").Data().(string)
	if !ok || schemaId == "" {
		log.Printf("[DEBUG] Schema ID not found or is invalid. Data was: %v", cont.S("reqDetails", "schemaId").Data())
		return fmt.Errorf("schema ID not found or is invalid")
	}

	req, err = msoClient.MakeRestRequest("GET", fmt.Sprintf("%s/%s", path, taskId), nil, true)
	if err != nil {
		log.Printf("[DEBUG] MakeRestRequest failed with err: %s.", err)
		return err
	}

	cont, resp, err = msoClient.DoWithRetryFunc(req, isTaskStatusPending)
	if err != nil && cont == nil {
		log.Printf("[DEBUG] Request failed with resp: %v. Err: %s.", resp, err)
		return err
	} else if cont != nil {
		taskStatusContainer := cont.S("operDetails", "taskStatus")
		if taskStatusContainer != nil {
			if status, ok := taskStatusContainer.Data().(string); ok && status == "Error" {
				errorMessage := "Could not determine specific deployment error message."
				firstErrorMessageContainer := cont.S("operDetails", "detailedStatus", "errMessage").Index(0)
				if message, ok := firstErrorMessageContainer.Data().(string); ok {
					errorMessage = message
				}
				return fmt.Errorf("Error on deploy: %s", errorMessage)
			}
		}
	}
	d.SetId(schemaId)
	log.Printf("[DEBUG] %s: Successful Template Deploy Execution", d.Id())
	return resourceNDOSchemaTemplateDeployRead(d, m)
}

func resourceNDOSchemaTemplateDeployRead(d *schema.ResourceData, m interface{}) error {
	d.Set("force_apply", "")
	return nil
}

func resourceNDOSchemaTemplateDeployDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Resource Deletion", d.Id())

	undeployOnDestroy := d.Get("undeploy_on_destroy").(bool)

	if undeployOnDestroy {
		return executeTemplateUndeploy(d, m)
	}
	return nil
}

func executeTemplateUndeploy(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Executing template undeploy")

	// Get template_id from state which is always set for non-application templates during deployment
	templateId, ok := d.GetOk("template_id")
	if !ok || templateId.(string) == "" {
		return fmt.Errorf("template_id not found in state. Resource must be deployed first before undeploying")
	}

	// Get schema_id from resource ID which is always set during deployment
	schemaId := d.Id()
	if schemaId == "" {
		return fmt.Errorf("schema_id which is the resource ID not found in state. Resource must be deployed first before undeploying")
	}

	log.Printf("[DEBUG] Using template_id from state: %s", templateId.(string))
	log.Printf("[DEBUG] Using schema_id from state: %s", schemaId)

	msoClient := m.(*client.Client)
	var err error
	var siteIds []string

	if siteIdsRaw, ok := d.GetOk("site_ids"); ok && len(siteIdsRaw.([]interface{})) > 0 {
		// User provided site IDs
		siteIdsList := siteIdsRaw.([]interface{})
		for _, siteId := range siteIdsList {
			siteIds = append(siteIds, siteId.(string))
		}
		log.Printf("[DEBUG] Using user-provided site_ids: %v", siteIds)
	} else {
		// Retrieve site IDs from API
		siteIds, err = GetDeployedSiteIds(msoClient, templateId.(string))
		if err != nil {
			return fmt.Errorf("failed to retrieve deployed site_ids: %w", err)
		}
		log.Printf("[DEBUG] Retrieved site_ids from API: %v", siteIds)
	}

	if len(siteIds) == 0 {
		return fmt.Errorf("no site IDs available for undeploy operation")
	}

	// Construct undeploy payload
	payloadStr := fmt.Sprintf(`{"schemaId": "%s", "templateId": "%s", "undeploy": [`, schemaId, templateId.(string))
	for i, siteId := range siteIds {
		if i > 0 {
			payloadStr += ","
		}
		payloadStr += fmt.Sprintf(`"%s"`, siteId)
	}
	payloadStr += "]}"

	log.Printf("[DEBUG] Undeploy payload: %s", payloadStr)

	payload, err := container.ParseJSON([]byte(payloadStr))
	if err != nil {
		return fmt.Errorf("failed to parse undeploy payload: %w", err)
	}

	path := "api/v1/task"
	req, err := msoClient.MakeRestRequest("POST", path, payload, true)
	if err != nil {
		return fmt.Errorf("failed to create undeploy request: %w", err)
	}

	cont, resp, err := msoClient.Do(req)
	if err != nil {
		return fmt.Errorf("undeploy request failed: %w", err)
	}

	if resp.StatusCode != 202 {
		return fmt.Errorf("unexpected status code %d for undeploy", resp.StatusCode)
	}

	// Wait for undeploy task completion
	taskId, ok := cont.S("id").Data().(string)
	if !ok || taskId == "" {
		return fmt.Errorf("task ID not found in undeploy response")
	}

	log.Printf("[DEBUG] Undeploy task ID: %s", taskId)

	req, err = msoClient.MakeRestRequest("GET", fmt.Sprintf("api/v1/task/%s", taskId), nil, true)
	if err != nil {
		return fmt.Errorf("failed to check undeploy task status: %w", err)
	}

	cont, resp, err = msoClient.DoWithRetryFunc(req, isTaskStatusPending)
	if err != nil && cont == nil {
		return fmt.Errorf("undeploy task monitoring failed: %w", err)
	}

	if cont != nil {
		taskStatusContainer := cont.S("operDetails", "taskStatus")
		if taskStatusContainer != nil {
			if status, ok := taskStatusContainer.Data().(string); ok && status == "Error" {
				errorMessage := "Could not determine specific undeploy error message."
				firstErrorMessageContainer := cont.S("operDetails", "detailedStatus", "errMessage").Index(0)
				if message, ok := firstErrorMessageContainer.Data().(string); ok {
					errorMessage = message
				}
				return fmt.Errorf("error on undeploy: %s", errorMessage)
			}
		}
	}

	log.Printf("[DEBUG] %s: Successful Template Undeploy", d.Id())
	return resourceNDOSchemaTemplateDeployRead(d, m)
}
