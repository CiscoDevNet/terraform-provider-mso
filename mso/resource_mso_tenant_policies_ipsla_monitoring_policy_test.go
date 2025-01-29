package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOTenantPoliciesIPSLAMonitoringPolicyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create IPSLA Monitoring Policy") },
				Config:    testAccMSOTenantPoliciesIPSLAMonitoringPolicyConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "name", "test_ipsla_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "description", "HTTP Type"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "sla_type", "http"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "destination_port", "80"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "http_version", "HTTP11"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "http_uri", "/example"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "sla_frequency", "120"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "detect_multiplier", "4"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "request_data_size", "64"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "type_of_service", "18"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "operation_timeout", "100"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "threshold", "100"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "ipv6_traffic_class", "255"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update IPSLA Monitoring Policy with ICMP") },
				Config:    testAccMSOTenantPoliciesIPSLAMonitoringPolicyConfigUpdateSlaTypeWithICMP(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "name", "test_ipsla_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "description", "ICMP Type"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "sla_type", "icmp"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "destination_port", "0"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "http_version", "HTTP11"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "http_uri", "/example"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "sla_frequency", "60"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "detect_multiplier", "4"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "request_data_size", "64"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "type_of_service", "18"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "operation_timeout", "200"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "threshold", "200"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "ipv6_traffic_class", "0"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update IPSLA Monitoring Policy with L2Ping") },
				Config:    testAccMSOTenantPoliciesIPSLAMonitoringPolicyConfigUpdateSlaTypeWithL2Ping(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "name", "test_ipsla_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "description", "L2Ping Type"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "sla_type", "l2ping"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "destination_port", "0"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "http_version", "HTTP11"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "http_uri", "/example"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "sla_frequency", "100"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "detect_multiplier", "3"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "request_data_size", "64"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "type_of_service", "18"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "operation_timeout", "150"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "threshold", "150"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "ipv6_traffic_class", "100"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update IPSLA Monitoring Policy with TCP") },
				Config:    testAccMSOTenantPoliciesIPSLAMonitoringPolicyConfigUpdateSlaTypeWithTCP(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "name", "test_ipsla_policy_tcp"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "description", "TCP Type"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "sla_type", "tcp"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "destination_port", "100"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "http_version", "HTTP11"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "http_uri", "/example"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "sla_frequency", "100"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "detect_multiplier", "3"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "request_data_size", "64"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "type_of_service", "18"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "operation_timeout", "150"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "threshold", "150"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy", "ipv6_traffic_class", "100"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import IPSLA Monitoring Policy") },
				ResourceName:      "mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithArguments("mso_tenant_policies_ipsla_monitoring_policy", "ipslaMonitoringPolicy"),
	})
}

func testAccMSOTenantPoliciesIPSLAMonitoringPolicyConfigCreate() string {
	return fmt.Sprintf(`%s
	resource "mso_tenant_policies_ipsla_monitoring_policy" "ipsla_policy" {
	    template_id        = mso_template.template_tenant.id
	    name               = "test_ipsla_policy"
	    description        = "HTTP Type"
	    sla_type           = "http"
	    destination_port   = 80
	    http_version       = "HTTP11"
	    http_uri           = "/example"
	    sla_frequency      = 120
	    detect_multiplier  = 4
	    request_data_size  = 64
	    type_of_service    = 18
	    operation_timeout  = 100
	    threshold          = 100
	    ipv6_traffic_class = 255
	}`, testAccMSOTemplateResourceTenantConfig())

}

func testAccMSOTenantPoliciesIPSLAMonitoringPolicyConfigUpdateSlaTypeWithICMP() string {
	return fmt.Sprintf(`%s
	resource "mso_tenant_policies_ipsla_monitoring_policy" "ipsla_policy" {
	    template_id        = mso_template.template_tenant.id
	    name               = "test_ipsla_policy"
	    description        = "ICMP Type"
	    sla_type           = "icmp"
	    operation_timeout  = 200
	    threshold          = 200
	    ipv6_traffic_class = 0
	    destination_port   = 80
	    detect_multiplier  = 4
	    sla_frequency      = 60
	    request_data_size  = 64
	}`, testAccMSOTemplateResourceTenantConfig())

}

func testAccMSOTenantPoliciesIPSLAMonitoringPolicyConfigUpdateSlaTypeWithL2Ping() string {
	return fmt.Sprintf(`%s
	resource "mso_tenant_policies_ipsla_monitoring_policy" "ipsla_policy" {
	    template_id        = mso_template.template_tenant.id
	    name               = "test_ipsla_policy"
	    description        = "L2Ping Type"
	    sla_type           = "l2ping"
	    detect_multiplier  = 3
	    sla_frequency      = 100
	    operation_timeout  = 150
	    threshold          = 150
	    ipv6_traffic_class = 100
	}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesIPSLAMonitoringPolicyConfigUpdateSlaTypeWithTCP() string {
	return fmt.Sprintf(`%s
	resource "mso_tenant_policies_ipsla_monitoring_policy" "ipsla_policy" {
	    template_id        = mso_template.template_tenant.id
	    name               = "test_ipsla_policy_tcp"
	    description        = "TCP Type"
	    sla_type           = "tcp"
	    destination_port   = 100
	    detect_multiplier  = 3
	    sla_frequency      = 100
	    operation_timeout  = 150
	    threshold          = 150
	    ipv6_traffic_class = 100
	}`, testAccMSOTemplateResourceTenantConfig())
}
