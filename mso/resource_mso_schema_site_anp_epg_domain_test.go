package mso

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccMSOSchemaSiteAnpEpgDomain_Basic(t *testing.T) {
	var ss SiteAnpEpgDomain
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpEpgDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSchemaSiteAnpEpgDomainConfig_basic("immediate"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteAnpEpgDomainExists("mso_schema_site_anp_epg_domain.site_anp_epg_domain", &ss),
					testAccCheckMSOSchemaSiteAnpEpgDomainAttributes("immediate", &ss),
				),
			},
		},
	})
}

func TestAccMSOSchemaSiteAnpEpgDomain_Update(t *testing.T) {
	var ss SiteAnpEpgDomain

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpEpgDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOSchemaSiteAnpEpgDomainConfig_basic("immediate"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteAnpEpgDomainExists("mso_schema_site_anp_epg_domain.site_anp_epg_domain", &ss),
					testAccCheckMSOSchemaSiteAnpEpgDomainAttributes("immediate", &ss),
				),
			},
			{
				Config: testAccCheckMSOSchemaSiteAnpEpgDomainConfig_basic("lazy"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaSiteAnpEpgDomainExists("mso_schema_site_anp_epg_domain.site_anp_epg_domain", &ss),
					testAccCheckMSOSchemaSiteAnpEpgDomainAttributes("lazy", &ss),
				),
			},
		},
	})
}

func testAccCheckMSOSchemaSiteAnpEpgDomainConfig_basic(immediacy string) string {
	return fmt.Sprintf(`
	resource "mso_schema_site_anp_epg_domain" "site_anp_epg_domain" {
		schema_id = "5c4d9fca270000a101f8094a"
		template_name = "Template1"
		site_id = "5c7c95b25100008f01c1ee3c"
		anp_name = "ANP"
		epg_name = "Web"
		domain_type = "vmmDomain"
		dn = "VMware-Vmm"
		deploy_immediacy = "%v"
		resolution_immediacy = "%v"
		vlan_encap_mode = "static"
		allow_micro_segmentation = true
		switching_mode = "native"
		switch_type = "default"
		micro_seg_vlan_type = "vlan"
		micro_seg_vlan = 46
		port_encap_vlan_type = "vlan"
		port_encap_vlan = 45
		enhanced_lagpolicy_name = "name"
		enhanced_lagpolicy_dn = "dn"
	  
	  }`, immediacy, immediacy)
}

func testAccCheckMSOSchemaSiteAnpEpgDomainExists(anpEpgDomainName string, ss *SiteAnpEpgDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs, err := s.RootModule().Resources[anpEpgDomainName]

		if !err {
			return fmt.Errorf("Site Anp Epg Domain %s not found", anpEpgDomainName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Domain Id was set")
		}

		cont, errs := client.GetViaURL("api/v1/schemas/5c4d9fca270000a101f8094a")
		if errs != nil {
			return errs
		}
		count, ers := cont.ArrayCount("sites")
		if ers != nil {
			return fmt.Errorf("No Sites found")
		}

		tp := SiteAnpEpgDomain{}
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
								domainCount, err := epgCont.ArrayCount("domainAssociations")
								if err != nil {
									return fmt.Errorf("Unable to get Domain Associations list")
								}
								for l := 0; l < domainCount; l++ {
									domainCont, err := epgCont.ArrayElement(l, "domainAssociations")
									if err != nil {
										return err
									}
									tempVar := strings.Split(models.StripQuotes(domainCont.S("dn").String()), "/")
									apiDomain := strings.SplitN(tempVar[2], "-", 2)
									if apiDomain[1] == "VMware-Vmm" {
										tp.dn = apiDomain[1]
										tp.siteId = apiSite
										tp.epgName = apiEPG
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
			return fmt.Errorf("Anp Epg Domain not found from API")
		}
		tp1 := &tp
		*ss = *tp1
		return nil
	}
}

func testAccCheckMSOSchemaSiteAnpEpgDomainDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_site_anp_epg_domain" {
			cont, err := client.GetViaURL("api/v1/schemas/5c6c16d7270000c710f8094d")
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
										domainCount, err := epgCont.ArrayCount("domainAssociations")
										if err != nil {
											return fmt.Errorf("Unable to get Domain Associations list")
										}
										for l := 0; l < domainCount; l++ {
											domainCont, err := epgCont.ArrayElement(l, "domainAssociations")
											if err != nil {
												return err
											}
											tempVar := strings.Split(models.StripQuotes(domainCont.S("dn").String()), "/")
											apiDomain := strings.SplitN(tempVar[2], "-", 2)
											if apiDomain[1] == "VMware-Vmm" {
												return fmt.Errorf("The Anp Epg Domain still exists")
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

func testAccCheckMSOSchemaSiteAnpEpgDomainAttributes(immediacy string, ss *SiteAnpEpgDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "5c7c95b25100008f01c1ee3c" != ss.siteId {
			return fmt.Errorf("Bad siteId %s", ss.siteId)
		}
		return nil
	}
}

type SiteAnpEpgDomain struct {
	dn      string
	siteId  string
	epgName string
}
