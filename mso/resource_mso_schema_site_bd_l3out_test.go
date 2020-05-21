package mso

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccMSOSchemaSiteBdL3out_Basic(t *testing.T) {
	var ss SchemaSiteBdL3out
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteBdL3outDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSchemaSiteBdL3outConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteBdL3outExists("mso_schema_site_bd_l3out.bdL3out", &ss),
					testAccCheckMSOSchemaSiteBdL3outAttributes(&ss),
				),
			},
		},
	})
}

func testAccCheckMSOSchemaSiteBdL3outConfig_basic() string {
	return fmt.Sprintf(`
	resource "mso_schema_site_bd_l3out" "bdL3out" {
		schema_id = "5d5dbf3f2e0000580553ccce"
		template_name = "Template1"
		site_id = "5c7c95b25100008f01c1ee3c"
		bd_name = "WebServer-Finance"
		l3out_name = "l3out1234"
	  }
	`)
}

func testAccCheckMSOSchemaSiteBdL3outExists(siteBdL3outName string, ss *SchemaSiteBdL3out) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs, err := s.RootModule().Resources[siteBdL3outName]

		if !err {
			return fmt.Errorf("Bd L3out %s not found", siteBdL3outName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3out id was set")
		}

		cont, errs := client.GetViaURL("api/v1/schemas/5d5dbf3f2e0000580553ccce")
		if errs != nil {
			return errs
		}
		count, ers := cont.ArrayCount("sites")
		if ers != nil {
			return fmt.Errorf("No Sites found")
		}

		tp := SchemaSiteBdL3out{}
		found := false

		for i := 0; i < count; i++ {
			tempCont, err := cont.ArrayElement(i, "sites")
			if err != nil {
				return err
			}
			apiSite := models.StripQuotes(tempCont.S("siteId").String())

			if apiSite == "5c7c95b25100008f01c1ee3c" {
				bdCount, err := tempCont.ArrayCount("bds")
				if err != nil {
					return fmt.Errorf("Unable to get Bd list")
				}
				for j := 0; j < bdCount; j++ {
					bdCont, err := tempCont.ArrayElement(j, "bds")
					if err != nil {
						return err
					}
					apiBdRef := models.StripQuotes(bdCont.S("bdRef").String())
					split := strings.Split(apiBdRef, "/")
					apiBd := split[6]
					if apiBd == "WebServer-Finance" {
						l3outCount, err := bdCont.ArrayCount("l3Outs")
						if err != nil {
							return fmt.Errorf("Unable to get l3Outs list")
						}
						for k := 0; k < l3outCount; k++ {
							l3outCont, err := bdCont.ArrayElement(k, "l3Outs")
							if err != nil {
								return err
							}
							tempVar := l3outCont.String()
							apiL3out := strings.Trim(tempVar, "\"")
							if apiL3out == "l3out1234" {
								tp.siteId = apiSite
								tp.bdName = apiBd
								tp.l3outName = apiL3out

								found = true
								break
							}
						}
					}
				}
			}
		}

		if !found {
			return fmt.Errorf("Bd L3out not found from API")
		}
		tp1 := &tp
		*ss = *tp1
		return nil
	}
}

func testAccCheckMSOSchemaSiteBdL3outDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_site_bd_l3out" {
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
						return fmt.Errorf("No Site exists")
					}
					apiSite := models.StripQuotes(tempCont.S("siteId").String())

					if apiSite == "5c7c95b25100008f01c1ee3c" {
						bdCount, err := tempCont.ArrayCount("bds")
						if err != nil {
							return fmt.Errorf("Unable to get Bd list")
						}
						for j := 0; j < bdCount; j++ {
							bdCont, err := tempCont.ArrayElement(j, "bds")
							if err != nil {
								return err
							}
							apiBdRef := models.StripQuotes(bdCont.S("bdRef").String())
							split := strings.Split(apiBdRef, "/")
							apiBd := split[6]
							if apiBd == "WebServer-Finance" {
								l3outCount, err := bdCont.ArrayCount("l3Outs")
								if err != nil {
									return fmt.Errorf("Unable to get l3Outs list")
								}
								for k := 0; k < l3outCount; k++ {
									l3outCont, err := bdCont.ArrayElement(k, "l3Outs")
									if err != nil {
										return err
									}
									tempVar := l3outCont.String()
									apiL3out := strings.Trim(tempVar, "\"")
									if apiL3out == "l3out1234" {
										return fmt.Errorf("The Bd L3out still exists")
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func testAccCheckMSOSchemaSiteBdL3outAttributes(ss *SchemaSiteBdL3out) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "WebServer-Finance" != ss.bdName {
			return fmt.Errorf("Bad Bd name %s", ss.bdName)
		}
		return nil
	}
}

type SchemaSiteBdL3out struct {
	siteId    string
	bdName    string
	l3outName string
}
