package client

import (
	"fmt"

	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
)

func (client *Client) GetDHCPOptionPolicyID(name string) (string, error) {
	path := "api/v1/policies/dhcp/option"
	cont, err := client.GetViaURL(path)
	if err != nil {
		return "", err
	}
	for _, policy := range cont.S("DhcpRelayPolicies").Data().([]interface{}) {
		if optionPol, ok := policy.(map[string]interface{}); ok {
			if name == optionPol["name"].(string) {
				return optionPol["id"].(string), nil
			}
		}
	}
	return "", fmt.Errorf("DHCP Option Policy with name: %s not found", name)
}

func (client *Client) CreateDHCPOptionPolicy(obj *models.DHCPOptionPolicy) (*container.Container, error) {
	path := "api/v1/policies/dhcp/option"
	cont, err := client.Save(path, obj)
	if err != nil {
		return nil, err
	}
	return cont, nil
}

func (client *Client) ReadDHCPOptionPolicy(id string) (*container.Container, error) {
	path := "api/v1/policies/dhcp/option/" + id
	cont, err := client.GetViaURL(path)
	if err != nil {
		return nil, err
	}
	return cont, nil
}

func (client *Client) UpdateDHCPOptionPolicy(id string, obj *models.DHCPOptionPolicy) (*container.Container, error) {
	path := "api/v1/policies/dhcp/option/" + id
	cont, err := client.Put(path, obj)
	if err != nil {
		return nil, err
	}
	return cont, nil
}

func (client *Client) DeleteDHCPOptionPolicy(id string) error {
	path := "api/v1/policies/dhcp/option/" + id
	err := client.DeletebyId(path)
	if err != nil {
		return err
	}
	return nil
}
