package provider

import (
	"fmt"
	"maps"
	"regexp"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-mso/internal/models"
	"github.com/CiscoDevNet/terraform-provider-mso/internal/ndoapi"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccMSOL3OutInterfaceGroupPolicyResource(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_interface_group")
	var templateID, l3outUUID string
	var keyRefs map[string]string
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccProviderPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig:        func() { fmt.Println("Test: Create L3Out interface group policy") },
				Config:           testAccMSOL3OutInterfaceGroupPolicyCreateConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCaptureResourceIdentifiers("mso_l3out.test", &templateID, &l3outUUID),
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "name", "test_interface_group"),
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "description", "initial"),
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "bfd.enabled", "false"),
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "bfd_multi_hop.enabled", "false"),
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "ospf.enabled", "false"),
				),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject a second interface group policy with the same name") },
				Config:      testAccMSOL3OutInterfaceGroupPolicyDuplicateConfig(siteName, l3outName),
				ExpectError: regexp.MustCompile(`L3Out Interface Group Policy Already Exists`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create a second independently managed interface group policy") },
				Config:    testAccMSOL3OutInterfaceGroupPolicyTwoGroupsConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "name", "test_interface_group"),
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.second", "name", "test_interface_group_second"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove the second interface group without removing the first") },
				Config:    testAccMSOL3OutInterfaceGroupPolicyCreateConfig(siteName, l3outName),
				Check:     resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "name", "test_interface_group"),
			},
			{
				PreConfig:        func() { fmt.Println("Test: Enable group BFD and OSPF without authentication") },
				Config:           testAccMSOL3OutInterfaceGroupPolicyEnabledConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "bfd.enabled", "true"),
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "bfd_multi_hop.enabled", "true"),
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "ospf.enabled", "true"),
				),
			},
			{
				PreConfig:        func() { fmt.Println("Test: Enable BFD and OSPF authentication with separate keys") },
				Config:           testAccMSOL3OutInterfaceGroupPolicyAuthenticatedConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "bfd.authentication_enabled", "true"),
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "bfd.key_id", "20"),
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "bfd_multi_hop.authentication_enabled", "true"),
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "bfd_multi_hop.key_id", "30"),
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "ospf.authentication_type", "md5"),
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "ospf.key_id", "10"),
					testAccCheckL3OutInterfaceGroupKeyRefs(&templateID, &l3outUUID),
				),
			},
			{
				PreConfig:        func() { fmt.Println("Test: Change BFD key ID without changing its key") },
				Config:           testAccMSOL3OutInterfaceGroupPolicyAuthenticationConfig(siteName, l3outName, 21, 30, 10, "md5"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "bfd.key_id", "21"),
					testAccCheckL3OutInterfaceGroupKeyRefs(&templateID, &l3outUUID),
				),
			},
			{
				PreConfig:        func() { fmt.Println("Test: Change multi-hop BFD key ID without changing its key") },
				Config:           testAccMSOL3OutInterfaceGroupPolicyAuthenticationConfig(siteName, l3outName, 21, 31, 10, "md5"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "bfd_multi_hop.key_id", "31"),
					testAccCheckL3OutInterfaceGroupKeyRefs(&templateID, &l3outUUID),
				),
			},
			{
				PreConfig:        func() { fmt.Println("Test: Change OSPF key ID without changing its key") },
				Config:           testAccMSOL3OutInterfaceGroupPolicyAuthenticationConfig(siteName, l3outName, 21, 31, 11, "md5"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "ospf.key_id", "11"),
					testAccCheckL3OutInterfaceGroupKeyRefs(&templateID, &l3outUUID),
				),
			},
			{
				PreConfig:        func() { fmt.Println("Test: Change OSPF authentication mode without changing its key") },
				Config:           testAccMSOL3OutInterfaceGroupPolicyAuthenticationConfig(siteName, l3outName, 21, 31, 11, "simple"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "ospf.authentication_type", "simple"),
					testAccCheckL3OutInterfaceGroupKeyRefs(&templateID, &l3outUUID),
					testAccCaptureL3OutInterfaceGroupKeyRefs(&templateID, &l3outUUID, &keyRefs),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import authenticated interface group without recovering its keys") },
				ResourceName: "mso_l3out_interface_group_policy.test",
				ImportState:  true,
				ImportStateIdFunc: func(_ *terraform.State) (string, error) {
					return templateID + "/" + l3outUUID + "/test_interface_group", nil
				},
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected one imported interface group policy, got %d", len(states))
					}
					attributes := states[0].Attributes
					for name, expected := range map[string]string{
						"bfd.authentication_enabled":           "true",
						"bfd.key_id":                           "21",
						"bfd_multi_hop.authentication_enabled": "true",
						"bfd_multi_hop.key_id":                 "31",
						"ospf.authentication_type":             "simple",
						"ospf.key_id":                          "11",
					} {
						if attributes[name] != expected {
							return fmt.Errorf("imported %s = %q, want %q", name, attributes[name], expected)
						}
					}
					for _, name := range []string{"bfd.key", "bfd_multi_hop.key", "ospf.key"} {
						if value, present := attributes[name]; present {
							return fmt.Errorf("imported %s unexpectedly contains %q", name, value)
						}
					}
					return nil
				},
			},
			{
				PreConfig: func() { fmt.Println("Test: Update description while keeping configured authentication") },
				Config:    testAccMSOL3OutInterfaceGroupPolicyAuthenticatedDescriptionUpdateConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "description", "authentication unchanged"),
					testAccCheckL3OutInterfaceGroupKeyRefsUnchanged(&templateID, &l3outUUID, &keyRefs),
				),
			},
			{
				PreConfig:        func() { fmt.Println("Test: Omit configured protocol objects during an unrelated update") },
				Config:           testAccMSOL3OutInterfaceGroupPolicyAuthenticatedOmittedConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "bfd.authentication_enabled", "true"),
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "bfd_multi_hop.authentication_enabled", "true"),
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "ospf.authentication_type", "simple"),
					testAccCheckL3OutInterfaceGroupKeyRefsUnchanged(&templateID, &l3outUUID, &keyRefs),
				),
			},
			{
				PreConfig:        func() { fmt.Println("Test: Disable authentication while leaving the group protocols enabled") },
				Config:           testAccMSOL3OutInterfaceGroupPolicyAuthenticationDisabledConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "bfd.enabled", "true"),
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "bfd.authentication_enabled", "false"),
					resource.TestCheckNoResourceAttr("mso_l3out_interface_group_policy.test", "bfd.key"),
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "bfd_multi_hop.authentication_enabled", "false"),
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "ospf.authentication_type", "none"),
					resource.TestCheckNoResourceAttr("mso_l3out_interface_group_policy.test", "ospf.key"),
				),
			},
			{
				PreConfig:        func() { fmt.Println("Test: Disable group protocols with empty objects") },
				Config:           testAccMSOL3OutInterfaceGroupPolicyEmptyConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "bfd.enabled", "false"),
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "bfd_multi_hop.enabled", "false"),
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "ospf.enabled", "false"),
				),
			},
			{
				PreConfig:        func() { fmt.Println("Test: Explicit false is equivalent to empty group protocol objects") },
				Config:           testAccMSOL3OutInterfaceGroupPolicyDisabledConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
			},
			{
				PreConfig: func() { fmt.Println("Test: Omit group protocol objects and update description") },
				Config:    testAccMSOL3OutInterfaceGroupPolicyOmittedConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "description", "protocols omitted"),
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "bfd.enabled", "false"),
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "ospf.enabled", "false"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate interface group policy deleted out of band")
					testAccDeleteL3OutInterfaceGroupPolicy(t, templateID, l3outUUID, "test_interface_group")
				},
				Config: testAccMSOL3OutInterfaceGroupPolicyOmittedConfig(siteName, l3outName),
				Check:  resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "name", "test_interface_group"),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import L3Out interface group policy") },
				ResourceName: "mso_l3out_interface_group_policy.test",
				ImportState:  true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					policy, ok := state.RootModule().Resources["mso_l3out_interface_group_policy.test"]
					if !ok {
						return "", fmt.Errorf("interface group policy not found in Terraform state")
					}
					return policy.Primary.Attributes["id"], nil
				},
				ImportStateVerify: true,
			},
			{
				PreConfig:       func() { fmt.Println("Test: Import interface group policy by resource identity") },
				ResourceName:    "mso_l3out_interface_group_policy.test",
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithResourceIdentity,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected one imported interface group policy, got %d", len(states))
					}
					attributes := states[0].Attributes
					if attributes["template_id"] != templateID || attributes["l3out_uuid"] != l3outUUID || attributes["name"] != "test_interface_group" {
						return fmt.Errorf("unexpected imported interface group policy identity: %q/%q/%q", attributes["template_id"], attributes["l3out_uuid"], attributes["name"])
					}
					return nil
				},
			},
			{
				PreConfig:     func() { fmt.Println("Test: Reject invalid interface group policy import ID") },
				ResourceName:  "mso_l3out_interface_group_policy.test",
				ImportState:   true,
				ImportStateId: "invalid-import-id",
				ExpectError:   regexp.MustCompile(`The import ID must be in the form <template_id>/<l3out_uuid>/<name>`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Replace interface group policy when its name changes") },
				Config:    testAccMSOL3OutInterfaceGroupPolicyRenamedConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("mso_l3out_interface_group_policy.test", plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "name", "test_interface_group_renamed"),
					resource.TestCheckResourceAttrPair("mso_l3out_interface_group_policy.test", "l3out_uuid", "mso_l3out.test", "uuid"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate interface group policy after its parent L3Out is deleted out of band")
					testAccDeleteL3Out(t, templateID, l3outUUID)
				},
				Config: testAccMSOL3OutInterfaceGroupPolicyRenamedConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCaptureResourceIdentifiers("mso_l3out.test", &templateID, &l3outUUID),
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "name", "test_interface_group_renamed"),
					resource.TestCheckResourceAttrPair("mso_l3out_interface_group_policy.test", "l3out_uuid", "mso_l3out.test", "uuid"),
				),
			},
		},
	})
}

