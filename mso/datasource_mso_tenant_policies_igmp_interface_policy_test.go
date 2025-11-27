package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOTenantPoliciesIGMPInterfacePolicyDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: IGMP Interface Policy Data Source") },
				Config:    testAccMSOTenantPoliciesIGMPInterfacePolicyDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_tenant_policies_igmp_interface_policy.igmp_policy", "name", "test_igmp_interface_policy_max"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_igmp_interface_policy.igmp_policy", "description", "Maximum Values Test"),
					resource.TestCheckResourceAttrSet("data.mso_tenant_policies_igmp_interface_policy.igmp_policy", "uuid"),
					resource.TestCheckResourceAttrSet("data.mso_tenant_policies_igmp_interface_policy.igmp_policy", "template_id"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_igmp_interface_policy.igmp_policy", "version3_asm", "true"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_igmp_interface_policy.igmp_policy", "fast_leave", "true"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_igmp_interface_policy.igmp_policy", "report_link_local_groups", "true"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_igmp_interface_policy.igmp_policy", "igmp_version", "v3"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_igmp_interface_policy.igmp_policy", "group_timeout", "65535"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_igmp_interface_policy.igmp_policy", "query_interval", "18000"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_igmp_interface_policy.igmp_policy", "query_response_interval", "25"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_igmp_interface_policy.igmp_policy", "last_member_count", "5"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_igmp_interface_policy.igmp_policy", "last_member_response_time", "25"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_igmp_interface_policy.igmp_policy", "startup_query_count", "10"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_igmp_interface_policy.igmp_policy", "startup_query_interval", "18000"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_igmp_interface_policy.igmp_policy", "querier_timeout", "65535"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_igmp_interface_policy.igmp_policy", "robustness_variable", "7"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_igmp_interface_policy.igmp_policy", "maximum_multicast_entries", "4294967295"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_igmp_interface_policy.igmp_policy", "reserved_multicast_entries", "4294967295"),
					resource.TestCheckResourceAttrSet("data.mso_tenant_policies_igmp_interface_policy.igmp_policy", "static_report_route_map_uuid"),
				),
			},
		},
	})
}

func testAccMSOTenantPoliciesIGMPInterfacePolicyDataSource() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_route_map_policy_multicast" "static_report" {
        template_id = mso_template.template_tenant.id
        name        = "test_static_report_route_map"
        description = "Static Report Route Map for IGMP"
		route_map_multicast_entries {
			order                   = 1
			group_ip                = "226.2.2.2/8"
			source_ip               = "1.1.1.1/1"
			rendezvous_point_ip     = "1.1.1.2"
			action                  = "permit"
		}
    }

    resource "mso_tenant_policies_igmp_interface_policy" "igmp_policy" {
        template_id                    = mso_template.template_tenant.id
        name                           = "test_igmp_interface_policy_max"
        description                    = "Maximum Values Test"
        version3_asm                   = true
        fast_leave                     = true
        report_link_local_groups       = true
        igmp_version                   = "v3"
        group_timeout                  = 65535
        query_interval                 = 18000
        query_response_interval        = 25
        last_member_count              = 5
        last_member_response_time      = 25
        startup_query_count            = 10
        startup_query_interval         = 18000
        querier_timeout                = 65535
        robustness_variable            = 7
        maximum_multicast_entries      = 4294967295
        reserved_multicast_entries     = 4294967295
        static_report_route_map_uuid   = mso_tenant_policies_route_map_policy_multicast.static_report.uuid
    }

    data "mso_tenant_policies_igmp_interface_policy" "igmp_policy" {
        template_id = mso_template.template_tenant.id
        name        = mso_tenant_policies_igmp_interface_policy.igmp_policy.name
    }`, testAccMSOTemplateResourceTenantConfig())
}
