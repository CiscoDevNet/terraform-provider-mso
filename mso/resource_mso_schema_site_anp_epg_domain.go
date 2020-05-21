package mso

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceMSOSchemaSiteAnpEpgDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteAnpEpgDomainCreate,
		Update: resourceMSOSchemaSiteAnpEpgDomainUpdate,
		Read:   resourceMSOSchemaSiteAnpEpgDomainRead,
		Delete: resourceMSOSchemaSiteAnpEpgDomainDelete,

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
				ForceNew:     true,
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

			"domain_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"dn": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"deploy_immediacy": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"resolution_immediacy": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
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
			"vlan_encap_mode": &schema.Schema{
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
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			"port_encap_vlan_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"port_encap_vlan": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			"enhanced_lagpolicy_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enhanced_lagpolicy_dn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func resourceMSOSchemaSiteAnpEpgDomainCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Anp Epg Domain: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	siteId := d.Get("site_id").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)
	domainType := d.Get("domain_type").(string)
	domainName := d.Get("dn").(string)
	deployImmediacy := d.Get("deploy_immediacy").(string)
	resolutionImmediacy := d.Get("resolution_immediacy").(string)

	var DN, microSegVlanType, portEncapVlanType, vlanEncapMode, switchingMode, switchType, enhancedLagpolicyName, enhancedLagpolicyDn string
	var microSegVlan, portEncapVlan float64
	var allowMicroSegmentation bool

	if domainType == "vmmDomain" {
		DN = fmt.Sprintf("uni/vmmp-VMware/dom-%s", domainName)

	} else if domainType == "l3ExtDomain" {
		DN = fmt.Sprintf("uni/l3dom-%s", domainName)

	} else if domainType == "l2ExtDomain" {
		DN = fmt.Sprintf("uni/l2dom-%s", domainName)

	} else if domainType == "physicalDomain" {
		DN = fmt.Sprintf("uni/phys-%s", domainName)

	} else if domainType == "fibreChannelDomain" {
		DN = fmt.Sprintf("uni/fc-%s", domainName)

	} else {
		DN = ""
	}

	vmmDomainPropertiesRefMap := make(map[string]interface{})

	if domainType == "vmmDomain" {
		if TempVar, ok := d.GetOk("micro_seg_vlan_type"); ok {
			microSegVlanType = TempVar.(string)
		}
		if TempVar, ok := d.GetOk("micro_seg_vlan"); ok {
			microSegVlan = TempVar.(float64)
		}
		if TempVar, ok := d.GetOk("port_encap_vlan_type"); ok {
			portEncapVlanType = TempVar.(string)
		}
		if TempVar, ok := d.GetOk("port_encap_vlan"); ok {
			portEncapVlan = TempVar.(float64)
		}
		if TempVar, ok := d.GetOk("enhanced_lagpolicy_name"); ok {
			enhancedLagpolicyName = TempVar.(string)
		}
		if TempVar, ok := d.GetOk("enhanced_lagpolicy_dn"); ok {
			enhancedLagpolicyDn = TempVar.(string)
		}

		vlanEncapMode = "dynamic"
		if TempVar, ok := d.GetOk("vlan_encap_mode"); ok {
			vlanEncapMode = TempVar.(string)
		}

		switchingMode = "native"
		if TempVar, ok := d.GetOk("switching_mode"); ok {
			switchingMode = TempVar.(string)
		}

		switchType = "default"
		if TempVar, ok := d.GetOk("switch_type"); ok {
			switchType = TempVar.(string)
		}
		if TempVar, ok := d.GetOk("allow_micro_segmentation"); ok {
			allowMicroSegmentation = TempVar.(bool)
		}

		if portEncapVlanType != "" && portEncapVlan != 0 {
			portEncapVlanRefMap := make(map[string]interface{})
			portEncapVlanRefMap["vlanType"] = portEncapVlanType
			portEncapVlanRefMap["vlan"] = portEncapVlan

			vmmDomainPropertiesRefMap["portEncapVlan"] = portEncapVlanRefMap
		}

		if microSegVlanType != "" && microSegVlan != 0 {
			microSegVlanRefMap := make(map[string]interface{})
			microSegVlanRefMap["vlanType"] = microSegVlanType
			microSegVlanRefMap["vlan"] = microSegVlan

			vmmDomainPropertiesRefMap["microSegVlan"] = microSegVlanRefMap
		}

		if enhancedLagpolicyName != "" && enhancedLagpolicyDn != "" {
			enhancedLagPolRefMap := make(map[string]interface{})
			enhancedLagPolRefMap["name"] = enhancedLagpolicyName
			enhancedLagPolRefMap["dn"] = enhancedLagpolicyDn

			epgLagPolRefMap := make(map[string]interface{})
			epgLagPolRefMap["enhancedLagPol"] = enhancedLagPolRefMap

			vmmDomainPropertiesRefMap["epgLagPol"] = epgLagPolRefMap
		}

		vmmDomainPropertiesRefMap["allowMicroSegmentation"] = allowMicroSegmentation
		vmmDomainPropertiesRefMap["switchingMode"] = switchingMode
		vmmDomainPropertiesRefMap["switchType"] = switchType
		vmmDomainPropertiesRefMap["vlanEncapMode"] = vlanEncapMode

	} else {
		log.Print("Passing Blank Value to the Model")
	}

	foundAnp := false
	foundEpg := false
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return fmt.Errorf("No Sites found")
	}
	//found := false

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())
		apiTemplate := models.StripQuotes(tempCont.S("templateName").String())

		if apiSite == siteId && apiTemplate == templateName {
			anpCount, err := tempCont.ArrayCount("anps")
			if err != nil {
				return fmt.Errorf("Unable to get Anp list")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")
				if err != nil {
					return err
				}

				anpRef := models.StripQuotes(anpCont.S("anpRef").String())

				re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/anps/(.*)")
				match := re.FindStringSubmatch(anpRef)

				if match[3] == anpName {

					foundAnp = true
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

						if apiEPG == epgName {
							foundEpg = true
							break
						}
					}

					if !foundEpg {
						log.Printf("[DEBUG] Site Anp Epg: Beginning Creation")
						anpEpgRefMap := make(map[string]interface{})
						anpEpgRefMap["schemaId"] = schemaId
						anpEpgRefMap["templateName"] = apiTemplate
						anpEpgRefMap["anpName"] = anpName
						anpEpgRefMap["epgName"] = epgName

						pathEpg := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/-", apiSite, apiTemplate, anpName)
						anpEpgStruct := models.NewSchemaSiteAnpEpg("add", pathEpg, anpEpgRefMap)

						_, ers := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpEpgStruct)
						if ers != nil {
							return ers
						}
						break

					}
				}
			}

			if !foundAnp {
				log.Printf("[DEBUG] Site Anp: Beginning Creation")

				anpRefMap := make(map[string]interface{})
				anpRefMap["schemaId"] = schemaId
				anpRefMap["templateName"] = apiTemplate
				anpRefMap["anpName"] = anpName

				pathAnp := fmt.Sprintf("/sites/%s-%s/anps/-", apiSite, apiTemplate)
				anpStruct := models.NewSchemaSiteAnp("add", pathAnp, anpRefMap)

				_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpStruct)
				if err != nil {
					return err
				}

				log.Printf("[DEBUG] Site Anp Epg: Beginning Creation")
				anpEpgRefMap := make(map[string]interface{})
				anpEpgRefMap["schemaId"] = schemaId
				anpEpgRefMap["templateName"] = apiTemplate
				anpEpgRefMap["anpName"] = anpName
				anpEpgRefMap["epgName"] = epgName

				pathEpg := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/-", apiSite, apiTemplate, anpName)
				anpEpgStruct := models.NewSchemaSiteAnpEpg("add", pathEpg, anpEpgRefMap)

				_, ers := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpEpgStruct)
				if ers != nil {
					return ers
				}

			}

		}
	}

	path := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/domainAssociations/-", siteId, templateName, anpName, epgName)
	anpEpgDomainStruct := models.NewSchemaSiteAnpEpgDomain("add", path, domainType, DN, deployImmediacy, resolutionImmediacy, vmmDomainPropertiesRefMap)

	_, errs := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpEpgDomainStruct)
	if errs != nil {
		return errs
	}

	d.SetId(DN)
	return resourceMSOSchemaSiteAnpEpgDomainRead(d, m)
}

