package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccMSOSchemaTemplateExternalepgSubnet_Basic(t *testing.T) {
	var ss TemplateExternalepgSubnet
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateExternalepgSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateExternalepgSubnetConfig_basic("sub1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateExternalepgSubnetExists("mso_schema_template_externalepg_subnet.subnet1", &ss),
					testAccCheckMSOSchemaTemplateExternalepgSubnetAttributes("sub1", &ss),
				),
			},
		},
	})
}

func TestAccMSOSchemaTemplateExternalepgSubnet_Update(t *testing.T) {
	var ss TemplateExternalepgSubnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateExternalepgSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateExternalepgSubnetConfig_basic("sub1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateExternalepgSubnetExists("mso_schema_template_externalepg_subnet.subnet1", &ss),
					testAccCheckMSOSchemaTemplateExternalepgSubnetAttributes("sub1", &ss),
				),
			},
			{
				Config: testAccCheckMSOTemplateExternalepgSubnetConfig_basic("sub2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateExternalepgSubnetExists("mso_schema_template_externalepg_subnet.subnet1", &ss),
					testAccCheckMSOSchemaTemplateExternalepgSubnetAttributes("sub2", &ss),
				),
			},
		},
	})
}

func testAccCheckMSOTemplateExternalepgSubnetConfig_basic(name string) string {
	return fmt.Sprintf(`
	resource "mso_schema_template_externalepg_subnet" "subnet1" {
		schema_id = "5ea809672c00003bc40a2799"
		template_name = "Template1"
		externalepg_name =  "UntitledExternalEPG1"
		ip = "10.101.100.0/25"
		name = "%v"
		scope = ["shared-rtctrl"]
		aggregate = ["shared-rtctrl"]
	  }
`, name)
}

func testAccCheckMSOSchemaTemplateExternalepgSubnetExists(externalepgName string, ss *TemplateExternalepgSubnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs1, err1 := s.RootModule().Resources[externalepgName]

		if !err1 {
			return fmt.Errorf("External Epg Subnet %s not found", externalepgName)
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
		tp := TemplateExternalepgSubnet{}
		found := false
		for i := 0; i < count; i++ {
			tempCont, err := cont.ArrayElement(i, "templates")
			if err != nil {
				return err
			}

			apiTemplate := models.StripQuotes(tempCont.S("name").String())
			if apiTemplate == "Template1" {
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
					if apiExternalepg == "UntitledExternalEPG1" {
						subnetCount, err := externalepgCont.ArrayCount("subnets")
						if err != nil {
							return fmt.Errorf("Unable to get Subnets list")
						}
						for k := 0; k < subnetCount; k++ {
							subnetsCont, err := externalepgCont.ArrayElement(k, "subnets")
							if err != nil {
								return err
							}
							apiIP := models.StripQuotes(subnetsCont.S("ip").String())
							if apiIP == "10.101.100.0/25" {
								tp.ip = apiIP
								tp.name = models.StripQuotes(subnetsCont.S("name").String())
							}
						}
						found = true
						break
					}
				}
			}
		}
		if !found {
			return fmt.Errorf("External Epg Subnet not found from API")
		}

		tp1 := &tp
		*ss = *tp1
		return nil
	}
}

func testAccCheckMSOSchemaTemplateExternalepgSubnetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_template_externalepg_subnet" {
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
							epgCont, err := tempCont.ArrayElement(j, "externalEpgs")
							if err != nil {
								return err
							}
							apiExternalepg := models.StripQuotes(epgCont.S("name").String())
							if apiExternalepg == "UntitledExternalEPG1" {
								subnetCount, err := epgCont.ArrayCount("subnets")
								if err != nil {
									return fmt.Errorf("Unable to get External Epg Subnets list")
								}
								for k := 0; k < subnetCount; k++ {
									subnetCont, err := epgCont.ArrayElement(k, "subnets")
									if err != nil {
										return err
									}
									apiIP := models.StripQuotes(subnetCont.S("ip").String())
									if apiIP == "10.101.100.0/25" {
										return fmt.Errorf("External Epg Subnet still exists")
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
func testAccCheckMSOSchemaTemplateExternalepgSubnetAttributes(name string, ss *TemplateExternalepgSubnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if name != ss.name {
			return fmt.Errorf("Bad Template External Epg Subnet Relationship Type %s", ss.name)
		}

		return nil
	}
}

type TemplateExternalepgSubnet struct {
	name string
	ip   string
}
