package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMSOSchemaTemplateBdDHCPPolicyResource(t *testing.T) {
	const relayPolicyResource = "mso_schema_template_bd_dhcp_policy.relay_policy"
	const relayPolicyWithoutOptionResource = "mso_schema_template_bd_dhcp_policy.relay_policy_without_option"
	const bdResource = "mso_schema_template_bd.dhcp_parent"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVersionCheck(t, "4.0")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Reject a standalone DHCP association with both relay identifiers") },
				Config:      testAccMSOSchemaTemplateBdDHCPPolicyConfigBothRelayIdentifiers(),
				ExpectError: regexp.MustCompile(`only one of .*name.*uuid.* can be specified`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Require a standalone DHCP association relay identifier") },
				Config:      testAccMSOSchemaTemplateBdDHCPPolicyConfigNoRelayIdentifier(),
				ExpectError: regexp.MustCompile(`one of .*name.*uuid.* must be specified`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject both DHCP Option Policy identifiers") },
				Config:      testAccMSOSchemaTemplateBdDHCPPolicyConfigBothOptionIdentifiers(),
				ExpectError: regexp.MustCompile(`conflicts with.*dhcp_option_policy_uuid`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create a standalone DHCP policy with an option policy") },
				Config:    testAccMSOSchemaTemplateBdDHCPPolicyCreateConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(relayPolicyResource, "name", msoTenantPoliciesDhcpRelayPolicyName),
					resource.TestCheckResourceAttr(relayPolicyResource, "dhcp_option_policy_name", msoTenantPoliciesDhcpOptionPolicyName),
				),
			},
			{
				PreConfig:          func() { fmt.Println("Test: Change the DHCP Option Policy association and detect external changes") },
				Config:             testAccMSOSchemaTemplateBdDHCPPolicySwitchOptionConfig(),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(relayPolicyResource, "name", msoTenantPoliciesDhcpRelayPolicyName),
					resource.TestCheckResourceAttr(relayPolicyResource, "dhcp_option_policy_name", msoTenantPoliciesDhcpOptionPolicyName2),
					testAccCheckMSOSchemaTemplateBdDHCPPolicyChangeOptionExternally(relayPolicyResource, msoTenantPoliciesDhcpOptionPolicyName),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add a second standalone DHCP policy without an option policy") },
				Config:    testAccMSOSchemaTemplateBdDHCPPolicyAddWithoutOptionConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(relayPolicyResource, "name", msoTenantPoliciesDhcpRelayPolicyName),
					resource.TestCheckResourceAttr(relayPolicyResource, "dhcp_option_policy_name", msoTenantPoliciesDhcpOptionPolicyName2),
					resource.TestCheckResourceAttr(relayPolicyWithoutOptionResource, "name", msoTenantPoliciesDhcpRelayPolicyName2),
					resource.TestCheckResourceAttr(relayPolicyWithoutOptionResource, "dhcp_option_policy_name", ""),
				),
			},
			{
				PreConfig:          func() { fmt.Println("Test: Detect an externally removed standalone DHCP policy association") },
				Config:             testAccMSOSchemaTemplateBdDHCPPolicyAddWithoutOptionConfig(),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMSOSchemaTemplateBdDHCPPolicyRemoveExternally(relayPolicyWithoutOptionResource),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Recreate an externally removed standalone DHCP policy association") },
				Config:    testAccMSOSchemaTemplateBdDHCPPolicyAddWithoutOptionConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(relayPolicyResource, "name", msoTenantPoliciesDhcpRelayPolicyName),
					resource.TestCheckResourceAttr(relayPolicyWithoutOptionResource, "name", msoTenantPoliciesDhcpRelayPolicyName2),
					testAccCheckMSOSchemaTemplateBdDHCPPolicyNames(bdResource, msoTenantPoliciesDhcpRelayPolicyName, msoTenantPoliciesDhcpRelayPolicyName2),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update the first standalone DHCP policy to remove its option policy") },
				Config:    testAccMSOSchemaTemplateBdDHCPPolicyUpdateConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(relayPolicyResource, "dhcp_option_policy_name", ""),
					resource.TestCheckResourceAttr(relayPolicyWithoutOptionResource, "name", msoTenantPoliciesDhcpRelayPolicyName2),
				),
			},
			{
				PreConfig:          func() { fmt.Println("Test: Verify standalone DHCP policies leave the unconfigured BD unchanged") },
				Config:             testAccMSOSchemaTemplateBdDHCPPolicyUpdateConfig(),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(bdResource, "dhcp_policies.#", "0"),
					resource.TestCheckResourceAttr(relayPolicyResource, "name", msoTenantPoliciesDhcpRelayPolicyName),
					resource.TestCheckResourceAttr(relayPolicyWithoutOptionResource, "name", msoTenantPoliciesDhcpRelayPolicyName2),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import BD while standalone DHCP policies are attached") },
				Config:       testAccMSOSchemaTemplateBdDHCPPolicyUpdateConfig(),
				ResourceName: bdResource,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					resource, ok := s.RootModule().Resources[bdResource]
					if !ok {
						return "", fmt.Errorf("BD resource not found in state")
					}
					attributes := resource.Primary.Attributes
					return fmt.Sprintf("%s/templates/%s/bds/%s", attributes["schema_id"], attributes["template_name"], attributes["name"]), nil
				},
				ImportStateVerify: true,
			},
			{
				PreConfig:          func() { fmt.Println("Test: Plan with one DHCP policy inline and another standalone") },
				Config:             testAccMSOSchemaTemplateBdDHCPPolicyMixedConfig(),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(relayPolicyResource, "name", msoTenantPoliciesDhcpRelayPolicyName),
					resource.TestCheckResourceAttr(relayPolicyWithoutOptionResource, "name", msoTenantPoliciesDhcpRelayPolicyName2),
					testAccCheckMSOSchemaTemplateBdDHCPPolicyNames(bdResource, msoTenantPoliciesDhcpRelayPolicyName, msoTenantPoliciesDhcpRelayPolicyName2),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import standalone DHCP policy") },
				Config:            testAccMSOSchemaTemplateBdDHCPPolicyUpdateConfig(),
				ResourceName:      relayPolicyResource,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig: func() { fmt.Println("Test: Delete one standalone DHCP policy and retain its sibling") },
				Config:    testAccMSOSchemaTemplateBdDHCPPolicyRemoveFirstConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(bdResource, "dhcp_policies.#", "0"),
					resource.TestCheckResourceAttr(relayPolicyWithoutOptionResource, "name", msoTenantPoliciesDhcpRelayPolicyName2),
					resource.TestCheckResourceAttr(relayPolicyWithoutOptionResource, "dhcp_option_policy_name", ""),
					testAccCheckMSOSchemaTemplateBdDHCPPolicyNames(bdResource, msoTenantPoliciesDhcpRelayPolicyName2),
				),
			},
			{
				PreConfig:          func() { fmt.Println("Test: Verify the remaining standalone DHCP policy has no diff") },
				Config:             testAccMSOSchemaTemplateBdDHCPPolicyRemoveFirstConfig(),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(bdResource, "dhcp_policies.#", "0"),
					resource.TestCheckResourceAttr(relayPolicyWithoutOptionResource, "name", msoTenantPoliciesDhcpRelayPolicyName2),
				),
			},
			{
				PreConfig:          func() { fmt.Println("Test: Detect out-of-band removal of the last standalone DHCP policy") },
				Config:             testAccMSOSchemaTemplateBdDHCPPolicyRemoveFirstConfig(),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMSOSchemaTemplateBdDHCPPolicyRemoveExternally(relayPolicyWithoutOptionResource),
					testAccCheckMSOSchemaTemplateBdDHCPPolicyNames(bdResource),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Recreate the last out-of-band removed DHCP policy") },
				Config:    testAccMSOSchemaTemplateBdDHCPPolicyRemoveFirstConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(relayPolicyWithoutOptionResource, "name", msoTenantPoliciesDhcpRelayPolicyName2),
					testAccCheckMSOSchemaTemplateBdDHCPPolicyNames(bdResource, msoTenantPoliciesDhcpRelayPolicyName2),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create standalone DHCP associations using UUID inputs") },
				Config:    testAccMSOSchemaTemplateBdDHCPPolicyUUIDCreateConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(relayPolicyResource, "name", msoTenantPoliciesDhcpRelayPolicyName),
					resource.TestCheckResourceAttrPair(relayPolicyResource, "uuid", "mso_tenant_policies_dhcp_relay_policy."+msoTenantPoliciesDhcpRelayPolicyName, "uuid"),
					resource.TestCheckResourceAttr(relayPolicyResource, "dhcp_option_policy_name", ""),
					resource.TestCheckResourceAttrPair(relayPolicyResource, "dhcp_option_policy_uuid", "mso_tenant_policies_dhcp_option_policy."+msoTenantPoliciesDhcpOptionPolicyName, "uuid"),
					resource.TestCheckResourceAttr(relayPolicyWithoutOptionResource, "name", msoTenantPoliciesDhcpRelayPolicyName2),
					resource.TestCheckResourceAttrPair(relayPolicyWithoutOptionResource, "uuid", "mso_tenant_policies_dhcp_relay_policy."+msoTenantPoliciesDhcpRelayPolicyName2, "uuid"),
					resource.TestCheckResourceAttr(relayPolicyWithoutOptionResource, "dhcp_option_policy_name", ""),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update the DHCP Option Policy using its UUID") },
				Config:    testAccMSOSchemaTemplateBdDHCPPolicyUUIDUpdateOptionConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(relayPolicyResource, "name", msoTenantPoliciesDhcpRelayPolicyName),
					resource.TestCheckResourceAttr(relayPolicyResource, "dhcp_option_policy_name", ""),
					resource.TestCheckResourceAttrPair(relayPolicyResource, "dhcp_option_policy_uuid", "mso_tenant_policies_dhcp_option_policy."+msoTenantPoliciesDhcpOptionPolicyName2, "uuid"),
					resource.TestCheckResourceAttr(relayPolicyWithoutOptionResource, "dhcp_option_policy_name", ""),
				),
			},
			{
				PreConfig:          func() { fmt.Println("Test: Verify UUID-managed DHCP associations have no diff") },
				Config:             testAccMSOSchemaTemplateBdDHCPPolicyUUIDUpdateOptionConfig(),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testAccCheckMSOSchemaTemplateBdDHCPPolicyChangeOptionExternally(resourceName, optionPolicyName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("DHCP policy resource not found in state: %s", resourceName)
		}
		attributes := resource.Primary.Attributes
		schemaID := attributes["schema_id"]
		templateName := attributes["template_name"]
		bdName := attributes["bd_name"]
		relayUUID := attributes["uuid"]

		response, err := msoClientTest.ReadTemplateBDDHCPPolicy(schemaID)
		if err != nil {
			return fmt.Errorf("failed to read Bridge Domain DHCP policies: %w", err)
		}
		template, err := response.S("templates").SearchInObjectList(func(item *container.Container) bool {
			return models.G(item, "name") == templateName
		})
		if err != nil {
			return err
		}
		bd, err := template.S("bds").SearchInObjectList(func(item *container.Container) bool {
			return models.G(item, "name") == bdName
		})
		if err != nil {
			return err
		}
		count, err := bd.ArrayCount("dhcpLabels")
		if err != nil {
			return err
		}
		policies, err := getDHCPPolicesNameByRef(count, schemaID, templateName, bd, msoClientTest)
		if err != nil {
			return err
		}
		policyIndex := -1
		for index, item := range policies {
			if item.(map[string]interface{})["uuid"] == relayUUID {
				policyIndex = index
				break
			}
		}
		if policyIndex == -1 {
			return fmt.Errorf("DHCP Relay Policy UUID %q not found on Bridge Domain %q", relayUUID, bdName)
		}

		tenantID, err := msoClientTest.GetTenantIDFromSchemaTemplate(schemaID, templateName)
		if err != nil {
			return err
		}
		optionUUID, err := getDHCPPolicyUUIDByName(msoClientTest, tenantID, "dhcpOption", optionPolicyName)
		if err != nil {
			return err
		}
		payload := container.New()
		payload.Array()
		value := map[string]interface{}{
			"ref":             relayUUID,
			"dhcpOptionLabel": map[string]interface{}{"ref": optionUUID},
		}
		path := fmt.Sprintf("/templates/%s/bds/%s/dhcpLabels/%d", templateName, bdName, policyIndex)
		if err := addPatchPayloadToContainer(payload, "replace", path, value); err != nil {
			return err
		}
		if err := doPatchRequest(msoClientTest, "api/v1/schemas/"+schemaID, payload); err != nil {
			return fmt.Errorf("failed to change DHCP Option Policy outside Terraform: %w", err)
		}
		return nil
	}
}

