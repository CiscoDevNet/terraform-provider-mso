package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOTenantPoliciesEndpointMACTagPolicyDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccVersionCheck(t, "5.1") },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Error when neither bd_uuid nor vrf_uuid is provided Data Source") },
				Config:      testAccMSOTenantPoliciesEndpointMACTagPolicyConfigErrorMissingScopeDataSource(),
				ExpectError: regexp.MustCompile(`Either 'bd_uuid' or 'vrf_uuid' must be specified to use Endpoint MAC Tag Policy Data Source`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Error when both bd_uuid and vrf_uuid are provided Data Source") },
				Config:      testAccMSOTenantPoliciesEndpointMACTagPolicyConfigErrorConflictingScopeDataSource(),
				ExpectError: regexp.MustCompile(`conflicts with`),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Endpoint MAC Tag Policy (BD scope) with multiple annotations and tags Data Source")
				},
				Config: testAccMSOTenantPoliciesEndpointMACTagPolicyConfigCreateBDWithMultipleTagsDataSource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "mac", "AA:BB:A1:B2:C3:D4"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "name", fmt.Sprintf("AA:BB:A1:B2:C3:D4-[%s]", msoSchemaTemplateBdName)),
					resource.TestCheckResourceAttrSet("data.mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "uuid"),
					resource.TestCheckResourceAttrSet("data.mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "bd_uuid"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "tag_annotations.#", "2"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "policy_tags.#", "2"),
					CustomTestCheckTypeSetElemAttrs("data.mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "tag_annotations",
						map[string]string{
							"key":   "annotation_key_1",
							"value": "annotation_value_1",
						},
					),
					CustomTestCheckTypeSetElemAttrs("data.mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "tag_annotations",
						map[string]string{
							"key":   "annotation_key_2",
							"value": "annotation_value_2",
						},
					),
					CustomTestCheckTypeSetElemAttrs("data.mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "policy_tags",
						map[string]string{
							"key":   "policy_key_1",
							"value": "policy_value_1",
						},
					),
					CustomTestCheckTypeSetElemAttrs("data.mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "policy_tags",
						map[string]string{
							"key":   "policy_key_2",
							"value": "policy_value_2",
						},
					),
				),
			},
		},
	})
}

func testAccMSOTenantPoliciesEndpointMACTagPolicyConfigErrorMissingScopeDataSource() string {
	return `
	data "mso_tenant_policies_endpoint_mac_tag_policy" "error" {
		template_id = "mso_template.error.id"
		mac         = "AA:BB:A1:B2:C3:D4"
	}`
}

func testAccMSOTenantPoliciesEndpointMACTagPolicyConfigErrorConflictingScopeDataSource() string {
	return `
	data "mso_tenant_policies_endpoint_mac_tag_policy" "error" {
		template_id = "mso_template.error.id"
		mac         = "AA:BB:A1:B2:C3:D4"
		bd_uuid     = "mso_schema_template_bd.error.uuid"
		vrf_uuid    = "mso_schema_template_vrf.error.uuid"
	}`
}

func testAccMSOTenantPoliciesEndpointMACTagPolicyConfigCreateBDWithMultipleTagsDataSource() string {
	return fmt.Sprintf(`%[1]s
	data "mso_tenant_policies_endpoint_mac_tag_policy" "endpoint_mac_bd" {
		template_id = mso_template.%[2]s.id
		mac         = "AA:BB:A1:B2:C3:D4"
		bd_uuid     = mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd.bd_uuid
	}`,
		testAccMSOTenantPoliciesEndpointMACTagPolicyConfigCreateBDWithMultipleTags(),
		msoTenantPolicyTemplateName,
	)
}
