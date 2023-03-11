package mso

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOSchemaSiteAnpEpgBulkStaticPort() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteAnpEpgBulkStaticPortCreate,
		Read:   resourceMSOSchemaSiteAnpEpgBulkStaticPortRead,
		Update: resourceMSOSchemaSiteAnpEpgBulkStaticPortUpdate,
		Delete: resourceMSOSchemaSiteAnpEpgBulkStaticPortDelete,

		// Importer: &schema.ResourceImporter{
		// 	State: resourceMSOSchemaSiteAnpEpgBulkStaticPortImport,
		// },

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
			"static_ports": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"port",
								"vpc",
								"dpc",
							}, false),
						},
						"pod": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"leaf": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"path": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"vlan": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"deployment_immediacy": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"immediate",
								"lazy",
							}, false),
						},
						"fex": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"micro_seg_vlan": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"native",
								"regular",
								"untagged",
							}, false),
						},
					},
				},
				Required: true,
			},
		}),
	}
}

// func resourceMSOSchemaSiteAnpEpgBulkStaticPortImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

// }

func resourceMSOSchemaSiteAnpEpgBulkStaticPortCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Static Port Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	stateSiteId := d.Get("site_id").(string)
	stateTemplateName := d.Get("template_name").(string)
	stateANPName := d.Get("anp_name").(string)
	stateEpgName := d.Get("epg_name").(string)

	var pathType, pod, leaf, path, deploymentImmediacy, mode, fex, portpath string
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

	// Make the above  values to map and update the values in static ports. copy the parent map to local staticport map(to avoid overriding thge values for seconfd static port)

	staticPortsList := make([]interface{}, 0, 1)
	log.Printf("CHECK get staticPort : %v", d.Get("static_ports"))
	if staticPortsValue, ok := d.GetOk("static_ports"); ok {
		log.Printf("CHECK IN IF ")
		staticPorts := staticPortsValue.([]interface{})
		// resest the values
		log.Printf("CHECK in IF staticPorts : %v", staticPorts)
		for _, staticPortValue := range staticPorts {
			log.Printf("CHECK IN FOR ")
			staticPort := staticPortValue.(map[string]interface{})
			log.Printf("CHECK staticPort : %v", staticPort)

			staticPortMap := make(map[string]interface{})
			var static_port_pod, static_port_leaf, static_port_path, static_port_fex string

			if staticPort["path_type"] != nil {
				staticPortMap["type"] = staticPort["path_type"]
			} else {
				staticPortMap["type"] = pathType
			}
			if staticPort["deployment_immediacy"] != nil {
				staticPortMap["deploymentImmediacy"] = staticPort["deployment_immediacyode"]
			} else {
				staticPortMap["deploymentImmediacy"] = deploymentImmediacy
			}
			if staticPort["mode"] != nil {
				staticPortMap["mode"] = staticPort["mode"]
			} else {
				staticPortMap["mode"] = mode
			}
			if staticPort["vlan"] != nil {
				staticPortMap["portEncapVlan"] = staticPort["vlan"]
			} else {
				staticPortMap["portEncapVlan"] = vlan
			}
			if staticPort["micro_seg_vlan"] != nil {
				staticPortMap["microSegVlan"] = staticPort["micro_seg_vlan"]
			} else if microsegvlan != 0 {
				staticPortMap["microSegVlan"] = microsegvlan
			}

			if staticPort["pod"] != nil {
				static_port_pod = staticPort["pod"].(string)
			} else {
				static_port_pod = pod
			}
			if staticPort["leaf"] != nil {
				static_port_leaf = staticPort["leaf"].(string)
			} else {
				static_port_leaf = leaf
			}
			if staticPort["path"] != nil {
				static_port_path = staticPort["path"].(string)
			} else {
				static_port_path = path
			}
			if staticPort["fex"] != nil {
				static_port_fex = staticPort["fex"].(string)
			} else {
				static_port_fex = fex
			}

			log.Printf("CHECK staticPort : %v", staticPort)

			if staticPortMap["path_type"] == "port" && fex != "" {
				portpath = fmt.Sprintf("topology/%s/paths-%s/extpaths-%s/pathep-[%s]", static_port_pod, static_port_leaf, static_port_fex, static_port_path)
			} else if staticPortMap["path_type"] == "vpc" {
				portpath = fmt.Sprintf("topology/%s/protpaths-%s/pathep-[%s]", static_port_pod, static_port_leaf, static_port_path)
			} else {
				portpath = fmt.Sprintf("topology/%s/paths-%s/pathep-[%s]", static_port_pod, static_port_leaf, static_port_path)
			}

			staticPortMap["path"] = portpath

			staticPortsList = append(staticPortsList, staticPortMap)
		}
		log.Printf("CHECK staticPortsList : %v", staticPortsList)

	}

	foundEpg := false
	foundAnp := false

	site, err := getSiteFromSiteIdAndTemplate(schemaId, stateSiteId, stateTemplateName, msoClient)
	if err != nil {
		return err
	} else {
		d.Set("site_id", stateSiteId)
		d.Set("template_name", stateTemplateName)
	}
	anpCont, err := getSiteAnp(stateANPName, site)
	if err != nil {
		return err
	} else {
		d.Set("anp_name", stateANPName)
		foundAnp = true
	}

	epgCont, err := getSiteEpg(stateEpgName, anpCont)
	if err != nil {
		return err
	} else {
		d.Set("epg_name", stateEpgName)
		foundEpg = true
	}
	log.Printf("CHECK READ EPGCONT : %v ", epgCont)

	if foundAnp == true && foundEpg == false {
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
	if foundAnp == false && foundEpg == false {
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

	pathsp := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/staticPorts/-", stateSiteId, stateTemplateName, stateANPName, stateEpgName)
	// staticStruct := models.NewSchemaSiteAnpEpgBulkStaticPort("add", pathsp, staticPortsList)
	staticStruct := models.NewSchemaSiteAnpEpgBulkStaticPort("add", pathsp, pathType, portpath, vlan, deploymentImmediacy, microsegvlan, mode)
	_, errs := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), staticStruct)
	if errs != nil {
		return errs
	}
	return resourceMSOSchemaSiteAnpEpgBulkStaticPortRead(d, m)
}