func testAccCheckMSOSchemaTemplateBdDHCPPolicyRemoveExternally(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("DHCP policy resource not found in state: %s", resourceName)
		}
		attributes := resource.Primary.Attributes
		schemaID := attributes["schema_id"]
		templateName := attributes["template_name"]
		bdName := attributes["bd_name"]
		relayUUID := attributes["uuid"]

		response, err := msoClientTest.ReadTemplateBDDHCPPolicy(schemaID)
		if err != nil {
			return fmt.Errorf("failed to read Bridge Domain DHCP policies: %w", err)
		}
		template, err := response.S("templates").SearchInObjectList(func(item *container.Container) bool {
			return models.G(item, "name") == templateName
		})
		if err != nil {
			return err
		}
		bd, err := template.S("bds").SearchInObjectList(func(item *container.Container) bool {
			return models.G(item, "name") == bdName
		})
		if err != nil {
			return err
		}
		count, err := bd.ArrayCount("dhcpLabels")
		if err != nil {
			return err
		}
		policyIndex := -1
		for index := 0; index < count; index++ {
			label, err := bd.ArrayElement(index, "dhcpLabels")
			if err != nil {
				return err
			}
			if models.StripQuotes(label.S("ref").String()) == relayUUID {
				policyIndex = index
				break
			}
		}
		if policyIndex == -1 {
			return fmt.Errorf("DHCP Relay Policy UUID %q not found on Bridge Domain %q", relayUUID, bdName)
		}

		path := fmt.Sprintf("/templates/%s/bds/%s/dhcpLabels/%d", templateName, bdName, policyIndex)
		if _, err := msoClientTest.PatchbyID("api/v1/schemas/"+schemaID, models.GetRemovePatchPayload(path)); err != nil {
			return fmt.Errorf("failed to remove DHCP policy association outside Terraform: %w", err)
		}
		return nil
	}
}

