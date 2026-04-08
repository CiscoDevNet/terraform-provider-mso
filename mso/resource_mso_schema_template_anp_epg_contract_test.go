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

// msoSchemaTemplateAnpEpgContractSchemaId is set during the first test step's Check to capture the dynamic schema ID for use in the manual deletion PreConfig step.
var msoSchemaTemplateAnpEpgContractSchemaId string

func TestAccMSOSchemaTemplateAnpEpgContractResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateAnpEpgContractResourceDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create EPG Contract as provider") },
				Config:    testAccMSOSchemaTemplateAnpEpgContractConfigProvider(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_anp_epg_contract."+msoSchemaTemplateContractName+"_provider", "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_contract."+msoSchemaTemplateContractName+"_provider", "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_contract."+msoSchemaTemplateContractName+"_provider", "anp_name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_contract."+msoSchemaTemplateContractName+"_provider", "epg_name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_contract."+msoSchemaTemplateContractName+"_provider", "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_contract."+msoSchemaTemplateContractName+"_provider", "relationship_type", "provider"),
					resource.TestCheckResourceAttrSet("mso_schema_template_anp_epg_contract."+msoSchemaTemplateContractName+"_provider", "contract_schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_contract."+msoSchemaTemplateContractName+"_provider", "contract_template_name", msoSchemaTemplateName),
					// Capture the dynamic schema ID from state for use in the manual deletion PreConfig step.
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["mso_schema_template_anp_epg_contract."+msoSchemaTemplateContractName+"_provider"]
						if !ok {
							return fmt.Errorf("EPG Contract resource not found in state")
						}
						msoSchemaTemplateAnpEpgContractSchemaId = rs.Primary.Attributes["schema_id"]
						return nil
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update EPG Contract relationship_type to consumer") },
				Config:    testAccMSOSchemaTemplateAnpEpgContractConfigConsumer(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_contract."+msoSchemaTemplateContractName+"_provider", "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_contract."+msoSchemaTemplateContractName+"_provider", "relationship_type", "consumer"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Reset EPG Contract relationship_type to provider") },
				Config:    testAccMSOSchemaTemplateAnpEpgContractConfigProvider(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_contract."+msoSchemaTemplateContractName+"_provider", "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_contract."+msoSchemaTemplateContractName+"_provider", "relationship_type", "provider"),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import EPG Contract") },
				ResourceName: "mso_schema_template_anp_epg_contract." + msoSchemaTemplateContractName + "_provider",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["mso_schema_template_anp_epg_contract."+msoSchemaTemplateContractName+"_provider"]
					if !ok {
						return "", fmt.Errorf("EPG Contract resource not found in state")
					}
					return fmt.Sprintf("%s/templates/%s/anps/%s/epgs/%s/contracts/%s/relationship_type/%s",
						rs.Primary.Attributes["schema_id"],
						rs.Primary.Attributes["template_name"],
						rs.Primary.Attributes["anp_name"],
						rs.Primary.Attributes["epg_name"],
						rs.Primary.Attributes["contract_name"],
						rs.Primary.Attributes["relationship_type"],
					), nil
				},
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate EPG Contract after manual deletion from NDO")
					msoClient := testAccProvider.Meta().(*client.Client)
					cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", msoSchemaTemplateAnpEpgContractSchemaId))
					if err != nil {
						t.Fatalf("Failed to get schema: %v", err)
					}
					index, _, err := getSchemaTemplateEPGContract(cont, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateContractName, msoSchemaTemplateAnpEpgContractSchemaId, msoSchemaTemplateName, "provider")
					if err != nil {
						t.Fatalf("Failed to fetch contract index: %v", err)
					}
					if index == -1 {
						t.Fatalf("Contract not found for manual deletion")
					}
					contractRemovePatchPayload := models.GetRemovePatchPayload(fmt.Sprintf("/templates/%s/anps/%s/epgs/%s/contractRelationships/%d", msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName, index))
					_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", msoSchemaTemplateAnpEpgContractSchemaId), contractRemovePatchPayload)
					if err != nil {
						t.Fatalf("Failed to manually delete contract: %v", err)
					}
				},
				Config: testAccMSOSchemaTemplateAnpEpgContractConfigProvider(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_contract."+msoSchemaTemplateContractName+"_provider", "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_contract."+msoSchemaTemplateContractName+"_provider", "relationship_type", "provider"),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateAnpEpgContractPrerequisiteConfig() string {
	return fmt.Sprintf(`%s%s%s%s%s%s`, testSiteConfigAnsibleTest(), testTenantConfig(), testSchemaConfig(), testSchemaTemplateAnpConfig(), testSchemaTemplateAnpEpgConfig(), testSchemaTemplateFilterEntryConfig()) + testSchemaTemplateContractConfig()
}

func testAccMSOSchemaTemplateAnpEpgContractConfigProvider() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_template_anp_epg_contract" "%[2]s_provider" {
	schema_id         = mso_schema.%[3]s.id
	template_name     = "%[4]s"
	anp_name          = "%[5]s"
	epg_name          = mso_schema_template_anp_epg.%[6]s.name
	contract_name     = mso_schema_template_contract.%[2]s.contract_name
	relationship_type = "provider"
}
`, testAccMSOSchemaTemplateAnpEpgContractPrerequisiteConfig(), msoSchemaTemplateContractName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName)
}

func testAccMSOSchemaTemplateAnpEpgContractConfigConsumer() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_template_anp_epg_contract" "%[2]s_provider" {
	schema_id         = mso_schema.%[3]s.id
	template_name     = "%[4]s"
	anp_name          = "%[5]s"
	epg_name          = mso_schema_template_anp_epg.%[6]s.name
	contract_name     = mso_schema_template_contract.%[2]s.contract_name
	relationship_type = "consumer"
}
`, testAccMSOSchemaTemplateAnpEpgContractPrerequisiteConfig(), msoSchemaTemplateContractName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName)
}

// testAccCheckMSOSchemaTemplateAnpEpgContractResourceDestroy verifies the contract relationship is removed after test.
func testAccCheckMSOSchemaTemplateAnpEpgContractResourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "mso_schema_template_anp_epg_contract" {
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
				anpCount, err := tempCont.ArrayCount("anps")
				if err != nil {
					return fmt.Errorf("No Anp found")
				}
				for j := 0; j < anpCount; j++ {
					anpCont, err := tempCont.ArrayElement(j, "anps")
					if err != nil {
						return err
					}
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return fmt.Errorf("Unable to get EPG list")
					}
					for k := 0; k < epgCount; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return err
						}
						crefCount, err := epgCont.ArrayCount("contractRelationships")
						if err != nil {
							return fmt.Errorf("Unable to get contract relationships list")
						}
						for l := 0; l < crefCount; l++ {
							crefCont, err := epgCont.ArrayElement(l, "contractRelationships")
							if err != nil {
								return err
							}
							contractRef := models.StripQuotes(crefCont.S("contractRef").String())
							re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/contracts/(.*)")
							match := re.FindStringSubmatch(contractRef)
							if len(match) > 3 && match[3] == rs.Primary.Attributes["contract_name"] {
								return fmt.Errorf("Schema Template ANP EPG Contract still exists")
							}
						}
					}
				}
			}
		}
	}
	return nil
}
