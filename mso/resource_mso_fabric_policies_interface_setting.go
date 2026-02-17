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

var (
	portChannelModeMap = map[string]string{
		"lacp_active":                   "active",
		"lacp_passive":                  "passive",
		"static_channel_mode_on":        "off",
		"mac_pinning":                   "mac-pin",
		"mac_pinning_physical_nic_load": "mac-pin-nicload",
		"use_explicit_failover_order":   "explicit-failover",
	}

	controlMap = map[string]string{
		"fast_sel_hot_stdby": "fast-sel-hot-stdby",
		"graceful_conv":      "graceful-conv",
		"susp_individual":    "susp-individual",
		"load_defer":         "load-defer",
		"symmetric_hash":     "symmetric-hash",
	}

	linkLevelFecMap = map[string]string{
		"inherit":       "inherit",
		"cl74_fc_fec":   "cl74-fc-fec",
		"cl91_rs_fec":   "cl91-rs-fec",
		"cons16_rs_fec": "cons16-rs-fec",
		"ieee_rs_fec":   "ieee-rs-fec",
		"kp_fec":        "kp-fec",
		"disable_fec":   "disable-fec",
	}

	l2InterfaceQinqMap = map[string]string{
		"double_q_tag_port": "doubleQtagPort",
		"core_port":         "corePort",
		"edge_port":         "edgePort",
		"disabled":          "disabled",
	}

	loadBalanceHashingMap = map[string]string{
		"destination_ip":         "dst-ip",
		"layer_4_destination_ip": "l4-dst-port",
		"layer_4_source_ip":      "l4-src-port",
		"source_ip":              "src-ip",
	}
)

