package mso

// Note: The mso_schema resource uses lifecycle { ignore_changes = [template] } to prevent drift.
// See resource_mso_schema_template_test.go for details.

// Data source tests require a two-step approach (create resource, then add data source) because
// terraform-plugin-sdk v1 reads data sources during the refresh walk before resources are applied.
// - Without depends_on: config values are wholly known so ReadDataSource is called immediately,
//   but the resource doesn't exist in the API yet, causing a "not found" error.
// - With depends_on: ForcePlanRead defers the read, but RemovePlannedResourceInstanceObjects
//   wipes the planned data source state after refresh, causing a non-empty plan error.
// Splitting into steps ensures the resource exists in the API before the data source is read.

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOSchemaTemplateDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Read non-existing template - not found error") },
				Config:      testAccMSOSchemaTemplateDatasourceNotFound(),
				ExpectError: regexp.MustCompile("Template of specified name not found"),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create template resource") },
				Config:    testAccMSOSchemaTemplateDatasourceCreateResource(),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read existing template via datasource") },
				Config:    testAccMSOSchemaTemplateDatasourceRead(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.mso_schema_template.template", "schema_id"),
					resource.TestCheckResourceAttr("data.mso_schema_template.template", "name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("data.mso_schema_template.template", "display_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttrSet("data.mso_schema_template.template", "tenant_id"),
					resource.TestCheckResourceAttr("data.mso_schema_template.template", "template_type", "aci_multi_site"),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateDatasourceNotFound() string {
	return fmt.Sprintf(`%[1]s
data "mso_schema_template" "template" {
	schema_id = mso_schema.%[2]s.id
	name      = "non_existing_template"
}
`, testAccMSOSchemaTemplatePrerequisiteConfig(), msoSchemaName)
}

func testAccMSOSchemaTemplateDatasourceCreateResource() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_template" "%[2]s" {
	schema_id    = mso_schema.%[3]s.id
	name         = "%[2]s"
	display_name = "%[2]s"
	tenant_id    = mso_tenant.%[4]s.id
}
`, testAccMSOSchemaTemplatePrerequisiteConfig(), msoSchemaTemplateName, msoSchemaName, msoTenantName)
}

func testAccMSOSchemaTemplateDatasourceRead() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_template" "%[2]s" {
	schema_id    = mso_schema.%[3]s.id
	name         = "%[2]s"
	display_name = "%[2]s"
	tenant_id    = mso_tenant.%[4]s.id
}

data "mso_schema_template" "template" {
	schema_id  = mso_schema.%[3]s.id
	name       = "%[2]s"
}
`, testAccMSOSchemaTemplatePrerequisiteConfig(), msoSchemaTemplateName, msoSchemaName, msoTenantName)
}
