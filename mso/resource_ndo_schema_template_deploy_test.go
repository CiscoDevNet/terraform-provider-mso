package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const msoTfTenantName = "tf_test_mso_tenant"

func TestAccNdoSchemaTemplateDeploy_Error(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Cross-template VRF/BD dependency (expecting deployment error)") },
				Config:      testAccMsoSchemaTemplateErrorCrossTemplateVrfBdConfig(),
				ExpectError: regexp.MustCompile(`^errors during apply: Error on deploy:`),
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
			{
				PreConfig: func() { fmt.Println("Test: Undeploy tenant template") },
				Config:    testAccNdoSchemaTemplateDeploy_TenantTemplateUndeploy(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_deploy_ndo.deploy", "template_id"),
					resource.TestCheckResourceAttr("mso_schema_template_deploy_ndo.deploy", "undeploy", "true"),
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

func testAccSingleTenantConfig() string {
	return fmt.Sprintf(`
    %s
    resource "mso_tenant" "%s" {
        name = "%s"
        display_name = "%s"
        site_associations { 
            site_id = data.mso_site.%s.id 
        }
    }
    `, testSiteConfigAnsibleTest(), msoTfTenantName, msoTfTenantName, msoTfTenantName, msoTemplateSiteName1)
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

func testAccNdoSchemaTemplateDeploy_TenantTemplateUndeploy() string {
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
        undeploy = true
    }`, testAccMSOTemplateResourceTenantConfig())
}
