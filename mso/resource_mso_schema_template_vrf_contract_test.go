package mso

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSchemaTemplateVrfContract_Basic(t *testing.T) {
	var s SchemaTemplateVrfContractTest
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMsoSchemaTemplateVrfContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMsoSchemaTemplateVrfContractConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMsoSchemaTemplateVrfContractExists("mso_schema_template_vrf_contract.acc_vrf", &s),
					testAccCheckMsoSchemaTemplateVrfContractAttributes(&s),
				),
			},
		},
	})
}

func testAccCheckMsoSchemaTemplateVrfContractConfig_basic() string {
	return fmt.Sprintf(`

	resource "mso_schema_template_vrf_contract" "acc_vrf" {
		schema_id              = "5eff091b0e00008318cff859"
		template_name          = "Template1"
		vrf_name               = "myVrf"
		relationship_type      = "provider"
		contract_name          = "hello"
	  }
	`)
}

func testAccCheckMsoSchemaTemplateVrfContractExists(schemaTemplateVrfName string, stvc *SchemaTemplateVrfContractTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs1, err1 := s.RootModule().Resources[schemaTemplateVrfName]

		if !err1 {
			return fmt.Errorf("Schema Template Vrf Contract record %s not found", schemaTemplateVrfName)
		}

		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Schema Template Vrf Contract id was set")
		}

		client := testAccProvider.Meta().(*client.Client)
		con, err := client.GetViaURL("api/v1/schemas/5eff091b0e00008318cff859")

		if err != nil {
			return err
		}

		stvt := SchemaTemplateVrfContractTest{}
		stvt.SchemaId = rs1.Primary.ID

		count, err := con.ArrayCount("templates")
		if err != nil {
			return err
		}
		found := false
		for i := 0; i < count; i++ {
			tempCont, err := con.ArrayElement(i, "templates")
			stvt.Template = models.StripQuotes(tempCont.S("name").String())
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
				if apiVRF == "myVrf" {
					stvt.VrfName = "myVrf"
					contractCount, err := vrfCont.ArrayCount(humanToApiType["provider"])
					if err != nil {
						return fmt.Errorf("Unable to get contract Relationships list")
					}
					for k := 0; k < contractCount; k++ {
						contractCont, err := vrfCont.ArrayElement(k, humanToApiType["provider"])
						if err != nil {
							return err
						}
						contractRef := models.StripQuotes(contractCont.S("contractRef").String())
						re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/contracts/(.*)")
						split := re.FindStringSubmatch(contractRef)
						if contractRef != "{}" && contractRef != "" {
							if "hello" == fmt.Sprintf("%s", split[3]) {
								stvt.ContractName = "hello"
								stvt.RelationType = "provider"
								found = true
								break
							}
						}
					}
				}
			}
		}

		if !found {
			return fmt.Errorf("Unable to get contract list")
		}

		log.Printf("hiiiiii %v", stvt)
		stv := &stvt
		*stvc = *stv

		return nil
	}
}

func testAccCheckMsoSchemaTemplateVrfContractDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_template_vrf" {
			con, err := client.GetViaURL("api/v1/schemas/5eff091b0e00008318cff859")
			if err != nil {
				return nil
			} else {
				count, err := con.ArrayCount("templates")
				if err != nil {
					return fmt.Errorf("No Template found")
				}
				for i := 0; i < count; i++ {
					tempCont, err := con.ArrayElement(i, "templates")
					if err != nil {
						return fmt.Errorf("No template exists")
					}
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
						if apiVRF == "myVrf" {
							contractCount, err := vrfCont.ArrayCount(humanToApiType["provider"])
							if err != nil {
								return fmt.Errorf("Unable to get contract Relationships list")
							}
							for k := 0; k < contractCount; k++ {
								contractCont, err := vrfCont.ArrayElement(k, humanToApiType["provider"])
								if err != nil {
									return err
								}
								contractRef := models.StripQuotes(contractCont.S("contractRef").String())
								re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/contracts/(.*)")
								split := re.FindStringSubmatch(contractRef)
								if contractRef != "{}" && contractRef != "" {
									if "hello" == fmt.Sprintf("%s", split[3]) {
										return fmt.Errorf("VRF contract still exist.")
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

func testAccCheckMsoSchemaTemplateVrfContractAttributes(stvc *SchemaTemplateVrfContractTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "hello" != stvc.ContractName {
			log.Printf("hjjjjj %v", stvc)
			return fmt.Errorf("Bad Schema Template Vrf Contract Name %s", stvc.ContractName)
		}
		return nil
	}
}

type SchemaTemplateVrfContractTest struct {
	Id           string `json:",omitempty"`
	SchemaId     string `json:",omitempty"`
	Template     string `json:",omitempty"`
	VrfName      string `json:",omitempty"`
	ContractName string `json:",omitempty"`
	RelationType string `json:",omitempty"`
}
