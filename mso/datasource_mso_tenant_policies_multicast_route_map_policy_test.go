package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOTenantPoliciesMcastRouteMapPolicyDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Multicast Route Map Policy Data Source") },
				Config:    testAccMSOTenantPoliciesMcastRouteMapPolicyDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "name", "test_multicast_route_map_policy"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "description", "test Route Map Policy for Multicast"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "multicast_route_map_entries.0.order", "1"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "multicast_route_map_entries.0.group_ip", "226.2.2.2/8"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "multicast_route_map_entries.0.source_ip", "1.1.1.1/1"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "multicast_route_map_entries.0.rp_ip", "1.1.1.2"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "multicast_route_map_entries.0.action", "permit"),
				),
			},
		},
	})
}

func testAccMSOTenantPoliciesMcastRouteMapPolicyDataSource() string {
	return fmt.Sprintf(`%s
	data "mso_tenant_policies_multicast_route_map_policy" "multicast_route_map_policy" {
	    template_id        = mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy.template_id
	    name               = "test_multicast_route_map_policy"
    }`, testAccMSOTenantPoliciesMcastRouteMapPolicyConfigCreate())
}
