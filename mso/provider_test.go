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
var siteNames = []string{"ansible_test"}
var tenantNames = []string{"acctest_crest"}
var validSchemaId = "6206831f1d000012864f99a8"
var inValidScheamaId = "620683151d0000f1854f99a4"

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
