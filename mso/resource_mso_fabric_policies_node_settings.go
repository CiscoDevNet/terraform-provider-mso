package mso

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
			"synce": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
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
			"ptp": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_domain": &schema.Schema{
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(24, 43),
						},
						"priority_2": &schema.Schema{
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
	synceMap := synce.(map[string]any)
	syncePayload := map[string]string{
		"adminState": synceMap["admin_state"].(string),
		"qlOption":   convertValueWithMap(synceMap["quality_level"].(string), synceQualityLevelOptionsMap),
	}
	return syncePayload
}

func getPtpPayload(ptp any) (map[string]int, error) {
	ptpMap := ptp.(map[string]any)
	domain, err := strconv.Atoi(ptpMap["node_domain"].(string))
	if err != nil {
		return nil, err
	}
	prio2, err := strconv.Atoi(ptpMap["priority_2"].(string))
	if err != nil {
		return nil, err
	}
	ptpPayload := map[string]int{
		"domain": domain,
		"prio2":  prio2,
		"prio1":  128,
	}
	return ptpPayload, nil
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
		synceMap := map[string]any{
			"admin_state":   models.StripQuotes(synce.S("adminState").String()),
			"quality_level": convertValueWithMap(models.StripQuotes(synce.S("qlOption").String()), synceQualityLevelOptionsMap),
		}
		d.Set("synce", synceMap)
	}

	if policy.Exists("ptp") {
		ptp := policy.S("ptp")
		ptpMap := map[string]any{
			"node_domain": ptp.S("domain").Data().(float64),
			"priority_2":  ptp.S("prio1").Data().(float64),
		}
		d.Set("ptp", ptpMap)
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

	setNodeSettingsData(d, msoClient, templateId, policyName)
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

	if synce, ok := d.GetOk("synce"); ok {
		syncePayload := getSyncePayload(synce)
		payload["synce"] = syncePayload
	}

	if ptp, ok := d.GetOk("ptp"); ok {
		ptpPayload, err := getPtpPayload(ptp)
		if err != nil {
			return err
		}
		payload["ptp"] = ptpPayload
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

	setNodeSettingsData(d, msoClient, templateId, policyName)
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
		if synce, ok := d.GetOk("synce"); ok {
			syncePayload := getSyncePayload(synce)
			err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/synce", updatePath), syncePayload)
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
		if ptp, ok := d.GetOk("ptp"); ok {
			ptpPayload, err := getPtpPayload(ptp)
			if err != nil {
				return err
			}
			err = addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/ptp", updatePath), ptpPayload)
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
