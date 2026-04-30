package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// msoSchemaTemplateExtEpgContractSchemaId is set during the first test step's Check to capture the dynamic schema ID for use in the manual deletion PreConfig step.
var msoSchemaTemplateExtEpgContractSchemaId string

func TestAccMSOSchemaTemplateExternalEpgContractResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateExtEpgContractDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create External EPG Contract as provider") },
				Config:    testAccMSOSchemaTemplateExtEpgContractConfigProvider(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_external_epg_contract."+msoSchemaTemplateContractName+"_provider", "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg_contract."+msoSchemaTemplateContractName+"_provider", "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg_contract."+msoSchemaTemplateContractName+"_provider", "external_epg_name", msoSchemaTemplateExtEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg_contract."+msoSchemaTemplateContractName+"_provider", "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg_contract."+msoSchemaTemplateContractName+"_provider", "relationship_type", "provider"),
					resource.TestCheckResourceAttrSet("mso_schema_template_external_epg_contract."+msoSchemaTemplateContractName+"_provider", "contract_schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg_contract."+msoSchemaTemplateContractName+"_provider", "contract_template_name", msoSchemaTemplateName),
					// Capture the dynamic schema ID from state for use in the manual deletion PreConfig step.
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["mso_schema_template_external_epg_contract."+msoSchemaTemplateContractName+"_provider"]
						if !ok {
							return fmt.Errorf("External EPG Contract resource not found in state")
						}
						msoSchemaTemplateExtEpgContractSchemaId = rs.Primary.Attributes["schema_id"]
						return nil
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update External EPG Contract relationship_type to consumer") },
				Config:    testAccMSOSchemaTemplateExtEpgContractConfigConsumer(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_external_epg_contract."+msoSchemaTemplateContractName+"_provider", "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg_contract."+msoSchemaTemplateContractName+"_provider", "relationship_type", "consumer"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Reset External EPG Contract relationship_type to provider") },
				Config:    testAccMSOSchemaTemplateExtEpgContractConfigProvider(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_external_epg_contract."+msoSchemaTemplateContractName+"_provider", "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg_contract."+msoSchemaTemplateContractName+"_provider", "relationship_type", "provider"),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import External EPG Contract") },
				ResourceName: "mso_schema_template_external_epg_contract." + msoSchemaTemplateContractName + "_provider",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["mso_schema_template_external_epg_contract."+msoSchemaTemplateContractName+"_provider"]
					if !ok {
						return "", fmt.Errorf("External EPG Contract resource not found in state")
					}
					return fmt.Sprintf("%s/templates/%s/externalEpgs/%s/contractRelationships/%s/%s",
						rs.Primary.Attributes["schema_id"],
						rs.Primary.Attributes["template_name"],
						rs.Primary.Attributes["external_epg_name"],
						rs.Primary.Attributes["contract_name"],
						rs.Primary.Attributes["relationship_type"],
					), nil
				},
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate External EPG Contract after manual deletion from NDO")
					msoClient := testAccProvider.Meta().(*client.Client)
					cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", msoSchemaTemplateExtEpgContractSchemaId))
					if err != nil {
						t.Fatalf("Failed to get schema: %v", err)
					}
					index, _, err := getSchemaTemplateExtEpgContract(cont, msoSchemaTemplateName, msoSchemaTemplateExtEpgName, msoSchemaTemplateContractName, msoSchemaTemplateExtEpgContractSchemaId, msoSchemaTemplateName, "provider")
					if err != nil {
						t.Fatalf("Failed to fetch contract index: %v", err)
					}
					if index == -1 {
						t.Fatalf("External EPG Contract not found for manual deletion")
					}
					contractRemovePatchPayload := models.GetRemovePatchPayload(fmt.Sprintf("/templates/%s/externalEpgs/%s/contractRelationships/%d", msoSchemaTemplateName, msoSchemaTemplateExtEpgName, index))
					_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", msoSchemaTemplateExtEpgContractSchemaId), contractRemovePatchPayload)
					if err != nil {
						t.Fatalf("Failed to manually delete External EPG Contract: %v", err)
					}
				},
				Config: testAccMSOSchemaTemplateExtEpgContractConfigProvider(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_external_epg_contract."+msoSchemaTemplateContractName+"_provider", "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_external_epg_contract."+msoSchemaTemplateContractName+"_provider", "relationship_type", "provider"),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateExtEpgContractPrerequisiteConfig() string {
	return fmt.Sprintf(`%s%s%s%s%s%s`,
		testSiteConfigAnsibleTest(),
		testTenantConfig(),
		testSchemaConfig(),
		testSchemaTemplateVrfConfig(),
		testSchemaTemplateExtEpgConfig(),
		testSchemaTemplateFilterEntryConfig(),
	) + testSchemaTemplateContractConfig()
}

func testAccMSOSchemaTemplateExtEpgContractConfigProvider() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_template_external_epg_contract" "%[2]s_provider" {
	schema_id         = mso_schema.%[3]s.id
	template_name     = "%[4]s"
	external_epg_name = mso_schema_template_external_epg.%[5]s.external_epg_name
	contract_name     = mso_schema_template_contract.%[2]s.contract_name
	relationship_type = "provider"
}
`, testAccMSOSchemaTemplateExtEpgContractPrerequisiteConfig(), msoSchemaTemplateContractName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateExtEpgName)
}

func testAccMSOSchemaTemplateExtEpgContractConfigConsumer() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_template_external_epg_contract" "%[2]s_provider" {
	schema_id         = mso_schema.%[3]s.id
	template_name     = "%[4]s"
	external_epg_name = mso_schema_template_external_epg.%[5]s.external_epg_name
	contract_name     = mso_schema_template_contract.%[2]s.contract_name
	relationship_type = "consumer"
}
`, testAccMSOSchemaTemplateExtEpgContractPrerequisiteConfig(), msoSchemaTemplateContractName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateExtEpgName)
}

// testAccCheckMSOSchemaTemplateExtEpgContractDestroy verifies the contract relationship is removed after test.
func testAccCheckMSOSchemaTemplateExtEpgContractDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "mso_schema_template_external_epg_contract" {
			schemaID := rs.Primary.Attributes["schema_id"]
			cont, err := client.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
			if err != nil {
				return nil
			}
			count, err := cont.ArrayCount("templates")
			if err != nil {
				return fmt.Errorf("No Template found")
			}
			for i := 0; i < count; i++ {
				tempCont, err := cont.ArrayElement(i, "templates")
				if err != nil {
					return fmt.Errorf("No template exists")
				}
				epgCount, err := tempCont.ArrayCount("externalEpgs")
				if err != nil {
					return fmt.Errorf("Unable to get External EPG list")
				}
				for j := 0; j < epgCount; j++ {
					epgCont, err := tempCont.ArrayElement(j, "externalEpgs")
					if err != nil {
						return err
					}
					crefCount, err := epgCont.ArrayCount("contractRelationships")
					if err != nil {
						return fmt.Errorf("Unable to get contract relationships list")
					}
					for k := 0; k < crefCount; k++ {
						crefCont, err := epgCont.ArrayElement(k, "contractRelationships")
						if err != nil {
							return err
						}
						contractRef := models.StripQuotes(crefCont.S("contractRef").String())
						re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/contracts/(.*)")
						match := re.FindStringSubmatch(contractRef)
						if len(match) > 3 && match[3] == rs.Primary.Attributes["contract_name"] {
							return fmt.Errorf("Schema Template External EPG Contract still exists")
						}
					}
				}
			}
		}
	}
	return nil
}
