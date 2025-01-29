package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOTenantPoliciesIPSLAMonitoringPolicyDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: IPSLA Monitoring Policy Data Source") },
				Config:    testAccMSOTenantPoliciesIPSLAMonitoringPolicyDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "name", "test_ipsla_policy"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "description", "HTTP Type"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "sla_type", "http"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "destination_port", "80"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "http_version", "HTTP11"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "http_uri", "/example"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "sla_frequency", "120"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "detect_multiplier", "4"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "request_data_size", "64"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "type_of_service", "18"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "operation_timeout", "100"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "threshold", "100"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "ipv6_traffic_class", "255"),
				),
			},
		},
	})
}

func testAccMSOTenantPoliciesIPSLAMonitoringPolicyDataSource() string {
	return fmt.Sprintf(`%s
	data "mso_tenant_policies_ipsla_monitoring_policy" "ipsla_policy" {
	    template_id        = mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy.template_id
	    name               = "test_ipsla_policy"
    }`, testAccMSOTenantPoliciesIPSLAMonitoringPolicyConfigCreate())
}
