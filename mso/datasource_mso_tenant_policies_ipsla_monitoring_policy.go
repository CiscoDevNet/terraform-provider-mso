package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSOIPSLAMonitoringPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOIPSLAMonitoringPolicyRead,

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
			"sla_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destination_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"http_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"http_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sla_frequency": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"detect_multiplier": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"request_data_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type_of_service": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"operation_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"threshold": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ipv6_traffic_class": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceMSOIPSLAMonitoringPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO IPSLA Monitoring Policy Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "ipslaMonitoringPolicies")
	if err != nil {
		return err
	}

	setIPSLAMonitoringPolicyData(d, policy, templateId)
	log.Printf("[DEBUG] MSO IPSLA Monitoring Policy Data Source - Read Complete : %v", d.Id())
	return nil
}
