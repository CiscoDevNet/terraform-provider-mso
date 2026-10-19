package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const msoTfTenantName = "tf_test_mso_tenant_app"

func TestAccNdoSchemaTemplateDeploy_Error(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Cross-template VRF/BD dependency (expecting deployment error)") },
				Config:      testAccMsoSchemaTemplateErrorCrossTemplateVrfBdConfig(),
				ExpectError: regexp.MustCompile(`Error on deploy:`),
			},
		},
	})
}

// TestAccNdoSchemaTemplateDeploy_OverlappingVlanError reproduces
// https://github.com/CiscoDevNet/terraform-provider-mso/issues/548: two EPGs
// statically bound to the same interface/pod/leaf with the same VLAN
// encapsulation cause the site deployment to fail. Previously the provider
// surfaced only the generic "all stretch fabrics failed to deploy" message;
// this test asserts that the underlying APIC validation error (e.g.
// "Validation failed" / "encapsulation") returned in
// operDetails.execSiteStatus is now included in the error message.
func TestAccNdoSchemaTemplateDeploy_OverlappingVlanError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Overlapping VLAN on two EPGs (expecting detailed deployment error)") },
				Config:      testAccNdoSchemaTemplateDeployOverlappingVlanConfig(),
				ExpectError: regexp.MustCompile(`(?i)Error on deploy:.*(validation failed|encapsulation)`),
			},
		},
	})
}

func TestAccNdoSchemaTemplateDeploy_WithCustomRetry(t *testing.T) {
	logFilePath := setupTestLogCapture(t, "TRACE")

	expectedLogs := []string{
		`\[TRACE\] Task status is \w+`,
		`\[DEBUG\] Custom retry function indicated a retry is needed for 2xx response`,
		`\[ERROR\] HTTP Request failed with status code 200, retrying\.\.\.`,
		`\[DEBUG\] Begining backoff method: attempts \d+ on \d+`,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Check Retry on Deploy") },
				Config:    testAccMsoSchemaTemplateVrfAndBdDeployWithRetry(),
				Check:     customTestCheckLogs(logFilePath, expectedLogs),
			},
		},
	})
}

func TestAccNdoSchemaTemplateDeploy_ValidationError_MissingSchemaId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Schema id validation") },
				Config:      testAccNdoSchemaTemplateDeploy_ErrorAppMissingSchemaId(),
				ExpectError: regexp.MustCompile("when 'template_id' is not provided, both 'schema_id' and 'template_name' must be set for template_type"),
			},
		},
	})
}

func TestAccNdoSchemaTemplateDeploy_ValidationError_MissingTemplateName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Template name validation") },
				Config:      testAccNdoSchemaTemplateDeploy_ErrorAppMissingTemplateName(),
				ExpectError: regexp.MustCompile("when 'template_id' is not provided, both 'schema_id' and 'template_name' must be set for template_type"),
			},
		},
	})
}

func TestAccNdoSchemaTemplateDeploy_ValidationError_NonAppMissingName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Template name validation with template_type tenant") },
				Config:      testAccNdoSchemaTemplateDeploy_ErrorNonAppMissingName(),
				ExpectError: regexp.MustCompile("when 'template_id' is not provided, 'template_name' must be set for template_type tenant"),
			},
		},
	})
}

func TestAccNdoSchemaTemplateDeploy_Success_WithTemplateId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Deploy with template id") },
				Config:    testAccNdoSchemaTemplateDeploy_IPSLAMonitoringPolicyWithTemplateId(),
				Check:     resource.TestCheckResourceAttrSet("mso_schema_template_deploy_ndo.deploy", "template_id"),
			},
		},
	})
}

func TestAccNdoSchemaTemplateDeploy_Success_WithoutTemplateId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Deploy without template id") },
				Config:    testAccNdoSchemaTemplateDeploy_IPSLAMonitoringPolicyWithoutTemplateId(),
				Check:     resource.TestCheckResourceAttrSet("mso_schema_template_deploy_ndo.deploy", "template_id"),
			},
		},
	})
}

func TestAccNdoSchemaTemplateDeploy_Undeploy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Deploy tenant template") },
				Config:    testAccNdoSchemaTemplateDeploy_TenantTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_deploy_ndo.deploy", "template_id"),
				),
			},
		},
	})
}

func TestAccNdoSchemaTemplateDeploy_UndeployWithSiteIds(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Deploy tenant template") },
				Config:    testAccNdoSchemaTemplateDeploy_TenantTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_deploy_ndo.deploy", "template_id"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Undeploy with explicit site_ids") },
				Config:    testAccNdoSchemaTemplateDeploy_TenantTemplateUndeployWithSiteIds(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_deploy_ndo.deploy", "undeploy", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_deploy_ndo.deploy", "site_ids.#", "1"),
				),
			},
		},
	})
}

