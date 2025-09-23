package mso

import (
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOServiceDeviceCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOServiceDeviceClusterCreate,
		Read:   resourceMSOServiceDeviceClusterRead,
		Update: resourceMSOServiceDeviceClusterUpdate,
		Delete: resourceMSOServiceDeviceClusterDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOServiceDeviceClusterImport,
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"device_mode": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"layer1", "layer2", "layer3",
				}, false),
			},
			"device_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"firewall", "load_balancer", "other",
				}, false),
			},
			"interface_properties": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"bd_uuid": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"external_epg_uuid": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ipsla_monitoring_policy_uuid": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"qos_policy_uuid": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"preferred_group": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"rewrite_source_mac": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"anycast": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"config_static_mac": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"is_backup_redirect_ip": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"load_balance_hashing": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "sourceDestinationAndProtocol",
							ValidateFunc: validation.StringInSlice([]string{
								"sourceDestinationAndProtocol", "sourceIP", "destinationIP",
							}, false),
						},
						"pod_aware_redirection": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"resilient_hashing": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"tag_based_sorting": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"min_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(0, 100),
						},
						"max_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(0, 100),
						},
						"threshold_down_action": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"permit", "deny", "bypass",
							}, false),
						},
					},
				},
			},
		},
		CustomizeDiff: func(d *schema.ResourceDiff, m interface{}) error {
			interfaceValue, ok := d.GetOk("interface_properties")
			if !ok {
				return nil
			}
			set := interfaceValue.(*schema.Set)
			for _, raw := range set.List() {
				listItem := raw.(map[string]interface{})
				count := 0
				if bdUUID, _ := listItem["bd_uuid"].(string); bdUUID != "" {
					count++
				}
				if extEpgUUID, _ := listItem["external_epg_uuid"].(string); extEpgUUID != "" {
					count++
				}
				if count != 1 {
					return fmt.Errorf("interface_properties: exactly one of bd_uuid or external_epg_uuid must be set")
				}
			}
			return nil
		},
	}
}

func buildServiceDeviceClusterInterfacesPayload(d *schema.ResourceData) []map[string]interface{} {
	interfacesSet := d.Get("interface_properties").(*schema.Set)
	interfaces := interfacesSet.List()

	payload := make([]map[string]interface{}, len(interfaces))

	for i, val := range interfaces {
		iface := val.(map[string]interface{})

		advancedConfig := make(map[string]interface{})

		if v, ok := iface["rewrite_source_mac"]; ok {
			advancedConfig["rewriteSourceMac"] = v.(bool)
		}
		if v, ok := iface["anycast"]; ok {
			advancedConfig["anycast"] = v.(bool)
		}
		if v, ok := iface["config_static_mac"]; ok {
			advancedConfig["configStaticMac"] = v.(bool)
		}
		if v, ok := iface["is_backup_redirect_ip"]; ok {
			advancedConfig["isBackupRedirectIP"] = v.(bool)
		}
		if v, ok := iface["pod_aware_redirection"]; ok {
			advancedConfig["podAwareRedirection"] = v.(bool)
		}
		if v, ok := iface["preferred_group"]; ok {
			advancedConfig["preferredGroup"] = v.(bool)
		}
		if v, ok := iface["resilient_hashing"]; ok {
			advancedConfig["resilientHashing"] = v.(bool)
		}
		if v, ok := iface["tag_based_sorting"]; ok {
			advancedConfig["tag"] = v.(bool)
		}

		if v, ok := iface["load_balance_hashing"].(string); ok && v != "" {
			advancedConfig["loadBalanceHashing"] = v
		}

		if v, ok := iface["qos_policy_uuid"].(string); ok && v != "" {
			advancedConfig["qosPolicyRef"] = v
		}

		thresholdConfig := make(map[string]interface{})
		if v := iface["max_threshold"].(int); v > 0 {
			thresholdConfig["maxThreshold"] = v
		}
		if v := iface["min_threshold"].(int); v > 0 {
			thresholdConfig["minThreshold"] = v
		}
		if v, ok := iface["threshold_down_action"].(string); ok && v != "" {
			thresholdConfig["thresholdDownAction"] = v
			advancedConfig["thresholdForRedirectDestination"] = true
		}

		if len(thresholdConfig) > 0 {
			advancedConfig["thresholdForRedirect"] = thresholdConfig
		}

		interfacePayload := map[string]interface{}{
			"name":                 iface["name"].(string),
			"isAdvancedIntfConfig": true,
			"redirect":             true,
		}

		if len(advancedConfig) > 0 {
			interfacePayload["advancedIntfConfig"] = advancedConfig
			interfacePayload["isAdvancedIntfConfig"] = true
		}

		if v, ok := iface["ipsla_monitoring_policy_uuid"].(string); ok && v != "" {
			interfacePayload["ipslaMonitoringRef"] = v
			advancedConfig["advancedTrackingOptions"] = true
		}

		if v, ok := iface["bd_uuid"].(string); ok && v != "" {
			interfacePayload["deviceInterfaceType"] = "bd"
			interfacePayload["bdRef"] = v
		} else if v, ok := iface["external_epg_uuid"].(string); ok && v != "" {
			interfacePayload["deviceInterfaceType"] = "l3out"
			interfacePayload["externalEpgRef"] = v
		}

		payload[i] = interfacePayload
	}
	return payload
}

