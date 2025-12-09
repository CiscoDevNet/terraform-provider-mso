package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOTenantPoliciesCustomQoSPolicyDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Custom QoS Policy Data Source - Basic") },
				Config:    testAccMSOTenantPoliciesCustomQoSPolicyDataSourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_custom_qos_policy.qos_policy", "name", "test_custom_qos_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_custom_qos_policy.qos_policy", "description", "Test Custom QoS Policy"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_custom_qos_policy.qos_policy", "uuid"),

					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_custom_qos_policy.qos_policy", "dscp_mappings", map[string]string{
						"dscp_from":    "af11",
						"dscp_to":      "af12",
						"dscp_target":  "af11",
						"target_cos":   "background",
						"qos_priority": "level1",
					}),
					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_custom_qos_policy.qos_policy", "dscp_mappings", map[string]string{
						"dscp_from":    "af21",
						"dscp_to":      "af22",
						"dscp_target":  "af21",
						"target_cos":   "best_effort",
						"qos_priority": "level2",
					}),

					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_custom_qos_policy.qos_policy", "cos_mappings", map[string]string{
						"dot1p_from":   "background",
						"dot1p_to":     "best_effort",
						"dscp_target":  "af11",
						"target_cos":   "background",
						"qos_priority": "level1",
					}),
					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_custom_qos_policy.qos_policy", "cos_mappings", map[string]string{
						"dot1p_from":   "excellent_effort",
						"dot1p_to":     "critical_applications",
						"dscp_target":  "af21",
						"target_cos":   "video",
						"qos_priority": "level2",
					}),
				),
			},
		},
	})
}

func testAccMSOTenantPoliciesCustomQoSPolicyDataSourceBasic() string {
	return fmt.Sprintf(`%s
    data "mso_tenant_policies_custom_qos_policy" "qos_policy" {
        template_id = mso_template.template_tenant.id
        name        = mso_tenant_policies_custom_qos_policy.qos_policy.name
    }`, testAccMSOTenantPoliciesCustomQoSPolicyConfigCreate())
}
