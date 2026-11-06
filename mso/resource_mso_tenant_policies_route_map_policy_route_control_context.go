package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceMSORouteMapPolicyContext() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSORouteMapPolicyContextCreate,
		Read:   resourceMSORouteMapPolicyContextRead,
		Update: resourceMSORouteMapPolicyContextUpdate,
		Delete: resourceMSORouteMapPolicyContextDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSORouteMapPolicyContextImport,
		},

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the tenant policy template.",
			},
			"route_map_policy_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The UUID of the parent Route Map Policy.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
				Description:  "The name of the Route Control context.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the context.",
			},
			// Default is explicitly set because the API requires 'order' in the payload
			// but does not assign a default value if omitted.
			"order": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 9),
				Description:  "The order of the context. Range: 0-9.",
			},
			// Default is explicitly set because the API requires 'action' in the payload
			// but does not assign a default value if omitted.
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "permit",
				ValidateFunc: validation.StringInSlice([]string{
					"permit", "deny",
				}, false),
				Description: "The action of the context.",
			},
			"set_rule_uuid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The UUID of the Set Rule Policy (max one).",
			},
			"match_rule_uuids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A set of Match Rule Policy UUIDs.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func buildMatchRuleUUIDsPayload(matchRulesRaw interface{}) []string {
	if matchRulesRaw == nil {
		return nil
	}
	set := matchRulesRaw.(*schema.Set).List()
	uuids := make([]string, len(set))
	for i, item := range set {
		uuids[i] = item.(string)
	}
	return uuids
}

func resourceMSORouteMapPolicyContextImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO Route Map Policy Context Resource - Beginning Import: %v", d.Id())
	importID := d.Id()
	parentID, contextName, err := ParseChildResourceId(importID, "/context/")
	if err != nil {
		return nil, err
	}
	templateID, err := GetTemplateIdFromResourceId(parentID)
	if err != nil {
		return nil, err
	}
	policyUUID, err := GetPolicyNameFromResourceId(parentID, "RouteMapPolicy")
	if err != nil {
		return nil, err
	}
	if err := d.Set("template_id", templateID); err != nil {
		return nil, err
	}
	if err := d.Set("route_map_policy_uuid", policyUUID); err != nil {
		return nil, err
	}
	if err := d.Set("name", contextName); err != nil {
		return nil, err
	}

	err = resourceMSORouteMapPolicyContextRead(d, m)
	if err != nil {
		return nil, err
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("cannot import mso_tenant_policies_route_map_policy_route_control_context %q: resource not found", importID)
	}
	log.Printf("[DEBUG] MSO Route Map Policy Context Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSORouteMapPolicyContextCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Route Map Policy Context Resource - Beginning Create")
	msoClient := m.(*client.Client)
	templateID := d.Get("template_id").(string)
	policyUUID := d.Get("route_map_policy_uuid").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", policyUUID, "tenantPolicyTemplate", "template", "routeMapPolicies")
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"order":       d.Get("order").(int),
		"action":      d.Get("action").(string),
	}

	if v, ok := d.GetOk("set_rule_uuid"); ok {
		payload["setRuleRef"] = v.(string)
	}

	if v, ok := d.GetOk("match_rule_uuids"); ok {
		payload["matchRules"] = buildMatchRuleUUIDsPayload(v)
	}

	path := fmt.Sprintf("/tenantPolicyTemplate/template/routeMapPolicies/%d/contexts/-", policyIndex)
	payloadModel := models.GetPatchPayload("add", path, payload)

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateID), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/RouteMapPolicy/%s/context/%s", templateID, policyUUID, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO Route Map Policy Context Resource - Create Complete: %v", d.Id())
	return resourceMSORouteMapPolicyContextRead(d, m)
}

func resourceMSORouteMapPolicyContextRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Route Map Policy Context Resource - Beginning Read")
	msoClient := m.(*client.Client)
	templateID := d.Get("template_id").(string)
	policyUUID := d.Get("route_map_policy_uuid").(string)

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(response, "uuid", policyUUID, "tenantPolicyTemplate", "template", "routeMapPolicies")
	if isPolicyNotFound(err) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	policy := response.S("tenantPolicyTemplate", "template", "routeMapPolicies").Index(policyIndex)
	match, err := getRouteMapPolicyContext(policy, d.Get("name").(string), policyUUID)
	if isPolicyNotFound(err) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}
	if err := setRouteMapPolicyContextData(d, templateID, policyUUID, match); err != nil {
		return err
	}
	log.Printf("[DEBUG] MSO Route Map Policy Context Resource - Read Complete: %v", d.Id())
	return nil
}

func getRouteMapPolicyContext(policy *container.Container, contextName, policyName string) (*container.Container, error) {
	match, err := GetPolicyByName(policy, contextName, "contexts")
	if isPolicyNotFound(err) {
		return nil, &policyNotFoundError{message: fmt.Sprintf("Route Map Policy Context %q not found in parent policy %q", contextName, policyName)}
	}
	return match, err
}

