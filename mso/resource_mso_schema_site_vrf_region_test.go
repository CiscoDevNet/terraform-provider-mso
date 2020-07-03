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

func TestAccMSOSchemaSiteVrfRegion_Basic(t *testing.T) {
	var ss SchemaSiteVrfRegion
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteVrfRegionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSchemaSiteVrfRegionConfig_basic(false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteVrfRegionExists("mso_schema_site_vrf_region.vrfRegion", &ss),
					testAccCheckMSOSchemaSiteVrfRegionAttributes(false, &ss),
				),
			},
		},
	})
}

func TestAccMSOSchemaSiteVrfRegion_Update(t *testing.T) {
	var ss SchemaSiteVrfRegion
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteVrfRegionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSchemaSiteVrfRegionConfig_basic(false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteVrfRegionExists("mso_schema_site_vrf_region.vrfRegion", &ss),
					testAccCheckMSOSchemaSiteVrfRegionAttributes(false, &ss),
				),
			},
			{
				Config: testAccCheckMSOSchemaSiteVrfRegionConfig_basic(true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteVrfRegionExists("mso_schema_site_vrf_region.vrfRegion", &ss),
					testAccCheckMSOSchemaSiteVrfRegionAttributes(true, &ss),
				),
			},
		},
	})
}

func testAccCheckMSOSchemaSiteVrfRegionConfig_basic(vpn bool) string {
	return fmt.Sprintf(`
	resource "mso_schema_site_vrf_region" "vrfRegion" {
		schema_id = "5efd6ea60f00005b0ebbd643"
		template_name = "Template1"
		site_id = "5efeb3c4190000cc12d05376"
		vrf_name = "Myvrf"
		region_name = "us-east-1"
		vpn_gateway = %v
		cidr {
		  cidr_ip = "2.2.2.2/10"
		  primary = true
		  subnet {
			ip = "1.20.30.4"
			zone = "us-east-1b"
			usage = "sdfkhsdkf"
		  }
		}
	  }
	`, vpn)
}

func testAccCheckMSOSchemaSiteVrfRegionExists(siteVrfRegionName string, ss *SchemaSiteVrfRegion) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs, err := s.RootModule().Resources[siteVrfRegionName]

		if !err {
			return fmt.Errorf("Vrf Region %s not found", siteVrfRegionName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Region id was set")
		}

		cont, errs := client.GetViaURL("api/v1/schemas/5efd6ea60f00005b0ebbd643")
		if errs != nil {
			return errs
		}
		count, ers := cont.ArrayCount("sites")
		if ers != nil {
			return fmt.Errorf("No Sites found")
		}

		tp := SchemaSiteVrfRegion{}
		found := false

		for i := 0; i < count; i++ {
			tempCont, err := cont.ArrayElement(i, "sites")
			if err != nil {
				return err
			}
			apiSite := models.StripQuotes(tempCont.S("siteId").String())

			if apiSite == "5efeb3c4190000cc12d05376" {
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
					if apiVrf == "Myvrf" {

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
							if apiRegion == "us-east-1" {
								tp.siteId = apiSite
								tp.vrfName = apiVrf
								tp.regionName = apiRegion
								tp.VPN = regionCont.S("isVpnGatewayRouter").Data().(bool)
								found = true
								break
							}
						}
					}
				}
			}
		}

		if !found {
			return fmt.Errorf("Vrf Region not found from API")
		}
		tp1 := &tp
		*ss = *tp1
		return nil
	}
}

func testAccCheckMSOSchemaSiteVrfRegionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_site_vrf_region" {
			cont, err := client.GetViaURL("api/v1/schemas/5efd6ea60f00005b0ebbd643")
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

					if apiSite == "5efeb3c4190000cc12d05376" {
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
							if apiVrf == "Myvrf" {

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
									if apiRegion == "us-east-1" {
										return fmt.Errorf("The Vrf Region still exists")
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
func testAccCheckMSOSchemaSiteVrfRegionAttributes(vpn bool, ss *SchemaSiteVrfRegion) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "Myvrf" != ss.vrfName {
			return fmt.Errorf("Bad Vrf name %s", ss.vrfName)
		}
		if vpn != ss.VPN {
			return fmt.Errorf("Bad VPN Gateway name %v", ss.VPN)
		}
		return nil
	}
}

type SchemaSiteVrfRegion struct {
	siteId     string
	vrfName    string
	regionName string
	VPN        bool
}
