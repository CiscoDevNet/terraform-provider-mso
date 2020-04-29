package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccMSOSchemaTemplate_Basic(t *testing.T) {
	var ss SchemaTemplateTest
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSchemaTemplateConfig_basic("Temp1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateExists("mso_schema.schema1", "mso_schema_template.sample1", &ss),
					testAccCheckMSOSchemaTemplateAttributes("Temp1", &ss),
				),
			},
		},
	})
}
func TestAccMSOSchemaTemplate_Update(t *testing.T) {
	var ss SchemaTemplateTest

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSchemaTemplateConfig_basic("Temp1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateExists("mso_schema.schema1", "mso_schema_template.sample1", &ss),
					testAccCheckMSOSchemaTemplateAttributes("Temp1", &ss),
				),
			},
			{
				Config: testAccCheckMSOSchemaTemplateConfig_basic("Temp2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateExists("mso_schema.schema1", "mso_schema_template.sample1", &ss),
					testAccCheckMSOSchemaTemplateAttributes("Temp2", &ss),
				),
			},
		},
	})
}

func testAccCheckMSOSchemaTemplateConfig_basic(displayName string) string {
	return fmt.Sprintf(`
	resource "mso_schema" "schema1" {
  name          = "shah2"
  template_name = "temp3"
  tenant_id     = "5e9d09482c000068500a269a"

}

resource "mso_schema_template" "sample1" {
  schema_id = "${mso_schema.schema1.id}"
  name = "Temp200"
  display_name = "%v"
  tenant_id = "5c4d9f3d2700007e01f80949"
  
}`, displayName)
}

func testAccCheckMSOSchemaTemplateExists(schemaName string, schemaTemplateName string, ss *SchemaTemplateTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs1, err1 := s.RootModule().Resources[schemaName]
		rs2, err2 := s.RootModule().Resources[schemaTemplateName]

		if !err1 {
			return fmt.Errorf("Schema %s not found", schemaName)
		}

		if !err2 {
			return fmt.Errorf("Schema template %s not found", schemaTemplateName)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Schema id was set")
		}
		if rs2.Primary.ID == "" {
			return fmt.Errorf("No Schema Template id was set")
		}

		cont, err := client.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", rs1.Primary.ID))
		if err != nil {
			return err
		}
		count, err := cont.ArrayCount("templates")
		if err != nil {
			return fmt.Errorf("No Template found")
		}
		tp := SchemaTemplateTest{}

		for i := 0; i < count; i++ {
			tempCont, err := cont.ArrayElement(i, "templates")
			if err != nil {
				return err
			}

			apiTenantId := models.StripQuotes(tempCont.S("tenantId").String())
			apiTemplateName := models.StripQuotes(tempCont.S("name").String())
			apiTemplateDisplayName := models.StripQuotes(tempCont.S("displayName").String())

			tp.SchemaId = rs1.Primary.ID
			tp.TenantId = apiTenantId
			tp.Name = apiTemplateName
			tp.DisplayName = apiTemplateDisplayName

		}
		tp1 := &tp

		*ss = *tp1
		return nil
	}
}

func testAccCheckMSOSchemaTemplateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	rs1, err1 := s.RootModule().Resources["mso_schema.schema1"]

	if !err1 {
		return fmt.Errorf("Schema %s not found", "mso_schema.schema1")
	}

	schemaid := rs1.Primary.ID
	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_template" {
			cont, err := client.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaid))
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
						return fmt.Errorf("No sites exists")
					}
					apiTemplateId := models.StripQuotes(tempCont.S("name").String())

					if rs.Primary.ID == apiTemplateId {
						return fmt.Errorf("Schema template record still exists")

					}

				}
			}
		}
	}
	return nil
}
func testAccCheckMSOSchemaTemplateAttributes(displayName string, ss *SchemaTemplateTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if displayName != ss.DisplayName {
			return fmt.Errorf("Bad Template display name %s", ss.DisplayName)
		}
		return nil
	}
}

type SchemaTemplateTest struct {
	SchemaId    string `json:",omitempty"`
	Name        string `json:",omitempty"`
	DisplayName string `json:",omitempty"`
	TenantId    string `json:",omitempty"`
}
