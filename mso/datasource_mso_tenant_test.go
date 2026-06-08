package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMSOTenantDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMsoTenantDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Read Tenant datasource not found error") },
				Config:      testAccMSOTenantDatasourceNotFound(),
				ExpectError: regexp.MustCompile("Tenant of specified name not found"),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read Tenant datasource") },
				Config:    testAccMSOTenantDatasource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_tenant.tenant", "name", msoTenantName),
					resource.TestCheckResourceAttr("data.mso_tenant.tenant", "display_name", msoTenantName),
					resource.TestCheckResourceAttr("data.mso_tenant.tenant", "description", "Terraform test tenant"),
					resource.TestCheckResourceAttr("data.mso_tenant.tenant", "site_associations.#", "1"),
				),
			},
		},
	})
}

func testAccMSOTenantDatasource() string {
	return fmt.Sprintf(`%s
	resource "mso_tenant" "tenant" {
		name         = "%s"
		display_name = "%s"
		description  = "Terraform test tenant"
		site_associations {
			site_id = data.mso_site.%s.id
		}
	}
	data "mso_tenant" "tenant" {
		name = mso_tenant.tenant.name
	}`, testSiteConfigAnsibleTest(), msoTenantName, msoTenantName, msoTemplateSiteName1)
}

func testAccMSOTenantDatasourceNotFound() string {
	return `
	data "mso_tenant" "tenant" {
		name = "non_existing_tenant_name"
	}`
}
