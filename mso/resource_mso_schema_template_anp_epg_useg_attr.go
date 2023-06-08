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

func resourceMSOSchemaTemplateAnpEpgUsegAttr() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaTemplateAnpEpgUsegAttrCreate,
		Update: resourceMSOSchemaTemplateAnpEpgUsegAttrUpdate,
		Read:   resourceMSOSchemaTemplateAnpEpgUsegAttrRead,
		Delete: resourceMSOSchemaTemplateAnpEpgUsegAttrDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaTemplateAnpEpgUsegAttrImport,
		},

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

			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"useg_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ip",
					"mac",
					"dns",
					"vm-name",      // Vm Name
					"rootContName", // VM data center
					"hv",           // Hypervisor
					"guest-os",     // Operating System
					"tag",
					"vm",     // Identifier
					"domain", // VMM domain
					"vnic",   // Vnic DN
				}, false),
			},

			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"operator": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"equals",
					"startsWith",
					"endsWith",
					"contains",
				}, false),
			},

			"category": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"value": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"useg_subnet": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func resourceMSOSchemaTemplateAnpEpgUsegAttrImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	msoClient := m.(*client.Client)
	get_attribute := strings.Split(d.Id(), "/")
	schemaId := get_attribute[0]
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return nil, err
	}
	d.Set("schema_id", schemaId)
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return nil, fmt.Errorf("No Template found")
	}

	templateName := get_attribute[2]
	anpName := get_attribute[4]
	epgName := get_attribute[6]
	name := get_attribute[8]
	found := false

	for i := 0; i < count; i++ {

		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return nil, err
		}
		currentTemplateName := models.StripQuotes(tempCont.S("name").String())

		if currentTemplateName == templateName {
			d.Set("template_name", currentTemplateName)
			anpCount, err := tempCont.ArrayCount("anps")

			if err != nil {
				return nil, fmt.Errorf("No Anp found")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")

				if err != nil {
					return nil, err
				}
				currentAnpName := models.StripQuotes(anpCont.S("name").String())
				if currentAnpName == anpName {
					d.Set("anp_name", currentAnpName)
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return nil, fmt.Errorf("No Epg found")
					}
					for k := 0; k < epgCount; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return nil, err
						}
						currentEpgName := models.StripQuotes(epgCont.S("name").String())
						if currentEpgName == epgName {
							d.Set("epg_name", currentEpgName)
							usegCount, err := epgCont.ArrayCount("uSegAttrs")
							if err != nil {
								return nil, fmt.Errorf("No usegAttrs found")
							}
							for s := 0; s < usegCount; s++ {
								usegCont, err := epgCont.ArrayElement(s, "uSegAttrs")
								if err != nil {
									return nil, err
								}
								currentName := models.StripQuotes(usegCont.S("name").String())
								if currentName == name {
									d.SetId(currentName)
									d.Set("name", currentName)
									d.Set("operator", models.StripQuotes(usegCont.S("operator").String()))
									d.Set("useg_type", models.StripQuotes(usegCont.S("type").String()))
									d.Set("value", models.StripQuotes(usegCont.S("value").String()))

									category := models.StripQuotes(usegCont.S("category").String())
									desc := models.StripQuotes(usegCont.S("description").String())

									if category != "{}" {
										d.Set("category", category)
									} else {
										d.Set("category", "")
									}

									if desc != "{}" {
										d.Set("description", desc)
									} else {
										d.Set("description", "")
									}

									if usegCont.Exists("fvSubnet") {
										usegSubnet, _ := strconv.ParseBool(models.StripQuotes(usegCont.S("fvSubnet").String()))
										d.Set("useg_subnet", usegSubnet)
									}

									found = true
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
		}
		if found {
			break
		}
	}

	if !found {

		d.SetId("")
		d.Set("name", "")
		d.Set("operator", "")
		d.Set("useg_type", "")
		d.Set("value", "")
		return nil, fmt.Errorf("Unable to find Schema template anp epg useg attribute %s", name)
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaTemplateAnpEpgUsegAttrCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema Template Anp Epg UsegAttr: Beginning Creation")
	msoClient := m.(*client.Client)

	var schemaId string
	if schema_id, ok := d.GetOk("schema_id"); ok {
		schemaId = schema_id.(string)
	}

	var templateName string
	if template, ok := d.GetOk("template_name"); ok {
		templateName = template.(string)
	}

	var anpName string
	if name, ok := d.GetOk("anp_name"); ok {
		anpName = name.(string)
	}

	var epgName string
	if name, ok := d.GetOk("epg_name"); ok {
		epgName = name.(string)
	}

	name := d.Get("name").(string)

	usegType := d.Get("useg_type").(string)

	var desc string

	if tmpVar, ok := d.GetOk("description"); ok {
		desc = tmpVar.(string)
	}

	var operator string

	if tmpVar, ok := d.GetOk("operator"); ok {
		operator = tmpVar.(string)
	} else {
		operator = "equals"
	}

	var category string

	if tmpVar, ok := d.GetOk("category"); ok {
		category = tmpVar.(string)
	}

	val := d.Get("value").(string)

	usegSubnet := false
	if tempVar, ok := d.GetOk("useg_subnet"); ok {
		usegSubnet = tempVar.(bool)
	}

	usegAttrMap := make(map[string]interface{})

	usegAttrMap["name"] = name
	usegAttrMap["displayName"] = name
	usegAttrMap["type"] = usegType
	usegAttrMap["operator"] = operator
	usegAttrMap["value"] = val
	if category != "" {
		usegAttrMap["category"] = category
	}
	if desc != "" {
		usegAttrMap["description"] = desc
	}
	if usegType == "ip" {
		usegAttrMap["fvSubnet"] = usegSubnet
	}

	if usegType == "ip" || usegType == "mac" || usegType == "dns" {
		usegAttrMap["operator"] = "equals"
	}

	path := fmt.Sprintf("/templates/%s/anps/%s/epgs/%s/uSegAttrs/-", templateName, anpName, epgName)
	usegAttrApp := models.NewSchemaTemplateAnpEpgUsegAttr("add", path, usegAttrMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), usegAttrApp)
	if err != nil {
		log.Println(err)
		return err
	}

	d.SetId(fmt.Sprintf("%v", name))
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceMSOSchemaTemplateAnpEpgUsegAttrRead(d, m)
}

func resourceMSOSchemaTemplateAnpEpgUsegAttrUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema Template Anp Epg Subnet: Beginning Updating")
	msoClient := m.(*client.Client)

	var schemaId string
	if schema_id, ok := d.GetOk("schema_id"); ok {
		schemaId = schema_id.(string)
	}

	var templateName string
	if template, ok := d.GetOk("template_name"); ok {
		templateName = template.(string)
	}

	var anpName string
	if name, ok := d.GetOk("anp_name"); ok {
		anpName = name.(string)
	}

	var epgName string
	if name, ok := d.GetOk("epg_name"); ok {
		epgName = name.(string)
	}

	name := d.Get("name").(string)

	usegType := d.Get("useg_type").(string)

	var desc string

	if tmpVar, ok := d.GetOk("description"); ok {
		desc = tmpVar.(string)
	}

	var operator string

	if tmpVar, ok := d.GetOk("operator"); ok {
		operator = tmpVar.(string)
	} else {
		operator = "equals"
	}

	var category string

	if tmpVar, ok := d.GetOk("category"); ok {
		category = tmpVar.(string)
	}

	val := d.Get("value").(string)

	usegSubnet := false
	if tempVar, ok := d.GetOk("useg_subnet"); ok {
		usegSubnet = tempVar.(bool)
	}

	usegAttrMap := make(map[string]interface{})

	usegAttrMap["name"] = name
	usegAttrMap["displayName"] = name
	usegAttrMap["type"] = usegType
	usegAttrMap["operator"] = operator
	usegAttrMap["value"] = val
	if category != "" {
		usegAttrMap["category"] = category
	}
	if desc != "" {
		usegAttrMap["description"] = desc
	}
	if usegType == "ip" {
		usegAttrMap["fvSubnet"] = usegSubnet
	}

	if usegType == "ip" || usegType == "mac" || usegType == "dns" {
		usegAttrMap["operator"] = "equals"
	}

	path := fmt.Sprintf("/templates/%s/anps/%s/epgs/%s/uSegAttrs/%s", templateName, anpName, epgName, name)
	usegAttrApp := models.NewSchemaTemplateAnpEpgUsegAttr("replace", path, usegAttrMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), usegAttrApp)
	if err != nil {
		log.Println(err)
		return err
	}

	d.SetId(fmt.Sprintf("%v", name))
	log.Printf("[DEBUG] %s: Updating finished successfully", d.Id())

	return resourceMSOSchemaTemplateAnpEpgUsegAttrRead(d, m)
}

