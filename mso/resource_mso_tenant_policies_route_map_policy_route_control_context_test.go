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

func TestAccMSORouteMapPolicyContextResource(t *testing.T) {
	resourceName := "mso_tenant_policies_route_map_policy_route_control_context.test_context"
	resourceName1 := "mso_tenant_policies_route_map_policy_route_control_context.ctx_1"
	resourceName2 := "mso_tenant_policies_route_map_policy_route_control_context.ctx_2"
	resourceName3 := "mso_tenant_policies_route_map_policy_route_control_context.ctx_3"
	resourceName4 := "mso_tenant_policies_route_map_policy_route_control_context.ctx_4"
	var templateID, parentUUID, contextID string
	// Capture the API identifiers for out-of-band changes and replacement checks.
	captureParent := func(s *terraform.State) error {
		context, ok := s.RootModule().Resources[resourceName]
		if !ok || context.Primary == nil {
			return fmt.Errorf("resource %s not found in state", resourceName)
		}
		templateID = context.Primary.Attributes["template_id"]
		parentUUID = context.Primary.Attributes["route_map_policy_uuid"]
		contextID = context.Primary.ID
		if contextID != fmt.Sprintf("templateId/%s/RouteMapPolicy/%s/context/%s", templateID, parentUUID, context.Primary.Attributes["name"]) {
			return fmt.Errorf("context ID does not contain its template ID, parent UUID and name")
		}
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVersionCheck(t, "5.1")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create Route Map Policy Context with Defaults") },
				Config:    testAccMSORouteMapPolicyContextConfig_Create(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "ctx_1"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "order", "0"),
					resource.TestCheckResourceAttr(resourceName, "action", "permit"),
					resource.TestCheckResourceAttr(resourceName, "set_rule_uuid", ""),
					resource.TestCheckResourceAttr(resourceName, "match_rule_uuids.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_tenant_policies_route_map_policy_route_control.test_policy", "template_id"),
					resource.TestCheckResourceAttrPair(resourceName, "route_map_policy_uuid", "mso_tenant_policies_route_map_policy_route_control.test_policy", "uuid"),
					captureParent,
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Route Map Policy Context Attributes") },
				Config:    testAccMSORouteMapPolicyContextConfig_UpdateAttrs(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
					resource.TestCheckResourceAttr(resourceName, "order", "2"),
					resource.TestCheckResourceAttr(resourceName, "action", "deny"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add Route Map Policy Context Match and Set Rules with Maximum Order") },
				Config:    testAccMSORouteMapPolicyContextConfig_WithRules(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "set_rule_uuid", "mso_tenant_policies_set_rule_policy.test_set_rule", "uuid"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "match_rule_uuids.*", "mso_tenant_policies_match_rule_policy.test_match_rule_1", "uuid"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "match_rule_uuids.*", "mso_tenant_policies_match_rule_policy.test_match_rule_2", "uuid"),
					resource.TestCheckResourceAttr(resourceName, "match_rule_uuids.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "order", "9"),
					resource.TestCheckResourceAttr(resourceName, "action", "deny"),
				),
			},
			{
				PreConfig:          func() { fmt.Println("Test: Plan Route Map Policy Context with Reordered Match Rule UUIDs") },
				Config:             testAccMSORouteMapPolicyContextConfigReorderedRules(),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove One Route Map Policy Context Match Rule") },
				Config:    testAccMSORouteMapPolicyContextConfigOneMatchRule(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "ctx_1"),
					resource.TestCheckResourceAttr(resourceName, "description", "Context with rules"),
					resource.TestCheckResourceAttr(resourceName, "order", "9"),
					resource.TestCheckResourceAttr(resourceName, "action", "deny"),
					resource.TestCheckResourceAttrPair(resourceName, "set_rule_uuid", "mso_tenant_policies_set_rule_policy.test_set_rule", "uuid"),
					resource.TestCheckResourceAttr(resourceName, "match_rule_uuids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "match_rule_uuids.*", "mso_tenant_policies_match_rule_policy.test_match_rule_2", "uuid"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Restore Route Map Policy Context Match Rules") },
				Config:    testAccMSORouteMapPolicyContextConfig_WithRules(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "ctx_1"),
					resource.TestCheckResourceAttr(resourceName, "description", "Context with rules"),
					resource.TestCheckResourceAttr(resourceName, "order", "9"),
					resource.TestCheckResourceAttr(resourceName, "action", "deny"),
					resource.TestCheckResourceAttrPair(resourceName, "set_rule_uuid", "mso_tenant_policies_set_rule_policy.test_set_rule", "uuid"),
					resource.TestCheckResourceAttr(resourceName, "match_rule_uuids.#", "2"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "match_rule_uuids.*", "mso_tenant_policies_match_rule_policy.test_match_rule_1", "uuid"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "match_rule_uuids.*", "mso_tenant_policies_match_rule_policy.test_match_rule_2", "uuid"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Restore Route Map Policy Context after Out-of-Band Field Changes")
					c := testAccPreCheck(t)
					response, err := c.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
					if err != nil {
						t.Fatal(err)
					}
					// Resolve current indices before changing the context outside Terraform.
					policyIndex, err := GetPolicyIndexByKeyAndValue(response, "uuid", parentUUID, "tenantPolicyTemplate", "template", "routeMapPolicies")
					if err != nil {
						t.Fatal(err)
					}
					policy := response.S("tenantPolicyTemplate", "template", "routeMapPolicies").Index(policyIndex)
					contextIndex, err := GetPolicyIndexByKeyAndValue(policy, "name", "ctx_1", "contexts")
					if err != nil {
						t.Fatal(err)
					}
					path := fmt.Sprintf("/tenantPolicyTemplate/template/routeMapPolicies/%d/contexts/%d", policyIndex, contextIndex)
					payload := container.New()
					if _, err := payload.Array(); err != nil {
						t.Fatal(err)
					}
					if err := addPatchPayloadToContainer(payload, "remove", fmt.Sprintf("%s/description", path), nil); err != nil {
						t.Fatal(err)
					}
					if err := addPatchPayloadToContainer(payload, "remove", fmt.Sprintf("%s/setRuleRef", path), nil); err != nil {
						t.Fatal(err)
					}
					if err := addPatchPayloadToContainer(payload, "remove", fmt.Sprintf("%s/matchRules", path), nil); err != nil {
						t.Fatal(err)
					}
					if err := addPatchPayloadToContainer(payload, "add", fmt.Sprintf("%s/order", path), 1); err != nil {
						t.Fatal(err)
					}
					if err := addPatchPayloadToContainer(payload, "add", fmt.Sprintf("%s/action", path), "permit"); err != nil {
						t.Fatal(err)
					}
					if err := doPatchRequest(c, fmt.Sprintf("api/v1/templates/%s", templateID), payload); err != nil {
						t.Fatal(err)
					}
				},
				Config: testAccMSORouteMapPolicyContextConfig_WithRules(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "ctx_1"),
					resource.TestCheckResourceAttr(resourceName, "description", "Context with rules"),
					resource.TestCheckResourceAttr(resourceName, "order", "9"),
					resource.TestCheckResourceAttr(resourceName, "action", "deny"),
					resource.TestCheckResourceAttrPair(resourceName, "set_rule_uuid", "mso_tenant_policies_set_rule_policy.test_set_rule", "uuid"),
					resource.TestCheckResourceAttr(resourceName, "match_rule_uuids.#", "2"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "match_rule_uuids.*", "mso_tenant_policies_match_rule_policy.test_match_rule_1", "uuid"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "match_rule_uuids.*", "mso_tenant_policies_match_rule_policy.test_match_rule_2", "uuid"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Reset Route Map Policy Context Optional Attributes with Null Values") },
				Config:    testAccMSORouteMapPolicyContextConfigNullValues(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "ctx_1"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "order", "0"),
					resource.TestCheckResourceAttr(resourceName, "action", "permit"),
					resource.TestCheckResourceAttr(resourceName, "set_rule_uuid", ""),
					resource.TestCheckResourceAttr(resourceName, "match_rule_uuids.#", "0"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Restore Route Map Policy Context Optional Attributes after Null Values") },
				Config:    testAccMSORouteMapPolicyContextConfig_WithRules(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "ctx_1"),
					resource.TestCheckResourceAttr(resourceName, "description", "Context with rules"),
					resource.TestCheckResourceAttr(resourceName, "order", "9"),
					resource.TestCheckResourceAttr(resourceName, "action", "deny"),
					resource.TestCheckResourceAttrPair(resourceName, "set_rule_uuid", "mso_tenant_policies_set_rule_policy.test_set_rule", "uuid"),
					resource.TestCheckResourceAttr(resourceName, "match_rule_uuids.#", "2"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "match_rule_uuids.*", "mso_tenant_policies_match_rule_policy.test_match_rule_1", "uuid"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "match_rule_uuids.*", "mso_tenant_policies_match_rule_policy.test_match_rule_2", "uuid"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Clear Route Map Policy Context Rules and Description with Explicit Empty Values")
				},
				Config: testAccMSORouteMapPolicyContextConfigClearRules(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "set_rule_uuid", ""),
					resource.TestCheckResourceAttr(resourceName, "match_rule_uuids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "order", "0"),
					resource.TestCheckResourceAttr(resourceName, "action", "permit"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Set Route Map Policy Context Rules and Description") },
				Config:    testAccMSORouteMapPolicyContextConfig_WithRules(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Context with rules"),
					resource.TestCheckResourceAttr(resourceName, "match_rule_uuids.#", "2"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove Route Map Policy Context Rules and Description") },
				Config:    testAccMSORouteMapPolicyContextConfig_Create(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "set_rule_uuid", ""),
					resource.TestCheckResourceAttr(resourceName, "match_rule_uuids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "order", "0"),
					resource.TestCheckResourceAttr(resourceName, "action", "permit"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate Route Map Policy Context after Out-of-Band Deletion")
					c := testAccPreCheck(t)
					response, err := c.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
					if err != nil {
						t.Fatal(err)
					}
					// Resolve current indices before deleting the context outside Terraform.
					policyIndex, err := GetPolicyIndexByKeyAndValue(response, "uuid", parentUUID, "tenantPolicyTemplate", "template", "routeMapPolicies")
					if err != nil {
						t.Fatal(err)
					}
					policy := response.S("tenantPolicyTemplate", "template", "routeMapPolicies").Index(policyIndex)
					contextIndex, err := GetPolicyIndexByKeyAndValue(policy, "name", "ctx_1", "contexts")
					if err != nil {
						t.Fatal(err)
					}
					path := fmt.Sprintf("/tenantPolicyTemplate/template/routeMapPolicies/%d/contexts/%d", policyIndex, contextIndex)
					if _, err := c.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateID), models.GetRemovePatchPayload(path)); err != nil {
						t.Fatal(err)
					}
					response, err = c.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
					if err != nil {
						t.Fatal(err)
					}
					policy = response.S("tenantPolicyTemplate", "template", "routeMapPolicies").Index(policyIndex)
					if _, err := GetPolicyByName(policy, "ctx_1", "contexts"); !isPolicyNotFound(err) {
						t.Fatalf("expected context to be deleted, got %v", err)
					}
				},
				Config: testAccMSORouteMapPolicyContextConfig_Create(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "ctx_1"),
					resource.TestCheckResourceAttr(resourceName, "order", "0"),
					resource.TestCheckResourceAttr(resourceName, "action", "permit"),
					captureParent,
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate Route Map Policy Context after Parent Out-of-Band Deletion")
					c := testAccPreCheck(t)
					response, err := c.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
					if err != nil {
						t.Fatal(err)
					}
					policyIndex, err := GetPolicyIndexByKeyAndValue(response, "uuid", parentUUID, "tenantPolicyTemplate", "template", "routeMapPolicies")
					if err != nil {
						t.Fatal(err)
					}
					path := fmt.Sprintf("/tenantPolicyTemplate/template/routeMapPolicies/%d", policyIndex)
					if _, err := c.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateID), models.GetRemovePatchPayload(path)); err != nil {
						t.Fatal(err)
					}
					response, err = c.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
					if err != nil {
						t.Fatal(err)
					}
					if _, err := GetPolicyIndexByKeyAndValue(response, "uuid", parentUUID, "tenantPolicyTemplate", "template", "routeMapPolicies"); !isPolicyNotFound(err) {
						t.Fatalf("expected parent policy to be deleted, got %v", err)
					}
				},
				Config: testAccMSORouteMapPolicyContextConfig_Create(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "ctx_1"),
					resource.TestCheckResourceAttr(resourceName, "order", "0"),
					resource.TestCheckResourceAttr(resourceName, "action", "permit"),
					resource.TestCheckResourceAttrPair(resourceName, "route_map_policy_uuid", "mso_tenant_policies_route_map_policy_route_control.test_policy", "uuid"),
					func(s *terraform.State) error {
						if s.RootModule().Resources[resourceName].Primary.Attributes["route_map_policy_uuid"] == parentUUID {
							return fmt.Errorf("recreating the parent retained its previous UUID")
						}
						return nil
					},
					captureParent,
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Recreate Route Map Policy Context with a New Name") },
				Config:    testAccMSORouteMapPolicyContextConfigRename(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "ctx_renamed"),
					resource.TestCheckResourceAttr(resourceName, "description", "Renamed context"),
					resource.TestCheckResourceAttr(resourceName, "order", "3"),
					resource.TestCheckResourceAttr(resourceName, "action", "deny"),
					func(s *terraform.State) error {
						response, err := testAccPreCheck(t).GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
						if err != nil {
							return err
						}
						policyIndex, err := GetPolicyIndexByKeyAndValue(response, "uuid", parentUUID, "tenantPolicyTemplate", "template", "routeMapPolicies")
						if err != nil {
							return err
						}
						policy := response.S("tenantPolicyTemplate", "template", "routeMapPolicies").Index(policyIndex)
						if _, err := GetPolicyByName(policy, "ctx_1", "contexts"); !isPolicyNotFound(err) {
							if err != nil {
								return err
							}
							return fmt.Errorf("context ctx_1 still exists in the original parent policy")
						}
						return nil
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Recreate Route Map Policy Context in Another Parent Policy") },
				Config:    testAccMSORouteMapPolicyContextConfigChangeParent(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "ctx_renamed"),
					resource.TestCheckResourceAttr(resourceName, "description", "Renamed context"),
					resource.TestCheckResourceAttr(resourceName, "order", "3"),
					resource.TestCheckResourceAttr(resourceName, "action", "deny"),
					resource.TestCheckResourceAttrPair(resourceName, "route_map_policy_uuid", "mso_tenant_policies_route_map_policy_route_control.second_policy", "uuid"),
					func(s *terraform.State) error {
						response, err := testAccPreCheck(t).GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
						if err != nil {
							return err
						}
						policyIndex, err := GetPolicyIndexByKeyAndValue(response, "uuid", parentUUID, "tenantPolicyTemplate", "template", "routeMapPolicies")
						if err != nil {
							return err
						}
						policy := response.S("tenantPolicyTemplate", "template", "routeMapPolicies").Index(policyIndex)
						if _, err := GetPolicyByName(policy, "ctx_renamed", "contexts"); !isPolicyNotFound(err) {
							if err != nil {
								return err
							}
							return fmt.Errorf("context ctx_renamed still exists in the original parent policy")
						}
						return nil
					},
					captureParent,
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Retain Route Map Policy Context when Renaming Its Parent Policy") },
				Config:    testAccMSORouteMapPolicyContextConfigRenameParent(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_route_map_policy_route_control.second_policy", "name", "second_policy_renamed"),
					resource.TestCheckResourceAttr(resourceName, "name", "ctx_renamed"),
					resource.TestCheckResourceAttr(resourceName, "description", "Renamed context"),
					resource.TestCheckResourceAttr(resourceName, "order", "3"),
					resource.TestCheckResourceAttr(resourceName, "action", "deny"),
					resource.TestCheckResourceAttrPair(resourceName, "route_map_policy_uuid", "mso_tenant_policies_route_map_policy_route_control.second_policy", "uuid"),
					func(s *terraform.State) error {
						context := s.RootModule().Resources[resourceName]
						if context == nil || context.Primary == nil {
							return fmt.Errorf("context not found in state")
						}
						if context.Primary.ID != contextID || context.Primary.Attributes["route_map_policy_uuid"] != parentUUID {
							return fmt.Errorf("renaming the parent changed its context ID or parent UUID")
						}
						response, err := testAccPreCheck(t).GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
						if err != nil {
							return err
						}
						policyIndex, err := GetPolicyIndexByKeyAndValue(response, "uuid", parentUUID, "tenantPolicyTemplate", "template", "routeMapPolicies")
						if err != nil {
							return err
						}
						count, err := response.S("tenantPolicyTemplate", "template", "routeMapPolicies").Index(policyIndex).ArrayCount("contexts")
						if err != nil {
							return err
						}
						if count != 1 {
							return fmt.Errorf("expected one context after parent rename, got %d", count)
						}
						return nil
					},
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Route Map Policy Context") },
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import Missing Route Map Policy Context") },
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					context := s.RootModule().Resources[resourceName]
					if context == nil || context.Primary == nil {
						return "", fmt.Errorf("context not found in state")
					}
					return fmt.Sprintf("templateId/%s/RouteMapPolicy/%s/context/tf_missing_context", context.Primary.Attributes["template_id"], context.Primary.Attributes["route_map_policy_uuid"]), nil
				},
				ExpectError: regexp.MustCompile(`(?s)cannot import .*: resource\s+not\s+found`),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import Route Map Policy Context with Missing Parent Policy") },
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					context := s.RootModule().Resources[resourceName]
					if context == nil || context.Primary == nil {
						return "", fmt.Errorf("context not found in state")
					}
					return fmt.Sprintf("templateId/%s/RouteMapPolicy/00000000-0000-0000-0000-000000000000/context/%s", context.Primary.Attributes["template_id"], context.Primary.Attributes["name"]), nil
				},
				ExpectError: regexp.MustCompile(`(?s)cannot import .*: resource\s+not\s+found`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Creating 3 Route Control Contexts") },
				Config:    testAccRouteControlContextMultipleCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName1, "name", "ctx_1"),
					resource.TestCheckResourceAttr(resourceName1, "description", "First context"),
					resource.TestCheckResourceAttr(resourceName1, "order", "0"),
					resource.TestCheckResourceAttr(resourceName1, "action", "permit"),
					resource.TestCheckResourceAttr(resourceName2, "name", "ctx_2"),
					resource.TestCheckResourceAttr(resourceName2, "description", "Second context"),
					resource.TestCheckResourceAttr(resourceName2, "order", "1"),
					resource.TestCheckResourceAttr(resourceName2, "action", "deny"),
					resource.TestCheckResourceAttr(resourceName3, "name", "ctx_3"),
					resource.TestCheckResourceAttr(resourceName3, "description", "Third context"),
					resource.TestCheckResourceAttr(resourceName3, "order", "2"),
					resource.TestCheckResourceAttr(resourceName3, "action", "permit"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Deleting Middle Context (ctx_2)") },
				Config:    testAccRouteControlContextMiddleDeleted(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName1, "name", "ctx_1"),
					resource.TestCheckResourceAttr(resourceName1, "description", "First context"),
					resource.TestCheckResourceAttr(resourceName1, "order", "0"),
					resource.TestCheckResourceAttr(resourceName1, "action", "permit"),
					resource.TestCheckResourceAttr(resourceName3, "name", "ctx_3"),
					resource.TestCheckResourceAttr(resourceName3, "description", "Third context"),
					resource.TestCheckResourceAttr(resourceName3, "order", "2"),
					resource.TestCheckResourceAttr(resourceName3, "action", "permit"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Updating Remaining Contexts After Middle Deletion") },
				Config:    testAccRouteControlContextUpdateAfterDeletion(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName1, "name", "ctx_1"),
					resource.TestCheckResourceAttr(resourceName1, "description", "First context updated"),
					resource.TestCheckResourceAttr(resourceName1, "order", "0"),
					resource.TestCheckResourceAttr(resourceName1, "action", "deny"),
					resource.TestCheckResourceAttr(resourceName3, "name", "ctx_3"),
					resource.TestCheckResourceAttr(resourceName3, "description", "Third context updated"),
					resource.TestCheckResourceAttr(resourceName3, "order", "5"),
					resource.TestCheckResourceAttr(resourceName3, "action", "deny"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Re-adding New Context (ctx_4) Alongside Existing Contexts") },
				Config:    testAccRouteControlContextReaddNew(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName1, "name", "ctx_1"),
					resource.TestCheckResourceAttr(resourceName1, "description", "First context updated"),
					resource.TestCheckResourceAttr(resourceName1, "order", "0"),
					resource.TestCheckResourceAttr(resourceName1, "action", "deny"),
					resource.TestCheckResourceAttr(resourceName3, "name", "ctx_3"),
					resource.TestCheckResourceAttr(resourceName3, "description", "Third context updated"),
					resource.TestCheckResourceAttr(resourceName3, "order", "5"),
					resource.TestCheckResourceAttr(resourceName3, "action", "deny"),
					resource.TestCheckResourceAttr(resourceName4, "name", "ctx_4"),
					resource.TestCheckResourceAttr(resourceName4, "description", "Fourth context new"),
					resource.TestCheckResourceAttr(resourceName4, "order", "3"),
					resource.TestCheckResourceAttr(resourceName4, "action", "permit"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Importing Context (ctx_3) After Index Shifts") },
				ResourceName:      resourceName3,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig: func() { fmt.Println("Test: Prepare Route Map Policy Context Template Replacement") },
				Config:    testAccMSORouteMapPolicyContextConfigTemplateChangeSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "route_map_policy_uuid", "mso_tenant_policies_route_map_policy_route_control.second_policy", "uuid"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_route_map_policy_route_control.replacement_policy", "uuid"),
					captureParent,
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Replace Route Map Policy Context in Another Template") },
				Config:    testAccMSORouteMapPolicyContextConfigTemplateChangeDestination(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant_replacement", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "route_map_policy_uuid", "mso_tenant_policies_route_map_policy_route_control.replacement_policy", "uuid"),
					resource.TestCheckResourceAttr(resourceName, "name", "ctx_renamed"),
					resource.TestCheckResourceAttr(resourceName, "description", "Renamed context"),
					resource.TestCheckResourceAttr(resourceName, "order", "3"),
					resource.TestCheckResourceAttr(resourceName, "action", "deny"),
					resource.TestCheckResourceAttr(resourceName, "set_rule_uuid", ""),
					resource.TestCheckResourceAttr(resourceName, "match_rule_uuids.#", "0"),
					func(state *terraform.State) error {
						context := state.RootModule().Resources[resourceName]
						if context.Primary.ID == contextID {
							return fmt.Errorf("moving the context to another template retained its previous ID")
						}
						response, err := testAccPreCheck(t).GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
						if err != nil {
							return err
						}
						policyIndex, err := GetPolicyIndexByKeyAndValue(response, "uuid", parentUUID, "tenantPolicyTemplate", "template", "routeMapPolicies")
						if err != nil {
							return err
						}
						policy := response.S("tenantPolicyTemplate", "template", "routeMapPolicies").Index(policyIndex)
						if _, err := GetPolicyByName(policy, "ctx_renamed", "contexts"); !isPolicyNotFound(err) {
							if err != nil {
								return err
							}
							return fmt.Errorf("context ctx_renamed still exists in the original template")
						}
						return nil
					},
				),
			},
		},
		CheckDestroy: testCheckRouteMapPolicyContextDestroy,
	})
}

func testAccMSORouteMapPolicyContextConfig_Create() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control_context" "test_context" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = "ctx_1"
}
`, testAccMSORouteMapPolicyContextConfigBase())
}

func testAccMSORouteMapPolicyContextConfig_UpdateAttrs() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control_context" "test_context" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = "ctx_1"
  description           = "Updated description"
  order                 = 2
  action                = "deny"
}
`, testAccMSORouteMapPolicyContextConfigBase())
}

func testAccMSORouteMapPolicyContextConfig_WithRules() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control_context" "test_context" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = "ctx_1"
  description           = "Context with rules"
  order                 = 9
  action                = "deny"
  set_rule_uuid         = mso_tenant_policies_set_rule_policy.test_set_rule.uuid
  match_rule_uuids      = [mso_tenant_policies_match_rule_policy.test_match_rule_1.uuid, mso_tenant_policies_match_rule_policy.test_match_rule_2.uuid]
}
`, testAccMSORouteMapPolicyContextConfigBase())
}

func testAccMSORouteMapPolicyContextConfigReorderedRules() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control_context" "test_context" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = "ctx_1"
  description           = "Context with rules"
  order                 = 9
  action                = "deny"
  set_rule_uuid         = mso_tenant_policies_set_rule_policy.test_set_rule.uuid
  match_rule_uuids      = [mso_tenant_policies_match_rule_policy.test_match_rule_2.uuid, mso_tenant_policies_match_rule_policy.test_match_rule_1.uuid]
}
`, testAccMSORouteMapPolicyContextConfigBase())
}

