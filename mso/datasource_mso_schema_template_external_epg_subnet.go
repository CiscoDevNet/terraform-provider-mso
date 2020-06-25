package mso

import (
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceMSOTemplateExternalEpgSubnet() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceMSOTemplateExternalEpgSubnetRead,

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
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"external_epg_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"ip": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"scope": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"aggregate": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		}),
	}
}

func dataSourceMSOTemplateExternalEpgSubnetRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No Template found")
	}
	stateTemplate := d.Get("template_name").(string)
	found := false
	stateExternalepg := d.Get("external_epg_name")
	stateIP := d.Get("ip")

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {
			externalepgCount, err := tempCont.ArrayCount("externalEpgs")
			if err != nil {
				return fmt.Errorf("Unable to get Externalepg list")
			}
			for j := 0; j < externalepgCount; j++ {
				externalepgCont, err := tempCont.ArrayElement(j, "externalEpgs")
				if err != nil {
					return err
				}
				apiExternalepg := models.StripQuotes(externalepgCont.S("name").String())
				if apiExternalepg == stateExternalepg {
					subnetCount, err := externalepgCont.ArrayCount("subnets")
					if err != nil {
						return fmt.Errorf("Unable to get subnets list")
					}
					for k := 0; k < subnetCount; k++ {
						subnetsCont, err := externalepgCont.ArrayElement(k, "subnets")
						if err != nil {
							return err
						}
						apiIP := models.StripQuotes(subnetsCont.S("ip").String())
						if apiIP == stateIP {
							d.Set("schema_id", schemaId)
							d.Set("template_name", apiTemplate)
							d.Set("external_epg_name", apiExternalepg)
							ip := models.StripQuotes(subnetsCont.S("ip").String())
							idSubnet := strings.Split(ip, "/")
							d.SetId(idSubnet[0])
							d.Set("ip", models.StripQuotes(subnetsCont.S("ip").String()))
							d.Set("name", models.StripQuotes(subnetsCont.S("name").String()))
							d.Set("scope", subnetsCont.S("scope").Data().([]interface{}))
							d.Set("aggregate", subnetsCont.S("aggregate").Data().([]interface{}))

							found = true
							break
						}
					}
				}
				if found {
					break
				}
			}
		}
		if found {
			break
		}
	}

	if !found {
		d.SetId("")
		return fmt.Errorf("External Epg Subnet Not Found")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
