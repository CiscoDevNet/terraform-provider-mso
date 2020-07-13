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

func TestAccMSOSchemaSiteAnpEpgSelector_Basic(t *testing.T) {
	var ss SiteAnpEpgSelectorTest
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpEpgSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSiteAnpEpgSelectorConfig_basic("one"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteAnpEpgSelectorExists("mso_schema_site_anp_epg_selector.selector1", &ss),
					testAccCheckMSOSchemaSiteAnpEpgSelectorAttributes("one", &ss),
				),
			},
		},
	})
}

func TestAccMSOSchemaSiteAnpEpgSelector_Update(t *testing.T) {
	var ss SiteAnpEpgSelectorTest

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpEpgSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSiteAnpEpgSelectorConfig_basic("one"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteAnpEpgSelectorExists("mso_schema_site_anp_epg_selector.selector1", &ss),
					testAccCheckMSOSchemaSiteAnpEpgSelectorAttributes("one", &ss),
				),
			},
			{
				Config: testAccCheckMSOSiteAnpEpgSelectorConfig_basic("two"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteAnpEpgSelectorExists("mso_schema_site_anp_epg_selector.selector1", &ss),
					testAccCheckMSOSchemaSiteAnpEpgSelectorAttributes("two", &ss),
				),
			},
		},
	})
}

func testAccCheckMSOSiteAnpEpgSelectorConfig_basic(key string) string {
	return fmt.Sprintf(`
	resource "mso_schema_site_anp_epg_selector" "selector1" {
		schema_id   = "5c4d5bb72700000401f80948"
		site_id     = "5c7c95b25100008f01c1ee3c"
		template_name    = "Template1"
		anp_name    = "ANP"
		epg_name    = "DB"
		name        = "test_check"
		expressions {
			key = "%s"
			operator = "keyExist"
			value = "32"
		}
	 }
`, key)
}

