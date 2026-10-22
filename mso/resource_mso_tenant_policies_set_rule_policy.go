package mso

import (
	"context"
	"fmt"
	"log"
	"sort"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceMSOSetRulePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSetRulePolicyCreate,
		Read:   resourceMSOSetRulePolicyRead,
		Update: resourceMSOSetRulePolicyUpdate,
		Delete: resourceMSOSetRulePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOSetRulePolicyImport,
		},
		CustomizeDiff: validateSetRulePolicy,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the tenant policy template.",
				ForceNew:    true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
				Description:  "The name of the Set Rule Policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the Set Rule Policy.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the Set Rule Policy.",
			},
			"set_route_tag": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The route tag action.",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 999999999),
							Description:  "The action value. Range: 0-999999999.",
						},
					},
				},
			},
			"set_weight": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The weight action.",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 9999),
							Description:  "The action value. Range: 0-9999.",
						},
					},
				},
			},
			"set_preference": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The preference action.",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 999999999),
							Description:  "The action value. Range: 0-999999999.",
						},
					},
				},
			},
			"set_metric": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 999999999),
				Description:  "The metric action. Range: 1-999999999.",
			},
			"set_next_hop": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
				Description:  "The next-hop IPv4 or IPv6 address.",
			},
			"set_metric_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"none", "ospf-type1", "ospf-type2"}, false),
				Description:  "The metric type. Valid values: none, ospf-type1, ospf-type2.",
			},
			"set_next_hop_propagate": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable next-hop propagation.",
			},
			"set_multi_path": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable multipath. Requires next-hop propagation.",
			},
			"set_community": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The community action.",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"community": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The community value. Required for append and replace.",
						},
						"criteria": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"none", "append", "replace"}, false),
							Description:  "The community criteria. Valid values: none, append, replace.",
						},
					},
				},
			},
			"set_additional_community": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The community values to append.",
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.StringIsNotEmpty},
			},
			"set_dampening": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The route dampening action.",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"half_life": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 60),
							Description:  "The half-life time in minutes. Range: 1-60.",
						},
						"reuse_limit": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 20000),
							Description:  "The reuse limit. Range: 1-20000.",
						},
						"suppress_limit": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 20000),
							Description:  "The suppress limit. Must be larger than reuse_limit. Range: 1-20000.",
						},
						"max_suppress_time": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 255),
							Description:  "The maximum suppress time in minutes. Range: 1-255.",
						},
					},
				},
			},
			"set_as_path": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "The AS-path prepend actions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prepend": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "The action to prepend the configured AS numbers.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path_asn": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "The AS numbers to prepend. Order values determine the sequence and must be unique.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"asn": {
													Type:         schema.TypeInt,
													Required:     true,
													ValidateFunc: validateUint32Range(1, 4294967295),
													Description:  "The AS number. Range: 1-4294967295.",
												},
												"order": {
													Type:         schema.TypeInt,
													Required:     true,
													ValidateFunc: validation.IntBetween(0, 31),
													Description:  "The AS-path order. Range: 0-31.",
												},
											},
										},
									},
								},
							},
						},
						"prepend_last_as": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "The action to prepend the last AS number.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"asn_count": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 10),
										Description:  "The number of times to prepend the last AS. Range: 1-10.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func validateSetRulePolicy(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	if d.NewValueKnown("set_next_hop_propagate") && d.NewValueKnown("set_route_tag") && d.Get("set_next_hop_propagate").(bool) && len(d.Get("set_route_tag").([]interface{})) > 0 {
		return fmt.Errorf("set_next_hop_propagate cannot be enabled when set_route_tag is configured")
	}
	if d.NewValueKnown("set_multi_path") && d.NewValueKnown("set_next_hop_propagate") && d.Get("set_multi_path").(bool) && !d.Get("set_next_hop_propagate").(bool) {
		return fmt.Errorf("set_multi_path requires set_next_hop_propagate to be enabled")
	}
	// Defer cross-field checks while the configured values are unknown.
	rawConfig := d.GetRawConfig()
	if !rawConfig.IsKnown() || rawConfig.IsNull() {
		return nil
	}
	if rawConfig.GetAttr("set_dampening").IsWhollyKnown() {
		blocks := d.Get("set_dampening").([]interface{})
		if len(blocks) > 0 {
			settings, ok := blocks[0].(map[string]interface{})
			if ok && settings["suppress_limit"].(int) <= settings["reuse_limit"].(int) {
				return fmt.Errorf("set_dampening suppress_limit must be larger than reuse_limit")
			}
		}
	}
	if rawConfig.GetAttr("set_community").IsWhollyKnown() {
		blocks := d.Get("set_community").([]interface{})
		if len(blocks) > 0 {
			settings, ok := blocks[0].(map[string]interface{})
			if ok && settings["criteria"].(string) != "none" && settings["community"].(string) == "" {
				return fmt.Errorf("set_community community is required for append and replace criteria")
			}
		}
	}
	if rawConfig.GetAttr("set_as_path").IsWhollyKnown() {
		blocks := d.Get("set_as_path").([]interface{})
		if len(blocks) > 0 {
			settings, ok := blocks[0].(map[string]interface{})
			if !ok || (len(settings["prepend"].([]interface{})) == 0 && len(settings["prepend_last_as"].([]interface{})) == 0) {
				return fmt.Errorf("set_as_path requires at least one of prepend or prepend_last_as")
			}
			if prepend := settings["prepend"].([]interface{}); len(prepend) > 0 {
				if prependSettings, ok := prepend[0].(map[string]interface{}); ok {
					orders := map[int]bool{}
					for _, item := range prependSettings["path_asn"].(*schema.Set).List() {
						order := item.(map[string]interface{})["order"].(int)
						if orders[order] {
							return fmt.Errorf("set_as_path prepend path_asn order %d is duplicated", order)
						}
						orders[order] = true
					}
				}
			}
		}
	}
	return nil
}

