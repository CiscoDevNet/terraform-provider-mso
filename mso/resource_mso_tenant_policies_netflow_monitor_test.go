package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMSOTenantPoliciesNetflowMonitorResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccVersionCheck(t, "5.1") },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create NetFlow Monitor") },
				Config:    testAccMSOTenantPoliciesNetflowMonitorConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_netflow_monitor.netflow_monitor", "name", "test_netflow_monitor"),
					resource.TestCheckResourceAttr("mso_tenant_policies_netflow_monitor.netflow_monitor", "description", "Test NetFlow Monitor"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_netflow_monitor.netflow_monitor", "uuid"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_netflow_monitor.netflow_monitor", "netflow_record_uuid"),
					resource.TestCheckResourceAttr("mso_tenant_policies_netflow_monitor.netflow_monitor", "netflow_exporter_uuids.#", "1"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update NetFlow Monitor") },
				Config:    testAccMSOTenantPoliciesNetflowMonitorConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_netflow_monitor.netflow_monitor", "name", "test_netflow_monitor_updated"),
					resource.TestCheckResourceAttr("mso_tenant_policies_netflow_monitor.netflow_monitor", "description", "Updated NetFlow Monitor"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_netflow_monitor.netflow_monitor", "uuid"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_netflow_monitor.netflow_monitor", "netflow_record_uuid"),
					resource.TestCheckResourceAttr("mso_tenant_policies_netflow_monitor.netflow_monitor", "netflow_exporter_uuids.#", "2"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update NetFlow Remove Description and Record UUID") },
				Config:    testAccMSOTenantPoliciesNetflowMonitorConfigRemoveOptionals(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_netflow_monitor.netflow_monitor", "name", "test_netflow_monitor_updated"),
					resource.TestCheckResourceAttr("mso_tenant_policies_netflow_monitor.netflow_monitor", "description", ""),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_netflow_monitor.netflow_monitor", "uuid"),
					resource.TestCheckResourceAttr("mso_tenant_policies_netflow_monitor.netflow_monitor", "netflow_record_uuid", ""),
					resource.TestCheckResourceAttr("mso_tenant_policies_netflow_monitor.netflow_monitor", "netflow_exporter_uuids.#", "2"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import NetFlow Monitor") },
				ResourceName:      "mso_tenant_policies_netflow_monitor.netflow_monitor",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments("mso_tenant_policies_netflow_monitor", "tenantPolicyTemplate", "template", "netFlowMonitors"),
	})
}

func testAccNeflowMonitorBaseConfig(extra_exporter bool) string {
	base := fmt.Sprintf(`%s
    resource "mso_tenant_policies_netflow_exporter" "netflow_exporter" {
        template_id = mso_template.template_tenant.id
        name        = "test_netflow_exporter"
    }
    resource "mso_tenant_policies_netflow_record" "netflow_record" {
        template_id = mso_template.template_tenant.id
        name        = "test_netflow_record"
    }`, testAccMSOTemplateResourceTenantConfig())

	if extra_exporter {
		return fmt.Sprintf(`%s
        resource "mso_tenant_policies_netflow_exporter" "netflow_exporter_2" {
            template_id = mso_template.template_tenant.id
            name        = "test_netflow_exporter_2"
            // depends_on ensures exporters are deleted sequentially to avoid index conflicts.
            depends_on = [mso_tenant_policies_netflow_exporter.netflow_exporter]
        }`, base)
	}

	return base
}

func testAccMSOTenantPoliciesNetflowMonitorConfigCreate() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_netflow_monitor" "netflow_monitor" {
        template_id                = mso_template.template_tenant.id
        name                       = "test_netflow_monitor"
        description                = "Test NetFlow Monitor"
        netflow_record_uuid        = mso_tenant_policies_netflow_record.netflow_record.uuid
        netflow_exporter_uuids     = [mso_tenant_policies_netflow_exporter.netflow_exporter.uuid]
    }`, testAccNeflowMonitorBaseConfig(false))
}

func testAccMSOTenantPoliciesNetflowMonitorConfigUpdate() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_netflow_monitor" "netflow_monitor" {
        template_id                = mso_template.template_tenant.id
        name                       = "test_netflow_monitor_updated"
        description                = "Updated NetFlow Monitor"
        netflow_record_uuid        = mso_tenant_policies_netflow_record.netflow_record.uuid
        netflow_exporter_uuids     = [mso_tenant_policies_netflow_exporter.netflow_exporter.uuid, mso_tenant_policies_netflow_exporter.netflow_exporter_2.uuid]
    }`, testAccNeflowMonitorBaseConfig(true))
}

func testAccMSOTenantPoliciesNetflowMonitorConfigRemoveOptionals() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_netflow_monitor" "netflow_monitor" {
        template_id                = mso_template.template_tenant.id
        name                       = "test_netflow_monitor_updated"
        netflow_exporter_uuids     = [mso_tenant_policies_netflow_exporter.netflow_exporter.uuid, mso_tenant_policies_netflow_exporter.netflow_exporter_2.uuid]
    }`, testAccNeflowMonitorBaseConfig(true))
}
