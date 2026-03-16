package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSONetflowExporter() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSONetflowExporterRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the tenant policy template.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the NetFlow Exporter.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the NetFlow Exporter.",
			},
		},
	}
}

func dataSourceMSONetflowExporterRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO NetFlow Exporter Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "netFlowExporters")
	if err != nil {
		return err
	}

	setNetflowExporterData(d, policy, templateId)
	log.Printf("[DEBUG] MSO NetFlow Exporter Data Source - Read Complete: %v", d.Id())
	return nil
}