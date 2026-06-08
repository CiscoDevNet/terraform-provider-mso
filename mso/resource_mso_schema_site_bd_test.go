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

// TestAccMSOSchemaSiteBdResource exercises the lifecycle of mso_schema_site_bd:
//   - attempt to create the site BD on a template that has no
//     mso_schema_site association (older NDO returns "Resource Not Found";
//     newer NDO silently drops the PATCH and the SDK reports "Provider
//     produced inconsistent result after apply" -- see the comment on the
//     ExpectError step for the underlying validation-engine cause)
//   - create the site BD and verify host_route/svi_mac defaults
//   - update host_route + svi_mac (in-place update path)
//   - update svi_mac only
//   - attach an mso_schema_site_bd_l3out child resource
//   - update the BD (svi_mac) while the L3Out is still attached and verify
//     the BD's l3Outs array on the API is not wiped by the BD PATCH
//   - drop svi_mac from config and flip host_route off
//   - change bd_name to a second template BD; because bd_name is
//     ForceNew, this exercises destroy+recreate (and verifies the new
//     bdRef resolves on the NDO side)
//   - import the site BD
//
// The lab must have the `ansible_test` and `ansible_test_2` sites onboarded.
func TestAccMSOSchemaSiteBdResource(t *testing.T) {
	siteBdResource := "mso_schema_site_bd." + msoSchemaTemplateBdName
	siteBdL3outResource := "mso_schema_site_bd_l3out." + msoSchemaTemplateBdName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteBdDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: Create site BD without mso_schema_site association (expect error)")
				},
				Config: testAccMSOSchemaSiteBdConfigNoSiteAssociation(),
				// Older MSO/NDO rejects both the "replace" and the fallback
				// "add" PATCH with "Resource Not Found", so Create returns
				// the raw error. Newer NDO appear to silently accept the
				// fallback "add" PATCH as a no-op, so Create succeeds, the
				// follow-up Read finds nothing, and the SDK raises
				// "Provider produced inconsistent result after apply".
				// Match either outcome.
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
				PreConfig: func() { fmt.Println("Test: Create site BD with defaults") },
				Config:    testAccMSOSchemaSiteBdConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(siteBdResource, "schema_id"),
					resource.TestCheckResourceAttrSet(siteBdResource, "site_id"),
					resource.TestCheckResourceAttr(siteBdResource, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr(siteBdResource, "bd_name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttrPair(
						siteBdResource, "site_id",
						"data.mso_site."+msoTemplateSiteName1, "id",
					),
					resource.TestCheckResourceAttr(siteBdResource, "id", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr(siteBdResource, "host_route", "false"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update host_route=true and svi_mac") },
				Config:    testAccMSOSchemaSiteBdConfigUpdateHostRouteAndSviMac(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(siteBdResource, "bd_name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr(siteBdResource, "host_route", "true"),
					resource.TestCheckResourceAttr(siteBdResource, "svi_mac", "00:00:5E:00:01:02"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update svi_mac only") },
				Config:    testAccMSOSchemaSiteBdConfigUpdateSviMacOnly(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(siteBdResource, "bd_name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr(siteBdResource, "host_route", "true"),
					resource.TestCheckResourceAttr(siteBdResource, "svi_mac", "00:00:5E:00:01:03"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Attach mso_schema_site_bd_l3out child resource") },
				Config:    testAccMSOSchemaSiteBdConfigWithL3Out(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(siteBdResource, "svi_mac", "00:00:5E:00:01:03"),
					resource.TestCheckResourceAttr(siteBdL3outResource, "l3out_name", msoSchemaTemplateL3outName),
					resource.TestCheckResourceAttr(siteBdL3outResource, "bd_name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr(siteBdL3outResource, "id", msoSchemaTemplateL3outName),
					testAccCheckMSOSchemaSiteBdHasL3Out(siteBdResource, msoSchemaTemplateL3outName),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update site BD svi_mac while L3Out attached -- l3Outs must survive")
				},
				Config: testAccMSOSchemaSiteBdConfigWithL3OutBdUpdated(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(siteBdResource, "svi_mac", "00:00:5E:00:01:04"),
					resource.TestCheckResourceAttr(siteBdResource, "host_route", "true"),
					resource.TestCheckResourceAttr(siteBdL3outResource, "l3out_name", msoSchemaTemplateL3outName),
					// Regression check: the BD Update PATCH must leave the
					// BD's l3Outs array intact -- this walks the schema on
					// MSO/NDO directly rather than trusting state alone.
					testAccCheckMSOSchemaSiteBdHasL3Out(siteBdResource, msoSchemaTemplateL3outName),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Flip host_route=false and drop svi_mac from config") },
				Config:    testAccMSOSchemaSiteBdConfigHostRouteFalse(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(siteBdResource, "bd_name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr(siteBdResource, "host_route", "false"),
					// svi_mac is Optional + Computed and the Update path is
					// gated on d.HasChange("svi_mac"). Dropping the attribute
					// from config does not trigger a diff (the SDK keeps the
					// previously applied value), so no `mac` PATCH is emitted
					// and the API still holds the value set in the prior step.
					resource.TestCheckResourceAttr(siteBdResource, "svi_mac", "00:00:5E:00:01:04"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Change bd_name to a second template BD (ForceNew destroy+recreate)")
				},
				Config: testAccMSOSchemaSiteBdConfigRename(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(siteBdResource, "bd_name", msoSchemaTemplateBdName+"_renamed"),
					resource.TestCheckResourceAttr(siteBdResource, "id", msoSchemaTemplateBdName+"_renamed"),
					resource.TestCheckResourceAttr(siteBdResource, "host_route", "false"),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import site BD") },
				ResourceName: siteBdResource,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[siteBdResource]
					if !ok {
						return "", fmt.Errorf("site BD resource not found in state: %s", siteBdResource)
					}
					return fmt.Sprintf("%s/%s/%s/%s",
						rs.Primary.Attributes["schema_id"],
						rs.Primary.Attributes["site_id"],
						rs.Primary.Attributes["template_name"],
						rs.Primary.Attributes["bd_name"],
					), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

// testAccMSOSchemaSiteBdPrerequisiteConfig emits the full scaffolding for a
// site BD: both data.mso_site blocks, the shared tenant attached to both
// sites, the schema with one template, the mso_schema_site association for
// site_1 (ansible_test, undeploy_on_destroy=false), the template VRF and
// the template BD that references the VRF.
func testAccMSOSchemaSiteBdPrerequisiteConfig() string {
	return fmt.Sprintf(`%s%s%s`,
		testSchemaWithSingleSiteAssociationConfig(),
		testSchemaTemplateVrfConfig(),
		testSchemaTemplateBdConfig(),
	)
}

// testAccMSOSchemaSiteBdNoSiteAssociationConfig is the prereq without the
// mso_schema_site block, used to drive the negative-path step.
func testAccMSOSchemaSiteBdNoSiteAssociationConfig() string {
	return fmt.Sprintf(`%s%s%s`,
		testSchemaWithBothSitesPrerequisiteConfig(),
		testSchemaTemplateVrfConfig(),
		testSchemaTemplateBdConfig(),
	)
}

func testAccMSOSchemaSiteBdConfigCreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_bd" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = mso_schema_site.%[4]s.site_id
		template_name = "%[5]s"
		bd_name       = mso_schema_template_bd.%[2]s.name
	}`,
		testAccMSOSchemaSiteBdPrerequisiteConfig(),
		msoSchemaTemplateBdName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
	)
}

func testAccMSOSchemaSiteBdConfigUpdateHostRouteAndSviMac() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_bd" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = mso_schema_site.%[4]s.site_id
		template_name = "%[5]s"
		bd_name       = mso_schema_template_bd.%[2]s.name
		host_route    = true
		svi_mac       = "00:00:5E:00:01:02"
	}`,
		testAccMSOSchemaSiteBdPrerequisiteConfig(),
		msoSchemaTemplateBdName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
	)
}

func testAccMSOSchemaSiteBdConfigUpdateSviMacOnly() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_bd" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = mso_schema_site.%[4]s.site_id
		template_name = "%[5]s"
		bd_name       = mso_schema_template_bd.%[2]s.name
		host_route    = true
		svi_mac       = "00:00:5E:00:01:03"
	}`,
		testAccMSOSchemaSiteBdPrerequisiteConfig(),
		msoSchemaTemplateBdName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
	)
}

// testAccMSOSchemaSiteBdConfigHostRouteFalse explicitly sets host_route=false
// (rather than relying on omission) because host_route is Optional + Computed:
// when omitted from config, the SDK plans no diff against the previously
// applied value, so step 5 would otherwise read back host_route=true. Setting
// the attribute back to false is what triggers the Update path.
func testAccMSOSchemaSiteBdConfigHostRouteFalse() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_bd" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = mso_schema_site.%[4]s.site_id
		template_name = "%[5]s"
		bd_name       = mso_schema_template_bd.%[2]s.name
		host_route    = false
	}`,
		testAccMSOSchemaSiteBdPrerequisiteConfig(),
		msoSchemaTemplateBdName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
	)
}

// testAccMSOSchemaSiteBdConfigRename adds a second mso_schema_template_bd
// (suffix "_renamed") and points the site BD at it. Because bd_name is
// ForceNew on mso_schema_site_bd, this drives a destroy+recreate of the
// site BD; the destination template BD must exist beforehand so the new
// bdRef resolves on the NDO side. NDO rejects multiple BDDelta entries
// for the same bd on the same fabric/template, which is why an in-place
// PATCH-based rename is not supported and ForceNew is required.
func testAccMSOSchemaSiteBdConfigRename() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_bd" "%[2]s_renamed" {
		schema_id              = mso_schema.%[3]s.id
		template_name          = "%[4]s"
		name                   = "%[2]s_renamed"
		display_name           = "%[2]s_renamed"
		layer2_unknown_unicast = "proxy"
		vrf_name               = mso_schema_template_vrf.%[5]s.name
	}

	resource "mso_schema_site_bd" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = mso_schema_site.%[6]s.site_id
		template_name = "%[4]s"
		bd_name       = mso_schema_template_bd.%[2]s_renamed.name
		host_route    = false
	}`,
		testAccMSOSchemaSiteBdPrerequisiteConfig(),
		msoSchemaTemplateBdName,
		msoSchemaName,
		msoSchemaTemplateName,
		msoSchemaTemplateVrfName,
		msoSchemaSiteResourceLabel1,
	)
}

func testAccMSOSchemaSiteBdConfigNoSiteAssociation() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_bd" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = data.mso_site.%[4]s.id
		template_name = "%[5]s"
		bd_name       = mso_schema_template_bd.%[2]s.name
	}`,
		testAccMSOSchemaSiteBdNoSiteAssociationConfig(),
		msoSchemaTemplateBdName,
		msoSchemaName,
		msoTemplateSiteName1,
		msoSchemaTemplateName,
	)
}

// testAccMSOSchemaSiteBdL3OutPrerequisiteConfig is the standard site BD
// prereq with the additional mso_schema_template_l3out so that a
// mso_schema_site_bd_l3out can reference it.
func testAccMSOSchemaSiteBdL3OutPrerequisiteConfig() string {
	return fmt.Sprintf(`%s%s`,
		testAccMSOSchemaSiteBdPrerequisiteConfig(),
		testSchemaTemplateL3outConfig(),
	)
}

// testAccMSOSchemaSiteBdConfigWithL3Out keeps the BD at host_route=true,
// svi_mac=:03 (matching the prior step) and attaches an
// mso_schema_site_bd_l3out child resource referencing the template L3Out.
func testAccMSOSchemaSiteBdConfigWithL3Out() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_bd" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = mso_schema_site.%[4]s.site_id
		template_name = "%[5]s"
		bd_name       = mso_schema_template_bd.%[2]s.name
		host_route    = true
		svi_mac       = "00:00:5E:00:01:03"
	}

	resource "mso_schema_site_bd_l3out" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = mso_schema_site.%[4]s.site_id
		template_name = "%[5]s"
		bd_name       = mso_schema_site_bd.%[2]s.bd_name
		l3out_name    = mso_schema_template_l3out.%[6]s.l3out_name
	}`,
		testAccMSOSchemaSiteBdL3OutPrerequisiteConfig(),
		msoSchemaTemplateBdName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
		msoSchemaTemplateL3outName,
	)
}

// testAccMSOSchemaSiteBdConfigWithL3OutBdUpdated is identical to
// testAccMSOSchemaSiteBdConfigWithL3Out except svi_mac is bumped to
// "00:00:5E:00:01:04". Applying it triggers the BD Update path while the
// L3Out child is still attached -- the regression check in the test step
// verifies the BD's l3Outs array is left intact by that PATCH.
func testAccMSOSchemaSiteBdConfigWithL3OutBdUpdated() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_bd" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = mso_schema_site.%[4]s.site_id
		template_name = "%[5]s"
		bd_name       = mso_schema_template_bd.%[2]s.name
		host_route    = true
		svi_mac       = "00:00:5E:00:01:04"
	}

	resource "mso_schema_site_bd_l3out" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = mso_schema_site.%[4]s.site_id
		template_name = "%[5]s"
		bd_name       = mso_schema_site_bd.%[2]s.bd_name
		l3out_name    = mso_schema_template_l3out.%[6]s.l3out_name
	}`,
		testAccMSOSchemaSiteBdL3OutPrerequisiteConfig(),
		msoSchemaTemplateBdName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
		msoSchemaTemplateL3outName,
	)
}

