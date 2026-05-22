package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccMSOSchemaSiteAnpEpgBulkStaticPortDatasource verifies the read-only
// datasource for mso_schema_site_anp_epg_bulk_staticport:
//   - lookup with a non-existent epg_name returns an error
//   - successful read returns two static ports and all attribute values
//     pair with the managing resource
func TestAccMSOSchemaSiteAnpEpgBulkStaticPortDatasource(t *testing.T) {
	bulkStaticPortResource := "mso_schema_site_anp_epg_bulk_staticport." + msoSchemaTemplateAnpEpgName
	bulkStaticPortDatasource := "data.mso_schema_site_anp_epg_bulk_staticport." + msoSchemaTemplateAnpEpgName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpEpgBulkStaticPortDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("DataSource: Lookup non-existent EPG (expect error)")
				},
				Config:      testAccMSOSchemaSiteAnpEpgBulkStaticPortDatasourceNotFoundConfig(),
				ExpectError: regexp.MustCompile(`EPG .* is not found in Site\.`),
			},
			{
				PreConfig: func() { fmt.Println("DataSource: Read existing bulk static port") },
				Config:    testAccMSOSchemaSiteAnpEpgBulkStaticPortDatasourceReadConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(bulkStaticPortDatasource, "schema_id", bulkStaticPortResource, "schema_id"),
					resource.TestCheckResourceAttrPair(bulkStaticPortDatasource, "site_id", bulkStaticPortResource, "site_id"),
					resource.TestCheckResourceAttrPair(bulkStaticPortDatasource, "template_name", bulkStaticPortResource, "template_name"),
					resource.TestCheckResourceAttrPair(bulkStaticPortDatasource, "anp_name", bulkStaticPortResource, "anp_name"),
					resource.TestCheckResourceAttrPair(bulkStaticPortDatasource, "epg_name", bulkStaticPortResource, "epg_name"),
					resource.TestCheckResourceAttr(bulkStaticPortDatasource, "static_ports.#", "2"),
					CustomTestCheckTypeSetElemAttrs(bulkStaticPortDatasource, "static_ports", map[string]string{
						"path_type":            "port",
						"pod":                  msoSchemaSiteAnpEpgStaticPortPod,
						"leaf":                 msoSchemaSiteAnpEpgStaticPortLeaf,
						"path":                 msoSchemaSiteAnpEpgStaticPortPath,
						"vlan":                 "200",
						"deployment_immediacy": "lazy",
						"mode":                 "regular",
						"micro_seg_vlan":       "300",
						"fex":                  "",
					}),
					CustomTestCheckTypeSetElemAttrs(bulkStaticPortDatasource, "static_ports", map[string]string{
						"path_type":            "port",
						"pod":                  msoSchemaSiteAnpEpgStaticPortPod,
						"leaf":                 msoSchemaSiteAnpEpgStaticPortLeaf,
						"path":                 msoSchemaSiteAnpEpgStaticPortPath2,
						"vlan":                 "201",
						"deployment_immediacy": "immediate",
						"mode":                 "untagged",
						"fex":                  msoSchemaSiteAnpEpgStaticPortFex,
					}),
				),
			},
		},
	})
}

// testAccMSOSchemaSiteAnpEpgBulkStaticPortDatasourceNotFoundConfig creates the
// bulk static port resource then queries the datasource for an epg_name that
// does not exist, exercising the "EPG X is not found in Site." error path in
// datasourceMSOSchemaSiteAnpEpgBulkStaticPortRead.
func testAccMSOSchemaSiteAnpEpgBulkStaticPortDatasourceNotFoundConfig() string {
	return fmt.Sprintf(`%[1]s
	data "mso_schema_site_anp_epg_bulk_staticport" "missing" {
		schema_id     = mso_schema.%[2]s.id
		site_id       = mso_schema_site.%[3]s.site_id
		template_name = "%[4]s"
		anp_name      = mso_schema_template_anp.%[5]s.name
		epg_name      = "nonexistent-epg"

		depends_on = [mso_schema_site_anp_epg_bulk_staticport.%[6]s]
	}`,
		testAccMSOSchemaSiteAnpEpgBulkStaticPortConfigCreate(),
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
		msoSchemaTemplateAnpName,
		msoSchemaTemplateAnpEpgName,
	)
}

// testAccMSOSchemaSiteAnpEpgBulkStaticPortDatasourceReadConfig creates the
// bulk static port resource (two ports) and reads it back via the datasource,
// verifying all computed attributes are populated correctly.
func testAccMSOSchemaSiteAnpEpgBulkStaticPortDatasourceReadConfig() string {
	return fmt.Sprintf(`%[1]s
	data "mso_schema_site_anp_epg_bulk_staticport" "%[2]s" {
		schema_id     = mso_schema_site_anp_epg_bulk_staticport.%[2]s.schema_id
		site_id       = mso_schema_site_anp_epg_bulk_staticport.%[2]s.site_id
		template_name = mso_schema_site_anp_epg_bulk_staticport.%[2]s.template_name
		anp_name      = mso_schema_site_anp_epg_bulk_staticport.%[2]s.anp_name
		epg_name      = mso_schema_site_anp_epg_bulk_staticport.%[2]s.epg_name
	}`,
		testAccMSOSchemaSiteAnpEpgBulkStaticPortConfigCreate(),
		msoSchemaTemplateAnpEpgName,
	)
}
