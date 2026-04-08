package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOSchemaTemplateVrfDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateVrfDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Read VRF datasource not found error") },
				Config:      testAccMSOSchemaTemplateVrfDatasourceNotFound(),
				ExpectError: regexp.MustCompile("Unable to find the VRF"),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read VRF datasource") },
				Config:    testAccMSOSchemaTemplateVrfDatasource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.mso_schema_template_vrf.vrf", "schema_id"),
					resource.TestCheckResourceAttr("data.mso_schema_template_vrf.vrf", "template", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("data.mso_schema_template_vrf.vrf", "name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr("data.mso_schema_template_vrf.vrf", "display_name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr("data.mso_schema_template_vrf.vrf", "description", "Terraform test VRF"),
					resource.TestCheckResourceAttr("data.mso_schema_template_vrf.vrf", "ip_data_plane_learning", "enabled"),
					resource.TestCheckResourceAttr("data.mso_schema_template_vrf.vrf", "vzany", "false"),
					resource.TestCheckResourceAttr("data.mso_schema_template_vrf.vrf", "preferred_group", "false"),
					resource.TestCheckResourceAttr("data.mso_schema_template_vrf.vrf", "site_aware_policy_enforcement", "false"),
					resource.TestCheckResourceAttr("data.mso_schema_template_vrf.vrf", "layer3_multicast", "true"),
					resource.TestCheckResourceAttr("data.mso_schema_template_vrf.vrf", "rendezvous_points.#", "1"),
					customTestCheckResourceTypeSetAttr("data.mso_schema_template_vrf.vrf", "rendezvous_points",
						map[string]string{
							"ip_address": "1.1.1.2",
							"type":       "static",
						},
					),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateVrfDatasourceConfig() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_vrf" "%[2]s" {
		schema_id        = mso_schema.%[3]s.id
		template         = "%[4]s"
		name             = "%[2]s"
		display_name     = "%[2]s"
		description      = "Terraform test VRF"
		layer3_multicast = true
		rendezvous_points {
			ip_address                      = "1.1.1.2"
			type                            = "static"
			route_map_policy_multicast_uuid = mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast.uuid
		}
	}`, testAccMSOSchemaTemplateVrfPrerequisiteConfig(), msoSchemaTemplateVrfName, msoSchemaName, msoSchemaTemplateName)
}

func testAccMSOSchemaTemplateVrfDatasource() string {
	return fmt.Sprintf(`%s
	data "mso_schema_template_vrf" "vrf" {
		schema_id = mso_schema.%[2]s.id
		template  = "%[3]s"
		name      = mso_schema_template_vrf.%[4]s.name
	}`, testAccMSOSchemaTemplateVrfDatasourceConfig(), msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfName)
}

func testAccMSOSchemaTemplateVrfDatasourceNotFound() string {
	return fmt.Sprintf(`%s
	data "mso_schema_template_vrf" "vrf" {
		schema_id = mso_schema.%[2]s.id
		template  = "%[3]s"
		name      = "non_existing_vrf_name"
	}`, testAccMSOSchemaTemplateVrfDatasourceConfig(), msoSchemaName, msoSchemaTemplateName)
}
