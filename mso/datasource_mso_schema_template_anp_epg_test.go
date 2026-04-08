package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOSchemaTemplateAnpEpgDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateAnpEpgDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Read EPG datasource not found error") },
				Config:      testAccMSOSchemaTemplateAnpEpgDatasourceNotFound(),
				ExpectError: regexp.MustCompile("Unable to find the ANP EPG"),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read EPG datasource") },
				Config:    testAccMSOSchemaTemplateAnpEpgDatasource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.mso_schema_template_anp_epg.epg", "schema_id"),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg.epg", "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg.epg", "anp_name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg.epg", "name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg.epg", "display_name", msoSchemaTemplateAnpEpgName),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateAnpEpgDatasource() string {
	return fmt.Sprintf(`%s
	data "mso_schema_template_anp_epg" "epg" {
		schema_id     = mso_schema.%[2]s.id
		template_name = "%[3]s"
		anp_name      = "%[4]s"
		name          = mso_schema_template_anp_epg.%[5]s.name
	}`, testAccMSOSchemaTemplateAnpEpgConfigCreate(), msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName)
}

func testAccMSOSchemaTemplateAnpEpgDatasourceNotFound() string {
	return fmt.Sprintf(`%s
	data "mso_schema_template_anp_epg" "epg" {
		schema_id     = mso_schema.%[2]s.id
		template_name = "%[3]s"
		anp_name      = "%[4]s"
		name          = "non_existing_epg_name"
	}`, testAccMSOSchemaTemplateAnpEpgConfigCreate(), msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName)
}