func testAccCheckMSOSchemaTemplateBdDHCPPolicyNames(bdResourceName string, expectedNames ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		bdResource, ok := s.RootModule().Resources[bdResourceName]
		if !ok {
			return fmt.Errorf("BD resource not found in state: %s", bdResourceName)
		}

		schemaID := bdResource.Primary.Attributes["schema_id"]
		templateName := bdResource.Primary.Attributes["template_name"]
		bdName := bdResource.Primary.Attributes["name"]
		response, err := msoClientTest.ReadTemplateBDDHCPPolicy(schemaID)
		if err != nil {
			return fmt.Errorf("failed to read Bridge Domain DHCP policies from schema %s: %w", schemaID, err)
		}
		template, err := response.S("templates").SearchInObjectList(func(item *container.Container) bool {
			return models.G(item, "name") == templateName
		})
		if err != nil {
			return fmt.Errorf("template %q not found in schema %s: %w", templateName, schemaID, err)
		}
		bd, err := template.S("bds").SearchInObjectList(func(item *container.Container) bool {
			return models.G(item, "name") == bdName
		})
		if err != nil {
			return fmt.Errorf("Bridge Domain %q not found in template %q: %w", bdName, templateName, err)
		}
		count := 0
		if bd.Exists("dhcpLabels") {
			count, err = bd.ArrayCount("dhcpLabels")
			if err != nil {
				return fmt.Errorf("failed to count DHCP policies on Bridge Domain %q: %w", bdName, err)
			}
		}
		if count != len(expectedNames) {
			return fmt.Errorf("expected DHCP policies %v on Bridge Domain %q, found %d associations", expectedNames, bdName, count)
		}

		policies, err := getDHCPPolicesNameByRef(count, schemaID, templateName, bd, msoClientTest)
		if err != nil {
			return fmt.Errorf("failed to resolve DHCP policy names on Bridge Domain %q: %w", bdName, err)
		}
		for _, expectedName := range expectedNames {
			found := false
			for _, policy := range policies {
				if policy.(map[string]interface{})["name"] == expectedName {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("expected DHCP policy %q to remain attached to Bridge Domain %q", expectedName, bdName)
			}
		}
		return nil
	}
}

