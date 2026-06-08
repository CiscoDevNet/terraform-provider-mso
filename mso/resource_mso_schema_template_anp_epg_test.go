package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// msoSchemaTemplateAnpEpgSchemaId is set during the first test step's Check to capture the dynamic schema ID for use in the manual deletion PreConfig step.
var msoSchemaTemplateAnpEpgSchemaId string

func TestAccMSOSchemaTemplateAnpEpgResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateAnpEpgDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create EPG with name and display_name") },
				Config:    testAccMSOSchemaTemplateAnpEpgConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "anp_name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "display_name", msoSchemaTemplateAnpEpgName),
					// Capture the dynamic schema ID from state for use in the manual deletion PreConfig step
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName]
						if !ok {
							return fmt.Errorf("EPG resource not found in state")
						}
						msoSchemaTemplateAnpEpgSchemaId = rs.Primary.Attributes["schema_id"]
						return nil
					},
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "description", ""),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "useg_epg", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "intersite_multicast_source", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "proxy_arp", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "preferred_group", "false"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update EPG display_name") },
				Config:    testAccMSOSchemaTemplateAnpEpgConfigUpdateDisplayName(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "anp_name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "display_name", msoSchemaTemplateAnpEpgName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "description", ""),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add EPG description") },
				Config:    testAccMSOSchemaTemplateAnpEpgConfigAddDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "display_name", msoSchemaTemplateAnpEpgName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "description", "Terraform test EPG"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove EPG description") },
				Config:    testAccMSOSchemaTemplateAnpEpgConfigRemoveDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "display_name", msoSchemaTemplateAnpEpgName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "description", ""),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update EPG useg_epg, intra_epg, intersite_multicast_source, proxy_arp, preferred_group")
				},
				Config: testAccMSOSchemaTemplateAnpEpgConfigUpdateAttributes(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "display_name", msoSchemaTemplateAnpEpgName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "useg_epg", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "intra_epg", "enforced"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "intersite_multicast_source", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "proxy_arp", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "preferred_group", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "bd_name", msoSchemaTemplateBdL3MulticastName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "vrf_name", msoSchemaTemplateVrfL3MulticastName),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add BD and VRF to EPG") },
				Config:    testAccMSOSchemaTemplateAnpEpgConfigWithBdAndVrf(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "display_name", msoSchemaTemplateAnpEpgName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "bd_name", msoSchemaTemplateBdL3MulticastName),
					resource.TestCheckResourceAttrSet("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "bd_schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "bd_template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "vrf_name", msoSchemaTemplateVrfL3MulticastName),
					resource.TestCheckResourceAttrSet("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "vrf_schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "vrf_template_name", msoSchemaTemplateName),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import EPG") },
				ResourceName: "mso_schema_template_anp_epg." + msoSchemaTemplateAnpEpgName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName]
					if !ok {
						return "", fmt.Errorf("EPG resource not found in state")
					}
					return fmt.Sprintf("%s/templates/%s/anps/%s/epgs/%s", rs.Primary.Attributes["schema_id"], rs.Primary.Attributes["template_name"], rs.Primary.Attributes["anp_name"], rs.Primary.Attributes["name"]), nil
				},
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate EPG after manual deletion from NDO")
					msoClient := testAccProvider.Meta().(*client.Client)
					epgRemovePatchPayload := models.GetRemovePatchPayload(fmt.Sprintf("/templates/%s/anps/%s/epgs/%s", msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName))
					_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", msoSchemaTemplateAnpEpgSchemaId), epgRemovePatchPayload)
					if err != nil {
						t.Fatalf("Failed to manually delete EPG: %v", err)
					}
				},
				Config: testAccMSOSchemaTemplateAnpEpgConfigWithBdAndVrf(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "display_name", msoSchemaTemplateAnpEpgName+" updated"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add contract and subnet children to EPG") },
				Config:    testAccMSOSchemaTemplateAnpEpgConfigWithChildren(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "display_name", msoSchemaTemplateAnpEpgName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_contract."+msoSchemaTemplateContractName+"_provider", "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_contract."+msoSchemaTemplateContractName+"_provider", "relationship_type", "provider"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "ip", msoSchemaTemplateAnpEpgSubnetIp),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "scope", "private"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update EPG description with children present") },
				Config:    testAccMSOSchemaTemplateAnpEpgConfigWithChildrenUpdateDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "display_name", msoSchemaTemplateAnpEpgName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "description", "Terraform test EPG with children"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_contract."+msoSchemaTemplateContractName+"_provider", "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_contract."+msoSchemaTemplateContractName+"_provider", "relationship_type", "provider"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "ip", msoSchemaTemplateAnpEpgSubnetIp),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "scope", "private"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove children from EPG") },
				Config:    testAccMSOSchemaTemplateAnpEpgConfigRemoveChildren(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "bd_name", msoSchemaTemplateBdL3MulticastName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "vrf_name", msoSchemaTemplateVrfL3MulticastName),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove BD from EPG") },
				Config:    testAccMSOSchemaTemplateAnpEpgConfigRemoveBd(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "display_name", msoSchemaTemplateAnpEpgName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "bd_name", ""),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "vrf_name", msoSchemaTemplateVrfL3MulticastName),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove VRF from EPG") },
				Config:    testAccMSOSchemaTemplateAnpEpgConfigRemoveVrf(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "display_name", msoSchemaTemplateAnpEpgName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "bd_name", ""),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "vrf_name", ""),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateAnpEpgPrerequisiteConfig() string {
	return fmt.Sprintf(`%s%s%s%s`, testSiteConfigAnsibleTest(), testTenantConfig(), testSchemaConfig(), testSchemaTemplateAnpConfig())
}

