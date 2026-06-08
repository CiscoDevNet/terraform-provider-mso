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

// TestAccMSOSchemaSiteVrfResource exercises the lifecycle of
// mso_schema_site_vrf:
//   - attempt to create the site VRF on a template that has no
//     mso_schema_site association (the Create PATCH targets
//     /sites/<siteId>-<template>/vrfs/<vrf> on a non-existent site
//     association, so MSO returns "Resource Not Found")
//   - create the site VRF with a proper schema_site association
//   - import the site VRF
//
// The lab must have the `ansible_test` and `ansible_test_2` sites onboarded.
func TestAccMSOSchemaSiteVrfResource(t *testing.T) {
	siteVrfResource := "mso_schema_site_vrf." + msoSchemaTemplateVrfName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteVrfDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: Create site VRF without mso_schema_site association (expect error)")
				},
				Config: testAccMSOSchemaSiteVrfConfigNoSiteAssociation(),
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
				PreConfig: func() { fmt.Println("Test: Create site VRF") },
				Config:    testAccMSOSchemaSiteVrfConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(siteVrfResource, "schema_id"),
					resource.TestCheckResourceAttrSet(siteVrfResource, "site_id"),
					resource.TestCheckResourceAttr(siteVrfResource, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr(siteVrfResource, "vrf_name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttrPair(
						siteVrfResource, "site_id",
						"data.mso_site."+msoTemplateSiteName1, "id",
					),
					resource.TestCheckResourceAttr(siteVrfResource, "id", msoSchemaTemplateVrfName),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import site VRF") },
				ResourceName: siteVrfResource,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[siteVrfResource]
					if !ok {
						return "", fmt.Errorf("site VRF resource not found in state: %s", siteVrfResource)
					}
					return fmt.Sprintf("%s/site/%s/vrf/%s",
						rs.Primary.Attributes["schema_id"],
						rs.Primary.Attributes["site_id"],
						rs.Primary.Attributes["vrf_name"],
					), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

// testAccMSOSchemaSiteVrfPrerequisiteConfig emits the full scaffolding for a
// site VRF: both data.mso_site blocks, the shared tenant attached to both
// sites, the schema with one template, the mso_schema_site association for
// site_1 (ansible_test, undeploy_on_destroy=false), and the template VRF.
func testAccMSOSchemaSiteVrfPrerequisiteConfig() string {
	return fmt.Sprintf(`%s%s`,
		testSchemaWithSingleSiteAssociationConfig(),
		testSchemaTemplateVrfConfig(),
	)
}

// testAccMSOSchemaSiteVrfNoSiteAssociationConfig is the prereq without the
// mso_schema_site block. It is used to drive the negative-path step: the
// site VRF API write succeeds but the subsequent Read cannot resolve the
// (siteId, templateName, vrfName) tuple because no site is attached to the
// template, producing "Provider produced inconsistent result after apply".
func testAccMSOSchemaSiteVrfNoSiteAssociationConfig() string {
	return fmt.Sprintf(`%s%s`,
		testSchemaWithBothSitesPrerequisiteConfig(),
		testSchemaTemplateVrfConfig(),
	)
}

func testAccMSOSchemaSiteVrfConfigCreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_vrf" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = mso_schema_site.%[4]s.site_id
		template_name = "%[5]s"
		vrf_name      = mso_schema_template_vrf.%[2]s.name
	}`,
		testAccMSOSchemaSiteVrfPrerequisiteConfig(),
		msoSchemaTemplateVrfName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
	)
}

func testAccMSOSchemaSiteVrfConfigNoSiteAssociation() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_vrf" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = data.mso_site.%[4]s.id
		template_name = "%[5]s"
		vrf_name      = mso_schema_template_vrf.%[2]s.name
	}`,
		testAccMSOSchemaSiteVrfNoSiteAssociationConfig(),
		msoSchemaTemplateVrfName,
		msoSchemaName,
		msoTemplateSiteName1,
		msoSchemaTemplateName,
	)
}

// testAccCheckMSOSchemaSiteVrfDestroy walks state for any mso_schema_site_vrf
// resources, fetches the schema, and asserts that no sites[].vrfs[] entry
// whose vrfRef resolves to vrf_name remains under the matching siteId +
// templateName. A missing schema or missing sites array is treated as a
// successful destroy.
func testAccCheckMSOSchemaSiteVrfDestroy(s *terraform.State) error {
	msoClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mso_schema_site_vrf" {
			continue
		}
		schemaId := rs.Primary.Attributes["schema_id"]
		stateSiteId := rs.Primary.Attributes["site_id"]
		stateTemplate := rs.Primary.Attributes["template_name"]
		stateVrf := rs.Primary.Attributes["vrf_name"]

		cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
		if err != nil {
			// Schema itself has been destroyed.
			return nil
		}
		count, err := cont.ArrayCount("sites")
		if err != nil {
			return nil
		}
		re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/vrfs/(.*)")
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
			vrfCount, err := siteCont.ArrayCount("vrfs")
			if err != nil {
				continue
			}
			for j := 0; j < vrfCount; j++ {
				vrfCont, err := siteCont.ArrayElement(j, "vrfs")
				if err != nil {
					return err
				}
				vrfRef := models.StripQuotes(vrfCont.S("vrfRef").String())
				match := re.FindStringSubmatch(vrfRef)
				if len(match) == 4 && match[3] == stateVrf {
					return fmt.Errorf("mso_schema_site_vrf (site=%s, template=%s, vrf=%s) still exists on schema %s", stateSiteId, stateTemplate, stateVrf, schemaId)
				}
			}
		}
	}
	return nil
}
