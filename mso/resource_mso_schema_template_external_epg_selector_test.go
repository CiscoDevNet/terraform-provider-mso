package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccMSOSchemaTemplateExternalEpgSelector_Basic(t *testing.T) {
	var ss EPGSelectorTest
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateExternalEpgSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateExternalEpgSelectorConfig_basic("1.2.3.4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateExternalEpgSelectorExists("mso_schema_template_external_epg_selector.selector1", &ss),
					testAccCheckMSOSchemaTemplateExternalEpgSelectorAttributes("1.2.3.4", &ss),
				),
			},
		},
	})
}

func TestAccMSOSchemaTemplateExternalEpgSelector_Update(t *testing.T) {
	var ss EPGSelectorTest

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateExternalEpgSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateExternalEpgSelectorConfig_basic("1.2.3.4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateExternalEpgSelectorExists("mso_schema_template_external_epg_selector.selector1", &ss),
					testAccCheckMSOSchemaTemplateExternalEpgSelectorAttributes("1.2.3.4", &ss),
				),
			},
			{
				Config: testAccCheckMSOTemplateExternalEpgSelectorConfig_basic("5.4.6.7"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateExternalEpgSelectorExists("mso_schema_template_external_epg_selector.selector1", &ss),
					testAccCheckMSOSchemaTemplateExternalEpgSelectorAttributes("5.4.6.7", &ss),
				),
			},
		},
	})
}

func testAccCheckMSOTemplateExternalEpgSelectorConfig_basic(ip string) string {
	return fmt.Sprintf(`
	resource "mso_schema_template_external_epg_selector" "selector1" {
		schema_id = "5ea809672c00003bc40a2799"
		template = "Template1"
		external_epg_name = "check_anp01"
		name = "test_check"
    	expressions {
      		value = "%s"
    	}
	}
`, ip)
}

func testAccCheckMSOSchemaTemplateExternalEpgSelectorExists(selectorName string, ss *EPGSelectorTest) resource.TestCheckFunc {
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
		tp := EPGSelectorTest{}
		found := false

		for i := 0; i < count; i++ {
			tempCont, err := cont.ArrayElement(i, "templates")
			if err != nil {
				return fmt.Errorf("Error fetching template")
			}

			tempName := models.StripQuotes(tempCont.S("name").String())
			if tempName == "Template1" {
				extrEpgCount, err := tempCont.ArrayCount("externalEpgs")
				if err != nil {
					return fmt.Errorf("no externalEpgs found")
				}

				for j := 0; j < extrEpgCount; j++ {
					extrEpgCont, err := tempCont.ArrayElement(j, "externalEpgs")
					if err != nil {
						return fmt.Errorf("Error fetching external Epg")
					}

					extrEpgName := models.StripQuotes(extrEpgCont.S("name").String())
					if extrEpgName == "check_anp01" {
						selectorCount, err := extrEpgCont.ArrayCount("selectors")
						if err != nil {
							return fmt.Errorf("No selectors found")
						}

						for k := 0; k < selectorCount; k++ {
							selectorCont, err := extrEpgCont.ArrayElement(k, "selectors")
							if err != nil {
								return fmt.Errorf("Error fetching selector")
							}

							selectorName := models.StripQuotes(selectorCont.S("name").String())
							if selectorName == "test_check" {
								found = true
								tp.Name = selectorName
								if selectorCont.Exists("expressions") {
									tp.Expressions = selectorCont.S("expressions").Data().([]interface{})
								}
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

func testAccCheckMSOSchemaTemplateExternalEpgSelectorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		cont, err := client.GetViaURL("api/v1/schemas/5ea809672c00003bc40a2799")
		if rs.Type == "mso_schema_template_external_epg_selector" {

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
					return fmt.Errorf("Error fetching template")
				}

				tempName := models.StripQuotes(tempCont.S("name").String())
				if tempName == "Template1" {
					extrEpgCount, err := tempCont.ArrayCount("externalEpgs")
					if err != nil {
						return fmt.Errorf("no externalEpgs found")
					}

					for j := 0; j < extrEpgCount; j++ {
						extrEpgCont, err := tempCont.ArrayElement(j, "externalEpgs")
						if err != nil {
							return fmt.Errorf("Error fetching external Epg")
						}

						extrEpgName := models.StripQuotes(extrEpgCont.S("name").String())
						if extrEpgName == "check_anp01" {
							selectorCount, err := extrEpgCont.ArrayCount("selectors")
							if err != nil {
								return fmt.Errorf("No selectors found")
							}

							for k := 0; k < selectorCount; k++ {
								selectorCont, err := extrEpgCont.ArrayElement(k, "selectors")
								if err != nil {
									return fmt.Errorf("Error fetching selector")
								}

								selectorName := models.StripQuotes(selectorCont.S("name").String())
								if selectorName == "test_check" {
									return fmt.Errorf("Schema Template external epg selector still exist")
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

func testAccCheckMSOSchemaTemplateExternalEpgSelectorAttributes(ip string, ss *EPGSelectorTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		val := ss.Expressions[0].(map[string]interface{})
		if ip != val["value"] {
			return fmt.Errorf("Bad Template anp epg Selector value %v", val["value"])
		}
		return nil
	}
}

type EPGSelectorTest struct {
	Id          string        `json:",omitempty"`
	SchemaId    string        `json:",omitempty"`
	Template    string        `json:",omitempty"`
	EpgName     string        `json:",omitempty"`
	Name        string        `json:",omitempty"`
	Expressions []interface{} `json:",omitempty"`
}
