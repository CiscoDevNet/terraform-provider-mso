package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOServiceDeviceClusterDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Service Device Cluster Data Source") },
				Config:    testAccMSOServiceDeviceClusterDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_service_device_cluster.cluster", "name", "test_device_cluster"),
					resource.TestCheckResourceAttr("data.mso_service_device_cluster.cluster", "device_mode", "layer3"),
					resource.TestCheckResourceAttr("data.mso_service_device_cluster.cluster", "device_type", "firewall"),
					resource.TestCheckResourceAttr("data.mso_service_device_cluster.cluster", "interface_properties.#", "2"),
					CustomTestCheckTypeSetElemAttrs("data.mso_service_device_cluster.cluster", "interface_properties", map[string]string{
						"name":                  "interface1",
						"load_balance_hashing":  "sourceIP",
						"min_threshold":         "10",
						"max_threshold":         "90",
						"threshold_down_action": "permit",
					}),
					CustomTestCheckTypeSetElemAttrs("data.mso_service_device_cluster.cluster", "interface_properties", map[string]string{
						"name":    "interface3",
						"anycast": "true",
					}),
				),
			},
		},
	})
}

func testAccMSOServiceDeviceClusterDataSource() string {
	return fmt.Sprintf(`%s
	data "mso_service_device_cluster" "cluster" {
		template_id = mso_service_device_cluster.cluster.template_id
		name        = mso_service_device_cluster.cluster.name
	  }`, testAccMSOServiceDeviceClusterConfigUpdateTwoInterfaces())
}
