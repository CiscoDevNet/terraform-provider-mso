package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOVlanPoolResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create VLAN Pool") },
				Config:    testAccMSOVlanPoolConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_vlan_pool.vlan_pool", "name", "tf_test_vlan_pool"),
					resource.TestCheckResourceAttr("mso_fabric_policies_vlan_pool.vlan_pool", "description", "Terraform test VLAN Pool"),
					resource.TestCheckResourceAttr("mso_fabric_policies_vlan_pool.vlan_pool", "allocation_mode", "static"),
					resource.TestCheckResourceAttr("mso_fabric_policies_vlan_pool.vlan_pool", "vlan_range.#", "1"),
					customTestCheckResourceTypeSetAttr("mso_fabric_policies_vlan_pool.vlan_pool", "vlan_range",
						map[string]string{
							"from":            "200",
							"to":              "202",
							"allocation_mode": "static",
						},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update VLAN Pool adding extra range") },
				Config:    testAccMSOVlanPoolConfigUpdateAddingExtraRange(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_vlan_pool.vlan_pool", "name", "tf_test_vlan_pool"),
					resource.TestCheckResourceAttr("mso_fabric_policies_vlan_pool.vlan_pool", "description", "Terraform test VLAN Pool adding extra range"),
					resource.TestCheckResourceAttr("mso_fabric_policies_vlan_pool.vlan_pool", "allocation_mode", "static"),
					resource.TestCheckResourceAttr("mso_fabric_policies_vlan_pool.vlan_pool", "vlan_range.#", "2"),
					customTestCheckResourceTypeSetAttr("mso_fabric_policies_vlan_pool.vlan_pool", "vlan_range",
						map[string]string{
							"from":            "200",
							"to":              "202",
							"allocation_mode": "static",
						},
					),
					customTestCheckResourceTypeSetAttr("mso_fabric_policies_vlan_pool.vlan_pool", "vlan_range",
						map[string]string{
							"from":            "204",
							"to":              "209",
							"allocation_mode": "static",
						},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update VLAN Pool removing extra range") },
				Config:    testAccMSOVlanPoolConfigUpdateRemovingExtraRange(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_vlan_pool.vlan_pool", "name", "tf_test_vlan_pool"),
					resource.TestCheckResourceAttr("mso_fabric_policies_vlan_pool.vlan_pool", "description", "Terraform test VLAN Pool removing extra range"),
					resource.TestCheckResourceAttr("mso_fabric_policies_vlan_pool.vlan_pool", "allocation_mode", "static"),
					resource.TestCheckResourceAttr("mso_fabric_policies_vlan_pool.vlan_pool", "vlan_range.#", "1"),
					customTestCheckResourceTypeSetAttr("mso_fabric_policies_vlan_pool.vlan_pool", "vlan_range",
						map[string]string{
							"from":            "200",
							"to":              "202",
							"allocation_mode": "static",
						},
					),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import VLAN Pool") },
				ResourceName:      "mso_fabric_policies_vlan_pool.vlan_pool",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithTemplateIdAndArguments("mso_fabric_policies_vlan_pool", "fabricPolicyTemplate", "template", "vlanPools"),
	})
}

func testAccMSOVlanPoolConfigCreate() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_vlan_pool" "vlan_pool" {
		template_id     = mso_template.template_fabric_policy.id
		name            = "tf_test_vlan_pool"
		description     = "Terraform test VLAN Pool"
		allocation_mode = "static"
		vlan_range {
			from            = 200
			to              = 202
			allocation_mode = "static"
		}
	}`, testAccMSOTemplateResourceFabricPolicyConfig())
}

func testAccMSOVlanPoolConfigUpdateAddingExtraRange() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_vlan_pool" "vlan_pool" {
		template_id     = mso_template.template_fabric_policy.id
		name            = "tf_test_vlan_pool"
		description     = "Terraform test VLAN Pool adding extra range"
		allocation_mode = "static"
		vlan_range {
			from            = 200
			to              = 202
			allocation_mode = "static"
		}
		vlan_range {
			from            = 204
			to              = 209
			allocation_mode = "static"
		}
	}`, testAccMSOTemplateResourceFabricPolicyConfig())
}

func testAccMSOVlanPoolConfigUpdateRemovingExtraRange() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_vlan_pool" "vlan_pool" {
		template_id     = mso_template.template_fabric_policy.id
		name            = "tf_test_vlan_pool"
		description     = "Terraform test VLAN Pool removing extra range"
		allocation_mode = "static"
		vlan_range {
			from            = 200
			to              = 202
			allocation_mode = "static"
		}
	}`, testAccMSOTemplateResourceFabricPolicyConfig())
}
