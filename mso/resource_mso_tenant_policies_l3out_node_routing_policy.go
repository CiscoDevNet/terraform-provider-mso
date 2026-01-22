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

func resourceMSOL3OutNodeRoutingPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOL3OutNodeRoutingPolicyCreate,
		Read:   resourceMSOL3OutNodeRoutingPolicyRead,
		Update: resourceMSOL3OutNodeRoutingPolicyUpdate,
		Delete: resourceMSOL3OutNodeRoutingPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOL3OutNodeRoutingPolicyImport,
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
				Description:  "The name of the L3Out Node Routing Policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the L3Out Node Routing Policy.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the L3Out Node Routing Policy.",
			},
			"as_path_multipath_relax": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "BGP Best Path Control - AS path multipath relax. Allows load balancing across paths with different AS paths.",
			},
			"bfd_multi_hop_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "BFD multi-hop configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"admin_state": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"enabled", "disabled",
							}, false),
							Description: "Administrative state. Default: enabled when unset during creation.",
						},
						"detection_multiplier": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(1, 50),
							Description:  "Detection multiplier. Default: 3 when unset during creation. Range: 1-50.",
						},
						"min_receive_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(250, 999),
							Description:  "Minimum receive interval in microseconds. Default: 250 when unset during creation. Range: 250-999.",
						},
						"min_transmit_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(250, 999),
							Description:  "Minimum transmit interval in microseconds. Default: 250 when unset during creation. Range: 250-999.",
						},
					},
				},
			},
			"bgp_node_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "BGP node configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"graceful_restart_helper": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Enable graceful restart helper mode. Default: true (enabled) when unset during creation.",
						},
						"keep_alive_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(0, 3600),
							Description:  "BGP keepalive interval in seconds. Default: 60 when unset during creation. Range: 0-3600.",
						},
						"hold_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntAtLeast(0),
							Description:  "BGP hold interval in seconds. Default: 180 when unset during creation. Must be 0 or between 3-3600.",
						},
						"stale_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(1, 3600),
							Description:  "BGP stale interval in seconds for graceful restart. Default: 300 when unset during creation. Range: 1-3600.",
						},
						"max_as_limit": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(0, 2000),
							Description:  "Maximum AS path limit. Default: 0 (no limit) when unset during creation. Range: 0-2000.",
						},
					},
				},
			},
		},
	}
}

func buildNodeBFDMultiHopPayload(bfdMultiHop []interface{}) map[string]interface{} {
	if len(bfdMultiHop) == 0 {
		return nil
	}

	settings := bfdMultiHop[0].(map[string]interface{})
	payload := make(map[string]interface{})

	if adminState, ok := settings["admin_state"].(string); ok && adminState != "" {
		payload["adminState"] = adminState
	}
	if detectionMultiplier, ok := settings["detection_multiplier"].(int); ok && detectionMultiplier != 0 {
		payload["detectionMultiplier"] = detectionMultiplier
	}
	if minRxInterval, ok := settings["min_receive_interval"].(int); ok && minRxInterval != 0 {
		payload["minRxInterval"] = minRxInterval
	}
	if minTxInterval, ok := settings["min_transmit_interval"].(int); ok && minTxInterval != 0 {
		payload["minTxInterval"] = minTxInterval
	}

	return payload
}

func buildBGPNodeSettingsPayload(bgpSettings []interface{}) map[string]interface{} {
	if len(bgpSettings) == 0 {
		return nil
	}

	settings := bgpSettings[0].(map[string]interface{})
	payload := make(map[string]interface{})

	if gracefulRestart, ok := settings["graceful_restart_helper"].(bool); ok {
		payload["gracefulRestartHelper"] = gracefulRestart
	}
	if keepAlive, ok := settings["keep_alive_interval"].(int); ok && keepAlive != 0 {
		payload["keepAliveInterval"] = keepAlive
	}
	if holdInterval, ok := settings["hold_interval"].(int); ok {
		payload["holdInterval"] = holdInterval
	}
	if staleInterval, ok := settings["stale_interval"].(int); ok && staleInterval != 0 {
		payload["staleInterval"] = staleInterval
	}
	if maxAsLimit, ok := settings["max_as_limit"].(int); ok {
		payload["maxAslimit"] = maxAsLimit
	}

	return payload
}

