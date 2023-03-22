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

func resourceMSOSchemaSiteAnpEpgBulkStaticPort() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteAnpEpgBulkStaticPortCreate,
		Read:   resourceMSOSchemaSiteAnpEpgBulkStaticPortRead,
		Update: resourceMSOSchemaSiteAnpEpgBulkStaticPortUpdate,
		Delete: resourceMSOSchemaSiteAnpEpgBulkStaticPortDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaSiteAnpEpgBulkStaticPortImport,
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
			"static_ports": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "port",
							ValidateFunc: validation.StringInSlice([]string{
								"port",
								"vpc",
								"dpc",
							}, false),
						},
						"pod": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
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
							Required: true,
							ForceNew: true,
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

func resourceMSOSchemaSiteAnpEpgBulkStaticPortImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())

	msoClient := m.(*client.Client)

	get_attribute := strings.Split(d.Id(), "/")
	schemaId := get_attribute[0]
	stateSite := get_attribute[2]
	stateTemplate := get_attribute[4]
	stateAnp := get_attribute[6]
	stateEpg := get_attribute[8]

	d.SetId(d.Id())
	d.Set("schema_id", schemaId)

	siteCont, err := getSiteFromSiteIdAndTemplate(schemaId, stateSite, stateTemplate, msoClient)
	if err != nil {
		return nil, err
	} else {
		d.Set("site_id", stateSite)
		d.Set("template_name", stateTemplate)
	}

	anpCont, err := getSiteAnp(stateAnp, siteCont)
	if err != nil {
		return nil, err
	} else {
		d.Set("anp_name", stateAnp)
	}

	epgCont, err := getSiteEpg(stateEpg, anpCont)
	if err != nil {
		return nil, err
	} else {
		d.Set("epg_name", stateEpg)
	}

	portCount, err := epgCont.ArrayCount("staticPorts")
	if err != nil {
		return nil, fmt.Errorf("Unable to get Static Port list")
	}

	portPath := regexp.MustCompile(`(topology\/(?P<podValue>.*)\/paths-(?P<leafValue>.*)\/extpaths-(?P<fexValue>.*)\/pathep-\[(?P<pathValue>.*)\])`)
	vpcPath := regexp.MustCompile(`(topology\/(?P<podValue>.*)\/protpaths-(?P<leafValue>.*)\/pathep-\[(?P<pathValue>.*)\])`)
	dpcPath := regexp.MustCompile(`(topology\/(?P<podValue>.*)\/paths-(?P<leafValue>.*)\/pathep-\[(?P<pathValue>.*)\])`)

	staticPortsList := make([]interface{}, 0, 1)
	for i := 0; i < portCount; i++ {
		portCont, err := epgCont.ArrayElement(i, "staticPorts")
		if err != nil {
			return nil, err
		}

		staticPortMap := make(map[string]interface{})

		if portCont.Exists("type") {
			staticPortMap["path_type"] = models.StripQuotes(portCont.S("type").String())
		}
		if portCont.Exists("portEncapVlan") {
			staticPortMap["vlan"], _ = strconv.Atoi(fmt.Sprintf("%v", portCont.S("portEncapVlan")))
		}
		if portCont.Exists("deploymentImmediacy") {
			staticPortMap["deployment_immediacy"] = models.StripQuotes(portCont.S("deploymentImmediacy").String())
		}
		if portCont.Exists("microSegVlan") {
			staticPortMap["micro_seg_vlan"], _ = strconv.Atoi(fmt.Sprintf("%v", portCont.S("microSegVlan")))
		}
		if portCont.Exists("mode") {
			staticPortMap["mode"] = models.StripQuotes(portCont.S("mode").String())
		}

		pathValue := models.StripQuotes(portCont.S("path").String())

		matchedMap := make(map[string]string)
		// put regex in variable and move it outside the for loop.

		if portPath.MatchString(pathValue) {
			matchedMap = getStaticPortPathValues(pathValue, portPath)
			staticPortMap["fex"] = matchedMap["fexValue"]
		} else if vpcPath.MatchString(pathValue) {
			matchedMap = getStaticPortPathValues(pathValue, vpcPath)
		} else if dpcPath.MatchString(pathValue) {
			matchedMap = getStaticPortPathValues(pathValue, dpcPath)
		}

		staticPortMap["pod"] = matchedMap["podValue"]
		staticPortMap["leaf"] = matchedMap["leafValue"]
		staticPortMap["path"] = matchedMap["pathValue"]

		staticPortsList = append(staticPortsList, staticPortMap)
	}
	d.Set("static_ports", staticPortsList)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaSiteAnpEpgBulkStaticPortCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Static Port Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	stateSiteId := d.Get("site_id").(string)
	stateTemplateName := d.Get("template_name").(string)
	stateANPName := d.Get("anp_name").(string)
	stateEpgName := d.Get("epg_name").(string)
	epgDn := fmt.Sprintf("%s/site/%s/template/%s/anp/%s/epg/%s", schemaId, stateSiteId, stateTemplateName, stateANPName, stateEpgName)
	staticPortsList := make([]interface{}, 0, 1)
	if staticPortsValue, ok := d.GetOk("static_ports"); ok {
		staticPorts := staticPortsValue.([]interface{})
		// resest the values
		for _, staticPortValue := range staticPorts {
			staticPort := staticPortValue.(map[string]interface{})
			staticPortMap := make(map[string]interface{})
			var static_port_pod, static_port_leaf, static_port_path, static_port_fex string

			if staticPort["path_type"] != nil {
				staticPortMap["type"] = staticPort["path_type"]
			}
			if staticPort["deployment_immediacy"] != nil {
				staticPortMap["deploymentImmediacy"] = staticPort["deployment_immediacy"]
			}
			if staticPort["mode"] != nil {
				staticPortMap["mode"] = staticPort["mode"]
			}
			if staticPort["vlan"] != nil {
				staticPortMap["portEncapVlan"] = staticPort["vlan"]
			}
			if staticPort["micro_seg_vlan"] != 0 {
				staticPortMap["microSegVlan"] = staticPort["micro_seg_vlan"]
			}
			if staticPort["pod"] != nil {
				static_port_pod = staticPort["pod"].(string)
			}
			if staticPort["leaf"] != nil {
				static_port_leaf = staticPort["leaf"].(string)
			}
			if staticPort["path"] != nil {
				static_port_path = staticPort["path"].(string)
			}
			if staticPort["fex"] != nil {
				static_port_fex = staticPort["fex"].(string)
			}

			var portpath string

			if staticPortMap["path_type"] == "port" && staticPort["fex"] != "" {
				portpath = fmt.Sprintf("topology/%s/paths-%s/extpaths-%s/pathep-[%s]", static_port_pod, static_port_leaf, static_port_fex, static_port_path)
			} else if staticPortMap["path_type"] == "vpc" {
				portpath = fmt.Sprintf("topology/%s/protpaths-%s/pathep-[%s]", static_port_pod, static_port_leaf, static_port_path)
			} else {
				portpath = fmt.Sprintf("topology/%s/paths-%s/pathep-[%s]", static_port_pod, static_port_leaf, static_port_path)
			}

			staticPortMap["path"] = portpath

			staticPortsList = append(staticPortsList, staticPortMap)
		}
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

	_, err = getSiteEpg(stateEpgName, anpCont)
	if err != nil {
		return err
	} else {
		d.Set("epg_name", stateEpgName)
		foundEpg = true
	}

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

	pathsp := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/staticPorts", stateSiteId, stateTemplateName, stateANPName, stateEpgName)
	staticStruct := models.NewSchemaSiteAnpEpgBulkStaticPort("add", pathsp, staticPortsList)
	_, errs := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), staticStruct)
	if errs != nil {
		return errs
	}

	d.SetId(epgDn)
	return resourceMSOSchemaSiteAnpEpgBulkStaticPortRead(d, m)
}

