package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSOServiceDeviceCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOServiceDeviceClusterRead,

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
			"device_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"device_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"interface_properties": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bd_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"external_epg_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipsla_monitoring_policy_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"qos_policy_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"preferred_group": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"rewrite_source_mac": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"anycast": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"config_static_mac": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_backup_redirect_ip": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"load_balance_hashing": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pod_aware_redirection": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"resilient_hashing": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"tag_based_sorting": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"min_threshold": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_threshold": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"threshold_down_action": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceMSOServiceDeviceClusterRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Service Device Cluster Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "deviceTemplate", "template", "devices")
	if err != nil {
		return err
	}

	setServiceDeviceClusterData(d, policy, templateId)
	log.Printf("[DEBUG] MSO Service Device Cluster Data Source - Read Complete: %v", d.Id())
	return nil
}
