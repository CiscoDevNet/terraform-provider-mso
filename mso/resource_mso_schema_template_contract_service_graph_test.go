package mso

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccMSOSchemaTemplateContractServiceGraph_Basic(t *testing.T) {
	var instance TemplateContractServiceGraph
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateContractServiceGraphDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateContractServiceGraphConfig_basic("BD1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateContractServiceGraphExists("mso_schema_template_contract_service_graph.sg1", &instance),
					testAccCheckMSOSchemaTemplateContractServiceGraphAttributes("BD1", &instance),
				),
			},
		},
	})
}

func TestAccMSOSchemaTemplateContractServiceGraph_Update(t *testing.T) {
	var instance TemplateContractServiceGraph

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateContractServiceGraphDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateContractServiceGraphConfig_basic("BD1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateContractServiceGraphExists("mso_schema_template_contract_service_graph.sg1", &instance),
					testAccCheckMSOSchemaTemplateContractServiceGraphAttributes("BD1", &instance),
				),
			},
			{
				Config: testAccCheckMSOTemplateContractServiceGraphConfig_basic("BD2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateContractServiceGraphExists("mso_schema_template_contract_service_graph.sg1", &instance),
					testAccCheckMSOSchemaTemplateContractServiceGraphAttributes("BD2", &instance),
				),
			},
		},
	})
}

func testAccCheckMSOTemplateContractServiceGraphConfig_basic(bdName string) string {
	return fmt.Sprintf(`
	resource "mso_schema_template_contract_service_graph" "sg1" {
		schema_id = "5f11b0e22c00001c4a812a2a"
		site_id = "5c7c95b25100008f01c1ee3c"
		template_name = "Template1"
		contract_name = "UntitledContract1"
		service_graph_name = "sg1"
		node_relationship {
		  provider_connector_bd_name = "%s"
		  consumer_connector_bd_name = "BD2"
		  provider_connector_cluster_interface = "test"
		  consumer_connector_cluster_interface = "test"
		}
	}`, bdName)
}

func testAccCheckMSOSchemaTemplateContractServiceGraphExists(ContracGraph string, tc *TemplateContractServiceGraph) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs, error := s.RootModule().Resources[ContracGraph]

		if !error {
			return fmt.Errorf("Contract  Serive Graph %s not found", ContracGraph)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Schema id was set")
		}

		cont, err := client.GetViaURL("api/v1/schemas/5f11b0e22c00001c4a812a2a")
		if err != nil {
			return err
		}
		count, err := cont.ArrayCount("templates")
		if err != nil {
			return fmt.Errorf("No Template found")
		}
		tp := TemplateContractServiceGraph{}
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

					if apiContract == "UntitledContract1" {
						if contractCont.Exists("serviceGraphRelationship") {
							graphRelation := contractCont.S("serviceGraphRelationship")

							graphRef := models.StripQuotes(graphRelation.S("serviceGraphRef").String())
							tokens := strings.Split(graphRef, "/")
							if tokens[len(tokens)-1] == "sg1" {
								tp.Name = tokens[len(tokens)-1]

								nodeCount, _ := graphRelation.ArrayCount("serviceNodesRelationship")
								for k := 0; k < nodeCount; k++ {
									node, err := graphRelation.ArrayElement(k, "serviceNodesRelationship")
									if err != nil {
										return fmt.Errorf("Unable to parse Node relationship for service graph")
									}

									probdRef := models.StripQuotes(node.S("providerConnector", "bdRef").String())
									probdRefTokens := strings.Split(probdRef, "/")
									tp.ProviderBD = probdRefTokens[len(probdRefTokens)-1]

									conbdRef := models.StripQuotes(node.S("consumerConnector", "bdRef").String())
									conbdRefTokens := strings.Split(conbdRef, "/")
									tp.ConsumerBD = conbdRefTokens[len(conbdRefTokens)-1]

									found = true
									break
								}
							}
						}
					}
				}
			}
		}

		if !found {
			return fmt.Errorf("Contract Service Graph not found from API")
		}

		tp1 := &tp

		*tc = *tp1
		return nil
	}
}

func testAccCheckMSOSchemaTemplateContractServiceGraphDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_template_contract_service_graph" {
			cont, err := client.GetViaURL("api/v1/schemas/5f11b0e22c00001c4a812a2a")
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

							if apiContract == "UntitledContract1" {
								if contractCont.Exists("serviceGraphRelationship") {
									graphRelation := contractCont.S("serviceGraphRelationship")

									graphRef := models.StripQuotes(graphRelation.S("serviceGraphRef").String())
									tokens := strings.Split(graphRef, "/")
									name := tokens[len(tokens)-1]

									if name == "sg1" {
										return fmt.Errorf("Contract Service Graph still exists")
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

func testAccCheckMSOSchemaTemplateContractServiceGraphAttributes(bdName string, tc *TemplateContractServiceGraph) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "sg1" != tc.Name {
			return fmt.Errorf("Bad Template Contract Service Graph name %v", tc.Name)
		}

		if bdName != tc.ProviderBD {
			return fmt.Errorf("Bad Template Contract Service Graph Provider BD %v", tc.ProviderBD)
		}

		if "BD2" != tc.ConsumerBD {
			return fmt.Errorf("Bad Template Contract Service Graph Consumer BD %v", tc.ConsumerBD)
		}
		return nil
	}
}

type TemplateContractServiceGraph struct {
	Name       string
	ProviderBD string
	ConsumerBD string
}
