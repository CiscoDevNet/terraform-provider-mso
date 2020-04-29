package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceMSOTenant() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOTenantCreate,
		Update: resourceMSOTenantUpdate,
		Read:   resourceMSOTenantRead,
		Delete: resourceMSOTenantDelete,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{

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

			"user_associations": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Computed: true,
			},

			"site_associations": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"site_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Computed: true,
			},
		}),
	}
}

func resourceMSOTenantCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Tenant: Beginning Creation")
	msoClient := m.(*client.Client)
	tenantAttr := models.TenantAttributes{}

	if name, ok := d.GetOk("name"); ok {
		tenantAttr.Name = name.(string)
	}

	if display_name, ok := d.GetOk("display_name"); ok {
		tenantAttr.DisplayName = display_name.(string)
	}

	if description, ok := d.GetOk("description"); ok {
		tenantAttr.Description = description.(string)
	}

	site_associations := make([]interface{}, 0, 1)
	if val, ok := d.GetOk("site_associations"); ok {
		siteList := val.(*schema.Set).List()
		for _, val := range siteList {

			mapSite := make(map[string]interface{})
			inner := val.(map[string]interface{})
			if inner["siteId"] != "" {
				mapSite["site_id"] = fmt.Sprintf("%v", inner["site_id"])
			}
			site_associations = append(site_associations, mapSite)
		}
	}
	tenantAttr.Sites = site_associations

	user_associations := make([]interface{}, 0, 1)
	if val, ok := d.GetOk("user_associations"); ok {
		userList := val.(*schema.Set).List()
		for _, val := range userList {

			mapUser := make(map[string]interface{})
			inner := val.(map[string]interface{})
			if inner["userId"] != "" {
				mapUser["user_id"] = fmt.Sprintf("%v", inner["user_id"])
			}
			user_associations = append(user_associations, mapUser)
		}
	}
	tenantAttr.Users = user_associations

	tenantApp := models.NewTenant(tenantAttr)

	cont, err := msoClient.Save("api/v1/tenants", tenantApp)
	if err != nil {
		log.Println(err)
		return err
	}

	id := models.StripQuotes(cont.S("id").String())
	d.SetId(fmt.Sprintf("%v", id))
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceMSOTenantRead(d, m)
}

func resourceMSOTenantUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Tenant: Beginning Update")

	msoClient := m.(*client.Client)

	tenantAttr := models.TenantAttributes{}

	if d.HasChange("name") {
		tenantAttr.Name = d.Get("name").(string)
	}

	if d.HasChange("display_name") {
		tenantAttr.DisplayName = d.Get("display_name").(string)
	}

	if d.HasChange("description") {
		tenantAttr.Description = d.Get("description").(string)
	}

	site_associations := make([]interface{}, 0, 1)
	if val, ok := d.GetOk("site_associations"); ok {

		siteList := val.(*schema.Set).List()
		for _, val := range siteList {

			mapSite := make(map[string]interface{})
			inner := val.(map[string]interface{})
			mapSite["userId"] = fmt.Sprintf("%v", inner["site_id"])
			site_associations = append(site_associations, mapSite)
		}
	}
	tenantAttr.Sites = site_associations

	user_associations := make([]interface{}, 0, 1)
	if val, ok := d.GetOk("user_associations"); ok {

		userList := val.(*schema.Set).List()
		for _, val := range userList {

			mapUser := make(map[string]interface{})
			inner := val.(map[string]interface{})
			mapUser["userId"] = fmt.Sprintf("%v", inner["user_id"])
			user_associations = append(user_associations, mapUser)
		}
	}
	tenantAttr.Users = user_associations

	tenantApp := models.NewTenant(tenantAttr)
	cont, err := msoClient.Put(fmt.Sprintf("api/v1/tenants/%s", d.Id()), tenantApp)
	if err != nil {
		return err
	}

	id := models.StripQuotes(cont.S("id").String())
	d.SetId(fmt.Sprintf("%v", id))
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceMSOTenantRead(d, m)
}

func resourceMSOTenantRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	msoClient := m.(*client.Client)

	dn := d.Id()

	con, err := msoClient.GetViaURL("api/v1/tenants/" + dn)
	if err != nil {
		return err
	}

	d.SetId(models.StripQuotes(con.S("id").String()))
	d.Set("name", models.StripQuotes(con.S("name").String()))
	d.Set("display_name", models.StripQuotes(con.S("displayName").String()))
	d.Set("description", models.StripQuotes(con.S("description").String()))

	count1, _ := con.ArrayCount("siteAssociations")
	site_associations := make([]interface{}, 0)
	for i := 0; i < count1; i++ {
		sitesCont, err := con.ArrayElement(i, "siteAssociations")
		if err != nil {
			return fmt.Errorf("Unable to parse the site associations list")
		}

		mapSite := make(map[string]interface{})
		mapSite["site_id"] = models.StripQuotes(sitesCont.S("siteId").String())
		site_associations = append(site_associations, mapSite)
	}

	d.Set("site_associations", site_associations)

	count2, err := con.ArrayCount("userAssociations")
	if err != nil {
		d.Set("user_assocoations", make([]interface{}, 0))
	}

	user_associations := make([]interface{}, 0)
	for i := 0; i < count2; i++ {
		usersCont, err := con.ArrayElement(i, "userAssociations")
		if err != nil {
			return fmt.Errorf("Unable to parse the user associations list")
		}

		mapUser := make(map[string]interface{})
		mapUser["user_id"] = models.StripQuotes(usersCont.S("userId").String())
		user_associations = append(user_associations, mapUser)
	}

	d.Set("user_associations", user_associations)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSOTenantDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	msoClient := m.(*client.Client)
	dn := d.Id()
	err := msoClient.DeletebyId(fmt.Sprintf("api/v1/tenants/%v", dn))
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
