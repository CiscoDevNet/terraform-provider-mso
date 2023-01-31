package client

import (
	"fmt"

	"github.com/ciscoecosystem/mso-go-client/models"
)

func (client *Client) CreateDHCPRelayPolicyProvider(obj *models.DHCPRelayPolicyProvider) error {
	relayPolicyId, err := client.GetDHCPRelayPolicyID(obj.PolicyName)
	if err != nil {
		return err
	}
	relayPolicyCont, err := client.ReadDHCPRelayPolicy(relayPolicyId)
	if err != nil {
		return err
	}
	DHCPRelay, err := models.DHCPRelayPolicyFromContainer(relayPolicyCont)
	if err != nil {
		return err
	}
	provider := models.DHCPProvider{
		ExternalEPG:       obj.ExternalEpgRef,
		EPG:               obj.EpgRef,
		DHCPServerAddress: obj.Addr,
		TenantID:          DHCPRelay.TenantID,
	}
	DHCPRelay.DHCPProvider = append(DHCPRelay.DHCPProvider, provider)
	_, err = client.UpdateDHCPRelayPolicy(relayPolicyId, DHCPRelay)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) UpdateDHCPRelayPolicyProvider(new *models.DHCPRelayPolicyProvider, old *models.DHCPRelayPolicyProvider) error {
	relayPolicyId, err := client.GetDHCPRelayPolicyID(old.PolicyName)
	if err != nil {
		return err
	}
	relayPolicyCont, err := client.ReadDHCPRelayPolicy(relayPolicyId)
	if err != nil {
		return err
	}
	DHCPRelay, err := models.DHCPRelayPolicyFromContainer(relayPolicyCont)
	if err != nil {
		return err
	}
	NewProviders := make([]models.DHCPProvider, 0, 1)
	NewProvider := models.DHCPProvider{
		ExternalEPG:       new.ExternalEpgRef,
		EPG:               new.EpgRef,
		DHCPServerAddress: new.Addr,
		TenantID:          DHCPRelay.TenantID,
	}
	for _, provider := range DHCPRelay.DHCPProvider {
		if provider.DHCPServerAddress != old.Addr && provider.EPG != old.EpgRef && old.ExternalEpgRef != new.ExternalEpgRef {
			NewProviders = append(NewProviders, provider)
		} else {
			NewProviders = append(NewProviders, NewProvider)
		}
	}
	DHCPRelay.DHCPProvider = NewProviders
	_, err = client.UpdateDHCPRelayPolicy(relayPolicyId, DHCPRelay)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) DeleteDHCPRelayPolicyProvider(obj *models.DHCPRelayPolicyProvider) error {
	relayPolicyId, err := client.GetDHCPRelayPolicyID(obj.PolicyName)
	if err != nil {
		return err
	}
	relayPolicyCont, err := client.ReadDHCPRelayPolicy(relayPolicyId)
	if err != nil {
		return err
	}
	DHCPRelay, err := models.DHCPRelayPolicyFromContainer(relayPolicyCont)
	if err != nil {
		return err
	}
	NewProviders := make([]models.DHCPProvider, 0, 1)
	for _, provider := range DHCPRelay.DHCPProvider {
		if provider.DHCPServerAddress == obj.Addr && provider.EPG == obj.EpgRef && provider.ExternalEPG == obj.ExternalEpgRef {
			provider.Operation = "remove"
		}
		NewProviders = append(NewProviders, provider)
	}
	DHCPRelay.DHCPProvider = NewProviders
	_, err = client.UpdateDHCPRelayPolicy(relayPolicyId, DHCPRelay)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) ReadDHCPRelayPolicyProvider(obj *models.DHCPRelayPolicyProvider) (*models.DHCPRelayPolicyProvider, error) {
	relayPolicyId, err := client.GetDHCPRelayPolicyID(obj.PolicyName)
	if err != nil {
		return nil, err
	}
	relayPolicyCont, err := client.ReadDHCPRelayPolicy(relayPolicyId)
	if err != nil {
		return nil, err
	}
	DHCPRelay, err := models.DHCPRelayPolicyFromContainer(relayPolicyCont)
	if err != nil {
		return nil, err
	}
	flag := false
	for _, provider := range DHCPRelay.DHCPProvider {
		if provider.DHCPServerAddress == obj.Addr && provider.EPG == obj.EpgRef && provider.ExternalEPG == obj.ExternalEpgRef {
			flag = true
			break
		}
	}
	if flag {
		return obj, nil
	}
	return nil, fmt.Errorf("no DHCP Relay Policy found")
}
