package mso

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSchemaTemplateVrf_Basic(t *testing.T) {
	var s SchemaTemplateVrfTest
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMsoSchemaTemplateVrfDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMsoSchemaTemplateVrfConfig_basic("vrf982"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMsoSchemaTemplateVrfExists("mso_schema.schema1", "mso_schema_template_vrf.vrf1", &s),
					testAccCheckMsoSchemaTemplateVrfAttributes("vrf982", &s),
				),
			},
		},
	})
}

func TestAccMsoSchemaTemplateVrf_Update(t *testing.T) {
	var s SchemaTemplateVrfTest

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMsoSchemaTemplateVrfDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMsoSchemaTemplateVrfConfig_basic("vrf982"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMsoSchemaTemplateVrfExists("mso_schema.schema1", "mso_schema_template_vrf.vrf1", &s),
					testAccCheckMsoSchemaTemplateVrfAttributes("vrf982", &s),
				),
			},
			{
				Config: testAccCheckMsoSchemaTemplateVrfConfig_basic("vrf983"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMsoSchemaTemplateVrfExists("mso_schema.schema1", "mso_schema_template_vrf.vrf1", &s),
					testAccCheckMsoSchemaTemplateVrfAttributes("vrf983", &s),
				),
			},
		},
	})
}

func testAccCheckMsoSchemaTemplateVrfConfig_basic(name string) string {
	return fmt.Sprintf(`

	resource "mso_schema" "schema1" {
		name          = "shah8"
		template_name = "temp5"
		tenant_id     = "5e9d09482c000068500a269a"
	  
	  }

	resource "mso_schema_template_vrf" "vrf1" {
		schema_id="${mso_schema.schema1.id}"
		template="temp5"
		name= "vrf982"
		display_name="%s"
		layer3_multicast=true
	  }
	`, name)
}

func testAccCheckMsoSchemaTemplateVrfExists(schemaName string, schemaTemplateVrfName string, stv *SchemaTemplateVrfTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs1, err1 := s.RootModule().Resources[schemaName]
		rs2, err2 := s.RootModule().Resources[schemaTemplateVrfName]

		if !err1 {
			return fmt.Errorf("Schema %s not found", schemaName)
		}
		if !err2 {
			return fmt.Errorf("Schema Template Vrf record %s not found", schemaTemplateVrfName)
		}

		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Schema id was set")
		}
		if rs2.Primary.ID == "" {
			return fmt.Errorf("No Schema Template Vrf id was set")
		}

		client := testAccProvider.Meta().(*client.Client)
		con, err := client.GetViaURL("api/v1/schemas/" + rs1.Primary.ID)

		if err != nil {
			return err
		}

		
		stvt := SchemaTemplateVrfTest{}
		stvt.SchemaId = rs1.Primary.ID
		
		count, err := con.ArrayCount("templates")
		if err != nil {
			return err
		}

		for i := 0; i < count; i++ {
			tempCont, err := con.ArrayElement(i, "templates")
			stvt.Template = models.StripQuotes(tempCont.S("name").String())
			vrfCount, err := tempCont.ArrayCount("vrfs")
			if err != nil {
				return fmt.Errorf("No Vrf found")
			}
			for j := 0; j < vrfCount; j++ {
				vrfCont, err := tempCont.ArrayElement(j, "vrfs")
				if err != nil {
					return err
				}
				if vrfCont.Exists("name") {
					stvt.Name = models.StripQuotes(vrfCont.S("name").String())

				}

				if vrfCont.Exists("displayName") {
					stvt.DisplayName = models.StripQuotes(vrfCont.S("displayName").String())
				}

				if vrfCont.Exists("l3MCast") {
					l3Mcast, _ := strconv.ParseBool(models.StripQuotes(vrfCont.S("l3MCast").String()))
					stvt.Layer3Multicast = l3Mcast
				}
			}
		}

		stvc := &stvt

		*stv = *stvc
		return nil
	}
}

func testAccCheckMsoSchemaTemplateVrfDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	rs1, err1 := s.RootModule().Resources["mso_schema.schema1"]

	if !err1 {
		return fmt.Errorf("Schema %s not found", "mso_schema.schema1")
	}
	
	schemaid := rs1.Primary.ID
	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_template_vrf" {
			con, err := client.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaid))
			if err != nil {
				return nil
			} else {
				count, err := con.ArrayCount("templates")
				if err != nil {
					return fmt.Errorf("No Template found")
				}
				for i := 0; i < count; i++ {
					tempCont, err := con.ArrayElement(i, "templates")
					if err != nil {
						return fmt.Errorf("No template exists")
					}
					vrfCount, err := tempCont.ArrayCount("vrfs")
					if err != nil {
						return fmt.Errorf("No Vrf found")
					}
					for j := 0; j < vrfCount; j++ {
						vrfCont, err := tempCont.ArrayElement(j, "vrfs")
						if err != nil {
							return err
						}
						name := models.StripQuotes(vrfCont.S("name").String())

						if rs.Primary.ID == name {
							return fmt.Errorf("Schema Template Vrf record still exists")

						}

					}
				}
			}
		}
	}
	return nil
}

func testAccCheckMsoSchemaTemplateVrfAttributes(name string, stv *SchemaTemplateVrfTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if name != stv.DisplayName {
			return fmt.Errorf("Bad Schema Template Vrf Name %s",stv.DisplayName)
		}
		return nil
	}
}

type SchemaTemplateVrfTest struct {
	Id              string `json:",omitempty"`
	SchemaId        string `json:",omitempty"`
	Template        string `json:",omitempty"`
	Name            string `json:",omitempty"`
	DisplayName     string `json:",omitempty"`
	Layer3Multicast bool   `json:",omitempty"`
}
