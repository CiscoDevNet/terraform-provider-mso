package provider

import (
	"fmt"
	"maps"
	"regexp"
	"strings"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-mso/internal/models"
	ndoapi "github.com/CiscoDevNet/terraform-provider-mso/internal/ndoapi"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccMSOL3OutResource(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_l3out")
	var templateID string
	var l3outUUID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccProviderPreCheck(t)
			testAccVersionCheck(t, "4.3")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create L3Out with direct configuration") },
				Config:    testAccMSOL3OutCreateConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCaptureResourceIdentifiers("mso_l3out.test", &templateID, &l3outUUID),
					resource.TestCheckResourceAttr("mso_l3out.test", "name", l3outName),
					testAccCheckL3OutCompositeID("mso_l3out.test"),
					resource.TestCheckResourceAttr("mso_l3out.test", "description", "initial L3Out description"),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"l3_domain",
						"mso_fabric_policies_l3_domain.l3out_test_domain_1",
						"name",
					),
					resource.TestCheckResourceAttr("mso_l3out.test", "target_dscp", "af11"),
					resource.TestCheckResourceAttr("mso_l3out.test", "pim_enabled", "true"),
					resource.TestCheckResourceAttr("mso_l3out.test", "import_route_control_enabled", "true"),
					resource.TestCheckResourceAttrSet("mso_l3out.test", "vrf_uuid"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Deploy standalone L3Out template") },
				Config: testAccMSOL3OutDeployConfig(
					testAccMSOL3OutCreateConfig(siteName, l3outName),
					"l3out_test",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"mso_schema_template_deploy_ndo.l3out_test",
						"template_id",
						"mso_template.l3out_test",
						"id",
					),
					resource.TestCheckResourceAttr(
						"mso_schema_template_deploy_ndo.l3out_test",
						"template_type",
						"l3out",
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Undeploy standalone L3Out template") },
				Config:    testAccMSOL3OutCreateConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "name", l3outName),
					resource.TestCheckResourceAttrSet("mso_l3out.test", "uuid"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out direct configuration") },
				Config:    testAccMSOL3OutUpdateConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "description", "updated L3Out description"),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"l3_domain",
						"mso_fabric_policies_l3_domain.l3out_test_domain_2",
						"name",
					),
					resource.TestCheckResourceAttr("mso_l3out.test", "target_dscp", "expedited_forwarding"),
					resource.TestCheckResourceAttr("mso_l3out.test", "pim_enabled", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "import_route_control_enabled", "false"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Omit Optional+Computed attributes and preserve remote values") },
				Config:    testAccMSOL3OutOmitOptionalComputedConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "description", "updated L3Out description"),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"l3_domain",
						"mso_fabric_policies_l3_domain.l3out_test_domain_2",
						"name",
					),
					resource.TestCheckResourceAttr("mso_l3out.test", "target_dscp", "expedited_forwarding"),
					resource.TestCheckResourceAttr("mso_l3out.test", "pim_enabled", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "import_route_control_enabled", "false"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Clear direct string attributes and reset DSCP") },
				Config:    testAccMSOL3OutClearConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "description", ""),
					resource.TestCheckResourceAttr("mso_l3out.test", "l3_domain", ""),
					resource.TestCheckResourceAttr("mso_l3out.test", "target_dscp", "unspecified"),
					resource.TestCheckResourceAttr("mso_l3out.test", "pim_enabled", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "import_route_control_enabled", "false"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate L3Out after out-of-band deletion")
					testAccDeleteL3Out(t, templateID, l3outUUID)
				},
				Config: testAccMSOL3OutCreateConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCaptureResourceIdentifiers("mso_l3out.test", &templateID, &l3outUUID),
					resource.TestCheckResourceAttr("mso_l3out.test", "name", l3outName),
					resource.TestCheckResourceAttr("mso_l3out.test", "description", "initial L3Out description"),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"l3_domain",
						"mso_fabric_policies_l3_domain.l3out_test_domain_1",
						"name",
					),
					resource.TestCheckResourceAttr("mso_l3out.test", "target_dscp", "af11"),
					resource.TestCheckResourceAttr("mso_l3out.test", "pim_enabled", "true"),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import existing L3Out") },
				Config:       testAccMSOL3OutImportConfig(siteName, l3outName),
				ResourceName: "mso_l3out.test",
				ImportState:  true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					resourceState, ok := state.RootModule().Resources["mso_l3out.test"]
					if !ok {
						return "", fmt.Errorf("mso_l3out.test not found in Terraform state")
					}
					return fmt.Sprintf(
						"%s/%s",
						resourceState.Primary.Attributes["template_id"],
						resourceState.Primary.Attributes["uuid"],
					), nil
				},
				ImportStateVerify: true,
			},
			{
				PreConfig:       func() { fmt.Println("Test: Import existing L3Out by resource identity") },
				ResourceName:    "mso_l3out.test",
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithResourceIdentity,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected one imported L3Out, got %d", len(states))
					}
					attributes := states[0].Attributes
					if attributes["template_id"] != templateID || attributes["uuid"] != l3outUUID {
						return fmt.Errorf(
							"unexpected imported L3Out identity: %q/%q",
							attributes["template_id"],
							attributes["uuid"],
						)
					}
					return nil
				},
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Replace L3Out template and verify the original L3Out is deleted")
				},
				Config: testAccMSOL3OutForceNewConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccAssertResourceIDChanged("mso_l3out.test", &l3outUUID),
					testAccCheckL3OutDeleted(&templateID, &l3outUUID),
					resource.TestCheckResourceAttr("mso_l3out.test", "name", l3outName),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"template_id",
						"mso_template.l3out_test_replacement",
						"id",
					),
					resource.TestCheckResourceAttrSet("mso_l3out.test", "vrf_uuid"),
				),
			},
			{
				PreConfig:     func() { fmt.Println("Test: Reject invalid L3Out import ID") },
				ResourceName:  "mso_l3out.test",
				ImportState:   true,
				ImportStateId: "invalid-l3out-import-id",
				ExpectError: regexp.MustCompile(
					`The import ID must be in the form <template_id>/<l3out_uuid>`,
				),
			},
		},
	})
}

func testAccMSOL3OutCreateConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id                  = mso_template.l3out_test.id
  name                         = local.l3out_name
  vrf_uuid                     = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description                  = "initial L3Out description"
  l3_domain                    = mso_fabric_policies_l3_domain.l3out_test_domain_1.name
  target_dscp                  = "af11"
  pim_enabled                  = true
  import_route_control_enabled = true
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutUpdateConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id                  = mso_template.l3out_test.id
  name                         = local.l3out_name
  vrf_uuid                     = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description                  = "updated L3Out description"
  l3_domain                    = mso_fabric_policies_l3_domain.l3out_test_domain_2.name
  target_dscp                  = "expedited_forwarding"
  pim_enabled                  = false
  import_route_control_enabled = false
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutOmitOptionalComputedConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  pim_enabled = false
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutClearConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id                  = mso_template.l3out_test.id
  name                         = local.l3out_name
  vrf_uuid                     = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description                  = ""
  l3_domain                    = ""
  target_dscp                  = "unspecified"
  pim_enabled                  = false
  import_route_control_enabled = false
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutImportConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description = "updated L3Out description"
  l3_domain   = mso_fabric_policies_l3_domain.l3out_test_domain_2.name
  target_dscp = "expedited_forwarding"
  pim_enabled = false
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutForceNewConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_template" "l3out_test_replacement" {
  template_name = "${local.l3out_name}_replacement_template"
  template_type = "l3out"
  tenant_id     = mso_tenant.l3out_test.id
  sites         = [data.mso_site.l3out_test.id]
}

resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test_replacement.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_2.uuid
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutDeployConfig(config, templateResource string) string {
	return fmt.Sprintf(`%s
resource "mso_schema_template_deploy_ndo" "l3out_test" {
  template_id         = mso_template.%s.id
  template_type       = "l3out"
  force_apply         = ""
  undeploy_on_destroy = true

  depends_on = [mso_l3out.test]
}
`, config, templateResource)
}

