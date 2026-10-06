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
)

func TestAccMSOL3OutResource(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_l3out")
	var templateID string
	var l3outUUID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccProviderPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create L3Out with direct configuration") },
				Config:    testAccMSOL3OutCreateConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCaptureResourceIdentifiers("mso_l3out.test", &templateID, &l3outUUID),
					resource.TestCheckResourceAttr("mso_l3out.test", "name", l3outName),
					testAccCheckL3OutCompositeID("mso_l3out.test"),
					resource.TestCheckResourceAttr("mso_l3out.test", "description", "initial L3Out description"),
					resource.TestCheckResourceAttrPair("mso_l3out.test", "l3_domain", "mso_fabric_policies_l3_domain.l3out_test_domain_1", "name"),
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
					resource.TestCheckResourceAttrPair("mso_l3out.test", "l3_domain", "mso_fabric_policies_l3_domain.l3out_test_domain_2", "name"),
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
					resource.TestCheckResourceAttrPair("mso_l3out.test", "l3_domain", "mso_fabric_policies_l3_domain.l3out_test_domain_2", "name"),
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
					resource.TestCheckResourceAttrPair("mso_l3out.test", "l3_domain", "mso_fabric_policies_l3_domain.l3out_test_domain_1", "name"),
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
					return fmt.Sprintf("%s/%s", resourceState.Primary.Attributes["template_id"], resourceState.Primary.Attributes["uuid"]), nil
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
						return fmt.Errorf("unexpected imported L3Out identity: %q/%q", attributes["template_id"], attributes["uuid"])
					}
					return nil
				},
			},
			{
				PreConfig: func() { fmt.Println("Test: Replace minimal L3Out when template_id changes") },
				Config:    testAccMSOL3OutForceNewConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccAssertResourceIDChanged("mso_l3out.test", &l3outUUID),
					resource.TestCheckResourceAttr("mso_l3out.test", "name", l3outName),
					resource.TestCheckResourceAttrPair("mso_l3out.test", "template_id", "mso_template.l3out_test_replacement", "id"),
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
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
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
`
}

func testAccMSOL3OutUpdateConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
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
`
}

func testAccMSOL3OutOmitOptionalComputedConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  pim_enabled = false
}
`
}

func testAccMSOL3OutClearConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
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
`
}

func testAccMSOL3OutImportConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description = "updated L3Out description"
  l3_domain   = mso_fabric_policies_l3_domain.l3out_test_domain_2.name
  target_dscp = "expedited_forwarding"
  pim_enabled = false
}
`
}

func testAccMSOL3OutForceNewConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
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
`
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

func TestAccMSOL3OutRoutingProtocol(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_l3out_routing_protocol")
	var templateID, l3outUUID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccProviderPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
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
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.suppress_forwarding_address_in_translated_lsa", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.originate_default_route_always", "false"),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "ospf"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out description with OSPF cost and controls omitted") },
				Config:    testAccMSOL3OutOSPFMinimalUpdateConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply:             []plancheck.PlanCheck{plancheck.ExpectResourceAction("mso_l3out.test", plancheck.ResourceActionUpdate)},
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "description", "OSPF default settings update"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.cost", "1"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.send_redistributed_lsa", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.originate_summary_lsa", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.suppress_forwarding_address_in_translated_lsa", "false"),
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
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.suppress_forwarding_address_in_translated_lsa", "true"),
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
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.suppress_forwarding_address_in_translated_lsa", "false"),
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
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.suppress_forwarding_address_in_translated_lsa", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "originate_default_route", "in_addition"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.originate_default_route_always", "true"),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "ospf"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update L3Out description while omitting OSPF") },
				Config:    testAccMSOL3OutOSPFOmittedUpdateConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction("mso_l3out.test", plancheck.ResourceActionUpdate)},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "description", "OSPF omission update"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.area_id", "0.0.0.2"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.area_type", "nssa"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.cost", "20"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.send_redistributed_lsa", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.originate_summary_lsa", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.suppress_forwarding_address_in_translated_lsa", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "originate_default_route", "in_addition"),
					resource.TestCheckResourceAttr("mso_l3out.test", "ospf.originate_default_route_always", "true"),
					testAccCheckL3OutRoutingProtocol(&templateID, &l3outUUID, "ospf"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Clear L3Out default-route origination with OSPF always omitted") },
				Config:    testAccMSOL3OutDefaultRouteEmptyConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction("mso_l3out.test", plancheck.ResourceActionUpdate)},
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
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction("mso_l3out.test", plancheck.ResourceActionUpdate)},
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
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction("mso_l3out.test", plancheck.ResourceActionUpdate)},
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
				PreConfig: func() { fmt.Println("Test: Update description while omitting both empty protocol objects") },
				Config:    testAccMSOL3OutBothProtocolsOmittedAfterEmptyConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction("mso_l3out.test", plancheck.ResourceActionUpdate)},
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
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction("mso_l3out.test", plancheck.ResourceActionUpdate)},
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

