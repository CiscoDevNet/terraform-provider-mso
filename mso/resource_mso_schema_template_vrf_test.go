package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// msoSchemaTemplateVrfSchemaId is set during the first test step's Check to capture the dynamic schema ID for use in the manual deletion PreConfig step.
var msoSchemaTemplateVrfSchemaId string

func TestAccMSOSchemaTemplateVrfResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateVrfDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create VRF with name and display_name") },
				Config:    testAccMSOSchemaTemplateVrfConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "template", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "display_name", msoSchemaTemplateVrfName),
					// Capture the dynamic schema ID from state for use in the manual deletion PreConfig step
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["mso_schema_template_vrf."+msoSchemaTemplateVrfName]
						if !ok {
							return fmt.Errorf("VRF resource not found in state")
						}
						msoSchemaTemplateVrfSchemaId = rs.Primary.Attributes["schema_id"]
						return nil
					},
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "description", ""),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "layer3_multicast", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "vzany", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "preferred_group", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "ip_data_plane_learning", "enabled"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "site_aware_policy_enforcement", "false"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update VRF display_name") },
				Config:    testAccMSOSchemaTemplateVrfConfigUpdateDisplayName(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "template", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "display_name", msoSchemaTemplateVrfName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "description", ""),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add VRF description") },
				Config:    testAccMSOSchemaTemplateVrfConfigAddDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "template", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "display_name", msoSchemaTemplateVrfName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "description", "Terraform test VRF"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove VRF description") },
				Config:    testAccMSOSchemaTemplateVrfConfigRemoveDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "template", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "display_name", msoSchemaTemplateVrfName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "description", ""),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Enable Layer3 Multicast and add RP") },
				Config:    testAccMSOSchemaTemplateVrfConfigEnableLayer3Multicast(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "template", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "display_name", msoSchemaTemplateVrfName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "layer3_multicast", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "rendezvous_points.#", "1"),
					customTestCheckResourceTypeSetAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "rendezvous_points",
						map[string]string{
							"ip_address": "1.1.1.2",
							"type":       "static",
						},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add extra RP") },
				Config:    testAccMSOSchemaTemplateVrfConfigAddExtraRp(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "template", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "display_name", msoSchemaTemplateVrfName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "layer3_multicast", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "rendezvous_points.#", "2"),
					customTestCheckResourceTypeSetAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "rendezvous_points",
						map[string]string{
							"ip_address": "1.1.1.2",
							"type":       "static",
						},
					),
					customTestCheckResourceTypeSetAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "rendezvous_points",
						map[string]string{
							"ip_address": "1.1.1.3",
							"type":       "fabric",
						},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove extra RP") },
				Config:    testAccMSOSchemaTemplateVrfConfigRemoveExtraRp(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "template", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "display_name", msoSchemaTemplateVrfName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "layer3_multicast", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "rendezvous_points.#", "1"),
					customTestCheckResourceTypeSetAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "rendezvous_points",
						map[string]string{
							"ip_address": "1.1.1.2",
							"type":       "static",
						},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove all RPs") },
				Config:    testAccMSOSchemaTemplateVrfConfigRemoveAllRps(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "template", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "display_name", msoSchemaTemplateVrfName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "layer3_multicast", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "rendezvous_points.#", "0"),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import VRF") },
				ResourceName: "mso_schema_template_vrf." + msoSchemaTemplateVrfName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["mso_schema_template_vrf."+msoSchemaTemplateVrfName]
					if !ok {
						return "", fmt.Errorf("VRF resource not found in state")
					}
					return fmt.Sprintf("%s/templates/%s/vrfs/%s", rs.Primary.Attributes["schema_id"], rs.Primary.Attributes["template"], rs.Primary.Attributes["name"]), nil
				},
				ImportStateVerify: true,
				// Description attribute is set to empty string on import but it is not provided in the config.
				ImportStateVerifyIgnore: []string{"description"},
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate VRF after manual deletion from NDO")
					msoClient := testAccProvider.Meta().(*client.Client)
					vrfRemovePatchPayload := models.GetRemovePatchPayload(fmt.Sprintf("/templates/%s/vrfs/%s", msoSchemaTemplateName, msoSchemaTemplateVrfName))
					_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", msoSchemaTemplateVrfSchemaId), vrfRemovePatchPayload)
					if err != nil {
						t.Fatalf("Failed to manually delete VRF: %v", err)
					}
				},
				Config: testAccMSOSchemaTemplateVrfConfigRemoveAllRps(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "display_name", msoSchemaTemplateVrfName+" updated"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add BD child to VRF") },
				Config:    testAccMSOSchemaTemplateVrfConfigWithBdChild(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "display_name", msoSchemaTemplateVrfName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "display_name", msoSchemaTemplateBdName),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update VRF description with BD child present") },
				Config:    testAccMSOSchemaTemplateVrfConfigWithBdChildUpdateDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "display_name", msoSchemaTemplateVrfName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "description", "Terraform test VRF with BD"),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "display_name", msoSchemaTemplateBdName),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateVrfRouteMapPolicyMulticastConfig() string {
	return fmt.Sprintf(`
	resource "mso_tenant_policies_route_map_policy_multicast" "route_map_policy_multicast" {
		template_id = mso_template.%[1]s.id
		name        = "tf_test_route_map_policy_multicast"
		description = "Terraform test Route Map Policy for Multicast"
		route_map_multicast_entries {
			order               = 1
			group_ip            = "226.2.2.2/8"
			source_ip           = "1.1.1.1/1"
			rendezvous_point_ip = "1.1.1.2"
			action              = "permit"
		}
	}`, msoTenantPolicyTemplateName)
}

