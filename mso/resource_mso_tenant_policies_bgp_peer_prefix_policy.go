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

func resourceMSOBGPPeerPrefixPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOBGPPeerPrefixPolicyCreate,
		Read:   resourceMSOBGPPeerPrefixPolicyRead,
		Update: resourceMSOBGPPeerPrefixPolicyUpdate,
		Delete: resourceMSOBGPPeerPrefixPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOBGPPeerPrefixPolicyImport,
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
				Description:  "The name of the BGP Peer Prefix Policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the BGP Peer Prefix Policy.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the BGP Peer Prefix Policy.",
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"log", "reject", "restart", "shutdown",
				}, false),
				Description: "The action of the BGP Peer Prefix Policy. Valid values are 'log', 'reject', 'restart', 'shutdown'.",
			},
			"max_number_of_prefixes": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 300000),
				Description:  "The maximum number of prefixes for the BGP Peer Prefix Policy. Value must be between 1 and 300000.",
			},
			"threshold_percentage": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 100),
				Description:  "The threshold percentage of the BGP Peer Prefix Policy. Value must be between 1 and 100.",
			},
			"restart_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
				Description:  "The restart time of the BGP Peer Prefix Policy in seconds. Value must be between 1 and 65535.",
			},
		},
	}
}

func setBGPPeerPrefixPolicyData(d *schema.ResourceData, response *container.Container, templateId string) error {
	d.SetId(fmt.Sprintf("templateId/%s/BGPPeerPrefixPolicy/%s", templateId, models.StripQuotes(response.S("name").String())))
	d.Set("template_id", templateId)
	d.Set("name", models.StripQuotes(response.S("name").String()))
	d.Set("description", models.StripQuotes(response.S("description").String()))
	d.Set("action", models.StripQuotes(response.S("action").String()))
	d.Set("max_number_of_prefixes", response.S("maxPrefixes").Data().(float64))
	d.Set("threshold_percentage", response.S("threshold").Data().(float64))
	d.Set("restart_time", response.S("restartTime").Data().(float64))
	d.Set("uuid", models.StripQuotes(response.S("uuid").String()))

	return nil
}

func resourceMSOBGPPeerPrefixPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO BGP Peer Prefix Policy Resource - Beginning Import: %v", d.Id())
	resourceMSOBGPPeerPrefixPolicyRead(d, m)
	log.Printf("[DEBUG] MSO BGP Peer Prefix Policy Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOBGPPeerPrefixPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO BGP Peer Prefix Policy Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	payload := map[string]interface{}{}

	payload["name"] = d.Get("name").(string)

	if description, ok := d.GetOk("description"); ok {
		payload["description"] = description.(string)
	}

	if action, ok := d.GetOk("action"); ok {
		payload["action"] = action.(string)
	}

	if maxNumberOfPrefixes, ok := d.GetOk("max_number_of_prefixes"); ok {
		payload["maxPrefixes"] = maxNumberOfPrefixes.(int)
	}

	if thresholdPercentage, ok := d.GetOk("threshold_percentage"); ok {
		payload["threshold"] = thresholdPercentage.(int)
	}

	if restartTime, ok := d.GetOk("restart_time"); ok {
		payload["restartTime"] = restartTime.(int)
	}

	payloadModel := models.GetPatchPayload("add", "/tenantPolicyTemplate/template/bgpPeerPrefixPolicies/-", payload)
	templateId := d.Get("template_id").(string)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/BGPPeerPrefixPolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO BGP Peer Prefix Policy Resource - Create Complete: %v", d.Id())
	return resourceMSOBGPPeerPrefixPolicyRead(d, m)
}

func resourceMSOBGPPeerPrefixPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO BGP Peer Prefix Policy Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return err
	}

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyName, err := GetPolicyNameFromResourceId(d.Id(), "BGPPeerPrefixPolicy")
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "bgpPeerPrefixPolicies")
	if err != nil {
		return err
	}

	setBGPPeerPrefixPolicyData(d, policy, templateId)
	log.Printf("[DEBUG] MSO BGP Peer Prefix Policy Resource - Read Complete: %v", d.Id())
	return nil
}

func resourceMSOBGPPeerPrefixPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO BGP Peer Prefix Policy Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "bgpPeerPrefixPolicies")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/tenantPolicyTemplate/template/bgpPeerPrefixPolicies/%d", policyIndex)

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

	if d.HasChange("action") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/action", updatePath), d.Get("action").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("max_number_of_prefixes") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/maxPrefixes", updatePath), d.Get("max_number_of_prefixes").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("threshold_percentage") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/threshold", updatePath), d.Get("threshold_percentage").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("restart_time") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/restartTime", updatePath), d.Get("restart_time").(int))
		if err != nil {
			return err
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/BGPPeerPrefixPolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO BGP Peer Prefix Policy Resource - Update Complete: %v", d.Id())
	return resourceMSOBGPPeerPrefixPolicyRead(d, m)
}

func resourceMSOBGPPeerPrefixPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO BGP Peer Prefix Policy Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "bgpPeerPrefixPolicies")
	if err != nil {
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/tenantPolicyTemplate/template/bgpPeerPrefixPolicies/%d", policyIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO BGP Peer Prefix Policy Resource - Delete Complete: %v", d.Id())
	return nil
}
