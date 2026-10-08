package mso

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMSOVlanPoolResource(t *testing.T) {
	resourceName := "mso_fabric_policies_vlan_pool.vlan_pool"
	var templateID, uuid string

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
					resource.TestCheckResourceAttr("mso_fabric_policies_vlan_pool.vlan_pool", "vlan_range.#", "1"),
					customTestCheckResourceTypeSetAttr("mso_fabric_policies_vlan_pool.vlan_pool", "vlan_range",
						map[string]string{
							"from": "200",
							"to":   "202",
						},
					),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok || rs.Primary == nil {
							return fmt.Errorf("resource %s not found in state", resourceName)
						}
						templateID = rs.Primary.Attributes["template_id"]
						uuid = rs.Primary.Attributes["uuid"]
						if templateID == "" || uuid == "" {
							return fmt.Errorf("resource %s is missing its template ID or UUID", resourceName)
						}
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate VLAN Pool after out-of-band deletion")
					if err := testAccDeletePolicyOutOfBand(testAccPreCheck(t), templateID, uuid, "fabricPolicyTemplate", "template", "vlanPools"); err != nil {
						t.Fatalf("delete %s out of band: %v", resourceName, err)
					}
				},
				Config: testAccMSOVlanPoolConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_vlan_pool.vlan_pool", "name", "tf_test_vlan_pool"),
					resource.TestCheckResourceAttr("mso_fabric_policies_vlan_pool.vlan_pool", "description", "Terraform test VLAN Pool"),
					resource.TestCheckResourceAttr("mso_fabric_policies_vlan_pool.vlan_pool", "vlan_range.#", "1"),
					customTestCheckResourceTypeSetAttr("mso_fabric_policies_vlan_pool.vlan_pool", "vlan_range",
						map[string]string{
							"from": "200",
							"to":   "202",
						},
					),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update VLAN Pool adding extra range") },
				Config:    testAccMSOVlanPoolConfigUpdateAddingExtraRange(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_vlan_pool.vlan_pool", "name", "tf_test_vlan_pool"),
					resource.TestCheckResourceAttr("mso_fabric_policies_vlan_pool.vlan_pool", "description", "Terraform test VLAN Pool adding extra range"),
					resource.TestCheckResourceAttr("mso_fabric_policies_vlan_pool.vlan_pool", "vlan_range.#", "2"),
					customTestCheckResourceTypeSetAttr("mso_fabric_policies_vlan_pool.vlan_pool", "vlan_range",
						map[string]string{
							"from": "200",
							"to":   "202",
						},
					),
					customTestCheckResourceTypeSetAttr("mso_fabric_policies_vlan_pool.vlan_pool", "vlan_range",
						map[string]string{
							"from": "204",
							"to":   "209",
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
					resource.TestCheckResourceAttr("mso_fabric_policies_vlan_pool.vlan_pool", "vlan_range.#", "1"),
					customTestCheckResourceTypeSetAttr("mso_fabric_policies_vlan_pool.vlan_pool", "vlan_range",
						map[string]string{
							"from": "200",
							"to":   "202",
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
			{
				PreConfig:    func() { fmt.Println("Test: Import missing VLAN Pool") },
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok || rs.Primary == nil {
						return "", fmt.Errorf("resource %s not found in state", resourceName)
					}
					separator := strings.LastIndex(rs.Primary.ID, "/")
					if separator < 0 || separator == len(rs.Primary.ID)-1 {
						return "", fmt.Errorf("resource %s has an invalid import ID", resourceName)
					}
					return rs.Primary.ID[:separator+1] + "tf_missing_oob", nil
				},
				ExpectError: regexp.MustCompile(`cannot import .*: resource not found`),
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments("mso_fabric_policies_vlan_pool", "fabricPolicyTemplate", "template", "vlanPools"),
	})
}

func testAccMSOVlanPoolConfigCreate() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_vlan_pool" "vlan_pool" {
		template_id     = mso_template.template_fabric_policy.id
		name            = "tf_test_vlan_pool"
		description     = "Terraform test VLAN Pool"
		vlan_range {
			from            = 200
			to              = 202
		}
	}`, testAccMSOTemplateResourceFabricPolicyConfig())
}

func testAccMSOVlanPoolConfigUpdateAddingExtraRange() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_vlan_pool" "vlan_pool" {
		template_id     = mso_template.template_fabric_policy.id
		name            = "tf_test_vlan_pool"
		description     = "Terraform test VLAN Pool adding extra range"
		vlan_range {
			from            = 200
			to              = 202
		}
		vlan_range {
			from            = 204
			to              = 209
		}
	}`, testAccMSOTemplateResourceFabricPolicyConfig())
}

func testAccMSOVlanPoolConfigUpdateRemovingExtraRange() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_vlan_pool" "vlan_pool" {
		template_id     = mso_template.template_fabric_policy.id
		name            = "tf_test_vlan_pool"
		description     = "Terraform test VLAN Pool removing extra range"
		vlan_range {
			from            = 200
			to              = 202
		}
	}`, testAccMSOTemplateResourceFabricPolicyConfig())
}
