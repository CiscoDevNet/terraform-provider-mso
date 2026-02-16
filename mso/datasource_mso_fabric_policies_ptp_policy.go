package mso

import (
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSOPtpPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePtpPolicyRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
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
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"admin_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"global_priority1": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"global_priority2": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"global_domain": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"fabric_profile_template": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fabric_announce_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"fabric_sync_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"fabric_delay_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"fabric_announce_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ptp_profiles": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"profile_template": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delay_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"sync_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"announce_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"announce_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"override_node_profile": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePtpPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO PTP Policy Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	setPtpPolicyData(d, msoClient, templateId, policyName)
	log.Printf("[DEBUG] MSO PTP Policy Data Source - Read Complete : %v", d.Id())
	return nil
}
