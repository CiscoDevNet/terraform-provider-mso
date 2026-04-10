package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// msoSchemaTemplateL3outSchemaId is set during the first test step's Check to capture the dynamic schema ID for use in the manual deletion PreConfig step.
var msoSchemaTemplateL3outSchemaId string

func TestAccMSOSchemaTemplateL3outResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateL3outDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create L3out with name and display_name") },
				Config:    testAccMSOSchemaTemplateL3outConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_l3out."+msoSchemaTemplateL3outName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_l3out."+msoSchemaTemplateL3outName, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_l3out."+msoSchemaTemplateL3outName, "l3out_name", msoSchemaTemplateL3outName),
					resource.TestCheckResourceAttr("mso_schema_template_l3out."+msoSchemaTemplateL3outName, "display_name", msoSchemaTemplateL3outName),
					resource.TestCheckResourceAttr("mso_schema_template_l3out."+msoSchemaTemplateL3outName, "vrf_name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttrSet("mso_schema_template_l3out."+msoSchemaTemplateL3outName, "vrf_schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_l3out."+msoSchemaTemplateL3outName, "vrf_template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_l3out."+msoSchemaTemplateL3outName, "description", ""),
					// Capture the dynamic schema ID from state for use in the manual deletion PreConfig step
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["mso_schema_template_l3out."+msoSchemaTemplateL3outName]
						if !ok {
							return fmt.Errorf("L3out resource not found in state")
						}
						msoSchemaTemplateL3outSchemaId = rs.Primary.Attributes["schema_id"]
						return nil
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add L3out description") },
				Config:    testAccMSOSchemaTemplateL3outConfigAddDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_l3out."+msoSchemaTemplateL3outName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_l3out."+msoSchemaTemplateL3outName, "l3out_name", msoSchemaTemplateL3outName),
					resource.TestCheckResourceAttr("mso_schema_template_l3out."+msoSchemaTemplateL3outName, "display_name", msoSchemaTemplateL3outName),
					resource.TestCheckResourceAttr("mso_schema_template_l3out."+msoSchemaTemplateL3outName, "description", "Terraform test L3out"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3out description") },
				Config:    testAccMSOSchemaTemplateL3outConfigUpdateDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_l3out."+msoSchemaTemplateL3outName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_l3out."+msoSchemaTemplateL3outName, "l3out_name", msoSchemaTemplateL3outName),
					resource.TestCheckResourceAttr("mso_schema_template_l3out."+msoSchemaTemplateL3outName, "description", "Terraform test L3out updated"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove L3out description") },
				Config:    testAccMSOSchemaTemplateL3outConfigRemoveDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_l3out."+msoSchemaTemplateL3outName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_l3out."+msoSchemaTemplateL3outName, "l3out_name", msoSchemaTemplateL3outName),
					resource.TestCheckResourceAttr("mso_schema_template_l3out."+msoSchemaTemplateL3outName, "description", ""),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3out VRF") },
				Config:    testAccMSOSchemaTemplateL3outConfigUpdateVrf(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_l3out."+msoSchemaTemplateL3outName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_l3out."+msoSchemaTemplateL3outName, "l3out_name", msoSchemaTemplateL3outName),
					resource.TestCheckResourceAttr("mso_schema_template_l3out."+msoSchemaTemplateL3outName, "vrf_name", msoSchemaTemplateVrfL3MulticastName),
					resource.TestCheckResourceAttrSet("mso_schema_template_l3out."+msoSchemaTemplateL3outName, "vrf_schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_l3out."+msoSchemaTemplateL3outName, "vrf_template_name", msoSchemaTemplateName),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import L3out") },
				ResourceName: "mso_schema_template_l3out." + msoSchemaTemplateL3outName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["mso_schema_template_l3out."+msoSchemaTemplateL3outName]
					if !ok {
						return "", fmt.Errorf("L3out resource not found in state")
					}
					return fmt.Sprintf("%s/templates/%s/intersiteL3outs/%s", rs.Primary.Attributes["schema_id"], rs.Primary.Attributes["template_name"], rs.Primary.Attributes["l3out_name"]), nil
				},
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate L3out after manual deletion from NDO")
					msoClient := testAccProvider.Meta().(*client.Client)
					l3outRemovePatchPayload := models.GetRemovePatchPayload(fmt.Sprintf("/templates/%s/intersiteL3outs/%s", msoSchemaTemplateName, msoSchemaTemplateL3outName))
					_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", msoSchemaTemplateL3outSchemaId), l3outRemovePatchPayload)
					if err != nil {
						t.Fatalf("Failed to manually delete L3out: %v", err)
					}
				},
				Config: testAccMSOSchemaTemplateL3outConfigUpdateVrf(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_l3out."+msoSchemaTemplateL3outName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_l3out."+msoSchemaTemplateL3outName, "l3out_name", msoSchemaTemplateL3outName),
					resource.TestCheckResourceAttr("mso_schema_template_l3out."+msoSchemaTemplateL3outName, "display_name", msoSchemaTemplateL3outName),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateL3outPrerequisiteConfig() string {
	return fmt.Sprintf(`%s%s%s%s`, testSiteConfigAnsibleTest(), testTenantConfig(), testSchemaConfig(), testSchemaTemplateVrfConfig())
}

