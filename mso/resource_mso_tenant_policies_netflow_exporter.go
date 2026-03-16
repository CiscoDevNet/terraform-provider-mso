package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSONetflowExporter() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSONetflowExporterCreate,
		Read:   resourceMSONetflowExporterRead,
		Update: resourceMSONetflowExporterUpdate,
		Delete: resourceMSONetflowExporterDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSONetflowExporterImport,
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the tenant policy template.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
				Description:  "The name of the NetFlow Exporter.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the NetFlow Exporter.",
			},
		},
	}
}

func setNetflowExporterData(d *schema.ResourceData, response *container.Container, templateId string) error {
	d.SetId(fmt.Sprintf("templateId/%s/NetflowExporter/%s", templateId, models.StripQuotes(response.S("name").String())))
	d.Set("template_id", templateId)
	d.Set("name", models.StripQuotes(response.S("name").String()))
	d.Set("uuid", models.StripQuotes(response.S("uuid").String()))

	return nil
}

func resourceMSONetflowExporterImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO NetFlow Exporter Resource - Beginning Import: %v", d.Id())
	resourceMSONetflowExporterRead(d, m)
	log.Printf("[DEBUG] MSO NetFlow Exporter Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSONetflowExporterCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO NetFlow Exporter Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	payload := map[string]interface{}{}

	payload["name"] = d.Get("name").(string)

	payloadModel := models.GetPatchPayload("add", "/tenantPolicyTemplate/template/netFlowExporters/-", payload)
	templateId := d.Get("template_id").(string)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/NetflowExporter/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO NetFlow Exporter Resource - Create Complete: %v", d.Id())
	return resourceMSONetflowExporterRead(d, m)
}

func resourceMSONetflowExporterRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO NetFlow Exporter Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return err
	}

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyName, err := GetPolicyNameFromResourceId(d.Id(), "NetflowExporter")
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "netFlowExporters")
	if err != nil {
		return err
	}

	setNetflowExporterData(d, policy, templateId)
	log.Printf("[DEBUG] MSO NetFlow Exporter Resource - Read Complete: %v", d.Id())
	return nil
}

func resourceMSONetflowExporterUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO NetFlow Exporter Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "netFlowExporters")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/tenantPolicyTemplate/template/netFlowExporters/%d", policyIndex)

	payloadCont := container.New()
	payloadCont.Array()

	if d.HasChange("name") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/name", updatePath), d.Get("name").(string))
		if err != nil {
			return err
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/NetflowExporter/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO NetFlow Exporter Resource - Update Complete: %v", d.Id())
	return resourceMSONetflowExporterRead(d, m)
}

func resourceMSONetflowExporterDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO NetFlow Exporter Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "netFlowExporters")
	if err != nil {
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/tenantPolicyTemplate/template/netFlowExporters/%d", policyIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO NetFlow Exporter Resource - Delete Complete: %v", d.Id())
	return nil
}