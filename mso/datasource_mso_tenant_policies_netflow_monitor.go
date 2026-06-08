package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceMSOTenantPoliciesNetflowMonitor() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOTenantPoliciesNetflowMonitorRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the tenant policy template.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the NetFlow Monitor.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the NetFlow Monitor.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the NetFlow Monitor.",
			},
			"netflow_record_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the NetFlow Record.",
			},
			"netflow_exporter_uuids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of NetFlow Exporter UUIDs.",
			},
		},
	}
}

func dataSourceMSOTenantPoliciesNetflowMonitorRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO NetFlow Monitor Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	monitorName := d.Get("name").(string)

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	monitor, err := GetPolicyByName(response, monitorName, "tenantPolicyTemplate", "template", "netFlowMonitors")
	if err != nil {
		return err
	}

	err = setNetflowMonitorData(d, monitor, templateId)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] MSO NetFlow Monitor Data Source - Read Complete: %v", d.Id())
	return nil
}