func testAccMSOL3OutInterfaceGroupPolicyBaseConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  bgp = {
    enabled = true
  }
  ospf = {
    enabled   = true
    area_id   = "0.0.0.7"
    area_type = "regular"
  }
}
`
}

func testAccMSOL3OutInterfaceGroupPolicyCreateConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupPolicyBaseConfig(siteName, l3outName) + `
resource "mso_l3out_interface_group_policy" "test" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  name        = "test_interface_group"
  description = "initial"
}
`
}

func testAccMSOL3OutInterfaceGroupPolicyRenamedConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupPolicyBaseConfig(siteName, l3outName) + `
resource "mso_l3out_interface_group_policy" "test" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  name        = "test_interface_group_renamed"
  description = "protocols omitted"
}
`
}

func testAccMSOL3OutInterfaceGroupPolicyDuplicateConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupPolicyCreateConfig(siteName, l3outName) + `
resource "mso_l3out_interface_group_policy" "duplicate" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  name        = "test_interface_group"
  depends_on  = [mso_l3out_interface_group_policy.test]
}
`
}

func testAccMSOL3OutInterfaceGroupPolicyTwoGroupsConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupPolicyCreateConfig(siteName, l3outName) + `
resource "mso_l3out_interface_group_policy" "second" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  name        = "test_interface_group_second"
  description = "second"
  depends_on  = [mso_l3out_interface_group_policy.test]
}
`
}

func testAccDeleteL3OutInterfaceGroupPolicy(t *testing.T, templateID, l3outUUID, name string) {
	t.Helper()
	template, err := ndoapi.GetTemplate(testAccAPIClient(), templateID)
	if err != nil {
		t.Fatal(err)
	}
	policyPath := models.NewL3OutPath(l3outUUID, "").WithSteps(ndoapi.PathStep{
		Field:    "interfaceGroups",
		Selector: &ndoapi.ObjectSelector{Keys: []ndoapi.ObjectIdentifier{{Field: "name", Value: name}}},
	})
	resolved, err := template.FindRequired(policyPath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := testAccAPIClient().PatchbyID(ndoapi.TemplateEndpoint(templateID), ndoapi.NewRemovePatchPayload(resolved.PatchPath())); err != nil {
		t.Fatal(err)
	}
}

func testAccCheckL3OutInterfaceGroupKeyRefs(templateID, l3outUUID *string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		_, err := testAccReadL3OutInterfaceGroupKeyRefs(*templateID, *l3outUUID)
		return err
	}
}

func testAccCaptureL3OutInterfaceGroupKeyRefs(templateID, l3outUUID *string, captured *map[string]string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		refs, err := testAccReadL3OutInterfaceGroupKeyRefs(*templateID, *l3outUUID)
		if err != nil {
			return err
		}
		*captured = refs
		return nil
	}
}

func testAccCheckL3OutInterfaceGroupKeyRefsUnchanged(templateID, l3outUUID *string, expected *map[string]string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if *expected == nil {
			return fmt.Errorf("authentication key references were not captured")
		}
		actual, err := testAccReadL3OutInterfaceGroupKeyRefs(*templateID, *l3outUUID)
		if err != nil {
			return err
		}
		if !maps.Equal(actual, *expected) {
			return fmt.Errorf("authentication key references changed during unrelated update: got %v, want %v", actual, *expected)
		}
		return nil
	}
}

func testAccReadL3OutInterfaceGroupKeyRefs(templateID, l3outUUID string) (map[string]string, error) {
	template, err := ndoapi.GetTemplate(testAccAPIClient(), templateID)
	if err != nil {
		return nil, err
	}
	policyPath := models.NewL3OutPath(l3outUUID, "").WithSteps(ndoapi.PathStep{
		Field:    "interfaceGroups",
		Selector: &ndoapi.ObjectSelector{Keys: []ndoapi.ObjectIdentifier{{Field: "name", Value: "test_interface_group"}}},
	})
	policy, err := template.FindRequired(policyPath)
	if err != nil {
		return nil, err
	}
	refs := make(map[string]string, 3)
	for _, field := range []string{"bfd", "bfdMultiHop", "ospf"} {
		protocol, _, err := ndoapi.MapField(policy.Object, field, ndoapi.RequiredField)
		if err != nil {
			return nil, fmt.Errorf("NDO %s: %w", field, err)
		}
		key, _, err := ndoapi.MapField(protocol, "key", ndoapi.RequiredField)
		if err != nil {
			return nil, fmt.Errorf("NDO %s.key: %w", field, err)
		}
		ref, _, err := ndoapi.StringField(key, "ref", ndoapi.RequiredNonEmptyField)
		if err != nil {
			return nil, fmt.Errorf("NDO %s.key: %w", field, err)
		}
		refs[field] = ref
	}
	return refs, nil
}

func testAccMSOL3OutInterfaceGroupPolicyEnabledConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupPolicyBaseConfig(siteName, l3outName) + `
resource "mso_l3out_interface_group_policy" "test" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  name        = "test_interface_group"
  description = "initial"
  bfd = {
    enabled = true
  }
  bfd_multi_hop = {
    enabled = true
  }
  ospf = {
    enabled = true
  }
}
`
}

func testAccMSOL3OutInterfaceGroupPolicyAuthenticatedConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupPolicyAuthenticationConfig(siteName, l3outName, 20, 30, 10, "md5")
}

func testAccMSOL3OutInterfaceGroupPolicyAuthenticationConfig(siteName, l3outName string, bfdKeyID, multiHopKeyID, ospfKeyID int, ospfAuthenticationType string) string {
	return testAccMSOL3OutInterfaceGroupPolicyBaseConfig(siteName, l3outName) + fmt.Sprintf(`