func testAccMSOL3OutOSPFEmptyConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  ospf        = {}
}
`
}

func testAccMSOL3OutOSPFDisabledConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  ospf = {
    enabled = false
  }
}
`
}

func testAccMSOL3OutBGPConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  bgp = {
    enabled = true
  }
  originate_default_route = "in_addition"
}
`
}

func testAccMSOL3OutBGPOmittedConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
resource "mso_l3out" "test" {
  template_id             = mso_template.l3out_test.id
  name                    = local.l3out_name
  vrf_uuid                = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  originate_default_route = "in_addition"
}
`
}

func testAccMSOL3OutBGPOmittedUpdateConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
resource "mso_l3out" "test" {
  template_id             = mso_template.l3out_test.id
  name                    = local.l3out_name
  vrf_uuid                = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description             = "BGP omission update"
  originate_default_route = "in_addition"
}
`
}

func testAccMSOL3OutBGPWithOSPFConfig(siteName, l3outName string) string {
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
    area_id   = "0.0.0.3"
    area_type = "regular"
  }
}
`
}

func testAccMSOL3OutBGPDisabledWithOSPFConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
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
`
}

func testAccMSOL3OutBothProtocolsEmptyConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
resource "mso_l3out" "test" {
  template_id             = mso_template.l3out_test.id
  name                    = local.l3out_name
  vrf_uuid                = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  bgp                     = {}
  ospf                    = {}
  originate_default_route = ""
}
`
}

func testAccMSOL3OutBothProtocolsOmittedAfterEmptyConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
resource "mso_l3out" "test" {
  template_id             = mso_template.l3out_test.id
  name                    = local.l3out_name
  vrf_uuid                = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description             = "protocols omitted after empty"
  originate_default_route = ""
}
`
}

func testAccMSOL3OutBothProtocolsDisabledConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
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
`
}

func testAccMSOL3OutBothProtocolsOmittedAfterFalseConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
resource "mso_l3out" "test" {
  template_id             = mso_template.l3out_test.id
  name                    = local.l3out_name
  vrf_uuid                = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description             = "protocols omitted after false"
  originate_default_route = ""
}
`
}

func testAccMSOL3OutOSPFMinimalConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
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
`
}

func testAccMSOL3OutOSPFMinimalUpdateConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
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
`
}

func testAccMSOL3OutOSPFEnabledConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
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
`
}

func testAccMSOL3OutOSPFUpdatedConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
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
`
}

func testAccMSOL3OutDefaultRouteEmptyConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
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
`
}

func testAccMSOL3OutOSPFOmittedConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
}
`
}

func testAccMSOL3OutOSPFOmittedUpdateConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description = "OSPF omission update"
}
`
}