func testAccMSOSchemaTemplateVrfPrerequisiteConfig() string {
	return fmt.Sprintf(`%s%s%s%s%s`, testSiteConfigAnsibleTest(), testTenantConfig(), testSchemaConfig(), testTenantPolicyTemplateConfig(), testAccMSOSchemaTemplateVrfRouteMapPolicyMulticastConfig())
}

func testAccMSOSchemaTemplateVrfConfigCreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_vrf" "%[2]s" {
		schema_id    = mso_schema.%[3]s.id
		template     = "%[4]s"
		name         = "%[2]s"
		display_name = "%[2]s"
	}`, testAccMSOSchemaTemplateVrfPrerequisiteConfig(), msoSchemaTemplateVrfName, msoSchemaName, msoSchemaTemplateName)
}

func testAccMSOSchemaTemplateVrfConfigUpdateDisplayName() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_vrf" "%[2]s" {
		schema_id    = mso_schema.%[3]s.id
		template     = "%[4]s"
		name         = "%[2]s"
		display_name = "%[2]s updated"
	}`, testAccMSOSchemaTemplateVrfPrerequisiteConfig(), msoSchemaTemplateVrfName, msoSchemaName, msoSchemaTemplateName)
}

func testAccMSOSchemaTemplateVrfConfigAddDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_vrf" "%[2]s" {
		schema_id    = mso_schema.%[3]s.id
		template     = "%[4]s"
		name         = "%[2]s"
		display_name = "%[2]s updated"
		description  = "Terraform test VRF"
	}`, testAccMSOSchemaTemplateVrfPrerequisiteConfig(), msoSchemaTemplateVrfName, msoSchemaName, msoSchemaTemplateName)
}

func testAccMSOSchemaTemplateVrfConfigRemoveDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_vrf" "%[2]s" {
		schema_id    = mso_schema.%[3]s.id
		template     = "%[4]s"
		name         = "%[2]s"
		display_name = "%[2]s updated"
		description  = ""
	}`, testAccMSOSchemaTemplateVrfPrerequisiteConfig(), msoSchemaTemplateVrfName, msoSchemaName, msoSchemaTemplateName)
}

