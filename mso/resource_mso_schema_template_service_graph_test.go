package mso

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSchemaTemplateServiceGraph_Basic(t *testing.T) {
	var s SchemaTemplateServiceGraphTest
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMsoSchemaTemplateServiceGraphDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMsoSchemaTemplateServiceGraphConfig_basic("acctest"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMsoSchemaTemplateServiceGraphExists("mso_schema_template_service_graph.test_sg", &s),
					testAccCheckMsoSchemaTemplateServiceGraphAttributes(&s, "acctest"),
				),
			},
		},
	})
}

func TestAccSchemaTemplateServiceGraph_Update(t *testing.T) {
	var s SchemaTemplateServiceGraphTest

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMsoSchemaTemplateServiceGraphDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMsoSchemaTemplateServiceGraphConfig_basic("acctest"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMsoSchemaTemplateServiceGraphExists("mso_schema_template_service_graph.test_sg", &s),
					testAccCheckMsoSchemaTemplateServiceGraphAttributes(&s, "acctest"),
				),
			},
			{
				Config: testAccCheckMsoSchemaTemplateServiceGraphConfig_basic("acctest_update"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMsoSchemaTemplateServiceGraphExists("mso_schema_template_service_graph.test_sg", &s),
					testAccCheckMsoSchemaTemplateServiceGraphAttributes(&s, "acctest_update"),
				),
			},
		},
	})
}

func testAccCheckMsoSchemaTemplateServiceGraphConfig_basic(desc string) string {
	return fmt.Sprintf(`
	resource "mso_schema_template_service_graph" "test_sg" {
		schema_id = "5f06a4c40f0000b63dbbd647"
		template_name = "Template1"
		service_graph_name = "acctestgraph"
		service_node_type = "firewall"
		description = "%s"
		site_nodes  {
			site_id = "5f05c69f1900002234d0537e"
			tenant_name = "NkAutomation"
			node_name = "nk-fw-2"
		}
	
	}
	`, desc)
}

func testAccCheckMsoSchemaTemplateServiceGraphExists(schemaTemplateVrfName string, stvc *SchemaTemplateServiceGraphTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs1, err1 := s.RootModule().Resources[schemaTemplateVrfName]

		if !err1 {
			return fmt.Errorf("Schema Template Service Graph record %s not found", schemaTemplateVrfName)
		}

		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Schema Template Service Graph id was set")
		}

		client := testAccProvider.Meta().(*client.Client)
		cont, err := client.GetViaURL("api/v1/schemas/5f06a4c40f0000b63dbbd647")

		if err != nil {
			return err
		}

		stvt := SchemaTemplateServiceGraphTest{}
		stvt.SchemaId = "5f06a4c40f0000b63dbbd647"

		sgCont, _, err := getTemplateServiceGraphCont(cont, "Template1", "acctestgraph")

		if err != nil {
			return err
		}

		stvt.GraphName = "acctestgraph"
		stvt.Template = "Template1"

		nodeId, err := getNodeIdFromName(client, "firewall")
		if err != nil {
			return err
		}

		_, _, err = getTemplateServiceNodeCont(sgCont, "tfnode1", nodeId)

		if err != nil {
			return err
		}
		stvt.NodeType = "firewall"

		stvt.Description = models.StripQuotes(sgCont.S("description").String())

		graphCont, _, err := getSiteServiceGraphCont(
			cont,
			"5f06a4c40f0000b63dbbd647",
			"Template1",
			"5f05c69f1900002234d0537e",
			"acctestgraph",
		)

		if err != nil {
			return err
		}

		nodeCont, _, err := getSiteServiceNodeCont(
			graphCont,
			"5f06a4c40f0000b63dbbd647",
			"Template1",
			"acctestgraph",
			"tfnode1",
		)

		if err != nil {
			return err
		}

		deviceDn := models.StripQuotes(nodeCont.S("device", "dn").String())

		dnSplit := strings.Split(deviceDn, "/")

		tnName := strings.Join(strings.Split(dnSplit[1], "-")[1:], "-")

		stvt.TenantName = tnName
		stvt.NodeName = strings.Join(strings.Split(dnSplit[2], "-")[1:], "-")
		stvt.SiteId = "5f05c69f1900002234d0537e"

		log.Printf("hiiiiii %v", stvt)
		stv := &stvt
		*stvc = *stv

		return nil
	}
}

func testAccCheckMsoSchemaTemplateServiceGraphDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_template_service_graph" {
			cont, err := client.GetViaURL("api/v1/schemas/5f06a4c40f0000b63dbbd647")
			if err != nil {
				return nil
			}

			_, ind, err := getTemplateServiceGraphCont(cont, "Template1", "acctestgraph")

			if ind != -1 {
				return fmt.Errorf("Service graph still exists")
			}

		}
	}
	return nil
}

func testAccCheckMsoSchemaTemplateServiceGraphAttributes(stvc *SchemaTemplateServiceGraphTest, desc string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "acctestgraph" != stvc.GraphName {
			log.Printf("hjjjjj %v", stvc)
			return fmt.Errorf("Bad Schema Template Service Graph Name %s", stvc.GraphName)
		}

		if desc != stvc.Description {
			return fmt.Errorf("Bad Schema Template Service Graph Description %s", desc)
		}
		return nil
	}
}

type SchemaTemplateServiceGraphTest struct {
	Id          string `json:",omitempty"`
	SchemaId    string `json:",omitempty"`
	Template    string `json:",omitempty"`
	GraphName   string `json:",omitempty"`
	NodeType    string `json:",omitempty"`
	SiteId      string `json:",omitempty"`
	NodeName    string `json:",omitempty"`
	TenantName  string `json:",omitempty"`
	Description string `json:",omitempty"`
}
