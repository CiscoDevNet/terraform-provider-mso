package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOTenantPoliciesMcastRouteMapPolicyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create Multicast Route Map Policy") },
				Config:    testAccMSOTenantPoliciesMcastRouteMapPolicyConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "name", "test_multicast_route_map_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "description", "test Route Map Policy for Multicast"),
					resource.TestCheckResourceAttr("mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "multicast_route_map_entries.0.order", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "multicast_route_map_entries.0.group_ip", "226.2.2.2/8"),
					resource.TestCheckResourceAttr("mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "multicast_route_map_entries.0.source_ip", "1.1.1.1/1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "multicast_route_map_entries.0.rp_ip", "1.1.1.2"),
					resource.TestCheckResourceAttr("mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "multicast_route_map_entries.0.action", "permit"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Multicast Route Map Policy adding extra entry") },
				Config:    testAccMSOTenantPoliciesMcastRouteMapPolicyConfigUpdateAddingExtraEntry(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "name", "test_multicast_route_map_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "description", "test Route Map Policy for Multicast adding extra entry"),
					resource.TestCheckResourceAttr("mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "multicast_route_map_entries.0.order", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "multicast_route_map_entries.0.group_ip", "226.2.2.2/8"),
					resource.TestCheckResourceAttr("mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "multicast_route_map_entries.0.source_ip", "1.1.1.1/1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "multicast_route_map_entries.0.rp_ip", "1.1.1.2"),
					resource.TestCheckResourceAttr("mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "multicast_route_map_entries.0.action", "permit"),
					resource.TestCheckResourceAttr("mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "multicast_route_map_entries.1.order", "2"),
					resource.TestCheckResourceAttr("mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "multicast_route_map_entries.1.group_ip", "226.3.3.3/24"),
					resource.TestCheckResourceAttr("mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "multicast_route_map_entries.1.action", "deny"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Multicast Route Map Policy removing extra entry") },
				Config:    testAccMSOTenantPoliciesMcastRouteMapPolicyConfigUpdateRemovingExtraEntry(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "name", "test_multicast_route_map_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "description", "test Route Map Policy for Multicast removing extra entry"),
					resource.TestCheckResourceAttr("mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "multicast_route_map_entries.0.order", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "multicast_route_map_entries.0.group_ip", "226.2.2.2/8"),
					resource.TestCheckResourceAttr("mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "multicast_route_map_entries.0.source_ip", "1.1.1.1/1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "multicast_route_map_entries.0.rp_ip", "1.1.1.2"),
					resource.TestCheckResourceAttr("mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy", "multicast_route_map_entries.0.action", "permit"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Multicast Route Map Policy") },
				ResourceName:      "mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithArguments("mso_tenant_policies_multicast_route_map_policy", "McastRouteMapPolicy"),
	})
}

func testAccMSOTenantPoliciesMcastRouteMapPolicyConfigCreate() string {
	return fmt.Sprintf(`%s
	resource "mso_tenant_policies_multicast_route_map_policy" "multicast_route_map_policy" {
		template_id = mso_template.tenant_template.id
		name        = "test_multicast_route_map_policy"
		description = "test Route Map Policy for Multicast"
		multicast_route_map_entries {
			order     = 1
			group_ip  = "226.2.2.2/8"
			source_ip = "1.1.1.1/1"
			rp_ip     = "1.1.1.2"
			action    = "permit"
		}
	}`, testAccMSOTemplateResourceTenantConfig())

}

func testAccMSOTenantPoliciesMcastRouteMapPolicyConfigUpdateAddingExtraEntry() string {
	return fmt.Sprintf(`%s
	resource "mso_tenant_policies_multicast_route_map_policy" "multicast_route_map_policy" {
		template_id = mso_template.tenant_template.id
		name        = "test_multicast_route_map_policy"
		description = "test Route Map Policy for Multicast adding extra entry"
		multicast_route_map_entries {
			order     = 1
			group_ip  = "226.2.2.2/8"
			source_ip = "1.1.1.1/1"
			rp_ip     = "1.1.1.2"
			action    = "permit"
		}
		multicast_route_map_entries {
			order     = 2
			group_ip  = "226.3.3.3/24"
			action    = "deny"
		}
	}`, testAccMSOTemplateResourceTenantConfig())

}

func testAccMSOTenantPoliciesMcastRouteMapPolicyConfigUpdateRemovingExtraEntry() string {
	return fmt.Sprintf(`%s
	resource "mso_tenant_policies_multicast_route_map_policy" "multicast_route_map_policy" {
		template_id = mso_template.tenant_template.id
		name        = "test_multicast_route_map_policy"
		description = "test Route Map Policy for Multicast removing extra entry"
		multicast_route_map_entries {
			order     = 1
			group_ip  = "226.2.2.2/8"
			source_ip = "1.1.1.1/1"
			rp_ip     = "1.1.1.2"
			action    = "permit"
		}
	}`, testAccMSOTemplateResourceTenantConfig())
}
