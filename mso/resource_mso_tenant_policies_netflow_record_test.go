package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOTenantPoliciesNetflowRecordResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
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
