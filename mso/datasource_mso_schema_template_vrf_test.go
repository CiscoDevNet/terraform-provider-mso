package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOSchemaTemplateVrfDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Schema Template VRF Data Source") },
				Config:    testAccMSOSchemaTemplateVrfDataSource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_schema_template_vrf.schema_template_vrf", "name", "tf_test_schema_template_vrf"),
					resource.TestCheckResourceAttr("data.mso_schema_template_vrf.schema_template_vrf", "description", "Terraform test Schema Template VRF"),
					resource.TestCheckResourceAttr("data.mso_schema_template_vrf.schema_template_vrf", "display_name", "tf_test_schema_template_vrf"),
					resource.TestCheckResourceAttr("data.mso_schema_template_vrf.schema_template_vrf", "ip_data_plane_learning", "enabled"),
					resource.TestCheckResourceAttr("data.mso_schema_template_vrf.schema_template_vrf", "vzany", "false"),
					resource.TestCheckResourceAttr("data.mso_schema_template_vrf.schema_template_vrf", "preferred_group", "false"),
					resource.TestCheckResourceAttr("data.mso_schema_template_vrf.schema_template_vrf", "site_aware_policy_enforcement", "false"),
					resource.TestCheckResourceAttr("data.mso_schema_template_vrf.schema_template_vrf", "layer3_multicast", "false"),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateVrfDataSource() string {
	return fmt.Sprintf(`%s
	data "mso_schema_template_vrf" "schema_template_vrf" {
		schema_id              = mso_schema_template_vrf.schema_template_vrf.schema_id
		template               = mso_schema_template_vrf.schema_template_vrf.template
		name                   = "tf_test_schema_template_vrf"
	}`, testAccMSOSchemaTemplateVrfConfigCreate())
}