func TestAccNdoSchemaTemplateDeploy_ValidationError_UndeployWithoutSiteIds(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Undeploy validation - site_ids required") },
				Config:      testAccNdoSchemaTemplateDeploy_ErrorUndeployWithoutSiteIds(),
				ExpectError: regexp.MustCompile("when 'undeploy=true', 'site_ids' must be provided"),
			},
		},
	})
}

func TestAccNdoSchemaTemplateDeploy_ApplicationTemplate_UndeployWithSiteIds(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Deploy application template") },
				Config:    testAccNdoSchemaTemplateDeploy_ApplicationTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_deploy_ndo.app_deploy", "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_deploy_ndo.app_deploy", "template_type", "application"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Undeploy application template with site_ids") },
				Config:    testAccNdoSchemaTemplateDeploy_ApplicationTemplateUndeployWithSiteIds(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_deploy_ndo.app_deploy", "undeploy", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_deploy_ndo.app_deploy", "site_ids.#", "1"),
				),
			},
		},
	})
}

func testAccSingleTenantConfig() string {
	return fmt.Sprintf(`
    %s
    resource "mso_tenant" "%s" {
        name = "%s"
        site_associations { 
            site_id = data.mso_site.%s.id 
        }
    }
    `, testSiteConfigAnsibleTest(), msoTfTenantName, msoTfTenantName, msoTemplateSiteName1)
}

func testAccMsoSchemaTemplateErrorCrossTemplateVrfBdConfig() string {
	return fmt.Sprintf(`%s
    resource "mso_schema" "schema_blocks" {
        name = "demo_schema_blocks"
        template {
            name         = "Template1"
            display_name = "TEMP1"
            tenant_id    = mso_tenant.%s.id
            template_type = "aci_multi_site"
        }
        template {
            name         = "Template2"
            display_name = "TEMP2"
            tenant_id    = mso_tenant.%s.id
            template_type = "aci_multi_site"
        }
    }

    resource "mso_schema_site" "schema_site_1" {
        schema_id     = mso_schema.schema_blocks.id
        site_id       = data.mso_site.%s.id
        template_name = tolist(mso_schema.schema_blocks.template)[0].name
    }

    resource "mso_schema_site" "schema_site_2" {
        schema_id     = mso_schema.schema_blocks.id
        site_id       = data.mso_site.%s.id
        template_name = tolist(mso_schema.schema_blocks.template)[1].name
    }

    resource "mso_schema_template_vrf" "vrf1" {
        schema_id       = mso_schema.schema_blocks.id
        template        = tolist(mso_schema.schema_blocks.template)[0].name
        name            = "vrf1"
        display_name    ="vrf"
        layer3_multicast=true
    }

    resource "mso_schema_template_bd" "bridgedomain" {
        schema_id              = mso_schema.schema_blocks.id
        template_name          = tolist(mso_schema.schema_blocks.template)[0].name
        name                   = "bd"
        display_name           = "test"
        vrf_name               = mso_schema_template_vrf.vrf1.name
        vrf_schema_id          = mso_schema.schema_blocks.id
        vrf_template_name      = tolist(mso_schema.schema_blocks.template)[0].name
        layer2_unknown_unicast = "proxy"
        intersite_bum_traffic  = false
        optimize_wan_bandwidth = true
        layer2_stretch         = true
        layer3_multicast       = true
    }

    resource "mso_schema_template_bd" "bridgedomain2" {
        schema_id              = mso_schema.schema_blocks.id
        template_name          = tolist(mso_schema.schema_blocks.template)[1].name
        name                   = "bd2"
        display_name           = "test"
        vrf_name               = mso_schema_template_vrf.vrf1.name
        vrf_schema_id          = mso_schema.schema_blocks.id
        vrf_template_name      = tolist(mso_schema.schema_blocks.template)[0].name
        layer2_unknown_unicast = "proxy"
        intersite_bum_traffic  = false
        optimize_wan_bandwidth = true
        layer2_stretch         = true
        layer3_multicast       = true
    }

    resource "mso_schema_template_deploy_ndo" "deploy_ndo2" {
        schema_id     = mso_schema_template_bd.bridgedomain2.schema_id
        template_name = tolist(mso_schema.schema_blocks.template)[1].name
    }
    `, testAccSingleTenantConfig(), msoTfTenantName, msoTfTenantName, msoTemplateSiteName1, msoTemplateSiteName1)
}

