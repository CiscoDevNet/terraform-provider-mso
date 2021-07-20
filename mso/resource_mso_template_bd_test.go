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
					testAccCheckMSOSchemaTemplateBDExists("mso_schema.schema1", "mso_schema_template_bd.bridge_domain", &ss),
					testAccCheckMSOSchemaTemplateBDAttributes("flood", &ss),
				),
			},
			{
				ResourceName:      "mso_schema_template_bd.bridge_domain",
				ImportState:       true,
				ImportStateVerify: true,
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
					testAccCheckMSOSchemaTemplateBDExists("mso_schema.schema1", "mso_schema_template_bd.bridge_domain", &ss),
					testAccCheckMSOSchemaTemplateBDAttributes("flood", &ss),
				),
			},
			{
				Config: testAccCheckMSOTemplateBDConfig_basic("proxy"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateBDExists("mso_schema.schema1", "mso_schema_template_bd.bridge_domain", &ss),
					testAccCheckMSOSchemaTemplateBDAttributes("proxy", &ss),
				),
			},
		},
	})
}

func testAccCheckMSOTemplateBDConfig_basic(unicast string) string {
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

	resource "mso_schema_template_bd" "bridge_domain" {
		schema_id = mso_schema.schema1.id
		template_name = "Template1"
		name = "BD"
		display_name = "bd1"
		vrf_name = mso_schema_template_vrf.vrf1.name
		intersite_bum_traffic=true
		optimize_wan_bandwidth=true
		layer2_stretch=true
		layer3_multicast=true
		layer2_unknown_unicast = "%s" 
		dhcp_policy ={
			name="dh"
			version=1
			dhcp_option_policy_name="dho"
			dhcp_option_policy_version=1
		}
	}
`, unicast)
}

func testAccCheckMSOSchemaTemplateBDExists(schemaName string, bdName string, ss *TemplateBD) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs1, err1 := s.RootModule().Resources[schemaName]
		rs2, err2 := s.RootModule().Resources[bdName]

		if !err1 {
			return fmt.Errorf("Schema %s not found", schemaName)
		}
		if !err2 {
			return fmt.Errorf("BD %s not found", bdName)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Schema id was set")
		}
		if rs2.Primary.ID == "" {
			return fmt.Errorf("No BD was set")
		}

		cont, err := client.GetViaURL("api/v1/schemas/" + rs1.Primary.ID)
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
					if apiBD == "BD" {
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
	rs1, err1 := s.RootModule().Resources["mso_schema.schema1"]
	if !err1 {
		return fmt.Errorf("Schema %s not found", "mso_schema.schema1")
	}
	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_template_bd" {
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
							if apiBD == "BD" {
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

		if "bd1" != ss.display_name {
			return fmt.Errorf("Bad Template BD display name %s", ss.display_name)
		}

		if "VRF" != ss.vrf_name {
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
