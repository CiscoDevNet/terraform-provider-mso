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

func TestAccMSOSchemaSchemaSiteExternalEpg(t *testing.T) {
	var ss SchemaSiteExternalEpg
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSchemaSiteExternalEpgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSchemaSiteExternalEpgConfig_basic("demo"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSchemaSiteExternalEpgExists("mso_schema_template_external_epg.template_externalepg", &ss),
					testAccCheckMSOSchemaSchemaSiteExternalEpgAttributes("demo", &ss),
				),
			},
		},
	})
}

func TestAccMSOSchemaSchemaSiteExternalEpg_Update(t *testing.T) {
	var ss SchemaSiteExternalEpg

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSchemaSiteExternalEpgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSchemaSiteExternalEpgConfig_basic("demo"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSchemaSiteExternalEpgExists("mso_schema_template_external_epg.template_externalepg", &ss),
					testAccCheckMSOSchemaSchemaSiteExternalEpgAttributes("demo", &ss),
				),
			},
			{
				Config: testAccCheckMSOSchemaSiteExternalEpgConfig_basic("vrf1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSchemaSiteExternalEpgExists("mso_schema_template_external_epg.template_externalepg", &ss),
					testAccCheckMSOSchemaSchemaSiteExternalEpgAttributes("vrf1", &ss),
				),
			},
		},
	})
}

func testAccCheckMSOSchemaSiteExternalEpgConfig_basic(name string) string {
	return fmt.Sprintf(`
	resource "mso_schema_template_external_epg" "template_externalepg" {
		schema_id = "5ea809672c00003bc40a2799"
		template_name = "Template1"
		external_epg_name = "external_epg12"
		l3out_name = "%v"
	  }
`, name)
}

func testAccCheckMSOSchemaSchemaSiteExternalEpgExists(externalepgName string, ss *SchemaSiteExternalEpg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs1, err1 := s.RootModule().Resources[externalepgName]

		if !err1 {
			return fmt.Errorf("External Epg %s not found", externalepgName)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Schema id was set")
		}

		cont, err := client.GetViaURL("api/v1/schemas/5ea809672c00003bc40a2799")
		if err != nil {
			return err
		}
		count, err := cont.ArrayCount("templates")
		if err != nil {
			return fmt.Errorf("No Template found")
		}
		tp := SchemaSiteExternalEpg{}
		found := false
		for i := 0; i < count; i++ {
			tempCont, err := cont.ArrayElement(i, "templates")
			if err != nil {
				return err
			}

			apiTemplateName := models.StripQuotes(tempCont.S("name").String())
			if apiTemplateName == "Template1" {
				externalepgCount, err := tempCont.ArrayCount("externalEpgs")
				if err != nil {
					return fmt.Errorf("Unable to get External Epg list")
				}
				for j := 0; j < externalepgCount; j++ {
					externalepgCont, err := tempCont.ArrayElement(j, "externalEpgs")
					if err != nil {
						return err
					}
					apiExternalepg := models.StripQuotes(externalepgCont.S("name").String())
					if apiExternalepg == "external_epg12" {
						l3outRef := models.StripQuotes(externalepgCont.S("l3outRef").String())
						re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/l3outs/(.*)")
						match := re.FindStringSubmatch(l3outRef)
						tp.l3out_name = match[3]
						found = true
						break
					}
				}
			}
		}
		if !found {
			return fmt.Errorf("External Epg not found from API")
		}

		tp1 := &tp
		*ss = *tp1
		return nil
	}
}

func testAccCheckMSOSchemaSchemaSiteExternalEpgDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_template_external_epg" {
			cont, err := client.GetViaURL("api/v1/schemas/5ea809672c00003bc40a2799")
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
						externalepgCount, err := tempCont.ArrayCount("externalEpgs")
						if err != nil {
							return fmt.Errorf("Unable to get External epg list")
						}
						for j := 0; j < externalepgCount; j++ {
							externalepgCont, err := tempCont.ArrayElement(j, "externalEpgs")
							if err != nil {
								return err
							}
							apiExternalepg := models.StripQuotes(externalepgCont.S("name").String())
							if apiExternalepg == "external_epg12" {
								return fmt.Errorf("template External Epg still exists.")
							}
						}
					}
				}
			}
		}
	}
	return nil
}
func testAccCheckMSOSchemaSchemaSiteExternalEpgAttributes(l3out_name string, ss *SchemaSiteExternalEpg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if l3out_name != ss.l3out_name {
			return fmt.Errorf("Bad Template External epg L3Out name %s", ss.l3out_name)
		}
		return nil
	}
}

type SchemaSiteExternalEpg struct {
	l3out_name string
}
