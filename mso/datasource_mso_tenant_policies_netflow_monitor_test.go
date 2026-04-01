package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOTenantPoliciesNetflowMonitorDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: NetFlow Monitor Data Source") },
				Config:    testAccMSOTenantPoliciesNetflowMonitorDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_tenant_policies_netflow_monitor.netflow_monitor", "name", "test_netflow_monitor"),
					resource.TestCheckResourceAttrSet("data.mso_tenant_policies_netflow_monitor.netflow_monitor", "uuid"),
					resource.TestCheckResourceAttrSet("data.mso_tenant_policies_netflow_monitor.netflow_monitor", "netflow_record_uuid"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_netflow_monitor.netflow_monitor", "netflow_exporter_uuids.#", "1"),
					resource.TestCheckResourceAttrSet("data.mso_tenant_policies_netflow_monitor.netflow_monitor", "netflow_exporter_uuids.0"),
				),
			},
		},
	})
}

func testAccMSOTenantPoliciesNetflowMonitorDataSource() string {
	return fmt.Sprintf(`%s
    data "mso_tenant_policies_netflow_monitor" "netflow_monitor" {
        template_id = mso_tenant_policies_netflow_monitor.netflow_monitor.template_id
        name        = "test_netflow_monitor"
    }`, testAccMSOTenantPoliciesNetflowMonitorConfigCreate())
}
