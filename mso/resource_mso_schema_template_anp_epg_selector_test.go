package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccMSOSchemaTemplateAnpEpgSelector_Basic(t *testing.T) {
	var ss SelectorTest
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateAnpEpgSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateAnpEpgSelectorConfig_basic("one"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateAnpEpgSelectorExists("mso_schema_template_anp_epg_selector.selector1", &ss),
					testAccCheckMSOSchemaTemplateAnpEpgSelectorAttributes("one", &ss),
				),
			},
		},
	})
}

func TestAccMSOSchemaTemplateAnpEpgSelector_Update(t *testing.T) {
	var ss SelectorTest

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateAnpEpgSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateAnpEpgSelectorConfig_basic("one"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateAnpEpgSelectorExists("mso_schema_template_anp_epg_selector.selector1", &ss),
					testAccCheckMSOSchemaTemplateAnpEpgSelectorAttributes("one", &ss),
				),
			},
			{
				Config: testAccCheckMSOTemplateAnpEpgSelectorConfig_basic("two"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateAnpEpgSelectorExists("mso_schema_template_anp_epg_selector.selector1", &ss),
					testAccCheckMSOSchemaTemplateAnpEpgSelectorAttributes("two", &ss),
				),
			},
		},
	})
}

func testAccCheckMSOTemplateAnpEpgSelectorConfig_basic(key string) string {
	return fmt.Sprintf(`
	resource "mso_schema_template_anp_epg_selector" "selector1" {
		schema_id = "5ea809672c00003bc40a2799"
		template_name = "Template1"
		anp_name = "ap1"
		epg_name = "epg1"
		name = "test_check"
		expressions {
		  key = "%s"
		  operator = "equals"
		  value = "1"
		}
	}
`, key)
}

func testAccCheckMSOSchemaTemplateAnpEpgSelectorExists(selectorName string, ss *SelectorTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs1, err1 := s.RootModule().Resources[selectorName]

		if !err1 {
			return fmt.Errorf("Selector %s not found", selectorName)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Selector id was set")
		}

		cont, err := client.GetViaURL("api/v1/schemas/5ea809672c00003bc40a2799")
		if err != nil {
			return err
		}
		count, err := cont.ArrayCount("templates")
		if err != nil {
			return fmt.Errorf("No Template found")
		}
		tp := SelectorTest{}
		found := false

		for i := 0; i < count; i++ {
			tempCont, err := cont.ArrayElement(i, "templates")
			if err != nil {
				return fmt.Errorf("No Template found")
			}

			apiTemplate := models.StripQuotes(tempCont.S("name").String())

			if apiTemplate == "Template1" {
				tp.Template = apiTemplate
				anpCount, err := tempCont.ArrayCount("anps")
				if err != nil {
					return fmt.Errorf("Unable to get ANP list")
				}
				for j := 0; j < anpCount; j++ {
					anpCont, err := tempCont.ArrayElement(j, "anps")
					if err != nil {
						return err
					}
					apiANP := models.StripQuotes(anpCont.S("name").String())
					if apiANP == "ap1" {
						tp.AnpName = apiANP
						epgCount, err := anpCont.ArrayCount("epgs")
						if err != nil {
							return fmt.Errorf("Unable to get EPG list")
						}
						for k := 0; k < epgCount; k++ {
							epgCont, err := anpCont.ArrayElement(k, "epgs")
							if err != nil {
								return err
							}
							apiEPG := models.StripQuotes(epgCont.S("name").String())
							if apiEPG == "epg1" {
								tp.EpgName = apiEPG

								selectorCount, err := epgCont.ArrayCount("selectors")
								if err != nil {
									return fmt.Errorf("Unable to get selectorlist")
								}

								for s := 0; s < selectorCount; s++ {
									selectorCont, err := epgCont.ArrayElement(s, "selectors")
									if err != nil {
										return err
									}

									selName := models.StripQuotes(selectorCont.S("name").String())

									if selName == "test_check" {
										tp.Name = selName
										if selectorCont.Exists("expressions") {
											tp.Expressions = selectorCont.S("expressions").Data().([]interface{})
										}
										found = true
										break
									}
								}
							}
						}
					}
				}
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

func testAccCheckMSOSchemaTemplateAnpEpgSelectorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		cont, err := client.GetViaURL("api/v1/schemas/5ea809672c00003bc40a2799")
		if rs.Type == "mso_schema_template_anp_epg_selector" {

			if err != nil {
				return err
			}
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
				apiTemplate := models.StripQuotes(tempCont.S("name").String())
				if apiTemplate == "Template1" {

					anpCount, err := tempCont.ArrayCount("anps")
					if err != nil {
						return fmt.Errorf("Unable to get ANP list")
					}
					for j := 0; j < anpCount; j++ {
						anpCont, err := tempCont.ArrayElement(j, "anps")
						if err != nil {
							return err
						}
						apiANP := models.StripQuotes(anpCont.S("name").String())
						if apiANP == "ap1" {
							epgCount, err := anpCont.ArrayCount("epgs")
							if err != nil {
								return fmt.Errorf("Unable to get Anp Epg list")
							}
							for k := 0; k < epgCount; k++ {
								epgCont, err := anpCont.ArrayElement(k, "epgs")
								if err != nil {
									return err
								}
								apiEPG := models.StripQuotes(epgCont.S("name").String())
								if apiEPG == "epg1" {
									selectorCount, err := epgCont.ArrayCount("selectors")
									if err != nil {
										return err
									}

									for s := 0; s < selectorCount; s++ {
										selectorCont, err := epgCont.ArrayElement(s, "selectors")
										if err != nil {
											return fmt.Errorf("Unable to find a selectors")
										}
										selName := models.StripQuotes(selectorCont.S("name").String())

										if selName == "test_check" {
											return fmt.Errorf("Schema Template Anp Epg Selector still exists")
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

func testAccCheckMSOSchemaTemplateAnpEpgSelectorAttributes(key string, ss *SelectorTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		val := ss.Expressions[0].(map[string]interface{})
		if key != val["key"] {
			return fmt.Errorf("Bad Template anp epg Selector key %v", val["key"])
		}
		return nil
	}
}

type SelectorTest struct {
	Id          string        `json:",omitempty"`
	SchemaId    string        `json:",omitempty"`
	Template    string        `json:",omitempty"`
	AnpName     string        `json:",omitempty"`
	EpgName     string        `json:",omitempty"`
	Name        string        `json:",omitempty"`
	Expressions []interface{} `json:",omitempty"`
}
