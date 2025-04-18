package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSOTenantPoliciesDHCPRelayPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOTenantPoliciesDHCPRelayPolicyRead,

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
			"dhcp_relay_providers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dhcp_server_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_epg_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"external_epg_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dhcp_server_vrf_preference": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceMSOTenantPoliciesDHCPRelayPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO DHCP Relay Policy Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "dhcpRelayPolicies")
	if err != nil {
		return err
	}

	setDHCPRelayPolicyData(d, policy, templateId)
	log.Printf("[DEBUG] MSO DHCP Relay Policy Data Source - Read Complete : %v", d.Id())
	return nil
}