// testAccNdoSchemaTemplateDeployOverlappingVlanConfig builds a single-site
// schema with a VRF/BD and two EPGs, both statically bound (via
// mso_schema_site_anp_epg_static_port) to the same pod/leaf/path with the
// same VLAN. NDO/APIC rejects deploying two EPGs with an identical
// port+VLAN encapsulation, which is the scenario reported in issue #548.
func testAccNdoSchemaTemplateDeployOverlappingVlanConfig() string {
	epgName1 := msoSchemaTemplateAnpEpgName + "1"
	epgName2 := msoSchemaTemplateAnpEpgName + "2"
	return fmt.Sprintf(`%[1]s
    resource "mso_schema_template_anp" "%[2]s" {
        name          = "%[2]s"
        display_name  = "%[2]s"
        schema_id     = mso_schema.%[3]s.id
        template      = tolist(mso_schema.%[3]s.template)[0].name
    }

    resource "mso_schema_template_vrf" "%[4]s" {
        name         = "%[4]s"
        display_name = "%[4]s"
        schema_id    = mso_schema.%[3]s.id
        template     = tolist(mso_schema.%[3]s.template)[0].name
    }

    resource "mso_schema_template_bd" "%[5]s" {
        schema_id              = mso_schema.%[3]s.id
        template_name          = tolist(mso_schema.%[3]s.template)[0].name
        name                   = "%[5]s"
        display_name           = "%[5]s"
        vrf_name               = mso_schema_template_vrf.%[4]s.name
        vrf_schema_id          = mso_schema.%[3]s.id
        vrf_template_name      = tolist(mso_schema.%[3]s.template)[0].name
        layer2_unknown_unicast = "proxy"
    }

    resource "mso_schema_template_anp_epg" "epg1" {
        name          = "%[6]s"
        display_name  = "%[6]s"
        anp_name      = mso_schema_template_anp.%[2]s.name
        schema_id     = mso_schema.%[3]s.id
        template_name = tolist(mso_schema.%[3]s.template)[0].name
        bd_name       = mso_schema_template_bd.%[5]s.name
    }

    resource "mso_schema_template_anp_epg" "epg2" {
        name          = "%[7]s"
        display_name  = "%[7]s"
        anp_name      = mso_schema_template_anp.%[2]s.name
        schema_id     = mso_schema.%[3]s.id
        template_name = tolist(mso_schema.%[3]s.template)[0].name
        bd_name       = mso_schema_template_bd.%[5]s.name
    }

    resource "mso_schema_site_anp_epg" "epg1_site" {
        schema_id     = mso_schema.%[3]s.id
        site_id       = mso_schema_site.%[8]s.site_id
        template_name = tolist(mso_schema.%[3]s.template)[0].name
        anp_name      = mso_schema_template_anp.%[2]s.name
        epg_name      = mso_schema_template_anp_epg.epg1.name
    }

    resource "mso_schema_site_anp_epg" "epg2_site" {
        schema_id     = mso_schema.%[3]s.id
        site_id       = mso_schema_site.%[8]s.site_id
        template_name = tolist(mso_schema.%[3]s.template)[0].name
        anp_name      = mso_schema_template_anp.%[2]s.name
        epg_name      = mso_schema_template_anp_epg.epg2.name
    }

    resource "mso_schema_site_anp_epg_static_port" "epg1_port" {
        schema_id            = mso_schema.%[3]s.id
        site_id              = mso_schema_site.%[8]s.site_id
        template_name        = tolist(mso_schema.%[3]s.template)[0].name
        anp_name             = mso_schema_template_anp.%[2]s.name
        epg_name             = mso_schema_site_anp_epg.epg1_site.epg_name
        path_type            = "port"
        pod                  = "%[9]s"
        leaf                 = "%[10]s"
        path                 = "%[11]s"
        vlan                 = %[12]d
        deployment_immediacy = "immediate"
        mode                 = "regular"
    }

    resource "mso_schema_site_anp_epg_static_port" "epg2_port" {
        schema_id            = mso_schema.%[3]s.id
        site_id              = mso_schema_site.%[8]s.site_id
        template_name        = tolist(mso_schema.%[3]s.template)[0].name
        anp_name             = mso_schema_template_anp.%[2]s.name
        epg_name             = mso_schema_site_anp_epg.epg2_site.epg_name
        path_type            = "port"
        pod                  = "%[9]s"
        leaf                 = "%[10]s"
        path                 = "%[11]s"
        vlan                 = %[12]d
        deployment_immediacy = "immediate"
        mode                 = "regular"
    }

    resource "mso_schema_template_deploy_ndo" "deploy_overlap" {
        schema_id     = mso_schema.%[3]s.id
        template_name = tolist(mso_schema.%[3]s.template)[0].name
        force_apply   = ""
        depends_on = [
            mso_schema_site_anp_epg_static_port.epg1_port,
            mso_schema_site_anp_epg_static_port.epg2_port,
        ]
    }
    `,
		testSchemaWithSingleSiteAssociationConfig(),
		msoSchemaTemplateAnpName,
		msoSchemaName,
		msoSchemaTemplateVrfName,
		msoSchemaTemplateBdName,
		epgName1,
		epgName2,
		msoSchemaSiteResourceLabel1,
		msoSchemaSiteAnpEpgStaticPortPod,
		msoSchemaSiteAnpEpgStaticPortLeaf,
		msoSchemaSiteAnpEpgStaticPortPath,
		999,
	)
}

