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

func resourceMSOSchemaSiteBd() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteBdCreate,
		Read:   resourceMSOSchemaSiteBdRead,
		Update: resourceMSOSchemaSiteBdUpdate,
		Delete: resourceMSOSchemaSiteBdDelete,

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
			"bd_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"host": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
		}),
	}
}

func resourceMSOSchemaSiteBdCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Bd: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	bdName := d.Get("bd_name").(string)

	var host bool

	if tempvar, ok := d.GetOk("host"); ok {
		host = tempvar.(bool)
	}

	var bd_schema_id, bd_template_name string
	bd_schema_id = schemaId
	bd_template_name = templateName

	bdRefMap := make(map[string]interface{})
	bdRefMap["schemaId"] = bd_schema_id
	bdRefMap["templateName"] = bd_template_name
	bdRefMap["bdName"] = bdName

	path := fmt.Sprintf("/sites/%s-%s/bds/-", siteId, templateName)
	bdStruct := models.NewSchemaSiteBd("add", path, bdRefMap, host)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), bdStruct)

	if err != nil {
		return err
	}
	return resourceMSOSchemaSiteBdRead(d, m)
}

func resourceMSOSchemaSiteBdRead(d *schema.ResourceData, m interface{}) error {
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
	statebd := d.Get("bd_name").(string)
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())

		if apiSite == stateSite {
			bdCount, err := tempCont.ArrayCount("bds")
			if err != nil {
				return fmt.Errorf("Unable to get bd list")
			}
			for j := 0; j < bdCount; j++ {
				bdCont, err := tempCont.ArrayElement(j, "bds")
				if err != nil {
					return err
				}
				bdRef := models.StripQuotes(bdCont.S("bdRef").String())
				re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/bds/(.*)")
				match := re.FindStringSubmatch(bdRef)
				if match[3] == statebd {
					d.SetId(match[3])
					d.Set("bd_name", match[3])
					d.Set("schema_id", match[1])
					d.Set("template_name", match[2])
					d.Set("site_id", apiSite)
					if bdCont.Exists("hostBasedRouting") {
						d.Set("host", bdCont.S("hostBasedRouting").Data().(bool))
					}
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

func resourceMSOSchemaSiteBdUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Bd: Beginning Updation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	bdName := d.Get("bd_name").(string)

	var host bool

	if tempvar, ok := d.GetOk("host"); ok {
		host = tempvar.(bool)
	}

	var bd_schema_id, bd_template_name string
	bd_schema_id = schemaId
	bd_template_name = templateName

	bdRefMap := make(map[string]interface{})
	bdRefMap["schemaId"] = bd_schema_id
	bdRefMap["templateName"] = bd_template_name
	bdRefMap["bdName"] = bdName

	path := fmt.Sprintf("/sites/%s-%s/bds/%s", siteId, templateName, bdName)
	bdStruct := models.NewSchemaSiteBd("replace", path, bdRefMap, host)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), bdStruct)

	if err != nil {
		return err
	}
	return resourceMSOSchemaSiteBdRead(d, m)
}

func resourceMSOSchemaSiteBdDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Bd: Beginning Deletion")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	bdName := d.Get("bd_name").(string)

	var host bool

	if tempvar, ok := d.GetOk("host"); ok {
		host = tempvar.(bool)
	}

	var bd_schema_id, bd_template_name string
	bd_schema_id = schemaId
	bd_template_name = templateName
	bdRefMap := make(map[string]interface{})
	bdRefMap["schemaId"] = bd_schema_id
	bdRefMap["templateName"] = bd_template_name
	bdRefMap["bdName"] = bdName

	path := fmt.Sprintf("/sites/%s-%s/bds/%s", siteId, templateName, bdName)
	bdStruct := models.NewSchemaSiteBd("remove", path, bdRefMap, host)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), bdStruct)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
