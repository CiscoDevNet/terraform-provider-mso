package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSOPhysicalInterface() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOPhysicalInterfaceRead,

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
			"nodes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"interfaces": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"interface_policy_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"breakout_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"interface_descriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interface": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_group_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMSOPhysicalInterfaceRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Physical Interface Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "fabricResourceTemplate", "template", "interfaceProfiles")
	if err != nil {
		return err
	}

	setPhysicalInterfaceData(d, policy, templateId)
	log.Printf("[DEBUG] MSO Physical Interface Data Source - Read Complete: %v", d.Id())
	return nil
}
