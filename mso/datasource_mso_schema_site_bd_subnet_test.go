package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccMSOSchemaSiteBdSubnetDatasource verifies the read-only datasource for
// mso_schema_site_bd_subnet:
//   - lookup with an IP that does not exist returns an error
//   - successful read returns the expected attribute set and all values pair
//     with the managing resource
func TestAccMSOSchemaSiteBdSubnetDatasource(t *testing.T) {
	siteBdSubnetResource := "mso_schema_site_bd_subnet." + msoSchemaTemplateBdName
	siteBdSubnetDatasource := "data.mso_schema_site_bd_subnet." + msoSchemaTemplateBdName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteBdSubnetDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("DataSource: Lookup non-existent site BD subnet IP (expect error)")
				},
				Config:      testAccMSOSchemaSiteBdSubnetDatasourceNotFoundConfig(),
				ExpectError: regexp.MustCompile(`Unable to find BD subnet entry with ip`),
			},
			{
				PreConfig: func() { fmt.Println("DataSource: Read existing site BD subnet") },
				Config:    testAccMSOSchemaSiteBdSubnetDatasourceReadConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(siteBdSubnetDatasource, "schema_id", siteBdSubnetResource, "schema_id"),
					resource.TestCheckResourceAttrPair(siteBdSubnetDatasource, "site_id", siteBdSubnetResource, "site_id"),
					resource.TestCheckResourceAttrPair(siteBdSubnetDatasource, "template_name", siteBdSubnetResource, "template_name"),
					resource.TestCheckResourceAttrPair(siteBdSubnetDatasource, "bd_name", siteBdSubnetResource, "bd_name"),
					resource.TestCheckResourceAttrPair(siteBdSubnetDatasource, "ip", siteBdSubnetResource, "ip"),
					resource.TestCheckResourceAttrPair(siteBdSubnetDatasource, "scope", siteBdSubnetResource, "scope"),
					resource.TestCheckResourceAttrPair(siteBdSubnetDatasource, "shared", siteBdSubnetResource, "shared"),
					resource.TestCheckResourceAttrPair(siteBdSubnetDatasource, "no_default_gateway", siteBdSubnetResource, "no_default_gateway"),
				),
			},
		},
	})
}

// testAccMSOSchemaSiteBdSubnetDatasourceNotFoundConfig creates the subnet
// resource then queries the datasource for an IP that does not exist,
// exercising the "Unable to find BD subnet entry" error path.
func testAccMSOSchemaSiteBdSubnetDatasourceNotFoundConfig() string {
	return fmt.Sprintf(`%[1]s
	data "mso_schema_site_bd_subnet" "missing" {
		schema_id     = mso_schema.%[2]s.id
		site_id       = mso_schema_site.%[3]s.site_id
		template_name = "%[4]s"
		bd_name       = mso_schema_site_bd_subnet.%[5]s.bd_name
		ip            = "192.168.99.1/24"

		depends_on = [mso_schema_site_bd_subnet.%[5]s]
	}`,
		testAccMSOSchemaSiteBdSubnetConfigUpdate(),
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
		msoSchemaTemplateBdName,
	)
}

// testAccMSOSchemaSiteBdSubnetDatasourceReadConfig creates the subnet resource
// (with updated attributes) and reads it back via the datasource.
func testAccMSOSchemaSiteBdSubnetDatasourceReadConfig() string {
	return fmt.Sprintf(`%[1]s
	data "mso_schema_site_bd_subnet" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = mso_schema_site.%[4]s.site_id
		template_name = "%[5]s"
		bd_name       = mso_schema_site_bd_subnet.%[2]s.bd_name
		ip            = mso_schema_site_bd_subnet.%[2]s.ip
	}`,
		testAccMSOSchemaSiteBdSubnetConfigUpdate(),
		msoSchemaTemplateBdName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
	)
}
