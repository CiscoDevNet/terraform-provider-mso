package mso

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMSOTenantPoliciesNetflowRecordResource(t *testing.T) {
	resourceName := "mso_tenant_policies_netflow_record.netflow_record"
	var templateID, uuid string

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccVersionCheck(t, "5.1") },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create NetFlow Record") },
				Config:    testAccMSOTenantPoliciesNetflowRecordConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_netflow_record.netflow_record", "name", "test_netflow_record"),
					resource.TestCheckResourceAttr("mso_tenant_policies_netflow_record.netflow_record", "description", "Test NetFlow Record"),
					resource.TestCheckResourceAttr("mso_tenant_policies_netflow_record.netflow_record", "match_parameters.#", "0"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_netflow_record.netflow_record", "uuid"),
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
					fmt.Println("Test: Recreate NetFlow Record after out-of-band deletion")
					if err := testAccDeletePolicyOutOfBand(testAccPreCheck(t), templateID, uuid, "tenantPolicyTemplate", "template", "netFlowRecords"); err != nil {
						t.Fatalf("delete %s out of band: %v", resourceName, err)
					}
				},
				Config: testAccMSOTenantPoliciesNetflowRecordConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_netflow_record.netflow_record", "name", "test_netflow_record"),
					resource.TestCheckResourceAttr("mso_tenant_policies_netflow_record.netflow_record", "description", "Test NetFlow Record"),
					resource.TestCheckResourceAttr("mso_tenant_policies_netflow_record.netflow_record", "match_parameters.#", "0"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_netflow_record.netflow_record", "uuid"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update NetFlow Record with Match Parameters") },
				Config:    testAccMSOTenantPoliciesNetflowRecordConfigUpdateMatchParams(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_netflow_record.netflow_record", "name", "test_netflow_record"),
					resource.TestCheckResourceAttr("mso_tenant_policies_netflow_record.netflow_record", "description", "Updated NetFlow Record"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_netflow_record.netflow_record", "uuid"),
					resource.TestCheckResourceAttr("mso_tenant_policies_netflow_record.netflow_record", "match_parameters.#", "2"),
					testCheckTypeSetStringElemAttr("mso_tenant_policies_netflow_record.netflow_record", "match_parameters", "ethertype"),
					testCheckTypeSetStringElemAttr("mso_tenant_policies_netflow_record.netflow_record", "match_parameters", "destination_mac"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update NetFlow Record Name") },
				Config:    testAccMSOTenantPoliciesNetflowRecordConfigUpdateName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_netflow_record.netflow_record", "name", "test_netflow_record_updated"),
					resource.TestCheckResourceAttr("mso_tenant_policies_netflow_record.netflow_record", "description", "Updated NetFlow Record"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_netflow_record.netflow_record", "uuid"),
					resource.TestCheckResourceAttr("mso_tenant_policies_netflow_record.netflow_record", "match_parameters.#", "2"),
					testCheckTypeSetStringElemAttr("mso_tenant_policies_netflow_record.netflow_record", "match_parameters", "ethertype"),
					testCheckTypeSetStringElemAttr("mso_tenant_policies_netflow_record.netflow_record", "match_parameters", "destination_mac"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update NetFlow Record Remove Match Params") },
				Config:    testAccMSOTenantPoliciesNetflowRecordConfigRemoveMatch(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_netflow_record.netflow_record", "name", "test_netflow_record_updated"),
					resource.TestCheckResourceAttr("mso_tenant_policies_netflow_record.netflow_record", "description", "Updated NetFlow Record (Removed Match Params)"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_netflow_record.netflow_record", "uuid"),
					resource.TestCheckResourceAttr("mso_tenant_policies_netflow_record.netflow_record", "match_parameters.#", "0"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import NetFlow Record") },
				ResourceName:      "mso_tenant_policies_netflow_record.netflow_record",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import missing NetFlow Record") },
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
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments("mso_tenant_policies_netflow_record", "tenantPolicyTemplate", "template", "netFlowRecords"),
	})
}

func testAccMSOTenantPoliciesNetflowRecordConfigCreate() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_netflow_record" "netflow_record" {
        template_id = mso_template.template_tenant.id
        name        = "test_netflow_record"
        description = "Test NetFlow Record"
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesNetflowRecordConfigUpdateMatchParams() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_netflow_record" "netflow_record" {
        template_id = mso_template.template_tenant.id
        name        = "test_netflow_record"
        description = "Updated NetFlow Record"
        match_parameters = ["ethertype", "destination_mac"]
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesNetflowRecordConfigUpdateName() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_netflow_record" "netflow_record" {
        template_id = mso_template.template_tenant.id
        name        = "test_netflow_record_updated"
        description = "Updated NetFlow Record"
        match_parameters = ["ethertype", "destination_mac"]
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesNetflowRecordConfigRemoveMatch() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_netflow_record" "netflow_record" {
        template_id = mso_template.template_tenant.id
        name        = "test_netflow_record_updated"
        description = "Updated NetFlow Record (Removed Match Params)"
        match_parameters = []
    }`, testAccMSOTemplateResourceTenantConfig())
}