resource "mso_l3out_interface_group_policy" "test" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  name        = "test_interface_group"
  description = "initial"
  bfd = {
    enabled                = true
    authentication_enabled = true
    key_id                 = %d
    key                    = "bfd-secret"
  }
  bfd_multi_hop = {
    enabled                = true
    authentication_enabled = true
    key_id                 = %d
    key                    = "multihop-secret"
  }
  ospf = {
    enabled             = true
    authentication_type = %q
    key_id              = %d
    key                 = "ospf-secret"
  }
}
`, bfdKeyID, multiHopKeyID, ospfAuthenticationType, ospfKeyID)
}

func testAccMSOL3OutInterfaceGroupPolicyAuthenticatedOmittedConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupPolicyBaseConfig(siteName, l3outName) + `
resource "mso_l3out_interface_group_policy" "test" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  name        = "test_interface_group"
  description = "authentication objects omitted"
}
`
}

func testAccMSOL3OutInterfaceGroupPolicyAuthenticatedDescriptionUpdateConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupPolicyBaseConfig(siteName, l3outName) + `
resource "mso_l3out_interface_group_policy" "test" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  name        = "test_interface_group"
  description = "authentication unchanged"
  bfd = {
    enabled                = true
    authentication_enabled = true
    key_id                 = 21
    key                    = "bfd-secret"
  }
  bfd_multi_hop = {
    enabled                = true
    authentication_enabled = true
    key_id                 = 31
    key                    = "multihop-secret"
  }
  ospf = {
    enabled             = true
    authentication_type = "simple"
    key_id              = 11
    key                 = "ospf-secret"
  }
}
`
}

func testAccMSOL3OutInterfaceGroupPolicyAuthenticationDisabledConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupPolicyBaseConfig(siteName, l3outName) + `
resource "mso_l3out_interface_group_policy" "test" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  name        = "test_interface_group"
  description = "authentication objects omitted"
  bfd = {
    enabled                = true
    authentication_enabled = false
  }
  bfd_multi_hop = {
    enabled                = true
    authentication_enabled = false
  }
  ospf = {
    enabled             = true
    authentication_type = "none"
  }
}
`
}

func testAccMSOL3OutInterfaceGroupPolicyEmptyConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupPolicyBaseConfig(siteName, l3outName) + `
resource "mso_l3out_interface_group_policy" "test" {
  template_id   = mso_template.l3out_test.id
  l3out_uuid    = mso_l3out.test.uuid
  name          = "test_interface_group"
  description   = "initial"
  bfd           = {}
  bfd_multi_hop = {}
  ospf          = {}
}
`
}

func testAccMSOL3OutInterfaceGroupPolicyDisabledConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupPolicyBaseConfig(siteName, l3outName) + `
resource "mso_l3out_interface_group_policy" "test" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  name        = "test_interface_group"
  description = "initial"
  bfd = {
    enabled = false
  }
  bfd_multi_hop = {
    enabled = false
  }
  ospf = {
    enabled = false
  }
}
`
}

func testAccMSOL3OutInterfaceGroupPolicyOmittedConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupPolicyBaseConfig(siteName, l3outName) + `
resource "mso_l3out_interface_group_policy" "test" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  name        = "test_interface_group"
  description = "protocols omitted"
}
`
}

func TestAccMSOL3OutInterfaceGroupPolicyInvalid(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_interface_group_invalid")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccProviderPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `resource "mso_l3out_interface_group_policy" "test" {
  template_id = "template"
  l3out_uuid  = "l3out"
  name        = "test"
  bfd = {
    enabled                = true
    authentication_enabled = true
  }
}`,
				ExpectError: regexp.MustCompile(`bfd.key must have a value when bfd.authentication_enabled is true`),
			},
			{
				Config: `resource "mso_l3out_interface_group_policy" "test" {
  template_id = "template"
  l3out_uuid  = "l3out"
  name        = "test"
  ospf = {
    enabled             = true
    authentication_type = "md5"
  }
}`,
				ExpectError: regexp.MustCompile(`ospf.key must have a value when ospf.authentication_type is one of simple,\s+md5`),
			},
			{
				Config: `resource "mso_l3out_interface_group_policy" "test" {
  template_id = "template"
  l3out_uuid  = "l3out"
  name        = "test"
  bfd_multi_hop = {
    enabled                = true
    authentication_enabled = true
    key                    = "secret"
  }
}`,
				ExpectError: regexp.MustCompile(`bfd_multi_hop.key_id must have a value when\s+bfd_multi_hop.authentication_enabled is true`),
			},
			{
				Config: `resource "mso_l3out_interface_group_policy" "test" {
  template_id = "template"
  l3out_uuid  = "l3out"
  name        = "test"
  bfd = {
    enabled                = true
    authentication_enabled = false
    key                    = "secret"
  }
}`,
				ExpectError: regexp.MustCompile(`key, key_id may be set only when authentication_enabled is one of true`),
			},
			{
				Config: `resource "mso_l3out_interface_group_policy" "test" {
  template_id = "template"
  l3out_uuid  = "l3out"
  name        = "test"
  ospf = {
    enabled             = true
    authentication_type = "none"
    key_id              = 10
  }
}`,
				ExpectError: regexp.MustCompile(`key, key_id may be set only when authentication_type is one of "simple",\s+"md5"`),
			},
			{
				Config: `resource "mso_l3out_interface_group_policy" "test" {
  template_id = "template"
  l3out_uuid  = "l3out"
  name        = "test"
  bfd = {
    enabled                = false
    authentication_enabled = true
  }
}`,
				ExpectError: regexp.MustCompile(`authentication_enabled requires enabled = true`),
			},
			{
				Config: `resource "mso_l3out_interface_group_policy" "test" {
  template_id = "template"
  l3out_uuid  = "l3out"
  name        = "test"
  netflow_monitor_uuids = {
    invalid = "monitor"
  }
}`,
				ExpectError: regexp.MustCompile(`invalid`),
			},
			{
				Config: `resource "mso_l3out_interface_group_policy" "test" {
  template_id = "template"
  l3out_uuid  = "l3out"
  name        = "test"
  netflow_monitor_uuids = {
    ipv4 = ""
  }
}`,
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Length`),
			},
			{
				Config: `resource "mso_l3out_interface_group_policy" "test" {
  template_id = "template"
  l3out_uuid  = "l3out"
  name        = "test"
  netflow_monitor_uuids = {
    ipv4 = null
  }
}`,
				ExpectError: regexp.MustCompile(`Null Map Value`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject interface group policy with missing template") },
				Config:      testAccMSOL3OutInterfaceGroupPolicyInvalidTemplateConfig(),
				ExpectError: regexp.MustCompile(`L3Out Template Not Found`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject interface group policy with missing L3Out") },
				Config:      testAccMSOL3OutInterfaceGroupPolicyMissingParentConfig(siteName, l3outName),
				ExpectError: regexp.MustCompile(`L3Out Not Found`),
			},
		},
	})
}

