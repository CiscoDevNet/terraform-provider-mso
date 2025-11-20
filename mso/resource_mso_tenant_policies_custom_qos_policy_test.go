package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOTenantPoliciesCustomQoSPolicyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create Custom QoS Policy") },
				Config:    testAccMSOTenantPoliciesCustomQoSPolicyConfigCreate(),
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
			{
				PreConfig: func() { fmt.Println("Test: Update Custom QoS Policy - Add More Mappings") },
				Config:    testAccMSOTenantPoliciesCustomQoSPolicyConfigUpdateAddMappings(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_custom_qos_policy.qos_policy", "name", "test_custom_qos_policy"),
					resource.TestCheckResourceAttr("mso_tenant_policies_custom_qos_policy.qos_policy", "description", "Updated Custom QoS Policy"),

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
					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_custom_qos_policy.qos_policy", "dscp_mappings", map[string]string{
						"dscp_from":    "af31",
						"dscp_to":      "af32",
						"dscp_target":  "af31",
						"target_cos":   "voice",
						"qos_priority": "level3",
					}),

					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_custom_qos_policy.qos_policy", "cos_mappings", map[string]string{
						"dot1p_from":   "voice",
						"dot1p_to":     "internetwork_control",
						"dscp_target":  "cs5",
						"target_cos":   "voice",
						"qos_priority": "level5",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Custom QoS Policy - Remove Some Mappings") },
				Config:    testAccMSOTenantPoliciesCustomQoSPolicyConfigUpdateRemoveMappings(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_custom_qos_policy.qos_policy", "name", "test_custom_qos_policy"),

					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_custom_qos_policy.qos_policy", "dscp_mappings", map[string]string{
						"dscp_from":    "af11",
						"dscp_to":      "af12",
						"dscp_target":  "af11",
						"target_cos":   "background",
						"qos_priority": "level1",
					}),

					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_custom_qos_policy.qos_policy", "cos_mappings", map[string]string{
						"dot1p_from":   "background",
						"dot1p_to":     "best_effort",
						"dscp_target":  "af11",
						"target_cos":   "background",
						"qos_priority": "level1",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Custom QoS Policy - Change Mapping Values") },
				Config:    testAccMSOTenantPoliciesCustomQoSPolicyConfigUpdateChangeMappings(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_custom_qos_policy.qos_policy", "name", "test_custom_qos_policy"),

					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_custom_qos_policy.qos_policy", "dscp_mappings", map[string]string{
						"dscp_from":    "cs0",
						"dscp_to":      "cs1",
						"dscp_target":  "cs0",
						"target_cos":   "network_control",
						"qos_priority": "level6",
					}),

					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_custom_qos_policy.qos_policy", "cos_mappings", map[string]string{
						"dot1p_from":   "network_control",
						"dot1p_to":     "unspecified",
						"dscp_target":  "expedited_forwarding",
						"target_cos":   "network_control",
						"qos_priority": "level6",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Custom QoS Policy with Default/Unspecified Values") },
				Config:    testAccMSOTenantPoliciesCustomQoSPolicyConfigUpdateWithDefaults(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_custom_qos_policy.qos_policy", "name", "test_custom_qos_policy"),

					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_custom_qos_policy.qos_policy", "dscp_mappings", map[string]string{
						"dscp_from":    "af11",
						"dscp_to":      "unspecified",
						"dscp_target":  "unspecified",
						"target_cos":   "unspecified",
						"qos_priority": "unspecified",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Custom QoS Policy Name") },
				Config:    testAccMSOTenantPoliciesCustomQoSPolicyConfigUpdateName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_custom_qos_policy.qos_policy", "name", "test_custom_qos_policy_renamed"),
					resource.TestCheckResourceAttr("mso_tenant_policies_custom_qos_policy.qos_policy", "description", "Renamed Policy"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Custom QoS Policy with Maximum Mappings") },
				Config:    testAccMSOTenantPoliciesCustomQoSPolicyConfigMaxMappings(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_custom_qos_policy.qos_policy", "name", "test_custom_qos_policy_renamed"),

					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_custom_qos_policy.qos_policy", "dscp_mappings", map[string]string{
						"dscp_from":    "af11",
						"dscp_to":      "af13",
						"qos_priority": "level1",
					}),
					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_custom_qos_policy.qos_policy", "dscp_mappings", map[string]string{
						"dscp_from":    "af21",
						"dscp_to":      "af23",
						"qos_priority": "level2",
					}),
					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_custom_qos_policy.qos_policy", "dscp_mappings", map[string]string{
						"dscp_from":    "af31",
						"dscp_to":      "af33",
						"qos_priority": "level3",
					}),
					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_custom_qos_policy.qos_policy", "dscp_mappings", map[string]string{
						"dscp_from":    "af41",
						"dscp_to":      "af43",
						"qos_priority": "level4",
					}),
					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_custom_qos_policy.qos_policy", "cos_mappings", map[string]string{
						"dot1p_from":   "background",
						"dot1p_to":     "best_effort",
						"qos_priority": "level1",
					}),
					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_custom_qos_policy.qos_policy", "cos_mappings", map[string]string{
						"dot1p_from":   "excellent_effort",
						"dot1p_to":     "critical_applications",
						"qos_priority": "level2",
					}),
					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_custom_qos_policy.qos_policy", "cos_mappings", map[string]string{
						"dot1p_from":   "video",
						"dot1p_to":     "voice",
						"qos_priority": "level3",
					}),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Custom QoS Policy") },
				ResourceName:      "mso_tenant_policies_custom_qos_policy.qos_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithArguments("mso_tenant_policies_custom_qos_policy", "qos"),
	})
}

