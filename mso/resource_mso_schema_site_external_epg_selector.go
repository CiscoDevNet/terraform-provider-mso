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

	cont, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), schemaSiteExternalEpgSelector)
	if err != nil {
		return err
	}

	index, _, err := checkselector(cont, schemaID, siteID, templateName, externalEpgName, name)
	if err != nil {
		return err
	}

	if index == -1 {
		d.SetId("")
	} else {
		d.SetId(name)
	}
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

	cont, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), schemaSiteExternalEpgSelector)
	if err != nil {
		return err
	}

	index, _, err := checkselector(cont, schemaID, siteID, templateName, externalEpgName, name)
	if err != nil {
		return err
	}

	if index == -1 {
		d.SetId("")
	} else {
		d.SetId(name)
	}
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

	_, err1 := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), schemaSiteExternalEpgSelector)
	if err1 != nil {
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
