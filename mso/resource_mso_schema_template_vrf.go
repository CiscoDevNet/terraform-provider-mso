package mso

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceMSOSchemaTemplateVrf() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaTemplateVrfCreate,
		Update: resourceMSOSchemaTemplateVrfUpdate,
		Read:   resourceMSOSchemaTemplateVrfRead,
		Delete: resourceMSOSchemaTemplateVrfDelete,

		Schema: (map[string]*schema.Schema{

			"schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"template": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"layer3_multicast": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func resourceMSOSchemaTemplateVrfCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema Template Vrf: Beginning Creation")
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

	var l3m bool
	if L3M, ok := d.GetOk("layer3_multicast"); ok {
		l3m = L3M.(bool)
	}

	schemaTemplateVrfApp := models.NewSchemaTemplateVrf("add", "/templates/"+templateName+"/vrfs/-", Name, displayName, l3m)

	cont, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), schemaTemplateVrfApp)
	if err != nil {
		log.Println(err)
		return err
	}

	id := models.StripQuotes(cont.S("id").String())
	d.SetId(fmt.Sprintf("%v", id))
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceMSOSchemaTemplateVrfRead(d, m)
}

func resourceMSOSchemaTemplateVrfUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema Template Vrf: Beginning Creation")
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

	var l3m bool
	if L3M, ok := d.GetOk("layer3_multicast"); ok {
		l3m = L3M.(bool)
	}

	schemaTemplateVrfApp := models.NewSchemaTemplateVrf("replace", "/templates/"+templateName+"/vrfs/"+Name, Name, displayName, l3m)

	cont, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), schemaTemplateVrfApp)
	if err != nil {
		log.Println(err)
		return err
	}

	id := models.StripQuotes(cont.S("id").String())
	d.SetId(fmt.Sprintf("%v", id))
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceMSOSchemaTemplateVrfRead(d, m)
}

func resourceMSOSchemaTemplateVrfRead(d *schema.ResourceData, m interface{}) error {
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
	vrfName := d.Get("name").(string)
	found := false

	for i := 0; i < count; i++ {

		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		currentTemplateName := models.StripQuotes(tempCont.S("name").String())

		if currentTemplateName == templateName {
			d.Set("template", currentTemplateName)
			vrfCount, err := tempCont.ArrayCount("vrfs")

			if err != nil {
				return fmt.Errorf("No Vrf found")
			}
			for j := 0; j < vrfCount; j++ {
				vrfCont, err := tempCont.ArrayElement(j, "vrfs")

				if err != nil {
					return err
				}
				currentVrfName := models.StripQuotes(vrfCont.S("name").String())
				log.Println("currentvrfname", currentVrfName)
				if currentVrfName == vrfName {
					log.Println("found correct vrfname")
					d.SetId(currentVrfName)
					d.Set("name", currentVrfName)
					if vrfCont.Exists("displayName") {
						d.Set("display_name", models.StripQuotes(vrfCont.S("displayName").String()))
					}
					if vrfCont.Exists("l3MCast") {
						l3Mcast, _ := strconv.ParseBool(models.StripQuotes(vrfCont.S("l3MCast").String()))
						d.Set("layer3_multicast", l3Mcast)
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

func resourceMSOSchemaTemplateVrfDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	template := d.Get("template").(string)
	name := d.Get("name").(string)
	var l3m bool
	if L3M, ok := d.GetOk("layer3_multicast"); ok {
		l3m = L3M.(bool)
	}
	schemaTemplateVrfApp := models.NewSchemaTemplateVrf("remove", "/templates/"+template+"/vrfs/"+name, "", "", l3m)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), schemaTemplateVrfApp)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return nil
}
