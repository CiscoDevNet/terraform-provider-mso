package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceMSOUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOUserCreate,
		Update: resourceMSOUserUpdate,
		Read:   resourceMSOUserRead,
		Delete: resourceMSOUserDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOUserImport,
		},

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"user_password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"first_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"last_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"email": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"phone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"roles": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"roleid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"access_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Required: true,
			},
		}),
	}
}

func resourceMSOUserImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] User: Beginning Import")

	msoClient := m.(*client.Client)
	platform := msoClient.GetPlatform()
	if platform == "nd" {
		return nil, fmt.Errorf("The mso_user resources is not supported on ND-based MSO/NDO. Use ND provider for manipulating users on ND-based MSO/NDO.")
	}
	con, err := msoClient.GetViaURL("api/v1/users" + d.Id())
	if err != nil {
		return nil, err
	}

	d.SetId(models.StripQuotes(con.S("id").String()))
	d.Set("username", models.StripQuotes(con.S("username").String()))
	d.Set("user_password", models.StripQuotes(con.S("password").String()))
	if con.Exists("firstName") {
		d.Set("first_name", models.StripQuotes(con.S("firstName").String()))
	}
	if con.Exists("lastName") {
		d.Set("last_name", models.StripQuotes(con.S("lastName").String()))
	}
	if con.Exists("emailAddress") {
		d.Set("email", models.StripQuotes(con.S("emailAddress").String()))
	}
	if con.Exists("phoneNumber") {
		d.Set("phone", models.StripQuotes(con.S("phoneNumber").String()))
	}
	if con.Exists("accountStatus") {
		d.Set("account_status", models.StripQuotes(con.S("accountStatus").String()))
	}
	if con.Exists("domain") {
		d.Set("domain", models.StripQuotes(con.S("domain").String()))
	}

	count, err := con.ArrayCount("roles")

	if err != nil {
		return nil, fmt.Errorf("No Roles found")
	}

	roles := make([]interface{}, 0)
	for i := 0; i < count; i++ {
		rolesCont, err := con.ArrayElement(i, "roles")

		if err != nil {
			return nil, fmt.Errorf("Unable to parse the roles list")
		}

		map1 := make(map[string]interface{})

		map1["roleid"] = models.StripQuotes(rolesCont.S("roleId").String())
		map1["access_type"] = models.StripQuotes(rolesCont.S("accessType").String())
		roles = append(roles, map1)
	}

	d.Set("roles", roles)

	log.Printf("[DEBUG] %s: User Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOUserCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] User: Beginning Creation")
	msoClient := m.(*client.Client)
	platform := msoClient.GetPlatform()
	if platform == "nd" {
		return fmt.Errorf("The mso_user resources is not supported on ND-based MSO/NDO. Use ND provider for manipulating users on ND-based MSO/NDO.")
	}

	var user string
	if username, ok := d.GetOk("username"); ok {
		user = username.(string)
	}

	var userPassword string
	if password, ok := d.GetOk("user_password"); ok {
		userPassword = password.(string)
	}

	var firstName string
	if firstname, ok := d.GetOk("first_name"); ok {
		firstName = firstname.(string)
	}

	var lastName string
	if lastname, ok := d.GetOk("last_name"); ok {
		lastName = lastname.(string)
	}

	var email string
	if emails, ok := d.GetOk("email"); ok {
		email = emails.(string)
	}

	var phone string
	if phones, ok := d.GetOk("phone"); ok {
		phone = phones.(string)
	}

	var accountStatus string
	if accountstatus, ok := d.GetOk("account_status"); ok {
		accountStatus = accountstatus.(string)
	}

	var domain string
	if Domain, ok := d.GetOk("domain"); ok {
		domain = Domain.(string)
	}

	roles := make([]interface{}, 0, 1)
	if val, ok := d.GetOk("roles"); ok {
		role_list := val.(*schema.Set).List()
		for _, val := range role_list {

			map_roles := make(map[string]interface{})
			inner_roles := val.(map[string]interface{})
			if inner_roles["roleid"] != "" {
				map_roles["roleId"] = fmt.Sprintf("%v", inner_roles["roleid"])
			}
			if inner_roles["access_type"] != "" {
				map_roles["accessType"] = fmt.Sprintf("%v", inner_roles["access_type"])
			}

			roles = append(roles, map_roles)
		}
	}

	userApp := models.NewUser("", user, userPassword, firstName, lastName, email, phone, accountStatus, domain, roles)
	cont, err := msoClient.Save("api/v1/users", userApp)
	if err != nil {
		return err
	}

	id := models.StripQuotes(cont.S("id").String())
	d.SetId(fmt.Sprintf("%v", id))
	log.Printf("[DEBUG] %s: User Creation finished successfully", d.Id())

	return resourceMSOUserRead(d, m)
}

func resourceMSOUserUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] User: Beginning Creation of resource")
	msoClient := m.(*client.Client)

	platform := msoClient.GetPlatform()
	if platform == "nd" {
		return fmt.Errorf("The mso_user resources is not supported on ND-based MSO/NDO. Use ND provider for manipulating users on ND-based MSO/NDO.")
	}

	var user string
	if username, ok := d.GetOk("username"); ok {
		user = username.(string)
	}

	var userPassword string
	if password, ok := d.GetOk("user_password"); ok {
		userPassword = password.(string)
	}

	var firstName string
	if firstname, ok := d.GetOk("first_name"); ok {
		firstName = firstname.(string)
	}

	var lastName string
	if lastname, ok := d.GetOk("last_name"); ok {
		lastName = lastname.(string)
	}

	var email string
	if emails, ok := d.GetOk("email"); ok {
		email = emails.(string)
	}

	var phone string
	if phones, ok := d.GetOk("phone"); ok {
		phone = phones.(string)
	}

	var accountStatus string
	if accountstatus, ok := d.GetOk("account_status"); ok {
		accountStatus = accountstatus.(string)
	}

	var domain string
	if Domain, ok := d.GetOk("domain"); ok {
		domain = Domain.(string)
	}

	roles := make([]interface{}, 0, 1)
	if val, ok := d.GetOk("roles"); ok {

		role_list := val.(*schema.Set).List()
		for _, val := range role_list {

			map_role := make(map[string]interface{})
			inner_role := val.(map[string]interface{})
			map_role["roleId"] = fmt.Sprintf("%v", inner_role["roleid"])

			if inner_role["access_type"] != "" {
				map_role["accessType"] = fmt.Sprintf("%v", inner_role["access_type"])
			}

			roles = append(roles, map_role)
		}

	}

	userApp := models.NewUser("", user, userPassword, firstName, lastName, email, phone, accountStatus, domain, roles)

	cont, err := msoClient.Put(fmt.Sprintf("api/v1/users/%s", d.Id()), userApp)
	if err != nil {
		return err
	}

	id := models.StripQuotes(cont.S("id").String())
	d.SetId(fmt.Sprintf("%v", id))
	log.Printf("[DEBUG] %s: Schema Creation finished successfully", d.Id())

	return resourceMSOUserRead(d, m)
	return nil

}

func resourceMSOUserRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	platform := msoClient.GetPlatform()
	if platform == "nd" {
		return fmt.Errorf("The mso_user resources is not supported on ND-based MSO/NDO. Use ND provider for manipulating users on ND-based MSO/NDO.")
	}

	dn := d.Id()
	con, err := msoClient.GetViaURL("api/v1/users/" + dn)

	if err != nil {
		return errorForObjectNotFound(err, dn, con, d)
	}

	d.SetId(models.StripQuotes(con.S("id").String()))
	d.Set("username", models.StripQuotes(con.S("username").String()))
	//d.Set("user_password", models.StripQuotes(con.S("password").String()))
	if con.Exists("firstName") {
		d.Set("first_name", models.StripQuotes(con.S("firstName").String()))
	}
	if con.Exists("lastName") {
		d.Set("last_name", models.StripQuotes(con.S("lastName").String()))
	}
	if con.Exists("emailAddress") {
		d.Set("email", models.StripQuotes(con.S("emailAddress").String()))
	}
	if con.Exists("phoneNumber") {
		d.Set("phone", models.StripQuotes(con.S("phoneNumber").String()))
	}
	if con.Exists("accountStatus") {
		d.Set("account_status", models.StripQuotes(con.S("accountStatus").String()))
	}
	if con.Exists("domain") {
		d.Set("domain", models.StripQuotes(con.S("domain").String()))
	}
	count, err := con.ArrayCount("roles")
	if err != nil {
		return fmt.Errorf("No Roles found")
	}

	roles := make([]interface{}, 0)
	for i := 0; i < count; i++ {
		rolesCont, err := con.ArrayElement(i, "roles")

		if err != nil {
			return fmt.Errorf("Unable to parse the roles list")
		}

		map_role := make(map[string]interface{})

		map_role["roleid"] = models.StripQuotes(rolesCont.S("roleId").String())
		map_role["access_type"] = models.StripQuotes(rolesCont.S("accessType").String())
		roles = append(roles, map_role)
	}
	d.Set("roles", roles)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSOUserDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	msoClient := m.(*client.Client)

	platform := msoClient.GetPlatform()
	if platform == "nd" {
		return fmt.Errorf("The mso_user resources is not supported on ND-based MSO/NDO. Use ND provider for manipulating users on ND-based MSO/NDO.")
	}

	dn := d.Id()
	err := msoClient.DeletebyId("api/v1/users/" + dn)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

// func toStringList(configured interface{}) []string {
// 	vs := make([]string, 0, 1)
// 	val, ok := configured.(string)
// 	if ok && val != "" {
// 		vs = append(vs, val)
// 	}
// 	return vs
// }
