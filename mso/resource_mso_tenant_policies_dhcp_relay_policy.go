package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOTenantPoliciesDHCPRelayPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOTenantPoliciesDHCPRelayPolicyCreate,
		Read:   resourceMSOTenantPoliciesDHCPRelayPolicyRead,
		Update: resourceMSOTenantPoliciesDHCPRelayPolicyUpdate,
		Delete: resourceMSOTenantPoliciesDHCPRelayPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOTenantPoliciesDHCPRelayPolicyImport,
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dhcp_relay_providers": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dhcp_server_address": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsIPAddress,
						},
						"application_epg_uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"external_epg_uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dhcp_server_vrf_preference": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func checkDHCPRelayPolicyProviders(d *schema.ResourceData) error {
	errors := make([]string, 0)
	for index, provider := range d.Get("dhcp_relay_providers").(*schema.Set).List() {
		providerMap := provider.(map[string]interface{})
		applicationEpgUuid := providerMap["application_epg_uuid"].(string)
		externalEpgUuid := providerMap["external_epg_uuid"].(string)
		if applicationEpgUuid != "" && externalEpgUuid != "" {
			errors = append(errors, fmt.Sprintf("\nError: Set either 'application_epg_uuid' or 'external_epg_uuid', not both for a provider at index position: %d\n", index))
		} else if applicationEpgUuid == "" && externalEpgUuid == "" {
			errors = append(errors, fmt.Sprintf("\nError: Please set either 'application_epg_uuid' or 'external_epg_uuid' for a provider at index position: %d\n", index))
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("%v", errors)
	}
	return nil
}

func setDHCPRelayPolicyData(d *schema.ResourceData, response *container.Container, templateId string) error {
	d.SetId(fmt.Sprintf("templateId/%s/DHCPRelayPolicy/%s", templateId, models.StripQuotes(response.S("name").String())))
	d.Set("template_id", templateId)
	d.Set("name", models.StripQuotes(response.S("name").String()))
	d.Set("description", models.StripQuotes(response.S("description").String()))
	d.Set("uuid", models.StripQuotes(response.S("uuid").String()))

	providersList, err := response.S("providers").Children()
	if err != nil {
		return err
	}
	dhcpRelayProviders := make([]interface{}, 0)
	for _, provider := range providersList {
		providerMap := map[string]interface{}{
			"dhcp_server_address":        models.StripQuotes(provider.S("ip").String()),
			"dhcp_server_vrf_preference": provider.S("useServerVrf").Data().(bool),
		}

		applicationEPG := models.StripQuotes(provider.S("epgRef").String())
		if applicationEPG != "{}" {
			providerMap["application_epg_uuid"] = applicationEPG
		}

		externalEPG := models.StripQuotes(provider.S("externalEpgRef").String())
		if externalEPG != "{}" {
			providerMap["external_epg_uuid"] = externalEPG
		}

		dhcpRelayProviders = append(dhcpRelayProviders, providerMap)
	}
	d.Set("dhcp_relay_providers", dhcpRelayProviders)
	return nil
}

func getProvidersMapList(dhcpRelayProviders []interface{}) []map[string]interface{} {
	providersMapList := make([]map[string]interface{}, 0)
	for _, provider := range dhcpRelayProviders {
		providerData := provider.(map[string]interface{})
		providerPayload := map[string]interface{}{
			"ip":           providerData["dhcp_server_address"].(string),
			"useServerVrf": providerData["dhcp_server_vrf_preference"].(bool),
		}

		applicationEpgUuid := providerData["application_epg_uuid"].(string)
		if applicationEpgUuid != "{}" {
			providerPayload["epgRef"] = applicationEpgUuid
		}

		externalEpgUuid := providerData["external_epg_uuid"].(string)
		if externalEpgUuid != "{}" {
			providerPayload["externalEpgRef"] = externalEpgUuid
		}

		providersMapList = append(providersMapList, providerPayload)
	}
	return providersMapList
}

func resourceMSOTenantPoliciesDHCPRelayPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO DHCP Relay Policy Resource - Beginning Create: %v", d.Id())

	err := checkDHCPRelayPolicyProviders(d)
	if err != nil {
		return err
	}

	msoClient := m.(*client.Client)

	payload := map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"providers":   []interface{}{},
	}

	payload["providers"] = getProvidersMapList(d.Get("dhcp_relay_providers").(*schema.Set).List())
	payloadModel := models.GetPatchPayload("add", "/tenantPolicyTemplate/template/dhcpRelayPolicies/-", payload)
	templateId := d.Get("template_id").(string)

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/DHCPRelayPolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO DHCP Relay Policy Resource - Create Complete: %v", d.Id())
	return resourceMSOTenantPoliciesDHCPRelayPolicyRead(d, m)
}

func resourceMSOTenantPoliciesDHCPRelayPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO DHCP Relay Policy Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return err
	}

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyName, err := GetPolicyNameFromResourceId(d.Id(), "DHCPRelayPolicy")
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "dhcpRelayPolicies")
	if err != nil {
		return err
	}

	setDHCPRelayPolicyData(d, policy, templateId)
	log.Printf("[DEBUG] MSO DHCP Relay Policy Resource - Read Complete : %v", d.Id())
	return nil
}

func resourceMSOTenantPoliciesDHCPRelayPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO DHCP Relay Policy Resource - Beginning Update: %v", d.Id())

	err := checkDHCPRelayPolicyProviders(d)
	if err != nil {
		return err
	}

	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "dhcpRelayPolicies")
	if err != nil {
		return err
	}

	payloadCon := container.New()
	payloadCon.Array()

	dhcpRelayPolicyPath := fmt.Sprintf("/tenantPolicyTemplate/template/dhcpRelayPolicies/%d", policyIndex)
	if d.HasChange("name") {
		err := addPatchPayloadToContainer(payloadCon, "replace", fmt.Sprintf("%s/name", dhcpRelayPolicyPath), d.Get("name"))
		if err != nil {
			return err
		}
	}

	if d.HasChange("description") {
		err := addPatchPayloadToContainer(payloadCon, "replace", fmt.Sprintf("%s/description", dhcpRelayPolicyPath), d.Get("description"))
		if err != nil {
			return err
		}
	}

	if d.HasChange("dhcp_relay_providers") {
		err := addPatchPayloadToContainer(payloadCon, "replace", fmt.Sprintf("%s/providers", dhcpRelayPolicyPath), getProvidersMapList(d.Get("dhcp_relay_providers").(*schema.Set).List()))
		if err != nil {
			return err
		}
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCon)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/DHCPRelayPolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO DHCP Relay Policy Resource - Update Complete: %v", d.Id())
	return resourceMSOTenantPoliciesDHCPRelayPolicyRead(d, m)
}

func resourceMSOTenantPoliciesDHCPRelayPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO DHCP Relay Policy Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", d.Get("template_id").(string)))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "dhcpRelayPolicies")
	if err != nil {
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/tenantPolicyTemplate/template/dhcpRelayPolicies/%d", policyIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", d.Get("template_id").(string)), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO DHCP Relay Policy Resource - Delete Complete: %v", d.Id())
	return nil
}

func resourceMSOTenantPoliciesDHCPRelayPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO DHCP Relay Policy Resource - Beginning Import: %v", d.Id())
	resourceMSOTenantPoliciesDHCPRelayPolicyRead(d, m)
	log.Printf("[DEBUG] MSO DHCP Relay Policy Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}
