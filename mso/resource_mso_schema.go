package mso

import (
	"fmt"
	"log"
	"reflect"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOSchema() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaCreate,
		Update: resourceMSOSchemaUpdate,
		Read:   resourceMSOSchemaRead,
		Delete: resourceMSOSchemaDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaImport,
		},

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
				Deprecated:   "use template by specifying a name, display_name and tenant_id instead",
			},
			"tenant_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
				Deprecated:   "use template by specifying a name, display_name and tenant_id instead",
			},
			"template": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"display_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"tenant_id": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
					},
				},
			},
		}),
	}
}

func resourceMSOSchemaCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema: Beginning Creation")
	msoClient := m.(*client.Client)
	name := d.Get("name").(string)
	var schemaApp *models.Schema
	if tempVar, ok := d.GetOk("template_name"); ok {
		templateName := tempVar.(string)
		tenantId := d.Get("tenant_id").(string)
		schemaApp, _ = models.NewSchema("", name, templateName, tenantId, make([]interface{}, 0, 1))

	} else {
		templates := make([]interface{}, 0, 1)
		if val, ok := d.GetOk("template"); ok {
			template_list := val.(*schema.Set).List()
			for _, val := range template_list {
				map_templates := make(map[string]interface{})
				inner_templates := val.(map[string]interface{})
				map_templates["name"] = inner_templates["name"]
				map_templates["displayName"] = inner_templates["display_name"]
				map_templates["tenantId"] = inner_templates["tenant_id"]
				templates = append(templates, map_templates)
			}
		}
		schemaApp, _ = models.NewSchema("", name, "", "", templates)
	}

	cont, err := msoClient.Save("api/v1/schemas", schemaApp)
	if err != nil {
		return err
	}

	id := models.StripQuotes(cont.S("id").String())
	d.SetId(fmt.Sprintf("%s", id))
	log.Printf("[DEBUG] %s: Schema Creation finished successfully", d.Id())

	return resourceMSOSchemaRead(d, m)
}

func resourceMSOSchemaImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] Schema: Beginning Import")
	msoClient := m.(*client.Client)
	con, err := msoClient.GetViaURL("api/v1/schemas/" + d.Id())
	if err != nil {
		return nil, err
	}
	d.SetId(models.StripQuotes(con.S("id").String()))
	d.Set("name", models.StripQuotes(con.S("displayName").String()))
	count, err := con.ArrayCount("templates")
	if err != nil {
		return nil, fmt.Errorf("No Template found")
	}
	templates := make([]interface{}, 0)
	for i := 0; i < count; i++ {
		templatesCont, err := con.ArrayElement(i, "templates")
		if err != nil {
			return nil, fmt.Errorf("Unable to parse the templates list")
		}
		map_template := make(map[string]interface{})
		map_template["name"] = models.StripQuotes(templatesCont.S("name").String())
		map_template["display_name"] = models.StripQuotes(templatesCont.S("displayName").String())
		map_template["tenant_id"] = models.StripQuotes(templatesCont.S("tenantId").String())
		templates = append(templates, map_template)

	}
	d.Set("template", templates)
	d.Set("template_name", "")
	d.Set("tenant_id", "")

	log.Printf("[DEBUG] %s: Schema Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema: Beginning Update")
	msoClient := m.(*client.Client)
	name := d.Get("name").(string)
	if _, ok := d.GetOk("template_name"); ok {
		old, new := d.GetChange("template_name")
		oldTemplate := old.(string)
		newTemplate := new.(string)
		if d.HasChange("tenant_id") {
			return fmt.Errorf("Tenant associated with Template cannot be changed.")
		}
		schemaNamePayload := fmt.Sprintf(`
	{ 
		"op": "replace",
		"path": "/displayName",
		"value": "%s"
	}
`, name)

		templateNamePayload := fmt.Sprintf(`
	{
		"op": "replace",
		"path": "/templates/%s/name",
		"value": "%s"
	}
`, oldTemplate, newTemplate)

		tempDisplayNamePayload := fmt.Sprintf(`
	{
		"op": "replace",
		"path": "/templates/%s/displayName",
		"value": "%s"
	}
`, newTemplate, newTemplate)

		jsonSchema, err := container.ParseJSON([]byte(schemaNamePayload))
		jsonTemplate, err := container.ParseJSON([]byte(templateNamePayload))
		jsonDispl, err := container.ParseJSON([]byte(tempDisplayNamePayload))
		payloadCon := container.New()

		payloadCon.Array()
		err = payloadCon.ArrayAppend(jsonSchema.Data())
		if err != nil {
			return err
		}
		payloadCon.ArrayAppend(jsonTemplate.Data())
		payloadCon.ArrayAppend(jsonDispl.Data())
		path := fmt.Sprintf("api/v1/schemas/%s", d.Id())

		req, err := msoClient.MakeRestRequest("PATCH", path, payloadCon, true)
		if err != nil {
			return err
		}
		cont, _, err := msoClient.Do(req)
		if err != nil {
			return err
		}

		err = client.CheckForErrors(cont, "PATCH")
		if err != nil {
			return err
		}
	} else {
		if d.HasChange("template") {
			old_templates, new_templates := d.GetChange("template")
			platform := msoClient.GetPlatform()

			//In non-ND based MSOs, tenant cannnot be changed by PATCH API, so we delete it explicitly.
			if platform != "nd" {
				//Get all the new maps
				getDifferenceNew := differenceInMaps(new_templates.(*schema.Set), old_templates.(*schema.Set))

				// Get old maps that have a change
				getDifferenceOld := differenceInMaps(old_templates.(*schema.Set), new_templates.(*schema.Set))

				for _, valueMapOld := range getDifferenceOld {
					valueOld := valueMapOld.(map[string]interface{})
					for _, valueMapNew := range getDifferenceNew {
						valueNew := valueMapNew.(map[string]interface{})

						//Tenant ID has been changed. Delete Template.
						if valueOld["name"] == valueNew["name"] && valueOld["display_name"] == valueNew["display_name"] {
							deleteTemplate(valueOld["name"].(string), valueOld["display_name"].(string), valueOld["tenant_id"].(string), d, m)
						}
					}
				}
			}

			templates := make([]interface{}, 0, 1)
			if val, ok := d.GetOk("template"); ok {
				template_list := val.(*schema.Set).List()
				for _, val := range template_list {
					map_templates := make(map[string]interface{})
					inner_templates := val.(map[string]interface{})
					map_templates["name"] = inner_templates["name"]
					map_templates["displayName"] = inner_templates["display_name"]
					map_templates["tenantId"] = inner_templates["tenant_id"]
					templates = append(templates, map_templates)
				}
			}
			dn := d.Id()
			_, schemaApp := models.NewSchema(dn, name, "", "", templates)
			_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", dn), schemaApp)
			if err != nil {
				return err
			}
		}
	}
	log.Printf("[DEBUG] %s: Schema Update finished successfully", d.Id())
	return resourceMSOSchemaRead(d, m)
}

func resourceMSOSchemaRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	dn := d.Id()
	con, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/" + dn))
	if err != nil {
		return err
	}

	d.SetId(models.StripQuotes(con.S("id").String()))
	d.Set("name", models.StripQuotes(con.S("displayName").String()))
	count, err := con.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No Template found")
	}
	stateTemplate := d.Get("template_name").(string)
	stateTenant := d.Get("tenant_id").(string)
	templates := make([]interface{}, 0)
	for i := 0; i < count; i++ {
		templatesCont, err := con.ArrayElement(i, "templates")
		if err != nil {
			return fmt.Errorf("Unable to parse the templates list")
		}
		map_template := make(map[string]interface{})
		map_template["name"] = models.StripQuotes(templatesCont.S("name").String())
		map_template["display_name"] = models.StripQuotes(templatesCont.S("displayName").String())
		map_template["tenant_id"] = models.StripQuotes(templatesCont.S("tenantId").String())
		templates = append(templates, map_template)

		apiTemplate := models.StripQuotes(templatesCont.S("name").String())
		apiTenant := models.StripQuotes(templatesCont.S("tenantId").String())
		if apiTemplate == stateTemplate && apiTenant == stateTenant {
			d.Set("template_name", apiTemplate)
			d.Set("tenant_id", apiTenant)
		} else {
			d.Set("template_name", "")
			d.Set("tenant_id", "")
		}
	}
	d.Set("template", templates)
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSOSchemaDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	msoClient := m.(*client.Client)
	dn := d.Id()
	err := msoClient.DeletebyId("api/v1/schemas/" + dn)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

//Helper function 1
func deleteTemplate(name, tenantId, displayName string, d *schema.ResourceData, m interface{}) error {
	msoClient := m.(*client.Client)
	schematemplateDelete := models.NewSchemaTemplate("remove", fmt.Sprintf("/templates/%s", name), tenantId, name, displayName)
	response, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", d.Id()), schematemplateDelete)
	if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return err
	}
	return nil
}

//Helper function 2
func differenceInMaps(mapSlice1, mapSlice2 *schema.Set) []interface{} {
	var difference []interface{}
	for i := 0; i < 1; i++ {
		for _, s1 := range mapSlice1.List() {
			found := false
			for _, s2 := range mapSlice2.List() {
				if reflect.DeepEqual(s1, s2) {
					found = true
					break
				}
			}
			if !found {
				difference = append(difference, s1)
			}
		}
		if i == 0 {
			mapSlice1, mapSlice2 = mapSlice2, mapSlice1
		}
	}
	return difference
}
