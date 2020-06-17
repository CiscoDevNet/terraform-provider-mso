package models

type TemplateContract struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewTemplateContract(ops, path, name, displayName, scope, filterType string, filterRelationships []interface{}) *TemplateContract {
	var contractMap map[string]interface{}
	contractMap = map[string]interface{}{
		"name":                                  name,
		"displayName":                           displayName,
		"scope":                                 scope,
		"filterType":                            filterType,
		"filterRelationships":                   filterRelationships,
		"filterRelationshipsProviderToConsumer": []interface{}{},
		"filterRelationshipsConsumerToProvider": []interface{}{},
	}

	if contractMap["filterType"] == "" {
		contractMap["filterType"] = "bothWay"
	}

	if contractMap["scope"] == "" {
		contractMap["scope"] = "context"
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
