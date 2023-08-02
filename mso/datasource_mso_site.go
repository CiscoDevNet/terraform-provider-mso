package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceMSOSite() *schema.Resource {
	return &schema.Resource{

		Read: datasourceMSOSiteRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"reprovision": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"proxy": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"sr_l3out": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"template_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"apic_site_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"labels": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"location": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lat": &schema.Schema{
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"long": &schema.Schema{
							Type:     schema.TypeFloat,
							Computed: true,
						},
					},
				},
			},
			"urls": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"cloud_providers": &schema.Schema{
				Type:     schema.TypeList,
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
	var path string
	platform := msoClient.GetPlatform()
	if platform == "nd" {
		path = "api/v2/sites"
	} else {
		path = "api/v1/sites"
	}
	con, err := msoClient.GetViaURL(path)
	if err != nil {
		return err
	}

	data := con.S("sites").Data().([]interface{})

	var flag bool
	var count int

	for _, info := range data {
		val := info.(map[string]interface{})
		if platform == "nd" {
			val = val["common"].(map[string]interface{})
		}
		if val["name"].(string) == name {
			flag = true
			break
		}

		count = count + 1
	}

	if flag != true {
		return fmt.Errorf("Site of specified name not found")
	}

	dataCon := con.S("sites").Index(count)
	d.SetId(models.StripQuotes(dataCon.S("id").String()))

	if platform == "nd" {

		if dataCon.Exists("isSrL3OutEnabled") {
			d.Set("sr_l3out", dataCon.S("isSrL3OutEnabled").Data().(bool))
		}

		if dataCon.Exists("templateCount") {
			d.Set("template_count", dataCon.S("templateCount").Data().(float64))
		}

		dataConAttr := dataCon.S("common")

		d.Set("name", models.StripQuotes(dataConAttr.S("name").String()))

		if dataConAttr.Exists("siteId") {
			d.Set("apic_site_id", models.StripQuotes(dataConAttr.S("siteId").String()))
		}

		if dataConAttr.Exists("urls") {
			d.Set("urls", dataConAttr.S("urls").Data().([]interface{}))
		}

		if dataConAttr.Exists("username") {
			d.Set("username", models.StripQuotes(dataConAttr.S("username").String()))
		}

		if dataConAttr.Exists("platformType") {
			d.Set("type", models.StripQuotes(dataConAttr.S("platformType").String()))
		}

		if dataConAttr.Exists("siteGroup") {
			d.Set("group_id", models.StripQuotes(dataConAttr.S("siteGroup").String()))
		}

		if dataConAttr.Exists("siteVersion") {
			d.Set("version", models.StripQuotes(dataConAttr.S("siteVersion").String()))
		}

		if dataConAttr.Exists("siteConnectivityStatus") {
			d.Set("status", models.StripQuotes(dataConAttr.S("siteConnectivityStatus").String()))
		}

		if dataConAttr.Exists("useProxy") {
			d.Set("proxy", dataConAttr.S("useProxy").Data().(bool))
		}

		if dataConAttr.Exists("needsReprovision") {
			d.Set("reprovision", dataConAttr.S("needsReprovision").Data().(bool))
		}

		if dataConAttr.Exists("latitude") || dataConAttr.Exists("longitude") {
			locset := make(map[string]interface{})
			locset["lat"] = models.StripQuotes(dataConAttr.S("latitude").String())
			locset["long"] = models.StripQuotes(dataConAttr.S("longitude").String())
			d.Set("location", locset)
		}

		dataCloud := dataCon.S("apic")
		if dataCloud.Exists("cApicType") {
			provider := [1]string{models.StripQuotes(dataCloud.S("cApicType").String())}
			d.Set("cloud_providers", provider)
		}

	} else {

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

		if dataCon.Exists("cloudProviders") {
			d.Set("cloud_providers", dataCon.S("cloudProviders").Data().([]interface{}))
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