func setL3OutNodeRoutingPolicyData(d *schema.ResourceData, response *container.Container, templateId string) error {
	d.SetId(fmt.Sprintf("templateId/%s/L3OutNodeRoutingPolicy/%s", templateId, models.StripQuotes(response.S("name").String())))
	d.Set("template_id", templateId)
	d.Set("name", models.StripQuotes(response.S("name").String()))
	d.Set("description", models.StripQuotes(response.S("description").String()))
	d.Set("uuid", models.StripQuotes(response.S("uuid").String()))

	if response.Exists("asPathPol") && response.S("asPathPol").Exists("asPathMultipathRelax") {
		if val, ok := response.S("asPathPol").S("asPathMultipathRelax").Data().(bool); ok {
			d.Set("as_path_multipath_relax", val)
		}
	} else {
		d.Set("as_path_multipath_relax", nil)
	}

	if response.Exists("bfdMultiHopPol") {
		bfdMultiHop := response.S("bfdMultiHopPol")
		bfdMultiHopMap := map[string]interface{}{}

		if bfdMultiHop.Exists("adminState") {
			bfdMultiHopMap["admin_state"] = models.StripQuotes(bfdMultiHop.S("adminState").String())
		}
		if bfdMultiHop.Exists("detectionMultiplier") {
			if val, ok := bfdMultiHop.S("detectionMultiplier").Data().(float64); ok {
				bfdMultiHopMap["detection_multiplier"] = int(val)
			}
		}
		if bfdMultiHop.Exists("minRxInterval") {
			if val, ok := bfdMultiHop.S("minRxInterval").Data().(float64); ok {
				bfdMultiHopMap["min_receive_interval"] = int(val)
			}
		}
		if bfdMultiHop.Exists("minTxInterval") {
			if val, ok := bfdMultiHop.S("minTxInterval").Data().(float64); ok {
				bfdMultiHopMap["min_transmit_interval"] = int(val)
			}
		}

		d.Set("bfd_multi_hop_settings", []interface{}{bfdMultiHopMap})
	} else {
		d.Set("bfd_multi_hop_settings", []interface{}{})
	}

	if response.Exists("bgpTimerPol") {
		bgpTimer := response.S("bgpTimerPol")
		bgpMap := map[string]interface{}{}

		if bgpTimer.Exists("gracefulRestartHelper") {
			if val, ok := bgpTimer.S("gracefulRestartHelper").Data().(bool); ok {
				bgpMap["graceful_restart_helper"] = val
			}
		}
		if bgpTimer.Exists("keepAliveInterval") {
			if val, ok := bgpTimer.S("keepAliveInterval").Data().(float64); ok {
				bgpMap["keep_alive_interval"] = int(val)
			}
		}
		if bgpTimer.Exists("holdInterval") {
			if val, ok := bgpTimer.S("holdInterval").Data().(float64); ok {
				bgpMap["hold_interval"] = int(val)
			}
		}
		if bgpTimer.Exists("staleInterval") {
			if val, ok := bgpTimer.S("staleInterval").Data().(float64); ok {
				bgpMap["stale_interval"] = int(val)
			}
		}
		if bgpTimer.Exists("maxAslimit") {
			if val, ok := bgpTimer.S("maxAslimit").Data().(float64); ok {
				bgpMap["max_as_limit"] = int(val)
			}
		}

		d.Set("bgp_node_settings", []interface{}{bgpMap})
	} else {
		d.Set("bgp_node_settings", []interface{}{})
	}

	return nil
}