func resourceMSOSchemaSiteAnpEpgBulkStaticPortRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	stateSite := d.Get("site_id").(string)
	stateTemplate := d.Get("template_name").(string)
	stateAnp := d.Get("anp_name").(string)
	stateEpg := d.Get("epg_name").(string)

	d.SetId(d.Id())
	d.Set("schema_id", schemaId)

	siteCont, err := getSiteFromSiteIdAndTemplate(schemaId, stateSite, stateTemplate, msoClient)
	if err != nil {
		return err
	} else {
		d.Set("site_id", stateSite)
		d.Set("template_name", stateTemplate)
	}

	anpCont, err := getSiteAnp(stateAnp, siteCont)
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

	portCount, err := epgCont.ArrayCount("staticPorts")
	if err != nil {
		return fmt.Errorf("Unable to get Static Port list")
	}

	portPath := regexp.MustCompile(`(topology\/(?P<podValue>.*)\/paths-(?P<leafValue>.*)\/extpaths-(?P<fexValue>.*)\/pathep-\[(?P<pathValue>.*)\])`)
	vpcPath := regexp.MustCompile(`(topology\/(?P<podValue>.*)\/protpaths-(?P<leafValue>.*)\/pathep-\[(?P<pathValue>.*)\])`)
	dpcPath := regexp.MustCompile(`(topology\/(?P<podValue>.*)\/paths-(?P<leafValue>.*)\/pathep-\[(?P<pathValue>.*)\])`)

	staticPortsList := make([]interface{}, 0, 1)
	for i := 0; i < portCount; i++ {
		portCont, err := epgCont.ArrayElement(i, "staticPorts")
		if err != nil {
			return err
		}

		staticPortMap := make(map[string]interface{})

		if portCont.Exists("type") {
			staticPortMap["path_type"] = models.StripQuotes(portCont.S("type").String())
		}
		if portCont.Exists("portEncapVlan") {
			staticPortMap["vlan"], _ = strconv.Atoi(fmt.Sprintf("%v", portCont.S("portEncapVlan")))
		}
		if portCont.Exists("deploymentImmediacy") {
			staticPortMap["deployment_immediacy"] = models.StripQuotes(portCont.S("deploymentImmediacy").String())
		}
		if portCont.Exists("microSegVlan") {
			staticPortMap["micro_seg_vlan"], _ = strconv.Atoi(fmt.Sprintf("%v", portCont.S("microSegVlan")))
		}
		if portCont.Exists("mode") {
			staticPortMap["mode"] = models.StripQuotes(portCont.S("mode").String())
		}

		pathValue := models.StripQuotes(portCont.S("path").String())

		matchedMap := make(map[string]string)

		if portPath.MatchString(pathValue) {
			matchedMap = getStaticPortPathValues(pathValue, portPath)
			staticPortMap["fex"] = matchedMap["fexValue"]
		} else if vpcPath.MatchString(pathValue) {
			matchedMap = getStaticPortPathValues(pathValue, vpcPath)
		} else if dpcPath.MatchString(pathValue) {
			matchedMap = getStaticPortPathValues(pathValue, dpcPath)
		}

		staticPortMap["pod"] = matchedMap["podValue"]
		staticPortMap["leaf"] = matchedMap["leafValue"]
		staticPortMap["path"] = matchedMap["pathValue"]

		staticPortsList = append(staticPortsList, staticPortMap)
	}
	d.Set("static_ports", staticPortsList)

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

	staticPortsList := make([]interface{}, 0, 1)
	if staticPortsValue, ok := d.GetOk("static_ports"); ok {
		staticPorts := staticPortsValue.([]interface{})
		for _, staticPortValue := range staticPorts {
			staticPort := staticPortValue.(map[string]interface{})
			staticPortMap := make(map[string]interface{})
			var static_port_pod, static_port_leaf, static_port_path, static_port_fex string

			if staticPort["path_type"] != nil {
				staticPortMap["type"] = staticPort["path_type"]
			}
			if staticPort["deployment_immediacy"] != nil {
				staticPortMap["deploymentImmediacy"] = staticPort["deployment_immediacy"]
			}
			if staticPort["mode"] != nil {
				staticPortMap["mode"] = staticPort["mode"]
			}
			if staticPort["vlan"] != nil {
				staticPortMap["portEncapVlan"] = staticPort["vlan"]
			}
			if staticPort["micro_seg_vlan"] != 0 {
				staticPortMap["microSegVlan"] = staticPort["micro_seg_vlan"]
			}
			if staticPort["pod"] != nil {
				static_port_pod = staticPort["pod"].(string)
			}
			if staticPort["leaf"] != nil {
				static_port_leaf = staticPort["leaf"].(string)
			}
			if staticPort["path"] != nil {
				static_port_path = staticPort["path"].(string)
			}
			if staticPort["fex"] != nil {
				static_port_fex = staticPort["fex"].(string)
			}

			var portpath string

			if staticPortMap["path_type"] == "port" && staticPort["fex"] != "" {
				portpath = fmt.Sprintf("topology/%s/paths-%s/extpaths-%s/pathep-[%s]", static_port_pod, static_port_leaf, static_port_fex, static_port_path)
			} else if staticPortMap["path_type"] == "vpc" {
				portpath = fmt.Sprintf("topology/%s/protpaths-%s/pathep-[%s]", static_port_pod, static_port_leaf, static_port_path)
			} else {
				portpath = fmt.Sprintf("topology/%s/paths-%s/pathep-[%s]", static_port_pod, static_port_leaf, static_port_path)
			}

			staticPortMap["path"] = portpath

			staticPortsList = append(staticPortsList, staticPortMap)
		}
	}

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
	}

	_, err = getSiteEpg(stateEpgName, anpCont)
	if err != nil {
		return err
	} else {
		d.Set("epg_name", stateEpgName)
	}

	pathsp := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/staticPorts", stateSiteId, stateTemplateName, stateANPName, stateEpgName)
	staticStruct := models.NewSchemaSiteAnpEpgBulkStaticPort("replace", pathsp, staticPortsList)
	_, errs := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), staticStruct)
	if errs != nil {
		return errs
	}

	return resourceMSOSchemaSiteAnpEpgBulkStaticPortRead(d, m)
}

