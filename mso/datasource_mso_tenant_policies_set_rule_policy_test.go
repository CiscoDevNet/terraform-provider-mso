package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMSOTenantPoliciesSetRulePolicyDataSource(t *testing.T) {
	dataSourceName := "data.mso_tenant_policies_set_rule_policy.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVersionCheck(t, "5.1")
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments("mso_tenant_policies_set_rule_policy", "tenantPolicyTemplate", "template", "setRulePolicies"),
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy Data Source - Non-Existing Policy (error)") },
				Config:      testAccMSOSetRulePolicyDataSourceMissing(),
				ExpectError: regexp.MustCompile(`Policy name non_existing_set_rule_policy not found`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Set Rule Policy Data Source - Read Policy") },
				Config:    testAccMSOSetRulePolicyDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "description", "Test Set Rule Policy"),
					resource.TestCheckResourceAttrSet(dataSourceName, "uuid"),
					resource.TestCheckResourceAttr(dataSourceName, "set_route_tag.0.value", "500"),
					resource.TestCheckResourceAttr(dataSourceName, "set_weight.0.value", "50"),
					resource.TestCheckResourceAttr(dataSourceName, "set_preference.0.value", "200"),
					resource.TestCheckResourceAttr(dataSourceName, "set_metric", "100"),
					resource.TestCheckResourceAttr(dataSourceName, "set_next_hop", "2001:db8::1"),
					resource.TestCheckResourceAttr(dataSourceName, "set_metric_type", "ospf-type2"),
					resource.TestCheckResourceAttr(dataSourceName, "set_next_hop_propagate", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "set_multi_path", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "set_community.0.community", "no-export"),
					resource.TestCheckResourceAttr(dataSourceName, "set_community.0.criteria", "replace"),
					resource.TestCheckResourceAttr(dataSourceName, "set_additional_community.#", "2"),
					testCheckTypeSetStringElemAttr(dataSourceName, "set_additional_community", "no-advertise"),
					testCheckTypeSetStringElemAttr(dataSourceName, "set_additional_community", "extended:color:35"),
					resource.TestCheckResourceAttr(dataSourceName, "set_dampening.0.half_life", "15"),
					resource.TestCheckResourceAttr(dataSourceName, "set_dampening.0.reuse_limit", "750"),
					resource.TestCheckResourceAttr(dataSourceName, "set_dampening.0.suppress_limit", "2000"),
					resource.TestCheckResourceAttr(dataSourceName, "set_dampening.0.max_suppress_time", "60"),
					resource.TestCheckResourceAttr(dataSourceName, "set_as_path.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "set_as_path.0.prepend.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "set_as_path.0.prepend_last_as.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "set_as_path.0.prepend.0.path_asn.#", "4"),
					CustomTestCheckTypeSetElemAttrs(dataSourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "10", "order": "1"}),
					CustomTestCheckTypeSetElemAttrs(dataSourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "20", "order": "3"}),
					CustomTestCheckTypeSetElemAttrs(dataSourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "30", "order": "2"}),
					CustomTestCheckTypeSetElemAttrs(dataSourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "30", "order": "4"}),
					resource.TestCheckResourceAttr(dataSourceName, "set_as_path.0.prepend_last_as.0.asn_count", "10"),
					resource.TestCheckResourceAttrPair(dataSourceName, "uuid", "mso_tenant_policies_set_rule_policy.test", "uuid"),
					resource.TestCheckResourceAttrPair(dataSourceName, "id", "mso_tenant_policies_set_rule_policy.test", "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "template_id", "mso_tenant_policies_set_rule_policy.test", "template_id"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "test_set_rule_policy"),
				),
			},
		},
	})
}

func testAccMSOSetRulePolicyDataSourceMissing() string {
	return fmt.Sprintf(`%s
data "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "non_existing_set_rule_policy"
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyDataSourceConfig() string {
	return fmt.Sprintf(`%s
data "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_tenant_policies_set_rule_policy.test.template_id
  name        = mso_tenant_policies_set_rule_policy.test.name
}`, testAccMSOSetRulePolicyConfigFull("test_set_rule_policy"))
}
