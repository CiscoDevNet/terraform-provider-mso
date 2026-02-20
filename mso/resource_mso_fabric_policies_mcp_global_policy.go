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

func resourceMSOMCPGlobalPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOMCPGlobalPolicyCreate,
		Read:   resourceMSOMCPGlobalPolicyRead,
		Update: resourceMSOMCPGlobalPolicyUpdate,
		Delete: resourceMSOMCPGlobalPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOMCPGlobalPolicyImport,
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"admin_state": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
			},
			"enable_mcp_pdu_per_vlan": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
			},
			"key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"loop_detect_multiplication_factor": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 255),
			},
			"port_disable_protection": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
			},
			"initial_delay_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 1800),
			},
			"transmission_frequency_sec": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 300),
			},
			"transmission_frequency_msec": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 999),
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func setMCPGlobalPolicyData(d *schema.ResourceData, response *container.Container, templateId string) error {
	d.SetId(fmt.Sprintf("templateId/%s/MCPGlobalPolicy/%s", templateId, models.StripQuotes(response.S("name").String())))
	d.Set("template_id", templateId)
	d.Set("name", models.StripQuotes(response.S("name").String()))
	d.Set("description", models.StripQuotes(response.S("description").String()))
	d.Set("uuid", models.StripQuotes(response.S("uuid").String()))
	d.Set("admin_state", models.StripQuotes(response.S("adminState").String()))

	if response.Exists("key") {
		d.Set("key", models.StripQuotes(response.S("key").String()))
	} else {
		d.Set("key", nil)
	}

	if response.Exists("enablePduPerVlan") {
		d.Set("enable_mcp_pdu_per_vlan", enabledDisabledMap[response.Data().(map[string]interface{})["enablePduPerVlan"].(bool)])
	}

	if response.Exists("loopDetectMultFactor") {
		d.Set("loop_detect_multiplication_factor", int(response.Data().(map[string]interface{})["loopDetectMultFactor"].(float64)))
	}

	if response.Exists("protectPortDisable") {
		d.Set("port_disable_protection", enabledDisabledMap[response.Data().(map[string]interface{})["protectPortDisable"].(bool)])
	}

	if response.Exists("initialDelayTime") {
		d.Set("initial_delay_time", int(response.Data().(map[string]interface{})["initialDelayTime"].(float64)))
	}

	if response.Exists("txFreq") {
		d.Set("transmission_frequency_sec", int(response.Data().(map[string]interface{})["txFreq"].(float64)))
	}

	if response.Exists("txFreqMsec") {
		d.Set("transmission_frequency_msec", int(response.Data().(map[string]interface{})["txFreqMsec"].(float64)))
	}

	return nil
}

func resourceMSOMCPGlobalPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO MCP Global Policy Resource - Beginning Import: %v", d.Id())
	resourceMSOMCPGlobalPolicyRead(d, m)
	log.Printf("[DEBUG] MSO MCP Global Policy Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOMCPGlobalPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO MCP Global Policy Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"name":       d.Get("name").(string),
		"templateId": models.StripQuotes(templateCont.S("templateId").String()),
	}

	if description, ok := d.GetOk("description"); ok {
		payload["description"] = description.(string)
	}

	if adminState, ok := d.GetOk("admin_state"); ok {
		payload["adminState"] = adminState.(string)
	}

	if enablePduPerVlan, ok := d.GetOk("enable_mcp_pdu_per_vlan"); ok {
		payload["enablePduPerVlan"] = enabledDisabledMap[enablePduPerVlan.(string)]
	}

	if key, ok := d.GetOk("key"); ok {
		payload["key"] = key.(string)
	}

	if loopDetectMultFactor, ok := d.GetOk("loop_detect_multiplication_factor"); ok {
		payload["loopDetectMultFactor"] = loopDetectMultFactor.(int)
	}

	if protectPortDisable, ok := d.GetOk("port_disable_protection"); ok {
		payload["protectPortDisable"] = enabledDisabledMap[protectPortDisable.(string)]
	}

	if initialDelayTime, ok := d.GetOk("initial_delay_time"); ok {
		payload["initialDelayTime"] = initialDelayTime.(int)
	}

	if txFreq, ok := d.GetOk("transmission_frequency_sec"); ok {
		payload["txFreq"] = txFreq.(int)
	}

	if txFreqMsec, ok := d.GetOk("transmission_frequency_msec"); ok {
		payload["txFreqMsec"] = txFreqMsec.(int)
	}

	payloadModel := models.GetPatchPayload("add", "/fabricPolicyTemplate/template/mcpGlobalPolicy", payload)

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/MCPGlobalPolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO MCP Global Policy Resource - Create Complete: %v", d.Id())
	return resourceMSOMCPGlobalPolicyRead(d, m)
}

func resourceMSOMCPGlobalPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO MCP Global Policy Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return err
	}

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyName, err := GetPolicyNameFromResourceId(d.Id(), "MCPGlobalPolicy")
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "fabricPolicyTemplate", "template", "mcpGlobalPolicy")
	if err != nil {
		return err
	}

	setMCPGlobalPolicyData(d, policy, templateId)
	log.Printf("[DEBUG] MSO MCP Global Policy Resource - Read Complete: %v", d.Id())
	return nil
}

func resourceMSOMCPGlobalPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO MCP Global Policy Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	updatePath := "/fabricPolicyTemplate/template/mcpGlobalPolicy"

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

	if d.HasChange("admin_state") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/adminState", updatePath), d.Get("admin_state").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("enable_mcp_pdu_per_vlan") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/enablePduPerVlan", updatePath), enabledDisabledMap[d.Get("enable_mcp_pdu_per_vlan").(string)])
		if err != nil {
			return err
		}
	}

	if d.HasChange("key") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/key", updatePath), d.Get("key").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("loop_detect_multiplication_factor") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/loopDetectMultFactor", updatePath), d.Get("loop_detect_multiplication_factor").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("port_disable_protection") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/protectPortDisable", updatePath), enabledDisabledMap[d.Get("port_disable_protection").(string)])
		if err != nil {
			return err
		}
	}

	if d.HasChange("initial_delay_time") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/initialDelayTime", updatePath), d.Get("initial_delay_time").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("transmission_frequency_sec") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/txFreq", updatePath), d.Get("transmission_frequency_sec").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("transmission_frequency_msec") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/txFreqMsec", updatePath), d.Get("transmission_frequency_msec").(int))
		if err != nil {
			return err
		}
	}

	err := doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/MCPGlobalPolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO MCP Global Policy Resource - Update Complete: %v", d.Id())
	return resourceMSOMCPGlobalPolicyRead(d, m)
}

func resourceMSOMCPGlobalPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO MCP Global Policy Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)

	payloadModel := models.GetRemovePatchPayload("/fabricPolicyTemplate/template/mcpGlobalPolicy")

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO MCP Global Policy Resource - Delete Complete: %v", d.Id())
	return nil
}
