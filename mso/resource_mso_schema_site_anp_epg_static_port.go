package mso

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOSchemaSiteAnpEpgStaticPort() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteAnpEpgStaticPortCreate,
		Read:   resourceMSOSchemaSiteAnpEpgStaticPortRead,
		Update: resourceMSOSchemaSiteAnpEpgStaticPortUpdate,
		Delete: resourceMSOSchemaSiteAnpEpgStaticPortDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaSiteAnpEpgStaticPortImport,
		},

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"schema_id": &schema.Schema{
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
			"template_name": &schema.Schema{
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
			"path_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"port",
					"vpc",
					"dpc",
				}, false),
			},
			"pod": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"leaf": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"path": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"vlan": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"deployment_immediacy": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"immediate",
					"lazy",
				}, false),
			},
			"fex": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"micro_seg_vlan": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"native",
					"regular",
					"untagged",
				}, false),
			},
		}),
	}
}

func resourceMSOSchemaSiteAnpEpgStaticPortImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	msoClient := m.(*client.Client)
	get_attribute := strings.Split(d.Id(), "/")
	import_attribute := regexp.MustCompile("(.*)/path/(.*)")
	import_split := import_attribute.FindStringSubmatch(d.Id())
	schemaId := get_attribute[0]
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return nil, err
	}
	d.Set("schema_id", schemaId)
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return nil, fmt.Errorf("No Sites found")
	}
	stateSite := get_attribute[2]
	found := false
	stateTemplate := get_attribute[4]
	stateAnp := get_attribute[6]
	stateEpg := get_attribute[8]
	statepod := get_attribute[10]
	stateleaf := get_attribute[12]
	pathType := get_attribute[14]
	fex := get_attribute[16]
	statepath := import_split[2]

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return nil, err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())
		apiTemplate := models.StripQuotes(tempCont.S("templateName").String())

		if apiSite == stateSite && apiTemplate == stateTemplate {
			d.Set("site_id", apiSite)
			d.Set("template_name", apiTemplate)
			anpCount, err := tempCont.ArrayCount("anps")
			if err != nil {
				return nil, fmt.Errorf("Unable to get Anp list")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")
				if err != nil {
					return nil, err
				}
				anpRef := models.StripQuotes(anpCont.S("anpRef").String())
				re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/anps/(.*)")
				match := re.FindStringSubmatch(anpRef)
				if match[3] == stateAnp {
					d.Set("anp_name", match[3])
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
						re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/epgs/(.*)")
						match := re.FindStringSubmatch(apiEpgRef)
						apiEPG := match[3]
						if apiEPG == stateEpg {
							d.Set("epg_name", apiEPG)
							portCount, err := epgCont.ArrayCount("staticPorts")
							if err != nil {
								return nil, fmt.Errorf("Unable to get Static Port list")
							}
							for l := 0; l < portCount; l++ {
								portCont, err := epgCont.ArrayElement(l, "staticPorts")
								if err != nil {
									return nil, err
								}
								var portpath string
								if pathType == "port" && fex != "" {
									portpath = fmt.Sprintf("topology/%s/paths-%s/extpaths-%s/pathep-[%s]", statepod, stateleaf, fex, statepath)
								} else if pathType == "vpc" {
									portpath = fmt.Sprintf("topology/%s/protpaths-%s/pathep-[%s]", statepod, stateleaf, statepath)
								} else {
									portpath = fmt.Sprintf("topology/%s/paths-%s/pathep-[%s]", statepod, stateleaf, statepath)
								}
								apiportpath := models.StripQuotes(portCont.S("path").String())
								apiType := models.StripQuotes(portCont.S("type").String())
								if portpath == apiportpath && pathType == apiType {
									d.SetId(apiportpath)
									if portCont.Exists("type") {
										d.Set("path_type", models.StripQuotes(portCont.S("type").String()))
									}
									if portCont.Exists("path") {
										d.Set("pod", statepod)
										d.Set("leaf", stateleaf)
										d.Set("path", statepath)
										d.Set("fex", fex)
									}
									if portCont.Exists("portEncapVlan") {
										tempvar, _ := strconv.Atoi(fmt.Sprintf("%v", portCont.S("portEncapVlan")))
										d.Set("vlan", tempvar)
									}
									if portCont.Exists("deploymentImmediacy") {
										d.Set("deployment_immediacy", models.StripQuotes(portCont.S("deploymentImmediacy").String()))
									}
									if portCont.Exists("microSegVlan") {
										tempvar1, _ := strconv.Atoi(fmt.Sprintf("%v", portCont.S("microSegVlan")))
										d.Set("micro_seg_vlan", tempvar1)
									}

									if portCont.Exists("mode") {
										d.Set("mode", models.StripQuotes(portCont.S("mode").String()))
									}
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
		d.SetId("")
		return nil, fmt.Errorf("Unable to find the static port entry")
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaSiteAnpEpgStaticPortCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Static Port Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	stateSiteId := d.Get("site_id").(string)
	stateTemplateName := d.Get("template_name").(string)
	stateANPName := d.Get("anp_name").(string)
	stateEpgName := d.Get("epg_name").(string)

	var pathType, pod, leaf, path, deploymentImmediacy, mode, fex string
	var vlan, microsegvlan int

	if tempVar, ok := d.GetOk("path_type"); ok {
		pathType = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("pod"); ok {
		pod = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("leaf"); ok {
		leaf = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("path"); ok {
		path = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("deployment_immediacy"); ok {
		deploymentImmediacy = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("mode"); ok {
		mode = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("vlan"); ok {
		vlan = tempVar.(int)
	}
	if tempVar, ok := d.GetOk("micro_seg_vlan"); ok {
		microsegvlan = tempVar.(int)
	}
	if tempVar, ok := d.GetOk("fex"); ok {
		fex = tempVar.(string)
	}

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return fmt.Errorf("No Sites found")
	}

	foundEpg := false
	foundAnp := false

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())
		apiTemplate := models.StripQuotes(tempCont.S("templateName").String())

		if apiSite == stateSiteId && apiTemplate == stateTemplateName {
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
				if match[3] == stateANPName {
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
						if apiEPG == stateEpgName {
							foundEpg = true
							break

						}
					}
					if !foundEpg {
						log.Printf("[DEBUG] Site Anp Epg: Beginning Creation")
						anpEpgRefMap := make(map[string]interface{})
						anpEpgRefMap["schemaId"] = schemaId
						anpEpgRefMap["templateName"] = stateTemplateName
						anpEpgRefMap["anpName"] = stateANPName
						anpEpgRefMap["epgName"] = stateEpgName

						pathEpg := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/-", apiSite, apiTemplate, stateANPName)
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
				anpRefMap["templateName"] = stateTemplateName
				anpRefMap["anpName"] = stateANPName

				pathAnp := fmt.Sprintf("/sites/%s-%s/anps/-", stateSiteId, stateTemplateName)
				anpStruct := models.NewSchemaSiteAnp("add", pathAnp, anpRefMap)

				_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpStruct)
				if err != nil {
					return err
				}

				log.Printf("[DEBUG] Site Anp Epg: Beginning Creation")

				anpEpgRefMap := make(map[string]interface{})
				anpEpgRefMap["schemaId"] = schemaId
				anpEpgRefMap["templateName"] = stateTemplateName
				anpEpgRefMap["anpName"] = stateANPName
				anpEpgRefMap["epgName"] = stateEpgName

				pathEpg := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/-", stateSiteId, stateTemplateName, stateANPName)
				//private_link_label argument used in resource site_anp_epg is set to nil here
				anpEpgStruct := models.NewSchemaSiteAnpEpg("add", pathEpg, nil, anpEpgRefMap)

				_, ers := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpEpgStruct)
				if ers != nil {
					return ers
				}
			}
		}
	}
	var portpath string
	if pathType == "port" && fex != "" {
		portpath = fmt.Sprintf("topology/%s/paths-%s/extpaths-%s/pathep-[%s]", pod, leaf, fex, path)
	} else if pathType == "vpc" {
		portpath = fmt.Sprintf("topology/%s/protpaths-%s/pathep-[%s]", pod, leaf, path)
	} else {
		portpath = fmt.Sprintf("topology/%s/paths-%s/pathep-[%s]", pod, leaf, path)
	}

	pathsp := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/staticPorts/-", stateSiteId, stateTemplateName, stateANPName, stateEpgName)
	staticStruct := models.NewSchemaSiteAnpEpgStaticPort("add", pathsp, pathType, portpath, vlan, deploymentImmediacy, microsegvlan, mode)
	_, errs := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), staticStruct)
	if errs != nil {
		return errs
	}
	return resourceMSOSchemaSiteAnpEpgStaticPortRead(d, m)
}

func resourceMSOSchemaSiteAnpEpgStaticPortRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	var fex, pathType string
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
	statepod := d.Get("pod").(string)
	stateleaf := d.Get("leaf").(string)
	statepath := d.Get("path").(string)
	if tempVar, ok := d.GetOk("fex"); ok {
		fex = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("path_type"); ok {
		pathType = tempVar.(string)
	}
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())
		apiTemplate := models.StripQuotes(tempCont.S("templateName").String())

		if apiSite == stateSite && apiTemplate == stateTemplate {
			d.Set("site_id", apiSite)
			d.Set("template_name", apiTemplate)
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
				if match[3] == stateAnp {
					d.Set("anp_name", match[3])
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
						if apiEPG == stateEpg {
							d.Set("epg_name", apiEPG)
							portCount, err := epgCont.ArrayCount("staticPorts")
							if err != nil {
								return fmt.Errorf("Unable to get Static Port list")
							}
							for l := 0; l < portCount; l++ {
								portCont, err := epgCont.ArrayElement(l, "staticPorts")
								if err != nil {
									return err
								}
								var portpath string
								if pathType == "port" && fex != "" {
									portpath = fmt.Sprintf("topology/%s/paths-%s/extpaths-%s/pathep-[%s]", statepod, stateleaf, fex, statepath)
								} else if pathType == "vpc" {
									portpath = fmt.Sprintf("topology/%s/protpaths-%s/pathep-[%s]", statepod, stateleaf, statepath)
								} else {
									portpath = fmt.Sprintf("topology/%s/paths-%s/pathep-[%s]", statepod, stateleaf, statepath)
								}
								apiportpath := models.StripQuotes(portCont.S("path").String())
								apiType := models.StripQuotes(portCont.S("type").String())
								if portpath == apiportpath && pathType == apiType {
									d.SetId(apiportpath)
									if portCont.Exists("type") {
										d.Set("type", models.StripQuotes(portCont.S("type").String()))
									}
									if portCont.Exists("path") {
										d.Set("pod", statepod)
										d.Set("leaf", stateleaf)
										d.Set("path", statepath)
										d.Set("fex", fex)
									}
									if portCont.Exists("portEncapVlan") {
										tempvar, err := strconv.Atoi(fmt.Sprintf("%v", portCont.S("portEncapVlan")))
										if err != nil {
											return err
										}
										d.Set("vlan", tempvar)
									}
									if portCont.Exists("deploymentImmediacy") {
										d.Set("deployment_immediacy", models.StripQuotes(portCont.S("deploymentImmediacy").String()))
									}
									if portCont.Exists("microSegVlan") {
										tempvar1, err := strconv.Atoi(fmt.Sprintf("%v", portCont.S("microSegVlan")))
										if err != nil {
											return err
										}
										d.Set("micro_seg_vlan", tempvar1)
									}
									if portCont.Exists("mode") {
										d.Set("mode", models.StripQuotes(portCont.S("mode").String()))
									}
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
		d.SetId("")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}

func resourceMSOSchemaSiteAnpEpgStaticPortUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Static Port: Beginning Updation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	stateSiteId := d.Get("site_id").(string)
	stateTemplateName := d.Get("template_name").(string)
	stateANPName := d.Get("anp_name").(string)
	stateEpgName := d.Get("epg_name").(string)

	var pathType, pod, leaf, path, deploymentImmediacy, mode, fex string
	var vlan, microsegvlan int

	if tempVar, ok := d.GetOk("path_type"); ok {
		pathType = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("pod"); ok {
		pod = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("leaf"); ok {
		leaf = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("path"); ok {
		path = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("deployment_immediacy"); ok {
		deploymentImmediacy = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("mode"); ok {
		mode = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("vlan"); ok {
		vlan = tempVar.(int)
	}
	if tempVar, ok := d.GetOk("micro_seg_vlan"); ok {
		microsegvlan = tempVar.(int)
	}
	if tempVar, ok := d.GetOk("fex"); ok {
		fex = tempVar.(string)
	}

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return fmt.Errorf("No Sites found")
	}

	found := false

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())
		apiTemplate := models.StripQuotes(tempCont.S("templateName").String())

		if apiSite == stateSiteId && apiTemplate == stateTemplateName {
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
				if match[3] == stateANPName {
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
						if apiEPG == stateEpgName {
							portCount, err := epgCont.ArrayCount("staticPorts")
							if err != nil {
								return fmt.Errorf("Unable to get Static Port list")
							}
							for l := 0; l < portCount; l++ {
								portCont, err := epgCont.ArrayElement(l, "staticPorts")
								if err != nil {
									return err
								}
								var portpath string
								if pathType == "port" && fex != "" {
									portpath = fmt.Sprintf("topology/%s/paths-%s/extpaths-%s/pathep-[%s]", pod, leaf, fex, path)
								} else if pathType == "vpc" {
									portpath = fmt.Sprintf("topology/%s/protpaths-%s/pathep-[%s]", pod, leaf, path)
								} else {
									portpath = fmt.Sprintf("topology/%s/paths-%s/pathep-[%s]", pod, leaf, path)
								}
								apiportpath := models.StripQuotes(portCont.S("path").String())
								if portpath == apiportpath {
									index := l
									path := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/staticPorts/%v", stateSiteId, stateTemplateName, stateANPName, stateEpgName, index)
									anpStruct := models.NewSchemaSiteAnpEpgStaticPort("replace", path, pathType, portpath, vlan, deploymentImmediacy, microsegvlan, mode)
									_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpStruct)

									if err != nil {
										return err
									}
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
		return fmt.Errorf("The specified parameters to update static port entry not found")
	}

	return resourceMSOSchemaSiteAnpEpgStaticPortRead(d, m)
}

func resourceMSOSchemaSiteAnpEpgStaticPortDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Anp: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	stateSite := d.Get("site_id").(string)
	stateTemplate := d.Get("template_name").(string)
	stateAnp := d.Get("anp_name").(string)
	stateEpg := d.Get("epg_name").(string)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return fmt.Errorf("No Sites found")
	}

	var pathType, pod, leaf, path, deploymentImmediacy, mode, fex string
	var vlan, microsegvlan int

	if tempVar, ok := d.GetOk("path_type"); ok {
		pathType = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("pod"); ok {
		pod = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("leaf"); ok {
		leaf = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("path"); ok {
		path = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("deployment_immediacy"); ok {
		deploymentImmediacy = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("mode"); ok {
		mode = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("vlan"); ok {
		vlan = tempVar.(int)
	}
	if tempVar, ok := d.GetOk("micro_seg_vlan"); ok {
		microsegvlan = tempVar.(int)
	}
	if tempVar, ok := d.GetOk("fex"); ok {
		fex = tempVar.(string)
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
				anpRef := models.StripQuotes(anpCont.S("anpRef").String())
				re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/anps/(.*)")
				match := re.FindStringSubmatch(anpRef)
				if match[3] == stateAnp {
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
						if apiEPG == stateEpg {
							portCount, err := epgCont.ArrayCount("staticPorts")
							if err != nil {
								return fmt.Errorf("Unable to get Static Port list")
							}
							for l := 0; l < portCount; l++ {
								portCont, err := epgCont.ArrayElement(l, "staticPorts")
								if err != nil {
									return err
								}
								var portpath string
								if pathType == "port" && fex != "" {
									portpath = fmt.Sprintf("topology/%s/paths-%s/extpaths-%s/pathep-[%s]", pod, leaf, fex, path)
								} else if pathType == "vpc" {
									portpath = fmt.Sprintf("topology/%s/protpaths-%s/pathep-[%s]", pod, leaf, path)
								} else {
									portpath = fmt.Sprintf("topology/%s/paths-%s/pathep-[%s]", pod, leaf, path)
								}
								apiportpath := models.StripQuotes(portCont.S("path").String())
								if portpath == apiportpath {
									index := l
									path := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/staticPorts/%v", stateSite, stateTemplate, stateAnp, stateEpg, index)
									anpStruct := models.NewSchemaSiteAnpEpgStaticPort("remove", path, pathType, portpath, vlan, deploymentImmediacy, microsegvlan, mode)
									response, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpStruct)

									// Ignoring Error with code 141: Resource Not Found when deleting
									if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
										return err
									}
									break
								}
							}
						}

					}
				}
			}
		}
	}
	d.SetId("")
	return resourceMSOSchemaSiteAnpEpgStaticPortRead(d, m)
}
