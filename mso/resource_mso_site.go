package mso

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
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

		Importer: &schema.ResourceImporter{
			State: resourceMSOSiteImport,
		},

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"username": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"password": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
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
				// Computed: true -> Removed, as it creates discrepancy for idempotency.
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
				Computed: true,
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

func resourceMSOSiteImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] Site: Beginning Import")
	msoClient := m.(*client.Client)

	dn := d.Id()

	var apiVersion string
	platform := msoClient.GetPlatform()
	if platform == "nd" {
		apiVersion = "v2"
	} else {
		apiVersion = "v1"
	}
	path := fmt.Sprintf("api/%v/sites/%v", apiVersion, dn)

	con, err := msoClient.GetViaURL(path)
	if err != nil {
		return nil, err
	}

	d.SetId(models.StripQuotes(con.S("id").String()))

	if platform == "nd" {
		dataConAttr := con.S("common")
		d.Set("name", models.StripQuotes(dataConAttr.S("name").String()))

		d.Set("username", models.StripQuotes(dataConAttr.S("username").String()))

		platform := models.StripQuotes(dataConAttr.S("platformType").String())
		if platform == "CloudApic" {
			d.Set("platform", "cloud")
		} else {
			d.Set("platform", "on-premise")
		}

		if dataConAttr.Exists("siteId") {
			d.Set("apic_site_id", models.StripQuotes(dataConAttr.S("siteId").String()))
		}

		if dataConAttr.Exists("labels") {
			d.Set("labels", dataConAttr.S("labels").Data().([]interface{}))
		} else {
			d.Set("labels", nil)
		}

		dataCloud := con.S("apic")
		if dataCloud.Exists("cApicType") {
			provider := [1]string{models.StripQuotes(dataCloud.S("cApicType").String())}
			d.Set("cloud_providers", provider)
		} else {
			d.Set("cloud_providers", nil)

		}
		d.Set("maintenance_mode", false)
		if dataConAttr.Exists("latitude") && dataConAttr.Exists("longitude") {
			locset := make(map[string]interface{})
			locset["lat"] = models.StripQuotes(dataConAttr.S("latitude").String())
			locset["long"] = models.StripQuotes(dataConAttr.S("longitude").String())
			d.Set("location", locset)
		}
		if dataConAttr.Exists("urls") {
			urls := dataConAttr.S("urls").Data().([]interface{})
			var urlsList []interface{}
			for _, url := range urls {
				split := strings.Split(url.(string), ":")
				urlsList = append(urlsList, fmt.Sprintf("%s:%s", split[0], split[1]))
			}
			d.Set("urls", urlsList)
		}
	} else {
		d.Set("name", models.StripQuotes(con.S("name").String()))
		d.Set("apic_site_id", models.StripQuotes(con.S("apicSiteId").String()))
		username := models.StripQuotes(con.S("username").String())
		d.Set("username", username)
		if username != "" {
			regex := regexp.MustCompile(`apic#(.*)\\{2}(.*)`)
			matches := regex.FindStringSubmatch(username)
			if len(matches) == 3 {
				d.Set("username", matches[2])
				d.Set("login_domain", matches[1])
			}

		}
		if con.Exists("labels") {
			d.Set("labels", con.S("labels").Data().([]interface{}))
		}
		if con.Exists("urls") {
			d.Set("urls", con.S("urls").Data().([]interface{}))
		}
		d.Set("platform", models.StripQuotes(con.S("platform").String()))
		if con.Exists("cloudProviders") {
			d.Set("cloud_providers", con.S("cloudProviders").Data().([]interface{}))
		} else {
			d.Set("cloud_providers", make([]interface{}, 0, 1))
		}
		if con.Exists("maintenanceMode") {
			d.Set("maintenance_mode", con.S("maintenanceMode").Data().(bool))
		}

		if con.Exists("location") {
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
		}
	}

	log.Printf("[DEBUG] %s: Site Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSiteCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site: Beginning Creation")
	msoClient := m.(*client.Client)
	siteAttr := models.SiteAttributes{}

	if name, ok := d.GetOk("name"); ok {
		siteAttr.Name = name.(string)
	}

	if apic_site_id, ok := d.GetOk("apic_site_id"); ok {
		siteAttr.ApicSiteId = apic_site_id.(string)
	}

	var apiVersion string
	var path string
	var id string
	apic_site_id := d.Get("apic_site_id").(string)
	platform := msoClient.GetPlatform()
	if platform == "nd" {
		apiVersion = "v2"
		path = fmt.Sprintf("api/%v/sites/manage", apiVersion)
		siteCont, err := GetSiteViaName(msoClient, siteAttr.Name)
		if err != nil {
			return err
		}
		siteData := siteCont.Data().(map[string]interface{})

		if siteData["id"] == "" {
			common := siteData["common"].(map[string]interface{})
			common["siteId"] = apic_site_id
			siteData["common"] = common
			payload, err := container.Consume(siteData)
			req, err := msoClient.MakeRestRequest("POST", path, payload, true)
			if err != nil {
				return err
			}
			res, _, err := msoClient.Do(req)
			if err != nil {
				return err
			}
			err = client.CheckForErrors(res, "POST")
			if err != nil {
				return err
			}
			id = models.StripQuotes(res.S("id").String())
		} else {
			id = siteData["id"].(string)
		}
	} else {
		apiVersion = "v1"
		path = fmt.Sprintf("api/v1/sites")

		if username, ok := d.GetOk("username"); ok {
			siteAttr.ApicUsername = username.(string)
		}

		if password, ok := d.GetOk("password"); ok {
			siteAttr.ApicPassword = password.(string)
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
		cont, err := msoClient.Save(path, siteApp)
		if err != nil {
			log.Println(err)
			return err
		}
		id = models.StripQuotes(cont.S("id").String())
	}

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

	if apic_site_id, ok := d.GetOk("apic_site_id"); ok {
		siteAttr.ApicSiteId = apic_site_id.(string)
	}

	var apiVersion string
	var path string
	platform := msoClient.GetPlatform()
	if platform == "nd" {
		apiVersion = "v2"
		path = fmt.Sprintf("api/%v/sites/manage", apiVersion)
	} else {
		apiVersion = "v1"
		path = "api/v1/sites"
		if username, ok := d.GetOk("username"); ok {
			siteAttr.ApicUsername = username.(string)
		}

		if password, ok := d.GetOk("password"); ok {
			siteAttr.ApicPassword = password.(string)
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
	}

	siteAttr.Platform = d.Get("platform").(string)
	siteApp := models.NewSite(siteAttr)

	cont, err := msoClient.Put(fmt.Sprintf("%v/%s", path, d.Id()), siteApp)
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

	var apiVersion string
	platform := msoClient.GetPlatform()
	if platform == "nd" {
		apiVersion = "v2"
	} else {
		apiVersion = "v1"
	}
	path := fmt.Sprintf("api/%v/sites/%v", apiVersion, dn)

	con, err := msoClient.GetViaURL(path)
	if err != nil {
		return errorForObjectNotFound(err, dn, con, d)
	}

	d.SetId(models.StripQuotes(con.S("id").String()))

	if platform == "nd" {
		dataConAttr := con.S("common")
		d.Set("name", models.StripQuotes(dataConAttr.S("name").String()))

		d.Set("username", models.StripQuotes(dataConAttr.S("username").String()))

		if dataConAttr.Exists("siteId") {
			d.Set("apic_site_id", models.StripQuotes(dataConAttr.S("siteId").String()))
		}

		if dataConAttr.Exists("labels") {
			d.Set("labels", dataConAttr.S("labels").Data().([]interface{}))
		} else {
			d.Set("labels", nil)
		}
		if dataConAttr.Exists("urls") {
			urls := dataConAttr.S("urls").Data().([]interface{})
			var urlsList []interface{}
			for _, url := range urls {
				split := strings.Split(url.(string), ":")
				urlsList = append(urlsList, fmt.Sprintf("%s:%s", split[0], split[1]))
			}
			d.Set("urls", urlsList)
		}

		platform := models.StripQuotes(dataConAttr.S("platformType").String())
		if platform == "CloudApic" {
			d.Set("platform", "cloud")
		} else {
			d.Set("platform", "on-premise")
		}

		dataCloud := con.S("apic")
		if dataCloud.Exists("cApicType") {
			provider := [1]string{models.StripQuotes(dataCloud.S("cApicType").String())}
			d.Set("cloud_providers", provider)
		} else {
			d.Set("cloud_providers", nil)

		}
		d.Set("maintenance_mode", false)
		if dataConAttr.Exists("latitude") && dataConAttr.Exists("longitude") {
			locset := make(map[string]interface{})
			locset["lat"] = models.StripQuotes(dataConAttr.S("latitude").String())
			locset["long"] = models.StripQuotes(dataConAttr.S("longitude").String())
			d.Set("location", locset)
		}
	} else {
		d.Set("name", models.StripQuotes(con.S("name").String()))
		d.Set("apic_site_id", models.StripQuotes(con.S("apicSiteId").String()))
		username := models.StripQuotes(con.S("username").String())
		d.Set("username", username)
		if username != "" {
			regex := regexp.MustCompile(`apic#(.*)\\{2}(.*)`)
			matches := regex.FindStringSubmatch(username)
			if len(matches) == 3 {
				d.Set("username", matches[2])
				d.Set("login_domain", matches[1])
			}

		}
		if con.Exists("labels") {
			d.Set("labels", con.S("labels").Data().([]interface{}))
		}
		if con.Exists("urls") {
			d.Set("urls", con.S("urls").Data().([]interface{}))
		}
		d.Set("platform", models.StripQuotes(con.S("platform").String()))
		if con.Exists("cloudProviders") {
			d.Set("cloud_providers", con.S("cloudProviders").Data().([]interface{}))
		} else {
			d.Set("cloud_providers", make([]interface{}, 0, 1))
		}
		if con.Exists("maintenanceMode") {
			d.Set("maintenance_mode", con.S("maintenanceMode").Data().(bool))
		}

		if con.Exists("location") {
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
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSOSiteDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	msoClient := m.(*client.Client)
	dn := d.Id()
	var apiVersion string
	var path string
	platform := msoClient.GetPlatform()
	if platform == "nd" {
		apiVersion = "v2"
		path = fmt.Sprintf("api/%v/sites/manage/%v", apiVersion, dn)
	} else {
		apiVersion = "v1"
		path = fmt.Sprintf("api/%v/sites/%v%s", apiVersion, dn, "?force=true")
	}
	err := msoClient.DeletebyId(path)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

func GetSiteViaName(msoClient *client.Client, name string) (*container.Container, error) {
	cont, err := msoClient.GetViaURL("api/v2/sites")
	if err != nil {
		return nil, err
	}
	sitesCount, err := cont.ArrayCount("sites")
	if err != nil {
		return nil, err
	}

	for i := 0; i < sitesCount; i++ {
		siteCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return nil, err
		}
		apiName := models.StripQuotes(siteCont.S("common").S("name").String())
		if apiName == name {
			return siteCont, nil
		}
	}
	return nil, fmt.Errorf(fmt.Sprintf("Site %v is not a valid Site configured at ND-level. Add Site to ND first.", name))
}