func resourceMSOInterfaceSetting() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOInterfaceSettingCreate,
		Read:   resourceMSOInterfaceSettingRead,
		Update: resourceMSOInterfaceSettingUpdate,
		Delete: resourceMSOInterfaceSettingDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOInterfaceSettingImport,
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"physical", "portchannel"}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cdp_admin_state": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
			},
			"lldp_receive_state": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
			},
			"lldp_transmit_state": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
			},
			"llfc_receive_state": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
			},
			"llfc_transmit_state": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
			},
			"pfc_admin_state": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"off", "on", "auto"}, false),
			},
			"l2_interface_qinq": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(getMapKeys(l2InterfaceQinqMap), false),
			},
			"l2_interface_reflective_relay": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
			},
			"vlan_scope": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"global", "portlocal"}, false),
			},
			"stp_bpdu_filter": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
			},
			"stp_bpdu_guard": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
			},
			"speed": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"100M", "1G", "10G", "25G", "40G", "50G", "100G", "200G", "400G", "inherit"}, false),
			},
			"auto_negotiation": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"on", "off", "on_enforce"}, false),
			},
			"link_level_bring_up_delay": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 10000),
			},
			"link_level_debounce_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 5000),
			},
			"link_level_fec": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(getMapKeys(linkLevelFecMap), false),
			},
			"mcp_admin_state": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
			},
			"mcp_strict_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
			},
			"mcp_initial_delay_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 1800),
			},
			"mcp_transmission_frequency_sec": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 300),
			},
			"mcp_transmission_frequency_msec": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 999),
			},
			"mcp_grace_period_sec": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 300),
			},
			"mcp_grace_period_msec": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 999),
			},
			"port_channel_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(getMapKeys(portChannelModeMap), false),
			},
			"controls": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(getMapKeys(controlMap), false),
				},
			},
			"port_channel_min_links": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 16),
			},
			"port_channel_max_links": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 64),
			},
			"load_balance_hashing": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(getMapKeys(loadBalanceHashingMap), false),
			},
			"synce": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domains": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"access_macsec_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func setInterfaceSettingData(d *schema.ResourceData, response *container.Container, templateId string) error {
	log.Printf("[DEBUG] MSO Interface Setting Resource - Beginning setInterfaceSettingData")

	name := models.StripQuotes(response.S("name").String())
	d.Set("name", name)
	uuid := models.StripQuotes(response.S("uuid").String())
	d.SetId(fmt.Sprintf("templateId/%s/InterfaceSetting/%s", templateId, name))
	d.Set("template_id", templateId)
	d.Set("uuid", uuid)
	d.Set("type", models.StripQuotes(response.S("type").String()))

	if response.S("description").String() != "" {
		d.Set("description", models.StripQuotes(response.S("description").String()))
	}

	if response.S("cdp").Data() != nil {
		if adminState := models.StripQuotes(response.S("cdp").S("adminState").String()); adminState != "" {
			d.Set("cdp_admin_state", adminState)
		}
	}

	if response.S("lldp").Data() != nil {
		if receiveState := models.StripQuotes(response.S("lldp").S("receiveState").String()); receiveState != "" {
			d.Set("lldp_receive_state", receiveState)
		}
		if transmitState := models.StripQuotes(response.S("lldp").S("transmitState").String()); transmitState != "" {
			d.Set("lldp_transmit_state", transmitState)
		}
	}

	if response.S("llfc").Data() != nil {
		if receiveState := models.StripQuotes(response.S("llfc").S("receiveState").String()); receiveState != "" {
			d.Set("llfc_receive_state", receiveState)
		}
		if transmitState := models.StripQuotes(response.S("llfc").S("transmitState").String()); transmitState != "" {
			d.Set("llfc_transmit_state", transmitState)
		}
	}

	if response.S("pfc").Data() != nil {
		if adminState := models.StripQuotes(response.S("pfc").S("adminState").String()); adminState != "" {
			d.Set("pfc_admin_state", adminState)
		}
	}

	if response.S("l2Interface").Data() != nil {
		if qinq := models.StripQuotes(response.S("l2Interface").S("qinq").String()); qinq != "" {
			d.Set("l2_interface_qinq", getKeyByValue(l2InterfaceQinqMap, qinq))
		}
		if reflectiveRelay := models.StripQuotes(response.S("l2Interface").S("reflectiveRelay").String()); reflectiveRelay != "" {
			d.Set("l2_interface_reflective_relay", reflectiveRelay)
		}
		if vlanScope := models.StripQuotes(response.S("l2Interface").S("vlanScope").String()); vlanScope != "" {
			d.Set("vlan_scope", vlanScope)
		}
	}

	if response.S("stp").Data() != nil {
		if bpduFilter := models.StripQuotes(response.S("stp").S("bpduFilterEnabled").String()); bpduFilter != "" {
			d.Set("stp_bpdu_filter", bpduFilter)
		}
		if bpduGuard := models.StripQuotes(response.S("stp").S("bpduGuardEnabled").String()); bpduGuard != "" {
			d.Set("stp_bpdu_guard", bpduGuard)
		}
	}

	if response.S("linkLevel").Data() != nil {
		if speed := models.StripQuotes(response.S("linkLevel").S("speed").String()); speed != "" {
			d.Set("speed", speed)
		}
		if autoNeg := models.StripQuotes(response.S("linkLevel").S("autoNegotiation").String()); autoNeg != "" {
			if autoNeg == "on-enforce" {
				autoNeg = "on_enforce"
			}
			d.Set("auto_negotiation", autoNeg)
		}
		if bringUpDelay := response.S("linkLevel").S("bringUpDelay").Data(); bringUpDelay != nil {
			d.Set("link_level_bring_up_delay", int(bringUpDelay.(float64)))
		}
		if debounceInterval := response.S("linkLevel").S("debounceInterval").Data(); debounceInterval != nil {
			d.Set("link_level_debounce_interval", int(debounceInterval.(float64)))
		}
		if fec := models.StripQuotes(response.S("linkLevel").S("fec").String()); fec != "" {
			d.Set("link_level_fec", getKeyByValue(linkLevelFecMap, fec))
		}
	}

	if response.S("mcp").Data() != nil {
		if adminState := models.StripQuotes(response.S("mcp").S("adminState").String()); adminState != "" {
			d.Set("mcp_admin_state", adminState)
		}
		if mcpMode := models.StripQuotes(response.S("mcp").S("mcpMode").String()); mcpMode != "" {
			d.Set("mcp_strict_mode", mcpMode)
		}
		if initialDelay := response.S("mcp").S("initialDelayTime").Data(); initialDelay != nil {
			d.Set("mcp_initial_delay_time", int(initialDelay.(float64)))
		}
		if txFreq := response.S("mcp").S("txFreq").Data(); txFreq != nil {
			d.Set("mcp_transmission_frequency_sec", int(txFreq.(float64)))
		}
		if txFreqMsec := response.S("mcp").S("txFreqMsec").Data(); txFreqMsec != nil {
			d.Set("mcp_transmission_frequency_msec", int(txFreqMsec.(float64)))
		}
		if gracePeriod := response.S("mcp").S("gracePeriod").Data(); gracePeriod != nil {
			d.Set("mcp_grace_period_sec", int(gracePeriod.(float64)))
		}
		if gracePeriodMsec := response.S("mcp").S("gracePeriodMsec").Data(); gracePeriodMsec != nil {
			d.Set("mcp_grace_period_msec", int(gracePeriodMsec.(float64)))
		}
	}

	if response.S("portChannelPolicy").Data() != nil {
		if mode := models.StripQuotes(response.S("portChannelPolicy").S("mode").String()); mode != "" {
			d.Set("port_channel_mode", getKeyByValue(portChannelModeMap, mode))
		}
		if controls := response.S("portChannelPolicy").S("control").Data(); controls != nil {
			controlList := make([]string, 0)
			for _, control := range controls.([]interface{}) {
				controlList = append(controlList, getKeyByValue(controlMap, control.(string)))
			}
			d.Set("controls", controlList)
		}
		if minLinks := response.S("portChannelPolicy").S("minLinks").Data(); minLinks != nil {
			d.Set("port_channel_min_links", int(minLinks.(float64)))
		}
		if maxLinks := response.S("portChannelPolicy").S("maxLinks").Data(); maxLinks != nil {
			d.Set("port_channel_max_links", int(maxLinks.(float64)))
		}
		if hashFields := models.StripQuotes(response.S("portChannelPolicy").S("hashFields").String()); hashFields != "" {
			d.Set("load_balance_hashing", getKeyByValue(loadBalanceHashingMap, hashFields))
		}
	}

	if syncE := models.StripQuotes(response.S("syncEthPolicy").String()); syncE != "" && syncE != "{}" {
		d.Set("synce", syncE)
	}

	if response.S("domains").Data() != nil {
		d.Set("domains", response.S("domains").Data())
	}

	if macsecPolicy := models.StripQuotes(response.S("accessMACsecPolicy").String()); macsecPolicy != "" && macsecPolicy != "{}" {
		d.Set("access_macsec_policy", macsecPolicy)
	}

	log.Printf("[DEBUG] MSO Interface Setting Resource - setInterfaceSettingData Complete")

	return nil
}

func resourceMSOInterfaceSettingImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO Interface Setting Resource - Beginning Import: %v", d.Id())
	resourceMSOInterfaceSettingRead(d, m)
	log.Printf("[DEBUG] MSO Interface Setting Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOInterfaceSettingCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Interface Setting Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)

	payload := buildInterfaceSettingPayload(d)

	payloadModel := models.GetPatchPayload("add", "/fabricPolicyTemplate/template/interfacePolicyGroups/-", payload)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/InterfaceSetting/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO Interface Setting Resource - Create Complete: %v", d.Id())
	return resourceMSOInterfaceSettingRead(d, m)
}

func resourceMSOInterfaceSettingUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Interface Setting Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	interfaceSettingIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "fabricPolicyTemplate", "template", "interfacePolicyGroups")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/fabricPolicyTemplate/template/interfacePolicyGroups/%d", interfaceSettingIndex)

	payloadCont := container.New()
	payloadCont.Array()

	if d.HasChange("type") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/type", updatePath), d.Get("type").(string))
		if err != nil {
			return err
		}
	}

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

	if d.HasChange("cdp_admin_state") {
		cdpAdminState, ok := d.GetOk("cdp_admin_state")
		if !ok {
			err := addPatchPayloadToContainer(payloadCont, "remove", fmt.Sprintf("%s/cdp", updatePath), nil)
			if err != nil {
				return err
			}
		} else {
			cdpPayload := map[string]interface{}{
				"adminState": cdpAdminState.(string),
			}
			err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/cdp", updatePath), cdpPayload)
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("lldp_receive_state") || d.HasChange("lldp_transmit_state") {
		lldpPayload := make(map[string]interface{})
		if lldpReceive, ok := d.GetOk("lldp_receive_state"); ok {
			lldpPayload["receiveState"] = lldpReceive.(string)
		}
		if lldpTransmit, ok := d.GetOk("lldp_transmit_state"); ok {
			lldpPayload["transmitState"] = lldpTransmit.(string)
		}
		if len(lldpPayload) == 0 {
			err := addPatchPayloadToContainer(payloadCont, "remove", fmt.Sprintf("%s/lldp", updatePath), nil)
			if err != nil {
				return err
			}
		} else {
			err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/lldp", updatePath), lldpPayload)
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("llfc_receive_state") || d.HasChange("llfc_transmit_state") {
		llfcPayload := make(map[string]interface{})
		if llfcReceive, ok := d.GetOk("llfc_receive_state"); ok {
			llfcPayload["receiveState"] = llfcReceive.(string)
		}
		if llfcTransmit, ok := d.GetOk("llfc_transmit_state"); ok {
			llfcPayload["transmitState"] = llfcTransmit.(string)
		}
		if len(llfcPayload) == 0 {
			err := addPatchPayloadToContainer(payloadCont, "remove", fmt.Sprintf("%s/llfc", updatePath), nil)
			if err != nil {
				return err
			}
		} else {
			err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/llfc", updatePath), llfcPayload)
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("pfc_admin_state") {
		pfcAdminState, ok := d.GetOk("pfc_admin_state")
		if !ok {
			err := addPatchPayloadToContainer(payloadCont, "remove", fmt.Sprintf("%s/pfc", updatePath), nil)
			if err != nil {
				return err
			}
		} else {
			pfcPayload := map[string]interface{}{
				"adminState": pfcAdminState.(string),
			}
			err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/pfc", updatePath), pfcPayload)
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("l2_interface_qinq") || d.HasChange("l2_interface_reflective_relay") || d.HasChange("vlan_scope") {
		l2InterfacePayload := make(map[string]interface{})
		if l2QinQ, ok := d.GetOk("l2_interface_qinq"); ok {
			qinqValue := l2QinQ.(string)
			if mappedValue, exists := l2InterfaceQinqMap[qinqValue]; exists {
				l2InterfacePayload["qinq"] = mappedValue
			} else {
				l2InterfacePayload["qinq"] = qinqValue
			}
		}
		if l2ReflectiveRelay, ok := d.GetOk("l2_interface_reflective_relay"); ok {
			l2InterfacePayload["reflectiveRelay"] = l2ReflectiveRelay.(string)
		}
		if vlanScope, ok := d.GetOk("vlan_scope"); ok {
			l2InterfacePayload["vlanScope"] = vlanScope.(string)
		}
		if len(l2InterfacePayload) == 0 {
			err := addPatchPayloadToContainer(payloadCont, "remove", fmt.Sprintf("%s/l2Interface", updatePath), nil)
			if err != nil {
				return err
			}
		} else {
			err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/l2Interface", updatePath), l2InterfacePayload)
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("stp_bpdu_filter") || d.HasChange("stp_bpdu_guard") {
		stpPayload := make(map[string]interface{})
		if stpBpduFilter, ok := d.GetOk("stp_bpdu_filter"); ok {
			stpPayload["bpduFilterEnabled"] = stpBpduFilter.(string)
		}
		if stpBpduGuard, ok := d.GetOk("stp_bpdu_guard"); ok {
			stpPayload["bpduGuardEnabled"] = stpBpduGuard.(string)
		}
		if len(stpPayload) == 0 {
			err := addPatchPayloadToContainer(payloadCont, "remove", fmt.Sprintf("%s/stp", updatePath), nil)
			if err != nil {
				return err
			}
		} else {
			err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/stp", updatePath), stpPayload)
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("speed") || d.HasChange("auto_negotiation") || d.HasChange("link_level_bring_up_delay") || d.HasChange("link_level_debounce_interval") || d.HasChange("link_level_fec") {
		linkLevelPayload := make(map[string]interface{})
		if speed, ok := d.GetOk("speed"); ok {
			linkLevelPayload["speed"] = speed.(string)
		}
		if autoNeg, ok := d.GetOk("auto_negotiation"); ok {
			if autoNeg == "on_enforce" {
				autoNeg = "on-enforce"
			}
			linkLevelPayload["autoNegotiation"] = autoNeg.(string)
		}
		if bringUpDelay, ok := d.GetOk("link_level_bring_up_delay"); ok {
			linkLevelPayload["bringUpDelay"] = bringUpDelay.(int)
		}
		if debounceInterval, ok := d.GetOk("link_level_debounce_interval"); ok {
			linkLevelPayload["debounceInterval"] = debounceInterval.(int)
		}
		if fec, ok := d.GetOk("link_level_fec"); ok {
			fecValue := fec.(string)
			if mappedValue, exists := linkLevelFecMap[fecValue]; exists {
				linkLevelPayload["fec"] = mappedValue
			} else {
				linkLevelPayload["fec"] = fecValue
			}
		}
		if len(linkLevelPayload) == 0 {
			err := addPatchPayloadToContainer(payloadCont, "remove", fmt.Sprintf("%s/linkLevel", updatePath), nil)
			if err != nil {
				return err
			}
		} else {
			err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/linkLevel", updatePath), linkLevelPayload)
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("mcp_admin_state") || d.HasChange("mcp_strict_mode") || d.HasChange("mcp_initial_delay_time") || d.HasChange("mcp_transmission_frequency_sec") || d.HasChange("mcp_transmission_frequency_msec") || d.HasChange("mcp_grace_period_sec") || d.HasChange("mcp_grace_period_msec") {
		mcpPayload := make(map[string]interface{})
		if mcpAdminState, ok := d.GetOk("mcp_admin_state"); ok {
			mcpPayload["adminState"] = mcpAdminState.(string)
		}
		if mcpStrictMode, ok := d.GetOk("mcp_strict_mode"); ok {
			mcpPayload["mcpMode"] = mcpStrictMode.(string)
		}
		if mcpInitialDelay, ok := d.GetOk("mcp_initial_delay_time"); ok {
			mcpPayload["initialDelayTime"] = mcpInitialDelay.(int)
		}
		if mcpTxFreq, ok := d.GetOk("mcp_transmission_frequency_sec"); ok {
			mcpPayload["txFreq"] = mcpTxFreq.(int)
		}
		if mcpTxFreqMsec, ok := d.GetOk("mcp_transmission_frequency_msec"); ok {
			mcpPayload["txFreqMsec"] = mcpTxFreqMsec.(int)
		}
		if mcpGracePeriod, ok := d.GetOk("mcp_grace_period_sec"); ok {
			mcpPayload["gracePeriod"] = mcpGracePeriod.(int)
		}
		if mcpGracePeriodMsec, ok := d.GetOk("mcp_grace_period_msec"); ok {
			mcpPayload["gracePeriodMsec"] = mcpGracePeriodMsec.(int)
		}
		if len(mcpPayload) == 0 {
			err := addPatchPayloadToContainer(payloadCont, "remove", fmt.Sprintf("%s/mcp", updatePath), nil)
			if err != nil {
				return err
			}
		} else {
			err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/mcp", updatePath), mcpPayload)
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("port_channel_mode") || d.HasChange("controls") || d.HasChange("port_channel_min_links") || d.HasChange("port_channel_max_links") || d.HasChange("load_balance_hashing") {
		portChannelPayload := make(map[string]interface{})
		if portChannelMode, ok := d.GetOk("port_channel_mode"); ok {
			modeValue := portChannelMode.(string)
			if mappedValue, exists := portChannelModeMap[modeValue]; exists {
				portChannelPayload["mode"] = mappedValue
			} else {
				portChannelPayload["mode"] = modeValue
			}
		}
		if controls, ok := d.GetOk("controls"); ok {
			controlsList := controls.(*schema.Set).List()
			mappedControls := make([]string, 0)
			for _, control := range controlsList {
				controlValue := control.(string)
				if mappedValue, exists := controlMap[controlValue]; exists {
					mappedControls = append(mappedControls, mappedValue)
				} else {
					mappedControls = append(mappedControls, controlValue)
				}
			}
			portChannelPayload["control"] = mappedControls
		}
		if portChannelMinLinks, ok := d.GetOk("port_channel_min_links"); ok {
			portChannelPayload["minLinks"] = portChannelMinLinks.(int)
		}
		if portChannelMaxLinks, ok := d.GetOk("port_channel_max_links"); ok {
			portChannelPayload["maxLinks"] = portChannelMaxLinks.(int)
		}
		if loadBalanceHashing, ok := d.GetOk("load_balance_hashing"); ok {
			hashingValue := loadBalanceHashing.(string)
			if mappedValue, exists := loadBalanceHashingMap[hashingValue]; exists {
				portChannelPayload["hashFields"] = mappedValue
			} else {
				portChannelPayload["hashFields"] = hashingValue
			}
		}
		if len(portChannelPayload) == 0 {
			err := addPatchPayloadToContainer(payloadCont, "remove", fmt.Sprintf("%s/portChannelPolicy", updatePath), nil)
			if err != nil {
				return err
			}
		} else {
			err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/portChannelPolicy", updatePath), portChannelPayload)
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("synce") {
		syncE, ok := d.GetOk("synce")
		if !ok {
			err := addPatchPayloadToContainer(payloadCont, "remove", fmt.Sprintf("%s/syncEthPolicy", updatePath), nil)
			if err != nil {
				return err
			}
		} else {
			err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/syncEthPolicy", updatePath), syncE.(string))
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("domains") {
		domains, ok := d.GetOk("domains")
		if !ok {
			err := addPatchPayloadToContainer(payloadCont, "remove", fmt.Sprintf("%s/domains", updatePath), nil)
			if err != nil {
				return err
			}
		} else {
			domainsList := domains.(*schema.Set).List()
			err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/domains", updatePath), domainsList)
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("access_macsec_policy") {
		macsecPolicy, ok := d.GetOk("access_macsec_policy")
		if !ok {
			err := addPatchPayloadToContainer(payloadCont, "remove", fmt.Sprintf("%s/accessMACsecPolicy", updatePath), nil)
			if err != nil {
				return err
			}
		} else {
			err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/accessMACsecPolicy", updatePath), macsecPolicy.(string))
			if err != nil {
				return err
			}
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/InterfaceSetting/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO Interface Setting Resource - Update Complete: %v", d.Id())
	return resourceMSOInterfaceSettingRead(d, m)
}

func resourceMSOInterfaceSettingRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Interface Setting Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return err
	}

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyName, err := GetPolicyNameFromResourceId(d.Id(), "InterfaceSetting")
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "fabricPolicyTemplate", "template", "interfacePolicyGroups")
	if err != nil {
		return err
	}

	err = setInterfaceSettingData(d, policy, templateId)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] MSO Interface Setting Resource - Read Complete: %v", d.Id())
	return nil
}

func resourceMSOInterfaceSettingDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Interface Setting Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	interfaceSettingIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "fabricPolicyTemplate", "template", "interfacePolicyGroups")
	if err != nil {
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/fabricPolicyTemplate/template/interfacePolicyGroups/%d", interfaceSettingIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO Interface Setting Resource - Delete Complete: %v", d.Id())
	return nil
}

func buildInterfaceSettingPayload(d *schema.ResourceData) map[string]interface{} {
	payload := make(map[string]interface{})

	if interfaceType, ok := d.GetOk("type"); ok {
		payload["type"] = interfaceType.(string)
	}

	if name, ok := d.GetOk("name"); ok {
		payload["name"] = name.(string)
	}

	if description, ok := d.GetOk("description"); ok {
		payload["description"] = description.(string)
	}

	// CDP Configuration
	if cdpAdminState, ok := d.GetOk("cdp_admin_state"); ok {
		payload["cdp"] = map[string]interface{}{
			"adminState": cdpAdminState.(string),
		}
	}

	// LLDP Configuration
	lldpPayload := make(map[string]interface{})
	if lldpReceive, ok := d.GetOk("lldp_receive_state"); ok {
		lldpPayload["receiveState"] = lldpReceive.(string)
	}
	if lldpTransmit, ok := d.GetOk("lldp_transmit_state"); ok {
		lldpPayload["transmitState"] = lldpTransmit.(string)
	}
	if len(lldpPayload) > 0 {
		payload["lldp"] = lldpPayload
	}

	// LLFC Configuration
	llfcPayload := make(map[string]interface{})
	if llfcReceive, ok := d.GetOk("llfc_receive_state"); ok {
		llfcPayload["receiveState"] = llfcReceive.(string)
	}
	if llfcTransmit, ok := d.GetOk("llfc_transmit_state"); ok {
		llfcPayload["transmitState"] = llfcTransmit.(string)
	}
	if len(llfcPayload) > 0 {
		payload["llfc"] = llfcPayload
	}

	// PFC Configuration
	if pfcAdminState, ok := d.GetOk("pfc_admin_state"); ok {
		payload["pfc"] = map[string]interface{}{
			"adminState": pfcAdminState.(string),
		}
	}

	// L2 Interface Configuration
	l2InterfacePayload := make(map[string]interface{})
	if l2QinQ, ok := d.GetOk("l2_interface_qinq"); ok {
		qinqValue := l2QinQ.(string)
		if mappedValue, exists := l2InterfaceQinqMap[qinqValue]; exists {
			l2InterfacePayload["qinq"] = mappedValue
		} else {
			l2InterfacePayload["qinq"] = qinqValue
		}
	}
	if l2ReflectiveRelay, ok := d.GetOk("l2_interface_reflective_relay"); ok {
		l2InterfacePayload["reflectiveRelay"] = l2ReflectiveRelay.(string)
	}
	if vlanScope, ok := d.GetOk("vlan_scope"); ok {
		l2InterfacePayload["vlanScope"] = vlanScope.(string)
	}
	if len(l2InterfacePayload) > 0 {
		payload["l2Interface"] = l2InterfacePayload
	}

	// STP Configuration
	stpPayload := make(map[string]interface{})
	if stpBpduFilter, ok := d.GetOk("stp_bpdu_filter"); ok {
		stpPayload["bpduFilterEnabled"] = stpBpduFilter.(string)
	}
	if stpBpduGuard, ok := d.GetOk("stp_bpdu_guard"); ok {
		stpPayload["bpduGuardEnabled"] = stpBpduGuard.(string)
	}
	if len(stpPayload) > 0 {
		payload["stp"] = stpPayload
	}

	// Link Level Configuration
	linkLevelPayload := make(map[string]interface{})
	if speed, ok := d.GetOk("speed"); ok {
		linkLevelPayload["speed"] = speed.(string)
	}
	if autoNeg, ok := d.GetOk("auto_negotiation"); ok {
		if autoNeg == "on_enforce" {
			autoNeg = "on-enforce"
		}
		linkLevelPayload["autoNegotiation"] = autoNeg.(string)
	}
	if bringUpDelay, ok := d.GetOk("link_level_bring_up_delay"); ok {
		linkLevelPayload["bringUpDelay"] = bringUpDelay.(int)
	}
	if debounceInterval, ok := d.GetOk("link_level_debounce_interval"); ok {
		linkLevelPayload["debounceInterval"] = debounceInterval.(int)
	}
	if fec, ok := d.GetOk("link_level_fec"); ok {
		fecValue := fec.(string)
		if mappedValue, exists := linkLevelFecMap[fecValue]; exists {
			linkLevelPayload["fec"] = mappedValue
		} else {
			linkLevelPayload["fec"] = fecValue
		}
	}
	if len(linkLevelPayload) > 0 {
		payload["linkLevel"] = linkLevelPayload
	}

	// MCP Configuration
	mcpPayload := make(map[string]interface{})
	if mcpAdminState, ok := d.GetOk("mcp_admin_state"); ok {
		mcpPayload["adminState"] = mcpAdminState.(string)
	}
	if mcpStrictMode, ok := d.GetOk("mcp_strict_mode"); ok {
		mcpPayload["mcpMode"] = mcpStrictMode.(string)
	}
	if mcpInitialDelay, ok := d.GetOk("mcp_initial_delay_time"); ok {
		mcpPayload["initialDelayTime"] = mcpInitialDelay.(int)
	}
	if mcpTxFreq, ok := d.GetOk("mcp_transmission_frequency_sec"); ok {
		mcpPayload["txFreq"] = mcpTxFreq.(int)
	}
	if mcpTxFreqMsec, ok := d.GetOk("mcp_transmission_frequency_msec"); ok {
		mcpPayload["txFreqMsec"] = mcpTxFreqMsec.(int)
	}
	if mcpGracePeriod, ok := d.GetOk("mcp_grace_period_sec"); ok {
		mcpPayload["gracePeriod"] = mcpGracePeriod.(int)
	}
	if mcpGracePeriodMsec, ok := d.GetOk("mcp_grace_period_msec"); ok {
		mcpPayload["gracePeriodMsec"] = mcpGracePeriodMsec.(int)
	}
	if len(mcpPayload) > 0 {
		payload["mcp"] = mcpPayload
	}

	// Port Channel Policy Configuration
	portChannelPayload := make(map[string]interface{})
	if portChannelMode, ok := d.GetOk("port_channel_mode"); ok {
		modeValue := portChannelMode.(string)
		if mappedValue, exists := portChannelModeMap[modeValue]; exists {
			portChannelPayload["mode"] = mappedValue
		} else {
			portChannelPayload["mode"] = modeValue
		}
	}
	if controls, ok := d.GetOk("controls"); ok {
		controlsList := controls.(*schema.Set).List()
		mappedControls := make([]string, 0)
		for _, control := range controlsList {
			controlValue := control.(string)
			if mappedValue, exists := controlMap[controlValue]; exists {
				mappedControls = append(mappedControls, mappedValue)
			} else {
				mappedControls = append(mappedControls, controlValue)
			}
		}
		portChannelPayload["control"] = mappedControls
	}
	if portChannelMinLinks, ok := d.GetOk("port_channel_min_links"); ok {
		portChannelPayload["minLinks"] = portChannelMinLinks.(int)
	}
	if portChannelMaxLinks, ok := d.GetOk("port_channel_max_links"); ok {
		portChannelPayload["maxLinks"] = portChannelMaxLinks.(int)
	}
	if loadBalanceHashing, ok := d.GetOk("load_balance_hashing"); ok {
		hashingValue := loadBalanceHashing.(string)
		if mappedValue, exists := loadBalanceHashingMap[hashingValue]; exists {
			portChannelPayload["hashFields"] = mappedValue
		} else {
			portChannelPayload["hashFields"] = hashingValue
		}
	}
	if len(portChannelPayload) > 0 || payload["type"] == "portchannel" {
		payload["portChannelPolicy"] = portChannelPayload
	}

	// SyncE Policy
	if syncE, ok := d.GetOk("synce"); ok {
		payload["syncEthPolicy"] = syncE.(string)
	}

	// Domains
	if domains, ok := d.GetOk("domains"); ok {
		domainsList := domains.(*schema.Set).List()
		if len(domainsList) > 0 {
			payload["domains"] = domainsList
		}
	}

	// Access MACsec Policy
	if macsecPolicy, ok := d.GetOk("access_macsec_policy"); ok {
		payload["accessMACsecPolicy"] = macsecPolicy.(string)
	}

	return payload
}
