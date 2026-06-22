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

// msoSchemaSiteContractServiceGraphSchemaId and
// msoSchemaSiteContractServiceGraphSiteId are set during the first test step's
// Check to capture dynamic IDs for use in the manual deletion PreConfig step.
var msoSchemaSiteContractServiceGraphSchemaId string
var msoSchemaSiteContractServiceGraphSiteId string

// TestAccMSOSchemaSiteContractServiceGraphResource tests the full lifecycle of
// the mso_schema_site_contract_service_graph resource:
//
//  1. Create — attach the service graph to the site contract with cluster
//     interface names for the provider and consumer connectors.
//  2. Error  — node_relationship count (2) does not match template node count (1).
//  3. Error  — provider_connector_redirect_policy set without its required
//     provider_connector_redirect_policy_tenant.
//  4. Update — swap cluster interfaces and set redirect policies for both
//     provider and consumer connectors.
//  5. Import — verify state round-trips through the import path.
//  6. Recreate — manually delete the serviceGraphRelationship via the API and
//     verify Terraform detects the drift and recreates it.
//
// Note: consumer_subnet_ips is only accepted by the API when the service node
// type is Load Balancer. This test uses a Firewall node, so subnet IPs are
// verified as empty in step 1 and omitted from the update config.
func TestAccMSOSchemaSiteContractServiceGraphResource(t *testing.T) {
	resourceRef := "mso_schema_site_contract_service_graph." + msoSchemaTemplateContractName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAPICPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteContractServiceGraphDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: Create site contract service graph with provider=internal / consumer=external")
				},
				Config: testAccMSOSchemaSiteContractServiceGraphConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceRef, "id"),
					resource.TestCheckResourceAttrSet(resourceRef, "schema_id"),
					resource.TestCheckResourceAttr(resourceRef, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttrSet(resourceRef, "site_id"),
					resource.TestCheckResourceAttr(resourceRef, "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr(resourceRef, "service_graph_name", msoSchemaTemplateServiceGraphName),
					resource.TestCheckResourceAttr(resourceRef, "node_relationship.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "node_relationship.0.provider_connector_cluster_interface", msoSchemaSiteContractServiceGraphProviderClusterInterface),
					resource.TestCheckResourceAttr(resourceRef, "node_relationship.0.consumer_connector_cluster_interface", msoSchemaSiteContractServiceGraphConsumerClusterInterface),
					resource.TestCheckResourceAttr(resourceRef, "node_relationship.0.provider_connector_redirect_policy", ""),
					resource.TestCheckResourceAttr(resourceRef, "node_relationship.0.provider_connector_redirect_policy_tenant", ""),
					resource.TestCheckResourceAttr(resourceRef, "node_relationship.0.consumer_connector_redirect_policy", ""),
					resource.TestCheckResourceAttr(resourceRef, "node_relationship.0.consumer_connector_redirect_policy_tenant", ""),
					resource.TestCheckResourceAttr(resourceRef, "node_relationship.0.consumer_subnet_ips.#", "0"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceRef]
						if !ok {
							return fmt.Errorf("resource not found in state: %s", resourceRef)
						}
						msoSchemaSiteContractServiceGraphSchemaId = rs.Primary.Attributes["schema_id"]
						msoSchemaSiteContractServiceGraphSiteId = rs.Primary.Attributes["site_id"]
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Expect error when node_relationship count does not match template graph node count")
				},
				Config:      testAccMSOSchemaSiteContractServiceGraphConfigNodeCountMismatch(),
				ExpectError: regexp.MustCompile(`service graph has 1 service node\(s\) in the template but 2 node_relationship\(s\) were provided`),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Expect error when provider_connector_redirect_policy is set without provider_connector_redirect_policy_tenant")
				},
				Config:      testAccMSOSchemaSiteContractServiceGraphConfigMissingRedirectPolicyTenant(),
				ExpectError: regexp.MustCompile(`provider redirect policy tenant is required`),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update — swap cluster interfaces and set redirect policies")
				},
				Config: testAccMSOSchemaSiteContractServiceGraphConfigUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "node_relationship.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "node_relationship.0.provider_connector_cluster_interface", msoSchemaSiteContractServiceGraphConsumerClusterInterface),
					resource.TestCheckResourceAttr(resourceRef, "node_relationship.0.consumer_connector_cluster_interface", msoSchemaSiteContractServiceGraphProviderClusterInterface),
					resource.TestCheckResourceAttr(resourceRef, "node_relationship.0.provider_connector_redirect_policy_tenant", msoTenantName),
					resource.TestCheckResourceAttr(resourceRef, "node_relationship.0.provider_connector_redirect_policy", msoSchemaSiteContractServiceGraphProviderRedirectPolicy),
					resource.TestCheckResourceAttr(resourceRef, "node_relationship.0.consumer_connector_redirect_policy_tenant", msoTenantName2),
					resource.TestCheckResourceAttr(resourceRef, "node_relationship.0.consumer_connector_redirect_policy", msoSchemaSiteContractServiceGraphConsumerRedirectPolicy),
					resource.TestCheckResourceAttr(resourceRef, "node_relationship.0.consumer_subnet_ips.#", "0"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import site contract service graph") },
				ResourceName:      resourceRef,
				ImportState:       true,
				ImportStateIdFunc: testAccMSOSchemaSiteContractServiceGraphImportStateId(resourceRef),
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate site contract service graph after manual deletion")
					msoClient := testAccProvider.Meta().(*client.Client)
					path := fmt.Sprintf("/sites/%s-%s/contracts/%s/serviceGraphRelationship",
						msoSchemaSiteContractServiceGraphSiteId,
						msoSchemaTemplateName,
						msoSchemaTemplateContractName,
					)
					_, err := msoClient.PatchbyID(
						fmt.Sprintf("api/v1/schemas/%s", msoSchemaSiteContractServiceGraphSchemaId),
						models.GetRemovePatchPayload(path),
					)
					if err != nil {
						t.Fatalf("Failed to manually delete site contract service graph relationship: %v", err)
					}
				},
				Config: testAccMSOSchemaSiteContractServiceGraphConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceRef, "id"),
					resource.TestCheckResourceAttr(resourceRef, "node_relationship.0.provider_connector_cluster_interface", msoSchemaSiteContractServiceGraphProviderClusterInterface),
					resource.TestCheckResourceAttr(resourceRef, "node_relationship.0.consumer_connector_cluster_interface", msoSchemaSiteContractServiceGraphConsumerClusterInterface),
				),
			},
		},
	})
}

