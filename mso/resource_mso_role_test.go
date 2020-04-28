package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccMSORole_Basic(t *testing.T) {
	var ss RoleTest
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSORoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSORoleConfig_basic("UserManager100"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSORoleExists("mso_role.role1", &ss),
					testAccCheckMSORoleAttributes("UserManager100", &ss),
				),
			},
		},
	})
}

func TestAccMSORole_Update(t *testing.T) {
	var ss RoleTest

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSORoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSORoleConfig_basic("UserManager100"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSORoleExists("mso_role.role1", &ss),
					testAccCheckMSORoleAttributes("UserManager100", &ss),
				),
			},
			{
				Config: testAccCheckMSORoleConfig_basic("UserManager101"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSORoleExists("mso_role.role1", &ss),
					testAccCheckMSORoleAttributes("UserManager101", &ss),
				),
			},
		},
	})
}

func testAccCheckMSORoleConfig_basic(name string) string {
	return fmt.Sprintf(`
	resource "mso_role" "role1" {
  	name = "%v"
  	display_name = "UserManager"
 	description = "helloo"
  	read_permissions = ["view-sites"]
  	write_permissions = ["manage-sites","manage-tenants"]
  
}`, name)

}

func testAccCheckMSORoleExists(roleName string, ss *RoleTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs1, err1 := s.RootModule().Resources[roleName]

		if !err1 {
			return fmt.Errorf("Role %s not found", roleName)
		}

		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Role id was set")
		}

		cont, err := client.GetViaURL(fmt.Sprintf("api/v1/roles/%s", rs1.Primary.ID))
		if err != nil {
			return err
		}

		tp := RoleTest{}

		tp.Id = models.StripQuotes(cont.S("id").String())
		tp.Name = models.StripQuotes(cont.S("name").String())
		tp.DisplayName = models.StripQuotes(cont.S("displayName").String())
		tp.Description = models.StripQuotes(cont.S("description").String())
		tp.ReadPermissions = cont.S("readPermissions").Data().([]interface{})
		tp.WritePermissions = cont.S("writePermissions").Data().([]interface{})

		tp1 := &tp

		*ss = *tp1
		return nil
	}
}

func testAccCheckMSORoleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "mso_role" {
			_, err := client.GetViaURL(fmt.Sprintf("api/v1/roles/%s", rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Role still exists")
			}
		} else {
		}
	}
	return nil
}

func testAccCheckMSORoleAttributes(name string, ss *RoleTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if name != ss.Name {
			return fmt.Errorf("Bad Role name %s", ss.Name)
		}
		return nil
	}
}

type RoleTest struct {
	Id          string `json:",omitempty"`
	Name        string `json:",omitempty`
	DisplayName string `json:",omitempty"`
	Description string `json:",omitempty"`

	ReadPermissions []interface{} `json:",omitempty"`

	WritePermissions []interface{} `json:",omitempty"`
}