func buildSetRulePolicyPayload(d *schema.ResourceData) map[string]interface{} {
	payload := map[string]interface{}{"name": d.Get("name").(string), "description": d.Get("description").(string)}
	if routeTag := d.Get("set_route_tag").([]interface{}); len(routeTag) > 0 {
		payload["setRouteTag"] = routeTag[0].(map[string]interface{})["value"].(int)
	}
	if weight := d.Get("set_weight").([]interface{}); len(weight) > 0 {
		payload["setWeight"] = weight[0].(map[string]interface{})["value"].(int)
	}
	if preference := d.Get("set_preference").([]interface{}); len(preference) > 0 {
		payload["setPreference"] = preference[0].(map[string]interface{})["value"].(int)
	}
	if value, ok := d.GetOk("set_metric"); ok {
		payload["setMetric"] = value.(int)
	}
	if nextHop := d.Get("set_next_hop").(string); nextHop != "" {
		payload["setNextHop"] = nextHop
	}
	if metricType := d.Get("set_metric_type").(string); metricType != "" {
		payload["setMetricType"] = metricType
	}
	if d.Get("set_next_hop_propagate").(bool) {
		payload["setNextHopPropagate"] = true
	}
	if d.Get("set_multi_path").(bool) {
		payload["setMultiPath"] = true
	}
	if blocks := d.Get("set_community").([]interface{}); len(blocks) > 0 {
		settings := blocks[0].(map[string]interface{})
		community := map[string]interface{}{"criteria": settings["criteria"]}
		if settings["community"].(string) != "" {
			community["community"] = settings["community"]
		}
		payload["setCommunity"] = community
	}
	if communities := d.Get("set_additional_community").(*schema.Set).List(); len(communities) > 0 {
		actions := make([]interface{}, 0, len(communities))
		for _, community := range communities {
			actions = append(actions, map[string]interface{}{"community": community, "criteria": "append"})
		}
		payload["setAddCommunities"] = actions
	}
	if blocks := d.Get("set_dampening").([]interface{}); len(blocks) > 0 {
		settings := blocks[0].(map[string]interface{})
		payload["setDampening"] = map[string]interface{}{
			"halfLife":        settings["half_life"],
			"reuseLimit":      settings["reuse_limit"],
			"suppressLimit":   settings["suppress_limit"],
			"maxSuppressTime": settings["max_suppress_time"],
		}
	}
	if blocks := d.Get("set_as_path").([]interface{}); len(blocks) > 0 {
		settings := blocks[0].(map[string]interface{})
		actions := make([]interface{}, 0, 2)
		if prepend := settings["prepend"].([]interface{}); len(prepend) > 0 {
			action := map[string]interface{}{"criteria": "prepend"}
			// An empty prepend block still configures the action and may have nil contents.
			if prependSettings, ok := prepend[0].(map[string]interface{}); ok {
				paths := prependSettings["path_asn"].(*schema.Set).List()
				sort.Slice(paths, func(i, j int) bool {
					return paths[i].(map[string]interface{})["order"].(int) < paths[j].(map[string]interface{})["order"].(int)
				})
				if len(paths) > 0 {
					action["pathASNs"] = paths
				}
			}
			actions = append(actions, action)
		}
		if lastAS := settings["prepend_last_as"].([]interface{}); len(lastAS) > 0 {
			actions = append(actions, map[string]interface{}{
				"criteria": "prepend-last-as",
				"asnCount": lastAS[0].(map[string]interface{})["asn_count"],
			})
		}
		payload["setAsPath"] = actions
	}
	return payload
}