func testAccMSOL3OutInterfaceGroupPolicyInvalidTemplateConfig() string {
	return `
resource "mso_l3out_interface_group_policy" "test" {
  template_id = "00000000-0000-0000-0000-000000000000"
  l3out_uuid  = "00000000-0000-0000-0000-000000000000"
  name        = "test"
}
`
}

func testAccMSOL3OutInterfaceGroupPolicyMissingParentConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupPolicyBaseConfig(siteName, l3outName) + `
resource "mso_l3out_interface_group_policy" "test" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = "00000000-0000-0000-0000-000000000000"
  name        = "test"
  depends_on  = [mso_l3out.test]
}
`
}

func TestAccMSOL3OutInterfaceGroupPolicyReferences(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_interface_group_refs")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccProviderPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig:        func() { fmt.Println("Test: Configure all four NetFlow monitor types and tenant policy references") },
				Config:           testAccMSOL3OutInterfaceGroupPolicyAllReferencesConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("mso_l3out_interface_group_policy.test", "interface_routing_policy_uuid", "mso_tenant_policies_l3out_interface_routing_policy.routing", "uuid"),
					resource.TestCheckResourceAttrPair("mso_l3out_interface_group_policy.test", "custom_qos_policy_uuid", "mso_tenant_policies_custom_qos_policy.qos", "uuid"),
					resource.TestCheckResourceAttrPair("mso_l3out_interface_group_policy.test", "netflow_monitor_uuids.ipv4", "mso_tenant_policies_netflow_monitor.ipv4", "uuid"),
					resource.TestCheckResourceAttrPair("mso_l3out_interface_group_policy.test", "netflow_monitor_uuids.ipv6", "mso_tenant_policies_netflow_monitor.ipv6", "uuid"),
					resource.TestCheckResourceAttrPair("mso_l3out_interface_group_policy.test", "netflow_monitor_uuids.ce", "mso_tenant_policies_netflow_monitor.ce", "uuid"),
					resource.TestCheckResourceAttrPair("mso_l3out_interface_group_policy.test", "netflow_monitor_uuids.unspecified", "mso_tenant_policies_netflow_monitor.unspecified", "uuid"),
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "qos_priority", "level6"),
				),
			},
			{
				PreConfig:        func() { fmt.Println("Test: Replace the NetFlow map with two references") },
				Config:           testAccMSOL3OutInterfaceGroupPolicyTwoReferencesConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "netflow_monitor_uuids.%", "2"),
					resource.TestCheckResourceAttrPair("mso_l3out_interface_group_policy.test", "netflow_monitor_uuids.ipv4", "mso_tenant_policies_netflow_monitor.ipv4", "uuid"),
					resource.TestCheckResourceAttrPair("mso_l3out_interface_group_policy.test", "netflow_monitor_uuids.ce", "mso_tenant_policies_netflow_monitor.ce", "uuid"),
				),
			},
			{
				PreConfig:        func() { fmt.Println("Test: Clear NetFlow and tenant policy references") },
				Config:           testAccMSOL3OutInterfaceGroupPolicyClearReferencesConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "netflow_monitor_uuids.%", "0"),
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "interface_routing_policy_uuid", ""),
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "custom_qos_policy_uuid", ""),
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "qos_priority", "unspecified"),
				),
			},
		},
	})
}

