package mso

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccMSOSchemaSiteAnpEpgResource exercises the lifecycle of
// mso_schema_site_anp_epg:
//   - attempt to create the site ANP EPG on a template that has no
//     mso_schema_site association (expect error)
//   - create the site ANP EPG with a proper schema_site association
//   - import the site ANP EPG
//
// private_link_label is intentionally skipped: it is a cloud-only service
// EPG attribute that requires a cloud site setup and is a deprecation
// candidate (see resource_mso_schema_site_anp_epg.go).
//
// The lab must have the `ansible_test` and `ansible_test_2` sites onboarded.
func TestAccMSOSchemaSiteAnpEpgResource(t *testing.T) {
	siteAnpEpgResource := "mso_schema_site_anp_epg." + msoSchemaTemplateAnpEpgName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpEpgDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: Create site ANP EPG without mso_schema_site association (expect error)")
				},
				Config:      testAccMSOSchemaSiteAnpEpgConfigNoSiteAssociation(),
				ExpectError: regexp.MustCompile(`Resource Not Found|Provider produced inconsistent result after apply`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create site ANP EPG") },
				Config:    testAccMSOSchemaSiteAnpEpgConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(siteAnpEpgResource, "schema_id"),
					resource.TestCheckResourceAttrSet(siteAnpEpgResource, "site_id"),
					resource.TestCheckResourceAttr(siteAnpEpgResource, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr(siteAnpEpgResource, "anp_name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr(siteAnpEpgResource, "epg_name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttrPair(
						siteAnpEpgResource, "site_id",
						"data.mso_site."+msoTemplateSiteName1, "id",
					),
					resource.TestCheckResourceAttr(siteAnpEpgResource, "id", msoSchemaTemplateAnpEpgName),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import site ANP EPG") },
				ResourceName: siteAnpEpgResource,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[siteAnpEpgResource]
					if !ok {
						return "", fmt.Errorf("site ANP EPG resource not found in state: %s", siteAnpEpgResource)
					}
					return fmt.Sprintf("%s/site/%s/template/%s/anp/%s/epg/%s",
						rs.Primary.Attributes["schema_id"],
						rs.Primary.Attributes["site_id"],
						rs.Primary.Attributes["template_name"],
						rs.Primary.Attributes["anp_name"],
						rs.Primary.Attributes["epg_name"],
					), nil
				},
				ImportStateVerify: true,
				// private_link_label is a cloud-only attribute that is not
				// exercised by these tests. Import sets it to "" because the
				// schema document carries an empty privateLinkLabel object,
				// while the post-create state leaves it unset. Skip the diff.
				ImportStateVerifyIgnore: []string{"private_link_label"},
			},
		},
	})
}

// testAccMSOSchemaSiteAnpEpgTemplateAnpEpgWithBdConfig emits a template ANP
// EPG with bd_name set. NDO requires every template EPG to be associated to
// a BD when the template is associated with an on-prem site, otherwise the
// PATCH that materializes the site EPG is rejected with:
//
//	"EPG: <name> ... must be associated to a BD while template is associated
//	 to an on-prem site."
//
// The shared testSchemaTemplateAnpEpgConfig does not set bd_name, so the
// site_anp_epg tests use this local variant instead.
func testAccMSOSchemaSiteAnpEpgTemplateAnpEpgWithBdConfig() string {
	return fmt.Sprintf(`
resource "mso_schema_template_anp_epg" "%[1]s" {
	name          = "%[1]s"
	display_name  = "%[1]s"
	anp_name      = mso_schema_template_anp.%[2]s.name
	schema_id     = mso_schema.%[3]s.id
	template_name = "%[4]s"
	bd_name       = mso_schema_template_bd.%[5]s.name
}
`, msoSchemaTemplateAnpEpgName, msoSchemaTemplateAnpName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateBdName)
}

