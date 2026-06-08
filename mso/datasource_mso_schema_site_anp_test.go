package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMSOSchemaSiteAnpDatasource(t *testing.T) {
	siteAnpResource := "mso_schema_site_anp." + msoSchemaTemplateAnpName
	siteAnpDatasource := "data.mso_schema_site_anp.anp"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Read site ANP datasource not-found error") },
				Config:      testAccMSOSchemaSiteAnpDatasourceNotFound(),
				ExpectError: regexp.MustCompile(`ANP .* is not found in Site`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read site ANP datasource") },
				Config:    testAccMSOSchemaSiteAnpDatasource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(siteAnpDatasource, "schema_id"),
					resource.TestCheckResourceAttrSet(siteAnpDatasource, "site_id"),
					resource.TestCheckResourceAttr(siteAnpDatasource, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr(siteAnpDatasource, "anp_name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttrPair(siteAnpDatasource, "schema_id", siteAnpResource, "schema_id"),
					resource.TestCheckResourceAttrPair(siteAnpDatasource, "site_id", siteAnpResource, "site_id"),
					resource.TestCheckResourceAttrPair(siteAnpDatasource, "template_name", siteAnpResource, "template_name"),
					resource.TestCheckResourceAttrPair(siteAnpDatasource, "anp_name", siteAnpResource, "anp_name"),
				),
			},
		},
	})
}

func testAccMSOSchemaSiteAnpDatasource() string {
	return fmt.Sprintf(`%[1]s
	data "mso_schema_site_anp" "anp" {
		schema_id     = mso_schema_site_anp.%[2]s.schema_id
		site_id       = mso_schema_site_anp.%[2]s.site_id
		template_name = mso_schema_site_anp.%[2]s.template_name
		anp_name      = mso_schema_site_anp.%[2]s.anp_name
	}`,
		testAccMSOSchemaSiteAnpConfigCreate(),
		msoSchemaTemplateAnpName,
	)
}

// testAccMSOSchemaSiteAnpDatasourceNotFound queries the datasource for an ANP
// name that does not exist on the site, exercising the getSiteAnp miss path
// ("ANP <name> is not found in Site.").
func testAccMSOSchemaSiteAnpDatasourceNotFound() string {
	return fmt.Sprintf(`%[1]s
	data "mso_schema_site_anp" "anp" {
		schema_id     = mso_schema_site_anp.%[2]s.schema_id
		site_id       = mso_schema_site_anp.%[2]s.site_id
		template_name = mso_schema_site_anp.%[2]s.template_name
		anp_name      = "non_existing_anp_name"
	}`,
		testAccMSOSchemaSiteAnpConfigCreate(),
		msoSchemaTemplateAnpName,
	)
}
