package mso

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMSOTenantPoliciesMatchRulePolicyResource(t *testing.T) {
	const resourceName = "mso_tenant_policies_match_rule_policy.test_policy"
	var uuid, templateID string
	// Save IDs from state for later API calls and UUID comparisons across test steps.
	capture := func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok || rs.Primary == nil {
			return fmt.Errorf("resource %s not found in state", resourceName)
		}
		uuid = rs.Primary.Attributes["uuid"]
		templateID = rs.Primary.Attributes["template_id"]
		if uuid == "" || templateID == "" {
			return fmt.Errorf("resource %s is missing its template ID or UUID", resourceName)
		}
		return nil
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVersionCheck(t, "5.1")
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments("mso_tenant_policies_match_rule_policy", "tenantPolicyTemplate", "template", "matchRulePolicies"),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create Match Rule Policy - Minimal Configuration") },
				Config:    testAccMSOTenantPoliciesMatchRulePolicyConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_match_rule_policy"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "match_prefix.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "match_community_term.#", "0"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Match Rule Policy - Configure All Settings") },
				Config:    testAccMSOTenantPoliciesMatchRulePolicyConfigFull(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_match_rule_policy"),
					resource.TestCheckResourceAttr(resourceName, "description", "Full match policy"),
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "match_prefix.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "10.10.10.1/24", "description": "IPv4 aggregate", "aggregate": "true", "from_prefix": "25", "to_prefix": "32"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "2001:db8::1/64", "description": "", "aggregate": "true", "from_prefix": "65", "to_prefix": "128"}),
					resource.TestCheckResourceAttr(resourceName, "match_community_term.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*", map[string]string{"name": "test_term", "description": "Community term", "match_community_factor.#": "2"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "no-export", "scope": "non-transitive", "description": "Export factor"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "no-advertise", "scope": "transitive", "description": "Advertise factor"}),
					capture,
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Match Rule Policy - Default Prefix and Community Settings") },
				Config:    testAccMSOTenantPoliciesMatchRulePolicyConfigDefaults(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_match_rule_policy"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "match_prefix.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "10.10.10.1/24", "description": "", "aggregate": "false", "from_prefix": "0", "to_prefix": "0"}),
					resource.TestCheckResourceAttr(resourceName, "match_community_term.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*", map[string]string{"name": "test_term", "description": "", "match_community_factor.#": "1"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "no-export", "scope": "transitive", "description": ""}),
					func(state *terraform.State) error {
						rs := state.RootModule().Resources[resourceName]
						if rs.Primary.Attributes["uuid"] != uuid || rs.Primary.ID != fmt.Sprintf("templateId/%s/MatchRulePolicy/%s", templateID, rs.Primary.Attributes["name"]) {
							return fmt.Errorf("updating the Match Rule Policy changed its ID or UUID")
						}
						return nil
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Match Rule Policy - Equal Aggregation Endpoints") },
				Config:    testAccMSOTenantPoliciesMatchRulePolicyConfigEqualAggregation(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_match_rule_policy"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "match_prefix.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "10.10.10.1/24", "description": "", "aggregate": "true", "from_prefix": "25", "to_prefix": "25"}),
					resource.TestCheckResourceAttr(resourceName, "match_community_term.#", "0"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Match Rule Policy - Zero Aggregation Endpoints") },
				Config:    testAccMSOTenantPoliciesMatchRulePolicyConfigZeroAggregation(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_match_rule_policy"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "match_prefix.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "10.10.10.1/24", "description": "", "aggregate": "true", "from_prefix": "0", "to_prefix": "0"}),
					resource.TestCheckResourceAttr(resourceName, "match_community_term.#", "0"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Match Rule Policy - Zero Minimum Aggregation Length") },
				Config:    testAccMSOTenantPoliciesMatchRulePolicyConfigZeroMinimumAggregation(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_match_rule_policy"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "match_prefix.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "10.10.10.1/24", "description": "", "aggregate": "true", "from_prefix": "0", "to_prefix": "25"}),
					resource.TestCheckResourceAttr(resourceName, "match_community_term.#", "0"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Match Rule Policy - Add Prefixes, Community Terms and Factors") },
				Config:    testAccMSOTenantPoliciesMatchRulePolicyConfigAddEntries(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_match_rule_policy"),
					resource.TestCheckResourceAttr(resourceName, "description", "Full match policy"),
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "match_prefix.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "10.10.10.1/24", "description": "IPv4 aggregate", "aggregate": "true", "from_prefix": "25", "to_prefix": "32"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "2001:db8::1/64", "description": "", "aggregate": "true", "from_prefix": "65", "to_prefix": "128"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "172.16.1.0/24", "description": "", "aggregate": "false", "from_prefix": "0", "to_prefix": "0"}),
					resource.TestCheckResourceAttr(resourceName, "match_community_term.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*", map[string]string{"name": "test_term", "description": "Community term", "match_community_factor.#": "2"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*", map[string]string{"name": "second_term", "description": "", "match_community_factor.#": "1"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*", map[string]string{"name": "empty_term", "description": "", "match_community_factor.#": "0"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "no-export", "scope": "non-transitive", "description": "Export factor"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "no-advertise", "scope": "transitive", "description": "Advertise factor"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "regular:as2-nn2:4:15", "scope": "transitive", "description": ""}),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update Match Rule Policy - Update Existing Prefix, Community Term and Factor")
				},
				Config: testAccMSOTenantPoliciesMatchRulePolicyConfigUpdateEntries(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_match_rule_policy"),
					resource.TestCheckResourceAttr(resourceName, "description", "Full match policy"),
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "match_prefix.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "10.10.10.1/24", "description": "Updated prefix", "aggregate": "true", "from_prefix": "26", "to_prefix": "31"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "2001:db8::1/64", "description": "", "aggregate": "true", "from_prefix": "65", "to_prefix": "128"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "172.16.1.0/24", "description": "", "aggregate": "false", "from_prefix": "0", "to_prefix": "0"}),
					resource.TestCheckResourceAttr(resourceName, "match_community_term.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*", map[string]string{"name": "test_term", "description": "Updated term", "match_community_factor.#": "2"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*", map[string]string{"name": "second_term", "description": "", "match_community_factor.#": "1"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*", map[string]string{"name": "empty_term", "description": "", "match_community_factor.#": "0"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "no-export", "scope": "non-transitive", "description": "Updated factor"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "no-advertise", "scope": "transitive", "description": "Advertise factor"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "regular:as2-nn2:4:15", "scope": "transitive", "description": ""}),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update Match Rule Policy - Remove Selected Prefixes, Community Terms and Factors")
				},
				Config: testAccMSOTenantPoliciesMatchRulePolicyConfigRemoveEntries(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_match_rule_policy"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "match_prefix.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "2001:db8::1/64", "description": "", "aggregate": "false", "from_prefix": "0", "to_prefix": "0"}),
					resource.TestCheckResourceAttr(resourceName, "match_community_term.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*", map[string]string{"name": "test_term", "description": "", "match_community_factor.#": "1"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "no-advertise", "scope": "transitive", "description": ""}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Match Rule Policy - Restore Prefixes, Community Terms and Factors") },
				Config:    testAccMSOTenantPoliciesMatchRulePolicyConfigUpdateEntries(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_match_rule_policy"),
					resource.TestCheckResourceAttr(resourceName, "description", "Full match policy"),
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "match_prefix.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "10.10.10.1/24", "description": "Updated prefix", "aggregate": "true", "from_prefix": "26", "to_prefix": "31"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "2001:db8::1/64", "description": "", "aggregate": "true", "from_prefix": "65", "to_prefix": "128"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "172.16.1.0/24", "description": "", "aggregate": "false", "from_prefix": "0", "to_prefix": "0"}),
					resource.TestCheckResourceAttr(resourceName, "match_community_term.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*", map[string]string{"name": "test_term", "description": "Updated term", "match_community_factor.#": "2"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*", map[string]string{"name": "second_term", "description": "", "match_community_factor.#": "1"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*", map[string]string{"name": "empty_term", "description": "", "match_community_factor.#": "0"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "no-export", "scope": "non-transitive", "description": "Updated factor"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "no-advertise", "scope": "transitive", "description": "Advertise factor"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "regular:as2-nn2:4:15", "scope": "transitive", "description": ""}),
				),
			},
			{
				PreConfig:          func() { fmt.Println("Test: Plan Match Rule Policy - Reorder Prefixes, Community Terms and Factors") },
				Config:             testAccMSOTenantPoliciesMatchRulePolicyConfigReordered(),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Match Rule Policy - Remove All Factors from a Community Term") },
				Config:    testAccMSOTenantPoliciesMatchRulePolicyConfigWithoutFactors(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_match_rule_policy"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "match_prefix.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "match_community_term.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*", map[string]string{"name": "test_term", "description": "", "match_community_factor.#": "0"}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Match Rule Policy - Remove All Prefix and Community Collections") },
				Config:    testAccMSOTenantPoliciesMatchRulePolicyConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_match_rule_policy"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "match_prefix.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "match_community_term.#", "0"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Match Rule Policy - Restore All Settings") },
				Config:    testAccMSOTenantPoliciesMatchRulePolicyConfigFull(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_match_rule_policy"),
					resource.TestCheckResourceAttr(resourceName, "description", "Full match policy"),
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "match_prefix.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "10.10.10.1/24", "description": "IPv4 aggregate", "aggregate": "true", "from_prefix": "25", "to_prefix": "32"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "2001:db8::1/64", "description": "", "aggregate": "true", "from_prefix": "65", "to_prefix": "128"}),
					resource.TestCheckResourceAttr(resourceName, "match_community_term.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*", map[string]string{"name": "test_term", "description": "Community term", "match_community_factor.#": "2"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "no-export", "scope": "non-transitive", "description": "Export factor"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "no-advertise", "scope": "transitive", "description": "Advertise factor"}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Match Rule Policy - Clear Attributes with Explicit Null Values") },
				Config:    testAccMSOTenantPoliciesMatchRulePolicyConfigNullValues(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_match_rule_policy"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "match_prefix.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "10.10.10.1/24", "description": "", "aggregate": "false", "from_prefix": "0", "to_prefix": "0"}),
					resource.TestCheckResourceAttr(resourceName, "match_community_term.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*", map[string]string{"name": "test_term", "description": "", "match_community_factor.#": "1"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "no-export", "scope": "transitive", "description": ""}),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update Match Rule Policy - Restore All Settings after Explicit Null Values")
				},
				Config: testAccMSOTenantPoliciesMatchRulePolicyConfigFull(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_match_rule_policy"),
					resource.TestCheckResourceAttr(resourceName, "description", "Full match policy"),
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "match_prefix.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "10.10.10.1/24", "description": "IPv4 aggregate", "aggregate": "true", "from_prefix": "25", "to_prefix": "32"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "2001:db8::1/64", "description": "", "aggregate": "true", "from_prefix": "65", "to_prefix": "128"}),
					resource.TestCheckResourceAttr(resourceName, "match_community_term.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*", map[string]string{"name": "test_term", "description": "Community term", "match_community_factor.#": "2"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "no-export", "scope": "non-transitive", "description": "Export factor"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "no-advertise", "scope": "transitive", "description": "Advertise factor"}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Rename Match Rule Policy while retaining UUID") },
				Config:    testAccMSOTenantPoliciesMatchRulePolicyConfigRename(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_match_rule_policy_renamed"),
					resource.TestCheckResourceAttr(resourceName, "description", "Full match policy"),
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "match_prefix.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "10.10.10.1/24", "description": "IPv4 aggregate", "aggregate": "true", "from_prefix": "25", "to_prefix": "32"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "2001:db8::1/64", "description": "", "aggregate": "true", "from_prefix": "65", "to_prefix": "128"}),
					resource.TestCheckResourceAttr(resourceName, "match_community_term.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*", map[string]string{"name": "test_term", "description": "Community term", "match_community_factor.#": "2"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "no-export", "scope": "non-transitive", "description": "Export factor"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "no-advertise", "scope": "transitive", "description": "Advertise factor"}),
					func(state *terraform.State) error {
						if state.RootModule().Resources[resourceName].Primary.Attributes["uuid"] != uuid || state.RootModule().Resources[resourceName].Primary.ID != fmt.Sprintf("templateId/%s/MatchRulePolicy/test_match_rule_policy_renamed", templateID) {
							return fmt.Errorf("rename changed Match Rule Policy ID or UUID")
						}
						response, err := testAccPreCheck(t).GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
						if err != nil {
							return err
						}
						if _, err := GetPolicyByName(response, "test_match_rule_policy", "tenantPolicyTemplate", "template", "matchRulePolicies"); !isPolicyNotFound(err) {
							if err != nil {
								return err
							}
							return fmt.Errorf("original Match Rule Policy name still exists after rename")
						}
						return nil
					},
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import Match Rule Policy") },
				ResourceName: resourceName, ImportState: true, ImportStateVerify: true,
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import missing Match Rule Policy") },
				ResourceName: resourceName, ImportState: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return fmt.Sprintf("templateId/%s/MatchRulePolicy/missing_match_rule_policy", state.RootModule().Resources[resourceName].Primary.Attributes["template_id"]), nil
				},
				ExpectError: regexp.MustCompile(`(?s)Match Rule Policy .*missing_match_rule_policy\s+not\s+found`),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update Match Rule Policy after Out-of-Band Changes")
					c := testAccPreCheck(t)
					// Resolve the current API array index before applying out-of-band changes.
					response, err := c.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
					if err != nil {
						t.Fatal(err)
					}
					index, err := GetPolicyIndexByKeyAndValue(response, "uuid", uuid, "tenantPolicyTemplate", "template", "matchRulePolicies")
					if err != nil {
						t.Fatal(err)
					}
					payload := container.New()
					payload.Array()
					path := fmt.Sprintf("/tenantPolicyTemplate/template/matchRulePolicies/%d", index)
					if err := addPatchPayloadToContainer(payload, "remove", fmt.Sprintf("%s/description", path), nil); err != nil {
						t.Fatal(err)
					}
					if err := addPatchPayloadToContainer(payload, "remove", fmt.Sprintf("%s/matchPrefixList", path), nil); err != nil {
						t.Fatal(err)
					}
					if err := doPatchRequest(c, fmt.Sprintf("api/v1/templates/%s", templateID), payload); err != nil {
						t.Fatal(err)
					}
				},
				Config: testAccMSOTenantPoliciesMatchRulePolicyConfigRename(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_match_rule_policy_renamed"),
					resource.TestCheckResourceAttr(resourceName, "description", "Full match policy"),
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "match_prefix.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "10.10.10.1/24", "description": "IPv4 aggregate", "aggregate": "true", "from_prefix": "25", "to_prefix": "32"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "2001:db8::1/64", "description": "", "aggregate": "true", "from_prefix": "65", "to_prefix": "128"}),
					resource.TestCheckResourceAttr(resourceName, "match_community_term.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*", map[string]string{"name": "test_term", "description": "Community term", "match_community_factor.#": "2"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "no-export", "scope": "non-transitive", "description": "Export factor"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "no-advertise", "scope": "transitive", "description": "Advertise factor"}),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate Match Rule Policy after Out-of-Band Deletion")
					c := testAccPreCheck(t)
					// Resolve the current API array index before deleting the policy outside Terraform.
					response, err := c.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
					if err != nil {
						t.Fatal(err)
					}
					index, err := GetPolicyIndexByKeyAndValue(response, "uuid", uuid, "tenantPolicyTemplate", "template", "matchRulePolicies")
					if err != nil {
						t.Fatal(err)
					}
					policy := response.S("tenantPolicyTemplate", "template", "matchRulePolicies").Index(index)
					name, ok := policy.S("name").Data().(string)
					if !ok {
						t.Fatal("policy to delete has no valid name")
					}
					path := fmt.Sprintf("/tenantPolicyTemplate/template/matchRulePolicies/%d", index)
					if _, err := c.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateID), models.GetRemovePatchPayload(path)); err != nil {
						t.Fatal(err)
					}
					response, err = c.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
					if err != nil {
						t.Fatal(err)
					}
					if _, err := GetPolicyByName(response, name, "tenantPolicyTemplate", "template", "matchRulePolicies"); !isPolicyNotFound(err) {
						if err != nil {
							t.Fatalf("verify out-of-band deletion: %s", err)
						}
						t.Fatalf("policy %q still exists after out-of-band deletion", name)
					}
				},
				Config: testAccMSOTenantPoliciesMatchRulePolicyConfigRename(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_match_rule_policy_renamed"),
					resource.TestCheckResourceAttr(resourceName, "description", "Full match policy"),
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "match_prefix.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "10.10.10.1/24", "description": "IPv4 aggregate", "aggregate": "true", "from_prefix": "25", "to_prefix": "32"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "2001:db8::1/64", "description": "", "aggregate": "true", "from_prefix": "65", "to_prefix": "128"}),
					resource.TestCheckResourceAttr(resourceName, "match_community_term.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*", map[string]string{"name": "test_term", "description": "Community term", "match_community_factor.#": "2"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "no-export", "scope": "non-transitive", "description": "Export factor"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "no-advertise", "scope": "transitive", "description": "Advertise factor"}),
					func(state *terraform.State) error {
						rs := state.RootModule().Resources[resourceName]
						if rs.Primary.Attributes["uuid"] == uuid {
							return fmt.Errorf("out-of-band deletion did not recreate Match Rule Policy")
						}
						return capture(state)
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create destination template for Match Rule Policy replacement") },
				Config:    testAccMSOTenantPoliciesMatchRulePolicyConfigTemplateChangeSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_match_rule_policy"),
					resource.TestCheckResourceAttr(resourceName, "description", "Full match policy"),
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "match_prefix.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "10.10.10.1/24", "description": "IPv4 aggregate", "aggregate": "true", "from_prefix": "25", "to_prefix": "32"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "2001:db8::1/64", "description": "", "aggregate": "true", "from_prefix": "65", "to_prefix": "128"}),
					resource.TestCheckResourceAttr(resourceName, "match_community_term.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*", map[string]string{"name": "test_term", "description": "Community term", "match_community_factor.#": "2"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "no-export", "scope": "non-transitive", "description": "Export factor"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "no-advertise", "scope": "transitive", "description": "Advertise factor"}),
					capture,
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Replace Match Rule Policy in another template") },
				Config:    testAccMSOTenantPoliciesMatchRulePolicyConfigTemplateChangeDestination(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_match_rule_policy"),
					resource.TestCheckResourceAttr(resourceName, "description", "Full match policy"),
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant_replacement", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "match_prefix.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "10.10.10.1/24", "description": "IPv4 aggregate", "aggregate": "true", "from_prefix": "25", "to_prefix": "32"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_prefix.*", map[string]string{"prefix": "2001:db8::1/64", "description": "", "aggregate": "true", "from_prefix": "65", "to_prefix": "128"}),
					resource.TestCheckResourceAttr(resourceName, "match_community_term.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*", map[string]string{"name": "test_term", "description": "Community term", "match_community_factor.#": "2"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "no-export", "scope": "non-transitive", "description": "Export factor"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "match_community_term.*.match_community_factor.*", map[string]string{"community": "no-advertise", "scope": "transitive", "description": "Advertise factor"}),
					func(state *terraform.State) error {
						if state.RootModule().Resources[resourceName].Primary.Attributes["uuid"] == uuid {
							return fmt.Errorf("template replacement retained Match Rule Policy UUID")
						}
						response, err := testAccPreCheck(t).GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
						if err != nil {
							return err
						}
						_, err = GetPolicyIndexByKeyAndValue(response, "uuid", uuid, "tenantPolicyTemplate", "template", "matchRulePolicies")
						if !isPolicyNotFound(err) {
							if err != nil {
								return err
							}
							return fmt.Errorf("original Match Rule Policy still exists after template replacement")
						}
						return nil
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create Multiple Match Rule Policies") },
				Config:    testAccMSOTenantPoliciesMatchRulePolicyConfigMultiple(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_match_rule_policy.first", "description", "First"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_match_rule_policy.first", "uuid"),
					resource.TestCheckResourceAttr("mso_tenant_policies_match_rule_policy.middle", "description", "Middle"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_match_rule_policy.middle", "uuid"),
					resource.TestCheckResourceAttr("mso_tenant_policies_match_rule_policy.last", "description", "Last"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_match_rule_policy.last", "uuid"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Multiple Match Rule Policies - Remove Middle Policy") },
				Config:    testAccMSOTenantPoliciesMatchRulePolicyConfigMultipleRemoveMiddle(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_match_rule_policy.first", "description", "First"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_match_rule_policy.first", "uuid"),
					resource.TestCheckResourceAttr("mso_tenant_policies_match_rule_policy.last", "description", "Last"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_match_rule_policy.last", "uuid"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Multiple Match Rule Policies - After Array Index Shift") },
				Config:    testAccMSOTenantPoliciesMatchRulePolicyConfigMultipleUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_match_rule_policy.first", "description", "Updated first"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_match_rule_policy.first", "uuid"),
					resource.TestCheckResourceAttr("mso_tenant_policies_match_rule_policy.last", "description", "Updated last"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_match_rule_policy.last", "uuid"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Multiple Match Rule Policies - Add Policy after Array Index Shift") },
				Config:    testAccMSOTenantPoliciesMatchRulePolicyConfigMultipleAdd(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_match_rule_policy.first", "description", "Updated first again"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_match_rule_policy.first", "uuid"),
					resource.TestCheckResourceAttr("mso_tenant_policies_match_rule_policy.last", "description", "Updated last again"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_match_rule_policy.last", "uuid"),
					resource.TestCheckResourceAttr("mso_tenant_policies_match_rule_policy.another", "description", "Another"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_match_rule_policy.another", "uuid"),
				),
			},
		},
	})
}

func TestAccMSOTenantPoliciesMatchRulePolicyResourceErrors(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVersionCheck(t, "5.1")
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments("mso_tenant_policies_match_rule_policy", "tenantPolicyTemplate", "template", "matchRulePolicies"),
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Reject Match Rule Policy without template ID") },
				Config:      testAccMSOTenantPoliciesMatchRulePolicyConfigErrorWithoutTemplateID(),
				ExpectError: regexp.MustCompile(`The argument "template_id" is required`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Match Rule Policy without name") },
				Config:      testAccMSOTenantPoliciesMatchRulePolicyConfigErrorWithoutName(),
				ExpectError: regexp.MustCompile(`The argument "name" is required`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Match Rule Policy with empty name") },
				Config:      testAccMSOTenantPoliciesMatchRulePolicyConfigErrorName(""),
				ExpectError: regexp.MustCompile(`expected length of name`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Match Rule Policy with long name") },
				Config:      testAccMSOTenantPoliciesMatchRulePolicyConfigErrorName(strings.Repeat("a", 65)),
				ExpectError: regexp.MustCompile(`expected length of name`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Match Rule Policy without prefix") },
				Config:      testAccMSOTenantPoliciesMatchRulePolicyConfigErrorWithoutPrefix(),
				ExpectError: regexp.MustCompile(`The argument "prefix" is required`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Match Rule Policy with invalid prefix") },
				Config:      testAccMSOTenantPoliciesMatchRulePolicyConfigErrorPrefix("invalid", false, 0, 0),
				ExpectError: regexp.MustCompile(`(?s)expected .* to be a valid CIDR`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Match Rule Policy with negative from prefix") },
				Config:      testAccMSOTenantPoliciesMatchRulePolicyConfigErrorPrefix("10.10.10.1/24", true, -1, 32),
				ExpectError: regexp.MustCompile(`(?s)expected .*from_prefix.*to\s+be\s+in\s+the\s+range`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Match Rule Policy with excessive from prefix") },
				Config:      testAccMSOTenantPoliciesMatchRulePolicyConfigErrorPrefix("2001:db8::/64", true, 129, 128),
				ExpectError: regexp.MustCompile(`(?s)expected .*from_prefix.*to\s+be\s+in\s+the\s+range`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Match Rule Policy with negative to prefix") },
				Config:      testAccMSOTenantPoliciesMatchRulePolicyConfigErrorPrefix("10.10.10.1/24", true, 25, -1),
				ExpectError: regexp.MustCompile(`(?s)expected .*to_prefix.*to\s+be\s+in\s+the\s+range`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Match Rule Policy with excessive to prefix") },
				Config:      testAccMSOTenantPoliciesMatchRulePolicyConfigErrorPrefix("2001:db8::/64", true, 65, 129),
				ExpectError: regexp.MustCompile(`(?s)expected .*to_prefix.*to\s+be\s+in\s+the\s+range`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Match Rule Policy with excessive IPv4 lengths") },
				Config:      testAccMSOTenantPoliciesMatchRulePolicyConfigErrorPrefix("10.10.10.1/24", true, 25, 33),
				ExpectError: regexp.MustCompile(`prefix lengths must be between 0 and 32`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Match Rule Policy with lengths but no aggregation") },
				Config:      testAccMSOTenantPoliciesMatchRulePolicyConfigErrorPrefix("10.10.10.1/24", false, 25, 32),
				ExpectError: regexp.MustCompile(`from_prefix and to_prefix must be 0 when aggregate is false`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Match Rule Policy with aggregation below mask") },
				Config:      testAccMSOTenantPoliciesMatchRulePolicyConfigErrorPrefix("10.10.10.1/24", true, 24, 32),
				ExpectError: regexp.MustCompile(`(?s)requires from_prefix to be 0 or greater than the prefix mask`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Match Rule Policy with reversed aggregation range") },
				Config:      testAccMSOTenantPoliciesMatchRulePolicyConfigErrorPrefix("10.10.10.1/24", true, 31, 30),
				ExpectError: regexp.MustCompile(`(?s)requires from_prefix to be 0 or greater than the prefix mask`),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Reject Match Rule Policy with nonzero minimum and zero maximum aggregation lengths")
				},
				Config:      testAccMSOTenantPoliciesMatchRulePolicyConfigErrorPrefix("10.10.10.1/24", true, 25, 0),
				ExpectError: regexp.MustCompile(`(?s)requires from_prefix to be 0 or greater than the prefix mask`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Match Rule Policy with duplicate prefix") },
				Config:      testAccMSOTenantPoliciesMatchRulePolicyConfigErrorDuplicatePrefix(),
				ExpectError: regexp.MustCompile(`match_prefix contains duplicate prefix`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Match Rule Policy without term name") },
				Config:      testAccMSOTenantPoliciesMatchRulePolicyConfigErrorWithoutTermName(),
				ExpectError: regexp.MustCompile(`The argument "name" is required`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Match Rule Policy with empty term name") },
				Config:      testAccMSOTenantPoliciesMatchRulePolicyConfigErrorTermName(""),
				ExpectError: regexp.MustCompile(`expected length of match_community_term`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Match Rule Policy with long term name") },
				Config:      testAccMSOTenantPoliciesMatchRulePolicyConfigErrorTermName(strings.Repeat("a", 65)),
				ExpectError: regexp.MustCompile(`expected length of match_community_term`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Match Rule Policy with duplicate term name") },
				Config:      testAccMSOTenantPoliciesMatchRulePolicyConfigErrorDuplicateTerm(),
				ExpectError: regexp.MustCompile(`match_community_term contains duplicate name`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Match Rule Policy with duplicate community") },
				Config:      testAccMSOTenantPoliciesMatchRulePolicyConfigErrorDuplicateFactor(),
				ExpectError: regexp.MustCompile(`contains duplicate community`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Match Rule Policy without community") },
				Config:      testAccMSOTenantPoliciesMatchRulePolicyConfigErrorWithoutCommunityFactor(),
				ExpectError: regexp.MustCompile(`The argument "community" is required`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Match Rule Policy with empty community") },
				Config:      testAccMSOTenantPoliciesMatchRulePolicyConfigErrorFactor("", "transitive"),
				ExpectError: regexp.MustCompile(`(?s)expected .*community.*not\s+be\s+an\s+empty\s+string`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Match Rule Policy with invalid community scope") },
				Config:      testAccMSOTenantPoliciesMatchRulePolicyConfigErrorFactor("no-export", "invalid"),
				ExpectError: regexp.MustCompile(`(?s)expected .*scope.*to\s+be\s+one\s+of`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create Match Rule Policy after configuration errors") },
				Config:    testAccMSOTenantPoliciesMatchRulePolicyConfigCreate(),
				Check:     resource.TestCheckResourceAttrSet("mso_tenant_policies_match_rule_policy.test_policy", "uuid"),
			},
		},
	})
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigCreate() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_match_rule_policy"

}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigFull() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_match_rule_policy"
  description = "Full match policy"
  match_prefix {
    prefix      = "10.10.10.1/24"
    description = "IPv4 aggregate"
    aggregate   = true
    from_prefix = 25
    to_prefix   = 32
  }
  match_prefix {
    prefix      = "2001:db8::1/64"
    aggregate   = true
    from_prefix = 65
    to_prefix   = 128
  }
  match_community_term {
    name        = "test_term"
    description = "Community term"
    match_community_factor {
      community   = "no-export"
      scope       = "non-transitive"
      description = "Export factor"
    }
    match_community_factor {
      community   = "no-advertise"
      scope       = "transitive"
      description = "Advertise factor"
    }
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigDefaults() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_match_rule_policy"
  match_prefix {
    prefix = "10.10.10.1/24"
  }
  match_community_term {
    name = "test_term"
    match_community_factor {
      community   = "no-export"
    }
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigAddEntries() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_match_rule_policy"
  description = "Full match policy"
  match_prefix {
    prefix      = "10.10.10.1/24"
    description = "IPv4 aggregate"
    aggregate   = true
    from_prefix = 25
    to_prefix   = 32
  }
  match_prefix {
    prefix      = "2001:db8::1/64"
    aggregate   = true
    from_prefix = 65
    to_prefix   = 128
  }
  match_community_term {
    name        = "test_term"
    description = "Community term"
    match_community_factor {
      community   = "no-export"
      scope       = "non-transitive"
      description = "Export factor"
    }
    match_community_factor {
      community   = "no-advertise"
      scope       = "transitive"
      description = "Advertise factor"
    }
  }
  match_prefix {
    prefix = "172.16.1.0/24"
  }
  match_community_term {
    name = "second_term"
    match_community_factor {
      community   = "regular:as2-nn2:4:15"
    }
  }
  match_community_term {
    name = "empty_term"
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigUpdateEntries() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_match_rule_policy"
  description = "Full match policy"
  match_prefix {
    prefix      = "10.10.10.1/24"
    description = "Updated prefix"
    aggregate   = true
    from_prefix = 26
    to_prefix   = 31
  }
  match_prefix {
    prefix      = "2001:db8::1/64"
    aggregate   = true
    from_prefix = 65
    to_prefix   = 128
  }
  match_community_term {
    name        = "test_term"
    description = "Updated term"
    match_community_factor {
      community   = "no-export"
      scope       = "non-transitive"
      description = "Updated factor"
    }
    match_community_factor {
      community   = "no-advertise"
      scope       = "transitive"
      description = "Advertise factor"
    }
  }
  match_prefix {
    prefix = "172.16.1.0/24"
  }
  match_community_term {
    name = "second_term"
    match_community_factor {
      community   = "regular:as2-nn2:4:15"
    }
  }
  match_community_term {
    name = "empty_term"
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigRemoveEntries() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_match_rule_policy"
  match_prefix {
    prefix = "2001:db8::1/64"
  }
  match_community_term {
    name = "test_term"
    match_community_factor {
      community   = "no-advertise"
    }
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigWithoutFactors() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_match_rule_policy"
  match_community_term {
    name = "test_term"
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigNullValues() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_match_rule_policy"
  description = null
  match_prefix {
    prefix      = "10.10.10.1/24"
    description = null
    aggregate   = null
    from_prefix = null
    to_prefix   = null
  }
  match_community_term {
    name        = "test_term"
    description = null
    match_community_factor {
      community   = "no-export"
      scope       = null
      description = null
    }
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigReordered() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_match_rule_policy"
  description = "Full match policy"
  match_community_term {
    name = "empty_term"
  }
  match_community_term {
    name = "second_term"
    match_community_factor {
      community   = "regular:as2-nn2:4:15"
    }
  }
  match_community_term {
    name        = "test_term"
    description = "Updated term"
    match_community_factor {
      community   = "no-advertise"
      scope       = "transitive"
      description = "Advertise factor"
    }
    match_community_factor {
      community   = "no-export"
      scope       = "non-transitive"
      description = "Updated factor"
    }
  }
  match_prefix {
    prefix = "172.16.1.0/24"
  }
  match_prefix {
    prefix      = "2001:db8::1/64"
    aggregate   = true
    from_prefix = 65
    to_prefix   = 128
  }
  match_prefix {
    prefix      = "10.10.10.1/24"
    description = "Updated prefix"
    aggregate   = true
    from_prefix = 26
    to_prefix   = 31
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigRename() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_match_rule_policy_renamed"
  description = "Full match policy"
  match_prefix {
    prefix      = "10.10.10.1/24"
    description = "IPv4 aggregate"
    aggregate   = true
    from_prefix = 25
    to_prefix   = 32
  }
  match_prefix {
    prefix      = "2001:db8::1/64"
    aggregate   = true
    from_prefix = 65
    to_prefix   = 128
  }
  match_community_term {
    name        = "test_term"
    description = "Community term"
    match_community_factor {
      community   = "no-export"
      scope       = "non-transitive"
      description = "Export factor"
    }
    match_community_factor {
      community   = "no-advertise"
      scope       = "transitive"
      description = "Advertise factor"
    }
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigTemplateChangeSource() string {
	return fmt.Sprintf(`%s
resource "mso_template" "template_tenant_replacement" {
  template_name = "test_template_tenant_match_replacement"
  template_type = "tenant"
  tenant_id     = mso_tenant.%s.id
  sites         = [data.mso_site.%s.id]
}

resource "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_match_rule_policy"
  description = "Full match policy"
  match_prefix {
    prefix      = "10.10.10.1/24"
    description = "IPv4 aggregate"
    aggregate   = true
    from_prefix = 25
    to_prefix   = 32
  }
  match_prefix {
    prefix      = "2001:db8::1/64"
    aggregate   = true
    from_prefix = 65
    to_prefix   = 128
  }
  match_community_term {
    name        = "test_term"
    description = "Community term"
    match_community_factor {
      community   = "no-export"
      scope       = "non-transitive"
      description = "Export factor"
    }
    match_community_factor {
      community   = "no-advertise"
      scope       = "transitive"
      description = "Advertise factor"
    }
  }
}`, testAccMSOTemplateResourceTenantConfig(), msoTemplateTenantName, msoTemplateSiteName1)
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigTemplateChangeDestination() string {
	return fmt.Sprintf(`%s
resource "mso_template" "template_tenant_replacement" {
  template_name = "test_template_tenant_match_replacement"
  template_type = "tenant"
  tenant_id     = mso_tenant.%s.id
  sites         = [data.mso_site.%s.id]
}

resource "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_template.template_tenant_replacement.id
  name        = "test_match_rule_policy"
  description = "Full match policy"
  match_prefix {
    prefix      = "10.10.10.1/24"
    description = "IPv4 aggregate"
    aggregate   = true
    from_prefix = 25
    to_prefix   = 32
  }
  match_prefix {
    prefix      = "2001:db8::1/64"
    aggregate   = true
    from_prefix = 65
    to_prefix   = 128
  }
  match_community_term {
    name        = "test_term"
    description = "Community term"
    match_community_factor {
      community   = "no-export"
      scope       = "non-transitive"
      description = "Export factor"
    }
    match_community_factor {
      community   = "no-advertise"
      scope       = "transitive"
      description = "Advertise factor"
    }
  }
}`, testAccMSOTemplateResourceTenantConfig(), msoTemplateTenantName, msoTemplateSiteName1)
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigMultiple() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "first" {
  template_id = mso_template.template_tenant.id
  name        = "first_match_rule_policy"
  description = "First"
}
resource "mso_tenant_policies_match_rule_policy" "middle" {
  depends_on  = [mso_tenant_policies_match_rule_policy.first]
  template_id = mso_template.template_tenant.id
  name        = "middle_match_rule_policy"
  description = "Middle"
}
resource "mso_tenant_policies_match_rule_policy" "last" {
  depends_on  = [mso_tenant_policies_match_rule_policy.middle]
  template_id = mso_template.template_tenant.id
  name        = "last_match_rule_policy"
  description = "Last"
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigMultipleRemoveMiddle() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "first" {
  template_id = mso_template.template_tenant.id
  name        = "first_match_rule_policy"
  description = "First"
}
resource "mso_tenant_policies_match_rule_policy" "last" {
  depends_on  = [mso_tenant_policies_match_rule_policy.first]
  template_id = mso_template.template_tenant.id
  name        = "last_match_rule_policy"
  description = "Last"
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigMultipleUpdate() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "first" {
  template_id = mso_template.template_tenant.id
  name        = "first_match_rule_policy"
  description = "Updated first"
}
resource "mso_tenant_policies_match_rule_policy" "last" {
  depends_on  = [mso_tenant_policies_match_rule_policy.first]
  template_id = mso_template.template_tenant.id
  name        = "last_match_rule_policy"
  description = "Updated last"
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigMultipleAdd() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "first" {
  template_id = mso_template.template_tenant.id
  name        = "first_match_rule_policy"
  description = "Updated first again"
}
resource "mso_tenant_policies_match_rule_policy" "last" {
  depends_on  = [mso_tenant_policies_match_rule_policy.first]
  template_id = mso_template.template_tenant.id
  name        = "last_match_rule_policy"
  description = "Updated last again"
}
resource "mso_tenant_policies_match_rule_policy" "another" {
  depends_on  = [mso_tenant_policies_match_rule_policy.last]
  template_id = mso_template.template_tenant.id
  name        = "another_match_rule_policy"
  description = "Another"
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigErrorWithoutTemplateID() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "test_policy" {
  name = "test_match_rule_policy"

}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigErrorWithoutName() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_template.template_tenant.id

}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigErrorName(name string) string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = %q

}`, testAccMSOTemplateResourceTenantConfig(), name)
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigErrorWithoutPrefix() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_match_rule_policy"
  match_prefix {
    description = "Missing prefix"
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigErrorPrefix(prefix string, aggregate bool, from, to int) string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_match_rule_policy"
  match_prefix {
    prefix      = %q
    aggregate   = %t
    from_prefix = %d
    to_prefix   = %d
  }
}`, testAccMSOTemplateResourceTenantConfig(), prefix, aggregate, from, to)
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigErrorDuplicatePrefix() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_match_rule_policy"
  match_prefix {
    prefix = "10.10.10.1/24"
  }
  match_prefix {
    prefix      = "10.10.10.1/24"
    description = "Duplicate"
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigErrorWithoutTermName() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_match_rule_policy"
  match_community_term {
    description = "Missing name"
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigErrorTermName(name string) string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_match_rule_policy"
  match_community_term {
    name = %q
  }
}`, testAccMSOTemplateResourceTenantConfig(), name)
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigErrorDuplicateTerm() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_match_rule_policy"
  match_community_term {
    name = "duplicate"
  }
  match_community_term {
    name        = "duplicate"
    description = "Duplicate"
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigErrorWithoutCommunityFactor() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_match_rule_policy"
  match_community_term {
    name = "test_term"
    match_community_factor {
      description = "Missing factor"
    }
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigErrorFactor(factor, scope string) string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_match_rule_policy"
  match_community_term {
    name = "test_term"
    match_community_factor {
      community   = %q
      scope       = %q
    }
  }
}`, testAccMSOTemplateResourceTenantConfig(), factor, scope)
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigEqualAggregation() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_match_rule_policy"
  match_prefix {
    prefix      = "10.10.10.1/24"
    aggregate   = true
    from_prefix = 25
    to_prefix   = 25
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigZeroAggregation() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_match_rule_policy"
  match_prefix {
    prefix      = "10.10.10.1/24"
    aggregate   = true
    from_prefix = 0
    to_prefix   = 0
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigErrorDuplicateFactor() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_match_rule_policy"
  match_community_term {
    name = "test_term"
    match_community_factor {
      community   = "no-export"
      scope       = "transitive"
    }
    match_community_factor {
      community   = "no-export"
      scope       = "non-transitive"
    }
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMatchRulePolicyConfigZeroMinimumAggregation() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_match_rule_policy"
  match_prefix {
    prefix      = "10.10.10.1/24"
    aggregate   = true
    from_prefix = 0
    to_prefix   = 25
  }
}`, testAccMSOTemplateResourceTenantConfig())
}
