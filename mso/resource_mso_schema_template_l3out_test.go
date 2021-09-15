package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccMSOSchemaTemplateL3out_Basic(t *testing.T) {
	var ss TemplateL3out
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateL3outDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateL3outConfig_basic("VRF"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateL3outExists("mso_schema.schema1", "mso_schema_template_l3out.template_l3out", &ss),
					testAccCheckMSOSchemaTemplateL3outAttributes("VRF", &ss),
				),
			},
			{
				ResourceName:      "mso_schema_template_l3out.template_l3out",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMSOSchemaTemplateL3out_Update(t *testing.T) {
	var ss TemplateL3out

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateL3outDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateL3outConfig_basic("VRF"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateL3outExists("mso_schema.schema1", "mso_schema_template_l3out.template_l3out", &ss),
					testAccCheckMSOSchemaTemplateL3outAttributes("VRF", &ss),
				),
			},
			{
				Config: testAccCheckMSOTemplateL3outConfig_basic("VRF2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateL3outExists("mso_schema.schema1", "mso_schema_template_l3out.template_l3out", &ss),
					testAccCheckMSOSchemaTemplateL3outAttributes("VRF2", &ss),
				),
			},
		},
	})
}

func testAccCheckMSOTemplateL3outConfig_basic(name string) string {
	return fmt.Sprintf(`
	resource "mso_schema" "schema1" {
		name          = "Schema2"
		template_name = "Template1"
		tenant_id     = "5fb5fed8520000452a9e8911"
	  
	  }

	  resource "mso_schema_template_vrf" "vrf1" {
		schema_id=mso_schema.schema1.id
		template=mso_schema.schema1.template_name
		name= "VRF"
		display_name="vrf1"
		layer3_multicast=true
		vzany=false
	  }

    resource "mso_schema_template_vrf" "vrf2" {
		schema_id=mso_schema.schema1.id
		template=mso_schema_template_vrf.vrf1.template
		name= "VRF2"
		display_name="vrf2"
		layer3_multicast=true
		vzany=false
	  }


	resource "mso_schema_template_l3out" "template_l3out" {
		schema_id = mso_schema.schema1.id
		template_name = mso_schema_template_vrf.vrf2.template
		l3out_name = "l3out3"
		display_name = "check2"
		vrf_name = "%s"
	}

`, name)
}

func testAccCheckMSOSchemaTemplateL3outExists(schemaName string, l3outName string, ss *TemplateL3out) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs1, err1 := s.RootModule().Resources[schemaName]
		rs2, err2 := s.RootModule().Resources[l3outName]
		if !err1 {
			return fmt.Errorf("Schema %s not found", schemaName)
		}
		if !err2 {
			return fmt.Errorf("L3out %s not found", l3outName)
		}
		if rs2.Primary.ID == "" {
			return fmt.Errorf("No Schema id was set")
		}

		cont, err := client.GetViaURL("api/v1/schemas/" + rs1.Primary.ID)
		if err != nil {
			return err
		}
		count, err := cont.ArrayCount("templates")
		if err != nil {
			return fmt.Errorf("No Template found")
		}
		tp := TemplateL3out{}
		found := false
		for i := 0; i < count; i++ {
			tempCont, err := cont.ArrayElement(i, "templates")
			if err != nil {
				return err
			}

			apiTemplateName := models.StripQuotes(tempCont.S("name").String())
			if apiTemplateName == "Template1" {
				l3outCount, err := tempCont.ArrayCount("intersiteL3outs")
				if err != nil {
					return fmt.Errorf("Unable to get L3out list")
				}
				for j := 0; j < l3outCount; j++ {
					l3outCont, err := tempCont.ArrayElement(j, "intersiteL3outs")
					if err != nil {
						return err
					}
					apiL3out := models.StripQuotes(l3outCont.S("name").String())
					if apiL3out == "l3out3" {
						tp.display_name = models.StripQuotes(l3outCont.S("displayName").String())
						vrfRef := models.StripQuotes(l3outCont.S("vrfRef").String())
						re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/vrfs/(.*)")
						match := re.FindStringSubmatch(vrfRef)
						tp.vrf_name = match[3]
						found = true
						break
					}
				}
			}
		}
		if !found {
			return fmt.Errorf("L3out not found from API")
		}

		tp1 := &tp
		*ss = *tp1
		return nil
	}
}

func testAccCheckMSOSchemaTemplateL3outDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	rs1, err1 := s.RootModule().Resources["mso_schema.schema1"]
	if !err1 {
		return fmt.Errorf("Schema %s not found", "mso_schema.schema1")
	}
	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_template_l3out" {
			cont, err := client.GetViaURL("api/v1/schemas/" + rs1.Primary.ID)
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
						return fmt.Errorf("No Template exists")
					}
					apiTemplateName := models.StripQuotes(tempCont.S("name").String())
					if apiTemplateName == "Template1" {
						l3outCount, err := tempCont.ArrayCount("intersiteL3outs")
						if err != nil {
							return fmt.Errorf("Unable to get L3out list")
						}
						for j := 0; j < l3outCount; j++ {
							l3outCont, err := tempCont.ArrayElement(j, "intersiteL3outs")
							if err != nil {
								return err
							}
							apiL3out := models.StripQuotes(l3outCont.S("name").String())
							if apiL3out == "l3out3" {
								return fmt.Errorf("template L3Out still exists.")
							}
						}
					}
				}
			}
		}
	}
	return nil
}
func testAccCheckMSOSchemaTemplateL3outAttributes(vrf_name string, ss *TemplateL3out) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "VRF" != ss.vrf_name {
			return fmt.Errorf("Bad Template VRF name %s", ss.vrf_name)
		}

		return nil
	}
}

type TemplateL3out struct {
	display_name string
	vrf_name     string
}
