package mso

import (
	"fmt"
	"log"
	"regexp"

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
							Optional:     true,
							Computed:     true,
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
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func datasourceMSOSchemaSiteAnpEpgBulkStaticPortRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Beginning Data source")

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	stateSite := d.Get("site_id").(string)
	stateTemplate := d.Get("template_name").(string)
	stateAnp := d.Get("anp_name").(string)
	stateEpg := d.Get("epg_name").(string)
	epgDn := fmt.Sprintf("%s/site/%s/template/%s/anp/%s/epg/%s", schemaId, stateSite, stateTemplate, stateAnp, stateEpg)

	d.SetId(epgDn)
	d.Set("schema_id", schemaId)

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

	portCount, err := epgCont.ArrayCount("staticPorts")
	if err != nil {
		return fmt.Errorf("Unable to get Static Port list")
	}

	log.Printf("CHECK DATA port Count : %v ", portCount)
	staticPortsList := make([]interface{}, 0, 1)
	for i := 0; i < portCount; i++ {
		portCont, err := epgCont.ArrayElement(i, "staticPorts")
		if err != nil {
			return err
		}
		log.Printf("CHECK DATA port portCont : %v ", portCont)

		staticPortMap := make(map[string]interface{})

		if portCont.Exists("type") {
			staticPortMap["type"] = models.StripQuotes(portCont.S("type").String())
		}

		if portCont.Exists("portEncapVlan") {
			staticPortMap["vlan"] = models.StripQuotes(portCont.S("portEncapVlan").String())
		}
		if portCont.Exists("deploymentImmediacy") {
			staticPortMap["deployment_immediacy"] = models.StripQuotes(portCont.S("deploymentImmediacy").String())
		}
		if portCont.Exists("microSegVlan") {

			staticPortMap["micro_seg_vlan"] = models.StripQuotes(portCont.S("microSegVlan").String())
		}
		if portCont.Exists("mode") {
			staticPortMap["mode"] = models.StripQuotes(portCont.S("mode").String())
		}

		pathValue := models.StripQuotes(portCont.S("path").String())

		matchedMap := make(map[string]string)

		if (regexp.MustCompile(`(topology\/(?P<podValue>.*)\/paths-(?P<leafValue>.*)\/extpaths-(?P<fexValue>.*)\/pathep-(?P<pathValue>.*))`)).MatchString(pathValue) {
			matchedMap = getStaticPortPathValues(pathValue, regexp.MustCompile(`(topology\/(?P<podValue>.*)\/paths-(?P<leafValue>.*)\/extpaths-(?P<fexValue>.*)\/pathep-(?P<pathValue>.*))`))
			staticPortMap["fex"] = matchedMap["fexValue"]
		} else if (regexp.MustCompile(`(topology\/(?P<podValue>.*)\/protpaths-(?P<leafValue>.*)\/pathep-(?P<pathValue>.*))`)).MatchString(pathValue) {
			matchedMap = getStaticPortPathValues(pathValue, regexp.MustCompile(`(topology\/(?P<podValue>.*)\/protpaths-(?P<leafValue>.*)\/pathep-(?P<pathValue>.*))`))
		} else if (regexp.MustCompile(`(topology\/(?P<podValue>.*)\/paths-(?P<leafValue>.*)\/pathep-(?P<pathValue>.*))`)).MatchString(pathValue) {
			matchedMap = getStaticPortPathValues(pathValue, (regexp.MustCompile(`(topology\/(?P<podValue>.*)\/paths-(?P<leafValue>.*)\/pathep-(?P<pathValue>.*))`)))
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
