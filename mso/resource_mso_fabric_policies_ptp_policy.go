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

func resourceMSOPtpPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOPtpPolicyCreate,
		Read:   resourceMSOPtpPolicyRead,
		Update: resourceMSOPtpPolicyUpdate,
		Delete: resourceMSOPtpPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOPtpPolicyImport,
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
			"admin_state": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"enabled", "disabled",
				}, false),
			},
			"fabric_profile_template": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"aes67", "default", "smpte",
				}, false),
			},
			"global_priority1": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: validation.IntBetween(0, 255),
			},
			"global_priority2": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: validation.IntBetween(0, 255),
			},
			"global_domain": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: validation.IntBetween(0, 128),
			},
			"fabric_sync_interval": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: validation.IntBetween(-4, 1),
			},
			"fabric_delay_interval": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: validation.IntBetween(-4, 5),
			},
			"fabric_announce_interval": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: validation.IntBetween(-3, 4),
			},
			"fabric_announce_timeout": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: validation.IntBetween(2, 10),
			},
			"ptp_profile": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"profile_template": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"aes67", "default", "smpte", "telecom",
							}, false),
						},
						"delay_interval": {
							Type:     schema.TypeInt,
							Required: true,
							ValidateFunc: validation.IntBetween(-4, 5),
						},
						"sync_interval": {
							Type:     schema.TypeInt,
							Required: true,
							ValidateFunc: validation.IntBetween(-4, 1),
						},
						"announce_interval": {
							Type:     schema.TypeInt,
							Required: true,
							ValidateFunc: validation.IntBetween(-3, 4),
						},
						"announce_timeout": {
							Type:     schema.TypeInt,
							Required: true,
							ValidateFunc: validation.IntBetween(2, 10),
						},
						"override_node_profile": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func setPtpProfiles(profileEntries *schema.Set) []map[string]any {
	profileList := profileEntries.List()
	profiles := make([]map[string]any, 0, 1)

	for _, val := range profileList {
		ptpProfileEntry := val.(map[string]any)
		template := ptpProfileEntry["profile_template"].(string)
		if (template == "telecom") {
			template = "telecomFullPath"
		}
		profile := map[string]any{
			"name":                  ptpProfileEntry["name"].(string),
			"description":           ptpProfileEntry["description"].(string),
			"profileTemplate":       template,
			"announceIntvl":         ptpProfileEntry["announce_interval"].(int),
			"syncIntvl":             ptpProfileEntry["sync_interval"].(int),
			"delayIntvl":            ptpProfileEntry["delay_interval"].(int),
			"announceTimeout":       ptpProfileEntry["announce_timeout"].(int),
			"nodeProfileOverride":   ptpProfileEntry["override_node_profile"].(bool),
			// TODO Maybe more
		}
		profiles = append(profiles, profile)
	}

	return profiles
}

func setPtpPolicyData(d *schema.ResourceData, msoClient *client.Client, templateId, policyName string) error {

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "fabricPolicyTemplate", "template", "ptpPolicy")
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/ptpPolicy/%s", templateId, models.StripQuotes(policy.S("name").String())))
	d.Set("template_id", templateId)
	d.Set("name", models.StripQuotes(policy.S("name").String()))
	d.Set("description", models.StripQuotes(policy.S("description").String()))
	d.Set("uuid", models.StripQuotes(policy.S("global").S("uuid").String()))

	globalCont := policy.S("global")
	d.Set("admin_state", models.StripQuotes(globalCont.S("adminState").String()))
	d.Set("global_priority1", globalCont.S("prio1").Data().(float64))
	d.Set("global_priority2", globalCont.S("prio2").Data().(float64))
	d.Set("global_domain", globalCont.S("globalDomain").Data().(float64))
	d.Set("fabric_profile_template", models.StripQuotes(globalCont.S("fabProfileTemplate").String()))
	d.Set("fabric_announce_interval", globalCont.S("fabAnnounceIntvl").Data().(float64))
	d.Set("fabric_sync_interval", globalCont.S("fabSyncIntvl").Data().(float64))
	d.Set("fabric_delay_interval", globalCont.S("fabDelayIntvl").Data().(float64))
	d.Set("fabric_announce_timeout", globalCont.S("fabAnnounceTimeout").Data().(float64))

	count, err := policy.ArrayCount("profiles")
	if err != nil {
		return fmt.Errorf("unable to count the number of PTP profiles: %s", err)
	}
	profiles := make([]any, 0)
	for i := range count {
		profilesCont, err := policy.ArrayElement(i, "profiles")
		if err != nil {
			return fmt.Errorf("unable to parse element %d from the list of PTP profiles: %s", i, err)
		}
		ptpProfileEntry := make(map[string]any)
		ptpProfileEntry["name"] = models.StripQuotes(profilesCont.S("name").String())
		ptpProfileEntry["description"] = models.StripQuotes(profilesCont.S("description").String())
		ptpProfileEntry["delay_interval"] = profilesCont.S("delayIntvl").Data().(float64)
		ptpProfileEntry["sync_interval"] = profilesCont.S("syncIntvl").Data().(float64)
		ptpProfileEntry["announce_interval"] = profilesCont.S("announceIntvl").Data().(float64)
		ptpProfileEntry["announce_timeout"] = profilesCont.S("announceTimeout").Data().(float64)
		template := models.StripQuotes(profilesCont.S("profileTemplate").String())
		if (template == "telecomFullPath") {
			template = "telecom"
		}
		ptpProfileEntry["profile_template"] = template
		if profilesCont.Exists("nodeProfileOverride") {
			ptpProfileEntry["override_node_profile"] = profilesCont.S("nodeProfileOverride").Data().(bool)
		}
		profiles = append(profiles, ptpProfileEntry)
	}
	d.Set("ptp_profile", profiles)

	return nil
}

