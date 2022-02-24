package mso

import (
	"fmt"
	"regexp"

	"testing"

	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSODHCPRelayPolicy_DataSource(t *testing.T) {
	var relayPolicy models.DHCPRelayPolicy
	resourceName := "mso_dhcp_relay_policy.test"
	dataSourceName := "data.mso_dhcp_relay_policy.test"
	tenant := tenantNames[0]
	epg := epg
	name := makeTestVariable(acctest.RandString(5))
	schemaName := makeTestVariable(acctest.RandString(5))
	templateName := makeTestVariable(acctest.RandString(5))
	displayName := makeTestVariable(acctest.RandString(5))
	dhcpServerAddress, _ := acctest.RandIpAddress("1.2.0.0/16")
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSODHCPRelayPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      MSODHCPRelayPolicyDataSourceWithoutName(tenant, name),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSODHCPRelayPolicyDataSourceAttr(tenant, name, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)+is not expected here.`),
			},
			{
				Config:      MSODHCPRelayPolicyDataSourceWithInvalidName(tenant, name),
				ExpectError: regexp.MustCompile(`DHCP Relay Policy with name: (.)+ not found`),
			},
			{
				Config: MSODHCPRelayPolicyWithEPGDataSourceConfig(tenant, name, dhcpServerAddress, epg),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPRelayPolicyExists(resourceName, &relayPolicy),
					resource.TestCheckResourceAttrPair(resourceName, "description", dataSourceName, "description"),
					resource.TestCheckResourceAttrPair(resourceName, "tenant_id", dataSourceName, "tenant_id"),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "dhcp_relay_policy_provider.#", dataSourceName, "dhcp_relay_policy_provider.#"),
					resource.TestCheckResourceAttrPair(resourceName, "dhcp_relay_policy_provider.0.epg", dataSourceName, "dhcp_relay_policy_provider.0.epg"),
					resource.TestCheckResourceAttrPair(resourceName, "dhcp_relay_policy_provider.0.dhcp_server_address", dataSourceName, "dhcp_relay_policy_provider.0.dhcp_server_address"),
				),
			},
			{
				Config: MSODHCPRelayPolicyWithExternalEPGDataSourceConfig(tenant, name, dhcpServerAddress, templateName, schemaName, displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPRelayPolicyExists(resourceName, &relayPolicy),
					resource.TestCheckResourceAttrPair(resourceName, "description", dataSourceName, "description"),
					resource.TestCheckResourceAttrPair(resourceName, "tenant_id", dataSourceName, "tenant_id"),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "dhcp_relay_policy_provider.#", dataSourceName, "dhcp_relay_policy_provider.#"),
					resource.TestCheckResourceAttrPair(resourceName, "dhcp_relay_policy_provider.0.external_epg", dataSourceName, "dhcp_relay_policy_provider.0.external_epg"),
					resource.TestCheckResourceAttrPair(resourceName, "dhcp_relay_policy_provider.0.dhcp_server_address", dataSourceName, "dhcp_relay_policy_provider.0.dhcp_server_address"),
				),
			},
			{
				Config:             MSODHCPRelayPolicyDataSourceWithUpdatedResource(tenant, name, "description", randomValue),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPRelayPolicyExists(resourceName, &relayPolicy),
					resource.TestCheckResourceAttrPair(resourceName, "description", dataSourceName, "description"),
				),
			},
		},
	})
}

func MSODHCPRelayPolicyWithEPGDataSourceConfig(tenant, name, dhcpServerAddress, epg string) string {
	resource := fmt.Sprintf(`
	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_relay_policy" "test" {
		tenant_id = data.mso_tenant.test.id
		name = "%s"
		dhcp_relay_policy_provider {
			epg = "%s"
			dhcp_server_address = "%s"
		}
	}
	data "mso_dhcp_relay_policy" "test" {
		name = mso_dhcp_relay_policy.test.name
	}
	`, tenant, tenant, name, epg, dhcpServerAddress)
	return resource
}

func MSODHCPRelayPolicyWithExternalEPGDataSourceConfig(tenant, name, dhcpServerAddress, templateName, schemaName, displayName string) string {
	resource := fmt.Sprintf(`
	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource mso_schema "test"{
		name = "%s"
		template_name = "%s"
		tenant_id = data.mso_tenant.test.id
	}
	resource mso_schema_template_vrf "test" {
		schema_id = mso_schema.test.id
		template= mso_schema.test.template_name
		name= "%s"
		display_name= "%s"
	}
	resource "mso_schema_template_external_epg" "test" {
		schema_id = mso_schema.test.id
		template_name = mso_schema.test.template_name
		external_epg_name = "%s"
		display_name = "%s"
		vrf_name = mso_schema_template_vrf.test.name
	}
	resource "mso_dhcp_relay_policy" "test" {
		tenant_id = data.mso_tenant.test.id
		name = "%s"
		dhcp_relay_policy_provider {
			external_epg = mso_schema_template_external_epg.test.id
			dhcp_server_address = "%s"
		}
	}
	data "mso_dhcp_relay_policy" "test" {
		name = mso_dhcp_relay_policy.test.name
	}
	`, tenant, tenant, schemaName, templateName, displayName, displayName, displayName, displayName, name, dhcpServerAddress)
	return resource
}

func MSODHCPRelayPolicyDataSourceWithoutName(tenant, name string) string {
	resource := fmt.Sprintf(`
	resource "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_relay_policy" "test" {
		tenant_id = mso_tenant.test.id
		name = "%s"
	}
	data "mso_dhcp_relay_policy" "test" {}
	`, tenant, tenant, name)
	return resource
}

func MSODHCPRelayPolicyDataSourceWithUpdatedResource(tenant, name, key, val string) string {
	resource := fmt.Sprintf(`
	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_relay_policy" "test" {
		tenant_id = data.mso_tenant.test.id
		name = "%s"
		%s = "%s"
	}
	data "mso_dhcp_relay_policy" "test" {
		name = mso_dhcp_relay_policy.test.name
		depends_on = [mso_dhcp_relay_policy.test]
	}
	`, tenant, tenant, name, key, val)
	return resource
}

func MSODHCPRelayPolicyDataSourceAttr(tenant, name, key, val string) string {
	resource := fmt.Sprintf(`
	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_relay_policy" "test" {
		tenant_id = data.mso_tenant.test.id
		name = "%s"
	}
	data "mso_dhcp_relay_policy" "test" {
		name = mso_dhcp_relay_policy.test.name
		%s = "%s"
	}
	`, tenant, tenant, name, key, val)
	return resource
}

func MSODHCPRelayPolicyDataSourceWithInvalidName(tenant, name string) string {
	resource := fmt.Sprintf(`
	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_relay_policy" "test" {
		tenant_id = data.mso_tenant.test.id
		name = "%s"
	}
	data "mso_dhcp_relay_policy" "test" {
		name = "${mso_dhcp_relay_policy.test.name}_invalid"
	}
	`, tenant, tenant, name)
	return resource
}
