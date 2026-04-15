package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOSchemaTemplateExternalEpgDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateExtEpgDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Read External EPG datasource not found error") },
				Config:      testAccMSOSchemaTemplateExtEpgDatasourceNotFound(),
				ExpectError: regexp.MustCompile("Unable to find the External Epg"),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read External EPG datasource") },
				Config:    testAccMSOSchemaTemplateExtEpgDatasource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.mso_schema_template_external_epg.ext_epg", "schema_id"),
					resource.TestCheckResourceAttr("data.mso_schema_template_external_epg.ext_epg", "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("data.mso_schema_template_external_epg.ext_epg", "external_epg_name", msoSchemaTemplateExtEpgName),
					resource.TestCheckResourceAttr("data.mso_schema_template_external_epg.ext_epg", "display_name", msoSchemaTemplateExtEpgName),
					resource.TestCheckResourceAttr("data.mso_schema_template_external_epg.ext_epg", "vrf_name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttrSet("data.mso_schema_template_external_epg.ext_epg", "vrf_schema_id"),
					resource.TestCheckResourceAttr("data.mso_schema_template_external_epg.ext_epg", "vrf_template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("data.mso_schema_template_external_epg.ext_epg", "external_epg_type", "on-premise"),
					resource.TestCheckResourceAttr("data.mso_schema_template_external_epg.ext_epg", "description", ""),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateExtEpgDatasource() string {
	return fmt.Sprintf(`%s
	data "mso_schema_template_external_epg" "ext_epg" {
		schema_id         = mso_schema.%[2]s.id
		template_name     = "%[3]s"
		external_epg_name = mso_schema_template_external_epg.%[4]s.external_epg_name
	}`, testAccMSOSchemaTemplateExtEpgConfigCreate(), msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateExtEpgName)
}

func testAccMSOSchemaTemplateExtEpgDatasourceNotFound() string {
	return fmt.Sprintf(`%s
	data "mso_schema_template_external_epg" "ext_epg" {
		schema_id         = mso_schema.%[2]s.id
		template_name     = "%[3]s"
		external_epg_name = "non_existing_ext_epg_name"
	}`, testAccMSOSchemaTemplateExtEpgConfigCreate(), msoSchemaName, msoSchemaTemplateName)
}
