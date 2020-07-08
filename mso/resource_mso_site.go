package mso

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOSite() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSiteCreate,
		Update: resourceMSOSiteUpdate,
		Read:   resourceMSOSiteRead,
		Delete: resourceMSOSiteDelete,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"apic_site_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"platform": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"login_domain": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"maintenance_mode": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"cloud_providers": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		}),
	}
}

func resourceMSOSiteCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site: Beginning Creation")
	msoClient := m.(*client.Client)
	siteAttr := models.SiteAttributes{}

	if name, ok := d.GetOk("name"); ok {
		siteAttr.Name = name.(string)
	}

	if username, ok := d.GetOk("username"); ok {
		siteAttr.ApicUsername = username.(string)
	}

	if password, ok := d.GetOk("password"); ok {
		siteAttr.ApicPassword = password.(string)
	}

	if apic_site_id, ok := d.GetOk("apic_site_id"); ok {
		siteAttr.ApicSiteId = apic_site_id.(string)
	}

	if labels, ok := d.GetOk("labels"); ok {
		siteAttr.Labels = labels.([]interface{})
	}

	if maintMode, ok := d.GetOk("maintenance_mode"); ok {
		siteAttr.MaintenanceMode = maintMode.(bool)
	}

	if domain, ok := d.GetOk("login_domain"); ok {
		domainStr := domain.(string)
		usrName := d.Get("username").(string)
		siteAttr.ApicUsername = fmt.Sprintf("apic#%s\\\\%s", domainStr, usrName)
		siteAttr.Domain = domainStr
		siteAttr.HasDomain = true
	}

	var loc *models.Location
	if location, ok := d.GetOk("location"); ok {
		loc = &models.Location{}
		loc_map := location.(map[string]interface{})
		loc.Lat, _ = strconv.ParseFloat(fmt.Sprintf("%v", loc_map["lat"]), 64)
		loc.Long, _ = strconv.ParseFloat(fmt.Sprintf("%v", loc_map["long"]), 64)
		siteAttr.Location = loc
	}

	if urls, ok := d.GetOk("urls"); ok {
		siteAttr.Url = urls.([]interface{})
	}

	siteApp := models.NewSite(siteAttr)

	cont, err := msoClient.Save("api/v1/sites", siteApp)
	if err != nil {
		log.Println(err)
		return err
	}

	id := models.StripQuotes(cont.S("id").String())
	d.SetId(fmt.Sprintf("%v", id))
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceMSOSiteRead(d, m)
}

func resourceMSOSiteUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site: Beginning Update")

	msoClient := m.(*client.Client)

	siteAttr := models.SiteAttributes{}

	if name, ok := d.GetOk("name"); ok {
		siteAttr.Name = name.(string)
	}

	if username, ok := d.GetOk("username"); ok {
		siteAttr.ApicUsername = username.(string)
	}

	if password, ok := d.GetOk("password"); ok {
		siteAttr.ApicPassword = password.(string)
	}

	if apic_site_id, ok := d.GetOk("apic_site_id"); ok {
		siteAttr.ApicSiteId = apic_site_id.(string)
	}

	if labels, ok := d.GetOk("labels"); ok {
		siteAttr.Labels = labels.([]interface{})
	}

	if maintMode, ok := d.GetOk("maintenance_mode"); ok {
		siteAttr.MaintenanceMode = maintMode.(bool)
	}

	if domain, ok := d.GetOk("login_domain"); ok {
		domainStr := domain.(string)
		usrName := d.Get("username").(string)
		siteAttr.ApicUsername = fmt.Sprintf("apic#%s\\\\%s", domainStr, usrName)
		siteAttr.Domain = domainStr
		siteAttr.HasDomain = true
	}

	var loc *models.Location
	if location, ok := d.GetOk("location"); ok {
		loc = &models.Location{}
		loc_map := location.(map[string]interface{})
		loc.Lat, _ = strconv.ParseFloat(fmt.Sprintf("%v", loc_map["lat"]), 64)
		loc.Long, _ = strconv.ParseFloat(fmt.Sprintf("%v", loc_map["long"]), 64)
		siteAttr.Location = loc
	}

	if urls, ok := d.GetOk("urls"); ok {
		siteAttr.Url = urls.([]interface{})
	}

	if cloudProviders, ok := d.GetOk("cloud_providers"); ok {
		siteAttr.CloudProviders = cloudProviders.([]interface{})
	}

	siteAttr.Platform = d.Get("platform").(string)
	siteApp := models.NewSite(siteAttr)

	cont, err := msoClient.Put(fmt.Sprintf("api/v1/sites/%s", d.Id()), siteApp)
	if err != nil {
		return err
	}

	id := models.StripQuotes(cont.S("id").String())
	d.SetId(fmt.Sprintf("%v", id))
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceMSOSiteRead(d, m)
}

func resourceMSOSiteRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	msoClient := m.(*client.Client)

	dn := d.Id()

	con, err := msoClient.GetViaURL("api/v1/sites/" + dn)
	if err != nil {
		return err
	}
	d.SetId(models.StripQuotes(con.S("id").String()))
	d.Set("name", models.StripQuotes(con.S("name").String()))
	d.Set("username", models.StripQuotes(con.S("username").String()))
	d.Set("apic_site_id", models.StripQuotes(con.S("apicSiteId").String()))
	d.Set("labels", con.S("labels").Data().([]interface{}))
	d.Set("urls", con.S("urls").Data().([]interface{}))
	d.Set("platform", models.StripQuotes(con.S("platform").String()))
	if con.Exists("cloudProviders") {
		d.Set("cloud_providers", con.S("cloudProviders").Data().([]interface{}))
	} else {
		d.Set("cloud_providers", make([]interface{}, 0, 1))
	}
	if con.Exists("maintenanceMode") {
		d.Set("maintenance_mode", con.S("maintenanceMode").Data().(bool))
	}
	if _, ok := d.GetOk("login_domain"); ok {
		regex := regexp.MustCompile(`apic#(.*)\\{4}(.*)`)
		unameStr := models.StripQuotes(con.S("username").String())
		matches := regex.FindStringSubmatch(unameStr)
		if len(matches) == 3 {
			d.Set("username", matches[2])
			d.Set("login_domain", matches[1])
		}

	}
	loc1 := con.S("location").Data()
	locset := make(map[string]interface{})
	if loc1 != nil {
		loc := loc1.(map[string]interface{})
		locset["lat"] = fmt.Sprintf("%v", loc["lat"])
		locset["long"] = fmt.Sprintf("%v", loc["long"])
	} else {
		locset = nil
	}
	d.Set("location", locset)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSOSiteDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	msoClient := m.(*client.Client)
	dn := d.Id()
	err := msoClient.DeletebyId(fmt.Sprintf("api/v1/sites/%v%s", dn, "?force=true"))
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
