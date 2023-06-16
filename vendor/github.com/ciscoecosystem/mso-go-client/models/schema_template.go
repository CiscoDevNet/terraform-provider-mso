package models

type SchemaTemplate struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewSchemaTemplate(ops, path, tenantId, templateName, templateDisplayName, description, templateType string, templateSubTypes []string) *SchemaTemplate {
	var templateMap map[string]interface{}
	if ops != "remove" {
		templateMap = map[string]interface{}{
			"tenantId":        tenantId,
			"name":            templateName,
			"displayName":     templateDisplayName,
			"templateType":    templateType,
			"templateSubType": templateSubTypes,
			"description":     description,
			"anps":            []interface{}{},
			"bds":             []interface{}{},
			"contracts":       []interface{}{},
			"externalEpgs":    []interface{}{},
			"filters":         []interface{}{},
			"serviceGraphs":   []interface{}{},
			"vrfs":            []interface{}{},
			"intersiteL3outs": []interface{}{},
		}
	} else {
		templateMap = nil
	}

	return &SchemaTemplate{
		Ops:   ops,
		Path:  path,
		Value: templateMap,
	}

}

func (schematemplateAttributes *SchemaTemplate) ToMap() (map[string]interface{}, error) {
	schematemplateAttributeMap := make(map[string]interface{})
	A(schematemplateAttributeMap, "op", schematemplateAttributes.Ops)
	A(schematemplateAttributeMap, "path", schematemplateAttributes.Path)
	if schematemplateAttributes.Value != nil {
		A(schematemplateAttributeMap, "value", schematemplateAttributes.Value)
	}

	return schematemplateAttributeMap, nil
}
