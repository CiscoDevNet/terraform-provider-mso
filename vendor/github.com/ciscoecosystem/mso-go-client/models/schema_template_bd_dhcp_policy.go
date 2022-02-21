package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/mso-go-client/container"
)

type TemplateBDDHCPPolicyOps struct {
	Ops   string                 `json:"op,omitempty"`
	Path  string                 `json:"path,omitempty"`
	Value map[string]interface{} `json:"value,omitempty"`
}

type TemplateBDDHCPPolicy struct {
	Name              string
	Version           int
	DHCPOptionName    string
	DHCPOptionVersion int
	BDName            string
	TemplateName      string
	SchemaID          string
}

func TemplateBDDHCPPolicyModelForCreation(bdDHCPPol *TemplateBDDHCPPolicy) *TemplateBDDHCPPolicyOps {
	opsMap := TemplateBDDHCPPolicyOps{
		Ops:  "add",
		Path: fmt.Sprintf("/templates/%s/bds/%s/dhcpLabels/-", bdDHCPPol.TemplateName, bdDHCPPol.BDName),
	}
	opsVal := map[string]interface{}{
		"name":    bdDHCPPol.Name,
		"version": bdDHCPPol.Version,
	}

	if bdDHCPPol.DHCPOptionName != "" {
		opsVal["dhcpOptionLabel"] = map[string]interface{}{
			"name": bdDHCPPol.DHCPOptionName,
			"version": func(v int) int {
				if v == 0 {
					return 1
				}
				return v
			}(bdDHCPPol.DHCPOptionVersion),
		}
	}
	opsMap.Value = opsVal
	return &opsMap
}

func TemplateBDDHCPPolicyModelForUpdate(bdDHCPPol *TemplateBDDHCPPolicy) *TemplateBDDHCPPolicyOps {
	opsMap := TemplateBDDHCPPolicyOps{
		Ops:  "replace",
		Path: fmt.Sprintf("/templates/%s/bds/%s/dhcpLabels/%s", bdDHCPPol.TemplateName, bdDHCPPol.BDName, bdDHCPPol.Name),
	}
	opsVal := map[string]interface{}{
		"name":    bdDHCPPol.Name,
		"version": bdDHCPPol.Version,
	}

	if bdDHCPPol.DHCPOptionName != "" {
		opsVal["dhcpOptionLabel"] = map[string]interface{}{
			"name": bdDHCPPol.DHCPOptionName,
			"version": func(v int) int {
				if v == 0 {
					return 1
				}
				return v
			}(bdDHCPPol.DHCPOptionVersion),
		}
	}
	opsMap.Value = opsVal
	return &opsMap
}

func TemplateBDDHCPPolicyModelForDeletion(bdDHCPPol *TemplateBDDHCPPolicy) *TemplateBDDHCPPolicyOps {
	opsMap := TemplateBDDHCPPolicyOps{
		Ops:  "remove",
		Path: fmt.Sprintf("/templates/%s/bds/%s/dhcpLabels/%s", bdDHCPPol.TemplateName, bdDHCPPol.BDName, bdDHCPPol.Name),
	}
	return &opsMap
}

func TemplateBDDHCPPolicyFromContainer(cont *container.Container, tf *TemplateBDDHCPPolicy) (*TemplateBDDHCPPolicy, error) {
	remoteBDDHCPPol := TemplateBDDHCPPolicy{}

	templateCont, err := cont.S("templates").SearchInObjectList(
		func(cont *container.Container) bool {
			return G(cont, "name") == tf.TemplateName
		},
	)
	if err != nil {
		return nil, err
	}

	bdCont, err := templateCont.S("bds").SearchInObjectList(
		func(cont *container.Container) bool {
			return G(cont, "name") == tf.BDName
		},
	)
	if err != nil {
		return nil, err
	}

	bdDHCPCont, err := bdCont.S("dhcpLabels").SearchInObjectList(
		func(cont *container.Container) bool {
			return G(cont, "name") == tf.Name
		},
	)
	if err != nil {
		return nil, err
	}

	remoteBDDHCPPol.Name = G(bdDHCPCont, "name")
	remoteBDDHCPPol.Version, err = strconv.Atoi(G(bdDHCPCont, "version"))
	if err != nil {
		return nil, err
	}
	if bdCont.Exists("dhcpOptionLabel") {
		remoteBDDHCPPol.DHCPOptionName = G(bdDHCPCont, "dhcpOptionLabel,name")
		remoteBDDHCPPol.DHCPOptionVersion, err = strconv.Atoi(G(bdDHCPCont, "dhcpOptionLabel,version"))
		if err != nil {
			return nil, err
		}
	}

	return &remoteBDDHCPPol, nil
}

func (templateBDDHCPPolicyOps *TemplateBDDHCPPolicyOps) ToMap() (map[string]interface{}, error) {
	templateBDDHCPPolicyOpsMap := make(map[string]interface{}, 0)
	A(templateBDDHCPPolicyOpsMap, "op", templateBDDHCPPolicyOps.Ops)
	A(templateBDDHCPPolicyOpsMap, "path", templateBDDHCPPolicyOps.Path)
	if templateBDDHCPPolicyOps.Value != nil {
		A(templateBDDHCPPolicyOpsMap, "value", templateBDDHCPPolicyOps.Value)
	}
	return templateBDDHCPPolicyOpsMap, nil
}
