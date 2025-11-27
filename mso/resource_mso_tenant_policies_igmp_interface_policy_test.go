package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOTenantPoliciesIGMPInterfacePolicyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create IGMP Interface Policy") },
				Config:    testAccMSOTenantPoliciesIGMPInterfacePolicyConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "name", "test_igmp_interface_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "description", "Test IGMP Interface Policy"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_igmp_interface_policy.igmp_policy", "uuid"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "version3_asm", "true"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "fast_leave", "true"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "report_link_local_groups", "true"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "igmp_version", "v3"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "group_timeout", "300"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "query_interval", "125"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "query_response_interval", "10"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "last_member_count", "2"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "last_member_response_time", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "startup_query_count", "2"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "startup_query_interval", "31"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "querier_timeout", "255"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "robustness_variable", "2"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "maximum_multicast_entries", "1000000"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "reserved_multicast_entries", "100000"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update IGMP Interface Policy - Change to v2") },
				Config:    testAccMSOTenantPoliciesIGMPInterfacePolicyConfigUpdateV2(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "name", "test_igmp_interface_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "description", "Updated to IGMP v2"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "version3_asm", "false"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "fast_leave", "false"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "report_link_local_groups", "false"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "igmp_version", "v2"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "group_timeout", "260"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "query_interval", "100"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update IGMP Interface Policy - With Route Maps") },
				Config:    testAccMSOTenantPoliciesIGMPInterfacePolicyConfigUpdateWithRouteMaps(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "name", "test_igmp_interface_policy"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_igmp_interface_policy.igmp_policy", "state_limit_route_map_uuid"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_igmp_interface_policy.igmp_policy", "report_policy_route_map_uuid"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_igmp_interface_policy.igmp_policy", "static_report_route_map_uuid"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update IGMP Interface Policy - Remove Route Maps") },
				Config:    testAccMSOTenantPoliciesIGMPInterfacePolicyConfigUpdateRemoveRouteMaps(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "name", "test_igmp_interface_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "state_limit_route_map_uuid", ""),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "report_policy_route_map_uuid", ""),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "static_report_route_map_uuid", ""),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update IGMP Interface Policy - Change Timer Values") },
				Config:    testAccMSOTenantPoliciesIGMPInterfacePolicyConfigUpdateTimers(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "name", "test_igmp_interface_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "group_timeout", "500"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "query_interval", "200"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "query_response_interval", "15"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "last_member_count", "3"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "last_member_response_time", "2"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "startup_query_count", "5"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "startup_query_interval", "50"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "querier_timeout", "300"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "robustness_variable", "3"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update IGMP Interface Policy - Maximum Values") },
				Config:    testAccMSOTenantPoliciesIGMPInterfacePolicyConfigUpdateMaxValues(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "name", "test_igmp_interface_policy_max"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "group_timeout", "65535"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "query_interval", "18000"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "query_response_interval", "25"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "last_member_count", "5"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "last_member_response_time", "25"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "startup_query_count", "10"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "startup_query_interval", "18000"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "querier_timeout", "65535"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "robustness_variable", "7"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "maximum_multicast_entries", "4294967295"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "reserved_multicast_entries", "4294967295"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update IGMP Interface Policy - Minimum Values") },
				Config:    testAccMSOTenantPoliciesIGMPInterfacePolicyConfigUpdateMinValues(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "name", "test_igmp_interface_policy_min"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "group_timeout", "3"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "query_interval", "2"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "query_response_interval", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "last_member_count", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "last_member_response_time", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "startup_query_count", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "startup_query_interval", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "querier_timeout", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "robustness_variable", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "maximum_multicast_entries", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "reserved_multicast_entries", "0"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update IGMP Interface Policy Name with UUID") },
				Config:    testAccMSOTenantPoliciesIGMPInterfacePolicyConfigUpdateName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "name", "test_igmp_interface_policy_renamed"),
					resource.TestCheckResourceAttr("mso_tenant_policies_igmp_interface_policy.igmp_policy", "description", "Renamed Policy"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import IGMP Interface Policy") },
				ResourceName:      "mso_tenant_policies_igmp_interface_policy.igmp_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithArguments("mso_tenant_policies_igmp_interface_policy", "igmpInterface"),
	})
}

