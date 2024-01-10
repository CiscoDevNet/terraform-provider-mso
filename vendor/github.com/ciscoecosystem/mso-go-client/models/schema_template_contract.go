package models

type TemplateContract struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewTemplateContract(ops, path, name, displayName, scope, filterType, targetDscp, priority, desc string, filterRelationships, filterRelationshipsProviderToConsumer, filterRelationshipsConsumerToProvider []interface{}) *TemplateContract {
	contractMap := map[string]interface{}{
		"name":                                  name,
		"displayName":                           displayName,
		"description":                           desc,
		"scope":                                 scope,
		"filterType":                            filterType,
		"filterRelationships":                   filterRelationships,
		"filterRelationshipsProviderToConsumer": filterRelationshipsProviderToConsumer,
		"filterRelationshipsConsumerToProvider": filterRelationshipsConsumerToProvider,
	}

	if contractMap["filterType"] == "" {
		contractMap["filterType"] = "bothWay"
	}

	if contractMap["scope"] == "" {
		contractMap["scope"] = "context"
	}

	if priority != "" {
		contractMap["prio"] = priority
	}

	if targetDscp != "" {
		contractMap["targetDscp"] = targetDscp
	}

	// If displayName is not set, set it to name because error will be raised if displayName is empty
	if displayName == "" {
		contractMap["displayName"] = name
	}

	return &TemplateContract{
		Ops:   ops,
		Path:  path,
		Value: contractMap,
	}

}

func (anpAttributes *TemplateContract) ToMap() (map[string]interface{}, error) {
	anpAttributesMap := make(map[string]interface{})
	A(anpAttributesMap, "op", anpAttributes.Ops)
	A(anpAttributesMap, "path", anpAttributes.Path)
	if anpAttributes.Value != nil {
		A(anpAttributesMap, "value", anpAttributes.Value)
	}

	return anpAttributesMap, nil
}
