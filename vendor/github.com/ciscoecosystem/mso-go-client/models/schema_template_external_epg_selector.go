package models

type SchemaTemplateExternalEPGSelector struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewSchemaTemplateExternalEPGSelector(ops, path string, extrepgSelectorMap map[string]interface{}) *SchemaTemplateExternalEPGSelector {
	var temp map[string]interface{}

	if ops != "remove" {
		temp = extrepgSelectorMap
	} else {
		temp = nil
	}

	return &SchemaTemplateExternalEPGSelector{
		Ops:   ops,
		Path:  path,
		Value: temp,
	}
}

func (schemaTempleteExtrEPGSelectorattr *SchemaTemplateExternalEPGSelector) ToMap() (map[string]interface{}, error) {
	schemaTempleteExtrEPGSelectorMap := make(map[string]interface{})

	A(schemaTempleteExtrEPGSelectorMap, "op", schemaTempleteExtrEPGSelectorattr.Ops)
	A(schemaTempleteExtrEPGSelectorMap, "path", schemaTempleteExtrEPGSelectorattr.Path)
	if schemaTempleteExtrEPGSelectorattr.Value != nil {
		A(schemaTempleteExtrEPGSelectorMap, "value", schemaTempleteExtrEPGSelectorattr.Value)
	}

	return schemaTempleteExtrEPGSelectorMap, nil
}