func TestAccMSOL3OutAnnotations(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_l3out_annotations")
	var templateID, l3outUUID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccProviderPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create L3Out with annotations sharing a value") },
				Config:    testAccMSOL3OutAnnotationsCreateConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCaptureResourceIdentifiers("mso_l3out.test", &templateID, &l3outUUID),
					resource.TestCheckResourceAttr("mso_l3out.test", "annotations.%", "2"),
					resource.TestCheckResourceAttr("mso_l3out.test", "annotations.foo", "shared"),
					resource.TestCheckResourceAttr("mso_l3out.test", "annotations.bar", "shared"),
					testAccCheckL3OutAnnotations(&templateID, &l3outUUID, map[string]string{"foo": "shared", "bar": "shared"}),
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
					testAccCheckL3OutAnnotations(&templateID, &l3outUUID, map[string]string{"foo": "shared", "bar": "shared"}),
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
					testAccCheckL3OutAnnotations(&templateID, &l3outUUID, map[string]string{"foo": "updated", "bar": "shared", "baz": "shared"}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove one L3Out annotation") },
				Config:    testAccMSOL3OutAnnotationsRemoveOneConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "annotations.%", "2"),
					resource.TestCheckResourceAttr("mso_l3out.test", "annotations.foo", "updated"),
					resource.TestCheckResourceAttr("mso_l3out.test", "annotations.baz", "shared"),
					testAccCheckL3OutAnnotations(&templateID, &l3outUUID, map[string]string{"foo": "updated", "baz": "shared"}),
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
					return fmt.Sprintf("%s/%s", resourceState.Primary.Attributes["template_id"], resourceState.Primary.Attributes["uuid"]), nil
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
					testAccCheckL3OutAnnotations(&templateID, &l3outUUID, map[string]string{"foo": "updated", "baz": "shared", "external": "shared"}),
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
		},
	})
}

func testAccMSOL3OutAnnotationsCreateConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
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
`
}

func testAccMSOL3OutAnnotationsReorderedConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
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
`
}

func testAccMSOL3OutAnnotationsUpdateConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
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
`
}

func testAccMSOL3OutAnnotationsRemoveOneConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
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
`
}

func testAccMSOL3OutAnnotationsOmittedConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description = "annotations omitted"
}
`
}

func testAccMSOL3OutAnnotationsClearConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description = "annotations cleared"
  annotations = {}
}
`
}

