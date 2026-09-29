package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var (
	testNetflowExporterSiteTemplateID string
	testNetflowExporterSiteSiteID     string
	testNetflowExporterSiteUUID       string
)

// TestAccMSOTenantPoliciesNetflowExporterSiteResource exercises the lifecycle
// of a site-specific NetFlow Exporter configuration. The schema prerequisite
// deploys the BD-backed application EPG, external EPG, and L3Out so NDO can
// provide the template-level objects and site association required by this
// resource. The resource adds the site exporter entry when NDO has not
// replicated it automatically.
func TestAccMSOTenantPoliciesNetflowExporterSiteResource(t *testing.T) {
	resourceName := "mso_tenant_policies_netflow_exporter_site.netflow_exporter_site"
	siteDataSource := "data.mso_site." + msoTemplateSiteName1
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccVersionCheck(t, "5.1") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOTenantPoliciesNetflowExporterSiteDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create NetFlow Exporter site configuration without EPG association") },
				Config:    testAccMSOTenantPoliciesNetflowExporterSiteConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "template_id"),
					resource.TestCheckResourceAttrSet(resourceName, "site_id"),
					resource.TestCheckResourceAttrSet(resourceName, "netflow_exporter_uuid"),
					resource.TestCheckResourceAttrPair(resourceName, "site_id", siteDataSource, "id"),
					resource.TestCheckResourceAttr(resourceName, "source_ip_address_type", "inbandMgmtIP"),
					resource.TestCheckResourceAttr(resourceName, "source_ip", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_port", "10"),
					resource.TestCheckResourceAttr(resourceName, "application_epg.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "external_epg.#", "0"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok {
							return fmt.Errorf("resource %s not found in state", resourceName)
						}
						testNetflowExporterSiteTemplateID = rs.Primary.Attributes["template_id"]
						testNetflowExporterSiteSiteID = rs.Primary.Attributes["site_id"]
						testNetflowExporterSiteUUID = rs.Primary.Attributes["netflow_exporter_uuid"]
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate NetFlow Exporter site configuration after out-of-band deletion")
					if err := manuallyDeleteNetflowExporterSiteFromTemplate(
						testNetflowExporterSiteTemplateID,
						testNetflowExporterSiteSiteID,
						testNetflowExporterSiteUUID,
					); err != nil {
						t.Fatalf("Failed to manually delete NetFlow Exporter site configuration: %v", err)
					}
				},
				Config: testAccMSOTenantPoliciesNetflowExporterSiteConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source_ip_address_type", "inbandMgmtIP"),
					resource.TestCheckResourceAttr(resourceName, "source_ip", ""),
					resource.TestCheckResourceAttr(resourceName, "application_epg.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "external_epg.#", "0"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update NetFlow Exporter site configuration with application EPG") },
				Config:    testAccMSOTenantPoliciesNetflowExporterSiteConfigApplicationEpg(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source_ip_address_type", "inbandMgmtIP"),
					resource.TestCheckResourceAttr(resourceName, "source_ip", ""),
					resource.TestCheckResourceAttr(resourceName, "application_epg.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "external_epg.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "application_epg.0.tenant_name", msoTenantName),
					resource.TestCheckResourceAttr(resourceName, "application_epg.0.anp_name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr(resourceName, "application_epg.0.epg_name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr(resourceName, "application_epg.0.vrf_name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr(resourceName, "application_epg.0.vrf_tenant_name", msoTenantName),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update NetFlow Exporter site configuration with external EPG") },
				Config:    testAccMSOTenantPoliciesNetflowExporterSiteConfigExternalEpgOobMgmtIp(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source_ip_address_type", "oobMgmtIP"),
					resource.TestCheckResourceAttr(resourceName, "source_ip", ""),
					resource.TestCheckResourceAttr(resourceName, "application_epg.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "external_epg.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "external_epg.0.tenant_name", msoTenantName),
					resource.TestCheckResourceAttr(resourceName, "external_epg.0.l3out_name", msoSchemaTemplateL3outName),
					resource.TestCheckResourceAttr(resourceName, "external_epg.0.external_epg_name", msoSchemaTemplateExtEpgName),
					resource.TestCheckResourceAttr(resourceName, "external_epg.0.vrf_name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr(resourceName, "external_epg.0.vrf_tenant_name", msoTenantName),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update NetFlow Exporter site configuration with application EPG using common VRF")
				},
				Config: testAccMSOTenantPoliciesNetflowExporterSiteConfigApplicationEpgCommonVrf(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source_ip_address_type", "inbandMgmtIP"),
					resource.TestCheckResourceAttr(resourceName, "application_epg.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "external_epg.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "application_epg.0.tenant_name", msoTenantName),
					resource.TestCheckResourceAttr(resourceName, "application_epg.0.anp_name", msoSchemaTemplateAnpWithVrfInCommonTenantName),
					resource.TestCheckResourceAttr(resourceName, "application_epg.0.epg_name", msoSchemaTemplateAnpEpgWithVrfInCommonTenantName),
					resource.TestCheckResourceAttr(resourceName, "application_epg.0.vrf_name", msoSchemaTemplateVrfInCommonTenantName),
					resource.TestCheckResourceAttr(resourceName, "application_epg.0.vrf_tenant_name", "common"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Return NetFlow Exporter site configuration to application-tenant VRF without explicit VRF tenant")
				},
				Config: testAccMSOTenantPoliciesNetflowExporterSiteConfigApplicationEpgSameTenantVrf(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source_ip_address_type", "inbandMgmtIP"),
					resource.TestCheckResourceAttr(resourceName, "application_epg.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "external_epg.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "application_epg.0.tenant_name", msoTenantName),
					resource.TestCheckResourceAttr(resourceName, "application_epg.0.anp_name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr(resourceName, "application_epg.0.epg_name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr(resourceName, "application_epg.0.vrf_name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr(resourceName, "application_epg.0.vrf_tenant_name", msoTenantName),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update NetFlow Exporter site configuration with external EPG using common VRF")
				},
				Config: testAccMSOTenantPoliciesNetflowExporterSiteConfigExternalEpgCommonVrf(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source_ip_address_type", "inbandMgmtIP"),
					resource.TestCheckResourceAttr(resourceName, "application_epg.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "external_epg.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "external_epg.0.tenant_name", msoTenantName),
					resource.TestCheckResourceAttr(resourceName, "external_epg.0.l3out_name", msoSchemaTemplateL3outWithVrfInCommonTenantName),
					resource.TestCheckResourceAttr(resourceName, "external_epg.0.external_epg_name", msoSchemaTemplateExternalEpgWithVrfInCommonTenantName),
					resource.TestCheckResourceAttr(resourceName, "external_epg.0.vrf_name", msoSchemaTemplateVrfInCommonTenantName),
					resource.TestCheckResourceAttr(resourceName, "external_epg.0.vrf_tenant_name", "common"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update NetFlow Exporter site configuration with PTEP source IP") },
				Config:    testAccMSOTenantPoliciesNetflowExporterSiteConfigExternalEpgPtep(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source_ip_address_type", "ptep"),
					resource.TestCheckResourceAttr(resourceName, "source_ip", ""),
					resource.TestCheckResourceAttr(resourceName, "application_epg.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "external_epg.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "external_epg.0.vrf_tenant_name", msoTenantName),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update NetFlow Exporter site configuration with custom source IP") },
				Config:    testAccMSOTenantPoliciesNetflowExporterSiteConfigCustomSourceIp(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source_ip_address_type", "customSrcIP"),
					resource.TestCheckResourceAttr(resourceName, "source_ip", "10.0.1.14/8"),
					resource.TestCheckResourceAttr(resourceName, "application_epg.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "external_epg.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "external_epg.0.vrf_tenant_name", msoTenantName),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove EPG associations and custom source IP") },
				Config:    testAccMSOTenantPoliciesNetflowExporterSiteConfigRemoveAssociations(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source_ip_address_type", "inbandMgmtIP"),
					resource.TestCheckResourceAttr(resourceName, "source_ip", ""),
					resource.TestCheckResourceAttr(resourceName, "application_epg.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "external_epg.#", "0"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import NetFlow Exporter site configuration") },
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig:   func() { fmt.Println("Test: Fail when the site is not attached to the tenant template") },
				Config:      testAccMSOTenantPoliciesNetflowExporterSiteConfigSiteNotAttached(),
				ExpectError: regexp.MustCompile(`(?i)site .* is not present in template`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Fail when the NetFlow Exporter does not exist in the template") },
				Config:      testAccMSOTenantPoliciesNetflowExporterSiteConfigExporterDoesNotExist(),
				ExpectError: regexp.MustCompile(`(?i)NetFlow Exporter .* does not exist in template`),
			},
		},
	})
}

func testAccMSOTenantPoliciesNetflowExporterSitePrerequisiteConfig() string {
	siteDataSource := "data.mso_site." + msoTemplateSiteName1
	siteConfiguration := fmt.Sprintf(`
resource "mso_schema_site_external_epg" "site_external_epg" {
	schema_id         = mso_schema.%[1]s.id
	template_name     = tolist(mso_schema.%[1]s.template)[0].name
	site_id           = %[2]s.id
	external_epg_name = mso_schema_template_external_epg.%[3]s.external_epg_name
	l3out_name        = mso_schema_template_l3out.%[4]s.l3out_name
}

`, msoSchemaName, siteDataSource, msoSchemaTemplateExtEpgName, msoSchemaTemplateL3outName)

	return fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s",
		testSchemaWithBothSitesPrerequisiteConfig(),
		testSchemaSiteConfig(msoSchemaSiteResourceLabel1, msoTemplateSiteName1, true),
		testSchemaTemplateVrfConfig(),
		testSchemaTemplateBdConfig(),
		testSchemaTemplateAnpConfig(),
		testAccMSOSchemaSiteAnpEpgTemplateAnpEpgWithBdConfig(),
		testSchemaTemplateL3outConfig(),
		testSchemaTemplateExtEpgConfig(),
		siteConfiguration,
		testSchemaTemplateDeployNdoConfig([]string{
			"mso_schema_site." + msoSchemaSiteResourceLabel1,
			"mso_schema_template_vrf." + msoSchemaTemplateVrfName,
			"mso_schema_template_bd." + msoSchemaTemplateBdName,
			"mso_schema_template_anp." + msoSchemaTemplateAnpName,
			"mso_schema_template_anp_epg." + msoSchemaTemplateAnpEpgName,
			"mso_schema_template_l3out." + msoSchemaTemplateL3outName,
			"mso_schema_template_external_epg." + msoSchemaTemplateExtEpgName,
			"mso_schema_site_external_epg.site_external_epg",
		}),
	)
}

func testAccMSOTenantPoliciesNetflowExporterSiteCommonVrfPrerequisiteConfig() string {
	siteDataSource := "data.mso_site." + msoTemplateSiteName1
	commonSchemaConfig := fmt.Sprintf(`
data "mso_tenant" "common" {
	name = "common"
}

resource "mso_schema" "common_tenant_vrf_schema" {
	name = "%[1]s"
	template {
		name         = "%[2]s"
		display_name = "%[2]s"
		tenant_id    = data.mso_tenant.common.id
	}
}

resource "mso_schema_site" "common_tenant_vrf_site" {
	schema_id           = mso_schema.common_tenant_vrf_schema.id
	site_id             = %[3]s.id
	template_name       = tolist(mso_schema.common_tenant_vrf_schema.template)[0].name
	undeploy_on_destroy = true
}

resource "mso_schema_template_vrf" "vrf_in_common_tenant" {
	schema_id    = mso_schema.common_tenant_vrf_schema.id
	template     = tolist(mso_schema.common_tenant_vrf_schema.template)[0].name
	name         = "%[4]s"
	display_name = "%[4]s"
}
`, msoSchemaTemplateName2, msoSchemaTemplateName2, siteDataSource, msoSchemaTemplateVrfInCommonTenantName)

	additionalCommonVrfConfig := fmt.Sprintf(`
resource "mso_schema_template_bd" "bd_with_vrf_in_common_tenant" {
	schema_id             = mso_schema.%[1]s.id
	template_name         = tolist(mso_schema.%[1]s.template)[0].name
	name                  = "%[2]s"
	display_name          = "%[2]s"
	layer2_unknown_unicast = "proxy"
	vrf_name              = mso_schema_template_vrf.vrf_in_common_tenant.name
	vrf_schema_id         = mso_schema.common_tenant_vrf_schema.id
	vrf_template_name     = tolist(mso_schema.common_tenant_vrf_schema.template)[0].name
}

resource "mso_schema_template_anp" "anp_with_vrf_in_common_tenant" {
	schema_id    = mso_schema.%[1]s.id
	template     = tolist(mso_schema.%[1]s.template)[0].name
	name         = "%[3]s"
	display_name = "%[3]s"
}

resource "mso_schema_template_anp_epg" "epg_with_vrf_in_common_tenant" {
	schema_id     = mso_schema.%[1]s.id
	template_name = tolist(mso_schema.%[1]s.template)[0].name
	anp_name      = mso_schema_template_anp.anp_with_vrf_in_common_tenant.name
	name          = "%[4]s"
	display_name  = "%[4]s"
	bd_name       = mso_schema_template_bd.bd_with_vrf_in_common_tenant.name
}

resource "mso_schema_template_l3out" "l3out_with_vrf_in_common_tenant" {
	schema_id          = mso_schema.%[1]s.id
	template_name      = tolist(mso_schema.%[1]s.template)[0].name
	l3out_name         = "%[5]s"
	display_name       = "%[5]s"
	vrf_name           = mso_schema_template_vrf.vrf_in_common_tenant.name
	vrf_schema_id      = mso_schema.common_tenant_vrf_schema.id
	vrf_template_name  = tolist(mso_schema.common_tenant_vrf_schema.template)[0].name
}

resource "mso_schema_template_external_epg" "external_epg_with_vrf_in_common_tenant" {
	schema_id           = mso_schema.%[1]s.id
	template_name       = tolist(mso_schema.%[1]s.template)[0].name
	external_epg_name   = "%[6]s"
	display_name        = "%[6]s"
	vrf_name            = mso_schema_template_vrf.vrf_in_common_tenant.name
	vrf_schema_id       = mso_schema.common_tenant_vrf_schema.id
	vrf_template_name   = tolist(mso_schema.common_tenant_vrf_schema.template)[0].name
	external_epg_type   = "on-premise"
	l3out_name          = mso_schema_template_l3out.l3out_with_vrf_in_common_tenant.l3out_name
	l3out_schema_id     = mso_schema_template_l3out.l3out_with_vrf_in_common_tenant.schema_id
	l3out_template_name = mso_schema_template_l3out.l3out_with_vrf_in_common_tenant.template_name
}

resource "mso_schema_site_external_epg" "site_external_epg_with_vrf_in_common_tenant" {
	schema_id         = mso_schema.%[1]s.id
	template_name     = tolist(mso_schema.%[1]s.template)[0].name
	site_id           = %[7]s.id
	external_epg_name = mso_schema_template_external_epg.external_epg_with_vrf_in_common_tenant.external_epg_name
	l3out_name        = mso_schema_template_l3out.l3out_with_vrf_in_common_tenant.l3out_name
}
`, msoSchemaName, msoSchemaTemplateBdWithVrfInCommonTenantName, msoSchemaTemplateAnpWithVrfInCommonTenantName, msoSchemaTemplateAnpEpgWithVrfInCommonTenantName, msoSchemaTemplateL3outWithVrfInCommonTenantName, msoSchemaTemplateExternalEpgWithVrfInCommonTenantName, siteDataSource)

	siteConfiguration := fmt.Sprintf(`
resource "mso_schema_site_external_epg" "site_external_epg" {
	schema_id         = mso_schema.%[1]s.id
	template_name     = tolist(mso_schema.%[1]s.template)[0].name
	site_id           = %[2]s.id
	external_epg_name = mso_schema_template_external_epg.%[3]s.external_epg_name
	l3out_name        = mso_schema_template_l3out.%[4]s.l3out_name
}

`, msoSchemaName, siteDataSource, msoSchemaTemplateExtEpgName, msoSchemaTemplateL3outName)

	deployConfig := fmt.Sprintf(`
resource "mso_schema_template_deploy_ndo" "common_tenant_vrf_deploy" {
	schema_id     = mso_schema.common_tenant_vrf_schema.id
	template_name = tolist(mso_schema.common_tenant_vrf_schema.template)[0].name
	force_apply   = ""
	depends_on = [
		mso_schema_site.common_tenant_vrf_site,
		mso_schema_template_vrf.vrf_in_common_tenant,
	]
}

resource "mso_schema_template_deploy_ndo" "deploy" {
	schema_id     = mso_schema.%[1]s.id
	template_name = tolist(mso_schema.%[1]s.template)[0].name
	force_apply   = ""
	depends_on = [
		mso_schema_template_deploy_ndo.common_tenant_vrf_deploy,
		mso_schema_site.%[2]s,
		mso_schema_template_vrf.%[3]s,
		mso_schema_template_bd.%[4]s,
		mso_schema_template_anp.%[5]s,
		mso_schema_template_anp_epg.%[6]s,
		mso_schema_template_l3out.%[7]s,
		mso_schema_template_external_epg.%[8]s,
		mso_schema_site_external_epg.site_external_epg,
		mso_schema_template_bd.bd_with_vrf_in_common_tenant,
		mso_schema_template_anp.anp_with_vrf_in_common_tenant,
		mso_schema_template_anp_epg.epg_with_vrf_in_common_tenant,
		mso_schema_template_l3out.l3out_with_vrf_in_common_tenant,
		mso_schema_template_external_epg.external_epg_with_vrf_in_common_tenant,
		mso_schema_site_external_epg.site_external_epg_with_vrf_in_common_tenant,
	]
}
`, msoSchemaName, msoSchemaSiteResourceLabel1, msoSchemaTemplateVrfName, msoSchemaTemplateBdName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateL3outName, msoSchemaTemplateExtEpgName)

	return fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%s%s",
		testSchemaWithBothSitesPrerequisiteConfig(),
		testSchemaSiteConfig(msoSchemaSiteResourceLabel1, msoTemplateSiteName1, true),
		commonSchemaConfig,
		testSchemaTemplateVrfConfig(),
		testSchemaTemplateBdConfig(),
		testSchemaTemplateAnpConfig(),
		testAccMSOSchemaSiteAnpEpgTemplateAnpEpgWithBdConfig(),
		testSchemaTemplateL3outConfig(),
		testSchemaTemplateExtEpgConfig(),
		siteConfiguration,
		additionalCommonVrfConfig,
		deployConfig,
	)
}

func testAccMSOTenantPoliciesNetflowExporterSiteBaseConfig() string {
	siteDataSource := "data.mso_site." + msoTemplateSiteName1
	return fmt.Sprintf(`%[1]s
	resource "mso_template" "template_tenant" {
		template_name = "%[2]s"
		template_type = "tenant"
		tenant_id     = mso_tenant.%[3]s.id
		sites         = [%[4]s.id]
	}

	resource "mso_tenant_policies_netflow_exporter" "netflow_exporter" {
		template_id = mso_template.template_tenant.id
		name        = "test_netflow_exporter"
	}
`, testAccMSOTenantPoliciesNetflowExporterSitePrerequisiteConfig(), msoTenantPolicyTemplateName, msoTenantName, siteDataSource)
}

func testAccMSOTenantPoliciesNetflowExporterSiteCommonVrfBaseConfig() string {
	siteDataSource := "data.mso_site." + msoTemplateSiteName1
	return fmt.Sprintf(`%[1]s
	resource "mso_template" "template_tenant" {
		template_name = "%[2]s"
		template_type = "tenant"
		tenant_id     = mso_tenant.%[3]s.id
		sites         = [%[4]s.id]
	}

	resource "mso_tenant_policies_netflow_exporter" "netflow_exporter" {
		template_id = mso_template.template_tenant.id
		name        = "test_netflow_exporter"
	}
`, testAccMSOTenantPoliciesNetflowExporterSiteCommonVrfPrerequisiteConfig(), msoTenantPolicyTemplateName, msoTenantName, siteDataSource)
}

func testAccMSOTenantPoliciesNetflowExporterSiteConfigCreate() string {
	siteDataSource := "data.mso_site." + msoTemplateSiteName1
	return fmt.Sprintf(`%[1]s
	resource "mso_tenant_policies_netflow_exporter_site" "netflow_exporter_site" {
		template_id           = mso_tenant_policies_netflow_exporter.netflow_exporter.template_id
		site_id               = %[2]s.id
		netflow_exporter_uuid = mso_tenant_policies_netflow_exporter.netflow_exporter.uuid

		source_ip_address_type = "inbandMgmtIP"
		destination_port       = "10"
		destination_ip         = "10.0.0.15"
		dscp                   = "cs1"

		depends_on = [mso_schema_template_deploy_ndo.deploy, mso_template.template_tenant]
	}
`, testAccMSOTenantPoliciesNetflowExporterSiteBaseConfig(), siteDataSource)
}

func testAccMSOTenantPoliciesNetflowExporterSiteConfigApplicationEpg() string {
	siteDataSource := "data.mso_site." + msoTemplateSiteName1
	return fmt.Sprintf(`%[1]s
	resource "mso_tenant_policies_netflow_exporter_site" "netflow_exporter_site" {
		template_id           = mso_tenant_policies_netflow_exporter.netflow_exporter.template_id
		site_id               = %[2]s.id
		netflow_exporter_uuid = mso_tenant_policies_netflow_exporter.netflow_exporter.uuid

		source_ip_address_type = "inbandMgmtIP"
		destination_port       = "10"
		destination_ip         = "10.0.0.15"
		dscp                   = "cs1"

		application_epg {
			tenant_name     = "%[3]s"
			anp_name        = "%[4]s"
			epg_name        = "%[5]s"
			vrf_name        = "%[6]s"
		}

		depends_on = [mso_schema_template_deploy_ndo.deploy, mso_template.template_tenant]
	}
`, testAccMSOTenantPoliciesNetflowExporterSiteBaseConfig(), siteDataSource, msoTenantName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateVrfName)
}

func testAccMSOTenantPoliciesNetflowExporterSiteConfigExternalEpgOobMgmtIp() string {
	siteDataSource := "data.mso_site." + msoTemplateSiteName1
	return fmt.Sprintf(`%[1]s
	resource "mso_tenant_policies_netflow_exporter_site" "netflow_exporter_site" {
		template_id           = mso_tenant_policies_netflow_exporter.netflow_exporter.template_id
		site_id               = %[2]s.id
		netflow_exporter_uuid = mso_tenant_policies_netflow_exporter.netflow_exporter.uuid

		source_ip_address_type = "oobMgmtIP"
		destination_port       = "10"
		destination_ip         = "10.0.0.15"
		dscp                   = "cs1"

		external_epg {
			tenant_name       = "%[3]s"
			l3out_name        = "%[4]s"
			external_epg_name = "%[5]s"
			vrf_name          = "%[6]s"
		}

		depends_on = [mso_schema_template_deploy_ndo.deploy, mso_template.template_tenant]
	}
`, testAccMSOTenantPoliciesNetflowExporterSiteBaseConfig(), siteDataSource, msoTenantName, msoSchemaTemplateL3outName, msoSchemaTemplateExtEpgName, msoSchemaTemplateVrfName)
}

func testAccMSOTenantPoliciesNetflowExporterSiteConfigApplicationEpgCommonVrf() string {
	siteDataSource := "data.mso_site." + msoTemplateSiteName1
	return fmt.Sprintf(`%[1]s
	resource "mso_tenant_policies_netflow_exporter_site" "netflow_exporter_site" {
		template_id           = mso_tenant_policies_netflow_exporter.netflow_exporter.template_id
		site_id               = %[2]s.id
		netflow_exporter_uuid = mso_tenant_policies_netflow_exporter.netflow_exporter.uuid

		source_ip_address_type = "inbandMgmtIP"
		destination_port       = "10"
		destination_ip         = "10.0.0.15"
		dscp                   = "cs1"

		application_epg {
			tenant_name     = "%[3]s"
			anp_name        = "%[4]s"
			epg_name        = "%[5]s"
			vrf_name        = mso_schema_template_vrf.vrf_in_common_tenant.name
			vrf_tenant_name = "common"
		}

		depends_on = [mso_schema_template_deploy_ndo.deploy, mso_template.template_tenant]
	}
`, testAccMSOTenantPoliciesNetflowExporterSiteCommonVrfBaseConfig(), siteDataSource, msoTenantName, msoSchemaTemplateAnpWithVrfInCommonTenantName, msoSchemaTemplateAnpEpgWithVrfInCommonTenantName)
}

func testAccMSOTenantPoliciesNetflowExporterSiteConfigApplicationEpgSameTenantVrf() string {
	siteDataSource := "data.mso_site." + msoTemplateSiteName1
	return fmt.Sprintf(`%[1]s
	resource "mso_tenant_policies_netflow_exporter_site" "netflow_exporter_site" {
		template_id           = mso_tenant_policies_netflow_exporter.netflow_exporter.template_id
		site_id               = %[2]s.id
		netflow_exporter_uuid = mso_tenant_policies_netflow_exporter.netflow_exporter.uuid

		source_ip_address_type = "inbandMgmtIP"
		destination_port       = "10"
		destination_ip         = "10.0.0.15"
		dscp                   = "cs1"

		application_epg {
			tenant_name = "%[3]s"
			anp_name    = "%[4]s"
			epg_name    = "%[5]s"
			vrf_name    = mso_schema_template_vrf.%[6]s.name
		}

		depends_on = [mso_schema_template_deploy_ndo.deploy, mso_template.template_tenant]
	}
`, testAccMSOTenantPoliciesNetflowExporterSiteCommonVrfBaseConfig(), siteDataSource, msoTenantName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateVrfName)
}

func testAccMSOTenantPoliciesNetflowExporterSiteConfigExternalEpgCommonVrf() string {
	siteDataSource := "data.mso_site." + msoTemplateSiteName1
	return fmt.Sprintf(`%[1]s
	resource "mso_tenant_policies_netflow_exporter_site" "netflow_exporter_site" {
		template_id           = mso_tenant_policies_netflow_exporter.netflow_exporter.template_id
		site_id               = %[2]s.id
		netflow_exporter_uuid = mso_tenant_policies_netflow_exporter.netflow_exporter.uuid

		source_ip_address_type = "inbandMgmtIP"
		destination_port       = "10"
		destination_ip         = "10.0.0.15"
		dscp                   = "cs1"

		external_epg {
			tenant_name       = "%[3]s"
			l3out_name        = "%[4]s"
			external_epg_name = "%[5]s"
			vrf_name          = mso_schema_template_vrf.vrf_in_common_tenant.name
			vrf_tenant_name   = "common"
		}

		depends_on = [mso_schema_template_deploy_ndo.deploy, mso_template.template_tenant]
	}
`, testAccMSOTenantPoliciesNetflowExporterSiteCommonVrfBaseConfig(), siteDataSource, msoTenantName, msoSchemaTemplateL3outWithVrfInCommonTenantName, msoSchemaTemplateExternalEpgWithVrfInCommonTenantName)
}

func testAccMSOTenantPoliciesNetflowExporterSiteConfigExternalEpgPtep() string {
	siteDataSource := "data.mso_site." + msoTemplateSiteName1
	return fmt.Sprintf(`%[1]s
	resource "mso_tenant_policies_netflow_exporter_site" "netflow_exporter_site" {
		template_id           = mso_tenant_policies_netflow_exporter.netflow_exporter.template_id
		site_id               = %[2]s.id
		netflow_exporter_uuid = mso_tenant_policies_netflow_exporter.netflow_exporter.uuid

		source_ip_address_type = "ptep"
		destination_port       = "10"
		destination_ip         = "10.0.0.15"
		dscp                   = "cs1"

		external_epg {
			tenant_name       = "%[3]s"
			l3out_name        = "%[4]s"
			external_epg_name = "%[5]s"
			vrf_name          = "%[6]s"
		}

		depends_on = [mso_schema_template_deploy_ndo.deploy, mso_template.template_tenant]
	}
`, testAccMSOTenantPoliciesNetflowExporterSiteBaseConfig(), siteDataSource, msoTenantName, msoSchemaTemplateL3outName, msoSchemaTemplateExtEpgName, msoSchemaTemplateVrfName)
}

func testAccMSOTenantPoliciesNetflowExporterSiteConfigCustomSourceIp() string {
	siteDataSource := "data.mso_site." + msoTemplateSiteName1
	return fmt.Sprintf(`%[1]s
	resource "mso_tenant_policies_netflow_exporter_site" "netflow_exporter_site" {
		template_id           = mso_tenant_policies_netflow_exporter.netflow_exporter.template_id
		site_id               = %[2]s.id
		netflow_exporter_uuid = mso_tenant_policies_netflow_exporter.netflow_exporter.uuid

		source_ip_address_type = "customSrcIP"
		source_ip              = "10.0.1.14/8"
		destination_port       = "10"
		destination_ip         = "10.0.0.15"
		dscp                   = "cs1"

		external_epg {
			tenant_name       = "%[3]s"
			l3out_name        = "%[4]s"
			external_epg_name = "%[5]s"
			vrf_name          = "%[6]s"
		}

		depends_on = [mso_schema_template_deploy_ndo.deploy, mso_template.template_tenant]
	}
`, testAccMSOTenantPoliciesNetflowExporterSiteBaseConfig(), siteDataSource, msoTenantName, msoSchemaTemplateL3outName, msoSchemaTemplateExtEpgName, msoSchemaTemplateVrfName)
}

func testAccMSOTenantPoliciesNetflowExporterSiteConfigRemoveAssociations() string {
	siteDataSource := "data.mso_site." + msoTemplateSiteName1
	return fmt.Sprintf(`%[1]s
	resource "mso_tenant_policies_netflow_exporter_site" "netflow_exporter_site" {
		template_id           = mso_tenant_policies_netflow_exporter.netflow_exporter.template_id
		site_id               = %[2]s.id
		netflow_exporter_uuid = mso_tenant_policies_netflow_exporter.netflow_exporter.uuid

		source_ip_address_type = "inbandMgmtIP"
		destination_port       = "10"
		destination_ip         = "10.0.0.15"
		dscp                   = "cs1"

		depends_on = [mso_schema_template_deploy_ndo.deploy, mso_template.template_tenant]
	}
`, testAccMSOTenantPoliciesNetflowExporterSiteBaseConfig(), siteDataSource)
}

func testAccMSOTenantPoliciesNetflowExporterSiteConfigSiteNotAttached() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_tenant_policies_netflow_exporter_site" "netflow_exporter_site" {
		template_id           = mso_tenant_policies_netflow_exporter.netflow_exporter.template_id
		site_id               = data.mso_site.%[2]s.id
		netflow_exporter_uuid = mso_tenant_policies_netflow_exporter.netflow_exporter.uuid

		source_ip_address_type = "inbandMgmtIP"
		destination_port       = "10"
		destination_ip         = "10.0.0.15"
		dscp                   = "cs1"

		depends_on = [mso_schema_template_deploy_ndo.deploy, mso_template.template_tenant]
	}
`, testAccMSOTenantPoliciesNetflowExporterSiteBaseConfig(), msoTemplateSiteName2)
}

func testAccMSOTenantPoliciesNetflowExporterSiteConfigExporterDoesNotExist() string {
	siteDataSource := "data.mso_site." + msoTemplateSiteName1
	return fmt.Sprintf(`%[1]s
	resource "mso_tenant_policies_netflow_exporter_site" "netflow_exporter_site" {
		template_id           = mso_tenant_policies_netflow_exporter.netflow_exporter.template_id
		site_id               = %[2]s.id
		netflow_exporter_uuid = "00000000-0000-0000-0000-000000000000"

		source_ip_address_type = "inbandMgmtIP"
		destination_port       = "10"
		destination_ip         = "10.0.0.15"
		dscp                   = "cs1"

		depends_on = [mso_schema_template_deploy_ndo.deploy, mso_template.template_tenant]
	}
`, testAccMSOTenantPoliciesNetflowExporterSiteBaseConfig(), siteDataSource)
}

func manuallyDeleteNetflowExporterSiteFromTemplate(templateID, siteID, exporterUUID string) error {
	msoClient := testAccProvider.Meta().(*client.Client)
	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
	if err != nil {
		return fmt.Errorf("manual delete: get template %s: %w", templateID, err)
	}
	siteIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "siteId", siteID, "tenantPolicyTemplate", "sites")
	if err != nil {
		return fmt.Errorf("manual delete: locate site %q: %w", siteID, err)
	}
	siteCont := templateCont.S("tenantPolicyTemplate", "sites").Index(siteIndex)
	exporterIndex, err := GetPolicyIndexByKeyAndValue(siteCont, "ref", exporterUUID, "netFlowExporters")
	if err != nil {
		return fmt.Errorf("manual delete: locate NetFlow Exporter %q on site %q: %w", exporterUUID, siteID, err)
	}
	path := fmt.Sprintf("/tenantPolicyTemplate/sites/%d/netFlowExporters/%d", siteIndex, exporterIndex)
	if _, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateID), models.GetRemovePatchPayload(path)); err != nil {
		return fmt.Errorf("manual delete: remove NetFlow Exporter %q from site %q: %w", exporterUUID, siteID, err)
	}
	return nil
}

func testAccCheckMSOTenantPoliciesNetflowExporterSiteDestroy(s *terraform.State) error {
	msoClient := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mso_tenant_policies_netflow_exporter_site" {
			continue
		}
		templateID := rs.Primary.Attributes["template_id"]
		siteID := rs.Primary.Attributes["site_id"]
		exporterUUID := rs.Primary.Attributes["netflow_exporter_uuid"]
		cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
		if err != nil {
			continue
		}
		siteIndex, err := GetPolicyIndexByKeyAndValue(cont, "siteId", siteID, "tenantPolicyTemplate", "sites")
		if err != nil {
			continue
		}
		siteCont := cont.S("tenantPolicyTemplate", "sites").Index(siteIndex)
		if _, err := GetPolicyIndexByKeyAndValue(siteCont, "ref", exporterUUID, "netFlowExporters"); err == nil {
			return fmt.Errorf("NetFlow Exporter %q is still configured on site %q", exporterUUID, siteID)
		}
	}
	return nil
}
