package mso

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMSOL3DomainResource(t *testing.T) {
	resourceName := "mso_fabric_policies_l3_domain.l3_domain"
	var templateID, uuid string

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create L3 Domain") },
				Config:    testAccMSOL3DomainConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_l3_domain.l3_domain", "name", "test_l3_domain"),
					resource.TestCheckResourceAttr("mso_fabric_policies_l3_domain.l3_domain", "description", "Test L3 Domain"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_l3_domain.l3_domain", "uuid"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_l3_domain.l3_domain", "vlan_pool_uuid"),
					resource.TestCheckResourceAttrPair(
						"mso_fabric_policies_l3_domain.l3_domain", "vlan_pool_uuid",
						"mso_fabric_policies_vlan_pool.vlan_pool", "uuid",
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
					fmt.Println("Test: Recreate L3 Domain after out-of-band deletion")
					if err := testAccDeletePolicyOutOfBand(testAccPreCheck(t), templateID, uuid, "fabricPolicyTemplate", "template", "l3Domains"); err != nil {
						t.Fatalf("delete %s out of band: %v", resourceName, err)
					}
				},
				Config: testAccMSOL3DomainConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_l3_domain.l3_domain", "name", "test_l3_domain"),
					resource.TestCheckResourceAttr("mso_fabric_policies_l3_domain.l3_domain", "description", "Test L3 Domain"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_l3_domain.l3_domain", "uuid"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_l3_domain.l3_domain", "vlan_pool_uuid"),
					resource.TestCheckResourceAttrPair(
						"mso_fabric_policies_l3_domain.l3_domain", "vlan_pool_uuid",
						"mso_fabric_policies_vlan_pool.vlan_pool", "uuid",
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3 Domain Description") },
				Config:    testAccMSOL3DomainConfigUpdateDescription(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_l3_domain.l3_domain", "name", "test_l3_domain"),
					resource.TestCheckResourceAttr("mso_fabric_policies_l3_domain.l3_domain", "description", "Updated L3 Domain Description"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_l3_domain.l3_domain", "vlan_pool_uuid"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3 Domain - Remove VLAN Pool") },
				Config:    testAccMSOL3DomainConfigRemoveVLANPool(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_l3_domain.l3_domain", "name", "test_l3_domain"),
					resource.TestCheckResourceAttr("mso_fabric_policies_l3_domain.l3_domain", "description", "L3 Domain without VLAN Pool"),
					resource.TestCheckResourceAttr("mso_fabric_policies_l3_domain.l3_domain", "vlan_pool_uuid", ""),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3 Domain - Re-add VLAN Pool") },
				Config:    testAccMSOL3DomainConfigReAddVLANPool(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_l3_domain.l3_domain", "name", "test_l3_domain"),
					resource.TestCheckResourceAttr("mso_fabric_policies_l3_domain.l3_domain", "description", "L3 Domain with VLAN Pool Re-added"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_l3_domain.l3_domain", "vlan_pool_uuid"),
					resource.TestCheckResourceAttrPair(
						"mso_fabric_policies_l3_domain.l3_domain", "vlan_pool_uuid",
						"mso_fabric_policies_vlan_pool.vlan_pool", "uuid",
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3 Domain Name with UUID") },
				Config:    testAccMSOL3DomainConfigUpdateName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_l3_domain.l3_domain", "name", "test_l3_domain_renamed"),
					resource.TestCheckResourceAttr("mso_fabric_policies_l3_domain.l3_domain", "description", "Renamed L3 Domain"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import L3 Domain") },
				ResourceName:      "mso_fabric_policies_l3_domain.l3_domain",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import missing L3 Domain") },
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
		CheckDestroy: testCheckResourceDestroyPolicyWithArguments("mso_fabric_policies_l3_domain", "l3Domain"),
	})
}

func testAccMSOL3DomainConfigCreate() string {
	return fmt.Sprintf(`%s
    resource "mso_fabric_policies_l3_domain" "l3_domain" {
        template_id    = mso_template.template_fabric_policy.id
        name           = "test_l3_domain"
        description    = "Test L3 Domain"
        vlan_pool_uuid = mso_fabric_policies_vlan_pool.vlan_pool.uuid
    }`, testAccMSOVlanPoolConfigCreate())
}

func testAccMSOL3DomainConfigUpdateDescription() string {
	return fmt.Sprintf(`%s
    resource "mso_fabric_policies_l3_domain" "l3_domain" {
        template_id    = mso_template.template_fabric_policy.id
        name           = "test_l3_domain"
        description    = "Updated L3 Domain Description"
        vlan_pool_uuid = mso_fabric_policies_vlan_pool.vlan_pool.uuid
    }`, testAccMSOVlanPoolConfigCreate())
}

func testAccMSOL3DomainConfigRemoveVLANPool() string {
	return fmt.Sprintf(`%s
    resource "mso_fabric_policies_l3_domain" "l3_domain" {
        template_id    = mso_template.template_fabric_policy.id
        name           = "test_l3_domain"
        description    = "L3 Domain without VLAN Pool"
        vlan_pool_uuid = ""
    }`, testAccMSOVlanPoolConfigCreate())
}

func testAccMSOL3DomainConfigReAddVLANPool() string {
	return fmt.Sprintf(`%s
    resource "mso_fabric_policies_l3_domain" "l3_domain" {
        template_id    = mso_template.template_fabric_policy.id
        name           = "test_l3_domain"
        description    = "L3 Domain with VLAN Pool Re-added"
        vlan_pool_uuid = mso_fabric_policies_vlan_pool.vlan_pool.uuid
    }`, testAccMSOVlanPoolConfigCreate())
}

func testAccMSOL3DomainConfigUpdateName() string {
	return fmt.Sprintf(`%s
    resource "mso_fabric_policies_l3_domain" "l3_domain" {
        template_id    = mso_template.template_fabric_policy.id
        name           = "test_l3_domain_renamed"
        description    = "Renamed L3 Domain"
        vlan_pool_uuid = mso_fabric_policies_vlan_pool.vlan_pool.uuid
    }`, testAccMSOVlanPoolConfigCreate())
}
