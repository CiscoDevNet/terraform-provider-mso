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

func resourceMSOPhysicalDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOPhysicalDomainCreate,
		Read:   resourceMSOPhysicalDomainRead,
		Update: resourceMSOPhysicalDomainUpdate,
		Delete: resourceMSOPhysicalDomainDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOPhysicalDomainImport,
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vlan_pool_uuid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func setPhysicalDomainData(d *schema.ResourceData, msoClient *client.Client, templateId, policyName string) error {

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "fabricPolicyTemplate", "template", "domains")
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/physicalDomain/%s", templateId, models.StripQuotes(policy.S("name").String())))
	d.Set("template_id", templateId)
	d.Set("name", models.StripQuotes(policy.S("name").String()))
	d.Set("description", models.StripQuotes(policy.S("description").String()))
	d.Set("uuid", models.StripQuotes(policy.S("uuid").String()))
	d.Set("vlan_pool_uuid", models.StripQuotes(policy.S("pool").String()))

	return nil
}

func resourceMSOPhysicalDomainImport(d *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO Physical Domain Resource - Beginning Import: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return nil, err
	}

	policyName, err := GetPolicyNameFromResourceId(d.Id(), "physicalDomain")
	if err != nil {
		return nil, err
	}

	setPhysicalDomainData(d, msoClient, templateId, policyName)
	log.Printf("[DEBUG] MSO Physical Domain Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOPhysicalDomainCreate(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO Physical Domain Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	payload := map[string]any{}

	payload["name"] = d.Get("name").(string)

	if description, ok := d.GetOk("description"); ok {
		payload["description"] = description.(string)
	}

	if vlan_pool_uuid, ok := d.GetOk("vlan_pool_uuid"); ok {
		payload["pool"] = vlan_pool_uuid.(string)
	}

	payloadModel := models.GetPatchPayload("add", "/fabricPolicyTemplate/template/domains/-", payload)
	templateId := d.Get("template_id").(string)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/physicalDomain/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO Physical Domain Resource - Create Complete: %v", d.Id())
	return resourceMSOPhysicalDomainRead(d, m)
}

func resourceMSOPhysicalDomainRead(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO Physical Domain Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	setPhysicalDomainData(d, msoClient, templateId, policyName)
	log.Printf("[DEBUG] MSO Physical Domain Resource - Read Complete : %v", d.Id())
	return nil
}

func resourceMSOPhysicalDomainUpdate(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO Physical Domain Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "fabricPolicyTemplate", "template", "domains")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/fabricPolicyTemplate/template/domains/%d", policyIndex)

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

	if d.HasChange("vlan_pool_uuid") {
		vlanPool := d.Get("vlan_pool_uuid").(string)
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/pool", updatePath), vlanPool)
		if err != nil {
			return err
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/physicalDomain/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO Physical Domain Resource - Update Complete: %v", d.Id())
	return resourceMSOPhysicalDomainRead(d, m)
}

func resourceMSOPhysicalDomainDelete(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO Physical Domain Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", d.Get("template_id").(string)))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "fabricPolicyTemplate", "template", "domains")
	if err != nil {
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/fabricPolicyTemplate/template/domains/%d", policyIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", d.Get("template_id").(string)), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO Physical Domain Resource - Delete Complete: %v", d.Id())
	return nil
}
