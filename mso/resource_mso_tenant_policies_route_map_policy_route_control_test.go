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

func TestAccMSORouteMapPolicyResource(t *testing.T) {
	resourceName := "mso_tenant_policies_route_map_policy_route_control.test_policy"
	var templateID, uuid string
	var originalTemplateID, originalUUID string
	// Save IDs from state for later API calls and UUID comparisons across test steps.
	capture := func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok || rs.Primary == nil {
			return fmt.Errorf("resource %s not found in state", resourceName)
		}
		templateID = rs.Primary.Attributes["template_id"]
		uuid = rs.Primary.Attributes["uuid"]
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
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create Route Map Policy without Description") },
				Config:    testAccMSORouteMapPolicyConfigWithoutDescription(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_policy"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "match_type", "routing_policy_only"),
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Route Map Policy - Set Description and Match Type") },
				Config:    testAccMSORouteMapPolicyConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Initial description"),
					resource.TestCheckResourceAttr(resourceName, "match_type", "prefix_and_routing_policy"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Route Map Policy - Null Description") },
				Config:    testAccMSORouteMapPolicyConfigNullDescription(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "match_type", "prefix_and_routing_policy"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Route Map Policy - Restore Description") },
				Config:    testAccMSORouteMapPolicyConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Initial description"),
					resource.TestCheckResourceAttr(resourceName, "match_type", "prefix_and_routing_policy"),
					capture,
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update Route Map Policy after Out-of-Band Description Removal")
					c := testAccPreCheck(t)
					// Resolve the current API array index before removing the description outside Terraform.
					response, err := c.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
					if err != nil {
						t.Fatal(err)
					}
					index, err := GetPolicyIndexByKeyAndValue(response, "uuid", uuid, "tenantPolicyTemplate", "template", "routeMapPolicies")
					if err != nil {
						t.Fatal(err)
					}
					payload := models.GetRemovePatchPayload(fmt.Sprintf("/tenantPolicyTemplate/template/routeMapPolicies/%d/description", index))
					if _, err := c.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateID), payload); err != nil {
						t.Fatal(err)
					}
				},
				Config: testAccMSORouteMapPolicyConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_policy"),
					resource.TestCheckResourceAttr(resourceName, "description", "Initial description"),
					resource.TestCheckResourceAttr(resourceName, "match_type", "prefix_and_routing_policy"),
					func(s *terraform.State) error {
						if s.RootModule().Resources[resourceName].Primary.Attributes["uuid"] != uuid {
							return fmt.Errorf("restoring the Route Map Policy description changed its UUID")
						}
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update Route Map Policy after Out-of-Band Match Type Change")
					c := testAccPreCheck(t)
					// Resolve the current API array index before changing the match type outside Terraform.
					response, err := c.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
					if err != nil {
						t.Fatal(err)
					}
					index, err := GetPolicyIndexByKeyAndValue(response, "uuid", uuid, "tenantPolicyTemplate", "template", "routeMapPolicies")
					if err != nil {
						t.Fatal(err)
					}
					payload := container.New()
					payload.Array()
					if err := addPatchPayloadToContainer(payload, "add", fmt.Sprintf("/tenantPolicyTemplate/template/routeMapPolicies/%d/routeMapType", index), "global"); err != nil {
						t.Fatal(err)
					}
					if err := doPatchRequest(c, fmt.Sprintf("api/v1/templates/%s", templateID), payload); err != nil {
						t.Fatal(err)
					}
				},
				Config: testAccMSORouteMapPolicyConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "match_type", "prefix_and_routing_policy"),
					func(s *terraform.State) error {
						if s.RootModule().Resources[resourceName].Primary.Attributes["uuid"] != uuid {
							return fmt.Errorf("restoring the Route Map Policy match type changed its UUID")
						}
						return nil
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Route Map Policy - Routing Policy Only Match Type") },
				Config:    testAccMSORouteMapPolicyConfigMatchType("routing_policy_only"),
				Check:     resource.TestCheckResourceAttr(resourceName, "match_type", "routing_policy_only"),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Route Map Policy - Prefix and Routing Policy Match Type") },
				Config:    testAccMSORouteMapPolicyConfigCreate(),
				Check:     resource.TestCheckResourceAttr(resourceName, "match_type", "prefix_and_routing_policy"),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Route Map Policy - Null Match Type") },
				Config:    testAccMSORouteMapPolicyConfigNullMatchType(),
				Check:     resource.TestCheckResourceAttr(resourceName, "match_type", "routing_policy_only"),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Route Map Policy - Restore Match Type") },
				Config:    testAccMSORouteMapPolicyConfigCreate(),
				Check:     resource.TestCheckResourceAttr(resourceName, "match_type", "prefix_and_routing_policy"),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Route Map Policy - Remove Description and Reset Match Type") },
				Config:    testAccMSORouteMapPolicyConfigWithoutDescription(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "match_type", "routing_policy_only"),
					capture,
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate Route Map Policy after out-of-band deletion")
					if err := testAccDeletePolicyOutOfBand(testAccPreCheck(t), templateID, uuid, "tenantPolicyTemplate", "template", "routeMapPolicies"); err != nil {
						t.Fatalf("delete %s out of band: %v", resourceName, err)
					}
				},
				Config: testAccMSORouteMapPolicyConfigWithoutDescription(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_policy"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "match_type", "routing_policy_only"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok || rs.Primary == nil {
							return fmt.Errorf("resource %s not found in state", resourceName)
						}
						if rs.Primary.Attributes["uuid"] == uuid {
							return fmt.Errorf("resource %s retained its UUID after out-of-band deletion", resourceName)
						}
						return nil
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Route Map Policy - Description") },
				Config:    testAccMSORouteMapPolicyConfigUpdateDescription(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
					resource.TestCheckResourceAttr(resourceName, "match_type", "prefix_and_routing_policy"),
					capture,
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Route Map Policy - Name") },
				Config:    testAccMSORouteMapPolicyConfigUpdateName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_policy_renamed"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
					resource.TestCheckResourceAttr(resourceName, "match_type", "prefix_and_routing_policy"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					func(s *terraform.State) error {
						rs := s.RootModule().Resources[resourceName]
						if rs.Primary.Attributes["uuid"] != uuid {
							return fmt.Errorf("renaming the Route Map Policy changed its UUID")
						}
						if rs.Primary.ID != fmt.Sprintf("templateId/%s/RouteMapPolicy/test_policy_renamed", templateID) {
							return fmt.Errorf("renaming the Route Map Policy did not update its ID")
						}
						response, err := testAccPreCheck(t).GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
						if err != nil {
							return err
						}
						if _, err := GetPolicyByName(response, "test_policy", "tenantPolicyTemplate", "template", "routeMapPolicies"); !isPolicyNotFound(err) {
							if err != nil {
								return err
							}
							return fmt.Errorf("the old Route Map Policy name still exists after renaming")
						}
						return nil
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Prepare Route Map Policy Template Replacement") },
				Config:    testAccMSORouteMapPolicyConfigTemplateChangeSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant", "id"),
					resource.TestCheckResourceAttrSet("mso_template.template_tenant_replacement", "id"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok || rs.Primary == nil {
							return fmt.Errorf("resource %s not found in state", resourceName)
						}
						originalTemplateID = rs.Primary.Attributes["template_id"]
						originalUUID = rs.Primary.Attributes["uuid"]
						if originalTemplateID == "" || originalUUID == "" {
							return fmt.Errorf("resource %s is missing its template ID or UUID", resourceName)
						}
						return nil
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Replace Route Map Policy - Template ID") },
				Config:    testAccMSORouteMapPolicyConfigTemplateChangeDestination(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "template_id", "mso_template.template_tenant_replacement", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "test_policy_renamed"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
					resource.TestCheckResourceAttr(resourceName, "match_type", "prefix_and_routing_policy"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok || rs.Primary == nil {
							return fmt.Errorf("resource %s not found in state", resourceName)
						}
						if rs.Primary.Attributes["uuid"] == originalUUID {
							return fmt.Errorf("resource %s retained its UUID after changing template ID", resourceName)
						}
						response, err := testAccPreCheck(t).GetViaURL(fmt.Sprintf("api/v1/templates/%s", originalTemplateID))
						if err != nil {
							return fmt.Errorf("read original template after replacement: %w", err)
						}
						if _, err := GetPolicyIndexByKeyAndValue(response, "uuid", originalUUID, "tenantPolicyTemplate", "template", "routeMapPolicies"); !isPolicyNotFound(err) {
							if err != nil {
								return fmt.Errorf("verify original policy was deleted during replacement: %w", err)
							}
							return fmt.Errorf("original policy with UUID %s still exists after replacement", originalUUID)
						}
						return nil
					},
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Route Map Policy") },
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import missing Route Map Policy") },
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok || rs.Primary == nil {
						return "", fmt.Errorf("resource %s not found in state", resourceName)
					}
					separator := strings.LastIndex(rs.Primary.ID, "/")
					if separator < 0 || separator == len(rs.Primary.ID)-1 {
						return "", fmt.Errorf("resource %s has an invalid import ID", resourceName)
					}
					return rs.Primary.ID[:separator+1] + "missing_route_map_policy", nil
				},
				ExpectError: regexp.MustCompile(`(?s)cannot import .*: resource\s+not\s+found`),
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments("mso_tenant_policies_route_map_policy_route_control", "tenantPolicyTemplate", "template", "routeMapPolicies"),
	})
}

func TestAccMSORouteMapPolicyResourceErrors(t *testing.T) {
	resourceName := "mso_tenant_policies_route_map_policy_route_control.test_policy"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVersionCheck(t, "5.1")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Reject Route Map Policy without template ID") },
				Config:      testAccMSORouteMapPolicyConfigErrorWithoutTemplateID(),
				ExpectError: regexp.MustCompile(`The argument "template_id" is required`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Route Map Policy without name") },
				Config:      testAccMSORouteMapPolicyConfigErrorWithoutName(),
				ExpectError: regexp.MustCompile(`The argument "name" is required`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject empty Route Map Policy name") },
				Config:      testAccMSORouteMapPolicyConfigErrorName(""),
				ExpectError: regexp.MustCompile(`expected length of name to be in the range \(1 - 64\)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Route Map Policy name longer than 64 characters") },
				Config:      testAccMSORouteMapPolicyConfigErrorName(strings.Repeat("a", 65)),
				ExpectError: regexp.MustCompile(`expected length of name to be in the range \(1 - 64\)`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Route Map Policy with invalid Match Type") },
				Config:      testAccMSORouteMapPolicyConfigMatchType("invalid"),
				ExpectError: regexp.MustCompile(`expected match_type to be one of`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject Route Map Policy with empty Match Type") },
				Config:      testAccMSORouteMapPolicyConfigMatchType(""),
				ExpectError: regexp.MustCompile(`expected match_type to be one of`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create Route Map Policy after configuration errors") },
				Config:    testAccMSORouteMapPolicyConfigCreate(),
				Check:     resource.TestCheckResourceAttrSet(resourceName, "uuid"),
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments("mso_tenant_policies_route_map_policy_route_control", "tenantPolicyTemplate", "template", "routeMapPolicies"),
	})
}

func testAccMSORouteMapPolicyConfigCreate() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_policy"
  description = "Initial description"
  match_type  = "prefix_and_routing_policy"
}
`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSORouteMapPolicyConfigWithoutDescription() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_policy"
}
`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSORouteMapPolicyConfigNullDescription() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_policy"
  description = null
  match_type  = "prefix_and_routing_policy"
}
`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSORouteMapPolicyConfigTemplateChangeSource() string {
	return fmt.Sprintf(`%s
resource "mso_template" "template_tenant_replacement" {
  template_name = "test_template_tenant_replacement"
  template_type = "tenant"
  tenant_id     = mso_tenant.%s.id
  sites         = [data.mso_site.%s.id]
}

resource "mso_tenant_policies_route_map_policy_route_control" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_policy_renamed"
  description = "Updated description"
  match_type  = "prefix_and_routing_policy"
}
`, testAccMSOTemplateResourceTenantConfig(), msoTemplateTenantName, msoTemplateSiteName1)
}

func testAccMSORouteMapPolicyConfigTemplateChangeDestination() string {
	return fmt.Sprintf(`%s
resource "mso_template" "template_tenant_replacement" {
  template_name = "test_template_tenant_replacement"
  template_type = "tenant"
  tenant_id     = mso_tenant.%s.id
  sites         = [data.mso_site.%s.id]
}

resource "mso_tenant_policies_route_map_policy_route_control" "test_policy" {
  template_id = mso_template.template_tenant_replacement.id
  name        = "test_policy_renamed"
  description = "Updated description"
  match_type  = "prefix_and_routing_policy"
}
`, testAccMSOTemplateResourceTenantConfig(), msoTemplateTenantName, msoTemplateSiteName1)
}

func testAccMSORouteMapPolicyConfigErrorWithoutTemplateID() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control" "test_policy" {
  name        = "test_policy"
  description = "Initial description"
}
`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSORouteMapPolicyConfigErrorWithoutName() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control" "test_policy" {
  template_id = mso_template.template_tenant.id
  description = "Initial description"
}
`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSORouteMapPolicyConfigErrorName(name string) string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = %q
  description = "Initial description"
}
`, testAccMSOTemplateResourceTenantConfig(), name)
}

func testAccMSORouteMapPolicyConfigUpdateDescription() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_policy"
  description = "Updated description"
  match_type  = "prefix_and_routing_policy"
}
`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSORouteMapPolicyConfigUpdateName() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_policy_renamed"
  description = "Updated description"
  match_type  = "prefix_and_routing_policy"
}
`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSORouteMapPolicyConfigMatchType(matchType string) string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_policy"
  description = "Initial description"
  match_type  = %q
}
`, testAccMSOTemplateResourceTenantConfig(), matchType)
}

func testAccMSORouteMapPolicyConfigNullMatchType() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_policy"
  description = "Initial description"
  match_type  = null
}
`, testAccMSOTemplateResourceTenantConfig())
}
