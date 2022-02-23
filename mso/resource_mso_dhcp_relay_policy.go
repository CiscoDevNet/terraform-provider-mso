package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSODHCPRelayPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSODHCPRelayPolicyCreate,
		Update: resourceMSODHCPRelayPolicyUpdate,
		Read:   resourceMSODHCPRelayPolicyRead,
		Delete: resourceMSODHCPRelayPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSODHCPRelayPolicyImport,
		},

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"tenant_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"dhcp_relay_policy_provider": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"epg": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"external_epg": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dhcp_server_address": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsIPAddress,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		}),
	}
}

func resourceMSODHCPRelayPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	msoClient := m.(*client.Client)
	id := d.Id()
	DHCPRelayPolicy, err := getDHCPRelayPolicy(msoClient, id)
	if err != nil {
		return nil, err
	}
	setDHCPRelayPolicy(DHCPRelayPolicy, d)
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func getDHCPRelayPolicy(client *client.Client, id string) (*models.DHCPRelayPolicy, error) {
	cont, err := client.ReadDHCPRelayPolicy(id)
	if err != nil {
		return nil, err
	}

	DHCPRelayPolicy, err := models.DHCPRelayPolicyFromContainer(cont)
	if err != nil {
		return nil, err
	}

	return DHCPRelayPolicy, nil
}

func setDHCPRelayPolicy(DHCPRelayPolicy *models.DHCPRelayPolicy, d *schema.ResourceData) {
	d.Set("description", DHCPRelayPolicy.Desc)
	d.Set("name", DHCPRelayPolicy.Name)
	d.Set("tenant_id", DHCPRelayPolicy.TenantID)
	providerList := make([]map[string]string, 0)
	for _, provider := range DHCPRelayPolicy.DHCPProvider {
		providerList = append(providerList, map[string]string{
			"external_epg":        provider.ExternalEPG,
			"epg":                 provider.EPG,
			"dhcp_server_address": provider.DHCPServerAddress,
			"tenant_id":           provider.TenantID,
		})
	}
	d.Set("dhcp_relay_policy_provider", providerList)
	d.SetId(DHCPRelayPolicy.ID)
}

func resourceMSODHCPRelayPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Create", d.Id())

	msoClient := m.(*client.Client)
	tenantId := d.Get("tenant_id").(string)

	DHCPRelayPolicy := models.DHCPRelayPolicy{
		TenantID: tenantId,
		Name:     d.Get("name").(string),
	}

	if desc, ok := d.GetOk("description"); ok {
		DHCPRelayPolicy.Desc = desc.(string)
	}

	if providerList, ok := d.GetOk("dhcp_relay_policy_provider"); ok {
		providerModelList := make([]models.DHCPProvider, 0)
		for _, provider := range providerList.([]interface{}) {
			providerMap := provider.(map[string]interface{})
			if providerMap["epg"] == "" && providerMap["external_epg"] == "" {
				return fmt.Errorf("expected any one of the epg or external_epg.")
			}
			if providerMap["epg"] != "" && providerMap["external_epg"] != "" {
				return fmt.Errorf("epg and external_epg both should not be set simultaneously")
			}
			providerModelList = append(providerModelList, models.DHCPProvider{
				ExternalEPG:       providerMap["external_epg"].(string),
				EPG:               providerMap["epg"].(string),
				DHCPServerAddress: providerMap["dhcp_server_address"].(string),
				TenantID:          tenantId,
			})
		}
		DHCPRelayPolicy.DHCPProvider = providerModelList
	}

	cont, err := msoClient.CreateDHCPRelayPolicy(&DHCPRelayPolicy)
	if err != nil {
		return err
	}
	d.SetId(models.StripQuotes(cont.S("id").String()))

	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceMSODHCPRelayPolicyRead(d, m)
}

func resourceMSODHCPRelayPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Policy Update", d.Id())

	msoClient := m.(*client.Client)

	tenantId := d.Get("tenant_id").(string)

	DHCPRelayPolicy := models.DHCPRelayPolicy{
		TenantID: tenantId,
		Name:     d.Get("name").(string),
	}

	if desc, ok := d.GetOk("description"); ok {
		DHCPRelayPolicy.Desc = desc.(string)
	}

	if providerList, ok := d.GetOk("dhcp_relay_policy_provider"); ok {
		providerModelList := make([]models.DHCPProvider, 0)
		for _, provider := range providerList.([]interface{}) {
			providerMap := provider.(map[string]interface{})
			if providerMap["epg"] == "" && providerMap["external_epg"] == "" {
				return fmt.Errorf("expected any one of the epg or external_epg.")
			}
			if providerMap["epg"] != "" && providerMap["external_epg"] != "" {
				return fmt.Errorf("epg and external_epg both should not be set simultaneously")
			}
			providerModelList = append(providerModelList, models.DHCPProvider{
				ExternalEPG:       providerMap["external_epg"].(string),
				EPG:               providerMap["epg"].(string),
				DHCPServerAddress: providerMap["dhcp_server_address"].(string),
				TenantID:          tenantId,
			})
		}
		DHCPRelayPolicy.DHCPProvider = providerModelList
	}

	_, err := msoClient.UpdateDHCPRelayPolicy(d.Id(), &DHCPRelayPolicy)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Policy Update finished successfully: %s", d.Id())

	return resourceMSODHCPRelayPolicyRead(d, m)
}

func resourceMSODHCPRelayPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	id := d.Id()
	log.Printf("id: %v\n", id)
	DHCPRelayPolicy, err := getDHCPRelayPolicy(msoClient, id)
	if err != nil {
		d.SetId("")
		return err
	}
	setDHCPRelayPolicy(DHCPRelayPolicy, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSODHCPRelayPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	msoClient := m.(*client.Client)
	id := d.Id()

	err := msoClient.DeleteDHCPRelayPolicy(id)
	if err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