func setRouteMapPolicyContextData(d *schema.ResourceData, templateID, policyUUID string, context *container.Container) error {
	name, ok := context.S("name").Data().(string)
	if !ok || name == "" {
		return fmt.Errorf("Route Map Policy Context response does not contain a valid name")
	}
	description, _ := context.S("description").Data().(string)
	setRuleUUID, _ := context.S("setRuleRef").Data().(string)
	action, _ := context.S("action").Data().(string)
	if action == "" {
		action = "permit"
	}
	order, _ := context.S("order").Data().(float64)
	matchRuleUUIDs := make([]string, 0)
	if context.S("matchRules").Data() != nil {
		matches, err := context.S("matchRules").Children()
		if err != nil {
			return fmt.Errorf("read Route Map Policy Context match rules: %w", err)
		}
		for _, match := range matches {
			uuid, ok := match.Data().(string)
			if !ok {
				return fmt.Errorf("Route Map Policy Context match rule response does not contain a valid UUID")
			}
			matchRuleUUIDs = append(matchRuleUUIDs, uuid)
		}
	}
	if err := d.Set("template_id", templateID); err != nil {
		return err
	}
	if err := d.Set("route_map_policy_uuid", policyUUID); err != nil {
		return err
	}
	if err := d.Set("name", name); err != nil {
		return err
	}
	if err := d.Set("description", description); err != nil {
		return err
	}
	if err := d.Set("set_rule_uuid", setRuleUUID); err != nil {
		return err
	}
	if err := d.Set("match_rule_uuids", matchRuleUUIDs); err != nil {
		return err
	}
	if err := d.Set("order", int(order)); err != nil {
		return err
	}
	if err := d.Set("action", action); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("templateId/%s/RouteMapPolicy/%s/context/%s", templateID, policyUUID, name))
	return nil
}

func resourceMSORouteMapPolicyContextUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Route Map Policy Context Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateID := d.Get("template_id").(string)
	policyUUID := d.Get("route_map_policy_uuid").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", policyUUID, "tenantPolicyTemplate", "template", "routeMapPolicies")
	if err != nil {
		return err
	}

	policy := templateCont.S("tenantPolicyTemplate", "template", "routeMapPolicies").Index(policyIndex)
	contextIndex, err := GetPolicyIndexByKeyAndValue(policy, "name", d.Get("name").(string), "contexts")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/tenantPolicyTemplate/template/routeMapPolicies/%d/contexts/%d", policyIndex, contextIndex)
	payloadCont := container.New()
	payloadCont.Array()

	if d.HasChange("description") {
		if err := addPatchPayloadToContainer(payloadCont, "add", fmt.Sprintf("%s/description", updatePath), d.Get("description").(string)); err != nil {
			return err
		}
	}
	if d.HasChange("order") {
		if err := addPatchPayloadToContainer(payloadCont, "add", fmt.Sprintf("%s/order", updatePath), d.Get("order").(int)); err != nil {
			return err
		}
	}
	if d.HasChange("action") {
		if err := addPatchPayloadToContainer(payloadCont, "add", fmt.Sprintf("%s/action", updatePath), d.Get("action").(string)); err != nil {
			return err
		}
	}
	if d.HasChange("set_rule_uuid") {
		if err := addPatchPayloadToContainer(payloadCont, "add", fmt.Sprintf("%s/setRuleRef", updatePath), d.Get("set_rule_uuid").(string)); err != nil {
			return err
		}
	}
	if d.HasChange("match_rule_uuids") {
		if err := addPatchPayloadToContainer(payloadCont, "add", fmt.Sprintf("%s/matchRules", updatePath), buildMatchRuleUUIDsPayload(d.Get("match_rule_uuids"))); err != nil {
			return err
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateID), payloadCont)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] MSO Route Map Policy Context Resource - Update Complete: %v", d.Id())
	return resourceMSORouteMapPolicyContextRead(d, m)
}

func resourceMSORouteMapPolicyContextDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Route Map Policy Context Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)
	templateID := d.Get("template_id").(string)
	policyUUID := d.Get("route_map_policy_uuid").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", policyUUID, "tenantPolicyTemplate", "template", "routeMapPolicies")
	if isPolicyNotFound(err) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	policy := templateCont.S("tenantPolicyTemplate", "template", "routeMapPolicies").Index(policyIndex)
	contextIndex, err := GetPolicyIndexByKeyAndValue(policy, "name", d.Get("name").(string), "contexts")
	if isPolicyNotFound(err) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/tenantPolicyTemplate/template/routeMapPolicies/%d/contexts/%d", policyIndex, contextIndex)
	payloadModel := models.GetRemovePatchPayload(path)

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateID), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO Route Map Policy Context Resource - Delete Complete")
	return nil
}
