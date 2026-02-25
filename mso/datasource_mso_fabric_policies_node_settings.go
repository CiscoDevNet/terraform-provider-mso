package mso

import (
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSONodeSettings() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNodeSettingsRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"synce": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"admin_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quality_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ptp": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_domain": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"priority_2": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNodeSettingsRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Node Settings Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	err := setNodeSettingsData(d, msoClient, templateId, policyName)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] MSO Node Settings Data Source - Read Complete : %v", d.Id())
	return nil
}
