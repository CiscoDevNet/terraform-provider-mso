package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccMSOSchemaTemplateBDSubnet_Basic(t *testing.T) {
	var ss TemplateBDSubnet
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateBDSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateBDSubnetConfig_basic(true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateBDSubnetExists("mso_schema.schema1", "mso_schema_template_bd_subnet.subnet", &ss),
					testAccCheckMSOSchemaTemplateBDSubnetAttributes(true, &ss),
				),
			},
			{
				ResourceName:      "mso_schema_template_bd_subnet.subnet",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMSOSchemaTemplateBDSubnet_Update(t *testing.T) {
	var ss TemplateBDSubnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateBDSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateBDSubnetConfig_basic(true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateBDSubnetExists("mso_schema.schema1", "mso_schema_template_bd_subnet.subnet", &ss),
					testAccCheckMSOSchemaTemplateBDSubnetAttributes(true, &ss),
				),
			},
			{
				Config: testAccCheckMSOTemplateBDSubnetConfig_basic(false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateBDSubnetExists("mso_schema.schema1", "mso_schema_template_bd_subnet.subnet", &ss),
					testAccCheckMSOSchemaTemplateBDSubnetAttributes(false, &ss),
				),
			},
		},
	})
}

func testAccCheckMSOTemplateBDSubnetConfig_basic(shared bool) string {
	return fmt.Sprintf(`

	resource "mso_schema" "schema1" {
		name          = "Schema2"
		template_name = "Template1"
		tenant_id     = "5fb5fed8520000452a9e8911"
	  
	  }

	  resource "mso_schema_template_vrf" "vrf1" {
		schema_id=mso_schema.schema1.id
		template=mso_schema.schema1.template_name
		name= "VRF"
		display_name="vrf1"
		layer3_multicast=true
		vzany=false
	  }

	resource "mso_schema_template_bd" "bridge_domain" {
		schema_id = mso_schema.schema1.id
		template_name = "Template1"
		name = "BD"
		display_name = "bd1"
		layer2_stretch = true
		intersite_bum_traffic = true
		vrf_name = mso_schema_template_vrf.vrf1.name
	}
	resource "mso_schema_template_bd_subnet" "subnet" {
		schema_id = mso_schema.schema1.id
		template_name = "Template1"
		bd_name = mso_schema_template_bd.bridge_domain.name
		ip = "13.1.1.0/8"
		scope = "private"
		shared = "%t"
		no_default_gateway = true
		querier = true
	}
`, shared)
}

func testAccCheckMSOSchemaTemplateBDSubnetExists(schemaName string, bdName string, ss *TemplateBDSubnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs1, err1 := s.RootModule().Resources[schemaName]
		rs2, err2 := s.RootModule().Resources[bdName]
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Schema id was set")
		}
		if !err1 {
			return fmt.Errorf("Schema %s not found", schemaName)
		}

		if !err2 {
			return fmt.Errorf("BD Subnet %s not found", bdName)
		}
		if rs2.Primary.ID == "" {
			return fmt.Errorf("No Subnet id was set")
		}

		cont, err := client.GetViaURL("api/v1/schemas/" + rs1.Primary.ID)
		if err != nil {
			return err
		}
		count, err := cont.ArrayCount("templates")
		if err != nil {
			return fmt.Errorf("No Template found")
		}
		tp := TemplateBDSubnet{}
		found := false
		for i := 0; i < count; i++ {
			tempCont, err := cont.ArrayElement(i, "templates")
			if err != nil {
				return err
			}

			apiTemplateName := models.StripQuotes(tempCont.S("name").String())
			if apiTemplateName == "Template1" {
				bdCount, err := tempCont.ArrayCount("bds")
				if err != nil {
					return fmt.Errorf("Unable to get BD list")
				}
				for j := 0; j < bdCount; j++ {
					bdCont, err := tempCont.ArrayElement(j, "bds")
					if err != nil {
						return err
					}
					apiBD := models.StripQuotes(bdCont.S("name").String())
					if apiBD == "BD" {
						bdsubnetCount, err := bdCont.ArrayCount("subnets")
						if err != nil {
							return fmt.Errorf("Unable to get BD subnet list")
						}
						for k := 0; k < bdsubnetCount; k++ {
							subnetCont, err := bdCont.ArrayElement(k, "subnets")
							if err != nil {
								return err
							}
							apiIP := models.StripQuotes(subnetCont.S("ip").String())
							if apiIP == "13.1.1.0/8" {

								tp.ip = apiIP
								tp.scope = models.StripQuotes(subnetCont.S("scope").String())
								tp.shared = subnetCont.S("shared").Data().(bool)
								tp.no_default_gateway = subnetCont.S("noDefaultGateway").Data().(bool)
								tp.querier = subnetCont.S("querier").Data().(bool)
								found = true
								break

							}
						}

					}
				}
			}
		}

		if !found {
			return fmt.Errorf("BD Subnet not found from API")
		}

		tp1 := &tp

		*ss = *tp1
		return nil
	}
}

func testAccCheckMSOSchemaTemplateBDSubnetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	rs1, err1 := s.RootModule().Resources["mso_schema.schema1"]
	if !err1 {
		return fmt.Errorf("Schema %s not found", "mso_schema.schema1")
	}

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_template_bd_subnet" {
			cont, err := client.GetViaURL("api/v1/schemas/" + rs1.Primary.ID)
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
						bdCount, err := tempCont.ArrayCount("bds")

						if err != nil {
							return fmt.Errorf("Unable to get BD list")
						}
						for j := 0; j < bdCount; j++ {
							bdCont, err := tempCont.ArrayElement(j, "bds")

							if err != nil {
								return err
							}
							apiBD := models.StripQuotes(bdCont.S("name").String())
							if apiBD == "BD" {
								bdsubnetCount, err := bdCont.ArrayCount("subnets")
								if err != nil {
									return fmt.Errorf("Unable to get BD subnet list")
								}
								for k := 0; k < bdsubnetCount; k++ {
									subnetCont, err := bdCont.ArrayElement(k, "subnets")
									if err != nil {
										return err
									}
									apiIP := models.StripQuotes(subnetCont.S("ip").String())
									if apiIP == "13.1.1.0/8" {
										return fmt.Errorf("The BD Subnet still exists")
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

func testAccCheckMSOSchemaTemplateBDSubnetAttributes(shared bool, ss *TemplateBDSubnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if shared != ss.shared {
			return fmt.Errorf("Bad Template BD Subnet shared value %v", ss.shared)
		}

		return nil
	}
}

type TemplateBDSubnet struct {
	ip                 string
	scope              string
	shared             bool
	no_default_gateway bool
	querier            bool
}
