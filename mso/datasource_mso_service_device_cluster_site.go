package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceMSOServiceDeviceClusterSite() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOServiceDeviceClusterSiteRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the service device template that contains the Service Device Cluster.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Service Device Cluster configured on the site.",
			},
			"site_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the site on which the Service Device Cluster is configured.",
			},

			"domain_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of domain associated with the Service Device Cluster on the site.",
			},
			"vmm_domain_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The VMM domain provider type when `domain_type` is `vmmDomain`.",
			},
			"domain_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the domain associated with the Service Device Cluster on the site.",
			},
			"domain_dn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The distinguished name of the domain associated with the Service Device Cluster on the site.",
			},

			"high_availability_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The high availability mode of the Service Device Cluster on the site.",
			},
			"promiscuous_mode": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether promiscuous mode is enabled on the Service Device Cluster on the site.",
			},
			"trunking_port": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the Service Device Cluster on the site uses a trunking port.",
			},
			"vlan": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The VLAN ID on the Service Device Cluster on the site. Populated at the device level when `high_availability_mode` is `activeStandby`.",
			},

			"interfaces": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of interfaces of the Service Device Cluster on the site.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the interface.",
						},
						"vlan": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The VLAN ID of the interface.",
						},
						"domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the physical domain associated with the interface, populated in `activeActive` high availability mode. Only physical domains are supported at interface scope.",
						},
						"domain_dn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The distinguished name of the physical domain associated with the interface, populated in `activeActive` high availability mode.",
						},
						"fabric_to_device_connectivity": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The list of fabric-to-device connectivity paths for the interface.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pod_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The pod ID of the fabric path.",
									},
									"node_id": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "The node ID(s) of the fabric path. A single element for `port_type` `port` and `dpc`, two elements for `port_type` `vpc`.",
									},
									"path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The path on the node.",
									},
									"port_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of port used for the fabric path.",
									},
									"tag": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The tag of the fabric path.",
									},
									"vlan": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The VLAN ID carried on this fabric path. Populated when `high_availability_mode` is `activeActive`.",
									},
								},
							},
						},
						"vm_information": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The list of VM information entries for the interface.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vm_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the VM.",
									},
									"vnic_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the vNIC on the VM.",
									},
									"pod_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The pod ID of the fabric path the VM interface attaches to.",
									},
									"node_id": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "The node ID(s) of the fabric path the VM interface attaches to.",
									},
									"path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The path on the node the VM interface attaches to.",
									},
									"port_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of port used for the VM interface's fabric path.",
									},
								},
							},
						},
						"enhanced_lag_policy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The UUID of the enhanced LAG policy associated with the interface.",
						},
						"pbr_destinations": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The list of policy-based redirect (PBR) destinations for the interface.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IP address of the PBR destination.",
									},
									"mac": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The MAC address of the PBR destination.",
									},
									"pod_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The pod ID of the PBR destination.",
									},
									"additional_tracking_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The additional IP address used for tracking the PBR destination.",
									},
									"weight": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The weight of the PBR destination.",
									},
									"is_backup": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether the PBR destination is a backup destination.",
									},
									"tag": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The tag of the PBR destination.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceMSOServiceDeviceClusterSiteRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Service Device Cluster Site Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	siteId := d.Get("site_id").(string)
	name := d.Get("name").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	siteIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "siteId", siteId, "deviceTemplate", "sites")
	if err != nil {
		return fmt.Errorf("site %q is not present in template %q: %v", siteId, templateId, err)
	}

	siteCont := templateCont.S("deviceTemplate", "sites").Index(siteIndex)

	deviceCont, err := GetPolicyByName(siteCont, name, "devices")
	if err != nil {
		return fmt.Errorf("device %q not found on site %q in template %q: %v", name, siteId, templateId, err)
	}

	if err := setServiceDeviceClusterSiteData(d, deviceCont, templateId, siteId); err != nil {
		return err
	}
	log.Printf("[DEBUG] MSO Service Device Cluster Site Data Source - Read Complete: %v", d.Id())
	return nil
}