func testAccMsoSchemaTemplateVrfAndBdDeployWithRetry() string {
	return fmt.Sprintf(`%s
    resource "mso_schema" "schema_blocks" {
        name = "demo_schema_blocks"
        template {
            name         = "Template1"
            display_name = "TEMP1"
            tenant_id    = mso_tenant.%s.id
            template_type = "aci_multi_site"
        }
    }

    resource "mso_schema_site" "schema_site_1" {
        schema_id     = mso_schema.schema_blocks.id
        site_id       = data.mso_site.%s.id
        template_name = tolist(mso_schema.schema_blocks.template)[0].name
        undeploy_on_destroy = true
    }

    resource "mso_schema_template_vrf" "vrf" {
        count = 50
        schema_id       = mso_schema.schema_blocks.id
        template        = tolist(mso_schema.schema_blocks.template)[0].name
        name            = "vrf${count.index + 1}"
        display_name    = "VRF-${count.index + 1}"
        layer3_multicast=true
      }

      resource "mso_schema_template_bd" "bridgedomain" {
          schema_id              = mso_schema.schema_blocks.id
          template_name          = tolist(mso_schema.schema_blocks.template)[0].name
          name                   = "bd"
          display_name           = "test"
          vrf_name               = mso_schema_template_vrf.vrf[0].name
          vrf_schema_id          = mso_schema.schema_blocks.id
          vrf_template_name      = tolist(mso_schema.schema_blocks.template)[0].name
          layer2_unknown_unicast = "proxy" 
          intersite_bum_traffic  = false
          optimize_wan_bandwidth = true
          layer2_stretch         = true
          layer3_multicast       = true  
    }

    resource "mso_schema_template_deploy_ndo" "deploy_ndo" {
        force_apply = ""
        schema_id     = mso_schema_template_bd.bridgedomain.schema_id
        template_name = tolist(mso_schema.schema_blocks.template)[0].name
    }
    `, testAccSingleTenantConfig(), msoTfTenantName, msoTemplateSiteName1)
}

func testAccNdoSchemaTemplateDeploy_ErrorAppMissingSchemaId() string {
	return `
    resource "mso_schema_template_deploy_ndo" "deploy_error" {
        template_name = "test_name"
        
    }
    `
}

func testAccNdoSchemaTemplateDeploy_ErrorAppMissingTemplateName() string {
	return `
    resource "mso_schema_template_deploy_ndo" "deploy_error" {
        schema_id = "schema_id"
    }
    `
}

func testAccNdoSchemaTemplateDeploy_ErrorNonAppMissingName() string {
	return `
    resource "mso_schema_template_deploy_ndo" "deploy_error" {
        template_type = "tenant"
    }
    `
}

