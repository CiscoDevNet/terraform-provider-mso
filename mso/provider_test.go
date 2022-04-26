package mso

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var siteNames = []string{"ansible_test", "azure_ansible_test"}
var tenantNames = []string{"acctest_crest", "crest_test_azure"}
var validSchemaId = "6206831f1d000012864f99a8"
var inValidScheamaId = "620683151d0000f1854f99a4"
var epg = "/schemas/621392f81d0000282a4f9d1c/templates/ACC_CREST/anps/UntitledAP1/epgs/test_epg"

func CreateServiceNodeResource(name string) string {
	resource := fmt.Sprintf(`
	resource "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
		  
	resource "mso_schema" "test" {
		name          = "%s"
		template_name = "%s"
		tenant_id     = mso_tenant.test.id
	}
	
	resource "mso_schema_template_service_graph" "test" {
	  schema_id          = mso_schema.test.id
	  template_name      = mso_schema.test.template_name
	  service_graph_name = "%s"
	  service_node_type  = "firewall"
	}
	
	resource "mso_schema_site_service_graph_node" "test" {
	  schema_id          = mso_schema_template_service_graph.test.schema_id
	  template_name      = mso_schema_template_service_graph.test.template_name
	  service_graph_name = mso_schema_template_service_graph.test.service_graph_name
	  service_node_type  = "firewall"
	}
	`, name, name, name, name, name)
	return resource
}

func CreatSchemaSiteConfig(site, tenant, name string) string {
	resource := fmt.Sprintf(`
	data "mso_site" "test" {
		name  = "%s"
	}
	  
 	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	  
	resource "mso_schema" "test" {
		name          = "%s"
		template_name = "%s"
		tenant_id     = data.mso_tenant.test.id
	}
			
	resource "mso_schema_site" "test" {
		schema_id       =  mso_schema.test.id
		site_id         =  data.mso_site.test.id
		template_name   =  "%s"
	}
	`, site, tenant, tenant, name, name, name)
	return resource
}

func CreateDHCPRelayPolicy(tenant, polname string) string {
	resource := fmt.Sprintf(`
	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_relay_policy" "test" {
		tenant_id = data.mso_tenant.test.id
		name = "%s"		
	}
	`, tenant, tenant, polname)
	return resource
}

func GetParentConfigBDDHCPPolicy(tenant, schema, name string) string {
	resource := fmt.Sprintf(`
	data "mso_tenant" "test" {
		name         = "%s"
		display_name = "%s"
	}
	  
	 resource "mso_schema" "test" {
		name          = "%s"
		template_name = "%s"
		tenant_id     = data.mso_tenant.test.id
	}
	  
	resource "mso_schema_template_vrf" "test" {
		schema_id        = mso_schema.test.id
		template         = mso_schema.test.template_name
		name             = "%s"
		display_name     = "%s"
	}
	  
	resource "mso_schema_template_bd" "test" {
		schema_id              = mso_schema.test.id
		template_name          = mso_schema.test.template_name
		name                   = "%s"
		display_name           = "%s"
		vrf_name               = mso_schema_template_vrf.test.name
	}
	  
	resource "mso_dhcp_relay_policy" "test" {
		tenant_id   = data.mso_tenant.test.id
		name        = "%s"
	}
	`, tenant, tenant, schema, schema, name, name, name, name, name)
	return resource
}

func CreateSchemaSiteAnpEpgConfig(site, tenant, templateName, anpName, epgName string) string {
	resource := CreatSchemaSiteConfig(site, tenant, templateName)
	resource += fmt.Sprintf(`
	resource "mso_schema_site_anp" "example" {
		schema_id     = mso_schema.example.id
		anp_name      = %s
		template_name = %s
		site_id       = data.mso_site.example.id
	}

	resource "mso_schema_site_anp_epg" "test" {
		schema_id     = mso_schema.test.id
		template_name = %s
		site_id       = data.mso_site.test.id
		anp_name      = mso_schema_site_anp.example.anp_name
		epg_name      = %s
	}
	`, anpName, templateName, templateName, epgName)
	return resource
}

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"mso": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	// We will use this function later on to make sure our test environment is valid.
	// For example, you can make sure here that some environment variables are set.
	if v := os.Getenv("MSO_USERNAME"); v == "" {
		t.Fatal("username variable must be set for acceptance tests")
	}

	if v := os.Getenv("MSO_PASSWORD"); v == "" {

		t.Fatal("password variable must be set for acceptance tests")
	}
	if v := os.Getenv("MSO_URL"); v == "" {
		t.Fatal("url variable must be set for acceptance tests")
	}

}
