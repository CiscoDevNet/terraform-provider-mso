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

// TestAccMSOSchemaSiteAnpEpgSubnetResource exercises the lifecycle of
// mso_schema_site_anp_epg_subnet:
//   - attempt to create a subnet on a site that has no mso_schema_site
//     association (expect error)
//   - create the subnet with all mutable attributes set
//   - update the subnet including clearing description back to ""
//   - import the subnet
//
// primary is intentionally skipped: fabric-local EPGs do
// not allow subnets to be marked as primary — NDO rejects the PATCH with
// "Fabric local EPG Subnet cannot be marked as primary".
//
// querier is intentionally skipped: it is only supported for Bridge Domain
// subnets — NDO rejects the PATCH with "'Querier' is only supported for
// Bridge Domain subnets". This attribute is a deprecation candidate.
//
// The lab must have the `ansible_test` and `ansible_test_2` sites onboarded.
func TestAccMSOSchemaSiteAnpEpgSubnetResource(t *testing.T) {
	subnetResource := "mso_schema_site_anp_epg_subnet." + msoSchemaTemplateAnpEpgName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpEpgSubnetDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: Create subnet without mso_schema_site association (expect error)")
				},
				Config: testAccMSOSchemaSiteAnpEpgSubnetConfigNoSiteAssociation(),
				// Older NDO rejects the PATCH with "Resource Not Found". Newer
				// NDO silently drops it so the follow-up Read finds nothing and
				// the SDK raises "Provider produced inconsistent result after
				// apply". Match either outcome.
				ExpectError: regexp.MustCompile(`Resource Not Found|Provider produced inconsistent result after apply`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create subnet") },
				Config:    testAccMSOSchemaSiteAnpEpgSubnetConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(subnetResource, "schema_id"),
					resource.TestCheckResourceAttrSet(subnetResource, "site_id"),
					resource.TestCheckResourceAttr(subnetResource, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr(subnetResource, "anp_name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr(subnetResource, "epg_name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr(subnetResource, "ip", msoSchemaSiteAnpEpgSubnetIp),
					resource.TestCheckResourceAttr(subnetResource, "scope", "private"),
					resource.TestCheckResourceAttr(subnetResource, "shared", "false"),
					resource.TestCheckResourceAttr(subnetResource, "no_default_gateway", "false"),
					resource.TestCheckResourceAttr(subnetResource, "description", "test description"),
					resource.TestCheckResourceAttrPair(
						subnetResource, "site_id",
						"data.mso_site."+msoTemplateSiteName1, "id",
					),
					resource.TestCheckResourceAttr(subnetResource, "id", msoSchemaSiteAnpEpgSubnetIp),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update subnet (attributes changed, description cleared)")
				},
				Config: testAccMSOSchemaSiteAnpEpgSubnetConfigUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(subnetResource, "ip", msoSchemaSiteAnpEpgSubnetIp),
					resource.TestCheckResourceAttr(subnetResource, "scope", "public"),
					resource.TestCheckResourceAttr(subnetResource, "shared", "true"),
					resource.TestCheckResourceAttr(subnetResource, "no_default_gateway", "true"),
					resource.TestCheckResourceAttr(subnetResource, "description", ""),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import subnet") },
				ResourceName: subnetResource,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[subnetResource]
					if !ok {
						return "", fmt.Errorf("subnet resource not found in state: %s", subnetResource)
					}
					// The importer uses a "(.*)/ip/(.*)" regex to handle the
					// "/" characters in the CIDR notation of the IP.
					return fmt.Sprintf("%s/site/%s/template/%s/anp/%s/epg/%s/ip/%s",
						rs.Primary.Attributes["schema_id"],
						rs.Primary.Attributes["site_id"],
						rs.Primary.Attributes["template_name"],
						rs.Primary.Attributes["anp_name"],
						rs.Primary.Attributes["epg_name"],
						rs.Primary.Attributes["ip"],
					), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMSOSchemaSiteAnpEpgSubnetConfigCreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_anp_epg_subnet" "%[2]s" {
		schema_id           = mso_schema.%[3]s.id
		site_id             = mso_schema_site.%[4]s.site_id
		template_name       = "%[5]s"
		anp_name            = mso_schema_template_anp.%[6]s.name
		epg_name            = mso_schema_site_anp_epg.%[2]s.epg_name
		ip                  = "%[7]s"
		scope               = "private"
		shared              = false
		no_default_gateway  = false
		description         = "test description"
	}`,
		testAccMSOSchemaSiteAnpEpgStaticLeafPrerequisiteConfig(),
		msoSchemaTemplateAnpEpgName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
		msoSchemaTemplateAnpName,
		msoSchemaSiteAnpEpgSubnetIp,
	)
}

func testAccMSOSchemaSiteAnpEpgSubnetConfigUpdate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_anp_epg_subnet" "%[2]s" {
		schema_id           = mso_schema.%[3]s.id
		site_id             = mso_schema_site.%[4]s.site_id
		template_name       = "%[5]s"
		anp_name            = mso_schema_template_anp.%[6]s.name
		epg_name            = mso_schema_site_anp_epg.%[2]s.epg_name
		ip                  = "%[7]s"
		scope               = "public"
		shared              = true
		no_default_gateway  = true
		description         = ""
	}`,
		testAccMSOSchemaSiteAnpEpgStaticLeafPrerequisiteConfig(),
		msoSchemaTemplateAnpEpgName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
		msoSchemaTemplateAnpName,
		msoSchemaSiteAnpEpgSubnetIp,
	)
}

