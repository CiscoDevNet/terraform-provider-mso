package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccMSOSchemaSite_Basic(t *testing.T) {
	var ss SchemaSiteTest
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSchemaSiteConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteExists("mso_schema.schema1", "mso_schema_site.schemasite1", &ss),
					testAccCheckMSOSchemaSiteAttributes(&ss),
				),
			},
		},
	})
}

func testAccCheckMSOSchemaSiteConfig_basic() string {
	return fmt.Sprintf(`
	resource "mso_schema" "schema1" {
  name          = "shah2"
  template_name = "temp3"
  tenant_id     = "5e9d09482c000068500a269a"

}

resource "mso_schema_site" "schemasite1" {
    schema_id = "${mso_schema.schema1.id}"
    template_name = "temp3"
    site_id = "5c7c95b25100008f01c1ee3c"
}
	

	`)
}

func testAccCheckMSOSchemaSiteExists(schemaName string, schemaSiteName string, ss *SchemaSiteTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs1, err1 := s.RootModule().Resources[schemaName]
		rs2, err2 := s.RootModule().Resources[schemaSiteName]

		if !err1 {
			return fmt.Errorf("Schema %s not found", schemaName)
		}

		if !err2 {
			return fmt.Errorf("Schema site %s not found", schemaSiteName)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Schema id was set")
		}
		if rs2.Primary.ID == "" {
			return fmt.Errorf("No Schema Site id was set")
		}

		cont, err := client.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", rs1.Primary.ID))
		if err != nil {
			return err
		}
		count, err := cont.ArrayCount("sites")
		if err != nil {
			return fmt.Errorf("No Template found")
		}
		tp := SchemaSiteTest{}

		for i := 0; i < count; i++ {
			tempCont, err := cont.ArrayElement(i, "sites")
			if err != nil {
				return err
			}

			apiSiteId := models.StripQuotes(tempCont.S("siteId").String())
			apiTemplate := models.StripQuotes(tempCont.S("templateName").String())

			tp.SchemaId = rs1.Primary.ID
			tp.SiteId = apiSiteId
			tp.TemplateName = apiTemplate

		}
		tp1 := &tp

		*ss = *tp1
		return nil
	}
}

func testAccCheckMSOSchemaSiteDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	rs1, err1 := s.RootModule().Resources["mso_schema.schema1"]

	if !err1 {
		return fmt.Errorf("Schema %s not found", "mso_schema.schema1")
	}

	schemaid := rs1.Primary.ID
	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_site" {
			cont, err := client.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaid))
			if err != nil {
				return nil
			} else {
				count, err := cont.ArrayCount("sites")
				if err != nil {
					return fmt.Errorf("No Template found")
				}
				for i := 0; i < count; i++ {
					tempCont, err := cont.ArrayElement(i, "sites")
					if err != nil {
						return fmt.Errorf("No sites exists")
					}
					apiSiteId := models.StripQuotes(tempCont.S("siteId").String())

					if rs.Primary.ID == apiSiteId {
						return fmt.Errorf("Schema site record still exists")

					}

				}
			}
		} else {

		}
	}
	return nil
}
func testAccCheckMSOSchemaSiteAttributes(ss *SchemaSiteTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "temp3" != ss.TemplateName {
			return fmt.Errorf("Bad Template name %s", ss.TemplateName)
		}
		return nil
	}
}

type SchemaSiteTest struct {
	SchemaId     string `json:",omitempty"`
	SiteId       string `json:",omitempty"`
	TemplateName string `json:",omitempty"`
}