func testAccMSOSchemaTemplateBdDHCPPolicyPrerequisiteConfig() string {
	return fmt.Sprintf("%s%s%s", testSiteConfigAnsibleTest(), testTenantConfig(), testSchemaConfig()) +
		testAccMSOSchemaTemplateBdDhcpPrerequisiteConfig() +
		testAccMSOSchemaTemplateBdDhcpRelayPolicyConfig() +
		testAccMSOSchemaTemplateBdDhcpOptionPolicyConfig()
}

func testAccMSOSchemaTemplateBdDHCPPolicyConfigBothRelayIdentifiers() string {
	return `resource "mso_schema_template_bd_dhcp_policy" "invalid" {
  schema_id     = "schema"
  template_name = "template"
  bd_name       = "bd"
  name          = "relay"
  uuid          = "00000000-0000-4000-8000-000000000001"
}`
}

func testAccMSOSchemaTemplateBdDHCPPolicyConfigNoRelayIdentifier() string {
	return `resource "mso_schema_template_bd_dhcp_policy" "invalid" {
  schema_id     = "schema"
  template_name = "template"
  bd_name       = "bd"
}`
}

func testAccMSOSchemaTemplateBdDHCPPolicyConfigBothOptionIdentifiers() string {
	return `resource "mso_schema_template_bd_dhcp_policy" "invalid" {
  schema_id               = "schema"
  template_name           = "template"
  bd_name                 = "bd"
  name                    = "relay"
  dhcp_option_policy_name = "option"
  dhcp_option_policy_uuid = "00000000-0000-4000-8000-000000000002"
}`
}