func resourceMSOSchemaSiteAnpEpgDomainRead(d *schema.ResourceData, m interface{}) error {
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
	stateAnp := d.Get("anp_name").(string)
	stateEpg := d.Get("epg_name").(string)
	domain := d.Get("dn").(string)
	domainType := d.Get("domain_type").(string)

	var stateDomain string

	if domainType == "vmmDomain" {
		stateDomain = fmt.Sprintf("uni/vmmp-VMware/dom-%s", domain)

	} else if domainType == "l3ExtDomain" {
		stateDomain = fmt.Sprintf("uni/l3dom-%s", domain)

	} else if domainType == "l2ExtDomain" {
		stateDomain = fmt.Sprintf("uni/l2dom-%s", domain)

	} else if domainType == "physicalDomain" {
		stateDomain = fmt.Sprintf("uni/phys-%s", domain)

	} else if domainType == "fibreChannelDomain" {
		stateDomain = fmt.Sprintf("uni/fc-%s", domain)

	} else {
		stateDomain = ""
	}

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())

		if apiSite == stateSite {
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
									d.Set("dn", domain)
									d.Set("deployment_immediacy", models.StripQuotes(domainCont.S("deployImmediacy").String()))
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
											d.Set("enhanced_lagpolicy_name", models.StripQuotes(domainCont.S("epgLagPol", "enhancedLagPol", "name").String()))
											d.Set("enhanced_lagpolicy_dn", models.StripQuotes(domainCont.S("epgLagPol", "enhancedLagPol", "dn").String()))
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
		d.SetId("")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
func resourceMSOSchemaSiteAnpEpgDomainUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Anp Epg Domain: Beginning Update")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	siteId := d.Get("site_id").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)
	domainType := d.Get("domain_type").(string)
	domainName := d.Get("dn").(string)
	deployImmediacy := d.Get("deploy_immediacy").(string)
	resolutionImmediacy := d.Get("resolution_immediacy").(string)

	var DN, microSegVlanType, portEncapVlanType, vlanEncapMode, switchingMode, switchType, enhancedLagpolicyName, enhancedLagpolicyDn string
	var microSegVlan, portEncapVlan float64
	var allowMicroSegmentation bool

	if domainType == "vmmDomain" {
		DN = fmt.Sprintf("uni/vmmp-VMware/dom-%s", domainName)

	} else if domainType == "l3ExtDomain" {
		DN = fmt.Sprintf("uni/l3dom-%s", domainName)

	} else if domainType == "l2ExtDomain" {
		DN = fmt.Sprintf("uni/l2dom-%s", domainName)

	} else if domainType == "physicalDomain" {
		DN = fmt.Sprintf("uni/phys-%s", domainName)

	} else if domainType == "fibreChannelDomain" {
		DN = fmt.Sprintf("uni/fc-%s", domainName)

	} else {
		DN = ""
	}

	vmmDomainPropertiesRefMap := make(map[string]interface{})
	if domainType == "vmmDomain" {
		if TempVar, ok := d.GetOk("micro_seg_vlan_type"); ok {
			microSegVlanType = TempVar.(string)
		}
		if TempVar, ok := d.GetOk("port_encap_vlan_type"); ok {
			portEncapVlanType = TempVar.(string)
		}
		if TempVar, ok := d.GetOk("vlan_encap_mode"); ok {
			vlanEncapMode = TempVar.(string)
		}
		if TempVar, ok := d.GetOk("switching_mode"); ok {
			switchingMode = TempVar.(string)
		}
		if TempVar, ok := d.GetOk("switch_type"); ok {
			switchType = TempVar.(string)
		}
		if TempVar, ok := d.GetOk("enhanced_lagpolicy_name"); ok {
			enhancedLagpolicyName = TempVar.(string)
		}
		if TempVar, ok := d.GetOk("enhanced_lagpolicy_dn"); ok {
			enhancedLagpolicyDn = TempVar.(string)
		}
		if TempVar, ok := d.GetOk("allow_micro_segmentation"); ok {
			allowMicroSegmentation = TempVar.(bool)
		}
		if TempVar, ok := d.GetOk("micro_seg_vlan"); ok {
			microSegVlan = TempVar.(float64)
		}
		if TempVar, ok := d.GetOk("port_encap_vlan"); ok {
			portEncapVlan = TempVar.(float64)
		}

		if portEncapVlanType != "" && portEncapVlan != 0 {
			portEncapVlanRefMap := make(map[string]interface{})
			portEncapVlanRefMap["vlanType"] = portEncapVlanType
			portEncapVlanRefMap["vlan"] = portEncapVlan

			vmmDomainPropertiesRefMap["portEncapVlan"] = portEncapVlanRefMap
		}

		if microSegVlanType != "" && microSegVlan != 0 {
			microSegVlanRefMap := make(map[string]interface{})
			microSegVlanRefMap["vlanType"] = microSegVlanType
			microSegVlanRefMap["vlan"] = microSegVlan

			vmmDomainPropertiesRefMap["microSegVlan"] = microSegVlanRefMap
		}

		if enhancedLagpolicyName != "" && enhancedLagpolicyDn != "" {
			enhancedLagPolRefMap := make(map[string]interface{})
			enhancedLagPolRefMap["name"] = enhancedLagpolicyName
			enhancedLagPolRefMap["dn"] = enhancedLagpolicyDn

			epgLagPolRefMap := make(map[string]interface{})
			epgLagPolRefMap["enhancedLagPol"] = enhancedLagPolRefMap

			vmmDomainPropertiesRefMap["epgLagPol"] = epgLagPolRefMap
		}

		vmmDomainPropertiesRefMap["allowMicroSegmentation"] = allowMicroSegmentation
		vmmDomainPropertiesRefMap["switchingMode"] = switchingMode
		vmmDomainPropertiesRefMap["switchType"] = switchType
		vmmDomainPropertiesRefMap["vlanEncapMode"] = vlanEncapMode

	} else {
		log.Print("Passing Blank Value to the Model")
	}

	id := d.Id()

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	index, err := indexCount(cont, siteId, anpName, epgName, id)
	if err != nil {
		return err
	}
	if index == -1 {
		fmt.Errorf("The given Anp Epg Domain is not found")
	}
	indexs := strconv.Itoa(index)

	path := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/domainAssociations/%s", siteId, templateName, anpName, epgName, indexs)
	anpEpgDomainStruct := models.NewSchemaSiteAnpEpgDomain("replace", path, domainType, DN, deployImmediacy, resolutionImmediacy, vmmDomainPropertiesRefMap)

	_, errs := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpEpgDomainStruct)
	if errs != nil {
		return errs
	}

	return resourceMSOSchemaSiteAnpEpgDomainRead(d, m)
}

