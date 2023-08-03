package mso

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceMSOSchemaSiteAnpEpgBulkStaticPort() *schema.Resource {
	return &schema.Resource{

		Read: datasourceMSOSchemaSiteAnpEpgBulkStaticPortRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"site_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_name": &schema.Schema{
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
			"static_ports": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pod": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"leaf": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vlan": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"deployment_immediacy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fex": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"micro_seg_vlan": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		}),
	}
}

func datasourceMSOSchemaSiteAnpEpgBulkStaticPortRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Beginning Data source")

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	anp := d.Get("anp_name").(string)
	epg := d.Get("epg_name").(string)
	epgDn := fmt.Sprintf("%s/sites/%s-%s/anps/%s/epgs/%s", schemaId, siteId, templateName, anp, epg)

	d.SetId(epgDn)
	d.Set("schema_id", schemaId)

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

	log.Printf("[DEBUG] %s: Data source finished successfully", epgDn)
	return nil
}
