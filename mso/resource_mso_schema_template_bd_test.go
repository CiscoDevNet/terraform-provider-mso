package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// msoSchemaTemplateBdSchemaId is set during the first test step's Check to capture the dynamic schema ID for use in the manual deletion PreConfig step.
var msoSchemaTemplateBdSchemaId string

func TestAccMSOSchemaTemplateBdResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create Schema Template BD") },
				Config:    testAccMSOSchemaTemplateBdConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_bd."+msoSchemaTemplateBdName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "display_name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "description", "Terraform test BD"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "layer2_unknown_unicast", "flood"),
					// Capture the dynamic schema ID from state for use in the manual deletion PreConfig step
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["mso_schema_template_bd."+msoSchemaTemplateBdName]
						if !ok {
							return fmt.Errorf("BD resource not found in state")
						}
						msoSchemaTemplateBdSchemaId = rs.Primary.Attributes["schema_id"]
						return nil
					},
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "unknown_multicast_flooding", "flood"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "multi_destination_flooding", "flood_in_bd"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "ipv6_unknown_multicast_flooding", "flood"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "arp_flooding", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "unicast_routing", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "layer3_multicast", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "intersite_bum_traffic", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "optimize_wan_bandwidth", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "layer2_stretch", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "virtual_mac_address", ""),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "ep_move_detection_mode", "none"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "vrf_name", msoSchemaTemplateVrfL3MulticastName),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Schema Template BD layer2_unknown_unicast to proxy") },
				Config:    testAccMSOSchemaTemplateBdConfigUpdateL2UnknownUnicast(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "layer2_unknown_unicast", "proxy"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "description", "Terraform test BD"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Schema Template BD description and enable layer3_multicast") },
				Config:    testAccMSOSchemaTemplateBdConfigUpdateDescriptionAndL3Multicast(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "description", "Terraform test BD updated"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "layer3_multicast", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "layer2_unknown_unicast", "proxy"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Schema Template BD flooding fields") },
				Config:    testAccMSOSchemaTemplateBdConfigUpdateFloodingFields(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "unknown_multicast_flooding", "optimized_flooding"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "multi_destination_flooding", "flood_in_encap"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "ipv6_unknown_multicast_flooding", "optimized_flooding"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Schema Template BD remaining attributes") },
				Config:    testAccMSOSchemaTemplateBdConfigUpdateRemainingAttributes(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "arp_flooding", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "unicast_routing", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "optimize_wan_bandwidth", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "virtual_mac_address", "00:00:5E:00:01:01"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "ep_move_detection_mode", "garp"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Schema Template BD by removing description and virtual_mac_address") },
				Config:    testAccMSOSchemaTemplateBdConfigRemoveDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "description", ""),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "arp_flooding", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "unicast_routing", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "optimize_wan_bandwidth", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "virtual_mac_address", ""),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "ep_move_detection_mode", "none"),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import BD") },
				ResourceName: "mso_schema_template_bd." + msoSchemaTemplateBdName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["mso_schema_template_bd."+msoSchemaTemplateBdName]
					if !ok {
						return "", fmt.Errorf("BD resource not found in state")
					}
					return fmt.Sprintf("%s/templates/%s/bds/%s", rs.Primary.Attributes["schema_id"], rs.Primary.Attributes["template_name"], rs.Primary.Attributes["name"]), nil
				},
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate BD after manual deletion from NDO")
					msoClient := testAccProvider.Meta().(*client.Client)
					bdRemovePatchPayload := models.GetRemovePatchPayload(fmt.Sprintf("/templates/%s/bds/%s", msoSchemaTemplateName, msoSchemaTemplateBdName))
					_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", msoSchemaTemplateBdSchemaId), bdRemovePatchPayload)
					if err != nil {
						t.Fatalf("Failed to manually delete BD: %v", err)
					}
				},
				Config: testAccMSOSchemaTemplateBdConfigRemoveDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_bd."+msoSchemaTemplateBdName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "display_name", msoSchemaTemplateBdName),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add subnet child to BD") },
				Config:    testAccMSOSchemaTemplateBdConfigWithChildren(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_bd."+msoSchemaTemplateBdName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "display_name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "ip", msoSchemaTemplateBdSubnetIp),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "scope", "private"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update BD description with children present") },
				Config:    testAccMSOSchemaTemplateBdConfigWithChildrenUpdateDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_bd."+msoSchemaTemplateBdName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "description", "Terraform test BD with children"),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "ip", msoSchemaTemplateBdSubnetIp),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "scope", "private"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove children from BD") },
				Config:    testAccMSOSchemaTemplateBdConfigRemoveChildren(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_bd."+msoSchemaTemplateBdName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "description", "Terraform test BD with children"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add one DHCP policy to BD") },
				Config:    testAccMSOSchemaTemplateBdConfigWithOneDhcpPolicy(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_bd."+msoSchemaTemplateBdName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "dhcp_policies.#", "1"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update BD to two DHCP policies") },
				Config:    testAccMSOSchemaTemplateBdConfigWithTwoDhcpPolicies(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_bd."+msoSchemaTemplateBdName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "dhcp_policies.#", "2"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update BD back to one DHCP policy") },
				Config:    testAccMSOSchemaTemplateBdConfigWithOneDhcpPolicy(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_bd."+msoSchemaTemplateBdName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "dhcp_policies.#", "1"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove all DHCP policies from BD") },
				Config:    testAccMSOSchemaTemplateBdConfigRemoveDhcpPolicies(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_bd."+msoSchemaTemplateBdName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "dhcp_policies.#", "0"),
				),
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithArguments("mso_schema_template_bd", "bd"),
	})
}

