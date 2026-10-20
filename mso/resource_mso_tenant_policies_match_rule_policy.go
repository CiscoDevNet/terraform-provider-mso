package mso

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceMSOMatchRulePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOMatchRulePolicyCreate,
		Read:   resourceMSOMatchRulePolicyRead,
		Update: resourceMSOMatchRulePolicyUpdate,
		Delete: resourceMSOMatchRulePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOMatchRulePolicyImport,
		},
		CustomizeDiff: validateMatchRulePolicy,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the tenant policy template.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
				Description:  "The name of the Match Rule Policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the Match Rule Policy.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the Match Rule Policy.",
			},
			"match_prefix": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The prefixes to match.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsCIDR,
							Description:  "The IPv4 or IPv6 prefix to match.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The description of the match prefix.",
						},
						"aggregate": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to match a range of prefix lengths. Default: false.",
						},
						"from_prefix": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(0, 128),
							Description:  "The minimum prefix length. Default: 0.",
						},
						"to_prefix": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(0, 128),
							Description:  "The maximum prefix length. Default: 0.",
						},
					},
				},
			},
			"match_community_term": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The community terms to match.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 64),
							Description:  "The name of the community term.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The description of the community term.",
						},
						"match_community_factor": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "The community factors within the term.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"community": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
										Description:  "The regular, extended, or well-known community to match.",
									},
									"scope": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "transitive",
										ValidateFunc: validation.StringInSlice([]string{"transitive", "non-transitive"}, false),
										Description:  "The community scope. Default: transitive.",
									},
									"description": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The description of the community factor.",
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

func validateMatchRulePolicy(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	// Defer cross-field checks while the configured values are unknown.
	rawConfig := d.GetRawConfig()
	if !rawConfig.IsKnown() || rawConfig.IsNull() {
		return nil
	}
	if rawConfig.GetAttr("match_prefix").IsWhollyKnown() {
		prefixes := make(map[string]bool)
		for _, item := range d.Get("match_prefix").(*schema.Set).List() {
			prefix := item.(map[string]interface{})
			value := prefix["prefix"].(string)
			if prefixes[value] {
				return fmt.Errorf("match_prefix contains duplicate prefix %q", value)
			}
			prefixes[value] = true
			_, network, err := net.ParseCIDR(value)
			if err != nil {
				// Invalid CIDRs are reported by the prefix attribute validator.
				continue
			}
			mask, max := network.Mask.Size()
			from := prefix["from_prefix"].(int)
			to := prefix["to_prefix"].(int)
			if from > max || to > max {
				return fmt.Errorf("match_prefix %q prefix lengths must be between 0 and %d", value, max)
			}
			if !prefix["aggregate"].(bool) {
				if from != 0 || to != 0 {
					return fmt.Errorf("match_prefix %q from_prefix and to_prefix must be 0 when aggregate is false", value)
				}
			} else if (from != 0 && from <= mask) || to < from {
				return fmt.Errorf("match_prefix %q requires from_prefix to be 0 or greater than the prefix mask, and less than or equal to to_prefix, when aggregate is true", value)
			}
		}
	}
	if rawConfig.GetAttr("match_community_term").IsWhollyKnown() {
		terms := make(map[string]bool)
		for _, item := range d.Get("match_community_term").(*schema.Set).List() {
			term := item.(map[string]interface{})
			name := term["name"].(string)
			if terms[name] {
				return fmt.Errorf("match_community_term contains duplicate name %q", name)
			}
			terms[name] = true
			factors := make(map[string]bool)
			for _, entry := range term["match_community_factor"].(*schema.Set).List() {
				factor := entry.(map[string]interface{})["community"].(string)
				if factors[factor] {
					return fmt.Errorf("match_community_term %q contains duplicate community %q", name, factor)
				}
				factors[factor] = true
			}
		}
	}
	return nil
}

func getMatchRulePolicyPrefixesPayload(prefixes *schema.Set) []map[string]interface{} {
	payload := make([]map[string]interface{}, 0, prefixes.Len())
	for _, item := range prefixes.List() {
		prefix := item.(map[string]interface{})
		payload = append(payload, map[string]interface{}{
			"prefix":      prefix["prefix"],
			"description": prefix["description"],
			"aggregate":   prefix["aggregate"],
			"fromPfxLen":  prefix["from_prefix"],
			"toPfxLen":    prefix["to_prefix"],
		})
	}
	return payload
}

