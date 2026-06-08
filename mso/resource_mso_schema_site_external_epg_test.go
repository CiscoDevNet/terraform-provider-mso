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

// TestAccMSOSchemaSiteExternalEpgResource exercises the lifecycle of
// mso_schema_site_external_epg:
//   - attempt to create the site external EPG on a template that has no
//     mso_schema_site association (older NDO returns "Resource Not Found";
//     newer NDO's always-on schema validation engine silently drops the
//     PATCH and the SDK reports "Provider produced inconsistent result
//     after apply")
//   - create the site external EPG with no L3Out attached
//   - update by attaching a template-managed L3Out
//   - update by switching the L3Out to l3out_on_apic = true (APIC-managed)
//   - update by detaching the L3Out (clear l3out_name)
//   - import the site external EPG
//
// Children of the site external EPG (selectors etc.) are intentionally not
// covered here.
//
// The lab must have the `ansible_test` and `ansible_test_2` sites onboarded.
func TestAccMSOSchemaSiteExternalEpgResource(t *testing.T) {
	siteEpgResource := "mso_schema_site_external_epg." + msoSchemaTemplateExtEpgName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteExternalEpgDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: Create site external EPG without mso_schema_site association (expect error)")
				},
				Config: testAccMSOSchemaSiteExternalEpgConfigNoSiteAssociation(),
				// Older NDO rejects both the "replace" and the fallback "add"
				// PATCH with "Resource Not Found", so Create returns the raw
				// error. Newer NDO's always-on schema validation engine
				// silently drops the PATCH as a no-op, so Create succeeds,
				// the follow-up Read finds nothing, and the SDK raises
				// "Provider produced inconsistent result after apply".
				// Match either outcome.
				ExpectError: regexp.MustCompile(`Resource Not Found|Provider produced inconsistent result after apply`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create site external EPG without L3Out") },
				Config:    testAccMSOSchemaSiteExternalEpgConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(siteEpgResource, "schema_id"),
					resource.TestCheckResourceAttrSet(siteEpgResource, "site_id"),
					resource.TestCheckResourceAttr(siteEpgResource, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr(siteEpgResource, "external_epg_name", msoSchemaTemplateExtEpgName),
					resource.TestCheckResourceAttrPair(
						siteEpgResource, "site_id",
						"data.mso_site."+msoTemplateSiteName1, "id",
					),
					resource.TestCheckResourceAttr(siteEpgResource, "id", msoSchemaTemplateExtEpgName),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Attach a template-managed L3Out") },
				Config:    testAccMSOSchemaSiteExternalEpgConfigWithTemplateL3Out(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(siteEpgResource, "external_epg_name", msoSchemaTemplateExtEpgName),
					resource.TestCheckResourceAttr(siteEpgResource, "l3out_name", msoSchemaTemplateL3outName),
					resource.TestCheckResourceAttr(siteEpgResource, "l3out_template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttrPair(
						siteEpgResource, "l3out_schema_id",
						"mso_schema."+msoSchemaName, "id",
					),
					resource.TestCheckResourceAttr(siteEpgResource, "l3out_on_apic", "false"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Switch L3Out to APIC-managed (l3out_on_apic=true)") },
				Config:    testAccMSOSchemaSiteExternalEpgConfigWithApicL3Out(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(siteEpgResource, "external_epg_name", msoSchemaTemplateExtEpgName),
					resource.TestCheckResourceAttr(siteEpgResource, "l3out_name", msoSchemaTemplateL3outName),
					resource.TestCheckResourceAttr(siteEpgResource, "l3out_on_apic", "true"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Detach the L3Out by clearing l3out_name") },
				Config:    testAccMSOSchemaSiteExternalEpgConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(siteEpgResource, "external_epg_name", msoSchemaTemplateExtEpgName),
					resource.TestCheckResourceAttr(siteEpgResource, "l3out_name", ""),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import site external EPG") },
				ResourceName: siteEpgResource,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[siteEpgResource]
					if !ok {
						return "", fmt.Errorf("site external EPG resource not found in state: %s", siteEpgResource)
					}
					// Importer reads index [0] (schema id), [2] (site id),
					// [4] (external EPG name).
					return fmt.Sprintf("%s/site/%s/externalEpg/%s",
						rs.Primary.Attributes["schema_id"],
						rs.Primary.Attributes["site_id"],
						rs.Primary.Attributes["external_epg_name"],
					), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

// testAccMSOSchemaSiteExternalEpgPrerequisiteConfig emits the standard
// schema + single-site-association + template VRF + template external EPG
// + template L3Out scaffolding required by the site external EPG steps.
func testAccMSOSchemaSiteExternalEpgPrerequisiteConfig() string {
	return fmt.Sprintf(`%s%s%s%s`,
		testSchemaWithSingleSiteAssociationConfig(),
		testSchemaTemplateVrfConfig(),
		testSchemaTemplateL3outConfig(),
		testSchemaTemplateExtEpgConfig(),
	)
}

// testAccMSOSchemaSiteExternalEpgNoSiteAssociationConfig is the prereq
// without the mso_schema_site block, used to drive the negative-path step.
func testAccMSOSchemaSiteExternalEpgNoSiteAssociationConfig() string {
	return fmt.Sprintf(`%s%s%s%s`,
		testSchemaWithBothSitesPrerequisiteConfig(),
		testSchemaTemplateVrfConfig(),
		testSchemaTemplateL3outConfig(),
		testSchemaTemplateExtEpgConfig(),
	)
}

func testAccMSOSchemaSiteExternalEpgConfigCreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_external_epg" "%[2]s" {
		schema_id         = mso_schema.%[3]s.id
		site_id           = mso_schema_site.%[4]s.site_id
		template_name     = "%[5]s"
		external_epg_name = mso_schema_template_external_epg.%[2]s.external_epg_name
	}`,
		testAccMSOSchemaSiteExternalEpgPrerequisiteConfig(),
		msoSchemaTemplateExtEpgName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
	)
}

func testAccMSOSchemaSiteExternalEpgConfigWithTemplateL3Out() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_external_epg" "%[2]s" {
		schema_id           = mso_schema.%[3]s.id
		site_id             = mso_schema_site.%[4]s.site_id
		template_name       = "%[5]s"
		external_epg_name   = mso_schema_template_external_epg.%[2]s.external_epg_name
		l3out_name          = mso_schema_template_l3out.%[6]s.l3out_name
		l3out_template_name = "%[5]s"
		l3out_schema_id     = mso_schema.%[3]s.id
	}`,
		testAccMSOSchemaSiteExternalEpgPrerequisiteConfig(),
		msoSchemaTemplateExtEpgName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
		msoSchemaTemplateL3outName,
	)
}

// testAccMSOSchemaSiteExternalEpgConfigWithApicL3Out points at an
// APIC-managed L3Out (l3out_on_apic = true). l3out_name still has to be
// supplied because the resource builds the l3outDn from it; the L3Out
// referenced does not need to exist on MSO/NDO -- only the matching
// uni/tn-<tenant>/out-<l3out_name> on the fabric -- but for the purposes
// of this acceptance test we reuse the same name as the template L3Out.
// l3out_template_name and l3out_schema_id must NOT be set when
// l3out_on_apic is true (ConflictsWith on the schema).
func testAccMSOSchemaSiteExternalEpgConfigWithApicL3Out() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_external_epg" "%[2]s" {
		schema_id         = mso_schema.%[3]s.id
		site_id           = mso_schema_site.%[4]s.site_id
		template_name     = "%[5]s"
		external_epg_name = mso_schema_template_external_epg.%[2]s.external_epg_name
		l3out_name        = "%[6]s"
		l3out_on_apic     = true
	}`,
		testAccMSOSchemaSiteExternalEpgPrerequisiteConfig(),
		msoSchemaTemplateExtEpgName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
		msoSchemaTemplateL3outName,
	)
}

func testAccMSOSchemaSiteExternalEpgConfigNoSiteAssociation() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_external_epg" "%[2]s" {
		schema_id         = mso_schema.%[3]s.id
		site_id           = data.mso_site.%[4]s.id
		template_name     = "%[5]s"
		external_epg_name = mso_schema_template_external_epg.%[2]s.external_epg_name
	}`,
		testAccMSOSchemaSiteExternalEpgNoSiteAssociationConfig(),
		msoSchemaTemplateExtEpgName,
		msoSchemaName,
		msoTemplateSiteName1,
		msoSchemaTemplateName,
	)
}

