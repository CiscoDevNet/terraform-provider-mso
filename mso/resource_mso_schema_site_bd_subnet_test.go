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

// TestAccMSOSchemaSiteBdSubnetResource exercises the lifecycle of
// mso_schema_site_bd_subnet:
//   - attempt to create without a mso_schema_site association (expect error)
//   - create the subnet with defaults and verify all attributes
//   - update scope, shared, description, and no_default_gateway in-place
//   - update ip (ForceNew — triggers destroy + recreate with the new CIDR)
//   - import the subnet
//
// The negative-path step omits the mso_schema_site association so that the
// PATCH target path does not exist in the schema, guaranteeing a rejection on
// all NDO versions. Without a site association, adding a subnet without a
// mso_schema_site_bd on newer NDO is silently accepted (the API either creates
// the BD entry implicitly or drops the error), so the "no site BD" condition
// alone is not a reliable error trigger.
//
// The lab must have the `ansible_test` and `ansible_test_2` sites onboarded.
func TestAccMSOSchemaSiteBdSubnetResource(t *testing.T) {
	siteBdSubnetResource := "mso_schema_site_bd_subnet." + msoSchemaTemplateBdName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteBdSubnetDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: Create site BD subnet without mso_schema_site association (expect error)")
				},
				Config: testAccMSOSchemaSiteBdSubnetConfigNoSiteAssociation(),
				// Older NDO rejects the PATCH with "Resource Not Found". Newer
				// NDO's always-on schema validation engine silently drops the
				// PATCH so Create succeeds, the follow-up Read finds nothing,
				// and the SDK raises "Provider produced inconsistent result
				// after apply". Match either outcome.
				ExpectError: regexp.MustCompile(`Resource Not Found|Provider produced inconsistent result after apply`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create site BD subnet with defaults") },
				Config:    testAccMSOSchemaSiteBdSubnetConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(siteBdSubnetResource, "schema_id"),
					resource.TestCheckResourceAttrSet(siteBdSubnetResource, "site_id"),
					resource.TestCheckResourceAttr(siteBdSubnetResource, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr(siteBdSubnetResource, "bd_name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr(siteBdSubnetResource, "ip", msoSchemaSiteBdSubnetIp),
					resource.TestCheckResourceAttrPair(
						siteBdSubnetResource, "site_id",
						"data.mso_site."+msoTemplateSiteName1, "id",
					),
					resource.TestCheckResourceAttr(siteBdSubnetResource, "id", msoSchemaSiteBdSubnetIp),
					resource.TestCheckResourceAttr(siteBdSubnetResource, "scope", "private"),
					resource.TestCheckResourceAttr(siteBdSubnetResource, "shared", "false"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update scope, shared, description, and no_default_gateway")
				},
				Config: testAccMSOSchemaSiteBdSubnetConfigUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(siteBdSubnetResource, "ip", msoSchemaSiteBdSubnetIp),
					resource.TestCheckResourceAttr(siteBdSubnetResource, "scope", "public"),
					resource.TestCheckResourceAttr(siteBdSubnetResource, "shared", "true"),
					resource.TestCheckResourceAttr(siteBdSubnetResource, "description", "updated subnet"),
					resource.TestCheckResourceAttr(siteBdSubnetResource, "no_default_gateway", "true"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Remove description (verifies empty string clears the field)")
				},
				Config: testAccMSOSchemaSiteBdSubnetConfigRemoveDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(siteBdSubnetResource, "ip", msoSchemaSiteBdSubnetIp),
					resource.TestCheckResourceAttr(siteBdSubnetResource, "description", ""),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update ip (ForceNew, destroy + recreate with new CIDR)")
				},
				Config: testAccMSOSchemaSiteBdSubnetConfigUpdateIp(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(siteBdSubnetResource, "ip", msoSchemaSiteBdSubnetIp2),
					resource.TestCheckResourceAttr(siteBdSubnetResource, "id", msoSchemaSiteBdSubnetIp2),
					resource.TestCheckResourceAttr(siteBdSubnetResource, "scope", "public"),
					resource.TestCheckResourceAttr(siteBdSubnetResource, "shared", "true"),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import site BD subnet") },
				ResourceName: siteBdSubnetResource,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[siteBdSubnetResource]
					if !ok {
						return "", fmt.Errorf("site BD subnet resource not found in state: %s", siteBdSubnetResource)
					}
					// Importer splits by "/" and uses a regex "(.*)/ip/(.*)"
					// to extract the IP (which itself contains "/").
					return fmt.Sprintf("%s/site/%s/template/%s/bd/%s/ip/%s",
						rs.Primary.Attributes["schema_id"],
						rs.Primary.Attributes["site_id"],
						rs.Primary.Attributes["template_name"],
						rs.Primary.Attributes["bd_name"],
						rs.Primary.Attributes["ip"],
					), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

// testAccMSOSchemaSiteBdSubnetPrerequisiteConfig extends the standard site BD
// prereqs with an mso_schema_site_bd, providing the minimum config
// required by mso_schema_site_bd_subnet.
func testAccMSOSchemaSiteBdSubnetPrerequisiteConfig() string {
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

func testAccMSOSchemaSiteBdSubnetConfigCreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_bd_subnet" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = mso_schema_site.%[4]s.site_id
		template_name = "%[5]s"
		bd_name       = mso_schema_site_bd.%[2]s.bd_name
		ip            = "%[6]s"
		shared        = false
	}`,
		testAccMSOSchemaSiteBdSubnetPrerequisiteConfig(),
		msoSchemaTemplateBdName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
		msoSchemaSiteBdSubnetIp,
	)
}

func testAccMSOSchemaSiteBdSubnetConfigUpdate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_bd_subnet" "%[2]s" {
		schema_id          = mso_schema.%[3]s.id
		site_id            = mso_schema_site.%[4]s.site_id
		template_name      = "%[5]s"
		bd_name            = mso_schema_site_bd.%[2]s.bd_name
		ip                 = "%[6]s"
		scope              = "public"
		shared             = true
		description        = "updated subnet"
		no_default_gateway = true
	}`,
		testAccMSOSchemaSiteBdSubnetPrerequisiteConfig(),
		msoSchemaTemplateBdName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
		msoSchemaSiteBdSubnetIp,
	)
}

// testAccMSOSchemaSiteBdSubnetConfigRemoveDescription omits the description
// field to verify that the resource can clear a previously-set description.
func testAccMSOSchemaSiteBdSubnetConfigRemoveDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_bd_subnet" "%[2]s" {
		schema_id          = mso_schema.%[3]s.id
		site_id            = mso_schema_site.%[4]s.site_id
		template_name      = "%[5]s"
		bd_name            = mso_schema_site_bd.%[2]s.bd_name
		ip                 = "%[6]s"
		scope              = "public"
		shared             = true
		no_default_gateway = true
	}`,
		testAccMSOSchemaSiteBdSubnetPrerequisiteConfig(),
		msoSchemaTemplateBdName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
		msoSchemaSiteBdSubnetIp,
	)
}

// testAccMSOSchemaSiteBdSubnetConfigUpdateIp changes the ip (ForceNew), which
// causes Terraform to destroy the existing subnet and recreate it with the new
// CIDR.
func testAccMSOSchemaSiteBdSubnetConfigUpdateIp() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_bd_subnet" "%[2]s" {
		schema_id          = mso_schema.%[3]s.id
		site_id            = mso_schema_site.%[4]s.site_id
		template_name      = "%[5]s"
		bd_name            = mso_schema_site_bd.%[2]s.bd_name
		ip                 = "%[6]s"
		scope              = "public"
		shared             = true
		description        = "updated subnet"
		no_default_gateway = true
	}`,
		testAccMSOSchemaSiteBdSubnetPrerequisiteConfig(),
		msoSchemaTemplateBdName,
		msoSchemaName,
		msoSchemaSiteResourceLabel1,
		msoSchemaTemplateName,
		msoSchemaSiteBdSubnetIp2,
	)
}

