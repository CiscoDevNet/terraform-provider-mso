package mso

import (
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOSchemaSiteAnpEpg() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteAnpEpgCreate,
		Read:   resourceMSOSchemaSiteAnpEpgRead,
		Delete: resourceMSOSchemaSiteAnpEpgDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaSiteAnpEpgImport,
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
		}),
	}
}

func resourceMSOSchemaSiteAnpEpgImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())

	msoClient := m.(*client.Client)
	get_attribute := strings.Split(d.Id(), "/")

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", get_attribute[0]))
	if err != nil {
		return nil, err
	}
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return nil, fmt.Errorf("No Sites found")
	}

	stateSite := get_attribute[2]
	found := false
	stateAnp := get_attribute[4]
	stateEpg := get_attribute[6]

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return nil, err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())

		if apiSite == stateSite {
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
							d.SetId(apiEPG)
							d.Set("site_id", apiSite)
							d.Set("schema_id", split[2])
							d.Set("template_name", split[4])
							d.Set("anp_name", split[6])
							d.Set("epg_name", apiEPG)

							found = true
							break
						}
					}
				}
			}
		}
	}

	if !found {
		return nil, fmt.Errorf("Unable to find the Site Anp Epg %s", stateEpg)
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil

}

func resourceMSOSchemaSiteAnpEpgCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Anp Epg: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)

	anpEpgRefMap := make(map[string]interface{})
	anpEpgRefMap["schemaId"] = schemaId
	anpEpgRefMap["templateName"] = templateName
	anpEpgRefMap["anpName"] = anpName
	anpEpgRefMap["epgName"] = epgName

	path := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/-", siteId, templateName, anpName)
	anpEpgStruct := models.NewSchemaSiteAnpEpg("add", path, anpEpgRefMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpEpgStruct)
	if err != nil {
		return err
	}
	return resourceMSOSchemaSiteAnpEpgRead(d, m)
}

func resourceMSOSchemaSiteAnpEpgRead(d *schema.ResourceData, m interface{}) error {
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
							d.SetId(apiEPG)
							d.Set("site_id", apiSite)
							d.Set("schema_id", split[2])
							d.Set("template_name", split[4])
							d.Set("anp_name", split[6])
							d.Set("epg_name", apiEPG)
							found = true
							break
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
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}

func resourceMSOSchemaSiteAnpEpgDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Anp Epg: Beginning Update")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)

	anpEpgRefMap := make(map[string]interface{})
	anpEpgRefMap["schemaId"] = schemaId
	anpEpgRefMap["templateName"] = templateName
	anpEpgRefMap["anpName"] = anpName
	anpEpgRefMap["epgName"] = epgName

	path := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s", siteId, templateName, anpName, epgName)
	anpEpgStruct := models.NewSchemaSiteAnpEpg("remove", path, anpEpgRefMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpEpgStruct)

	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