func getMatchRulePolicyCommunityTermsPayload(terms *schema.Set, existing *container.Container) ([]map[string]interface{}, error) {
	payload := make([]map[string]interface{}, 0, terms.Len())
	for _, item := range terms.List() {
		term := item.(map[string]interface{})
		value := make(map[string]interface{})
		// Retain NDO-generated metadata for community terms that remain in the policy.
		if existing != nil {
			oldTerm, err := GetPolicyByName(existing, term["name"].(string), "matchCommunityTermsList")
			if err != nil && !isPolicyNotFound(err) {
				return nil, err
			}
			if err == nil {
				for key, entry := range oldTerm.Data().(map[string]interface{}) {
					value[key] = entry
				}
			}
		}
		factors := term["match_community_factor"].(*schema.Set)
		communities := make([]map[string]interface{}, 0, factors.Len())
		for _, entry := range factors.List() {
			factor := entry.(map[string]interface{})
			communities = append(communities, map[string]interface{}{
				"community":   factor["community"],
				"scope":       factor["scope"],
				"description": factor["description"],
			})
		}
		value["name"] = term["name"]
		value["description"] = term["description"]
		value["matchCommunityList"] = communities
		payload = append(payload, value)
	}
	return payload, nil
}

func setMatchRulePolicyData(d *schema.ResourceData, policy *container.Container, templateID string) error {
	name, ok := policy.S("name").Data().(string)
	if !ok || name == "" {
		return fmt.Errorf("Match Rule Policy response does not contain a valid name")
	}
	description, _ := policy.S("description").Data().(string)
	uuid, ok := policy.S("uuid").Data().(string)
	if !ok || uuid == "" {
		return fmt.Errorf("Match Rule Policy response does not contain a valid UUID")
	}
	prefixes := make([]map[string]interface{}, 0)
	if policy.S("matchPrefixList").Data() != nil {
		children, err := policy.S("matchPrefixList").Children()
		if err != nil {
			return fmt.Errorf("read Match Rule Policy prefixes: %w", err)
		}
		for _, child := range children {
			prefix, ok := child.S("prefix").Data().(string)
			if !ok {
				return fmt.Errorf("Match Rule Policy prefix response does not contain a valid prefix")
			}
			description, _ := child.S("description").Data().(string)
			aggregate, _ := child.S("aggregate").Data().(bool)
			from, _ := child.S("fromPfxLen").Data().(float64)
			to, _ := child.S("toPfxLen").Data().(float64)
			prefixes = append(prefixes, map[string]interface{}{
				"prefix":      prefix,
				"description": description,
				"aggregate":   aggregate,
				"from_prefix": int(from),
				"to_prefix":   int(to),
			})
		}
	}
	terms := make([]map[string]interface{}, 0)
	if policy.S("matchCommunityTermsList").Data() != nil {
		children, err := policy.S("matchCommunityTermsList").Children()
		if err != nil {
			return fmt.Errorf("read Match Rule Policy community terms: %w", err)
		}
		for _, child := range children {
			termName, ok := child.S("name").Data().(string)
			if !ok {
				return fmt.Errorf("Match Rule Policy community term response does not contain a valid name")
			}
			termDescription, _ := child.S("description").Data().(string)
			factors := make([]map[string]interface{}, 0)
			if child.S("matchCommunityList").Data() != nil {
				communities, err := child.S("matchCommunityList").Children()
				if err != nil {
					return fmt.Errorf("read Match Rule Policy community factors: %w", err)
				}
				for _, community := range communities {
					value, ok := community.S("community").Data().(string)
					if !ok {
						return fmt.Errorf("Match Rule Policy community factor response does not contain a valid community")
					}
					scope, _ := community.S("scope").Data().(string)
					if scope == "" {
						scope = "transitive"
					}
					description, _ := community.S("description").Data().(string)
					factors = append(factors, map[string]interface{}{
						"community":   value,
						"scope":       scope,
						"description": description,
					})
				}
			}
			terms = append(terms, map[string]interface{}{
				"name":                   termName,
				"description":            termDescription,
				"match_community_factor": factors,
			})
		}
	}
	if err := d.Set("template_id", templateID); err != nil {
		return err
	}
	if err := d.Set("name", name); err != nil {
		return err
	}
	if err := d.Set("description", description); err != nil {
		return err
	}
	if err := d.Set("uuid", uuid); err != nil {
		return err
	}
	if err := d.Set("match_prefix", prefixes); err != nil {
		return err
	}
	if err := d.Set("match_community_term", terms); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("templateId/%s/MatchRulePolicy/%s", templateID, name))
	return nil
}

func resourceMSOMatchRulePolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	if err := resourceMSOMatchRulePolicyRead(d, m); err != nil {
		return nil, err
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("Match Rule Policy %s not found", importID)
	}
	return []*schema.ResourceData{d}, nil
}

func resourceMSOMatchRulePolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Match Rule Policy Resource - Beginning Create")
	msoClient := m.(*client.Client)
	templateID := d.Get("template_id").(string)
	terms, err := getMatchRulePolicyCommunityTermsPayload(d.Get("match_community_term").(*schema.Set), nil)
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"name":                    d.Get("name").(string),
		"description":             d.Get("description").(string),
		"matchPrefixList":         getMatchRulePolicyPrefixesPayload(d.Get("match_prefix").(*schema.Set)),
		"matchCommunityTermsList": terms,
	}
	if _, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateID), models.GetPatchPayload("add", "/tenantPolicyTemplate/template/matchRulePolicies/-", payload)); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("templateId/%s/MatchRulePolicy/%s", templateID, d.Get("name").(string)))
	return resourceMSOMatchRulePolicyRead(d, m)
}

func resourceMSOMatchRulePolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Match Rule Policy Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)
	templateID, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return err
	}
	name, err := GetPolicyNameFromResourceId(d.Id(), "MatchRulePolicy")
	if err != nil {
		return err
	}
	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
	if err != nil {
		return err
	}
	policy, err := GetPolicyByName(response, name, "tenantPolicyTemplate", "template", "matchRulePolicies")
	if err != nil {
		if isPolicyNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	return setMatchRulePolicyData(d, policy, templateID)
}

func resourceMSOMatchRulePolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Match Rule Policy Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateID := d.Get("template_id").(string)
	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
	if err != nil {
		return err
	}
	index, err := GetPolicyIndexByKeyAndValue(response, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "matchRulePolicies")
	if err != nil {
		return err
	}
	policy, err := response.ArrayElement(index, "tenantPolicyTemplate", "template", "matchRulePolicies")
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/tenantPolicyTemplate/template/matchRulePolicies/%d", index)
	payload := container.New()
	payload.Array()
	if d.HasChange("name") {
		if err := addPatchPayloadToContainer(payload, "add", fmt.Sprintf("%s/name", path), d.Get("name").(string)); err != nil {
			return err
		}
	}
	if d.HasChange("description") {
		if err := addPatchPayloadToContainer(payload, "add", fmt.Sprintf("%s/description", path), d.Get("description").(string)); err != nil {
			return err
		}
	}
	if d.HasChange("match_prefix") {
		if err := addPatchPayloadToContainer(payload, "add", fmt.Sprintf("%s/matchPrefixList", path), getMatchRulePolicyPrefixesPayload(d.Get("match_prefix").(*schema.Set))); err != nil {
			return err
		}
	}
	if d.HasChange("match_community_term") {
		terms, err := getMatchRulePolicyCommunityTermsPayload(d.Get("match_community_term").(*schema.Set), policy)
		if err != nil {
			return err
		}
		if err := addPatchPayloadToContainer(payload, "add", fmt.Sprintf("%s/matchCommunityTermsList", path), terms); err != nil {
			return err
		}
	}
	if err := doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateID), payload); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("templateId/%s/MatchRulePolicy/%s", templateID, d.Get("name").(string)))
	return resourceMSOMatchRulePolicyRead(d, m)
}

func resourceMSOMatchRulePolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Match Rule Policy Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)
	templateID := d.Get("template_id").(string)
	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
	if err != nil {
		return err
	}
	index, err := GetPolicyIndexByKeyAndValue(response, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "matchRulePolicies")
	if err != nil {
		if isPolicyNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	payload := models.GetRemovePatchPayload(fmt.Sprintf("/tenantPolicyTemplate/template/matchRulePolicies/%d", index))
	if _, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateID), payload); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