func testAccMSORouteMapPolicyContextConfigOneMatchRule() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control_context" "test_context" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = "ctx_1"
  description           = "Context with rules"
  order                 = 9
  action                = "deny"
  set_rule_uuid         = mso_tenant_policies_set_rule_policy.test_set_rule.uuid
  match_rule_uuids      = [mso_tenant_policies_match_rule_policy.test_match_rule_2.uuid]
}
`, testAccMSORouteMapPolicyContextConfigBase())
}

func testAccMSORouteMapPolicyContextConfigNullValues() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control_context" "test_context" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = "ctx_1"
  description           = null
  order                 = null
  action                = null
  set_rule_uuid         = null
  match_rule_uuids      = null
}
`, testAccMSORouteMapPolicyContextConfigBase())
}

func testAccMSORouteMapPolicyContextConfigClearRules() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control_context" "test_context" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = "ctx_1"
  description           = ""
  order                 = 0
  action                = "permit"
  set_rule_uuid         = ""
  match_rule_uuids      = []
}
`, testAccMSORouteMapPolicyContextConfigBase())
}

func testAccMSORouteMapPolicyContextConfigRename() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control_context" "test_context" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = "ctx_renamed"
  description           = "Renamed context"
  order                 = 3
  action                = "deny"
}
`, testAccMSORouteMapPolicyContextConfigBase())
}