func resourceMSOSchemaSiteAnpEpgBulkStaticPortDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Anp: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	stateSiteId := d.Get("site_id").(string)
	stateTemplateName := d.Get("template_name").(string)
	stateANPName := d.Get("anp_name").(string)
	stateEpgName := d.Get("epg_name").(string)

	// Make the above  values to map and update the values in static ports. copy the parent map to local staticport map(to avoid overriding thge values for seconfd static port)

	staticPortsList := make([]interface{}, 0, 1)
	if staticPortsValue, ok := d.GetOk("static_ports"); ok {
		staticPorts := staticPortsValue.([]interface{})
		// resest the values
		for _, staticPortValue := range staticPorts {
			staticPort := staticPortValue.(map[string]interface{})
			staticPortMap := make(map[string]interface{})
			var static_port_pod, static_port_leaf, static_port_path, static_port_fex string

			if staticPort["path_type"] != nil {
				staticPortMap["type"] = staticPort["path_type"]
			}
			if staticPort["deployment_immediacy"] != nil {
				staticPortMap["deploymentImmediacy"] = staticPort["deployment_immediacy"]
			}
			if staticPort["mode"] != nil {
				staticPortMap["mode"] = staticPort["mode"]
			}
			if staticPort["vlan"] != nil {
				staticPortMap["portEncapVlan"] = staticPort["vlan"]
			}
			if staticPort["micro_seg_vlan"] != 0 {
				staticPortMap["microSegVlan"] = staticPort["micro_seg_vlan"]
			}
			if staticPort["pod"] != nil {
				static_port_pod = staticPort["pod"].(string)
			}
			if staticPort["leaf"] != nil {
				static_port_leaf = staticPort["leaf"].(string)
			}
			if staticPort["path"] != nil {
				static_port_path = staticPort["path"].(string)
			}
			if staticPort["fex"] != nil {
				static_port_fex = staticPort["fex"].(string)
			}

			var portpath string

			if staticPortMap["path_type"] == "port" && staticPort["fex"] != "" {
				portpath = fmt.Sprintf("topology/%s/paths-%s/extpaths-%s/pathep-[%s]", static_port_pod, static_port_leaf, static_port_fex, static_port_path)
			} else if staticPortMap["path_type"] == "vpc" {
				portpath = fmt.Sprintf("topology/%s/protpaths-%s/pathep-[%s]", static_port_pod, static_port_leaf, static_port_path)
			} else {
				portpath = fmt.Sprintf("topology/%s/paths-%s/pathep-[%s]", static_port_pod, static_port_leaf, static_port_path)
			}

			staticPortMap["path"] = portpath

			staticPortsList = append(staticPortsList, staticPortMap)
		}
	}

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
	}

	_, err = getSiteEpg(stateEpgName, anpCont)
	if err != nil {
		return err
	} else {
		d.Set("epg_name", stateEpgName)
	}

	pathsp := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/staticPorts", stateSiteId, stateTemplateName, stateANPName, stateEpgName)
	staticStruct := models.NewSchemaSiteAnpEpgBulkStaticPort("remove", pathsp, staticPortsList)
	response, errs := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), staticStruct)

	if errs != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return err
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

func getStaticPortPathValues(pathValue string, re *regexp.Regexp) map[string]string {
	match := re.FindStringSubmatch(pathValue) //list of matched strings
	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	return result
}