// testAccCheckMSOSchemaSiteExternalEpgDestroy walks state for any
// mso_schema_site_external_epg resources and asserts that no
// sites[].externalEpgs[] entry whose externalEpgRef resolves to
// external_epg_name remains under the matching siteId + templateName.
// A missing schema or missing sites array is treated as a successful
// destroy.
func testAccCheckMSOSchemaSiteExternalEpgDestroy(s *terraform.State) error {
	msoClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mso_schema_site_external_epg" {
			continue
		}
		schemaId := rs.Primary.Attributes["schema_id"]
		stateSiteId := rs.Primary.Attributes["site_id"]
		stateTemplate := rs.Primary.Attributes["template_name"]
		stateEpg := rs.Primary.Attributes["external_epg_name"]

		cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
		if err != nil {
			return nil
		}
		count, err := cont.ArrayCount("sites")
		if err != nil {
			return nil
		}
		re := regexp.MustCompile("/schemas/(.*?)/templates/(.*?)/externalEpgs/(.*)")
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
			epgCount, err := siteCont.ArrayCount("externalEpgs")
			if err != nil {
				continue
			}
			for j := 0; j < epgCount; j++ {
				epgCont, err := siteCont.ArrayElement(j, "externalEpgs")
				if err != nil {
					return err
				}
				epgRef := models.StripQuotes(epgCont.S("externalEpgRef").String())
				match := re.FindStringSubmatch(epgRef)
				if len(match) == 4 && match[3] == stateEpg {
					return fmt.Errorf("mso_schema_site_external_epg (site=%s, template=%s, epg=%s) still exists on schema %s",
						stateSiteId, stateTemplate, stateEpg, schemaId)
				}
			}
		}
	}
	return nil
}