func testAccMSOSchemaTemplateBdPrerequisiteConfig() string {
	return fmt.Sprintf(`%s%s%s%s`, testSiteConfigAnsibleTest(), testTenantConfig(), testSchemaConfig(), testSchemaTemplateVrfL3MulticastConfig())
}

func testAccMSOSchemaTemplateBdConfigCreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_bd" "%[2]s" {
		schema_id              = mso_schema.%[3]s.id
		template_name          = "%[4]s"
		name                   = "%[2]s"
		display_name           = "%[2]s"
		description            = "Terraform test BD"
		layer2_unknown_unicast = "flood"
		layer2_stretch         = true
		intersite_bum_traffic  = true
		arp_flooding           = true
		vrf_name               = mso_schema_template_vrf.%[5]s.name
	}`, testAccMSOSchemaTemplateBdPrerequisiteConfig(), msoSchemaTemplateBdName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfL3MulticastName)
}

func testAccMSOSchemaTemplateBdConfigUpdateL2UnknownUnicast() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_bd" "%[2]s" {
		schema_id              = mso_schema.%[3]s.id
		template_name          = "%[4]s"
		name                   = "%[2]s"
		display_name           = "%[2]s"
		description            = "Terraform test BD"
		layer2_unknown_unicast = "proxy"
		vrf_name               = mso_schema_template_vrf.%[5]s.name
	}`, testAccMSOSchemaTemplateBdPrerequisiteConfig(), msoSchemaTemplateBdName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfL3MulticastName)
}

func testAccMSOSchemaTemplateBdConfigUpdateDescriptionAndL3Multicast() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_bd" "%[2]s" {
		schema_id              = mso_schema.%[3]s.id
		template_name          = "%[4]s"
		name                   = "%[2]s"
		display_name           = "%[2]s"
		description            = "Terraform test BD updated"
		layer2_unknown_unicast = "proxy"
		layer3_multicast       = true
		vrf_name               = mso_schema_template_vrf.%[5]s.name
	}`, testAccMSOSchemaTemplateBdPrerequisiteConfig(), msoSchemaTemplateBdName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfL3MulticastName)
}

func testAccMSOSchemaTemplateBdConfigUpdateFloodingFields() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_bd" "%[2]s" {
		schema_id                       = mso_schema.%[3]s.id
		template_name                   = "%[4]s"
		name                            = "%[2]s"
		display_name                    = "%[2]s"
		description                     = "Terraform test BD updated"
		layer2_unknown_unicast          = "proxy"
		layer3_multicast                = true
		unknown_multicast_flooding      = "optimized_flooding"
		multi_destination_flooding      = "flood_in_encap"
		ipv6_unknown_multicast_flooding = "optimized_flooding"
		layer2_stretch                  = false
		intersite_bum_traffic           = false
		vrf_name                        = mso_schema_template_vrf.%[5]s.name
	}`, testAccMSOSchemaTemplateBdPrerequisiteConfig(), msoSchemaTemplateBdName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfL3MulticastName)
}

