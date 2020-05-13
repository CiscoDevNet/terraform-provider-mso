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

func TestAccMSOSchema_Basic(t *testing.T) {
	var s SchemaTest
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSchemaConfig_basic("nkp1003"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaExists("mso_schema.schema1", &s),
					testAccCheckMSOSchemaAttributes("nkp1003", &s),
				),
			},
		},
	})
}

func TestAccMSOSchema_Update(t *testing.T) {
	var s SchemaTest

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSchemaConfig_basic("nkp1003"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaExists("mso_schema.schema1", &s),
					testAccCheckMSOSchemaAttributes("nkp1003", &s),
				),
			},
			{
				Config: testAccCheckMSOSchemaConfig_basic("nkp1004"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaExists("mso_schema.schema1", &s),
					testAccCheckMSOSchemaAttributes("nkp1004", &s),
				),
			},
		},
	})
}

func testAccCheckMSOSchemaConfig_basic(name string) string {
	return fmt.Sprintf(`
	resource "mso_schema" "schema1" {
		name          = "nkp1003"
		template_name = "temp1"
		tenant_id     = "5e9d09482c000068500a269a"

	  }
	`)
}

func testAccCheckMSOSchemaExists(schemaName string, st *SchemaTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs1, err1 := s.RootModule().Resources[schemaName]

		if !err1 {
			return fmt.Errorf("Schema %s not found", schemaName)
		}

		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Schema id was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.GetViaURL("api/v1/schemas/" + rs1.Primary.ID)

		if err != nil {
			return err
		}

		sts, _ := schemaFromcontainer(cont)

		*st = *sts
		return nil
	}
}

func schemaFromcontainer(con *container.Container) (*SchemaTest, error) {

	s := SchemaTest{}

	s.Name = models.StripQuotes(con.S("displayName").String())
	count, err := con.ArrayCount("templates")
	if err != nil {
		return nil, fmt.Errorf("No Template found")
	}

	for i := 0; i < count; i++ {
		tempCont, err := con.ArrayElement(i, "templates")

		if err != nil {
			return nil, fmt.Errorf("Unable to parse the template list")
		}
		s.TemplateName = models.StripQuotes(tempCont.S("name").String())
		s.TenantId = models.StripQuotes(tempCont.S("tenantId").String())

	}

	return &s, nil
}

func testAccCheckMSOSchemaDestroy(s *terraform.State) error {
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

func testAccCheckMSOSchemaAttributes(name string, st *SchemaTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "temp1" != st.TemplateName {
			return fmt.Errorf("Bad Template name %s", st.TemplateName)
		}

		return nil
	}
}

type SchemaTest struct {
	Id   string `json:",omitempty"`
	Name string `json:",omitempty"`

	TemplateName string `json:",omitempty"`

	TenantId string `json:",omitempty"`
}
