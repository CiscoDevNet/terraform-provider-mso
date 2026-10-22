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

func TestAccMSOTenantPoliciesSetRulePolicyResource(t *testing.T) {
	resourceName := "mso_tenant_policies_set_rule_policy.test"
	var templateID, uuid, originalTemplateID, originalUUID string
	// Save IDs from state for later API calls and UUID comparisons across test steps.
	capture := func(s *terraform.State) error {
		state, ok := s.RootModule().Resources[resourceName]
		if !ok || state.Primary == nil {
			return fmt.Errorf("resource %s not found in state", resourceName)
		}
		templateID = state.Primary.Attributes["template_id"]
		uuid = state.Primary.Attributes["uuid"]
		if templateID == "" || uuid == "" {
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
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments("mso_tenant_policies_set_rule_policy", "tenantPolicyTemplate", "template", "setRulePolicies"),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create Set Rule Policy - Minimal Configuration") },
				Config:    testAccMSOSetRulePolicyConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_set_rule_policy"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "set_route_tag.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "set_next_hop_propagate", "false"),
					resource.TestCheckResourceAttr(resourceName, "set_multi_path", "false"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.#", "0"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create Set Rule Policy - Zero Numeric Actions and Literal Description") },
				Config:    testAccMSOSetRulePolicyConfigCreateZeros(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_set_rule_policy.zero", "description", "{}"),
					resource.TestCheckResourceAttr("mso_tenant_policies_set_rule_policy.zero", "set_route_tag.0.value", "0"),
					resource.TestCheckResourceAttr("mso_tenant_policies_set_rule_policy.zero", "set_weight.0.value", "0"),
					resource.TestCheckResourceAttr("mso_tenant_policies_set_rule_policy.zero", "set_preference.0.value", "0"),
					resource.TestCheckResourceAttr("mso_tenant_policies_set_rule_policy.zero", "set_metric", "1"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Set Rule Policy - Configure All Actions") },
				Config:    testAccMSOSetRulePolicyConfigFull("test_set_rule_policy"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Test Set Rule Policy"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "set_route_tag.0.value", "500"),
					resource.TestCheckResourceAttr(resourceName, "set_weight.0.value", "50"),
					resource.TestCheckResourceAttr(resourceName, "set_preference.0.value", "200"),
					resource.TestCheckResourceAttr(resourceName, "set_metric", "100"),
					resource.TestCheckResourceAttr(resourceName, "set_next_hop", "2001:db8::1"),
					resource.TestCheckResourceAttr(resourceName, "set_metric_type", "ospf-type2"),
					resource.TestCheckResourceAttr(resourceName, "set_next_hop_propagate", "false"),
					resource.TestCheckResourceAttr(resourceName, "set_multi_path", "false"),
					resource.TestCheckResourceAttr(resourceName, "set_community.0.community", "no-export"),
					resource.TestCheckResourceAttr(resourceName, "set_community.0.criteria", "replace"),
					resource.TestCheckResourceAttr(resourceName, "set_additional_community.#", "2"),
					testCheckTypeSetStringElemAttr(resourceName, "set_additional_community", "no-advertise"),
					testCheckTypeSetStringElemAttr(resourceName, "set_additional_community", "extended:color:35"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.half_life", "15"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.reuse_limit", "750"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.suppress_limit", "2000"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.max_suppress_time", "60"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend_last_as.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend.0.path_asn.#", "4"),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "10", "order": "1"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "20", "order": "3"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "30", "order": "2"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "30", "order": "4"}),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend_last_as.0.asn_count", "10"),
					capture,
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Refresh Set Rule Policy after Out-of-Band AS Path Reordering")
					msoClient := testAccPreCheck(t)
					// Read the API arrays so this step can reorder them outside Terraform.
					response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
					if err != nil {
						t.Fatal(err)
					}
					index, err := GetPolicyIndexByKeyAndValue(response, "uuid", uuid, "tenantPolicyTemplate", "template", "setRulePolicies")
					if err != nil {
						t.Fatal(err)
					}
					policy := response.S("tenantPolicyTemplate", "template", "setRulePolicies").Index(index)
					actions, ok := policy.S("setAsPath").Data().([]interface{})
					if !ok {
						t.Fatal("AS-path actions are not a list")
					}
					for _, action := range actions {
						settings, ok := action.(map[string]interface{})
						if !ok {
							t.Fatal("AS-path action is not an object")
						}
						if settings["criteria"] == "prepend" {
							paths, ok := settings["pathASNs"].([]interface{})
							if !ok {
								t.Fatal("AS-path ASN entries are not a list")
							}
							for left, right := 0, len(paths)-1; left < right; left, right = left+1, right-1 {
								paths[left], paths[right] = paths[right], paths[left]
							}
						}
					}
					for left, right := 0, len(actions)-1; left < right; left, right = left+1, right-1 {
						actions[left], actions[right] = actions[right], actions[left]
					}
					payload := container.New()
					payload.Array()
					path := fmt.Sprintf("/tenantPolicyTemplate/template/setRulePolicies/%d/setAsPath", index)
					if err := addPatchPayloadToContainer(payload, "replace", path, actions); err != nil {
						t.Fatal(err)
					}
					if err := doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateID), payload); err != nil {
						t.Fatal(err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Test Set Rule Policy"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "set_route_tag.0.value", "500"),
					resource.TestCheckResourceAttr(resourceName, "set_weight.0.value", "50"),
					resource.TestCheckResourceAttr(resourceName, "set_preference.0.value", "200"),
					resource.TestCheckResourceAttr(resourceName, "set_metric", "100"),
					resource.TestCheckResourceAttr(resourceName, "set_next_hop", "2001:db8::1"),
					resource.TestCheckResourceAttr(resourceName, "set_metric_type", "ospf-type2"),
					resource.TestCheckResourceAttr(resourceName, "set_next_hop_propagate", "false"),
					resource.TestCheckResourceAttr(resourceName, "set_multi_path", "false"),
					resource.TestCheckResourceAttr(resourceName, "set_community.0.community", "no-export"),
					resource.TestCheckResourceAttr(resourceName, "set_community.0.criteria", "replace"),
					resource.TestCheckResourceAttr(resourceName, "set_additional_community.#", "2"),
					testCheckTypeSetStringElemAttr(resourceName, "set_additional_community", "no-advertise"),
					testCheckTypeSetStringElemAttr(resourceName, "set_additional_community", "extended:color:35"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.half_life", "15"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.reuse_limit", "750"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.suppress_limit", "2000"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.max_suppress_time", "60"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend_last_as.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend.0.path_asn.#", "4"),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "10", "order": "1"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "20", "order": "3"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "30", "order": "2"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "30", "order": "4"}),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend_last_as.0.asn_count", "10"),
				),
			},
			{
				PreConfig:          func() { fmt.Println("Test: Plan Set Rule Policy - Reorder AS Path Blocks and ASN Entries") },
				Config:             testAccMSOSetRulePolicyConfigASPathReordered(),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Set Rule Policy - Update One AS Path Action") },
				Config:    testAccMSOSetRulePolicyConfigASPathUpdated(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "set_as_path.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend_last_as.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend.0.path_asn.#", "4"),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "10", "order": "1"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "20", "order": "3"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "30", "order": "2"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "30", "order": "4"}),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend_last_as.0.asn_count", "3"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Set Rule Policy - Remove Selected ASN and Community Entries") },
				Config:    testAccMSOSetRulePolicyConfigPartialCollections(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Test Set Rule Policy"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "set_route_tag.0.value", "500"),
					resource.TestCheckResourceAttr(resourceName, "set_weight.0.value", "50"),
					resource.TestCheckResourceAttr(resourceName, "set_preference.0.value", "200"),
					resource.TestCheckResourceAttr(resourceName, "set_metric", "100"),
					resource.TestCheckResourceAttr(resourceName, "set_next_hop", "2001:db8::1"),
					resource.TestCheckResourceAttr(resourceName, "set_metric_type", "ospf-type2"),
					resource.TestCheckResourceAttr(resourceName, "set_next_hop_propagate", "false"),
					resource.TestCheckResourceAttr(resourceName, "set_multi_path", "false"),
					resource.TestCheckResourceAttr(resourceName, "set_community.0.community", "no-export"),
					resource.TestCheckResourceAttr(resourceName, "set_community.0.criteria", "replace"),
					resource.TestCheckResourceAttr(resourceName, "set_additional_community.#", "1"),
					testCheckTypeSetStringElemAttr(resourceName, "set_additional_community", "extended:color:35"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.half_life", "15"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.reuse_limit", "750"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.suppress_limit", "2000"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.max_suppress_time", "60"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend_last_as.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend.0.path_asn.#", "2"),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "10", "order": "1"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "30", "order": "2"}),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend_last_as.0.asn_count", "3"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Set Rule Policy - Restore Selected ASN and Community Entries") },
				Config:    testAccMSOSetRulePolicyConfigASPathUpdated(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Test Set Rule Policy"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "set_route_tag.0.value", "500"),
					resource.TestCheckResourceAttr(resourceName, "set_weight.0.value", "50"),
					resource.TestCheckResourceAttr(resourceName, "set_preference.0.value", "200"),
					resource.TestCheckResourceAttr(resourceName, "set_metric", "100"),
					resource.TestCheckResourceAttr(resourceName, "set_next_hop", "2001:db8::1"),
					resource.TestCheckResourceAttr(resourceName, "set_metric_type", "ospf-type2"),
					resource.TestCheckResourceAttr(resourceName, "set_next_hop_propagate", "false"),
					resource.TestCheckResourceAttr(resourceName, "set_multi_path", "false"),
					resource.TestCheckResourceAttr(resourceName, "set_community.0.community", "no-export"),
					resource.TestCheckResourceAttr(resourceName, "set_community.0.criteria", "replace"),
					resource.TestCheckResourceAttr(resourceName, "set_additional_community.#", "2"),
					testCheckTypeSetStringElemAttr(resourceName, "set_additional_community", "no-advertise"),
					testCheckTypeSetStringElemAttr(resourceName, "set_additional_community", "extended:color:35"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.half_life", "15"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.reuse_limit", "750"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.suppress_limit", "2000"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.max_suppress_time", "60"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend_last_as.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend.0.path_asn.#", "4"),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "10", "order": "1"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "20", "order": "3"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "30", "order": "2"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "30", "order": "4"}),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend_last_as.0.asn_count", "3"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Set Rule Policy - Remove One AS Path Action") },
				Config:    testAccMSOSetRulePolicyConfigASPathLastOnly(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "set_as_path.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend_last_as.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend_last_as.0.asn_count", "3"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Set Rule Policy - Empty Prepend AS Path Action") },
				Config:    testAccMSOSetRulePolicyConfigASPathEmptyPrepend(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "set_as_path.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend_last_as.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend.0.path_asn.#", "0"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Set Rule Policy - Restore Both AS Path Actions") },
				Config:    testAccMSOSetRulePolicyConfigFull("test_set_rule_policy"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Test Set Rule Policy"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "set_route_tag.0.value", "500"),
					resource.TestCheckResourceAttr(resourceName, "set_weight.0.value", "50"),
					resource.TestCheckResourceAttr(resourceName, "set_preference.0.value", "200"),
					resource.TestCheckResourceAttr(resourceName, "set_metric", "100"),
					resource.TestCheckResourceAttr(resourceName, "set_next_hop", "2001:db8::1"),
					resource.TestCheckResourceAttr(resourceName, "set_metric_type", "ospf-type2"),
					resource.TestCheckResourceAttr(resourceName, "set_next_hop_propagate", "false"),
					resource.TestCheckResourceAttr(resourceName, "set_multi_path", "false"),
					resource.TestCheckResourceAttr(resourceName, "set_community.0.community", "no-export"),
					resource.TestCheckResourceAttr(resourceName, "set_community.0.criteria", "replace"),
					resource.TestCheckResourceAttr(resourceName, "set_additional_community.#", "2"),
					testCheckTypeSetStringElemAttr(resourceName, "set_additional_community", "no-advertise"),
					testCheckTypeSetStringElemAttr(resourceName, "set_additional_community", "extended:color:35"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.half_life", "15"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.reuse_limit", "750"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.suppress_limit", "2000"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.max_suppress_time", "60"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend_last_as.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend.0.path_asn.#", "4"),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "10", "order": "1"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "20", "order": "3"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "30", "order": "2"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "30", "order": "4"}),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend_last_as.0.asn_count", "10"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Set Rule Policy - Clear Attributes with Explicit Null Values") },
				Config:    testAccMSOSetRulePolicyConfigNullValues(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_set_rule_policy"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "set_route_tag.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "set_weight.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "set_preference.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "set_metric", "0"),
					resource.TestCheckResourceAttr(resourceName, "set_next_hop", ""),
					resource.TestCheckResourceAttr(resourceName, "set_metric_type", ""),
					resource.TestCheckResourceAttr(resourceName, "set_next_hop_propagate", "false"),
					resource.TestCheckResourceAttr(resourceName, "set_multi_path", "false"),
					resource.TestCheckResourceAttr(resourceName, "set_community.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "set_additional_community.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.#", "0"),
					func(s *terraform.State) error {
						if s.RootModule().Resources[resourceName].Primary.Attributes["uuid"] != uuid {
							return fmt.Errorf("clearing attributes with explicit null values changed the policy UUID")
						}
						return nil
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Set Rule Policy - Restore All Actions after Explicit Null Values") },
				Config:    testAccMSOSetRulePolicyConfigFull("test_set_rule_policy"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Test Set Rule Policy"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "set_route_tag.0.value", "500"),
					resource.TestCheckResourceAttr(resourceName, "set_weight.0.value", "50"),
					resource.TestCheckResourceAttr(resourceName, "set_preference.0.value", "200"),
					resource.TestCheckResourceAttr(resourceName, "set_metric", "100"),
					resource.TestCheckResourceAttr(resourceName, "set_next_hop", "2001:db8::1"),
					resource.TestCheckResourceAttr(resourceName, "set_metric_type", "ospf-type2"),
					resource.TestCheckResourceAttr(resourceName, "set_next_hop_propagate", "false"),
					resource.TestCheckResourceAttr(resourceName, "set_multi_path", "false"),
					resource.TestCheckResourceAttr(resourceName, "set_community.0.community", "no-export"),
					resource.TestCheckResourceAttr(resourceName, "set_community.0.criteria", "replace"),
					resource.TestCheckResourceAttr(resourceName, "set_additional_community.#", "2"),
					testCheckTypeSetStringElemAttr(resourceName, "set_additional_community", "no-advertise"),
					testCheckTypeSetStringElemAttr(resourceName, "set_additional_community", "extended:color:35"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.half_life", "15"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.reuse_limit", "750"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.suppress_limit", "2000"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.max_suppress_time", "60"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend_last_as.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend.0.path_asn.#", "4"),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "10", "order": "1"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "20", "order": "3"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "30", "order": "2"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "30", "order": "4"}),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend_last_as.0.asn_count", "10"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Set Rule Policy - Rename") },
				Config:    testAccMSOSetRulePolicyConfigFull("test_set_rule_policy_renamed"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_set_rule_policy_renamed"),
					func(s *terraform.State) error {
						state := s.RootModule().Resources[resourceName]
						if state.Primary.Attributes["uuid"] != uuid {
							return fmt.Errorf("renaming the policy changed its UUID")
						}
						if state.Primary.ID != fmt.Sprintf("templateId/%s/SetRulePolicy/test_set_rule_policy_renamed", templateID) {
							return fmt.Errorf("policy ID did not reflect its new name")
						}
						return nil
					},
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Set Rule Policy") },
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import Non-Existing Set Rule Policy (error)") },
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					return fmt.Sprintf("templateId/%s/SetRulePolicy/non_existing_set_rule_policy", templateID), nil
				},
				ExpectError: regexp.MustCompile(`(?s)Set Rule Policy.*non_existing_set_rule_policy.*not found`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Set Rule Policy - Maximum Values and Explicit AS Path Only") },
				Config:    testAccMSOSetRulePolicyConfigMaxValues(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "set_route_tag.0.value", "999999999"),
					resource.TestCheckResourceAttr(resourceName, "set_weight.0.value", "9999"),
					resource.TestCheckResourceAttr(resourceName, "set_preference.0.value", "999999999"),
					resource.TestCheckResourceAttr(resourceName, "set_metric", "999999999"),
					resource.TestCheckResourceAttr(resourceName, "set_metric_type", "none"),
					resource.TestCheckResourceAttr(resourceName, "set_community.0.criteria", "none"),
					resource.TestCheckResourceAttr(resourceName, "set_community.0.community", ""),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.half_life", "60"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.reuse_limit", "19999"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.suppress_limit", "20000"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.max_suppress_time", "255"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend_last_as.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend.0.path_asn.#", "1"),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "4294967295", "order": "31"}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Set Rule Policy - Zero and Minimum Values") },
				Config:    testAccMSOSetRulePolicyConfigMinValues(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_set_rule_policy"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "set_route_tag.0.value", "0"),
					resource.TestCheckResourceAttr(resourceName, "set_weight.0.value", "0"),
					resource.TestCheckResourceAttr(resourceName, "set_preference.0.value", "0"),
					resource.TestCheckResourceAttr(resourceName, "set_metric", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_next_hop", "192.0.2.1"),
					resource.TestCheckResourceAttr(resourceName, "set_metric_type", "ospf-type1"),
					resource.TestCheckResourceAttr(resourceName, "set_community.0.community", "no-advertise"),
					resource.TestCheckResourceAttr(resourceName, "set_community.0.criteria", "append"),
					resource.TestCheckResourceAttr(resourceName, "set_additional_community.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.half_life", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.reuse_limit", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.suppress_limit", "2"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.max_suppress_time", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend_last_as.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend.0.path_asn.#", "1"),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "1", "order": "0"}),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend_last_as.0.asn_count", "1"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Set Rule Policy - Remove All Actions") },
				Config:    testAccMSOSetRulePolicyConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "set_route_tag.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "set_weight.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "set_preference.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "set_metric", "0"),
					resource.TestCheckResourceAttr(resourceName, "set_next_hop", ""),
					resource.TestCheckResourceAttr(resourceName, "set_metric_type", ""),
					resource.TestCheckResourceAttr(resourceName, "set_community.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "set_additional_community.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.#", "0"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Set Rule Policy - Enable Next-Hop Propagation and Multipath") },
				Config:    testAccMSOSetRulePolicyConfigPropagation(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "set_next_hop_propagate", "true"),
					resource.TestCheckResourceAttr(resourceName, "set_multi_path", "true"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Set Rule Policy - Disable Next-Hop Propagation and Multipath") },
				Config:    testAccMSOSetRulePolicyConfigDisabled(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "set_next_hop_propagate", "false"),
					resource.TestCheckResourceAttr(resourceName, "set_multi_path", "false"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Set Rule Policy - Re-enable Next-Hop Propagation and Multipath") },
				Config:    testAccMSOSetRulePolicyConfigPropagation(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "set_next_hop_propagate", "true"),
					resource.TestCheckResourceAttr(resourceName, "set_multi_path", "true"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update Set Rule Policy - Remove Next-Hop Propagation and Multipath with Explicit Null Values")
				},
				Config: testAccMSOSetRulePolicyConfigNullValues(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "set_next_hop_propagate", "false"),
					resource.TestCheckResourceAttr(resourceName, "set_multi_path", "false"),
					resource.TestCheckResourceAttr(resourceName, "set_route_tag.#", "0"),
					func(s *terraform.State) error {
						if s.RootModule().Resources[resourceName].Primary.Attributes["uuid"] != uuid {
							return fmt.Errorf("removing propagation and multipath with explicit null values changed the policy UUID")
						}
						return nil
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Set Rule Policy - Re-add All Actions") },
				Config:    testAccMSOSetRulePolicyConfigFull("test_set_rule_policy"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Test Set Rule Policy"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "set_route_tag.0.value", "500"),
					resource.TestCheckResourceAttr(resourceName, "set_weight.0.value", "50"),
					resource.TestCheckResourceAttr(resourceName, "set_preference.0.value", "200"),
					resource.TestCheckResourceAttr(resourceName, "set_metric", "100"),
					resource.TestCheckResourceAttr(resourceName, "set_next_hop", "2001:db8::1"),
					resource.TestCheckResourceAttr(resourceName, "set_metric_type", "ospf-type2"),
					resource.TestCheckResourceAttr(resourceName, "set_next_hop_propagate", "false"),
					resource.TestCheckResourceAttr(resourceName, "set_multi_path", "false"),
					resource.TestCheckResourceAttr(resourceName, "set_community.0.community", "no-export"),
					resource.TestCheckResourceAttr(resourceName, "set_community.0.criteria", "replace"),
					resource.TestCheckResourceAttr(resourceName, "set_additional_community.#", "2"),
					testCheckTypeSetStringElemAttr(resourceName, "set_additional_community", "no-advertise"),
					testCheckTypeSetStringElemAttr(resourceName, "set_additional_community", "extended:color:35"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.half_life", "15"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.reuse_limit", "750"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.suppress_limit", "2000"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.max_suppress_time", "60"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend_last_as.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend.0.path_asn.#", "4"),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "10", "order": "1"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "20", "order": "3"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "30", "order": "2"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "30", "order": "4"}),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend_last_as.0.asn_count", "10"),
					capture,
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update Set Rule Policy after Out-of-Band Changes")
					msoClient := testAccPreCheck(t)
					// Resolve the current API array index before applying out-of-band changes.
					response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
					if err != nil {
						t.Fatal(err)
					}
					index, err := GetPolicyIndexByKeyAndValue(response, "uuid", uuid, "tenantPolicyTemplate", "template", "setRulePolicies")
					if err != nil {
						t.Fatal(err)
					}
					payload := container.New()
					payload.Array()
					path := fmt.Sprintf("/tenantPolicyTemplate/template/setRulePolicies/%d", index)
					if err := addPatchPayloadToContainer(payload, "remove", path+"/description", nil); err != nil {
						t.Fatal(err)
					}
					if err := addPatchPayloadToContainer(payload, "add", path+"/setMetric", 321); err != nil {
						t.Fatal(err)
					}
					if err := addPatchPayloadToContainer(payload, "add", path+"/setCommunity", map[string]interface{}{"community": "no-advertise", "criteria": "append"}); err != nil {
						t.Fatal(err)
					}
					if err := doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateID), payload); err != nil {
						t.Fatal(err)
					}
				},
				Config: testAccMSOSetRulePolicyConfigFull("test_set_rule_policy"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Test Set Rule Policy"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "set_route_tag.0.value", "500"),
					resource.TestCheckResourceAttr(resourceName, "set_weight.0.value", "50"),
					resource.TestCheckResourceAttr(resourceName, "set_preference.0.value", "200"),
					resource.TestCheckResourceAttr(resourceName, "set_metric", "100"),
					resource.TestCheckResourceAttr(resourceName, "set_next_hop", "2001:db8::1"),
					resource.TestCheckResourceAttr(resourceName, "set_metric_type", "ospf-type2"),
					resource.TestCheckResourceAttr(resourceName, "set_next_hop_propagate", "false"),
					resource.TestCheckResourceAttr(resourceName, "set_multi_path", "false"),
					resource.TestCheckResourceAttr(resourceName, "set_community.0.community", "no-export"),
					resource.TestCheckResourceAttr(resourceName, "set_community.0.criteria", "replace"),
					resource.TestCheckResourceAttr(resourceName, "set_additional_community.#", "2"),
					testCheckTypeSetStringElemAttr(resourceName, "set_additional_community", "no-advertise"),
					testCheckTypeSetStringElemAttr(resourceName, "set_additional_community", "extended:color:35"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.half_life", "15"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.reuse_limit", "750"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.suppress_limit", "2000"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.max_suppress_time", "60"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend_last_as.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend.0.path_asn.#", "4"),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "10", "order": "1"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "20", "order": "3"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "30", "order": "2"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "30", "order": "4"}),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend_last_as.0.asn_count", "10"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate Set Rule Policy after Out-of-Band Deletion")
					msoClient := testAccPreCheck(t)
					// Resolve the current API array index before deleting the policy outside Terraform.
					response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
					if err != nil {
						t.Fatal(err)
					}
					index, err := GetPolicyIndexByKeyAndValue(response, "uuid", uuid, "tenantPolicyTemplate", "template", "setRulePolicies")
					if err != nil {
						t.Fatal(err)
					}
					policy := response.S("tenantPolicyTemplate", "template", "setRulePolicies").Index(index)
					name, ok := policy.S("name").Data().(string)
					if !ok {
						t.Fatal("policy to delete has no valid name")
					}
					path := fmt.Sprintf("/tenantPolicyTemplate/template/setRulePolicies/%d", index)
					if _, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateID), models.GetRemovePatchPayload(path)); err != nil {
						t.Fatal(err)
					}
					// Confirm the API deletion succeeded before Terraform attempts recreation.
					response, err = msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
					if err != nil {
						t.Fatal(err)
					}
					if _, err := GetPolicyByName(response, name, "tenantPolicyTemplate", "template", "setRulePolicies"); !isPolicyNotFound(err) {
						if err != nil {
							t.Fatalf("verify out-of-band deletion: %v", err)
						}
						t.Fatalf("policy %q still exists after out-of-band deletion", name)
					}
				},
				Config: testAccMSOSetRulePolicyConfigFull("test_set_rule_policy"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Test Set Rule Policy"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "set_route_tag.0.value", "500"),
					resource.TestCheckResourceAttr(resourceName, "set_weight.0.value", "50"),
					resource.TestCheckResourceAttr(resourceName, "set_preference.0.value", "200"),
					resource.TestCheckResourceAttr(resourceName, "set_metric", "100"),
					resource.TestCheckResourceAttr(resourceName, "set_next_hop", "2001:db8::1"),
					resource.TestCheckResourceAttr(resourceName, "set_metric_type", "ospf-type2"),
					resource.TestCheckResourceAttr(resourceName, "set_next_hop_propagate", "false"),
					resource.TestCheckResourceAttr(resourceName, "set_multi_path", "false"),
					resource.TestCheckResourceAttr(resourceName, "set_community.0.community", "no-export"),
					resource.TestCheckResourceAttr(resourceName, "set_community.0.criteria", "replace"),
					resource.TestCheckResourceAttr(resourceName, "set_additional_community.#", "2"),
					testCheckTypeSetStringElemAttr(resourceName, "set_additional_community", "no-advertise"),
					testCheckTypeSetStringElemAttr(resourceName, "set_additional_community", "extended:color:35"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.half_life", "15"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.reuse_limit", "750"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.suppress_limit", "2000"),
					resource.TestCheckResourceAttr(resourceName, "set_dampening.0.max_suppress_time", "60"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend_last_as.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend.0.path_asn.#", "4"),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "10", "order": "1"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "20", "order": "3"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "30", "order": "2"}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "set_as_path.0.prepend.0.path_asn", map[string]string{"asn": "30", "order": "4"}),
					resource.TestCheckResourceAttr(resourceName, "set_as_path.0.prepend_last_as.0.asn_count", "10"),
					func(s *terraform.State) error {
						if s.RootModule().Resources[resourceName].Primary.Attributes["uuid"] == uuid {
							return fmt.Errorf("recreated policy retained its deleted UUID")
						}
						return capture(s)
					},
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Replace Set Rule Policy - Change Template")
					originalTemplateID, originalUUID = templateID, uuid
				},
				Config: testAccMSOSetRulePolicyConfigTemplateReplacement(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant_replacement", "id"),
					func(s *terraform.State) error {
						state := s.RootModule().Resources[resourceName]
						if state.Primary.Attributes["uuid"] == originalUUID {
							return fmt.Errorf("template replacement retained the original UUID")
						}
						// Check the original template to verify replacement removed the old policy.
						response, err := testAccPreCheck(t).GetViaURL(fmt.Sprintf("api/v1/templates/%s", originalTemplateID))
						if err != nil {
							return err
						}
						_, err = GetPolicyIndexByKeyAndValue(response, "uuid", originalUUID, "tenantPolicyTemplate", "template", "setRulePolicies")
						if !isPolicyNotFound(err) {
							if err != nil {
								return err
							}
							return fmt.Errorf("original policy still exists after template replacement")
						}
						return nil
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create Multiple Set Rule Policies") },
				Config:    testAccMSOSetRulePolicyConfigMultiple(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_set_rule_policy.first", "set_weight.0.value", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_set_rule_policy.middle", "set_weight.0.value", "2"),
					resource.TestCheckResourceAttr("mso_tenant_policies_set_rule_policy.last", "set_weight.0.value", "3"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Multiple Set Rule Policies - Remove Middle Policy") },
				Config:    testAccMSOSetRulePolicyConfigMultipleRemoveMiddle(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_set_rule_policy.first", "set_weight.0.value", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_set_rule_policy.last", "set_weight.0.value", "3"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Multiple Set Rule Policies - After Array Index Shift") },
				Config:    testAccMSOSetRulePolicyConfigMultipleUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_set_rule_policy.first", "set_weight.0.value", "10"),
					resource.TestCheckResourceAttr("mso_tenant_policies_set_rule_policy.last", "set_weight.0.value", "30"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Multiple Set Rule Policies - Add Policy after Array Index Shift") },
				Config:    testAccMSOSetRulePolicyConfigMultipleAdd(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_set_rule_policy.first", "set_weight.0.value", "11"),
					resource.TestCheckResourceAttr("mso_tenant_policies_set_rule_policy.last", "set_weight.0.value", "31"),
					resource.TestCheckResourceAttr("mso_tenant_policies_set_rule_policy.new", "set_weight.0.value", "4"),
				),
			},
		},
	})
}

func TestAccMSOTenantPoliciesSetRulePolicyResourceErrors(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVersionCheck(t, "5.1")
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments("mso_tenant_policies_set_rule_policy", "tenantPolicyTemplate", "template", "setRulePolicies"),
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Missing Template (error)") },
				Config:      testAccMSOSetRulePolicyConfigMissingTemplate(),
				ExpectError: regexp.MustCompile(`(?s)(template_id.*required|required.*template_id)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Missing Name (error)") },
				Config:      testAccMSOSetRulePolicyConfigMissingName(),
				ExpectError: regexp.MustCompile(`(?s)(name.*required|required.*name)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Invalid Name (error)") },
				Config:      testAccMSOSetRulePolicyConfigInvalidName(""),
				ExpectError: regexp.MustCompile(`(?s)name.*(length|range|between)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Long Name (error)") },
				Config:      testAccMSOSetRulePolicyConfigInvalidName(strings.Repeat("a", 65)),
				ExpectError: regexp.MustCompile(`(?s)name.*(length|range|between)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - set_route_tag Below Range (error)") },
				Config:      testAccMSOSetRulePolicyConfigInvalidSetRouteTag(-1),
				ExpectError: regexp.MustCompile(`(?s)set_route_tag.*value.*(range|between)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - set_route_tag Above Range (error)") },
				Config:      testAccMSOSetRulePolicyConfigInvalidSetRouteTag(1000000000),
				ExpectError: regexp.MustCompile(`(?s)set_route_tag.*value.*(range|between)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - set_route_tag Missing Value (error)") },
				Config:      testAccMSOSetRulePolicyConfigMissingSetRouteTagValue(),
				ExpectError: regexp.MustCompile(`(?s)(value.*required|required.*value)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - set_route_tag Multiple Blocks (error)") },
				Config:      testAccMSOSetRulePolicyConfigMultipleSetRouteTag(),
				ExpectError: regexp.MustCompile(`(?s)(Too many set_route_tag blocks|set_route_tag.*maximum)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - set_weight Below Range (error)") },
				Config:      testAccMSOSetRulePolicyConfigInvalidSetWeight(-1),
				ExpectError: regexp.MustCompile(`(?s)set_weight.*value.*(range|between)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - set_weight Above Range (error)") },
				Config:      testAccMSOSetRulePolicyConfigInvalidSetWeight(10000),
				ExpectError: regexp.MustCompile(`(?s)set_weight.*value.*(range|between)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - set_weight Missing Value (error)") },
				Config:      testAccMSOSetRulePolicyConfigMissingSetWeightValue(),
				ExpectError: regexp.MustCompile(`(?s)(value.*required|required.*value)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - set_weight Multiple Blocks (error)") },
				Config:      testAccMSOSetRulePolicyConfigMultipleSetWeight(),
				ExpectError: regexp.MustCompile(`(?s)(Too many set_weight blocks|set_weight.*maximum)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - set_preference Below Range (error)") },
				Config:      testAccMSOSetRulePolicyConfigInvalidSetPreference(-1),
				ExpectError: regexp.MustCompile(`(?s)set_preference.*value.*(range|between)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - set_preference Above Range (error)") },
				Config:      testAccMSOSetRulePolicyConfigInvalidSetPreference(1000000000),
				ExpectError: regexp.MustCompile(`(?s)set_preference.*value.*(range|between)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - set_preference Missing Value (error)") },
				Config:      testAccMSOSetRulePolicyConfigMissingSetPreferenceValue(),
				ExpectError: regexp.MustCompile(`(?s)(value.*required|required.*value)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - set_preference Multiple Blocks (error)") },
				Config:      testAccMSOSetRulePolicyConfigMultipleSetPreference(),
				ExpectError: regexp.MustCompile(`(?s)(Too many set_preference blocks|set_preference.*maximum)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - set_metric Below Range (error)") },
				Config:      testAccMSOSetRulePolicyConfigInvalidSetMetric(0),
				ExpectError: regexp.MustCompile(`(?s)set_metric.*(range|between)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - set_metric Above Range (error)") },
				Config:      testAccMSOSetRulePolicyConfigInvalidSetMetric(1000000000),
				ExpectError: regexp.MustCompile(`(?s)set_metric.*(range|between)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Invalid Next Hop (error)") },
				Config:      testAccMSOSetRulePolicyConfigInvalidNextHop(),
				ExpectError: regexp.MustCompile(`(?s)set_next_hop.*(IP|address|valid)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Empty Next Hop (error)") },
				Config:      testAccMSOSetRulePolicyConfigEmptyNextHop(),
				ExpectError: regexp.MustCompile(`(?s)set_next_hop.*(IP|address|valid)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Invalid Metric Type (error)") },
				Config:      testAccMSOSetRulePolicyConfigInvalidMetricType(),
				ExpectError: regexp.MustCompile(`(?s)set_metric_type.*(one of|expected)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Empty Metric Type (error)") },
				Config:      testAccMSOSetRulePolicyConfigEmptyMetricType(),
				ExpectError: regexp.MustCompile(`(?s)set_metric_type.*(one of|expected)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - set_community Empty Community (error)") },
				Config:      testAccMSOSetRulePolicyConfigEmptySetCommunity(),
				ExpectError: regexp.MustCompile(`(?s)set_community.*required.*append.*replace`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - set_community Invalid Criteria (error)") },
				Config:      testAccMSOSetRulePolicyConfigInvalidSetCommunityCriteria(),
				ExpectError: regexp.MustCompile(`(?s)criteria.*(one of|expected)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - set_community Missing Community (error)") },
				Config:      testAccMSOSetRulePolicyConfigMissingSetCommunityCommunity(),
				ExpectError: regexp.MustCompile(`(?s)set_community.*required.*append.*replace`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - set_community Missing Criteria (error)") },
				Config:      testAccMSOSetRulePolicyConfigMissingSetCommunityCriteria(),
				ExpectError: regexp.MustCompile(`(?s)(criteria.*required|required.*criteria)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Empty Additional Community (error)") },
				Config:      testAccMSOSetRulePolicyConfigEmptyAdditionalCommunity(),
				ExpectError: regexp.MustCompile(`(?s)set_additional_community.*(empty|not be empty)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Multiple Community Blocks (error)") },
				Config:      testAccMSOSetRulePolicyConfigMultipleCommunity(),
				ExpectError: regexp.MustCompile(`(?s)(Too many set_community blocks|set_community.*maximum)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Dampening half_life Below Range (error)") },
				Config:      testAccMSOSetRulePolicyConfigInvalidDampeningHalfLife(0),
				ExpectError: regexp.MustCompile(`(?s)half_life.*(range|between)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Dampening half_life Above Range (error)") },
				Config:      testAccMSOSetRulePolicyConfigInvalidDampeningHalfLife(61),
				ExpectError: regexp.MustCompile(`(?s)half_life.*(range|between)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Dampening Missing half_life (error)") },
				Config:      testAccMSOSetRulePolicyConfigMissingDampeningHalfLife(),
				ExpectError: regexp.MustCompile(`(?s)(half_life.*required|required.*half_life)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Dampening reuse_limit Below Range (error)") },
				Config:      testAccMSOSetRulePolicyConfigInvalidDampeningReuseLimit(0),
				ExpectError: regexp.MustCompile(`(?s)reuse_limit.*(range|between)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Dampening reuse_limit Above Range (error)") },
				Config:      testAccMSOSetRulePolicyConfigInvalidDampeningReuseLimit(20001),
				ExpectError: regexp.MustCompile(`(?s)reuse_limit.*(range|between)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Dampening Missing reuse_limit (error)") },
				Config:      testAccMSOSetRulePolicyConfigMissingDampeningReuseLimit(),
				ExpectError: regexp.MustCompile(`(?s)(reuse_limit.*required|required.*reuse_limit)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Dampening suppress_limit Below Range (error)") },
				Config:      testAccMSOSetRulePolicyConfigInvalidDampeningSuppressLimit(0),
				ExpectError: regexp.MustCompile(`(?s)suppress_limit.*(range|between)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Dampening suppress_limit Above Range (error)") },
				Config:      testAccMSOSetRulePolicyConfigInvalidDampeningSuppressLimit(20001),
				ExpectError: regexp.MustCompile(`(?s)suppress_limit.*(range|between)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Dampening Missing suppress_limit (error)") },
				Config:      testAccMSOSetRulePolicyConfigMissingDampeningSuppressLimit(),
				ExpectError: regexp.MustCompile(`(?s)(suppress_limit.*required|required.*suppress_limit)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Dampening max_suppress_time Below Range (error)") },
				Config:      testAccMSOSetRulePolicyConfigInvalidDampeningMaxSuppressTime(0),
				ExpectError: regexp.MustCompile(`(?s)max_suppress_time.*(range|between)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Dampening max_suppress_time Above Range (error)") },
				Config:      testAccMSOSetRulePolicyConfigInvalidDampeningMaxSuppressTime(256),
				ExpectError: regexp.MustCompile(`(?s)max_suppress_time.*(range|between)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Dampening Missing max_suppress_time (error)") },
				Config:      testAccMSOSetRulePolicyConfigMissingDampeningMaxSuppressTime(),
				ExpectError: regexp.MustCompile(`(?s)(max_suppress_time.*required|required.*max_suppress_time)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Dampening Equal Suppress and Reuse Limits (error)") },
				Config:      testAccMSOSetRulePolicyConfigDampeningEqualLimits(),
				ExpectError: regexp.MustCompile(`(?s)set_dampening.*suppress_limit.*larger than.*reuse_limit`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Dampening Suppress Limit Below Reuse Limit (error)") },
				Config:      testAccMSOSetRulePolicyConfigDampeningReversedLimits(),
				ExpectError: regexp.MustCompile(`(?s)set_dampening.*suppress_limit.*larger than.*reuse_limit`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Multiple Dampening Blocks (error)") },
				Config:      testAccMSOSetRulePolicyConfigMultipleDampening(),
				ExpectError: regexp.MustCompile(`(?s)(Too many set_dampening blocks|set_dampening.*maximum)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Empty AS Path Settings (error)") },
				Config:      testAccMSOSetRulePolicyConfigEmptyASPath(),
				ExpectError: regexp.MustCompile(`(?s)set_as_path.*(prepend|at least one)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Multiple AS Path Settings Blocks (error)") },
				Config:      testAccMSOSetRulePolicyConfigMultipleASPath(),
				ExpectError: regexp.MustCompile(`(?s)(Too many set_as_path blocks|set_as_path.*maximum)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Missing ASN Count for Prepend Last AS (error)") },
				Config:      testAccMSOSetRulePolicyConfigMissingASNCount(),
				ExpectError: regexp.MustCompile(`(?s)(asn_count.*required|required.*asn_count)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Multiple Prepend Blocks (error)") },
				Config:      testAccMSOSetRulePolicyConfigMultipleASPathPrepend(),
				ExpectError: regexp.MustCompile(`(?s)(Too many prepend blocks|prepend.*maximum)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Multiple Prepend Last AS Blocks (error)") },
				Config:      testAccMSOSetRulePolicyConfigMultipleASPathPrependLastAS(),
				ExpectError: regexp.MustCompile(`(?s)(Too many prepend_last_as blocks|prepend_last_as.*maximum)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - ASN Count Below Range (error)") },
				Config:      testAccMSOSetRulePolicyConfigInvalidASNCount(0),
				ExpectError: regexp.MustCompile(`(?s)asn_count.*(range|between)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - ASN Count Above Range (error)") },
				Config:      testAccMSOSetRulePolicyConfigInvalidASNCount(11),
				ExpectError: regexp.MustCompile(`(?s)asn_count.*(range|between)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Path asn Below Range (error)") },
				Config:      testAccMSOSetRulePolicyConfigInvalidPathASN(0),
				ExpectError: regexp.MustCompile(`(?s)asn.*(range|between)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Path asn Above Range (error)") },
				Config:      testAccMSOSetRulePolicyConfigInvalidPathASN(4294967296),
				ExpectError: regexp.MustCompile(`(?s)asn.*(range|between)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Missing Path asn (error)") },
				Config:      testAccMSOSetRulePolicyConfigMissingPathASN(),
				ExpectError: regexp.MustCompile(`(?s)(asn.*required|required.*asn)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Path order Below Range (error)") },
				Config:      testAccMSOSetRulePolicyConfigInvalidPathOrder(-1),
				ExpectError: regexp.MustCompile(`(?s)order.*(range|between)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Path order Above Range (error)") },
				Config:      testAccMSOSetRulePolicyConfigInvalidPathOrder(32),
				ExpectError: regexp.MustCompile(`(?s)order.*(range|between)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Missing Path order (error)") },
				Config:      testAccMSOSetRulePolicyConfigMissingPathOrder(),
				ExpectError: regexp.MustCompile(`(?s)(order.*required|required.*order)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Duplicate AS Path Orders (error)") },
				Config:      testAccMSOSetRulePolicyConfigDuplicateASPathOrders(),
				ExpectError: regexp.MustCompile(`(?s)path_asn order 0 is duplicated`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Multipath Without Propagation (error)") },
				Config:      testAccMSOSetRulePolicyConfigMultipathWithoutPropagation(),
				ExpectError: regexp.MustCompile(`(?s)set_multi_path.*requires.*set_next_hop_propagate`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Set Rule Policy - Propagation With Zero Route Tag (error)") },
				Config:      testAccMSOSetRulePolicyConfigPropagationWithZeroRouteTag(),
				ExpectError: regexp.MustCompile(`(?s)set_next_hop_propagate.*cannot.*set_route_tag`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create Set Rule Policy - Valid Configuration") },
				Config:    testAccMSOSetRulePolicyConfigMinimal(),
				Check:     resource.TestCheckResourceAttrSet("mso_tenant_policies_set_rule_policy.test", "uuid"),
			},
		},
	})
}

func testAccMSOSetRulePolicyConfigMinimal() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigFull(name string) string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "%s"
  description = "Test Set Rule Policy"
  set_route_tag {
    value = 500
  }
  set_weight {
    value = 50
  }
  set_preference {
    value = 200
  }
  set_metric      = 100
  set_next_hop    = "2001:db8::1"
  set_metric_type = "ospf-type2"
  set_community {
    community = "no-export"
    criteria  = "replace"
  }
  set_additional_community = ["no-advertise", "extended:color:35"]
  set_dampening {
    half_life         = 15
    reuse_limit       = 750
    suppress_limit    = 2000
    max_suppress_time = 60
  }
  set_as_path {
    prepend {
      path_asn {
        asn   = 10
        order = 1
      }
      path_asn {
        asn   = 20
        order = 3
      }
      path_asn {
        asn   = 30
        order = 2
      }
      path_asn {
        asn   = 30
        order = 4
      }
    }
    prepend_last_as {
      asn_count = 10
    }
  }
}`, testAccMSOTemplateResourceTenantConfig(), name)
}

func testAccMSOSetRulePolicyConfigNullValues() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id              = mso_template.template_tenant.id
  name                     = "test_set_rule_policy"
  description              = null
  set_metric               = null
  set_next_hop             = null
  set_metric_type          = null
  set_next_hop_propagate   = null
  set_multi_path           = null
  set_additional_community = null
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigMinValues() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_route_tag {
    value = 0
  }
  set_weight {
    value = 0
  }
  set_preference {
    value = 0
  }
  set_metric      = 1
  set_next_hop    = "192.0.2.1"
  set_metric_type = "ospf-type1"
  set_community {
    community = "no-advertise"
    criteria  = "append"
  }
  set_dampening {
    half_life         = 1
    reuse_limit       = 1
    suppress_limit    = 2
    max_suppress_time = 1
  }
  set_as_path {
    prepend {
      path_asn {
        asn   = 1
        order = 0
      }
    }
    prepend_last_as {
      asn_count = 1
    }
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigMaxValues() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_route_tag {
    value = 999999999
  }
  set_weight {
    value = 9999
  }
  set_preference {
    value = 999999999
  }
  set_metric      = 999999999
  set_metric_type = "none"
  set_community {
    criteria = "none"
  }
  set_dampening {
    half_life         = 60
    reuse_limit       = 19999
    suppress_limit    = 20000
    max_suppress_time = 255
  }
  set_as_path {
    prepend {
      path_asn {
        asn   = 4294967295
        order = 31
      }
    }
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigPropagation() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id            = mso_template.template_tenant.id
  name                   = "test_set_rule_policy"
  set_next_hop_propagate = true
  set_multi_path         = true
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigDisabled() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id            = mso_template.template_tenant.id
  name                   = "test_set_rule_policy"
  set_next_hop_propagate = false
  set_multi_path         = false
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigTemplateReplacement() string {
	return fmt.Sprintf(`%s
resource "mso_template" "template_tenant_replacement" {
  template_name = "test_template_tenant_replacement"
  template_type = "tenant"
  tenant_id     = mso_tenant.%s.id
}
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant_replacement.id
  name        = "test_set_rule_policy"
}`, testAccMSOTemplateResourceTenantConfig(), msoTemplateTenantName)
}

func testAccMSOSetRulePolicyConfigMultiple() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "first" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy_first"
  set_weight {
    value = 1
  }
}
resource "mso_tenant_policies_set_rule_policy" "middle" {
  # Order writes to the shared API array and destroy policies in reverse order.
  depends_on  = [mso_tenant_policies_set_rule_policy.first]
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy_middle"
  set_weight {
    value = 2
  }
}
resource "mso_tenant_policies_set_rule_policy" "last" {
  depends_on  = [mso_tenant_policies_set_rule_policy.middle]
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy_last"
  set_weight {
    value = 3
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigMultipleRemoveMiddle() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "first" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy_first"
  set_weight {
    value = 1
  }
}
resource "mso_tenant_policies_set_rule_policy" "last" {
  depends_on  = [mso_tenant_policies_set_rule_policy.first]
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy_last"
  set_weight {
    value = 3
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigMultipleUpdate() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "first" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy_first"
  set_weight {
    value = 10
  }
}
resource "mso_tenant_policies_set_rule_policy" "last" {
  depends_on  = [mso_tenant_policies_set_rule_policy.first]
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy_last"
  set_weight {
    value = 30
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigMultipleAdd() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "first" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy_first"
  set_weight {
    value = 11
  }
}
resource "mso_tenant_policies_set_rule_policy" "last" {
  depends_on  = [mso_tenant_policies_set_rule_policy.first]
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy_last"
  set_weight {
    value = 31
  }
}
resource "mso_tenant_policies_set_rule_policy" "new" {
  depends_on  = [mso_tenant_policies_set_rule_policy.last]
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy_new"
  set_weight {
    value = 4
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigMissingTemplate() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  name = "test_set_rule_policy"
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigMissingName() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigInvalidName(name string) string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = %q
}`, testAccMSOTemplateResourceTenantConfig(), name)
}

func testAccMSOSetRulePolicyConfigInvalidSetRouteTag(value int) string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_route_tag {
    value = %d
  }
}`, testAccMSOTemplateResourceTenantConfig(), value)
}

func testAccMSOSetRulePolicyConfigMissingSetRouteTagValue() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_route_tag {
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigMultipleSetRouteTag() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_route_tag {
    value = 1
  }
  set_route_tag {
    value = 2
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigInvalidSetWeight(value int) string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_weight {
    value = %d
  }
}`, testAccMSOTemplateResourceTenantConfig(), value)
}

func testAccMSOSetRulePolicyConfigMissingSetWeightValue() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_weight {
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigMultipleSetWeight() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_weight {
    value = 1
  }
  set_weight {
    value = 2
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigInvalidSetPreference(value int) string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_preference {
    value = %d
  }
}`, testAccMSOTemplateResourceTenantConfig(), value)
}

func testAccMSOSetRulePolicyConfigMissingSetPreferenceValue() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_preference {
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigMultipleSetPreference() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_preference {
    value = 1
  }
  set_preference {
    value = 2
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigInvalidSetMetric(value int) string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_metric  = %d
}`, testAccMSOTemplateResourceTenantConfig(), value)
}

func testAccMSOSetRulePolicyConfigInvalidNextHop() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id  = mso_template.template_tenant.id
  name         = "test_set_rule_policy"
  set_next_hop = "invalid"
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigEmptyNextHop() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id  = mso_template.template_tenant.id
  name         = "test_set_rule_policy"
  set_next_hop = ""
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigInvalidMetricType() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id     = mso_template.template_tenant.id
  name            = "test_set_rule_policy"
  set_metric_type = "invalid"
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigEmptyMetricType() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id     = mso_template.template_tenant.id
  name            = "test_set_rule_policy"
  set_metric_type = ""
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigEmptySetCommunity() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_community {
    community = ""
    criteria  = "append"
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigInvalidSetCommunityCriteria() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_community {
    community = "no-export"
    criteria  = "invalid"
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigMissingSetCommunityCommunity() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_community {
    criteria = "append"
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigMissingSetCommunityCriteria() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_community {
    community = "no-export"
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigEmptyAdditionalCommunity() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id              = mso_template.template_tenant.id
  name                     = "test_set_rule_policy"
  set_additional_community = [""]
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigMultipleCommunity() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_community {
    community = "no-export"
    criteria  = "replace"
  }
  set_community {
    community = "no-advertise"
    criteria  = "append"
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigInvalidDampeningHalfLife(value int) string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_dampening {
    half_life         = %d
    reuse_limit       = 1
    suppress_limit    = 2
    max_suppress_time = 1
  }
}`, testAccMSOTemplateResourceTenantConfig(), value)
}

func testAccMSOSetRulePolicyConfigMissingDampeningHalfLife() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_dampening {
    reuse_limit       = 1
    suppress_limit    = 2
    max_suppress_time = 1
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigInvalidDampeningReuseLimit(value int) string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_dampening {
    half_life         = 1
    reuse_limit       = %d
    suppress_limit    = 2
    max_suppress_time = 1
  }
}`, testAccMSOTemplateResourceTenantConfig(), value)
}

func testAccMSOSetRulePolicyConfigMissingDampeningReuseLimit() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_dampening {
    half_life         = 1
    suppress_limit    = 2
    max_suppress_time = 1
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigInvalidDampeningSuppressLimit(value int) string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_dampening {
    half_life         = 1
    reuse_limit       = 1
    suppress_limit    = %d
    max_suppress_time = 1
  }
}`, testAccMSOTemplateResourceTenantConfig(), value)
}

func testAccMSOSetRulePolicyConfigMissingDampeningSuppressLimit() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_dampening {
    half_life         = 1
    reuse_limit       = 1
    max_suppress_time = 1
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigInvalidDampeningMaxSuppressTime(value int) string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_dampening {
    half_life         = 1
    reuse_limit       = 1
    suppress_limit    = 2
    max_suppress_time = %d
  }
}`, testAccMSOTemplateResourceTenantConfig(), value)
}

func testAccMSOSetRulePolicyConfigMissingDampeningMaxSuppressTime() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_dampening {
    half_life      = 1
    reuse_limit    = 1
    suppress_limit = 2
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigMultipleDampening() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_dampening {
    half_life         = 1
    reuse_limit       = 1
    suppress_limit    = 2
    max_suppress_time = 1
  }
  set_dampening {
    half_life         = 1
    reuse_limit       = 1
    suppress_limit    = 2
    max_suppress_time = 1
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigEmptyASPath() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_as_path {
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigInvalidASNCount(value int) string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_as_path {
    prepend_last_as {
      asn_count = %d
    }
  }
}`, testAccMSOTemplateResourceTenantConfig(), value)
}

func testAccMSOSetRulePolicyConfigInvalidPathASN(value int64) string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_as_path {
    prepend {
      path_asn {
        asn   = %d
        order = 0
      }
    }
  }
}`, testAccMSOTemplateResourceTenantConfig(), value)
}

func testAccMSOSetRulePolicyConfigMissingPathASN() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_as_path {
    prepend {
      path_asn {
        order = 0
      }
    }
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigInvalidPathOrder(value int) string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_as_path {
    prepend {
      path_asn {
        asn   = 65001
        order = %d
      }
    }
  }
}`, testAccMSOTemplateResourceTenantConfig(), value)
}

func testAccMSOSetRulePolicyConfigMissingPathOrder() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_as_path {
    prepend {
      path_asn {
        asn = 65001
      }
    }
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigDuplicateASPathOrders() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_as_path {
    prepend {
      path_asn {
        asn   = 65001
        order = 0
      }
      path_asn {
        asn   = 65002
        order = 0
      }
    }
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigMultipathWithoutPropagation() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id    = mso_template.template_tenant.id
  name           = "test_set_rule_policy"
  set_multi_path = true
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigPropagationWithZeroRouteTag() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id            = mso_template.template_tenant.id
  name                   = "test_set_rule_policy"
  set_next_hop_propagate = true
  set_route_tag {
    value = 0
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigCreateZeros() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
}
resource "mso_tenant_policies_set_rule_policy" "zero" {
  depends_on  = [mso_tenant_policies_set_rule_policy.test]
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy_zero"
  description = "{}"
  set_route_tag {
    value = 0
  }
  set_weight {
    value = 0
  }
  set_preference {
    value = 0
  }
  set_metric = 1
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigDampeningEqualLimits() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_dampening {
    half_life         = 1
    reuse_limit       = 1
    suppress_limit    = 1
    max_suppress_time = 1
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigDampeningReversedLimits() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_dampening {
    half_life         = 1
    reuse_limit       = 2
    suppress_limit    = 1
    max_suppress_time = 1
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigASPathReordered() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  description = "Test Set Rule Policy"
  set_route_tag {
    value = 500
  }
  set_weight {
    value = 50
  }
  set_preference {
    value = 200
  }
  set_metric      = 100
  set_next_hop    = "2001:db8::1"
  set_metric_type = "ospf-type2"
  set_community {
    community = "no-export"
    criteria  = "replace"
  }
  set_additional_community = ["no-advertise", "extended:color:35"]
  set_dampening {
    half_life         = 15
    reuse_limit       = 750
    suppress_limit    = 2000
    max_suppress_time = 60
  }
  set_as_path {
    prepend_last_as {
      asn_count = 10
    }
    prepend {
      path_asn {
        asn   = 30
        order = 4
      }
      path_asn {
        asn   = 30
        order = 2
      }
      path_asn {
        asn   = 20
        order = 3
      }
      path_asn {
        asn   = 10
        order = 1
      }
    }
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigASPathUpdated() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  description = "Test Set Rule Policy"
  set_route_tag {
    value = 500
  }
  set_weight {
    value = 50
  }
  set_preference {
    value = 200
  }
  set_metric      = 100
  set_next_hop    = "2001:db8::1"
  set_metric_type = "ospf-type2"
  set_community {
    community = "no-export"
    criteria  = "replace"
  }
  set_additional_community = ["no-advertise", "extended:color:35"]
  set_dampening {
    half_life         = 15
    reuse_limit       = 750
    suppress_limit    = 2000
    max_suppress_time = 60
  }
  set_as_path {
    prepend {
      path_asn {
        asn   = 10
        order = 1
      }
      path_asn {
        asn   = 20
        order = 3
      }
      path_asn {
        asn   = 30
        order = 2
      }
      path_asn {
        asn   = 30
        order = 4
      }
    }
    prepend_last_as {
      asn_count = 3
    }
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigPartialCollections() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  description = "Test Set Rule Policy"
  set_route_tag {
    value = 500
  }
  set_weight {
    value = 50
  }
  set_preference {
    value = 200
  }
  set_metric      = 100
  set_next_hop    = "2001:db8::1"
  set_metric_type = "ospf-type2"
  set_community {
    community = "no-export"
    criteria  = "replace"
  }
  set_additional_community = ["extended:color:35"]
  set_dampening {
    half_life         = 15
    reuse_limit       = 750
    suppress_limit    = 2000
    max_suppress_time = 60
  }
  set_as_path {
    prepend {
      path_asn {
        asn   = 10
        order = 1
      }
      path_asn {
        asn   = 30
        order = 2
      }
    }
    prepend_last_as {
      asn_count = 3
    }
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigASPathLastOnly() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  description = "Test Set Rule Policy"
  set_route_tag {
    value = 500
  }
  set_weight {
    value = 50
  }
  set_preference {
    value = 200
  }
  set_metric      = 100
  set_next_hop    = "2001:db8::1"
  set_metric_type = "ospf-type2"
  set_community {
    community = "no-export"
    criteria  = "replace"
  }
  set_additional_community = ["no-advertise", "extended:color:35"]
  set_dampening {
    half_life         = 15
    reuse_limit       = 750
    suppress_limit    = 2000
    max_suppress_time = 60
  }
  set_as_path {
    prepend_last_as {
      asn_count = 3
    }
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigASPathEmptyPrepend() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  description = "Test Set Rule Policy"
  set_route_tag {
    value = 500
  }
  set_weight {
    value = 50
  }
  set_preference {
    value = 200
  }
  set_metric      = 100
  set_next_hop    = "2001:db8::1"
  set_metric_type = "ospf-type2"
  set_community {
    community = "no-export"
    criteria  = "replace"
  }
  set_additional_community = ["no-advertise", "extended:color:35"]
  set_dampening {
    half_life         = 15
    reuse_limit       = 750
    suppress_limit    = 2000
    max_suppress_time = 60
  }
  set_as_path {
    prepend {
    }
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigMissingASNCount() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_as_path {
    prepend_last_as {
    }
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigMultipleASPathPrepend() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_as_path {
    prepend {
    }
    prepend {
    }
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigMultipleASPathPrependLastAS() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_as_path {
    prepend_last_as {
      asn_count = 1
    }
    prepend_last_as {
      asn_count = 2
    }
  }
}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOSetRulePolicyConfigMultipleASPath() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_set_rule_policy"
  set_as_path {
    prepend {
    }
  }
  set_as_path {
    prepend_last_as {
      asn_count = 1
    }
  }
}`, testAccMSOTemplateResourceTenantConfig())
}