func resourceMSOSchemaTemplateAnpEpgUsegAttrRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return errorForObjectNotFound(err, d.Id(), cont, d)
	}

	count, err := cont.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No Template found")
	}

	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)
	name := d.Get("name").(string)
	found := false

	for i := 0; i < count; i++ {

		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		currentTemplateName := models.StripQuotes(tempCont.S("name").String())

		if currentTemplateName == templateName {
			d.Set("template_name", currentTemplateName)
			anpCount, err := tempCont.ArrayCount("anps")

			if err != nil {
				return fmt.Errorf("No Anp found")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")

				if err != nil {
					return err
				}
				currentAnpName := models.StripQuotes(anpCont.S("name").String())
				if currentAnpName == anpName {
					d.Set("anp_name", currentAnpName)
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return fmt.Errorf("No Epg found")
					}
					for k := 0; k < epgCount; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return err
						}
						currentEpgName := models.StripQuotes(epgCont.S("name").String())
						if currentEpgName == epgName {
							d.Set("epg_name", currentEpgName)
							usegCount, err := epgCont.ArrayCount("uSegAttrs")
							if err != nil {
								return fmt.Errorf("No usegAttrs found")
							}
							for s := 0; s < usegCount; s++ {
								usegCont, err := epgCont.ArrayElement(s, "uSegAttrs")
								if err != nil {
									return err
								}
								currentName := models.StripQuotes(usegCont.S("name").String())
								if currentName == name {
									d.SetId(currentName)
									d.Set("name", currentName)
									d.Set("operator", models.StripQuotes(usegCont.S("operator").String()))
									d.Set("useg_type", models.StripQuotes(usegCont.S("type").String()))
									d.Set("value", models.StripQuotes(usegCont.S("value").String()))

									category := models.StripQuotes(usegCont.S("category").String())
									desc := models.StripQuotes(usegCont.S("description").String())

									if category != "{}" {
										d.Set("category", category)
									} else {
										d.Set("category", "")
									}

									if desc != "{}" {
										d.Set("description", desc)
									} else {
										d.Set("description", "")
									}

									if usegCont.Exists("fvSubnet") {
										usegSubnet, _ := strconv.ParseBool(models.StripQuotes(usegCont.S("fvSubnet").String()))
										d.Set("useg_subnet", usegSubnet)
									}

									found = true
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
		}
		if found {
			break
		}
	}

	if !found {
		d.SetId("")
		d.Set("name", "")
		d.Set("operator", "")
		d.Set("useg_type", "")
		d.Set("value", "")
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSOSchemaTemplateAnpEpgUsegAttrDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	msoClient := m.(*client.Client)
	var schemaId string
	if schema_id, ok := d.GetOk("schema_id"); ok {
		schemaId = schema_id.(string)
	}

	var templateName string
	if template, ok := d.GetOk("template_name"); ok {
		templateName = template.(string)
	}

	var anpName string
	if name, ok := d.GetOk("anp_name"); ok {
		anpName = name.(string)
	}

	var epgName string
	if name, ok := d.GetOk("epg_name"); ok {
		epgName = name.(string)
	}

	name := d.Get("name").(string)

	usegType := d.Get("useg_type").(string)

	var desc string

	if tmpVar, ok := d.GetOk("description"); ok {
		desc = tmpVar.(string)
	}

	var operator string

	if tmpVar, ok := d.GetOk("operator"); ok {
		operator = tmpVar.(string)
	} else {
		operator = "equals"
	}

	var category string

	if tmpVar, ok := d.GetOk("category"); ok {
		category = tmpVar.(string)
	}

	val := d.Get("value").(string)

	usegSubnet := false
	if tempVar, ok := d.GetOk("useg_subnet"); ok {
		usegSubnet = tempVar.(bool)
	}

	usegAttrMap := make(map[string]interface{})

	usegAttrMap["name"] = name
	usegAttrMap["displayName"] = name
	usegAttrMap["type"] = usegType
	usegAttrMap["operator"] = operator
	usegAttrMap["value"] = val
	if category != "" {
		usegAttrMap["category"] = category
	}
	if desc != "" {
		usegAttrMap["description"] = desc
	}
	if usegType == "ip" {
		usegAttrMap["fvSubnet"] = usegSubnet
	}

	if usegType == "ip" || usegType == "mac" || usegType == "dns" {
		usegAttrMap["operator"] = "equals"
	}

	path := fmt.Sprintf("/templates/%s/anps/%s/epgs/%s/uSegAttrs/%s", templateName, anpName, epgName, name)
	usegAttrApp := models.NewSchemaTemplateAnpEpgUsegAttr("remove", path, usegAttrMap)
	response, errs := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), usegAttrApp)

	// Ignoring Error with code 141: Resource Not Found when deleting
	if errs != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return errs
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return nil
}
