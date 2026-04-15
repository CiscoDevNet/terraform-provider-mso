package mso

// Note: The mso_schema resource uses lifecycle { ignore_changes = [template] } to prevent drift.
// This is necessary because mso_schema.Read() reads all templates from the API and writes them to state,
// which conflicts with mso_schema_template managing the same templates. The ignore_changes directive
// tells Terraform to ignore externally added templates, avoiding perpetual plan diffs.

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// msoSchemaTemplateSchemaId captures the schema ID for the manual delete/recreate step.
var msoSchemaTemplateSchemaId string

func TestAccMSOSchemaTemplateResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateTestDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create template with name and display_name") },
				Config:    testAccMSOSchemaTemplateConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template."+msoSchemaTemplateName, "schema_id"),
					resource.TestCheckResourceAttrSet("mso_schema_template."+msoSchemaTemplateName, "tenant_id"),
					resource.TestCheckResourceAttr("mso_schema_template."+msoSchemaTemplateName, "name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template."+msoSchemaTemplateName, "display_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template."+msoSchemaTemplateName, "template_type", "aci_multi_site"),
					resource.TestCheckResourceAttr("mso_schema_template."+msoSchemaTemplateName, "description", ""),
					// Capture the schema ID for the manual delete/recreate step.
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["mso_schema_template."+msoSchemaTemplateName]
						if !ok {
							return fmt.Errorf("mso_schema_template.%s not found in state", msoSchemaTemplateName)
						}
						msoSchemaTemplateSchemaId = rs.Primary.Attributes["schema_id"]
						return nil
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update display_name") },
				Config:    testAccMSOSchemaTemplateConfigUpdateDisplayName(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template."+msoSchemaTemplateName, "name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template."+msoSchemaTemplateName, "display_name", msoSchemaTemplateName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template."+msoSchemaTemplateName, "description", ""),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add description") },
				Config:    testAccMSOSchemaTemplateConfigAddDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template."+msoSchemaTemplateName, "name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template."+msoSchemaTemplateName, "display_name", msoSchemaTemplateName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template."+msoSchemaTemplateName, "description", "Terraform test template"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove description") },
				Config:    testAccMSOSchemaTemplateConfigRemoveDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template."+msoSchemaTemplateName, "name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template."+msoSchemaTemplateName, "display_name", msoSchemaTemplateName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template."+msoSchemaTemplateName, "description", ""),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import template") },
				ResourceName:      "mso_schema_template." + msoSchemaTemplateName,
				ImportState:       true,
				ImportStateIdFunc: testAccMSOSchemaTemplateImportStateIdFunc,
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Verify recreation of manually deleted template")
					msoClient := testAccProvider.Meta().(*client.Client)
					path := fmt.Sprintf("/templates/%s", msoSchemaTemplateName)
					_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", msoSchemaTemplateSchemaId), models.GetRemovePatchPayload(path))
					if err != nil {
						panic(fmt.Sprintf("Failed to delete template %s via API: %s", msoSchemaTemplateName, err))
					}
				},
				Config: testAccMSOSchemaTemplateConfigRemoveDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template."+msoSchemaTemplateName, "name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template."+msoSchemaTemplateName, "display_name", msoSchemaTemplateName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template."+msoSchemaTemplateName, "description", ""),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add children (VRF, BD, ANP, EPG)") },
				Config:    testAccMSOSchemaTemplateConfigWithChildren(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template."+msoSchemaTemplateName, "name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template."+msoSchemaTemplateName, "display_name", msoSchemaTemplateName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "name", msoSchemaTemplateAnpEpgName),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update description with children present") },
				Config:    testAccMSOSchemaTemplateConfigWithChildrenUpdateDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template."+msoSchemaTemplateName, "name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template."+msoSchemaTemplateName, "display_name", msoSchemaTemplateName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template."+msoSchemaTemplateName, "description", "Terraform test template updated"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf."+msoSchemaTemplateVrfName, "name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr("mso_schema_template_bd."+msoSchemaTemplateBdName, "name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("mso_schema_template_anp."+msoSchemaTemplateAnpName, "name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "name", msoSchemaTemplateAnpEpgName),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove children") },
				Config:    testAccMSOSchemaTemplateConfigRemoveChildren(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template."+msoSchemaTemplateName, "name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template."+msoSchemaTemplateName, "display_name", msoSchemaTemplateName+" updated"),
					resource.TestCheckResourceAttr("mso_schema_template."+msoSchemaTemplateName, "description", "Terraform test template updated"),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateImportStateIdFunc(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources["mso_schema_template."+msoSchemaTemplateName]
	if !ok {
		return "", fmt.Errorf("mso_schema_template.%s not found in state", msoSchemaTemplateName)
	}
	return fmt.Sprintf("%s/templates/%s", rs.Primary.Attributes["schema_id"], msoSchemaTemplateName), nil
}

// testAccMSOSchemaTemplatePrerequisiteConfig returns the Terraform config for prerequisites (site + tenant + schema).
// The schema uses lifecycle { ignore_changes = [template] } to prevent drift from templates
// added by mso_schema_template. See testSchemaConfigIgnoreTemplates() for details.
func testAccMSOSchemaTemplatePrerequisiteConfig() string {
	return fmt.Sprintf(`%s%s%s`, testSiteConfigAnsibleTest(), testTenantConfig(), testSchemaConfigIgnoreTemplates())
}

func testAccMSOSchemaTemplateConfigCreate() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_template" "%[2]s" {
	schema_id    = mso_schema.%[3]s.id
	name         = "%[2]s"
	display_name = "%[2]s"
	tenant_id    = mso_tenant.%[4]s.id
}
`, testAccMSOSchemaTemplatePrerequisiteConfig(), msoSchemaTemplateName, msoSchemaName, msoTenantName)
}

func testAccMSOSchemaTemplateConfigUpdateDisplayName() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_template" "%[2]s" {
	schema_id    = mso_schema.%[3]s.id
	name         = "%[2]s"
	display_name = "%[2]s updated"
	tenant_id    = mso_tenant.%[4]s.id
}
`, testAccMSOSchemaTemplatePrerequisiteConfig(), msoSchemaTemplateName, msoSchemaName, msoTenantName)
}

