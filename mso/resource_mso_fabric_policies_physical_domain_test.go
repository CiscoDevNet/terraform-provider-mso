package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOPhysicalDomainResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create Physical Domain") },
				Config:    testAccMSOPhysicalDomainConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_physical_domain.physical_domain", "name", "tf_test_physical_domain"),
					resource.TestCheckResourceAttr("mso_fabric_policies_physical_domain.physical_domain", "description", "Terraform test Physical Domain"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_physical_domain.physical_domain", "vlan_pool_uuid"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Physical Domain") },
				Config:    testAccMSOPhysicalDomainConfigUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_physical_domain.physical_domain", "name", "tf_test_physical_domain"),
					resource.TestCheckResourceAttr("mso_fabric_policies_physical_domain.physical_domain", "description", "Terraform test Physical Domain Updated"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_physical_domain.physical_domain", "vlan_pool_uuid"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Physical Domain") },
				ResourceName:      "mso_fabric_policies_physical_domain.physical_domain",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments("mso_fabric_policies_physical_domain", "fabricPolicyTemplate", "template", "physicalDomain"),
	})
}

func testAccMSOPhysicalDomainConfigCreate() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_physical_domain" "physical_domain" {
		template_id     = mso_template.template_fabric_policy.id
		name            = "tf_test_physical_domain"
		description     = "Terraform test Physical Domain"
		vlan_pool_uuid  = mso_fabric_policies_vlan_pool.vlan_pool.uuid
	}`, testAccMSOVlanPoolConfigCreate())
}

func testAccMSOPhysicalDomainConfigUpdate() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_physical_domain" "physical_domain" {
		template_id     = mso_template.template_fabric_policy.id
		name            = "tf_test_physical_domain"
		description     = "Terraform test Physical Domain Updated"
		vlan_pool_uuid  = mso_fabric_policies_vlan_pool.vlan_pool.uuid
	}`, testAccMSOVlanPoolConfigCreate())
}
