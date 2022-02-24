package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSODHCPRelayPolicyProvider_DataSource(t *testing.T) {
	var provider models.DHCPRelayPolicyProvider
	resourceName := "mso_dhcp_relay_policy_provider.test"
	dataSourceName := "data.mso_dhcp_relay_policy_provider.test"
	name := makeTestVariable(acctest.RandString(5))
	nameOther := makeTestVariable(acctest.RandString(5))
	addr, _ := acctest.RandIpAddress("10.6.0.0/16")
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSODHCPRelayPolicyProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config:      MSODHCPRelayPolicyProviderDataSourceWithoutRequired(tenantNames[0], name, addr, epg, "dhcp_relay_policy_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSODHCPRelayPolicyProviderDataSourceWithoutRequired(tenantNames[0], name, addr, epg, "dhcp_server_address"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSODHCPRelayPolicyProviderDataSourceWithRequired(tenantNames[0], name, addr, epg),
				ExpectError: regexp.MustCompile(`one of (.)+ must be specified`),
			},
			{
				Config:      MSODHCPRelayPolicyProviderDataSourceWithEpgExtEpg(tenantNames[0], name, addr, epg, randomValue),
				ExpectError: regexp.MustCompile(`(.)+ conflicts with external_epg_ref`),
			},
			{
				Config:      MSODHCPRelayPolicyProviderDataSourceAttr(tenantNames[0], name, addr, epg, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)+is not expected here.`),
			},
			{
				Config:      MSODHCPRelayPolicyProviderDataSourceWithInvalidParentResourceName(tenantNames[0], name, addr, epg),
				ExpectError: regexp.MustCompile(`DHCP Relay Policy with name(.)+ not found`),
			},
			{
				Config: MSODHCPRelayPolicyProviderDataSourceWithExtEpg(tenantNames[0], nameOther, addr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPRelayPolicyProviderExists(resourceName, &provider),
					resource.TestCheckResourceAttrPair(resourceName, "dhcp_relay_policy_name", dataSourceName, "dhcp_relay_policy_name"),
					resource.TestCheckResourceAttrPair(resourceName, "dhcp_server_address", dataSourceName, "dhcp_server_address"),
					resource.TestCheckResourceAttrPair(resourceName, "epg_ref", dataSourceName, "epg_ref"),
					resource.TestCheckResourceAttrPair(resourceName, "external_epg_ref", dataSourceName, "external_epg_ref"),
				),
			},
		},
	})
}

func MSODHCPRelayPolicyProviderDataSourceWithExtEpg(tenant, name, addr string) string {
	resource := MSODHCPRelayPolicyProviderWithExtEpg(tenant, name, name, addr)
	resource += fmt.Sprintln(`
	data "mso_dhcp_relay_policy_provider" "test" {
		dhcp_relay_policy_name = mso_dhcp_relay_policy_provider.test.dhcp_relay_policy_name
		dhcp_server_address = mso_dhcp_relay_policy_provider.test.dhcp_server_address
		external_epg_ref = mso_dhcp_relay_policy_provider.test.external_epg_ref
	}
	`)
	return resource
}

func MSODHCPRelayPolicyProviderDataSourceWithEpg(tenant, name, addr, epg string) string {
	resource := MSODHCPRelayPolicyProviderWithEpg(tenant, name, addr, epg)
	resource += fmt.Sprintln(`
	data "mso_dhcp_relay_policy_provider" "test" {
		dhcp_relay_policy_name = mso_dhcp_relay_policy_provider.test.dhcp_relay_policy_name
		dhcp_server_address = mso_dhcp_relay_policy_provider.test.dhcp_server_address
		epg_ref = mso_dhcp_relay_policy_provider.test.epg_ref
	}
	`)
	return resource
}

func MSODHCPRelayPolicyProviderDataSourceWithInvalidParentResourceName(tenant, name, addr, epg string) string {
	resource := MSODHCPRelayPolicyProviderWithEpg(tenant, name, addr, epg)
	resource += fmt.Sprintln(`
	data "mso_dhcp_relay_policy_provider" "test" {
		dhcp_relay_policy_name = "${mso_dhcp_relay_policy_provider.test.dhcp_relay_policy_name}_invalid"
		dhcp_server_address = mso_dhcp_relay_policy_provider.test.dhcp_server_address
		epg_ref = mso_dhcp_relay_policy_provider.test.epg_ref
	}
	`)
	return resource
}

func MSODHCPRelayPolicyProviderDataSourceAttr(tenant, name, addr, epg, key, value string) string {
	resource := MSODHCPRelayPolicyProviderWithEpg(tenant, name, addr, epg)
	resource += fmt.Sprintf(`
	data "mso_dhcp_relay_policy_provider" "test" {
		dhcp_relay_policy_name = mso_dhcp_relay_policy_provider.test.dhcp_relay_policy_name
		dhcp_server_address = mso_dhcp_relay_policy_provider.test.dhcp_server_address
		epg_ref = mso_dhcp_relay_policy_provider.test.epg_ref
		%s = "%s"
	}
	`, key, value)
	return resource
}

func MSODHCPRelayPolicyProviderDataSourceWithEpgExtEpg(tenant, name, addr, epg, val string) string {
	resource := MSODHCPRelayPolicyProviderWithEpg(tenant, name, addr, epg)
	resource += fmt.Sprintf(`
	data "mso_dhcp_relay_policy_provider" "test" {
		dhcp_relay_policy_name = mso_dhcp_relay_policy_provider.test.dhcp_relay_policy_name
		dhcp_server_address = mso_dhcp_relay_policy_provider.test.dhcp_server_address
		epg_ref = mso_dhcp_relay_policy_provider.test.epg_ref
		external_epg_ref = "%s"
	}
	`, val)
	return resource
}

func MSODHCPRelayPolicyProviderDataSourceWithRequired(tenant, name, addr, epg string) string {
	resource := MSODHCPRelayPolicyProviderWithEpg(tenant, name, addr, epg)
	resource += fmt.Sprintln(`
	data "mso_dhcp_relay_policy_provider" "test" {
		dhcp_relay_policy_name = mso_dhcp_relay_policy_provider.test.dhcp_relay_policy_name
		dhcp_server_address = mso_dhcp_relay_policy_provider.test.dhcp_server_address
	}
	`)
	return resource
}

func MSODHCPRelayPolicyProviderDataSourceWithoutRequired(tenant, name, addr, epg, attr string) string {
	rBlock := MSODHCPRelayPolicyProviderWithEpg(tenant, name, addr, epg)
	switch attr {
	case "dhcp_relay_policy_name":
		rBlock += `
		data "mso_dhcp_relay_policy_provider" "test" {
		#	dhcp_relay_policy_name = mso_dhcp_relay_policy_provider.test.dhcp_relay_policy_name
			dhcp_server_address = mso_dhcp_relay_policy_provider.test.dhcp_server_address
		}
		`
	case "dhcp_server_address":
		rBlock += `
		data "mso_dhcp_relay_policy_provider" "test" {
			dhcp_relay_policy_name = mso_dhcp_relay_policy_provider.test.dhcp_relay_policy_name
		#	dhcp_server_address = mso_dhcp_relay_policy_provider.test.dhcp_server_address
		}
		`
	}
	return fmt.Sprintln(rBlock)
}