func testAccMSORouteMapPolicyContextConfigChangeParent() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control" "second_policy" {
  template_id = mso_template.template_tenant.id
  name        = "second_policy"
  description = "Another parent policy"
  depends_on  = [mso_tenant_policies_route_map_policy_route_control.test_policy]
}
resource "mso_tenant_policies_route_map_policy_route_control_context" "test_context" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.second_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.second_policy.uuid
  name                  = "ctx_renamed"
  description           = "Renamed context"
  order                 = 3
  action                = "deny"
}
`, testAccMSORouteMapPolicyContextConfigBase())
}

func testAccMSORouteMapPolicyContextConfigRenameParent() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control" "second_policy" {
  template_id = mso_template.template_tenant.id
  name        = "second_policy_renamed"
  description = "Another parent policy"
  depends_on  = [mso_tenant_policies_route_map_policy_route_control.test_policy]
}
resource "mso_tenant_policies_route_map_policy_route_control_context" "test_context" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.second_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.second_policy.uuid
  name                  = "ctx_renamed"
  description           = "Renamed context"
  order                 = 3
  action                = "deny"

  # Fail if renaming the parent would replace this context.
  lifecycle {
    prevent_destroy = true
  }
}
`, testAccMSORouteMapPolicyContextConfigBase())
}

