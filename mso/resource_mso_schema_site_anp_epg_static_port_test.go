package mso

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccMSOSchemaSiteAnpEpgStaticPort_Basic(t *testing.T) {
	var ss SchemaSiteAnpEpgStaticPort
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpEpgStaticPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSiteAnpEpgStaticPortConfig_basic("untagged"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteAnpEpgStaticPortExists("mso_schema_site_anp_epg_static_port.static_port", &ss),
					testAccCheckMSOSchemaSiteAnpEpgStaticPortAttributes("untagged", &ss),
				),
			},
		},
	})
}

func TestAccMSOSchemaSiteAnpEpgStaticPort_Update(t *testing.T) {
	var ss SchemaSiteAnpEpgStaticPort

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpEpgStaticPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSiteAnpEpgStaticPortConfig_basic("untagged"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteAnpEpgStaticPortExists("mso_schema_site_anp_epg_static_port.static_port", &ss),
					testAccCheckMSOSchemaSiteAnpEpgStaticPortAttributes("untagged", &ss),
				),
			},
			{
				Config: testAccCheckMSOSiteAnpEpgStaticPortConfig_basic("regular"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteAnpEpgStaticPortExists("mso_schema_site_anp_epg_static_port.static_port", &ss),
					testAccCheckMSOSchemaSiteAnpEpgStaticPortAttributes("regular", &ss),
				),
			},
		},
	})
}

func testAccCheckMSOSiteAnpEpgStaticPortConfig_basic(mode string) string {
	return fmt.Sprintf(`
   resource "mso_schema_site_anp_epg_static_port" "static_port" {
   schema_id = "5c4d5bb72700000401f80948"
   site_id = "5c7c95b25100008f01c1ee3c"
   template_name = "Template1"
   anp_name = "ANP"
   epg_name = "DB"
   path_type = "port"
   deployment_immediacy = "lazy"
   pod = "pod-9"
   leaf = "112"
   path = "eth1/10"
   vlan = 50
   mode = "%s"

  
}

`, mode)
}

func testAccCheckMSOSchemaSiteAnpEpgStaticPortExists(portName string, ss *SchemaSiteAnpEpgStaticPort) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs1, err1 := s.RootModule().Resources[portName]

		if !err1 {
			return fmt.Errorf("Entry %s not found", portName)
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
		tp := SchemaSiteAnpEpgStaticPort{}
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
								portCount, err := epgCont.ArrayCount("staticPorts")
								if err != nil {
									return fmt.Errorf("Unable to get Static Port list")
								}
								for l := 0; l < portCount; l++ {
									portCont, err := epgCont.ArrayElement(l, "staticPorts")
									if err != nil {
										return err
									}
									portpath := fmt.Sprintf("topology/pod-9/paths-112/pathep-[eth1/10]")
									apiportpath := models.StripQuotes(portCont.S("path").String())
									if portpath == apiportpath {
										if portCont.Exists("portEncapVlan") {
											tempvar, _ := strconv.Atoi(fmt.Sprintf("%v", portCont.S("portEncapVlan")))
											tp.portencapvlan = tempvar
										}
										tp.deploymentimmediacy = models.StripQuotes(portCont.S("deploymentImmediacy").String())
										tp.mode = models.StripQuotes(portCont.S("mode").String())
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
			return fmt.Errorf("Static Port Entry not found from API")
		}

		tp1 := &tp

		*ss = *tp1
		return nil
	}
}

func testAccCheckMSOSchemaSiteAnpEpgStaticPortDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_site_anp_epg_static_port" {
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
										portCount, err := epgCont.ArrayCount("staticPorts")
										if err != nil {
											return fmt.Errorf("Unable to get Static Port list")
										}
										for l := 0; l < portCount; l++ {
											portCont, err := epgCont.ArrayElement(l, "staticPorts")
											if err != nil {
												return err
											}
											portpath := fmt.Sprintf("topology/pod-9/paths-112/pathep-[eth1/10]")
											apiportpath := models.StripQuotes(portCont.S("path").String())
											if portpath == apiportpath {
												return fmt.Errorf("The static port entry still exists")
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

func testAccCheckMSOSchemaSiteAnpEpgStaticPortAttributes(ethertype string, ss *SchemaSiteAnpEpgStaticPort) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if 50 != ss.portencapvlan {
			return fmt.Errorf("Bad Static Port Encap Vlan value %v", ss.portencapvlan)
		}

		if "lazy" != ss.deploymentimmediacy {
			return fmt.Errorf("Bad Static Port Deployment Immediacy value %s", ss.deploymentimmediacy)
		}
		return nil
	}
}

type SchemaSiteAnpEpgStaticPort struct {
	portencapvlan       int
	deploymentimmediacy string
	mode                string
}