// testAccCheckMSOSchemaSiteBdHasL3Out walks the schema on MSO/NDO and
// asserts that the BD identified by the named mso_schema_site_bd resource
// in state has wantL3Out present in its sites[].bds[].l3Outs[] array. Used
// to detect regressions where a BD update would clobber the BD's L3Out
// list on the API side without leaving any trace in Terraform state.
func testAccCheckMSOSchemaSiteBdHasL3Out(siteBdResourceName, wantL3Out string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[siteBdResourceName]
		if !ok {
			return fmt.Errorf("site BD resource %s not found in state", siteBdResourceName)
		}
		schemaId := rs.Primary.Attributes["schema_id"]
		stateSiteId := rs.Primary.Attributes["site_id"]
		stateTemplate := rs.Primary.Attributes["template_name"]
		stateBd := rs.Primary.Attributes["bd_name"]

		msoClient := testAccProvider.Meta().(*client.Client)
		cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
		if err != nil {
			return fmt.Errorf("unable to fetch schema %s: %v", schemaId, err)
		}
		siteCount, err := cont.ArrayCount("sites")
		if err != nil {
			return fmt.Errorf("schema %s has no sites array", schemaId)
		}
		bdRefRe := regexp.MustCompile("/schemas/(.*)/templates/(.*)/bds/(.*)")
		for i := 0; i < siteCount; i++ {
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
			bdCount, err := siteCont.ArrayCount("bds")
			if err != nil {
				return fmt.Errorf("site %s/%s has no bds array", stateSiteId, stateTemplate)
			}
			for j := 0; j < bdCount; j++ {
				bdCont, err := siteCont.ArrayElement(j, "bds")
				if err != nil {
					return err
				}
				match := bdRefRe.FindStringSubmatch(models.StripQuotes(bdCont.S("bdRef").String()))
				if len(match) != 4 || match[3] != stateBd {
					continue
				}
				l3outCount, err := bdCont.ArrayCount("l3Outs")
				if err != nil {
					return fmt.Errorf("site BD %s has no l3Outs array (expected %s)", stateBd, wantL3Out)
				}
				for k := 0; k < l3outCount; k++ {
					l3outCont, err := bdCont.ArrayElement(k, "l3Outs")
					if err != nil {
						return err
					}
					if strings.Trim(l3outCont.String(), "\"") == wantL3Out {
						return nil
					}
				}
				return fmt.Errorf("L3Out %s not found in l3Outs of site BD %s on schema %s", wantL3Out, stateBd, schemaId)
			}
			return fmt.Errorf("site BD %s not found under site %s/%s on schema %s", stateBd, stateSiteId, stateTemplate, schemaId)
		}
		return fmt.Errorf("site %s/%s not found on schema %s", stateSiteId, stateTemplate, schemaId)
	}
}

