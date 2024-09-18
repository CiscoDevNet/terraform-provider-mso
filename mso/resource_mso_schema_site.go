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

func resourceMSOSchemaSite() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteCreate,
		Read:   resourceMSOSchemaSiteRead,
		Update: resourceMSOSchemaSiteUpdate,
		Delete: resourceMSOSchemaSiteDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaSiteImport,
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

			"undeploy_on_destroy": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
		}),
	}
}

func resourceMSOSchemaSiteImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	msoClient := m.(*client.Client)
	get_attribute := strings.Split(d.Id(), "/")
	schemaId := get_attribute[0]
	name := get_attribute[2]
	con, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/sites"))
	if err != nil {
		return nil, err
	}

	data := con.S("sites").Data().([]interface{})
	var flag bool
	var count int
	for _, info := range data {
		val := info.(map[string]interface{})
		if val["name"].(string) == name {
			flag = true
			break
		}
		count = count + 1
	}
	if flag != true {
		return nil, fmt.Errorf("Site of specified name not found")
	}

	dataCon := con.S("sites").Index(count)
	stateSiteId := models.StripQuotes(dataCon.S("id").String())

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return nil, err
	}

	countSites, err := cont.ArrayCount("sites")
	if err != nil {
		return nil, fmt.Errorf("No Template found")
	}

	found := false

	for i := 0; i < countSites; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return nil, err
		}
		apiSiteId := models.StripQuotes(tempCont.S("siteId").String())
		apiTemplate := models.StripQuotes(tempCont.S("templateName").String())

		if apiSiteId == stateSiteId {
			d.SetId(apiSiteId)
			d.Set("schema_id", schemaId)
			d.Set("site_id", apiSiteId)
			d.Set("template_name", apiTemplate)
			found = true
		}

	}

	if !found {
		d.SetId("")
		d.Set("schema_id", "")
		d.Set("site_id", "")
		d.Set("template_name", "")
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaSiteCreate(d *schema.ResourceData, m interface{}) error {
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	siteId := d.Get("site_id").(string)

	schemasite := models.NewSchemaSite("add", "/sites/-", siteId, templateName)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), schemasite)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%v", siteId))
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceMSOSchemaSiteRead(d, m)
}

func resourceMSOSchemaSiteRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return errorForObjectNotFound(err, d.Id(), cont, d)
	}

	count, err := cont.ArrayCount("sites")
	if err != nil {
		return fmt.Errorf("No Template found")
	}

	stateSiteId := d.Get("site_id").(string)
	stateTemplate := d.Get("template_name").(string)
	found := false

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSiteId := models.StripQuotes(tempCont.S("siteId").String())
		apiTemplate := models.StripQuotes(tempCont.S("templateName").String())

		if apiSiteId == stateSiteId && apiTemplate == stateTemplate {
			d.SetId(apiSiteId)
			d.Set("schema_id", schemaId)
			d.Set("site_id", apiSiteId)
			d.Set("template_name", apiTemplate)
			found = true
		}

	}

	if !found {
		d.SetId("")
		d.Set("schema_id", "")
		d.Set("site_id", "")
		d.Set("template_name", "")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSOSchemaSiteUpdate(d *schema.ResourceData, m interface{}) error {
	d.Set("undeploy_on_destroy", d.Get("undeploy_on_destroy").(bool))
	return nil
}

func resourceMSOSchemaSiteDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)

	if d.Get("undeploy_on_destroy").(bool) {
		req, err := msoClient.MakeRestRequest("GET", fmt.Sprintf("mso/api/v1/deploy/status/schema/%s/template/%s", schemaId, templateName), nil, true)
		if err != nil {
			log.Printf("[DEBUG] MakeRestRequest failed with err: %s.", err)
			return err
		}
		obj, resp, err := msoClient.Do(req)
		if err != nil || resp.StatusCode != 200 {
			log.Printf("[DEBUG] Request failed with resp: %v. Err: %s.", resp, err)
			return err
		}
		if deployMap, ok := obj.Data().(map[string]interface{}); ok {
			if statusList, ok := deployMap["status"].([]map[string]interface{}); ok && len(statusList) > 0 {
				for _, statusMap := range statusList {
					if statusMap["siteId"] == siteId && statusMap["status"].(map[string]interface{})["siteStatus"] == "Succeeded" {
						versionInt, err := msoClient.CompareVersion("3.7.0.0")
						if versionInt == -1 {
							payload, err := container.ParseJSON([]byte(fmt.Sprintf(`{"schemaId": "%s", "templateName": "%s", "undeploy": ["%s"]}`, schemaId, templateName, siteId)))
							if err != nil {
								log.Printf("[DEBUG] Parse of JSON failed with err: %s.", err)
								return err
							}
							req, err := msoClient.MakeRestRequest("POST", "api/v1/task", payload, true)
							if err != nil {
								log.Printf("[DEBUG] MakeRestRequest failed with err: %s.", err)
								return err
							}
							_, resp, err := msoClient.Do(req)
							if err != nil || resp.StatusCode != 202 {
								log.Printf("[DEBUG] Request failed with resp: %v. Err: %s.", resp, err)
								return err
							}
						} else if err == nil {
							_, err := msoClient.GetViaURL(fmt.Sprintf("/api/v1/execute/schema/%s/template/%s?undeploy=%s", schemaId, templateName, siteId))
							if err != nil {
								return err
							}
						} else {
							log.Printf("[WARNING] Failed to compare version. Template could not be undeployed prior to schema site deletion. Err: %s.", err)
						}
						break
					}
				}
			}
		}
	}

	schemasite := models.NewSchemaSite("remove", fmt.Sprintf("/sites/%s-%s", siteId, templateName), siteId, templateName)

	response, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), schemasite)

	// Ignoring Error with code 141: Resource Not Found when deleting
	if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return nil
}
