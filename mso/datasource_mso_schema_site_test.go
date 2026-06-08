package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMSOSchemaSiteDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteDestroy,
		Steps: []resource.TestStep{
			// Step ordering: error -> error -> success, with each step's
			// PreConfig calling cleanupOrphanSchemaSiteTestResources to remove
			// any orphaned tenant/schema left on MSO from a previous step's
			// rolled-back apply.
			//
			// Why each step needs cleanup: SDK v1 rolls back state on any
			// apply error (including data source read errors). The two
			// ExpectError steps therefore leave their prereq tenant + schema
			// on the MSO server but absent from Terraform state. Without
			// cleanup, the next step would try to recreate the same names
			// and fail with "Tenant/Schema: '...' already exists".
			{
				PreConfig: func() {
					fmt.Println("Test: Read schema_site datasource site-name-not-onboarded error")
					cleanupOrphanSchemaSiteTestResources(t)
				},
				Config:      testAccMSOSchemaSiteDatasourceSiteNameNotOnboarded(),
				ExpectError: regexp.MustCompile("Site of specified name not found"),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Read schema_site datasource site-not-attached-to-template error")
					cleanupOrphanSchemaSiteTestResources(t)
				},
				Config:      testAccMSOSchemaSiteDatasourceSiteNotAttached(),
				ExpectError: regexp.MustCompile("Site-Template association for .* is not found"),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Read schema_site datasource")
					cleanupOrphanSchemaSiteTestResources(t)
				},
				Config: testAccMSOSchemaSiteDatasource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.mso_schema_site.site", "schema_id"),
					resource.TestCheckResourceAttrSet("data.mso_schema_site.site", "site_id"),
					resource.TestCheckResourceAttr("data.mso_schema_site.site", "name", msoTemplateSiteName1),
					resource.TestCheckResourceAttrPair(
						"data.mso_schema_site.site", "site_id",
						"data.mso_site."+msoTemplateSiteName1, "id",
					),
					resource.TestCheckResourceAttrPair(
						"data.mso_schema_site.site", "template_name",
						"mso_schema_site."+msoSchemaSiteResourceLabel1, "template_name",
					),
				),
			},
		},
	})
}

// testAccMSOSchemaSiteDatasource exercises the success flow: the schema has
// site_1 (ansible_test) attached to its template, and the datasource looks up
// that site by name.
//
// Note: schema_id / template_name / name are sourced from the
// mso_schema_site.site_1 resource attributes (rather than from
// mso_schema/data.mso_site directly with a `depends_on`). In SDK v1 a
// `depends_on` on a data source forces a deferred read and produces a
// non-empty plan after apply ("plan was not empty" failure). Implicit
// references on resource attributes give us the same ordering guarantee
// without the deferral.
func testAccMSOSchemaSiteDatasource() string {
	return fmt.Sprintf(`%s
	data "mso_schema_site" "site" {
		schema_id     = mso_schema_site.%[2]s.schema_id
		template_name = mso_schema_site.%[2]s.template_name
		name          = data.mso_site.%[3]s.name
	}`, testSchemaWithSingleSiteAssociationConfig(), msoSchemaSiteResourceLabel1, msoTemplateSiteName1)
}

// testAccMSOSchemaSiteDatasourceSiteNameNotOnboarded targets the first error
// path in datasourceMSOSchemaSiteRead: GET /api/v1/sites is scanned for a
// matching display name, and a miss returns "Site of specified name not
// found" before the schema/template association is even checked.
//
// Uses the single-site association prereq (rather than the bare
// testSchemaWithBothSitesPrerequisiteConfig) so this step does not
// destroy/recreate mso_schema_site relative to the surrounding steps in
// TestAccMSOSchemaSiteDatasource -- keeping the prereq stable across all
// steps avoids state churn and lets the ExpectError step plan as a no-op
// for the existing resources.
func testAccMSOSchemaSiteDatasourceSiteNameNotOnboarded() string {
	return fmt.Sprintf(`%s
	data "mso_schema_site" "site" {
		schema_id     = mso_schema_site.%[2]s.schema_id
		template_name = mso_schema_site.%[2]s.template_name
		name          = "non_existing_site_name"
	}`, testSchemaWithSingleSiteAssociationConfig(), msoSchemaSiteResourceLabel1)
}

// testAccMSOSchemaSiteDatasourceSiteNotAttached targets the second error
// path: the site exists in the MSO inventory (ansible_test_2 is onboarded)
// but is not associated with the schema/template, so
// getSiteFromSiteIdAndTemplate returns "Site-Template association ... is not
// found."
func testAccMSOSchemaSiteDatasourceSiteNotAttached() string {
	return fmt.Sprintf(`%s
	data "mso_schema_site" "site" {
		schema_id     = mso_schema_site.%[2]s.schema_id
		template_name = mso_schema_site.%[2]s.template_name
		name          = data.mso_site.%[3]s.name
	}`, testSchemaWithSingleSiteAssociationConfig(), msoSchemaSiteResourceLabel1, msoTemplateSiteName2)
}
