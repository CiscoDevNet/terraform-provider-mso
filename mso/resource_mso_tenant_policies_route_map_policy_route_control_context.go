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
			"parent_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Terraform ID of the parent Route Map Policy.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the Route Control context.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the context.",
			},
			"order": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 9),
				Description:  "The order of the context. Range: 0-9.",
			},
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
	parentId, contextName, err := ParseChildResourceId(d.Id(), "/context/")
	if err != nil {
		return nil, err
	}
	d.Set("parent_id", parentId)
	d.Set("name", contextName)

	err = resourceMSORouteMapPolicyContextRead(d, m)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] MSO Route Map Policy Context Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSORouteMapPolicyContextCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Route Map Policy Context Resource - Beginning Create")
	msoClient := m.(*client.Client)
	parentId := d.Get("parent_id").(string)

	templateId, err := GetTemplateIdFromResourceId(parentId)
	if err != nil {
		return err
	}
	policyName, err := GetPolicyNameFromResourceId(parentId, "RouteMapPolicy")
	if err != nil {
		return err
	}

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "name", policyName, "tenantPolicyTemplate", "template", "routeMapPolicies")
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

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/context/%s", parentId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO Route Map Policy Context Resource - Create Complete: %v", d.Id())
	return resourceMSORouteMapPolicyContextRead(d, m)
}

func resourceMSORouteMapPolicyContextRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Route Map Policy Context Resource - Beginning Read")
	msoClient := m.(*client.Client)
	parentId := d.Get("parent_id").(string)

	templateId, err := GetTemplateIdFromResourceId(parentId)
	if err != nil {
		return err
	}
	policyName, err := GetPolicyNameFromResourceId(parentId, "RouteMapPolicy")
	if err != nil {
		return err
	}

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "routeMapPolicies")
	if err != nil {
		return err
	}

	contextName := d.Get("name").(string)
	contexts := policy.S("contexts")
	count, _ := contexts.ArrayCount()

	var match *container.Container
	for i := 0; i < count; i++ {
		c := contexts.Index(i)
		if models.StripQuotes(c.S("name").String()) == contextName {
			match = c
			break
		}
	}

	if match == nil {
		log.Printf("[DEBUG] MSO Route Map Policy Context Resource - Context '%s' not found in parent policy, clearing ID", contextName)
		d.SetId("")
		return nil
	}

	if match.Exists("description") {
		descValue := models.StripQuotes(match.S("description").String())
		if descValue == "{}" {
			d.Set("description", "")
		} else {
			d.Set("description", descValue)
		}
	} else {
		d.Set("description", "")
	}

	if match.Exists("order") {
		d.Set("order", int(match.S("order").Data().(float64)))
	}

	if match.Exists("action") {
		d.Set("action", models.StripQuotes(match.S("action").String()))
	}

	if match.Exists("setRuleRef") {
		setRuleValue := models.StripQuotes(match.S("setRuleRef").String())
		if setRuleValue == "{}" {
			d.Set("set_rule_uuid", "")
		} else {
			d.Set("set_rule_uuid", setRuleValue)
		}
	} else {
		d.Set("set_rule_uuid", "")
	}

	if match.Exists("matchRules") {
		apiList, _ := match.S("matchRules").Data().([]interface{})
		uuids := make([]string, len(apiList))
		for i, uuid := range apiList {
			uuids[i] = models.StripQuotes(fmt.Sprintf("%v", uuid))
		}
		d.Set("match_rule_uuids", uuids)
	} else {
		d.Set("match_rule_uuids", []string{})
	}

	log.Printf("[DEBUG] MSO Route Map Policy Context Resource - Read Complete: %v", d.Id())
	return nil
}

func resourceMSORouteMapPolicyContextUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Route Map Policy Context Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	parentId := d.Get("parent_id").(string)

	templateId, err := GetTemplateIdFromResourceId(parentId)
	if err != nil {
		return err
	}
	policyName, err := GetPolicyNameFromResourceId(parentId, "RouteMapPolicy")
	if err != nil {
		return err
	}

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "name", policyName, "tenantPolicyTemplate", "template", "routeMapPolicies")
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
		addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/description", updatePath), d.Get("description").(string))
	}
	if d.HasChange("order") {
		addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/order", updatePath), d.Get("order").(int))
	}
	if d.HasChange("action") {
		addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/action", updatePath), d.Get("action").(string))
	}
	if d.HasChange("set_rule_uuid") {
		addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/setRuleRef", updatePath), d.Get("set_rule_uuid").(string))
	}
	if d.HasChange("match_rule_uuids") {
		_, newVal := d.GetChange("match_rule_uuids")
		uuids := buildMatchRuleUUIDsPayload(newVal)

		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/matchRules", updatePath), uuids)
		if err != nil {
			return err
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] MSO Route Map Policy Context Resource - Update Complete: %v", d.Id())
	return resourceMSORouteMapPolicyContextRead(d, m)
}

func resourceMSORouteMapPolicyContextDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Route Map Policy Context Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)
	parentId := d.Get("parent_id").(string)

	templateId, err := GetTemplateIdFromResourceId(parentId)
	if err != nil {
		return err
	}
	policyName, err := GetPolicyNameFromResourceId(parentId, "RouteMapPolicy")
	if err != nil {
		return err
	}

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "name", policyName, "tenantPolicyTemplate", "template", "routeMapPolicies")
	if err != nil {
		return err
	}

	policy := templateCont.S("tenantPolicyTemplate", "template", "routeMapPolicies").Index(policyIndex)
	contextIndex, err := GetPolicyIndexByKeyAndValue(policy, "name", d.Get("name").(string), "contexts")
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/tenantPolicyTemplate/template/routeMapPolicies/%d/contexts/%d", policyIndex, contextIndex)
	payloadModel := models.GetRemovePatchPayload(path)

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO Route Map Policy Context Resource - Delete Complete")
	return nil
}
