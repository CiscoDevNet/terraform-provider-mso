package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// msoSchemaTemplateExtEpgSchemaId is set during the first test step's Check to capture the dynamic schema ID for use in the manual deletion PreConfig step.
var msoSchemaTemplateExtEpgSchemaId string

func TestAccMSOSchemaTemplateExternalEpgResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateExtEpgDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create External EPG with name and display_name") },
				Config:    testAccMSOSchemaTemplateExtEpgConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "external_epg_name", msoSchemaTemplateExtEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "display_name", msoSchemaTemplateExtEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "vrf_name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttrSet("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "vrf_schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "vrf_template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "external_epg_type", "on-premise"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "description", ""),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "include_in_preferred_group", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "l3out_name", ""),
					// Capture the dynamic schema ID from state for use in the manual deletion PreConfig step
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName]
						if !ok {
							return fmt.Errorf("External EPG resource not found in state")
						}
						msoSchemaTemplateExtEpgSchemaId = rs.Primary.Attributes["schema_id"]
						return nil
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update External EPG display_name") },
				Config:    testAccMSOSchemaTemplateExtEpgConfigUpdateDisplayName(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "external_epg_name", msoSchemaTemplateExtEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "display_name", msoSchemaTemplateExtEpgName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "description", ""),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add External EPG description") },
				Config:    testAccMSOSchemaTemplateExtEpgConfigAddDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "external_epg_name", msoSchemaTemplateExtEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "display_name", msoSchemaTemplateExtEpgName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "description", "Terraform test External EPG"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove External EPG description") },
				Config:    testAccMSOSchemaTemplateExtEpgConfigRemoveDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "external_epg_name", msoSchemaTemplateExtEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "display_name", msoSchemaTemplateExtEpgName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "description", ""),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Set include_in_preferred_group to true") },
				Config:    testAccMSOSchemaTemplateExtEpgConfigUpdatePreferredGroup(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "external_epg_name", msoSchemaTemplateExtEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "display_name", msoSchemaTemplateExtEpgName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "include_in_preferred_group", "true"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Reset include_in_preferred_group to false") },
				Config:    testAccMSOSchemaTemplateExtEpgConfigResetPreferredGroup(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "external_epg_name", msoSchemaTemplateExtEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "display_name", msoSchemaTemplateExtEpgName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "include_in_preferred_group", "false"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add L3Out to External EPG") },
				Config:    testAccMSOSchemaTemplateExtEpgConfigWithL3out(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "external_epg_name", msoSchemaTemplateExtEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "display_name", msoSchemaTemplateExtEpgName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "l3out_name", msoSchemaTemplateL3outName),
					resource.TestCheckResourceAttrSet("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "l3out_schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "l3out_template_name", msoSchemaTemplateName),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove L3Out from External EPG") },
				Config:    testAccMSOSchemaTemplateExtEpgConfigRemoveL3out(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "external_epg_name", msoSchemaTemplateExtEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "display_name", msoSchemaTemplateExtEpgName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "l3out_name", ""),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import External EPG") },
				ResourceName: "mso_schema_template_external_epg." + msoSchemaTemplateExtEpgName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName]
					if !ok {
						return "", fmt.Errorf("External EPG resource not found in state")
					}
					return fmt.Sprintf("%s/templates/%s/externalEpgs/%s", rs.Primary.Attributes["schema_id"], rs.Primary.Attributes["template_name"], rs.Primary.Attributes["external_epg_name"]), nil
				},
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate External EPG after manual deletion from NDO")
					msoClient := testAccProvider.Meta().(*client.Client)
					epgRemovePatchPayload := models.GetRemovePatchPayload(fmt.Sprintf("/templates/%s/externalEpgs/%s", msoSchemaTemplateName, msoSchemaTemplateExtEpgName))
					_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", msoSchemaTemplateExtEpgSchemaId), epgRemovePatchPayload)
					if err != nil {
						t.Fatalf("Failed to manually delete External EPG: %v", err)
					}
				},
				Config: testAccMSOSchemaTemplateExtEpgConfigRemoveL3out(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "external_epg_name", msoSchemaTemplateExtEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "display_name", msoSchemaTemplateExtEpgName+" updated"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add contract and subnet children to External EPG") },
				Config:    testAccMSOSchemaTemplateExtEpgConfigWithChildren(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "external_epg_name", msoSchemaTemplateExtEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "display_name", msoSchemaTemplateExtEpgName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg_contract."+msoSchemaTemplateContractName+"_provider", "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg_contract."+msoSchemaTemplateContractName+"_provider", "relationship_type", "provider"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg_subnet."+msoSchemaTemplateExtEpgName+"_subnet", "ip", msoSchemaTemplateExtEpgSubnetIp),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update External EPG description with children present") },
				Config:    testAccMSOSchemaTemplateExtEpgConfigWithChildrenUpdateDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "external_epg_name", msoSchemaTemplateExtEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "display_name", msoSchemaTemplateExtEpgName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "description", "Terraform test External EPG with children"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg_contract."+msoSchemaTemplateContractName+"_provider", "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg_contract."+msoSchemaTemplateContractName+"_provider", "relationship_type", "provider"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg_subnet."+msoSchemaTemplateExtEpgName+"_subnet", "ip", msoSchemaTemplateExtEpgSubnetIp),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove children from External EPG") },
				Config:    testAccMSOSchemaTemplateExtEpgConfigRemoveChildren(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "external_epg_name", msoSchemaTemplateExtEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "display_name", msoSchemaTemplateExtEpgName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg."+msoSchemaTemplateExtEpgName, "vrf_name", msoSchemaTemplateVrfName),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateExtEpgPrerequisiteConfig() string {
	return fmt.Sprintf(`%s%s%s%s`, testSiteConfigAnsibleTest(), testTenantConfig(), testSchemaConfig(), testSchemaTemplateVrfConfig())
}