func testAccMSOSchemaTemplateL3outPrerequisiteWithL3MulticastVrfConfig() string {
	return fmt.Sprintf(`%s%s`, testAccMSOSchemaTemplateL3outPrerequisiteConfig(), testSchemaTemplateVrfL3MulticastConfig())
}

func testAccMSOSchemaTemplateL3outConfigCreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_l3out" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		template_name = "%[4]s"
		l3out_name    = "%[2]s"
		display_name  = "%[2]s"
		vrf_name      = mso_schema_template_vrf.%[5]s.name
	}`, testAccMSOSchemaTemplateL3outPrerequisiteConfig(), msoSchemaTemplateL3outName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfName)
}

func testAccMSOSchemaTemplateL3outConfigAddDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_l3out" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		template_name = "%[4]s"
		l3out_name    = "%[2]s"
		display_name  = "%[2]s"
		vrf_name      = mso_schema_template_vrf.%[5]s.name
		description   = "Terraform test L3out"
	}`, testAccMSOSchemaTemplateL3outPrerequisiteConfig(), msoSchemaTemplateL3outName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfName)
}

func testAccMSOSchemaTemplateL3outConfigUpdateDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_l3out" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		template_name = "%[4]s"
		l3out_name    = "%[2]s"
		display_name  = "%[2]s"
		vrf_name      = mso_schema_template_vrf.%[5]s.name
		description   = "Terraform test L3out updated"
	}`, testAccMSOSchemaTemplateL3outPrerequisiteConfig(), msoSchemaTemplateL3outName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfName)
}

func testAccMSOSchemaTemplateL3outConfigRemoveDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_l3out" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		template_name = "%[4]s"
		l3out_name    = "%[2]s"
		display_name  = "%[2]s"
		vrf_name      = mso_schema_template_vrf.%[5]s.name
		description   = ""
	}`, testAccMSOSchemaTemplateL3outPrerequisiteConfig(), msoSchemaTemplateL3outName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfName)
}

func testAccMSOSchemaTemplateL3outConfigUpdateVrf() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_l3out" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		template_name = "%[4]s"
		l3out_name    = "%[2]s"
		display_name  = "%[2]s"
		vrf_name      = mso_schema_template_vrf.%[5]s.name
	}`, testAccMSOSchemaTemplateL3outPrerequisiteWithL3MulticastVrfConfig(), msoSchemaTemplateL3outName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfL3MulticastName)
}

func testAccCheckMSOSchemaTemplateL3outDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "mso_schema_template_l3out" {
			schemaID := rs.Primary.Attributes["schema_id"]
			cont, err := client.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
			if err != nil {
				return nil
			}
			count, err := cont.ArrayCount("templates")
			if err != nil {
				return fmt.Errorf("No Template found")
			}
			for i := 0; i < count; i++ {
				tempCont, err := cont.ArrayElement(i, "templates")
				if err != nil {
					return fmt.Errorf("No template exists")
				}
				l3outCount, err := tempCont.ArrayCount("intersiteL3outs")
				if err != nil {
					return fmt.Errorf("Unable to get L3out list")
				}
				for j := 0; j < l3outCount; j++ {
					l3outCont, err := tempCont.ArrayElement(j, "intersiteL3outs")
					if err != nil {
						return err
					}
					name := models.StripQuotes(l3outCont.S("name").String())
					if rs.Primary.Attributes["l3out_name"] == name {
						return fmt.Errorf("Schema Template L3out record still exists")
					}
				}
			}
		}
	}
	return nil
}
