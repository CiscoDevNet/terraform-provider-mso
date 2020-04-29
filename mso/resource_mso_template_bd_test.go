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

func TestAccMSOSchemaTemplateBD_Basic(t *testing.T) {
	var ss TemplateBD
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateBDDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateBDConfig_basic("flood"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateBDExists("mso_schema_template_bd.bridge_domain", &ss),
					testAccCheckMSOSchemaTemplateBDAttributes("flood", &ss),
				),
			},
		},
	})
}

func TestAccMSOSchemaTemplateBD_Update(t *testing.T) {
	var ss TemplateBD

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateBDDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateBDConfig_basic("flood"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateBDExists("mso_schema_template_bd.bridge_domain", &ss),
					testAccCheckMSOSchemaTemplateBDAttributes("flood", &ss),
				),
			},
			{
				Config: testAccCheckMSOTemplateBDConfig_basic("proxy"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateBDExists("mso_schema_template_bd.bridge_domain", &ss),
					testAccCheckMSOSchemaTemplateBDAttributes("proxy", &ss),
				),
			},
		},
	})
}

func testAccCheckMSOTemplateBDConfig_basic(unicast string) string {
	return fmt.Sprintf(`
	resource "mso_schema_template_bd" "bridge_domain" {
		schema_id = "5ea809672c00003bc40a2799"
		template_name = "Template1"
		name = "testAccBD"
		display_name = "testAcc"
		vrf_name = "demo"
		layer2_unknown_unicast = "%s" 
	}
`, unicast)
}

func testAccCheckMSOSchemaTemplateBDExists(bdName string, ss *TemplateBD) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs1, err1 := s.RootModule().Resources[bdName]

		if !err1 {
			return fmt.Errorf("BD %s not found", bdName)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Schema id was set")
		}

		cont, err := client.GetViaURL("api/v1/schemas/5ea809672c00003bc40a2799")
		if err != nil {
			return err
		}
		count, err := cont.ArrayCount("templates")
		if err != nil {
			return fmt.Errorf("No Template found")
		}
		tp := TemplateBD{}
		found := false
		for i := 0; i < count; i++ {
			tempCont, err := cont.ArrayElement(i, "templates")
			if err != nil {
				return err
			}

			apiTemplateName := models.StripQuotes(tempCont.S("name").String())
			if apiTemplateName == "Template1" {
				bdCount, err := tempCont.ArrayCount("bds")
				if err != nil {
					return fmt.Errorf("Unable to get BD list")
				}
				for j := 0; j < bdCount; j++ {
					bdCont, err := tempCont.ArrayElement(j, "bds")
					if err != nil {
						return err
					}
					apiBD := models.StripQuotes(bdCont.S("name").String())
					if apiBD == "testAccBD" {
						tp.display_name = models.StripQuotes(bdCont.S("displayName").String())
						tp.layer2_unknown_unicast = models.StripQuotes(bdCont.S("l2UnknownUnicast").String())
						vrfRef := models.StripQuotes(bdCont.S("vrfRef").String())
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
			return fmt.Errorf("BD not found from API")
		}

		tp1 := &tp

		*ss = *tp1
		return nil
	}
}

func testAccCheckMSOSchemaTemplateBDDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_template_bd" {
			cont, err := client.GetViaURL("api/v1/schemas/5ea809672c00003bc40a2799")
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
						bdCount, err := tempCont.ArrayCount("bds")
						if err != nil {
							return fmt.Errorf("Unable to get BD list")
						}
						for j := 0; j < bdCount; j++ {
							bdCont, err := tempCont.ArrayElement(j, "bds")
							if err != nil {
								return err
							}
							apiBD := models.StripQuotes(bdCont.S("name").String())
							if apiBD == "testAccBD" {
								return fmt.Errorf("template bridge domain still exists.")
							}
						}
					}

				}
			}
		}
	}
	return nil
}
func testAccCheckMSOSchemaTemplateBDAttributes(layer2_unknown_unicast string, ss *TemplateBD) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if layer2_unknown_unicast != ss.layer2_unknown_unicast {
			return fmt.Errorf("Bad Template BD layer2_unknown_unicast %s", ss.layer2_unknown_unicast)
		}

		if "testAcc" != ss.display_name {
			return fmt.Errorf("Bad Template BD display name %s", ss.display_name)
		}

		if "demo" != ss.vrf_name {
			return fmt.Errorf("Bad Template BD VRF name %s", ss.vrf_name)
		}
		return nil
	}
}

type TemplateBD struct {
	display_name           string
	vrf_name               string
	layer2_unknown_unicast string
}