func testAccMSOSchemaTemplateVrfConfigEnableLayer3Multicast() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_vrf" "%[2]s" {
		schema_id        = mso_schema.%[3]s.id
		template         = "%[4]s"
		name             = "%[2]s"
		display_name     = "%[2]s updated"
		description      = ""
		layer3_multicast = true
		rendezvous_points {
			ip_address                      = "1.1.1.2"
			type                            = "static"
			route_map_policy_multicast_uuid = mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast.uuid
		}
	}`, testAccMSOSchemaTemplateVrfPrerequisiteConfig(), msoSchemaTemplateVrfName, msoSchemaName, msoSchemaTemplateName)
}

func testAccMSOSchemaTemplateVrfConfigAddExtraRp() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_vrf" "%[2]s" {
		schema_id        = mso_schema.%[3]s.id
		template         = "%[4]s"
		name             = "%[2]s"
		display_name     = "%[2]s updated"
		description      = ""
		layer3_multicast = true
		rendezvous_points {
			ip_address                      = "1.1.1.2"
			type                            = "static"
			route_map_policy_multicast_uuid = mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast.uuid
		}
		rendezvous_points {
			ip_address = "1.1.1.3"
			type       = "fabric"
		}
	}`, testAccMSOSchemaTemplateVrfPrerequisiteConfig(), msoSchemaTemplateVrfName, msoSchemaName, msoSchemaTemplateName)
}

func testAccMSOSchemaTemplateVrfConfigRemoveExtraRp() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_vrf" "%[2]s" {
		schema_id        = mso_schema.%[3]s.id
		template         = "%[4]s"
		name             = "%[2]s"
		display_name     = "%[2]s updated"
		description      = ""
		layer3_multicast = true
		rendezvous_points {
			ip_address                      = "1.1.1.2"
			type                            = "static"
			route_map_policy_multicast_uuid = mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast.uuid
		}
	}`, testAccMSOSchemaTemplateVrfPrerequisiteConfig(), msoSchemaTemplateVrfName, msoSchemaName, msoSchemaTemplateName)
}

func testAccMSOSchemaTemplateVrfConfigRemoveAllRps() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_vrf" "%[2]s" {
		schema_id        = mso_schema.%[3]s.id
		template         = "%[4]s"
		name             = "%[2]s"
		display_name     = "%[2]s updated"
		description      = ""
		layer3_multicast = true
	}`, testAccMSOSchemaTemplateVrfPrerequisiteConfig(), msoSchemaTemplateVrfName, msoSchemaName, msoSchemaTemplateName)
}

func testAccMSOSchemaTemplateVrfConfigWithBdChild() string {
	return testAccMSOSchemaTemplateVrfConfigUpdateDisplayName() + testSchemaTemplateBdConfig()
}

func testAccMSOSchemaTemplateVrfConfigWithBdChildUpdateDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_vrf" "%[2]s" {
		schema_id    = mso_schema.%[3]s.id
		template     = "%[4]s"
		name         = "%[2]s"
		display_name = "%[2]s updated"
		description  = "Terraform test VRF with BD"
	}`, testAccMSOSchemaTemplateVrfPrerequisiteConfig(), msoSchemaTemplateVrfName, msoSchemaName, msoSchemaTemplateName) + testSchemaTemplateBdConfig()
}

func testAccCheckMSOSchemaTemplateVrfDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "mso_schema_template_vrf" {
			schemaID := rs.Primary.Attributes["schema_id"]
			con, err := client.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
			if err != nil {
				return nil
			}
			count, err := con.ArrayCount("templates")
			if err != nil {
				return fmt.Errorf("No Template found")
			}
			for i := 0; i < count; i++ {
				tempCont, err := con.ArrayElement(i, "templates")
				if err != nil {
					return fmt.Errorf("No template exists")
				}
				vrfCount, err := tempCont.ArrayCount("vrfs")
				if err != nil {
					return fmt.Errorf("No Vrf found")
				}
				for j := 0; j < vrfCount; j++ {
					vrfCont, err := tempCont.ArrayElement(j, "vrfs")
					if err != nil {
						return err
					}
					name := models.StripQuotes(vrfCont.S("name").String())
					if rs.Primary.ID == name {
						return fmt.Errorf("Schema Template Vrf record still exists")
					}
				}
			}
		}
	}
	return nil
}
