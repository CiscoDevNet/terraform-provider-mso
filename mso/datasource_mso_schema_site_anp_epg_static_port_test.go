package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccMSOSchemaSiteAnpEpgStaticPortDatasource verifies the read-only
// datasource for mso_schema_site_anp_epg_static_port:
//   - lookup with a path that does not exist returns an error
//   - successful read returns the expected attribute set and all values pair
//     with the managing resource
func TestAccMSOSchemaSiteAnpEpgStaticPortDatasource(t *testing.T) {
	staticPortResource := "mso_schema_site_anp_epg_static_port." + msoSchemaTemplateAnpEpgName
	staticPortDatasource := "data.mso_schema_site_anp_epg_static_port." + msoSchemaTemplateAnpEpgName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpEpgStaticPortDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: Lookup non-existent static port path (expect error)")
				},
				Config:      testAccMSOSchemaSiteAnpEpgStaticPortDatasourceNotFoundConfig(),
				ExpectError: regexp.MustCompile(`Unable to find static port entry`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read existing static port") },
				Config:    testAccMSOSchemaSiteAnpEpgStaticPortDatasourceReadConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(staticPortDatasource, "schema_id", staticPortResource, "schema_id"),
					resource.TestCheckResourceAttrPair(staticPortDatasource, "site_id", staticPortResource, "site_id"),
					resource.TestCheckResourceAttrPair(staticPortDatasource, "template_name", staticPortResource, "template_name"),
					resource.TestCheckResourceAttrPair(staticPortDatasource, "anp_name", staticPortResource, "anp_name"),
					resource.TestCheckResourceAttrPair(staticPortDatasource, "epg_name", staticPortResource, "epg_name"),
					resource.TestCheckResourceAttrPair(staticPortDatasource, "path_type", staticPortResource, "path_type"),
					resource.TestCheckResourceAttrPair(staticPortDatasource, "pod", staticPortResource, "pod"),
					resource.TestCheckResourceAttrPair(staticPortDatasource, "leaf", staticPortResource, "leaf"),
					resource.TestCheckResourceAttrPair(staticPortDatasource, "path", staticPortResource, "path"),
					resource.TestCheckResourceAttrPair(staticPortDatasource, "vlan", staticPortResource, "vlan"),
					resource.TestCheckResourceAttrPair(staticPortDatasource, "deployment_immediacy", staticPortResource, "deployment_immediacy"),
					resource.TestCheckResourceAttrPair(staticPortDatasource, "mode", staticPortResource, "mode"),
					resource.TestCheckResourceAttrPair(staticPortDatasource, "fex", staticPortResource, "fex"),
					resource.TestCheckResourceAttrPair(staticPortDatasource, "micro_seg_vlan", staticPortResource, "micro_seg_vlan"),
				),
			},
		},
	})
}

// testAccMSOSchemaSiteAnpEpgStaticPortDatasourceNotFoundConfig creates the
// static port resource then queries the datasource for a path that does not
// exist, exercising the "Unable to find static port entry" error path in
// datasourceMSOSchemaSiteAnpEpgStaticPortRead.
func testAccMSOSchemaSiteAnpEpgStaticPortDatasourceNotFoundConfig() string {
	return fmt.Sprintf(`%[1]s
	data "mso_schema_site_anp_epg_static_port" "missing" {
		schema_id     = mso_schema.%[2]s.id
		site_id       = mso_schema_site.%[3]s.site_id
		template_name = "%[4]s"
		anp_name      = mso_schema_template_anp.%[5]s.name
		epg_name      = mso_schema_site_anp_epg_static_port.%[6]s.epg_name
		path_type     = "port"
		pod           = "%[7]s"
		leaf          = "%[8]s"
		path          = "eth1/99"

		depends_on = [mso_schema_site_anp_epg_static_port.%[6]s]
	}`,
		testAccMSOSchemaSiteAnpEpgStaticPortConfigCreate(),
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
		msoSchemaTemplateAnpName,
		msoSchemaTemplateAnpEpgName,
		msoSchemaSiteAnpEpgStaticPortPod,
		msoSchemaSiteAnpEpgStaticPortLeaf,
	)
}

// testAccMSOSchemaSiteAnpEpgStaticPortDatasourceReadConfig creates the static
// port resource and reads it back via the datasource, verifying all computed
// attributes are populated correctly.
func testAccMSOSchemaSiteAnpEpgStaticPortDatasourceReadConfig() string {
	return fmt.Sprintf(`%[1]s
	data "mso_schema_site_anp_epg_static_port" "%[2]s" {
		schema_id     = mso_schema_site_anp_epg_static_port.%[2]s.schema_id
		site_id       = mso_schema_site_anp_epg_static_port.%[2]s.site_id
		template_name = mso_schema_site_anp_epg_static_port.%[2]s.template_name
		anp_name      = mso_schema_site_anp_epg_static_port.%[2]s.anp_name
		epg_name      = mso_schema_site_anp_epg_static_port.%[2]s.epg_name
		path_type     = mso_schema_site_anp_epg_static_port.%[2]s.path_type
		pod           = mso_schema_site_anp_epg_static_port.%[2]s.pod
		leaf          = mso_schema_site_anp_epg_static_port.%[2]s.leaf
		path          = mso_schema_site_anp_epg_static_port.%[2]s.path
		fex           = mso_schema_site_anp_epg_static_port.%[2]s.fex
	}`,
		testAccMSOSchemaSiteAnpEpgStaticPortConfigCreate(),
		msoSchemaTemplateAnpEpgName,
	)
}
