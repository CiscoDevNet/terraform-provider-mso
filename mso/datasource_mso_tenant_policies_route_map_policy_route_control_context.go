package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSORouteMapPolicyContext() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSORouteMapPolicyContextRead,

		Schema: map[string]*schema.Schema{
			"parent_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"order": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"match_rule_uuids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"set_rule_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMSORouteMapPolicyContextRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Route Map Policy Context Data Source - Beginning Read")
	msoClient := m.(*client.Client)
	parentId := d.Get("parent_id").(string)
	contextName := d.Get("name").(string)

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
		return fmt.Errorf("Route Map Policy Context '%s' not found in parent policy '%s'", contextName, policyName)
	}

	d.SetId(fmt.Sprintf("%s/context/%s", parentId, models.StripQuotes(match.S("name").String())))
	d.Set("parent_id", parentId)
	d.Set("name", contextName)

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

	log.Printf("[DEBUG] MSO Route Map Policy Context Data Source - Read Complete: %v", d.Id())
	return nil
}
