package mso

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// TestAccMSOSchemaSiteAnpEpgBulkStaticPortResource exercises the full lifecycle
// of mso_schema_site_anp_epg_bulk_staticport:
//   - attempt to create without a mso_schema_site association (expect error)
//   - create two static ports and verify all attributes via TypeSet helpers
//   - update: modify one port's attributes and drop the second (verifies removal)
//   - import the resource
//
// The test requires msoSchemaSiteAnpEpgStaticPortPod/Leaf/Path/Path2 to
// correspond to real interfaces on a leaf switch onboarded to the ansible_test site.
//
// The lab must have the `ansible_test` and `ansible_test_2` sites onboarded.
func TestAccMSOSchemaSiteAnpEpgBulkStaticPortResource(t *testing.T) {
	bulkStaticPortResource := "mso_schema_site_anp_epg_bulk_staticport." + msoSchemaTemplateAnpEpgName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpEpgBulkStaticPortDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: Create bulk static port without mso_schema_site association (expect error)")
				},
				Config:      testAccMSOSchemaSiteAnpEpgBulkStaticPortConfigNoSiteAssociation(),
				ExpectError: regexp.MustCompile(`Site-Template association for .* is not found\.`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create bulk static port (two ports)") },
				Config:    testAccMSOSchemaSiteAnpEpgBulkStaticPortConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(bulkStaticPortResource, "schema_id"),
					resource.TestCheckResourceAttrSet(bulkStaticPortResource, "site_id"),
					resource.TestCheckResourceAttr(bulkStaticPortResource, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr(bulkStaticPortResource, "anp_name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr(bulkStaticPortResource, "epg_name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr(bulkStaticPortResource, "static_ports.#", "2"),
					CustomTestCheckTypeSetElemAttrs(bulkStaticPortResource, "static_ports", map[string]string{
						"path_type":            "port",
						"pod":                  msoSchemaSiteAnpEpgStaticPortPod,
						"leaf":                 msoSchemaSiteAnpEpgStaticPortLeaf,
						"path":                 msoSchemaSiteAnpEpgStaticPortPath,
						"vlan":                 "200",
						"deployment_immediacy": "lazy",
						"mode":                 "regular",
						"micro_seg_vlan":       "300",
					}),
					CustomTestCheckTypeSetElemAttrs(bulkStaticPortResource, "static_ports", map[string]string{
						"path_type":            "port",
						"pod":                  msoSchemaSiteAnpEpgStaticPortPod,
						"leaf":                 msoSchemaSiteAnpEpgStaticPortLeaf,
						"path":                 msoSchemaSiteAnpEpgStaticPortPath2,
						"vlan":                 "201",
						"deployment_immediacy": "immediate",
						"mode":                 "untagged",
						"fex":                  msoSchemaSiteAnpEpgStaticPortFex,
					}),
					resource.TestCheckResourceAttrPair(
						bulkStaticPortResource, "site_id",
						"data.mso_site."+msoTemplateSiteName1, "id",
					),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import bulk static port (two ports)") },
				ResourceName: bulkStaticPortResource,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[bulkStaticPortResource]
					if !ok {
						return "", fmt.Errorf("bulk static port resource not found in state: %s", bulkStaticPortResource)
					}
					// Import ID format (used by resourceMSOSchemaSiteAnpEpgBulkStaticPortImport):
					// {schemaId}/site/{siteId}/template/{templateName}/anp/{anpName}/epg/{epgName}
					return fmt.Sprintf("%s/site/%s/template/%s/anp/%s/epg/%s",
						rs.Primary.Attributes["schema_id"],
						rs.Primary.Attributes["site_id"],
						rs.Primary.Attributes["template_name"],
						rs.Primary.Attributes["anp_name"],
						rs.Primary.Attributes["epg_name"],
					), nil
				},
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update bulk static port (one port, changed vlan/mode/immediacy)")
				},
				Config: testAccMSOSchemaSiteAnpEpgBulkStaticPortConfigUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(bulkStaticPortResource, "static_ports.#", "1"),
					CustomTestCheckTypeSetElemAttrs(bulkStaticPortResource, "static_ports", map[string]string{
						"path_type":            "port",
						"pod":                  msoSchemaSiteAnpEpgStaticPortPod,
						"leaf":                 msoSchemaSiteAnpEpgStaticPortLeaf,
						"path":                 msoSchemaSiteAnpEpgStaticPortPath,
						"vlan":                 "202",
						"deployment_immediacy": "immediate",
						"mode":                 "native",
					}),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update bulk static port (remove all ports)")
				},
				Config: testAccMSOSchemaSiteAnpEpgBulkStaticPortConfigEmpty(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(bulkStaticPortResource, "static_ports.#", "0"),
				),
			},
		},
	})
}

