package provider

import (
	"fmt"
	"regexp"
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
