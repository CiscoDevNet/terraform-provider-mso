package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceMSOSchemaSiteContractServiceGraphListener() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOSchemaSiteContractServiceGraphListenerRead,
		Schema: map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"contract_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"site_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"service_node_index": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"listener_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"protocol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"security_policy": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssl_certificates": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_dn": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_store": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"frontend_ip_dn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"rules": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"floating_ip": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"priority": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"host": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"action": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						// TODO: Should be uncommented once condition is configured through UI
						// "condition": &schema.Schema{
						// 	Type:     schema.TypeString,
						// 	Computed: true,
						// },
						"action_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"content_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"protocol": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider_epg_ref": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"schema_id": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"template_name": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"anp_name": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"epg_name": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"url_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"custom_url": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"redirect_host_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"redirect_path": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"redirect_query": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"response_code": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"response_body": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"redirect_protocol": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"redirect_port": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"redirect_code": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": &schema.Schema{
										Type:     schema.TypeInt,
										Computed: true,
									},
									"protocol": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"path": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"interval": &schema.Schema{
										Type:     schema.TypeInt,
										Computed: true,
									},
									"timeout": &schema.Schema{
										Type:     schema.TypeInt,
										Computed: true,
									},
									"unhealthy_threshold": &schema.Schema{
										Type:     schema.TypeInt,
										Computed: true,
									},
									"use_host_from_rule": &schema.Schema{
										Type:     schema.TypeBool,
										Computed: true,
									},
									"success_code": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"host": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"target_ip_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceMSOSchemaSiteContractServiceGraphListenerRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Beginning datasource Read")

	schemaID := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	contractName := d.Get("contract_name").(string)

	msoClient := m.(*client.Client)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return err
	}

	siteID := d.Get("site_id").(string)
	serviceNodeIndex := d.Get("service_node_index").(int)
	listenerName := d.Get("listener_name").(string)

	d.SetId(fmt.Sprintf("%s/sites/%s/templates/%s/contracts/%s/serviceNodes/%d/listeners/%s", schemaID, siteID, templateName, contractName, serviceNodeIndex, listenerName))

	if err != nil {
		return errorForObjectNotFound(err, d.Id(), cont, d)
	}

	err = setSchemaSiteContractServiceGraphListenerAttrs(cont, d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Datasource read finished successfully", d.Id())
	return nil
}
