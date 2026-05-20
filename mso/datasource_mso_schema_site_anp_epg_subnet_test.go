package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccMSOSchemaSiteAnpEpgSubnetDatasource verifies the read-only datasource
// for mso_schema_site_anp_epg_subnet:
//   - lookup with an IP that does not exist returns an error
//   - successful read returns the expected attribute set and all values pair
//     with the managing resource
func TestAccMSOSchemaSiteAnpEpgSubnetDatasource(t *testing.T) {
	subnetResource := "mso_schema_site_anp_epg_subnet." + msoSchemaTemplateAnpEpgName
	subnetDatasource := "data.mso_schema_site_anp_epg_subnet." + msoSchemaTemplateAnpEpgName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpEpgSubnetDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("DataSource: Lookup non-existent subnet IP (expect error)")
				},
				Config:      testAccMSOSchemaSiteAnpEpgSubnetDatasourceNotFoundConfig(),
				ExpectError: regexp.MustCompile(`Unable to find subnet entry with ip`),
			},
			{
				PreConfig: func() { fmt.Println("DataSource: Read existing subnet") },
				Config:    testAccMSOSchemaSiteAnpEpgSubnetDatasourceReadConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(subnetDatasource, "schema_id", subnetResource, "schema_id"),
					resource.TestCheckResourceAttrPair(subnetDatasource, "site_id", subnetResource, "site_id"),
					resource.TestCheckResourceAttrPair(subnetDatasource, "template_name", subnetResource, "template_name"),
					resource.TestCheckResourceAttrPair(subnetDatasource, "anp_name", subnetResource, "anp_name"),
					resource.TestCheckResourceAttrPair(subnetDatasource, "epg_name", subnetResource, "epg_name"),
					resource.TestCheckResourceAttrPair(subnetDatasource, "ip", subnetResource, "ip"),
					resource.TestCheckResourceAttrPair(subnetDatasource, "scope", subnetResource, "scope"),
					resource.TestCheckResourceAttrPair(subnetDatasource, "shared", subnetResource, "shared"),
					resource.TestCheckResourceAttrPair(subnetDatasource, "no_default_gateway", subnetResource, "no_default_gateway"),
					resource.TestCheckResourceAttrPair(subnetDatasource, "querier", subnetResource, "querier"),
					resource.TestCheckResourceAttrPair(subnetDatasource, "primary", subnetResource, "primary"),
					resource.TestCheckResourceAttrPair(subnetDatasource, "description", subnetResource, "description"),
				),
			},
		},
	})
}

// testAccMSOSchemaSiteAnpEpgSubnetDatasourceNotFoundConfig creates the subnet
// resource then queries the datasource for an IP that does not exist,
// exercising the "Unable to find subnet entry with ip" error path in
// datasourceMSOSchemaSiteAnpEpgSubnetRead.
func testAccMSOSchemaSiteAnpEpgSubnetDatasourceNotFoundConfig() string {
	return fmt.Sprintf(`%[1]s
	data "mso_schema_site_anp_epg_subnet" "missing" {
		schema_id     = mso_schema.%[2]s.id
		site_id       = mso_schema_site.%[3]s.site_id
		template_name = "%[4]s"
		anp_name      = mso_schema_template_anp.%[5]s.name
		epg_name      = mso_schema_site_anp_epg_subnet.%[6]s.epg_name
		ip            = "192.0.2.1/32"

		depends_on = [mso_schema_site_anp_epg_subnet.%[6]s]
	}`,
		testAccMSOSchemaSiteAnpEpgSubnetConfigCreate(),
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
		msoSchemaTemplateAnpName,
		msoSchemaTemplateAnpEpgName,
	)
}

// testAccMSOSchemaSiteAnpEpgSubnetDatasourceReadConfig creates the subnet
// resource and reads it back via the datasource, verifying all computed
// attributes are populated correctly.
func testAccMSOSchemaSiteAnpEpgSubnetDatasourceReadConfig() string {
	return fmt.Sprintf(`%[1]s
	data "mso_schema_site_anp_epg_subnet" "%[2]s" {
		schema_id     = mso_schema_site_anp_epg_subnet.%[2]s.schema_id
		site_id       = mso_schema_site_anp_epg_subnet.%[2]s.site_id
		template_name = mso_schema_site_anp_epg_subnet.%[2]s.template_name
		anp_name      = mso_schema_site_anp_epg_subnet.%[2]s.anp_name
		epg_name      = mso_schema_site_anp_epg_subnet.%[2]s.epg_name
		ip            = mso_schema_site_anp_epg_subnet.%[2]s.ip
	}`,
		testAccMSOSchemaSiteAnpEpgSubnetConfigCreate(),
		msoSchemaTemplateAnpEpgName,
	)
}