// testAccMSOSchemaSiteBdSubnetConfigNoSiteAssociation creates a site BD Subnet
// without a prior mso_schema_site association, exercising the negative path.
func testAccMSOSchemaSiteBdSubnetConfigNoSiteAssociation() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_site_bd_subnet" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		site_id       = data.mso_site.%[4]s.id
		template_name = "%[5]s"
		bd_name       = mso_schema_template_bd.%[2]s.name
		ip            = "%[6]s"
		scope         = "private"
		shared        = false
	}`,
		fmt.Sprintf(`%s%s%s`,
			testSchemaWithBothSitesPrerequisiteConfig(),
			testSchemaTemplateVrfConfig(),
			testSchemaTemplateBdConfig(),
		),
		msoSchemaTemplateBdName,
		msoSchemaName,
		msoTemplateSiteName1,
		msoSchemaTemplateName,
		msoSchemaSiteBdSubnetIp,
	)
}

// testAccCheckMSOSchemaSiteBdSubnetDestroy walks state for any
// mso_schema_site_bd_subnet resources and asserts that no
// sites[].bds[].subnets[] entry matching ip remains under the matching
// siteId + templateName + bd_name. A missing schema or missing sites array
// is treated as a successful destroy.
func testAccCheckMSOSchemaSiteBdSubnetDestroy(s *terraform.State) error {
	msoClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mso_schema_site_bd_subnet" {
			continue
		}
		schemaId := rs.Primary.Attributes["schema_id"]
		stateSiteId := rs.Primary.Attributes["site_id"]
		stateTemplate := rs.Primary.Attributes["template_name"]
		stateBd := rs.Primary.Attributes["bd_name"]
		stateIp := rs.Primary.Attributes["ip"]

		cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
		if err != nil {
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
			if models.StripQuotes(siteCont.S("siteId").String()) != stateSiteId {
				continue
			}
			if models.StripQuotes(siteCont.S("templateName").String()) != stateTemplate {
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
				match := re.FindStringSubmatch(models.StripQuotes(bdCont.S("bdRef").String()))
				if len(match) != 4 || match[3] != stateBd {
					continue
				}
				subnetCount, err := bdCont.ArrayCount("subnets")
				if err != nil {
					continue
				}
				for l := 0; l < subnetCount; l++ {
					subnetCont, err := bdCont.ArrayElement(l, "subnets")
					if err != nil {
						return err
					}
					if strings.TrimSpace(models.StripQuotes(subnetCont.S("ip").String())) == stateIp {
						return fmt.Errorf("mso_schema_site_bd_subnet (site=%s, template=%s, bd=%s, ip=%s) still exists on schema %s",
							stateSiteId, stateTemplate, stateBd, stateIp, schemaId)
					}
				}
			}
		}
	}
	return nil
}