func testAccMSOSchemaTemplateBdConfigUpdateRemainingAttributes() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_bd" "%[2]s" {
		schema_id                       = mso_schema.%[3]s.id
		template_name                   = "%[4]s"
		name                            = "%[2]s"
		display_name                    = "%[2]s"
		description                     = "Terraform test BD updated"
		layer2_unknown_unicast          = "proxy"
		layer3_multicast                = true
		unknown_multicast_flooding      = "optimized_flooding"
		multi_destination_flooding      = "flood_in_encap"
		ipv6_unknown_multicast_flooding = "optimized_flooding"
		layer2_stretch                  = false
		intersite_bum_traffic           = false
		arp_flooding                    = true
		unicast_routing                 = true
		optimize_wan_bandwidth          = true
		virtual_mac_address             = "00:00:5E:00:01:01"
		ep_move_detection_mode          = "garp"
		vrf_name                        = mso_schema_template_vrf.%[5]s.name
	}`, testAccMSOSchemaTemplateBdPrerequisiteConfig(), msoSchemaTemplateBdName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfL3MulticastName)
}

func testAccMSOSchemaTemplateBdConfigRemoveDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_bd" "%[2]s" {
		schema_id                       = mso_schema.%[3]s.id
		template_name                   = "%[4]s"
		name                            = "%[2]s"
		display_name                    = "%[2]s"
		description                     = ""
		layer2_unknown_unicast          = "proxy"
		layer3_multicast                = true
		unknown_multicast_flooding      = "optimized_flooding"
		multi_destination_flooding      = "flood_in_encap"
		ipv6_unknown_multicast_flooding = "optimized_flooding"
		layer2_stretch                  = false
		intersite_bum_traffic           = false
		arp_flooding                    = true
		unicast_routing                 = false
		optimize_wan_bandwidth          = false
		virtual_mac_address             = ""
		ep_move_detection_mode          = "none"
		vrf_name                        = mso_schema_template_vrf.%[5]s.name
	}`, testAccMSOSchemaTemplateBdPrerequisiteConfig(), msoSchemaTemplateBdName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfL3MulticastName)
}

func testAccMSOSchemaTemplateBdConfigWithChildren() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_bd" "%[2]s" {
		schema_id                       = mso_schema.%[3]s.id
		template_name                   = "%[4]s"
		name                            = "%[2]s"
		display_name                    = "%[2]s"
		description                     = ""
		layer2_unknown_unicast          = "proxy"
		layer3_multicast                = true
		unknown_multicast_flooding      = "optimized_flooding"
		multi_destination_flooding      = "flood_in_bd"
		ipv6_unknown_multicast_flooding = "optimized_flooding"
		layer2_stretch                  = true
		intersite_bum_traffic           = true
		arp_flooding                    = true
		unicast_routing                 = false
		vrf_name                        = mso_schema_template_vrf.%[5]s.name
	}%[6]s`, testAccMSOSchemaTemplateBdPrerequisiteConfig(), msoSchemaTemplateBdName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfL3MulticastName, testSchemaTemplateBdSubnetConfig())
}

func testAccMSOSchemaTemplateBdConfigWithChildrenUpdateDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_bd" "%[2]s" {
		schema_id                       = mso_schema.%[3]s.id
		template_name                   = "%[4]s"
		name                            = "%[2]s"
		display_name                    = "%[2]s"
		description                     = "Terraform test BD with children"
		layer2_unknown_unicast          = "proxy"
		layer3_multicast                = true
		unknown_multicast_flooding      = "optimized_flooding"
		multi_destination_flooding      = "flood_in_bd"
		ipv6_unknown_multicast_flooding = "optimized_flooding"
		layer2_stretch                  = true
		intersite_bum_traffic           = true
		arp_flooding                    = true
		unicast_routing                 = false
		vrf_name                        = mso_schema_template_vrf.%[5]s.name
	}%[6]s`, testAccMSOSchemaTemplateBdPrerequisiteConfig(), msoSchemaTemplateBdName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfL3MulticastName, testSchemaTemplateBdSubnetConfig())
}

