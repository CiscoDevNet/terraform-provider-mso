package mso

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceMSOSchemaSiteAnpEpgStaticPort() *schema.Resource {
	return &schema.Resource{

		Read: datasourceMSOSchemaSiteAnpEpgStaticPortRead,

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
			"path_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"port",
					"vpc",
					"dpc",
				}, false),
			},
			"pod": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"leaf": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"path": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"fex": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"vlan": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"deployment_immediacy": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"micro_seg_vlan": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}
}

func datasourceMSOSchemaSiteAnpEpgStaticPortRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	anp := d.Get("anp_name").(string)
	epg := d.Get("epg_name").(string)
	pod := d.Get("pod").(string)
	leaf := d.Get("leaf").(string)
	path := d.Get("path").(string)
	pathType := d.Get("path_type").(string)
	var fex string
	if tempVar, ok := d.GetOk("fex"); ok {
		fex = tempVar.(string)
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

	portCount, err := epgCont.ArrayCount("staticPorts")
	if err != nil {
		return fmt.Errorf("Unable to get Static Port list")
	}

	found := false
	var portPath string
	for l := 0; l < portCount; l++ {
		portCont, err := epgCont.ArrayElement(l, "staticPorts")
		if err != nil {
			return err
		}
		portPath = createPortPath(pathType, pod, leaf, fex, path)
		if portPath == models.StripQuotes(portCont.S("path").String()) && pathType == models.StripQuotes(portCont.S("type").String()) {
			found = true
			d.SetId(fmt.Sprintf("%s/sites/%s-%s/anps/%s/epgs/%s/staticPorts/%s", schemaId, siteId, templateName, anp, epg, portPath))
			if portCont.Exists("type") {
				d.Set("type", models.StripQuotes(portCont.S("type").String()))
			}
			if portCont.Exists("path") {
				d.Set("pod", pod)
				d.Set("leaf", leaf)
				d.Set("path", path)
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
			break
		}
	}

	if !found {
		d.SetId("")
		return fmt.Errorf("Unable to find static port entry: %s", portPath)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
