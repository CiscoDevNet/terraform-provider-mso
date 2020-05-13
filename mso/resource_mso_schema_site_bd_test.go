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

func TestAccMSOSchemaSiteBd_Basic(t *testing.T) {
	var ss SiteBd
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteBdDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSchemaSiteBdConfig_basic(false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteBdExists("mso_schema_site_bd.bd1", &ss),
					testAccCheckMSOSchemaSiteBdAttributes(false, &ss),
				),
			},
		},
	})
}

func TestAccMSOSchemaSiteBd_Update(t *testing.T) {
	var ss SiteBd

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteBdDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSchemaSiteBdConfig_basic(false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteBdExists("mso_schema_site_bd.bd1", &ss),
					testAccCheckMSOSchemaSiteBdAttributes(false, &ss),
				),
			},
			{
				Config: testAccCheckMSOSchemaSiteBdConfig_basic(true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteBdExists("mso_schema_site_bd.bd1", &ss),
					testAccCheckMSOSchemaSiteBdAttributes(true, &ss),
				),
			},
		},
	})
}
func testAccCheckMSOSchemaSiteBdConfig_basic(host bool) string {
	return fmt.Sprintf(`
	resource "mso_schema_site_bd" "bd1" {
		schema_id = "5d5dbf3f2e0000580553ccce"
		bd_name = "bd4"
		template_name = "Template1"
		site_id = "5c7c95b25100008f01c1ee3c"
	    host = %v
	  }`, host)
}

func testAccCheckMSOSchemaSiteBdExists(anpName string, ss *SiteBd) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs1, err1 := s.RootModule().Resources[anpName]

		if !err1 {
			return fmt.Errorf("Site Bd %s not found", anpName)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Bd Id was set")
		}

		cont, err := client.GetViaURL("api/v1/schemas/5d5dbf3f2e0000580553ccce")
		if err != nil {
			return err
		}
		count, err := cont.ArrayCount("sites")
		if err != nil {
			return fmt.Errorf("No Sites found")
		}
		tp := SiteBd{}

		found := false

		for i := 0; i < count; i++ {
			tempCont, err := cont.ArrayElement(i, "sites")
			if err != nil {
				return err
			}
			apiSite := models.StripQuotes(tempCont.S("siteId").String())

			if apiSite == "5c7c95b25100008f01c1ee3c" {
				tp.siteId = apiSite
				bdCount, err := tempCont.ArrayCount("bds")
				if err != nil {
					return fmt.Errorf("Unable to get bd list")
				}
				for j := 0; j < bdCount; j++ {
					bdCont, err := tempCont.ArrayElement(j, "bds")
					if err != nil {
						return err
					}
					bdRef := models.StripQuotes(bdCont.S("bdRef").String())
					re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/bds/(.*)")
					match := re.FindStringSubmatch(bdRef)
					if match[3] == "bd4" {
						tp.name = match[3]
						tp.host = bdCont.S("hostBasedRouting").Data().(bool)
						found = true
						break
					}
				}
			}
		}

		if !found {
			return fmt.Errorf("bd Epg not found from API")
		}
		tp1 := &tp
		*ss = *tp1
		return nil
	}
}

func testAccCheckMSOSchemaSiteBdDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_site_bd" {
			cont, err := client.GetViaURL("api/v1/schemas/5d5dbf3f2e0000580553ccce")
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

					if apiSite == "5c7c95b25100008f01c1ee3c" {

						bdCount, err := tempCont.ArrayCount("bds")
						if err != nil {
							return fmt.Errorf("Unable to get bd list")
						}
						for j := 0; j < bdCount; j++ {
							bdCont, err := tempCont.ArrayElement(j, "bds")
							if err != nil {
								return err
							}
							bdRef := models.StripQuotes(bdCont.S("bdRef").String())
							re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/bds/(.*)")
							match := re.FindStringSubmatch(bdRef)
							if match[3] == "bd4" {
								return fmt.Errorf("bd Still exists")
							}
						}
					}
				}

			}
		}
	}

	return nil

}

func testAccCheckMSOSchemaSiteBdAttributes(preferred_group bool, ss *SiteBd) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "5c7c95b25100008f01c1ee3c" != ss.siteId {
			return fmt.Errorf("Bad siteId %s", ss.siteId)
		}
		if preferred_group != ss.host {
			return fmt.Errorf("Bad Host %v", ss.host)
		}
		return nil
	}
}

type SiteBd struct {
	name   string
	siteId string
	host   bool
}
