package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccMSOSchemaTemplateContract_Basic(t *testing.T) {
	var tc TemplateContract
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateContractConfig_basic("bothWay"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateContractExists("mso_schema_template_contract.template_contract", &tc),
					testAccCheckMSOSchemaTemplateContractAttributes("bothWay", &tc),
				),
			},
		},
	})
}

func TestAccMSOSchemaTemplateContract_Update(t *testing.T) {
	var tc TemplateContract

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateContractConfig_basic("bothWay"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateContractExists("mso_schema_template_contract.template_contract", &tc),
					testAccCheckMSOSchemaTemplateContractAttributes("bothWay", &tc),
				),
			},
			{
				Config: testAccCheckMSOTemplateContractConfig_basic("oneWay"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateContractExists("mso_schema_template_contract.template_contract", &tc),
					testAccCheckMSOSchemaTemplateContractAttributes("oneWay", &tc),
				),
			},
		},
	})
}

func testAccCheckMSOTemplateContractConfig_basic(filter_type string) string {
	return fmt.Sprintf(`
	resource "mso_schema_template_contract" "template_contract" {
		schema_id = "5c4d5bb72700000401f80948"
		template_name = "Template1"
		contract_name = "C1"
		display_name = "C1"
		filter_type = "%v"
		scope = "context"
		filter_relationships = {
			filter_schema_id = "5c4d5bb72700000401f80948"
    		filter_template_name = "Template1"
		  	filter_name = "Any"
		}
		directives = ["log"]
	  }
`, filter_type)
}

func testAccCheckMSOSchemaTemplateContractExists(contractName string, tc *TemplateContract) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs, error := s.RootModule().Resources[contractName]

		if !error {
			return fmt.Errorf("Contract %s not found", contractName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Schema id was set")
		}

		cont, err := client.GetViaURL("api/v1/schemas/5c4d5bb72700000401f80948")
		if err != nil {
			return err
		}
		count, err := cont.ArrayCount("templates")
		if err != nil {
			return fmt.Errorf("No Template found")
		}
		tp := TemplateContract{}
		found := false
		for i := 0; i < count; i++ {
			tempCont, err := cont.ArrayElement(i, "templates")
			if err != nil {
				return err
			}

			apiTemplateName := models.StripQuotes(tempCont.S("name").String())
			if apiTemplateName == "Template1" {
				contractCount, err := tempCont.ArrayCount("contracts")
				if err != nil {
					return fmt.Errorf("Unable to get Contract list")
				}
				for j := 0; j < contractCount; j++ {
					contractCont, err := tempCont.ArrayElement(j, "contracts")
					if err != nil {
						return err
					}
					apiContract := models.StripQuotes(contractCont.S("name").String())
					if apiContract == "C1" {
						tp.display_name = models.StripQuotes(contractCont.S("displayName").String())
						tp.filter_type = models.StripQuotes(contractCont.S("filterType").String())
						tp.scope = models.StripQuotes(contractCont.S("scope").String())

						found = true
						break

					}
				}
			}
		}

		if !found {
			return fmt.Errorf("Contract not found from API")
		}

		tp1 := &tp

		*tc = *tp1
		return nil
	}
}

func testAccCheckMSOSchemaTemplateContractDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_template_contract" {
			cont, err := client.GetViaURL("api/v1/schemas/5c4d5bb72700000401f80948")
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
						contractCount, err := tempCont.ArrayCount("contracts")
						if err != nil {
							return fmt.Errorf("Unable to get Contract list")
						}
						for j := 0; j < contractCount; j++ {
							contractCont, err := tempCont.ArrayElement(j, "contracts")
							if err != nil {
								return err
							}
							apiContract := models.StripQuotes(contractCont.S("name").String())
							if apiContract == "C1" {
								return fmt.Errorf("template contract still exists.")
							}
						}
					}
				}
			}
		}
	}
	return nil
}
func testAccCheckMSOSchemaTemplateContractAttributes(filter_type string, tc *TemplateContract) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if filter_type != tc.filter_type {
			return fmt.Errorf("Bad Template Contract filter_type %v", tc.filter_type)
		}

		if "C1" != tc.display_name {
			return fmt.Errorf("Bad Template Contract display name %s", tc.display_name)
		}

		if "context" != tc.scope {
			return fmt.Errorf("Bad Template Contract Scope name %s", tc.scope)
		}
		return nil
	}
}

type TemplateContract struct {
	display_name string
	scope        string
	filter_type  string
}