func testAccRouteControlContextMultipleCreate() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control" "second_policy" {
  template_id = mso_template.template_tenant.id
  name        = "second_policy_renamed"
  description = "Another parent policy"
  depends_on  = [mso_tenant_policies_route_map_policy_route_control.test_policy]
}
resource "mso_tenant_policies_route_map_policy_route_control_context" "test_context" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.second_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.second_policy.uuid
  name                  = "ctx_renamed"
  description           = "Renamed context"
  order                 = 3
  action                = "deny"
}
resource "mso_tenant_policies_route_map_policy_route_control_context" "ctx_1" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = "ctx_1"
  depends_on            = [mso_tenant_policies_route_map_policy_route_control_context.test_context]
  description           = "First context"
  order                 = 0
  action                = "permit"
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "ctx_2" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = "ctx_2"
  depends_on            = [mso_tenant_policies_route_map_policy_route_control_context.ctx_1]
  description           = "Second context"
  order                 = 1
  action                = "deny"
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "ctx_3" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = "ctx_3"
  depends_on            = [mso_tenant_policies_route_map_policy_route_control_context.ctx_2]
  description           = "Third context"
  order                 = 2
  action                = "permit"
}
`, testAccMSORouteMapPolicyContextConfigBase())
}

func testAccRouteControlContextMiddleDeleted() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control" "second_policy" {
  template_id = mso_template.template_tenant.id
  name        = "second_policy_renamed"
  description = "Another parent policy"
  depends_on  = [mso_tenant_policies_route_map_policy_route_control.test_policy]
}
resource "mso_tenant_policies_route_map_policy_route_control_context" "test_context" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.second_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.second_policy.uuid
  name                  = "ctx_renamed"
  description           = "Renamed context"
  order                 = 3
  action                = "deny"
}
resource "mso_tenant_policies_route_map_policy_route_control_context" "ctx_1" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = "ctx_1"
  depends_on            = [mso_tenant_policies_route_map_policy_route_control_context.test_context]
  description           = "First context"
  order                 = 0
  action                = "permit"
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "ctx_3" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = "ctx_3"
  depends_on            = [mso_tenant_policies_route_map_policy_route_control_context.ctx_1]
  description           = "Third context"
  order                 = 2
  action                = "permit"
}
`, testAccMSORouteMapPolicyContextConfigBase())
}

