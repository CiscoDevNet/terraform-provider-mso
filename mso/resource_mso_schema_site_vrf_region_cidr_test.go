package mso

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccMSOSchemaSiteVrfRegionCidr_Basic(t *testing.T) {
	var ss SiteVrfRegionCidr
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteVrfRegionCidrDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSchemaSiteVrfRegionCidrConfig_basic(true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteVrfRegionCidrExists("mso_schema_site_vrf_region_cidr.vrfRegionCidr", &ss),
					testAccCheckMSOSchemaSiteVrfRegionCidrAttributes(true, &ss),
				),
			},
		},
	})
}

func TestAccMSOSchemaSiteVrfRegionCidr_Update(t *testing.T) {
	var ss SiteVrfRegionCidr

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteVrfRegionCidrDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSchemaSiteVrfRegionCidrConfig_basic(true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteVrfRegionCidrExists("mso_schema_site_vrf_region_cidr.vrfRegionCidr", &ss),
					testAccCheckMSOSchemaSiteVrfRegionCidrAttributes(true, &ss),
				),
			},
			{
				Config: testAccCheckMSOSchemaSiteVrfRegionCidrConfig_basic(false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteVrfRegionCidrExists("mso_schema_site_vrf_region_cidr.vrfRegionCidr", &ss),
					testAccCheckMSOSchemaSiteVrfRegionCidrAttributes(false, &ss),
				),
			},
		},
	})
}

func testAccCheckMSOSchemaSiteVrfRegionCidrConfig_basic(primary bool) string {
	return fmt.Sprintf(`
	resource "mso_schema_site_vrf_region_cidr" "vrfRegionCidr" {
		schema_id = "5d5dbf3f2e0000580553ccce"
		template_name = "Template1"
		site_id = "5ce2de773700006a008a2678"
		vrf_name = "Campus"
		region_name = "region1"
		ip = "3.3.2.2/2"
		primary = %v
	  }`, primary)
}

func testAccCheckMSOSchemaSiteVrfRegionCidrExists(VrfRegionCidrName string, ss *SiteVrfRegionCidr) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs, err := s.RootModule().Resources[VrfRegionCidrName]

		if !err {
			return fmt.Errorf("Site Vrf Region Cidr %s not found", VrfRegionCidrName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cidr Id was set")
		}

		cont, errs := client.GetViaURL("api/v1/schemas/5d5dbf3f2e0000580553ccce")
		if errs != nil {
			return errs
		}
		count, ers := cont.ArrayCount("sites")
		if ers != nil {
			return fmt.Errorf("No Sites found")
		}

		tp := SiteVrfRegionCidr{}
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
							if apiRegion == "region1" {
								cidrCount, err := regionCont.ArrayCount("cidrs")
								if err != nil {
									return fmt.Errorf("Unable to get Cidr list")
								}
								for l := 0; l < cidrCount; l++ {
									cidrCont, err := regionCont.ArrayElement(l, "cidrs")
									if err != nil {
										return err
									}
									apiIp := models.StripQuotes(cidrCont.S("ip").String())
									if apiIp == "3.3.2.2/2" {
										tp.ip = apiIp
										tp.primary = cidrCont.S("primary").Data().(bool)
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

		if !found {
			return fmt.Errorf("Vrf Region Cidr not found from API")
		}
		tp1 := &tp
		*ss = *tp1
		return nil
	}
}

func testAccCheckMSOSchemaSiteVrfRegionCidrDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_site_vrf_region_cidr" {
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
									if apiRegion == "region1" {
										cidrCount, err := regionCont.ArrayCount("cidrs")
										if err != nil {
											return fmt.Errorf("Unable to get Cidr list")
										}
										for l := 0; l < cidrCount; l++ {
											cidrCont, err := regionCont.ArrayElement(l, "cidrs")
											if err != nil {
												return err
											}
											apiIp := models.StripQuotes(cidrCont.S("ip").String())
											if apiIp == "3.3.2.2/2" {
												return fmt.Errorf("Vrf Region Cidr still Exist.")
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

func testAccCheckMSOSchemaSiteVrfRegionCidrAttributes(primary bool, ss *SiteVrfRegionCidr) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if primary != ss.primary {
			return fmt.Errorf("Bad primary")
		}
		if "Campus" != ss.vrfName {
			return fmt.Errorf("Bad Vrf Name %s", ss.vrfName)
		}
		return nil
	}
}

type SiteVrfRegionCidr struct {
	ip         string
	vrfName    string
	regionName string
	primary    bool
}
