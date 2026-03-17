package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOTenantPoliciesEndpointMACTagPolicyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Error when neither bd_uuid nor vrf_uuid is provided") },
				Config:      testAccMSOTenantPoliciesEndpointMACTagPolicyConfigErrorMissingScope(),
				ExpectError: regexp.MustCompile(`Either 'bd_uuid' or 'vrf_uuid' must be specified for creating Endpoint MAC Tag Policy`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Error when both bd_uuid and vrf_uuid are provided") },
				Config:      testAccMSOTenantPoliciesEndpointMACTagPolicyConfigErrorConflictingScope(),
				ExpectError: regexp.MustCompile(`conflicts with`),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Create Endpoint MAC Tag Policy (BD scope) with multiple annotations and tags")
				},
				Config: testAccMSOTenantPoliciesEndpointMACTagPolicyConfigCreateBDWithMultipleTags(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "mac", "AA:BB:A1:B2:C3:D4"),
					resource.TestCheckResourceAttr("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "name", fmt.Sprintf("AA:BB:A1:B2:C3:D4-[%s]", msoSchemaTemplateBdName)),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "uuid"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "bd_uuid"),
					resource.TestCheckResourceAttr("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "tag_annotations.#", "2"),
					resource.TestCheckResourceAttr("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "policy_tags.#", "2"),
					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "tag_annotations",
						map[string]string{
							"key":   "annotation_key_1",
							"value": "annotation_value_1",
						},
					),
					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "tag_annotations",
						map[string]string{
							"key":   "annotation_key_2",
							"value": "annotation_value_2",
						},
					),
					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "policy_tags",
						map[string]string{
							"key":   "policy_key_1",
							"value": "policy_value_1",
						},
					),
					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "policy_tags",
						map[string]string{
							"key":   "policy_key_2",
							"value": "policy_value_2",
						},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Endpoint MAC Tag Policy from BD scope to VRF scope") },
				Config:    testAccMSOTenantPoliciesEndpointMACTagPolicyConfigUpdateBDToVRF(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "mac", "AA:BB:A1:B2:C3:D4"),
					resource.TestCheckResourceAttr("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "name", "AA:BB:A1:B2:C3:D4-[*]"),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "vrf_uuid"),
					resource.TestCheckResourceAttr("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "tag_annotations.#", "2"),
					resource.TestCheckResourceAttr("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "policy_tags.#", "2"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Revert Endpoint MAC Tag Policy from VRF scope back to BD scope") },
				Config:    testAccMSOTenantPoliciesEndpointMACTagPolicyConfigUpdateVRFToBD(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "mac", "AA:BB:A1:B2:C3:D4"),
					resource.TestCheckResourceAttr("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "name", fmt.Sprintf("AA:BB:A1:B2:C3:D4-[%s]", msoSchemaTemplateBdName)),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "bd_uuid"),
					resource.TestCheckResourceAttr("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "tag_annotations.#", "2"),
					resource.TestCheckResourceAttr("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "policy_tags.#", "2"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Endpoint MAC Tag Policy removing one annotation and one tag") },
				Config:    testAccMSOTenantPoliciesEndpointMACTagPolicyConfigUpdateRemoveOneTagAndAnnotation(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "mac", "AA:BB:A1:B2:C3:D4"),
					resource.TestCheckResourceAttr("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "name", fmt.Sprintf("AA:BB:A1:B2:C3:D4-[%s]", msoSchemaTemplateBdName)),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "bd_uuid"),
					resource.TestCheckResourceAttr("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "tag_annotations.#", "1"),
					resource.TestCheckResourceAttr("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "policy_tags.#", "1"),
					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "tag_annotations",
						map[string]string{
							"key":   "annotation_key_2",
							"value": "annotation_value_2",
						},
					),
					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "policy_tags",
						map[string]string{
							"key":   "policy_key_2",
							"value": "policy_value_2",
						},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Endpoint MAC Tag Policy removing all annotations and tags") },
				Config:    testAccMSOTenantPoliciesEndpointMACTagPolicyConfigUpdateRemoveAllTagsAndAnnotations(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "mac", "AA:BB:A1:B2:C3:D4"),
					resource.TestCheckResourceAttr("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "name", fmt.Sprintf("AA:BB:A1:B2:C3:D4-[%s]", msoSchemaTemplateBdName)),
					resource.TestCheckResourceAttrSet("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "bd_uuid"),
					resource.TestCheckResourceAttr("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "tag_annotations.#", "0"),
					resource.TestCheckResourceAttr("mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd", "policy_tags.#", "0"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Endpoint MAC Tag Policy") },
				ResourceName:      "mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments("mso_tenant_policies_endpoint_mac_tag_policy", "tenantPolicyTemplate", "template", "endpointMacTagPolicies"),
	})
}

var endpointMACTagPolicyPreConfig = testSiteConfigAnsibleTest() + testTenantConfig() + testSchemaConfig() + testSchemaTemplateVrfConfig() + testSchemaTemplateBdConfig() + testTenantPolicyTemplateConfig()

func testAccMSOTenantPoliciesEndpointMACTagPolicyConfigErrorMissingScope() string {
	return `
	resource "mso_tenant_policies_endpoint_mac_tag_policy" "error" {
		template_id = "mso_template.error.id"
		mac         = "AA:BB:A1:B2:C3:D4"
	}`
}

func testAccMSOTenantPoliciesEndpointMACTagPolicyConfigErrorConflictingScope() string {
	return `
	resource "mso_tenant_policies_endpoint_mac_tag_policy" "error" {
		template_id = "mso_template.error.id"
		mac         = "AA:BB:A1:B2:C3:D4"
		bd_uuid     = "mso_schema_template_bd.error.uuid"
		vrf_uuid    = "mso_schema_template_vrf.error.uuid"
	}`
}

func testAccMSOTenantPoliciesEndpointMACTagPolicyConfigCreateBDWithMultipleTags() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_tenant_policies_endpoint_mac_tag_policy" "endpoint_mac_bd" {
		template_id = mso_template.%[2]s.id
		mac         = "AA:BB:A1:B2:C3:D4"
		bd_uuid     = mso_schema_template_bd.%[3]s.uuid

		tag_annotations {
			key   = "annotation_key_1"
			value = "annotation_value_1"
		}

		tag_annotations {
			key   = "annotation_key_2"
			value = "annotation_value_2"
		}

		policy_tags {
			key   = "policy_key_1"
			value = "policy_value_1"
		}

		policy_tags {
			key   = "policy_key_2"
			value = "policy_value_2"
		}
	}`,
		endpointMACTagPolicyPreConfig,
		msoTenantPolicyTemplateName,
		msoSchemaTemplateBdName,
	)
}

func testAccMSOTenantPoliciesEndpointMACTagPolicyConfigUpdateBDToVRF() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_tenant_policies_endpoint_mac_tag_policy" "endpoint_mac_bd" {
		template_id = mso_template.%[2]s.id
		mac         = "AA:BB:A1:B2:C3:D4"
		vrf_uuid    = mso_schema_template_vrf.%[3]s.uuid

		tag_annotations {
			key   = "annotation_key_1"
			value = "annotation_value_1"
		}

		tag_annotations {
			key   = "annotation_key_2"
			value = "annotation_value_2"
		}

		policy_tags {
			key   = "policy_key_1"
			value = "policy_value_1"
		}

		policy_tags {
			key   = "policy_key_2"
			value = "policy_value_2"
		}
	}`,
		endpointMACTagPolicyPreConfig,
		msoTenantPolicyTemplateName,
		msoSchemaTemplateVrfName,
	)
}

func testAccMSOTenantPoliciesEndpointMACTagPolicyConfigUpdateVRFToBD() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_tenant_policies_endpoint_mac_tag_policy" "endpoint_mac_bd" {
		template_id = mso_template.%[2]s.id
		mac         = "AA:BB:A1:B2:C3:D4"
		bd_uuid     = mso_schema_template_bd.%[3]s.uuid

		tag_annotations {
			key   = "annotation_key_1"
			value = "annotation_value_1"
		}

		tag_annotations {
			key   = "annotation_key_2"
			value = "annotation_value_2"
		}

		policy_tags {
			key   = "policy_key_1"
			value = "policy_value_1"
		}

		policy_tags {
			key   = "policy_key_2"
			value = "policy_value_2"
		}
	}`,
		endpointMACTagPolicyPreConfig,
		msoTenantPolicyTemplateName,
		msoSchemaTemplateBdName,
	)
}

func testAccMSOTenantPoliciesEndpointMACTagPolicyConfigUpdateRemoveOneTagAndAnnotation() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_tenant_policies_endpoint_mac_tag_policy" "endpoint_mac_bd" {
		template_id = mso_template.%[2]s.id
		bd_uuid     = mso_schema_template_bd.%[3]s.uuid
		mac         = "AA:BB:A1:B2:C3:D4"

		tag_annotations {
			key   = "annotation_key_2"
			value = "annotation_value_2"
		}

		policy_tags {
			key   = "policy_key_2"
			value = "policy_value_2"
		}
	}`,
		endpointMACTagPolicyPreConfig,
		msoTenantPolicyTemplateName,
		msoSchemaTemplateBdName,
	)
}

func testAccMSOTenantPoliciesEndpointMACTagPolicyConfigUpdateRemoveAllTagsAndAnnotations() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_tenant_policies_endpoint_mac_tag_policy" "endpoint_mac_bd" {
		template_id = mso_template.%[2]s.id
		bd_uuid     = mso_schema_template_bd.%[3]s.uuid
		mac         = "AA:BB:A1:B2:C3:D4"
	}`,
		endpointMACTagPolicyPreConfig,
		msoTenantPolicyTemplateName,
		msoSchemaTemplateBdName,
	)
}
