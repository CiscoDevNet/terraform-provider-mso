package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccMSOSchemaTemplateAnpEpgUsegAttr_Basic(t *testing.T) {
	var ss UsegAttrTest
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateAnpEpgUsegAttrDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateAnpEpgUsegAttrConfig_basic("10.2.3.4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateAnpEpgUsegAttrExists("mso_schema_template_anp_epg_useg_attr.useg_attrs", &ss),
					testAccCheckMSOSchemaTemplateAnpEpgUsegAttrAttributes("10.2.3.4", &ss),
				),
			},
		},
	})
}

func TestAccMSOSchemaTemplateAnpEpgUsegAttr_Update(t *testing.T) {
	var ss UsegAttrTest

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateAnpEpgUsegAttrDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateAnpEpgUsegAttrConfig_basic("10.2.3.4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateAnpEpgUsegAttrExists("mso_schema_template_anp_epg_useg_attr.useg_attrs", &ss),
					testAccCheckMSOSchemaTemplateAnpEpgUsegAttrAttributes("10.2.3.4", &ss),
				),
			},
			{
				Config: testAccCheckMSOTemplateAnpEpgUsegAttrConfig_basic("10.2.3.5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateAnpEpgUsegAttrExists("mso_schema_template_anp_epg_useg_attr.useg_attrs", &ss),
					testAccCheckMSOSchemaTemplateAnpEpgUsegAttrAttributes("10.2.3.5", &ss),
				),
			},
		},
	})
}

func testAccCheckMSOTemplateAnpEpgUsegAttrConfig_basic(val string) string {
	return fmt.Sprintf(`
	resource "mso_schema_template_anp_epg_useg_attr" "useg_attrs" {
		schema_id     = "5eafca7d2c000052860a2902"
		anp_name      = "sanp1"
		epg_name      = "nkuseg"
		template_name = "stemplate1"
		name          = "usg_acc_test"
		useg_type     = "tag"
		operator      = "startsWith"
		category      = "tagger"
		value         = "%s"
		useg_subnet   = true
		
	}
`, val)
}

func testAccCheckMSOSchemaTemplateAnpEpgUsegAttrExists(usegName string, ss *UsegAttrTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs1, err1 := s.RootModule().Resources[usegName]

		if !err1 {
			return fmt.Errorf("Useg %s not found", usegName)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Useg id was set")
		}

		cont, err := client.GetViaURL("api/v1/schemas/5eafca7d2c000052860a2902")
		if err != nil {
			return err
		}
		count, err := cont.ArrayCount("templates")
		if err != nil {
			return fmt.Errorf("No Template found")
		}
		tp := UsegAttrTest{}
		found := false

		for i := 0; i < count; i++ {
			tempCont, err := cont.ArrayElement(i, "templates")
			if err != nil {
				return fmt.Errorf("No Template found")
			}

			apiTemplate := models.StripQuotes(tempCont.S("name").String())

			if apiTemplate == "stemplate1" {
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
					if apiANP == "sanp1" {
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
							if apiEPG == "nkuseg" {
								tp.EpgName = apiEPG

								usegCount, err := epgCont.ArrayCount("uSegAttrs")
								if err != nil {
									return fmt.Errorf("Unable to get useg Attrs")
								}

								for s := 0; s < usegCount; s++ {
									usegCont, err := epgCont.ArrayElement(s, "uSegAttrs")
									if err != nil {
										return err
									}

									apiName := models.StripQuotes(usegCont.S("name").String())

									if apiName == "usg_acc_test" {
										tp.Name = apiName
										tp.Operator = models.StripQuotes(usegCont.S("operator").String())
										tp.Value = models.StripQuotes(usegCont.S("value").String())

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
			return fmt.Errorf("UsegAttr not found from API")
		}
		tp1 := &tp
		*ss = *tp1
		return nil
	}
}

func testAccCheckMSOSchemaTemplateAnpEpgUsegAttrDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	cont, err := client.GetViaURL("api/v1/schemas/5eafca7d2c000052860a2902")
	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_template_anp_epg_useg_attr" {

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
				if apiTemplate == "stemplate1" {

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
						if apiANP == "sanp1" {
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
								if apiEPG == "nkuseg" {
									usegCount, err := epgCont.ArrayCount("uSegAttrs")
									if err != nil {
										return err
									}

									for s := 0; s < usegCount; s++ {
										usegCont, err := epgCont.ArrayElement(s, "uSegAttrs")
										if err != nil {
											return fmt.Errorf("Unable to find a useg Attrs")
										}
										currentIp := models.StripQuotes(usegCont.S("name").String())

										if currentIp == "usg_acc_test" {
											return fmt.Errorf("Schema Template Anp Epg useg Attrs still exists")
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

func testAccCheckMSOSchemaTemplateAnpEpgUsegAttrAttributes(val string, ss *UsegAttrTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != ss.Value {
			return fmt.Errorf("Bad Useg Attrs value %v", ss.Value)
		}
		return nil
	}
}

type UsegAttrTest struct {
	Id       string `json:",omitempty"`
	SchemaId string `json:",omitempty"`
	Template string `json:",omitempty"`
	AnpName  string `json:",omitempty"`
	EpgName  string `json:",omitempty"`
	Name     string `json:",omitempty"`
	Operator string `json:",omitempty"`
	Value    string `json:",omitempty"`
}
