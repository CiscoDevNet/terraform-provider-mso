package mso

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMSODHCPOptionPolicyResource(t *testing.T) {
	resourceName := "mso_tenant_policies_dhcp_option_policy.dhcp_policy"
	var templateID, uuid string

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create DHCP Option Policy") },
				Config:    testAccMSODHCPOptionPolicyConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_dhcp_option_policy.dhcp_policy", "name", "test_dhcp_option_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_dhcp_option_policy.dhcp_policy", "description", "Test DHCP Option Policy"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_dhcp_option_policy.dhcp_policy", "uuid"),

					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_dhcp_option_policy.dhcp_policy", "options", map[string]string{
						"name": "option_1",
						"id":   "1",
						"data": "data_1",
					}),
					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_dhcp_option_policy.dhcp_policy", "options", map[string]string{
						"name": "option_2",
						"id":   "2",
						"data": "data_2",
					}),
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
					fmt.Println("Test: Recreate DHCP Option Policy after out-of-band deletion")
					if err := testAccDeletePolicyOutOfBand(testAccPreCheck(t), templateID, uuid, "tenantPolicyTemplate", "template", "dhcpOptionPolicies"); err != nil {
						t.Fatalf("delete %s out of band: %v", resourceName, err)
					}
				},
				Config: testAccMSODHCPOptionPolicyConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_dhcp_option_policy.dhcp_policy", "name", "test_dhcp_option_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_dhcp_option_policy.dhcp_policy", "description", "Test DHCP Option Policy"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_dhcp_option_policy.dhcp_policy", "uuid"),

					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_dhcp_option_policy.dhcp_policy", "options", map[string]string{
						"name": "option_1",
						"id":   "1",
						"data": "data_1",
					}),
					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_dhcp_option_policy.dhcp_policy", "options", map[string]string{
						"name": "option_2",
						"id":   "2",
						"data": "data_2",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update DHCP Option Policy - Add Option") },
				Config:    testAccMSODHCPOptionPolicyConfigAddOption(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_dhcp_option_policy.dhcp_policy", "name", "test_dhcp_option_policy"),
					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_dhcp_option_policy.dhcp_policy", "options", map[string]string{
						"name": "option_1",
						"id":   "1",
						"data": "data_1",
					}),
					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_dhcp_option_policy.dhcp_policy", "options", map[string]string{
						"name": "option_2",
						"id":   "2",
						"data": "data_2",
					}),
					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_dhcp_option_policy.dhcp_policy", "options", map[string]string{
						"name": "option_3",
						"id":   "3",
						"data": "data_3",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update DHCP Option Policy - Remove Option") },
				Config:    testAccMSODHCPOptionPolicyConfigRemoveOption(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_dhcp_option_policy.dhcp_policy", "name", "test_dhcp_option_policy"),
					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_dhcp_option_policy.dhcp_policy", "options", map[string]string{
						"name": "option_1",
						"id":   "1",
						"data": "data_1",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update DHCP Option Policy - Option Without ID") },
				Config:    testAccMSODHCPOptionPolicyConfigOptionWithoutId(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_dhcp_option_policy.dhcp_policy", "name", "test_dhcp_option_policy"),
					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_dhcp_option_policy.dhcp_policy", "options", map[string]string{
						"name": "option_no_id",
						"data": "data_no_id",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update DHCP Option Policy - Option Without Data") },
				Config:    testAccMSODHCPOptionPolicyConfigOptionWithoutData(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_dhcp_option_policy.dhcp_policy", "name", "test_dhcp_option_policy"),
					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_dhcp_option_policy.dhcp_policy", "options", map[string]string{
						"name": "option_only_name_id",
						"id":   "100",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update DHCP Option Policy Name") },
				Config:    testAccMSODHCPOptionPolicyConfigUpdateName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_dhcp_option_policy.dhcp_policy", "name", "test_dhcp_option_policy_renamed"),
					resource.TestCheckResourceAttr("mso_tenant_policies_dhcp_option_policy.dhcp_policy", "description", "Renamed Policy"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import DHCP Option Policy") },
				ResourceName:      "mso_tenant_policies_dhcp_option_policy.dhcp_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import missing DHCP Option Policy") },
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
		CheckDestroy: testCheckResourceDestroyPolicyWithArguments("mso_tenant_policies_dhcp_option_policy", "dhcpOption"),
	})
}

func testAccMSODHCPOptionPolicyConfigCreate() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_dhcp_option_policy" "dhcp_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_dhcp_option_policy"
        description = "Test DHCP Option Policy"
        
        options {
            name = "option_1"
            id   = 1
            data = "data_1"
        }
        
        options {
            name = "option_2"
            id   = 2
            data = "data_2"
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSODHCPOptionPolicyConfigAddOption() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_dhcp_option_policy" "dhcp_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_dhcp_option_policy"
        description = "Test DHCP Option Policy"
        
        options {
            name = "option_1"
            id   = 1
            data = "data_1"
        }
        
        options {
            name = "option_2"
            id   = 2
            data = "data_2"
        }
        
        options {
            name = "option_3"
            id   = 3
            data = "data_3"
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSODHCPOptionPolicyConfigRemoveOption() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_dhcp_option_policy" "dhcp_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_dhcp_option_policy"
        description = "Test DHCP Option Policy"
        
        options {
            name = "option_1"
            id   = 1
            data = "data_1"
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSODHCPOptionPolicyConfigOptionWithoutId() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_dhcp_option_policy" "dhcp_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_dhcp_option_policy"
        description = "Test DHCP Option Policy"
        
        options {
            name = "option_no_id"
            data = "data_no_id"
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSODHCPOptionPolicyConfigOptionWithoutData() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_dhcp_option_policy" "dhcp_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_dhcp_option_policy"
        description = "Test DHCP Option Policy"
        
        options {
            name = "option_only_name_id"
            id   = 100
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSODHCPOptionPolicyConfigUpdateName() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_dhcp_option_policy" "dhcp_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_dhcp_option_policy_renamed"
        description = "Renamed Policy"
        
        options {
            name = "option_1"
            id   = 1
            data = "data_1"
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}
