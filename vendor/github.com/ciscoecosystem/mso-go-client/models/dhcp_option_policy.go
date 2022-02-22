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

func PrepareDHCPOptionPolicyModelForUpdate(remotePolicyCont *container.Container, newPolicy *DHCPOptionPolicy) (*DHCPOptionPolicy, error) {
	remotePolicy := DHCPOptionPolicy{}
	err := json.Unmarshal(remotePolicyCont.Bytes(), &remotePolicy)
	if err != nil {
		return nil, err
	}

	newOptionList := make([]DHCPOption, 0)

	for _, newOption := range newPolicy.DHCPOption {
		if newOption.ID != "remove" {
			newOptionList = append(newOptionList, newOption)
		}
	}

	for _, remoteOption := range remotePolicy.DHCPOption {
		found := false
		for _, newOption := range newPolicy.DHCPOption {
			if newOption.Name == remoteOption.Name {
				found = true
			}
		}
		if !found {
			newOptionList = append(newOptionList, remoteOption)
		}
	}
	fmt.Printf("newOptionList: %v\n", newOptionList)
	newPolicy.DHCPOption = newOptionList
	return newPolicy, nil
}
