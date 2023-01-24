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
	Operation         string `json:"-"`
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

func PrepareDHCPRelayPolicyModelForUpdate(remotePolicyCont *container.Container, newPolicy *DHCPRelayPolicy) (*DHCPRelayPolicy, error) {
	remotePolicy := DHCPRelayPolicy{}
	err := json.Unmarshal(remotePolicyCont.Bytes(), &remotePolicy)
	if err != nil {
		return nil, err
	}

	newProviderList := make([]DHCPProvider, 0)

	for _, newProvider := range newPolicy.DHCPProvider {
		if newProvider.Operation != "remove" {
			newProviderList = append(newProviderList, newProvider)
		}
	}

	for _, remoteProvider := range remotePolicy.DHCPProvider {
		found := false
		for _, newProvider := range newPolicy.DHCPProvider {
			if remoteProvider.DHCPServerAddress == newProvider.DHCPServerAddress && remoteProvider.EPG == newProvider.EPG && remoteProvider.ExternalEPG == newProvider.ExternalEPG {
				found = true
			}
		}
		if !found {
			newProviderList = append(newProviderList, remoteProvider)
		}
	}

	newPolicy.DHCPProvider = newProviderList
	return newPolicy, nil
}
