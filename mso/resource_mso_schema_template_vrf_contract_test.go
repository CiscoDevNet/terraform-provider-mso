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

// msoSchemaTemplateVrfContractSchemaId is set during the first test step's Check to capture the dynamic schema ID for use in the manual deletion PreConfig step.
var msoSchemaTemplateVrfContractSchemaId string

func TestAccMSOSchemaTemplateVrfContractResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateVrfContractDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create VRF contract as provider") },
				Config:    testAccMSOSchemaTemplateVrfContractConfigProvider(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_vrf_contract.vrf_contract", "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf_contract.vrf_contract", "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf_contract.vrf_contract", "vrf_name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf_contract.vrf_contract", "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf_contract.vrf_contract", "relationship_type", "provider"),
					resource.TestCheckResourceAttrSet("mso_schema_template_vrf_contract.vrf_contract", "contract_schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf_contract.vrf_contract", "contract_template_name", msoSchemaTemplateName),
					// Capture the dynamic schema ID from state for use in the manual deletion PreConfig step
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["mso_schema_template_vrf_contract.vrf_contract"]
						if !ok {
							return fmt.Errorf("VRF contract resource not found in state")
						}
						msoSchemaTemplateVrfContractSchemaId = rs.Primary.Attributes["schema_id"]
						return nil
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: ForceNew replacement - change relationship_type to consumer") },
				Config:    testAccMSOSchemaTemplateVrfContractConfigConsumer(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_vrf_contract.vrf_contract", "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf_contract.vrf_contract", "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf_contract.vrf_contract", "vrf_name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf_contract.vrf_contract", "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf_contract.vrf_contract", "relationship_type", "consumer"),
					resource.TestCheckResourceAttrSet("mso_schema_template_vrf_contract.vrf_contract", "contract_schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf_contract.vrf_contract", "contract_template_name", msoSchemaTemplateName),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import VRF contract") },
				ResourceName: "mso_schema_template_vrf_contract.vrf_contract",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["mso_schema_template_vrf_contract.vrf_contract"]
					if !ok {
						return "", fmt.Errorf("VRF contract resource not found in state")
					}
					return fmt.Sprintf("%s/templates/%s/vrfs/%s/contracts/%s/relationship_type/%s",
						rs.Primary.Attributes["schema_id"],
						rs.Primary.Attributes["template_name"],
						rs.Primary.Attributes["vrf_name"],
						rs.Primary.Attributes["contract_name"],
						rs.Primary.Attributes["relationship_type"],
					), nil
				},
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate VRF contract after manual deletion from NDO")
					msoClient := testAccProvider.Meta().(*client.Client)
					vrfContractRemovePatchPayload := models.GetRemovePatchPayload(fmt.Sprintf("/templates/%s/vrfs/%s/%s/0", msoSchemaTemplateName, msoSchemaTemplateVrfName, humanToApiType["consumer"]))
					_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", msoSchemaTemplateVrfContractSchemaId), vrfContractRemovePatchPayload)
					if err != nil {
						t.Fatalf("Failed to manually delete VRF contract: %v", err)
					}
				},
				Config: testAccMSOSchemaTemplateVrfContractConfigConsumer(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_vrf_contract.vrf_contract", "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_vrf_contract.vrf_contract", "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf_contract.vrf_contract", "vrf_name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf_contract.vrf_contract", "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_vrf_contract.vrf_contract", "relationship_type", "consumer"),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateVrfContractPrerequisiteConfig() string {
	return fmt.Sprintf(`%s%s%s%s%s%s`, testSiteConfigAnsibleTest(), testTenantConfig(), testSchemaConfig(), testSchemaTemplateVrfConfig(), testSchemaTemplateFilterEntryConfig(), testSchemaTemplateContractConfig())
}

func testAccMSOSchemaTemplateVrfContractConfigProvider() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_vrf_contract" "vrf_contract" {
		schema_id         = mso_schema.%[2]s.id
		template_name     = "%[3]s"
		vrf_name          = mso_schema_template_vrf.%[4]s.name
		contract_name     = mso_schema_template_contract.%[5]s.contract_name
		relationship_type = "provider"
	}`, testAccMSOSchemaTemplateVrfContractPrerequisiteConfig(), msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfName, msoSchemaTemplateContractName)
}

func testAccMSOSchemaTemplateVrfContractConfigConsumer() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_vrf_contract" "vrf_contract" {
		schema_id         = mso_schema.%[2]s.id
		template_name     = "%[3]s"
		vrf_name          = mso_schema_template_vrf.%[4]s.name
		contract_name     = mso_schema_template_contract.%[5]s.contract_name
		relationship_type = "consumer"
	}`, testAccMSOSchemaTemplateVrfContractPrerequisiteConfig(), msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfName, msoSchemaTemplateContractName)
}

func testAccCheckMSOSchemaTemplateVrfContractDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "mso_schema_template_vrf_contract" {
			schemaID := rs.Primary.Attributes["schema_id"]
			con, err := client.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
			if err != nil {
				return nil
			}
			count, err := con.ArrayCount("templates")
			if err != nil {
				return fmt.Errorf("No Template found")
			}
			templateName := rs.Primary.Attributes["template_name"]
			vrfName := rs.Primary.Attributes["vrf_name"]
			contractName := rs.Primary.Attributes["contract_name"]
			relationshipType := rs.Primary.Attributes["relationship_type"]
			for i := 0; i < count; i++ {
				tempCont, err := con.ArrayElement(i, "templates")
				if err != nil {
					return fmt.Errorf("No template exists")
				}
				apiTemplate := models.StripQuotes(tempCont.S("name").String())
				if apiTemplate == templateName {
					vrfCount, err := tempCont.ArrayCount("vrfs")
					if err != nil {
						return fmt.Errorf("No Vrf found")
					}
					for j := 0; j < vrfCount; j++ {
						vrfCont, err := tempCont.ArrayElement(j, "vrfs")
						if err != nil {
							return err
						}
						apiVRF := models.StripQuotes(vrfCont.S("name").String())
						if apiVRF == vrfName {
							contractCount, err := vrfCont.ArrayCount(humanToApiType[relationshipType])
							if err != nil {
								return fmt.Errorf("Unable to get contract Relationships list")
							}
							for k := 0; k < contractCount; k++ {
								contractCont, err := vrfCont.ArrayElement(k, humanToApiType[relationshipType])
								if err != nil {
									return err
								}
								contractRef := models.StripQuotes(contractCont.S("contractRef").String())
								re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/contracts/(.*)")
								split := re.FindStringSubmatch(contractRef)
								if contractRef != "{}" && contractRef != "" {
									if split[3] == contractName {
										return fmt.Errorf("VRF contract still exists")
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return nil
}
