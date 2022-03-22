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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOSchemaSiteAnpEpgDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteAnpEpgDomainCreate,
		Update: resourceMSOSchemaSiteAnpEpgDomainUpdate,
		Read:   resourceMSOSchemaSiteAnpEpgDomainRead,
		Delete: resourceMSOSchemaSiteAnpEpgDomainDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaSiteAnpEpgDomainImport,
		},

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
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == "" {
						return true
					} else {
						return false
					}

				},
			},
			"dn": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringLenBetween(1, 1000),
				ConflictsWith: []string{"domain_name"},
				Deprecated:    "use domain_dn alone or domain_name in association with domain_type and vmm_domain_type if applicable.",
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

func resourceMSOSchemaSiteAnpEpgDomainImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())

	msoClient := m.(*client.Client)
	get_dn := d.Id()
	get_attribute := strings.Split(d.Id(), "/")
	schemaId := get_attribute[0]

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return nil, err
	}
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return nil, fmt.Errorf("No Sites found")
	}
	templateInfo := regexp.MustCompile("-").Split(get_attribute[2], 2)
	site_template := strings.Split(get_attribute[2], "-")
	stateSite := site_template[0]
	stateTemplate := templateInfo[1]
	found := false
	stateAnp := get_attribute[4]
	stateEpg := get_attribute[6]

	var stateDomain string

	vmmp_match, _ := regexp.MatchString(".*/uni/vmmp-.*", get_dn)
	l3dom_match, _ := regexp.MatchString(".*/uni/l3dom-.*", get_dn)
	l2dom_match, _ := regexp.MatchString(".*/uni/l2dom-.*", get_dn)
	phys_match, _ := regexp.MatchString(".*/uni/phys-.*", get_dn)
	fc_match, _ := regexp.MatchString(".*/uni/fc-.*", get_dn)
	re_domain := regexp.MustCompile("(.*)/uni/(.*)-(.*)")
	match_domain := re_domain.FindStringSubmatch(get_dn)
	d.Set("domain_name", match_domain[3])
	if vmmp_match {
		re_vmmDomain := regexp.MustCompile("uni/vmmp-(.*)/dom-(.*)")
		match_vmmDomain := re_vmmDomain.FindStringSubmatch(get_dn)
		d.Set("vmm_domain_type", match_vmmDomain[1])
		d.Set("domain_name", match_vmmDomain[2])
		stateDomain = match_vmmDomain[0]
	} else if l2dom_match {
		re_domain := regexp.MustCompile("uni/l2dom-(.*)")
		match_domain := re_domain.FindStringSubmatch(get_dn)
		stateDomain = match_domain[0]
	} else if l3dom_match {
		re_domain := regexp.MustCompile("uni/l3dom-(.*)")
		match_domain := re_domain.FindStringSubmatch(get_dn)
		stateDomain = match_domain[0]
	} else if phys_match {
		re_domain := regexp.MustCompile("uni/phys(.*)")
		match_domain := re_domain.FindStringSubmatch(get_dn)
		stateDomain = match_domain[0]
	} else if fc_match {
		re_domain := regexp.MustCompile("uni/fc(.*)")
		match_domain := re_domain.FindStringSubmatch(get_dn)
		stateDomain = match_domain[0]
	} else {
		stateDomain = ""
	}

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return nil, err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())
		apiTemplate := models.StripQuotes(tempCont.S("templateName").String())

		if apiSite == stateSite && apiTemplate == stateTemplate {
			anpCount, err := tempCont.ArrayCount("anps")
			if err != nil {
				return nil, fmt.Errorf("Unable to get Anp list")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")
				if err != nil {
					return nil, err
				}
				apiAnpRef := models.StripQuotes(anpCont.S("anpRef").String())
				split := strings.Split(apiAnpRef, "/")
				apiAnp := split[6]
				if apiAnp == stateAnp {
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return nil, fmt.Errorf("Unable to get EPG list")
					}
					for k := 0; k < epgCount; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return nil, err
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
								return nil, fmt.Errorf("Unable to get Domain Associations list")
							}
							for l := 0; l < domainCount; l++ {
								domainCont, err := epgCont.ArrayElement(l, "domainAssociations")
								if err != nil {
									return nil, err
								}
								apiDomain := models.StripQuotes(domainCont.S("dn").String())

								if apiDomain == stateDomain {
									d.SetId(apiDomain)
									d.Set("site_id", apiSite)
									d.Set("domain_type", models.StripQuotes(domainCont.S("domainType").String()))
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
		return nil, fmt.Errorf("Unable to find the Site Anp Epg Domain %s", stateDomain)
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
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
	vmmDomainType := d.Get("vmm_domain_type").(string)
	domainNameDnOld := d.Get("dn").(string)
	deployImmediacy := d.Get("deploy_immediacy").(string)
	resolutionImmediacy := d.Get("resolution_immediacy").(string)

	var DN, microSegVlanType, portEncapVlanType, vlanEncapMode, switchingMode, switchType, enhancedLagpolicyName, enhancedLagpolicyDn, domainName string
	var microSegVlan, portEncapVlan float64
	var allowMicroSegmentation, checkDomainTypeFromDN bool

	_, ok_oldName := d.GetOk("dn")
	tempVarName, ok_name := d.GetOk("domain_name")
	tempVarDn, ok_dn := d.GetOk("domain_dn")

	if !ok_oldName && !ok_name && !ok_dn {
		return fmt.Errorf("domain_dn or domain_name in association with domain_type and vmm_domain_type if applicable are required.")
	}

	if ok_name {
		domainName = tempVarName.(string)
	} else {
		domainName = domainNameDnOld
	}

	if ok_dn {
		DN = tempVarDn.(string)
		vmmp_match, _ := regexp.MatchString("uni/vmmp-.*", DN)
		checkDomainTypeFromDN = vmmp_match
		l3dom_match, _ := regexp.MatchString("uni/l3dom-.*", DN)
		l2dom_match, _ := regexp.MatchString("uni/l2dom-.*", DN)
		phys_match, _ := regexp.MatchString("uni/phys-.*", DN)
		fc_match, _ := regexp.MatchString("uni/fc-.*", DN)
		domainType = "vmmDomain"
		if l2dom_match {
			domainType = "l2ExtDomain"
		} else if l3dom_match {
			domainType = "l3ExtDomain"
		} else if phys_match {
			domainType = "physicalDomain"
		} else if fc_match {
			domainType = "fibreChannelDomain"
		}
	} else {
		if domainType == "vmmDomain" {
			DN = fmt.Sprintf("uni/vmmp-%s/dom-%s", vmmDomainType, domainName)

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
	}

	vmmDomainPropertiesRefMap := make(map[string]interface{})

	if domainType == "vmmDomain" || checkDomainTypeFromDN {
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
		if TempVar, ok := d.GetOk("enhanced_lag_policy_name"); ok {
			enhancedLagpolicyName = TempVar.(string)
		}
		if TempVar, ok := d.GetOk("enhanced_lag_policy_dn"); ok {
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
						//private_link_label argument used in resource site_anp_epg is set to nil here
						anpEpgStruct := models.NewSchemaSiteAnpEpg("add", pathEpg, nil, anpEpgRefMap)

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
				//private_link_label argument used in resource site_anp_epg is set to nil here
				anpEpgStruct := models.NewSchemaSiteAnpEpg("add", pathEpg, nil, anpEpgRefMap)

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
	stateTemplate := d.Get("template_name").(string)
	stateAnp := d.Get("anp_name").(string)
	stateEpg := d.Get("epg_name").(string)
	domainNameDnOld := d.Get("dn").(string)
	domainType := d.Get("domain_type").(string)
	vmmDomainType := d.Get("vmm_domain_type").(string)
	stateDomain := d.Get("domain_dn").(string)

	var domainName string

	if tempVar, ok := d.GetOk("domain_name"); ok {
		domainName = tempVar.(string)
	} else {
		domainName = domainNameDnOld
	}

	if tempVar, ok := d.GetOk("domain_dn"); ok {
		stateDomain = tempVar.(string)
	} else {
		if domainType == "vmmDomain" {
			stateDomain = fmt.Sprintf("uni/vmmp-%s/dom-%s", vmmDomainType, domainName)

		} else if domainType == "l3ExtDomain" {
			stateDomain = fmt.Sprintf("uni/l3dom-%s", domainName)

		} else if domainType == "l2ExtDomain" {
			stateDomain = fmt.Sprintf("uni/l2dom-%s", domainName)

		} else if domainType == "physicalDomain" {
			stateDomain = fmt.Sprintf("uni/phys-%s", domainName)

		} else if domainType == "fibreChannelDomain" {
			stateDomain = fmt.Sprintf("uni/fc-%s", domainName)

		} else {
			stateDomain = ""
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
									d.Set("domain_dn", apiDomain)
									if _, ok := d.GetOk("domain_dn"); !ok {
										d.Set("domain_type", models.StripQuotes(domainCont.S("domainType").String()))
										vmmp_match, _ := regexp.MatchString("uni/vmmp-.*", apiDomain)
										if vmmp_match {
											re_vmmDomain := regexp.MustCompile("uni/vmmp-(.*)/dom-(.*)")
											match_vmmDomain := re_vmmDomain.FindStringSubmatch(apiDomain)
											d.Set("vmm_domain_type", match_vmmDomain[1])
										}
									}

									if tempVar, ok := d.GetOk("domain_name"); ok {
										d.Set("domain_name", tempVar.(string))

									} else if tempVar, ok := d.GetOk("dn"); ok {
										d.Set("dn", tempVar.(string))
									}
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
	vmmDomainType := d.Get("domain_type_name").(string)
	domainNameDnOld := d.Get("dn").(string)
	deployImmediacy := d.Get("deploy_immediacy").(string)
	resolutionImmediacy := d.Get("resolution_immediacy").(string)

	var DN, microSegVlanType, portEncapVlanType, vlanEncapMode, switchingMode, switchType, enhancedLagpolicyName, enhancedLagpolicyDn, domainName string
	var microSegVlan, portEncapVlan float64
	var allowMicroSegmentation, checkDomainTypeFromDN bool

	_, ok_oldName := d.GetOk("dn")
	tempVarName, ok_name := d.GetOk("domain_name")
	tempVarDn, ok_dn := d.GetOk("domain_dn")

	if !ok_oldName && !ok_name && !ok_dn {
		return fmt.Errorf("domain_dn or domain_name in association with domain_type and vmm_domain_type if applicable are required.")
	}

	if ok_name {
		domainName = tempVarName.(string)
	} else {
		domainName = domainNameDnOld
	}

	if ok_dn {
		DN = tempVarDn.(string)
		vmmp_match, _ := regexp.MatchString("uni/vmmp-.*", DN)
		checkDomainTypeFromDN = vmmp_match
		l3dom_match, _ := regexp.MatchString("uni/l3dom-.*", DN)
		l2dom_match, _ := regexp.MatchString("uni/l2dom-.*", DN)
		phys_match, _ := regexp.MatchString("uni/phys-.*", DN)
		fc_match, _ := regexp.MatchString("uni/fc-.*", DN)
		domainType = "vmmDomain"
		if l2dom_match {
			domainType = "l2ExtDomain"
		} else if l3dom_match {
			domainType = "l3ExtDomain"
		} else if phys_match {
			domainType = "physicalDomain"
		} else if fc_match {
			domainType = "fibreChannelDomain"
		}
	} else {
		if domainType == "vmmDomain" {
			DN = fmt.Sprintf("uni/vmmp-%s/dom-%s", vmmDomainType, domainName)

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
	}

	vmmDomainPropertiesRefMap := make(map[string]interface{})
	if domainType == "vmmDomain" || checkDomainTypeFromDN {
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
		if TempVar, ok := d.GetOk("enhanced_lag_policy_name"); ok {
			enhancedLagpolicyName = TempVar.(string)
		}
		if TempVar, ok := d.GetOk("enhanced_lag_policy_dn"); ok {
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
	vmmDomainType := d.Get("vmm_domain_type").(string)
	domainNameDnOld := d.Get("dn").(string)
	deployImmediacy := d.Get("deploy_immediacy").(string)
	resolutionImmediacy := d.Get("resolution_immediacy").(string)

	var DN, microSegVlanType, portEncapVlanType, vlanEncapMode, switchingMode, switchType, enhancedLagpolicyName, enhancedLagpolicyDn, domainName string
	var microSegVlan, portEncapVlan float64
	var allowMicroSegmentation, checkDomainTypeFromDN bool

	_, ok_oldName := d.GetOk("dn")
	tempVarName, ok_name := d.GetOk("domain_name")
	tempVarDn, ok_dn := d.GetOk("domain_dn")

	if !ok_oldName && !ok_name && !ok_dn {
		return fmt.Errorf("domain_dn or domain_name in association with domain_type and vmm_domain_type if applicable are required.")
	}

	if ok_name {
		domainName = tempVarName.(string)
	} else {
		domainName = domainNameDnOld
	}

	if ok_dn {
		DN = tempVarDn.(string)
		vmmp_match, _ := regexp.MatchString("uni/vmmp-.*", DN)
		checkDomainTypeFromDN = vmmp_match
		l3dom_match, _ := regexp.MatchString("uni/l3dom-.*", DN)
		l2dom_match, _ := regexp.MatchString("uni/l2dom-.*", DN)
		phys_match, _ := regexp.MatchString("uni/phys-.*", DN)
		fc_match, _ := regexp.MatchString("uni/fc-.*", DN)
		domainType = "vmmDomain"
		if l2dom_match {
			domainType = "l2ExtDomain"
		} else if l3dom_match {
			domainType = "l3ExtDomain"
		} else if phys_match {
			domainType = "physicalDomain"
		} else if fc_match {
			domainType = "fibreChannelDomain"
		}
	} else {
		if domainType == "vmmDomain" {
			DN = fmt.Sprintf("uni/vmmp-%s/dom-%s", vmmDomainType, domainName)

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
	}

	vmmDomainPropertiesRefMap := make(map[string]interface{})
	if domainType == "vmmDomain" || checkDomainTypeFromDN {
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
		if TempVar, ok := d.GetOk("enhanced_lag_policy_name"); ok {
			enhancedLagpolicyName = TempVar.(string)
		}
		if TempVar, ok := d.GetOk("enhanced_lag_policy_dn"); ok {
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
		d.SetId("")
		return nil
	}
	indexs := strconv.Itoa(index)

	path := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/domainAssociations/%s", siteId, templateName, anpName, epgName, indexs)
	anpEpgDomainStruct := models.NewSchemaSiteAnpEpgDomain("remove", path, domainType, DN, deployImmediacy, resolutionImmediacy, vmmDomainPropertiesRefMap)

	response, errs := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpEpgDomainStruct)

	// Ignoring Error with code 141: Resource Not Found when deleting
	if errs != nil && !(response.Exists("code") && response.S("code").String() == "141") {
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
