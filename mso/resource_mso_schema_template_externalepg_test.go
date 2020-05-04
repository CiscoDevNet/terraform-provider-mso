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

func TestAccMSOSchemaTemplateExternalepg_Basic(t *testing.T) {
	var ss TemplateExternalepg
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateExternalepgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateExternalepgConfig_basic("demo"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateExternalepgExists("mso_schema_template_externalepg.template_externalepg", &ss),
					testAccCheckMSOSchemaTemplateExternalepgAttributes("demo", &ss),
				),
			},
		},
	})
}

func TestAccMSOSchemaTemplateExternalepg_Update(t *testing.T) {
	var ss TemplateExternalepg

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateExternalepgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateExternalepgConfig_basic("demo"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateExternalepgExists("mso_schema_template_externalepg.template_externalepg", &ss),
					testAccCheckMSOSchemaTemplateExternalepgAttributes("demo", &ss),
				),
			},
			{
				Config: testAccCheckMSOTemplateExternalepgConfig_basic("vrf1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateExternalepgExists("mso_schema_template_externalepg.template_externalepg", &ss),
					testAccCheckMSOSchemaTemplateExternalepgAttributes("vrf1", &ss),
				),
			},
		},
	})
}

func testAccCheckMSOTemplateExternalepgConfig_basic(name string) string {
	return fmt.Sprintf(`
	resource "mso_schema_template_externalepg" "template_externalepg" {
		schema_id = "5ea809672c00003bc40a2799"
		template_name = "Template1"
		externalepg_name = "external_epg12"
		display_name = "external_epg12"
		vrf_name = "%v"
	  }
`, name)
}

func testAccCheckMSOSchemaTemplateExternalepgExists(externalepgName string, ss *TemplateExternalepg) resource.TestCheckFunc {
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
		tp := TemplateExternalepg{}
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
						tp.display_name = models.StripQuotes(externalepgCont.S("displayName").String())
						vrfRef := models.StripQuotes(externalepgCont.S("vrfRef").String())
						re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/vrfs/(.*)")
						match := re.FindStringSubmatch(vrfRef)
						tp.vrf_name = match[3]
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

func testAccCheckMSOSchemaTemplateExternalepgDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_template_externalepg" {
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
func testAccCheckMSOSchemaTemplateExternalepgAttributes(vrf_name string, ss *TemplateExternalepg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "external_epg12" != ss.display_name {
			return fmt.Errorf("Bad Template External epg display name %s", ss.display_name)
		}

		if vrf_name != ss.vrf_name {
			return fmt.Errorf("Bad Template External epg VRF name %s", ss.vrf_name)
		}
		return nil
	}
}

type TemplateExternalepg struct {
	display_name string
	vrf_name     string
}
