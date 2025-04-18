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
				PreConfig:   func() { fmt.Println("Test: Create DHCP Relay Policy without dhcp_relay_providers") },
				Config:      testAccMSOTenantPoliciesDHCPRelayPolicyConfigCreateErrorWithoutProviders(name),
				ExpectError: regexp.MustCompile(`config is invalid: "dhcp_relay_providers": required field is not set`),
			},
			{
				PreConfig:                 func() { fmt.Println("Test: Create DHCP Relay Policy with invalid dhcp_relay_providers") },
				Config:                    testAccMSOTenantPoliciesDHCPRelayPolicyConfigCreateErrorWithInvalidProviders(name),
				Destroy:                   false,
				PreventPostDestroyRefresh: true,
				ExpectError:               regexp.MustCompile(`[\nError: Set either 'application_epg_uuid' or 'external_epg_uuid', not both for a provider at index position: 0\n\nError: Please set either 'application_epg_uuid' or 'external_epg_uuid' for a provider at index position: 1\n]`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create DHCP Relay Policy") },
				Config:    testAccMSOTenantPoliciesDHCPRelayPolicyConfigCreate(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "name", name),
					resource.TestCheckResourceAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "description", ""),
					resource.TestCheckResourceAttrSet(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "template_id"),
					resource.TestCheckResourceAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "dhcp_relay_providers.#", "2"),
					customTestCheckResourceTypeSetAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "dhcp_relay_providers",
						map[string]string{
							"application_epg_uuid":       fmt.Sprintf("mso_schema_template_anp_epg.%s.uuid", msoSchemaTemplateAnpEpgName),
							"dhcp_server_address":        "1.1.1.1",
							"dhcp_server_vrf_preference": "false",
							"external_epg_uuid":          "",
						},
					),
					customTestCheckResourceTypeSetAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "dhcp_relay_providers",
						map[string]string{
							"application_epg_uuid":       "",
							"dhcp_server_address":        "2.2.2.2",
							"dhcp_server_vrf_preference": "true",
							"external_epg_uuid":          fmt.Sprintf("mso_schema_template_external_epg.%s.uuid", msoSchemaTemplateExtEpgName),
						},
					),
				),
			},
			{
				PreConfig:                 func() { fmt.Println("Test: Update DHCP Relay Policy with invalid dhcp_relay_providers") },
				Config:                    testAccMSOTenantPoliciesDHCPRelayPolicyConfigCreateErrorWithInvalidProviders(name),
				Destroy:                   false,
				PreventPostDestroyRefresh: true,
				ExpectError:               regexp.MustCompile(`[\nError: Set either 'application_epg_uuid' or 'external_epg_uuid', not both for a provider at index position: 0\n\nError: Please set either 'application_epg_uuid' or 'external_epg_uuid' for a provider at index position: 1\n]`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update DHCP Relay Policy - Providers property") },
				Config:    testAccMSOTenantPoliciesDHCPRelayPolicyConfigUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "name", name),
					resource.TestCheckResourceAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "description", "Updated DHCP Relay Policy"),
					resource.TestCheckResourceAttrSet(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "template_id"),
					resource.TestCheckResourceAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "dhcp_relay_providers.#", "2"),
					customTestCheckResourceTypeSetAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "dhcp_relay_providers",
						map[string]string{
							"application_epg_uuid":       fmt.Sprintf("mso_schema_template_anp_epg.%s.uuid", msoSchemaTemplateAnpEpgName),
							"dhcp_server_address":        "1.1.1.1",
							"dhcp_server_vrf_preference": "true",
							"external_epg_uuid":          "",
						},
					),
					customTestCheckResourceTypeSetAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "dhcp_relay_providers",
						map[string]string{
							"application_epg_uuid":       "",
							"dhcp_server_address":        "2.2.2.2",
							"dhcp_server_vrf_preference": "false",
							"external_epg_uuid":          fmt.Sprintf("mso_schema_template_external_epg.%s.uuid", msoSchemaTemplateExtEpgName),
						},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update DHCP Relay Policy - Remove One Provider") },
				Config:    testAccMSOTenantPoliciesDHCPRelayPolicyConfigUpdateRemoveProvider(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "name", name),
					resource.TestCheckResourceAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "description", "Updated DHCP Relay Policy"),
					resource.TestCheckResourceAttrSet(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "template_id"),
					resource.TestCheckResourceAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "dhcp_relay_providers.#", "1"),
					customTestCheckResourceTypeSetAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "dhcp_relay_providers",
						map[string]string{
							"application_epg_uuid":       fmt.Sprintf("mso_schema_template_anp_epg.%s.uuid", msoSchemaTemplateAnpEpgName),
							"dhcp_server_address":        "1.1.1.1",
							"dhcp_server_vrf_preference": "false",
							"external_epg_uuid":          "",
						},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update DHCP Relay Policy - Add the provider back to the list") },
				Config:    testAccMSOTenantPoliciesDHCPRelayPolicyConfigUpdateAddProvider(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "name", name),
					resource.TestCheckResourceAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "description", "Updated DHCP Relay Policy"),
					resource.TestCheckResourceAttrSet(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "template_id"),
					resource.TestCheckResourceAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "dhcp_relay_providers.#", "2"),
					customTestCheckResourceTypeSetAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "dhcp_relay_providers",
						map[string]string{
							"application_epg_uuid":       fmt.Sprintf("mso_schema_template_anp_epg.%s.uuid", msoSchemaTemplateAnpEpgName),
							"dhcp_server_address":        "1.1.1.2",
							"dhcp_server_vrf_preference": "true",
							"external_epg_uuid":          "",
						},
					),
					customTestCheckResourceTypeSetAttr(fmt.Sprintf("mso_tenant_policies_dhcp_relay_policy.%s", name), "dhcp_relay_providers",
						map[string]string{
							"application_epg_uuid":       "",
							"dhcp_server_address":        "2.2.2.2",
							"dhcp_server_vrf_preference": "false",
							"external_epg_uuid":          fmt.Sprintf("mso_schema_template_external_epg.%s.uuid", msoSchemaTemplateExtEpgName),
						},
					),
				),
			},
			{
				PreConfig:                 func() { fmt.Println("Test: Update DHCP Relay Policy with invalid DHCP server address") },
				Config:                    testAccMSOTenantPoliciesDHCPRelayPolicyConfigUpdateAddProviderWithInvalidDHCPServerAddress(name),
				Destroy:                   false,
				PreventPostDestroyRefresh: true,
				ExpectError:               regexp.MustCompile(`expected dhcp_relay_providers.0.dhcp_server_address to contain a valid IP, got: google.com`),
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

func testAccMSOTenantPoliciesDHCPRelayPolicyConfigCreateErrorWithoutProviders(name string) string {
	return fmt.Sprintf(`%[1]s
resource "mso_tenant_policies_dhcp_relay_policy" "%[2]s" {
	name        = "%[2]s"
	template_id = mso_template.%[3]s.id
}
`, dhcpRelayPolicyParentConfig, name, msoTenantPolicyTemplateName)
}

func testAccMSOTenantPoliciesDHCPRelayPolicyConfigCreateErrorWithInvalidProviders(name string) string {
	return fmt.Sprintf(`%[1]s
resource "mso_tenant_policies_dhcp_relay_policy" "%[2]s" {
	name        = "%[2]s"
	template_id = mso_template.%[3]s.id
	dhcp_relay_providers {
		dhcp_server_address  = "1.1.1.1"
		application_epg_uuid = mso_schema_template_anp_epg.%[4]s.uuid
		external_epg_uuid    = mso_schema_template_external_epg.%[5]s.uuid
	}
	dhcp_relay_providers {
		dhcp_server_address        = "2.2.2.2"
		dhcp_server_vrf_preference = true
	}
}
`, dhcpRelayPolicyParentConfig, name, msoTenantPolicyTemplateName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateExtEpgName)
}

func testAccMSOTenantPoliciesDHCPRelayPolicyConfigCreate(name string) string {
	return fmt.Sprintf(`%[1]s
resource "mso_tenant_policies_dhcp_relay_policy" "%[2]s" {
	name        = "%[2]s"
	template_id = mso_template.%[3]s.id
	dhcp_relay_providers {
		dhcp_server_address  = "1.1.1.1"
		application_epg_uuid = mso_schema_template_anp_epg.%[4]s.uuid
	}
	dhcp_relay_providers {
		dhcp_server_address        = "2.2.2.2"
		external_epg_uuid          = mso_schema_template_external_epg.%[5]s.uuid
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
	dhcp_relay_providers {
		dhcp_server_address        = "1.1.1.1"
		application_epg_uuid       = mso_schema_template_anp_epg.%[4]s.uuid
		dhcp_server_vrf_preference = true
	}
	dhcp_relay_providers {
		dhcp_server_address        = "2.2.2.2"
		external_epg_uuid          = mso_schema_template_external_epg.%[5]s.uuid
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
	dhcp_relay_providers {
		dhcp_server_address  = "1.1.1.1"
		application_epg_uuid = mso_schema_template_anp_epg.%[4]s.uuid
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
	dhcp_relay_providers {
		dhcp_server_address        = "2.2.2.2"
		external_epg_uuid          = mso_schema_template_external_epg.%[4]s.uuid
		dhcp_server_vrf_preference = false
	}
	dhcp_relay_providers {
		dhcp_server_address        = "1.1.1.2"
		application_epg_uuid       = mso_schema_template_anp_epg.%[5]s.uuid
		dhcp_server_vrf_preference = true
	}
}
`, dhcpRelayPolicyParentConfig, name, msoTenantPolicyTemplateName, msoSchemaTemplateExtEpgName, msoSchemaTemplateAnpEpgName)
}

func testAccMSOTenantPoliciesDHCPRelayPolicyConfigUpdateAddProviderWithInvalidDHCPServerAddress(name string) string {
	return fmt.Sprintf(`%[1]s
resource "mso_tenant_policies_dhcp_relay_policy" "%[2]s" {
	name        = "%[2]s"
	template_id = mso_template.%[3]s.id
	description = "Updated DHCP Relay Policy"
	dhcp_relay_providers {
		dhcp_server_address        = "google.com"
		external_epg_uuid          = mso_schema_template_external_epg.%[4]s.uuid
		dhcp_server_vrf_preference = false
	}
}
`, dhcpRelayPolicyParentConfig, name, msoTenantPolicyTemplateName, msoSchemaTemplateExtEpgName)
}