func testAccNdoSchemaTemplateDeploy_IPSLAMonitoringPolicyWithTemplateId() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_ipsla_monitoring_policy" "ipsla_policy" {
        template_id        = mso_template.template_tenant.id
        name               = "test_ipsla_policy"
        description        = "HTTP Type"
        sla_type           = "http"
        destination_port   = 80
        http_version       = "HTTP11"
        http_uri           = "/example"
        sla_frequency      = 120
        detect_multiplier  = 4
        request_data_size  = 64
        type_of_service    = 18
        operation_timeout  = 100
        threshold          = 100
        ipv6_traffic_class = 255
    }

    resource "mso_schema_template_deploy_ndo" "deploy" {
        force_apply = ""
        template_id = mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy.template_id
        template_type = "tenant"
        undeploy_on_destroy = true
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccNdoSchemaTemplateDeploy_IPSLAMonitoringPolicyWithoutTemplateId() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_ipsla_monitoring_policy" "ipsla_policy" {
        template_id        = mso_template.template_tenant.id
        name               = "test_ipsla_policy2"
        description        = "HTTP Type"
        sla_type           = "http"
        destination_port   = 80
        http_version       = "HTTP11"
        http_uri           = "/example"
        sla_frequency      = 120
        detect_multiplier  = 4
        request_data_size  = 64
        type_of_service    = 18
        operation_timeout  = 100
        threshold          = 100
        ipv6_traffic_class = 255
    }

    resource "mso_schema_template_deploy_ndo" "deploy" {
        depends_on = [mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy]
        force_apply = ""
        template_name = mso_template.template_tenant.template_name
        template_type = "tenant"
        undeploy_on_destroy = true
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccNdoSchemaTemplateDeploy_TenantTemplate() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_ipsla_monitoring_policy" "ipsla_policy" {
        template_id        = mso_template.template_tenant.id
        name               = "test_ipsla_undeploy"
        description        = "HTTP Type"
        sla_type           = "http"
        destination_port   = 80
        http_version       = "HTTP11"
        http_uri           = "/example"
        sla_frequency      = 120
        detect_multiplier  = 4
        request_data_size  = 64
        type_of_service    = 18
        operation_timeout  = 100
        threshold          = 100
        ipv6_traffic_class = 255
    }

    resource "mso_schema_template_deploy_ndo" "deploy" {
        depends_on = [mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy]
        force_apply = ""
        template_name = mso_template.template_tenant.template_name
        template_type = "tenant"
        undeploy_on_destroy = true
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccNdoSchemaTemplateDeploy_TenantTemplateUndeployWithSiteIds() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_ipsla_monitoring_policy" "ipsla_policy" {
        template_id        = mso_template.template_tenant.id
        name               = "test_ipsla_undeploy"
        description        = "HTTP Type"
        sla_type           = "http"
        destination_port   = 80
        http_version       = "HTTP11"
        http_uri           = "/example"
        sla_frequency      = 120
        detect_multiplier  = 4
        request_data_size  = 64
        type_of_service    = 18
        operation_timeout  = 100
        threshold          = 100
        ipv6_traffic_class = 255
    }

    resource "mso_schema_template_deploy_ndo" "deploy" {
        depends_on = [mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy]
        force_apply = ""
        template_name = mso_template.template_tenant.template_name
        template_type = "tenant"
        site_ids = [data.mso_site.%s.id]
        undeploy = true
    }`, testAccMSOTemplateResourceTenantConfig(), msoTemplateSiteName1)
}

func testAccNdoSchemaTemplateDeploy_ErrorUndeployWithoutSiteIds() string {
	return `
    resource "mso_schema_template_deploy_ndo" "deploy_error" {
        template_id   = "test_template_id"
        template_type = "tenant"
        undeploy      = true
    }
    `
}

func testAccNdoSchemaTemplateDeploy_ApplicationTemplate() string {
	return fmt.Sprintf(`%s
    resource "mso_schema" "app_schema" {
        name = "test_app_schema"
        template {
            name          = "AppTemplate1"
            display_name  = "Application Template 1"
            tenant_id     = mso_tenant.%s.id
            template_type = "aci_multi_site"
        }
    }

    resource "mso_schema_site" "app_site" {
        schema_id     = mso_schema.app_schema.id
        site_id       = data.mso_site.%s.id
        template_name = tolist(mso_schema.app_schema.template)[0].name
    }

    resource "mso_schema_template_deploy_ndo" "app_deploy" {
        force_apply = ""
        schema_id     = mso_schema.app_schema.id
        template_name = tolist(mso_schema.app_schema.template)[0].name
    }`, testAccSingleTenantConfig(), msoTfTenantName, msoTemplateSiteName1)
}

func testAccNdoSchemaTemplateDeploy_ApplicationTemplateUndeployWithSiteIds() string {
	return fmt.Sprintf(`%s
    resource "mso_schema" "app_schema" {
        name = "test_app_schema"
        template {
            name          = "AppTemplate1"
            display_name  = "Application Template 1"
            tenant_id     = mso_tenant.%s.id
            template_type = "aci_multi_site"
        }
    }

    resource "mso_schema_site" "app_site" {
        schema_id     = mso_schema.app_schema.id
        site_id       = data.mso_site.%s.id
        template_name = tolist(mso_schema.app_schema.template)[0].name
    }

    resource "mso_schema_template_deploy_ndo" "app_deploy" {
        force_apply = ""
        schema_id     = mso_schema_site.app_site.schema_id
        template_name = tolist(mso_schema.app_schema.template)[0].name
        undeploy      = true
        site_ids      = [data.mso_site.%s.id]
    }`, testAccSingleTenantConfig(), msoTfTenantName, msoTemplateSiteName1, msoTemplateSiteName1)
}
