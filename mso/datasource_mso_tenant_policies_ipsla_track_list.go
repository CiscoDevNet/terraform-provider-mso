package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSOIPSLATrackList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOIPSLATrackListRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"threshold_up": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"threshold_down": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"members": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipsla_monitoring_policy_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scope_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scope_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"weight": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceMSOIPSLATrackListRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO IPSLA Track List Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "ipslaTrackLists")
	if err != nil {
		return err
	}

	setIPSLATrackListData(d, policy, templateId)
	log.Printf("[DEBUG] MSO IPSLA Track List Data Source - Read Complete : %v", d.Id())
	return nil
}
