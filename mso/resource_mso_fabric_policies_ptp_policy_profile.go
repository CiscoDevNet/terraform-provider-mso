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

func resourceMSOPtpPolicyProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOPtpPolicyProfileCreate,
		Read:   resourceMSOPtpPolicyProfileRead,
		Update: resourceMSOPtpPolicyProfileUpdate,
		Delete: resourceMSOPtpPolicyProfileDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOPtpPolicyProfileImport,
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: validation.StringLenBetween(1, 16),
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
			"profile_template": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"aes67", "default", "smpte", "telecom",
				}, false),
			},
			"delay_interval": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(-4, 5),
			},
			"sync_interval": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(-4, 1),
			},
			"announce_interval": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(-3, 4),
			},
			"announce_timeout": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(2, 10),
			},
			"override_node_profile": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func setPtpPolicyProfileData(d *schema.ResourceData, msoClient *client.Client, templateId, policyName string) error {
	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "fabricPolicyTemplate", "template", "ptpPolicy", "profiles")
	if err != nil {
		return err
	}

	name := models.StripQuotes(policy.S("name").String())
	d.SetId(fmt.Sprintf("templateId/%s/ptpPolicyProfile/%s", templateId, name))
	d.Set("template_id", templateId)
	d.Set("name", name)
	d.Set("description", models.StripQuotes(policy.S("description").String()))
	d.Set("uuid", models.StripQuotes(policy.S("uuid").String()))
	d.Set("delay_interval", policy.S("delayIntvl").Data().(float64))
	d.Set("sync_interval", policy.S("syncIntvl").Data().(float64))
	d.Set("announce_interval", policy.S("announceIntvl").Data().(float64))
	d.Set("announce_timeout", policy.S("announceTimeout").Data().(float64))
	template := models.StripQuotes(policy.S("profileTemplate").String())
	if template == "telecomFullPath" {
		template = "telecom"
	}
	d.Set("profile_template", template)
	if policy.Exists("nodeProfileOverride") {
		d.Set("override_node_profile", policy.S("nodeProfileOverride").Data().(bool))
	}

	return nil
}

func resourceMSOPtpPolicyProfileImport(d *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO PTP Policy Profile Resource - Beginning Import: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return nil, err
	}

	policyName, err := GetPolicyNameFromResourceId(d.Id(), "ptpPolicyProfile")
	if err != nil {
		return nil, err
	}

	setPtpPolicyProfileData(d, msoClient, templateId, policyName)
	log.Printf("[DEBUG] MSO PTP Policy Profile Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOPtpPolicyProfileCreate(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO PTP Policy Profile Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	payload := map[string]any{}

	payload["name"] = d.Get("name").(string)

	if description, ok := d.GetOk("description"); ok {
		payload["description"] = description.(string)
	}

	if profile_template, ok := d.GetOk("profile_template"); ok {
		if profile_template == "telecom" {
			profile_template = "telecomFullPath"
		}
		payload["profileTemplate"] = profile_template.(string)
	}

	if announce_interval, ok := d.GetOk("announce_interval"); ok {
		payload["announceIntvl"] = announce_interval.(int)
	}

	if sync_interval, ok := d.GetOk("sync_interval"); ok {
		payload["syncIntvl"] = sync_interval.(int)
	}

	if delay_interval, ok := d.GetOk("delay_interval"); ok {
		payload["delayIntvl"] = delay_interval.(int)
	}

	if announce_timeout, ok := d.GetOk("announce_timeout"); ok {
		payload["announceTimeout"] = announce_timeout.(int)
	}

	if override_node_profile, ok := d.GetOk("override_node_profile"); ok {
		if override_node_profile.(bool) {
			payload["announceIntvl"] = override_node_profile.(bool)
		}
	}

	payloadModel := models.GetPatchPayload("add", "/fabricPolicyTemplate/template/ptpPolicy/profiles/-", payload)
	templateId := d.Get("template_id").(string)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/ptpPolicyProfile/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO PTP Policy Profile Resource - Create Complete: %v", d.Id())
	return resourceMSOPtpPolicyProfileRead(d, m)
}

func resourceMSOPtpPolicyProfileRead(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO PTP Policy Profile Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	setPtpPolicyProfileData(d, msoClient, templateId, policyName)
	log.Printf("[DEBUG] MSO PTP Policy Profile Resource - Read Complete : %v", d.Id())
	return nil
}

func resourceMSOPtpPolicyProfileUpdate(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO PTP Policy Profile Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "fabricPolicyTemplate", "template", "ptpPolicy", "profiles")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/fabricPolicyTemplate/template/ptpPolicy/profiles/%d", policyIndex)

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

	if d.HasChange("profile_template") {
		profile_template := d.Get("profile_template").(string)
		if (profile_template == "telecom") {
			profile_template = "telecomFullPath"
		}
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/profileTemplate", updatePath), profile_template)
		if err != nil {
			return err
		}
	}

	if d.HasChange("announce_interval") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/announceIntvl", updatePath), d.Get("announce_interval").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("sync_interval") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/syncIntvl", updatePath), d.Get("sync_interval").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("delay_interval") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/delayIntvl", updatePath), d.Get("delay_interval").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("announce_timeout") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/announceTimeout", updatePath), d.Get("announce_timeout").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("override_node_profile") {
		override := d.Get("override_node_profile").(bool)
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/nodeProfileOverride", updatePath), override)
		if err != nil {
			return err
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/ptpPolicyProfile/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO PTP Policy Profile Resource - Update Complete: %v", d.Id())
	return resourceMSOPtpPolicyProfileRead(d, m)
}

func resourceMSOPtpPolicyProfileDelete(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO PTP Policy Profile Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", d.Get("template_id").(string)))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "fabricPolicyTemplate", "template", "ptpPolicy", "profiles")
	if err != nil {
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/fabricPolicyTemplate/template/ptpPolicy/profiles/%d", policyIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", d.Get("template_id").(string)), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO PTP Policy Profile Resource - Delete Complete: %v", d.Id())
	return nil
}
