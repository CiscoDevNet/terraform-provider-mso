package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceMSOTemplateSubnetBD() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOTemplateSubnetBDRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"bd_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"ip": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"scope": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"shared": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"primary": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"no_default_gateway": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"querier": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		}),
	}
}

func dataSourceMSOTemplateSubnetBDRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No TemplateSubnet found")
	}
	stateTemplateSubnet := d.Get("template_name").(string)
	stateBD := d.Get("bd_name")
	stateIP := d.Get("ip")

	found := false
	for i := 0; i < count && !found; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		apiTemplateSubnet := models.StripQuotes(tempCont.S("name").String())

		if apiTemplateSubnet == stateTemplateSubnet {
			bdCount, err := tempCont.ArrayCount("bds")
			if err != nil {
				return fmt.Errorf("Unable to get BD list")
			}
			for j := 0; j < bdCount && !found; j++ {
				bdCont, err := tempCont.ArrayElement(j, "bds")
				if err != nil {
					return err
				}

				apiBD := models.StripQuotes(bdCont.S("name").String())
				if apiBD == stateBD {
					count1, err := bdCont.ArrayCount("subnets")
					if err != nil {
						return fmt.Errorf("Unable to get Subnet List")
					}
					for k := 0; k < count1; k++ {
						dataCon, err := bdCont.ArrayElement(k, "subnets")
						if err != nil {
							return fmt.Errorf("Unable to parse the subnets list")
						}

						apiIP := models.StripQuotes(dataCon.S("ip").String())
						if apiIP == stateIP {
							d.SetId(fmt.Sprintf("%s/templates/%s/bds/%s/subnets/%s", schemaId, stateTemplateSubnet, stateBD, stateIP))
							d.Set("schema_id", schemaId)
							d.Set("template_name", apiTemplateSubnet)
							d.Set("bd_name", apiBD)
							d.Set("ip", models.StripQuotes(dataCon.S("ip").String()))
							d.Set("scope", models.StripQuotes(dataCon.S("scope").String()))
							d.Set("description", models.StripQuotes(dataCon.S("description").String()))
							d.Set("shared", dataCon.S("shared").Data().(bool))
							if dataCon.Exists("noDefaultGateway") {
								d.Set("no_default_gateway", dataCon.S("noDefaultGateway").Data().(bool))
							}
							if dataCon.Exists("querier") {
								d.Set("querier", dataCon.S("querier").Data().(bool))
							}
							if dataCon.Exists("primary") {
								d.Set("primary", dataCon.S("primary").Data().(bool))
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
		return fmt.Errorf("Unable to find the BD Subnet %s in Template %s of Schema Id %s ", stateIP, stateTemplateSubnet, schemaId)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