// testAccCheckMSOSchemaSiteContractServiceGraphDestroy verifies that the
// serviceGraphRelationship has been removed from the site contract in the
// schema after the test.
func testAccCheckMSOSchemaSiteContractServiceGraphDestroy(s *terraform.State) error {
	msoClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mso_schema_site_contract_service_graph" {
			continue
		}

		cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", rs.Primary.Attributes["schema_id"]))
		if err != nil {
			// Schema itself is gone — the relationship cannot exist.
			return nil
		}

		siteCount, err := cont.ArrayCount("sites")
		if err != nil {
			return nil
		}
		for i := 0; i < siteCount; i++ {
			siteCont, err := cont.ArrayElement(i, "sites")
			if err != nil {
				continue
			}
			if models.StripQuotes(siteCont.S("siteId").String()) != rs.Primary.Attributes["site_id"] ||
				models.StripQuotes(siteCont.S("templateName").String()) != rs.Primary.Attributes["template_name"] {
				continue
			}
			contractCount, err := siteCont.ArrayCount("contracts")
			if err != nil {
				continue
			}
			for j := 0; j < contractCount; j++ {
				contractCont, err := siteCont.ArrayElement(j, "contracts")
				if err != nil {
					continue
				}
				contractRef := models.StripQuotes(contractCont.S("contractRef").String())
				contractTokens := strings.Split(contractRef, "/")
				if contractTokens[len(contractTokens)-1] != rs.Primary.Attributes["contract_name"] {
					continue
				}
				if contractCont.Exists("serviceGraphRelationship") {
					return fmt.Errorf("mso_schema_site_contract_service_graph (%s) still exists", rs.Primary.ID)
				}
			}
		}
	}

	return nil
}

// testAccMSOSchemaSiteContractServiceGraphImportStateId builds the import ID
// from the resource's state attributes.
// Import ID format: {schema_id}/sites/{site_id}/templates/{template_name}/contracts/{contract_name}
func testAccMSOSchemaSiteContractServiceGraphImportStateId(resourceRef string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceRef]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceRef)
		}
		return fmt.Sprintf("%s/sites/%s/templates/%s/contracts/%s",
			rs.Primary.Attributes["schema_id"],
			rs.Primary.Attributes["site_id"],
			rs.Primary.Attributes["template_name"],
			rs.Primary.Attributes["contract_name"],
		), nil
	}
}

