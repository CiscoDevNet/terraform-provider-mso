package models

import (
	"encoding/json"

	"github.com/ciscoecosystem/mso-go-client/container"
)

type DHCPRelayPolicy struct {
	ID            string         `json:"id,omitempty"`
	Name          string         `json:"name"`
	PolicyType    string         `json:"policyType,omitempty"`
	PolicySubtype string         `json:"policySubtype,omitempty"`
	Desc          string         `json:"desc"`
	TenantID      string         `json:"tenantId"`
	DHCPProvider  []DHCPProvider `json:"provider"`
}

func NewDHCPRelayPolicy(policy DHCPRelayPolicy) *DHCPRelayPolicy {
	newDHCPRelayPolicy := policy
	return &newDHCPRelayPolicy
}

type DHCPProvider struct {
	ExternalEPG       string `json:"externalEpgRef"`
	EPG               string `json:"epgRef"`
	DHCPServerAddress string `json:"addr"`
	TenantID          string `json:"tenantId"`
}

func (model *DHCPRelayPolicy) ToMap() (map[string]interface{}, error) {
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

func DHCPRelayPolicyFromContainer(cont *container.Container) (*DHCPRelayPolicy, error) {
	policy := DHCPRelayPolicy{}
	err := json.Unmarshal(cont.EncodeJSON(), &policy)
	if err != nil {
		return nil, err
	}
	return &policy, nil
}
