package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceMSOSchemaTemplateAnpEpg() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOTemplateAnpEpgRead,

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
			"anp_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"bd_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"bd_schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"bd_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vrf_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vrf_schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vrf_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"useg_epg": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"intra_epg": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"intersite_multicast_source": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"proxy_arp": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"preferred_group": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"epg_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"deployment_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_service_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}
}

func dataSourceMSOTemplateAnpEpgRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}

	stateTemplate := d.Get("template_name").(string)
	stateANP := d.Get("anp_name").(string)
	stateEPG := d.Get("name").(string)

	err = resourceMSOSchemaTemplateAnpEpgSetAttr(stateTemplate, stateANP, stateEPG, cont, d)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