// testAccMSOSchemaSiteContractServiceGraphPrereqConfig builds the full
// prerequisite stack needed by both the create and update configs:
//   - schema + site association (msoTenantName, one site)
//   - template service graph (1 firewall node)
//   - VRF, BD, filter entry, contract, VRF contract provider
//
// testAccMSOSchemaSiteContractServiceGraphPrereqConfig builds the self-contained
// prerequisite stack for both resource and datasource tests:
//   - schema + tenant (msoTenantName) + consumer tenant (for cross-tenant redirect policy)
//   - template service graph (firewall node)
//   - VRF, BD, filter, contract, VRF-contract binding
//   - template contract service graph (BD connector for provider and consumer)
//   - site service graph (assigns the firewall device to the site-level node)
func testAccMSOSchemaSiteContractServiceGraphPrereqConfig() string {
	return fmt.Sprintf(`%[1]s
%[14]s
resource "mso_schema_template_service_graph" "%[2]s" {
  schema_id          = mso_schema.%[3]s.id
  template_name      = "%[4]s"
  service_graph_name = "%[2]s"

  service_node {
    type = "firewall"
  }
}

%[5]s%[6]s%[7]s%[8]s%[9]s

resource "mso_schema_template_contract_service_graph" "%[10]s" {
  schema_id          = mso_schema.%[3]s.id
  template_name      = "%[4]s"
  contract_name      = mso_schema_template_contract.%[10]s.contract_name
  service_graph_name = mso_schema_template_service_graph.%[2]s.service_graph_name

  node_relationship {
    provider_connector_bd_name = mso_schema_template_bd.%[11]s.name
    consumer_connector_bd_name = mso_schema_template_bd.%[11]s.name
  }
}

resource "mso_schema_site_service_graph" "%[2]s" {
  schema_id          = mso_schema.%[3]s.id
  template_name      = "%[4]s"
  site_id            = mso_schema_site.%[12]s.site_id
  service_graph_name = "%[2]s"

  service_node {
    device_dn = "%[13]s"
  }

  depends_on = [mso_schema_template_service_graph.%[2]s]
}
`,
		testSchemaWithSingleSiteAssociationConfig(), // %[1]s
		msoSchemaTemplateServiceGraphName,           // %[2]s
		msoSchemaName,                               // %[3]s
		msoSchemaTemplateName,                       // %[4]s
		testSchemaTemplateVrfConfig(),               // %[5]s
		testSchemaTemplateBdConfig(),                // %[6]s
		testSchemaTemplateFilterEntryConfig(),       // %[7]s
		testSchemaTemplateContractConfig(),          // %[8]s
		testSchemaTemplateVrfContractConfig(),       // %[9]s
		msoSchemaTemplateContractName,               // %[10]s
		msoSchemaTemplateBdName,                     // %[11]s
		msoSchemaSiteResourceLabel1,                 // %[12]s
		msoSchemaSiteContractServiceGraphDeviceDn,   // %[13]s
		testTenantConfigOneSite(msoTenantName2),     // %[14]s
	)
}

func testAccMSOSchemaSiteContractServiceGraphConfigCreate() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_site_contract_service_graph" "%[2]s" {
  schema_id          = mso_schema.%[3]s.id
  site_id            = mso_schema_site.%[4]s.site_id
  template_name      = "%[5]s"
  contract_name      = mso_schema_template_contract.%[2]s.contract_name
  service_graph_name = mso_schema_site_service_graph.%[6]s.service_graph_name

  node_relationship {
    provider_connector_cluster_interface = "%[7]s"
    consumer_connector_cluster_interface = "%[8]s"
  }
}
`,
		testAccMSOSchemaSiteContractServiceGraphPrereqConfig(),    // %[1]s
		msoSchemaTemplateContractName,                             // %[2]s
		msoSchemaName,                                             // %[3]s
		msoSchemaSiteResourceLabel1,                               // %[4]s
		msoSchemaTemplateName,                                     // %[5]s
		msoSchemaTemplateServiceGraphName,                         // %[6]s
		msoSchemaSiteContractServiceGraphProviderClusterInterface, // %[7]s
		msoSchemaSiteContractServiceGraphConsumerClusterInterface, // %[8]s
	)
}

func testAccMSOSchemaSiteContractServiceGraphConfigUpdate() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_site_contract_service_graph" "%[2]s" {
  schema_id          = mso_schema.%[3]s.id
  site_id            = mso_schema_site.%[4]s.site_id
  template_name      = "%[5]s"
  contract_name      = mso_schema_template_contract.%[2]s.contract_name
  service_graph_name = mso_schema_site_service_graph.%[6]s.service_graph_name

  node_relationship {
    provider_connector_cluster_interface      = "%[7]s"
    provider_connector_redirect_policy_tenant = "%[9]s"
    provider_connector_redirect_policy        = "%[10]s"
    consumer_connector_cluster_interface      = "%[8]s"
    consumer_connector_redirect_policy_tenant = "%[11]s"
    consumer_connector_redirect_policy        = "%[12]s"
  }
}
`,
		testAccMSOSchemaSiteContractServiceGraphPrereqConfig(),    // %[1]s
		msoSchemaTemplateContractName,                             // %[2]s
		msoSchemaName,                                             // %[3]s
		msoSchemaSiteResourceLabel1,                               // %[4]s
		msoSchemaTemplateName,                                     // %[5]s
		msoSchemaTemplateServiceGraphName,                         // %[6]s
		msoSchemaSiteContractServiceGraphConsumerClusterInterface, // %[7]s
		msoSchemaSiteContractServiceGraphProviderClusterInterface, // %[8]s
		msoTenantName,                                             // %[9]s
		msoSchemaSiteContractServiceGraphProviderRedirectPolicy,   // %[10]s
		msoTenantName2,                                            // %[11]s
		msoSchemaSiteContractServiceGraphConsumerRedirectPolicy,   // %[12]s
	)
}

