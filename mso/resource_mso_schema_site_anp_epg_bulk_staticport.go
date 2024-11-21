package mso

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
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
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"leaf": {
							Type:     schema.TypeString,
							Optional: true,
							// Remove computed because when a user updates the list and causes index shifts
							//  the leaf state value will be used at the location of the list index when not provided in config.
							// Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"path": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"vlan": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"deployment_immediacy": {
							Type:     schema.TypeString,
							Optional: true,
							// Remove computed because when a user updates the list and causes index shifts
							//  the deployment_immediacy state value will be used at the location of the list index when not provided in config.
							// Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{
								"immediate",
								"lazy",
							}, false),
						},
						"fex": {
							Type:     schema.TypeString,
							Optional: true,
							// Remove computed because when a user updates the list and causes index shifts
							//  the fex state value will be used at the location of the list index when not provided in config.
							// Computed:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"micro_seg_vlan": {
							Type:     schema.TypeInt,
							Optional: true,
							// Remove computed because when a user updates the list and causes index shifts
							//  the micro_seg_vlan state value will be used at the location of the list index when not provided in config.
							// Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Optional: true,
							// Remove computed because when a user updates the list and causes index shifts
							//  the mode state value will be used at the location of the list index when not provided in config.
							// Computed:     true,
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
	siteId := get_attribute[2]
	templateName := get_attribute[4]
	anp := get_attribute[6]
	epg := get_attribute[8]

	d.Set("schema_id", schemaId)

	siteCont, err := getSiteFromSiteIdAndTemplate(schemaId, siteId, templateName, msoClient)
	if err != nil {
		return nil, err
	} else {
		d.Set("site_id", siteId)
		d.Set("template_name", templateName)
	}

	anpCont, err := getSiteAnp(anp, siteCont)
	if err != nil {
		return nil, err
	} else {
		d.Set("anp_name", anp)
	}

	epgCont, err := getSiteEpg(epg, anpCont)
	if err != nil {
		return nil, err
	} else {
		d.Set("epg_name", epg)
	}

	portCount, err := epgCont.ArrayCount("staticPorts")
	if err != nil {
		return nil, fmt.Errorf("Unable to get Static Port list")
	}

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

		setValuesFromPortPath(staticPortMap, models.StripQuotes(portCont.S("path").String()))

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
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	anp := d.Get("anp_name").(string)
	epg := d.Get("epg_name").(string)
	epgDn := fmt.Sprintf("%s/site/%s/template/%s/anp/%s/epg/%s", schemaId, siteId, templateName, anp, epg)
	staticPortsList := make([]interface{}, 0, 1)
	if staticPortsValue, ok := d.GetOk("static_ports"); ok {
		staticPorts := staticPortsValue.([]interface{})
		for _, staticPortValue := range staticPorts {
			staticPort := staticPortValue.(map[string]interface{})
			staticPortMap := make(map[string]interface{})
			var static_port_pod, static_port_leaf, static_port_path, static_port_fex string

			if staticPort["path_type"] != nil {
				staticPortMap["type"] = staticPort["path_type"].(string)
			}
			if staticPort["deployment_immediacy"] != nil {
				staticPortMap["deploymentImmediacy"] = staticPort["deployment_immediacy"].(string)
			}
			if staticPort["mode"] != nil {
				staticPortMap["mode"] = staticPort["mode"].(string)
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

			staticPortMap["path"] = createPortPath(staticPortMap["type"].(string), static_port_pod, static_port_leaf, static_port_fex, static_port_path)

			staticPortsList = append(staticPortsList, staticPortMap)
		}
	}

	foundEpg := false
	foundAnp := false

	site, err := getSiteFromSiteIdAndTemplate(schemaId, siteId, templateName, msoClient)
	if err != nil {
		return err
	} else {
		d.Set("site_id", siteId)
		d.Set("template_name", templateName)
	}

	anpCont, err := getSiteAnp(anp, site)
	if err != nil {
		return err
	} else {
		d.Set("anp_name", anp)
		foundAnp = true
	}

	_, err = getSiteEpg(epg, anpCont)
	if err != nil {
		return err
	} else {
		d.Set("epg_name", epg)
		foundEpg = true
	}

	if !foundAnp {
		log.Printf("[DEBUG] Site Anp: Beginning Creation")

		anpRefMap := make(map[string]interface{})
		anpRefMap["schemaId"] = schemaId
		anpRefMap["templateName"] = templateName
		anpRefMap["anpName"] = anp

		pathAnp := fmt.Sprintf("/sites/%s-%s/anps/-", siteId, templateName)
		anpStruct := models.NewSchemaSiteAnp("add", pathAnp, anpRefMap)

		_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpStruct)
		if err != nil {
			return err
		}
	}

	if !foundEpg {
		log.Printf("[DEBUG] Site Anp Epg: Beginning Creation")
		anpEpgRefMap := make(map[string]interface{})
		anpEpgRefMap["schemaId"] = schemaId
		anpEpgRefMap["templateName"] = templateName
		anpEpgRefMap["anpName"] = anp
		anpEpgRefMap["epgName"] = epg

		pathEpg := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/-", siteId, templateName, anp)
		//private_link_label argument used in resource site_anp_epg is set to nil here
		anpEpgStruct := models.NewSchemaSiteAnpEpg("add", pathEpg, nil, anpEpgRefMap)

		_, ers := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpEpgStruct)
		if ers != nil {
			return ers
		}
	}

	pathsp := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/staticPorts", siteId, templateName, anp, epg)
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
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	anp := d.Get("anp_name").(string)
	epg := d.Get("epg_name").(string)

	d.SetId(d.Id())
	d.Set("schema_id", schemaId)

	siteCont, err := getSiteFromSiteIdAndTemplate(schemaId, siteId, templateName, msoClient)
	if err != nil {
		return err
	} else {
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

	portCount, err := epgCont.ArrayCount("staticPorts")
	if err != nil {
		return fmt.Errorf("Unable to get Static Port list")
	}

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

		setValuesFromPortPath(staticPortMap, models.StripQuotes(portCont.S("path").String()))

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
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	anp := d.Get("anp_name").(string)
	epg := d.Get("epg_name").(string)

	staticPortsList := make([]interface{}, 0, 1)
	if staticPortsValue, ok := d.GetOk("static_ports"); ok {
		staticPorts := staticPortsValue.([]interface{})
		for _, staticPortValue := range staticPorts {
			staticPort := staticPortValue.(map[string]interface{})
			staticPortMap := make(map[string]interface{})
			var static_port_pod, static_port_leaf, static_port_path, static_port_fex string

			if staticPort["path_type"] != nil {
				staticPortMap["type"] = staticPort["path_type"].(string)
			}
			if staticPort["deployment_immediacy"] != nil {
				staticPortMap["deploymentImmediacy"] = staticPort["deployment_immediacy"].(string)
			}
			if staticPort["mode"] != nil {
				staticPortMap["mode"] = staticPort["mode"].(string)
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

			staticPortMap["path"] = createPortPath(staticPortMap["type"].(string), static_port_pod, static_port_leaf, static_port_fex, static_port_path)

			staticPortsList = append(staticPortsList, staticPortMap)
		}
	}

	site, err := getSiteFromSiteIdAndTemplate(schemaId, siteId, templateName, msoClient)
	if err != nil {
		return err
	} else {
		d.Set("site_id", siteId)
		d.Set("template_name", templateName)
	}

	anpCont, err := getSiteAnp(anp, site)
	if err != nil {
		return err
	} else {
		d.Set("anp_name", anp)
	}

	_, err = getSiteEpg(epg, anpCont)
	if err != nil {
		return err
	} else {
		d.Set("epg_name", epg)
	}

	pathsp := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/staticPorts", siteId, templateName, anp, epg)
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
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	anp := d.Get("anp_name").(string)
	epg := d.Get("epg_name").(string)

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

			staticPortMap["path"] = createPortPath(staticPortMap["type"].(string), static_port_pod, static_port_leaf, static_port_fex, static_port_path)

			staticPortsList = append(staticPortsList, staticPortMap)
		}
	}

	site, err := getSiteFromSiteIdAndTemplate(schemaId, siteId, templateName, msoClient)
	if err != nil {
		return err
	} else {
		d.Set("site_id", siteId)
		d.Set("template_name", templateName)
	}

	anpCont, err := getSiteAnp(anp, site)
	if err != nil {
		return err
	} else {
		d.Set("anp_name", anp)
	}

	_, err = getSiteEpg(epg, anpCont)
	if err != nil {
		return err
	} else {
		d.Set("epg_name", epg)
	}

	pathsp := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/staticPorts", siteId, templateName, anp, epg)
	staticStruct := models.NewSchemaSiteAnpEpgBulkStaticPort("remove", pathsp, staticPortsList)
	response, errs := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), staticStruct)

	if errs != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return err
	}

	d.SetId("")
	return resourceMSOSchemaSiteAnpEpgBulkStaticPortRead(d, m)
}
