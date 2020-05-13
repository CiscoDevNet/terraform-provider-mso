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

func TestAccMSOSchemaSiteAnpEpgSubnet_Basic(t *testing.T) {
	var ss SchemaSiteAnpEpgSubnet
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpEpgSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSiteAnpEpgSubnetConfig_basic("private"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteAnpEpgSubnetExists("mso_schema_site_anp_epg_subnet.subnet1", &ss),
					testAccCheckMSOSchemaSiteAnpEpgSubnetAttributes("private", &ss),
				),
			},
		},
	})
}

func TestAccMSOSchemaSiteAnpEpgSubnet_Update(t *testing.T) {
	var ss SchemaSiteAnpEpgSubnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpEpgSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSiteAnpEpgSubnetConfig_basic("private"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteAnpEpgSubnetExists("mso_schema_site_anp_epg_subnet.subnet1", &ss),
					testAccCheckMSOSchemaSiteAnpEpgSubnetAttributes("private", &ss),
				),
			},
			{
				Config: testAccCheckMSOSiteAnpEpgSubnetConfig_basic("public"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteAnpEpgSubnetExists("mso_schema_site_anp_epg_subnet.subnet1", &ss),
					testAccCheckMSOSchemaSiteAnpEpgSubnetAttributes("public", &ss),
				),
			},
		},
	})
}

func testAccCheckMSOSiteAnpEpgSubnetConfig_basic(scope string) string {
	return fmt.Sprintf(`
   resource "mso_schema_site_anp_epg_subnet" "subnet1" {
  schema_id = "5c4d5bb72700000401f80948"
  site_id = "5c7c95b25100008f01c1ee3c"
  template_name = "Template1"
  anp_name = "ANP"
  epg_name = "DB"
  ip = "10.8.0.1/8"
  description = "SubnetEntry"
  scope = "%s"
  shared = true
  
}
`, scope)
}