func testAccMSOSchemaTemplateConfigAddDescription() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_template" "%[2]s" {
	schema_id    = mso_schema.%[3]s.id
	name         = "%[2]s"
	display_name = "%[2]s updated"
	description  = "Terraform test template"
	tenant_id    = mso_tenant.%[4]s.id
}
`, testAccMSOSchemaTemplatePrerequisiteConfig(), msoSchemaTemplateName, msoSchemaName, msoTenantName)
}

func testAccMSOSchemaTemplateConfigRemoveDescription() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_template" "%[2]s" {
	schema_id    = mso_schema.%[3]s.id
	name         = "%[2]s"
	display_name = "%[2]s updated"
	tenant_id    = mso_tenant.%[4]s.id
}
`, testAccMSOSchemaTemplatePrerequisiteConfig(), msoSchemaTemplateName, msoSchemaName, msoTenantName)
}

func testAccMSOSchemaTemplateConfigWithChildren() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_template" "%[2]s" {
	schema_id    = mso_schema.%[3]s.id
	name         = "%[2]s"
	display_name = "%[2]s updated"
	tenant_id    = mso_tenant.%[4]s.id
}

resource "mso_schema_template_vrf" "%[5]s" {
	name         = "%[5]s"
	display_name = "%[5]s"
	schema_id    = mso_schema.%[3]s.id
	template     = mso_schema_template.%[2]s.name
}

resource "mso_schema_template_bd" "%[6]s" {
	schema_id              = mso_schema.%[3]s.id
	template_name          = mso_schema_template.%[2]s.name
	name                   = "%[6]s"
	display_name           = "%[6]s"
	layer2_unknown_unicast = "proxy"
	vrf_name               = mso_schema_template_vrf.%[5]s.name
}

