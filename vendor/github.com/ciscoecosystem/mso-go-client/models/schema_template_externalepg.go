package models

type TemplateExternalepg struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewTemplateExternalepg(ops, path, name, displayName, externalEpgType, desc string, preferredGroup bool, vrfRef map[string]interface{}, l3outRef map[string]interface{}, anpRef map[string]interface{}, selector []interface{}) *TemplateExternalepg {
	var externalepgMap map[string]interface{}
	externalepgMap = map[string]interface{}{
		"name":           name,
		"displayName":    displayName,
		"description":    desc,
		"vrfRef":         vrfRef,
		"extEpgType":     externalEpgType,
		"preferredGroup": preferredGroup,
	}

	if l3outRef != nil {
		externalepgMap["l3outRef"] = l3outRef
	}

	if anpRef != nil {
		externalepgMap["anpRef"] = anpRef
	}

	if selector != nil {
		externalepgMap["selectors"] = selector
	}

	return &TemplateExternalepg{
		Ops:   ops,
		Path:  path,
		Value: externalepgMap,
	}

}

func (externalepgAttributes *TemplateExternalepg) ToMap() (map[string]interface{}, error) {
	externalepgAttributesMap := make(map[string]interface{})
	A(externalepgAttributesMap, "op", externalepgAttributes.Ops)
	A(externalepgAttributesMap, "path", externalepgAttributes.Path)
	if externalepgAttributes.Value != nil {
		A(externalepgAttributesMap, "value", externalepgAttributes.Value)
	}

	return externalepgAttributesMap, nil
}
