package mso

import (
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSODHCPOptionPolicy() *schema.Resource {
	return &schema.Resource{

		Read: datasourceMSODHCPOptionPolicyRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"option": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		}),
	}
}

func datasourceMSODHCPOptionPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] dhcp_option_policy: Beginning Import")
	msoClient := m.(*client.Client)
	name := d.Get("name").(string)
	id, err := msoClient.GetDHCPOptionPolicyID(name)
	if err != nil {
		return err
	}
	DHCPOptionPolicy, err := getDHCPOptionPolicy(msoClient, id)
	if err != nil {
		return err
	}
	setDHCPOptionPolicy(DHCPOptionPolicy, d)
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return nil
}
