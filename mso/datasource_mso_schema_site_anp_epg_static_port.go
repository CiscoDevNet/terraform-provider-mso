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

func datasourceMSOSchemaSiteAnpEpgStaticPort() *schema.Resource {
	return &schema.Resource{

		Read: datasourceMSOSchemaSiteAnpEpgStaticPortRead,

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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
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
				Optional: true,
				Computed: true,
			},
			"deployment_immediacy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"micro_seg_vlan": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fex": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		}),
	}
}

func datasourceMSOSchemaSiteAnpEpgStaticPortRead(d *schema.ResourceData, m interface{}) error {
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
		return fmt.Errorf("Unable to find the static port entry")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
