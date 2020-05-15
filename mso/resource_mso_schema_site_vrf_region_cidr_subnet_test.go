package mso

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccMSOSchemaSiteVrfRegionCidrSubnet_Basic(t *testing.T) {
	var ss SchemaSiteVrfRegionCidrSubnet
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteVrfRegionCidrSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSchemaSiteVrfRegionCidrSubnetConfig_basic("gateway"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteVrfRegionCidrSubnetExists("mso_schema_site_vrf_region_cidr_subnet.sub1", &ss),
					testAccCheckMSOSchemaSiteVrfRegionCidrSubnetAttributes("gateway", &ss),
				),
			},
		},
	})
}

func TestAccMSOSchemaSiteVrfRegionCidrSubnet_Update(t *testing.T) {
	var ss SchemaSiteVrfRegionCidrSubnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteVrfRegionCidrSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSchemaSiteVrfRegionCidrSubnetConfig_basic("gateway"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteVrfRegionCidrSubnetExists("mso_schema_site_vrf_region_cidr_subnet.sub1", &ss),
					testAccCheckMSOSchemaSiteVrfRegionCidrSubnetAttributes("gateway", &ss),
				),
			},
			{
				Config: testAccCheckMSOSchemaSiteVrfRegionCidrSubnetConfig_basic("gateways"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteVrfRegionCidrSubnetExists("mso_schema_site_vrf_region_cidr_subnet.sub1", &ss),
					testAccCheckMSOSchemaSiteVrfRegionCidrSubnetAttributes("gateways", &ss),
				),
			},
		},
	})
}

func testAccCheckMSOSchemaSiteVrfRegionCidrSubnetConfig_basic(usage string) string {
	return fmt.Sprintf(`
	resource "mso_schema_site_vrf_region_cidr_subnet" "sub1" {
		schema_id = "5d5dbf3f2e0000580553ccce"
		template_name = "Template1"
		site_id = "5ce2de773700006a008a2678"
		vrf_name = "Campus"
		region_name = "westus"
		cidr_ip = "1.1.1.1/24"
		ip = "203.168.240.1/24"
		zone = "West"
		usage = "%s"
	  
	  }
	`, usage)
}