func testAccRouteControlContextUpdateAfterDeletion() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control" "second_policy" {
  template_id = mso_template.template_tenant.id
  name        = "second_policy_renamed"
  description = "Another parent policy"
  depends_on  = [mso_tenant_policies_route_map_policy_route_control.test_policy]
}
resource "mso_tenant_policies_route_map_policy_route_control_context" "test_context" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.second_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.second_policy.uuid
  name                  = "ctx_renamed"
  description           = "Renamed context"
  order                 = 3
  action                = "deny"
}
resource "mso_tenant_policies_route_map_policy_route_control_context" "ctx_1" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = "ctx_1"
  depends_on            = [mso_tenant_policies_route_map_policy_route_control_context.test_context]
  description           = "First context updated"
  order                 = 0
  action                = "deny"
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "ctx_3" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = "ctx_3"
  depends_on            = [mso_tenant_policies_route_map_policy_route_control_context.ctx_1]
  description           = "Third context updated"
  order                 = 5
  action                = "deny"
}
`, testAccMSORouteMapPolicyContextConfigBase())
}

func testAccRouteControlContextReaddNew() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control" "second_policy" {
  template_id = mso_template.template_tenant.id
  name        = "second_policy_renamed"
  description = "Another parent policy"
  depends_on  = [mso_tenant_policies_route_map_policy_route_control.test_policy]
}
resource "mso_tenant_policies_route_map_policy_route_control_context" "test_context" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.second_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.second_policy.uuid
  name                  = "ctx_renamed"
  description           = "Renamed context"
  order                 = 3
  action                = "deny"
}
resource "mso_tenant_policies_route_map_policy_route_control_context" "ctx_1" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = "ctx_1"
  depends_on            = [mso_tenant_policies_route_map_policy_route_control_context.test_context]
  description           = "First context updated"
  order                 = 0
  action                = "deny"
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "ctx_3" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = "ctx_3"
  depends_on            = [mso_tenant_policies_route_map_policy_route_control_context.ctx_1]
  description           = "Third context updated"
  order                 = 5
  action                = "deny"
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "ctx_4" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = "ctx_4"
  depends_on            = [mso_tenant_policies_route_map_policy_route_control_context.ctx_3]
  description           = "Fourth context new"
  order                 = 3
  action                = "permit"
}
`, testAccMSORouteMapPolicyContextConfigBase())
}