func testAccMSOSchemaTemplateExtEpgPrerequisiteWithL3outConfig() string {
	return fmt.Sprintf(`%s%s`, testAccMSOSchemaTemplateExtEpgPrerequisiteConfig(), testSchemaTemplateL3outConfig())
}

func testAccMSOSchemaTemplateExtEpgPrerequisiteWithVrfPreferredGroupConfig() string {
	return fmt.Sprintf(`%[1]s%[2]s%[3]s
resource "mso_schema_template_vrf" "%[4]s" {
	name            = "%[4]s"
	display_name    = "%[4]s"
	schema_id       = mso_schema.%[5]s.id
	template        = "%[6]s"
	preferred_group = true
}
`, testSiteConfigAnsibleTest(), testTenantConfig(), testSchemaConfig(), msoSchemaTemplateVrfName, msoSchemaName, msoSchemaTemplateName)
}

func testAccMSOSchemaTemplateExtEpgConfigCreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_external_epg" "%[2]s" {
		schema_id         = mso_schema.%[3]s.id
		template_name     = "%[4]s"
		external_epg_name = "%[2]s"
		display_name      = "%[2]s"
		vrf_name          = mso_schema_template_vrf.%[5]s.name
	}`, testAccMSOSchemaTemplateExtEpgPrerequisiteConfig(), msoSchemaTemplateExtEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfName)
}

func testAccMSOSchemaTemplateExtEpgConfigUpdateDisplayName() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_external_epg" "%[2]s" {
		schema_id         = mso_schema.%[3]s.id
		template_name     = "%[4]s"
		external_epg_name = "%[2]s"
		display_name      = "%[2]s updated"
		vrf_name          = mso_schema_template_vrf.%[5]s.name
	}`, testAccMSOSchemaTemplateExtEpgPrerequisiteConfig(), msoSchemaTemplateExtEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfName)
}

func testAccMSOSchemaTemplateExtEpgConfigAddDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_external_epg" "%[2]s" {
		schema_id         = mso_schema.%[3]s.id
		template_name     = "%[4]s"
		external_epg_name = "%[2]s"
		display_name      = "%[2]s updated"
		vrf_name          = mso_schema_template_vrf.%[5]s.name
		description       = "Terraform test External EPG"
	}`, testAccMSOSchemaTemplateExtEpgPrerequisiteConfig(), msoSchemaTemplateExtEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfName)
}

func testAccMSOSchemaTemplateExtEpgConfigRemoveDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_external_epg" "%[2]s" {
		schema_id         = mso_schema.%[3]s.id
		template_name     = "%[4]s"
		external_epg_name = "%[2]s"
		display_name      = "%[2]s updated"
		vrf_name          = mso_schema_template_vrf.%[5]s.name
		description       = ""
	}`, testAccMSOSchemaTemplateExtEpgPrerequisiteConfig(), msoSchemaTemplateExtEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfName)
}

