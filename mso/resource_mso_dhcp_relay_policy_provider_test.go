package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccMSODHCPRelayPolicyProvider_Basic(t *testing.T) {
	// var prov1 models.IntersiteL3outs
	// var prov2 models.IntersiteL3outs
	resourceName := "mso_dhcp_relay_policy_provider.test"
	addr, _ := acctest.RandIpAddress("10.1.0.0/16")
	// addrother,_:=acctest.RandIpAddress("10.2.0.0/16")
	polName := "need_to_update"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSODHCPRelayPolicyProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config:      MSODHCPRelayPolicyProviderWithoutRequired(polName, addr, "dhcp_relay_policy_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSODHCPRelayPolicyProviderWithoutRequired(polName, addr, "dhcp_server_address"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSODHCPRelayPolicyProviderWithRequired(polName, addr),
				ExpectError: regexp.MustCompile(`one of (.)+ must be specified`),
			},
			//TODO: case when both epg and external epg are defined
			//TODO: case with epg
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			//TODO: case with external epg
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			//TODO: policy name update forcenew check if applicable
			//TODO: addr update forcenew check
		},
	})
}

func TestAccMSODHCPRelayPolicyProvider_Negative(t *testing.T) {
	// var prov1 models.IntersiteL3outs
	// var prov2 models.IntersiteL3outs
	addr, _ := acctest.RandIpAddress("10.3.0.0/16")
	polName := "need_to_update"
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSODHCPRelayPolicyProviderDestroy,
		Steps: []resource.TestStep{
			//TODO: correct resource config
			//TODO: incorrect dhcp relay policy name
			//TODO: incorrect address
			//TODO: incorrect epg_ref
			//TODO: incorrect external_epg_ref
			{
				Config:      MSODHCPRelayPolicyProviderAttr(polName, addr, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)+is not expected here.`),
			},
			//TODO: correct resource config
		},
	})
}

func TestAccMSODHCPRelayPolicyProvider_MultipleCreateDelete(t *testing.T) {
	// polName := "need_to_update"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSODHCPRelayPolicyProviderDestroy,
		Steps:        []resource.TestStep{
			//TODO: config for multiple create delete
		},
	})
}

func testAccCheckMSODHCPRelayPolicyProviderDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_dhcp_relay_policy_provider" {
			id := rs.Primary.ID
			prov, err := DHCPRelayPolicyProviderIdtoModel(id)
			if err != nil {
				return err
			}
			_, err = client.ReadDHCPRelayPolicyProvider(prov)
			if err == nil {
				return fmt.Errorf("DHCP Relay Policy Provider still exists")
			}
		}
	}
	return nil
}

func testAccCheckMSODHCPRelayPolicyProviderExists(providerName string, m *models.DHCPRelayPolicyProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs, ok := s.RootModule().Resources[providerName]
		if !ok {
			return fmt.Errorf("DHCP Relay Policy Provider %s not found", providerName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No DHCP Relay Policy Provider Id was set")
		}
		provider, err := DHCPRelayPolicyProviderIdtoModel(rs.Primary.ID)
		if err != nil {
			return err
		}
		var read *models.DHCPRelayPolicyProvider
		read, err = client.ReadDHCPRelayPolicyProvider(provider)
		if err != nil {
			return err
		}
		*m = *read
		return nil
	}
}

func testAccCheckMSODHCPRelayPolicyProviderIdNotEqual(m1, m2 *models.DHCPRelayPolicyProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		id1 := DHCPRelayPolicyProviderModeltoId(m1)
		id2 := DHCPRelayPolicyProviderModeltoId(m2)
		if id1 == id2 {
			return fmt.Errorf("DHCP Relay Policy Provider Ids are equal")
		}
		return nil
	}
}

func MSODHCPRelayPolicyProviderWithoutRequired(polname, addr, attr string) string {
	rBlock := ""
	switch attr {
	case "dhcp_relay_policy_name":
		rBlock += `
		resource "mso_dhcp_relay_policy_provider" "test" {
		#	dhcp_relay_policy_name = "%s"
			dhcp_server_address = "%s"
		}
		`
	case "dhcp_server_address":
		rBlock += `
		resource "mso_dhcp_relay_policy_provider" "test" {
			dhcp_relay_policy_name = "%s"
		#	dhcp_server_address = "%s"
		}
		`
	}
	return fmt.Sprintf(rBlock, polname, addr)
}

func MSODHCPRelayPolicyProviderWithRequired(polname, addr string) string {
	resource := fmt.Sprintf(`
	resource "mso_dhcp_relay_policy_provider" "test" {
		dhcp_relay_policy_name = "%s"
		dhcp_server_address = "%s"
	}
	`, polname, addr)
	return resource
}

func MSODHCPRelayPolicyProviderAttr(polname, addr, key, value string) string {
	resource := fmt.Sprintf(`
	resource "mso_dhcp_relay_policy_provider" "test" {
		dhcp_relay_policy_name = "%s"
		dhcp_server_address = "%s"
		%s = "%s"
	}
	`, polname, addr, key, value)
	return resource
}
