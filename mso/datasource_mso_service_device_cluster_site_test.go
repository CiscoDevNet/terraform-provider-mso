package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMSOServiceDeviceClusterSiteDataSource(t *testing.T) {
	dataSourceName := "data.mso_service_device_cluster_site.cluster_site"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: Service Device Cluster Site Data Source returns an error when the device does not exist")
				},
				Config:      testAccMSOServiceDeviceClusterSiteDataSourceConfigNotFound(),
				ExpectError: regexp.MustCompile(`not found`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Service Device Cluster Site Data Source") },
				Config:    testAccMSOServiceDeviceClusterSiteDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", testServiceDeviceClusterSiteClusterName),
					resource.TestCheckResourceAttrSet(dataSourceName, "template_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "site_id"),
					resource.TestCheckResourceAttr(dataSourceName, "domain_type", "physicalDomain"),
					resource.TestCheckResourceAttr(dataSourceName, "domain_name", "test_physical_domain_for_device"),
					resource.TestCheckResourceAttr(dataSourceName, "interfaces.#", "1"),
					CustomTestCheckCollectionElemAttrsByKeys(dataSourceName, "interfaces", map[string]string{
						"name": "interface1",
					}, map[string]string{
						"vlan":                                      "210",
						"fabric_to_device_connectivity.#":           "1",
						"fabric_to_device_connectivity.0.pod_id":    "1",
						"fabric_to_device_connectivity.0.node_id.#": "1",
						"fabric_to_device_connectivity.0.node_id.0": "101",
						"fabric_to_device_connectivity.0.path":      "eth1/10",
						"fabric_to_device_connectivity.0.port_type": "port",
					}),
				),
			},
		},
	})
}

func testAccMSOServiceDeviceClusterSiteDataSourceConfig() string {
	return fmt.Sprintf(`%s
    data "mso_service_device_cluster_site" "cluster_site" {
        template_id = mso_service_device_cluster_site.cluster_site.template_id
        name        = mso_service_device_cluster_site.cluster_site.name
        site_id     = mso_service_device_cluster_site.cluster_site.site_id
    }
    `, testAccMSOServiceDeviceClusterSiteConfigOneInterface(testServiceDeviceClusterSiteClusterName))
}

// testAccMSOServiceDeviceClusterSiteDataSourceConfigNotFound wraps the
// positive data source config and adds a second data source query against a
// name that was never created on the site bucket. Used as the first step so
// the expected lookup failure fires immediately; the positive step that
// follows leaves the full deploy chain in place for the final destroy walk to
// clean up in reverse order.
func testAccMSOServiceDeviceClusterSiteDataSourceConfigNotFound() string {
	return fmt.Sprintf(`%s
    data "mso_service_device_cluster_site" "missing" {
        template_id = mso_template.device_template.id
        site_id     = data.mso_site.%[2]s.id
        name        = "does_not_exist_cluster_site"
    }
    `, testAccMSOServiceDeviceClusterSiteDataSourceConfig(), msoTemplateSiteName1)
}
