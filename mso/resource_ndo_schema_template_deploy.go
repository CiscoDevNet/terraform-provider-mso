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
		}),
	}
}

func resourceNDOSchemaTemplateDeployExecute(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Template Deploy Execution", d.Id())
	templateId, templateIdProvided := d.GetOk("template_id")
	templateName, templateNameProvided := d.GetOk("template_name")
	templateType := d.Get("template_type").(string)
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
		var resolvedTemplateId string
		if templateIdProvided {
			resolvedTemplateId = templateId.(string)
		} else {
			if !templateNameProvided {
				return fmt.Errorf("when 'template_id' is not provided, 'template_name' must be set for template_type %s", templateType)
			}
			templateId, err := GetTemplateIdByNameAndType(msoClient, templateName.(string), templateType)
			if err != nil {
				return err
			}
			resolvedTemplateId = templateId
			if err := d.Set("template_id", resolvedTemplateId); err != nil {
				return fmt.Errorf("error setting resolved template_id in state: %w", err)
			}
		}
		payloadStr = fmt.Sprintf(`{"templateId": "%s", "isRedeploy": %v}`, resolvedTemplateId, d.Get("re_deploy").(bool))
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
	return nil
}
