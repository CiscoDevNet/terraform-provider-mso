package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOTenantPoliciesL3OutNodeRoutingPolicyDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: L3Out Node Routing Policy Data Source") },
				Config:    testAccMSOTenantPoliciesL3OutNodeRoutingPolicyDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_node_routing_policy.node_policy", "name", "test_node_routing_policy"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_node_routing_policy.node_policy", "description", "Test L3Out Node Routing Policy"),
					resource.TestCheckResourceAttrSet("data.mso_tenant_policies_l3out_node_routing_policy.node_policy", "uuid"),
					resource.TestCheckResourceAttrSet("data.mso_tenant_policies_l3out_node_routing_policy.node_policy", "template_id"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_node_routing_policy.node_policy", "as_path_multipath_relax", "false"),

					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_node_routing_policy.node_policy", "bfd_multi_hop_settings.#", "1"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_node_routing_policy.node_policy", "bfd_multi_hop_settings.0.admin_state", "enabled"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_node_routing_policy.node_policy", "bfd_multi_hop_settings.0.detection_multiplier", "3"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_node_routing_policy.node_policy", "bfd_multi_hop_settings.0.min_receive_interval", "250"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_node_routing_policy.node_policy", "bfd_multi_hop_settings.0.min_transmit_interval", "250"),

					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.#", "1"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.graceful_restart_helper", "true"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.keep_alive_interval", "60"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.hold_interval", "180"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.stale_interval", "300"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_l3out_node_routing_policy.node_policy", "bgp_node_settings.0.max_as_limit", "0"),
				),
			},
		},
	})
}

func testAccMSOTenantPoliciesL3OutNodeRoutingPolicyDataSource() string {
	return fmt.Sprintf(`%s
    data "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
        template_id = mso_tenant_policies_l3out_node_routing_policy.node_policy.template_id
        name        = "test_node_routing_policy"
    }`, testAccMSOTenantPoliciesL3OutNodeRoutingPolicyConfigCreate())
}
