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

func resourceMSOCustomQoSPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOCustomQoSPolicyCreate,
		Read:   resourceMSOCustomQoSPolicyRead,
		Update: resourceMSOCustomQoSPolicyUpdate,
		Delete: resourceMSOCustomQoSPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOCustomQoSPolicyImport,
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
				Description:  "The name of the Custom QoS Policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the Custom QoS Policy.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the Custom QoS Policy.",
			},
			"dscp_mappings": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "The DSCP mappings of the Custom QoS Policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dscp_from": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "unspecified",
							ValidateFunc: validation.StringInSlice([]string{
								"af11", "af12", "af13", "af21", "af22", "af23", "af31", "af32", "af33",
								"af41", "af42", "af43", "cs0", "cs1", "cs2", "cs3", "cs4", "cs5", "cs6",
								"cs7", "expedited_forwarding", "unspecified", "voice_admit",
							}, false),
							Description: "The starting encoding point of the DSCP range.",
						},
						"dscp_to": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "unspecified",
							ValidateFunc: validation.StringInSlice([]string{
								"af11", "af12", "af13", "af21", "af22", "af23", "af31", "af32", "af33",
								"af41", "af42", "af43", "cs0", "cs1", "cs2", "cs3", "cs4", "cs5", "cs6",
								"cs7", "expedited_forwarding", "unspecified", "voice_admit",
							}, false),
							Description: "The ending encoding point of the DSCP range.",
						},
						"dscp_target": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "unspecified",
							ValidateFunc: validation.StringInSlice([]string{
								"af11", "af12", "af13", "af21", "af22", "af23", "af31", "af32", "af33",
								"af41", "af42", "af43", "cs0", "cs1", "cs2", "cs3", "cs4", "cs5", "cs6",
								"cs7", "expedited_forwarding", "unspecified", "voice_admit",
							}, false),
							Description: "The DSCP target encoding point for egressing traffic.",
						},
						"target_cos": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "unspecified",
							ValidateFunc: validation.StringInSlice([]string{
								"background", "best_effort", "excellent_effort",
								"critical_applications", "video", "voice",
								"internetwork_control", "network_control", "unspecified",
							}, false),
							Description: "The target CoS value/traffic type for egressing traffic.",
						},
						"qos_priority": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "unspecified",
							ValidateFunc: validation.StringInSlice([]string{
								"level1", "level2", "level3", "level4", "level5", "level6", "unspecified",
							}, false),
							Description: "The QoS priority level.",
						},
					},
				},
			},
			"cos_mappings": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "The CoS mappings of the Custom QoS Policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dot1p_from": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "unspecified",
							ValidateFunc: validation.StringInSlice([]string{
								"background", "best_effort", "excellent_effort",
								"critical_applications", "video", "voice",
								"internetwork_control", "network_control", "unspecified",
							}, false),
							Description: "The starting value/traffic type of the CoS range.",
						},
						"dot1p_to": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "unspecified",
							ValidateFunc: validation.StringInSlice([]string{
								"background", "best_effort", "excellent_effort",
								"critical_applications", "video", "voice",
								"internetwork_control", "network_control", "unspecified",
							}, false),
							Description: "The ending value/traffic type of the CoS range.",
						},
						"dscp_target": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "unspecified",
							ValidateFunc: validation.StringInSlice([]string{
								"af11", "af12", "af13", "af21", "af22", "af23", "af31", "af32", "af33",
								"af41", "af42", "af43", "cs0", "cs1", "cs2", "cs3", "cs4", "cs5", "cs6",
								"cs7", "expedited_forwarding", "unspecified", "voice_admit",
							}, false),
							Description: "The DSCP target encoding point for egressing traffic.",
						},
						"target_cos": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "unspecified",
							ValidateFunc: validation.StringInSlice([]string{
								"background", "best_effort", "excellent_effort",
								"critical_applications", "video", "voice",
								"internetwork_control", "network_control", "unspecified",
							}, false),
							Description: "The target CoS value/traffic type for egressing traffic.",
						},
						"qos_priority": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "unspecified",
							ValidateFunc: validation.StringInSlice([]string{
								"level1", "level2", "level3", "level4", "level5", "level6", "unspecified",
							}, false),
							Description: "The QoS priority level.",
						},
					},
				},
			},
		},
	}
}

