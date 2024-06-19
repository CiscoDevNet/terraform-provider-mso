package mso

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOSchemaSiteBd() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteBdCreate,
		Read:   resourceMSOSchemaSiteBdRead,
		Update: resourceMSOSchemaSiteBdUpdate,
		Delete: resourceMSOSchemaSiteBdDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaSiteBdImport,
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
			"bd_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"host_route": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"svi_mac": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func resourceMSOSchemaSiteBdImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
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
	stateSite := get_attribute[1]
	stateTemplate := get_attribute[2]
	found := false
	statebd := get_attribute[3]
	for i := 0; i < count && !found; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return nil, err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())
		apiTemplate := models.StripQuotes(tempCont.S("templateName").String())

		if apiSite == stateSite && apiTemplate == stateTemplate {
			bdCount, err := tempCont.ArrayCount("bds")
			if err != nil {
				return nil, fmt.Errorf("Unable to get bd list")
			}
			for j := 0; j < bdCount; j++ {
				bdCont, err := tempCont.ArrayElement(j, "bds")
				if err != nil {
					return nil, err
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
						d.Set("host_route", bdCont.S("hostBasedRouting").Data().(bool))
					}
					if bdCont.Exists("mac") {
						d.Set("svi_mac", models.StripQuotes(bdCont.S("mac").String()))
					}
					found = true
					break
				}
			}
		}
	}

	if !found {
		d.SetId("")
		return nil, fmt.Errorf("Unable to find the given Schema Site Bd")
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaSiteBdCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Bd: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	bdName := d.Get("bd_name").(string)

	var host bool
	var mac string

	if tempvar, ok := d.GetOk("host_route"); ok {
		host = tempvar.(bool)
	}

	if tempvar, ok := d.GetOk("svi_mac"); ok {
		mac = tempvar.(string)
	}

	var bd_schema_id, bd_template_name string
	bd_schema_id = schemaId
	bd_template_name = templateName

	bdRefMap := make(map[string]interface{})
	bdRefMap["schemaId"] = bd_schema_id
	bdRefMap["templateName"] = bd_template_name
	bdRefMap["bdName"] = bdName

	versionInt, err := msoClient.CompareVersion("4.0.0.0")
	if err != nil {
		return err
	}

	if versionInt != 1 {
		path := fmt.Sprintf("/sites/%s-%s/bds/%s", siteId, templateName, bdName)
		bdStruct := models.NewSchemaSiteBd("replace", path, mac, bdRefMap, host)
		_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), bdStruct)
	}

	if versionInt == 1 || err != nil {
		path := fmt.Sprintf("/sites/%s-%s/bds/-", siteId, templateName)
		bdStruct := models.NewSchemaSiteBd("add", path, mac, bdRefMap, host)
		_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), bdStruct)
	}

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
		return errorForObjectNotFound(err, d.Id(), cont, d)
	}
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return fmt.Errorf("No Sites found")
	}
	stateSite := d.Get("site_id").(string)
	stateTemplate := d.Get("template_name").(string)
	found := false
	statebd := d.Get("bd_name").(string)
	for i := 0; i < count && !found; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())
		apiTemplate := models.StripQuotes(tempCont.S("templateName").String())

		if apiSite == stateSite && apiTemplate == stateTemplate {
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
						d.Set("host_route", bdCont.S("hostBasedRouting").Data().(bool))
					}
					if bdCont.Exists("mac") {
						d.Set("svi_mac", models.StripQuotes(bdCont.S("mac").String()))
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
	var mac string

	if tempvar, ok := d.GetOk("host_route"); ok {
		host = tempvar.(bool)
	}

	if tempvar, ok := d.GetOk("svi_mac"); ok {
		mac = tempvar.(string)
	}

	var bd_schema_id, bd_template_name string
	bd_schema_id = schemaId
	bd_template_name = templateName

	bdRefMap := make(map[string]interface{})
	bdRefMap["schemaId"] = bd_schema_id
	bdRefMap["templateName"] = bd_template_name
	bdRefMap["bdName"] = bdName

	payloadCon := container.New()
	payloadCon.Array()

	err := setPatchPayloadToContainer(payloadCon, "replace", fmt.Sprintf("/sites/%s-%s/bds/%s/bdRef", siteId, templateName, bdName), bdRefMap)
	if err != nil {
		return err
	}

	err = setPatchPayloadToContainer(payloadCon, "replace", fmt.Sprintf("/sites/%s-%s/bds/%s/hostBasedRouting", siteId, templateName, bdName), host)
	if err != nil {
		return err
	}

	if mac != "" {
		err := setPatchPayloadToContainer(payloadCon, "replace", fmt.Sprintf("/sites/%s-%s/bds/%s/mac", siteId, templateName, bdName), mac)
		if err != nil {
			return err
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/schemas/%s", schemaId), payloadCon)
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
	var mac string

	if tempvar, ok := d.GetOk("host_route"); ok {
		host = tempvar.(bool)
	}

	if tempvar, ok := d.GetOk("svi_mac"); ok {
		mac = tempvar.(string)
	}

	var bd_schema_id, bd_template_name string
	bd_schema_id = schemaId
	bd_template_name = templateName
	bdRefMap := make(map[string]interface{})
	bdRefMap["schemaId"] = bd_schema_id
	bdRefMap["templateName"] = bd_template_name
	bdRefMap["bdName"] = bdName

	path := fmt.Sprintf("/sites/%s-%s/bds/%s", siteId, templateName, bdName)
	bdStruct := models.NewSchemaSiteBd("remove", path, mac, bdRefMap, host)

	response, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), bdStruct)

	// Ignoring Error with code 141: Resource Not Found when deleting
	if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return err
	}
	d.SetId("")
	return nil
}
