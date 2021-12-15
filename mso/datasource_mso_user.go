package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSOUser() *schema.Resource {
	return &schema.Resource{

		Read: datasourceMSOUserRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"user_password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
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
				Optional: true,
			},
		}),
	}
}

func datasourceMSOUserRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	username := d.Get("username").(string)
	var path string
	platform := msoClient.GetPlatform()
	if platform == "nd" {
		path = "api/v2/users"
	} else {
		path = "api/v1/users"
	}
	con, err := msoClient.GetViaURL(path)
	if err != nil {
		return err
	}

	var data []interface{}
	var usernameKey string
	if platform == "nd" {
		data = con.Data().([]interface{})
		usernameKey = "loginID"
	} else {
		data = con.S("users").Data().([]interface{})
		usernameKey = "username"
	}

	var flag bool
	var cnt int
	for _, info := range data {
		val := info.(map[string]interface{})
		if val[usernameKey].(string) == username {
			flag = true
			break
		}
		cnt = cnt + 1
	}
	if flag != true {
		return fmt.Errorf("User of specified name not found")
	}

	var dataCon *container.Container
	if platform == "nd" {
		dataCon = con.Index(cnt)
		d.SetId(models.StripQuotes(dataCon.S("userID").String()))
	} else {
		dataCon = con.S("users").Index(cnt)
		d.SetId(models.StripQuotes(dataCon.S("id").String()))
	}

	d.Set("username", models.StripQuotes(dataCon.S(usernameKey).String()))
	d.Set("user_password", models.StripQuotes(dataCon.S("password").String()))
	if dataCon.Exists("firstName") {
		d.Set("first_name", models.StripQuotes(dataCon.S("firstName").String()))
	}
	if dataCon.Exists("lastName") {
		d.Set("last_name", models.StripQuotes(dataCon.S("lastName").String()))
	}
	if dataCon.Exists("emailAddress") {
		d.Set("email", models.StripQuotes(dataCon.S("emailAddress").String()))
	}
	if dataCon.Exists("phoneNumber") {
		d.Set("phone", models.StripQuotes(dataCon.S("phoneNumber").String()))
	}
	if dataCon.Exists("accountStatus") {
		d.Set("account_status", models.StripQuotes(dataCon.S("accountStatus").String()))
	}
	if dataCon.Exists("domain") {
		d.Set("domain", models.StripQuotes(dataCon.S("domain").String()))
	}

	var roles []interface{}
	var userRbac []interface{}
	if platform == "nd" {
		roles = make([]interface{}, 0)
		userRbac = make([]interface{}, 0)
		if dataCon.Exists("userRbac") {
			for name, _ := range dataCon.S("userRbac").Data().(map[string]interface{}) {
				map1 := make(map[string]interface{})

				map1["roleid"] = models.StripQuotes(name)
				map1["access_type"] = models.StripQuotes(dataCon.S("userRbac").S(name).S("userPriv").String())
				roles = append(roles, map1)

				map2 := make(map[string]interface{})

				map2["name"] = models.StripQuotes(name)
				map2["user_priv"] = models.StripQuotes(dataCon.S("userRbac").S(name).S("userPriv").String())
				userRbac = append(userRbac, map2)

			}
		}

		d.Set("roles", roles)
		d.Set("user_rbac", userRbac)

	} else {
		count, err := dataCon.ArrayCount("roles")

		if err != nil {
			return fmt.Errorf("No Roles found")
		}

		roles = make([]interface{}, 0)
		for i := 0; i < count; i++ {
			rolesCont, err := dataCon.ArrayElement(i, "roles")

			if err != nil {
				return fmt.Errorf("Unable to parse the roles list")
			}

			map1 := make(map[string]interface{})

			map1["roleid"] = models.StripQuotes(rolesCont.S("roleId").String())
			map1["access_type"] = models.StripQuotes(rolesCont.S("accessType").String())
			roles = append(roles, map1)
		}

		d.Set("roles", roles)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
