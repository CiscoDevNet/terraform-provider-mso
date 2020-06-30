package models

type SchemaTemplateAnpEpgSelector struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewSchemaTemplateAnpEpgSelector(ops, path string, selectorMap map[string]interface{}) *SchemaTemplateAnpEpgSelector {
	var temp map[string]interface{}

	if ops != "remove" {
		temp = selectorMap
	} else {
		temp = nil
	}

	return &SchemaTemplateAnpEpgSelector{
		Ops:   ops,
		Path:  path,
		Value: temp,
	}
}

func (schematemplateanpepgselectorattr *SchemaTemplateAnpEpgSelector) ToMap() (map[string]interface{}, error) {
	schematemplateanpepgselectorMap := make(map[string]interface{})

	A(schematemplateanpepgselectorMap, "op", schematemplateanpepgselectorattr.Ops)
	A(schematemplateanpepgselectorMap, "path", schematemplateanpepgselectorattr.Path)
	if schematemplateanpepgselectorattr.Value != nil {
		A(schematemplateanpepgselectorMap, "value", schematemplateanpepgselectorattr.Value)
	}

	return schematemplateanpepgselectorMap, nil
}
