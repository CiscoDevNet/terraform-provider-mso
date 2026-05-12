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

func resourceMSOMacsecPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOMacsecPolicyCreate,
		Read:   resourceMSOMacsecPolicyRead,
		Update: resourceMSOMacsecPolicyUpdate,
		Delete: resourceMSOMacsecPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOMacsecPolicyImport,
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
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"enabled", "disabled",
				}, false),
			},
			"interface_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"fabric", "access",
				}, false),
			},
			"cipher_suite": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"128GcmAes", "128GcmAesXpn", "256GcmAes", "256GcmAesXpn",
				}, false),
			},
			"window_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"security_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"shouldSecure", "mustSecure",
				}, false),
			},
			"sak_expire_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"confidentiality_offset": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"offset0", "offset30", "offset50",
				}, false),
			},
			"key_server_priority": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"macsec_keys": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"psk": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func setMacsecKeys(macsecKeyEntries *schema.Set) []map[string]any {
	macsecKeyList := macsecKeyEntries.List()
	macsecKeys := make([]map[string]any, 0, 1)

	for _, val := range macsecKeyList {
		macsecKeyEntry := val.(map[string]any)
		macsecKey := map[string]any{
			"keyname": macsecKeyEntry["key_name"].(string),
			"psk":     macsecKeyEntry["psk"].(string),
			"start":   macsecKeyEntry["start_time"].(string),
			"end":     macsecKeyEntry["end_time"].(string),
		}
		macsecKeys = append(macsecKeys, macsecKey)
	}

	return macsecKeys
}

func setMacsecPolicyData(d *schema.ResourceData, msoClient *client.Client, templateId, policyName string) error {

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "fabricPolicyTemplate", "template", "macsecPolicies")
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/macsecPolicy/%s", templateId, models.StripQuotes(policy.S("name").String())))
	d.Set("template_id", templateId)
	d.Set("name", models.StripQuotes(policy.S("name").String()))
	d.Set("description", models.StripQuotes(policy.S("description").String()))
	d.Set("uuid", models.StripQuotes(policy.S("uuid").String()))
	d.Set("admin_state", models.StripQuotes(policy.S("adminState").String()))
	d.Set("interface_type", models.StripQuotes(policy.S("type").String()))
	d.Set("cipher_suite", models.StripQuotes(policy.S("macsecParams", "cipherSuite").String()))
	if v, ok := policy.S("macsecParams", "windowSize").Data().(float64); ok {
		d.Set("window_size", v)
	}
	d.Set("security_policy", models.StripQuotes(policy.S("macsecParams", "securityPol").String()))
	if v, ok := policy.S("macsecParams", "sakExpiryTime").Data().(float64); ok {
		d.Set("sak_expire_time", v)
	}
	if policy.S("macsecParams", "confOffSet").Data() != nil {
		d.Set("confidentiality_offset", models.StripQuotes(policy.S("macsecParams", "confOffSet").String()))
	}
	if v, ok := policy.S("macsecParams", "keyServerPrio").Data().(float64); ok {
		d.Set("key_server_priority", v)
	}

	existingPskByKeyName := map[string]string{}
	if existingKeys, ok := d.GetOk("macsec_keys"); ok {
		for _, v := range existingKeys.(*schema.Set).List() {
			key := v.(map[string]any)
			keyName := key["key_name"].(string)
			psk := key["psk"].(string)
			if psk != "" {
				existingPskByKeyName[keyName] = psk
			}
		}
	}

	count, err := policy.ArrayCount("macsecKeys")
	if err != nil {
		return fmt.Errorf("unable to count the number of macsec keys: %s", err)
	}
	macsecKeys := make([]any, 0)
	for i := range count {
		macsecKeysCont, err := policy.ArrayElement(i, "macsecKeys")
		if err != nil {
			return fmt.Errorf("unable to parse element %d from the list of macsec keys: %s", i, err)
		}
		macsecKeyEntry := make(map[string]any)
		keyName := models.StripQuotes(macsecKeysCont.S("keyname").String())
		macsecKeyEntry["key_name"] = keyName
		psk := models.StripQuotes(macsecKeysCont.S("psk").String())
		// The API returns an encrypted form of the PSK rather than the original plaintext.
		// Always prefer the plaintext value from state to avoid perpetual diffs.
		if existingPsk, exists := existingPskByKeyName[keyName]; exists {
			psk = existingPsk
		}
		macsecKeyEntry["psk"] = psk
		macsecKeyEntry["start_time"] = models.StripQuotes(macsecKeysCont.S("start").String())
		macsecKeyEntry["end_time"] = models.StripQuotes(macsecKeysCont.S("end").String())
		macsecKeys = append(macsecKeys, macsecKeyEntry)
	}
	d.Set("macsec_keys", macsecKeys)

	return nil

}

