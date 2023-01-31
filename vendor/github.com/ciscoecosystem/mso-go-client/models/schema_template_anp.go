package models

type SchemaTemplateAnp struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewSchemaTemplateAnp(ops, path, Name, displayName string) *SchemaTemplateAnp {
	var VrfMap map[string]interface{}

	if ops != "remove" {
		VrfMap = map[string]interface{}{
			"displayName": displayName,
			"name":        Name,
			"epgs":        []interface{}{},
		}
	} else {

		VrfMap = nil
	}

	return &SchemaTemplateAnp{
		Ops:   ops,
		Path:  path,
		Value: VrfMap,
	}

}

func (schematemplateanpAttributes *SchemaTemplateAnp) ToMap() (map[string]interface{}, error) {
	schematemplateanpAttributeMap := make(map[string]interface{})
	A(schematemplateanpAttributeMap, "op", schematemplateanpAttributes.Ops)
	A(schematemplateanpAttributeMap, "path", schematemplateanpAttributes.Path)
	if schematemplateanpAttributes.Value != nil {
		A(schematemplateanpAttributeMap, "value", schematemplateanpAttributes.Value)
	}

	return schematemplateanpAttributeMap, nil
}
