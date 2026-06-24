package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOServiceDeviceClusterSiteDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Service Device Cluster Site Data Source") },
				Config:    testAccMSOServiceDeviceClusterSiteDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_service_device_cluster_site.cluster_site", "name", "test_device_cluster"),
					resource.TestCheckResourceAttrSet("data.mso_service_device_cluster_site.cluster_site", "template_id"),
					resource.TestCheckResourceAttrSet("data.mso_service_device_cluster_site.cluster_site", "site_id"),
				),
			},
		},
	})
}

func testAccMSOServiceDeviceClusterSiteDataSource() string {
	return fmt.Sprintf(`%s
    data "mso_service_device_cluster_site" "cluster_site" {
      template_id = mso_service_device_cluster_site.cluster_site.template_id
      name        = mso_service_device_cluster_site.cluster_site.name
      site_id     = mso_service_device_cluster_site.cluster_site.site_id
    }
    `, testAccMSOServiceDeviceClusterSiteConfigCreate())
}