func testAccCheckMSOSchemaSiteAnpEpgSelectorExists(selectorName string, ss *SiteAnpEpgSelectorTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs1, err1 := s.RootModule().Resources[selectorName]

		if !err1 {
			return fmt.Errorf("Selector %s not found", selectorName)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Selector id was set")
		}

		cont, err := client.GetViaURL(fmt.Sprintf("api/v1/schemas/5c4d5bb72700000401f80948"))
		if err != nil {
			return err
		}

		siteCount, err := cont.ArrayCount("sites")
		if err != nil {
			return fmt.Errorf("No Sites found")
		}

		tp := SiteAnpEpgSelectorTest{}
		found := false
		for i := 0; i < siteCount; i++ {
			siteCont, err := cont.ArrayElement(i, "sites")
			if err != nil {
				return err
			}

			currentSite := models.StripQuotes(siteCont.S("siteId").String())
			currentTemp := models.StripQuotes(siteCont.S("templateName").String())

			if currentTemp == "Template1" && currentSite == "5c7c95b25100008f01c1ee3c" {
				anpCount, err := siteCont.ArrayCount("anps")
				if err != nil {
					return fmt.Errorf("No Anp found")
				}

				for j := 0; j < anpCount; j++ {
					anpCont, err := siteCont.ArrayElement(j, "anps")
					if err != nil {
						return err
					}

					anpRef := models.StripQuotes(anpCont.S("anpRef").String())
					tokens := strings.Split(anpRef, "/")
					currentAnpName := tokens[len(tokens)-1]
					if currentAnpName == "ANP" {
						epgCount, err := anpCont.ArrayCount("epgs")
						if err != nil {
							return fmt.Errorf("No Epg found")
						}

						for k := 0; k < epgCount; k++ {
							epgCont, err := anpCont.ArrayElement(k, "epgs")
							if err != nil {
								return err
							}

							epgRef := models.StripQuotes(epgCont.S("epgRef").String())
							tokensEpg := strings.Split(epgRef, "/")
							currentEpgName := tokensEpg[len(tokensEpg)-1]
							if currentEpgName == "DB" {
								selectorCount, err := epgCont.ArrayCount("selectors")
								if err != nil {
									return fmt.Errorf("No selectors found")
								}

								for s := 0; s < selectorCount; s++ {
									selectorCont, err := epgCont.ArrayElement(s, "selectors")
									if err != nil {
										return err
									}

									currentName := models.StripQuotes(selectorCont.S("name").String())
									if currentName == "test_check" {
										found = true
										tp.Name = currentName
										exps := selectorCont.S("expressions").Data().([]interface{})

										expressionsList := make([]interface{}, 0, 1)
										for _, val := range exps {
											temp := val.(map[string]interface{})
											expressionsMap := make(map[string]interface{})

											expressionsMap["key"] = temp["key"]

											expressionsMap["operator"] = temp["operator"]

											if temp["value"] != nil {
												expressionsMap["value"] = temp["value"]
											}
											expressionsList = append(expressionsList, expressionsMap)
										}
										tp.Expressions = expressionsList
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

func testAccCheckMSOSchemaSiteAnpEpgSelectorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		cont, err := client.GetViaURL(fmt.Sprintf("api/v1/schemas/5c4d5bb72700000401f80948"))
		if rs.Type == "mso_schema_site_anp_epg_selector" {
			if err != nil {
				return err
			}
		} else {
			siteCount, err := cont.ArrayCount("sites")
			if err != nil {
				return fmt.Errorf("No Sites found")
			}

			for i := 0; i < siteCount; i++ {
				siteCont, err := cont.ArrayElement(i, "sites")
				if err != nil {
					return err
				}

				currentSite := models.StripQuotes(siteCont.S("siteId").String())
				currentTemp := models.StripQuotes(siteCont.S("templateName").String())

				if currentTemp == "Template1" && currentSite == "5c7c95b25100008f01c1ee3c" {
					anpCount, err := siteCont.ArrayCount("anps")
					if err != nil {
						return fmt.Errorf("No Anp found")
					}

					for j := 0; j < anpCount; j++ {
						anpCont, err := siteCont.ArrayElement(j, "anps")
						if err != nil {
							return err
						}

						anpRef := models.StripQuotes(anpCont.S("anpRef").String())
						tokens := strings.Split(anpRef, "/")
						currentAnpName := tokens[len(tokens)-1]
						if currentAnpName == "ANP" {
							epgCount, err := anpCont.ArrayCount("epgs")
							if err != nil {
								return fmt.Errorf("No Epg found")
							}

							for k := 0; k < epgCount; k++ {
								epgCont, err := anpCont.ArrayElement(k, "epgs")
								if err != nil {
									return err
								}

								epgRef := models.StripQuotes(epgCont.S("epgRef").String())
								tokensEpg := strings.Split(epgRef, "/")
								currentEpgName := tokensEpg[len(tokensEpg)-1]
								if currentEpgName == "DB" {
									selectorCount, err := epgCont.ArrayCount("selectors")
									if err != nil {
										return fmt.Errorf("No selectors found")
									}

									for s := 0; s < selectorCount; s++ {
										selectorCont, err := epgCont.ArrayElement(s, "selectors")
										if err != nil {
											return err
										}

										currentName := models.StripQuotes(selectorCont.S("name").String())
										if currentName == "test_check" {
											return fmt.Errorf("Schema Site Anp Epg Selector still exist")
										}
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

func testAccCheckMSOSchemaSiteAnpEpgSelectorAttributes(key string, ss *SiteAnpEpgSelectorTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		val := ss.Expressions[0].(map[string]interface{})
		if key != val["key"] {
			return fmt.Errorf("Bad Site anp epg Selector key %v", val["key"])
		}
		return nil
	}
}

type SiteAnpEpgSelectorTest struct {
	Id          string        `json:",omitempty"`
	SchemaId    string        `json:",omitempty"`
	SiteId      string        `json:",omitempty"`
	Template    string        `json:",omitempty"`
	AnpName     string        `json:",omitempty"`
	EpgName     string        `json:",omitempty"`
	Name        string        `json:",omitempty"`
	Expressions []interface{} `json:",omitempty"`
}