func testAccMSOSchemaSiteAnpEpgBulkStaticPortConfigCreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_anp_epg_bulk_staticport" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = mso_schema_site.%[4]s.site_id
		template_name = "%[5]s"
		anp_name      = mso_schema_template_anp.%[6]s.name
		epg_name      = mso_schema_site_anp_epg.%[2]s.epg_name

		static_ports {
			path_type            = "port"
			pod                  = "%[7]s"
			leaf                 = "%[8]s"
			path                 = "%[9]s"
			vlan                 = 200
			deployment_immediacy = "lazy"
			mode                 = "regular"
			micro_seg_vlan       = 300
		}

		static_ports {
			path_type            = "port"
			pod                  = "%[7]s"
			leaf                 = "%[8]s"
			path                 = "%[10]s"
			vlan                 = 201
			deployment_immediacy = "immediate"
			mode                 = "untagged"
			fex                  = "%[11]s"
		}
	}`,
		testAccMSOSchemaSiteAnpEpgStaticLeafPrerequisiteConfig(),
		msoSchemaTemplateAnpEpgName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
		msoSchemaTemplateAnpName,
		msoSchemaSiteAnpEpgStaticPortPod,
		msoSchemaSiteAnpEpgStaticPortLeaf,
		msoSchemaSiteAnpEpgStaticPortPath,
		msoSchemaSiteAnpEpgStaticPortPath2,
		msoSchemaSiteAnpEpgStaticPortFex,
	)
}

func testAccMSOSchemaSiteAnpEpgBulkStaticPortConfigUpdate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_anp_epg_bulk_staticport" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = mso_schema_site.%[4]s.site_id
		template_name = "%[5]s"
		anp_name      = mso_schema_template_anp.%[6]s.name
		epg_name      = mso_schema_site_anp_epg.%[2]s.epg_name

		static_ports {
			path_type            = "port"
			pod                  = "%[7]s"
			leaf                 = "%[8]s"
			path                 = "%[9]s"
			vlan                 = 202
			deployment_immediacy = "immediate"
			mode                 = "native"
		}
	}`,
		testAccMSOSchemaSiteAnpEpgStaticLeafPrerequisiteConfig(),
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

func testAccMSOSchemaSiteAnpEpgBulkStaticPortConfigEmpty() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_anp_epg_bulk_staticport" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = mso_schema_site.%[4]s.site_id
		template_name = "%[5]s"
		anp_name      = mso_schema_template_anp.%[6]s.name
		epg_name      = mso_schema_site_anp_epg.%[2]s.epg_name
	}`,
		testAccMSOSchemaSiteAnpEpgStaticLeafPrerequisiteConfig(),
		msoSchemaTemplateAnpEpgName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
		msoSchemaTemplateAnpName,
	)
}

// testAccMSOSchemaSiteAnpEpgBulkStaticPortConfigNoSiteAssociation creates a
// bulk static port without a prior mso_schema_site association. Because
// resourceMSOSchemaSiteAnpEpgBulkStaticPortCreate calls
// getSiteFromSiteIdAndTemplate before any PATCH, NDO returns
// "Site-Template association for X-Y is not found." immediately.
func testAccMSOSchemaSiteAnpEpgBulkStaticPortConfigNoSiteAssociation() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_anp_epg_bulk_staticport" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = data.mso_site.%[4]s.id
		template_name = "%[5]s"
		anp_name      = mso_schema_template_anp.%[6]s.name
		epg_name      = mso_schema_template_anp_epg.%[2]s.name

		static_ports {
			path_type            = "port"
			pod                  = "%[7]s"
			leaf                 = "%[8]s"
			path                 = "%[9]s"
			vlan                 = 200
			deployment_immediacy = "lazy"
			mode                 = "regular"
			micro_seg_vlan       = 300
		}
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
		msoSchemaSiteAnpEpgStaticPortPod,
		msoSchemaSiteAnpEpgStaticPortLeaf,
		msoSchemaSiteAnpEpgStaticPortPath,
	)
}

// testAccCheckMSOSchemaSiteAnpEpgBulkStaticPortDestroy walks state for any
// mso_schema_site_anp_epg_bulk_staticport resources, fetches the schema, and
// asserts that the EPG's staticPorts[] array is empty. A missing schema or
// missing sites array is treated as a successful destroy.
func testAccCheckMSOSchemaSiteAnpEpgBulkStaticPortDestroy(s *terraform.State) error {
	msoClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mso_schema_site_anp_epg_bulk_staticport" {
			continue
		}
		schemaId := rs.Primary.Attributes["schema_id"]
		stateSiteId := rs.Primary.Attributes["site_id"]
		stateTemplate := rs.Primary.Attributes["template_name"]
		stateAnp := rs.Primary.Attributes["anp_name"]
		stateEpg := rs.Primary.Attributes["epg_name"]

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
						// No staticPorts key — already empty.
						continue
					}
					if portCount > 0 {
						return fmt.Errorf(
							"mso_schema_site_anp_epg_bulk_staticport (site=%s, template=%s, anp=%s, epg=%s) still has %d static port(s) on schema %s",
							stateSiteId, stateTemplate, stateAnp, stateEpg, portCount, schemaId,
						)
					}
				}
			}
		}
	}
	return nil
}
