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
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "name", "tf_test_route_map_policy_multicast"),
					resource.TestCheckResourceAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "description", "Terraform test Route Map Policy for Multicast adding extra entry"),
					resource.TestCheckResourceAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "route_map_multicast_entries.#", "2"),
					customTestCheckResourceTypeSetAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "route_map_multicast_entries",
						map[string]string{
							"order":               "1",
							"group_ip":            "226.2.2.2/8",
							"source_ip":           "1.1.1.1/1",
							"rendezvous_point_ip": "1.1.1.2",
							"action":              "permit",
						},
					),
					customTestCheckResourceTypeSetAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "route_map_multicast_entries",
						map[string]string{
							"order":               "2",
							"group_ip":            "226.3.3.3/24",
							"source_ip":           "2.2.2.2/2",
							"rendezvous_point_ip": "2.2.2.3",
							"action":              "deny",
						},
					),
				),
			},
		},
	})
}

func testAccMSOTenantPoliciesMcastRouteMapPolicyDataSource() string {
	return fmt.Sprintf(`%s
	data "mso_tenant_policies_route_map_policy_multicast" "route_map_policy_multicast" {
	    template_id        = mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast.template_id
	    name               = "tf_test_route_map_policy_multicast"
    }`, testAccMSOTenantPoliciesMcastRouteMapPolicyConfigUpdateAddingExtraEntry())
}
