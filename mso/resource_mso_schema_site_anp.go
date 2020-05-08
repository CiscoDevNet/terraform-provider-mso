package mso

import (
	"fmt"
	"log"
	"regexp"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceMSOSchemaSiteAnp() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteAnpCreate,
		Read:   resourceMSOSchemaSiteAnpRead,
		Update: resourceMSOSchemaSiteAnpUpdate,
		Delete: resourceMSOSchemaSiteAnpDelete,

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
			"anp_schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"anp_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func resourceMSOSchemaSiteAnpCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Anp: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)

	var anp_schema_id, anp_template_name string

	if tempVar, ok := d.GetOk("anp_schema_id"); ok {
		anp_schema_id = tempVar.(string)
	} else {
		anp_schema_id = schemaId
	}
	if tempVar, ok := d.GetOk("anp_template_name"); ok {
		anp_template_name = tempVar.(string)
	} else {
		anp_template_name = templateName
	}

	anpRefMap := make(map[string]interface{})
	anpRefMap["schemaId"] = anp_schema_id
	anpRefMap["templateName"] = anp_template_name
	anpRefMap["anpName"] = anpName

	path := fmt.Sprintf("/sites/%s-%s/anps/-", siteId, templateName)
	anpStruct := models.NewSchemaSiteAnp("add", path, anpRefMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpStruct)

	if err != nil {
		return err
	}
	return resourceMSOSchemaSiteAnpRead(d, m)
}

func resourceMSOSchemaSiteAnpRead(d *schema.ResourceData, m interface{}) error {
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
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())

		if apiSite == stateSite {
			d.Set("site_id", apiSite)
			d.Set("template_name", models.StripQuotes(tempCont.S("templateName").String()))
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
					d.SetId(match[3])
					d.Set("anp_name", match[3])
					d.Set("anp_schema_id", match[1])
					d.Set("anp_template_name", match[2])
					found = true
					break
				}
			}
		}
	}

	if !found {
		d.SetId("")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}

func resourceMSOSchemaSiteAnpUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Anp: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)

	var anp_schema_id, anp_template_name string

	if tempVar, ok := d.GetOk("anp_schema_id"); ok {
		anp_schema_id = tempVar.(string)
	} else {
		anp_schema_id = schemaId
	}
	if tempVar, ok := d.GetOk("anp_template_name"); ok {
		anp_template_name = tempVar.(string)
	} else {
		anp_template_name = templateName
	}

	anpRefMap := make(map[string]interface{})
	anpRefMap["schemaId"] = anp_schema_id
	anpRefMap["templateName"] = anp_template_name
	anpRefMap["anpName"] = anpName

	path := fmt.Sprintf("/sites/%s-%s/anps/%s", siteId, templateName, anpName)
	anpStruct := models.NewSchemaSiteAnp("replace", path, anpRefMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpStruct)

	if err != nil {
		return err
	}
	return resourceMSOSchemaSiteAnpRead(d, m)
}

func resourceMSOSchemaSiteAnpDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Anp: Beginning Deletion")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)

	var anp_schema_id, anp_template_name string

	if tempVar, ok := d.GetOk("anp_schema_id"); ok {
		anp_schema_id = tempVar.(string)
	} else {
		anp_schema_id = schemaId
	}
	if tempVar, ok := d.GetOk("anp_template_name"); ok {
		anp_template_name = tempVar.(string)
	} else {
		anp_template_name = templateName
	}

	anpRefMap := make(map[string]interface{})
	anpRefMap["schemaId"] = anp_schema_id
	anpRefMap["templateName"] = anp_template_name
	anpRefMap["anpName"] = anpName

	path := fmt.Sprintf("/sites/%s-%s/anps/%s", siteId, templateName, anpName)
	anpStruct := models.NewSchemaSiteAnp("remove", path, anpRefMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpStruct)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