// testAccMSOSchemaSiteAnpEpgPrerequisiteConfig emits the full scaffolding:
// both data.mso_site blocks, the shared tenant attached to both sites, the
// schema with one template, the mso_schema_site association for site_1, the
// template VRF + BD (so the EPG can satisfy the on-prem BD-association
// requirement), the template ANP, and the template ANP EPG bound to the BD.
func testAccMSOSchemaSiteAnpEpgPrerequisiteConfig() string {
	return fmt.Sprintf(`%s%s%s%s%s`,
		testSchemaWithSingleSiteAssociationConfig(),
		testSchemaTemplateVrfConfig(),
		testSchemaTemplateBdConfig(),
		testSchemaTemplateAnpConfig(),
		testAccMSOSchemaSiteAnpEpgTemplateAnpEpgWithBdConfig(),
	)
}

// testAccMSOSchemaSiteAnpEpgNoSiteAssociationConfig is the prereq without
// the mso_schema_site block. Used to drive the negative-path step.
func testAccMSOSchemaSiteAnpEpgNoSiteAssociationConfig() string {
	return fmt.Sprintf(`%s%s%s%s%s`,
		testSchemaWithBothSitesPrerequisiteConfig(),
		testSchemaTemplateVrfConfig(),
		testSchemaTemplateBdConfig(),
		testSchemaTemplateAnpConfig(),
		testAccMSOSchemaSiteAnpEpgTemplateAnpEpgWithBdConfig(),
	)
}

func testAccMSOSchemaSiteAnpEpgConfigCreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_anp_epg" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = mso_schema_site.%[4]s.site_id
		template_name = "%[5]s"
		anp_name      = mso_schema_template_anp.%[6]s.name
		epg_name      = mso_schema_template_anp_epg.%[2]s.name
	}`,
		testAccMSOSchemaSiteAnpEpgPrerequisiteConfig(),
		msoSchemaTemplateAnpEpgName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
		msoSchemaTemplateAnpName,
	)
}

func testAccMSOSchemaSiteAnpEpgConfigNoSiteAssociation() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_anp_epg" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = data.mso_site.%[4]s.id
		template_name = "%[5]s"
		anp_name      = mso_schema_template_anp.%[6]s.name
		epg_name      = mso_schema_template_anp_epg.%[2]s.name
	}`,
		testAccMSOSchemaSiteAnpEpgNoSiteAssociationConfig(),
		msoSchemaTemplateAnpEpgName,
		msoSchemaName,
		msoTemplateSiteName1,
		msoSchemaTemplateName,
		msoSchemaTemplateAnpName,
	)
}

// testAccCheckMSOSchemaSiteAnpEpgDestroy walks state for any
// mso_schema_site_anp_epg resources, fetches the schema, and asserts that no
// sites[].anps[].epgs[] entry whose epgRef resolves to epg_name remains
// under the matching siteId + anpName. A missing schema or missing sites
// array is treated as a successful destroy.
func testAccCheckMSOSchemaSiteAnpEpgDestroy(s *terraform.State) error {
	msoClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mso_schema_site_anp_epg" {
			continue
		}
		schemaId := rs.Primary.Attributes["schema_id"]
		stateSiteId := rs.Primary.Attributes["site_id"]
		stateTemplate := rs.Primary.Attributes["template_name"]
		stateAnp := rs.Primary.Attributes["anp_name"]
		stateEpg := rs.Primary.Attributes["epg_name"]

		cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
		if err != nil {
			return nil
		}
		count, err := cont.ArrayCount("sites")
		if err != nil {
			return nil
		}
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
				anpSplit := strings.Split(anpRef, "/")
				if len(anpSplit) < 7 || anpSplit[6] != stateAnp {
					continue
				}
				epgCount, err := anpCont.ArrayCount("epgs")
				if err != nil {
					continue
				}
				for k := 0; k < epgCount; k++ {
					epgCont, err := anpCont.ArrayElement(k, "epgs")
					if err != nil {
						return err
					}
					epgRef := models.StripQuotes(epgCont.S("epgRef").String())
					epgSplit := strings.Split(epgRef, "/")
					if len(epgSplit) >= 9 && epgSplit[8] == stateEpg {
						return fmt.Errorf("mso_schema_site_anp_epg (site=%s, template=%s, anp=%s, epg=%s) still exists on schema %s", stateSiteId, stateTemplate, stateAnp, stateEpg, schemaId)
					}
				}
			}
		}
	}
	return nil
}
