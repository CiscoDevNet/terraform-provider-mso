package mso

import (
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceMSOSchemaSiteAnpEpgDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOSchemaSiteAnpEpgDomainRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"site_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"anp_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"epg_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"domain_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"domain_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"vmmDomain",
					"l3ExtDomain",
					"l2ExtDomain",
					"physicalDomain",
					"fibreChannelDomain",
				}, false),
			},
			"vmm_domain_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"VMware",
					"Microsoft",
					"Redhat",
				}, false),
			},
			"domain_dn": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringLenBetween(1, 1000),
				ConflictsWith: []string{"domain_name", "vmm_domain_type", "domain_type"},
			},
			"dn": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringLenBetween(1, 1000),
				ConflictsWith: []string{"domain_name", "domain_dn"},
				Deprecated:    "use domain_dn alone or domain_name in association with domain_type and vmm_domain_type when it is applicable.",
			},
			"deployment_immediacy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resolution_immediacy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"micro_seg_vlan_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"micro_seg_vlan": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"port_encap_vlan_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"port_encap_vlan": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"vlan_encap_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"allow_micro_segmentation": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"switching_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"switch_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enhanced_lag_policy_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enhanced_lag_policy_dn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceMSOSchemaSiteAnpEpgDomainRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return fmt.Errorf("No Sites found")
	}
	stateSite := d.Get("site_id").(string)
	found := false
	stateTemplate := d.Get("template_name").(string)
	stateAnp := d.Get("anp_name").(string)
	stateEpg := d.Get("epg_name").(string)
	domainNameDnOld := d.Get("dn").(string)
	domainNameNew := d.Get("domain_name").(string)
	domainType := d.Get("domain_type").(string)
	DN := d.Get("domain_dn").(string)

	var domainName, stateDomain string

	if domainNameNew == "" && DN == "" && domainNameDnOld == "" {
		return fmt.Errorf("domain_dn or domain_name in association with domain_type and vmm_domain_type when it is applicable are required.")
	}

	if domainNameNew != "" && DN == "" && domainNameDnOld == "" {
		domainName = domainNameNew
	} else if domainNameNew == "" && DN == "" && domainNameDnOld != "" {
		domainName = domainNameDnOld
	} else if domainNameNew == "" && DN != "" && domainNameDnOld == "" {
		stateDomain = DN
	}
	if stateDomain == "" {
		if domainType == "vmmDomain" {
			vmmDomainType := d.Get("vmm_domain_type").(string)
			stateDomain = fmt.Sprintf("uni/vmmp-%s/dom-%s", vmmDomainType, domainName)

		} else if domainType == "l3ExtDomain" {
			stateDomain = fmt.Sprintf("uni/l3dom-%s", domainName)

		} else if domainType == "l2ExtDomain" {
			stateDomain = fmt.Sprintf("uni/l2dom-%s", domainName)

		} else if domainType == "physicalDomain" {
			stateDomain = fmt.Sprintf("uni/phys-%s", domainName)

		} else if domainType == "fibreChannelDomain" {
			stateDomain = fmt.Sprintf("uni/fc-%s", domainName)

		}
	}

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())
		apiTemplate := models.StripQuotes(tempCont.S("templateName").String())

		if apiSite == stateSite && apiTemplate == stateTemplate {
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
				if apiAnp == stateAnp {
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
						if apiEPG == stateEpg {
							d.Set("schema_id", split[2])
							d.Set("template_name", split[4])
							d.Set("anp_name", split[6])
							d.Set("epg_name", apiEPG)

							domainCount, err := epgCont.ArrayCount("domainAssociations")
							if err != nil {
								return fmt.Errorf("Unable to get Domain Associations list")
							}
							for l := 0; l < domainCount; l++ {
								domainCont, err := epgCont.ArrayElement(l, "domainAssociations")
								if err != nil {
									return err
								}
								apiDomain := models.StripQuotes(domainCont.S("dn").String())

								if apiDomain == stateDomain {
									d.SetId(apiDomain)
									d.Set("site_id", apiSite)
									d.Set("domain_type", models.StripQuotes(domainCont.S("domainType").String()))
									d.Set("domain_dn", apiDomain)

									d.Set("deploy_immediacy", models.StripQuotes(domainCont.S("deployImmediacy").String()))
									d.Set("resolution_immediacy", models.StripQuotes(domainCont.S("resolutionImmediacy").String()))

									if domainCont.Exists("switchingMode") {
										d.Set("switching_mode", models.StripQuotes(domainCont.S("switchingMode").String()))
									}

									if domainCont.Exists("switchType") {
										d.Set("switch_type", models.StripQuotes(domainCont.S("switchType").String()))
									}

									if domainCont.Exists("vlanEncapMode") {
										d.Set("vlan_encap_mode", models.StripQuotes(domainCont.S("vlanEncapMode").String()))
									}

									if domainCont.Exists("allowMicroSegmentation") {
										d.Set("allow_micro_segmentation", domainCont.S("allowMicroSegmentation").Data().(bool))
									}

									if domainCont.Exists("portEncapVlan") {
										d.Set("port_encap_vlan", domainCont.S("portEncapVlan", "vlan").Data().(float64))
										d.Set("port_encap_vlan_type", models.StripQuotes(domainCont.S("portEncapVlan", "vlanType").String()))
									}

									if domainCont.Exists("microSegVlan") {
										d.Set("micro_seg_vlan", domainCont.S("microSegVlan", "vlan").Data().(float64))
										d.Set("micro_seg_vlan_type", models.StripQuotes(domainCont.S("microSegVlan", "vlanType").String()))
									}

									if domainCont.Exists("epgLagPol") {
										if domainCont.Exists("epgLagPol", "enhancedLagPol") {
											d.Set("enhanced_lag_policy_name", models.StripQuotes(domainCont.S("epgLagPol", "enhancedLagPol", "name").String()))
											d.Set("enhanced_lag_policy_dn", models.StripQuotes(domainCont.S("epgLagPol", "enhancedLagPol", "dn").String()))
										}
									}
									found = true
									break
								}
							}
						}
						if found {
							break
						}
					}
				}
				if found {
					break
				}
			}
		}
		if found {
			break
		}
	}

	if !found {
		return fmt.Errorf("Unable to find the Site Anp Epg Domain %s", stateDomain)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
