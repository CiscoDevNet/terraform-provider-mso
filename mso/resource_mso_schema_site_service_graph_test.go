package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// TestAccMSOSchemaSiteServiceGraphResource tests the full lifecycle of the
// mso_schema_site_service_graph resource:
//
//  1. Create — attach firewall1 device with default connector types (none/none)
//  2. Error  — service_node count mismatch against template node count
//  3. Update — change device to firewall2 and set connector types to redir/redir
//  4. Import — verify state round-trips through the import path
func TestAccMSOSchemaSiteServiceGraphResource(t *testing.T) {
	resourceRef := "mso_schema_site_service_graph." + msoSchemaTemplateServiceGraphName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteServiceGraphDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: Create Site Service Graph with firewall1 device and default connector types (none/none)")
				},
				Config: testAccMSOSchemaSiteServiceGraphConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceRef, "schema_id"),
					resource.TestCheckResourceAttr(resourceRef, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttrSet(resourceRef, "site_id"),
					resource.TestCheckResourceAttr(resourceRef, "service_graph_name", msoSchemaTemplateServiceGraphName),
					resource.TestCheckResourceAttr(resourceRef, "service_node.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "service_node.0.device_dn", msoSchemaSiteServiceGraphDeviceDn),
					resource.TestCheckResourceAttr(resourceRef, "service_node.0.consumer_connector_type", "none"),
					resource.TestCheckResourceAttr(resourceRef, "service_node.0.provider_connector_type", "none"),
					// consumer_interface and provider_interface are only applicable for cloud sites;
					// they default to empty for on-prem sites used in this test.
					resource.TestCheckResourceAttr(resourceRef, "service_node.0.consumer_interface", ""),
					resource.TestCheckResourceAttr(resourceRef, "service_node.0.provider_interface", ""),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Expect error when service_node count does not match template graph node count")
				},
				Config:      testAccMSOSchemaSiteServiceGraphConfigNodeCountMismatch(),
				ExpectError: regexp.MustCompile(`service graph has 1 service node\(s\) in the template but 2 service node\(s\) were provided`),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update Site Service Graph - change device to firewall2 and set connector types to redir/redir")
				},
				Config: testAccMSOSchemaSiteServiceGraphConfigUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "service_node.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "service_node.0.device_dn", msoSchemaSiteServiceGraphDeviceDn2),
					resource.TestCheckResourceAttr(resourceRef, "service_node.0.consumer_connector_type", "redir"),
					resource.TestCheckResourceAttr(resourceRef, "service_node.0.provider_connector_type", "redir"),
					resource.TestCheckResourceAttr(resourceRef, "service_node.0.consumer_interface", ""),
					resource.TestCheckResourceAttr(resourceRef, "service_node.0.provider_interface", ""),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Site Service Graph") },
				ResourceName:      resourceRef,
				ImportState:       true,
				ImportStateIdFunc: testAccMSOSchemaSiteServiceGraphImportStateId(resourceRef),
				ImportStateVerify: true,
			},
		},
	})
}

// testAccCheckMSOSchemaSiteServiceGraphDestroy verifies that the site service
// graph has been removed from the schema after the test. The delete operation
// on this resource is a no-op in the API, but the graph is removed implicitly
// when the template service graph is destroyed during test teardown.
func testAccCheckMSOSchemaSiteServiceGraphDestroy(s *terraform.State) error {
	msoClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mso_schema_site_service_graph" {
			continue
		}

		cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", rs.Primary.Attributes["schema_id"]))
		if err != nil {
			// Schema itself is gone — the site service graph cannot exist.
			return nil
		}

		_, _, err = getSiteServiceGraphCont(
			cont,
			rs.Primary.Attributes["schema_id"],
			rs.Primary.Attributes["template_name"],
			rs.Primary.Attributes["site_id"],
			rs.Primary.Attributes["service_graph_name"],
		)
		if err == nil {
			return fmt.Errorf("mso_schema_site_service_graph (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

// testAccMSOSchemaSiteServiceGraphImportStateId builds the import ID string
// from the resource's state attributes.
// Import ID format: {schema_id}/sites/{site_id}/template/{template_name}/serviceGraphs/{graph_name}
func testAccMSOSchemaSiteServiceGraphImportStateId(resourceRef string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceRef]
		if !ok {
			return "", fmt.Errorf("Resource %s not found in state", resourceRef)
		}
		return fmt.Sprintf("%s/sites/%s/template/%s/serviceGraphs/%s",
			rs.Primary.Attributes["schema_id"],
			rs.Primary.Attributes["site_id"],
			rs.Primary.Attributes["template_name"],
			rs.Primary.Attributes["service_graph_name"],
		), nil
	}
}