func testAccMSOSchemaTemplateBdDHCPPolicyCreateConfig() string {
	return testAccMSOSchemaTemplateBdDHCPPolicyPrerequisiteConfig() + testAccMSOSchemaTemplateBdDHCPPolicyBDConfig() + fmt.Sprintf(`
resource "mso_schema_template_bd_dhcp_policy" "relay_policy" {
  schema_id               = mso_schema_template_bd.dhcp_parent.schema_id
  template_name           = mso_schema_template_bd.dhcp_parent.template_name
  bd_name                 = mso_schema_template_bd.dhcp_parent.name
  name                    = mso_tenant_policies_dhcp_relay_policy.%s.name
  dhcp_option_policy_name = mso_tenant_policies_dhcp_option_policy.%s.name
}
`, msoTenantPoliciesDhcpRelayPolicyName, msoTenantPoliciesDhcpOptionPolicyName)
}

func testAccMSOSchemaTemplateBdDHCPPolicyUUIDCreateConfig() string {
	return testAccMSOSchemaTemplateBdDHCPPolicyPrerequisiteConfigWithTwoOptionPolicies() + testAccMSOSchemaTemplateBdDHCPPolicyBDConfig() + fmt.Sprintf(`
resource "mso_schema_template_bd_dhcp_policy" "relay_policy" {
  depends_on                = [mso_schema_template_bd.dhcp_parent]
  schema_id                 = mso_schema_template_bd.dhcp_parent.schema_id
  template_name             = mso_schema_template_bd.dhcp_parent.template_name
  bd_name                   = mso_schema_template_bd.dhcp_parent.name
  uuid                      = mso_tenant_policies_dhcp_relay_policy.%s.uuid
  dhcp_option_policy_uuid   = mso_tenant_policies_dhcp_option_policy.%s.uuid
}
resource "mso_schema_template_bd_dhcp_policy" "relay_policy_without_option" {
  depends_on    = [mso_schema_template_bd_dhcp_policy.relay_policy]
  schema_id     = mso_schema_template_bd.dhcp_parent.schema_id
  template_name = mso_schema_template_bd.dhcp_parent.template_name
  bd_name       = mso_schema_template_bd.dhcp_parent.name
  uuid          = mso_tenant_policies_dhcp_relay_policy.%s.uuid
}
`, msoTenantPoliciesDhcpRelayPolicyName, msoTenantPoliciesDhcpOptionPolicyName, msoTenantPoliciesDhcpRelayPolicyName2)
}

