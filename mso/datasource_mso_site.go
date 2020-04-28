package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func datasourceMSOSite() *schema.Resource {
	return &schema.Resource{

		Read: datasourceMSOSiteRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"username": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"apic_site_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"labels": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},

			"location": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"lat": &schema.Schema{
							Type:     schema.TypeFloat,
							Optional: true,
							Computed: true,
						},
						"long": &schema.Schema{
							Type:     schema.TypeFloat,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"urls": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		}),
	}
}

func datasourceMSOSiteRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	name := d.Get("name").(string)
	con, err := msoClient.GetViaURL("api/v1/sites")
	if err != nil {
		return err
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
		return fmt.Errorf("Sie of specified name not found")
	}

	dataCon := con.S("sites").Index(count)

	d.SetId(models.StripQuotes(dataCon.S("id").String()))

	d.Set("name", models.StripQuotes(dataCon.S("name").String()))

	if dataCon.Exists("username") {
		d.Set("username", models.StripQuotes(dataCon.S("username").String()))
	}

	if dataCon.Exists("password") {
		d.Set("password", models.StripQuotes(dataCon.S("password").String()))
	}

	if dataCon.Exists("apicSiteId") {
		d.Set("apic_site_id", models.StripQuotes(dataCon.S("apicSiteId").String()))
	}

	loc1 := dataCon.S("location").Data()
	locset := make(map[string]interface{})
	if loc1 != nil {
		loc := loc1.(map[string]interface{})
		locset["lat"] = fmt.Sprintf("%v", loc["lat"])
		locset["long"] = fmt.Sprintf("%v", loc["long"])
	} else {
		locset = nil
	}
	d.Set("location", locset)

	if dataCon.Exists("labels") {
		d.Set("labels", dataCon.S("labels").Data().([]interface{}))
	}

	if dataCon.Exists("urls") {
		d.Set("urls", dataCon.S("urls").Data().([]interface{}))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
