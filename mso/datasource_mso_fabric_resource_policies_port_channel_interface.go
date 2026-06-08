package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceMSOPortChannelInterface() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOPortChannelInterfaceRead,

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
			"node": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"interfaces": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"interface_policy_group_uuid": {
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
		},
	}
}

func dataSourceMSOPortChannelInterfaceRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Port Channel Interface Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "fabricResourceTemplate", "template", "portChannels")
	if err != nil {
		return err
	}

	err = setPortChannelInterfaceData(d, policy, templateId)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] MSO Port Channel Interface Data Source - Read Complete: %v", d.Id())
	return nil
}
