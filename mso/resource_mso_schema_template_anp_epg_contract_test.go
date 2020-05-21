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

func TestAccMSOSchemaTemplateAnpEpgContract_Basic(t *testing.T) {
	var ss TemplateAnpEpgContract
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateAnpEpgContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateAnpEpgContractConfig_basic("provider"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateAnpEpgContractExists("mso_schema_template_anp_epg_contract.contract", &ss),
					testAccCheckMSOSchemaTemplateAnpEpgContractAttributes("provider", &ss),
				),
			},
		},
	})
}

func TestAccMSOSchemaTemplateAnpEpgContract_Update(t *testing.T) {
	var ss TemplateAnpEpgContract

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateAnpEpgContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateAnpEpgContractConfig_basic("provider"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateAnpEpgContractExists("mso_schema_template_anp_epg_contract.contract", &ss),
					testAccCheckMSOSchemaTemplateAnpEpgContractAttributes("provider", &ss),
				),
			},
			{
				Config: testAccCheckMSOTemplateAnpEpgContractConfig_basic("consumer"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateAnpEpgContractExists("mso_schema_template_anp_epg_contract.contract", &ss),
					testAccCheckMSOSchemaTemplateAnpEpgContractAttributes("consumer", &ss),
				),
			},
		},
	})
}

func testAccCheckMSOTemplateAnpEpgContractConfig_basic(relationshiptype string) string {
	return fmt.Sprintf(`
	resource "mso_schema_template_anp_epg_contract" "contract" {
  schema_id = "5c6c16d7270000c710f8094d"
  template_name = "Template1"
  anp_name = "WoS-Cloud-Only-2"
  epg_name = "DB"
  contract_name = "Internet-access"
  relationship_type = "%v"
  
}
`, relationshiptype)
}

func testAccCheckMSOSchemaTemplateAnpEpgContractExists(contractName string, ss *TemplateAnpEpgContract) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs1, err1 := s.RootModule().Resources[contractName]

		if !err1 {
			return fmt.Errorf("Contract %s not found", contractName)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Contract id was set")
		}

		cont, err := client.GetViaURL("api/v1/schemas/5c6c16d7270000c710f8094d")
		if err != nil {
			return err
		}
		count, err := cont.ArrayCount("templates")
		if err != nil {
			return fmt.Errorf("No Template found")
		}
		tp := TemplateAnpEpgContract{}
		found := false
		for i := 0; i < count; i++ {
			tempCont, err := cont.ArrayElement(i, "templates")
			if err != nil {
				return err
			}
			apiTemplateName := models.StripQuotes(tempCont.S("name").String())
			if apiTemplateName == "Template1" {
				anpCount, err := tempCont.ArrayCount("anps")
				if err != nil {
					return fmt.Errorf("Unable to get ANP list")
				}
				for j := 0; j < anpCount; j++ {
					anpCont, err := tempCont.ArrayElement(j, "anps")
					if err != nil {
						return err
					}
					apiANP := models.StripQuotes(anpCont.S("name").String())
					if apiANP == "WoS-Cloud-Only-2" {
						epgCount, err := anpCont.ArrayCount("epgs")
						if err != nil {
							return fmt.Errorf("Unable to get EPG list")
						}
						for k := 0; k < epgCount; k++ {
							epgCont, err := anpCont.ArrayElement(k, "epgs")
							if err != nil {
								return err
							}
							apiEPG := models.StripQuotes(epgCont.S("name").String())
							if apiEPG == "DB" {
								crefCount, err := epgCont.ArrayCount("contractRelationships")
								if err != nil {
									return fmt.Errorf("Unable to get the contract relationships list")
								}
								for l := 0; l < crefCount; l++ {
									crefCont, err := epgCont.ArrayElement(l, "contractRelationships")
									if err != nil {
										return err
									}
									contractRef := models.StripQuotes(crefCont.S("contractRef").String())
									re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/contracts/(.*)")
									match := re.FindStringSubmatch(contractRef)
									apiContract := match[3]
									if apiContract == "Internet-access" {
										tp.relationship_type = models.StripQuotes(crefCont.S("relationshipType").String())
										tp.contract_name = apiContract
										found = true
										break
									}
								}
							}

						}
					}
				}
			}
		}

		if !found {
			return fmt.Errorf("Contract not found from API")
		}

		tp1 := &tp

		*ss = *tp1
		return nil
	}
}

func testAccCheckMSOSchemaTemplateAnpEpgContractDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_template_anp_epg_contract" {
			cont, err := client.GetViaURL("api/v1/schemas/5ea809672c00003bc40a2799")
			if err != nil {
				return nil
			} else {
				count, err := cont.ArrayCount("templates")
				if err != nil {
					return fmt.Errorf("No Template found")
				}
				for i := 0; i < count; i++ {
					tempCont, err := cont.ArrayElement(i, "templates")
					if err != nil {
						return fmt.Errorf("No Template exists")
					}
					apiTemplateName := models.StripQuotes(tempCont.S("name").String())
					if apiTemplateName == "Template1" {
						anpCount, err := tempCont.ArrayCount("anps")
						if err != nil {
							return fmt.Errorf("Unable to get ANP list")
						}
						for j := 0; j < anpCount; j++ {
							anpCont, err := tempCont.ArrayElement(j, "anps")
							if err != nil {
								return err
							}
							apiANP := models.StripQuotes(anpCont.S("name").String())
							if apiANP == "WoS-Cloud-Only-2" {
								epgCount, err := anpCont.ArrayCount("epgs")
								if err != nil {
									return fmt.Errorf("Unable to get EPG list")
								}
								for k := 0; k < epgCount; k++ {
									epgCont, err := anpCont.ArrayElement(k, "epgs")
									if err != nil {
										return err
									}
									apiEPG := models.StripQuotes(epgCont.S("name").String())
									if apiEPG == "DB" {
										crefCount, err := epgCont.ArrayCount("contractRelationships")
										if err != nil {
											return fmt.Errorf("Unable to get the contract relationships list")
										}
										for l := 0; l < crefCount; l++ {
											crefCont, err := epgCont.ArrayElement(l, "contractRelationships")
											if err != nil {
												return err
											}
											contractRef := models.StripQuotes(crefCont.S("contractRef").String())
											re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/contracts/(.*)")
											match := re.FindStringSubmatch(contractRef)
											apiContract := match[3]
											if apiContract == "Internet-access" {
												return fmt.Errorf("Contract still exists.")
											}
										}
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
func testAccCheckMSOSchemaTemplateAnpEpgContractAttributes(relationship_type string, ss *TemplateAnpEpgContract) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if relationship_type != ss.relationship_type {
			return fmt.Errorf("Bad Contract Relationship Type %s", ss.relationship_type)
		}
		return nil
	}
}

type TemplateAnpEpgContract struct {
	contract_name     string
	relationship_type string
}
