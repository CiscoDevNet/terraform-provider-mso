package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOSchemaTemplateVrfResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create Schema Template VRF") },
				Config:    testAccMSOSchemaTemplateVrfConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "name", "tf_test_schema_template_vrf"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "description", "Terraform test Schema Template VRF"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "display_name", "tf_test_schema_template_vrf"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "ip_data_plane_learning", "enabled"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "vzany", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "preferred_group", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "site_aware_policy_enforcement", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "layer3_multicast", "false"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Create Schema Template VRF by enabling Layer3 Multicast") },
				Config:    testAccMSOSchemaTemplateVrfConfigUpdateEnablingLayer3Multicast(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "name", "tf_test_schema_template_vrf"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "display_name", "tf_test_schema_template_vrf"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "ip_data_plane_learning", "enabled"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "description", "Terraform test Schema Template VRF with Layer3 Multicast"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "vzany", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "preferred_group", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "site_aware_policy_enforcement", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "layer3_multicast", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "rendezvous_points.#", "1"),
					customTestCheckResourceTypeSetAttr("mso_schema_template_vrf.schema_template_vrf", "rendezvous_points",
						map[string]string{
							"ip_address": "1.1.1.2",
							"type":       "static",
						},
					),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update Create Schema Template VRF by enabling Layer3 Multicast and adding an extra RP")
				},
				Config: testAccMSOSchemaTemplateVrfConfigUpdateAddingExtraRp(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "name", "tf_test_schema_template_vrf"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "display_name", "tf_test_schema_template_vrf"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "ip_data_plane_learning", "enabled"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "description", "Terraform test Schema Template VRF with Layer3 Multicast"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "vzany", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "preferred_group", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "site_aware_policy_enforcement", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "layer3_multicast", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "rendezvous_points.#", "2"),
					customTestCheckResourceTypeSetAttr("mso_schema_template_vrf.schema_template_vrf", "rendezvous_points",
						map[string]string{
							"ip_address": "1.1.1.2",
							"type":       "static",
						},
					),
					customTestCheckResourceTypeSetAttr("mso_schema_template_vrf.schema_template_vrf", "rendezvous_points",
						map[string]string{
							"ip_address": "1.1.1.3",
							"type":       "fabric",
						},
					),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update Create Schema Template VRF by enabling Layer3 Multicast and removing an extra RP")
				},
				Config: testAccMSOSchemaTemplateVrfConfigUpdateRemovingExtraRp(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "name", "tf_test_schema_template_vrf"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "display_name", "tf_test_schema_template_vrf"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "ip_data_plane_learning", "enabled"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "description", "Terraform test Schema Template VRF with Layer3 Multicast"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "vzany", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "preferred_group", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "site_aware_policy_enforcement", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "layer3_multicast", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf.schema_template_vrf", "rendezvous_points.#", "1"),
					customTestCheckResourceTypeSetAttr("mso_schema_template_vrf.schema_template_vrf", "rendezvous_points",
						map[string]string{
							"ip_address": "1.1.1.2",
							"type":       "static",
						},
					),
				),
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithArguments("mso_schema_template_vrf", "vrf"),
	})
}

func testAccMSOSchemaTemplateForVrfConfig() string {
	return fmt.Sprintf(`%s
	resource "mso_schema" "schema_for_vrf" {
		name = "tf_test_schema_for_vrf"
		template {
			name         = "tf_test_schema_template"
			display_name = "tf_test_schema_template"
			tenant_id    = mso_tenant.%s.id
			template_type = "aci_multi_site"
		}
} `, testAccMSOTenantPoliciesMcastRouteMapPolicyConfigCreate(), msoTemplateTenantName)
}

func testAccMSOSchemaTemplateVrfConfigCreate() string {
	return fmt.Sprintf(`%s
	resource "mso_schema_template_vrf" "schema_template_vrf" {
		schema_id                     = mso_schema.schema_for_vrf.id
		template                      = tolist(mso_schema.schema_for_vrf.template)[0].name
		name                          = "tf_test_schema_template_vrf"
		display_name                  = "tf_test_schema_template_vrf"
		description                   = "Terraform test Schema Template VRF"
		ip_data_plane_learning        = "enabled"
		vzany                         = false
		preferred_group               = false
		site_aware_policy_enforcement = false
		layer3_multicast              = false
	}`, testAccMSOSchemaTemplateForVrfConfig())
}

func testAccMSOSchemaTemplateVrfConfigUpdateEnablingLayer3Multicast() string {
	return fmt.Sprintf(`%s
	resource "mso_schema_template_vrf" "schema_template_vrf" {
		schema_id                     = mso_schema.schema_for_vrf.id
		template                      = tolist(mso_schema.schema_for_vrf.template)[0].name
		name                          = "tf_test_schema_template_vrf"
		display_name                  = "tf_test_schema_template_vrf"
		description                   = "Terraform test Schema Template VRF with Layer3 Multicast"
		ip_data_plane_learning        = "enabled"
		vzany                         = false
		preferred_group               = false
		site_aware_policy_enforcement = false
		layer3_multicast              = true
		rendezvous_points {
			ip_address                      = "1.1.1.2"
			type                            = "static"
			route_map_policy_multicast_uuid = mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast.uuid
		}
	}`, testAccMSOSchemaTemplateForVrfConfig())
}

func testAccMSOSchemaTemplateVrfConfigUpdateAddingExtraRp() string {
	return fmt.Sprintf(`%s
	resource "mso_schema_template_vrf" "schema_template_vrf" {
		schema_id                     = mso_schema.schema_for_vrf.id
		template                      = tolist(mso_schema.schema_for_vrf.template)[0].name
		name                          = "tf_test_schema_template_vrf"
		display_name                  = "tf_test_schema_template_vrf"
		description                   = "Terraform test Schema Template VRF with Layer3 Multicast"
		ip_data_plane_learning        = "enabled"
		vzany                         = false
		preferred_group               = false
		site_aware_policy_enforcement = false
		layer3_multicast              = true
		rendezvous_points {
			ip_address                      = "1.1.1.2"
			type                            = "static"
			route_map_policy_multicast_uuid = mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast.uuid
		}
		rendezvous_points {
			ip_address                      = "1.1.1.3"
			type                            = "fabric"
		}
	}`, testAccMSOSchemaTemplateForVrfConfig())
}

func testAccMSOSchemaTemplateVrfConfigUpdateRemovingExtraRp() string {
	return fmt.Sprintf(`%s
	resource "mso_schema_template_vrf" "schema_template_vrf" {
		schema_id                     = mso_schema.schema_for_vrf.id
		template                      = tolist(mso_schema.schema_for_vrf.template)[0].name
		name                          = "tf_test_schema_template_vrf"
		display_name                  = "tf_test_schema_template_vrf"
		description                   = "Terraform test Schema Template VRF with Layer3 Multicast"
		ip_data_plane_learning        = "enabled"
		vzany                         = false
		preferred_group               = false
		site_aware_policy_enforcement = false
		layer3_multicast              = true
		rendezvous_points {
			ip_address                      = "1.1.1.2"
			type                            = "static"
			route_map_policy_multicast_uuid = mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast.uuid
		}
	}`, testAccMSOSchemaTemplateForVrfConfig())
}
