package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOTenantPoliciesDHCPRelayPolicyResource(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	resource.Test(t, resource.TestCase{
		PreCheck:            func() { testAccPreCheck(t) },
		Providers:           testAccProviders,
		DisableBinaryDriver: true,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Create DHCP Relay Policy without providers") },
				Config:      testAccMSOTenantPoliciesDHCPRelayPolicyConfigCreateError(name),
				ExpectError: regexp.MustCompile(`config is invalid: "providers": required field is not set`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create DHCP Relay Policy") },
				Config:    testAccMSOTenantPoliciesDHCPRelayPolicyConfigCreate(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					customTestCheckResourceAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name),
						map[string]string{
							"name":                                   name,
							"description":                            "",
							"template_id":                            "reference",
							"providers.#":                            "2",
							"providers.0.application_epg":            "reference",
							"providers.0.dhcp_server_address":        "1.1.1.1",
							"providers.0.dhcp_server_vrf_preference": "false",
							"providers.0.external_epg":               "",
							"providers.1.application_epg":            "",
							"providers.1.dhcp_server_address":        "2.2.2.2",
							"providers.1.dhcp_server_vrf_preference": "true",
							"providers.1.external_epg":               "reference",
						},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update DHCP Relay Policy - Providers property") },
				Config:    testAccMSOTenantPoliciesDHCPRelayPolicyConfigUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					customTestCheckResourceAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), map[string]string{
						"name":                                   name,
						"description":                            "Updated DHCP Relay Policy",
						"template_id":                            "reference",
						"providers.#":                            "2",
						"providers.0.application_epg":            "reference",
						"providers.0.dhcp_server_address":        "1.1.1.2",
						"providers.0.dhcp_server_vrf_preference": "true",
						"providers.0.external_epg":               "",
						"providers.1.application_epg":            "",
						"providers.1.dhcp_server_address":        "2.2.2.2",
						"providers.1.dhcp_server_vrf_preference": "false",
						"providers.1.external_epg":               "reference",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update DHCP Relay Policy - Remove One Provider") },
				Config:    testAccMSOTenantPoliciesDHCPRelayPolicyConfigUpdateRemoveProvider(name),
				Check: resource.ComposeTestCheckFunc(

					customTestCheckResourceAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), map[string]string{
						"name":                                   name,
						"description":                            "Updated DHCP Relay Policy",
						"template_id":                            "reference",
						"providers.#":                            "1",
						"providers.0.application_epg":            "reference",
						"providers.0.dhcp_server_address":        "1.1.1.2",
						"providers.0.dhcp_server_vrf_preference": "true",
						"providers.0.external_epg":               "",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update DHCP Relay Policy - Add the provider back to the list") },
				Config:    testAccMSOTenantPoliciesDHCPRelayPolicyConfigUpdateAddProvider(name),
				Check: resource.ComposeTestCheckFunc(

					customTestCheckResourceAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), map[string]string{
						"name":                                   name,
						"description":                            "Updated DHCP Relay Policy",
						"template_id":                            "reference",
						"providers.#":                            "2",
						"providers.0.application_epg":            "",
						"providers.0.dhcp_server_address":        "2.2.2.2",
						"providers.0.dhcp_server_vrf_preference": "false",
						"providers.0.external_epg":               "reference",
						"providers.1.application_epg":            "reference",
						"providers.1.dhcp_server_address":        "1.1.1.2",
						"providers.1.dhcp_server_vrf_preference": "true",
						"providers.1.external_epg":               "",
					}),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import DHCP Relay Policy") },
				ResourceName:      fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithArguments("mso_tenant_policies_dhcp_relay_policy", "dhcpRelay"),
	})
}

var dhcpRelayPolicyParentConfig = testSiteConfigAnsibleTest() + testTenantConfig() + testSchemaConfig() + testSchemaTemplateAnpConfig() + testSchemaTemplateAnpEpgConfig() + testSchemaTemplateVrfConfig() + testSchemaTemplateExtEpgConfig() + testTenantPolicyTemplateConfig()

func testAccMSOTenantPoliciesDHCPRelayPolicyConfigCreateError(name string) string {
	return fmt.Sprintf(`%[1]s
resource "mso_tenant_policies_dhcp_relay_policy" "%[2]s" {
	name        = "%[2]s"
	template_id = mso_template.%[3]s.id
}
`, dhcpRelayPolicyParentConfig, name, msoTenantPolicyTemplateName)
}

func testAccMSOTenantPoliciesDHCPRelayPolicyConfigCreate(name string) string {
	return fmt.Sprintf(`%[1]s
resource "mso_tenant_policies_dhcp_relay_policy" "%[2]s" {
	name        = "%[2]s"
	template_id = mso_template.%[3]s.id
	providers {
		dhcp_server_address = "1.1.1.1"
		application_epg     = mso_schema_template_anp_epg.%[4]s.uuid
	}
	providers {
		dhcp_server_address        = "2.2.2.2"
		external_epg               = mso_schema_template_external_epg.%[5]s.uuid
		dhcp_server_vrf_preference = true
	}
}
`, dhcpRelayPolicyParentConfig, name, msoTenantPolicyTemplateName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateExtEpgName)
}

func testAccMSOTenantPoliciesDHCPRelayPolicyConfigUpdate(name string) string {
	return fmt.Sprintf(`%[1]s
resource "mso_tenant_policies_dhcp_relay_policy" "%[2]s" {
	name        = "%[2]s"
	template_id = mso_template.%[3]s.id
	description = "Updated DHCP Relay Policy"
	providers {
		dhcp_server_address        = "1.1.1.2"
		application_epg            = mso_schema_template_anp_epg.%[4]s.uuid
		dhcp_server_vrf_preference = true
	}
	providers {
		dhcp_server_address        = "2.2.2.2"
		external_epg               = mso_schema_template_external_epg.%[5]s.uuid
		dhcp_server_vrf_preference = false
	}
}
`, dhcpRelayPolicyParentConfig, name, msoTenantPolicyTemplateName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateExtEpgName)
}

func testAccMSOTenantPoliciesDHCPRelayPolicyConfigUpdateRemoveProvider(name string) string {
	return fmt.Sprintf(`%[1]s
resource "mso_tenant_policies_dhcp_relay_policy" "%[2]s" {
	name        = "%[2]s"
	template_id = mso_template.%[3]s.id
	description = "Updated DHCP Relay Policy"
	providers {
		dhcp_server_address        = "1.1.1.2"
		application_epg            = mso_schema_template_anp_epg.%[4]s.uuid
		dhcp_server_vrf_preference = true
	}
}
`, dhcpRelayPolicyParentConfig, name, msoTenantPolicyTemplateName, msoSchemaTemplateAnpEpgName)
}

func testAccMSOTenantPoliciesDHCPRelayPolicyConfigUpdateAddProvider(name string) string {
	return fmt.Sprintf(`%[1]s
resource "mso_tenant_policies_dhcp_relay_policy" "%[2]s" {
	name        = "%[2]s"
	template_id = mso_template.%[3]s.id
	description = "Updated DHCP Relay Policy"
	providers {
		dhcp_server_address        = "2.2.2.2"
		external_epg               = mso_schema_template_external_epg.%[4]s.uuid
		dhcp_server_vrf_preference = false
	}
	providers {
		dhcp_server_address        = "1.1.1.2"
		application_epg            = mso_schema_template_anp_epg.%[5]s.uuid
		dhcp_server_vrf_preference = true
	}
}
`, dhcpRelayPolicyParentConfig, name, msoTenantPolicyTemplateName, msoSchemaTemplateExtEpgName, msoSchemaTemplateAnpEpgName)
}