func testAccMSOL3OutInterfaceGroupPolicyReferencePrerequisites(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupPolicyBaseConfig(siteName, l3outName) + testAccMSOL3OutInterfaceGroupPolicyTenantReferencesConfig()
}

func testAccMSOL3OutInterfaceGroupPolicyTenantReferencesConfig() string {
	return `
resource "mso_template" "tenant_policy" {
  template_name = "${local.l3out_name}_tenant_policy"
  template_type = "tenant"
  tenant_id     = mso_tenant.l3out_test.id
  sites         = [data.mso_site.l3out_test.id]
}

resource "mso_tenant_policies_custom_qos_policy" "qos" {
  template_id = mso_template.tenant_policy.id
  name        = "${local.l3out_name}_qos"
}

resource "mso_tenant_policies_l3out_interface_routing_policy" "routing" {
  template_id = mso_template.tenant_policy.id
  name        = "${local.l3out_name}_routing"
  depends_on  = [mso_tenant_policies_custom_qos_policy.qos]

  bfd_settings {
    admin_state = "enabled"
  }
}

resource "mso_tenant_policies_netflow_exporter" "exporter" {
  template_id = mso_template.tenant_policy.id
  name        = "${local.l3out_name}_exporter"
  depends_on  = [mso_tenant_policies_l3out_interface_routing_policy.routing]
}

resource "mso_tenant_policies_netflow_record" "record" {
  template_id = mso_template.tenant_policy.id
  name        = "${local.l3out_name}_record"
  depends_on  = [mso_tenant_policies_netflow_exporter.exporter]
}

resource "mso_tenant_policies_netflow_monitor" "ipv4" {
  template_id            = mso_template.tenant_policy.id
  name                   = "${local.l3out_name}_ipv4"
  netflow_record_uuid    = mso_tenant_policies_netflow_record.record.uuid
  netflow_exporter_uuids = [mso_tenant_policies_netflow_exporter.exporter.uuid]
}

resource "mso_tenant_policies_netflow_monitor" "ipv6" {
  template_id            = mso_template.tenant_policy.id
  name                   = "${local.l3out_name}_ipv6"
  netflow_record_uuid    = mso_tenant_policies_netflow_record.record.uuid
  netflow_exporter_uuids = [mso_tenant_policies_netflow_exporter.exporter.uuid]
  depends_on             = [mso_tenant_policies_netflow_monitor.ipv4]
}

resource "mso_tenant_policies_netflow_monitor" "ce" {
  template_id            = mso_template.tenant_policy.id
  name                   = "${local.l3out_name}_ce"
  netflow_record_uuid    = mso_tenant_policies_netflow_record.record.uuid
  netflow_exporter_uuids = [mso_tenant_policies_netflow_exporter.exporter.uuid]
  depends_on             = [mso_tenant_policies_netflow_monitor.ipv6]
}

resource "mso_tenant_policies_netflow_monitor" "unspecified" {
  template_id            = mso_template.tenant_policy.id
  name                   = "${local.l3out_name}_unspecified"
  netflow_record_uuid    = mso_tenant_policies_netflow_record.record.uuid
  netflow_exporter_uuids = [mso_tenant_policies_netflow_exporter.exporter.uuid]
  depends_on             = [mso_tenant_policies_netflow_monitor.ce]
}
`
}