resource "mso_schema_template_anp" "%[7]s" {
	name         = "%[7]s"
	display_name = "%[7]s"
	schema_id    = mso_schema.%[3]s.id
	template     = mso_schema_template.%[2]s.name
	depends_on   = [mso_schema_template_bd.%[6]s]
}

resource "mso_schema_template_anp_epg" "%[8]s" {
	name          = "%[8]s"
	display_name  = "%[8]s"
	anp_name      = "%[7]s"
	schema_id     = mso_schema.%[3]s.id
	template_name = mso_schema_template.%[2]s.name
	depends_on    = [mso_schema_template_anp.%[7]s]
}
`, testAccMSOSchemaTemplatePrerequisiteConfig(), msoSchemaTemplateName, msoSchemaName, msoTenantName,
		msoSchemaTemplateVrfName, msoSchemaTemplateBdName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName)
}

func testAccMSOSchemaTemplateConfigWithChildrenUpdateDescription() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_template" "%[2]s" {
	schema_id    = mso_schema.%[3]s.id
	name         = "%[2]s"
	display_name = "%[2]s updated"
	description  = "Terraform test template updated"
	tenant_id    = mso_tenant.%[4]s.id
}

resource "mso_schema_template_vrf" "%[5]s" {
	name         = "%[5]s"
	display_name = "%[5]s"
	schema_id    = mso_schema.%[3]s.id
	template     = mso_schema_template.%[2]s.name
}

resource "mso_schema_template_bd" "%[6]s" {
	schema_id              = mso_schema.%[3]s.id
	template_name          = mso_schema_template.%[2]s.name
	name                   = "%[6]s"
	display_name           = "%[6]s"
	layer2_unknown_unicast = "proxy"
	vrf_name               = mso_schema_template_vrf.%[5]s.name
}

resource "mso_schema_template_anp" "%[7]s" {
	name         = "%[7]s"
	display_name = "%[7]s"
	schema_id    = mso_schema.%[3]s.id
	template     = mso_schema_template.%[2]s.name
	depends_on   = [mso_schema_template_bd.%[6]s]
}

resource "mso_schema_template_anp_epg" "%[8]s" {
	name          = "%[8]s"
	display_name  = "%[8]s"
	anp_name      = "%[7]s"
	schema_id     = mso_schema.%[3]s.id
	template_name = mso_schema_template.%[2]s.name
	depends_on    = [mso_schema_template_anp.%[7]s]
}
`, testAccMSOSchemaTemplatePrerequisiteConfig(), msoSchemaTemplateName, msoSchemaName, msoTenantName,
		msoSchemaTemplateVrfName, msoSchemaTemplateBdName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName)
}

// testAccMSOSchemaTemplateConfigRemoveChildren returns the same template config as testAccMSOSchemaTemplateConfigWithChildrenUpdateDescription
// but without children. Keeping the template unchanged avoids a concurrent PATCH to NDO during child deletion.
func testAccMSOSchemaTemplateConfigRemoveChildren() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_template" "%[2]s" {
	schema_id    = mso_schema.%[3]s.id
	name         = "%[2]s"
	display_name = "%[2]s updated"
	description  = "Terraform test template updated"
	tenant_id    = mso_tenant.%[4]s.id
}
`, testAccMSOSchemaTemplatePrerequisiteConfig(), msoSchemaTemplateName, msoSchemaName, msoTenantName)
}

// testAccCheckMSOSchemaTemplateTestDestroy verifies that the template is deleted from the schema.
func testAccCheckMSOSchemaTemplateTestDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "mso_schema_template" {
			schemaId := rs.Primary.Attributes["schema_id"]
			cont, err := client.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
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
					return fmt.Errorf("No Template exists")
				}
				apiTemplateName := models.StripQuotes(tempCont.S("name").String())
				if rs.Primary.ID == apiTemplateName {
					return fmt.Errorf("Schema template record still exists")
				}
			}
		}
	}
	return nil
}