func resourceMSOPtpPolicyImport(d *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO PTP Policy Resource - Beginning Import: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return nil, err
	}

	policyName, err := GetPolicyNameFromResourceId(d.Id(), "ptpPolicy")
	if err != nil {
		return nil, err
	}

	setPtpPolicyData(d, msoClient, templateId, policyName)
	log.Printf("[DEBUG] MSO PTP Policy Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOPtpPolicyCreate(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO PTP Policy Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	payload := map[string]any{}

	payload["name"] = d.Get("name").(string)

	if description, ok := d.GetOk("description"); ok {
		payload["description"] = description.(string)
	}

	globalParams := make(map[string]any)

	if adminState, ok := d.GetOk("admin_state"); ok {
		globalParams["adminState"] = adminState.(string)
	}

	if profileTemplate, ok := d.GetOk("fabric_profile_template"); ok {
		globalParams["fabProfileTemplate"] = profileTemplate.(string)
	}

	if prio1, ok := d.GetOk("global_priority1"); ok {
		globalParams["prio1"] = prio1.(int)
	}

	if prio2, ok := d.GetOk("global_priority2"); ok {
		globalParams["prio2"] = prio2.(int)
	}

	if globalDomain, ok := d.GetOk("global_domain"); ok {
		globalParams["globalDomain"] = globalDomain.(int)
	}

	if fabAnnounceIntvl, ok := d.GetOk("fabric_announce_interval"); ok {
		globalParams["fabAnnounceIntvl"] = fabAnnounceIntvl.(int)
	}

	if fabSyncIntvl, ok := d.GetOk("fabric_sync_interval"); ok {
		globalParams["fabSyncIntvl"] = fabSyncIntvl.(int)
	}

	if fabDelayIntvl, ok := d.GetOk("fabric_delay_interval"); ok {
		globalParams["fabDelayIntvl"] = fabDelayIntvl.(int)
	}

	if fabAnnounceTimeout, ok := d.GetOk("fabric_announce_timeout"); ok {
		globalParams["fabAnnounceTimeout"] = fabAnnounceTimeout.(int)
	}

	if len(globalParams) > 0 {
		payload["global"] = globalParams
	}

	if profileEntries, ok := d.GetOk("ptp_profile"); ok {
		payload["profiles"] = setPtpProfiles(profileEntries.(*schema.Set))
	}

	payloadModel := models.GetPatchPayload("add", "/fabricPolicyTemplate/template/ptpPolicy", payload)
	templateId := d.Get("template_id").(string)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/ptpPolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO PTP Policy Resource - Create Complete: %v", d.Id())
	return resourceMSOPtpPolicyRead(d, m)
}

func resourceMSOPtpPolicyRead(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO PTP Policy Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	setPtpPolicyData(d, msoClient, templateId, policyName)
	log.Printf("[DEBUG] MSO PTP Policy Resource - Read Complete : %v", d.Id())
	return nil
}

func resourceMSOPtpPolicyUpdate(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO PTP Policy Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	updatePath := "/fabricPolicyTemplate/template/ptpPolicy"

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
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/global/adminState", updatePath), d.Get("admin_state").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("fabric_profile_template") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/global/fabProfileTemplate", updatePath), d.Get("fabric_profile_template").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("global_priority1") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/global/prio1", updatePath), d.Get("global_priority1").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("global_priority2") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/global/prio2", updatePath), d.Get("global_priority2").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("global_domain") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/global/globalDomain", updatePath), d.Get("global_domain").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("fabric_announce_interval") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/global/fabAnnounceIntvl", updatePath), d.Get("fabric_announce_interval").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("fabric_sync_interval") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/global/fabSyncIntvl", updatePath), d.Get("fabric_sync_interval").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("fabric_delay_interval") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/global/fabDelayIntvl", updatePath), d.Get("fabric_delay_interval").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("fabric_announce_timeout") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/global/fabAnnounceTimeout", updatePath), d.Get("fabric_announce_timeout").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("ptp_profile") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/profiles", updatePath), setPtpProfiles(d.Get("ptp_profile").(*schema.Set)))
		if err != nil {
			return err
		}
	}

	err := doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/ptpPolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO PTP Policy Resource - Update Complete: %v", d.Id())
	return resourceMSOPtpPolicyRead(d, m)
}

func resourceMSOPtpPolicyDelete(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO PTP Policy Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	payloadModel := models.GetRemovePatchPayload("/fabricPolicyTemplate/template/ptpPolicy")

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", d.Get("template_id").(string)), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO PTP Policy Resource - Delete Complete: %v", d.Id())
	return nil
}
