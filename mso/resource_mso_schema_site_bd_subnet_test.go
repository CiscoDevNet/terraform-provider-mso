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

func TestAccMSOSchemaSiteBdSubnet_Basic(t *testing.T) {
	var ss SchemaSiteBdSubnet
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteBdSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSiteBdSubnetConfig_basic("private"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteBdSubnetExists("mso_schema_site_bd_subnet.sub1", &ss),
					testAccCheckMSOSchemaSiteBdSubnetAttributes("private", &ss),
				),
			},
		},
	})
}

func TestAccMSOSchemaSiteBdSubnet_Update(t *testing.T) {
	var ss SchemaSiteBdSubnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteBdSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSiteBdSubnetConfig_basic("private"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteBdSubnetExists("mso_schema_site_bd_subnet.sub1", &ss),
					testAccCheckMSOSchemaSiteBdSubnetAttributes("private", &ss),
				),
			},
			{
				Config: testAccCheckMSOSiteBdSubnetConfig_basic("public"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteBdSubnetExists("mso_schema_site_bd_subnet.sub1", &ss),
					testAccCheckMSOSchemaSiteBdSubnetAttributes("public", &ss),
				),
			},
		},
	})
}

func testAccCheckMSOSiteBdSubnetConfig_basic(scope string) string {
	return fmt.Sprintf(`
	resource "mso_schema_site_bd_subnet" "sub1" {
		schema_id = "5d5dbf3f2e0000580553ccce"
		template_name = "Template1"
		site_id = "5c7c95b25100008f01c1ee3c"
		bd_name = "WebServer-Finance"
		ip = "200.168.240.1/24"
		description = "Subnet 1"
		shared = false
		scope = "%s"
		querier = false
		no_default_gateway = false
	  
	  }
`, scope)
}

func testAccCheckMSOSchemaSiteBdSubnetExists(subnetName string, ss *SchemaSiteBdSubnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs1, err1 := s.RootModule().Resources[subnetName]

		if !err1 {
			return fmt.Errorf("Entry %s not found", subnetName)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Schema id was set")
		}

		cont, err := client.GetViaURL("api/v1/schemas/5d5dbf3f2e0000580553ccce")
		if err != nil {
			return err
		}

		count, err := cont.ArrayCount("sites")
		if err != nil {
			return fmt.Errorf("No Site found")
		}
		tp := SchemaSiteBdSubnet{}
		found := false
		for i := 0; i < count; i++ {
			tempCont, err := cont.ArrayElement(i, "sites")
			if err != nil {
				return err
			}

			apisiteId := models.StripQuotes(tempCont.S("siteId").String())
			apiTemplateName := models.StripQuotes(tempCont.S("templateName").String())
			if apiTemplateName == "Template1" && apisiteId == "5c7c95b25100008f01c1ee3c" {
				bdCount, err := tempCont.ArrayCount("bds")
				if err != nil {
					return fmt.Errorf("Unable to get bd list")
				}
				for j := 0; j < bdCount; j++ {
					bdCont, err := tempCont.ArrayElement(j, "bds")
					if err != nil {
						return err
					}
					bdRef := models.StripQuotes(bdCont.S("bdRef").String())
					re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/bds/(.*)")
					match := re.FindStringSubmatch(bdRef)
					if match[3] == "WebServer-Finance" {
						subnetCount, err := bdCont.ArrayCount("subnets")
						if err != nil {
							return fmt.Errorf("Unable to get Static subnet list")
						}
						for l := 0; l < subnetCount; l++ {
							subnetCont, err := bdCont.ArrayElement(l, "subnets")
							if err != nil {
								return err
							}
							subnetip := "200.168.240.1/24"
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

		if !found {
			return fmt.Errorf("Subnet Entry not found from API")
		}

		tp1 := &tp

		*ss = *tp1
		return nil
	}
}

func testAccCheckMSOSchemaSiteBdSubnetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_site_bd_subnet" {
			cont, err := client.GetViaURL("api/v1/schemas/5d5dbf3f2e0000580553ccce")
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
						bdCount, err := tempCont.ArrayCount("bds")
						if err != nil {
							return fmt.Errorf("Unable to get bd list")
						}
						for j := 0; j < bdCount; j++ {
							bdCont, err := tempCont.ArrayElement(j, "bds")
							if err != nil {
								return err
							}
							bdRef := models.StripQuotes(bdCont.S("bdRef").String())
							re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/bds/(.*)")
							match := re.FindStringSubmatch(bdRef)
							if match[3] == "WebServer-Finance" {
								subnetCount, err := bdCont.ArrayCount("subnets")
								if err != nil {
									return fmt.Errorf("Unable to get Static subnet list")
								}
								for l := 0; l < subnetCount; l++ {
									subnetCont, err := bdCont.ArrayElement(l, "subnets")
									if err != nil {
										return err
									}
									subnetip := "200.168.240.1/24"
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

	return nil
}

func testAccCheckMSOSchemaSiteBdSubnetAttributes(scope string, ss *SchemaSiteBdSubnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "Subnet 1" != ss.description {
			return fmt.Errorf("Bad Subnet Description value %s", ss.description)
		}

		if false != ss.shared {
			return fmt.Errorf("Bad Subnet Shared value %v", ss.shared)
		}
		return nil
	}
}

type SchemaSiteBdSubnet struct {
	description string
	scope       string
	shared      bool
}