func testAccMSOL3OutInterfaceGroupPolicyAllReferencesConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupPolicyReferencePrerequisites(siteName, l3outName) + `
resource "mso_l3out_interface_group_policy" "test" {
  template_id                   = mso_template.l3out_test.id
  l3out_uuid                    = mso_l3out.test.uuid
  name                          = "test_interface_group"
  interface_routing_policy_uuid = mso_tenant_policies_l3out_interface_routing_policy.routing.uuid
  custom_qos_policy_uuid        = mso_tenant_policies_custom_qos_policy.qos.uuid
  qos_priority                  = "level6"
  netflow_monitor_uuids = {
    ipv4        = mso_tenant_policies_netflow_monitor.ipv4.uuid
    ipv6        = mso_tenant_policies_netflow_monitor.ipv6.uuid
    ce          = mso_tenant_policies_netflow_monitor.ce.uuid
    unspecified = mso_tenant_policies_netflow_monitor.unspecified.uuid
  }
}

`
}

func testAccMSOL3OutInterfaceGroupPolicyTwoReferencesConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupPolicyReferencePrerequisites(siteName, l3outName) + `
resource "mso_l3out_interface_group_policy" "test" {
  template_id                   = mso_template.l3out_test.id
  l3out_uuid                    = mso_l3out.test.uuid
  name                          = "test_interface_group"
  interface_routing_policy_uuid = mso_tenant_policies_l3out_interface_routing_policy.routing.uuid
  custom_qos_policy_uuid        = mso_tenant_policies_custom_qos_policy.qos.uuid
  qos_priority                  = "level6"
  netflow_monitor_uuids = {
    ipv4 = mso_tenant_policies_netflow_monitor.ipv4.uuid
    ce   = mso_tenant_policies_netflow_monitor.ce.uuid
  }
}
`
}

