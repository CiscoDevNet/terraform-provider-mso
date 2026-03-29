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

func resourceMSORouteMapPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSORouteMapPolicyCreate,
		Read:   resourceMSORouteMapPolicyRead,
		Update: resourceMSORouteMapPolicyUpdate,
		Delete: resourceMSORouteMapPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSORouteMapPolicyImport,
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the tenant policy template.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
				Description:  "The name of the Route Map Policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the Route Map Policy.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the Route Map Policy.",
			},
		},
	}
}

func setRouteMapPolicyData(d *schema.ResourceData, response *container.Container, templateId string) error {
	d.SetId(fmt.Sprintf("templateId/%s/RouteMapPolicy/%s", templateId, models.StripQuotes(response.S("name").String())))
	d.Set("template_id", templateId)
	d.Set("name", models.StripQuotes(response.S("name").String()))
	d.Set("description", models.StripQuotes(response.S("description").String()))
	d.Set("uuid", models.StripQuotes(response.S("uuid").String()))

	return nil
}

func resourceMSORouteMapPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO Route Map Policy Resource - Beginning Import: %v", d.Id())
	err := resourceMSORouteMapPolicyRead(d, m)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] MSO Route Map Policy Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSORouteMapPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Route Map Policy Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	payload := map[string]interface{}{
		"name": d.Get("name").(string),
	}

	if description, ok := d.GetOk("description"); ok {
		payload["description"] = description.(string)
	}

	payloadModel := models.GetPatchPayload("add", "/tenantPolicyTemplate/template/routeMapPolicies/-", payload)
	templateId := d.Get("template_id").(string)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/RouteMapPolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO Route Map Policy Resource - Create Complete: %v", d.Id())
	return resourceMSORouteMapPolicyRead(d, m)
}

func resourceMSORouteMapPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Route Map Policy Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return err
	}

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyName, err := GetPolicyNameFromResourceId(d.Id(), "RouteMapPolicy")
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "routeMapPolicies")
	if err != nil {
		return err
	}

	setRouteMapPolicyData(d, policy, templateId)
	log.Printf("[DEBUG] MSO Route Map Policy Resource - Read Complete: %v", d.Id())
	return nil
}

func resourceMSORouteMapPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Route Map Policy Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "routeMapPolicies")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/tenantPolicyTemplate/template/routeMapPolicies/%d", policyIndex)

	payloadCont := container.New()
	payloadCont.Array()

	if d.HasChange("name") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/name", updatePath), d.Get("name").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("description") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/description", updatePath), d.Get("description").(string))
		if err != nil {
			return err
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/RouteMapPolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO Route Map Policy Resource - Update Complete: %v", d.Id())
	return resourceMSORouteMapPolicyRead(d, m)
}

func resourceMSORouteMapPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Route Map Policy Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "routeMapPolicies")
	if err != nil {
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/tenantPolicyTemplate/template/routeMapPolicies/%d", policyIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO Route Map Policy Resource - Delete Complete: %v", d.Id())
	return nil
}
