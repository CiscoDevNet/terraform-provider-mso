package client

import (
	"fmt"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/models"
)

func (client *Client) CreateDHCPOptionPolicyOption(obj *models.DHCPOptionPolicyOption) error {
	optionPolicyID, err := client.GetDHCPOptionPolicyID(obj.PolicyName)
	if err != nil {
		return err
	}
	optionPolicyCont, err := client.ReadDHCPOptionPolicy(optionPolicyID)
	if err != nil {
		return err
	}
	DHCPOptionPolicy, err := models.DHCPOptionPolicyFromContainer(optionPolicyCont)
	if err != nil {
		return err
	}

	option := models.DHCPOption{
		Data: obj.Data,
		ID:   obj.ID,
		Name: obj.Name,
	}

	DHCPOptionPolicy.DHCPOption = append(DHCPOptionPolicy.DHCPOption, option)
	_, err = client.UpdateDHCPOptionPolicy(optionPolicyID, DHCPOptionPolicy)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) ReadDHCPOptionPolicyOption(id string) (*models.DHCPOptionPolicyOption, error) {
	idSplit := strings.Split(id, "/")
	optionPolicyID, err := client.GetDHCPOptionPolicyID(idSplit[0])
	if err != nil {
		return nil, err
	}
	optionPolicyCont, err := client.ReadDHCPOptionPolicy(optionPolicyID)
	if err != nil {
		return nil, err
	}
	DHCPOptionPolicy, err := models.DHCPOptionPolicyFromContainer(optionPolicyCont)
	if err != nil {
		return nil, err
	}

	flag := false
	dhcpOption := models.DHCPOptionPolicyOption{}
	for _, option := range DHCPOptionPolicy.DHCPOption {
		if option.Name == idSplit[1] {
			flag = true
			dhcpOption.Name = option.Name
			dhcpOption.ID = option.ID
			dhcpOption.Data = option.Data
			dhcpOption.PolicyName = DHCPOptionPolicy.Name
			break
		}
	}
	if flag {
		return &dhcpOption, nil
	}
	return nil, fmt.Errorf("No DHCP Option Policy found")
}

func (client *Client) UpdateDHCPOptionPolicyOption(obj *models.DHCPOptionPolicyOption) error {
	optionPolicyID, err := client.GetDHCPOptionPolicyID(obj.PolicyName)
	if err != nil {
		return err
	}
	optionPolicyCont, err := client.ReadDHCPOptionPolicy(optionPolicyID)
	if err != nil {
		return err
	}
	DHCPOptionPolicy, err := models.DHCPOptionPolicyFromContainer(optionPolicyCont)
	if err != nil {
		return err
	}

	NewOptions := make([]models.DHCPOption, 0, 1)
	NewOption := models.DHCPOption{
		Data: obj.Data,
		ID:   obj.ID,
		Name: obj.Name,
	}

	for _, option := range DHCPOptionPolicy.DHCPOption {
		if option.Name != obj.Name {
			NewOptions = append(NewOptions, option)
		} else {
			NewOptions = append(NewOptions, NewOption)
		}
	}

	DHCPOptionPolicy.DHCPOption = NewOptions
	_, err = client.UpdateDHCPOptionPolicy(optionPolicyID, DHCPOptionPolicy)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) DeleteDHCPOptionPolicyOption(id string) error {
	idSplit := strings.Split(id, "/")
	optionPolicyID, err := client.GetDHCPOptionPolicyID(idSplit[0])
	if err != nil {
		return err
	}
	optionPolicyCont, err := client.ReadDHCPOptionPolicy(optionPolicyID)
	if err != nil {
		return err
	}
	DHCPOptionPolicy, err := models.DHCPOptionPolicyFromContainer(optionPolicyCont)
	if err != nil {
		return err
	}
	NewOptions := make([]models.DHCPOption, 0, 1)
	for _, option := range DHCPOptionPolicy.DHCPOption {
		if option.Name != idSplit[1] {
			NewOptions = append(NewOptions, option)
		}
	}
	DHCPOptionPolicy.DHCPOption = NewOptions
	_, err = client.UpdateDHCPOptionPolicy(optionPolicyID, DHCPOptionPolicy)
	if err != nil {
		return err
	}
	return nil
}
