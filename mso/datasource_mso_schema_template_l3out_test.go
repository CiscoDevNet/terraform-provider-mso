package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOSchemaTemplateL3outDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateL3outDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Read L3out datasource not found error") },
				Config:      testAccMSOSchemaTemplateL3outDatasourceNotFound(),
				ExpectError: regexp.MustCompile("Unable to find the L3out"),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read L3out datasource") },
				Config:    testAccMSOSchemaTemplateL3outDatasource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.mso_schema_template_l3out.l3out", "schema_id"),
					resource.TestCheckResourceAttr("data.mso_schema_template_l3out.l3out", "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("data.mso_schema_template_l3out.l3out", "l3out_name", msoSchemaTemplateL3outName),
					resource.TestCheckResourceAttr("data.mso_schema_template_l3out.l3out", "display_name", msoSchemaTemplateL3outName),
					resource.TestCheckResourceAttr("data.mso_schema_template_l3out.l3out", "vrf_name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttrSet("data.mso_schema_template_l3out.l3out", "vrf_schema_id"),
					resource.TestCheckResourceAttr("data.mso_schema_template_l3out.l3out", "vrf_template_name", msoSchemaTemplateName),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateL3outDatasource() string {
	return fmt.Sprintf(`%s
	data "mso_schema_template_l3out" "l3out" {
		schema_id     = mso_schema.%[2]s.id
		template_name = "%[3]s"
		l3out_name    = mso_schema_template_l3out.%[4]s.l3out_name
	}`, testAccMSOSchemaTemplateL3outConfigCreate(), msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateL3outName)
}

func testAccMSOSchemaTemplateL3outDatasourceNotFound() string {
	return fmt.Sprintf(`%s
	data "mso_schema_template_l3out" "l3out" {
		schema_id     = mso_schema.%[2]s.id
		template_name = "%[3]s"
		l3out_name    = "non_existing_l3out_name"
	}`, testAccMSOSchemaTemplateL3outConfigCreate(), msoSchemaName, msoSchemaTemplateName)
}