func testAccMSOTenantPoliciesCustomQoSPolicyConfigCreate() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_custom_qos_policy" "qos_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_custom_qos_policy"
        description = "Test Custom QoS Policy"
        
        dscp_mappings {
            dscp_from    = "af11"
            dscp_to      = "af12"
            dscp_target  = "af11"
            target_cos   = "background"
            qos_priority = "level1"
        }
        
        dscp_mappings {
            dscp_from    = "af21"
            dscp_to      = "af22"
            dscp_target  = "af21"
            target_cos   = "best_effort"
            qos_priority = "level2"
        }
        
        cos_mappings {
            dot1p_from   = "background"
            dot1p_to     = "best_effort"
            dscp_target  = "af11"
            target_cos   = "background"
            qos_priority = "level1"
        }
        
        cos_mappings {
            dot1p_from   = "excellent_effort"
            dot1p_to     = "critical_applications"
            dscp_target  = "af21"
            target_cos   = "video"
            qos_priority = "level2"
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesCustomQoSPolicyConfigUpdateAddMappings() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_custom_qos_policy" "qos_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_custom_qos_policy"
        description = "Updated Custom QoS Policy"
        
        dscp_mappings {
            dscp_from    = "af11"
            dscp_to      = "af12"
            dscp_target  = "af11"
            target_cos   = "background"
            qos_priority = "level1"
        }
        
        dscp_mappings {
            dscp_from    = "af21"
            dscp_to      = "af22"
            dscp_target  = "af21"
            target_cos   = "best_effort"
            qos_priority = "level2"
        }
        
        dscp_mappings {
            dscp_from    = "af31"
            dscp_to      = "af32"
            dscp_target  = "af31"
            target_cos   = "voice"
            qos_priority = "level3"
        }
        
        cos_mappings {
            dot1p_from   = "voice"
            dot1p_to     = "internetwork_control"
            dscp_target  = "cs5"
            target_cos   = "voice"
            qos_priority = "level5"
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesCustomQoSPolicyConfigUpdateRemoveMappings() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_custom_qos_policy" "qos_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_custom_qos_policy"
        description = "Updated Custom QoS Policy"
        
        dscp_mappings {
            dscp_from    = "af11"
            dscp_to      = "af12"
            dscp_target  = "af11"
            target_cos   = "background"
            qos_priority = "level1"
        }
        
        cos_mappings {
            dot1p_from   = "background"
            dot1p_to     = "best_effort"
            dscp_target  = "af11"
            target_cos   = "background"
            qos_priority = "level1"
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesCustomQoSPolicyConfigUpdateChangeMappings() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_custom_qos_policy" "qos_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_custom_qos_policy"
        description = "Updated Custom QoS Policy"
        
        dscp_mappings {
            dscp_from    = "cs0"
            dscp_to      = "cs1"
            dscp_target  = "cs0"
            target_cos   = "network_control"
            qos_priority = "level6"
        }
        
        cos_mappings {
            dot1p_from   = "network_control"
            dot1p_to     = "unspecified"
            dscp_target  = "expedited_forwarding"
            target_cos   = "network_control"
            qos_priority = "level6"
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesCustomQoSPolicyConfigUpdateWithDefaults() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_custom_qos_policy" "qos_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_custom_qos_policy"
        description = "Testing Default Values"
        
        dscp_mappings {
            dscp_from    = "af11"
            dscp_to      = "unspecified"
            dscp_target  = "unspecified"
            target_cos   = "unspecified"
            qos_priority = "unspecified"
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesCustomQoSPolicyConfigUpdateName() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_custom_qos_policy" "qos_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_custom_qos_policy_renamed"
        description = "Renamed Policy"
        
        dscp_mappings {
            dscp_from    = "af11"
            dscp_to      = "af12"
            dscp_target  = "af11"
            target_cos   = "background"
            qos_priority = "level1"
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesCustomQoSPolicyConfigMaxMappings() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_custom_qos_policy" "qos_policy" {
        template_id = mso_template.template_tenant.id
        name        = "test_custom_qos_policy_renamed"
        description = "Maximum Mappings Test"
        
        dscp_mappings {
            dscp_from    = "af11"
            dscp_to      = "af13"
            qos_priority = "level1"
        }
        
        dscp_mappings {
            dscp_from    = "af21"
            dscp_to      = "af23"
            qos_priority = "level2"
        }
        
        dscp_mappings {
            dscp_from    = "af31"
            dscp_to      = "af33"
            qos_priority = "level3"
        }
        
        dscp_mappings {
            dscp_from    = "af41"
            dscp_to      = "af43"
            qos_priority = "level4"
        }
        
        cos_mappings {
            dot1p_from   = "background"
            dot1p_to     = "best_effort"
            qos_priority = "level1"
        }
        
        cos_mappings {
            dot1p_from   = "excellent_effort"
            dot1p_to     = "critical_applications"
            qos_priority = "level2"
        }
        
        cos_mappings {
            dot1p_from   = "video"
            dot1p_to     = "voice"
            qos_priority = "level3"
        }
    }`, testAccMSOTemplateResourceTenantConfig())
}
