package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// msoSchemaTemplateAnpSchemaId is set during the first test step's Check to capture the dynamic schema ID for use in the manual deletion PreConfig step.
var msoSchemaTemplateAnpSchemaId string

func TestAccMSOSchemaTemplateAnpResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateAnpDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create ANP with name and display_name") },
				Config:    testAccMSOSchemaTemplateAnpConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_anp."+msoSchemaTemplateAnpName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "template", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "display_name", msoSchemaTemplateAnpName),
					// Capture the dynamic schema ID from state for use in the manual deletion PreConfig step
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["mso_schema_template_anp."+msoSchemaTemplateAnpName]
						if !ok {
							return fmt.Errorf("ANP resource not found in state")
						}
						msoSchemaTemplateAnpSchemaId = rs.Primary.Attributes["schema_id"]
						return nil
					},
					// The resource description does not exist because the schema is initialized without a value and set to null
					resource.TestCheckNoResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "description"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update ANP display_name") },
				Config:    testAccMSOSchemaTemplateAnpConfigUpdateDisplayName(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_anp."+msoSchemaTemplateAnpName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "template", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "display_name", msoSchemaTemplateAnpName+" updated"),
					// The resource description does not exist because the schema is initialized without a value and set to null
					resource.TestCheckNoResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "description"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add ANP description") },
				Config:    testAccMSOSchemaTemplateAnpConfigAddDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_anp."+msoSchemaTemplateAnpName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "template", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "display_name", msoSchemaTemplateAnpName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "description", "Terraform test ANP"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove ANP description") },
				Config:    testAccMSOSchemaTemplateAnpConfigRemoveDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_anp."+msoSchemaTemplateAnpName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "template", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "display_name", msoSchemaTemplateAnpName+" updated"),
					// The resource description exist because the schema when set before will unset to empty string
					resource.TestCheckResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "description", ""),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import ANP") },
				ResourceName: "mso_schema_template_anp." + msoSchemaTemplateAnpName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["mso_schema_template_anp."+msoSchemaTemplateAnpName]
					if !ok {
						return "", fmt.Errorf("ANP resource not found in state")
					}
					return fmt.Sprintf("%s/templates/%s/anps/%s", rs.Primary.Attributes["schema_id"], rs.Primary.Attributes["template"], rs.Primary.Attributes["name"]), nil
				},
				ImportStateVerify: true,
				// Description attribute is set to empty string on import but it is not provided in the config.
				ImportStateVerifyIgnore: []string{"description"},
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate ANP after manual deletion from NDO")
					msoClient := testAccProvider.Meta().(*client.Client)
					anpRemovePatchPayload := models.GetRemovePatchPayload(fmt.Sprintf("/templates/%s/anps/%s", msoSchemaTemplateName, msoSchemaTemplateAnpName))
					_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", msoSchemaTemplateAnpSchemaId), anpRemovePatchPayload)
					if err != nil {
						t.Fatalf("Failed to manually delete ANP: %v", err)
					}
				},
				Config: testAccMSOSchemaTemplateAnpConfigRemoveDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_anp."+msoSchemaTemplateAnpName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "display_name", msoSchemaTemplateAnpName+" updated"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add EPG child to ANP") },
				Config:    testAccMSOSchemaTemplateAnpConfigWithEpgChild(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_anp."+msoSchemaTemplateAnpName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "display_name", msoSchemaTemplateAnpName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "display_name", msoSchemaTemplateAnpEpgName),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update ANP description with EPG child present") },
				Config:    testAccMSOSchemaTemplateAnpConfigWithEpgChildUpdateDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_anp."+msoSchemaTemplateAnpName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "display_name", msoSchemaTemplateAnpName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "description", "Terraform test ANP with EPG"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "display_name", msoSchemaTemplateAnpEpgName),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateAnpPrerequisiteConfig() string {
	return fmt.Sprintf(`%s%s%s`, testSiteConfigAnsibleTest(), testTenantConfig(), testSchemaConfig())
}

func testAccMSOSchemaTemplateAnpConfigCreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp" "%[2]s" {
		schema_id    = mso_schema.%[3]s.id
		template     = "%[4]s"
		name         = "%[2]s"
		display_name = "%[2]s"
	}`, testAccMSOSchemaTemplateAnpPrerequisiteConfig(), msoSchemaTemplateAnpName, msoSchemaName, msoSchemaTemplateName)
}

func testAccMSOSchemaTemplateAnpConfigUpdateDisplayName() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp" "%[2]s" {
		schema_id    = mso_schema.%[3]s.id
		template     = "%[4]s"
		name         = "%[2]s"
		display_name = "%[2]s updated"
	}`, testAccMSOSchemaTemplateAnpPrerequisiteConfig(), msoSchemaTemplateAnpName, msoSchemaName, msoSchemaTemplateName)
}

func testAccMSOSchemaTemplateAnpConfigAddDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp" "%[2]s" {
		schema_id    = mso_schema.%[3]s.id
		template     = "%[4]s"
		name         = "%[2]s"
		display_name = "%[2]s updated"
		description  = "Terraform test ANP"
	}`, testAccMSOSchemaTemplateAnpPrerequisiteConfig(), msoSchemaTemplateAnpName, msoSchemaName, msoSchemaTemplateName)
}

func testAccMSOSchemaTemplateAnpConfigRemoveDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp" "%[2]s" {
		schema_id    = mso_schema.%[3]s.id
		template     = "%[4]s"
		name         = "%[2]s"
		display_name = "%[2]s updated"
		description  = ""
	}`, testAccMSOSchemaTemplateAnpPrerequisiteConfig(), msoSchemaTemplateAnpName, msoSchemaName, msoSchemaTemplateName)
}

func testAccMSOSchemaTemplateAnpConfigWithEpgChild() string {
	return testAccMSOSchemaTemplateAnpConfigUpdateDisplayName() + testSchemaTemplateAnpEpgConfig()
}

func testAccMSOSchemaTemplateAnpConfigWithEpgChildUpdateDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp" "%[2]s" {
		schema_id    = mso_schema.%[3]s.id
		template     = "%[4]s"
		name         = "%[2]s"
		display_name = "%[2]s updated"
		description  = "Terraform test ANP with EPG"
	}`, testAccMSOSchemaTemplateAnpPrerequisiteConfig(), msoSchemaTemplateAnpName, msoSchemaName, msoSchemaTemplateName) + testSchemaTemplateAnpEpgConfig()
}

func testAccCheckMSOSchemaTemplateAnpDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "mso_schema_template_anp" {
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
					name := models.StripQuotes(anpCont.S("name").String())
					if rs.Primary.ID == name {
						return fmt.Errorf("Schema Template Anp record still exists")
					}
				}
			}
		}
	}
	return nil
}
