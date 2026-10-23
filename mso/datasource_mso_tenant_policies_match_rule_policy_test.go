package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMSOTenantPoliciesMatchRulePolicyDataSource(t *testing.T) {
	const dataSourceName = "data.mso_tenant_policies_match_rule_policy.test_policy"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVersionCheck(t, "5.1")
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments("mso_tenant_policies_match_rule_policy", "tenantPolicyTemplate", "template", "matchRulePolicies"),
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Match Rule Policy Data Source - Non-Existing Policy (error)") },
				Config:      testAccMSOTenantPoliciesMatchRulePolicyDataSourceMissing(),
				ExpectError: regexp.MustCompile(`(?s)Policy name missing_match_rule_policy not found`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Match Rule Policy Data Source - Read Policy") },
				Config:    testAccMSOTenantPoliciesMatchRulePolicyDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "description", "Full match policy"),
					resource.TestCheckResourceAttrSet(dataSourceName, "uuid"),
					resource.TestCheckResourceAttr(dataSourceName, "match_prefix.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "match_prefix.*", map[string]string{"prefix": "10.10.10.1/24", "description": "IPv4 aggregate", "aggregate": "true", "from_prefix": "25", "to_prefix": "32"}),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "match_prefix.*", map[string]string{"prefix": "2001:db8::1/64", "description": "", "aggregate": "true", "from_prefix": "65", "to_prefix": "128"}),
					resource.TestCheckResourceAttr(dataSourceName, "match_community_term.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "match_community_term.*", map[string]string{"name": "test_term", "description": "Community term", "match_community_factor.#": "2"}),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "no-export", "scope": "non-transitive", "description": "Export factor"}),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "no-advertise", "scope": "transitive", "description": "Advertise factor"}),
					resource.TestCheckResourceAttrPair(dataSourceName, "uuid", "mso_tenant_policies_match_rule_policy.test_policy", "uuid"),
					resource.TestCheckResourceAttrPair(dataSourceName, "id", "mso_tenant_policies_match_rule_policy.test_policy", "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "template_id", "mso_tenant_policies_match_rule_policy.test_policy", "template_id"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "test_match_rule_policy"),
				),
			},
		},
	})
}

func testAccMSOTenantPoliciesMatchRulePolicyDataSourceMissing() string {
	return fmt.Sprintf(`%s
data "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "missing_match_rule_policy"
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMatchRulePolicyDataSource() string {
	return fmt.Sprintf(`%s
data "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_tenant_policies_match_rule_policy.test_policy.template_id
  name        = mso_tenant_policies_match_rule_policy.test_policy.name
}`, testAccMSOTenantPoliciesMatchRulePolicyConfigFull())
}