func testAccMSOSchemaTemplateBdConfigRemoveChildren() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_bd" "%[2]s" {
		schema_id                       = mso_schema.%[3]s.id
		template_name                   = "%[4]s"
		name                            = "%[2]s"
		display_name                    = "%[2]s"
		description                     = "Terraform test BD with children"
		layer2_unknown_unicast          = "proxy"
		layer3_multicast                = true
		unknown_multicast_flooding      = "optimized_flooding"
		multi_destination_flooding      = "flood_in_bd"
		ipv6_unknown_multicast_flooding = "optimized_flooding"
		layer2_stretch                  = true
		intersite_bum_traffic           = true
		arp_flooding                    = true
		unicast_routing                 = false
		vrf_name                        = mso_schema_template_vrf.%[5]s.name
	}`, testAccMSOSchemaTemplateBdPrerequisiteConfig(), msoSchemaTemplateBdName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfL3MulticastName)
}

func testAccMSOSchemaTemplateBdDhcpPrerequisiteConfig() string {
	return fmt.Sprintf(`%[1]s%[2]s%[3]s%[4]s%[5]s`,
		testSchemaTemplateAnpConfig(),
		testSchemaTemplateAnpEpgConfig(),
		testSchemaTemplateVrfConfig(),
		testSchemaTemplateExtEpgConfig(),
		testTenantPolicyTemplateConfig(),
	)
}

func testAccMSOSchemaTemplateBdDhcpRelayPolicyConfig() string {
	return fmt.Sprintf(`
resource "mso_tenant_policies_dhcp_relay_policy" "%[1]s" {
	name        = "%[1]s"
	template_id = mso_template.%[3]s.id
	dhcp_relay_providers {
		dhcp_server_address  = "1.1.1.1"
		application_epg_uuid = mso_schema_template_anp_epg.%[4]s.uuid
	}
}
resource "mso_tenant_policies_dhcp_relay_policy" "%[2]s" {
	depends_on  = [mso_tenant_policies_dhcp_relay_policy.%[1]s]
	name        = "%[2]s"
	template_id = mso_template.%[3]s.id
	dhcp_relay_providers {
		dhcp_server_address  = "2.2.2.2"
		application_epg_uuid = mso_schema_template_anp_epg.%[4]s.uuid
	}
}
`, msoTenantPoliciesDhcpRelayPolicyName, msoTenantPoliciesDhcpRelayPolicyName2, msoTenantPolicyTemplateName, msoSchemaTemplateAnpEpgName)
}

func testAccMSOSchemaTemplateBdDhcpOptionPolicyConfig() string {
	return fmt.Sprintf(`
