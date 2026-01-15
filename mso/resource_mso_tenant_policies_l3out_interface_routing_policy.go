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

func resourceMSOL3OutInterfaceRoutingPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOL3OutInterfaceRoutingPolicyCreate,
		Read:   resourceMSOL3OutInterfaceRoutingPolicyRead,
		Update: resourceMSOL3OutInterfaceRoutingPolicyUpdate,
		Delete: resourceMSOL3OutInterfaceRoutingPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOL3OutInterfaceRoutingPolicyImport,
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
				Description:  "The name of the L3Out Interface Routing Policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the L3Out Interface Routing Policy.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the L3Out Interface Routing Policy.",
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
			"bfd_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "BFD configuration.",
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
							ValidateFunc: validation.IntBetween(50, 999),
							Description:  "Minimum receive interval in microseconds. Default: 50 when unset during creation. Range: 50-999.",
						},
						"min_transmit_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(50, 999),
							Description:  "Minimum transmit interval in microseconds. Default: 50 when unset during creation. Range: 50-999.",
						},
						"echo_receive_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(50, 999),
							Description:  "Echo receive interval in microseconds. Default: 50 when unset during creation. Range: 50-999.",
						},
						"echo_admin_state": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"enabled", "disabled",
							}, false),
							Description: "Echo administrative state. Default: enabled when unset during creation.",
						},
						"interface_control": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Interface control. Default: false (disabled) when unset during creation.",
						},
					},
				},
			},
			"ospf_interface_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "OSPF interface configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"broadcast", "point_to_point",
							}, false),
							Description: "Network type. Default: broadcast when unset during creation.",
						},
						"priority": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(0, 255),
							Description:  "OSPF priority. Default: 1 when unset during creation. Range: 0-255.",
						},
						"cost_of_interface": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(0, 65535),
							Description:  "OSPF cost. Default: 0 when unset during creation. Range: 0-65535.",
						},
						"hello_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
							Description:  "Hello interval in seconds. Default: 10 when unset during creation. Range: 1-65535.",
						},
						"dead_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
							Description:  "Dead interval in seconds. Default: 40 when unset during creation. Range: 1-65535.",
						},
						"retransmit_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
							Description:  "Retransmit interval in seconds. Default: 5 when unset during creation. Range: 1-65535.",
						},
						"transmit_delay": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(1, 450),
							Description:  "Transmit delay in seconds. Default: 1 when unset during creation. Range: 1-450.",
						},
						"advertise_subnet": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Advertise subnet. Default: false (disabled) when unset during creation.",
						},
						"bfd": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Enable BFD. Default: false (disabled) when unset during creation.",
						},
						"mtu_ignore": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Ignore MTU. Default: false (disabled) when unset during creation.",
						},
						"passive_participation": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Passive participation. Default: false (disabled) when unset during creation.",
						},
					},
				},
			},
		},
	}
}

