package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
	//"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceMSOSchema() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaCreate,
		Update: resourceMSOSchemaUpdate,
		Read:   resourceMSOSchemaRead,
		Delete: resourceMSOSchemaDelete,

		// Importer: &schema.ResourceImporter{
		// 	State: resourceMSOSchemaImport,
		// },

		SchemaVersion: 1,

		Schema: (map[string]*schema.Schema{
			"schema": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"templates": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				
			},

			"sites": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				//Computed: true,
			},

		}),
	}
}
// func getRemoteCloudApplicationcontainer(client *client.Client, dn string) (*models.CloudApplicationcontainer, error) {
// 	cloudAppCont, err := client.Get(dn)
// 	if err != nil {
// 		return nil, err
// 	}

// 	cloudApp := models.CloudApplicationcontainerFromContainer(cloudAppCont)

// 	if cloudApp.DistinguishedName == "" {
// 		return nil, fmt.Errorf("CloudApplicationcontainer %s not found", cloudApp.DistinguishedName)
// 	}

// 	return cloudApp, nil
// }

// func setCloudApplicationcontainerAttributes(cloudApp *models.CloudApplicationcontainer, d *schema.ResourceData) *schema.ResourceData {
// 	d.SetId(cloudApp.DistinguishedName)
// 	d.Set("description", cloudApp.Description)
// 	d.Set("tenant_dn", GetParentDn(cloudApp.DistinguishedName))
// 	cloudAppMap, _ := cloudApp.ToMap()

// 	d.Set("name", cloudAppMap["name"])

// 	d.Set("annotation", cloudAppMap["annotation"])
// 	d.Set("name_alias", cloudAppMap["nameAlias"])
// 	return d
// }

// func resourceMSOSchemaImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
// 	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
// 	aciClient := m.(*client.Client)

// 	dn := d.Id()

// 	cloudApp, err := getRemoteCloudApplicationcontainer(aciClient, dn)

// 	if err != nil {
// 		return nil, err
// 	}
// 	schemaFilled := setCloudApplicationcontainerAttributes(cloudApp, d)

// 	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

// 	return []*schema.ResourceData{schemaFilled}, nil
// }

func resourceMSOSchemaCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema: Beginning Creation")
	msoClient := m.(*client.Client)
	schemaAttr := models.SchemaAttributes{}
	if schema,ok := d.GetOk("schema");ok{
		schemaAttr.Schema=schema.(string)

	}
    if templates,ok := d.GetOk("templates");ok{
		templateList := toStringList(templates.([]interface{}))
		schemaAttr.Templates=templateList
	}
	
	if sites, ok := d.GetOk("sites"); ok {
		siteList := toStringList(sites.([]interface{}))
		schemaAttr.Sites=siteList
	
	}
	schemaApp := models.NewSchemacontainer(schemaAttr)

	cont,err := msoClient.Save("https://173.36.219.193/api/v1/schemas",schemaApp)
	if err != nil {
		return err
	}

	id:=cont.S("id")
	log.Println("Id value", id)
	d.SetId(fmt.Sprintf("%v",id))
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceMSOSchemaRead(d, m)
}

func resourceMSOSchemaUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] CloudApplicationcontainer: Beginning Update")

	msoClient := m.(*client.Client)

	schemaAttr := models.SchemaAttributes{}

	if d.HasChange("schema") {
		schemaAttr.Schema = d.Get("schema").(string)
	}

	if d.HasChange("templates") {
		if templates,ok := d.GetOk("templates");ok{
			templateList := toStringList(templates.([]interface{}))
			schemaAttr.Templates=templateList
		}
	}

	if d.HasChange("sites") {
		if sites, ok := d.GetOk("sites"); ok {
			siteList := toStringList(sites.([]interface{}))
			schemaAttr.Sites=siteList
		
		}
	}
   schemaApp := models.NewSchemacontainer(schemaAttr)
  cont,err := msoClient.PatchbyID("https://173.36.219.193/api/v1/schemas/"+d.Id(),schemaApp)

	if err != nil {
		return err
	}
	
	id:=cont.S("id")
	log.Println("Id value", id)
	d.SetId(fmt.Sprintf("%v",id))
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceMSOSchemaRead(d, m)

}

func resourceMSOSchemaRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	dn := d.Id()
	
   con, err := msoClient.GetViaURL("https://173.36.219.193/api/v1/schemas/" + dn)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%v", con.S("id")))
	d.Set("schema", con.S("schema"))
	d.Set("templates", con.S("templates").String())
	d.Set("sites", con.S("sites").String())
	
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSOSchemaDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	msoClient := m.(*client.Client)
	dn := d.Id()
	err := msoClient.DeletebyId("https://173.36.219.193/api/v1/schemas/" + dn)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}


func toStringList(configured interface{}) []string {
    vs := make([]string, 0, 1)
    val, ok := configured.(string)
    if ok && val != "" {
        vs = append(vs, val)
    }
    return vs
}