func testAccMSORouteMapPolicyContextConfigTemplateChangeSource() string {
	return fmt.Sprintf(`%s
resource "mso_template" "template_tenant_replacement" {
  template_name = "test_template_tenant_replacement"
  template_type = "tenant"
  tenant_id     = mso_tenant.%s.id
  sites         = [data.mso_site.%s.id]
}
resource "mso_tenant_policies_route_map_policy_route_control" "replacement_policy" {
  template_id = mso_template.template_tenant_replacement.id
  name        = "replacement_policy"
  description = "Replacement parent policy"
  depends_on  = [mso_tenant_policies_route_map_policy_route_control.second_policy]
}
`, testAccRouteControlContextReaddNew(), msoTemplateTenantName, msoTemplateSiteName1)
}

func testAccMSORouteMapPolicyContextConfigTemplateChangeDestination() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control" "second_policy" {
  template_id = mso_template.template_tenant.id
  name        = "second_policy_renamed"
  description = "Another parent policy"
  depends_on  = [mso_tenant_policies_route_map_policy_route_control.test_policy]
}
resource "mso_tenant_policies_route_map_policy_route_control_context" "test_context" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.replacement_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.replacement_policy.uuid
  name                  = "ctx_renamed"
  description           = "Renamed context"
  order                 = 3
  action                = "deny"
}
resource "mso_tenant_policies_route_map_policy_route_control_context" "ctx_1" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = "ctx_1"
  depends_on            = [mso_tenant_policies_route_map_policy_route_control_context.test_context]
  description           = "First context updated"
  order                 = 0
  action                = "deny"
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "ctx_3" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = "ctx_3"
  depends_on            = [mso_tenant_policies_route_map_policy_route_control_context.ctx_1]
  description           = "Third context updated"
  order                 = 5
  action                = "deny"
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "ctx_4" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = "ctx_4"
  depends_on            = [mso_tenant_policies_route_map_policy_route_control_context.ctx_3]
  description           = "Fourth context new"
  order                 = 3
  action                = "permit"
}
resource "mso_template" "template_tenant_replacement" {
  template_name = "test_template_tenant_replacement"
  template_type = "tenant"
  tenant_id     = mso_tenant.%s.id
  sites         = [data.mso_site.%s.id]
}
resource "mso_tenant_policies_route_map_policy_route_control" "replacement_policy" {
  template_id = mso_template.template_tenant_replacement.id
  name        = "replacement_policy"
  description = "Replacement parent policy"
  depends_on  = [mso_tenant_policies_route_map_policy_route_control.second_policy]
}
`, testAccMSORouteMapPolicyContextConfigBase(), msoTemplateTenantName, msoTemplateSiteName1)
}

func TestAccMSORouteMapPolicyContextResourceErrors(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVersionCheck(t, "5.1")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Reject Route Map Policy Context without Template ID") },
				Config:      testAccMSORouteMapPolicyContextConfigMissingTemplate(),
				ExpectError: regexp.MustCompile(`The argument "template_id" is required`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Route Map Policy Context without Parent Policy UUID") },
				Config:      testAccMSORouteMapPolicyContextConfigMissingParent(),
				ExpectError: regexp.MustCompile(`The argument "route_map_policy_uuid" is required`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Route Map Policy Context without Name") },
				Config:      testAccMSORouteMapPolicyContextConfigMissingName(),
				ExpectError: regexp.MustCompile(`The argument "name" is required`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Route Map Policy Context with Empty Name") },
				Config:      testAccMSORouteMapPolicyContextConfigName(""),
				ExpectError: regexp.MustCompile(`expected length of name to be in the range \(1 - 64\)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Route Map Policy Context with Name Longer than 64 Characters") },
				Config:      testAccMSORouteMapPolicyContextConfigName(strings.Repeat("x", 65)),
				ExpectError: regexp.MustCompile(`expected length of name to be in the range \(1 - 64\)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Route Map Policy Context with Order Below Minimum") },
				Config:      testAccMSORouteMapPolicyContextConfigOrder(-1),
				ExpectError: regexp.MustCompile(`expected order to be in the range \(0 - 9\)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Route Map Policy Context with Order Above Maximum") },
				Config:      testAccMSORouteMapPolicyContextConfigOrder(10),
				ExpectError: regexp.MustCompile(`expected order to be in the range \(0 - 9\)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Route Map Policy Context with Invalid Action") },
				Config:      testAccMSORouteMapPolicyContextConfigInvalidAction(),
				ExpectError: regexp.MustCompile(`expected action to be one of`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Route Map Policy Context with Empty Match Rule UUID") },
				Config:      testAccMSORouteMapPolicyContextConfigEmptyMatchRule(),
				ExpectError: regexp.MustCompile(`(?s)expected "match_rule_uuids\.[^"]+" to not be an empty string`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create Route Map Policy Context with Valid Configuration") },
				Config:    testAccMSORouteMapPolicyContextConfigValid(),
				Check:     resource.TestCheckResourceAttrSet("mso_tenant_policies_route_map_policy_route_control_context.test_context", "id"),
			},
		},
		CheckDestroy: testCheckRouteMapPolicyContextDestroy,
	})
}