func buildBFDMultiHopPayload(bfdMultiHop []interface{}) map[string]interface{} {
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

func buildBFDPayload(bfdSettings []interface{}) map[string]interface{} {
	if len(bfdSettings) == 0 {
		return nil
	}

	settings := bfdSettings[0].(map[string]interface{})
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
	if echoRxInterval, ok := settings["echo_receive_interval"].(int); ok && echoRxInterval != 0 {
		payload["echoRxInterval"] = echoRxInterval
	}
	if echoAdminState, ok := settings["echo_admin_state"].(string); ok && echoAdminState != "" {
		payload["echoAdminState"] = echoAdminState
	}
	if ifControl, ok := settings["interface_control"].(bool); ok {
		payload["ifControl"] = ifControl
	}

	return payload
}

func buildOSPFInterfacePayload(ospfSettings []interface{}) map[string]interface{} {
	if len(ospfSettings) == 0 {
		return nil
	}

	settings := ospfSettings[0].(map[string]interface{})
	payload := make(map[string]interface{})
	ifControl := make(map[string]interface{})

	if networkType, ok := settings["network_type"].(string); ok && networkType != "" {
		if networkType == "point_to_point" {
			payload["networkType"] = "pointToPoint"
		} else {
			payload["networkType"] = networkType
		}
	}
	if priority, ok := settings["priority"].(int); ok {
		payload["prio"] = priority
	}
	if cost, ok := settings["cost_of_interface"].(int); ok {
		payload["cost"] = cost
	}
	if helloInterval, ok := settings["hello_interval"].(int); ok && helloInterval != 0 {
		payload["helloInterval"] = helloInterval
	}
	if deadInterval, ok := settings["dead_interval"].(int); ok && deadInterval != 0 {
		payload["deadInterval"] = deadInterval
	}
	if retransmitInterval, ok := settings["retransmit_interval"].(int); ok && retransmitInterval != 0 {
		payload["retransmitInterval"] = retransmitInterval
	}
	if transmitDelay, ok := settings["transmit_delay"].(int); ok && transmitDelay != 0 {
		payload["transmitDelay"] = transmitDelay
	}

	if advertiseSubnet, ok := settings["advertise_subnet"].(bool); ok {
		ifControl["advertiseSubnet"] = advertiseSubnet
	}
	if bfd, ok := settings["bfd"].(bool); ok {
		ifControl["bfd"] = bfd
	}
	if mtuIgnore, ok := settings["mtu_ignore"].(bool); ok {
		ifControl["ignoreMtu"] = mtuIgnore
	}
	if passiveParticipation, ok := settings["passive_participation"].(bool); ok {
		ifControl["passiveParticipation"] = passiveParticipation
	}

	if len(ifControl) > 0 {
		payload["ifControl"] = ifControl
	}

	return payload
}

func setL3OutInterfaceRoutingPolicyData(d *schema.ResourceData, response *container.Container, templateId string) error {
	d.SetId(fmt.Sprintf("templateId/%s/L3OutInterfaceRoutingPolicy/%s", templateId, models.StripQuotes(response.S("name").String())))
	d.Set("template_id", templateId)
	d.Set("name", models.StripQuotes(response.S("name").String()))
	d.Set("description", models.StripQuotes(response.S("description").String()))
	d.Set("uuid", models.StripQuotes(response.S("uuid").String()))

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

	if response.Exists("bfdPol") {
		bfdPol := response.S("bfdPol")
		bfdMap := map[string]interface{}{}

		if bfdPol.Exists("adminState") {
			bfdMap["admin_state"] = models.StripQuotes(bfdPol.S("adminState").String())
		}
		if bfdPol.Exists("detectionMultiplier") {
			if val, ok := bfdPol.S("detectionMultiplier").Data().(float64); ok {
				bfdMap["detection_multiplier"] = int(val)
			}
		}
		if bfdPol.Exists("minRxInterval") {
			if val, ok := bfdPol.S("minRxInterval").Data().(float64); ok {
				bfdMap["min_receive_interval"] = int(val)
			}
		}
		if bfdPol.Exists("minTxInterval") {
			if val, ok := bfdPol.S("minTxInterval").Data().(float64); ok {
				bfdMap["min_transmit_interval"] = int(val)
			}
		}
		if bfdPol.Exists("echoRxInterval") {
			if val, ok := bfdPol.S("echoRxInterval").Data().(float64); ok {
				bfdMap["echo_receive_interval"] = int(val)
			}
		}
		if bfdPol.Exists("echoAdminState") {
			bfdMap["echo_admin_state"] = models.StripQuotes(bfdPol.S("echoAdminState").String())
		}
		if bfdPol.Exists("ifControl") {
			if val, ok := bfdPol.S("ifControl").Data().(bool); ok {
				bfdMap["interface_control"] = val
			}
		}

		d.Set("bfd_settings", []interface{}{bfdMap})
	} else {
		d.Set("bfd_settings", []interface{}{})
	}

	if response.Exists("ospfIntfPol") {
		ospfPol := response.S("ospfIntfPol")
		ospfMap := map[string]interface{}{}

		if ospfPol.Exists("networkType") {
			networkType := models.StripQuotes(ospfPol.S("networkType").String())
			if networkType == "pointToPoint" {
				ospfMap["network_type"] = "point_to_point"
			} else {
				ospfMap["network_type"] = networkType
			}
		}
		if ospfPol.Exists("prio") {
			if val, ok := ospfPol.S("prio").Data().(float64); ok {
				ospfMap["priority"] = int(val)
			}
		}
		if ospfPol.Exists("cost") {
			if val, ok := ospfPol.S("cost").Data().(float64); ok {
				ospfMap["cost_of_interface"] = int(val)
			}
		}
		if ospfPol.Exists("helloInterval") {
			if val, ok := ospfPol.S("helloInterval").Data().(float64); ok {
				ospfMap["hello_interval"] = int(val)
			}
		}
		if ospfPol.Exists("deadInterval") {
			if val, ok := ospfPol.S("deadInterval").Data().(float64); ok {
				ospfMap["dead_interval"] = int(val)
			}
		}
		if ospfPol.Exists("retransmitInterval") {
			if val, ok := ospfPol.S("retransmitInterval").Data().(float64); ok {
				ospfMap["retransmit_interval"] = int(val)
			}
		}
		if ospfPol.Exists("transmitDelay") {
			if val, ok := ospfPol.S("transmitDelay").Data().(float64); ok {
				ospfMap["transmit_delay"] = int(val)
			}
		}

		if ospfPol.Exists("ifControl") {
			ifControl := ospfPol.S("ifControl")

			if ifControl.Exists("advertiseSubnet") {
				if val, ok := ifControl.S("advertiseSubnet").Data().(bool); ok {
					ospfMap["advertise_subnet"] = val
				}
			}
			if ifControl.Exists("bfd") {
				if val, ok := ifControl.S("bfd").Data().(bool); ok {
					ospfMap["bfd"] = val
				}
			}
			if ifControl.Exists("ignoreMtu") {
				if val, ok := ifControl.S("ignoreMtu").Data().(bool); ok {
					ospfMap["mtu_ignore"] = val
				}
			}
			if ifControl.Exists("passiveParticipation") {
				if val, ok := ifControl.S("passiveParticipation").Data().(bool); ok {
					ospfMap["passive_participation"] = val
				}
			}
		}

		d.Set("ospf_interface_settings", []interface{}{ospfMap})
	} else {
		d.Set("ospf_interface_settings", []interface{}{})
	}

	return nil
}

func resourceMSOL3OutInterfaceRoutingPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO L3Out Interface Routing Policy Resource - Beginning Import: %v", d.Id())
	resourceMSOL3OutInterfaceRoutingPolicyRead(d, m)
	log.Printf("[DEBUG] MSO L3Out Interface Routing Policy Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOL3OutInterfaceRoutingPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO L3Out Interface Routing Policy Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	payload := map[string]interface{}{
		"name": d.Get("name").(string),
	}

	if description, ok := d.GetOk("description"); ok {
		payload["description"] = description.(string)
	}

	if bfdMultiHopRaw, ok := d.GetOk("bfd_multi_hop_settings"); ok {
		if bfdMultiHopPayload := buildBFDMultiHopPayload(bfdMultiHopRaw.([]interface{})); bfdMultiHopPayload != nil {
			payload["bfdMultiHopPol"] = bfdMultiHopPayload
		}
	}

	if bfdRaw, ok := d.GetOk("bfd_settings"); ok {
		if bfdPayload := buildBFDPayload(bfdRaw.([]interface{})); bfdPayload != nil {
			payload["bfdPol"] = bfdPayload
		}
	}

	if ospfRaw, ok := d.GetOk("ospf_interface_settings"); ok {
		if ospfPayload := buildOSPFInterfacePayload(ospfRaw.([]interface{})); ospfPayload != nil {
			payload["ospfIntfPol"] = ospfPayload
		}
	}

	payloadModel := models.GetPatchPayload("add", "/tenantPolicyTemplate/template/l3OutIntfPolGroups/-", payload)
	templateId := d.Get("template_id").(string)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/L3OutInterfaceRoutingPolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO L3Out Interface Routing Policy Resource - Create Complete: %v", d.Id())
	return resourceMSOL3OutInterfaceRoutingPolicyRead(d, m)
}