func buildServiceDeviceClusterPayload(d *schema.ResourceData) map[string]interface{} {
	payload := map[string]interface{}{}
	payload["name"] = d.Get("name").(string)
	payload["description"] = d.Get("description").(string)
	payload["deviceMode"] = d.Get("device_mode").(string)
	payload["deviceLocation"] = "onPremise"

	if d.Get("device_type").(string) == "load_balancer" {
		payload["deviceType"] = "loadBalancer"
	} else {
		payload["deviceType"] = d.Get("device_type").(string)
	}

	interfaces := buildServiceDeviceClusterInterfacesPayload(d)
	payload["interfaces"] = interfaces

	numInterfaces := len(interfaces)
	if numInterfaces == 1 {
		payload["connectivityMode"] = "oneArm"
	} else if numInterfaces == 2 {
		payload["connectivityMode"] = "twoArm"
	} else {
		payload["connectivityMode"] = "advanced"
	}

	return payload
}

func setServiceDeviceClusterData(d *schema.ResourceData, response *container.Container, templateId string) error {
	d.SetId(fmt.Sprintf("templateId/%s/ServiceDeviceCluster/%s", templateId, models.StripQuotes(response.S("name").String())))
	d.Set("template_id", templateId)
	d.Set("name", models.StripQuotes(response.S("name").String()))
	d.Set("uuid", models.StripQuotes(response.S("uuid").String()))

	if response.Exists("description") {
		d.Set("description", models.StripQuotes(response.S("description").String()))
	}

	d.Set("device_mode", models.StripQuotes(response.S("deviceMode").String()))

	deviceType := models.StripQuotes(response.S("deviceType").String())
	if deviceType == "loadBalancer" {
		d.Set("device_type", "load_balancer")
	} else {
		d.Set("device_type", deviceType)
	}

	interfaces, err := response.S("interfaces").Children()
	if err != nil {
		return nil
	}

	var interfaceProperties []map[string]interface{}
	for _, iface := range interfaces {
		prop := make(map[string]interface{})
		prop["name"] = models.StripQuotes(iface.S("name").String())

		if iface.Exists("bdRef") {
			prop["bd_uuid"] = models.StripQuotes(iface.S("bdRef").String())
		}
		if iface.Exists("externalEpgRef") {
			prop["external_epg_uuid"] = models.StripQuotes(iface.S("externalEpgRef").String())
		}
		if iface.Exists("ipslaMonitoringRef") {
			prop["ipsla_monitoring_policy_uuid"] = models.StripQuotes(iface.S("ipslaMonitoringRef").String())
		}

		if iface.Exists("advancedIntfConfig") {
			advancedConfig := iface.S("advancedIntfConfig")

			if advancedConfig.Exists("qosPolicyRef") {
				qosPolicyRef := models.StripQuotes(advancedConfig.S("qosPolicyRef").String())
				if qosPolicyRef == "{}" {
					prop["qos_policy_uuid"] = ""
				} else {
					prop["qos_policy_uuid"] = qosPolicyRef
				}
			}

			if advancedConfig.Exists("preferredGroup") {
				prop["preferred_group"] = advancedConfig.S("preferredGroup").Data().(bool)
			}
			if advancedConfig.Exists("rewriteSourceMac") {
				prop["rewrite_source_mac"] = advancedConfig.S("rewriteSourceMac").Data().(bool)
			}
			if advancedConfig.Exists("anycast") {
				prop["anycast"] = advancedConfig.S("anycast").Data().(bool)
			}
			if advancedConfig.Exists("configStaticMac") {
				prop["config_static_mac"] = advancedConfig.S("configStaticMac").Data().(bool)
			}
			if advancedConfig.Exists("isBackupRedirectIP") {
				prop["is_backup_redirect_ip"] = advancedConfig.S("isBackupRedirectIP").Data().(bool)
			}
			if advancedConfig.Exists("podAwareRedirection") {
				prop["pod_aware_redirection"] = advancedConfig.S("podAwareRedirection").Data().(bool)
			}
			if advancedConfig.Exists("resilientHashing") {
				prop["resilient_hashing"] = advancedConfig.S("resilientHashing").Data().(bool)
			}
			if advancedConfig.Exists("tag") {
				prop["tag_based_sorting"] = advancedConfig.S("tag").Data().(bool)
			}

			if advancedConfig.Exists("loadBalanceHashing") {
				prop["load_balance_hashing"] = models.StripQuotes(advancedConfig.S("loadBalanceHashing").String())
			}

			if advancedConfig.Exists("thresholdForRedirect") {
				thresholdConfig := advancedConfig.S("thresholdForRedirect")
				if thresholdConfig.Exists("minThreshold") {
					prop["min_threshold"] = int(thresholdConfig.S("minThreshold").Data().(float64))
				}
				if thresholdConfig.Exists("maxThreshold") {
					prop["max_threshold"] = int(thresholdConfig.S("maxThreshold").Data().(float64))
				}
				if thresholdConfig.Exists("thresholdDownAction") {
					action := models.StripQuotes(thresholdConfig.S("thresholdDownAction").String())
					if action == "{}" {
						prop["threshold_down_action"] = ""
					} else {
						prop["threshold_down_action"] = action
					}
				}
			}
		}
		interfaceProperties = append(interfaceProperties, prop)
	}
	d.Set("interface_properties", interfaceProperties)

	return nil
}

func resourceMSOServiceDeviceClusterImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO Service Device Cluster Resource - Beginning Import: %v", d.Id())
	resourceMSOServiceDeviceClusterRead(d, m)
	log.Printf("[DEBUG] MSO Service Device Cluster Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOServiceDeviceClusterCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Service Device Cluster Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	payload := buildServiceDeviceClusterPayload(d)
	payloadModel := models.GetPatchPayload("add", "/deviceTemplate/template/devices/-", payload)
	templateId := d.Get("template_id").(string)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/ServiceDeviceCluster/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO Service Device Cluster Resource - Create Complete: %v", d.Id())
	return resourceMSOServiceDeviceClusterRead(d, m)
}

func resourceMSOServiceDeviceClusterRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Service Device Cluster Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	var templateId, policyName string
	var err error

	if strings.Contains(d.Id(), "ServiceDeviceCluster") {
		templateId, err = GetTemplateIdFromResourceId(d.Id())
		if err != nil {
			return err
		}
		policyName, err = GetPolicyNameFromResourceId(d.Id(), "ServiceDeviceCluster")
		if err != nil {
			return err
		}
	} else {
		templateId = d.Get("template_id").(string)
		policyName = d.Get("name").(string)
	}

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		log.Printf("[DEBUG] MSO Service Device Cluster not found, removing from state: %v", d.Id())
		d.SetId("")
		return nil
	}

	policy, err := GetPolicyByName(response, policyName, "deviceTemplate", "template", "devices")
	if err != nil {
		log.Printf("[DEBUG] MSO Service Device Cluster not found, removing from state: %v", d.Id())
		d.SetId("")
		return nil
	}

	setServiceDeviceClusterData(d, policy, templateId)
	log.Printf("[DEBUG] MSO Service Device Cluster Resource - Read Complete: %v", d.Id())
	return nil
}

func resourceMSOServiceDeviceClusterUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Service Device Cluster Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "deviceTemplate", "template", "devices")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/deviceTemplate/template/devices/%d", policyIndex)
	payloadCont := container.New()
	payloadCont.Array()

	if d.HasChange("name") {
		addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/name", updatePath), d.Get("name").(string))
	}
	if d.HasChange("description") {
		addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/description", updatePath), d.Get("description").(string))
	}
	if d.HasChange("device_mode") {
		addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/deviceMode", updatePath), d.Get("device_mode").(string))
	}
	if d.HasChange("device_type") {
		var deviceType string
		if d.Get("device_type").(string) == "load_balancer" {
			deviceType = "loadBalancer"
		} else {
			deviceType = d.Get("device_type").(string)
		}
		addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/deviceType", updatePath), deviceType)
	}

	if d.HasChange("interface_properties") {
		log.Printf(" HERE Detected change in interface_properties. Replacing the entire block.")
		interfaces := buildServiceDeviceClusterInterfacesPayload(d)
		addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/interfaces", updatePath), interfaces)

		var connectivityMode string
		numInterfaces := len(interfaces)
		if numInterfaces == 1 {
			connectivityMode = "oneArm"
		} else if numInterfaces == 2 {
			connectivityMode = "twoArm"
		} else {
			connectivityMode = "advanced"
		}
		addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/connectivityMode", updatePath), connectivityMode)
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/ServiceDeviceCluster/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO Service Device Cluster Resource - Update Complete: %v", d.Id())
	return resourceMSOServiceDeviceClusterRead(d, m)
}

func resourceMSOServiceDeviceClusterDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Service Device Cluster Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		log.Printf("[DEBUG] Template not found during delete for resource: %v", d.Id())
		d.SetId("")
		return nil
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "deviceTemplate", "template", "devices")
	if err != nil {
		log.Printf("[DEBUG] Service Device Cluster not found in template during delete for resource: %v", d.Id())
		d.SetId("")
		return nil
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/deviceTemplate/template/devices/%d", policyIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", d.Get("template_id").(string)), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO Service Device Cluster Resource - Delete Complete: %v", d.Id())
	return nil
}