func testCheckRouteMapPolicyContextDestroy(s *terraform.State) error {
	msoClient := testAccPreCheck(nil)
	summaries, err := msoClient.GetViaURL("api/v1/templates/summaries")
	if err != nil {
		return fmt.Errorf("error fetching template summaries during Route Map Policy Context cleanup check: %w", err)
	}
	templateIDs := make(map[string]bool)
	if summaries.Data() != nil {
		templates, err := summaries.Children()
		if err != nil {
			return fmt.Errorf("error reading template summaries during Route Map Policy Context cleanup check: %w", err)
		}
		for _, template := range templates {
			id, ok := template.S("templateId").Data().(string)
			if !ok {
				return fmt.Errorf("template summary has no valid templateId")
			}
			templateIDs[id] = true
		}
	}
	for name, rs := range s.RootModule().Resources {
		if rs.Type == "mso_tenant_policies_route_map_policy_route_control_context" {
			templateID := rs.Primary.Attributes["template_id"]
			parentUUID := rs.Primary.Attributes["route_map_policy_uuid"]
			contextName := rs.Primary.Attributes["name"]
			if !templateIDs[templateID] {
				continue
			}
			response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
			if err != nil {
				return fmt.Errorf("error fetching template for resource '%s' with ID '%s': %w", name, rs.Primary.ID, err)
			}
			parentIndex, err := GetPolicyIndexByKeyAndValue(response, "uuid", parentUUID, "tenantPolicyTemplate", "template", "routeMapPolicies")
			if isPolicyNotFound(err) {
				continue
			}
			if err != nil {
				return err
			}
			parent := response.S("tenantPolicyTemplate", "template", "routeMapPolicies").Index(parentIndex)
			_, err = GetPolicyByName(parent, contextName, "contexts")
			if isPolicyNotFound(err) {
				continue
			}
			if err != nil {
				return err
			}
			return fmt.Errorf("terraform destroy was unsuccessful. The Route Map Policy Context '%s' with name '%s' still exists under parent UUID '%s'", name, contextName, parentUUID)
		}
	}
	return nil
}

