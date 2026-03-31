package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSONetflowRecord() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSONetflowRecordRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the tenant policy template.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the NetFlow Record.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the NetFlow Record.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the NetFlow Record.",
			},
			"match_parameters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The match parameters of the NetFlow Record.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceMSONetflowRecordRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO NetFlow Record Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "netFlowRecords")
	if err != nil {
		return err
	}

	setNetflowRecordData(d, policy, templateId)
	log.Printf("[DEBUG] MSO NetFlow Record Data Source - Read Complete: %v", d.Id())
	return nil
}