func testAccMSOL3OutInterfaceGroupPolicyClearReferencesConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupPolicyReferencePrerequisites(siteName, l3outName) + `
resource "mso_l3out_interface_group_policy" "test" {
  template_id                   = mso_template.l3out_test.id
  l3out_uuid                    = mso_l3out.test.uuid
  name                          = "test_interface_group"
  interface_routing_policy_uuid = ""
  custom_qos_policy_uuid        = ""
  qos_priority                  = "unspecified"
  netflow_monitor_uuids         = {}
}
`
}

func TestAccMSOL3OutInterfaceGroupPolicyOwnership(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_interface_group_ownership")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccProviderPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig:        func() { fmt.Println("Test: Manage a standalone group while the L3Out map is omitted") },
				Config:           testAccMSOL3OutInterfaceGroupPolicyCreateConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
				Check:            resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.test", "name", "test_interface_group"),
			},
			{
				PreConfig:        func() { fmt.Println("Test: Update the L3Out without managing its interface groups") },
				Config:           testAccMSOL3OutInterfaceGroupPolicyOmittedParentUpdateConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "description", "parent updated"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.%", "1"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.test_interface_group.description", "initial"),
				),
			},
		},
	})
}

func testAccMSOL3OutInterfaceGroupPolicyOmittedParentUpdateConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  description = "parent updated"
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  bgp = {
    enabled = true
  }
  ospf = {
    enabled   = true
    area_id   = "0.0.0.7"
    area_type = "regular"
  }
}

