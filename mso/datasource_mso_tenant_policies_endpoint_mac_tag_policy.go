package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSOEndpointMACTagPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOEndpointMACTagPolicyRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mac": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bd_uuid": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"vrf_uuid"},
			},
			"vrf_uuid": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"bd_uuid"},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag_annotations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"policy_tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceMSOEndpointMACTagPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Endpoint MAC Tag Policy Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	mac := d.Get("mac").(string)
	bdUUID := d.Get("bd_uuid").(string)
	vrfUUID := d.Get("vrf_uuid").(string)

	if bdUUID == "" && vrfUUID == "" {
		return fmt.Errorf("Either 'bd_uuid' or 'vrf_uuid' must be specified to use Endpoint MAC Tag Policy Data Source")
	}

	name, dataSourceId, err := setEndpointMACTagPolicyId(m, templateId, "", mac, bdUUID, vrfUUID)
	if err != nil {
		return err
	}
	d.SetId(dataSourceId)

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, name, "tenantPolicyTemplate", "template", "endpointMacTagPolicies")
	if err != nil {
		return err
	}

	err = setEndpointMACTagPolicyData(d, policy, templateId, m)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] MSO Endpoint MAC Tag Policy Data Source - Read Complete : %v", d.Id())
	return nil
}
