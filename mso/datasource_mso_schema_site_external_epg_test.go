package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccMSOSchemaSiteExternalEpgDataSource verifies the read-only datasource
// for mso_schema_site_external_epg:
//   - lookup by a name that does not exist returns an error
//   - successful read returns the expected attribute set, including the
//     attached L3Out reference (template-managed L3Out scenario)
func TestAccMSOSchemaSiteExternalEpgDataSource(t *testing.T) {
	siteEpgResource := "mso_schema_site_external_epg." + msoSchemaTemplateExtEpgName
	siteEpgDataSource := "data.mso_schema_site_external_epg." + msoSchemaTemplateExtEpgName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteExternalEpgDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("DataSource: Lookup non-existent site external EPG (expect error)")
				},
				Config:      testAccMSOSchemaSiteExternalEpgDataSourceNotFoundConfig(),
				ExpectError: regexp.MustCompile(`External EPG .* is not found in Site`),
			},
			{
				PreConfig: func() { fmt.Println("DataSource: Read existing site external EPG with template L3Out") },
				Config:    testAccMSOSchemaSiteExternalEpgDataSourceReadConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(siteEpgDataSource, "schema_id", siteEpgResource, "schema_id"),
					resource.TestCheckResourceAttrPair(siteEpgDataSource, "site_id", siteEpgResource, "site_id"),
					resource.TestCheckResourceAttrPair(siteEpgDataSource, "template_name", siteEpgResource, "template_name"),
					resource.TestCheckResourceAttrPair(siteEpgDataSource, "external_epg_name", siteEpgResource, "external_epg_name"),
					resource.TestCheckResourceAttrPair(siteEpgDataSource, "l3out_name", siteEpgResource, "l3out_name"),
					resource.TestCheckResourceAttrPair(siteEpgDataSource, "l3out_template_name", siteEpgResource, "l3out_template_name"),
					resource.TestCheckResourceAttrPair(siteEpgDataSource, "l3out_schema_id", siteEpgResource, "l3out_schema_id"),
					resource.TestCheckResourceAttrSet(siteEpgDataSource, "l3out_dn"),
				),
			},
		},
	})
}

// testAccMSOSchemaSiteExternalEpgDataSourceNotFoundConfig stands up the
// prereqs and the resource, then queries the datasource for an external EPG
// name that is not present under the site -- the read should fail.
func testAccMSOSchemaSiteExternalEpgDataSourceNotFoundConfig() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_external_epg" "%[2]s" {
		schema_id         = mso_schema.%[3]s.id
		site_id           = mso_schema_site.%[4]s.site_id
		template_name     = "%[5]s"
		external_epg_name = mso_schema_template_external_epg.%[2]s.external_epg_name
	}

	data "mso_schema_site_external_epg" "missing" {
		schema_id         = mso_schema.%[3]s.id
		site_id           = mso_schema_site.%[4]s.site_id
		template_name     = "%[5]s"
		external_epg_name = "does_not_exist"

		depends_on = [mso_schema_site_external_epg.%[2]s]
	}`,
		testAccMSOSchemaSiteExternalEpgPrerequisiteConfig(),
		msoSchemaTemplateExtEpgName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
	)
}

// testAccMSOSchemaSiteExternalEpgDataSourceReadConfig creates the site
// external EPG with a template-managed L3Out and reads it back via the
// datasource.
func testAccMSOSchemaSiteExternalEpgDataSourceReadConfig() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_external_epg" "%[2]s" {
		schema_id           = mso_schema.%[3]s.id
		site_id             = mso_schema_site.%[4]s.site_id
		template_name       = "%[5]s"
		external_epg_name   = mso_schema_template_external_epg.%[2]s.external_epg_name
		l3out_name          = mso_schema_template_l3out.%[6]s.l3out_name
		l3out_template_name = "%[5]s"
		l3out_schema_id     = mso_schema.%[3]s.id
	}

	data "mso_schema_site_external_epg" "%[2]s" {
		schema_id         = mso_schema.%[3]s.id
		site_id           = mso_schema_site.%[4]s.site_id
		template_name     = "%[5]s"
		external_epg_name = mso_schema_site_external_epg.%[2]s.external_epg_name
	}`,
		testAccMSOSchemaSiteExternalEpgPrerequisiteConfig(),
		msoSchemaTemplateExtEpgName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
		msoSchemaTemplateL3outName,
	)
}
