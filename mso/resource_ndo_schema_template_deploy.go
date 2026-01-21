package mso

import (
	"errors"
	"fmt"
	"log"
	"strings"

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

			// Validate undeploy configuration
			undeploy := diff.Get("undeploy").(bool)
			siteIds := diff.Get("site_ids").([]interface{})

			if undeploy && len(siteIds) == 0 {
				return errors.New("when 'undeploy=true', 'site_ids' must be provided. To undeploy the template from all associated sites, set undeploy_on_destroy=true, apply that change to save the value of undeploy_on_destroy to state and then destroy the resource")
			}
			return nil
		},

		Schema: (map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"template_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
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
				Description: "List of site IDs where template is deployed.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"undeploy": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set to true to undeploy the template from select sites provided in 'site_ids' without destroying the resource.",
			},

			"undeploy_on_destroy": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Undeploys the template from all the sites before the Terraform resource is destroyed.",
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

	var templateInfo *TemplateInfo
	if templateIdProvided {
		idOfTemplate := d.Get("template_id").(string)

		info, err := GetTemplateInfo(msoClient, idOfTemplate, "", "")
		if err != nil {
			return fmt.Errorf("failed to retrieve template info: %w", err)
		}
		templateInfo = info

		if templateInfo.TemplateType != templateType {
			return fmt.Errorf("template_type mismatch: template_id '%s' is associated with template_type '%s', but you provided template_type '%s'. Please change template_type to '%s'",
				idOfTemplate, templateInfo.TemplateType, templateType, templateInfo.TemplateType)
		}

		if !templateNameProvided {
			if err := d.Set("template_name", templateInfo.TemplateName); err != nil {
				return fmt.Errorf("error setting template_name in state: %w", err)
			}
			templateName = templateInfo.TemplateName
			templateNameProvided = true
		}

		if d.Get("schema_id").(string) == "" {
			if err := d.Set("schema_id", templateInfo.SchemaId); err != nil {
				return fmt.Errorf("error setting schema_id in state: %w", err)
			}
		}
	}

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

			info, err := GetTemplateInfo(msoClient, "", templateName.(string), templateType)
			if err != nil {
				return err
			}
			templateInfo = info

			if err := d.Set("template_id", templateInfo.TemplateId); err != nil {
				return fmt.Errorf("error setting resolved template_id in state: %w", err)
			}

			if d.Get("schema_id").(string) == "" {
				if err := d.Set("schema_id", templateInfo.SchemaId); err != nil {
					return fmt.Errorf("error setting schema_id in state: %w", err)
				}
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

	d.SetId(fmt.Sprintf("%s/template/%s", schemaId, d.Get("template_name").(string)))
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

	templateType := d.Get("template_type").(string)
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	if schemaId == "" {
		return fmt.Errorf("schema_id not found in state. Resource must be deployed first before undeploying")
	}

	var err error
	var siteIds []string
	var payloadStr string
	var templateInfo *TemplateInfo

	if siteIdsRaw, ok := d.GetOk("site_ids"); ok && len(siteIdsRaw.([]interface{})) > 0 {
		siteIdsList := siteIdsRaw.([]interface{})
		for _, siteId := range siteIdsList {
			siteIds = append(siteIds, siteId.(string))
		}
		log.Printf("[DEBUG] Using user-provided site_ids: %v", siteIds)
	}

	if templateType == "application" {
		templateName := d.Get("template_name").(string)
		if templateName == "" {
			return fmt.Errorf("template_name not found in state for application template undeploy")
		}

		if len(siteIds) == 0 {
			siteIds, err = GetDeployedSiteIdsForApplicationTemplate(msoClient, schemaId, templateName)
			if err != nil {
				return fmt.Errorf("failed to retrieve deployed site_ids for application template: %w", err)
			}
			log.Printf("[DEBUG] Retrieved site_ids from API: %v", siteIds)
		}

		if len(siteIds) == 0 {
			return fmt.Errorf("no site IDs available for undeploy operation")
		}

		payloadStr = fmt.Sprintf(`{"schemaId":"%s","templateName":"%s","undeploy":[`, schemaId, templateName)
	} else {
		templateId, ok := d.GetOk("template_id")
		if !ok || templateId.(string) == "" {
			return fmt.Errorf("template_id not found in state. Resource must be deployed first before undeploying")
		}

		if len(siteIds) == 0 {
			templateInfo, err = GetTemplateInfo(msoClient, templateId.(string), "", "")
			if err != nil {
				return fmt.Errorf("failed to retrieve template info: %w", err)
			}

			if templateInfo.TemplateStatus != "DEPLOYMENT_SUCCESSFUL" {
				return fmt.Errorf("template is not in DEPLOYMENT_SUCCESSFUL status (current: %s)", templateInfo.TemplateStatus)
			}

			if len(templateInfo.DeployedSiteIds) == 0 {
				return fmt.Errorf("no successfully deployed sites found for template")
			}

			siteIds = templateInfo.DeployedSiteIds
			log.Printf("[DEBUG] Retrieved site_ids from API: %v", siteIds)
		}

		if len(siteIds) == 0 {
			return fmt.Errorf("no site IDs available for undeploy operation")
		}

		payloadStr = fmt.Sprintf(`{"schemaId":"%s","templateId":"%s","undeploy":[`, schemaId, templateId.(string))
	}

	quotedSiteIds := make([]string, len(siteIds))
	for i, siteId := range siteIds {
		quotedSiteIds[i] = fmt.Sprintf(`"%s"`, siteId)
	}
	payloadStr += strings.Join(quotedSiteIds, ",") + "]}"

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
