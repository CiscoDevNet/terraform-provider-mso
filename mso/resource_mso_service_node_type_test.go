package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccMsoServiceNodeType_Basic(t *testing.T) {
	var s ServiceNodeTypeTest
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMsoServiceNodeTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMsoServiceNodeTypeConfig_basic("acctest"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMsoServiceNodeTypeExists("mso_service_node_type.node_type", &s),
					testAccCheckMsoServiceNodeTypeAttributes("acctest", &s),
				),
			},
		},
	})
}

func testAccCheckMsoServiceNodeTypeConfig_basic(name string) string {
	return fmt.Sprintf(`
	resource "mso_service_node_type" "node_type" {
		name = "%v"
		display_name = "%v"
	  }
	`, name, name)
}

func testAccCheckMsoServiceNodeTypeExists(tenantName string, st *ServiceNodeTypeTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs1, err1 := s.RootModule().Resources[tenantName]

		if !err1 {
			return fmt.Errorf("Service Node Type record %s not found", tenantName)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Service Node Type record id was set")
		}
		typeId := rs1.Primary.ID
		found := false

		stvt := ServiceNodeTypeTest{}
		msoClient := testAccProvider.Meta().(*client.Client)

		cont, err := msoClient.GetViaURL("api/v1/schemas/service-node-types")
		if err != nil {
			return err
		}

		nodesCount, err := cont.ArrayCount("serviceNodeTypes")
		if err != nil {
			return err
		}

		for i := 0; i < nodesCount; i++ {
			nodeCont, err := cont.ArrayElement(i, "serviceNodeTypes")
			if err != nil {
				return err
			}

			apiId := models.StripQuotes(nodeCont.S("id").String())

			if apiId == typeId {
				stvt.Id = apiId
				stvt.Name = models.StripQuotes(nodeCont.S("name").String())
				stvt.DisplayName = models.StripQuotes(nodeCont.S("displayName").String())
				found = true
			}
		}
		if !found {
			return fmt.Errorf("Unable to find service node type %s", typeId)
		}

		stv := &stvt
		*st = *stv
		return nil
	}
}

func testAccCheckMsoServiceNodeTypeDestroy(s *terraform.State) error {
	msoClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_service_node_type" {
			typeId := rs.Primary.ID
			cont, err := msoClient.GetViaURL("api/v1/schemas/service-node-types")
			if err != nil {
				return err
			}

			nodesCount, err := cont.ArrayCount("serviceNodeTypes")
			if err != nil {
				return err
			}

			for i := 0; i < nodesCount; i++ {
				nodeCont, err := cont.ArrayElement(i, "serviceNodeTypes")
				if err != nil {
					return err
				}

				apiId := models.StripQuotes(nodeCont.S("id").String())

				if apiId == typeId {

					return fmt.Errorf("Service Node Type still exists %s", typeId)

				}
			}

		} else {
			continue
		}
	}
	return nil
}
func testAccCheckMsoServiceNodeTypeAttributes(name string, st *ServiceNodeTypeTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if name != st.Name {
			return fmt.Errorf("Bad Service Type Name %s", st.Name)
		}
		return nil
	}
}

type ServiceNodeTypeTest struct {
	Id          string
	Name        string
	DisplayName string
}
