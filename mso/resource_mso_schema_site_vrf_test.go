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

func TestAccMSOSchemaSiteVrf_Basic(t *testing.T) {
	var ss SiteVrf
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteVrfDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSchemaSiteVrfConfig_basic("Template1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteVrfExists("mso_schema_site_vrf.vrf1", &ss),
					testAccCheckMSOSchemaSiteVrfAttributes("Template1", &ss),
				),
			},
		},
	})
}

func testAccCheckMSOSchemaSiteVrfConfig_basic(preferred_group string) string {
	return fmt.Sprintf(`
	resource "mso_schema_site_vrf" "vrf1" {
        template_name = "Template1"
		site_id = "5c7c95d9510000cf01c1ee3d"
		schema_id ="5c6c16d7270000c710f8094d"
		vrf_name = "vrf3"
		
	  }`)
}

func testAccCheckMSOSchemaSiteVrfExists(vrfName string, ss *SiteVrf) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs1, err1 := s.RootModule().Resources[vrfName]

		if !err1 {
			return fmt.Errorf("Site Vrf %s not found", vrfName)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Vrf Id was set")
		}

		cont, err := client.GetViaURL("api/v1/schemas/5c6c16d7270000c710f8094d")
		if err != nil {
			return err
		}
		count, err := cont.ArrayCount("sites")
		if err != nil {
			return fmt.Errorf("No Sites found")
		}
		tp := SiteVrf{}

		found := false

		for i := 0; i < count; i++ {
			tempCont, err := cont.ArrayElement(i, "sites")
			if err != nil {
				return err
			}
			apiSite := models.StripQuotes(tempCont.S("siteId").String())

			if apiSite == "5c7c95d9510000cf01c1ee3d" {
				tp.siteId = apiSite
				vrfCount, err := tempCont.ArrayCount("vrfs")
				if err != nil {
					return fmt.Errorf("Unable to get Vrf list")
				}
				for j := 0; j < vrfCount; j++ {
					vrfCont, err := tempCont.ArrayElement(j, "vrfs")
					if err != nil {
						return err
					}
					vrfRef := models.StripQuotes(vrfCont.S("vrfRef").String())
					re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/vrfs/(.*)")
					match := re.FindStringSubmatch(vrfRef)
					if match[3] == "vrf3" {
						tp.name = match[3]
						found = true
						break
					}
				}
			}
		}

		if !found {
			return fmt.Errorf("Vrf not found from API")
		}
		tp1 := &tp
		*ss = *tp1
		return nil
	}
}

func testAccCheckMSOSchemaSiteVrfDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_site_vrf" {
			cont, err := client.GetViaURL("api/v1/schemas/5c6c16d7270000c710f8094d")
			if err != nil {
				return err
			} else {
				count, err := cont.ArrayCount("sites")
				if err != nil {
					return fmt.Errorf("No Sites found")
				}

				for i := 0; i < count; i++ {
					tempCont, err := cont.ArrayElement(i, "sites")
					if err != nil {
						return err
					}
					apiSite := models.StripQuotes(tempCont.S("siteId").String())

					if apiSite == "5c7c95d9510000cf01c1ee3d" {

						vrfCount, err := tempCont.ArrayCount("vrfs")
						if err != nil {
							return fmt.Errorf("Unable to get Vrf list")
						}
						for j := 0; j < vrfCount; j++ {
							vrfCont, err := tempCont.ArrayElement(j, "vrfs")
							if err != nil {
								return err
							}
							vrfRef := models.StripQuotes(vrfCont.S("vrfRef").String())
							re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/vrfs/(.*)")
							match := re.FindStringSubmatch(vrfRef)
							if match[3] == "vrf3" {
								return fmt.Errorf("Vrf Still exists")
							}
						}
					}
				}

			}
		}
	}

	return nil

}

func testAccCheckMSOSchemaSiteVrfAttributes(preferred_group string, ss *SiteVrf) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "5c7c95d9510000cf01c1ee3d" != ss.siteId {
			return fmt.Errorf("Bad siteId %s", ss.siteId)
		}
		return nil
	}
}

type SiteVrf struct {
	name   string
	siteId string
}
