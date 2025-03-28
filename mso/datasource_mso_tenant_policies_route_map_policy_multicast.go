package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSOMcastRouteMapPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOMcastRouteMapPolicyRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"route_map_entries_multicast": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"order": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"group_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rp_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceMSOMcastRouteMapPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Route Map Policy for Multicast Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "mcastRouteMapPolicies")
	if err != nil {
		return err
	}

	setMcastRouteMapPolicyData(d, policy, templateId)
	log.Printf("[DEBUG] MSO Route Map Policy for Multicast Data Source - Read Complete : %v", d.Id())
	return nil
}