func setSetRulePolicyData(d *schema.ResourceData, policy *container.Container, templateID string) error {
	name, ok := policy.S("name").Data().(string)
	if !ok || name == "" {
		return fmt.Errorf("Set Rule Policy name is missing or is not a string")
	}
	uuid, ok := policy.S("uuid").Data().(string)
	if !ok || uuid == "" {
		return fmt.Errorf("Set Rule Policy UUID is missing or is not a string")
	}
	description, _ := policy.S("description").Data().(string)
	routeTag := []interface{}{}
	if value := policy.S("setRouteTag").Data(); value != nil {
		number, ok := value.(float64)
		if !ok {
			return fmt.Errorf("Set Rule Policy setRouteTag is not a number")
		}
		routeTag = append(routeTag, map[string]interface{}{"value": int(number)})
	}
	weight := []interface{}{}
	if value := policy.S("setWeight").Data(); value != nil {
		number, ok := value.(float64)
		if !ok {
			return fmt.Errorf("Set Rule Policy setWeight is not a number")
		}
		weight = append(weight, map[string]interface{}{"value": int(number)})
	}
	preference := []interface{}{}
	if value := policy.S("setPreference").Data(); value != nil {
		number, ok := value.(float64)
		if !ok {
			return fmt.Errorf("Set Rule Policy setPreference is not a number")
		}
		preference = append(preference, map[string]interface{}{"value": int(number)})
	}
	metric := 0
	if value := policy.S("setMetric").Data(); value != nil {
		number, ok := value.(float64)
		if !ok {
			return fmt.Errorf("Set Rule Policy setMetric is not a number")
		}
		metric = int(number)
	}
	nextHop, _ := policy.S("setNextHop").Data().(string)
	metricType, _ := policy.S("setMetricType").Data().(string)
	nextHopPropagate, _ := policy.S("setNextHopPropagate").Data().(bool)
	multiPath, _ := policy.S("setMultiPath").Data().(bool)
	community := []interface{}{}
	if value := policy.S("setCommunity").Data(); value != nil {
		settings, ok := value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("Set Rule Policy setCommunity is not an object")
		}
		if len(settings) > 0 {
			community = append(community, settings)
		}
	}
	communities := []interface{}{}
	if value := policy.S("setAddCommunities").Data(); value != nil {
		actions, ok := value.([]interface{})
		if !ok {
			return fmt.Errorf("Set Rule Policy setAddCommunities is not a list")
		}
		for _, action := range actions {
			entry, ok := action.(map[string]interface{})
			if !ok {
				return fmt.Errorf("Set Rule Policy setAddCommunities entry is not an object")
			}
			community, ok := entry["community"].(string)
			if !ok {
				return fmt.Errorf("Set Rule Policy setAddCommunities community is not a string")
			}
			communities = append(communities, community)
		}
	}
	dampening := []interface{}{}
	if value := policy.S("setDampening").Data(); value != nil {
		settings, ok := value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("Set Rule Policy setDampening is not an object")
		}
		if len(settings) > 0 {
			halfLife, ok := settings["halfLife"].(float64)
			if !ok {
				return fmt.Errorf("Set Rule Policy setDampening.halfLife is not a number")
			}
			reuseLimit, ok := settings["reuseLimit"].(float64)
			if !ok {
				return fmt.Errorf("Set Rule Policy setDampening.reuseLimit is not a number")
			}
			suppressLimit, ok := settings["suppressLimit"].(float64)
			if !ok {
				return fmt.Errorf("Set Rule Policy setDampening.suppressLimit is not a number")
			}
			maxSuppressTime, ok := settings["maxSuppressTime"].(float64)
			if !ok {
				return fmt.Errorf("Set Rule Policy setDampening.maxSuppressTime is not a number")
			}
			dampening = append(dampening, map[string]interface{}{
				"half_life":         int(halfLife),
				"reuse_limit":       int(reuseLimit),
				"suppress_limit":    int(suppressLimit),
				"max_suppress_time": int(maxSuppressTime),
			})
		}
	}
	asPath := []interface{}{}
	if value := policy.S("setAsPath").Data(); value != nil {
		actions, ok := value.([]interface{})
		if !ok {
			return fmt.Errorf("Set Rule Policy setAsPath is not a list")
		}
		settings := map[string]interface{}{"prepend": []interface{}{}, "prepend_last_as": []interface{}{}}
		for _, action := range actions {
			entry, ok := action.(map[string]interface{})
			if !ok {
				return fmt.Errorf("Set Rule Policy setAsPath entry is not an object")
			}
			switch entry["criteria"] {
			case "prepend":
				paths := []interface{}{}
				if pathValue := entry["pathASNs"]; pathValue != nil {
					paths, ok = pathValue.([]interface{})
					if !ok {
						return fmt.Errorf("Set Rule Policy setAsPath.pathASNs is not a list")
					}
				}
				settings["prepend"] = []interface{}{map[string]interface{}{"path_asn": paths}}
			case "prepend-last-as":
				count, ok := entry["asnCount"].(float64)
				if !ok {
					return fmt.Errorf("Set Rule Policy setAsPath.asnCount is not a number")
				}
				settings["prepend_last_as"] = []interface{}{map[string]interface{}{"asn_count": int(count)}}
			default:
				return fmt.Errorf("Set Rule Policy has unsupported AS-path criteria %v", entry["criteria"])
			}
		}
		if len(actions) > 0 {
			asPath = append(asPath, settings)
		}
	}
	if err := d.Set("template_id", templateID); err != nil {
		return fmt.Errorf("set Set Rule Policy template_id: %w", err)
	}
	if err := d.Set("name", name); err != nil {
		return fmt.Errorf("set Set Rule Policy name: %w", err)
	}
	if err := d.Set("description", description); err != nil {
		return fmt.Errorf("set Set Rule Policy description: %w", err)
	}
	if err := d.Set("uuid", uuid); err != nil {
		return fmt.Errorf("set Set Rule Policy uuid: %w", err)
	}
	if err := d.Set("set_route_tag", routeTag); err != nil {
		return fmt.Errorf("set Set Rule Policy set_route_tag: %w", err)
	}
	if err := d.Set("set_weight", weight); err != nil {
		return fmt.Errorf("set Set Rule Policy set_weight: %w", err)
	}
	if err := d.Set("set_preference", preference); err != nil {
		return fmt.Errorf("set Set Rule Policy set_preference: %w", err)
	}
	if err := d.Set("set_metric", metric); err != nil {
		return fmt.Errorf("set Set Rule Policy set_metric: %w", err)
	}
	if err := d.Set("set_next_hop", nextHop); err != nil {
		return fmt.Errorf("set Set Rule Policy set_next_hop: %w", err)
	}
	if err := d.Set("set_metric_type", metricType); err != nil {
		return fmt.Errorf("set Set Rule Policy set_metric_type: %w", err)
	}
	if err := d.Set("set_next_hop_propagate", nextHopPropagate); err != nil {
		return fmt.Errorf("set Set Rule Policy set_next_hop_propagate: %w", err)
	}
	if err := d.Set("set_multi_path", multiPath); err != nil {
		return fmt.Errorf("set Set Rule Policy set_multi_path: %w", err)
	}
	if err := d.Set("set_community", community); err != nil {
		return fmt.Errorf("set Set Rule Policy set_community: %w", err)
	}
	if err := d.Set("set_additional_community", communities); err != nil {
		return fmt.Errorf("set Set Rule Policy set_additional_community: %w", err)
	}
	if err := d.Set("set_dampening", dampening); err != nil {
		return fmt.Errorf("set Set Rule Policy set_dampening: %w", err)
	}
	if err := d.Set("set_as_path", asPath); err != nil {
		return fmt.Errorf("set Set Rule Policy set_as_path: %w", err)
	}
	d.SetId(fmt.Sprintf("templateId/%s/SetRulePolicy/%s", templateID, name))
	return nil
}

func resourceMSOSetRulePolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	if err := resourceMSOSetRulePolicyRead(d, m); err != nil {
		return nil, err
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("Set Rule Policy %s not found", importID)
	}
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSetRulePolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Set Rule Policy Resource - Beginning Create")
	msoClient := m.(*client.Client)
	templateID := d.Get("template_id").(string)
	payload := models.GetPatchPayload("add", "/tenantPolicyTemplate/template/setRulePolicies/-", buildSetRulePolicyPayload(d))
	if _, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateID), payload); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("templateId/%s/SetRulePolicy/%s", templateID, d.Get("name").(string)))
	return resourceMSOSetRulePolicyRead(d, m)
}

func resourceMSOSetRulePolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Set Rule Policy Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)
	templateID, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return err
	}
	name, err := GetPolicyNameFromResourceId(d.Id(), "SetRulePolicy")
	if err != nil {
		return err
	}
	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
	if err != nil {
		return err
	}
	policy, err := GetPolicyByName(response, name, "tenantPolicyTemplate", "template", "setRulePolicies")
	if err != nil {
		if isPolicyNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	return setSetRulePolicyData(d, policy, templateID)
}

func resourceMSOSetRulePolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Set Rule Policy Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateID := d.Get("template_id").(string)
	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
	if err != nil {
		return err
	}
	index, err := GetPolicyIndexByKeyAndValue(response, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "setRulePolicies")
	if err != nil {
		return err
	}
	policy := response.S("tenantPolicyTemplate", "template", "setRulePolicies").Index(index)
	path := fmt.Sprintf("/tenantPolicyTemplate/template/setRulePolicies/%d", index)
	desired := buildSetRulePolicyPayload(d)
	payload := container.New()
	payload.Array()
	if d.HasChange("name") {
		if err := addPatchPayloadToContainer(payload, "add", fmt.Sprintf("%s/name", path), desired["name"]); err != nil {
			return err
		}
	}
	if d.HasChange("description") {
		if err := addPatchPayloadToContainer(payload, "add", fmt.Sprintf("%s/description", path), desired["description"]); err != nil {
			return err
		}
	}
	if d.HasChange("set_route_tag") {
		attributePath := fmt.Sprintf("%s/setRouteTag", path)
		if value, ok := desired["setRouteTag"]; ok {
			if err := addPatchPayloadToContainer(payload, "add", attributePath, value); err != nil {
				return err
			}
		} else if policy.S("setRouteTag").Data() != nil {
			if err := addPatchPayloadToContainer(payload, "remove", attributePath, nil); err != nil {
				return err
			}
		}
	}
	if d.HasChange("set_weight") {
		attributePath := fmt.Sprintf("%s/setWeight", path)
		if value, ok := desired["setWeight"]; ok {
			if err := addPatchPayloadToContainer(payload, "add", attributePath, value); err != nil {
				return err
			}
		} else if policy.S("setWeight").Data() != nil {
			if err := addPatchPayloadToContainer(payload, "remove", attributePath, nil); err != nil {
				return err
			}
		}
	}
	if d.HasChange("set_preference") {
		attributePath := fmt.Sprintf("%s/setPreference", path)
		if value, ok := desired["setPreference"]; ok {
			if err := addPatchPayloadToContainer(payload, "add", attributePath, value); err != nil {
				return err
			}
		} else if policy.S("setPreference").Data() != nil {
			if err := addPatchPayloadToContainer(payload, "remove", attributePath, nil); err != nil {
				return err
			}
		}
	}
	if d.HasChange("set_metric") {
		attributePath := fmt.Sprintf("%s/setMetric", path)
		if value, ok := desired["setMetric"]; ok {
			if err := addPatchPayloadToContainer(payload, "add", attributePath, value); err != nil {
				return err
			}
		} else if policy.S("setMetric").Data() != nil {
			if err := addPatchPayloadToContainer(payload, "remove", attributePath, nil); err != nil {
				return err
			}
		}
	}
	if d.HasChange("set_next_hop") {
		attributePath := fmt.Sprintf("%s/setNextHop", path)
		if value, ok := desired["setNextHop"]; ok {
			if err := addPatchPayloadToContainer(payload, "add", attributePath, value); err != nil {
				return err
			}
		} else if policy.S("setNextHop").Data() != nil {
			if err := addPatchPayloadToContainer(payload, "remove", attributePath, nil); err != nil {
				return err
			}
		}
	}
	if d.HasChange("set_metric_type") {
		attributePath := fmt.Sprintf("%s/setMetricType", path)
		if value, ok := desired["setMetricType"]; ok {
			if err := addPatchPayloadToContainer(payload, "add", attributePath, value); err != nil {
				return err
			}
		} else if policy.S("setMetricType").Data() != nil {
			if err := addPatchPayloadToContainer(payload, "remove", attributePath, nil); err != nil {
				return err
			}
		}
	}
	if d.HasChange("set_next_hop_propagate") {
		attributePath := fmt.Sprintf("%s/setNextHopPropagate", path)
		if value, ok := desired["setNextHopPropagate"]; ok {
			if err := addPatchPayloadToContainer(payload, "add", attributePath, value); err != nil {
				return err
			}
		} else if policy.S("setNextHopPropagate").Data() != nil {
			if err := addPatchPayloadToContainer(payload, "remove", attributePath, nil); err != nil {
				return err
			}
		}
	}
	if d.HasChange("set_multi_path") {
		attributePath := fmt.Sprintf("%s/setMultiPath", path)
		if value, ok := desired["setMultiPath"]; ok {
			if err := addPatchPayloadToContainer(payload, "add", attributePath, value); err != nil {
				return err
			}
		} else if policy.S("setMultiPath").Data() != nil {
			if err := addPatchPayloadToContainer(payload, "remove", attributePath, nil); err != nil {
				return err
			}
		}
	}
	if d.HasChange("set_community") {
		attributePath := fmt.Sprintf("%s/setCommunity", path)
		if value, ok := desired["setCommunity"]; ok {
			if err := addPatchPayloadToContainer(payload, "add", attributePath, value); err != nil {
				return err
			}
		} else if policy.S("setCommunity").Data() != nil {
			if err := addPatchPayloadToContainer(payload, "remove", attributePath, nil); err != nil {
				return err
			}
		}
	}
	if d.HasChange("set_additional_community") {
		attributePath := fmt.Sprintf("%s/setAddCommunities", path)
		if value, ok := desired["setAddCommunities"]; ok {
			if err := addPatchPayloadToContainer(payload, "add", attributePath, value); err != nil {
				return err
			}
		} else if policy.S("setAddCommunities").Data() != nil {
			if err := addPatchPayloadToContainer(payload, "remove", attributePath, nil); err != nil {
				return err
			}
		}
	}
	if d.HasChange("set_dampening") {
		attributePath := fmt.Sprintf("%s/setDampening", path)
		if value, ok := desired["setDampening"]; ok {
			if err := addPatchPayloadToContainer(payload, "add", attributePath, value); err != nil {
				return err
			}
		} else if policy.S("setDampening").Data() != nil {
			if err := addPatchPayloadToContainer(payload, "remove", attributePath, nil); err != nil {
				return err
			}
		}
	}
	if d.HasChange("set_as_path") {
		attributePath := fmt.Sprintf("%s/setAsPath", path)
		if value, ok := desired["setAsPath"]; ok {
			if err := addPatchPayloadToContainer(payload, "add", attributePath, value); err != nil {
				return err
			}
		} else if policy.S("setAsPath").Data() != nil {
			if err := addPatchPayloadToContainer(payload, "remove", attributePath, nil); err != nil {
				return err
			}
		}
	}
	if err := doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateID), payload); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("templateId/%s/SetRulePolicy/%s", templateID, d.Get("name").(string)))
	return resourceMSOSetRulePolicyRead(d, m)
}

func resourceMSOSetRulePolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Set Rule Policy Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)
	templateID := d.Get("template_id").(string)
	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
	if err != nil {
		return err
	}
	index, err := GetPolicyIndexByKeyAndValue(response, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "setRulePolicies")
	if err != nil {
		if isPolicyNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	payload := models.GetRemovePatchPayload(fmt.Sprintf("/tenantPolicyTemplate/template/setRulePolicies/%d", index))
	if _, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateID), payload); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
