package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTenant_Basic(t *testing.T) {
	var s TenantTest
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMsoTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMsoTenantConfig_basic("Mso123"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMsoTenantExists("mso_tenant.tenant1", &s),
					testAccCheckMsoTenantAttributes("Mso123", &s),
				),
			},
		},
	})
}

func TestAccMsoTenant_Update(t *testing.T) {
	var s TenantTest

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMsoTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMsoTenantConfig_basic("Mso123"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMsoTenantExists("mso_tenant.tenant1", &s),
					testAccCheckMsoTenantAttributes("Mso123", &s),
				),
			},
			{
				Config: testAccCheckMsoTenantConfig_basic("Mso234"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMsoTenantExists("mso_tenant.tenant1", &s),
					testAccCheckMsoTenantAttributes("Mso234", &s),
				),
			},
		},
	})
}

func testAccCheckMsoTenantConfig_basic(name string) string {
	return fmt.Sprintf(`
	resource "mso_tenant" "tenant1" {
		name = "%v"
		display_name = "%v"
	  }
	`, name, name)
}

func testAccCheckMsoTenantExists(tenantName string, st *TenantTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs1, err1 := s.RootModule().Resources[tenantName]

		if !err1 {
			return fmt.Errorf("Tenant record %s not found", tenantName)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Tenant record id was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		resp, err := client.GetViaURL("api/v1/tenants/" + rs1.Primary.ID)

		if err != nil {
			return err
		}

		tp, _ := tenantfromcontainer(resp)

		*st = *tp
		return nil
	}
}

func tenantfromcontainer(con *container.Container) (*TenantTest, error) {

	s := TenantTest{}
	s.Name = models.StripQuotes(con.S("name").String())
	s.DisplayName = models.StripQuotes(con.S("display_name").String())
	s.Description = models.StripQuotes(con.S("description").String())

	return &s, nil
}

func testAccCheckMsoTenantDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_tenant" {
			_, err := client.GetViaURL("api/v1/tenants/" + rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Tenant still exists")
			}
		} else {
			continue
		}
	}
	return nil
}
func testAccCheckMsoTenantAttributes(name string, st *TenantTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if name != st.Name {
			return fmt.Errorf("Bad Tenant Name %s", st.Name)
		}
		return nil
	}
}

type TenantTest struct {
	Id          string        `json:",omitempty"`
	Name        string        `json:",omitempty"`
	DisplayName string        `json:",omitempty"`
	Description string        `json:",omitempty"`
	Users       []interface{} `json:",omitempty"`
	Sites       []interface{} `json:",omitempty"`
}
