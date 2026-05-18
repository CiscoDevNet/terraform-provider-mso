package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccMSOSchemaSiteAnpEpgStaticLeafDatasource verifies the read-only
// datasource for mso_schema_site_anp_epg_static_leaf:
//   - lookup with a path that does not exist returns an error
//   - successful read returns the expected attribute set and all values pair
//     with the managing resource
func TestAccMSOSchemaSiteAnpEpgStaticLeafDatasource(t *testing.T) {
	staticLeafResource := "mso_schema_site_anp_epg_static_leaf." + msoSchemaTemplateAnpEpgName
	staticLeafDatasource := "data.mso_schema_site_anp_epg_static_leaf." + msoSchemaTemplateAnpEpgName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpEpgStaticLeafDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("DataSource: Lookup non-existent static leaf path (expect error)")
				},
				Config:      testAccMSOSchemaSiteAnpEpgStaticLeafDatasourceNotFoundConfig(),
				ExpectError: regexp.MustCompile(`Unable to find the Site ANP EPG Static Leaf`),
			},
			{
				PreConfig: func() { fmt.Println("DataSource: Read existing static leaf") },
				Config:    testAccMSOSchemaSiteAnpEpgStaticLeafDatasourceReadConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(staticLeafDatasource, "schema_id", staticLeafResource, "schema_id"),
					resource.TestCheckResourceAttrPair(staticLeafDatasource, "site_id", staticLeafResource, "site_id"),
					resource.TestCheckResourceAttrPair(staticLeafDatasource, "template_name", staticLeafResource, "template_name"),
					resource.TestCheckResourceAttrPair(staticLeafDatasource, "anp_name", staticLeafResource, "anp_name"),
					resource.TestCheckResourceAttrPair(staticLeafDatasource, "epg_name", staticLeafResource, "epg_name"),
					resource.TestCheckResourceAttrPair(staticLeafDatasource, "path", staticLeafResource, "path"),
					resource.TestCheckResourceAttrPair(staticLeafDatasource, "port_encap_vlan", staticLeafResource, "port_encap_vlan"),
				),
			},
		},
	})
}

// testAccMSOSchemaSiteAnpEpgStaticLeafDatasourceNotFoundConfig creates the
// static leaf resource then queries the datasource for a path that does not
// exist, exercising the "Unable to find the Site ANP EPG Static Leaf" error
// path in dataSourceMSOSchemaSiteAnpEpgStaticleafRead.
func testAccMSOSchemaSiteAnpEpgStaticLeafDatasourceNotFoundConfig() string {
	return fmt.Sprintf(`%[1]s
	data "mso_schema_site_anp_epg_static_leaf" "missing" {
		schema_id     = mso_schema.%[2]s.id
		site_id       = mso_schema_site.%[3]s.site_id
		template_name = "%[4]s"
		anp_name      = mso_schema_template_anp.%[5]s.name
		epg_name      = mso_schema_site_anp_epg_static_leaf.%[6]s.epg_name
		path          = "topology/pod-1/node-999"

		depends_on = [mso_schema_site_anp_epg_static_leaf.%[6]s]
	}`,
		testAccMSOSchemaSiteAnpEpgStaticLeafConfigCreate(),
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
		msoSchemaTemplateAnpName,
		msoSchemaTemplateAnpEpgName,
	)
}

// testAccMSOSchemaSiteAnpEpgStaticLeafDatasourceReadConfig creates the static
// leaf resource and reads it back via the datasource, verifying all computed
// attributes are populated correctly.
func testAccMSOSchemaSiteAnpEpgStaticLeafDatasourceReadConfig() string {
	return fmt.Sprintf(`%[1]s
	data "mso_schema_site_anp_epg_static_leaf" "%[2]s" {
		schema_id     = mso_schema_site_anp_epg_static_leaf.%[2]s.schema_id
		site_id       = mso_schema_site_anp_epg_static_leaf.%[2]s.site_id
		template_name = mso_schema_site_anp_epg_static_leaf.%[2]s.template_name
		anp_name      = mso_schema_site_anp_epg_static_leaf.%[2]s.anp_name
		epg_name      = mso_schema_site_anp_epg_static_leaf.%[2]s.epg_name
		path          = mso_schema_site_anp_epg_static_leaf.%[2]s.path
	}`,
		testAccMSOSchemaSiteAnpEpgStaticLeafConfigCreate(),
		msoSchemaTemplateAnpEpgName,
	)
}
