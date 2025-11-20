package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSOCustomQoSPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOCustomQoSPolicyRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the tenant policy template.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Custom QoS Policy.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the Custom QoS Policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the Custom QoS Policy.",
			},
			"dscp_mappings": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The DSCP mappings of the Custom QoS Policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dscp_from": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The starting encoding point of the DSCP range.",
						},
						"dscp_to": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ending encoding point of the DSCP range.",
						},
						"dscp_target": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The DSCP target encoding point for egressing traffic.",
						},
						"target_cos": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The target CoS value/traffic type for egressing traffic.",
						},
						"qos_priority": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The QoS priority level.",
						},
					},
				},
			},
			"cos_mappings": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The CoS mappings of the Custom QoS Policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dot1p_from": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The starting value/traffic type of the CoS range.",
						},
						"dot1p_to": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ending value/traffic type of the CoS range.",
						},
						"dscp_target": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The DSCP target encoding point for egressing traffic.",
						},
						"target_cos": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The target CoS value/traffic type for egressing traffic.",
						},
						"qos_priority": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The QoS priority level.",
						},
					},
				},
			},
		},
	}
}

func dataSourceMSOCustomQoSPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Custom QoS Policy Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "qosPolicies")
	if err != nil {
		return err
	}

	setCustomQoSPolicyData(d, policy, templateId)
	log.Printf("[DEBUG] MSO Custom QoS Policy Data Source - Read Complete: %v", d.Id())
	return nil
}
