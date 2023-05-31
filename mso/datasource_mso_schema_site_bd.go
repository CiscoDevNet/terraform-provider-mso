package mso

import (
	"fmt"
	"log"
	"regexp"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceMSOSchemaSiteBd() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceMSOSchemaSiteBdRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"site_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"bd_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"host_route": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"svi_mac": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}
}

func dataSourceMSOSchemaSiteBdRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return fmt.Errorf("No Sites found")
	}
	stateSite := d.Get("site_id").(string)
	stateTemplate := d.Get("template_name").(string)
	found := false
	statebd := d.Get("bd_name").(string)
	for i := 0; i < count && !found; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())
		apiTemplate := models.StripQuotes(tempCont.S("templateName").String())

		if apiSite == stateSite && apiTemplate == stateTemplate {
			bdCount, err := tempCont.ArrayCount("bds")
			if err != nil {
				return fmt.Errorf("Unable to get bd list")
			}
			for j := 0; j < bdCount && !found; j++ {
				bdCont, err := tempCont.ArrayElement(j, "bds")
				if err != nil {
					return err
				}
				bdRef := models.StripQuotes(bdCont.S("bdRef").String())
				re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/bds/(.*)")
				match := re.FindStringSubmatch(bdRef)
				if match[3] == statebd {
					d.SetId(match[3])
					d.Set("bd_name", match[3])
					d.Set("schema_id", match[1])
					d.Set("template_name", match[2])
					d.Set("site_id", apiSite)
					if bdCont.Exists("hostBasedRouting") {
						d.Set("host_route", bdCont.S("hostBasedRouting").Data().(bool))
					}
					if bdCont.Exists("mac") {
						d.Set("svi_mac", models.StripQuotes(bdCont.S("mac").String()))
					}
					found = true
					break
				}
			}
		}
	}

	if !found {
		d.SetId("")
		return fmt.Errorf("Unable to find the given Schema Site Bd")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
