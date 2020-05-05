package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccMSOSchemaTemplateL3out_Basic(t *testing.T) {
	var ss TemplateL3out
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateL3outDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateL3outConfig_basic("WoS_Cloud_VRF2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateL3outExists("mso_schema_template_l3out.template_l3out", &ss),
					testAccCheckMSOSchemaTemplateL3outAttributes("WoS_Cloud_VRF2", &ss),
				),
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
				Config: testAccCheckMSOTemplateL3outConfig_basic("WoS_Cloud_VRF2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateL3outExists("mso_schema_template_l3out.template_l3out", &ss),
					testAccCheckMSOSchemaTemplateL3outAttributes("WoS_Cloud_VRF2", &ss),
				),
			},
			{
				Config: testAccCheckMSOTemplateL3outConfig_basic("vrf589"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateL3outExists("mso_schema_template_l3out.template_l3out", &ss),
					testAccCheckMSOSchemaTemplateL3outAttributes("vrf589", &ss),
				),
			},
		},
	})
}

func testAccCheckMSOTemplateL3outConfig_basic(name string) string {
	return fmt.Sprintf(`
	resource "mso_schema_template_l3out" "template_l3out" {
		schema_id = "5c6c16d7270000c710f8094d"
		template_name = "Template1"
		l3out_name = "l3out1"
		display_name = "l3out1"
		vrf_name = "%v"
	}
`, name)
}

func testAccCheckMSOSchemaTemplateL3outExists(l3outName string, ss *TemplateL3out) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs1, err1 := s.RootModule().Resources[l3outName]

		if !err1 {
			return fmt.Errorf("L3out %s not found", l3outName)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Schema id was set")
		}

		cont, err := client.GetViaURL("api/v1/schemas/5c6c16d7270000c710f8094d")
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
					if apiL3out == "l3out1" {
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

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_template_l3out" {
			cont, err := client.GetViaURL("api/v1/schemas/5c6c16d7270000c710f8094d")
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
							if apiL3out == "l3out1" {
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
		if "l3out1" != ss.display_name {
			return fmt.Errorf("Bad Template L3out display name %s", ss.display_name)
		}

		if vrf_name != ss.vrf_name {
			return fmt.Errorf("Bad Template L3out VRF name %s", ss.vrf_name)
		}
		return nil
	}
}

type TemplateL3out struct {
	display_name string
	vrf_name     string
}
