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

func resourceMSOVlanPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOVlanPoolCreate,
		Read:   resourceMSOVlanPoolRead,
		Update: resourceMSOVlanPoolUpdate,
		Delete: resourceMSOVlanPoolDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOVlanPoolImport,
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:         schema.TypeString,
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
			"allocation_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"dynamic", "static",
				}, false),
			},
			"vlan_range": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"to": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"allocation_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"dynamic", "static",
							}, false),
						},
					},
				},
			},
		},
	}
}

func setVlanRange(rangeEntries *schema.Set) []map[string]any {
	rangeEntryList := rangeEntries.List()
	VlanRange := make([]map[string]any, len(rangeEntryList))

	for i, val := range rangeEntryList {
		rangeEntry := val.(map[string]any)
		VlanRange[i] = map[string]any{
			"range": map[string]any{
				"from":      rangeEntry["from"].(int),
				"to":        rangeEntry["to"].(int),
				"allocMode": rangeEntry["allocation_mode"].(string),
			},
		}
	}

	return VlanRange
}

func setVlanPoolData(d *schema.ResourceData, response *container.Container, templateId string) error {

	d.SetId(fmt.Sprintf("templateId/%s/VlanPool/%s", templateId, models.StripQuotes(response.S("name").String())))
	d.Set("template_id", templateId)
	d.Set("name", models.StripQuotes(response.S("name").String()))
	d.Set("description", models.StripQuotes(response.S("description").String()))
	d.Set("allocation_mode", models.StripQuotes(response.S("allocMode").String()))
	d.Set("uuid", models.StripQuotes(response.S("uuid").String()))

	count, _ := response.ArrayCount("encapBlocks")
	VlanRange := make([]any, 0)
	for i := range count {
		encapBlocksCont, err := response.ArrayElement(i, "encapBlocks")
		if err != nil {
			return fmt.Errorf("unable to parse encapsulation blocks for vlan range")
		}
		rangeEntry := make(map[string]any)
		rangeEntry["from"] = encapBlocksCont.S("range", "from").Data().(float64)
		rangeEntry["to"] = encapBlocksCont.S("range", "to").Data().(float64)
		rangeEntry["allocation_mode"] = models.StripQuotes(encapBlocksCont.S("range", "allocMode").String())
		VlanRange = append(VlanRange, rangeEntry)
	}
	d.Set("vlan_range", VlanRange)

	return nil

}

func resourceMSOVlanPoolImport(d *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO VLAN Pool Resource - Beginning Import: %v", d.Id())
	resourceMSOVlanPoolRead(d, m)
	log.Printf("[DEBUG] MSO VLAN Pool Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOVlanPoolCreate(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO VLAN Pool Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	payload := map[string]any{}

	payload["name"] = d.Get("name").(string)

	if description, ok := d.GetOk("description"); ok {
		payload["description"] = description.(string)
	}

	if rangeEntries, ok := d.GetOk("vlan_range"); ok {
		payload["encapBlocks"] = setVlanRange(rangeEntries.(*schema.Set))
	}

	payloadModel := models.GetPatchPayload("add", "/fabricPolicyTemplate/template/vlanPools/-", payload)
	templateId := d.Get("template_id").(string)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/VlanPool/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO VLAN Pool Resource - Create Complete: %v", d.Id())
	return resourceMSOVlanPoolRead(d, m)
}

func resourceMSOVlanPoolRead(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO VLAN Pool Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return err
	}

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyName, err := GetPolicyNameFromResourceId(d.Id(), "VlanPool")
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "fabricPolicyTemplate", "template", "vlanPools")
	if err != nil {
		return err
	}

	setVlanPoolData(d, policy, templateId)
	log.Printf("[DEBUG] MSO VLAN Pool Resource - Read Complete : %v", d.Id())
	return nil
}

func resourceMSOVlanPoolUpdate(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO VLAN Pool Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "fabricPolicyTemplate", "template", "vlanPools")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/fabricPolicyTemplate/template/vlanPools/%d", policyIndex)

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

	if d.HasChange("vlan_range") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/encapBlocks", updatePath), setVlanRange(d.Get("vlan_range").(*schema.Set)))
		if err != nil {
			return err
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/VlanPool/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO VLAN Pool Resource - Update Complete: %v", d.Id())
	return resourceMSOVlanPoolRead(d, m)
}

func resourceMSOVlanPoolDelete(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO VLAN Pool Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", d.Get("template_id").(string)))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "fabricPolicyTemplate", "template", "vlanPools")
	if err != nil {
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/fabricPolicyTemplate/template/vlanPools/%d", policyIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", d.Get("template_id").(string)), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO VLAN Pool Resource - Delete Complete: %v", d.Id())
	return nil
}
