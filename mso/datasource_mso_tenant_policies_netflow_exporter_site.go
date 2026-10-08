package mso

import (
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceMSONetflowExporterSite() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSONetflowExporterSiteRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the tenant policy template.",
			},
			"site_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the site on which the NetFlow Exporter is configured.",
			},
			"netflow_exporter_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the template-level NetFlow Exporter.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The site-specific description of the NetFlow Exporter.",
			},
			"source_ip_address_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source IP address type.",
			},
			"source_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The custom source IP address or prefix.",
			},
			"destination_port": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The destination port of the NetFlow Exporter.",
			},
			"destination_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IPv4 or IPv6 destination address.",
			},
			"dscp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The QoS DSCP value.",
			},
			"application_epg": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The application EPG associated with the NetFlow Exporter.",
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"tenant_name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The name of the tenant.",
					},
					"anp_name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The name of the application profile.",
					},
					"epg_name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The name of the EPG.",
					},
					"vrf_name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The name of the VRF.",
					},
					"vrf_tenant_name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The name of the tenant that contains the VRF.",
					},
				}},
			},
			"external_epg": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The external EPG associated with the NetFlow Exporter.",
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"tenant_name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The name of the tenant.",
					},
					"l3out_name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The name of the L3Out.",
					},
					"external_epg_name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The name of the external EPG.",
					},
					"vrf_name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The name of the VRF.",
					},
					"vrf_tenant_name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The name of the tenant that contains the VRF.",
					},
				}},
			},
		},
	}
}

func dataSourceMSONetflowExporterSiteRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO NetFlow Exporter Site Data Source - Beginning Read")
	msoClient := m.(*client.Client)
	templateID := d.Get("template_id").(string)
	siteID := d.Get("site_id").(string)
	exporterUUID := d.Get("netflow_exporter_uuid").(string)

	exporterCont, err := getNetflowExporterSite(msoClient, templateID, siteID, exporterUUID)
	if err != nil {
		return err
	}

	if err := setNetflowExporterSiteData(d, exporterCont, templateID, siteID, exporterUUID); err != nil {
		return err
	}
	log.Printf("[DEBUG] MSO NetFlow Exporter Site Data Source - Read Complete: %v", d.Id())
	return nil
}
