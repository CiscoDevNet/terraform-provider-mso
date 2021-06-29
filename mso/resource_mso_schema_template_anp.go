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

func resourceMSOSchemaTemplateAnp() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaTemplateAnpCreate,
		Update: resourceMSOSchemaTemplateAnpUpdate,
		Read:   resourceMSOSchemaTemplateAnpRead,
		Delete: resourceMSOSchemaTemplateAnpDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaTemplateAnpImport,
		},

		Schema: (map[string]*schema.Schema{

			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"template": &schema.Schema{
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

			"display_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		}),
	}
}

func resourceMSOSchemaTemplateAnpImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
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
	found := false

	for i := 0; i < count; i++ {

		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return nil, err
		}
		currentTemplateName := models.StripQuotes(tempCont.S("name").String())

		if currentTemplateName == templateName {
			d.Set("template", currentTemplateName)
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
				log.Println("currentanpname", currentAnpName)
				if currentAnpName == anpName {
					log.Println("found correct anpname")
					d.SetId(get_attribute[4])
					d.Set("name", currentAnpName)
					if anpCont.Exists("displayName") {
						d.Set("display_name", models.StripQuotes(anpCont.S("displayName").String()))
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
	if !found {
		d.SetId("")
		return nil, fmt.Errorf("The ANP is not found")
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaTemplateAnpCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema Template Anp: Beginning Creation")
	msoClient := m.(*client.Client)

	var schemaId string
	if schema_id, ok := d.GetOk("schema_id"); ok {
		schemaId = schema_id.(string)
	}

	var templateName string
	if template, ok := d.GetOk("template"); ok {
		templateName = template.(string)
	}

	var Name string
	if name, ok := d.GetOk("name"); ok {
		Name = name.(string)
	}

	var displayName string
	if display_name, ok := d.GetOk("display_name"); ok {
		displayName = display_name.(string)
	}

	schemaTemplateAnpApp := models.NewSchemaTemplateAnp("add", "/templates/"+templateName+"/anps/-", Name, displayName)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), schemaTemplateAnpApp)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("[DEBUG] %s: Creation finished successfully", Name)

	return resourceMSOSchemaTemplateAnpRead(d, m)
}

func resourceMSOSchemaTemplateAnpUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema Template Anp: Beginning Updating")
	msoClient := m.(*client.Client)

	var schemaId string
	if schema_id, ok := d.GetOk("schema_id"); ok {
		schemaId = schema_id.(string)
	}

	var templateName string
	if template, ok := d.GetOk("template"); ok {
		templateName = template.(string)
	}

	var Name string
	if name, ok := d.GetOk("name"); ok {
		Name = name.(string)
	}

	var displayName string
	if display_name, ok := d.GetOk("display_name"); ok {
		displayName = display_name.(string)
	}

	schemaTemplateAnpApp := models.NewSchemaTemplateAnp("replace", "/templates/"+templateName+"/anps/"+Name, Name, displayName)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), schemaTemplateAnpApp)
	if err != nil {
		log.Println(err)
		return err
	}

	id := models.StripQuotes(Name)
	d.SetId(fmt.Sprintf("%v", id))
	log.Printf("[DEBUG] %s: Updating finished successfully", Name)

	return resourceMSOSchemaTemplateAnpRead(d, m)
}

func resourceMSOSchemaTemplateAnpRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}

	count, err := cont.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No Template found")
	}

	templateName := d.Get("template").(string)
	anpName := d.Get("name").(string)
	found := false

	for i := 0; i < count; i++ {

		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		currentTemplateName := models.StripQuotes(tempCont.S("name").String())

		if currentTemplateName == templateName {
			d.Set("template", currentTemplateName)
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
				log.Println("currentanpname", currentAnpName)
				if currentAnpName == anpName {
					log.Println("found correct anpname")
					d.SetId(schemaId + "/templates/" + currentTemplateName + "/anps/" + currentAnpName)
					d.Set("name", currentAnpName)
					if anpCont.Exists("displayName") {
						d.Set("display_name", models.StripQuotes(anpCont.S("displayName").String()))
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
	if !found {
		d.SetId("")
		d.Set("name", "")
		d.Set("display_name", "")
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSOSchemaTemplateAnpDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	template := d.Get("template").(string)
	name := d.Get("name").(string)
	schemaTemplateAnpApp := models.NewSchemaTemplateAnp("remove", "/templates/"+template+"/anps/"+name, "", "")
	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), schemaTemplateAnpApp)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return nil
}
