package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func datasourceMSOTenant() *schema.Resource {
	return &schema.Resource{

		Read: datasourceMSOTenantRead,

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

func datasourceMSOTenantRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	name := d.Get("name").(string)
	con, err := msoClient.GetViaURL("api/v1/tenants")
	if err != nil {
		return err
	}

	data := con.S("tenants").Data().([]interface{})
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
		return fmt.Errorf("Tenant of specified name not found")
	}

	dataCon := con.S("tenants").Index(count)

	d.SetId(models.StripQuotes(dataCon.S("id").String()))

	d.Set("name", models.StripQuotes(dataCon.S("name").String()))

	d.Set("display_name", models.StripQuotes(dataCon.S("displayName").String()))

	if dataCon.Exists("description") {
		d.Set("description", models.StripQuotes(dataCon.S("description").String()))
	}

	count1, _ := dataCon.ArrayCount("siteAssociations")
	site_associations := make([]interface{}, 0)
	for i := 0; i < count1; i++ {
		sitesCont, err := dataCon.ArrayElement(i, "siteAssociations")
		if err != nil {
			return fmt.Errorf("Unable to parse the site associations list")
		}

		map1 := make(map[string]interface{})
		map1["site_id"] = models.StripQuotes(sitesCont.S("siteId").String())
		site_associations = append(site_associations, map1)
	}

	d.Set("site_associations", site_associations)

	count2, err := con.ArrayCount("userAssociations")
	if err != nil {
		d.Set("user_assocoations", make([]interface{}, 0))
	}
	user_associations := make([]interface{}, 0)
	for i := 0; i < count2; i++ {
		usersCont, err := dataCon.ArrayElement(i, "userAssociations")
		if err != nil {
			return fmt.Errorf("Unable to parse the user associations list")
		}

		map1 := make(map[string]interface{})
		map1["user_id"] = models.StripQuotes(usersCont.S("userId").String())
		user_associations = append(user_associations, map1)
	}

	d.Set("user_associations", user_associations)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
