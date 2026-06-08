package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccMSOSchemaSiteAnpResource exercises the lifecycle of
// mso_schema_site_anp:
//   - attempt to create the site ANP on a template that has no
//     mso_schema_site association (the Create PATCH targets
//     /sites/<siteId>-<template>/anps/<anp> on a non-existent site
//     association, so MSO returns "Resource Not Found")
//   - create the site ANP with a proper schema_site association
//   - import the site ANP
//
// The lab must have the `ansible_test` and `ansible_test_2` sites onboarded.
func TestAccMSOSchemaSiteAnpResource(t *testing.T) {
	siteAnpResource := "mso_schema_site_anp." + msoSchemaTemplateAnpName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: Create site ANP without mso_schema_site association (expect error)")
				},
				Config: testAccMSOSchemaSiteAnpConfigNoSiteAssociation(),
				// Older MSO/NDO rejects both the "replace" and the fallback
				// "add" PATCH with "Resource Not Found", so Create returns
				// the raw error. Newer NDO appear to silently accept the fallback
				// "add" PATCH as a no-op, so Create succeeds, the follow-up
				// Read finds nothing, and the SDK raises "Provider produced
				// inconsistent result after apply". Match either outcome.
				//
				// Note: this divergence is driven by NDO's schema validation
				// engine, which is now always-on (the validation flag is no
				// longer a configurable option). Older releases let the
				// PATCH bypass validation and inject the entry into the
				// schema document; newer releases run the validator on
				// every PATCH and silently drop changes that reference a
				// site/template association that does not exist.
				ExpectError: regexp.MustCompile(`Resource Not Found|Provider produced inconsistent result after apply`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create site ANP") },
				Config:    testAccMSOSchemaSiteAnpConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(siteAnpResource, "schema_id"),
					resource.TestCheckResourceAttrSet(siteAnpResource, "site_id"),
					resource.TestCheckResourceAttr(siteAnpResource, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr(siteAnpResource, "anp_name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttrPair(
						siteAnpResource, "site_id",
						"data.mso_site."+msoTemplateSiteName1, "id",
					),
					resource.TestCheckResourceAttr(siteAnpResource, "id", msoSchemaTemplateAnpName),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import site ANP") },
				ResourceName: siteAnpResource,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[siteAnpResource]
					if !ok {
						return "", fmt.Errorf("site ANP resource not found in state: %s", siteAnpResource)
					}
					return fmt.Sprintf("%s/site/%s/template/%s/anp/%s",
						rs.Primary.Attributes["schema_id"],
						rs.Primary.Attributes["site_id"],
						rs.Primary.Attributes["template_name"],
						rs.Primary.Attributes["anp_name"],
					), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

// testAccMSOSchemaSiteAnpPrerequisiteConfig emits the full scaffolding for a
// site ANP: both data.mso_site blocks, the shared tenant attached to both
// sites, the schema with one template, the mso_schema_site association for
// site_1 (ansible_test, undeploy_on_destroy=false), and the template ANP.
func testAccMSOSchemaSiteAnpPrerequisiteConfig() string {
	return fmt.Sprintf(`%s%s`,
		testSchemaWithSingleSiteAssociationConfig(),
		testSchemaTemplateAnpConfig(),
	)
}

// testAccMSOSchemaSiteAnpNoSiteAssociationConfig is the prereq without the
// mso_schema_site block. It is used to drive the negative-path step: the
// site ANP API write succeeds but the subsequent Read cannot resolve the
// (siteId, templateName, anpName) tuple because no site is attached to the
// template, producing "Provider produced inconsistent result after apply".
func testAccMSOSchemaSiteAnpNoSiteAssociationConfig() string {
	return fmt.Sprintf(`%s%s`,
		testSchemaWithBothSitesPrerequisiteConfig(),
		testSchemaTemplateAnpConfig(),
	)
}

func testAccMSOSchemaSiteAnpConfigCreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_anp" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = mso_schema_site.%[4]s.site_id
		template_name = "%[5]s"
		anp_name      = mso_schema_template_anp.%[2]s.name
	}`,
		testAccMSOSchemaSiteAnpPrerequisiteConfig(),
		msoSchemaTemplateAnpName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
	)
}

func testAccMSOSchemaSiteAnpConfigNoSiteAssociation() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_anp" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = data.mso_site.%[4]s.id
		template_name = "%[5]s"
		anp_name      = mso_schema_template_anp.%[2]s.name
	}`,
		testAccMSOSchemaSiteAnpNoSiteAssociationConfig(),
		msoSchemaTemplateAnpName,
		msoSchemaName,
		msoTemplateSiteName1,
		msoSchemaTemplateName,
	)
}

// testAccCheckMSOSchemaSiteAnpDestroy walks state for any mso_schema_site_anp
// resources, fetches the schema, and asserts that no sites[].anps[] entry
// whose anpRef resolves to anp_name remains under the matching siteId +
// templateName. A missing schema or missing sites array is treated as a
// successful destroy.
func testAccCheckMSOSchemaSiteAnpDestroy(s *terraform.State) error {
	msoClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mso_schema_site_anp" {
			continue
		}
		schemaId := rs.Primary.Attributes["schema_id"]
		stateSiteId := rs.Primary.Attributes["site_id"]
		stateTemplate := rs.Primary.Attributes["template_name"]
		stateAnp := rs.Primary.Attributes["anp_name"]

		cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
		if err != nil {
			// Schema itself has been destroyed.
			return nil
		}
		count, err := cont.ArrayCount("sites")
		if err != nil {
			return nil
		}
		re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/anps/(.*)")
		for i := 0; i < count; i++ {
			siteCont, err := cont.ArrayElement(i, "sites")
			if err != nil {
				return err
			}
			apiSiteId := models.StripQuotes(siteCont.S("siteId").String())
			apiTemplate := models.StripQuotes(siteCont.S("templateName").String())
			if apiSiteId != stateSiteId || apiTemplate != stateTemplate {
				continue
			}
			anpCount, err := siteCont.ArrayCount("anps")
			if err != nil {
				continue
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := siteCont.ArrayElement(j, "anps")
				if err != nil {
					return err
				}
				anpRef := models.StripQuotes(anpCont.S("anpRef").String())
				match := re.FindStringSubmatch(anpRef)
				if len(match) == 4 && match[3] == stateAnp {
					return fmt.Errorf("mso_schema_site_anp (site=%s, template=%s, anp=%s) still exists on schema %s", stateSiteId, stateTemplate, stateAnp, schemaId)
				}
			}
		}
	}
	return nil
}
