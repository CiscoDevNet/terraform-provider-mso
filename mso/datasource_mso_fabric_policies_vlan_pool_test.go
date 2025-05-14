package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOVlanPoolDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: VLAN Pool Data Source") },
				Config:    testAccMSOVlanPoolDataSource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_fabric_policies_vlan_pool.vlan_pool", "name", "tf_test_vlan_pool"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_vlan_pool.vlan_pool", "description", "Terraform test VLAN Pool"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_vlan_pool.vlan_pool", "vlan_range.#", "1"),
					customTestCheckResourceTypeSetAttr("data.mso_fabric_policies_vlan_pool.vlan_pool", "vlan_range",
						map[string]string{
							"from": "200",
							"to":   "202",
						},
					),
				),
			},
		},
	})
}

func testAccMSOVlanPoolDataSource() string {
	return fmt.Sprintf(`%s
	data "mso_fabric_policies_vlan_pool" "vlan_pool" {
	    template_id        = mso_fabric_policies_vlan_pool.vlan_pool.template_id
	    name               = "tf_test_vlan_pool"
    }`, testAccMSOVlanPoolConfigCreate())
}