func testAccCheckMSOSchemaSiteVrfRegionCidrSubnetExists(siteVrfRegionName string, ss *SchemaSiteVrfRegionCidrSubnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs, err := s.RootModule().Resources[siteVrfRegionName]

		if !err {
			return fmt.Errorf("Vrf Region %s not found", siteVrfRegionName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Region id was set")
		}

		cont, errs := client.GetViaURL("api/v1/schemas/5d5dbf3f2e0000580553ccce")
		if errs != nil {
			return errs
		}
		count, ers := cont.ArrayCount("sites")
		if ers != nil {
			return fmt.Errorf("No Sites found")
		}

		tp := SchemaSiteVrfRegionCidrSubnet{}
		found := false

		for i := 0; i < count; i++ {
			tempCont, err := cont.ArrayElement(i, "sites")
			if err != nil {
				return err
			}
			apiSite := models.StripQuotes(tempCont.S("siteId").String())

			if apiSite == "5ce2de773700006a008a2678" {
				vrfCount, err := tempCont.ArrayCount("vrfs")
				if err != nil {
					return fmt.Errorf("Unable to get Vrf list")
				}
				for j := 0; j < vrfCount; j++ {
					vrfCont, err := tempCont.ArrayElement(j, "vrfs")
					if err != nil {
						return err
					}
					apiVrfRef := models.StripQuotes(vrfCont.S("vrfRef").String())
					split := strings.Split(apiVrfRef, "/")
					apiVrf := split[6]
					if apiVrf == "Campus" {

						regionCount, err := vrfCont.ArrayCount("regions")
						if err != nil {
							return fmt.Errorf("Unable to get Regions list")
						}
						for k := 0; k < regionCount; k++ {
							regionCont, err := vrfCont.ArrayElement(k, "regions")
							if err != nil {
								return err
							}
							apiRegion := models.StripQuotes(regionCont.S("name").String())
							if apiRegion == "westus" {
								cidrCount, err := regionCont.ArrayCount("cidrs")
								if err != nil {
									return fmt.Errorf("Unable to get Cidr list")
								}
								for l := 0; l < cidrCount; l++ {
									cidrCont, err := regionCont.ArrayElement(l, "cidrs")
									if err != nil {
										return err
									}
									apiCidr := models.StripQuotes(cidrCont.S("ip").String())
									log.Println("Current Cidr Ip", apiCidr)
									if apiCidr == "1.1.1.1/24" {

										subnetCount, err := cidrCont.ArrayCount("subnets")
										if err != nil {
											return fmt.Errorf("Unable to get Subnet list")
										}
										for m := 0; m < subnetCount; m++ {
											subnetCont, err := cidrCont.ArrayElement(m, "subnets")

											if err != nil {
												return err
											}
											apiIp := models.StripQuotes(subnetCont.S("ip").String())

											if apiIp == "203.168.240.1/24" {
												tp.siteId = apiSite
												tp.regionName = apiRegion
												tp.vrfName = apiVrf
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
			}
		}

		if !found {
			return fmt.Errorf("Vrf Region Cidr Subnet not found from API")
		}
		tp1 := &tp
		*ss = *tp1
		return nil
	}
}

func testAccCheckMSOSchemaSiteVrfRegionCidrSubnetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_site_vrf_region_cidr_subnet" {
			cont, err := client.GetViaURL("api/v1/schemas/5d5dbf3f2e0000580553ccce")
			if err != nil {
				return err
			} else {
				count, err := cont.ArrayCount("sites")
				if err != nil {
					return fmt.Errorf("No Sites found")
				}
				for i := 0; i < count; i++ {
					tempCont, err := cont.ArrayElement(i, "sites")
					if err != nil {
						return fmt.Errorf("No Site exists")
					}
					apiSite := models.StripQuotes(tempCont.S("siteId").String())

					if apiSite == "5ce2de773700006a008a2678" {
						vrfCount, err := tempCont.ArrayCount("vrfs")
						if err != nil {
							return fmt.Errorf("Unable to get Vrf list")
						}
						for j := 0; j < vrfCount; j++ {
							vrfCont, err := tempCont.ArrayElement(j, "vrfs")
							if err != nil {
								return err
							}
							apiVrfRef := models.StripQuotes(vrfCont.S("vrfRef").String())
							split := strings.Split(apiVrfRef, "/")
							apiVrf := split[6]
							if apiVrf == "Campus" {

								regionCount, err := vrfCont.ArrayCount("regions")
								if err != nil {
									return fmt.Errorf("Unable to get Regions list")
								}
								for k := 0; k < regionCount; k++ {
									regionCont, err := vrfCont.ArrayElement(k, "regions")
									if err != nil {
										return err
									}
									apiRegion := models.StripQuotes(regionCont.S("name").String())
									if apiRegion == "westus" {
										cidrCount, err := regionCont.ArrayCount("cidrs")
										if err != nil {
											return fmt.Errorf("Unable to get Cidr list")
										}
										for l := 0; l < cidrCount; l++ {
											cidrCont, err := regionCont.ArrayElement(l, "cidrs")
											if err != nil {
												return err
											}
											apiCidr := models.StripQuotes(cidrCont.S("ip").String())
											log.Println("Current Cidr Ip", apiCidr)
											if apiCidr == "1.1.1.1/24" {

												subnetCount, err := cidrCont.ArrayCount("subnets")
												if err != nil {
													return fmt.Errorf("Unable to get Subnet list")
												}
												for m := 0; m < subnetCount; m++ {
													subnetCont, err := cidrCont.ArrayElement(m, "subnets")

													if err != nil {
														return err
													}
													apiIp := models.StripQuotes(subnetCont.S("ip").String())

													if apiIp == "203.168.240.1/24" {
														return fmt.Errorf("Vrf Region Cidr Subnet exists")

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
		}
	}
	return nil

}
func testAccCheckMSOSchemaSiteVrfRegionCidrSubnetAttributes(usage string, ss *SchemaSiteVrfRegionCidrSubnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "Campus" != ss.vrfName {
			return fmt.Errorf("Bad Vrf name %s", ss.vrfName)
		}
		return nil
	}
}

type SchemaSiteVrfRegionCidrSubnet struct {
	siteId     string
	vrfName    string
	regionName string
}