func testAccMSOSchemaTemplateBdDHCPPolicyUUIDUpdateOptionConfig() string {
	return testAccMSOSchemaTemplateBdDHCPPolicyPrerequisiteConfigWithTwoOptionPolicies() + testAccMSOSchemaTemplateBdDHCPPolicyBDConfig() + fmt.Sprintf(`
resource "mso_schema_template_bd_dhcp_policy" "relay_policy" {
  depends_on                = [mso_schema_template_bd.dhcp_parent]
  schema_id                 = mso_schema_template_bd.dhcp_parent.schema_id
  template_name             = mso_schema_template_bd.dhcp_parent.template_name
  bd_name                   = mso_schema_template_bd.dhcp_parent.name
  uuid                      = mso_tenant_policies_dhcp_relay_policy.%s.uuid
  dhcp_option_policy_uuid   = mso_tenant_policies_dhcp_option_policy.%s.uuid
}
resource "mso_schema_template_bd_dhcp_policy" "relay_policy_without_option" {
  depends_on    = [mso_schema_template_bd_dhcp_policy.relay_policy]
  schema_id     = mso_schema_template_bd.dhcp_parent.schema_id
  template_name = mso_schema_template_bd.dhcp_parent.template_name
  bd_name       = mso_schema_template_bd.dhcp_parent.name
  uuid          = mso_tenant_policies_dhcp_relay_policy.%s.uuid
}
`, msoTenantPoliciesDhcpRelayPolicyName, msoTenantPoliciesDhcpOptionPolicyName2, msoTenantPoliciesDhcpRelayPolicyName2)
}

func testAccMSOSchemaTemplateBdDHCPPolicyPrerequisiteConfigWithTwoOptionPolicies() string {
	return fmt.Sprintf("%s%s%s", testSiteConfigAnsibleTest(), testTenantConfig(), testSchemaConfig()) +
		testAccMSOSchemaTemplateBdDhcpPrerequisiteConfig() +
		testAccMSOSchemaTemplateBdDhcpRelayPolicyConfig() +
		testAccMSOSchemaTemplateBdDhcpOptionPolicyConfig() +
		testAccMSOSchemaTemplateBdDhcpOptionPolicyConfig2()
}

func testAccMSOSchemaTemplateBdDHCPPolicySwitchOptionConfig() string {
	return testAccMSOSchemaTemplateBdDHCPPolicyPrerequisiteConfigWithTwoOptionPolicies() + testAccMSOSchemaTemplateBdDHCPPolicyBDConfig() + fmt.Sprintf(`
resource "mso_schema_template_bd_dhcp_policy" "relay_policy" {
  schema_id               = mso_schema_template_bd.dhcp_parent.schema_id
  template_name           = mso_schema_template_bd.dhcp_parent.template_name
  bd_name                 = mso_schema_template_bd.dhcp_parent.name
  name                    = mso_tenant_policies_dhcp_relay_policy.%s.name
  dhcp_option_policy_name = mso_tenant_policies_dhcp_option_policy.%s.name
}
`, msoTenantPoliciesDhcpRelayPolicyName, msoTenantPoliciesDhcpOptionPolicyName2)
}

func testAccMSOSchemaTemplateBdDHCPPolicyAddWithoutOptionConfig() string {
	return testAccMSOSchemaTemplateBdDHCPPolicySwitchOptionConfig() + fmt.Sprintf(`
resource "mso_schema_template_bd_dhcp_policy" "relay_policy_without_option" {
  depends_on    = [mso_schema_template_bd_dhcp_policy.relay_policy]
  schema_id     = mso_schema_template_bd.dhcp_parent.schema_id
  template_name = mso_schema_template_bd.dhcp_parent.template_name
  bd_name       = mso_schema_template_bd.dhcp_parent.name
  name          = mso_tenant_policies_dhcp_relay_policy.%s.name
}
`, msoTenantPoliciesDhcpRelayPolicyName2)
}

