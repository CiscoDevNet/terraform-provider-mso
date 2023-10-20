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

func TestAccSchemaSiteServiceGraphNode_Basic(t *testing.T) {
	var s SchemaTemplateServiceGraphTest
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMsoSchemaSiteServiceGraphNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMsoSchemaSiteServiceGraphNodeConfig_basic("nk-fw-2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMsoSchemaSiteServiceGraphNodeExists("mso_schema_site_service_graph_node.test_sg", &s),
					testAccCheckMsoSchemaSiteTemplateServiceGraphNodeAttributes(&s, "nk-fw-2"),
				),
			},
		},
	})
}

func TestAccSchemaSiteServiceGraphNode_Update(t *testing.T) {
	var s SchemaTemplateServiceGraphTest

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMsoSchemaSiteServiceGraphNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMsoSchemaSiteServiceGraphNodeConfig_basic("nk-fw-2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMsoSchemaSiteServiceGraphNodeExists("mso_schema_site_service_graph_node.test_sg", &s),
					testAccCheckMsoSchemaSiteTemplateServiceGraphNodeAttributes(&s, "nk-fw-2"),
				),
			},
			{
				Config: testAccCheckMsoSchemaSiteServiceGraphNodeConfig_basic("nk-fw-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMsoSchemaSiteServiceGraphNodeExists("mso_schema_site_service_graph_node.test_sg", &s),
					testAccCheckMsoSchemaSiteTemplateServiceGraphNodeAttributes(&s, "nk-fw-1"),
				),
			},
		},
	})
}

func testAccCheckMsoSchemaSiteServiceGraphNodeConfig_basic(desc string) string {
	return fmt.Sprintf(`
	resource "mso_schema_site_service_graph_node" "test_sg" {
		schema_id = "5f06a4c40f0000b63dbbd647"
		template_name = "Template1"
		service_graph_name = "sgtf"
		service_node_type = "firewall"
		site_nodes  {
			site_id = "5f05c69f1900002234d0537e"
			tenant_name = "NkAutomation"
			node_name = "%s"
		}
	
	}
	`, desc)
}

func testAccCheckMsoSchemaSiteServiceGraphNodeExists(schemaTemplateVrfName string, stvc *SchemaTemplateServiceGraphTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs1, err1 := s.RootModule().Resources[schemaTemplateVrfName]
		nodeTf := rs1.Primary.ID
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

		sgCont, _, err := getTemplateServiceGraphCont(cont, "Template1", "sgtf")

		if err != nil {
			return err
		}

		stvt.GraphName = "sgtf"
		stvt.Template = "Template1"

		nodesCount, err := cont.ArrayCount("serviceNodeTypes")
		if err != nil {
			return err
		}

		nodeId, err := getNodeIdFromName(cont, nodesCount, "firewall")
		if err != nil {
			return err
		}

		_, _, err = getTemplateServiceNodeCont(sgCont, nodeTf, nodeId)

		if err != nil {
			return err
		}
		stvt.NodeType = "firewall"

		graphCont, _, err := getSiteServiceGraphCont(
			cont,
			"5f06a4c40f0000b63dbbd647",
			"Template1",
			"5f05c69f1900002234d0537e",
			"sgtf",
		)

		if err != nil {
			return err
		}

		nodeCont, _, err := getSiteServiceNodeCont(
			graphCont,
			"5f06a4c40f0000b63dbbd647",
			"Template1",
			"sgtf",
			nodeTf,
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

func testAccCheckMsoSchemaSiteServiceGraphNodeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_site_service_graph_node" {
			cont, err := client.GetViaURL("api/v1/schemas/5f06a4c40f0000b63dbbd647")
			if err != nil {
				return nil
			}
			nodeTf := rs.Primary.ID
			graphCont, _, err := getSiteServiceGraphCont(
				cont,
				"5f06a4c40f0000b63dbbd647",
				"Template1",
				"5f05c69f1900002234d0537e",
				"sgtf",
			)

			if err != nil {
				return nil
			}

			_, ind, err := getSiteServiceNodeCont(
				graphCont,
				"5f06a4c40f0000b63dbbd647",
				"Template1",
				"sgtf",
				nodeTf,
			)

			if ind != -1 {
				return fmt.Errorf("Service graph Node still exists")
			}

		}
	}
	return nil
}

func testAccCheckMsoSchemaSiteTemplateServiceGraphNodeAttributes(stvc *SchemaTemplateServiceGraphTest, desc string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "sgtf" != stvc.GraphName {
			log.Printf("hjjjjj %v", stvc)
			return fmt.Errorf("Bad Schema Site Service Graph  Name %s", stvc.GraphName)
		}

		if desc != stvc.NodeName {
			return fmt.Errorf("Bad Schema Site Service Graph Node Name %s", desc)
		}
		return nil
	}
}