func resourceMSOL3OutInterfaceRoutingPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO L3Out Interface Routing Policy Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return err
	}

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyName, err := GetPolicyNameFromResourceId(d.Id(), "L3OutInterfaceRoutingPolicy")
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "l3OutIntfPolGroups")
	if err != nil {
		return err
	}

	setL3OutInterfaceRoutingPolicyData(d, policy, templateId)
	log.Printf("[DEBUG] MSO L3Out Interface Routing Policy Resource - Read Complete: %v", d.Id())
	return nil
}

func resourceMSOL3OutInterfaceRoutingPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO L3Out Interface Routing Policy Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "l3OutIntfPolGroups")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/tenantPolicyTemplate/template/l3OutIntfPolGroups/%d", policyIndex)

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

	if d.HasChange("bfd_multi_hop_settings") {
		bfdMultiHopRaw := d.Get("bfd_multi_hop_settings").([]interface{})

		if len(bfdMultiHopRaw) == 0 {
			err := addPatchPayloadToContainer(payloadCont, "remove", fmt.Sprintf("%s/bfdMultiHopPol", updatePath), nil)
			if err != nil {
				return err
			}
		} else {
			if bfdMultiHopPayload := buildBFDMultiHopPayload(bfdMultiHopRaw); bfdMultiHopPayload != nil {
				err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/bfdMultiHopPol", updatePath), bfdMultiHopPayload)
				if err != nil {
					return err
				}
			}
		}
	}

	if d.HasChange("bfd_settings") {
		bfdRaw := d.Get("bfd_settings").([]interface{})

		if len(bfdRaw) == 0 {
			err := addPatchPayloadToContainer(payloadCont, "remove", fmt.Sprintf("%s/bfdPol", updatePath), nil)
			if err != nil {
				return err
			}
		} else {
			if bfdPayload := buildBFDPayload(bfdRaw); bfdPayload != nil {
				err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/bfdPol", updatePath), bfdPayload)
				if err != nil {
					return err
				}
			}
		}
	}

	if d.HasChange("ospf_interface_settings") {
		ospfRaw := d.Get("ospf_interface_settings").([]interface{})

		if len(ospfRaw) == 0 {
			err := addPatchPayloadToContainer(payloadCont, "remove", fmt.Sprintf("%s/ospfIntfPol", updatePath), nil)
			if err != nil {
				return err
			}
		} else {
			if ospfPayload := buildOSPFInterfacePayload(ospfRaw); ospfPayload != nil {
				err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/ospfIntfPol", updatePath), ospfPayload)
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

	d.SetId(fmt.Sprintf("templateId/%s/L3OutInterfaceRoutingPolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO L3Out Interface Routing Policy Resource - Update Complete: %v", d.Id())
	return resourceMSOL3OutInterfaceRoutingPolicyRead(d, m)
}

func resourceMSOL3OutInterfaceRoutingPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO L3Out Interface Routing Policy Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "l3OutIntfPolGroups")
	if err != nil {
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/tenantPolicyTemplate/template/l3OutIntfPolGroups/%d", policyIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO L3Out Interface Routing Policy Resource - Delete Complete: %v", d.Id())
	return nil
}
