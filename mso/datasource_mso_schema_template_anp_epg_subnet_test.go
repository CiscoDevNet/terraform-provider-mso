package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOSchemaTemplateAnpEpgSubnetDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateAnpEpgSubnetDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Read EPG Subnet datasource not found error") },
				Config:      testAccMSOSchemaTemplateAnpEpgSubnetDatasourceNotFound(),
				ExpectError: regexp.MustCompile("Unable to find the ANP EPG Subnet"),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read EPG Subnet datasource") },
				Config:    testAccMSOSchemaTemplateAnpEpgSubnetDatasource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.mso_schema_template_anp_epg_subnet.subnet", "schema_id"),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg_subnet.subnet", "template", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg_subnet.subnet", "anp_name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg_subnet.subnet", "epg_name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg_subnet.subnet", "ip", msoSchemaTemplateAnpEpgSubnetIp),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg_subnet.subnet", "scope", "private"),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg_subnet.subnet", "shared", "false"),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg_subnet.subnet", "querier", "false"),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg_subnet.subnet", "no_default_gateway", "false"),
					resource.TestCheckResourceAttr("data.mso_schema_template_anp_epg_subnet.subnet", "primary", "false"),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateAnpEpgSubnetDatasource() string {
	return fmt.Sprintf(`%s
	data "mso_schema_template_anp_epg_subnet" "subnet" {
		schema_id = mso_schema.%[2]s.id
		template  = "%[3]s"
		anp_name  = "%[4]s"
		epg_name  = "%[5]s"
		ip        = mso_schema_template_anp_epg_subnet.%[6]s_subnet.ip
	}`, testAccMSOSchemaTemplateAnpEpgSubnetConfigCreate(), msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateAnpEpgName)
}

func testAccMSOSchemaTemplateAnpEpgSubnetDatasourceNotFound() string {
	return fmt.Sprintf(`%s
	data "mso_schema_template_anp_epg_subnet" "subnet" {
		schema_id = mso_schema.%[2]s.id
		template  = "%[3]s"
		anp_name  = "%[4]s"
		epg_name  = "%[5]s"
		ip        = "99.99.99.99/32"
	}`, testAccMSOSchemaTemplateAnpEpgSubnetConfigCreate(), msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName)
}