func testAccMSOSchemaTemplateAnpEpgPrerequisiteWithBdVrfConfig() string {
	return fmt.Sprintf(`%s%s%s`, testAccMSOSchemaTemplateAnpEpgPrerequisiteConfig(), testSchemaTemplateVrfConfig(), testSchemaTemplateBdConfig())
}

func testAccMSOSchemaTemplateAnpEpgConfigCreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp_epg" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		template_name = "%[4]s"
		anp_name      = mso_schema_template_anp.%[5]s.name
		name          = "%[2]s"
		display_name  = "%[2]s"
	}`, testAccMSOSchemaTemplateAnpEpgPrerequisiteConfig(), msoSchemaTemplateAnpEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName)
}

func testAccMSOSchemaTemplateAnpEpgConfigUpdateDisplayName() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp_epg" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		template_name = "%[4]s"
		anp_name      = mso_schema_template_anp.%[5]s.name
		name          = "%[2]s"
		display_name  = "%[2]s updated"
	}`, testAccMSOSchemaTemplateAnpEpgPrerequisiteConfig(), msoSchemaTemplateAnpEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName)
}

func testAccMSOSchemaTemplateAnpEpgConfigAddDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp_epg" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		template_name = "%[4]s"
		anp_name      = mso_schema_template_anp.%[5]s.name
		name          = "%[2]s"
		display_name  = "%[2]s updated"
		description   = "Terraform test EPG"
	}`, testAccMSOSchemaTemplateAnpEpgPrerequisiteConfig(), msoSchemaTemplateAnpEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName)
}

func testAccMSOSchemaTemplateAnpEpgConfigRemoveDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp_epg" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		template_name = "%[4]s"
		anp_name      = mso_schema_template_anp.%[5]s.name
		name          = "%[2]s"
		display_name  = "%[2]s updated"
		description   = ""
	}`, testAccMSOSchemaTemplateAnpEpgPrerequisiteConfig(), msoSchemaTemplateAnpEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName)
}

func testAccMSOSchemaTemplateAnpEpgPrerequisiteWithBdL3MulticastConfig() string {
	return fmt.Sprintf(`%s%s%s`, testAccMSOSchemaTemplateAnpEpgPrerequisiteConfig(), testSchemaTemplateVrfL3MulticastConfig(), testSchemaTemplateBdL3MulticastConfig())
}

func testAccMSOSchemaTemplateAnpEpgConfigUpdateAttributes() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp_epg" "%[2]s" {
		schema_id                  = mso_schema.%[3]s.id
		template_name              = "%[4]s"
		anp_name                   = mso_schema_template_anp.%[5]s.name
		name                       = "%[2]s"
		display_name               = "%[2]s updated"
		useg_epg                   = true
		intra_epg                  = "enforced"
		intersite_multicast_source = true
		proxy_arp                  = true
		preferred_group            = true
		bd_name                    = mso_schema_template_bd.%[6]s.name
		vrf_name                   = mso_schema_template_vrf.%[7]s.name
	}`, testAccMSOSchemaTemplateAnpEpgPrerequisiteWithBdL3MulticastConfig(), msoSchemaTemplateAnpEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateBdL3MulticastName, msoSchemaTemplateVrfL3MulticastName)
}

func testAccMSOSchemaTemplateAnpEpgConfigWithBdAndVrf() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp_epg" "%[2]s" {
		schema_id                  = mso_schema.%[3]s.id
		template_name              = "%[4]s"
		anp_name                   = mso_schema_template_anp.%[5]s.name
		name                       = "%[2]s"
		display_name               = "%[2]s updated"
		bd_name                    = mso_schema_template_bd.%[6]s.name
		vrf_name                   = mso_schema_template_vrf.%[7]s.name
		preferred_group            = false
		intersite_multicast_source = false
		proxy_arp                  = false
		useg_epg                   = false
		intra_epg                  = "unenforced"
	}`, testAccMSOSchemaTemplateAnpEpgPrerequisiteWithBdL3MulticastConfig(), msoSchemaTemplateAnpEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateBdL3MulticastName, msoSchemaTemplateVrfL3MulticastName)
}

func testAccMSOSchemaTemplateAnpEpgChildrenConfig() string {
	return fmt.Sprintf(`%s%s%s`, testSchemaTemplateFilterEntryConfig(), testSchemaTemplateContractConfig(), testSchemaTemplateAnpEpgContractConfig()) + testSchemaTemplateAnpEpgSubnetConfig()
}

