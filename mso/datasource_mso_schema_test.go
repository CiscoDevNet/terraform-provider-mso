package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMSOSchemaDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Read Schema datasource not found error") },
				Config:      testAccMSOSchemaDatasourceNotFound(),
				ExpectError: regexp.MustCompile("Schema of specified name not found"),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read Schema datasource") },
				Config:    testAccMSOSchemaDatasource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_schema.schema", "name", msoSchemaName),
					resource.TestCheckResourceAttr("data.mso_schema.schema", "description", ""),
					resource.TestCheckResourceAttr("data.mso_schema.schema", "template.#", "1"),
					CustomTestCheckTypeSetElemAttrs("data.mso_schema.schema", "template", map[string]string{
						"name":          msoSchemaTemplateName,
						"display_name":  msoSchemaTemplateName,
						"template_type": "aci_multi_site",
						"description":   "",
					}),
				),
			},
		},
	})
}

func testAccMSOSchemaDatasource() string {
	return fmt.Sprintf(`%s
	data "mso_schema" "schema" {
		name = mso_schema.%s.name
	}`, testAccMSOSchemaConfigCreate(), msoSchemaName)
}

func testAccMSOSchemaDatasourceNotFound() string {
	return `
	data "mso_schema" "schema" {
		name = "non_existing_schema_name"
	}`
}
