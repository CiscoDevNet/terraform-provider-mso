package mso

import (
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSOPtpPolicyProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePtpPolicyProfileRead,

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
			"profile_template": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delay_interval": {
				Type:         schema.TypeInt,
				Computed: true,
			},
			"sync_interval": {
				Type:         schema.TypeInt,
				Computed: true,
			},
			"announce_interval": {
				Type:         schema.TypeInt,
				Computed: true,
			},
			"announce_timeout": {
				Type:         schema.TypeInt,
				Computed: true,
			},
			"override_node_profile": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourcePtpPolicyProfileRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO PTP Policy Profile Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	setPtpPolicyProfileData(d, msoClient, templateId, policyName)
	log.Printf("[DEBUG] MSO PTP Policy Profile Data Source - Read Complete : %v", d.Id())
	return nil
}
