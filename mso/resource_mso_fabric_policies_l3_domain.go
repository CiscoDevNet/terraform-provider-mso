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

func resourceMSOL3Domain() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOL3DomainCreate,
		Read:   resourceMSOL3DomainRead,
		Update: resourceMSOL3DomainUpdate,
		Delete: resourceMSOL3DomainDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOL3DomainImport,
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the fabric policy template.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
				Description:  "The name of the L3 Domain.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the L3 Domain.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the L3 Domain.",
			},
			"vlan_pool_uuid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The UUID of the VLAN Pool. Providing an empty string will remove the VLAN pool reference.",
			},
		},
	}
}

func setL3DomainData(d *schema.ResourceData, response *container.Container, templateId string) error {
	d.SetId(fmt.Sprintf("templateId/%s/L3Domain/%s", templateId, models.StripQuotes(response.S("name").String())))
	d.Set("template_id", templateId)
	d.Set("name", models.StripQuotes(response.S("name").String()))
	d.Set("description", models.StripQuotes(response.S("description").String()))
	d.Set("uuid", models.StripQuotes(response.S("uuid").String()))

	poolValue := models.StripQuotes(response.S("pool").String())
	if poolValue == "{}" || poolValue == "" {
		d.Set("vlan_pool_uuid", "")
	} else {
		d.Set("vlan_pool_uuid", poolValue)
	}

	return nil
}

func resourceMSOL3DomainImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO L3 Domain Resource - Beginning Import: %v", d.Id())
	resourceMSOL3DomainRead(d, m)
	log.Printf("[DEBUG] MSO L3 Domain Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOL3DomainCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO L3 Domain Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"name":       d.Get("name").(string),
		"templateId": models.StripQuotes(templateCont.S("templateId").String()),
		"schemaId":   models.StripQuotes(templateCont.S("schemaId").String()),
	}

	if description, ok := d.GetOk("description"); ok {
		payload["description"] = description.(string)
	}

	if vlanPoolUUID, ok := d.GetOk("vlan_pool_uuid"); ok && vlanPoolUUID.(string) != "" {
		payload["pool"] = vlanPoolUUID.(string)
	}

	payloadModel := models.GetPatchPayload("add", "/fabricPolicyTemplate/template/l3Domains/-", payload)

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/L3Domain/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO L3 Domain Resource - Create Complete: %v", d.Id())
	return resourceMSOL3DomainRead(d, m)
}

func resourceMSOL3DomainRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO L3 Domain Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return err
	}

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	domainName, err := GetPolicyNameFromResourceId(d.Id(), "L3Domain")
	if err != nil {
		return err
	}

	domain, err := GetPolicyByName(response, domainName, "fabricPolicyTemplate", "template", "l3Domains")
	if err != nil {
		return err
	}

	setL3DomainData(d, domain, templateId)
	log.Printf("[DEBUG] MSO L3 Domain Resource - Read Complete: %v", d.Id())
	return nil
}

func resourceMSOL3DomainUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO L3 Domain Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	domainIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "fabricPolicyTemplate", "template", "l3Domains")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/fabricPolicyTemplate/template/l3Domains/%d", domainIndex)

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
		uuid := d.Get("vlan_pool_uuid").(string)
		if uuid != "" {
			err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/pool", updatePath), uuid)
			if err != nil {
				return err
			}
		} else {
			err := addPatchPayloadToContainer(payloadCont, "remove", fmt.Sprintf("%s/pool", updatePath), nil)
			if err != nil {
				return err
			}
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/L3Domain/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO L3 Domain Resource - Update Complete: %v", d.Id())
	return resourceMSOL3DomainRead(d, m)
}

func resourceMSOL3DomainDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO L3 Domain Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	domainIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "fabricPolicyTemplate", "template", "l3Domains")
	if err != nil {
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/fabricPolicyTemplate/template/l3Domains/%d", domainIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO L3 Domain Resource - Delete Complete: %v", d.Id())
	return nil
}
