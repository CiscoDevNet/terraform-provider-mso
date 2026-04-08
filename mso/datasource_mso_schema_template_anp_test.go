package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOSchemaTemplateAnpDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateAnpDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Read ANP datasource not found error") },
				Config:      testAccMSOSchemaTemplateAnpDatasourceNotFound(),
				ExpectError: regexp.MustCompile("Unable to find the ANP"),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read ANP datasource") },
				Config:    testAccMSOSchemaTemplateAnpDatasource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.mso_schema_template_anp.anp", "schema_id"),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp.anp", "template", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp.anp", "name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp.anp", "display_name", msoSchemaTemplateAnpName),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateAnpDatasource() string {
	return fmt.Sprintf(`%s
	data "mso_schema_template_anp" "anp" {
		schema_id = mso_schema.%[2]s.id
		template  = "%[3]s"
		name      = mso_schema_template_anp.%[4]s.name
	}`, testAccMSOSchemaTemplateAnpConfigCreate(), msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName)
}

func testAccMSOSchemaTemplateAnpDatasourceNotFound() string {
	return fmt.Sprintf(`%s
	data "mso_schema_template_anp" "anp" {
		schema_id = mso_schema.%[2]s.id
		template  = "%[3]s"
		name      = "non_existing_anp_name"
	}`, testAccMSOSchemaTemplateAnpConfigCreate(), msoSchemaName, msoSchemaTemplateName)
}
