package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccMSOSchemaSiteBdL3outDatasource verifies the read-only datasource for
// mso_schema_site_bd_l3out:
//   - lookup with an l3out name that does not exist returns an error
//   - successful read returns the expected attribute set, including
//     l3out_schema_id and l3out_template_name
func TestAccMSOSchemaSiteBdL3outDatasource(t *testing.T) {
	siteBdL3outResource := "mso_schema_site_bd_l3out." + msoSchemaTemplateBdName
	siteBdL3outDatasource := "data.mso_schema_site_bd_l3out." + msoSchemaTemplateBdName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteBdL3outDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("DataSource: Lookup non-existent site BD L3out (expect error)")
				},
				Config:      testAccMSOSchemaSiteBdL3outDatasourceNotFoundConfig(),
				ExpectError: regexp.MustCompile(`Unable to find the Site BD L3out`),
			},
			{
				PreConfig: func() { fmt.Println("DataSource: Read existing site BD L3out") },
				Config:    testAccMSOSchemaSiteBdL3outDatasourceReadConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(siteBdL3outDatasource, "schema_id", siteBdL3outResource, "schema_id"),
					resource.TestCheckResourceAttrPair(siteBdL3outDatasource, "site_id", siteBdL3outResource, "site_id"),
					resource.TestCheckResourceAttrPair(siteBdL3outDatasource, "template_name", siteBdL3outResource, "template_name"),
					resource.TestCheckResourceAttrPair(siteBdL3outDatasource, "bd_name", siteBdL3outResource, "bd_name"),
					resource.TestCheckResourceAttrPair(siteBdL3outDatasource, "l3out_name", siteBdL3outResource, "l3out_name"),
					resource.TestCheckResourceAttr(siteBdL3outDatasource, "l3out_template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttrSet(siteBdL3outDatasource, "l3out_schema_id"),
				),
			},
		},
	})
}

// testAccMSOSchemaSiteBdL3outDatasourceNotFoundConfig stands up the prereqs
// and the resource, then queries the datasource for an l3out name that is not
// present, exercising the "Unable to find the Site BD L3out" error path.
func testAccMSOSchemaSiteBdL3outDatasourceNotFoundConfig() string {
	return fmt.Sprintf(`%[1]s
	data "mso_schema_site_bd_l3out" "missing" {
		schema_id     = mso_schema.%[2]s.id
		site_id       = mso_schema_site.%[3]s.site_id
		template_name = "%[4]s"
		bd_name       = mso_schema_site_bd_l3out.%[5]s.bd_name
		l3out_name    = "does_not_exist"

		depends_on = [mso_schema_site_bd_l3out.%[5]s]
	}`,
		testAccMSOSchemaSiteBdL3outConfigCreate(),
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
		msoSchemaTemplateBdName,
	)
}

// testAccMSOSchemaSiteBdL3outDatasourceReadConfig creates the site BD L3out
// resource and reads it back via the datasource.
func testAccMSOSchemaSiteBdL3outDatasourceReadConfig() string {
	return fmt.Sprintf(`%[1]s
	data "mso_schema_site_bd_l3out" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = mso_schema_site.%[4]s.site_id
		template_name = "%[5]s"
		bd_name       = mso_schema_site_bd_l3out.%[2]s.bd_name
		l3out_name    = mso_schema_site_bd_l3out.%[2]s.l3out_name
	}`,
		testAccMSOSchemaSiteBdL3outConfigCreate(),
		msoSchemaTemplateBdName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
	)
}
