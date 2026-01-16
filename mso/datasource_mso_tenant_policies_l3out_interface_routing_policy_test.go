package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOL3OutInterfaceRoutingPolicyDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: L3Out Interface Routing Policy Data Source") },
				Config:    testAccMSOL3OutInterfaceRoutingPolicyDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "name", "test_routing_policy"),
					resource.TestCheckResourceAttrSet("data.mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "uuid"),
					resource.TestCheckResourceAttrSet("data.mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "template_id"),

					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_multi_hop_settings.#", "1"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_multi_hop_settings.0.admin_state", "enabled"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_multi_hop_settings.0.detection_multiplier", "3"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_multi_hop_settings.0.min_receive_interval", "250"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_multi_hop_settings.0.min_transmit_interval", "250"),

					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_settings.#", "1"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_settings.0.admin_state", "enabled"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_settings.0.detection_multiplier", "3"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_settings.0.min_receive_interval", "50"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_settings.0.min_transmit_interval", "50"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_settings.0.echo_receive_interval", "50"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_settings.0.echo_admin_state", "enabled"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "bfd_settings.0.interface_control", "false"),

					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.#", "1"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.network_type", "broadcast"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.priority", "1"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_interface_routing_policy.routing_policy", "ospf_interface_settings.0.cost_of_interface", "0"),
				),
			},
		},
	})
}

func testAccMSOL3OutInterfaceRoutingPolicyDataSource() string {
	return fmt.Sprintf(`%s
    data "mso_tenant_policies_l3out_interface_routing_policy" "routing_policy" {
        template_id = mso_tenant_policies_l3out_interface_routing_policy.routing_policy.template_id
        name        = "test_routing_policy"
    }`, testAccMSOL3OutInterfaceRoutingPolicyConfigAddOSPF())
}
