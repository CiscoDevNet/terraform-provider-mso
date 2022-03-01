package mso

import (
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOSchemaSiteExternalEpgSelector() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteExternalEpgSelectorCreate,
		Update: resourceMSOSchemaSiteExternalEpgSelectorUpdate,
		Read:   resourceMSOSchemaSiteExternalEpgSelectorRead,
		Delete: resourceMSOSchemaSiteExternalEpgSelectorDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaSiteExternalEpgSelectorImport,
		},

		Schema: map[string]*schema.Schema{
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

			"external_epg_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"ip": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		},
	}
}

func resourceMSOSchemaSiteExternalEpgSelectorImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	msoClient := m.(*client.Client)
	get_attribute := strings.Split(d.Id(), "/")
	dn := get_attribute[8]
	schemaID := get_attribute[0]
	siteID := get_attribute[2]
	templateName := get_attribute[4]
	externalEpgName := get_attribute[6]

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return nil, err
	}

	found := false

	count, err := cont.ArrayCount("sites")
	if err != nil {
		return nil, err
	}

	for i := 0; i < count; i++ {
		siteCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return nil, err
		}

		currSite := models.StripQuotes(siteCont.S("siteId").String())
		currTemplate := models.StripQuotes(siteCont.S("templateName").String())

		if currSite == siteID && currTemplate == templateName {
			extEpgCount, err := siteCont.ArrayCount("externalEpgs")
			if err != nil {
				return nil, err
			}

			for j := 0; j < extEpgCount; j++ {
				extEpgCont, err := siteCont.ArrayElement(j, "externalEpgs")
				if err != nil {
					return nil, err
				}

				extEpgRef := models.StripQuotes(extEpgCont.S("externalEpgRef").String())
				tokens := strings.Split(extEpgRef, "/")
				extEpgName := tokens[len(tokens)-1]
				if extEpgName == externalEpgName {
					subnetCount, err := extEpgCont.ArrayCount("subnets")
					if err != nil {
						return nil, err
					}

					for k := 0; k < subnetCount; k++ {
						subnetCont, err := extEpgCont.ArrayElement(k, "subnets")
						if err != nil {
							return nil, err
						}

						subnetName := models.StripQuotes(subnetCont.S("name").String())
						if subnetName == dn {
							found = true
							d.SetId(dn)
							d.Set("name", subnetName)
							d.Set("ip", models.StripQuotes(subnetCont.S("ip").String()))
							break
						}
					}
				}
				if found {
					d.Set("external_epg_name", extEpgName)
					break
				}
			}
		}
		if found {
			d.Set("site_id", siteID)
			d.Set("template_name", templateName)
			break
		}
	}

	if found {
		d.Set("schema_id", schemaID)
	} else {
		d.SetId("")
		return nil, fmt.Errorf("Selector of specified name not found")
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaSiteExternalEpgSelectorCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema Site External EPG Selector: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	siteID := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	externalEpgName := d.Get("external_epg_name").(string)
	name := d.Get("name").(string)
	ip := d.Get("ip").(string)

	selectorMap := make(map[string]interface{})
	selectorMap["name"] = name
	selectorMap["ip"] = ip

	contGet, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return err
	}

	check, err := checkEpg(contGet, schemaID, siteID, templateName, externalEpgName)
	if err != nil {
		return err
	}
	if !check {
		return fmt.Errorf("No site External EPG available of given name")
	}

	path := fmt.Sprintf("/sites/%s-%s/externalEpgs/%s/subnets/-", siteID, templateName, externalEpgName)

	schemaSiteExternalEpgSelector := models.NewSchemaSiteExternalEpgSelector("add", path, selectorMap)

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), schemaSiteExternalEpgSelector)
	if err != nil {
		return err
	}

	d.SetId(name)
	log.Printf("[DEBUG] Schema Site External EPG Selector: Creation Completed")
	return resourceMSOSchemaSiteExternalEpgSelectorRead(d, m)
}

func resourceMSOSchemaSiteExternalEpgSelectorUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Update", d.Id())
	msoClient := m.(*client.Client)

	dn := d.Id()
	schemaID := d.Get("schema_id").(string)
	siteID := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	externalEpgName := d.Get("external_epg_name").(string)
	name := d.Get("name").(string)
	ip := d.Get("ip").(string)

	selectorMap := make(map[string]interface{})
	selectorMap["name"] = name
	selectorMap["ip"] = ip

	contGet, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return err
	}

	indexGet, _, err := checkselector(contGet, schemaID, siteID, templateName, externalEpgName, dn)
	if err != nil {
		return nil
	}
	if indexGet == -1 {
		return fmt.Errorf("No Selectors found")
	}

	path := fmt.Sprintf("/sites/%s-%s/externalEpgs/%s/subnets/%v", siteID, templateName, externalEpgName, indexGet)

	schemaSiteExternalEpgSelector := models.NewSchemaSiteExternalEpgSelector("replace", path, selectorMap)

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), schemaSiteExternalEpgSelector)
	if err != nil {
		return err
	}

	d.SetId(name)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceMSOSchemaSiteExternalEpgSelectorRead(d, m)
}

