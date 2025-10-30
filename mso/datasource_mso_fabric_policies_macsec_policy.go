package mso

import (
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMacsecPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMacsecPolicyRead,

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
			"admin_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"interface_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cipher_suite": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"window_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"security_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sak_expire_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"confidentiality_offset": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_server_priority": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"macsec_key": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"psk": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceMacsecPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO MACsec Policy Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	setMacsecPolicyData(d, msoClient, templateId, policyName)
	log.Printf("[DEBUG] MSO MACsec Policy Data Source - Read Complete : %v", d.Id())
	return nil
}