func resourceMSOSchemaSiteAnpEpgDomainDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Anp Epg Domain: Beginning Update")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	siteId := d.Get("site_id").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)
	domainType := d.Get("domain_type").(string)
	domainName := d.Get("dn").(string)
	deployImmediacy := d.Get("deploy_immediacy").(string)
	resolutionImmediacy := d.Get("resolution_immediacy").(string)

	var DN, microSegVlanType, portEncapVlanType, vlanEncapMode, switchingMode, switchType, enhancedLagpolicyName, enhancedLagpolicyDn string
	var microSegVlan, portEncapVlan float64
	var allowMicroSegmentation bool

	if domainType == "vmmDomain" {
		DN = fmt.Sprintf("uni/vmmp-VMware/dom-%s", domainName)

	} else if domainType == "l3ExtDomain" {
		DN = fmt.Sprintf("uni/l3dom-%s", domainName)

	} else if domainType == "l2ExtDomain" {
		DN = fmt.Sprintf("uni/l2dom-%s", domainName)

	} else if domainType == "physicalDomain" {
		DN = fmt.Sprintf("uni/phys-%s", domainName)

	} else if domainType == "fibreChannelDomain" {
		DN = fmt.Sprintf("uni/fc-%s", domainName)

	} else {
		DN = ""
	}

	vmmDomainPropertiesRefMap := make(map[string]interface{})
	if domainType == "vmmDomain" {
		if TempVar, ok := d.GetOk("micro_seg_vlan_type"); ok {
			microSegVlanType = TempVar.(string)
		}
		if TempVar, ok := d.GetOk("port_encap_vlan_type"); ok {
			portEncapVlanType = TempVar.(string)
		}
		if TempVar, ok := d.GetOk("vlan_encap_mode"); ok {
			vlanEncapMode = TempVar.(string)
		}
		if TempVar, ok := d.GetOk("switching_mode"); ok {
			switchingMode = TempVar.(string)
		}
		if TempVar, ok := d.GetOk("switch_type"); ok {
			switchType = TempVar.(string)
		}
		if TempVar, ok := d.GetOk("enhanced_lagpolicy_name"); ok {
			enhancedLagpolicyName = TempVar.(string)
		}
		if TempVar, ok := d.GetOk("enhanced_lagpolicy_dn"); ok {
			enhancedLagpolicyDn = TempVar.(string)
		}
		if TempVar, ok := d.GetOk("allow_micro_segmentation"); ok {
			allowMicroSegmentation = TempVar.(bool)
		}
		if TempVar, ok := d.GetOk("micro_seg_vlan"); ok {
			microSegVlan = TempVar.(float64)
		}
		if TempVar, ok := d.GetOk("port_encap_vlan"); ok {
			portEncapVlan = TempVar.(float64)
		}

		if portEncapVlanType != "" && portEncapVlan != 0 {
			portEncapVlanRefMap := make(map[string]interface{})
			portEncapVlanRefMap["vlanType"] = portEncapVlanType
			portEncapVlanRefMap["vlan"] = portEncapVlan

			vmmDomainPropertiesRefMap["portEncapVlan"] = portEncapVlanRefMap
		}

		if microSegVlanType != "" && microSegVlan != 0 {
			microSegVlanRefMap := make(map[string]interface{})
			microSegVlanRefMap["vlanType"] = microSegVlanType
			microSegVlanRefMap["vlan"] = microSegVlan

			vmmDomainPropertiesRefMap["microSegVlan"] = microSegVlanRefMap
		}

		if enhancedLagpolicyName != "" && enhancedLagpolicyDn != "" {
			enhancedLagPolRefMap := make(map[string]interface{})
			enhancedLagPolRefMap["name"] = enhancedLagpolicyName
			enhancedLagPolRefMap["dn"] = enhancedLagpolicyDn

			epgLagPolRefMap := make(map[string]interface{})
			epgLagPolRefMap["enhancedLagPol"] = enhancedLagPolRefMap

			vmmDomainPropertiesRefMap["epgLagPol"] = epgLagPolRefMap
		}

		vmmDomainPropertiesRefMap["allowMicroSegmentation"] = allowMicroSegmentation
		vmmDomainPropertiesRefMap["switchingMode"] = switchingMode
		vmmDomainPropertiesRefMap["switchType"] = switchType
		vmmDomainPropertiesRefMap["vlanEncapMode"] = vlanEncapMode

	} else {
		log.Print("Passing Blank Value to the Model")
	}

	id := d.Id()
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	index, err := indexCount(cont, siteId, anpName, epgName, id)
	if err != nil {
		return err
	}
	if index == -1 {
		fmt.Errorf("The given Anp Epg Domain is not found")
	}
	indexs := strconv.Itoa(index)

	path := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/domainAssociations/%s", siteId, templateName, anpName, epgName, indexs)
	anpEpgDomainStruct := models.NewSchemaSiteAnpEpgDomain("remove", path, domainType, DN, deployImmediacy, resolutionImmediacy, vmmDomainPropertiesRefMap)

	_, errs := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpEpgDomainStruct)
	if errs != nil {
		return errs
	}
	d.SetId("")
	return nil
}