// testAccCheckMSOSchemaSiteBdDestroy walks state for any mso_schema_site_bd
// resources, fetches the schema, and asserts that no sites[].bds[] entry
// whose bdRef resolves to bd_name remains under the matching siteId +
// templateName. A missing schema or missing sites array is treated as a
// successful destroy.
func testAccCheckMSOSchemaSiteBdDestroy(s *terraform.State) error {
	msoClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mso_schema_site_bd" {
			continue
		}
		schemaId := rs.Primary.Attributes["schema_id"]
		stateSiteId := rs.Primary.Attributes["site_id"]
		stateTemplate := rs.Primary.Attributes["template_name"]
		stateBd := rs.Primary.Attributes["bd_name"]

		cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
		if err != nil {
			// Schema itself has been destroyed.
			return nil
		}
		count, err := cont.ArrayCount("sites")
		if err != nil {
			return nil
		}
		re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/bds/(.*)")
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
			bdCount, err := siteCont.ArrayCount("bds")
			if err != nil {
				continue
			}
			for j := 0; j < bdCount; j++ {
				bdCont, err := siteCont.ArrayElement(j, "bds")
				if err != nil {
					return err
				}
				bdRef := models.StripQuotes(bdCont.S("bdRef").String())
				match := re.FindStringSubmatch(bdRef)
				if len(match) == 4 && match[3] == stateBd {
					return fmt.Errorf("mso_schema_site_bd (site=%s, template=%s, bd=%s) still exists on schema %s", stateSiteId, stateTemplate, stateBd, schemaId)
				}
			}
		}
	}
	return nil
}
