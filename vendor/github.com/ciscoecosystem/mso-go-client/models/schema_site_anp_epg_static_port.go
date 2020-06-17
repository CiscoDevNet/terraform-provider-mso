package models

type SchemaSiteAnpEpgStaticPort struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewSchemaSiteAnpEpgStaticPort(ops, path, Type, portPath string, vlan int, deploymentImmediacy string, microsegVlan int, mode string) *SchemaSiteAnpEpgStaticPort {
	var anpepgMap map[string]interface{}
	anpepgMap = map[string]interface{}{
		"type":                Type,
		"path":                portPath,
		"portEncapVlan":       vlan,
		"deploymentImmediacy": deploymentImmediacy,
		"microSegVlan":        microsegVlan,
		"mode":                mode,
	}

	if anpepgMap["deploymentImmediacy"] == "" {
		anpepgMap["deploymentImmediacy"] = "lazy"
	}

	if anpepgMap["mode"] == "" {
		anpepgMap["mode"] = "untagged"
	}

	if anpepgMap["microSegVlan"] == 0 {
		delete(anpepgMap, "microSegVlan")
	}

	return &SchemaSiteAnpEpgStaticPort{
		Ops:   ops,
		Path:  path,
		Value: anpepgMap,
	}

}

func (anpAttributes *SchemaSiteAnpEpgStaticPort) ToMap() (map[string]interface{}, error) {
	anpAttributesMap := make(map[string]interface{})
	A(anpAttributesMap, "op", anpAttributes.Ops)
	A(anpAttributesMap, "path", anpAttributes.Path)
	if anpAttributes.Value != nil {
		A(anpAttributesMap, "value", anpAttributes.Value)
	}

	return anpAttributesMap, nil
}