func testAccCheckMSOSchemaSiteAnpEpgSubnetExists(subnetName string, ss *SchemaSiteAnpEpgSubnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs1, err1 := s.RootModule().Resources[subnetName]

		if !err1 {
			return fmt.Errorf("Entry %s not found", subnetName)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Schema id was set")
		}

		cont, err := client.GetViaURL("api/v1/schemas/5c4d5bb72700000401f80948")
		if err != nil {
			return err
		}

		count, err := cont.ArrayCount("sites")
		if err != nil {
			return fmt.Errorf("No Site found")
		}
		tp := SchemaSiteAnpEpgSubnet{}
		found := false
		for i := 0; i < count; i++ {
			tempCont, err := cont.ArrayElement(i, "sites")
			if err != nil {
				return err
			}

			apisiteId := models.StripQuotes(tempCont.S("siteId").String())
			apiTemplateName := models.StripQuotes(tempCont.S("templateName").String())
			if apiTemplateName == "Template1" && apisiteId == "5c7c95b25100008f01c1ee3c" {
				anpCount, err := tempCont.ArrayCount("anps")
				if err != nil {
					return fmt.Errorf("Unable to get ANP list")
				}
				for j := 0; j < anpCount; j++ {
					anpCont, err := tempCont.ArrayElement(j, "anps")
					if err != nil {
						return err
					}
					anpRef := models.StripQuotes(anpCont.S("anpRef").String())
					re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/anps/(.*)")
					match := re.FindStringSubmatch(anpRef)
					if match[3] == "ANP" {
						epgCount, err := anpCont.ArrayCount("epgs")
						if err != nil {
							return fmt.Errorf("Unable to get EPG list")
						}
						for k := 0; k < epgCount; k++ {
							epgCont, err := anpCont.ArrayElement(k, "epgs")
							if err != nil {
								return err
							}
							apiEpgRef := models.StripQuotes(epgCont.S("epgRef").String())
							re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/epgs/(.*)")
							match := re.FindStringSubmatch(apiEpgRef)
							apiEPG := match[3]
							if apiEPG == "DB" {
								subnetCount, err := epgCont.ArrayCount("subnets")
								if err != nil {
									return fmt.Errorf("Unable to get Static subnet list")
								}
								for l := 0; l < subnetCount; l++ {
									subnetCont, err := epgCont.ArrayElement(l, "subnets")
									if err != nil {
										return err
									}
									subnetip := "10.8.0.1/8"
									apisubnetip := models.StripQuotes(subnetCont.S("ip").String())
									if subnetip == apisubnetip {
										if subnetCont.Exists("description") {
											tp.description = models.StripQuotes(subnetCont.S("description").String())
										}
										if subnetCont.Exists("scope") {
											tp.scope = models.StripQuotes(subnetCont.S("scope").String())
										}
										if subnetCont.Exists("shared") {
											tp.shared = (subnetCont.S("shared").Data().(bool))
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
			return fmt.Errorf("Subnet Entry not found from API")
		}

		tp1 := &tp

		*ss = *tp1
		return nil
	}
}

func testAccCheckMSOSchemaSiteAnpEpgSubnetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_site_anp_epg_subnet" {
			cont, err := client.GetViaURL("api/v1/schemas/5c4d5bb72700000401f80948")
			if err != nil {
				return nil
			} else {
				count, err := cont.ArrayCount("sites")
				if err != nil {
					return fmt.Errorf("No Site found")
				}

				for i := 0; i < count; i++ {
					tempCont, err := cont.ArrayElement(i, "sites")
					if err != nil {
						return err
					}
					apisiteId := models.StripQuotes(tempCont.S("siteId").String())
					apiTemplateName := models.StripQuotes(tempCont.S("templateName").String())
					if apiTemplateName == "Template1" && apisiteId == "5c7c95b25100008f01c1ee3c" {
						anpCount, err := tempCont.ArrayCount("anps")
						if err != nil {
							return fmt.Errorf("Unable to get ANP list")
						}
						for j := 0; j < anpCount; j++ {
							anpCont, err := tempCont.ArrayElement(j, "anps")
							if err != nil {
								return err
							}
							anpRef := models.StripQuotes(anpCont.S("anpRef").String())
							re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/anps/(.*)")
							match := re.FindStringSubmatch(anpRef)
							if match[3] == "ANP" {
								epgCount, err := anpCont.ArrayCount("epgs")
								if err != nil {
									return fmt.Errorf("Unable to get EPG list")
								}
								for k := 0; k < epgCount; k++ {
									epgCont, err := anpCont.ArrayElement(k, "epgs")
									if err != nil {
										return err
									}
									apiEpgRef := models.StripQuotes(epgCont.S("epgRef").String())
									re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/epgs/(.*)")
									match := re.FindStringSubmatch(apiEpgRef)
									apiEPG := match[3]
									if apiEPG == "DB" {
										subnetCount, err := epgCont.ArrayCount("subnets")
										if err != nil {
											return fmt.Errorf("Unable to get Static subnet list")
										}
										for l := 0; l < subnetCount; l++ {
											subnetCont, err := epgCont.ArrayElement(l, "subnets")
											if err != nil {
												return err
											}
											subnetip := "10.8.0.1/8"
											apisubnetip := models.StripQuotes(subnetCont.S("ip").String())
											if subnetip == apisubnetip {
												return fmt.Errorf("The Subnet entry still exists")
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
	}
	return nil
}

func testAccCheckMSOSchemaSiteAnpEpgSubnetAttributes(scope string, ss *SchemaSiteAnpEpgSubnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "SubnetEntry" != ss.description {
			return fmt.Errorf("Bad Subnet Description value %s", ss.description)
		}

		if true != ss.shared {
			return fmt.Errorf("Bad Subnet Shared value %v", ss.shared)
		}
		return nil
	}
}

type SchemaSiteAnpEpgSubnet struct {
	description string
	scope       string
	shared      bool
}
