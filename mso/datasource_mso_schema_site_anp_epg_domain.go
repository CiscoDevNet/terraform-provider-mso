package mso

import (
	"fmt"
	"log"

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
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"anp_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"epg_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"domain_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"domain_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"VMware",
					"Microsoft",
					"Redhat",
				}, false),
			},
			"domain_dn": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringLenBetween(1, 1000),
				ConflictsWith: []string{"domain_name", "vmm_domain_type", "domain_type"},
			},
			"dn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"deployment_immediacy": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"resolution_immediacy": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"micro_seg_vlan_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"micro_seg_vlan": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"port_encap_vlan_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"port_encap_vlan": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vlan_encap_mode": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"allow_micro_segmentation": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"switching_mode": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"switch_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"enhanced_lag_policy_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"enhanced_lag_policy_dn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}
}

func dataSourceMSOSchemaSiteAnpEpgDomainRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	anp := d.Get("anp_name").(string)
	epg := d.Get("epg_name").(string)
	domainNameDnOld := d.Get("dn").(string)
	domainNameNew := d.Get("domain_name").(string)
	domainType := d.Get("domain_type").(string)
	dn := d.Get("domain_dn").(string)

	var domainName, stateDomain string

	if domainNameNew == "" && dn == "" && domainNameDnOld == "" {
		return fmt.Errorf("domain_dn or domain_name in association with domain_type and vmm_domain_type when it is applicable are required.")
	}

	if domainNameNew != "" && dn == "" && domainNameDnOld == "" {
		domainName = domainNameNew
	} else if domainNameNew == "" && dn == "" && domainNameDnOld != "" {
		domainName = domainNameDnOld
	} else if domainNameNew == "" && dn != "" && domainNameDnOld == "" {
		stateDomain = dn
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

	siteCont, err := getSiteFromSiteIdAndTemplate(schemaId, siteId, templateName, msoClient)
	if err != nil {
		return err
	} else {
		d.Set("schema_id", schemaId)
		d.Set("site_id", siteId)
		d.Set("template_name", templateName)
	}

	anpCont, err := getSiteAnp(anp, siteCont)
	if err != nil {
		return err
	} else {
		d.Set("anp_name", anp)
	}

	epgCont, err := getSiteEpg(epg, anpCont)
	if err != nil {
		return err
	} else {
		d.Set("epg_name", epg)
	}

	domainCount, err := epgCont.ArrayCount("domainAssociations")
	if err != nil {
		return fmt.Errorf("Unable to get Domain Associations list")
	}

	found := false
	for l := 0; l < domainCount; l++ {
		domainCont, err := epgCont.ArrayElement(l, "domainAssociations")
		if err != nil {
			return err
		}
		domain := models.StripQuotes(domainCont.S("dn").String())

		if domain == stateDomain {
			found = true
			d.SetId(fmt.Sprintf("%s/sites/%s-%s/anps/%s/epgs/%s/domainAssociations/%s", schemaId, siteId, templateName, anp, epg, stateDomain))
			d.Set("domain_type", models.StripQuotes(domainCont.S("domainType").String()))
			d.Set("domain_dn", stateDomain)
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
			break
		}
	}

	if !found {
		d.SetId("")
		return fmt.Errorf("Unable to find the Site ANP EPG Domain %s", stateDomain)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