resource "mso_l3out_interface_group_policy" "test" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  name        = "test_interface_group"
  description = "initial"
}
`
}

func TestAccMSOL3OutInterfaceGroupPolicyOwnershipConflict(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_interface_group_conflict")
	var templateID, l3outUUID string
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccProviderPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Configure one interface group on the L3Out") },
				Config:    testAccMSOL3OutInterfaceGroupPolicyOwnershipParentConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCaptureResourceIdentifiers("mso_l3out.test", &templateID, &l3outUUID),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.%", "1"),
					testAccCheckL3OutInterfaceGroups(&templateID, &l3outUUID, map[string]string{"parent": "parent group"}),
				),
			},
			{
				PreConfig:          func() { fmt.Println("Test: Expect a parent plan diff when a child adds a group to a managed map") },
				Config:             testAccMSOL3OutInterfaceGroupPolicyOwnershipConflictConfig(siteName, l3outName),
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("mso_l3out.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(
							"mso_l3out.test",
							tfjsonpath.New("interface_groups"),
							knownvalue.MapExact(map[string]knownvalue.Check{
								"parent": knownvalue.ObjectPartial(map[string]knownvalue.Check{
									"description": knownvalue.StringExact("parent group"),
								}),
							}),
						),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out_interface_group_policy.child", "name", "child"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.parent.description", "parent group"),
					testAccCheckL3OutInterfaceGroups(&templateID, &l3outUUID, map[string]string{
						"parent": "parent group",
						"child":  "standalone group",
					}),
				),
			},
		},
	})
}

func testAccMSOL3OutInterfaceGroupPolicyOwnershipConflictConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupPolicyOwnershipParentConfig(siteName, l3outName) + `
resource "mso_l3out_interface_group_policy" "child" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  name        = "child"
  description = "standalone group"
}
`
}

func testAccMSOL3OutInterfaceGroupPolicyOwnershipParentConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  interface_groups = {
    parent = {
      description = "parent group"
    }
  }
}
`
}

func testAccCheckL3OutInterfaceGroups(templateID, l3outUUID *string, expected map[string]string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		template, err := ndoapi.GetTemplate(testAccAPIClient(), *templateID)
		if err != nil {
			return err
		}
		l3out, err := template.FindRequired(models.NewL3OutPath(*l3outUUID, ""))
		if err != nil {
			return err
		}
		groups, _, err := ndoapi.ListField(l3out.Object, "interfaceGroups", ndoapi.RequiredField)
		if err != nil {
			return err
		}
		if len(groups) != len(expected) {
			return fmt.Errorf("unexpected number of remote interface groups: got %d, want %d", len(groups), len(expected))
		}
		actual := make(map[string]string, len(groups))
		for index, item := range groups {
			group, ok := item.(map[string]any)
			if !ok {
				return fmt.Errorf("NDO interfaceGroups[%d] has unexpected type %T", index, item)
			}
			name, _, err := ndoapi.StringField(group, "name", ndoapi.RequiredNonEmptyField)
			if err != nil {
				return fmt.Errorf("NDO interfaceGroups[%d]: %w", index, err)
			}
			description, _, err := ndoapi.StringField(group, "description", ndoapi.OptionalField)
			if err != nil {
				return fmt.Errorf("NDO interfaceGroups[%d]: %w", index, err)
			}
			actual[name] = description
		}
		if !maps.Equal(actual, expected) {
			return fmt.Errorf("unexpected remote interface groups: got %v, want %v", actual, expected)
		}
		return nil
	}
}
