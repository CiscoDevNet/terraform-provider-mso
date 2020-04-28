package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"

	"github.com/ciscoecosystem/mso-go-client/models"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccMSOUser_Basic(t *testing.T) {
	var s UserTest
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOUserConfig_basic("first"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOUserExists("mso_user.user1", &s),
					testAccCheckMSOUserAttributes("first", &s),
				),
			},
		},
	})
}

func TestAccMSOUser_Update(t *testing.T) {
	var s UserTest

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOUserConfig_basic("first"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOUserExists("mso_user.user1", &s),
					testAccCheckMSOUserAttributes("first", &s),
				),
			},
			{
				Config: testAccCheckMSOUserConfig_basic("first1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOUserExists("mso_user.user1", &s),
					testAccCheckMSOUserAttributes("first1", &s),
				),
			},
		},
	})
}

func testAccCheckMSOUserConfig_basic(name string) string {
	return fmt.Sprintf(`
	resource "mso_user" "user1" {
		username      = "user"
		user_password  = "user@123412341234"
		first_name="%s"
		last_name="last"
	  
		email="email@gmail.com"     
		phone="123456789150"
		account_status="inactive"
	
		  roles{
		  roleid="5ea2bf5a2f0000610b82aa5d"
		  access_type="readWrite"
		
		  }
		
	  }
	`, name)
}

func testAccCheckMSOUserExists(userName string, st *UserTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs1, err1 := s.RootModule().Resources[userName]

		if !err1 {
			return fmt.Errorf("User %s not found", userName)
		}

		if rs1.Primary.ID == "" {
			return fmt.Errorf("No User id was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.GetViaURL("api/v1/users/" + rs1.Primary.ID)

		if err != nil {
			return err
		}

		sts, _ := userFromcontainer(cont)

		*st = *sts
		return nil
	}
}

func userFromcontainer(con *container.Container) (*UserTest, error) {

	s := UserTest{}

	s.User = models.StripQuotes(con.S("username").String())
	s.UserPassword = models.StripQuotes(con.S("password").String())
	if con.Exists("firstName") {
		s.FirstName = models.StripQuotes(con.S("firstName").String())
	}
	if con.Exists("lastName") {
		s.LastName = models.StripQuotes(con.S("lastName").String())
	}
	if con.Exists("emailAddress") {
		s.Email = models.StripQuotes(con.S("emailAddress").String())
	}
	if con.Exists("phoneNumber") {
		s.Phone = models.StripQuotes(con.S("phoneNumber").String())
	}
	if con.Exists("accountStatus") {
		s.AccountStatus = models.StripQuotes(con.S("accountStatus").String())
	}
	if con.Exists("domain") {
		s.Domain = models.StripQuotes(con.S("domain").String())
	}
	count, err := con.ArrayCount("roles")
	if err != nil {
		return nil, fmt.Errorf("No Roles found")
	}

	roles := make([]interface{}, 0)
	for i := 0; i < count; i++ {
		rolesCont, err := con.ArrayElement(i, "roles")

		if err != nil {
			return nil, fmt.Errorf("Unable to parse the roles list")
		}

		map1 := make(map[string]interface{})

		map1["roleid"] = models.StripQuotes(rolesCont.S("roleId").String())
		map1["access_type"] = models.StripQuotes(rolesCont.S("accessType").String())
		roles = append(roles, map1)
	}
	s.Roles = roles

	return &s, nil
}

func testAccCheckMSOUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema" {
			_, err := client.GetViaURL("api/v1/schemas/" + rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Schema still exists")
			}
		} else {
			continue
		}

	}
	return nil
}

func testAccCheckMSOUserAttributes(name string, st *UserTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "last" != st.LastName {
			return fmt.Errorf("Bad Lastname %s", st.LastName)
		}
		if name != st.FirstName {
			return fmt.Errorf("%s Bad Firstname %s", name, st.FirstName)
		}

		return nil
	}
}

type UserTest struct {
	Id           string `json:",omitempty"`
	User         string `json:",omitempty"`
	UserPassword string `json:",omitempty"`

	FirstName string `json:",omitempty"`

	LastName      string        `json:",omitempty"`
	Email         string        `json:",omitempty"`
	Phone         string        `json:",omitempty"`
	AccountStatus string        `json:",omitempty"`
	Domain        string        `json:",omitempty"`
	Roles         []interface{} `json:",omitempty"`
}