func testAccMSOTenantPoliciesIGMPInterfacePolicyConfigCreate() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_igmp_interface_policy" "igmp_policy" {
        template_id                 = mso_template.template_tenant.id
        name                        = "test_igmp_interface_policy"
        description                 = "Test IGMP Interface Policy"
        version3_asm                = true
        fast_leave                  = true
        report_link_local_groups    = true
        igmp_version                = "v3"
        group_timeout               = 300
        query_interval              = 125
        query_response_interval     = 10
        last_member_count           = 2
        last_member_response_time   = 1
        startup_query_count         = 2
        startup_query_interval      = 31
        querier_timeout             = 255
        robustness_variable         = 2
        maximum_multicast_entries   = 1000000
        reserved_multicast_entries  = 100000
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesIGMPInterfacePolicyConfigUpdateV2() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_igmp_interface_policy" "igmp_policy" {
        template_id                 = mso_template.template_tenant.id
        name                        = "test_igmp_interface_policy"
        description                 = "Updated to IGMP v2"
        version3_asm                = false
        fast_leave                  = false
        report_link_local_groups    = false
        igmp_version                = "v2"
        group_timeout               = 260
        query_interval              = 100
        query_response_interval     = 10
        last_member_count           = 2
        last_member_response_time   = 1
        startup_query_count         = 2
        startup_query_interval      = 31
        querier_timeout             = 255
        robustness_variable         = 2
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesIGMPInterfacePolicyConfigUpdateWithRouteMaps() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_route_map_policy_multicast" "state_limit" {
		template_id = mso_template.template_tenant.id
		name        = "tf_test_state_limit"
		description = "Terraform test Route Map Policy for Multicast"
		route_map_multicast_entries {
			order                   = 1
			group_ip                = "226.2.2.2/8"
			source_ip               = "1.1.1.1/1"
			rendezvous_point_ip     = "1.1.1.2"
			action                  = "permit"
		}
	}

    resource "mso_tenant_policies_route_map_policy_multicast" "report_policy" {
		template_id = mso_tenant_policies_route_map_policy_multicast.state_limit.template_id
		name        = "tf_test_report_policy"
		description = "Terraform test Route Map Policy for Multicast"
		route_map_multicast_entries {
			order                   = 1
			group_ip                = "226.2.2.2/8"
			source_ip               = "1.1.1.1/1"
			rendezvous_point_ip     = "1.1.1.2"
			action                  = "permit"
		}
	}

	resource "mso_tenant_policies_route_map_policy_multicast" "static_report" {
		template_id = mso_tenant_policies_route_map_policy_multicast.report_policy.template_id
		name        = "tf_test_static_report"
		description = "Terraform test Route Map Policy for Multicast"
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
        name                           = "test_igmp_interface_policy"
        description                    = "With Route Maps"
        igmp_version                   = "v3"
        state_limit_route_map_uuid     = mso_tenant_policies_route_map_policy_multicast.state_limit.uuid
        report_policy_route_map_uuid   = mso_tenant_policies_route_map_policy_multicast.report_policy.uuid
        static_report_route_map_uuid   = mso_tenant_policies_route_map_policy_multicast.static_report.uuid
        maximum_multicast_entries      = 5000000
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesIGMPInterfacePolicyConfigUpdateRemoveRouteMaps() string {
	return fmt.Sprintf(`%s
	resource "mso_tenant_policies_route_map_policy_multicast" "state_limit" {
		template_id = mso_template.template_tenant.id
		name        = "tf_test_state_limit"
		description = "Terraform test Route Map Policy for Multicast"
		route_map_multicast_entries {
			order                   = 1
			group_ip                = "226.2.2.2/8"
			source_ip               = "1.1.1.1/1"
			rendezvous_point_ip     = "1.1.1.2"
			action                  = "permit"
		}
	}

    resource "mso_tenant_policies_route_map_policy_multicast" "report_policy" {
		template_id = mso_tenant_policies_route_map_policy_multicast.state_limit.template_id
		name        = "tf_test_report_policy"
		description = "Terraform test Route Map Policy for Multicast"
		route_map_multicast_entries {
			order                   = 1
			group_ip                = "226.2.2.2/8"
			source_ip               = "1.1.1.1/1"
			rendezvous_point_ip     = "1.1.1.2"
			action                  = "permit"
		}
	}

	resource "mso_tenant_policies_route_map_policy_multicast" "static_report" {
		template_id = mso_tenant_policies_route_map_policy_multicast.report_policy.template_id
		name        = "tf_test_static_report"
		description = "Terraform test Route Map Policy for Multicast"
		route_map_multicast_entries {
			order                   = 1
			group_ip                = "226.2.2.2/8"
			source_ip               = "1.1.1.1/1"
			rendezvous_point_ip     = "1.1.1.2"
			action                  = "permit"
		}
	}
    resource "mso_tenant_policies_igmp_interface_policy" "igmp_policy" {
        template_id                 = mso_template.template_tenant.id
        name                        = "test_igmp_interface_policy"
        description                 = "Route Maps Removed"
        igmp_version                = "v3"
        state_limit_route_map_uuid  = ""
        report_policy_route_map_uuid = ""
        static_report_route_map_uuid = ""
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesIGMPInterfacePolicyConfigUpdateTimers() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_igmp_interface_policy" "igmp_policy" {
        template_id                 = mso_template.template_tenant.id
        name                        = "test_igmp_interface_policy"
        description                 = "Updated Timers"
        igmp_version                = "v3"
        group_timeout               = 500
        query_interval              = 200
        query_response_interval     = 15
        last_member_count           = 3
        last_member_response_time   = 2
        startup_query_count         = 5
        startup_query_interval      = 50
        querier_timeout             = 300
        robustness_variable         = 3
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesIGMPInterfacePolicyConfigUpdateMaxValues() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_igmp_interface_policy" "igmp_policy" {
        template_id                 = mso_template.template_tenant.id
        name                        = "test_igmp_interface_policy_max"
        description                 = "Maximum Values Test"
        version3_asm                = true
        fast_leave                  = true
        report_link_local_groups    = true
        igmp_version                = "v3"
        group_timeout               = 65535
        query_interval              = 18000
        query_response_interval     = 25
        last_member_count           = 5
        last_member_response_time   = 25
        startup_query_count         = 10
        startup_query_interval      = 18000
        querier_timeout             = 65535
        robustness_variable         = 7
        maximum_multicast_entries   = 4294967295
        reserved_multicast_entries  = 4294967295
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesIGMPInterfacePolicyConfigUpdateMinValues() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_igmp_interface_policy" "igmp_policy" {
        template_id                 = mso_template.template_tenant.id
        name                        = "test_igmp_interface_policy_min"
        description                 = "Minimum Values Test"
        version3_asm                = false
        fast_leave                  = false
        report_link_local_groups    = false
        igmp_version                = "v2"
        group_timeout               = 3
        query_interval              = 2
        query_response_interval     = 1
        last_member_count           = 1
        last_member_response_time   = 1
        startup_query_count         = 1
        startup_query_interval      = 1
        querier_timeout             = 1
        robustness_variable         = 1
        maximum_multicast_entries   = 1
        reserved_multicast_entries  = 0
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesIGMPInterfacePolicyConfigUpdateName() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_igmp_interface_policy" "igmp_policy" {
        template_id                 = mso_template.template_tenant.id
        name                        = "test_igmp_interface_policy_renamed"
        description                 = "Renamed Policy"
        igmp_version                = "v3"
        group_timeout               = 300
    }`, testAccMSOTemplateResourceTenantConfig())
}
