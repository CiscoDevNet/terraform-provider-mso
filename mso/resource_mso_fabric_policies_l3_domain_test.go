package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOL3DomainResource(t *testing.T) {
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
