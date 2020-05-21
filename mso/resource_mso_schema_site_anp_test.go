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

func TestAccMSOSchemaSiteAnp_Basic(t *testing.T) {
	var ss SiteAnp
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSchemaSiteAnpConfig_basic("Template1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteAnpExists("mso_schema_site_anp.anp1", &ss),
					testAccCheckMSOSchemaSiteAnpAttributes("Template1", &ss),
				),
			},
		},
	})
}

func testAccCheckMSOSchemaSiteAnpConfig_basic(preferred_group string) string {
	return fmt.Sprintf(`
	resource "mso_schema_site_anp" "anp1" {
		schema_id = "5c6c16d7270000c710f8094d"
		anp_name = "AP1234"
		template_name = "Template1"
		site_id = "5c7c95d9510000cf01c1ee3d"
		
	  }`)
}

func testAccCheckMSOSchemaSiteAnpExists(anpName string, ss *SiteAnp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs1, err1 := s.RootModule().Resources[anpName]

		if !err1 {
			return fmt.Errorf("Site Anp %s not found", anpName)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Anp Id was set")
		}

		cont, err := client.GetViaURL("api/v1/schemas/5c6c16d7270000c710f8094d")
		if err != nil {
			return err
		}
		count, err := cont.ArrayCount("sites")
		if err != nil {
			return fmt.Errorf("No Sites found")
		}
		tp := SiteAnp{}

		found := false

		for i := 0; i < count; i++ {
			tempCont, err := cont.ArrayElement(i, "sites")
			if err != nil {
				return err
			}
			apiSite := models.StripQuotes(tempCont.S("siteId").String())

			if apiSite == "5c7c95d9510000cf01c1ee3d" {
				tp.siteId = apiSite
				anpCount, err := tempCont.ArrayCount("anps")
				if err != nil {
					return fmt.Errorf("Unable to get Anp list")
				}
				for j := 0; j < anpCount; j++ {
					anpCont, err := tempCont.ArrayElement(j, "anps")
					if err != nil {
						return err
					}
					anpRef := models.StripQuotes(anpCont.S("anpRef").String())
					re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/anps/(.*)")
					match := re.FindStringSubmatch(anpRef)
					if match[3] == "AP1234" {
						tp.name = match[3]
						found = true
						break
					}
				}
			}
		}

		if !found {
			return fmt.Errorf("Anp Epg not found from API")
		}
		tp1 := &tp
		*ss = *tp1
		return nil
	}
}

func testAccCheckMSOSchemaSiteAnpDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_site_anp" {
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

						anpCount, err := tempCont.ArrayCount("anps")
						if err != nil {
							return fmt.Errorf("Unable to get Anp list")
						}
						for j := 0; j < anpCount; j++ {
							anpCont, err := tempCont.ArrayElement(j, "anps")
							if err != nil {
								return err
							}
							anpRef := models.StripQuotes(anpCont.S("anpRef").String())
							re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/anps/(.*)")
							match := re.FindStringSubmatch(anpRef)
							if match[3] == "AP1234" {
								return fmt.Errorf("ANP Still exists")
							}
						}
					}
				}

			}
		}
	}

	return nil

}

func testAccCheckMSOSchemaSiteAnpAttributes(preferred_group string, ss *SiteAnp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "5c7c95d9510000cf01c1ee3d" != ss.siteId {
			return fmt.Errorf("Bad siteId %s", ss.siteId)
		}
		return nil
	}
}

type SiteAnp struct {
	name   string
	siteId string
}