func resourceMSOL3OutNodeRoutingPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO L3Out Node Routing Policy Resource - Beginning Import: %v", d.Id())
	resourceMSOL3OutNodeRoutingPolicyRead(d, m)
	log.Printf("[DEBUG] MSO L3Out Node Routing Policy Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOL3OutNodeRoutingPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO L3Out Node Routing Policy Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	hasBFDMultiHop := len(d.Get("bfd_multi_hop_settings").([]interface{})) > 0
	hasBGP := len(d.Get("bgp_node_settings").([]interface{})) > 0
	hasASPath := d.Get("as_path_multipath_relax") != nil

	if !hasBFDMultiHop && !hasBGP && !hasASPath {
		return fmt.Errorf("at least one of 'bfd_multi_hop_settings', 'bgp_node_settings', or 'as_path_multipath_relax' must be specified")
	}

	payload := map[string]interface{}{
		"name": d.Get("name").(string),
	}

	if description, ok := d.GetOk("description"); ok {
		payload["description"] = description.(string)
	}

	if asPath, ok := d.GetOk("as_path_multipath_relax"); ok {
		payload["asPathPol"] = map[string]interface{}{
			"asPathMultipathRelax": asPath.(bool),
		}
	}

	if bfdMultiHopRaw, ok := d.GetOk("bfd_multi_hop_settings"); ok {
		if bfdMultiHopPayload := buildNodeBFDMultiHopPayload(bfdMultiHopRaw.([]interface{})); bfdMultiHopPayload != nil {
			payload["bfdMultiHopPol"] = bfdMultiHopPayload
		}
	}

	if bgpRaw, ok := d.GetOk("bgp_node_settings"); ok {
		if bgpPayload := buildBGPNodeSettingsPayload(bgpRaw.([]interface{})); bgpPayload != nil {
			payload["bgpTimerPol"] = bgpPayload
		}
	}

	payloadModel := models.GetPatchPayload("add", "/tenantPolicyTemplate/template/l3OutNodePolGroups/-", payload)
	templateId := d.Get("template_id").(string)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/L3OutNodeRoutingPolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO L3Out Node Routing Policy Resource - Create Complete: %v", d.Id())
	return resourceMSOL3OutNodeRoutingPolicyRead(d, m)
}

func resourceMSOL3OutNodeRoutingPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO L3Out Node Routing Policy Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return err
	}

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyName, err := GetPolicyNameFromResourceId(d.Id(), "L3OutNodeRoutingPolicy")
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "l3OutNodePolGroups")
	if err != nil {
		return err
	}

	setL3OutNodeRoutingPolicyData(d, policy, templateId)
	log.Printf("[DEBUG] MSO L3Out Node Routing Policy Resource - Read Complete: %v", d.Id())
	return nil
}

func resourceMSOL3OutNodeRoutingPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO L3Out Node Routing Policy Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "l3OutNodePolGroups")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/tenantPolicyTemplate/template/l3OutNodePolGroups/%d", policyIndex)

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

	if d.HasChange("as_path_multipath_relax") {
		_, newVal := d.GetChange("as_path_multipath_relax")

		if newVal == nil {
			err := addPatchPayloadToContainer(payloadCont, "remove", fmt.Sprintf("%s/asPathPol", updatePath), nil)
			if err != nil {
				return err
			}
		} else {
			asPathPayload := map[string]interface{}{
				"asPathMultipathRelax": newVal.(bool),
			}
			err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/asPathPol", updatePath), asPathPayload)
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("bfd_multi_hop_settings") {
		_, newVal := d.GetChange("bfd_multi_hop_settings")
		bfdMultiHopRaw := newVal.([]interface{})

		if len(bfdMultiHopRaw) == 0 {
			err := addPatchPayloadToContainer(payloadCont, "remove", fmt.Sprintf("%s/bfdMultiHopPol", updatePath), nil)
			if err != nil {
				return err
			}
		} else {
			if bfdMultiHopPayload := buildNodeBFDMultiHopPayload(bfdMultiHopRaw); bfdMultiHopPayload != nil {
				err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/bfdMultiHopPol", updatePath), bfdMultiHopPayload)
				if err != nil {
					return err
				}
			}
		}
	}

	if d.HasChange("bgp_node_settings") {
		_, newVal := d.GetChange("bgp_node_settings")
		bgpRaw := newVal.([]interface{})

		if len(bgpRaw) == 0 {
			err := addPatchPayloadToContainer(payloadCont, "remove", fmt.Sprintf("%s/bgpTimerPol", updatePath), nil)
			if err != nil {
				return err
			}
		} else {
			if bgpPayload := buildBGPNodeSettingsPayload(bgpRaw); bgpPayload != nil {
				err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/bgpTimerPol", updatePath), bgpPayload)
				if err != nil {
					return err
				}
			}
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/L3OutNodeRoutingPolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO L3Out Node Routing Policy Resource - Update Complete: %v", d.Id())
	return resourceMSOL3OutNodeRoutingPolicyRead(d, m)
}

func resourceMSOL3OutNodeRoutingPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO L3Out Node Routing Policy Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "l3OutNodePolGroups")
	if err != nil {
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/tenantPolicyTemplate/template/l3OutNodePolGroups/%d", policyIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO L3Out Node Routing Policy Resource - Delete Complete: %v", d.Id())
	return nil
}