func resourceMSOSchemaSiteAnpEpgBulkStaticPortRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("CHECK READ ")
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	var fex, pathType string
	found := false

	schemaId := d.Get("schema_id").(string)
	stateSite := d.Get("site_id").(string)
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

	site, err := getSiteFromSiteIdAndTemplate(schemaId, stateSite, stateTemplate, msoClient)
	if err != nil {
		return err
	} else {
		d.Set("site_id", stateSite)
		d.Set("template_name", stateTemplate)
	}
	anpCont, err := getSiteAnp(stateAnp, site)
	if err != nil {
		return err
	} else {
		d.Set("anp_name", stateAnp)
	}

	epgCont, err := getSiteEpg(stateEpg, anpCont)
	if err != nil {
		return err
	} else {
		d.Set("epg_name", stateEpg)
	}
	log.Printf("CHECK READ EPGCONT : %v ", epgCont)

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

	if !found {
		d.SetId("")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}

func resourceMSOSchemaSiteAnpEpgBulkStaticPortUpdate(d *schema.ResourceData, m interface{}) error {
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
									anpStruct := models.NewSchemaSiteAnpEpgBulkStaticPort("replace", path, pathType, portpath, vlan, deploymentImmediacy, microsegvlan, mode)
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

	return resourceMSOSchemaSiteAnpEpgBulkStaticPortRead(d, m)
}

func resourceMSOSchemaSiteAnpEpgBulkStaticPortDelete(d *schema.ResourceData, m interface{}) error {
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
									anpStruct := models.NewSchemaSiteAnpEpgBulkStaticPort("remove", path, pathType, portpath, vlan, deploymentImmediacy, microsegvlan, mode)
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
	return resourceMSOSchemaSiteAnpEpgBulkStaticPortRead(d, m)
}

func getSiteFromSiteIdAndTemplate(schemaId, siteId, templateName string, msoClient *client.Client) (*container.Container, error) {
	schemaObject, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return nil, err
	}
	siteCount, err := schemaObject.ArrayCount("sites")
	if err != nil {
		return nil, fmt.Errorf("No Sites found")
	}
	// found := false

	for i := 0; i < siteCount; i++ {
		site, err := schemaObject.ArrayElement(i, "sites")
		if err != nil {
			return nil, err
		}

		if models.G(site, "siteId") == siteId && models.G(site, "templateName") == templateName {
			return site, nil
		}
	}
	return nil, fmt.Errorf("Site-Template association for %v-%v is not found.", siteId, templateName)
}

func getSiteAnp(anpName string, site *container.Container) (*container.Container, error) {

	anpCount, err := site.ArrayCount("anps")
	if err != nil {
		return nil, fmt.Errorf("Unable to get Anp list")
	}
	for i := 0; i < anpCount; i++ {
		anpCont, err := site.ArrayElement(i, "anps")
		if err != nil {
			return nil, err
		}
		anpRef := models.G(anpCont, "anpRef")
		re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/anps/(.*)")
		match := re.FindStringSubmatch(anpRef)
		if match[3] == anpName {
			return anpCont, nil
		}
	}
	return nil, fmt.Errorf("ANP %v is not found in Site.", anpName)
}

func getSiteEpg(epgName string, anpCont *container.Container) (*container.Container, error) {

	epgCount, err := anpCont.ArrayCount("epgs")
	if err != nil {
		return nil, fmt.Errorf("Unable to get EPG list")
	}
	for i := 0; i < epgCount; i++ {
		epgCont, err := anpCont.ArrayElement(i, "epgs")
		if err != nil {
			return nil, err
		}
		epgRef := models.G(epgCont, "epgRef")
		re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/epgs/(.*)")
		match := re.FindStringSubmatch(epgRef)
		if match[3] == epgName {
			return epgCont, nil
		}
	}
	return nil, fmt.Errorf("EPG %v is not found in Site.", epgName)
}
