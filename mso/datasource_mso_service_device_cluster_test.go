package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMSOServiceDeviceClusterDataSource(t *testing.T) {
	dataSourceName := "data.mso_service_device_cluster.cluster"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: Service Device Cluster Data Source returns an error when the cluster does not exist")
				},
				Config:      testAccMSOServiceDeviceClusterDataSourceConfigNotFound(),
				ExpectError: regexp.MustCompile(`Policy name .* not found`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Service Device Cluster Data Source") },
				Config:    testAccMSOServiceDeviceClusterDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", "test_device_cluster"),
					resource.TestCheckResourceAttr(dataSourceName, "device_mode", "layer3"),
					resource.TestCheckResourceAttr(dataSourceName, "device_type", "firewall"),
					resource.TestCheckResourceAttrSet(dataSourceName, "template_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "uuid"),
					resource.TestCheckResourceAttr(dataSourceName, "interface_properties.#", "2"),
					CustomTestCheckCollectionElemAttrsByKeys(dataSourceName, "interface_properties", map[string]string{
						"name": "interface1",
					}, map[string]string{
						"load_balance_hashing":         "sourceIP",
						"min_threshold":                "10",
						"max_threshold":                "90",
						"threshold_down_action":        "permit",
						"external_epg_uuid":            "mso_schema_template_external_epg.epg1.uuid",
						"ipsla_monitoring_policy_uuid": "mso_tenant_policies_ipsla_monitoring_policy.ipsla1.uuid",
					}),
					CustomTestCheckCollectionElemAttrsByKeys(dataSourceName, "interface_properties", map[string]string{
						"name": "interface3",
					}, map[string]string{
						"bd_uuid": "mso_schema_template_bd.bd2.uuid",
					}),
				),
			},
		},
	})
}

func testAccMSOServiceDeviceClusterDataSourceConfig() string {
	return fmt.Sprintf(`%s
	data "mso_service_device_cluster" "cluster" {
		template_id = mso_service_device_cluster.cluster.template_id
		name        = mso_service_device_cluster.cluster.name
	}`, testAccMSOServiceDeviceClusterConfigUpdateTwoInterfaces())
}

// testAccMSOServiceDeviceClusterDataSourceConfigNotFound wraps the positive
// data source config and adds a second data source query against a name that
// was never created. Used as the first step so the expected lookup failure
// fires immediately; the positive step that follows leaves the dependency
// graph in place for the final destroy walk to clean up in reverse order.
func testAccMSOServiceDeviceClusterDataSourceConfigNotFound() string {
	return fmt.Sprintf(`%s
	data "mso_service_device_cluster" "missing" {
		template_id = mso_template.device_template.id
		name        = "does_not_exist_cluster"
	}`, testAccMSOServiceDeviceClusterDataSourceConfig())
}