// testAccMSOSchemaSiteContractServiceGraphConfigMissingRedirectPolicyTenant
// sets provider_connector_redirect_policy without the required
// provider_connector_redirect_policy_tenant, exercising the cross-field
// validation guard in getSiteServiceNodesRelationshipObject.
func testAccMSOSchemaSiteContractServiceGraphConfigMissingRedirectPolicyTenant() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_site_contract_service_graph" "%[2]s" {
  schema_id          = mso_schema.%[3]s.id
  site_id            = mso_schema_site.%[4]s.site_id
  template_name      = "%[5]s"
  contract_name      = mso_schema_template_contract.%[2]s.contract_name
  service_graph_name = mso_schema_site_service_graph.%[6]s.service_graph_name

  node_relationship {
    provider_connector_cluster_interface = "%[7]s"
    provider_connector_redirect_policy   = "%[8]s"
    consumer_connector_cluster_interface = "%[9]s"
  }
}
`,
		testAccMSOSchemaSiteContractServiceGraphPrereqConfig(),    // %[1]s
		msoSchemaTemplateContractName,                             // %[2]s
		msoSchemaName,                                             // %[3]s
		msoSchemaSiteResourceLabel1,                               // %[4]s
		msoSchemaTemplateName,                                     // %[5]s
		msoSchemaTemplateServiceGraphName,                         // %[6]s
		msoSchemaSiteContractServiceGraphProviderClusterInterface, // %[7]s
		msoSchemaSiteContractServiceGraphProviderRedirectPolicy,   // %[8]s
		msoSchemaSiteContractServiceGraphConsumerClusterInterface, // %[9]s
	)
}

// testAccMSOSchemaSiteContractServiceGraphConfigNodeCountMismatch configures a
// site contract service graph with 2 node_relationship blocks against a
// template service graph that only has 1 node. This exercises the count-mismatch
// guard in postSiteContractServiceGraphConfig.
func testAccMSOSchemaSiteContractServiceGraphConfigNodeCountMismatch() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_site_contract_service_graph" "%[2]s" {
  schema_id          = mso_schema.%[3]s.id
  site_id            = mso_schema_site.%[4]s.site_id
  template_name      = "%[5]s"
  contract_name      = mso_schema_template_contract.%[2]s.contract_name
  service_graph_name = mso_schema_site_service_graph.%[6]s.service_graph_name

  node_relationship {
    provider_connector_cluster_interface = "%[7]s"
    consumer_connector_cluster_interface = "%[8]s"
  }

  node_relationship {
    provider_connector_cluster_interface = "%[7]s"
    consumer_connector_cluster_interface = "%[8]s"
  }
}
`,
		testAccMSOSchemaSiteContractServiceGraphPrereqConfig(),    // %[1]s
		msoSchemaTemplateContractName,                             // %[2]s
		msoSchemaName,                                             // %[3]s
		msoSchemaSiteResourceLabel1,                               // %[4]s
		msoSchemaTemplateName,                                     // %[5]s
		msoSchemaTemplateServiceGraphName,                         // %[6]s
		msoSchemaSiteContractServiceGraphProviderClusterInterface, // %[7]s
		msoSchemaSiteContractServiceGraphConsumerClusterInterface, // %[8]s
	)
}
