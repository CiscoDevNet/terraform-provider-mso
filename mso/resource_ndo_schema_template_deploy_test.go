package mso

import (
	"fmt"
	"os"
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
	logFile, err := os.CreateTemp("", "tf-acc-test-*.log")
	if err != nil {
		t.Fatalf("Failed to create temp log file: %v", err)
	}

	t.Cleanup(func() {
		logFile.Close()
		os.Remove(logFile.Name())
	})

	t.Setenv("TF_LOG", "TRACE")
	t.Setenv("TF_LOG_PATH", logFile.Name())

	expectedLogs := []string{
		`\[TRACE\] Task status is \w+`,
		`\[DEBUG\] Custom retry function indicated a retry is needed for 2xx response`,
		`\[ERROR\] HTTP Request failed with status code 200, retrying\.\.\.`,
		`\[DEBUG\] Begining backoff method: attempts \d+ on \d+`,
		`\[DEBUG\] Exit from backoff method with return value false`,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Check Retry on Deploy") },
				Config:    testAccMsoSchemaTemplateVrfAndBdDeployWithRetry(),
				Check:     customTestCheckLogs(logFile.Name(), expectedLogs),
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
