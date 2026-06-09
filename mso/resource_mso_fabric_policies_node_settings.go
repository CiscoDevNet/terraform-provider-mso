package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceMSONodeSettings() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSONodeSettingsCreate,
		Read:   resourceMSONodeSettingsRead,
		Update: resourceMSONodeSettingsUpdate,
		Delete: resourceMSONodeSettingsDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSONodeSettingsImport,
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
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"synce": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"admin_state": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"enabled", "disabled",
							}, false),
						},
						"quality_level": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"option_1", "option_2_generation_1", "option_2_generation_2",
							}, false),
						},
					},
				},
			},
			"ptp": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_domain": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(24, 43),
						},
						"priority_2": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 255),
						},
					},
				},
			},
		},
	}
}

func getSyncePayload(synce any) map[string]string {
	synceMap := synce.([]interface{})[0].(map[string]interface{})
	return map[string]string{
		"adminState": synceMap["admin_state"].(string),
		"qlOption":   convertValueWithMap(synceMap["quality_level"].(string), synceQualityLevelOptionsMap),
	}
}

func getPtpPayload(ptp any) map[string]int {
	ptpMap := ptp.([]interface{})[0].(map[string]interface{})
	return map[string]int{
		"domain": ptpMap["node_domain"].(int),
		"prio2":  ptpMap["priority_2"].(int),
		"prio1":  128,
	}
}

func setNodeSettingsData(d *schema.ResourceData, msoClient *client.Client, templateId, policyName string) error {
	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "fabricPolicyTemplate", "template", "nodePolicyGroups")
	if err != nil {
		return err
	}

	name := models.StripQuotes(policy.S("name").String())
	d.SetId(fmt.Sprintf("templateId/%s/nodeSettings/%s", templateId, name))
	d.Set("template_id", templateId)
	d.Set("name", name)
	d.Set("description", models.StripQuotes(policy.S("description").String()))
	d.Set("uuid", models.StripQuotes(policy.S("uuid").String()))

	if policy.Exists("synce") {
		synce := policy.S("synce")
		d.Set("synce", []interface{}{map[string]interface{}{
			"admin_state":   models.StripQuotes(synce.S("adminState").String()),
			"quality_level": convertValueWithMap(models.StripQuotes(synce.S("qlOption").String()), synceQualityLevelOptionsMap),
		}})
	}

	if policy.Exists("ptp") {
		ptp := policy.S("ptp")
		d.Set("ptp", []interface{}{map[string]interface{}{
			"node_domain": int(ptp.S("domain").Data().(float64)),
			"priority_2":  int(ptp.S("prio2").Data().(float64)),
		}})
	}

	return nil
}

func resourceMSONodeSettingsImport(d *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO Node Settings Resource - Beginning Import: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return nil, err
	}

	policyName, err := GetPolicyNameFromResourceId(d.Id(), "nodeSettings")
	if err != nil {
		return nil, err
	}

	err = setNodeSettingsData(d, msoClient, templateId, policyName)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] MSO Node Settings Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSONodeSettingsCreate(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO Node Settings Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	payload := map[string]any{}

	payload["name"] = d.Get("name").(string)

	if description, ok := d.GetOk("description"); ok {
		payload["description"] = description.(string)
	}

	if synceList := d.Get("synce").([]interface{}); len(synceList) > 0 {
		payload["synce"] = getSyncePayload(synceList)
	}

	if ptpList := d.Get("ptp").([]interface{}); len(ptpList) > 0 {
		payload["ptp"] = getPtpPayload(ptpList)
	}

	payloadModel := models.GetPatchPayload("add", "/fabricPolicyTemplate/template/nodePolicyGroups/-", payload)
	templateId := d.Get("template_id").(string)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/nodeSettings/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO Node Settings Resource - Create Complete: %v", d.Id())
	return resourceMSONodeSettingsRead(d, m)
}

func resourceMSONodeSettingsRead(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO Node Settings Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	err := setNodeSettingsData(d, msoClient, templateId, policyName)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] MSO Node Settings Resource - Read Complete : %v", d.Id())
	return nil
}

func resourceMSONodeSettingsUpdate(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO Node Settings Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "fabricPolicyTemplate", "template", "nodePolicyGroups")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/fabricPolicyTemplate/template/nodePolicyGroups/%d", policyIndex)

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

	if d.HasChange("synce") {
		if synceList := d.Get("synce").([]interface{}); len(synceList) > 0 {
			err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/synce", updatePath), getSyncePayload(synceList))
			if err != nil {
				return err
			}
		} else {
			err := addPatchPayloadToContainer(payloadCont, "remove", fmt.Sprintf("%s/synce", updatePath), nil)
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("ptp") {
		if ptpList := d.Get("ptp").([]interface{}); len(ptpList) > 0 {
			err = addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/ptp", updatePath), getPtpPayload(ptpList))
			if err != nil {
				return err
			}
		} else {
			err := addPatchPayloadToContainer(payloadCont, "remove", fmt.Sprintf("%s/ptp", updatePath), nil)
			if err != nil {
				return err
			}
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/nodeSettings/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO Node Settings Resource - Update Complete: %v", d.Id())
	return resourceMSONodeSettingsRead(d, m)
}

func resourceMSONodeSettingsDelete(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO Node Settings Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", d.Get("template_id").(string)))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "fabricPolicyTemplate", "template", "nodePolicyGroups")
	if err != nil {
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/fabricPolicyTemplate/template/nodePolicyGroups/%d", policyIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", d.Get("template_id").(string)), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO Node Settings Resource - Delete Complete: %v", d.Id())
	return nil
}
