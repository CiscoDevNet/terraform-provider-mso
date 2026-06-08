package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMSOSchemaTemplateExternalEpgSubnetDatasource(t *testing.T) {
	dataSourceName := "data.mso_schema_template_external_epg_subnet.ext_epg_subnet"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateExtEpgSubnetDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Read External EPG Subnet datasource not found error") },
				Config:      testAccMSOSchemaTemplateExtEpgSubnetDatasourceNotFound(),
				ExpectError: regexp.MustCompile("Unable to find the External Epg"),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read External EPG Subnet datasource") },
				Config:    testAccMSOSchemaTemplateExtEpgSubnetDatasource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "schema_id"),
					resource.TestCheckResourceAttr(dataSourceName, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr(dataSourceName, "external_epg_name", msoSchemaTemplateExtEpgName),
					resource.TestCheckResourceAttr(dataSourceName, "ip", msoSchemaTemplateExtEpgSubnetIp),
					resource.TestCheckResourceAttr(dataSourceName, "name", msoSchemaTemplateExtEpgSubnetName),
					resource.TestCheckResourceAttr(dataSourceName, "scope.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "scope.0", "shared-rtctrl"),
					resource.TestCheckResourceAttr(dataSourceName, "scope.1", "export-rtctrl"),
					resource.TestCheckResourceAttr(dataSourceName, "aggregate.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "aggregate.0", "shared-rtctrl"),
					resource.TestCheckResourceAttr(dataSourceName, "aggregate.1", "export-rtctrl"),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateExtEpgSubnetDatasourceConfigCreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_external_epg_subnet" "%[2]s_subnet" {
		schema_id         = mso_schema.%[3]s.id
		template_name     = "%[4]s"
		external_epg_name = mso_schema_template_external_epg.%[2]s.external_epg_name
		ip                = "%[5]s"
		name              = "%[6]s"
		scope             = ["shared-rtctrl", "export-rtctrl"]
		aggregate         = ["shared-rtctrl", "export-rtctrl"]
	}`, testAccMSOSchemaTemplateExtEpgSubnetPrerequisiteConfig(), msoSchemaTemplateExtEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateExtEpgSubnetIp, msoSchemaTemplateExtEpgSubnetName)
}

func testAccMSOSchemaTemplateExtEpgSubnetDatasource() string {
	return fmt.Sprintf(`%s
	data "mso_schema_template_external_epg_subnet" "ext_epg_subnet" {
		schema_id         = mso_schema.%[2]s.id
		template_name     = "%[3]s"
		external_epg_name = mso_schema_template_external_epg.%[4]s.external_epg_name
		ip                = mso_schema_template_external_epg_subnet.%[4]s_subnet.ip
	}`, testAccMSOSchemaTemplateExtEpgSubnetDatasourceConfigCreate(), msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateExtEpgName)
}

func testAccMSOSchemaTemplateExtEpgSubnetDatasourceNotFound() string {
	return fmt.Sprintf(`%s
	data "mso_schema_template_external_epg_subnet" "ext_epg_subnet" {
		schema_id         = mso_schema.%[2]s.id
		template_name     = "%[3]s"
		external_epg_name = mso_schema_template_external_epg.%[4]s.external_epg_name
		ip                = "10.99.99.1/32"
	}`, testAccMSOSchemaTemplateExtEpgSubnetDatasourceConfigCreate(), msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateExtEpgName)
}
