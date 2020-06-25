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

func TestAccMSOSchemaSiteAnpEpgStaticleaf_Basic(t *testing.T) {
	var ss SchemaSiteAnpEpgStaticleaf
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpEpgStaticleafDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSchemaSiteAnpEpgStaticleafConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteAnpEpgStaticleafExists("mso_schema_site_anp_epg_static_leaf.staticleaf1", &ss),
					testAccCheckMSOSchemaSiteAnpEpgStaticleafAttributes(&ss),
				),
			},
		},
	})
}

func testAccCheckMSOSchemaSiteAnpEpgStaticleafConfig_basic() string {
	return fmt.Sprintf(`
	resource "mso_schema_site_anp_epg_static_leaf" "staticleaf1" {
		schema_id = "5c4d9fca270000a101f8094a"
		template_name = "Template1"
		site_id = "5c7c95b25100008f01c1ee3c"
		anp_name = "ANP"
		epg_name = "Web"
		path= "topology/pod-1/paths-103/pathep-[eth1/111]"
		port_encap_vlan = 100
	  }
	`)
}

func testAccCheckMSOSchemaSiteAnpEpgStaticleafExists(siteAnpEpgName string, ss *SchemaSiteAnpEpgStaticleaf) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs, err := s.RootModule().Resources[siteAnpEpgName]

		if !err {
			return fmt.Errorf("Anp Epg StaticLeaf %s not found", siteAnpEpgName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No StaticLeaf id was set")
		}

		cont, errs := client.GetViaURL("api/v1/schemas/5c4d9fca270000a101f8094a")
		if errs != nil {
			return errs
		}
		count, ers := cont.ArrayCount("sites")
		if ers != nil {
			return fmt.Errorf("No Sites found")
		}

		tp := SchemaSiteAnpEpgStaticleaf{}
		found := false

		for i := 0; i < count; i++ {
			tempCont, err := cont.ArrayElement(i, "sites")
			if err != nil {
				return err
			}
			apiSite := models.StripQuotes(tempCont.S("siteId").String())

			if apiSite == "5c7c95b25100008f01c1ee3c" {
				anpCount, err := tempCont.ArrayCount("anps")
				if err != nil {
					return fmt.Errorf("Unable to get Anp list")
				}
				for j := 0; j < anpCount; j++ {
					anpCont, err := tempCont.ArrayElement(j, "anps")
					if err != nil {
						return err
					}
					apiAnpRef := models.StripQuotes(anpCont.S("anpRef").String())
					split := strings.Split(apiAnpRef, "/")
					apiAnp := split[6]
					if apiAnp == "ANP" {
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
							split := strings.Split(apiEpgRef, "/")
							apiEPG := split[8]
							if apiEPG == "Web" {
								staticLeafCount, err := epgCont.ArrayCount("staticLeafs")
								if err != nil {
									return fmt.Errorf("Unable to get Static Leaf list")
								}
								for s := 0; s < staticLeafCount; s++ {
									staticLeafCont, err := epgCont.ArrayElement(s, "staticLeafs")
									if err != nil {
										return err
									}
									apiPath := models.StripQuotes(staticLeafCont.S("path").String())
									if apiPath == "topology/pod-1/paths-103/pathep-[eth1/111]" {
										tp.epgName = apiEPG
										tp.schemaId = split[2]
										tp.templateName = split[4]
										tp.path = apiPath
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
			return fmt.Errorf("Anp Epg StaticLeaf not found from API")
		}
		tp1 := &tp
		*ss = *tp1
		return nil
	}
}

func testAccCheckMSOSchemaSiteAnpEpgStaticleafDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_site_anp_epg_static_leaf" {
			cont, err := client.GetViaURL("api/v1/schemas/5c4d9fca270000a101f8094a")
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

					if apiSite == "5c7c95b25100008f01c1ee3c" {
						anpCount, err := tempCont.ArrayCount("anps")
						if err != nil {
							return fmt.Errorf("Unable to get Anp list")
						}
						for j := 0; j < anpCount; j++ {
							anpCont, err := tempCont.ArrayElement(j, "anps")
							if err != nil {
								return err
							}
							apiAnpRef := models.StripQuotes(anpCont.S("anpRef").String())
							split := strings.Split(apiAnpRef, "/")
							apiAnp := split[6]
							if apiAnp == "ANP" {
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
									split := strings.Split(apiEpgRef, "/")
									apiEPG := split[8]
									if apiEPG == "Web" {
										staticLeafCount, err := epgCont.ArrayCount("staticLeafs")
										if err != nil {
											return fmt.Errorf("Unable to get Static Leaf list")
										}
										for s := 0; s < staticLeafCount; s++ {
											staticLeafCont, err := epgCont.ArrayElement(s, "staticLeafs")
											if err != nil {
												return err
											}
											apiPath := models.StripQuotes(staticLeafCont.S("path").String())
											if apiPath == "topology/pod-1/paths-103/pathep-[eth1/111]" {
												return fmt.Errorf("StaticLeaf still exists")
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
func testAccCheckMSOSchemaSiteAnpEpgStaticleafAttributes(ss *SchemaSiteAnpEpgStaticleaf) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "Template1" != ss.templateName {
			return fmt.Errorf("Bad Template name %s", ss.templateName)
		}
		return nil
	}
}

type SchemaSiteAnpEpgStaticleaf struct {
	schemaId     string
	templateName string
	epgName      string
	path         string
}