func setCustomQoSPolicyData(d *schema.ResourceData, response *container.Container, templateId string) error {
	d.SetId(fmt.Sprintf("templateId/%s/CustomQoSPolicy/%s", templateId, models.StripQuotes(response.S("name").String())))
	d.Set("template_id", templateId)
	d.Set("name", models.StripQuotes(response.S("name").String()))
	d.Set("description", models.StripQuotes(response.S("description").String()))
	d.Set("uuid", models.StripQuotes(response.S("uuid").String()))

	dscpMappingsCount, _ := response.S("dscpMappings").ArrayCount()
	dscpMappings := make([]interface{}, 0, dscpMappingsCount)
	for i := 0; i < dscpMappingsCount; i++ {
		dscpMapping := response.S("dscpMappings").Index(i)
		mapping := map[string]interface{}{
			"dscp_from":    convertValueWithMap(models.StripQuotes(dscpMapping.S("dscpFrom").String()), targetDscpMap),
			"dscp_to":      convertValueWithMap(models.StripQuotes(dscpMapping.S("dscpTo").String()), targetDscpMap),
			"dscp_target":  convertValueWithMap(models.StripQuotes(dscpMapping.S("dscpTarget").String()), targetDscpMap),
			"target_cos":   convertValueWithMap(models.StripQuotes(dscpMapping.S("targetCos").String()), targetCosMap),
			"qos_priority": models.StripQuotes(dscpMapping.S("priority").String()),
		}
		dscpMappings = append(dscpMappings, mapping)
	}
	d.Set("dscp_mappings", dscpMappings)

	cosMappingsCount, _ := response.S("cosMappings").ArrayCount()
	cosMappings := make([]interface{}, 0, cosMappingsCount)
	for i := 0; i < cosMappingsCount; i++ {
		cosMapping := response.S("cosMappings").Index(i)
		mapping := map[string]interface{}{
			"dot1p_from":   convertValueWithMap(models.StripQuotes(cosMapping.S("dot1pFrom").String()), targetCosMap),
			"dot1p_to":     convertValueWithMap(models.StripQuotes(cosMapping.S("dot1pTo").String()), targetCosMap),
			"dscp_target":  convertValueWithMap(models.StripQuotes(cosMapping.S("dscpTarget").String()), targetDscpMap),
			"target_cos":   convertValueWithMap(models.StripQuotes(cosMapping.S("targetCos").String()), targetCosMap),
			"qos_priority": models.StripQuotes(cosMapping.S("priority").String()),
		}
		cosMappings = append(cosMappings, mapping)
	}
	d.Set("cos_mappings", cosMappings)

	return nil
}

func resourceMSOCustomQoSPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO Custom QoS Policy Resource - Beginning Import: %v", d.Id())
	resourceMSOCustomQoSPolicyRead(d, m)
	log.Printf("[DEBUG] MSO Custom QoS Policy Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOCustomQoSPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Custom QoS Policy Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	payload := map[string]interface{}{}

	payload["name"] = d.Get("name").(string)

	if description, ok := d.GetOk("description"); ok {
		payload["description"] = description.(string)
	} else {
		payload["description"] = ""
	}

	if dscpMappingsRaw, ok := d.GetOk("dscp_mappings"); ok {
		dscpMappingsSet := dscpMappingsRaw.(*schema.Set)
		dscpMappingsList := dscpMappingsSet.List()
		dscpMappings := make([]interface{}, 0, len(dscpMappingsList))

		for _, item := range dscpMappingsList {
			mapping := item.(map[string]interface{})
			dscpMapping := make(map[string]interface{})

			if val, ok := mapping["dscp_from"].(string); ok {
				dscpMapping["dscpFrom"] = convertValueWithMap(val, targetDscpMap)
			}

			if val, ok := mapping["dscp_to"].(string); ok {
				dscpMapping["dscpTo"] = convertValueWithMap(val, targetDscpMap)
			}

			if val, ok := mapping["dscp_target"].(string); ok {
				dscpMapping["dscpTarget"] = convertValueWithMap(val, targetDscpMap)
			}

			if val, ok := mapping["target_cos"].(string); ok {
				dscpMapping["targetCos"] = convertValueWithMap(val, targetCosMap)
			}

			if val, ok := mapping["qos_priority"].(string); ok {
				dscpMapping["priority"] = val
			}

			dscpMappings = append(dscpMappings, dscpMapping)
		}
		payload["dscpMappings"] = dscpMappings
	}

	if cosMappingsRaw, ok := d.GetOk("cos_mappings"); ok {
		cosMappingsSet := cosMappingsRaw.(*schema.Set)
		cosMappingsList := cosMappingsSet.List()
		cosMappings := make([]interface{}, 0, len(cosMappingsList))

		for _, item := range cosMappingsList {
			mapping := item.(map[string]interface{})
			cosMapping := make(map[string]interface{})

			if val, ok := mapping["dot1p_from"].(string); ok {
				cosMapping["dot1pFrom"] = convertValueWithMap(val, targetCosMap)
			}

			if val, ok := mapping["dot1p_to"].(string); ok {
				cosMapping["dot1pTo"] = convertValueWithMap(val, targetCosMap)
			}

			if val, ok := mapping["dscp_target"].(string); ok {
				cosMapping["dscpTarget"] = convertValueWithMap(val, targetDscpMap)
			}

			if val, ok := mapping["target_cos"].(string); ok {
				cosMapping["targetCos"] = convertValueWithMap(val, targetCosMap)
			}

			if val, ok := mapping["qos_priority"].(string); ok {
				cosMapping["priority"] = val
			}

			cosMappings = append(cosMappings, cosMapping)
		}
		payload["cosMappings"] = cosMappings
	}

	payloadModel := models.GetPatchPayload("add", "/tenantPolicyTemplate/template/qosPolicies/-", payload)
	templateId := d.Get("template_id").(string)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/CustomQoSPolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO Custom QoS Policy Resource - Create Complete: %v", d.Id())
	return resourceMSOCustomQoSPolicyRead(d, m)
}

func resourceMSOCustomQoSPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Custom QoS Policy Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return err
	}

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyName, err := GetPolicyNameFromResourceId(d.Id(), "CustomQoSPolicy")
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "qosPolicies")
	if err != nil {
		return err
	}

	setCustomQoSPolicyData(d, policy, templateId)
	log.Printf("[DEBUG] MSO Custom QoS Policy Resource - Read Complete: %v", d.Id())
	return nil
}

func resourceMSOCustomQoSPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Custom QoS Policy Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "qosPolicies")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/tenantPolicyTemplate/template/qosPolicies/%d", policyIndex)

	payloadCont := container.New()
	payloadCont.Array()

	if d.HasChange("name") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/name", updatePath), d.Get("name").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("description") {
		desc := d.Get("description").(string)
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/description", updatePath), desc)
		if err != nil {
			return err
		}
	}

	if d.HasChange("dscp_mappings") {
		dscpMappingsSet := d.Get("dscp_mappings").(*schema.Set)
		dscpMappingsList := dscpMappingsSet.List()
		dscpMappings := make([]interface{}, 0, len(dscpMappingsList))

		for _, item := range dscpMappingsList {
			mapping := item.(map[string]interface{})
			dscpMapping := make(map[string]interface{})

			if val, ok := mapping["dscp_from"].(string); ok {
				dscpMapping["dscpFrom"] = convertValueWithMap(val, targetDscpMap)
			}

			if val, ok := mapping["dscp_to"].(string); ok {
				dscpMapping["dscpTo"] = convertValueWithMap(val, targetDscpMap)
			}

			if val, ok := mapping["dscp_target"].(string); ok {
				dscpMapping["dscpTarget"] = convertValueWithMap(val, targetDscpMap)
			}

			if val, ok := mapping["target_cos"].(string); ok {
				dscpMapping["targetCos"] = convertValueWithMap(val, targetCosMap)
			}

			if val, ok := mapping["qos_priority"].(string); ok {
				dscpMapping["priority"] = val
			}

			dscpMappings = append(dscpMappings, dscpMapping)
		}
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/dscpMappings", updatePath), dscpMappings)
		if err != nil {
			return err
		}
	}

	if d.HasChange("cos_mappings") {
		cosMappingsSet := d.Get("cos_mappings").(*schema.Set)
		cosMappingsList := cosMappingsSet.List()
		cosMappings := make([]interface{}, 0, len(cosMappingsList))

		for _, item := range cosMappingsList {
			mapping := item.(map[string]interface{})
			cosMapping := make(map[string]interface{})

			if val, ok := mapping["dot1p_from"].(string); ok {
				cosMapping["dot1pFrom"] = convertValueWithMap(val, targetCosMap)
			}

			if val, ok := mapping["dot1p_to"].(string); ok {
				cosMapping["dot1pTo"] = convertValueWithMap(val, targetCosMap)
			}

			if val, ok := mapping["dscp_target"].(string); ok {
				cosMapping["dscpTarget"] = convertValueWithMap(val, targetDscpMap)
			}

			if val, ok := mapping["target_cos"].(string); ok {
				cosMapping["targetCos"] = convertValueWithMap(val, targetCosMap)
			}

			if val, ok := mapping["qos_priority"].(string); ok {
				cosMapping["priority"] = val
			}

			cosMappings = append(cosMappings, cosMapping)
		}
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/cosMappings", updatePath), cosMappings)
		if err != nil {
			return err
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/CustomQoSPolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO Custom QoS Policy Resource - Update Complete: %v", d.Id())
	return resourceMSOCustomQoSPolicyRead(d, m)
}

func resourceMSOCustomQoSPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Custom QoS Policy Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "qosPolicies")
	if err != nil {
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/tenantPolicyTemplate/template/qosPolicies/%d", policyIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO Custom QoS Policy Resource - Delete Complete: %v", d.Id())
	return nil
}
