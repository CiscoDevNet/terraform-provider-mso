package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccMSOSchemaTemplateAnpEpg_Basic(t *testing.T) {
	var ss TemplateAnpEpg
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateAnpEpgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateAnpEpgConfig_basic(true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateAnpEpgExists("mso_schema_template_anp_epg.anp_epg", &ss),
					testAccCheckMSOSchemaTemplateAnpEpgAttributes(true, &ss),
				),
			},
		},
	})
}

func TestAccMSOSchemaTemplateAnpEpg_Update(t *testing.T) {
	var ss TemplateAnpEpg

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateAnpEpgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateAnpEpgConfig_basic(true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateAnpEpgExists("mso_schema_template_anp_epg.anp_epg", &ss),
					testAccCheckMSOSchemaTemplateAnpEpgAttributes(true, &ss),
				),
			},
			{
				Config: testAccCheckMSOTemplateAnpEpgConfig_basic(false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateAnpEpgExists("mso_schema_template_anp_epg.anp_epg", &ss),
					testAccCheckMSOSchemaTemplateAnpEpgAttributes(false, &ss),
				),
			},
		},
	})
}

func testAccCheckMSOTemplateAnpEpgConfig_basic(preferred_group bool) string {
	return fmt.Sprintf(`
	resource "mso_schema_template_anp_epg" "anp_epg" {
		schema_id = "5c4d5bb72700000401f80948"
		template_name = "Template1"
		anp_name = "ANP"
		name = "mso_epg16"
		bd_name = "BD1"
		vrf_name = "DEVNET-VRF"
		preferred_group = %v
		display_name = "mso_epg16"
	}
`, preferred_group)
}

func testAccCheckMSOSchemaTemplateAnpEpgExists(anpName string, ss *TemplateAnpEpg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs1, err1 := s.RootModule().Resources[anpName]

		if !err1 {
			return fmt.Errorf("Anp Epg %s not found", anpName)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No EPG id was set")
		}

		cont, err := client.GetViaURL("api/v1/schemas/5c4d5bb72700000401f80948")
		if err != nil {
			return err
		}
		count, err := cont.ArrayCount("templates")
		if err != nil {
			return fmt.Errorf("No Template found")
		}
		tp := TemplateAnpEpg{}
		found := false

		for i := 0; i < count; i++ {
			tempCont, err := cont.ArrayElement(i, "templates")
			if err != nil {
				return err
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
					if apiANP == "ANP" {
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
							if apiEPG == "mso_epg16" {
								tp.name = apiEPG
								tp.displayName = models.StripQuotes(epgCont.S("displayName").String())
								tp.uSegEpg = epgCont.S("uSegEpg").Data().(bool)
								tp.preferredGroup = epgCont.S("preferredGroup").Data().(bool)
								found = true
								break
							}
						}
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

func testAccCheckMSOSchemaTemplateAnpEpgDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_template_anp_epg" {
			cont, err := client.GetViaURL("api/v1/schemas/5c4d5bb72700000401f80948")
			if err != nil {
				return err
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
							if apiANP == "ANP" {
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
									if apiEPG == "mso_epg16" {
										return fmt.Errorf("The Anp Epg still exists")
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

func testAccCheckMSOSchemaTemplateAnpEpgAttributes(preferred_group bool, ss *TemplateAnpEpg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if preferred_group != ss.preferredGroup {
			return fmt.Errorf("Bad Template Anp Epg preferred group value %v", ss.preferredGroup)
		}
		return nil
	}
}

type TemplateAnpEpg struct {
	name           string
	displayName    string
	uSegEpg        bool
	preferredGroup bool
}