func testAccMSOSchemaTemplateExtEpgConfigUpdatePreferredGroup() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_external_epg" "%[2]s" {
		schema_id                  = mso_schema.%[3]s.id
		template_name              = "%[4]s"
		external_epg_name          = "%[2]s"
		display_name               = "%[2]s updated"
		vrf_name                   = mso_schema_template_vrf.%[5]s.name
		include_in_preferred_group = true
	}`, testAccMSOSchemaTemplateExtEpgPrerequisiteWithVrfPreferredGroupConfig(), msoSchemaTemplateExtEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfName)
}

func testAccMSOSchemaTemplateExtEpgConfigResetPreferredGroup() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_external_epg" "%[2]s" {
		schema_id                  = mso_schema.%[3]s.id
		template_name              = "%[4]s"
		external_epg_name          = "%[2]s"
		display_name               = "%[2]s updated"
		vrf_name                   = mso_schema_template_vrf.%[5]s.name
		include_in_preferred_group = false
	}`, testAccMSOSchemaTemplateExtEpgPrerequisiteConfig(), msoSchemaTemplateExtEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfName)
}

func testAccMSOSchemaTemplateExtEpgConfigWithL3out() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_external_epg" "%[2]s" {
		schema_id                  = mso_schema.%[3]s.id
		template_name              = "%[4]s"
		external_epg_name          = "%[2]s"
		display_name               = "%[2]s updated"
		vrf_name                   = mso_schema_template_vrf.%[5]s.name
		include_in_preferred_group = false
		l3out_name                 = mso_schema_template_l3out.%[6]s.l3out_name
	}`, testAccMSOSchemaTemplateExtEpgPrerequisiteWithL3outConfig(), msoSchemaTemplateExtEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfName, msoSchemaTemplateL3outName)
}

func testAccMSOSchemaTemplateExtEpgConfigRemoveL3out() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_external_epg" "%[2]s" {
		schema_id                  = mso_schema.%[3]s.id
		template_name              = "%[4]s"
		external_epg_name          = "%[2]s"
		display_name               = "%[2]s updated"
		vrf_name                   = mso_schema_template_vrf.%[5]s.name
		include_in_preferred_group = false
	}`, testAccMSOSchemaTemplateExtEpgPrerequisiteWithL3outConfig(), msoSchemaTemplateExtEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfName)
}

func testAccMSOSchemaTemplateExtEpgChildrenConfig() string {
	return fmt.Sprintf(`%s%s%s`, testSchemaTemplateFilterEntryConfig(), testSchemaTemplateContractConfig(), testSchemaTemplateExtEpgContractConfig()) + testSchemaTemplateExtEpgSubnetConfig()
}

func testAccMSOSchemaTemplateExtEpgConfigWithChildren() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_external_epg" "%[2]s" {
		schema_id                  = mso_schema.%[3]s.id
		template_name              = "%[4]s"
		external_epg_name          = "%[2]s"
		display_name               = "%[2]s updated"
		vrf_name                   = mso_schema_template_vrf.%[5]s.name
		include_in_preferred_group = false
	}`, testAccMSOSchemaTemplateExtEpgPrerequisiteConfig(), msoSchemaTemplateExtEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfName) + testAccMSOSchemaTemplateExtEpgChildrenConfig()
}

func testAccMSOSchemaTemplateExtEpgConfigWithChildrenUpdateDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_external_epg" "%[2]s" {
		schema_id                  = mso_schema.%[3]s.id
		template_name              = "%[4]s"
		external_epg_name          = "%[2]s"
		display_name               = "%[2]s updated"
		vrf_name                   = mso_schema_template_vrf.%[5]s.name
		description                = "Terraform test External EPG with children"
		include_in_preferred_group = false
	}`, testAccMSOSchemaTemplateExtEpgPrerequisiteConfig(), msoSchemaTemplateExtEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfName) + testAccMSOSchemaTemplateExtEpgChildrenConfig()
}

func testAccMSOSchemaTemplateExtEpgConfigRemoveChildren() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_external_epg" "%[2]s" {
		schema_id                  = mso_schema.%[3]s.id
		template_name              = "%[4]s"
		external_epg_name          = "%[2]s"
		display_name               = "%[2]s updated"
		vrf_name                   = mso_schema_template_vrf.%[5]s.name
		description                = "Terraform test External EPG with children"
		include_in_preferred_group = false
	}`, testAccMSOSchemaTemplateExtEpgPrerequisiteConfig(), msoSchemaTemplateExtEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfName)
}

func testAccCheckMSOSchemaTemplateExtEpgDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "mso_schema_template_external_epg" {
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
				externalEpgCount, err := tempCont.ArrayCount("externalEpgs")
				if err != nil {
					return fmt.Errorf("Unable to get External EPG list")
				}
				for j := 0; j < externalEpgCount; j++ {
					externalEpgCont, err := tempCont.ArrayElement(j, "externalEpgs")
					if err != nil {
						return err
					}
					name := models.StripQuotes(externalEpgCont.S("name").String())
					if rs.Primary.Attributes["external_epg_name"] == name {
						return fmt.Errorf("Schema Template External EPG record still exists")
					}
				}
			}
		}
	}
	return nil
}
