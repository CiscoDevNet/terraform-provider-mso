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
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"epg": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: epgRefValidation(),
						},
						"external_epg": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: externalEpgRefValidation(),
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
	tfProviderList := make([]map[string]string, 0)
	if _, ok := d.GetOk("tenant_id"); ok {
		providerList := d.Get("dhcp_relay_policy_provider").(*schema.Set).List()
		for _, provider := range providerList {
			providerMap := provider.(map[string]interface{})
			for _, remoteProvider := range DHCPRelayPolicy.DHCPProvider {
				if providerMap["external_epg"].(string) == remoteProvider.ExternalEPG && providerMap["epg"] == remoteProvider.EPG && providerMap["dhcp_server_address"] == remoteProvider.DHCPServerAddress {
					tfProviderList = append(tfProviderList, map[string]string{
						"epg":                 remoteProvider.EPG,
						"external_epg":        remoteProvider.ExternalEPG,
						"dhcp_server_address": remoteProvider.DHCPServerAddress,
						"tenant_id":           remoteProvider.TenantID,
					})
				}
			}
		}
	} else {
		for _, provider := range DHCPRelayPolicy.DHCPProvider {
			tfProviderList = append(tfProviderList, map[string]string{
				"external_epg":        provider.ExternalEPG,
				"epg":                 provider.EPG,
				"dhcp_server_address": provider.DHCPServerAddress,
				"tenant_id":           provider.TenantID,
			})
		}
	}
	d.Set("tenant_id", DHCPRelayPolicy.TenantID)
	d.Set("dhcp_relay_policy_provider", tfProviderList)
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
		err := ValidateProviderList(providerList.([]interface{}))
		if err != nil {
			return err
		}
		providerModelList := make([]models.DHCPProvider, 0)
		for _, provider := range providerList.(*schema.Set).List() {
			providerMap := provider.(map[string]interface{})
			if providerMap["epg"] == "" && providerMap["external_epg"] == "" {
				return fmt.Errorf("expected any one of the epg or external_epg")
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

	providerModelList := make([]models.DHCPProvider, 0)
	if d.HasChange("dhcp_relay_policy_provider") {
		providerList := d.Get("dhcp_relay_policy_provider").(*schema.Set).List()
		for _, provider := range providerList {
			providerMap := provider.(map[string]interface{})
			if providerMap["epg"] == "" && providerMap["external_epg"] == "" {
				return fmt.Errorf("expected any one of the epg or external_epg")
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
		err := ValidateProviderList(providerList)
		if err != nil {
			return err
		}
		oldProviders, newProviders := d.GetChange("dhcp_relay_policy_provider")
		oldProvidersList := oldProviders.(*schema.Set).List()
		newProvidersList := newProviders.(*schema.Set).List()

		providerModelList := make([]models.DHCPProvider, 0)
		oldProviderHashMap := make(map[string]int, 0)

		for i, v := range oldProvidersList {
			val := v.(map[string]interface{})
			oldProviderHashMap[fmt.Sprintf("%s%s%s", val["dhcp_server_address"].(string), val["epg"].(string), val["external_epg"].(string))] = i
		}
		newProviderHashMap := make(map[string]int, 0)
		for i, v := range newProvidersList {
			val := v.(map[string]interface{})
			newProviderHashMap[fmt.Sprintf("%s%s%s", val["dhcp_server_address"].(string), val["epg"].(string), val["external_epg"].(string))] = i
		}

		for k, i := range oldProviderHashMap {
			if _, ok := newProviderHashMap[k]; !ok {
				providerMap := oldProvidersList[i].(map[string]interface{})
				providerModelList = append(providerModelList, models.DHCPProvider{
					ExternalEPG:       providerMap["external_epg"].(string),
					EPG:               providerMap["epg"].(string),
					DHCPServerAddress: providerMap["dhcp_server_address"].(string),
					TenantID:          providerMap["tenant_id"].(string),
					Operation:         "remove",
				})
			}
		}

		for _, provider := range newProvidersList {
			providerMap := provider.(map[string]interface{})
			providerModelList = append(providerModelList, models.DHCPProvider{
				ExternalEPG:       providerMap["external_epg"].(string),
				EPG:               providerMap["epg"].(string),
				DHCPServerAddress: providerMap["dhcp_server_address"].(string),
				TenantID:          providerMap["tenant_id"].(string),
			})
		}

		log.Printf("providerModelList: %v\n", providerModelList)
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

func ValidateProviderList(providers []interface{}) error {
	idMap := make(map[string]bool, 0)
	for _, provider := range providers {
		providerMap := provider.(map[string]interface{})
		id := fmt.Sprintf("ip-%s/epg-%s/extepg-%s", providerMap["dhcp_server_address"].(string), providerMap["epg"].(string), providerMap["external_epg"].(string))
		if _, ok := idMap[id]; ok {
			return fmt.Errorf("duplicate provider configuration")
		}
		idMap[id] = true
	}
	return nil
}
