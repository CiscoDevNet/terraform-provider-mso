package mso

import (
	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceMSODHCPRelayPolicyProvider() *schema.Resource {
	return &schema.Resource{
		Read:          dataSourceMSODHCPRelayPolicyProviderRead,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"dhcp_relay_policy_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"dhcp_server_address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"epg_ref": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"external_epg_ref"},
				AtLeastOneOf:  []string{"epg_ref", "external_epg_ref"},
			},
			"external_epg_ref": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceMSODHCPRelayPolicyProviderRead(d *schema.ResourceData, m interface{}) error {
	msoClient := m.(*client.Client)
	dhcpRelayPolicyProvider := models.DHCPRelayPolicyProvider{
		Addr:       d.Get("dhcp_server_address").(string),
		PolicyName: d.Get("dhcp_relay_policy_name").(string),
	}
	if epg, ok := d.GetOk("epg_ref"); ok {
		dhcpRelayPolicyProvider.EpgRef = epg.(string)
	}
	if externalEpg, ok := d.GetOk("external_epg_ref"); ok {
		dhcpRelayPolicyProvider.ExternalEpgRef = externalEpg.(string)
	}
	remoteDHCPRelayPolicyProvider, err := msoClient.ReadDHCPRelayPolicyProvider(&dhcpRelayPolicyProvider)
	if err != nil {
		return err
	}
	setDHCPRelayPolicyProviderAttr(d, remoteDHCPRelayPolicyProvider)
	d.SetId(DHCPRelayPolicyProviderModeltoId(remoteDHCPRelayPolicyProvider))
	return nil
}
