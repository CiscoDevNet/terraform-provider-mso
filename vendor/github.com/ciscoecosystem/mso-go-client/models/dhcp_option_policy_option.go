package models

import (
	"encoding/json"

	"github.com/ciscoecosystem/mso-go-client/container"
)

type DHCPOptionPolicyOption struct {
	ID         string
	Name       string
	Data       string
	PolicyName string
}

func DHCPOptionPolicyOptionFromContainer(cont *container.Container) (*DHCPOptionPolicyOption, error) {
	option := DHCPOptionPolicyOption{}

	err := json.Unmarshal(cont.EncodeJSON(), &option)
	if err != nil {
		return nil, err
	}

	return &option, nil
}