func TestAccMSOL3OutInterfaceGroups(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_l3out_interface_groups")
	var templateID string
	var l3outUUID string
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccProviderPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig:        func() { fmt.Println("Test: Create L3Out with four interface groups") },
				Config:           testAccMSOL3OutInterfaceGroupsCreateConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCaptureResourceIdentifiers("mso_l3out.test", &templateID, &l3outUUID),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.%", "4"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.edge.bfd.key_id", "20"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.edge.bfd_multi_hop.key_id", "30"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.edge.bfd.key", "bfd-secret"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.edge.bfd_multi_hop.key", "multi-hop-secret"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.edge.ospf.key", "ospf-secret"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.edge.ospf.authentication_type", "md5"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.backup.description", "backup group"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.middle.description", "middle group"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.tail.description", "tail group"),
				),
			},
			{
				PreConfig:        func() { fmt.Println("Test: Update interface group attributes while retaining authentication keys") },
				Config:           testAccMSOL3OutInterfaceGroupsUpdateConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.edge.bfd.key_id", "21"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.edge.bfd_multi_hop.key_id", "31"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.edge.bfd.key", "bfd-secret"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.edge.bfd_multi_hop.key", "multi-hop-secret"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.edge.ospf.key", "ospf-secret"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.edge.ospf.authentication_type", "simple"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.backup.description", "updated backup"),
					testAccCheckL3OutInterfaceGroupOrder(&templateID, &l3outUUID, []string{"backup", "edge", "middle", "tail"}),
				),
			},
			{
				PreConfig:        func() { fmt.Println("Test: Update and add interface groups while removing two nonadjacent groups") },
				Config:           testAccMSOL3OutInterfaceGroupsMixedPatchConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.%", "3"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.edge.description", "edge group after mixed patch"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.tail.description", "tail group after mixed patch"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.new.description", "new group"),
					testAccCheckL3OutInterfaceGroups(&templateID, &l3outUUID, map[string]string{
						"edge": "edge group after mixed patch",
						"new":  "new group",
						"tail": "tail group after mixed patch",
					}),
				),
			},
			{
				PreConfig:        func() { fmt.Println("Test: Retain one interface group") },
				Config:           testAccMSOL3OutInterfaceGroupsRemoveOneConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.%", "1"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.edge.bfd.key_id", "21"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.edge.bfd_multi_hop.key_id", "31"),
				),
			},
			{
				PreConfig:        func() { fmt.Println("Test: Clear all L3Out interface groups") },
				Config:           testAccMSOL3OutInterfaceGroupsEmptyConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
				Check:            resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.%", "0"),
			},
			{
				PreConfig:        func() { fmt.Println("Test: Add interface groups after clearing the collection") },
				Config:           testAccMSOL3OutInterfaceGroupsCreateConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
				Check:            resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.%", "4"),
			},
			{
				PreConfig:        func() { fmt.Println("Test: Omit configured interface groups during an unrelated update") },
				Config:           testAccMSOL3OutInterfaceGroupsOmittedConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "description", "interface groups omitted"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.%", "4"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.edge.bfd_multi_hop.key_id", "30"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.edge.bfd_multi_hop.key", "multi-hop-secret"),
				),
			},
			{
				PreConfig:        func() { fmt.Println("Test: Manage an interface group without optional settings") },
				Config:           testAccMSOL3OutInterfaceGroupsEmptyGroupConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.%", "1"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.empty.bfd.enabled", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.empty.bfd_multi_hop.enabled", "false"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.empty.ospf.enabled", "false"),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import L3Out with an interface group") },
				Config:       testAccMSOL3OutInterfaceGroupsEmptyGroupConfig(siteName, l3outName),
				ResourceName: "mso_l3out.test",
				ImportState:  true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					l3out, ok := state.RootModule().Resources["mso_l3out.test"]
					if !ok {
						return "", fmt.Errorf("mso_l3out.test not found in Terraform state")
					}
					return l3out.Primary.Attributes["id"], nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckL3OutInterfaceGroupOrder(templateID, l3outUUID *string, expected []string) resource.TestCheckFunc {
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
		for index, item := range groups {
			group, ok := item.(map[string]any)
			if !ok {
				return fmt.Errorf("NDO interfaceGroups[%d] has unexpected type %T", index, item)
			}
			name, _, err := ndoapi.StringField(group, "name", ndoapi.RequiredNonEmptyField)
			if err != nil {
				return fmt.Errorf("NDO interfaceGroups[%d]: %w", index, err)
			}
			if name != expected[index] {
				return fmt.Errorf("unexpected remote interface group at index %d: got %q, want %q", index, name, expected[index])
			}
		}
		return nil
	}
}