func testAccMSOSchemaTemplateBdDHCPPolicyUpdateConfig() string {
	return testAccMSOSchemaTemplateBdDHCPPolicyPrerequisiteConfigWithTwoOptionPolicies() + testAccMSOSchemaTemplateBdDHCPPolicyBDConfig() + fmt.Sprintf(`
resource "mso_schema_template_bd_dhcp_policy" "relay_policy" {
  schema_id     = mso_schema_template_bd.dhcp_parent.schema_id
  template_name = mso_schema_template_bd.dhcp_parent.template_name
  bd_name       = mso_schema_template_bd.dhcp_parent.name
  name          = mso_tenant_policies_dhcp_relay_policy.%s.name
}
resource "mso_schema_template_bd_dhcp_policy" "relay_policy_without_option" {
  depends_on    = [mso_schema_template_bd_dhcp_policy.relay_policy]
  schema_id     = mso_schema_template_bd.dhcp_parent.schema_id
  template_name = mso_schema_template_bd.dhcp_parent.template_name
  bd_name       = mso_schema_template_bd.dhcp_parent.name
  name          = mso_tenant_policies_dhcp_relay_policy.%s.name
}
`, msoTenantPoliciesDhcpRelayPolicyName, msoTenantPoliciesDhcpRelayPolicyName2)
}

func testAccMSOSchemaTemplateBdDHCPPolicyRemoveFirstConfig() string {
	return testAccMSOSchemaTemplateBdDHCPPolicyPrerequisiteConfig() + testAccMSOSchemaTemplateBdDHCPPolicyBDConfig() + fmt.Sprintf(`
resource "mso_schema_template_bd_dhcp_policy" "relay_policy_without_option" {
  schema_id     = mso_schema_template_bd.dhcp_parent.schema_id
  template_name = mso_schema_template_bd.dhcp_parent.template_name
  bd_name       = mso_schema_template_bd.dhcp_parent.name
  name          = mso_tenant_policies_dhcp_relay_policy.%s.name
}
`, msoTenantPoliciesDhcpRelayPolicyName2)
}

func testAccMSOSchemaTemplateBdDHCPPolicyBDConfig() string {
	return fmt.Sprintf(`
resource "mso_schema_template_bd" "dhcp_parent" {
  schema_id              = mso_schema.%s.id
  template_name          = "%s"
  name                   = "%s_dhcp"
  display_name           = "%s_dhcp"
  vrf_name               = mso_schema_template_vrf.%s.name
  layer2_stretch         = true
  layer2_unknown_unicast = "proxy"
  intersite_bum_traffic  = true
}
`, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateBdName, msoSchemaTemplateBdName, msoSchemaTemplateVrfName)
}

func testAccMSOSchemaTemplateBdDHCPPolicyMixedConfig() string {
	return testAccMSOSchemaTemplateBdDHCPPolicyPrerequisiteConfig() + fmt.Sprintf(`
resource "mso_schema_template_bd" "dhcp_parent" {
  schema_id              = mso_schema.%s.id
  template_name          = "%s"
  name                   = "%s_dhcp"
  display_name           = "%s_dhcp"
  vrf_name               = mso_schema_template_vrf.%s.name
  layer2_stretch         = true
  layer2_unknown_unicast = "proxy"
  intersite_bum_traffic  = true
  dhcp_policies {
    name = mso_tenant_policies_dhcp_relay_policy.%s.name
  }
}
resource "mso_schema_template_bd_dhcp_policy" "relay_policy_without_option" {
  schema_id     = mso_schema_template_bd.dhcp_parent.schema_id
  template_name = mso_schema_template_bd.dhcp_parent.template_name
  bd_name       = mso_schema_template_bd.dhcp_parent.name
  name          = mso_tenant_policies_dhcp_relay_policy.%s.name
}
`, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateBdName, msoSchemaTemplateBdName, msoSchemaTemplateVrfName, msoTenantPoliciesDhcpRelayPolicyName, msoTenantPoliciesDhcpRelayPolicyName2)
}
