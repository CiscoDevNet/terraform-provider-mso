package mso

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOTemplateL3out() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOTemplateL3outCreate,
		Read:   resourceMSOTemplateL3outRead,
		Update: resourceMSOTemplateL3outUpdate,
		Delete: resourceMSOTemplateL3outDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOTemplateL3outImport,
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
			"l3out_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"display_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"vrf_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"vrf_schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vrf_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func resourceMSOTemplateL3outImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())

	msoClient := m.(*client.Client)
	get_attribute := strings.Split(d.Id(), "/")
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", get_attribute[0]))
	if err != nil {
		return nil, err
	}
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return nil, fmt.Errorf("No Template found")
	}
	stateTemplate := get_attribute[2]
	found := false
	stateL3out := get_attribute[4]
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return nil, err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {
			l3outCount, err := tempCont.ArrayCount("intersiteL3outs")
			if err != nil {
				return nil, fmt.Errorf("Unable to get L3out list")
			}
			for j := 0; j < l3outCount; j++ {
				l3outCont, err := tempCont.ArrayElement(j, "intersiteL3outs")
				if err != nil {
					return nil, err
				}
				apiL3out := models.StripQuotes(l3outCont.S("name").String())
				if apiL3out == stateL3out {
					d.SetId(apiL3out)
					d.Set("l3out_name", get_attribute[4])
					d.Set("schema_id", get_attribute[0])
					d.Set("template_name", apiTemplate)
					d.Set("display_name", models.StripQuotes(l3outCont.S("displayName").String()))
					d.Set("description", models.StripQuotes(l3outCont.S("description").String()))

					vrfRef := models.StripQuotes(l3outCont.S("vrfRef").String())
					re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/vrfs/(.*)")
					match := re.FindStringSubmatch(vrfRef)
					d.Set("vrf_name", match[3])
					d.Set("vrf_schema_id", match[1])
					d.Set("vrf_template_name", match[2])

					found = true
					break
				}
			}
		}
	}

	if !found {
		d.SetId("")
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil

}

func resourceMSOTemplateL3outCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template L3out: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	l3outName := d.Get("l3out_name").(string)
	displayName := d.Get("display_name").(string)
	templateName := d.Get("template_name").(string)
	vrfName := d.Get("vrf_name").(string)

	var vrf_schema_id, vrf_template_name string

	if tempVar, ok := d.GetOk("vrf_schema_id"); ok {
		vrf_schema_id = tempVar.(string)
	} else {
		vrf_schema_id = schemaID
	}
	if tempVar, ok := d.GetOk("vrf_template_name"); ok {
		vrf_template_name = tempVar.(string)
	} else {
		vrf_template_name = templateName
	}

	var description string
	if tempVar, ok := d.GetOk("description"); ok {
		description = tempVar.(string)
	}

	vrfRefMap := make(map[string]interface{})
	vrfRefMap["schemaId"] = vrf_schema_id
	vrfRefMap["templateName"] = vrf_template_name
	vrfRefMap["vrfName"] = vrfName

	path := fmt.Sprintf("/templates/%s/intersiteL3outs/-", templateName)
	l3outStruct := models.NewTemplateL3out("add", path, l3outName, displayName, description, vrfRefMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), l3outStruct)

	if err != nil {
		return err
	}
	return resourceMSOTemplateL3outRead(d, m)
}

func resourceMSOTemplateL3outRead(d *schema.ResourceData, m interface{}) error {
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
	stateTemplate := d.Get("template_name").(string)
	found := false
	stateL3out := d.Get("l3out_name")
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {
			l3outCount, err := tempCont.ArrayCount("intersiteL3outs")
			if err != nil {
				return fmt.Errorf("Unable to get L3out list")
			}
			for j := 0; j < l3outCount; j++ {
				l3outCont, err := tempCont.ArrayElement(j, "intersiteL3outs")
				if err != nil {
					return err
				}
				apiL3out := models.StripQuotes(l3outCont.S("name").String())
				if apiL3out == stateL3out {
					d.SetId(apiL3out)
					d.Set("l3out_name", apiL3out)
					d.Set("schema_id", schemaId)
					d.Set("template_name", apiTemplate)
					d.Set("display_name", models.StripQuotes(l3outCont.S("displayName").String()))
					d.Set("description", models.StripQuotes(l3outCont.S("description").String()))

					vrfRef := models.StripQuotes(l3outCont.S("vrfRef").String())
					re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/vrfs/(.*)")
					match := re.FindStringSubmatch(vrfRef)
					d.Set("vrf_name", match[3])
					d.Set("vrf_schema_id", match[1])
					d.Set("vrf_template_name", match[2])

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

func resourceMSOTemplateL3outUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template L3out: Beginning Update")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	l3outName := d.Get("l3out_name").(string)
	displayName := d.Get("display_name").(string)
	templateName := d.Get("template_name").(string)
	vrfName := d.Get("vrf_name").(string)

	var vrf_schema_id, vrf_template_name string

	if tempVar, ok := d.GetOk("vrf_schema_id"); ok {
		vrf_schema_id = tempVar.(string)
	} else {
		vrf_schema_id = schemaID
	}
	if tempVar, ok := d.GetOk("vrf_template_name"); ok {
		vrf_template_name = tempVar.(string)
	} else {
		vrf_template_name = templateName
	}

	var description string
	if tempVar, ok := d.GetOk("description"); ok {
		description = tempVar.(string)
	}

	vrfRefMap := make(map[string]interface{})
	vrfRefMap["schemaId"] = vrf_schema_id
	vrfRefMap["templateName"] = vrf_template_name
	vrfRefMap["vrfName"] = vrfName

	path := fmt.Sprintf("/templates/%s/intersiteL3outs/%s", templateName, l3outName)
	l3outStruct := models.NewTemplateL3out("replace", path, l3outName, displayName, description, vrfRefMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), l3outStruct)

	if err != nil {
		return err
	}
	return resourceMSOTemplateL3outRead(d, m)
}

func resourceMSOTemplateL3outDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template L3out: Beginning Update")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	l3outName := d.Get("l3out_name").(string)
	displayName := d.Get("display_name").(string)
	templateName := d.Get("template_name").(string)
	vrfName := d.Get("vrf_name").(string)

	var vrf_schema_id, vrf_template_name string

	if tempVar, ok := d.GetOk("vrf_schema_id"); ok {
		vrf_schema_id = tempVar.(string)
	} else {
		vrf_schema_id = schemaID
	}
	if tempVar, ok := d.GetOk("vrf_template_name"); ok {
		vrf_template_name = tempVar.(string)
	} else {
		vrf_template_name = templateName
	}

	var description string
	if tempVar, ok := d.GetOk("description"); ok {
		description = tempVar.(string)
	}

	vrfRefMap := make(map[string]interface{})
	vrfRefMap["schemaId"] = vrf_schema_id
	vrfRefMap["templateName"] = vrf_template_name
	vrfRefMap["vrfName"] = vrfName

	path := fmt.Sprintf("/templates/%s/intersiteL3outs/%s", templateName, l3outName)
	l3outStruct := models.NewTemplateL3out("remove", path, l3outName, displayName, description, vrfRefMap)

	response, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), l3outStruct)

	// Ignoring Error with code 141: Resource Not Found when deleting
	if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return err
	}
	d.SetId("")
	return nil
}
