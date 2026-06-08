package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMSOSchemaTemplateBdSubnetDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateBdSubnetDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: BD Subnet Data Source - Not Found") },
				Config:      testAccMSOSchemaTemplateBdSubnetDataSourceNotFound(),
				ExpectError: regexp.MustCompile(`Unable to find the BD Subnet`),
			},
			{
				PreConfig: func() { fmt.Println("Test: BD Subnet Data Source") },
				Config:    testAccMSOSchemaTemplateBdSubnetDataSource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "schema_id"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "bd_name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "ip", msoSchemaTemplateBdSubnetIp),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "scope", "private"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "shared", "false"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "querier", "false"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "no_default_gateway", "false"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "primary", "false"),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateBdSubnetDataSource() string {
	return fmt.Sprintf(`%[1]s
data "mso_schema_template_bd_subnet" "%[2]s_subnet" {
	schema_id     = mso_schema_template_bd_subnet.%[2]s_subnet.schema_id
	template_name = mso_schema_template_bd_subnet.%[2]s_subnet.template_name
	bd_name       = mso_schema_template_bd_subnet.%[2]s_subnet.bd_name
	ip            = mso_schema_template_bd_subnet.%[2]s_subnet.ip
}`, testAccMSOSchemaTemplateBdSubnetConfigCreate(), msoSchemaTemplateBdName)
}

func testAccMSOSchemaTemplateBdSubnetDataSourceNotFound() string {
	return fmt.Sprintf(`%[1]s
data "mso_schema_template_bd_subnet" "%[2]s_subnet" {
	schema_id     = mso_schema_template_bd_subnet.%[2]s_subnet.schema_id
	template_name = mso_schema_template_bd_subnet.%[2]s_subnet.template_name
	bd_name       = mso_schema_template_bd_subnet.%[2]s_subnet.bd_name
	ip            = "99.99.99.99/32"
}`, testAccMSOSchemaTemplateBdSubnetConfigCreate(), msoSchemaTemplateBdName)
}