func testAccMSOSchemaTemplateAnpEpgConfigWithChildren() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp_epg" "%[2]s" {
		schema_id                  = mso_schema.%[3]s.id
		template_name              = "%[4]s"
		anp_name                   = mso_schema_template_anp.%[5]s.name
		name                       = "%[2]s"
		display_name               = "%[2]s updated"
		bd_name                    = mso_schema_template_bd.%[6]s.name
		vrf_name                   = mso_schema_template_vrf.%[7]s.name
		preferred_group            = false
		intersite_multicast_source = false
	}`, testAccMSOSchemaTemplateAnpEpgPrerequisiteWithBdL3MulticastConfig(), msoSchemaTemplateAnpEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateBdL3MulticastName, msoSchemaTemplateVrfL3MulticastName) + testAccMSOSchemaTemplateAnpEpgChildrenConfig()
}

func testAccMSOSchemaTemplateAnpEpgConfigWithChildrenUpdateDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp_epg" "%[2]s" {
		schema_id                  = mso_schema.%[3]s.id
		template_name              = "%[4]s"
		anp_name                   = mso_schema_template_anp.%[5]s.name
		name                       = "%[2]s"
		display_name               = "%[2]s updated"
		description                = "Terraform test EPG with children"
		bd_name                    = mso_schema_template_bd.%[6]s.name
		vrf_name                   = mso_schema_template_vrf.%[7]s.name
		preferred_group            = false
		intersite_multicast_source = false
	}`, testAccMSOSchemaTemplateAnpEpgPrerequisiteWithBdL3MulticastConfig(), msoSchemaTemplateAnpEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateBdL3MulticastName, msoSchemaTemplateVrfL3MulticastName) + testAccMSOSchemaTemplateAnpEpgChildrenConfig()
}

func testAccMSOSchemaTemplateAnpEpgConfigRemoveChildren() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp_epg" "%[2]s" {
		schema_id                  = mso_schema.%[3]s.id
		template_name              = "%[4]s"
		anp_name                   = mso_schema_template_anp.%[5]s.name
		name                       = "%[2]s"
		display_name               = "%[2]s updated"
		description                = "Terraform test EPG with children"
		bd_name                    = mso_schema_template_bd.%[6]s.name
		vrf_name                   = mso_schema_template_vrf.%[7]s.name
		preferred_group            = false
		intersite_multicast_source = false
	}`, testAccMSOSchemaTemplateAnpEpgPrerequisiteWithBdL3MulticastConfig(), msoSchemaTemplateAnpEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateBdL3MulticastName, msoSchemaTemplateVrfL3MulticastName)
}

func testAccMSOSchemaTemplateAnpEpgConfigRemoveBd() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp_epg" "%[2]s" {
		schema_id                  = mso_schema.%[3]s.id
		template_name              = "%[4]s"
		anp_name                   = mso_schema_template_anp.%[5]s.name
		name                       = "%[2]s"
		display_name               = "%[2]s updated"
		description                = "Terraform test EPG with children"
		bd_name                    = ""
		vrf_name                   = mso_schema_template_vrf.%[6]s.name
		preferred_group            = false
		intersite_multicast_source = false
		proxy_arp                  = false
		useg_epg                   = false
		intra_epg                  = "unenforced"
	}`, testAccMSOSchemaTemplateAnpEpgPrerequisiteWithBdL3MulticastConfig(), msoSchemaTemplateAnpEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateVrfL3MulticastName)
}

func testAccMSOSchemaTemplateAnpEpgConfigRemoveVrf() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp_epg" "%[2]s" {
		schema_id                  = mso_schema.%[3]s.id
		template_name              = "%[4]s"
		anp_name                   = mso_schema_template_anp.%[5]s.name
		name                       = "%[2]s"
		display_name               = "%[2]s updated"
		description                = "Terraform test EPG with children"
		bd_name                    = ""
		vrf_name                   = ""
		preferred_group            = false
		intersite_multicast_source = false
		proxy_arp                  = false
		useg_epg                   = false
		intra_epg                  = "unenforced"
	}`, testAccMSOSchemaTemplateAnpEpgPrerequisiteWithBdL3MulticastConfig(), msoSchemaTemplateAnpEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName)
}

func testAccCheckMSOSchemaTemplateAnpEpgDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "mso_schema_template_anp_epg" {
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
				anpCount, err := tempCont.ArrayCount("anps")
				if err != nil {
					return fmt.Errorf("No Anp found")
				}
				for j := 0; j < anpCount; j++ {
					anpCont, err := tempCont.ArrayElement(j, "anps")
					if err != nil {
						return err
					}
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return fmt.Errorf("Unable to get EPG list")
					}
					for k := 0; k < epgCount; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return err
						}
						name := models.StripQuotes(epgCont.S("name").String())
						if rs.Primary.Attributes["name"] == name {
							return fmt.Errorf("Schema Template Anp Epg record still exists")
						}
					}
				}
			}
		}
	}
	return nil
}
