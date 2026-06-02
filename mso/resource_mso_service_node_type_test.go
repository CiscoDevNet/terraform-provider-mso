package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var msoServiceNodeTypeId string

func TestAccMSOServiceNodeTypeResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOServiceNodeTypeDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create Service Node Type") },
				Config:    testAccMSOServiceNodeTypeConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_service_node_type.test", "name", msoServiceNodeTypeName),
					// display_name defaults to name when not explicitly set
					resource.TestCheckResourceAttr("mso_service_node_type.test", "display_name", msoServiceNodeTypeName),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["mso_service_node_type.test"]
						if !ok {
							return fmt.Errorf("mso_service_node_type.test not found in state")
						}
						msoServiceNodeTypeId = rs.Primary.ID
						return nil
					},
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import Service Node Type") },
				ResourceName: "mso_service_node_type.test",
				ImportState:  true,
				// The importer looks up by name, so we must pass the name as the import ID.
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["mso_service_node_type.test"]
					if !ok {
						return "", fmt.Errorf("mso_service_node_type.test not found in state")
					}
					return rs.Primary.Attributes["name"], nil
				},
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate Service Node Type with explicit display name after manual deletion")
					c := testAccProvider.Meta().(*client.Client)
					if err := c.DeletebyId("api/v1/schemas/service-node-types/" + msoServiceNodeTypeId); err != nil {
						t.Fatalf("Failed to delete service node type %s via API: %s", msoServiceNodeTypeId, err)
					}
				},
				Config: testAccMSOServiceNodeTypeConfigCreateWithDisplayName(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_service_node_type.test", "name", msoServiceNodeTypeName),
					resource.TestCheckResourceAttr("mso_service_node_type.test", "display_name", msoServiceNodeTypeName+" display"),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import Service Node Type with explicit display name") },
				ResourceName: "mso_service_node_type.test",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["mso_service_node_type.test"]
					if !ok {
						return "", fmt.Errorf("mso_service_node_type.test not found in state")
					}
					return rs.Primary.Attributes["name"], nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMSOServiceNodeTypeDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*client.Client)
	cont, err := c.GetViaURL("api/v1/schemas/service-node-types")
	if err != nil {
		return err
	}
	nodesCount, err := cont.ArrayCount("serviceNodeTypes")
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mso_service_node_type" {
			continue
		}
		for i := 0; i < nodesCount; i++ {
			nodeCont, err := cont.ArrayElement(i, "serviceNodeTypes")
			if err != nil {
				return err
			}
			if models.StripQuotes(nodeCont.S("id").String()) == rs.Primary.ID {
				return fmt.Errorf("Service Node Type %s still exists", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccMSOServiceNodeTypeConfigCreate() string {
	return fmt.Sprintf(`
resource "mso_service_node_type" "test" {
  name = "%s"
}
`, msoServiceNodeTypeName)
}

func testAccMSOServiceNodeTypeConfigCreateWithDisplayName() string {
	return fmt.Sprintf(`
resource "mso_service_node_type" "test" {
  name         = "%s"
  display_name = "%s display"
}
`, msoServiceNodeTypeName, msoServiceNodeTypeName)
}
