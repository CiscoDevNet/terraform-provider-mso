package mso

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceMSOSite() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSiteCreate,
		Update: resourceMSOSiteUpdate,
		Read:   resourceMSOSiteRead,
		Delete: resourceMSOSiteDelete,

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

			"site": &schema.Schema{
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
				Type:     schema.TypeSet,
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

	if site, ok := d.GetOk("site"); ok {
		siteAttr.Site = site.(string)
	}

	if labels, ok := d.GetOk("labels"); ok {
		siteAttr.Labels = labels.([]interface{})
	}

	var loc *models.Location
	if location, ok := d.GetOk("location"); ok {
		loc = &models.Location{}
		tp := location.(map[string]interface{})
		loc.Latitude, _ = strconv.ParseFloat(fmt.Sprintf("%v", tp["lat"]), 64)
		loc.Longitude, _ = strconv.ParseFloat(fmt.Sprintf("%v", tp["long"]), 64)
		if loc != nil {
			siteAttr.Location = loc
		} else {
			siteAttr.Location = nil
		}
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

	if d.HasChange("name") {
		siteAttr.Name = d.Get("name").(string)
	}

	if d.HasChange("username") {
		siteAttr.ApicUsername = d.Get("username").(string)
	}

	if d.HasChange("password") {
		siteAttr.ApicPassword = d.Get("password").(string)
	}

	if d.HasChange("apic_site_id") {
		siteAttr.ApicSiteId = d.Get("apic_site_id").(string)
	}

	if d.HasChange("site") {
		siteAttr.Site = d.Get("site").(string)
	}

	if d.HasChange("labels") {
		siteAttr.Labels = d.Get("labels").([]interface{})
	}

	if d.HasChange("urls") {
		siteAttr.Url = d.Get("urls").([]interface{})
	}

	var loc *models.Location
	if location, ok := d.GetOk("location"); ok {
		loc = &models.Location{}
		tp := location.(map[string]interface{})
		loc.Latitude, _ = strconv.ParseFloat(fmt.Sprintf("%v", tp["lat"]), 64)
		loc.Longitude, _ = strconv.ParseFloat(fmt.Sprintf("%v", tp["long"]), 64)
		if loc != nil {
			siteAttr.Location = loc
		} else {
			siteAttr.Location = nil
		}
	}

	siteApp := models.NewSite(siteAttr)
	cont, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/sites/%s", d.Id()), siteApp)
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
	d.Set("name", con.S("name").String())
	d.Set("username", con.S("username").String())
	d.Set("password", con.S("password").String())
	d.Set("apic_site_id", con.S("apic_site_id").String())
	d.Set("site", con.S("site").String())

	count1, err := con.ArrayCount("labels")
	if err != nil {
		return fmt.Errorf("No Label found")
	}
	found1 := false
	for i := 0; i < count1; i++ {

		temp, err := con.ArrayElement(i, "labels")
		log.Println(temp)
		d.Set("labels", temp)
		found1 = true
		if err != nil {
			return fmt.Errorf("Unable to parse the label list")
		}
	}
	if !found1 {
		d.Set("labels", "")
	}

	count2, err := con.ArrayCount("urls")
	if err != nil {
		return fmt.Errorf("No URL found")
	}
	found2 := false
	for i := 0; i < count2; i++ {

		temp, err := con.ArrayElement(i, "urls")
		log.Println(temp)
		d.Set("urls", temp)
		found2 = true
		if err != nil {
			return fmt.Errorf("Unable to parse the url list")
		}
	}
	if !found2 {
		d.Set("urls", "")
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
	err := msoClient.DeletebyId("api/v1/sites" + dn)
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
