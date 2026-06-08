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

// TestAccMSOSchemaSiteAnpEpgStaticPortResource exercises the full lifecycle of
// mso_schema_site_anp_epg_static_port:
//   - attempt to create without a mso_schema_site association (expect error)
//   - create the static port and verify all attributes
//   - update the mutable fields (vlan, mode) and verify changes
//   - import the static port
//
// The test requires msoSchemaSiteAnpEpgStaticPortPod/Leaf/Path to correspond
// to a real interface on a leaf switch onboarded to the ansible_test site.
//
// The lab must have the `ansible_test` and `ansible_test_2` sites onboarded.
func TestAccMSOSchemaSiteAnpEpgStaticPortResource(t *testing.T) {
	staticPortResource := "mso_schema_site_anp_epg_static_port." + msoSchemaTemplateAnpEpgName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpEpgStaticPortDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: Create static port without mso_schema_site association (expect error)")
				},
				Config: testAccMSOSchemaSiteAnpEpgStaticPortConfigNoSiteAssociation(),
				// Older NDO rejects the PATCH with "Resource Not Found". Newer
				// NDO's always-on schema validation engine silently drops the
				// PATCH so Create succeeds, the follow-up Read finds nothing,
				// and the SDK raises "Provider produced inconsistent result
				// after apply". Match either outcome.
				ExpectError: regexp.MustCompile(`Resource Not Found|Provider produced inconsistent result after apply`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create static port") },
				Config:    testAccMSOSchemaSiteAnpEpgStaticPortConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(staticPortResource, "schema_id"),
					resource.TestCheckResourceAttrSet(staticPortResource, "site_id"),
					resource.TestCheckResourceAttr(staticPortResource, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr(staticPortResource, "anp_name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr(staticPortResource, "epg_name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr(staticPortResource, "path_type", "port"),
					resource.TestCheckResourceAttr(staticPortResource, "pod", msoSchemaSiteAnpEpgStaticPortPod),
					resource.TestCheckResourceAttr(staticPortResource, "leaf", msoSchemaSiteAnpEpgStaticPortLeaf),
					resource.TestCheckResourceAttr(staticPortResource, "path", msoSchemaSiteAnpEpgStaticPortPath),
					resource.TestCheckResourceAttr(staticPortResource, "vlan", "200"),
					resource.TestCheckResourceAttr(staticPortResource, "micro_seg_vlan", "300"),
					resource.TestCheckResourceAttr(staticPortResource, "deployment_immediacy", "lazy"),
					resource.TestCheckResourceAttr(staticPortResource, "mode", "regular"),
					resource.TestCheckResourceAttr(staticPortResource, "fex", msoSchemaSiteAnpEpgStaticPortFex),
					resource.TestCheckResourceAttrPair(
						staticPortResource, "site_id",
						"data.mso_site."+msoTemplateSiteName1, "id",
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update static port (vlan, micro_seg_vlan, mode, deployment_immediacy)") },
				Config:    testAccMSOSchemaSiteAnpEpgStaticPortConfigUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(staticPortResource, "vlan", "201"),
					resource.TestCheckResourceAttr(staticPortResource, "micro_seg_vlan", "301"),
					resource.TestCheckResourceAttr(staticPortResource, "mode", "untagged"),
					resource.TestCheckResourceAttr(staticPortResource, "deployment_immediacy", "immediate"),
					// ForceNew identity fields must be unchanged after an in-place update.
					resource.TestCheckResourceAttr(staticPortResource, "fex", msoSchemaSiteAnpEpgStaticPortFex),
					resource.TestCheckResourceAttr(staticPortResource, "path_type", "port"),
					resource.TestCheckResourceAttr(staticPortResource, "pod", msoSchemaSiteAnpEpgStaticPortPod),
					resource.TestCheckResourceAttr(staticPortResource, "leaf", msoSchemaSiteAnpEpgStaticPortLeaf),
					resource.TestCheckResourceAttr(staticPortResource, "path", msoSchemaSiteAnpEpgStaticPortPath),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import static port") },
				ResourceName: staticPortResource,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[staticPortResource]
					if !ok {
						return "", fmt.Errorf("static port resource not found in state: %s", staticPortResource)
					}
					// Import ID format (indices used by resourceMSOSchemaSiteAnpEpgStaticPortImport):
					// {schemaId}/site/{siteId}/template/{templateName}/anp/{anpName}/epg/{epgName}/
					// pod/{pod}/leaf/{leaf}/path_type/{pathType}/fex/{fex}/path/{path}
					// When fex is empty the segment becomes "fex//path/...".
					return fmt.Sprintf("%s/site/%s/template/%s/anp/%s/epg/%s/pod/%s/leaf/%s/path_type/%s/fex/%s/path/%s",
						rs.Primary.Attributes["schema_id"],
						rs.Primary.Attributes["site_id"],
						rs.Primary.Attributes["template_name"],
						rs.Primary.Attributes["anp_name"],
						rs.Primary.Attributes["epg_name"],
						rs.Primary.Attributes["pod"],
						rs.Primary.Attributes["leaf"],
						rs.Primary.Attributes["path_type"],
						rs.Primary.Attributes["fex"],
						rs.Primary.Attributes["path"],
					), nil
				},
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Remove fex (ForceNew, destroy + recreate)")
				},
				Config: testAccMSOSchemaSiteAnpEpgStaticPortConfigRemoveFex(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(staticPortResource, "fex", ""),
					resource.TestCheckResourceAttr(staticPortResource, "vlan", "201"),
					resource.TestCheckResourceAttr(staticPortResource, "mode", "untagged"),
				),
			},
		},
	})
}