func TestAccMSOL3OutFullyConfigured(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_l3out_fully_configured")
	var templateID, l3outUUID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccProviderPreCheck(t)
			testAccVersionCheck(t, "4.3")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: Create fully configured L3Out with both protocols and all route-map references")
				},
				Config: testAccMSOL3OutFullyConfiguredConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCaptureResourceIdentifiers("mso_l3out.test", &templateID, &l3outUUID),
					testAccCheckL3OutCompositeID("mso_l3out.test"),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"template_id",
						"mso_template.l3out_test",
						"id",
					),
					resource.TestCheckResourceAttr("mso_l3out.test", "name", l3outName),
					resource.TestCheckResourceAttr("mso_l3out.test", "description", "fully configured L3Out"),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"vrf_uuid",
						"mso_schema_template_vrf.l3out_test_vrf_1",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"l3_domain",
						"mso_fabric_policies_l3_domain.l3out_test_domain_1",
						"name",
					),
					resource.TestCheckResourceAttr("mso_l3out.test", "target_dscp", "expedited_forwarding"),
					resource.TestCheckResourceAttr("mso_l3out.test", "pim_enabled", "true"),
					resource.TestCheckResourceAttr("mso_l3out.test", "import_route_control_enabled", "true"),
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.enabled", "true"),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"bgp.route_dampening_ipv4_uuid",
						"mso_tenant_policies_route_map_policy_route_control.dampening_v4",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"bgp.route_dampening_ipv6_uuid",
						"mso_tenant_policies_route_map_policy_route_control.dampening_v6",
						"uuid",
					),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.enabled", "true"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.area_id", "0.0.0.10"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.area_type", "regular"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.cost", "10"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.send_redistributed_lsa", "true"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.originate_summary_lsa", "true"),
					resource.TestCheckResourceAttr(
						"mso_l3out.test",
						"ospf.suppress_forwarding_address_in_translated_lsa",
						"true",
					),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.originate_default_route_always", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "originate_default_route", "in_addition"),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.outbound_uuid",
						"mso_tenant_policies_route_map_policy_route_control.outbound",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.interleak_uuid",
						"mso_tenant_policies_route_map_policy_route_control.interleak",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.static_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.static",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.connected_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.connected",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.attached_host_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.attached",
						"uuid",
					),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "bgpOspf"),
				),
			},
			{
				SkipFunc: func() (bool, error) {
					result, err := testAccAPIClient().CompareVersion("5.1")
					if err != nil || result > 0 {
						return result > 0, err
					}
					fmt.Println("Test: Configure route-map match rules, set rules and contexts")
					return false, nil
				},
				Config: testAccMSOL3OutFullyConfiguredWithRouteMapRulesConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					func(state *terraform.State) error {
						l3out := state.RootModule().Resources["mso_l3out.test"]
						if l3out == nil || l3out.Primary == nil || l3out.Primary.Attributes["uuid"] != l3outUUID {
							return fmt.Errorf("L3Out UUID changed while configuring route-map rules")
						}
						return nil
					},
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.outbound_uuid",
						"mso_tenant_policies_route_map_policy_route_control.outbound",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_tenant_policies_route_map_policy_route_control_context.outbound",
						"route_map_policy_uuid",
						"mso_tenant_policies_route_map_policy_route_control.outbound",
						"uuid",
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_route_map_policy_route_control_context.outbound",
						"match_rule_uuids.#",
						"1",
					),
					resource.TestCheckTypeSetElemAttrPair(
						"mso_tenant_policies_route_map_policy_route_control_context.outbound",
						"match_rule_uuids.*",
						"mso_tenant_policies_match_rule_policy.outbound",
						"uuid",
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_route_map_policy_route_control_context.outbound",
						"action",
						"permit",
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_route_map_policy_route_control_context.outbound",
						"order",
						"0",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.interleak_uuid",
						"mso_tenant_policies_route_map_policy_route_control.interleak",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_tenant_policies_route_map_policy_route_control_context.interleak",
						"route_map_policy_uuid",
						"mso_tenant_policies_route_map_policy_route_control.interleak",
						"uuid",
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_route_map_policy_route_control_context.interleak",
						"match_rule_uuids.#",
						"1",
					),
					resource.TestCheckTypeSetElemAttrPair(
						"mso_tenant_policies_route_map_policy_route_control_context.interleak",
						"match_rule_uuids.*",
						"mso_tenant_policies_match_rule_policy.interleak",
						"uuid",
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_route_map_policy_route_control_context.interleak",
						"action",
						"permit",
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_route_map_policy_route_control_context.interleak",
						"order",
						"0",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.static_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.static",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_tenant_policies_route_map_policy_route_control_context.static",
						"route_map_policy_uuid",
						"mso_tenant_policies_route_map_policy_route_control.static",
						"uuid",
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_route_map_policy_route_control_context.static",
						"match_rule_uuids.#",
						"1",
					),
					resource.TestCheckTypeSetElemAttrPair(
						"mso_tenant_policies_route_map_policy_route_control_context.static",
						"match_rule_uuids.*",
						"mso_tenant_policies_match_rule_policy.static",
						"uuid",
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_route_map_policy_route_control_context.static",
						"action",
						"permit",
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_route_map_policy_route_control_context.static",
						"order",
						"0",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.connected_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.connected",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_tenant_policies_route_map_policy_route_control_context.connected",
						"route_map_policy_uuid",
						"mso_tenant_policies_route_map_policy_route_control.connected",
						"uuid",
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_route_map_policy_route_control_context.connected",
						"match_rule_uuids.#",
						"1",
					),
					resource.TestCheckTypeSetElemAttrPair(
						"mso_tenant_policies_route_map_policy_route_control_context.connected",
						"match_rule_uuids.*",
						"mso_tenant_policies_match_rule_policy.connected",
						"uuid",
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_route_map_policy_route_control_context.connected",
						"action",
						"permit",
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_route_map_policy_route_control_context.connected",
						"order",
						"0",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.attached_host_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.attached",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_tenant_policies_route_map_policy_route_control_context.attached",
						"route_map_policy_uuid",
						"mso_tenant_policies_route_map_policy_route_control.attached",
						"uuid",
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_route_map_policy_route_control_context.attached",
						"match_rule_uuids.#",
						"1",
					),
					resource.TestCheckTypeSetElemAttrPair(
						"mso_tenant_policies_route_map_policy_route_control_context.attached",
						"match_rule_uuids.*",
						"mso_tenant_policies_match_rule_policy.attached",
						"uuid",
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_route_map_policy_route_control_context.attached",
						"action",
						"permit",
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_route_map_policy_route_control_context.attached",
						"order",
						"0",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"bgp.route_dampening_ipv4_uuid",
						"mso_tenant_policies_route_map_policy_route_control.dampening_v4",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_tenant_policies_route_map_policy_route_control_context.dampening_v4",
						"route_map_policy_uuid",
						"mso_tenant_policies_route_map_policy_route_control.dampening_v4",
						"uuid",
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_route_map_policy_route_control_context.dampening_v4",
						"match_rule_uuids.#",
						"1",
					),
					resource.TestCheckTypeSetElemAttrPair(
						"mso_tenant_policies_route_map_policy_route_control_context.dampening_v4",
						"match_rule_uuids.*",
						"mso_tenant_policies_match_rule_policy.dampening_v4",
						"uuid",
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_route_map_policy_route_control_context.dampening_v4",
						"action",
						"permit",
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_route_map_policy_route_control_context.dampening_v4",
						"order",
						"0",
					),
					resource.TestCheckResourceAttrPair(
						"mso_tenant_policies_route_map_policy_route_control_context.dampening_v4",
						"set_rule_uuid",
						"mso_tenant_policies_set_rule_policy.dampening_v4",
						"uuid",
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						"mso_tenant_policies_match_rule_policy.dampening_v4",
						"match_prefix.*",
						map[string]string{"prefix": "203.0.113.0/24"},
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_set_rule_policy.dampening_v4",
						"set_dampening.0.half_life",
						"15",
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_set_rule_policy.dampening_v4",
						"set_dampening.0.reuse_limit",
						"750",
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_set_rule_policy.dampening_v4",
						"set_dampening.0.suppress_limit",
						"2000",
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_set_rule_policy.dampening_v4",
						"set_dampening.0.max_suppress_time",
						"60",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"bgp.route_dampening_ipv6_uuid",
						"mso_tenant_policies_route_map_policy_route_control.dampening_v6",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_tenant_policies_route_map_policy_route_control_context.dampening_v6",
						"route_map_policy_uuid",
						"mso_tenant_policies_route_map_policy_route_control.dampening_v6",
						"uuid",
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_route_map_policy_route_control_context.dampening_v6",
						"match_rule_uuids.#",
						"1",
					),
					resource.TestCheckTypeSetElemAttrPair(
						"mso_tenant_policies_route_map_policy_route_control_context.dampening_v6",
						"match_rule_uuids.*",
						"mso_tenant_policies_match_rule_policy.dampening_v6",
						"uuid",
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_route_map_policy_route_control_context.dampening_v6",
						"action",
						"permit",
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_route_map_policy_route_control_context.dampening_v6",
						"order",
						"0",
					),
					resource.TestCheckResourceAttrPair(
						"mso_tenant_policies_route_map_policy_route_control_context.dampening_v6",
						"set_rule_uuid",
						"mso_tenant_policies_set_rule_policy.dampening_v6",
						"uuid",
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						"mso_tenant_policies_match_rule_policy.dampening_v6",
						"match_prefix.*",
						map[string]string{"prefix": "2001:db8:10::/48"},
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_set_rule_policy.dampening_v6",
						"set_dampening.0.half_life",
						"15",
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_set_rule_policy.dampening_v6",
						"set_dampening.0.reuse_limit",
						"750",
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_set_rule_policy.dampening_v6",
						"set_dampening.0.suppress_limit",
						"2000",
					),
					resource.TestCheckResourceAttr(
						"mso_tenant_policies_set_rule_policy.dampening_v6",
						"set_dampening.0.max_suppress_time",
						"60",
					),
				),
			},
		},
	})
}

func testAccMSOL3OutFullyConfiguredConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s%s
resource "mso_l3out" "test" {
  template_id                  = mso_template.l3out_test.id
  name                         = local.l3out_name
  vrf_uuid                     = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description                  = "fully configured L3Out"
  l3_domain                    = mso_fabric_policies_l3_domain.l3out_test_domain_1.name
  target_dscp                  = "expedited_forwarding"
  pim_enabled                  = true
  import_route_control_enabled = true
  originate_default_route      = "in_addition"

  bgp = {
    enabled                   = true
    route_dampening_ipv4_uuid = mso_tenant_policies_route_map_policy_route_control.dampening_v4.uuid
    route_dampening_ipv6_uuid = mso_tenant_policies_route_map_policy_route_control.dampening_v6.uuid
  }

  ospf = {
    enabled                                       = true
    area_id                                       = "0.0.0.10"
    area_type                                     = "regular"
    cost                                          = 10
    send_redistributed_lsa                        = true
    originate_summary_lsa                         = true
    suppress_forwarding_address_in_translated_lsa = true
    originate_default_route_always                = false
  }

  route_maps = {
    outbound_uuid                           = mso_tenant_policies_route_map_policy_route_control.outbound.uuid
    interleak_uuid                          = mso_tenant_policies_route_map_policy_route_control.interleak.uuid
    static_route_redistribution_uuid        = mso_tenant_policies_route_map_policy_route_control.static.uuid
    connected_route_redistribution_uuid     = mso_tenant_policies_route_map_policy_route_control.connected.uuid
    attached_host_route_redistribution_uuid = mso_tenant_policies_route_map_policy_route_control.attached.uuid
  }
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName), testAccMSOL3OutRouteMapPolicyPrerequisites())
}

func testAccMSOL3OutFullyConfiguredWithRouteMapRulesConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_match_rule_policy" "outbound" {
  template_id = mso_template.l3out_test_route_map_policies.id
  name        = "${local.l3out_name}_outbound_match"

  match_prefix {
    prefix = "203.0.113.0/24"
  }
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "outbound" {
  template_id           = mso_template.l3out_test_route_map_policies.id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.outbound.uuid
  name                  = "outbound_context"
  order                 = 0
  action                = "permit"
  match_rule_uuids      = [mso_tenant_policies_match_rule_policy.outbound.uuid]
}

