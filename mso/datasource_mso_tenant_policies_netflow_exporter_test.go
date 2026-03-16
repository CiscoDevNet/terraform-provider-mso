package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOTenantPoliciesNetflowExporterDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: NetFlow Exporter Data Source") },
				Config:    testAccMSOTenantPoliciesNetflowExporterDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_tenant_policies_netflow_exporter.netflow_exporter", "name", "test_netflow_exporter"),
					resource.TestCheckResourceAttrSet("data.mso_tenant_policies_netflow_exporter.netflow_exporter", "uuid"),
				),
			},
		},
	})
}

func testAccMSOTenantPoliciesNetflowExporterDataSource() string {
	return fmt.Sprintf(`%s
    data "mso_tenant_policies_netflow_exporter" "netflow_exporter" {
        template_id = mso_tenant_policies_netflow_exporter.netflow_exporter.template_id
        name        = "test_netflow_exporter"
    }`, testAccMSOTenantPoliciesNetflowExporterConfigCreate())
}