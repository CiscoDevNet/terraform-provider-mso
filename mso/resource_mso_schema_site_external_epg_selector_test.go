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

func TestAccMSOSchemaSiteExternalEpgSelector_Basic(t *testing.T) {
	var ss SiteEPGSelectorTest
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteExternalEpgSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSiteExternalEpgSelectorConfig_basic("1.2.3.4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteExternalEpgSelectorExists("mso_schema_site_external_epg_selector.selector1", &ss),
					testAccCheckMSOSchemaSiteExternalEpgSelectorAttributes("1.2.3.4", &ss),
				),
			},
		},
	})
}

func TestAccMSOSchemaSiteExternalEpgSelector_Update(t *testing.T) {
	var ss SiteEPGSelectorTest

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteExternalEpgSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSiteExternalEpgSelectorConfig_basic("1.2.3.4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteExternalEpgSelectorExists("mso_schema_site_external_epg_selector.selector1", &ss),
					testAccCheckMSOSchemaSiteExternalEpgSelectorAttributes("1.2.3.4", &ss),
				),
			},
			{
				Config: testAccCheckMSOSiteExternalEpgSelectorConfig_basic("5.4.6.7"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteExternalEpgSelectorExists("mso_schema_site_external_epg_selector.selector1", &ss),
					testAccCheckMSOSchemaSiteExternalEpgSelectorAttributes("5.4.6.7", &ss),
				),
			},
		},
	})
}

func testAccCheckMSOSiteExternalEpgSelectorConfig_basic(ip string) string {
	return fmt.Sprintf(`
	resource "mso_schema_site_external_epg_selector" "selector1" {
		schema_id = "5f043b3b2c0000f47e812a0b"
		template_name = "Template1"
		site_id = "5c7c95d9510000cf01c1ee3d"
		external_epg_name = "test_epg"
		name = "test_selector"
    	ip = "%s"
	}
`, ip)
}

func testAccCheckMSOSchemaSiteExternalEpgSelectorExists(selectorName string, ss *SiteEPGSelectorTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs1, err1 := s.RootModule().Resources[selectorName]

		if !err1 {
			return fmt.Errorf("Selector %s not found", selectorName)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Selector id was set")
		}

		cont, err := client.GetViaURL("api/v1/schemas/5f043b3b2c0000f47e812a0b")
		if err != nil {
			return err
		}
		count, err := cont.ArrayCount("sites")
		if err != nil {
			return fmt.Errorf("No Sites found")
		}
		tp := SiteEPGSelectorTest{}
		found := false

		for i := 0; i < count; i++ {
			siteCont, err := cont.ArrayElement(i, "sites")
			if err != nil {
				return fmt.Errorf("Error fetching site")
			}

			site := models.StripQuotes(siteCont.S("siteId").String())
			tempName := models.StripQuotes(siteCont.S("templateName").String())

			if tempName == "Template1" && site == "5c7c95d9510000cf01c1ee3d" {
				extrEpgCount, err := siteCont.ArrayCount("externalEpgs")
				if err != nil {
					return fmt.Errorf("no externalEpgs found")
				}

				for j := 0; j < extrEpgCount; j++ {
					extrEpgCont, err := siteCont.ArrayElement(j, "externalEpgs")
					if err != nil {
						return fmt.Errorf("Error fetching external Epg")
					}

					extEpgRef := models.StripQuotes(extrEpgCont.S("externalEpgRef").String())
					tokens := strings.Split(extEpgRef, "/")
					extEpgName := tokens[len(tokens)-1]
					if extEpgName == "test_epg" {
						selectorCount, err := extrEpgCont.ArrayCount("subnets")
						if err != nil {
							return fmt.Errorf("No selectors found")
						}

						for k := 0; k < selectorCount; k++ {
							selectorCont, err := extrEpgCont.ArrayElement(k, "subnets")
							if err != nil {
								return fmt.Errorf("Error fetching selector")
							}

							selectorName := models.StripQuotes(selectorCont.S("name").String())
							if selectorName == "test_selector" {
								found = true
								tp.Name = selectorName
								tp.Ip = models.StripQuotes(selectorCont.S("ip").String())
								break
							}
						}
					}
					if found {
						break
					}
				}
			}
			if found {
				break
			}
		}

		if !found {
			return fmt.Errorf("Selector not found from API")
		}
		tp1 := &tp
		*ss = *tp1
		return nil
	}
}

func testAccCheckMSOSchemaSiteExternalEpgSelectorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		cont, err := client.GetViaURL("api/v1/schemas/5f043b3b2c0000f47e812a0b")
		if rs.Type == "mso_schema_site_external_epg_selector" {

			if err != nil {
				return err
			}
		} else {
			count, err := cont.ArrayCount("sites")
			if err != nil {
				return fmt.Errorf("No sites found")
			}
			for i := 0; i < count; i++ {
				siteCont, err := cont.ArrayElement(i, "sites")
				if err != nil {
					return fmt.Errorf("Error fetching site")
				}

				site := models.StripQuotes(siteCont.S("siteId").String())
				tempName := models.StripQuotes(siteCont.S("templateName").String())
				if tempName == "Template1" && site == "5c7c95d9510000cf01c1ee3d" {
					extrEpgCount, err := siteCont.ArrayCount("externalEpgs")
					if err != nil {
						return fmt.Errorf("no externalEpgs found")
					}

					for j := 0; j < extrEpgCount; j++ {
						extrEpgCont, err := siteCont.ArrayElement(j, "externalEpgs")
						if err != nil {
							return fmt.Errorf("Error fetching external Epg")
						}

						extEpgRef := models.StripQuotes(extrEpgCont.S("externalEpgRef").String())
						tokens := strings.Split(extEpgRef, "/")
						extEpgName := tokens[len(tokens)-1]
						if extEpgName == "test_epg" {
							selectorCount, err := extrEpgCont.ArrayCount("subnets")
							if err != nil {
								return fmt.Errorf("No selectors found")
							}

							for k := 0; k < selectorCount; k++ {
								selectorCont, err := extrEpgCont.ArrayElement(k, "subnets")
								if err != nil {
									return fmt.Errorf("Error fetching selector")
								}

								selectorName := models.StripQuotes(selectorCont.S("name").String())
								if selectorName == "test_selector" {
									return fmt.Errorf("Schema Site external epg selector still exist")
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

func testAccCheckMSOSchemaSiteExternalEpgSelectorAttributes(ip string, ss *SiteEPGSelectorTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if ip != ss.Ip {
			return fmt.Errorf("Bad site external epg Selector ip %v", ss.Ip)
		}
		return nil
	}
}

type SiteEPGSelectorTest struct {
	SchemaId string `json:",omitempty"`
	Template string `json:",omitempty"`
	EpgName  string `json:",omitempty"`
	Name     string `json:",omitempty"`
	Ip       string `json:",omitempty"`
}
