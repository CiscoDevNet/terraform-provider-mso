package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSOVlanPool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOVlanPoolRead,

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
			"allocation_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vlan_range": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"to": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"allocation_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceMSOVlanPoolRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO VLAN Pool Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "fabricPolicyTemplate", "template", "vlanPools")
	if err != nil {
		return err
	}

	setVlanPoolData(d, policy, templateId)
	log.Printf("[DEBUG] MSO VLAN Pool Data Source - Read Complete : %v", d.Id())
	return nil
}
