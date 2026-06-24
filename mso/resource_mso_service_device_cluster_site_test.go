package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOServiceDeviceClusterSiteResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create Service Device Cluster Site") },
				Config:    testAccMSOServiceDeviceClusterSiteConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_service_device_cluster_site.cluster_site", "name", "test_device_cluster"),
					resource.TestCheckResourceAttrSet("mso_service_device_cluster_site.cluster_site", "template_id"),
					resource.TestCheckResourceAttrSet("mso_service_device_cluster_site.cluster_site", "site_id"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Service Device Cluster Site") },
				ResourceName:      "mso_service_device_cluster_site.cluster_site",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMSOServiceDeviceClusterSiteConfigCreate() string {
	return fmt.Sprintf(`%s
    resource "mso_service_device_cluster_site" "cluster_site" {
      template_id = mso_template.device_template.id
      name        = "test_device_cluster"
      site_id     = data.mso_site.%s.id
    }
    `, testAccMSOServiceDeviceClusterDependencies(), msoTemplateSiteName1)
}