func testAccMSORouteMapPolicyContextConfigMissingParent() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control_context" "test_context" {
  template_id = mso_template.template_tenant.id
  name        = "ctx_1"
}
`, testAccMSORouteMapPolicyConfigCreate())
}

func testAccMSORouteMapPolicyContextConfigMissingName() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control_context" "test_context" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
}
`, testAccMSORouteMapPolicyConfigCreate())
}

func testAccMSORouteMapPolicyContextConfigName(name string) string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control_context" "test_context" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = %q
}
`, testAccMSORouteMapPolicyConfigCreate(), name)
}

func testAccMSORouteMapPolicyContextConfigOrder(order int) string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control_context" "test_context" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = "ctx_1"
  order                 = %d
}
`, testAccMSORouteMapPolicyConfigCreate(), order)
}

func testAccMSORouteMapPolicyContextConfigInvalidAction() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control_context" "test_context" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = "ctx_1"
  action                = "invalid"
}
`, testAccMSORouteMapPolicyConfigCreate())
}

func testAccMSORouteMapPolicyContextConfigEmptyMatchRule() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control_context" "test_context" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = "ctx_1"
  match_rule_uuids      = [""]
}
`, testAccMSORouteMapPolicyConfigCreate())
}

func testAccMSORouteMapPolicyContextConfigMissingTemplate() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control_context" "test_context" {
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = "ctx_1"
}
`, testAccMSORouteMapPolicyConfigCreate())
}

func testAccMSORouteMapPolicyContextConfigValid() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control_context" "test_context" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = "ctx_1"
}
`, testAccMSORouteMapPolicyConfigCreate())
}

func testAccMSORouteMapPolicyContextConfigBase() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_set_rule_policy" "test_set_rule" {
  template_id = mso_template.template_tenant.id
  name        = "tf_test_context_set_rule"
  set_metric  = 10

  set_route_tag {
    value = 0
  }
  set_weight {
    value = 0
  }
  set_preference {
    value = 0
  }
}
resource "mso_tenant_policies_match_rule_policy" "test_match_rule_1" {
  template_id = mso_template.template_tenant.id
  name        = "tf_test_context_match_rule_1"
  match_prefix {
    prefix      = "10.0.0.0/8"
    aggregate   = true
    from_prefix = 9
    to_prefix   = 32
  }
  depends_on = [mso_tenant_policies_set_rule_policy.test_set_rule]
}
resource "mso_tenant_policies_match_rule_policy" "test_match_rule_2" {
  template_id = mso_template.template_tenant.id
  name        = "tf_test_context_match_rule_2"
  match_prefix {
    prefix      = "192.168.0.0/16"
    aggregate   = true
    from_prefix = 17
    to_prefix   = 32
  }
  depends_on = [mso_tenant_policies_match_rule_policy.test_match_rule_1]
}
resource "mso_tenant_policies_route_map_policy_route_control" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_policy"
  description = "Initial description"
  depends_on  = [mso_tenant_policies_match_rule_policy.test_match_rule_2]
}
`, testAccMSOTemplateResourceTenantConfig())
}
