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

// TestAccMSOSchemaSiteAnpEpgStaticLeafResource exercises the lifecycle of
// mso_schema_site_anp_epg_static_leaf:
//   - attempt to create without a mso_schema_site association (expect error)
//   - create the static leaf and verify all attributes
//   - import the static leaf
//
// All resource attributes are ForceNew so there is no update step.
// The test requires msoSchemaSiteAnpEpgStaticLeafPath to correspond to a real
// leaf switch onboarded to the ansible_test site.
//
// The lab must have the `ansible_test` and `ansible_test_2` sites onboarded.
func TestAccMSOSchemaSiteAnpEpgStaticLeafResource(t *testing.T) {
	staticLeafResource := "mso_schema_site_anp_epg_static_leaf." + msoSchemaTemplateAnpEpgName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpEpgStaticLeafDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: Create static leaf without mso_schema_site association (expect error)")
				},
				Config: testAccMSOSchemaSiteAnpEpgStaticLeafConfigNoSiteAssociation(),
				// Older NDO rejects the PATCH with "Resource Not Found". Newer
				// NDO's always-on schema validation engine silently drops the
				// PATCH so Create succeeds, the follow-up Read finds nothing,
				// and the SDK raises "Provider produced inconsistent result
				// after apply". Match either outcome.
				ExpectError: regexp.MustCompile(`Resource Not Found|Provider produced inconsistent result after apply`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create static leaf") },
				Config:    testAccMSOSchemaSiteAnpEpgStaticLeafConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(staticLeafResource, "schema_id"),
					resource.TestCheckResourceAttrSet(staticLeafResource, "site_id"),
					resource.TestCheckResourceAttr(staticLeafResource, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr(staticLeafResource, "anp_name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr(staticLeafResource, "epg_name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr(staticLeafResource, "path", msoSchemaSiteAnpEpgStaticLeafPath),
					resource.TestCheckResourceAttr(staticLeafResource, "port_encap_vlan", "100"),
					resource.TestCheckResourceAttrPair(
						staticLeafResource, "site_id",
						"data.mso_site."+msoTemplateSiteName1, "id",
					),
					resource.TestCheckResourceAttr(staticLeafResource, "id", msoSchemaSiteAnpEpgStaticLeafPath),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Change port_encap_vlan (ForceNew - expect destroy+recreate)") },
				Config:    testAccMSOSchemaSiteAnpEpgStaticLeafConfigRecreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// port_encap_vlan is ForceNew so a change destroys the old entry and
					// creates a new one. Verify the new value is present and the identity
					// fields (path, id) remain unchanged.
					resource.TestCheckResourceAttr(staticLeafResource, "port_encap_vlan", "101"),
					resource.TestCheckResourceAttr(staticLeafResource, "path", msoSchemaSiteAnpEpgStaticLeafPath),
					resource.TestCheckResourceAttr(staticLeafResource, "id", msoSchemaSiteAnpEpgStaticLeafPath),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import static leaf") },
				ResourceName: staticLeafResource,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[staticLeafResource]
					if !ok {
						return "", fmt.Errorf("static leaf resource not found in state: %s", staticLeafResource)
					}
					// Importer splits by "/" and uses a regex "(.*)/path/(.*)"
					// to extract the path (which itself contains "/").
					return fmt.Sprintf("%s/site/%s/template/%s/anp/%s/epg/%s/path/%s",
						rs.Primary.Attributes["schema_id"],
						rs.Primary.Attributes["site_id"],
						rs.Primary.Attributes["template_name"],
						rs.Primary.Attributes["anp_name"],
						rs.Primary.Attributes["epg_name"],
						rs.Primary.Attributes["path"],
					), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

// testAccMSOSchemaSiteAnpEpgStaticLeafPrerequisiteConfig extends the site ANP
// EPG config with an explicit mso_schema_site_anp_epg so the static leaf
// PATCH has a well-formed parent EPG entry in the schema document.
func testAccMSOSchemaSiteAnpEpgStaticLeafPrerequisiteConfig() string {
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

func testAccMSOSchemaSiteAnpEpgStaticLeafConfigRecreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_anp_epg_static_leaf" "%[2]s" {
		schema_id       = mso_schema.%[3]s.id
		site_id         = mso_schema_site.%[4]s.site_id
		template_name   = "%[5]s"
		anp_name        = mso_schema_template_anp.%[6]s.name
		epg_name        = mso_schema_site_anp_epg.%[2]s.epg_name
		path            = "%[7]s"
		port_encap_vlan = 101
	}`,
		testAccMSOSchemaSiteAnpEpgStaticLeafPrerequisiteConfig(),
		msoSchemaTemplateAnpEpgName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
		msoSchemaTemplateAnpName,
		msoSchemaSiteAnpEpgStaticLeafPath,
	)
}

func testAccMSOSchemaSiteAnpEpgStaticLeafConfigCreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_anp_epg_static_leaf" "%[2]s" {
		schema_id       = mso_schema.%[3]s.id
		site_id         = mso_schema_site.%[4]s.site_id
		template_name   = "%[5]s"
		anp_name        = mso_schema_template_anp.%[6]s.name
		epg_name        = mso_schema_site_anp_epg.%[2]s.epg_name
		path            = "%[7]s"
		port_encap_vlan = 100
	}`,
		testAccMSOSchemaSiteAnpEpgStaticLeafPrerequisiteConfig(),
		msoSchemaTemplateAnpEpgName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
		msoSchemaTemplateAnpName,
		msoSchemaSiteAnpEpgStaticLeafPath,
	)
}

// testAccMSOSchemaSiteAnpEpgStaticLeafConfigNoSiteAssociation creates a static
// leaf without a prior mso_schema_site association, exercising the negative
// path. The PATCH targets a site that has no entry in the schema document, so
// NDO rejects it with "Resource Not Found".
func testAccMSOSchemaSiteAnpEpgStaticLeafConfigNoSiteAssociation() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_anp_epg_static_leaf" "%[2]s" {
		schema_id       = mso_schema.%[3]s.id
		site_id         = data.mso_site.%[4]s.id
		template_name   = "%[5]s"
		anp_name        = mso_schema_template_anp.%[6]s.name
		epg_name        = mso_schema_template_anp_epg.%[7]s.name
		path            = "%[8]s"
		port_encap_vlan = 100
	}`,
		fmt.Sprintf(`%s%s%s%s`,
			testSchemaWithBothSitesPrerequisiteConfig(),
			testSchemaTemplateVrfConfig(),
			testSchemaTemplateBdConfig(),
			testSchemaTemplateAnpConfig(),
		)+testAccMSOSchemaSiteAnpEpgTemplateAnpEpgWithBdConfig(),
		msoSchemaTemplateAnpEpgName,
		msoSchemaName,
		msoTemplateSiteName1,
		msoSchemaTemplateName,
		msoSchemaTemplateAnpName,
		msoSchemaTemplateAnpEpgName,
		msoSchemaSiteAnpEpgStaticLeafPath,
	)
}

// testAccCheckMSOSchemaSiteAnpEpgStaticLeafDestroy walks state for any
// mso_schema_site_anp_epg_static_leaf resources, fetches the schema, and
// asserts that no sites[].anps[].epgs[].staticLeafs[] entry whose path
// matches still exists. A missing schema or missing sites array is treated
// as a successful destroy.
func testAccCheckMSOSchemaSiteAnpEpgStaticLeafDestroy(s *terraform.State) error {
	msoClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mso_schema_site_anp_epg_static_leaf" {
			continue
		}
		schemaId := rs.Primary.Attributes["schema_id"]
		stateSiteId := rs.Primary.Attributes["site_id"]
		stateTemplate := rs.Primary.Attributes["template_name"]
		stateAnp := rs.Primary.Attributes["anp_name"]
		stateEpg := rs.Primary.Attributes["epg_name"]
		statePath := rs.Primary.Attributes["path"]

		cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
		if err != nil {
			// Schema itself has been destroyed.
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
			if models.StripQuotes(siteCont.S("siteId").String()) != stateSiteId {
				continue
			}
			if models.StripQuotes(siteCont.S("templateName").String()) != stateTemplate {
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
					if len(epgSplit) < 9 || epgSplit[8] != stateEpg {
						continue
					}
					leafCount, err := epgCont.ArrayCount("staticLeafs")
					if err != nil {
						continue
					}
					for l := 0; l < leafCount; l++ {
						leafCont, err := epgCont.ArrayElement(l, "staticLeafs")
						if err != nil {
							return err
						}
						if models.StripQuotes(leafCont.S("path").String()) == statePath {
							return fmt.Errorf(
								"mso_schema_site_anp_epg_static_leaf (site=%s, template=%s, anp=%s, epg=%s, path=%s) still exists on schema %s",
								stateSiteId, stateTemplate, stateAnp, stateEpg, statePath, schemaId,
							)
						}
					}
				}
			}
		}
	}
	return nil
}
