package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceMSOSchemaTemplateAnpEpgUsegAttr() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOSchemaTemplateAnpEpgUsegAttrRead,

		Schema: (map[string]*schema.Schema{
			"schema_id": &schema.Schema{
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
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"useg_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"operator": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"category": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"useg_subnet": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		}),
	}
}

func dataSourceMSOSchemaTemplateAnpEpgUsegAttrRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)
	name := d.Get("name").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}

	usegCont, err := findUsegAttrContainer(cont, templateName, anpName, epgName, name)
	if err != nil {
		return err
	}
	if usegCont == nil {
		d.SetId("")
		return fmt.Errorf("Unable to find the ANP EPG uSeg Attribute %s in Template %s of Schema Id %s ", name, templateName, schemaId)
	}

	d.Set("template_name", templateName)
	d.Set("anp_name", anpName)
	d.Set("epg_name", epgName)
	setUsegAttrAttributes(d, usegCont)
	d.SetId(fmt.Sprintf("%s/templates/%s/anps/%s/epgs/%s/uSegAttrs/%s", schemaId, templateName, anpName, epgName, name))

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
