package mso

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMSOPhysicalDomainResource(t *testing.T) {
	resourceName := "mso_fabric_policies_physical_domain.physical_domain"
	var templateID, uuid string

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
					fmt.Println("Test: Recreate Physical Domain after out-of-band deletion")
					if err := testAccDeletePolicyOutOfBand(testAccPreCheck(t), templateID, uuid, "fabricPolicyTemplate", "template", "domains"); err != nil {
						t.Fatalf("delete %s out of band: %v", resourceName, err)
					}
				},
				Config: testAccMSOPhysicalDomainConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_physical_domain.physical_domain", "name", "tf_test_physical_domain"),
					resource.TestCheckResourceAttr("mso_fabric_policies_physical_domain.physical_domain", "description", "Terraform test Physical Domain"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_physical_domain.physical_domain", "vlan_pool_uuid"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
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
			{
				PreConfig:    func() { fmt.Println("Test: Import missing Physical Domain") },
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
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments("mso_fabric_policies_physical_domain", "fabricPolicyTemplate", "template", "domains"),
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
