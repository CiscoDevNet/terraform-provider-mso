package provider

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-mso/internal/models"
	ndoapi "github.com/CiscoDevNet/terraform-provider-mso/internal/ndoapi"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccMSOL3OutResource(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_l3out")
	var templateID string
	var l3outUUID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccProviderPreCheck(t) },
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
		PreCheck:                 func() { testAccProviderPreCheck(t) },
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

func TestAccMSOL3OutRoutingProtocol(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_l3out_routing_protocol")
	var templateID, l3outUUID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccProviderPreCheck(t) },
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
		PreCheck:                 func() { testAccProviderPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Configure common route maps without a routing protocol") },
				Config:    testAccMSOL3OutRouteMapsResourceConfig(siteName, l3outName, "none"),
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
				Config: testAccMSOL3OutRouteMapsResourceConfig(siteName, l3outName, "bgp"),
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
				Config: testAccMSOL3OutRouteMapsResourceConfig(siteName, l3outName, "bgp"),
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
				PreConfig: func() { fmt.Println("Test: Configure OSPF route maps and clear BGP dampening references") },
				Config:    testAccMSOL3OutRouteMapsResourceConfig(siteName, l3outName, "ospf"),
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
				Config:    testAccMSOL3OutRouteMapsResourceConfig(siteName, l3outName, "both"),
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
				Config:            testAccMSOL3OutRouteMapsResourceConfig(siteName, l3outName, "both"),
				ResourceName:      "mso_l3out.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Clear interleak while configuring the other route-map references")
				},
				Config: testAccMSOL3OutRouteMapsResourceConfig(siteName, l3outName, "clear_interleak"),
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
				Config:    testAccMSOL3OutRouteMapsResourceConfig(siteName, l3outName, "both"),
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

func testAccPolicyUUID(name string) string {
	return fmt.Sprintf("mso_tenant_policies_route_map_policy_route_control.%s.uuid", name)
}

func testAccMSOL3OutRouteMapsResourceConfig(siteName, l3outName, protocol string) string {
	bgpConfig := ""
	ospfConfig := ""
	includeOutbound := protocol == "ospf" || protocol == "both" || protocol == "clear_interleak"
	switch protocol {
	case "bgp", "both", "clear_interleak":
		bgpConfig = fmt.Sprintf(`
  bgp = {
    enabled                = true
    route_dampening_ipv4_uuid = %s
    route_dampening_ipv6_uuid = %s
  }
`, testAccPolicyUUID("dampening_v4"), testAccPolicyUUID("dampening_v6"))
	case "ospf":
		bgpConfig = `
  bgp = {
    enabled                  = false
    route_dampening_ipv4_uuid = ""
    route_dampening_ipv6_uuid = ""
  }
`
	}
	if protocol == "ospf" || protocol == "both" || protocol == "clear_interleak" {
		ospfConfig = `
  ospf = {
    enabled   = true
    area_id   = "0.0.0.10"
    area_type = "regular"
  }
`
	}
	outbound := ""
	if includeOutbound {
		outbound = fmt.Sprintf("    outbound_uuid = %s\n", testAccPolicyUUID("outbound"))
	}
	interleak := fmt.Sprintf("    interleak_uuid = %s\n", testAccPolicyUUID("interleak"))
	if protocol == "clear_interleak" {
		interleak = "    interleak_uuid = \"\"\n"
	}

	return fmt.Sprintf(`%s%s
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid

  route_maps = {
%s%s    static_route_redistribution_uuid       = %s
    connected_route_redistribution_uuid    = %s
    attached_host_route_redistribution_uuid = %s
  }
%s%s}
`,
		testAccMSOL3OutPrerequisites(siteName, l3outName),
		testAccMSOL3OutRouteMapPolicyPrerequisites(),
		outbound,
		interleak,
		testAccPolicyUUID("static"),
		testAccPolicyUUID("connected"),
		testAccPolicyUUID("attached"),
		bgpConfig,
		ospfConfig,
	)
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
	dampeningValidationMessage := `route_dampening_ipv4_uuid, route_dampening_ipv6_uuid may be set only when enabled is one of true`
	dampeningValidationError := regexp.MustCompile(
		strings.ReplaceAll(regexp.QuoteMeta(dampeningValidationMessage), " ", `\s+`),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccProviderPreCheck(t) },
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
				PreConfig:   func() { fmt.Println("Test: Reject OSPF settings while disabled") },
				Config:      testAccMSOL3OutOSPFDisabledWithAreaConfig(siteName, l3outName),
				ExpectError: regexp.MustCompile(`requires enabled = true`),
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
				PreConfig:   func() { fmt.Println("Test: Reject enabled OSPF without an area ID") },
				Config:      testAccMSOL3OutInvalidOSPFWithoutAreaIDConfig(),
				ExpectError: regexp.MustCompile(`area_id must be configured when enabled is true`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject enabled OSPF without an area type") },
				Config:      testAccMSOL3OutInvalidOSPFWithoutAreaTypeConfig(),
				ExpectError: regexp.MustCompile(`area_type must be configured when enabled is true`),
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
