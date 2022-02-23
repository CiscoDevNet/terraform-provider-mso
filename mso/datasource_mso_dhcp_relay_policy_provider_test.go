package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSODHCPRelayPolicyProvider_DataSource(t *testing.T) {
	// var provider models.DHCPRelayPolicyProvider
	// resourceName := "mso_dhcp_relay_policy_provider.test"
	// dataSourceName := "data.mso_dhcp_relay_policy_provider.test"
	polName := "need_to_update"
	addr,_:=acctest.RandIpAddress("10.6.0.0/16")
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSODHCPRelayPolicyProviderDestroy,
		Steps: []resource.TestStep{
			{
				//TODO: need to update resource block
				Config:      MSODHCPRelayPolicyProviderDataSourceWithoutRequired(polName,addr, "dhcp_relay_policy_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSODHCPRelayPolicyProviderDataSourceWithoutRequired(polName,addr, "dhcp_server_address"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			//TODO: case when both epg and external epg are defined
			{
				Config:MSODHCPRelayPolicyProviderDataSourceAttr(polName,addr,randomParameter,randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)+is not expected here.`),
			},
			{
				//TODO: config with invalid dhcp relay policy name
				ExpectError: regexp.MustCompile(`DHCP Relay Policy with name(.)+ not found`),
			},
			//TODO: case with epg
			//TODO: case with external epg
		},
	})
}

func MSODHCPRelayPolicyProviderDataSourceWithoutRequired(polname, addr, attr string) string {
	rBlock := `
	resource "mso_dhcp_relay_policy_provider" "test" {
		dhcp_relay_policy_name = "%s"
		dhcp_server_address = "%s"
	}
	`
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
	return fmt.Sprintf(rBlock, polname, addr)
}

func MSODHCPRelayPolicyProviderDataSourceAttr(polname,addr,key,value string) string{
	resource:=fmt.Sprintf(`
	resource "mso_dhcp_relay_policy_provider" "test" {
		dhcp_relay_policy_name = "%s"
		dhcp_server_address = "%s"
	}
	data "mso_dhcp_relay_policy_provider" "test" {
		dhcp_relay_policy_name = mso_dhcp_relay_policy_provider.test.dhcp_relay_policy_name
		dhcp_server_address = mso_dhcp_relay_policy_provider.test.dhcp_server_address
		%s = "%s"
	}
	`,polname,addr,key,value)
	return resource
}
