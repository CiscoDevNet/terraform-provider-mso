package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMSOSchemaTemplateAnpEpgUsegAttrDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateAnpEpgUsegAttrDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Read EPG UsegAttr datasource not found error") },
				Config:      testAccMSOSchemaTemplateAnpEpgUsegAttrDatasourceNotFound(),
				ExpectError: regexp.MustCompile("Unable to find the ANP EPG uSeg Attribute"),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read EPG UsegAttr datasource") },
				Config:    testAccMSOSchemaTemplateAnpEpgUsegAttrDatasource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.mso_schema_template_anp_epg_useg_attr.useg_attr", "schema_id"),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg_useg_attr.useg_attr", "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg_useg_attr.useg_attr", "anp_name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg_useg_attr.useg_attr", "epg_name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg_useg_attr.useg_attr", "name", msoSchemaTemplateAnpEpgUsegAttrName),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg_useg_attr.useg_attr", "useg_type", "ip"),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg_useg_attr.useg_attr", "value", msoSchemaTemplateAnpEpgUsegAttrIp),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg_useg_attr.useg_attr", "operator", ""),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg_useg_attr.useg_attr", "category", "test_category"),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg_useg_attr.useg_attr", "description", "test useg"),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg_useg_attr.useg_attr", "useg_subnet", "true"),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateAnpEpgUsegAttrDatasource() string {
	return fmt.Sprintf(`%s
	data "mso_schema_template_anp_epg_useg_attr" "useg_attr" {
		schema_id     = mso_schema.%[2]s.id
		template_name = "%[3]s"
		anp_name      = "%[4]s"
		epg_name      = "%[5]s"
		name          = mso_schema_template_anp_epg_useg_attr.%[6]s.name
	}`, testAccMSOSchemaTemplateAnpEpgUsegAttrConfigUpdateAll(), msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateAnpEpgUsegAttrName)
}

func testAccMSOSchemaTemplateAnpEpgUsegAttrDatasourceNotFound() string {
	return fmt.Sprintf(`%s
	data "mso_schema_template_anp_epg_useg_attr" "useg_attr" {
		schema_id     = mso_schema.%[2]s.id
		template_name = "%[3]s"
		anp_name      = "%[4]s"
		epg_name      = "%[5]s"
		name          = "non_existent_useg_attr"
	}`, testAccMSOSchemaTemplateAnpEpgUsegAttrConfigUpdateAll(), msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName)
}
