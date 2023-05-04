package mso

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

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
				Deprecated:   "use template block with name, display_name and tenant_id instead",
			},
			"tenant_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
				Deprecated:   "use template block with name, display_name and tenant_id instead",
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
	tempVarTemplateName, ok_template_name := d.GetOk("template_name")
	tempVarTemplates, ok_templates := d.GetOk("template")

	if !ok_template_name && !ok_templates {
		return fmt.Errorf("template_name or a template block with its name, tenant_id and display_name are required.")
	}

	var schemaApp *models.Schema
	if ok_template_name {
		tempVarTenantId, ok := d.GetOk("tenant_id")
		if !ok {
			return fmt.Errorf("tenant_id is required when using template_name.")
		}
		templateName := tempVarTemplateName.(string)
		tenantId := tempVarTenantId.(string)
		schemaApp = models.NewSchema("", name, templateName, tenantId, make([]interface{}, 0, 1))

	} else {
		templates := make([]interface{}, 0, 1)
		if ok_templates {
			template_list := tempVarTemplates.(*schema.Set).List()
			for _, val := range template_list {
				map_templates := make(map[string]interface{})
				inner_templates := val.(map[string]interface{})
				map_templates["name"] = inner_templates["name"]
				map_templates["displayName"] = inner_templates["display_name"]
				map_templates["tenantId"] = inner_templates["tenant_id"]
				templates = append(templates, map_templates)
			}
		}
		schemaApp = models.NewSchema("", name, "", "", templates)
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
	/* When importing a schema with a single template, there is no way of knowing which template format(single or block) the user is expecting to be populated. Since template_name and tenant_id are deprecated, and are going to be removed in a future release,
	   template_name and tenant_id are set to "" in the import function. */
	d.Set("template_name", "")
	d.Set("tenant_id", "")

	log.Printf("[DEBUG] %s: Schema Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema: Beginning Update")
	msoClient := m.(*client.Client)
	name := d.Get("name").(string)
	_, ok_template_name := d.GetOk("template_name")
	_, ok_templates := d.GetOk("template")

	if !ok_template_name && !ok_templates {
		return fmt.Errorf("template_name or a template block with its name, tenant_id and display_name are required.")
	} else if ok_template_name {
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
		`, oldTemplate, newTemplate)

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
		listAttributesToChange := make([]string, 0)
		if d.HasChange("name") {
			listAttributesToChange = append(listAttributesToChange, fmt.Sprintf(`
			{ 
				"op": "replace",
				"path": "/displayName",
				"value": "%s"
			}
		`, name))
		}
		if d.HasChange("template") {
			// This keeps a track of new maps whose values have been changed (new values)
			listMapsReplaced := make([]interface{}, 0)

			// This keeps a track of old maps whose values will be changed (old values)
			listMapsToReplace := make([]interface{}, 0)

			old_templates, new_templates := d.GetChange("template")

			//Get all the new maps
			getDifferenceNew := differenceInMaps(new_templates.(*schema.Set), old_templates.(*schema.Set))

			// Get old maps that have a change
			getDifferenceOld := differenceInMaps(old_templates.(*schema.Set), new_templates.(*schema.Set))

			for _, valueMapOld := range getDifferenceOld {
				valueOld := valueMapOld.(map[string]interface{})
				for _, valueMapNew := range getDifferenceNew {
					valueNew := valueMapNew.(map[string]interface{})
					// Tenant Id of template has been changed
					if valueOld["name"] == valueNew["name"] && valueOld["tenant_id"] != valueNew["tenant_id"] {
						listMapsReplaced = append(listMapsReplaced, valueNew)
						listMapsToReplace = append(listMapsToReplace, valueOld)
						listAttributesToChange = append(listAttributesToChange, fmt.Sprintf(`
							{
								"op": "replace",
								"path": "/templates/%s/tenantId",
								"value": "%s"
							}
					`, valueOld["name"].(string), valueNew["tenant_id"].(string)))
					}
					// Display name of template has been changed
					if valueOld["name"] == valueNew["name"] && valueOld["display_name"] != valueNew["display_name"] {
						listMapsReplaced = append(listMapsReplaced, valueNew)
						listMapsToReplace = append(listMapsToReplace, valueOld)
						listAttributesToChange = append(listAttributesToChange, fmt.Sprintf(`
							{
								"op": "replace",
								"path": "/templates/%s/displayName",
								"value": "%s"
							}
					`, valueOld["name"].(string), valueNew["display_name"].(string)))
					}
					// Name of template has been changed
					if valueOld["name"] != valueNew["name"] && valueOld["display_name"] == valueNew["display_name"] && valueOld["tenant_id"] == valueNew["tenant_id"] {
						listMapsReplaced = append(listMapsReplaced, valueNew)
						listMapsToReplace = append(listMapsToReplace, valueOld)
						listAttributesToChange = append(listAttributesToChange, fmt.Sprintf(`
							{
								"op": "replace",
								"path": "/templates/%s/name",
								"value": "%s"
							}
						`, valueOld["name"].(string), valueNew["name"].(string)))

					}
				}
			}

			// New templates have been added to the block.
			listMapsToAdd := differenceInLists(getDifferenceNew, listMapsReplaced)
			for _, MapToAdd := range listMapsToAdd {
				map_add, _ := json.Marshal(MapToAdd)
				map_values := strings.Replace(strings.Replace(string(map_add), "display_name", "displayName", 1), "tenant_id", "tenantId", 1)
				listAttributesToChange = append(listAttributesToChange, fmt.Sprintf(`
							{
								"op": "add",
								"path": "/templates/-",
								"value": %s
							}
						`, map_values))
			}

			// templates have been removed from the block
			listMapsToRemove := differenceInLists(getDifferenceOld, listMapsToReplace)
			for _, MapToRemove := range listMapsToRemove {
				valueRemove := MapToRemove.(map[string]interface{})
				map_remove, _ := json.Marshal(valueRemove)
				map_values := strings.Replace(strings.Replace(string(map_remove), "display_name", "displayName", 1), "tenant_id", "tenantId", 1)
				listAttributesToChange = append(listAttributesToChange, fmt.Sprintf(`
							{
								"op": "remove",
								"path": "/templates/%s",
								"value": %s
							}
						`, valueRemove["name"], map_values))
			}

		}

		// Construction of complete payload for PATCH
		if len(listAttributesToChange) != 0 {
			payloadCon := container.New()
			payloadCon.Array()
			jsonAttributes, err := container.ParseJSON([]byte(fmt.Sprintf(`[` + strings.Join(listAttributesToChange, ",") + `]`)))
			if err != nil {
				return err
			}
			payloadCon.ArrayAppend(jsonAttributes.Data())

			path := fmt.Sprintf("api/v1/schemas/%s", d.Id())

			req, err := msoClient.MakeRestRequest("PATCH", path, payloadCon.Index(0), true)
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
		if con.S("code").String() == "404" {
			d.SetId("")
			return nil
		} else {
			return err
		}
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
		}
	}
	if _, ok := d.GetOk("template_name"); !ok {
		d.Set("template", templates)
		d.Set("template_name", "")
		d.Set("tenant_id", "")
	}
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

// Helper function 1 for sets
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

// Helper function 2 for lists
func differenceInLists(mapSlice1, mapSlice2 []interface{}) []interface{} {
	var difference []interface{}
	for i := 0; i < 1; i++ {
		for _, s1 := range mapSlice1 {
			found := false
			for _, s2 := range mapSlice2 {
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
