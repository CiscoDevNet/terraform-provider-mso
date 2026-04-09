package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOSchemaTemplateBdDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Schema Template BD Data Source - Not Found") },
				Config:      testAccMSOSchemaTemplateBdDataSourceNotFound(),
				ExpectError: regexp.MustCompile(`Unable to find the BD`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Schema Template BD Data Source") },
				Config:    testAccMSOSchemaTemplateBdDataSource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "display_name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "description", "Terraform test BD"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "layer2_unknown_unicast", "flood"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "unknown_multicast_flooding", "flood"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "multi_destination_flooding", "flood_in_bd"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "ipv6_unknown_multicast_flooding", "flood"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "arp_flooding", "true"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "unicast_routing", "false"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "layer3_multicast", "false"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "intersite_bum_traffic", "true"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "optimize_wan_bandwidth", "false"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "layer2_stretch", "true"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "virtual_mac_address", ""),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "ep_move_detection_mode", "none"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "vrf_name", msoSchemaTemplateVrfL3MulticastName),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateBdDataSource() string {
	return fmt.Sprintf(`%[1]s
	data "mso_schema_template_bd" "%[2]s" {
		schema_id     = mso_schema_template_bd.%[2]s.schema_id
		template_name = mso_schema_template_bd.%[2]s.template_name
		name          = "%[2]s"
	}`, testAccMSOSchemaTemplateBdConfigCreate(), msoSchemaTemplateBdName)
}

func testAccMSOSchemaTemplateBdDataSourceNotFound() string {
	return fmt.Sprintf(`%[1]s
	data "mso_schema_template_bd" "%[2]s" {
		schema_id     = mso_schema_template_bd.%[2]s.schema_id
		template_name = mso_schema_template_bd.%[2]s.template_name
		name          = "non_existing_bd"
	}`, testAccMSOSchemaTemplateBdConfigCreate(), msoSchemaTemplateBdName)
}