func testAccMSOL3OutInterfaceGroupsBaseConfig(siteName, l3outName string) string {
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
`
}

func testAccMSOL3OutInterfaceGroupsCreateConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupsBaseConfig(siteName, l3outName) + `
  interface_groups = {
    edge = {
      description = "edge group"
      bfd = {
        enabled                = true
        authentication_enabled = true
        key_id                 = 20
        key                    = "bfd-secret"
      }
      bfd_multi_hop = {
        enabled                = true
        authentication_enabled = true
        key_id                 = 30
        key                    = "multi-hop-secret"
      }
      ospf = {
        enabled             = true
        authentication_type = "md5"
        key_id              = 10
        key                 = "ospf-secret"
      }
    }
    backup = {
      description = "backup group"
    }
    middle = {
      description = "middle group"
    }
    tail = {
      description = "tail group"
    }
  }
}
`
}

func testAccMSOL3OutInterfaceGroupsUpdateConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupsBaseConfig(siteName, l3outName) + `
  interface_groups = {
    edge = {
      description = "edge group"
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
        key                    = "multi-hop-secret"
      }
      ospf = {
        enabled             = true
        authentication_type = "simple"
        key_id              = 10
        key                 = "ospf-secret"
      }
    }
    backup = {
      description = "updated backup"
    }
    middle = {
      description = "middle group"
    }
    tail = {
      description = "tail group"
    }
  }
}
`
}

func testAccMSOL3OutInterfaceGroupsMixedPatchConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupsBaseConfig(siteName, l3outName) + `
  interface_groups = {
    edge = {
      description = "edge group after mixed patch"
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
        key                    = "multi-hop-secret"
      }
      ospf = {
        enabled             = true
        authentication_type = "simple"
        key_id              = 10
        key                 = "ospf-secret"
      }
    }
    new = {
      description = "new group"
    }
    tail = {
      description = "tail group after mixed patch"
    }
  }
}
`
}

func testAccMSOL3OutInterfaceGroupsRemoveOneConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupsBaseConfig(siteName, l3outName) + `
  interface_groups = {
    edge = {
      description = "edge group"
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
        key                    = "multi-hop-secret"
      }
      ospf = {
        enabled             = true
        authentication_type = "simple"
        key_id              = 10
        key                 = "ospf-secret"
      }
    }
  }
}
`
}

func testAccMSOL3OutInterfaceGroupsEmptyConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupsBaseConfig(siteName, l3outName) + `
  interface_groups = {}
}
`
}

func testAccMSOL3OutInterfaceGroupsEmptyGroupConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupsBaseConfig(siteName, l3outName) + `
  interface_groups = {
    empty = {}
  }
}
`
}

func testAccMSOL3OutInterfaceGroupsOmittedConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupsBaseConfig(siteName, l3outName) + `
  description = "interface groups omitted"
}
`
}

func TestAccMSOL3OutInterfaceGroupReferences(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_l3out_group_refs")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccProviderPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig:        func() { fmt.Println("Test: Configure L3Out interface group tenant references and four NetFlow types") },
				Config:           testAccMSOL3OutInterfaceGroupAllReferencesConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("mso_l3out.test", "interface_groups.edge.interface_routing_policy_uuid", "mso_tenant_policies_l3out_interface_routing_policy.routing", "uuid"),
					resource.TestCheckResourceAttrPair("mso_l3out.test", "interface_groups.edge.custom_qos_policy_uuid", "mso_tenant_policies_custom_qos_policy.qos", "uuid"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.edge.netflow_monitor_uuids.%", "4"),
					resource.TestCheckResourceAttrPair("mso_l3out.test", "interface_groups.edge.netflow_monitor_uuids.ipv4", "mso_tenant_policies_netflow_monitor.ipv4", "uuid"),
					resource.TestCheckResourceAttrPair("mso_l3out.test", "interface_groups.edge.netflow_monitor_uuids.ipv6", "mso_tenant_policies_netflow_monitor.ipv6", "uuid"),
					resource.TestCheckResourceAttrPair("mso_l3out.test", "interface_groups.edge.netflow_monitor_uuids.ce", "mso_tenant_policies_netflow_monitor.ce", "uuid"),
					resource.TestCheckResourceAttrPair("mso_l3out.test", "interface_groups.edge.netflow_monitor_uuids.unspecified", "mso_tenant_policies_netflow_monitor.unspecified", "uuid"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.edge.qos_priority", "level6"),
				),
			},
			{
				PreConfig:        func() { fmt.Println("Test: Replace interface group NetFlow map with two references") },
				Config:           testAccMSOL3OutInterfaceGroupTwoReferencesConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.edge.netflow_monitor_uuids.%", "2"),
					resource.TestCheckResourceAttrPair("mso_l3out.test", "interface_groups.edge.netflow_monitor_uuids.ipv4", "mso_tenant_policies_netflow_monitor.ipv4", "uuid"),
					resource.TestCheckResourceAttrPair("mso_l3out.test", "interface_groups.edge.netflow_monitor_uuids.ce", "mso_tenant_policies_netflow_monitor.ce", "uuid"),
				),
			},
			{
				PreConfig:        func() { fmt.Println("Test: Clear parent interface group tenant references and NetFlow monitors") },
				Config:           testAccMSOL3OutInterfaceGroupClearReferencesConfig(siteName, l3outName),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.edge.interface_routing_policy_uuid", ""),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.edge.custom_qos_policy_uuid", ""),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.edge.qos_priority", "unspecified"),
					resource.TestCheckResourceAttr("mso_l3out.test", "interface_groups.edge.netflow_monitor_uuids.%", "0"),
				),
			},
		},
	})
}

func testAccMSOL3OutInterfaceGroupReferenceBaseConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + testAccMSOL3OutInterfaceGroupPolicyTenantReferencesConfig() + `
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  interface_groups = {
    edge = {
      interface_routing_policy_uuid = mso_tenant_policies_l3out_interface_routing_policy.routing.uuid
      custom_qos_policy_uuid        = mso_tenant_policies_custom_qos_policy.qos.uuid
      qos_priority                  = "level6"
`
}

func testAccMSOL3OutInterfaceGroupAllReferencesConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupReferenceBaseConfig(siteName, l3outName) + `
      netflow_monitor_uuids = {
        ipv4        = mso_tenant_policies_netflow_monitor.ipv4.uuid
        ipv6        = mso_tenant_policies_netflow_monitor.ipv6.uuid
        ce          = mso_tenant_policies_netflow_monitor.ce.uuid
        unspecified = mso_tenant_policies_netflow_monitor.unspecified.uuid
      }
    }
  }
}
`
}

func testAccMSOL3OutInterfaceGroupTwoReferencesConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupReferenceBaseConfig(siteName, l3outName) + `
      netflow_monitor_uuids = {
        ipv4 = mso_tenant_policies_netflow_monitor.ipv4.uuid
        ce   = mso_tenant_policies_netflow_monitor.ce.uuid
      }
    }
  }
}
`
}

func testAccMSOL3OutInterfaceGroupClearReferencesConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + testAccMSOL3OutInterfaceGroupPolicyTenantReferencesConfig() + `
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  interface_groups = {
    edge = {
      interface_routing_policy_uuid = ""
      custom_qos_policy_uuid        = ""
      qos_priority                  = "unspecified"
      netflow_monitor_uuids         = {}
    }
  }
}
`
}

func TestAccMSOL3OutInvalid(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_l3out_invalid")
	defaultRouteValidationMessage := `originate_default_route must be one of "only", "in_addition" when ospf.originate_default_route_always is true`
	defaultRouteValidationError := regexp.MustCompile(strings.ReplaceAll(regexp.QuoteMeta(defaultRouteValidationMessage), " ", `\s+`))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccProviderPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Reject OSPF always-originate when the default-route mode is omitted") },
				Config:      testAccMSOL3OutDefaultRouteModeMissingConfig(siteName, l3outName),
				ExpectError: defaultRouteValidationError,
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject OSPF always-originate with an empty default-route mode") },
				Config:      testAccMSOL3OutDefaultRouteModeEmptyConfig(siteName, l3outName),
				ExpectError: defaultRouteValidationError,
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject OSPF settings while disabled") },
				Config:      testAccMSOL3OutOSPFDisabledWithAreaConfig(siteName, l3outName),
				ExpectError: regexp.MustCompile(`requires enabled = true`),
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
				PreConfig:   func() { fmt.Println("Test: Reject authenticated interface group BFD without a key") },
				Config:      testAccMSOL3OutInvalidInterfaceGroupBFDWithoutKeyConfig(),
				ExpectError: regexp.MustCompile(`interface_groups\["edge"\]\.bfd\.key must have a value when\s+interface_groups\["edge"\]\.bfd\.authentication_enabled is true`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject interface group OSPF authentication without a key ID") },
				Config:      testAccMSOL3OutInvalidInterfaceGroupOSPFWithoutKeyIDConfig(),
				ExpectError: regexp.MustCompile(`interface_groups\["edge"\]\.ospf\.key_id must have a value when\s+interface_groups\["edge"\]\.ospf\.authentication_type is one of simple, md5`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject an empty interface group name") },
				Config:      testAccMSOL3OutInvalidInterfaceGroupEmptyNameConfig(),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Length`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject a null interface group") },
				Config:      testAccMSOL3OutInvalidInterfaceGroupNullConfig(),
				ExpectError: regexp.MustCompile(`Null Map Value`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject an empty NetFlow monitor UUID in an interface group") },
				Config:      testAccMSOL3OutInvalidInterfaceGroupEmptyNetFlowUUIDConfig(),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Length`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject a null NetFlow monitor UUID in an interface group") },
				Config:      testAccMSOL3OutInvalidInterfaceGroupNullNetFlowUUIDConfig(),
				ExpectError: regexp.MustCompile(`Null Map Value`),
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

func testAccMSOL3OutInvalidOSPFWithoutAreaIDConfig() string {
	return `
resource "mso_l3out" "test" {
  template_id = "template"
  name        = "test"
  vrf_uuid    = "vrf"
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
  template_id = "template"
  name        = "test"
  vrf_uuid    = "vrf"
  ospf = {
    enabled = true
    area_id = "0.0.0.1"
  }
}
`
}

func testAccMSOL3OutInvalidInterfaceGroupBFDWithoutKeyConfig() string {
	return `
resource "mso_l3out" "test" {
  template_id = "template"
  name        = "test"
  vrf_uuid    = "vrf"
  interface_groups = {
    edge = {
      bfd = {
        enabled                = true
        authentication_enabled = true
      }
    }
  }
}
`
}

func testAccMSOL3OutInvalidInterfaceGroupOSPFWithoutKeyIDConfig() string {
	return `
resource "mso_l3out" "test" {
  template_id = "template"
  name        = "test"
  vrf_uuid    = "vrf"
  interface_groups = {
    edge = {
      ospf = {
        enabled             = true
        authentication_type = "md5"
        key                 = "secret"
      }
    }
  }
}
`
}

func testAccMSOL3OutInvalidInterfaceGroupEmptyNameConfig() string {
	return `
resource "mso_l3out" "test" {
  template_id = "template"
  name        = "test"
  vrf_uuid    = "vrf"
  interface_groups = {
    "" = {}
  }
}
`
}

func testAccMSOL3OutInvalidInterfaceGroupNullConfig() string {
	return `
resource "mso_l3out" "test" {
  template_id = "template"
  name        = "test"
  vrf_uuid    = "vrf"
  interface_groups = {
    edge = null
  }
}
`
}

func testAccMSOL3OutInvalidInterfaceGroupEmptyNetFlowUUIDConfig() string {
	return `
resource "mso_l3out" "test" {
  template_id = "template"
  name        = "test"
  vrf_uuid    = "vrf"
  interface_groups = {
    edge = {
      netflow_monitor_uuids = {
        ipv4 = ""
      }
    }
  }
}
`
}

func testAccMSOL3OutInvalidInterfaceGroupNullNetFlowUUIDConfig() string {
	return `
resource "mso_l3out" "test" {
  template_id = "template"
  name        = "test"
  vrf_uuid    = "vrf"
  interface_groups = {
    edge = {
      netflow_monitor_uuids = {
        ipv4 = null
      }
    }
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
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
}
`
}

func testAccMSOL3OutDefaultRouteModeMissingConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
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
`
}

func testAccMSOL3OutDefaultRouteModeEmptyConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
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
`
}

func testAccMSOL3OutOSPFDisabledWithAreaConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  ospf = {
    enabled = false
    area_id = "0.0.0.1"
  }
}
`
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
		expected := attributes["template_id"] + "/" + attributes["uuid"]
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

	if _, err := msoClient.PatchbyID(ndoapi.TemplateEndpoint(templateID), ndoapi.NewRemovePatchPayload(resolved.PatchPath())); err != nil {
		t.Fatalf("unable to delete L3Out out of band: %v", err)
	}
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
		ndoapi.NewPatchPayload("replace", resolved.PatchPath()+"/tagAnnotations", annotations.ToPayload()),
	)
	if err != nil {
		t.Fatalf("unable to replace L3Out annotations out of band: %v", err)
	}
}
