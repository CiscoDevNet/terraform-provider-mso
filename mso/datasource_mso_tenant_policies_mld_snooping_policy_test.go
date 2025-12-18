package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOTenantPoliciesMLDSnoopingPolicyDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: MLD Snooping Policy Data Source") },
				Config:    testAccMSOTenantPoliciesMLDSnoopingPolicyDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_tenant_policies_mld_snooping_policy.mld_policy", "name", "test_mld_snooping_policy"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_mld_snooping_policy.mld_policy", "description", "Test MLD Snooping Policy"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_mld_snooping_policy.mld_policy", "admin_state", "enabled"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_mld_snooping_policy.mld_policy", "fast_leave_control", "true"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_mld_snooping_policy.mld_policy", "querier_control", "true"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_mld_snooping_policy.mld_policy", "querier_version", "v2"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_mld_snooping_policy.mld_policy", "query_interval", "125"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_mld_snooping_policy.mld_policy", "query_response_interval", "10"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_mld_snooping_policy.mld_policy", "last_member_query_interval", "1"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_mld_snooping_policy.mld_policy", "start_query_interval", "31"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_mld_snooping_policy.mld_policy", "start_query_count", "2"),
				),
			},
		},
	})
}

func testAccMSOTenantPoliciesMLDSnoopingPolicyDataSource() string {
	return fmt.Sprintf(`%s
    data "mso_tenant_policies_mld_snooping_policy" "mld_policy" {
        template_id = mso_tenant_policies_mld_snooping_policy.mld_policy.template_id
        name        = "test_mld_snooping_policy"
    }`, testAccMSOTenantPoliciesMLDSnoopingPolicyConfigCreate())
}
