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

func dataSourceMSOSchemaSiteBdSubnet() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceMSOSchemaSiteBdSubnetRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
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

			"ip": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scope": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"shared": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"no_default_gateway": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"querier": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceMSOSchemaSiteBdSubnetRead(d *schema.ResourceData, m interface{}) error {
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
	found := false
	stateBd := d.Get("bd_name").(string)
	stateIp := d.Get("ip").(string)
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())

		if apiSite == stateSite {
			d.Set("site_id", apiSite)
			d.Set("template_name", models.StripQuotes(tempCont.S("templateName").String()))
			bdCount, err := tempCont.ArrayCount("bds")
			if err != nil {
				return fmt.Errorf("Unable to get Bd list")
			}
			for j := 0; j < bdCount; j++ {
				bdCont, err := tempCont.ArrayElement(j, "bds")
				if err != nil {
					return err
				}
				bdRef := models.StripQuotes(bdCont.S("bdRef").String())
				re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/bds/(.*)")
				match := re.FindStringSubmatch(bdRef)
				if match[3] == stateBd {
					d.Set("bd_name", match[3])
					subnetCount, err := bdCont.ArrayCount("subnets")
					if err != nil {
						return fmt.Errorf("Unable to get Subnet list")
					}
					for l := 0; l < subnetCount; l++ {
						subnetCont, err := bdCont.ArrayElement(l, "subnets")
						if err != nil {
							return err
						}
						apiIP := models.StripQuotes(subnetCont.S("ip").String())
						if stateIp == apiIP {
							d.SetId(apiIP)
							if subnetCont.Exists("ip") {
								d.Set("ip", models.StripQuotes(subnetCont.S("ip").String()))
							}
							if subnetCont.Exists("description") {
								d.Set("description", models.StripQuotes(subnetCont.S("description").String()))
							}
							if subnetCont.Exists("scope") {
								d.Set("scope", models.StripQuotes(subnetCont.S("scope").String()))
							}
							if subnetCont.Exists("shared") {
								d.Set("shared", subnetCont.S("shared").Data().(bool))
							}
							if subnetCont.Exists("noDefaultGateway") {
								d.Set("no_default_gateway", subnetCont.S("noDefaultGateway").Data().(bool))
							}
							if subnetCont.Exists("querier") {
								d.Set("querier", subnetCont.S("querier").Data().(bool))
							}
							found = true
							break
						}
					}
				}

			}
		}
	}
	if !found {
		d.SetId("")
		return fmt.Errorf("Unable to find the Site BD Subnet")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