resource "mso_tenant_policies_dhcp_option_policy" "%[1]s" {
	template_id = mso_template.%[2]s.id
	name        = "%[1]s"
	options {
		name = "option_1"
		id   = 1
		data = "data_1"
	}
}
`, msoTenantPoliciesDhcpOptionPolicyName, msoTenantPolicyTemplateName)
}

func testAccMSOSchemaTemplateBdConfigWithOneDhcpPolicy() string {
	return fmt.Sprintf(`%[1]s%[2]s%[3]s%[4]s
	resource "mso_schema_template_bd" "%[5]s" {
		schema_id                       = mso_schema.%[6]s.id
		template_name                   = "%[7]s"
		name                            = "%[5]s"
		display_name                    = "%[5]s"
		description                     = "Terraform test BD with children"
		layer2_unknown_unicast          = "proxy"
		layer3_multicast                = true
		unknown_multicast_flooding      = "optimized_flooding"
		multi_destination_flooding      = "flood_in_bd"
		ipv6_unknown_multicast_flooding = "optimized_flooding"
		layer2_stretch                  = true
		intersite_bum_traffic           = true
		arp_flooding                    = true
		unicast_routing                 = false
		vrf_name                        = mso_schema_template_vrf.%[8]s.name
		dhcp_policies {
			name                    = mso_tenant_policies_dhcp_relay_policy.%[9]s.name
			dhcp_option_policy_name = mso_tenant_policies_dhcp_option_policy.%[11]s.name
		}
	}`, testAccMSOSchemaTemplateBdPrerequisiteConfig(), testAccMSOSchemaTemplateBdDhcpPrerequisiteConfig(), testAccMSOSchemaTemplateBdDhcpRelayPolicyConfig(), testAccMSOSchemaTemplateBdDhcpOptionPolicyConfig(),
		msoSchemaTemplateBdName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfL3MulticastName,
		msoTenantPoliciesDhcpRelayPolicyName, msoTenantPoliciesDhcpRelayPolicyName2, msoTenantPoliciesDhcpOptionPolicyName)
}

func testAccMSOSchemaTemplateBdConfigWithTwoDhcpPolicies() string {
	return fmt.Sprintf(`%[1]s%[2]s%[3]s%[4]s
	resource "mso_schema_template_bd" "%[5]s" {
		schema_id                       = mso_schema.%[6]s.id
		template_name                   = "%[7]s"
		name                            = "%[5]s"
		display_name                    = "%[5]s"
		description                     = "Terraform test BD with children"
		layer2_unknown_unicast          = "proxy"
		layer3_multicast                = true
		unknown_multicast_flooding      = "optimized_flooding"
		multi_destination_flooding      = "flood_in_bd"
		ipv6_unknown_multicast_flooding = "optimized_flooding"
		layer2_stretch                  = true
		intersite_bum_traffic           = true
		arp_flooding                    = true
		unicast_routing                 = false
		vrf_name                        = mso_schema_template_vrf.%[8]s.name
		dhcp_policies {
			name                    = mso_tenant_policies_dhcp_relay_policy.%[9]s.name
			dhcp_option_policy_name = mso_tenant_policies_dhcp_option_policy.%[11]s.name
		}
		dhcp_policies {
			name = mso_tenant_policies_dhcp_relay_policy.%[10]s.name
		}
	}`, testAccMSOSchemaTemplateBdPrerequisiteConfig(), testAccMSOSchemaTemplateBdDhcpPrerequisiteConfig(), testAccMSOSchemaTemplateBdDhcpRelayPolicyConfig(), testAccMSOSchemaTemplateBdDhcpOptionPolicyConfig(),
		msoSchemaTemplateBdName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfL3MulticastName,
		msoTenantPoliciesDhcpRelayPolicyName, msoTenantPoliciesDhcpRelayPolicyName2, msoTenantPoliciesDhcpOptionPolicyName)
}

func testAccMSOSchemaTemplateBdConfigRemoveDhcpPolicies() string {
	return fmt.Sprintf(`%[1]s%[2]s%[3]s%[4]s
	resource "mso_schema_template_bd" "%[5]s" {
		schema_id                       = mso_schema.%[6]s.id
		template_name                   = "%[7]s"
		name                            = "%[5]s"
		display_name                    = "%[5]s"
		description                     = "Terraform test BD with children"
		layer2_unknown_unicast          = "proxy"
		layer3_multicast                = true
		unknown_multicast_flooding      = "optimized_flooding"
		multi_destination_flooding      = "flood_in_bd"
		ipv6_unknown_multicast_flooding = "optimized_flooding"
		layer2_stretch                  = true
		intersite_bum_traffic           = true
		arp_flooding                    = true
		unicast_routing                 = false
		vrf_name                        = mso_schema_template_vrf.%[8]s.name
	}`, testAccMSOSchemaTemplateBdPrerequisiteConfig(), testAccMSOSchemaTemplateBdDhcpPrerequisiteConfig(), testAccMSOSchemaTemplateBdDhcpRelayPolicyConfig(), testAccMSOSchemaTemplateBdDhcpOptionPolicyConfig(),
		msoSchemaTemplateBdName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfL3MulticastName)
}