resource "mso_tenant_policies_match_rule_policy" "interleak" {
  template_id = mso_template.l3out_test_route_map_policies.id
  name        = "${local.l3out_name}_interleak_match"

  match_prefix {
    prefix = "10.10.0.0/16"
  }
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "interleak" {
  template_id           = mso_template.l3out_test_route_map_policies.id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.interleak.uuid
  name                  = "interleak_context"
  order                 = 0
  action                = "permit"
  match_rule_uuids      = [mso_tenant_policies_match_rule_policy.interleak.uuid]
}

resource "mso_tenant_policies_match_rule_policy" "static" {
  template_id = mso_template.l3out_test_route_map_policies.id
  name        = "${local.l3out_name}_static_match"

  match_prefix {
    prefix = "198.51.100.0/24"
  }
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "static" {
  template_id           = mso_template.l3out_test_route_map_policies.id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.static.uuid
  name                  = "static_context"
  order                 = 0
  action                = "permit"
  match_rule_uuids      = [mso_tenant_policies_match_rule_policy.static.uuid]
}

resource "mso_tenant_policies_match_rule_policy" "connected" {
  template_id = mso_template.l3out_test_route_map_policies.id
  name        = "${local.l3out_name}_connected_match"

  match_prefix {
    prefix = "192.0.2.0/24"
  }
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "connected" {
  template_id           = mso_template.l3out_test_route_map_policies.id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.connected.uuid
  name                  = "connected_context"
  order                 = 0
  action                = "permit"
  match_rule_uuids      = [mso_tenant_policies_match_rule_policy.connected.uuid]
}

resource "mso_tenant_policies_match_rule_policy" "attached" {
  template_id = mso_template.l3out_test_route_map_policies.id
  name        = "${local.l3out_name}_attached_match"

  match_prefix {
    prefix = "198.51.100.10/32"
  }
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "attached" {
  template_id           = mso_template.l3out_test_route_map_policies.id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.attached.uuid
  name                  = "attached_context"
  order                 = 0
  action                = "permit"
  match_rule_uuids      = [mso_tenant_policies_match_rule_policy.attached.uuid]
}

resource "mso_tenant_policies_match_rule_policy" "dampening_v4" {
  template_id = mso_template.l3out_test_route_map_policies.id
  name        = "${local.l3out_name}_dampening_v4_match"

  match_prefix {
    prefix = "203.0.113.0/24"
  }
}

resource "mso_tenant_policies_set_rule_policy" "dampening_v4" {
  template_id = mso_template.l3out_test_route_map_policies.id
  name        = "${local.l3out_name}_dampening_v4_set"
  set_metric  = 100

  set_route_tag {
    value = 500
  }
  set_weight {
    value = 50
  }
  set_preference {
    value = 200
  }
  set_dampening {
    half_life         = 15
    reuse_limit       = 750
    suppress_limit    = 2000
    max_suppress_time = 60
  }
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "dampening_v4" {
  template_id           = mso_template.l3out_test_route_map_policies.id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.dampening_v4.uuid
  name                  = "dampening_v4_context"
  order                 = 0
  action                = "permit"
  match_rule_uuids      = [mso_tenant_policies_match_rule_policy.dampening_v4.uuid]
  set_rule_uuid         = mso_tenant_policies_set_rule_policy.dampening_v4.uuid
}

resource "mso_tenant_policies_match_rule_policy" "dampening_v6" {
  template_id = mso_template.l3out_test_route_map_policies.id
  name        = "${local.l3out_name}_dampening_v6_match"

  match_prefix {
    prefix = "2001:db8:10::/48"
  }
}

resource "mso_tenant_policies_set_rule_policy" "dampening_v6" {
  template_id = mso_template.l3out_test_route_map_policies.id
  name        = "${local.l3out_name}_dampening_v6_set"
  set_metric  = 100

  set_route_tag {
    value = 500
  }
  set_weight {
    value = 50
  }
  set_preference {
    value = 200
  }
  set_dampening {
    half_life         = 15
    reuse_limit       = 750
    suppress_limit    = 2000
    max_suppress_time = 60
  }
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "dampening_v6" {
  template_id           = mso_template.l3out_test_route_map_policies.id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.dampening_v6.uuid
  name                  = "dampening_v6_context"
  order                 = 0
  action                = "permit"
  match_rule_uuids      = [mso_tenant_policies_match_rule_policy.dampening_v6.uuid]
  set_rule_uuid         = mso_tenant_policies_set_rule_policy.dampening_v6.uuid
}
`, testAccMSOL3OutFullyConfiguredConfig(siteName, l3outName))
}

func TestAccMSOL3OutRoutingProtocol(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_l3out_routing_protocol")
	var templateID, l3outUUID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccProviderPreCheck(t)
			testAccVersionCheck(t, "4.3")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create L3Out with OSPF explicitly disabled") },
				Config:    testAccMSOL3OutOSPFEmptyConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCaptureResourceIdentifiers("mso_l3out.test", &templateID, &l3outUUID),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.enabled", "false"),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "none"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Enable OSPF with only its required area settings") },
				Config:    testAccMSOL3OutOSPFMinimalConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.enabled", "true"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.area_id", "0.0.0.1"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.area_type", "regular"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.cost", "1"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.send_redistributed_lsa", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.originate_summary_lsa", "false"),
					resource.TestCheckResourceAttr(
						"mso_l3out.test",
						"ospf.suppress_forwarding_address_in_translated_lsa",
						"false",
					),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.originate_default_route_always", "false"),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "ospf"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out description with OSPF cost and controls omitted") },
				Config:    testAccMSOL3OutOSPFMinimalUpdateConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("mso_l3out.test", plancheck.ResourceActionUpdate),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "description", "OSPF default settings update"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.cost", "1"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.send_redistributed_lsa", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.originate_summary_lsa", "false"),
					resource.TestCheckResourceAttr(
						"mso_l3out.test",
						"ospf.suppress_forwarding_address_in_translated_lsa",
						"false",
					),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "ospf"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Enable OSPF with all attributes configured") },
				Config:    testAccMSOL3OutOSPFEnabledConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.enabled", "true"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.area_id", "0.0.0.1"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.area_type", "regular"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.cost", "10"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.send_redistributed_lsa", "true"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.originate_summary_lsa", "true"),
					resource.TestCheckResourceAttr(
						"mso_l3out.test",
						"ospf.suppress_forwarding_address_in_translated_lsa",
						"true",
					),
					resource.TestCheckResourceAttr("mso_l3out.test", "originate_default_route", "only"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.originate_default_route_always", "false"),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "ospf"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update OSPF and L3Out default route") },
				Config:    testAccMSOL3OutOSPFUpdatedConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.area_id", "0.0.0.2"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.area_type", "nssa"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.cost", "20"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.send_redistributed_lsa", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.originate_summary_lsa", "false"),
					resource.TestCheckResourceAttr(
						"mso_l3out.test",
						"ospf.suppress_forwarding_address_in_translated_lsa",
						"false",
					),
					resource.TestCheckResourceAttr("mso_l3out.test", "originate_default_route", "in_addition"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.originate_default_route_always", "true"),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "ospf"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Omit OSPF without changing L3Out configuration") },
				Config:    testAccMSOL3OutOSPFOmittedConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.area_id", "0.0.0.2"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.area_type", "nssa"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.cost", "20"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.send_redistributed_lsa", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.originate_summary_lsa", "false"),
					resource.TestCheckResourceAttr(
						"mso_l3out.test",
						"ospf.suppress_forwarding_address_in_translated_lsa",
						"false",
					),
					resource.TestCheckResourceAttr("mso_l3out.test", "originate_default_route", "in_addition"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.originate_default_route_always", "true"),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "ospf"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out description while omitting OSPF") },
				Config:    testAccMSOL3OutOSPFOmittedUpdateConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("mso_l3out.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "description", "OSPF omission update"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.area_id", "0.0.0.2"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.area_type", "nssa"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.cost", "20"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.send_redistributed_lsa", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.originate_summary_lsa", "false"),
					resource.TestCheckResourceAttr(
						"mso_l3out.test",
						"ospf.suppress_forwarding_address_in_translated_lsa",
						"false",
					),
					resource.TestCheckResourceAttr("mso_l3out.test", "originate_default_route", "in_addition"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.originate_default_route_always", "true"),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "ospf"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Clear L3Out default-route origination with OSPF always omitted")
				},
				Config: testAccMSOL3OutDefaultRouteEmptyConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("mso_l3out.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.area_id", "0.0.0.2"),
					resource.TestCheckResourceAttr("mso_l3out.test", "originate_default_route", ""),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.originate_default_route_always", "false"),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "ospf"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Disable OSPF with an empty object") },
				Config:    testAccMSOL3OutOSPFEmptyConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("mso_l3out.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.enabled", "false"),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "none"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Explicitly disabled OSPF is equivalent to an empty object") },
				Config:    testAccMSOL3OutOSPFDisabledConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.enabled", "false"),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "none"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Configure BGP without OSPF") },
				Config:    testAccMSOL3OutBGPConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.enabled", "true"),
					resource.TestCheckResourceAttr("mso_l3out.test", "originate_default_route", "in_addition"),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "bgp"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Omit BGP without changing its state") },
				Config:    testAccMSOL3OutBGPOmittedConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.enabled", "true"),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "bgp"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out description while omitting enabled BGP") },
				Config:    testAccMSOL3OutBGPOmittedUpdateConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("mso_l3out.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "description", "BGP omission update"),
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.enabled", "true"),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "bgp"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Configure BGP and OSPF together") },
				Config:    testAccMSOL3OutBGPWithOSPFConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.enabled", "true"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.area_id", "0.0.0.3"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.area_type", "regular"),
					resource.TestCheckResourceAttr("mso_l3out.test", "originate_default_route", "in_addition"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.originate_default_route_always", "false"),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "bgpOspf"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Disable BGP while retaining OSPF") },
				Config:    testAccMSOL3OutBGPDisabledWithOSPFConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.enabled", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.area_id", "0.0.0.3"),
					resource.TestCheckResourceAttr("mso_l3out.test", "originate_default_route", "in_addition"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.originate_default_route_always", "false"),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "ospf"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Disable both protocols with empty objects") },
				Config:    testAccMSOL3OutBothProtocolsEmptyConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.enabled", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.enabled", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "originate_default_route", ""),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "none"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update description while omitting both empty protocol objects")
				},
				Config: testAccMSOL3OutBothProtocolsOmittedAfterEmptyConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("mso_l3out.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "description", "protocols omitted after empty"),
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.enabled", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.enabled", "false"),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "none"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Explicit false matches empty objects in state") },
				Config:    testAccMSOL3OutBothProtocolsDisabledConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.enabled", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.enabled", "false"),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "none"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update description while omitting both explicitly disabled protocol objects")
				},
				Config: testAccMSOL3OutBothProtocolsOmittedAfterFalseConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("mso_l3out.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "description", "protocols omitted after false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.enabled", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.enabled", "false"),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "none"),
				),
			},
		},
	})
}

func TestAccMSOL3OutRouteMaps(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_l3out_route_maps")
	var templateID, l3outUUID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccProviderPreCheck(t)
			testAccVersionCheck(t, "4.3")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Configure common route maps without a routing protocol") },
				Config:    testAccMSOL3OutRouteMapsNoProtocolConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCaptureResourceIdentifiers("mso_l3out.test", &templateID, &l3outUUID),
					resource.TestCheckNoResourceAttr("mso_l3out.test", "route_maps.outbound_uuid"),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.interleak_uuid",
						"mso_tenant_policies_route_map_policy_route_control.interleak",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.static_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.static",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.connected_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.connected",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.attached_host_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.attached",
						"uuid",
					),
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.enabled", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.enabled", "false"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Configure BGP dampening route maps without an outbound route map")
				},
				Config: testAccMSOL3OutRouteMapsBGPConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.enabled", "true"),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"bgp.route_dampening_ipv4_uuid",
						"mso_tenant_policies_route_map_policy_route_control.dampening_v4",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"bgp.route_dampening_ipv6_uuid",
						"mso_tenant_policies_route_map_policy_route_control.dampening_v6",
						"uuid",
					),
					resource.TestCheckNoResourceAttr("mso_l3out.test", "route_maps.outbound_uuid"),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "bgp"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Clear IPv4 dampening while preserving omitted IPv6 and common route maps")
				},
				Config: testAccMSOL3OutBGPDampeningIPv4ClearConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.enabled", "true"),
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.route_dampening_ipv4_uuid", ""),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"bgp.route_dampening_ipv6_uuid",
						"mso_tenant_policies_route_map_policy_route_control.dampening_v6",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.interleak_uuid",
						"mso_tenant_policies_route_map_policy_route_control.interleak",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.static_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.static",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.connected_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.connected",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.attached_host_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.attached",
						"uuid",
					),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "bgp"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Restore both BGP dampening route maps before changing protocols")
				},
				Config: testAccMSOL3OutRouteMapsBGPConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.enabled", "true"),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"bgp.route_dampening_ipv4_uuid",
						"mso_tenant_policies_route_map_policy_route_control.dampening_v4",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"bgp.route_dampening_ipv6_uuid",
						"mso_tenant_policies_route_map_policy_route_control.dampening_v6",
						"uuid",
					),
					resource.TestCheckNoResourceAttr("mso_l3out.test", "route_maps.outbound_uuid"),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "bgp"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Disable BGP with an empty object and clear both dampening references")
				},
				Config: testAccMSOL3OutRouteMapsBGPEmptyConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.enabled", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.enabled", "false"),
					resource.TestCheckNoResourceAttr("mso_l3out.test", "bgp.route_dampening_ipv4_uuid"),
					resource.TestCheckNoResourceAttr("mso_l3out.test", "bgp.route_dampening_ipv6_uuid"),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.interleak_uuid",
						"mso_tenant_policies_route_map_policy_route_control.interleak",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.static_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.static",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.connected_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.connected",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.attached_host_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.attached",
						"uuid",
					),
					resource.TestCheckNoResourceAttr("mso_l3out.test", "route_maps.outbound_uuid"),
					testAccCheckL3OutBGPDampeningCleared(&templateID, &l3outUUID),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "none"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Restore both BGP dampening references after an empty BGP object")
				},
				Config: testAccMSOL3OutRouteMapsBGPConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.enabled", "true"),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"bgp.route_dampening_ipv4_uuid",
						"mso_tenant_policies_route_map_policy_route_control.dampening_v4",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"bgp.route_dampening_ipv6_uuid",
						"mso_tenant_policies_route_map_policy_route_control.dampening_v6",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.interleak_uuid",
						"mso_tenant_policies_route_map_policy_route_control.interleak",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.static_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.static",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.connected_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.connected",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.attached_host_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.attached",
						"uuid",
					),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.enabled", "false"),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "bgp"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Disable BGP explicitly and clear both dampening references")
				},
				Config: testAccMSOL3OutRouteMapsBGPDisabledConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.enabled", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.enabled", "false"),
					resource.TestCheckNoResourceAttr("mso_l3out.test", "bgp.route_dampening_ipv4_uuid"),
					resource.TestCheckNoResourceAttr("mso_l3out.test", "bgp.route_dampening_ipv6_uuid"),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.interleak_uuid",
						"mso_tenant_policies_route_map_policy_route_control.interleak",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.static_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.static",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.connected_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.connected",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.attached_host_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.attached",
						"uuid",
					),
					resource.TestCheckNoResourceAttr("mso_l3out.test", "route_maps.outbound_uuid"),
					testAccCheckL3OutBGPDampeningCleared(&templateID, &l3outUUID),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "none"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Restore both dampening references and enable OSPF before the mixed clear")
				},
				Config: testAccMSOL3OutRouteMapsBGPWithOSPFConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.enabled", "true"),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"bgp.route_dampening_ipv4_uuid",
						"mso_tenant_policies_route_map_policy_route_control.dampening_v4",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"bgp.route_dampening_ipv6_uuid",
						"mso_tenant_policies_route_map_policy_route_control.dampening_v6",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.interleak_uuid",
						"mso_tenant_policies_route_map_policy_route_control.interleak",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.static_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.static",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.connected_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.connected",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.attached_host_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.attached",
						"uuid",
					),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.enabled", "true"),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.outbound_uuid",
						"mso_tenant_policies_route_map_policy_route_control.outbound",
						"uuid",
					),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "bgpOspf"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Disable BGP with IPv4 explicitly cleared while retaining OSPF and common route maps")
				},
				Config: testAccMSOL3OutRouteMapsBGPDisabledIPv4ClearConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.enabled", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.enabled", "true"),
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.route_dampening_ipv4_uuid", ""),
					resource.TestCheckNoResourceAttr("mso_l3out.test", "bgp.route_dampening_ipv6_uuid"),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.interleak_uuid",
						"mso_tenant_policies_route_map_policy_route_control.interleak",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.static_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.static",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.connected_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.connected",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.attached_host_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.attached",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.outbound_uuid",
						"mso_tenant_policies_route_map_policy_route_control.outbound",
						"uuid",
					),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.area_id", "0.0.0.10"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.area_type", "regular"),
					testAccCheckL3OutBGPDampeningCleared(&templateID, &l3outUUID),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "ospf"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Restore BGP while retaining OSPF after the mixed dampening clear")
				},
				Config: testAccMSOL3OutRouteMapsBGPWithOSPFConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.enabled", "true"),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"bgp.route_dampening_ipv4_uuid",
						"mso_tenant_policies_route_map_policy_route_control.dampening_v4",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"bgp.route_dampening_ipv6_uuid",
						"mso_tenant_policies_route_map_policy_route_control.dampening_v6",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.interleak_uuid",
						"mso_tenant_policies_route_map_policy_route_control.interleak",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.static_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.static",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.connected_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.connected",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.attached_host_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.attached",
						"uuid",
					),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.enabled", "true"),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.outbound_uuid",
						"mso_tenant_policies_route_map_policy_route_control.outbound",
						"uuid",
					),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "bgpOspf"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Configure OSPF route maps and clear BGP dampening references") },
				Config:    testAccMSOL3OutRouteMapsOSPFConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.enabled", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.route_dampening_ipv4_uuid", ""),
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.route_dampening_ipv6_uuid", ""),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.enabled", "true"),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.outbound_uuid",
						"mso_tenant_policies_route_map_policy_route_control.outbound",
						"uuid",
					),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "ospf"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Configure all route-map references with BGP and OSPF") },
				Config:    testAccMSOL3OutRouteMapsBGPWithOSPFConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.enabled", "true"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.enabled", "true"),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.outbound_uuid",
						"mso_tenant_policies_route_map_policy_route_control.outbound",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"bgp.route_dampening_ipv4_uuid",
						"mso_tenant_policies_route_map_policy_route_control.dampening_v4",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"bgp.route_dampening_ipv6_uuid",
						"mso_tenant_policies_route_map_policy_route_control.dampening_v6",
						"uuid",
					),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "bgpOspf"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import L3Out with all route-map references") },
				Config:            testAccMSOL3OutRouteMapsBGPWithOSPFConfig(siteName, l3outName),
				ResourceName:      "mso_l3out.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Clear interleak while configuring the other route-map references")
				},
				Config: testAccMSOL3OutRouteMapsInterleakClearConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "route_maps.interleak_uuid", ""),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.outbound_uuid",
						"mso_tenant_policies_route_map_policy_route_control.outbound",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.static_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.static",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.connected_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.connected",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.attached_host_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.attached",
						"uuid",
					),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "bgpOspf"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Restore all route-map references before replacing a policy") },
				Config:    testAccMSOL3OutRouteMapsBGPWithOSPFConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.interleak_uuid",
						"mso_tenant_policies_route_map_policy_route_control.interleak",
						"uuid",
					),
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.enabled", "true"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.enabled", "true"),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.outbound_uuid",
						"mso_tenant_policies_route_map_policy_route_control.outbound",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"bgp.route_dampening_ipv4_uuid",
						"mso_tenant_policies_route_map_policy_route_control.dampening_v4",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"bgp.route_dampening_ipv6_uuid",
						"mso_tenant_policies_route_map_policy_route_control.dampening_v6",
						"uuid",
					),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "bgpOspf"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Replace interleak policy while preserving the other route maps and protocols")
				},
				Config: testAccMSOL3OutRouteMapsReplacementConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.interleak_uuid",
						"mso_tenant_policies_route_map_policy_route_control.static",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.outbound_uuid",
						"mso_tenant_policies_route_map_policy_route_control.outbound",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.static_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.static",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.connected_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.connected",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.attached_host_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.attached",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"bgp.route_dampening_ipv4_uuid",
						"mso_tenant_policies_route_map_policy_route_control.dampening_v4",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"bgp.route_dampening_ipv6_uuid",
						"mso_tenant_policies_route_map_policy_route_control.dampening_v6",
						"uuid",
					),
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.enabled", "true"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.enabled", "true"),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "bgpOspf"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Clear interleak while preserving omitted route maps and protocols")
				},
				Config: testAccMSOL3OutRouteMapsPartialInterleakClearConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "route_maps.interleak_uuid", ""),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.outbound_uuid",
						"mso_tenant_policies_route_map_policy_route_control.outbound",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.static_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.static",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.connected_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.connected",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"route_maps.attached_host_route_redistribution_uuid",
						"mso_tenant_policies_route_map_policy_route_control.attached",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"bgp.route_dampening_ipv4_uuid",
						"mso_tenant_policies_route_map_policy_route_control.dampening_v4",
						"uuid",
					),
					resource.TestCheckResourceAttrPair(
						"mso_l3out.test",
						"bgp.route_dampening_ipv6_uuid",
						"mso_tenant_policies_route_map_policy_route_control.dampening_v6",
						"uuid",
					),
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.enabled", "true"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.enabled", "true"),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "bgpOspf"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Clear all route-map references while retaining BGP and OSPF") },
				Config:    testAccMSOL3OutRouteMapsClearConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "route_maps.outbound_uuid", ""),
					resource.TestCheckResourceAttr("mso_l3out.test", "route_maps.interleak_uuid", ""),
					resource.TestCheckResourceAttr("mso_l3out.test", "route_maps.static_route_redistribution_uuid", ""),
					resource.TestCheckResourceAttr(
						"mso_l3out.test",
						"route_maps.connected_route_redistribution_uuid",
						"",
					),
					resource.TestCheckResourceAttr(
						"mso_l3out.test",
						"route_maps.attached_host_route_redistribution_uuid",
						"",
					),
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.route_dampening_ipv4_uuid", ""),
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.route_dampening_ipv6_uuid", ""),
					resource.TestCheckResourceAttr("mso_l3out.test", "bgp.enabled", "true"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.enabled", "true"),
					testAccCheckL3OutRouteMapsCleared(&templateID, &l3outUUID),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "bgpOspf"),
				),
			},
		},
	})
}

func testAccMSOL3OutRouteMapsNoProtocolConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid

  route_maps = {
    interleak_uuid                          = mso_tenant_policies_route_map_policy_route_control.interleak.uuid
    static_route_redistribution_uuid        = mso_tenant_policies_route_map_policy_route_control.static.uuid
    connected_route_redistribution_uuid     = mso_tenant_policies_route_map_policy_route_control.connected.uuid
    attached_host_route_redistribution_uuid = mso_tenant_policies_route_map_policy_route_control.attached.uuid
  }
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName), testAccMSOL3OutRouteMapPolicyPrerequisites())
}

func testAccMSOL3OutRouteMapsBGPConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid

  route_maps = {
    interleak_uuid                          = mso_tenant_policies_route_map_policy_route_control.interleak.uuid
    static_route_redistribution_uuid        = mso_tenant_policies_route_map_policy_route_control.static.uuid
    connected_route_redistribution_uuid     = mso_tenant_policies_route_map_policy_route_control.connected.uuid
    attached_host_route_redistribution_uuid = mso_tenant_policies_route_map_policy_route_control.attached.uuid
  }

  bgp = {
    enabled                   = true
    route_dampening_ipv4_uuid = mso_tenant_policies_route_map_policy_route_control.dampening_v4.uuid
    route_dampening_ipv6_uuid = mso_tenant_policies_route_map_policy_route_control.dampening_v6.uuid
  }
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName), testAccMSOL3OutRouteMapPolicyPrerequisites())
}

func testAccMSOL3OutRouteMapsOSPFConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid

  route_maps = {
    outbound_uuid                          = mso_tenant_policies_route_map_policy_route_control.outbound.uuid
    interleak_uuid                          = mso_tenant_policies_route_map_policy_route_control.interleak.uuid
    static_route_redistribution_uuid        = mso_tenant_policies_route_map_policy_route_control.static.uuid
    connected_route_redistribution_uuid     = mso_tenant_policies_route_map_policy_route_control.connected.uuid
    attached_host_route_redistribution_uuid = mso_tenant_policies_route_map_policy_route_control.attached.uuid
  }

  bgp = {
    enabled                  = false
    route_dampening_ipv4_uuid = ""
    route_dampening_ipv6_uuid = ""
  }

  ospf = {
    enabled   = true
    area_id   = "0.0.0.10"
    area_type = "regular"
  }
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName), testAccMSOL3OutRouteMapPolicyPrerequisites())
}

func testAccMSOL3OutRouteMapsBGPWithOSPFConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid

  route_maps = {
    outbound_uuid                          = mso_tenant_policies_route_map_policy_route_control.outbound.uuid
    interleak_uuid                          = mso_tenant_policies_route_map_policy_route_control.interleak.uuid
    static_route_redistribution_uuid        = mso_tenant_policies_route_map_policy_route_control.static.uuid
    connected_route_redistribution_uuid     = mso_tenant_policies_route_map_policy_route_control.connected.uuid
    attached_host_route_redistribution_uuid = mso_tenant_policies_route_map_policy_route_control.attached.uuid
  }

  bgp = {
    enabled                   = true
    route_dampening_ipv4_uuid = mso_tenant_policies_route_map_policy_route_control.dampening_v4.uuid
    route_dampening_ipv6_uuid = mso_tenant_policies_route_map_policy_route_control.dampening_v6.uuid
  }

  ospf = {
    enabled   = true
    area_id   = "0.0.0.10"
    area_type = "regular"
  }
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName), testAccMSOL3OutRouteMapPolicyPrerequisites())
}

func testAccMSOL3OutRouteMapsInterleakClearConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid

  route_maps = {
    outbound_uuid                          = mso_tenant_policies_route_map_policy_route_control.outbound.uuid
    interleak_uuid                          = ""
    static_route_redistribution_uuid        = mso_tenant_policies_route_map_policy_route_control.static.uuid
    connected_route_redistribution_uuid     = mso_tenant_policies_route_map_policy_route_control.connected.uuid
    attached_host_route_redistribution_uuid = mso_tenant_policies_route_map_policy_route_control.attached.uuid
  }

  bgp = {
    enabled                   = true
    route_dampening_ipv4_uuid = mso_tenant_policies_route_map_policy_route_control.dampening_v4.uuid
    route_dampening_ipv6_uuid = mso_tenant_policies_route_map_policy_route_control.dampening_v6.uuid
  }

  ospf = {
    enabled   = true
    area_id   = "0.0.0.10"
    area_type = "regular"
  }
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName), testAccMSOL3OutRouteMapPolicyPrerequisites())
}

func testAccMSOL3OutRouteMapsBGPEmptyConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid

  depends_on = [
    mso_tenant_policies_route_map_policy_route_control.outbound,
    mso_tenant_policies_route_map_policy_route_control.interleak,
    mso_tenant_policies_route_map_policy_route_control.static,
    mso_tenant_policies_route_map_policy_route_control.connected,
    mso_tenant_policies_route_map_policy_route_control.attached,
    mso_tenant_policies_route_map_policy_route_control.dampening_v4,
    mso_tenant_policies_route_map_policy_route_control.dampening_v6,
  ]

  bgp = {}
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName), testAccMSOL3OutRouteMapPolicyPrerequisites())
}

func testAccMSOL3OutRouteMapsBGPDisabledConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid

  depends_on = [
    mso_tenant_policies_route_map_policy_route_control.outbound,
    mso_tenant_policies_route_map_policy_route_control.interleak,
    mso_tenant_policies_route_map_policy_route_control.static,
    mso_tenant_policies_route_map_policy_route_control.connected,
    mso_tenant_policies_route_map_policy_route_control.attached,
    mso_tenant_policies_route_map_policy_route_control.dampening_v4,
    mso_tenant_policies_route_map_policy_route_control.dampening_v6,
  ]

  bgp = {
    enabled = false
  }
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName), testAccMSOL3OutRouteMapPolicyPrerequisites())
}

func testAccMSOL3OutRouteMapsBGPDisabledIPv4ClearConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid

  depends_on = [
    mso_tenant_policies_route_map_policy_route_control.outbound,
    mso_tenant_policies_route_map_policy_route_control.interleak,
    mso_tenant_policies_route_map_policy_route_control.static,
    mso_tenant_policies_route_map_policy_route_control.connected,
    mso_tenant_policies_route_map_policy_route_control.attached,
    mso_tenant_policies_route_map_policy_route_control.dampening_v4,
    mso_tenant_policies_route_map_policy_route_control.dampening_v6,
  ]

  bgp = {
    enabled                   = false
    route_dampening_ipv4_uuid = ""
  }
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName), testAccMSOL3OutRouteMapPolicyPrerequisites())
}

func testAccMSOL3OutRouteMapsPartialInterleakClearConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid

  depends_on = [
    mso_tenant_policies_route_map_policy_route_control.outbound,
    mso_tenant_policies_route_map_policy_route_control.interleak,
    mso_tenant_policies_route_map_policy_route_control.static,
    mso_tenant_policies_route_map_policy_route_control.connected,
    mso_tenant_policies_route_map_policy_route_control.attached,
    mso_tenant_policies_route_map_policy_route_control.dampening_v4,
    mso_tenant_policies_route_map_policy_route_control.dampening_v6,
  ]

  route_maps = {
    interleak_uuid = ""
  }
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName), testAccMSOL3OutRouteMapPolicyPrerequisites())
}

func testAccMSOL3OutBGPDampeningIPv4ClearConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid

  depends_on = [
    mso_tenant_policies_route_map_policy_route_control.outbound,
    mso_tenant_policies_route_map_policy_route_control.interleak,
    mso_tenant_policies_route_map_policy_route_control.static,
    mso_tenant_policies_route_map_policy_route_control.connected,
    mso_tenant_policies_route_map_policy_route_control.attached,
    mso_tenant_policies_route_map_policy_route_control.dampening_v4,
    mso_tenant_policies_route_map_policy_route_control.dampening_v6,
  ]

  bgp = {
    enabled                   = true
    route_dampening_ipv4_uuid = ""
  }
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName), testAccMSOL3OutRouteMapPolicyPrerequisites())
}

func testAccMSOL3OutRouteMapsReplacementConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid

  route_maps = {
    interleak_uuid = mso_tenant_policies_route_map_policy_route_control.static.uuid
  }
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName), testAccMSOL3OutRouteMapPolicyPrerequisites())
}

func testAccMSOL3OutRouteMapsClearConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid

  depends_on = [
    mso_tenant_policies_route_map_policy_route_control.outbound,
    mso_tenant_policies_route_map_policy_route_control.interleak,
    mso_tenant_policies_route_map_policy_route_control.static,
    mso_tenant_policies_route_map_policy_route_control.connected,
    mso_tenant_policies_route_map_policy_route_control.attached,
    mso_tenant_policies_route_map_policy_route_control.dampening_v4,
    mso_tenant_policies_route_map_policy_route_control.dampening_v6,
  ]

  route_maps = {
    outbound_uuid                           = ""
    interleak_uuid                          = ""
    static_route_redistribution_uuid        = ""
    connected_route_redistribution_uuid     = ""
    attached_host_route_redistribution_uuid = ""
  }

  bgp = {
    enabled                   = true
    route_dampening_ipv4_uuid = ""
    route_dampening_ipv6_uuid = ""
  }
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName), testAccMSOL3OutRouteMapPolicyPrerequisites())
}

func testAccMSOL3OutRouteMapPolicyPrerequisites() string {
	var config strings.Builder
	config.WriteString(`
resource "mso_template" "l3out_test_route_map_policies" {
  template_name = "${local.l3out_name}_route_map_policies"
  template_type = "tenant"
  tenant_id     = mso_tenant.l3out_test.id
  sites         = [data.mso_site.l3out_test.id]
}
`)
	names := []string{"outbound", "interleak", "static", "connected", "attached", "dampening_v4", "dampening_v6"}
	for index, name := range names {
		dependency := ""
		if index > 0 {
			dependency = fmt.Sprintf(
				"\n  depends_on = [mso_tenant_policies_route_map_policy_route_control.%s]\n",
				names[index-1],
			)
		}
		fmt.Fprintf(&config, `
resource "mso_tenant_policies_route_map_policy_route_control" "%s" {
  template_id = mso_template.l3out_test_route_map_policies.id
  name        = "${local.l3out_name}_%s"
%s
}
`, name, name, dependency)
	}
	return config.String()
}

func testAccMSOL3OutOSPFEmptyConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  ospf        = {}
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutOSPFDisabledConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  ospf = {
    enabled = false
  }
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutBGPConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  bgp = {
    enabled = true
  }
  originate_default_route = "in_addition"
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutBGPOmittedConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id             = mso_template.l3out_test.id
  name                    = local.l3out_name
  vrf_uuid                = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  originate_default_route = "in_addition"
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutBGPOmittedUpdateConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id             = mso_template.l3out_test.id
  name                    = local.l3out_name
  vrf_uuid                = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description             = "BGP omission update"
  originate_default_route = "in_addition"
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutBGPWithOSPFConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  bgp = {
    enabled = true
  }
  ospf = {
    enabled   = true
    area_id   = "0.0.0.3"
    area_type = "regular"
  }
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutBGPDisabledWithOSPFConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  bgp = {
    enabled = false
  }
  ospf = {
    enabled   = true
    area_id   = "0.0.0.3"
    area_type = "regular"
  }
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutBothProtocolsEmptyConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id             = mso_template.l3out_test.id
  name                    = local.l3out_name
  vrf_uuid                = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  bgp                     = {}
  ospf                    = {}
  originate_default_route = ""
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutBothProtocolsOmittedAfterEmptyConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id             = mso_template.l3out_test.id
  name                    = local.l3out_name
  vrf_uuid                = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description             = "protocols omitted after empty"
  originate_default_route = ""
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutBothProtocolsDisabledConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id             = mso_template.l3out_test.id
  name                    = local.l3out_name
  vrf_uuid                = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description             = "protocols omitted after empty"
  originate_default_route = ""
  bgp = {
    enabled = false
  }
  ospf = {
    enabled = false
  }
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutBothProtocolsOmittedAfterFalseConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id             = mso_template.l3out_test.id
  name                    = local.l3out_name
  vrf_uuid                = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description             = "protocols omitted after false"
  originate_default_route = ""
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutOSPFMinimalConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  ospf = {
    enabled   = true
    area_id   = "0.0.0.1"
    area_type = "regular"
  }
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutOSPFMinimalUpdateConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description = "OSPF default settings update"
  ospf = {
    enabled   = true
    area_id   = "0.0.0.1"
    area_type = "regular"
  }
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutOSPFEnabledConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  ospf = {
    enabled = true
    area_id = "0.0.0.1"
    area_type = "regular"
    cost = 10
    send_redistributed_lsa = true
    originate_summary_lsa = true
    suppress_forwarding_address_in_translated_lsa = true
    originate_default_route_always = false
  }
  originate_default_route = "only"
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutOSPFUpdatedConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  ospf = {
    enabled = true
    area_id = "0.0.0.2"
    area_type = "nssa"
    cost = 20
    send_redistributed_lsa = false
    originate_summary_lsa = false
    suppress_forwarding_address_in_translated_lsa = false
    originate_default_route_always = true
  }
  originate_default_route = "in_addition"
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutDefaultRouteEmptyConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  ospf = {
    enabled = true
    area_id = "0.0.0.2"
    area_type = "nssa"
    cost = 20
    send_redistributed_lsa = false
    originate_summary_lsa = false
    suppress_forwarding_address_in_translated_lsa = false
  }
  originate_default_route = ""
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutOSPFOmittedConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutOSPFOmittedUpdateConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description = "OSPF omission update"
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccCheckL3OutRoutingProtocol(templateID, l3outUUID *string, expected string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		template, err := ndoapi.GetTemplate(testAccAPIClient(), *templateID)
		if err != nil {
			return err
		}
		resolved, found, err := template.Find(models.NewL3OutPath(*l3outUUID, ""))
		if err != nil {
			return err
		}
		if !found {
			return fmt.Errorf("L3Out %q was not found in template %q", *l3outUUID, *templateID)
		}
		actual, _, err := ndoapi.StringField(resolved.Object, "routingProtocol", ndoapi.OptionalField)
		if err != nil {
			return err
		}
		if actual != expected {
			return fmt.Errorf("expected L3Out routingProtocol %q, got %q", expected, actual)
		}
		return nil
	}
}

func testAccCheckL3OutDeleted(templateID, l3outUUID *string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		// Require the original template to remain readable so its deletion cannot mask an orphaned L3Out.
		template, err := ndoapi.GetTemplate(testAccAPIClient(), *templateID)
		if err != nil {
			return err
		}
		_, found, err := template.Find(models.NewL3OutPath(*l3outUUID, ""))
		if err != nil {
			return err
		}
		if found {
			return fmt.Errorf("L3Out %q still exists in template %q after deletion", *l3outUUID, *templateID)
		}
		return nil
	}
}

func testAccCheckL3OutBGPDampeningCleared(templateID, l3outUUID *string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		template, err := ndoapi.GetTemplate(testAccAPIClient(), *templateID)
		if err != nil {
			return err
		}
		resolved, found, err := template.Find(models.NewL3OutPath(*l3outUUID, ""))
		if err != nil {
			return err
		}
		if !found {
			return fmt.Errorf("L3Out %q was not found in template %q", *l3outUUID, *templateID)
		}
		advanced, _, err := ndoapi.MapField(resolved.Object, "advancedRouteMapRefs", ndoapi.OptionalField)
		if err != nil {
			return err
		}
		l3out, exists := state.RootModule().Resources["mso_l3out.test"]
		if !exists {
			return fmt.Errorf("mso_l3out.test was not found in Terraform state")
		}
		outbound, _, err := ndoapi.StringField(resolved.Object, "exportRouteMapRef", ndoapi.OptionalField)
		if err != nil {
			return err
		}
		if outbound != l3out.Primary.Attributes["route_maps.outbound_uuid"] {
			return fmt.Errorf("L3Out %q did not preserve its outbound route map: %s", *l3outUUID, outbound)
		}
		for _, field := range []string{"routeDampeningV4Ref", "routeDampeningV6Ref"} {
			if value, exists := advanced[field]; exists && value != nil {
				return fmt.Errorf("L3Out %q still has %s after BGP was disabled: %v", *l3outUUID, field, value)
			}
		}
		for _, reference := range []struct{ field, policy string }{
			{"interleakRef", "interleak"},
			{"staticRouteRedistRef", "static"},
			{"connectedRouteRedistRef", "connected"},
			{"attachedHostRouteRedistRef", "attached"},
		} {
			resourceName := fmt.Sprintf("mso_tenant_policies_route_map_policy_route_control.%s", reference.policy)
			policy, exists := state.RootModule().Resources[resourceName]
			if !exists {
				return fmt.Errorf("%s was not found in Terraform state", resourceName)
			}
			if advanced[reference.field] != policy.Primary.Attributes["uuid"] {
				return fmt.Errorf("L3Out %q did not preserve %s: %v", *l3outUUID, reference.field, advanced[reference.field])
			}
		}
		return nil
	}
}

func testAccCheckL3OutRouteMapsCleared(templateID, l3outUUID *string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		template, err := ndoapi.GetTemplate(testAccAPIClient(), *templateID)
		if err != nil {
			return err
		}
		resolved, found, err := template.Find(models.NewL3OutPath(*l3outUUID, ""))
		if err != nil {
			return err
		}
		if !found {
			return fmt.Errorf("L3Out %q was not found in template %q", *l3outUUID, *templateID)
		}
		for _, field := range []string{"exportRouteMapRef", "advancedRouteMapRefs"} {
			if value, exists := resolved.Object[field]; exists && value != nil {
				return fmt.Errorf(
					"L3Out %q still has %s after clearing all route-map references: %v",
					*l3outUUID,
					field,
					value,
				)
			}
		}
		return nil
	}
}

// testAccCaptureResourceIdentifiers copies identifiers from Terraform state into
// Go variables for use by later test callbacks. Unlike TestCheckResourceAttrSet,
// which only validates an attribute during the current step, this preserves the
// values for operations such as out-of-band deletion and replacement checks.
func testAccCaptureResourceIdentifiers(resourceName string, templateID, resourceID *string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		resourceState, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("%s not found in Terraform state", resourceName)
		}

		if templateID != nil {
			*templateID = resourceState.Primary.Attributes["template_id"]
		}
		if resourceID != nil {
			*resourceID = resourceState.Primary.Attributes["uuid"]
		}
		if (templateID != nil && *templateID == "") || (resourceID != nil && *resourceID == "") {
			return fmt.Errorf("%s did not have the requested identifiers in state", resourceName)
		}
		return nil
	}
}

func testAccCheckL3OutCompositeID(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		resourceState, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("%s not found in Terraform state", resourceName)
		}
		attributes := resourceState.Primary.Attributes
		expected := fmt.Sprintf("%s/%s", attributes["template_id"], attributes["uuid"])
		if attributes["id"] != expected {
			return fmt.Errorf("unexpected L3Out ID: got %q, want %q", attributes["id"], expected)
		}
		return nil
	}
}

// testAccAssertResourceIDChanged verifies that a replacement produced a new
// remote resource ID rather than updating the existing resource in place.
func testAccAssertResourceIDChanged(resourceName string, previousID *string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		resourceState, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("%s not found in Terraform state", resourceName)
		}
		currentID := resourceState.Primary.Attributes["uuid"]
		if currentID == "" {
			return fmt.Errorf("%s did not have a uuid in state", resourceName)
		}
		if *previousID == "" {
			return fmt.Errorf("no previous uuid was captured for %s", resourceName)
		}
		if currentID == *previousID {
			return fmt.Errorf("%s uuid did not change after template_id replacement", resourceName)
		}
		return nil
	}
}

// testAccDeleteL3Out removes the L3Out directly through NDO to simulate
// out-of-band deletion before Terraform refreshes the resource.
func testAccDeleteL3Out(t *testing.T, templateID, l3outUUID string) {
	t.Helper()
	if templateID == "" || l3outUUID == "" {
		t.Fatal("cannot delete L3Out out of band without template_id and uuid")
	}

	msoClient := testAccAPIClient()
	template, err := ndoapi.GetTemplate(msoClient, templateID)
	if err != nil {
		t.Fatalf("unable to read L3Out template for out-of-band deletion: %v", err)
	}

	resolved, found, err := template.Find(models.NewL3OutPath(l3outUUID, ""))
	if err != nil {
		t.Fatalf("unable to resolve L3Out for out-of-band deletion: %v", err)
	}
	if !found {
		t.Fatalf("L3Out %q was not found for out-of-band deletion", l3outUUID)
	}

	if _, err := msoClient.PatchbyID(
		ndoapi.TemplateEndpoint(templateID),
		ndoapi.NewRemovePatchPayload(resolved.PatchPath()),
	); err != nil {
		t.Fatalf("unable to delete L3Out out of band: %v", err)
	}
}

func TestAccMSOL3OutInvalid(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_l3out_invalid")
	defaultRouteValidationMessage := `originate_default_route must be one of "only", "in_addition" when ospf.originate_default_route_always is true`
	defaultRouteValidationError := regexp.MustCompile(
		strings.ReplaceAll(regexp.QuoteMeta(defaultRouteValidationMessage), " ", `\s+`),
	)
	dampeningValidationMessage := fmt.Sprintf(
		"%s, %s may be set only when %s is one of true",
		"bgp.route_dampening_ipv4_uuid",
		"bgp.route_dampening_ipv6_uuid",
		"bgp.enabled",
	)
	dampeningValidationError := regexp.MustCompile(
		strings.ReplaceAll(regexp.QuoteMeta(dampeningValidationMessage), " ", `\s+`),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccProviderPreCheck(t)
			testAccVersionCheck(t, "4.3")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: Reject OSPF always-originate when the default-route mode is omitted")
				},
				Config:      testAccMSOL3OutDefaultRouteModeMissingConfig(siteName, l3outName),
				ExpectError: defaultRouteValidationError,
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Reject OSPF always-originate with an empty default-route mode")
				},
				Config:      testAccMSOL3OutDefaultRouteModeEmptyConfig(siteName, l3outName),
				ExpectError: defaultRouteValidationError,
			},
			{
				PreConfig: func() { fmt.Println("Test: Reject OSPF settings while disabled") },
				Config:    testAccMSOL3OutOSPFDisabledWithAreaConfig(siteName, l3outName),
				ExpectError: regexp.MustCompile(
					`ospf\.area_id\s+requires\s+ospf\.enabled\s+=\s+true\.`,
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Reject non-empty dampening references when BGP is explicitly disabled")
				},
				Config:      testAccMSOL3OutBGPDisabledWithDampeningConfig(),
				ExpectError: dampeningValidationError,
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Reject non-empty dampening references when BGP enabled defaults to false")
				},
				Config:      testAccMSOL3OutBGPDefaultDisabledWithDampeningConfig(),
				ExpectError: dampeningValidationError,
			},
			{
				PreConfig: func() { fmt.Println("Test: Reject enabled OSPF without an area ID") },
				Config:    testAccMSOL3OutInvalidOSPFWithoutAreaIDConfig(),
				ExpectError: regexp.MustCompile(
					`ospf\.area_id\s+must\s+be\s+configured\s+when\s+ospf\.enabled\s+is\s+true\.`,
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Reject enabled OSPF without an area type") },
				Config:    testAccMSOL3OutInvalidOSPFWithoutAreaTypeConfig(),
				ExpectError: regexp.MustCompile(
					`ospf\.area_type\s+must\s+be\s+configured\s+when\s+ospf\.enabled\s+is\s+true\.`,
				),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject L3Out creation in a nonexistent template") },
				Config:      testAccMSOL3OutInvalidTemplateCreateConfig(),
				ExpectError: regexp.MustCompile(`Failed to Create L3Out`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create L3Out for invalid import reference") },
				Config:    testAccMSOL3OutInvalidImportParentConfig(siteName, l3outName),
			},
			{
				PreConfig:     func() { fmt.Println("Test: Reject import when the L3Out template does not exist") },
				ResourceName:  "mso_l3out.test",
				ImportState:   true,
				ImportStateId: "00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000",
				ExpectError: regexp.MustCompile(
					`(?i)(non-existent|not found|does not exist)`,
				),
			},
		},
	})
}

func testAccMSOL3OutDefaultRouteModeMissingConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  ospf = {
    enabled                         = true
    area_id                         = "0.0.0.1"
    area_type                       = "regular"
    originate_default_route_always = true
  }
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutDefaultRouteModeEmptyConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  ospf = {
    enabled                         = true
    area_id                         = "0.0.0.1"
    area_type                       = "regular"
    originate_default_route_always = true
  }
  originate_default_route = ""
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutOSPFDisabledWithAreaConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  ospf = {
    enabled = false
    area_id = "0.0.0.1"
  }
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutBGPDisabledWithDampeningConfig() string {
	return `
resource "mso_l3out" "test" {
  template_id = "00000000-0000-0000-0000-000000000000"
  name        = "terraform_invalid_bgp_dampening"
  vrf_uuid    = "00000000-0000-0000-0000-000000000000"
  bgp = {
    enabled                   = false
    route_dampening_ipv4_uuid = "00000000-0000-0000-0000-000000000001"
    route_dampening_ipv6_uuid = "00000000-0000-0000-0000-000000000002"
  }
}
`
}

func testAccMSOL3OutBGPDefaultDisabledWithDampeningConfig() string {
	return `
resource "mso_l3out" "test" {
  template_id = "00000000-0000-0000-0000-000000000000"
  name        = "terraform_invalid_bgp_dampening_default"
  vrf_uuid    = "00000000-0000-0000-0000-000000000000"
  bgp = {
    route_dampening_ipv6_uuid = "00000000-0000-0000-0000-000000000002"
  }
}
`
}

func testAccMSOL3OutInvalidOSPFWithoutAreaIDConfig() string {
	return `
resource "mso_l3out" "test" {
  template_id = "00000000-0000-0000-0000-000000000000"
  name        = "terraform_invalid_ospf"
  vrf_uuid    = "00000000-0000-0000-0000-000000000000"
  ospf = {
    enabled   = true
    area_type = "regular"
  }
}
`
}

func testAccMSOL3OutInvalidOSPFWithoutAreaTypeConfig() string {
	return `
resource "mso_l3out" "test" {
  template_id = "00000000-0000-0000-0000-000000000000"
  name        = "terraform_invalid_ospf"
  vrf_uuid    = "00000000-0000-0000-0000-000000000000"
  ospf = {
    enabled = true
    area_id = "0.0.0.1"
  }
}
`
}

func testAccMSOL3OutInvalidTemplateCreateConfig() string {
	return `
resource "mso_l3out" "test" {
  template_id = "00000000-0000-0000-0000-000000000000"
  name        = "terraform_missing_template"
  vrf_uuid    = "00000000-0000-0000-0000-000000000000"
}
`
}

func testAccMSOL3OutInvalidImportParentConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutPrerequisites(siteName, l3outName string) string {
	return fmt.Sprintf(`
locals {
	 l3out_name = %q
}

data "mso_site" "l3out_test" {
  name = %q
}

resource "mso_tenant" "l3out_test" {
  name = "${local.l3out_name}_tenant"

  site_associations {
    site_id = data.mso_site.l3out_test.id
  }
}

resource "mso_schema" "l3out_test" {
  name = "${local.l3out_name}_schema"

  template {
    name         = "${local.l3out_name}_schema_template"
    display_name = "${local.l3out_name}_schema_template"
    tenant_id    = mso_tenant.l3out_test.id
  }
}

resource "mso_schema_site" "l3out_test" {
  schema_id     = mso_schema.l3out_test.id
  site_id       = data.mso_site.l3out_test.id
  template_name = "${local.l3out_name}_schema_template"
}

resource "mso_schema_template_vrf" "l3out_test_vrf_1" {
  name         = "${local.l3out_name}_vrf_1"
  display_name = "${local.l3out_name}_vrf_1"
  schema_id    = mso_schema.l3out_test.id
  template     = "${local.l3out_name}_schema_template"
  layer3_multicast = true

  depends_on = [mso_schema_site.l3out_test]
}

resource "mso_schema_template_vrf" "l3out_test_vrf_2" {
  name         = "${local.l3out_name}_vrf_2"
  display_name = "${local.l3out_name}_vrf_2"
  schema_id    = mso_schema.l3out_test.id
  template     = "${local.l3out_name}_schema_template"
  layer3_multicast = true

  depends_on = [mso_schema_site.l3out_test]
}

resource "mso_schema_template_deploy_ndo" "l3out_test_vrf" {
  schema_id           = mso_schema.l3out_test.id
  template_name       = "${local.l3out_name}_schema_template"
  force_apply         = ""
  undeploy_on_destroy = true

  depends_on = [
    mso_schema_template_vrf.l3out_test_vrf_1,
    mso_schema_template_vrf.l3out_test_vrf_2,
  ]
}

resource "mso_template" "l3out_fabric_policy" {
  template_name = "${local.l3out_name}_fabric_policy"
  template_type = "fabric_policy"
  sites         = [data.mso_site.l3out_test.id]
}

resource "mso_fabric_policies_l3_domain" "l3out_test_domain_1" {
  template_id = mso_template.l3out_fabric_policy.id
  name        = "${local.l3out_name}_l3_domain_1"
}

resource "mso_fabric_policies_l3_domain" "l3out_test_domain_2" {
  template_id = mso_template.l3out_fabric_policy.id
  name        = "${local.l3out_name}_l3_domain_2"
  depends_on  = [mso_fabric_policies_l3_domain.l3out_test_domain_1]
}

resource "mso_template" "l3out_test" {
  template_name = "${local.l3out_name}_template"
  template_type = "l3out"
  tenant_id     = mso_tenant.l3out_test.id
  sites         = [data.mso_site.l3out_test.id]

  depends_on = [mso_schema_template_deploy_ndo.l3out_test_vrf]
}
	`, l3outName, siteName,
	)
}

func testAccCheckL3OutAnnotations(templateID, l3outUUID *string, expected map[string]string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if *templateID == "" || *l3outUUID == "" {
			return fmt.Errorf("L3Out template ID or UUID has not been captured")
		}
		template, err := ndoapi.GetTemplate(testAccAPIClient(), *templateID)
		if err != nil {
			return err
		}
		resolved, found, err := template.Find(models.NewL3OutPath(*l3outUUID, ""))
		if err != nil {
			return err
		}
		if !found {
			return fmt.Errorf("L3Out %q not found in template %q", *l3outUUID, *templateID)
		}
		var actual models.AnnotationsModel
		if err := actual.SetFromNDOObject(resolved.Object); err != nil {
			return err
		}
		if !maps.Equal(map[string]string(actual), expected) {
			return fmt.Errorf("unexpected remote annotations: got %v, want %v", actual, expected)
		}
		return nil
	}
}

func testAccReplaceL3OutAnnotations(t *testing.T, templateID, l3outUUID string, annotations models.AnnotationsModel) {
	t.Helper()
	template, err := ndoapi.GetTemplate(testAccAPIClient(), templateID)
	if err != nil {
		t.Fatalf("unable to read L3Out template: %v", err)
	}
	resolved, found, err := template.Find(models.NewL3OutPath(l3outUUID, ""))
	if err != nil {
		t.Fatalf("unable to resolve L3Out for annotation replacement: %v", err)
	}
	if !found {
		t.Fatalf("L3Out %q not found in template %q", l3outUUID, templateID)
	}
	_, err = testAccAPIClient().PatchbyID(
		ndoapi.TemplateEndpoint(templateID),
		ndoapi.NewPatchPayload(
			"replace",
			fmt.Sprintf("%s/tagAnnotations", resolved.PatchPath()),
			annotations.ToPayload(),
		),
	)
	if err != nil {
		t.Fatalf("unable to replace L3Out annotations out of band: %v", err)
	}
}
func TestAccMSOL3OutAnnotations(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_l3out_annotations")
	var templateID, l3outUUID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccProviderPreCheck(t)
			testAccVersionCheck(t, "4.3")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create L3Out with annotations sharing a value") },
				Config:    testAccMSOL3OutAnnotationsCreateConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCaptureResourceIdentifiers("mso_l3out.test", &templateID, &l3outUUID),
					resource.TestCheckResourceAttr("mso_l3out.test", "annotations.%", "2"),
					resource.TestCheckResourceAttr("mso_l3out.test", "annotations.foo", "shared"),
					resource.TestCheckResourceAttr("mso_l3out.test", "annotations.bar", "shared"),
					testAccCheckL3OutAnnotations(
						&templateID,
						&l3outUUID,
						map[string]string{"foo": "shared", "bar": "shared"},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Reorder L3Out annotation keys without changing the plan") },
				Config:    testAccMSOL3OutAnnotationsReorderedConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "annotations.%", "2"),
					testAccCheckL3OutAnnotations(
						&templateID,
						&l3outUUID,
						map[string]string{"foo": "shared", "bar": "shared"},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update and add L3Out annotations") },
				Config:    testAccMSOL3OutAnnotationsUpdateConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "annotations.%", "3"),
					resource.TestCheckResourceAttr("mso_l3out.test", "annotations.foo", "updated"),
					resource.TestCheckResourceAttr("mso_l3out.test", "annotations.bar", "shared"),
					resource.TestCheckResourceAttr("mso_l3out.test", "annotations.baz", "shared"),
					testAccCheckL3OutAnnotations(
						&templateID,
						&l3outUUID,
						map[string]string{"foo": "updated", "bar": "shared", "baz": "shared"},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove one L3Out annotation") },
				Config:    testAccMSOL3OutAnnotationsRemoveOneConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "annotations.%", "2"),
					resource.TestCheckResourceAttr("mso_l3out.test", "annotations.foo", "updated"),
					resource.TestCheckResourceAttr("mso_l3out.test", "annotations.baz", "shared"),
					testAccCheckL3OutAnnotations(
						&templateID,
						&l3outUUID,
						map[string]string{"foo": "updated", "baz": "shared"},
					),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import L3Out with annotations") },
				Config:       testAccMSOL3OutAnnotationsRemoveOneConfig(siteName, l3outName),
				ResourceName: "mso_l3out.test",
				ImportState:  true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					resourceState, ok := state.RootModule().Resources["mso_l3out.test"]
					if !ok {
						return "", fmt.Errorf("mso_l3out.test not found in Terraform state")
					}
					return fmt.Sprintf(
						"%s/%s",
						resourceState.Primary.Attributes["template_id"],
						resourceState.Primary.Attributes["uuid"],
					), nil
				},
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Omit L3Out annotations and preserve remote values")
					testAccReplaceL3OutAnnotations(t, templateID, l3outUUID, models.AnnotationsModel{
						"foo": "updated", "baz": "shared", "external": "shared",
					})
				},
				Config: testAccMSOL3OutAnnotationsOmittedConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "annotations.%", "3"),
					resource.TestCheckResourceAttr("mso_l3out.test", "annotations.foo", "updated"),
					resource.TestCheckResourceAttr("mso_l3out.test", "annotations.baz", "shared"),
					resource.TestCheckResourceAttr("mso_l3out.test", "annotations.external", "shared"),
					testAccCheckL3OutAnnotations(
						&templateID,
						&l3outUUID,
						map[string]string{"foo": "updated", "baz": "shared", "external": "shared"},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove all L3Out annotations with an empty map") },
				Config:    testAccMSOL3OutAnnotationsClearConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "annotations.%", "0"),
					testAccCheckL3OutAnnotations(&templateID, &l3outUUID, map[string]string{}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Populate annotations from a computed map") },
				Config: testAccMSOL3OutAnnotationsComputedConfig(siteName, l3outName, `{
    owner = "computed-network"
    purpose = "external-routing"
  }`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("mso_l3out.test", tfjsonpath.New("annotations")),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: testAccCheckL3OutAnnotations(
					&templateID,
					&l3outUUID,
					map[string]string{"owner": "computed-network", "purpose": "external-routing"},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update and remove annotations through a computed map") },
				Config: testAccMSOL3OutAnnotationsComputedConfig(siteName, l3outName, `{
    owner = "computed-operations"
  }`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("mso_l3out.test", tfjsonpath.New("annotations")),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: testAccCheckL3OutAnnotations(
					&templateID,
					&l3outUUID,
					map[string]string{"owner": "computed-operations"},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Clear annotations through a computed empty map") },
				Config:    testAccMSOL3OutAnnotationsComputedConfig(siteName, l3outName, "{}"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("mso_l3out.test", tfjsonpath.New("annotations")),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: testAccCheckL3OutAnnotations(&templateID, &l3outUUID, map[string]string{}),
			},
		},
	})
}

func testAccMSOL3OutAnnotationsCreateConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description = "initial annotations"
  annotations = {
    foo = "shared"
    bar = "shared"
  }
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutAnnotationsReorderedConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description = "initial annotations"
  annotations = {
    bar = "shared"
    foo = "shared"
  }
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutAnnotationsUpdateConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description = "updated annotations"
  annotations = {
    bar = "shared"
    foo = "updated"
    baz = "shared"
  }
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutAnnotationsRemoveOneConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description = "updated annotations"
  annotations = {
    foo = "updated"
    baz = "shared"
  }
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutAnnotationsOmittedConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description = "annotations omitted"
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutAnnotationsClearConfig(siteName, l3outName string) string {
	return fmt.Sprintf(`%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description = "annotations cleared"
  annotations = {}
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName))
}

func testAccMSOL3OutAnnotationsComputedConfig(siteName, l3outName, annotations string) string {
	return fmt.Sprintf(`%s
resource "terraform_data" "annotations" {
  input = %s
}

resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description = "computed annotations"
  annotations = terraform_data.annotations.output
}
`, testAccMSOL3OutPrerequisites(siteName, l3outName), annotations)
}
