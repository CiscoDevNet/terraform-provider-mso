package mso

import (
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceMSODHCPRelayPolicy() *schema.Resource {
	return &schema.Resource{

		Read: datasourceMSODHCPRelayPolicyRead,

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

			"dhcp_relay_policy_provider": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"epg": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"external_epg": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dhcp_server_address": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsIPAddress,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		}),
	}
}

func datasourceMSODHCPRelayPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] dhcp_relay_policy: Beginning Import")
	msoClient := m.(*client.Client)
	name := d.Get("name").(string)
	id, err := msoClient.GetDHCPRelayPolicyID(name)
	if err != nil {
		return err
	}
	DHCPRelayPolicy, err := getDHCPRelayPolicy(msoClient, id)
	if err != nil {
		return err
	}
	setDHCPRelayPolicy(DHCPRelayPolicy, d)
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return nil
}
