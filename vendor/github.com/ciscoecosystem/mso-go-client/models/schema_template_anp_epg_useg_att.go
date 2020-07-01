package models

type SchemaTemplateAnpEpgUsegAttr struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewSchemaTemplateAnpEpgUsegAttr(ops, path string, selectorMap map[string]interface{}) *SchemaTemplateAnpEpgUsegAttr {
	var temp map[string]interface{}

	if ops != "remove" {
		temp = selectorMap
	} else {
		temp = nil
	}

	return &SchemaTemplateAnpEpgUsegAttr{
		Ops:   ops,
		Path:  path,
		Value: temp,
	}
}

func (schematemplateanpepgusegattr *SchemaTemplateAnpEpgUsegAttr) ToMap() (map[string]interface{}, error) {
	schematemplateanpepgUsegAttrMap := make(map[string]interface{})

	A(schematemplateanpepgUsegAttrMap, "op", schematemplateanpepgusegattr.Ops)
	A(schematemplateanpepgUsegAttrMap, "path", schematemplateanpepgusegattr.Path)
	if schematemplateanpepgusegattr.Value != nil {
		A(schematemplateanpepgUsegAttrMap, "value", schematemplateanpepgusegattr.Value)
	}

	return schematemplateanpepgUsegAttrMap, nil
}
