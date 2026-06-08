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

// TestAccMSOSchemaSiteResource exercises the lifecycle of mso_schema_site:
//   - add a single site to a schema template
//   - add a second site to the same template
//   - remove the second site while the template is NOT deployed
//   - re-add the second site and deploy the template via mso_schema_template_deploy_ndo
//   - attempt to remove the deployed site with undeploy_on_destroy=false and
//     expect MSO to reject the destroy with a "must first be undeployed" error
//   - flip undeploy_on_destroy=true and remove the second site while the
//     template IS deployed (exercising the undeploy branch of
//     resourceMSOSchemaSiteDelete)
//   - import the remaining site
//
// The lab must have both `ansible_test` and `ansible_test_2` sites onboarded.
func TestAccMSOSchemaSiteResource(t *testing.T) {
	site1Resource := "mso_schema_site." + msoSchemaSiteResourceLabel1
	site2Resource := "mso_schema_site." + msoSchemaSiteResourceLabel2

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Add site_1 to schema template (not deployed)") },
				Config:    testSchemaWithSingleSiteAssociationConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(site1Resource, "schema_id"),
					resource.TestCheckResourceAttrSet(site1Resource, "site_id"),
					resource.TestCheckResourceAttr(site1Resource, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr(site1Resource, "undeploy_on_destroy", "false"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add site_2 to the same schema template (not deployed)") },
				Config:    testAccMSOSchemaSiteConfigTwoSites(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(site1Resource, "site_id"),
					resource.TestCheckResourceAttr(site1Resource, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttrSet(site2Resource, "site_id"),
					resource.TestCheckResourceAttr(site2Resource, "template_name", msoSchemaTemplateName),
					testAccCheckMSOSchemaSiteCount(site1Resource, 2),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove site_2 while template is NOT deployed") },
				Config:    testSchemaWithSingleSiteAssociationConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(site1Resource, "site_id"),
					resource.TestCheckResourceAttr(site1Resource, "template_name", msoSchemaTemplateName),
					testAccCheckMSOSchemaSiteCount(site1Resource, 1),
					testAccCheckMSOSchemaSiteAbsent(site1Resource, msoTemplateSiteName2),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Re-add site_2 with undeploy_on_destroy=false, add VRF and deploy template")
				},
				Config: testAccMSOSchemaSiteConfigTwoSitesDeployedSite2Undeploy(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(site1Resource, "site_id"),
					resource.TestCheckResourceAttrSet(site2Resource, "site_id"),
					resource.TestCheckResourceAttr(site2Resource, "undeploy_on_destroy", "false"),
					resource.TestCheckResourceAttrSet("mso_schema_template_deploy_ndo.deploy", "schema_id"),
					testAccCheckMSOSchemaSiteCount(site1Resource, 2),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Remove deployed site_2 with undeploy_on_destroy=false (expect error)")
				},
				Config:      testSchemaWithSingleSiteAssociationDeployedConfig(),
				ExpectError: regexp.MustCompile(`cannot be deleted; template .* must first be undeployed`),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Flip site_2 undeploy_on_destroy=true to allow undeploy on destroy")
				},
				Config: testAccMSOSchemaSiteConfigTwoSitesDeployedSite2Undeploy(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(site2Resource, "undeploy_on_destroy", "true"),
					resource.TestCheckResourceAttrSet("mso_schema_template_deploy_ndo.deploy", "schema_id"),
					testAccCheckMSOSchemaSiteCount(site1Resource, 2),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove site_2 while template IS deployed (undeploy_on_destroy=true)") },
				Config:    testSchemaWithSingleSiteAssociationDeployedConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(site1Resource, "site_id"),
					resource.TestCheckResourceAttr(site1Resource, "undeploy_on_destroy", "true"),
					resource.TestCheckResourceAttrSet("mso_schema_template_deploy_ndo.deploy", "schema_id"),
					testAccCheckMSOSchemaSiteCount(site1Resource, 1),
					testAccCheckMSOSchemaSiteAbsent(site1Resource, msoTemplateSiteName2),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import remaining schema_site (site_1)") },
				ResourceName: site1Resource,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[site1Resource]
					if !ok {
						return "", fmt.Errorf("schema_site resource not found in state: %s", site1Resource)
					}
					return fmt.Sprintf("%s/sites/%s/templates/%s/undeploy_on_destroy/%s",
						rs.Primary.Attributes["schema_id"],
						msoTemplateSiteName1,
						rs.Primary.Attributes["template_name"],
						rs.Primary.Attributes["undeploy_on_destroy"],
					), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

// testAccMSOSchemaSiteConfigTwoSites adds site_2 to the single-site config to
// exercise adding a second site association to the same template (no deploy).
func testAccMSOSchemaSiteConfigTwoSites() string {
	return fmt.Sprintf(`%s%s`,
		testSchemaWithSingleSiteAssociationConfig(),
		testSchemaSiteConfig(msoSchemaSiteResourceLabel2, msoTemplateSiteName2, false),
	)
}

// testAccMSOSchemaSiteConfigTwoSitesDeployedSite2Undeploy emits the deployed
// two-site configuration where site_1 has `undeploy_on_destroy=true` and the
// caller controls the value on site_2. This lets the lifecycle test exercise
// both the rejected-destroy (false) and the undeploy-on-destroy (true) paths.
func testAccMSOSchemaSiteConfigTwoSitesDeployedSite2Undeploy(site2UndeployOnDestroy bool) string {
	return fmt.Sprintf(`%s%s%s%s%s`,
		testSchemaWithBothSitesPrerequisiteConfig(),
		testSchemaSiteConfig(msoSchemaSiteResourceLabel1, msoTemplateSiteName1, true),
		testSchemaSiteConfig(msoSchemaSiteResourceLabel2, msoTemplateSiteName2, site2UndeployOnDestroy),
		testSchemaTemplateVrfConfig(),
		testSchemaTemplateDeployNdoConfig([]string{
			"mso_schema_site." + msoSchemaSiteResourceLabel1,
			"mso_schema_site." + msoSchemaSiteResourceLabel2,
			"mso_schema_template_vrf." + msoSchemaTemplateVrfName,
		}),
	)
}

// testAccCheckMSOSchemaSiteCount asserts the number of `sites` entries on the
// schema referenced by the supplied resource name matches the expected count.
func testAccCheckMSOSchemaSiteCount(resourceName string, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", resourceName)
		}
		schemaId := rs.Primary.Attributes["schema_id"]
		msoClient := testAccProvider.Meta().(*client.Client)
		cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
		if err != nil {
			return fmt.Errorf("failed to GET schema %s: %v", schemaId, err)
		}
		count, err := cont.ArrayCount("sites")
		if err != nil {
			if expected == 0 {
				return nil
			}
			return fmt.Errorf("expected %d sites on schema %s but `sites` array missing: %v", expected, schemaId, err)
		}
		if count != expected {
			return fmt.Errorf("expected %d sites on schema %s, got %d", expected, schemaId, count)
		}
		return nil
	}
}

// testAccCheckMSOSchemaSiteAbsent asserts that the schema does not contain a
// site association whose `siteId` resolves to the supplied site name.
func testAccCheckMSOSchemaSiteAbsent(resourceName, siteName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", resourceName)
		}
		schemaId := rs.Primary.Attributes["schema_id"]
		msoClient := testAccProvider.Meta().(*client.Client)

		// Resolve siteName -> siteId.
		sitesCont, err := msoClient.GetViaURL("api/v1/sites")
		if err != nil {
			return fmt.Errorf("failed to GET sites: %v", err)
		}
		var targetSiteId string
		sitesData, ok := sitesCont.S("sites").Data().([]interface{})
		if !ok {
			return fmt.Errorf("unexpected payload for api/v1/sites")
		}
		for _, raw := range sitesData {
			entry, ok := raw.(map[string]interface{})
			if !ok {
				continue
			}
			if name, _ := entry["name"].(string); name == siteName {
				targetSiteId, _ = entry["id"].(string)
				break
			}
		}
		if targetSiteId == "" {
			return fmt.Errorf("site %q not found on MSO", siteName)
		}

		schemaCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
		if err != nil {
			return fmt.Errorf("failed to GET schema %s: %v", schemaId, err)
		}
		count, err := schemaCont.ArrayCount("sites")
		if err != nil {
			return nil // no sites at all
		}
		for i := 0; i < count; i++ {
			elem, err := schemaCont.ArrayElement(i, "sites")
			if err != nil {
				return err
			}
			if models.StripQuotes(elem.S("siteId").String()) == targetSiteId {
				return fmt.Errorf("expected site %q (id=%s) absent from schema %s but it is still present", siteName, targetSiteId, schemaId)
			}
		}
		return nil
	}
}

func testAccCheckMSOSchemaSiteDestroy(s *terraform.State) error {
	msoClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mso_schema_site" {
			continue
		}
		schemaId := rs.Primary.Attributes["schema_id"]
		stateSiteId := rs.Primary.Attributes["site_id"]
		stateTemplate := rs.Primary.Attributes["template_name"]

		cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
		if err != nil {
			// Schema itself has also been destroyed, which implicitly removes
			// any site associations.
			return nil
		}
		count, err := cont.ArrayCount("sites")
		if err != nil {
			return nil
		}
		for i := 0; i < count; i++ {
			elem, err := cont.ArrayElement(i, "sites")
			if err != nil {
				return err
			}
			apiSiteId := models.StripQuotes(elem.S("siteId").String())
			apiTemplate := models.StripQuotes(elem.S("templateName").String())
			if apiSiteId == stateSiteId && apiTemplate == stateTemplate {
				return fmt.Errorf("mso_schema_site (site=%s, template=%s) still exists on schema %s", stateSiteId, stateTemplate, schemaId)
			}
		}
	}
	return nil
}