func resourceMSOSchemaSiteExternalEpgSelectorRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	msoClient := m.(*client.Client)

	dn := d.Id()
	schemaID := d.Get("schema_id").(string)
	siteID := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	externalEpgName := d.Get("external_epg_name").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return err
	}

	found := false

	count, err := cont.ArrayCount("sites")
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		siteCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}

		currSite := models.StripQuotes(siteCont.S("siteId").String())
		currTemplate := models.StripQuotes(siteCont.S("templateName").String())

		if currSite == siteID && currTemplate == templateName {
			extEpgCount, err := siteCont.ArrayCount("externalEpgs")
			if err != nil {
				return err
			}

			for j := 0; j < extEpgCount; j++ {
				extEpgCont, err := siteCont.ArrayElement(j, "externalEpgs")
				if err != nil {
					return err
				}

				extEpgRef := models.StripQuotes(extEpgCont.S("externalEpgRef").String())
				tokens := strings.Split(extEpgRef, "/")
				extEpgName := tokens[len(tokens)-1]
				if extEpgName == externalEpgName {
					subnetCount, err := extEpgCont.ArrayCount("subnets")
					if err != nil {
						return err
					}

					for k := 0; k < subnetCount; k++ {
						subnetCont, err := extEpgCont.ArrayElement(k, "subnets")
						if err != nil {
							return err
						}

						subnetName := models.StripQuotes(subnetCont.S("name").String())
						if subnetName == dn {
							found = true
							d.SetId(dn)
							d.Set("name", subnetName)
							d.Set("ip", models.StripQuotes(subnetCont.S("ip").String()))
							break
						}
					}
				}
				if found {
					d.Set("external_epg_name", extEpgName)
					break
				}
			}
		}
		if found {
			d.Set("site_id", siteID)
			d.Set("template_name", templateName)
			break
		}
	}

	if found {
		d.Set("schema_id", schemaID)
	} else {
		d.SetId("")
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSOSchemaSiteExternalEpgSelectorDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Delete", d.Id())

	msoClient := m.(*client.Client)

	dn := d.Id()
	schemaID := d.Get("schema_id").(string)
	siteID := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	extrEpgName := d.Get("external_epg_name").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return err
	}

	index, _, err := checkselector(cont, schemaID, siteID, templateName, extrEpgName, dn)
	if err != nil {
		return err
	}
	if index == -1 {
		d.SetId("")
		return nil
	}

	path := fmt.Sprintf("/sites/%s-%s/externalEpgs/%s/subnets/%v", siteID, templateName, extrEpgName, index)

	schemaSiteExternalEpgSelector := models.NewSchemaSiteExternalEpgSelector("remove", path, nil)

	response, err1 := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), schemaSiteExternalEpgSelector)

	// Ignoring Error with code 141: Resource Not Found when deleting
	if err1 != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return err1
	}

	d.SetId("")
	log.Printf("[DEBUG] %s: Delete finished successfully", d.Id())
	return nil
}

func checkselector(cont *container.Container, schema, site, template, epgName, name string) (int, int, error) {
	found := false
	index := -1
	subnetCounter := 0

	count, err := cont.ArrayCount("sites")
	if err != nil {
		return index, subnetCounter, err
	}

	for i := 0; i < count; i++ {
		siteCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return index, subnetCounter, err
		}

		currSite := models.StripQuotes(siteCont.S("siteId").String())
		currTemplate := models.StripQuotes(siteCont.S("templateName").String())

		if currSite == site && currTemplate == template {
			extEpgCount, err := siteCont.ArrayCount("externalEpgs")
			if err != nil {
				return index, subnetCounter, err
			}

			for j := 0; j < extEpgCount; j++ {
				extEpgCont, err := siteCont.ArrayElement(j, "externalEpgs")
				if err != nil {
					return index, subnetCounter, err
				}

				extEpgRef := models.StripQuotes(extEpgCont.S("externalEpgRef").String())
				tokens := strings.Split(extEpgRef, "/")
				extEpgName := tokens[len(tokens)-1]
				if extEpgName == epgName {
					subnetCount, err := extEpgCont.ArrayCount("subnets")
					if err != nil {
						return index, subnetCounter, err
					}
					subnetCounter = subnetCount

					for k := 0; k < subnetCount; k++ {
						subnetCont, err := extEpgCont.ArrayElement(k, "subnets")
						if err != nil {
							return index, subnetCounter, err
						}

						subnetName := models.StripQuotes(subnetCont.S("name").String())
						if subnetName == name {
							found = true
							index = k
							break
						}
					}
				}
				if found {
					break
				}
			}
		}
		if found {
			break
		}
	}
	return index, subnetCounter, nil
}

func checkEpg(cont *container.Container, schema, site, template, epg string) (bool, error) {
	found := false

	count, err := cont.ArrayCount("sites")
	if err != nil {
		return found, err
	}

	for i := 0; i < count; i++ {
		siteCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return found, err
		}

		currSite := models.StripQuotes(siteCont.S("siteId").String())
		currTemplate := models.StripQuotes(siteCont.S("templateName").String())

		if currSite == site && currTemplate == template {
			extEpgCount, err := siteCont.ArrayCount("externalEpgs")
			if err != nil {
				return found, err
			}

			for j := 0; j < extEpgCount; j++ {
				extEpgCont, err := siteCont.ArrayElement(j, "externalEpgs")
				if err != nil {
					return found, err
				}

				extEpgRef := models.StripQuotes(extEpgCont.S("externalEpgRef").String())
				tokens := strings.Split(extEpgRef, "/")
				extEpgName := tokens[len(tokens)-1]
				if extEpgName == epg {
					found = true
					break
				}
			}
		}
		if found {
			break
		}
	}
	return found, nil
}
