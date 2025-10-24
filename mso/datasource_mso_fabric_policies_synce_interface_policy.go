package mso

import (
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSOSyncEInterfacePolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOSyncEInterfacePolicyRead,

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
			"sync_state_msg": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"selection_input": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"src_priority": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"wait_to_restore": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceMSOSyncEInterfacePolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO VLAN Pool Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	setSyncEInterfacePolicyData(d, msoClient, templateId, policyName)
	log.Printf("[DEBUG] MSO VLAN Pool Data Source - Read Complete : %v", d.Id())
	return nil
}
