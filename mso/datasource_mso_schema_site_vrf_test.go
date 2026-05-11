package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOSchemaSiteVrfDatasource(t *testing.T) {
	siteVrfResource := "mso_schema_site_vrf." + msoSchemaTemplateVrfName
	siteVrfDatasource := "data.mso_schema_site_vrf.vrf"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteVrfDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Read site VRF datasource not-found error") },
				Config:      testAccMSOSchemaSiteVrfDatasourceNotFound(),
				ExpectError: regexp.MustCompile(`VRF .* is not found in Site`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read site VRF datasource") },
				Config:    testAccMSOSchemaSiteVrfDatasource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(siteVrfDatasource, "schema_id"),
					resource.TestCheckResourceAttrSet(siteVrfDatasource, "site_id"),
					resource.TestCheckResourceAttr(siteVrfDatasource, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr(siteVrfDatasource, "vrf_name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttrPair(siteVrfDatasource, "schema_id", siteVrfResource, "schema_id"),
					resource.TestCheckResourceAttrPair(siteVrfDatasource, "site_id", siteVrfResource, "site_id"),
					resource.TestCheckResourceAttrPair(siteVrfDatasource, "template_name", siteVrfResource, "template_name"),
					resource.TestCheckResourceAttrPair(siteVrfDatasource, "vrf_name", siteVrfResource, "vrf_name"),
				),
			},
		},
	})
}

// testAccMSOSchemaSiteVrfDatasourceConfig emits the prereq plus the site VRF
// resource so the datasource has something to read.
func testAccMSOSchemaSiteVrfDatasourceConfig() string {
	return testAccMSOSchemaSiteVrfConfigCreate()
}

func testAccMSOSchemaSiteVrfDatasource() string {
	return fmt.Sprintf(`%[1]s
	data "mso_schema_site_vrf" "vrf" {
		schema_id     = mso_schema_site_vrf.%[2]s.schema_id
		site_id       = mso_schema_site_vrf.%[2]s.site_id
		template_name = mso_schema_site_vrf.%[2]s.template_name
		vrf_name      = mso_schema_site_vrf.%[2]s.vrf_name
	}`,
		testAccMSOSchemaSiteVrfDatasourceConfig(),
		msoSchemaTemplateVrfName,
	)
}

// testAccMSOSchemaSiteVrfDatasourceNotFound queries the datasource for a VRF
// name that does not exist on the site, exercising the getSiteVrf miss path
// ("VRF <name> is not found in Site.").
func testAccMSOSchemaSiteVrfDatasourceNotFound() string {
	return fmt.Sprintf(`%[1]s
	data "mso_schema_site_vrf" "vrf" {
		schema_id     = mso_schema_site_vrf.%[2]s.schema_id
		site_id       = mso_schema_site_vrf.%[2]s.site_id
		template_name = mso_schema_site_vrf.%[2]s.template_name
		vrf_name      = "non_existing_vrf_name"
	}`,
		testAccMSOSchemaSiteVrfDatasourceConfig(),
		msoSchemaTemplateVrfName,
	)
}