// testAccMSOSchemaSiteAnpEpgSubnetConfigNoSiteAssociation creates a subnet
// without a prior mso_schema_site association, exercising the negative path.
func testAccMSOSchemaSiteAnpEpgSubnetConfigNoSiteAssociation() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_anp_epg_subnet" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = data.mso_site.%[4]s.id
		template_name = "%[5]s"
		anp_name      = mso_schema_template_anp.%[6]s.name
		epg_name      = mso_schema_template_anp_epg.%[2]s.name
		ip            = "%[7]s"
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
		msoSchemaSiteAnpEpgSubnetIp,
	)
}

// testAccCheckMSOSchemaSiteAnpEpgSubnetDestroy walks state for any
// mso_schema_site_anp_epg_subnet resources, fetches the schema, and asserts
// that no sites[].anps[].epgs[].subnets[] entry with the matching IP remains.
// A missing schema or missing sites array is treated as a successful destroy.
func testAccCheckMSOSchemaSiteAnpEpgSubnetDestroy(s *terraform.State) error {
	msoClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mso_schema_site_anp_epg_subnet" {
			continue
		}
		schemaId := rs.Primary.Attributes["schema_id"]
		stateSiteId := rs.Primary.Attributes["site_id"]
		stateTemplate := rs.Primary.Attributes["template_name"]
		stateAnp := rs.Primary.Attributes["anp_name"]
		stateEpg := rs.Primary.Attributes["epg_name"]
		stateIp := rs.Primary.Attributes["ip"]

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
					subnetCount, err := epgCont.ArrayCount("subnets")
					if err != nil {
						continue
					}
					for l := 0; l < subnetCount; l++ {
						subnetCont, err := epgCont.ArrayElement(l, "subnets")
						if err != nil {
							return err
						}
						if models.StripQuotes(subnetCont.S("ip").String()) == stateIp {
							return fmt.Errorf(
								"mso_schema_site_anp_epg_subnet (site=%s, template=%s, anp=%s, epg=%s, ip=%s) still exists on schema %s",
								stateSiteId, stateTemplate, stateAnp, stateEpg, stateIp, schemaId,
							)
						}
					}
				}
			}
		}
	}
	return nil
}
