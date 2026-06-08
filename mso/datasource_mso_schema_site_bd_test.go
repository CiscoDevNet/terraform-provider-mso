package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMSOSchemaSiteBdDatasource(t *testing.T) {
	siteBdResource := "mso_schema_site_bd." + msoSchemaTemplateBdName
	siteBdDatasource := "data.mso_schema_site_bd.bd"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteBdDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Read site BD datasource not-found error") },
				Config:      testAccMSOSchemaSiteBdDatasourceNotFound(),
				ExpectError: regexp.MustCompile(`BD .* is not found in Site`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read site BD datasource") },
				Config:    testAccMSOSchemaSiteBdDatasource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(siteBdDatasource, "schema_id"),
					resource.TestCheckResourceAttrSet(siteBdDatasource, "site_id"),
					resource.TestCheckResourceAttr(siteBdDatasource, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr(siteBdDatasource, "bd_name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttrPair(siteBdDatasource, "schema_id", siteBdResource, "schema_id"),
					resource.TestCheckResourceAttrPair(siteBdDatasource, "site_id", siteBdResource, "site_id"),
					resource.TestCheckResourceAttrPair(siteBdDatasource, "template_name", siteBdResource, "template_name"),
					resource.TestCheckResourceAttrPair(siteBdDatasource, "bd_name", siteBdResource, "bd_name"),
					resource.TestCheckResourceAttrPair(siteBdDatasource, "host_route", siteBdResource, "host_route"),
					resource.TestCheckResourceAttrPair(siteBdDatasource, "svi_mac", siteBdResource, "svi_mac"),
				),
			},
		},
	})
}

func testAccMSOSchemaSiteBdDatasource() string {
	return fmt.Sprintf(`%[1]s
	data "mso_schema_site_bd" "bd" {
		schema_id     = mso_schema_site_bd.%[2]s.schema_id
		site_id       = mso_schema_site_bd.%[2]s.site_id
		template_name = mso_schema_site_bd.%[2]s.template_name
		bd_name       = mso_schema_site_bd.%[2]s.bd_name
	}`,
		testAccMSOSchemaSiteBdConfigUpdateHostRouteAndSviMac(),
		msoSchemaTemplateBdName,
	)
}

// testAccMSOSchemaSiteBdDatasourceNotFound queries the datasource for a BD
// name that does not exist on the site, exercising the getSiteBd miss path
// ("BD <name> is not found in Site.").
func testAccMSOSchemaSiteBdDatasourceNotFound() string {
	return fmt.Sprintf(`%[1]s
	data "mso_schema_site_bd" "bd" {
		schema_id     = mso_schema_site_bd.%[2]s.schema_id
		site_id       = mso_schema_site_bd.%[2]s.site_id
		template_name = mso_schema_site_bd.%[2]s.template_name
		bd_name       = "non_existing_bd_name"
	}`,
		testAccMSOSchemaSiteBdConfigUpdateHostRouteAndSviMac(),
		msoSchemaTemplateBdName,
	)
}