// testAccMSOSchemaSiteAnpEpgStaticPortPrerequisiteConfig extends the site ANP
// EPG config with an explicit mso_schema_site_anp_epg so the static port
// PATCH has a well-formed parent EPG entry in the schema document.
func testAccMSOSchemaSiteAnpEpgStaticPortPrerequisiteConfig() string {
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

func testAccMSOSchemaSiteAnpEpgStaticPortConfigCreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_anp_epg_static_port" "%[2]s" {
		schema_id            = mso_schema.%[3]s.id
		site_id              = mso_schema_site.%[4]s.site_id
		template_name        = "%[5]s"
		anp_name             = mso_schema_template_anp.%[6]s.name
		epg_name             = mso_schema_site_anp_epg.%[2]s.epg_name
		path_type            = "port"
		pod                  = "%[7]s"
		leaf                 = "%[8]s"
		path                 = "%[9]s"
		vlan                 = 200
		micro_seg_vlan       = 300
		deployment_immediacy = "lazy"
		mode                 = "regular"
		fex                  = "%[10]s"
	}`,
		testAccMSOSchemaSiteAnpEpgStaticPortPrerequisiteConfig(),
		msoSchemaTemplateAnpEpgName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
		msoSchemaTemplateAnpName,
		msoSchemaSiteAnpEpgStaticPortPod,
		msoSchemaSiteAnpEpgStaticPortLeaf,
		msoSchemaSiteAnpEpgStaticPortPath,
		msoSchemaSiteAnpEpgStaticPortFex,
	)
}

func testAccMSOSchemaSiteAnpEpgStaticPortConfigUpdate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_anp_epg_static_port" "%[2]s" {
		schema_id            = mso_schema.%[3]s.id
		site_id              = mso_schema_site.%[4]s.site_id
		template_name        = "%[5]s"
		anp_name             = mso_schema_template_anp.%[6]s.name
		epg_name             = mso_schema_site_anp_epg.%[2]s.epg_name
		path_type            = "port"
		pod                  = "%[7]s"
		leaf                 = "%[8]s"
		path                 = "%[9]s"
		vlan                 = 201
		micro_seg_vlan       = 301
		deployment_immediacy = "immediate"
		mode                 = "untagged"
		fex                  = "%[10]s"
	}`,
		testAccMSOSchemaSiteAnpEpgStaticPortPrerequisiteConfig(),
		msoSchemaTemplateAnpEpgName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
		msoSchemaTemplateAnpName,
		msoSchemaSiteAnpEpgStaticPortPod,
		msoSchemaSiteAnpEpgStaticPortLeaf,
		msoSchemaSiteAnpEpgStaticPortPath,
		msoSchemaSiteAnpEpgStaticPortFex,
	)
}

