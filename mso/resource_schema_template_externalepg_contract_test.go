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

func TestAccMSOSchemaTemplateExternalepgContract_Basic(t *testing.T) {
	var ss TemplateExternalepgContract
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateExternalepgContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateExternalepgContractConfig_basic("provider"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateExternalepgContractExists("mso_schema_template_externalepg_contract.c1", &ss),
					testAccCheckMSOSchemaTemplateExternalepgContractAttributes("provider", &ss),
				),
			},
		},
	})
}

func TestAccMSOSchemaTemplateExternalepgContract_Update(t *testing.T) {
	var ss TemplateExternalepgContract

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateExternalepgContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateExternalepgContractConfig_basic("provider"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateExternalepgContractExists("mso_schema_template_externalepg_contract.c1", &ss),
					testAccCheckMSOSchemaTemplateExternalepgContractAttributes("provider", &ss),
				),
			},
			{
				Config: testAccCheckMSOTemplateExternalepgContractConfig_basic("consumer"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateExternalepgContractExists("mso_schema_template_externalepg_contract.c1", &ss),
					testAccCheckMSOSchemaTemplateExternalepgContractAttributes("consumer", &ss),
				),
			},
		},
	})
}

func testAccCheckMSOTemplateExternalepgContractConfig_basic(name string) string {
	return fmt.Sprintf(`
	resource "mso_schema_template_externalepg_contract" "c1" {
		schema_id = "5ea809672c00003bc40a2799"
		template_name = "Template1"
		contract_name = "contract9999"
		external_epg_name = "UntitledExternalEPG1"
		relationship_type = "%s"
	}
`, name)
}

func testAccCheckMSOSchemaTemplateExternalepgContractExists(externalepgName string, ss *TemplateExternalepgContract) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs1, err1 := s.RootModule().Resources[externalepgName]

		if !err1 {
			return fmt.Errorf("External Epg Contract %s not found", externalepgName)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Schema id was set")
		}

		cont, err := client.GetViaURL("api/v1/schemas/5ea809672c00003bc40a2799")
		if err != nil {
			return err
		}
		count, err := cont.ArrayCount("templates")
		if err != nil {
			return fmt.Errorf("No Template found")
		}
		tp := TemplateExternalepgContract{}
		found := false
		for i := 0; i < count; i++ {
			tempCont, err := cont.ArrayElement(i, "templates")
			if err != nil {
				return err
			}

			apiTemplateName := models.StripQuotes(tempCont.S("name").String())
			if apiTemplateName == "Template1" {
				externalepgCount, err := tempCont.ArrayCount("externalEpgs")
				if err != nil {
					return fmt.Errorf("Unable to get External Epg list")
				}
				for j := 0; j < externalepgCount; j++ {
					epgCont, err := tempCont.ArrayElement(j, "externalEpgs")
					if err != nil {
						return err
					}
					apiExternalepg := models.StripQuotes(epgCont.S("name").String())
					if apiExternalepg == "UntitledExternalEPG1" {
						contractCount, err := epgCont.ArrayCount("contractRelationships")
						if err != nil {
							return fmt.Errorf("Unable to get contract Relationships list")
						}
						for k := 0; k < contractCount; k++ {
							contractCont, err := epgCont.ArrayElement(k, "contractRelationships")
							if err != nil {

								return err
							}
							contractRef := models.StripQuotes(contractCont.S("contractRef").String())
							re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/contracts/(.*)")
							split := re.FindStringSubmatch(contractRef)
							if "contract9999" == fmt.Sprintf("%s", split[3]) {
								tp.name = fmt.Sprintf("%s", split[3])
								tp.relation = models.StripQuotes(contractCont.S("relationshipType").String())
								found = true
								break
							}
						}
					}
				}
			}
		}
		if !found {
			return fmt.Errorf("External Epg Contract not found from API")
		}

		tp1 := &tp
		*ss = *tp1
		return nil
	}
}

func testAccCheckMSOSchemaTemplateExternalepgContractDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_template_externalepg" {
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
						externalepgCount, err := tempCont.ArrayCount("externalEpgs")
						if err != nil {
							return fmt.Errorf("Unable to get External epg list")
						}
						for j := 0; j < externalepgCount; j++ {
							epgCont, err := tempCont.ArrayElement(j, "externalEpgs")
							if err != nil {
								return err
							}
							apiExternalepg := models.StripQuotes(epgCont.S("name").String())
							if apiExternalepg == "UntitledExternalEPG1" {
								contractCount, err := epgCont.ArrayCount("contractRelationships")
								if err != nil {
									return fmt.Errorf("Unable to get contract Relationships list")
								}
								for k := 0; k < contractCount; k++ {
									contractCont, err := epgCont.ArrayElement(k, "contractRelationships")
									if err != nil {

										return err
									}
									contractRef := models.StripQuotes(contractCont.S("contractRef").String())
									re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/contracts/(.*)")
									split := re.FindStringSubmatch(contractRef)
									if "contract9999" == fmt.Sprintf("%s", split[3]) {
										return fmt.Errorf("External Epg Contract still exists")
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
func testAccCheckMSOSchemaTemplateExternalepgContractAttributes(name string, ss *TemplateExternalepgContract) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if name != ss.relation {
			return fmt.Errorf("Bad Template External epg Contract Relationship Type %s", ss.relation)
		}

		return nil
	}
}

type TemplateExternalepgContract struct {
	name     string
	relation string
}
