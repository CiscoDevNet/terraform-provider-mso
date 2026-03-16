package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOTenantPoliciesNetflowExporterResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create NetFlow Exporter") },
				Config:    testAccMSOTenantPoliciesNetflowExporterConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_netflow_exporter.netflow_exporter", "name", "test_netflow_exporter"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_netflow_exporter.netflow_exporter", "uuid"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update NetFlow Exporter Name") },
				Config:    testAccMSOTenantPoliciesNetflowExporterConfigUpdateName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_netflow_exporter.netflow_exporter", "name", "test_netflow_exporter_updated"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_netflow_exporter.netflow_exporter", "uuid"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import NetFlow Exporter") },
				ResourceName:      "mso_tenant_policies_netflow_exporter.netflow_exporter",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments("mso_tenant_policies_netflow_exporter", "tenantPolicyTemplate", "template", "netFlowExporters"),
	})
}

func testAccMSOTenantPoliciesNetflowExporterConfigCreate() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_netflow_exporter" "netflow_exporter" {
        template_id = mso_template.template_tenant.id
        name        = "test_netflow_exporter"
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesNetflowExporterConfigUpdateName() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_netflow_exporter" "netflow_exporter" {
        template_id = mso_template.template_tenant.id
        name        = "test_netflow_exporter_updated"
    }`, testAccMSOTemplateResourceTenantConfig())
}