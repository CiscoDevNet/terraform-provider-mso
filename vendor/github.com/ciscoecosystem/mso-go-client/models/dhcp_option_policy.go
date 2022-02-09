package models

import (
	"encoding/json"

	"github.com/ciscoecosystem/mso-go-client/container"
)

type DHCPOptionPolicy struct {
	ID            string       `json:"id,omitempty"`
	Name          string       `json:"name"`
	PolicyType    string       `json:"policyType,omitempty"`
	PolicySubtype string       `json:"policySubtype,omitempty"`
	Desc          string       `json:"desc"`
	TenantID      string       `json:"tenantId"`
	DHCPOption    []DHCPOption `json:"dhcpOption"`
}

func NewDHCPOptionPolicy(policy DHCPOptionPolicy) *DHCPOptionPolicy {
	newDHCPOptionPolicy := policy
	return &newDHCPOptionPolicy
}

type DHCPOption struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Data string `json:"data"`
}

func (model *DHCPOptionPolicy) ToMap() (map[string]interface{}, error) {
	objMap := make(map[string]interface{})

	jsonObj, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonObj, &objMap)
	if err != nil {
		return nil, err
	}

	return objMap, nil
}

func DHCPOptionPolicyFromContainer(cont *container.Container) (*DHCPOptionPolicy, error) {
	policy := DHCPOptionPolicy{}

	err := json.Unmarshal(cont.EncodeJSON(), &policy)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}
