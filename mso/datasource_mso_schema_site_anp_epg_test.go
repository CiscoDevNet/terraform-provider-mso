package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOSchemaSiteAnpEpgDatasource(t *testing.T) {
	siteAnpEpgResource := "mso_schema_site_anp_epg." + msoSchemaTemplateAnpEpgName
	siteAnpEpgDatasource := "data.mso_schema_site_anp_epg.epg"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpEpgDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Read site ANP EPG datasource not-found error") },
				Config:      testAccMSOSchemaSiteAnpEpgDatasourceNotFound(),
				ExpectError: regexp.MustCompile(`EPG .* is not found in Site`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read site ANP EPG datasource") },
				Config:    testAccMSOSchemaSiteAnpEpgDatasource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(siteAnpEpgDatasource, "schema_id"),
					resource.TestCheckResourceAttrSet(siteAnpEpgDatasource, "site_id"),
					resource.TestCheckResourceAttr(siteAnpEpgDatasource, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr(siteAnpEpgDatasource, "anp_name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr(siteAnpEpgDatasource, "epg_name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttrPair(siteAnpEpgDatasource, "schema_id", siteAnpEpgResource, "schema_id"),
					resource.TestCheckResourceAttrPair(siteAnpEpgDatasource, "site_id", siteAnpEpgResource, "site_id"),
					resource.TestCheckResourceAttrPair(siteAnpEpgDatasource, "template_name", siteAnpEpgResource, "template_name"),
					resource.TestCheckResourceAttrPair(siteAnpEpgDatasource, "anp_name", siteAnpEpgResource, "anp_name"),
					resource.TestCheckResourceAttrPair(siteAnpEpgDatasource, "epg_name", siteAnpEpgResource, "epg_name"),
				),
			},
		},
	})
}

func testAccMSOSchemaSiteAnpEpgDatasource() string {
	return fmt.Sprintf(`%[1]s
	data "mso_schema_site_anp_epg" "epg" {
		schema_id     = mso_schema_site_anp_epg.%[2]s.schema_id
		site_id       = mso_schema_site_anp_epg.%[2]s.site_id
		template_name = mso_schema_site_anp_epg.%[2]s.template_name
		anp_name      = mso_schema_site_anp_epg.%[2]s.anp_name
		epg_name      = mso_schema_site_anp_epg.%[2]s.epg_name
	}`,
		testAccMSOSchemaSiteAnpEpgConfigCreate(),
		msoSchemaTemplateAnpEpgName,
	)
}

// testAccMSOSchemaSiteAnpEpgDatasourceNotFound queries the datasource for an
// EPG name that does not exist on the site, exercising the getSiteEpg miss
// path ("EPG <name> is not found in Site.").
func testAccMSOSchemaSiteAnpEpgDatasourceNotFound() string {
	return fmt.Sprintf(`%[1]s
	data "mso_schema_site_anp_epg" "epg" {
		schema_id     = mso_schema_site_anp_epg.%[2]s.schema_id
		site_id       = mso_schema_site_anp_epg.%[2]s.site_id
		template_name = mso_schema_site_anp_epg.%[2]s.template_name
		anp_name      = mso_schema_site_anp_epg.%[2]s.anp_name
		epg_name      = "non_existing_epg_name"
	}`,
		testAccMSOSchemaSiteAnpEpgConfigCreate(),
		msoSchemaTemplateAnpEpgName,
	)
}