func indexCount(cont *container.Container, stateSite, stateAnp, stateEpg, stateDomain string) (int, error) {
	found := false
	index := -1
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return index, fmt.Errorf("No Sites found")
	}

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return index, err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())

		if apiSite == stateSite {
			anpCount, err := tempCont.ArrayCount("anps")
			if err != nil {
				return index, fmt.Errorf("Unable to get Anp list")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")
				if err != nil {
					return index, err
				}
				apiAnpRef := models.StripQuotes(anpCont.S("anpRef").String())
				split := strings.Split(apiAnpRef, "/")
				apiAnp := split[6]
				if apiAnp == stateAnp {
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return index, fmt.Errorf("Unable to get EPG list")
					}
					for k := 0; k < epgCount; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return index, err
						}
						apiEpgRef := models.StripQuotes(epgCont.S("epgRef").String())
						split := strings.Split(apiEpgRef, "/")
						apiEPG := split[8]
						if apiEPG == stateEpg {
							domainCount, err := epgCont.ArrayCount("domainAssociations")
							if err != nil {
								return index, fmt.Errorf("Unable to get Domain Associations list")
							}
							for l := 0; l < domainCount; l++ {
								domainCont, err := epgCont.ArrayElement(l, "domainAssociations")
								if err != nil {
									return index, err
								}
								apiDomain := models.StripQuotes(domainCont.S("dn").String())
								if apiDomain == stateDomain {
									log.Println("found correct domain")
									index = l
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
	return index, nil

}
