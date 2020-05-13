package mso

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccMSOSchemaTemplateAnpEpgSubnet_Basic(t *testing.T) {
	var ss SubnetTest
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateAnpEpgSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateAnpEpgSubnetConfig_basic(true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateAnpEpgSubnetExists("mso_schema_template_anp_epg_subnet.subnet1", &ss),
					testAccCheckMSOSchemaTemplateAnpEpgSubnetAttributes(true, &ss),
				),
			},
		},
	})
}

func TestAccMSOSchemaTemplateAnpEpgSubnet_Update(t *testing.T) {
	var ss SubnetTest

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateAnpEpgSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateAnpEpgSubnetConfig_basic(true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateAnpEpgSubnetExists("mso_schema_template_anp_epg_subnet.subnet1", &ss),
					testAccCheckMSOSchemaTemplateAnpEpgSubnetAttributes(true, &ss),
				),
			},
			{
				Config: testAccCheckMSOTemplateAnpEpgSubnetConfig_basic(false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateAnpEpgSubnetExists("mso_schema_template_anp_epg_subnet.subnet1", &ss),
					testAccCheckMSOSchemaTemplateAnpEpgSubnetAttributes(false, &ss),
				),
			},
		},
	})
}

func testAccCheckMSOTemplateAnpEpgSubnetConfig_basic(shared bool) string {
	return fmt.Sprintf(`
	resource "mso_schema_template_anp_epg_subnet" "subnet1" {
		schema_id = "5c6c16d7270000c710f8094d"
		anp_name = "WoS-Cloud-Only-2"
		epg_name ="DB"
		template = "Template1"
		ip = "99.101.102.0/8"
		scope = "private"
		shared = "%v"
		}
`, shared)
}

func testAccCheckMSOSchemaTemplateAnpEpgSubnetExists(subnetName string, ss *SubnetTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs1, err1 := s.RootModule().Resources[subnetName]

		if !err1 {
			return fmt.Errorf("Subnet %s not found", subnetName)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Subnet id was set")
		}

		cont, err := client.GetViaURL("api/v1/schemas/5c6c16d7270000c710f8094d")
		if err != nil {
			return err
		}
		count, err := cont.ArrayCount("templates")
		if err != nil {
			return fmt.Errorf("No Template found")
		}
		tp := SubnetTest{}
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
					if apiANP == "WoS-Cloud-Only-2" {
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
							if apiEPG == "DB" {
								tp.EpgName = apiEPG

								subnetCount, err := epgCont.ArrayCount("subnets")
								if err != nil {
									return fmt.Errorf("Unable to get subnetlist")
								}

								for s := 0; s < subnetCount; s++ {
									subnetCont, err := epgCont.ArrayElement(s, "subnets")
									if err != nil {
										return err
									}

									apiIp := models.StripQuotes(subnetCont.S("ip").String())

									if apiIp == "99.101.102.0/8" {
										tp.Ip = apiIp
										if subnetCont.Exists("scope") {
											tp.Scope = models.StripQuotes(subnetCont.S("scope").String())
										}
										if subnetCont.Exists("shared") {
											shared, _ := strconv.ParseBool(models.StripQuotes(subnetCont.S("shared").String()))
											tp.Shared = shared
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
			return fmt.Errorf("Subnet not found from API")
		}
		tp1 := &tp
		*ss = *tp1
		return nil
	}
}

func testAccCheckMSOSchemaTemplateAnpEpgSubnetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		cont, err := client.GetViaURL("api/v1/schemas/5c4d5bb72700000401f80948")
		if rs.Type == "mso_schema_template_anp_epg_subnet" {

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
						if apiANP == "WoS-Cloud-Only-2" {
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
								if apiEPG == "DB" {
									subnetCount, err := epgCont.ArrayCount("subnets")
									if err != nil {
										return err
									}

									for s := 0; s < subnetCount; s++ {
										subnetCont, err := epgCont.ArrayElement(s, "subnets")
										if err != nil {
											return fmt.Errorf("Unable to find a subnets")
										}
										currentIp := models.StripQuotes(subnetCont.S("ip").String())

										if currentIp == "99.101.102.0/8" {
											return fmt.Errorf("Schema Template Anp Epg Ip still exists")
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

func testAccCheckMSOSchemaTemplateAnpEpgSubnetAttributes(shared bool, ss *SubnetTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if shared != ss.Shared {
			return fmt.Errorf("Bad Template Subnet shared value %v", ss.Shared)
		}
		return nil
	}
}

type SubnetTest struct {
	Id       string `json:",omitempty"`
	SchemaId string `json:",omitempty"`
	Template string `json:",omitempty"`
	AnpName  string `json:",omitempty"`
	EpgName  string `json:",omitempty"`
	Ip       string `json:",omitempty"`
	Scope    string `json:",omitempty"`
	Shared   bool   `json:",omitempty"`
}
