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

// TestAccMSOSchemaSiteBdL3outResource exercises the lifecycle of
// mso_schema_site_bd_l3out:
//   - attempt to create without a mso_schema_site association (expect error)
//   - create the site BD L3out and verify all attributes
//   - import the site BD L3out
//
// The lab must have the `ansible_test` and `ansible_test_2` sites onboarded.
func TestAccMSOSchemaSiteBdL3outResource(t *testing.T) {
	siteBdL3outResource := "mso_schema_site_bd_l3out." + msoSchemaTemplateBdName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteBdL3outDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: Create site BD L3out without mso_schema_site association (expect error)")
				},
				Config: testAccMSOSchemaSiteBdL3outConfigNoSiteAssociation(),
				// Older NDO rejects the PATCH with "Resource Not Found". Newer
				// NDO's always-on schema validation engine silently drops the
				// PATCH as a no-op, so Create succeeds, the follow-up Read
				// finds nothing, and the SDK raises "Provider produced
				// inconsistent result after apply". Match either outcome.
				ExpectError: regexp.MustCompile(`Resource Not Found|Provider produced inconsistent result after apply`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create site BD L3out") },
				Config:    testAccMSOSchemaSiteBdL3outConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(siteBdL3outResource, "schema_id"),
					resource.TestCheckResourceAttrSet(siteBdL3outResource, "site_id"),
					resource.TestCheckResourceAttr(siteBdL3outResource, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr(siteBdL3outResource, "bd_name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr(siteBdL3outResource, "l3out_name", msoSchemaTemplateL3outName),
					resource.TestCheckResourceAttrPair(
						siteBdL3outResource, "site_id",
						"data.mso_site."+msoTemplateSiteName1, "id",
					),
					resource.TestCheckResourceAttr(siteBdL3outResource, "id", msoSchemaTemplateL3outName),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import site BD L3out") },
				ResourceName: siteBdL3outResource,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[siteBdL3outResource]
					if !ok {
						return "", fmt.Errorf("site BD L3out resource not found in state: %s", siteBdL3outResource)
					}
					// Importer splits on "/" and reads indices [0] (schema_id),
					// [2] (site_id), [4] (bd_name), [6] (l3out_name).
					return fmt.Sprintf("%s/site/%s/bd/%s/l3out/%s",
						rs.Primary.Attributes["schema_id"],
						rs.Primary.Attributes["site_id"],
						rs.Primary.Attributes["bd_name"],
						rs.Primary.Attributes["l3out_name"],
					), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

// testAccMSOSchemaSiteBdL3outPrerequisiteConfig extends the standard site BD
// L3out template prereqs with an mso_schema_site_bd resource, providing the
// minimum scaffolding required by the mso_schema_site_bd_l3out resource.
func testAccMSOSchemaSiteBdL3outPrerequisiteConfig() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_bd" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = mso_schema_site.%[4]s.site_id
		template_name = "%[5]s"
		bd_name       = mso_schema_template_bd.%[2]s.name
	}`,
		testAccMSOSchemaSiteBdL3OutPrerequisiteConfig(),
		msoSchemaTemplateBdName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
	)
}

func testAccMSOSchemaSiteBdL3outConfigCreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_bd_l3out" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = mso_schema_site.%[4]s.site_id
		template_name = "%[5]s"
		bd_name       = mso_schema_site_bd.%[2]s.bd_name
		l3out_name    = mso_schema_template_l3out.%[6]s.l3out_name
	}`,
		testAccMSOSchemaSiteBdL3outPrerequisiteConfig(),
		msoSchemaTemplateBdName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
		msoSchemaTemplateL3outName,
	)
}

// testAccMSOSchemaSiteBdL3outConfigNoSiteAssociation creates a site BD L3out
// without a prior mso_schema_site association, exercising the negative path.
func testAccMSOSchemaSiteBdL3outConfigNoSiteAssociation() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_bd_l3out" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = data.mso_site.%[4]s.id
		template_name = "%[5]s"
		bd_name       = mso_schema_template_bd.%[2]s.name
		l3out_name    = mso_schema_template_l3out.%[6]s.l3out_name
	}`,
		fmt.Sprintf(`%s%s%s%s`,
			testSchemaWithBothSitesPrerequisiteConfig(),
			testSchemaTemplateVrfConfig(),
			testSchemaTemplateBdConfig(),
			testSchemaTemplateL3outConfig(),
		),
		msoSchemaTemplateBdName,
		msoSchemaName,
		msoTemplateSiteName1,
		msoSchemaTemplateName,
		msoSchemaTemplateL3outName,
	)
}

// testAccCheckMSOSchemaSiteBdL3outDestroy walks state for any
// mso_schema_site_bd_l3out resources and asserts that no
// sites[].bds[].l3Outs[] entry matching l3out_name remains under the matching
// siteId + templateName + bd_name. A missing schema or missing sites array is
// treated as a successful destroy.
func testAccCheckMSOSchemaSiteBdL3outDestroy(s *terraform.State) error {
	msoClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mso_schema_site_bd_l3out" {
			continue
		}
		schemaId := rs.Primary.Attributes["schema_id"]
		stateSiteId := rs.Primary.Attributes["site_id"]
		stateTemplate := rs.Primary.Attributes["template_name"]
		stateBd := rs.Primary.Attributes["bd_name"]
		stateL3out := rs.Primary.Attributes["l3out_name"]

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
			bdCount, err := siteCont.ArrayCount("bds")
			if err != nil {
				continue
			}
			for j := 0; j < bdCount; j++ {
				bdCont, err := siteCont.ArrayElement(j, "bds")
				if err != nil {
					return err
				}
				apiBdRef := models.StripQuotes(bdCont.S("bdRef").String())
				split := strings.Split(apiBdRef, "/")
				if len(split) < 7 || split[6] != stateBd {
					continue
				}
				l3outCount, err := bdCont.ArrayCount("l3Outs")
				if err != nil {
					continue
				}
				for k := 0; k < l3outCount; k++ {
					l3outCont, err := bdCont.ArrayElement(k, "l3Outs")
					if err != nil {
						return err
					}
					if strings.Trim(l3outCont.String(), "\"") == stateL3out {
						return fmt.Errorf("mso_schema_site_bd_l3out (site=%s, template=%s, bd=%s, l3out=%s) still exists on schema %s",
							stateSiteId, stateTemplate, stateBd, stateL3out, schemaId)
					}
				}
			}
		}
	}
	return nil
}