// testAccMSOSchemaSiteAnpEpgStaticPortConfigRemoveFex is the same as the
// update config but omits fex, triggering a destroy + recreate because fex is
// ForceNew. Confirms NDO correctly removes the FEX path on recreate.
func testAccMSOSchemaSiteAnpEpgStaticPortConfigRemoveFex() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_anp_epg_static_port" "%[2]s" {
		schema_id            = mso_schema.%[3]s.id
		site_id              = mso_schema_site.%[4]s.site_id
		template_name        = "%[5]s"
		anp_name             = mso_schema_template_anp.%[6]s.name
		epg_name             = mso_schema_site_anp_epg.%[2]s.epg_name
		path_type            = "port"
		pod                  = "%[7]s"
		leaf                 = "%[8]s"
		path                 = "%[9]s"
		vlan                 = 201
		micro_seg_vlan       = 301
		deployment_immediacy = "immediate"
		mode                 = "untagged"
	}`,
		testAccMSOSchemaSiteAnpEpgStaticPortPrerequisiteConfig(),
		msoSchemaTemplateAnpEpgName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
		msoSchemaTemplateAnpName,
		msoSchemaSiteAnpEpgStaticPortPod,
		msoSchemaSiteAnpEpgStaticPortLeaf,
		msoSchemaSiteAnpEpgStaticPortPath,
	)
}

// testAccMSOSchemaSiteAnpEpgStaticPortConfigNoSiteAssociation creates a static
// port without a prior mso_schema_site association, exercising the negative
// path. The PATCH targets a site that has no entry in the schema document, so
// NDO rejects it with "Resource Not Found".
func testAccMSOSchemaSiteAnpEpgStaticPortConfigNoSiteAssociation() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_anp_epg_static_port" "%[2]s" {
		schema_id            = mso_schema.%[3]s.id
		site_id              = data.mso_site.%[4]s.id
		template_name        = "%[5]s"
		anp_name             = mso_schema_template_anp.%[6]s.name
		epg_name             = mso_schema_template_anp_epg.%[7]s.name
		path_type            = "port"
		pod                  = "%[8]s"
		leaf                 = "%[9]s"
		path                 = "%[10]s"
		vlan                 = 200
		deployment_immediacy = "lazy"
		mode                 = "regular"
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
		msoSchemaSiteAnpEpgStaticPortPod,
		msoSchemaSiteAnpEpgStaticPortLeaf,
		msoSchemaSiteAnpEpgStaticPortPath,
	)
}

// testAccCheckMSOSchemaSiteAnpEpgStaticPortDestroy walks state for any
// mso_schema_site_anp_epg_static_port resources, fetches the schema, and
// asserts that no sites[].anps[].epgs[].staticPorts[] entry whose path
// matches the portPath still exists. A missing schema or missing sites array
// is treated as a successful destroy.
func testAccCheckMSOSchemaSiteAnpEpgStaticPortDestroy(s *terraform.State) error {
	msoClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mso_schema_site_anp_epg_static_port" {
			continue
		}
		schemaId := rs.Primary.Attributes["schema_id"]
		stateSiteId := rs.Primary.Attributes["site_id"]
		stateTemplate := rs.Primary.Attributes["template_name"]
		stateAnp := rs.Primary.Attributes["anp_name"]
		stateEpg := rs.Primary.Attributes["epg_name"]
		// Reconstruct the portPath that identifies the entry in the array.
		portPath := createPortPath(
			rs.Primary.Attributes["path_type"],
			rs.Primary.Attributes["pod"],
			rs.Primary.Attributes["leaf"],
			rs.Primary.Attributes["fex"],
			rs.Primary.Attributes["path"],
		)

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
					portCount, err := epgCont.ArrayCount("staticPorts")
					if err != nil {
						continue
					}
					for l := 0; l < portCount; l++ {
						portCont, err := epgCont.ArrayElement(l, "staticPorts")
						if err != nil {
							return err
						}
						if models.StripQuotes(portCont.S("path").String()) == portPath {
							return fmt.Errorf(
								"mso_schema_site_anp_epg_static_port (site=%s, template=%s, anp=%s, epg=%s, portPath=%s) still exists on schema %s",
								stateSiteId, stateTemplate, stateAnp, stateEpg, portPath, schemaId,
							)
						}
					}
				}
			}
		}
	}
	return nil
}
