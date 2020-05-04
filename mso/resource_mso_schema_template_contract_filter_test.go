package mso

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccMSOSchemaTemplateContractFilter_Basic(t *testing.T) {
	var tc TemplateContractFilter
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateContractFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateContractFilterConfig_basic("Many"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateContractFilterExists("mso_schema_template_contract_filter.filter1", &tc),
					testAccCheckMSOSchemaTemplateContractFilterAttributes("Many", &tc),
				),
			},
		},
	})
}

func TestAccMSOSchemaTemplateContractFilter_Update(t *testing.T) {
	var tc TemplateContractFilter

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateContractFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateContractFilterConfig_basic("Many"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateContractFilterExists("mso_schema_template_contract_filter.filter1", &tc),
					testAccCheckMSOSchemaTemplateContractFilterAttributes("Many", &tc),
				),
			},
			{
				Config: testAccCheckMSOTemplateContractFilterConfig_basic("Any"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateContractFilterExists("mso_schema_template_contract_filter.filter1", &tc),
					testAccCheckMSOSchemaTemplateContractFilterAttributes("Any", &tc),
				),
			},
		},
	})
}

func testAccCheckMSOTemplateContractFilterConfig_basic(filter_type string) string {
	return fmt.Sprintf(`
						resource "mso_schema_template_contract_filter" "filter1" {
							schema_id = "5c4d5bb72700000401f80948"
							template_name = "Template1"
							contract_name = "C300"
							display_name = "C300"
							filter_type = "bothWay"
							scope = "context"
							filter_relationships_procon = {
						   procon_schema_id = "5c4d5bb72700000401f80948"
						   procon_template_name = "Template1"
							procon_name = "%s"
							}
							procon_directives = ["log","none","log"]
							filter_relationships_conpro = {
							    conpro_schema_id = "5c4d5bb72700000401f80948"
                              conpro_template_name = "Template1"
							 conpro_name = "Many"
							}
							conpro_directives = ["log","none"]
						   
						  }
					`, filter_type)
}

func testAccCheckMSOSchemaTemplateContractFilterExists(contractName string, tc *TemplateContractFilter) resource.TestCheckFunc {
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
		tp := TemplateContractFilter{}
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
					
					if apiContract == "C300" {
						tp.display_name = models.StripQuotes(contractCont.S("displayName").String())
						tp.filter_type = models.StripQuotes(contractCont.S("filterType").String())
						tp.scope = models.StripQuotes(contractCont.S("scope").String())
						if contractCont.Exists("filterRelationshipsProviderToConsumer") {
							count, _ := contractCont.ArrayCount("filterRelationshipsProviderToConsumer")

							for i := 0; i < count; i++ {
								filterCont, err := contractCont.ArrayElement(i, "filterRelationshipsProviderToConsumer")
								if err != nil {
									return fmt.Errorf("Unable to parse the filter Relationships Provider to Consumer list")
								}

								if filterCont.Exists("filterRef") {
									filRef := filterCont.S("filterRef").Data()
									split := strings.Split(filRef.(string), "/")
									tp.Name = fmt.Sprintf("%s", split[6])
								}
							}
						}
						

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

func testAccCheckMSOSchemaTemplateContractFilterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_template_contract_filter" {
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
							if apiContract == "C300" {
								return fmt.Errorf("template contract Filter still exists.")
							}
						}
					}
				}
			}
		}
	}
	return nil
}
func testAccCheckMSOSchemaTemplateContractFilterAttributes(filter_type string, tc *TemplateContractFilter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if filter_type != tc.Name {
			return fmt.Errorf("Bad Template Contract filter_type %v", tc.Name)
		}

		if "C300" != tc.display_name {
			return fmt.Errorf("Bad Template Contract display name %s", tc.display_name)
		}

		if "context" != tc.scope {
			return fmt.Errorf("Bad Template Contract Scope name %s", tc.scope)
		}
		return nil
	}
}

type TemplateContractFilter struct {
	display_name string
	scope        string
	filter_type  string
	Name         string
}
