package mso

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceMSOSchemaSiteAnpEpgStaticleaf() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteAnpEpgStaticleafCreate,
		Read:   resourceMSOSchemaSiteAnpEpgStaticleafRead,
		Delete: resourceMSOSchemaSiteAnpEpgStaticleafDelete,

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
			"path": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"port_encap_vlan": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		}),
	}
}

func resourceMSOSchemaSiteAnpEpgStaticleafCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Anp Epg StaticLeaf: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)
	paths := d.Get("path").(string)
	portEncapVlan := d.Get("port_encap_vlan").(int)

	path := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/staticLeafs/-", siteId, templateName, anpName, epgName)
	anpEpgStaticStruct := models.NewSchemaSiteAnpEpgStaticleaf("add", path, paths, portEncapVlan)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpEpgStaticStruct)
	if err != nil {
		return err
	}
	return resourceMSOSchemaSiteAnpEpgStaticleafRead(d, m)
}

func resourceMSOSchemaSiteAnpEpgStaticleafRead(d *schema.ResourceData, m interface{}) error {
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
	statePath := d.Get("path").(string)

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
							staticLeafCount, err := epgCont.ArrayCount("staticLeafs")
							if err != nil {
								return fmt.Errorf("Unable to get Static Leaf list")
							}
							for s := 0; s < staticLeafCount; s++ {
								staticLeafCont, err := epgCont.ArrayElement(s, "staticLeafs")
								if err != nil {
									return err
								}
								apiPath := models.StripQuotes(staticLeafCont.S("path").String())
								if apiPath == statePath {
									d.SetId(apiPath)
									d.Set("path", apiPath)
									d.Set("site_id", apiSite)
									d.Set("schema_id", split[2])
									d.Set("template_name", split[4])
									d.Set("anp_name", split[6])
									d.Set("epg_name", apiEPG)
									apiPort, _ := strconv.Atoi(staticLeafCont.S("portEncapVlan").String())
									d.Set("port_encap_vlan", apiPort)
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
		d.Set("schema_id", "")
		d.Set("site_id", "")
		d.Set("template_name", "")
		d.Set("epg_name", "")
		d.Set("anp_name", "")
		d.Set("path", "")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}

func resourceMSOSchemaSiteAnpEpgStaticleafDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Anp Epg StaticLeaf: Beginning Deletion")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)
	paths := d.Get("path").(string)
	portEncapVlan := d.Get("port_encap_vlan").(int)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return fmt.Errorf("No Sites found")
	}

	index := -1
	stateSite := d.Get("site_id").(string)
	stateAnp := d.Get("anp_name").(string)
	stateEpg := d.Get("epg_name").(string)
	statePath := d.Get("path").(string)

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
							staticLeafCount, err := epgCont.ArrayCount("staticLeafs")
							if err != nil {
								return fmt.Errorf("Unable to get Static Leaf list")
							}
							for s := 0; s < staticLeafCount; s++ {
								staticLeafCont, err := epgCont.ArrayElement(s, "staticLeafs")
								if err != nil {
									return err
								}
								apiPath := models.StripQuotes(staticLeafCont.S("path").String())
								if apiPath == statePath {
									index = s
									break
								}
							}
						}
					}
				}
			}
		}
	}

	if index == -1 {
		return fmt.Errorf("The given staticLeaf is not found")
	}

	indexs := strconv.Itoa(index)
	path := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/staticLeafs/%s", siteId, templateName, anpName, epgName, indexs)
	anpEpgStaticStruct := models.NewSchemaSiteAnpEpgStaticleaf("remove", path, paths, portEncapVlan)
	_, err1 := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpEpgStaticStruct)

	if err1 != nil {
		return err1
	}
	d.SetId("")
	return nil
}
