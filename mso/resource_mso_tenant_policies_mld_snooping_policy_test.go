package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOTenantPoliciesMLDSnoopingPolicyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create MLD Snooping Policy") },
				Config:    testAccMSOTenantPoliciesMLDSnoopingPolicyConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "name", "test_mld_snooping_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "description", "Test MLD Snooping Policy"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_mld_snooping_policy.mld_policy", "uuid"),
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "admin_state", "enabled"),
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "fast_leave_control", "true"),
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "querier_control", "true"),
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "querier_version", "v2"),
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "query_interval", "125"),
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "query_response_interval", "10"),
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "last_member_query_interval", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "start_query_interval", "31"),
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "start_query_count", "2"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update MLD Snooping Policy - Change to v1") },
				Config:    testAccMSOTenantPoliciesMLDSnoopingPolicyConfigUpdateV1(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "name", "test_mld_snooping_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "description", "Updated to MLD v1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "admin_state", "enabled"),
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "fast_leave_control", "false"),
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "querier_control", "true"),
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "querier_version", "v1"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update MLD Snooping Policy - Disable Admin State") },
				Config:    testAccMSOTenantPoliciesMLDSnoopingPolicyConfigDisabled(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "name", "test_mld_snooping_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "admin_state", "disabled"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update MLD Snooping Policy - Maximum Values") },
				Config:    testAccMSOTenantPoliciesMLDSnoopingPolicyConfigMaxValues(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "name", "test_mld_snooping_policy_max"),
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "query_interval", "18000"),
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "query_response_interval", "25"),
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "last_member_query_interval", "25"),
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "start_query_interval", "18000"),
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "start_query_count", "10"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update MLD Snooping Policy - Minimum Values") },
				Config:    testAccMSOTenantPoliciesMLDSnoopingPolicyConfigMinValues(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "name", "test_mld_snooping_policy_min"),
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "query_interval", "2"),
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "query_response_interval", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "last_member_query_interval", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "start_query_interval", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_mld_snooping_policy.mld_policy", "start_query_count", "1"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import MLD Snooping Policy") },
				ResourceName:      "mso_tenant_policies_mld_snooping_policy.mld_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithArguments("mso_tenant_policies_mld_snooping_policy", "mldSnoop"),
	})
}

func testAccMSOTenantPoliciesMLDSnoopingPolicyConfigCreate() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_mld_snooping_policy" "mld_policy" {
        template_id                = mso_template.template_tenant.id
        name                       = "test_mld_snooping_policy"
        description                = "Test MLD Snooping Policy"
        admin_state                = "enabled"
        fast_leave_control         = true
        querier_control            = true
        querier_version            = "v2"
        query_interval             = 125
        query_response_interval    = 10
        last_member_query_interval = 1
        start_query_interval       = 31
        start_query_count          = 2
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMLDSnoopingPolicyConfigUpdateV1() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_mld_snooping_policy" "mld_policy" {
        template_id                = mso_template.template_tenant.id
        name                       = "test_mld_snooping_policy"
        description                = "Updated to MLD v1"
        admin_state                = "enabled"
        fast_leave_control         = false
        querier_control            = true
        querier_version            = "v1"
        query_interval             = 100
        query_response_interval    = 10
        last_member_query_interval = 1
        start_query_interval       = 31
        start_query_count          = 2
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMLDSnoopingPolicyConfigDisabled() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_mld_snooping_policy" "mld_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_mld_snooping_policy"
        description = "MLD Snooping Disabled"
        admin_state = "disabled"
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMLDSnoopingPolicyConfigMaxValues() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_mld_snooping_policy" "mld_policy" {
        template_id                = mso_template.template_tenant.id
        name                       = "test_mld_snooping_policy_max"
        description                = "Maximum Values Test"
        admin_state                = "enabled"
        fast_leave_control         = true
        querier_control            = true
        querier_version            = "v2"
        query_interval             = 18000
        query_response_interval    = 25
        last_member_query_interval = 25
        start_query_interval       = 18000
        start_query_count          = 10
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMLDSnoopingPolicyConfigMinValues() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_mld_snooping_policy" "mld_policy" {
        template_id                = mso_template.template_tenant.id
        name                       = "test_mld_snooping_policy_min"
        description                = "Minimum Values Test"
        admin_state                = "enabled"
        fast_leave_control         = false
        querier_control            = true
        querier_version            = "v1"
        query_interval             = 2
        query_response_interval    = 1
        last_member_query_interval = 1
        start_query_interval       = 1
        start_query_count          = 1
    }`, testAccMSOTemplateResourceTenantConfig())
}