func resourceMSOMacsecPolicyImport(d *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO MACSec Policy Resource - Beginning Import: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return nil, err
	}

	policyName, err := GetPolicyNameFromResourceId(d.Id(), "macsecPolicy")
	if err != nil {
		return nil, err
	}

	setMacsecPolicyData(d, msoClient, templateId, policyName)
	log.Printf("[DEBUG] MSO MACSec Policy Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOMacsecPolicyCreate(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO MACSec Policy Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	payload := map[string]any{}

	payload["name"] = d.Get("name").(string)

	if description, ok := d.GetOk("description"); ok {
		payload["description"] = description.(string)
	}

	if adminState, ok := d.GetOk("admin_state"); ok {
		payload["adminState"] = adminState.(string)
	}

	interfaceType, ok := d.GetOk("interface_type")
	if !ok {
		return fmt.Errorf("interface_type is required when creating a MACsec Policy")
	}
	payload["type"] = interfaceType.(string)

	macsecParams := make(map[string]any)

	if cipherSuite, ok := d.GetOk("cipher_suite"); ok {
		macsecParams["cipherSuite"] = cipherSuite.(string)
	}

	if windowSize, ok := d.GetOk("window_size"); ok {
		macsecParams["windowSize"] = windowSize.(int)
	}

	if securityPol, ok := d.GetOk("security_policy"); ok {
		macsecParams["securityPol"] = securityPol.(string)
	}
	if sakExpiryTime, ok := d.GetOk("sak_expire_time"); ok {
		macsecParams["sakExpiryTime"] = sakExpiryTime.(int)
	}

	if confOffSet, ok := d.GetOk("confidentiality_offset"); ok {
		macsecParams["confOffSet"] = confOffSet.(string)
	}

	if keyServerPrio, ok := d.GetOk("key_server_priority"); ok {
		macsecParams["keyServerPrio"] = keyServerPrio.(int)
	}

	if len(macsecParams) > 0 {
		payload["macsecParams"] = macsecParams
	}

	if macsecKeyEntries, ok := d.GetOk("macsec_keys"); ok {
		payload["macsecKeys"] = setMacsecKeys(macsecKeyEntries.(*schema.Set))
	}

	payloadModel := models.GetPatchPayload("add", "/fabricPolicyTemplate/template/macsecPolicies/-", payload)
	templateId := d.Get("template_id").(string)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/macsecPolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO MACSec Policy Resource - Create Complete: %v", d.Id())
	return resourceMSOMacsecPolicyRead(d, m)
}

func resourceMSOMacsecPolicyRead(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO MACSec Policy Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	setMacsecPolicyData(d, msoClient, templateId, policyName)
	log.Printf("[DEBUG] MSO MACSec Policy Resource - Read Complete : %v", d.Id())
	return nil
}

func resourceMSOMacsecPolicyUpdate(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO MACSec Policy Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "fabricPolicyTemplate", "template", "macsecPolicies")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/fabricPolicyTemplate/template/macsecPolicies/%d", policyIndex)

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

	if d.HasChange("interface_type") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/type", updatePath), d.Get("interface_type").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("cipher_suite") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/macsecParams/cipherSuite", updatePath), d.Get("cipher_suite").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("window_size") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/macsecParams/windowSize", updatePath), d.Get("window_size").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("security_policy") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/macsecParams/securityPol", updatePath), d.Get("security_policy").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("sak_expire_time") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/macsecParams/sakExpiryTime", updatePath), d.Get("sak_expire_time").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("confidentiality_offset") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/macsecParams/confOffSet", updatePath), d.Get("confidentiality_offset").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("key_server_priority") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/macsecParams/keyServerPrio", updatePath), d.Get("key_server_priority").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("macsec_keys") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/macsecKeys", updatePath), setMacsecKeys(d.Get("macsec_keys").(*schema.Set)))
		if err != nil {
			return err
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/macsecPolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO MACSec Policy Resource - Update Complete: %v", d.Id())
	return resourceMSOMacsecPolicyRead(d, m)
}

func resourceMSOMacsecPolicyDelete(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] MSO MACSec Policy Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", d.Get("template_id").(string)))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "fabricPolicyTemplate", "template", "macsecPolicies")
	if err != nil {
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/fabricPolicyTemplate/template/macsecPolicies/%d", policyIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", d.Get("template_id").(string)), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO MACSec Policy Resource - Delete Complete: %v", d.Id())
	return nil
}
