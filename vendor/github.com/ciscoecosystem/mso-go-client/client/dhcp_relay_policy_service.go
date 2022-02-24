package client

import (
	"fmt"

	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
)

func (client *Client) GetDHCPRelayPolicyID(name string) (string, error) {
	path := "api/v1/policies/dhcp/relay"
	cont, err := client.GetViaURL(path)
	if err != nil {
		return "", err
	}
	for _, policy := range cont.S("DhcpRelayPolicies").Data().([]interface{}) {
		if relayPol, ok := policy.(map[string]interface{}); ok {
			if name == relayPol["name"].(string) {
				return relayPol["id"].(string), nil
			}
		}
	}
	return "", fmt.Errorf("DHCP Relay Policy with name: %s not found", name)
}

func (client *Client) CreateDHCPRelayPolicy(obj *models.DHCPRelayPolicy) (*container.Container, error) {
	path := "api/v1/policies/dhcp/relay"
	cont, err := client.Save(path, obj)
	if err != nil {
		return nil, err
	}
	return cont, nil
}

func (client *Client) ReadDHCPRelayPolicy(id string) (*container.Container, error) {
	path := "api/v1/policies/dhcp/relay/" + id
	cont, err := client.GetViaURL(path)
	if err != nil {
		return nil, err
	}
	return cont, nil
}

func (client *Client) UpdateDHCPRelayPolicy(id string, obj *models.DHCPRelayPolicy) (*container.Container, error) {
	remotePolicy, err := client.ReadDHCPRelayPolicy(id)
	if err != nil {
		return nil, err
	}

	payloadModel, err := models.PrepareDHCPRelayPolicyModelForUpdate(remotePolicy, obj)
	if err != nil {
		return nil, err
	}
	path := "api/v1/policies/dhcp/relay/" + id
	cont, err := client.Put(path, payloadModel)
	if err != nil {
		return nil, err
	}
	return cont, nil
}

func (client *Client) DeleteDHCPRelayPolicy(id string) error {
	path := "api/v1/policies/dhcp/relay/" + id
	err := client.DeletebyId(path)
	if err != nil {
		return err
	}
	return nil
}
