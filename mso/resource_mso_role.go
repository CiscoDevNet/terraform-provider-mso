package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceMSORole() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSORoleCreate,
		Update: resourceMSORoleUpdate,
		Read:   resourceMSORoleRead,
		Delete: resourceMSORoleDelete,

		SchemaVersion: 1,

		Schema: (map[string]*schema.Schema{
			"role_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"read_permissions": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"write_permissions": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		}),
	}
}

func resourceMSORoleCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema: Beginning Creation")
	msoClient := m.(*client.Client)
	siteAttr := models.RoleAttributes{}

	siteAttr.Id = ""

	if name, ok := d.GetOk("name"); ok {
		siteAttr.Name = name.(string)
	}

	if disp, ok := d.GetOk("display_name"); ok {
		siteAttr.DisplayName = disp.(string)
	}

	if desc, ok := d.GetOk("description"); ok {
		siteAttr.Description = desc.(string)
	}

	if readp, ok := d.GetOk("read_permissions"); ok {
		siteAttr.ReadPermissions = readp.([]interface{})
	}

	if writep, ok := d.GetOk("write_permissions"); ok {
		siteAttr.WritePermissions = writep.([]interface{})
	}

	roleApp := models.NewRole(siteAttr)

	cont, err := msoClient.Save("api/v1/roles", roleApp)
	if err != nil {
		log.Println(err)
		return err
	}

	id := models.StripQuotes(cont.S("id").String())
	d.SetId(fmt.Sprintf("%v", id))
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceMSORoleRead(d, m)

}

func resourceMSORoleUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema: Beginning Creation")
	msoClient := m.(*client.Client)
	siteAttr := models.RoleAttributes{}

	id1 := d.Id()
	siteAttr.Id = id1

	if name, ok := d.GetOk("name"); ok {
		siteAttr.Name = name.(string)
	}

	if disp, ok := d.GetOk("display_name"); ok {
		siteAttr.DisplayName = disp.(string)
	}

	if desc, ok := d.GetOk("description"); ok {
		siteAttr.Description = desc.(string)
	}

	if readp, ok := d.GetOk("read_permissions"); ok {
		siteAttr.ReadPermissions = readp.([]interface{})
	}

	if writep, ok := d.GetOk("write_permissions"); ok {
		siteAttr.WritePermissions = writep.([]interface{})
	}

	roleApp := models.NewRole(siteAttr)

	cont, err := msoClient.Put(fmt.Sprintf("api/v1/roles/%s", d.Id()), roleApp)
	if err != nil {
		return err
	}

	id := models.StripQuotes(cont.S("id").String())
	d.SetId(fmt.Sprintf("%v", id))
	log.Printf("[DEBUG] %s: Schema Creation finished successfully", d.Id())

	return resourceMSORoleRead(d, m)
	return nil

}

func resourceMSORoleRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	msoClient := m.(*client.Client)

	dn := d.Id()

	con, err := msoClient.GetViaURL("api/v1/roles/" + dn)
	if err != nil {
		return err
	}

	d.SetId(models.StripQuotes(con.S("id").String()))
	d.Set("name", models.StripQuotes(con.S("name").String()))
	d.Set("display_name", models.StripQuotes(con.S("displayName").String()))
	d.Set("description", models.StripQuotes(con.S("description").String()))

	count1, err := con.ArrayCount("readPermissions")
	if err != nil {
		return fmt.Errorf("No Read Permission found")
	}
	found1 := false
	for i := 0; i < count1; i++ {

		temp, err := con.ArrayElement(i, "readPermissions")
		log.Println(temp)
		d.Set("read_permissions", temp)
		found1 = true
		if err != nil {
			return fmt.Errorf("Unable to parse the read permissions list")
		}
	}
	if !found1 {
		d.Set("read_permissions", "")
	}

	count2, err := con.ArrayCount("writePermissions")
	if err != nil {
		return fmt.Errorf("No write permission found")
	}
	found2 := false
	for i := 0; i < count2; i++ {

		temp, err := con.ArrayElement(i, "writePermissions")
		log.Println(temp)
		d.Set("write_permissions", temp)
		found2 = true
		if err != nil {
			return fmt.Errorf("Unable to parse the write permissions list")
		}
	}
	if !found2 {
		d.Set("write_permissions", "")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSORoleDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	msoClient := m.(*client.Client)
	dn := d.Id()
	err := msoClient.DeletebyId("api/v1/roles/" + dn)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