// testAccMSOSchemaSiteServiceGraphPrereqConfig returns the prerequisite HCL
// shared by both the create and update test steps: schema + site association +
// template service graph with a single firewall node.
func testAccMSOSchemaSiteServiceGraphPrereqConfig() string {
	return fmt.Sprintf(`%s
resource "mso_schema_template_service_graph" "%[2]s" {
  schema_id          = mso_schema.%[3]s.id
  template_name      = "%[4]s"
  service_graph_name = "%[2]s"

  service_node {
    type = "firewall"
  }
}
`, testSchemaWithAnsibleTestTenantAndSingleSiteConfig(), msoSchemaTemplateServiceGraphName, msoSchemaName, msoSchemaTemplateName)
}

func testAccMSOSchemaSiteServiceGraphConfigCreate() string {
	return fmt.Sprintf(`%s
resource "mso_schema_site_service_graph" "%[2]s" {
  schema_id          = mso_schema.%[3]s.id
  template_name      = "%[4]s"
  site_id            = mso_schema_site.%[5]s.site_id
  service_graph_name = "%[2]s"

  service_node {
    device_dn = "%[6]s"
  }
}
`, testAccMSOSchemaSiteServiceGraphPrereqConfig(), msoSchemaTemplateServiceGraphName, msoSchemaName, msoSchemaTemplateName, msoSchemaSiteResourceLabel1, msoSchemaSiteServiceGraphDeviceDn)
}

func testAccMSOSchemaSiteServiceGraphConfigUpdate() string {
	return fmt.Sprintf(`%s
resource "mso_schema_site_service_graph" "%[2]s" {
  schema_id          = mso_schema.%[3]s.id
  template_name      = "%[4]s"
  site_id            = mso_schema_site.%[5]s.site_id
  service_graph_name = "%[2]s"

  service_node {
    device_dn               = "%[6]s"
    consumer_connector_type = "redir"
    provider_connector_type = "redir"
  }
}
`, testAccMSOSchemaSiteServiceGraphPrereqConfig(), msoSchemaTemplateServiceGraphName, msoSchemaName, msoSchemaTemplateName, msoSchemaSiteResourceLabel1, msoSchemaSiteServiceGraphDeviceDn2)
}

// testAccMSOSchemaSiteServiceGraphConfigNodeCountMismatch configures a site
// service graph with 2 service_node blocks against a template graph that only
// has 1 node. This exercises the count-mismatch guard in createSiteServiceNodeList.
func testAccMSOSchemaSiteServiceGraphConfigNodeCountMismatch() string {
	return fmt.Sprintf(`%s
resource "mso_schema_site_service_graph" "%[2]s" {
  schema_id          = mso_schema.%[3]s.id
  template_name      = "%[4]s"
  site_id            = mso_schema_site.%[5]s.site_id
  service_graph_name = "%[2]s"

  service_node {
    device_dn = "%[6]s"
  }

  service_node {
    device_dn = "%[6]s"
  }
}
`, testAccMSOSchemaSiteServiceGraphPrereqConfig(), msoSchemaTemplateServiceGraphName, msoSchemaName, msoSchemaTemplateName, msoSchemaSiteResourceLabel1, msoSchemaSiteServiceGraphDeviceDn)
}
