package models

type TemplateL3out struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewTemplateL3out(ops, path, name, displayName, desc string, vrfRef map[string]interface{}) *TemplateL3out {
	var l3outMap map[string]interface{}
	l3outMap = map[string]interface{}{
		"name":        name,
		"displayName": displayName,
		"description": desc,
		"vrfRef":      vrfRef,
	}

	return &TemplateL3out{
		Ops:   ops,
		Path:  path,
		Value: l3outMap,
	}

}

func (l3outAttributes *TemplateL3out) ToMap() (map[string]interface{}, error) {
	l3outAttributesMap := make(map[string]interface{})
	A(l3outAttributesMap, "op", l3outAttributes.Ops)
	A(l3outAttributesMap, "path", l3outAttributes.Path)
	if l3outAttributes.Value != nil {
		A(l3outAttributesMap, "value", l3outAttributes.Value)
	}

	return l3outAttributesMap, nil
}
