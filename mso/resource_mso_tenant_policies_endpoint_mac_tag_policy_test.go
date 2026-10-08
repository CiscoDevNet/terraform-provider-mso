package mso

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMSOTenantPoliciesEndpointMACTagPolicyResource(t *testing.T) {
	resourceName := "mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_bd"
	var templateID, uuid string

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccVersionCheck(t, "5.1") },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Error when neither bd_uuid nor vrf_uuid is provided") },
				Config:      testAccMSOTenantPoliciesEndpointMACTagPolicyConfigErrorMissingScope(),
				ExpectError: regexp.MustCompile(`BdRef and VrfRef cannot both be empty in Endpoint Mac Tag Policy AA:BB:A1:B2:C3:D4. Either BdRef or VrfRef must be populated.`),
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
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok || rs.Primary == nil {
							return fmt.Errorf("resource %s not found in state", resourceName)
						}
						templateID = rs.Primary.Attributes["template_id"]
						uuid = rs.Primary.Attributes["uuid"]
						if templateID == "" || uuid == "" {
							return fmt.Errorf("resource %s is missing its template ID or UUID", resourceName)
						}
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate Endpoint MAC Tag Policy after out-of-band deletion")
					if err := testAccDeletePolicyOutOfBand(testAccPreCheck(t), templateID, uuid, "tenantPolicyTemplate", "template", "endpointMacTagPolicies"); err != nil {
						t.Fatalf("delete %s out of band: %v", resourceName, err)
					}
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
			{
				PreConfig:    func() { fmt.Println("Test: Import missing Endpoint MAC Tag Policy") },
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok || rs.Primary == nil {
						return "", fmt.Errorf("resource %s not found in state", resourceName)
					}
					separator := strings.LastIndex(rs.Primary.ID, "/")
					if separator < 0 || separator == len(rs.Primary.ID)-1 {
						return "", fmt.Errorf("resource %s has an invalid import ID", resourceName)
					}
					return rs.Primary.ID[:separator+1] + "tf_missing_oob", nil
				},
				ExpectError: regexp.MustCompile(`cannot import .*: resource not found`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Create Endpoint MAC Tag Policy with duplicate object") },
				Config:      testAccMSOTenantPoliciesEndpointMACTagPolicyConfigWithDuplicateObject(),
				ExpectError: regexp.MustCompile(regexp.QuoteMeta(fmt.Sprintf("Multiple endpointMacTag policies are using the name: AA:BB:A1:B2:C3:D4-[%s]", msoSchemaTemplateBdName))),
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments("mso_tenant_policies_endpoint_mac_tag_policy", "tenantPolicyTemplate", "template", "endpointMacTagPolicies"),
	})
}

var endpointMACTagPolicyPreConfig = testSiteConfigAnsibleTest() + testTenantConfig() + testSchemaConfig() + testSchemaTemplateVrfConfig() + testSchemaTemplateBdConfig() + testTenantPolicyTemplateConfig()

func testAccMSOTenantPoliciesEndpointMACTagPolicyConfigErrorMissingScope() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_tenant_policies_endpoint_mac_tag_policy" "error" {
		template_id = mso_template.%[2]s.id
		mac         = "AA:BB:A1:B2:C3:D4"
	}`, endpointMACTagPolicyPreConfig, msoTenantPolicyTemplateName)
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

func testAccMSOTenantPoliciesEndpointMACTagPolicyConfigWithDuplicateObject() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_tenant_policies_endpoint_mac_tag_policy" "endpoint_mac_bd_duplicate" {
		template_id = mso_template.%[2]s.id
		bd_uuid     = mso_schema_template_bd.%[3]s.uuid
		mac         = "AA:BB:A1:B2:C3:D4"
	}`,
		endpointMACTagPolicyPreConfig,
		msoTenantPolicyTemplateName,
		msoSchemaTemplateBdName,
	)
}
