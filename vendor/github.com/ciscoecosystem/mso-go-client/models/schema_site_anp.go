package models

type SiteAnp struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewSchemaSiteAnp(ops, path string, anpRef map[string]interface{}) *SiteAnp {
	var externalepgMap map[string]interface{}
	externalepgMap = map[string]interface{}{
		"anpRef": anpRef,
		"epgs":   []interface{}{},
	}

	return &SiteAnp{
		Ops:   ops,
		Path:  path,
		Value: externalepgMap,
	}

}

func (externalepgAttributes *SiteAnp) ToMap() (map[string]interface{}, error) {
	externalepgAttributesMap := make(map[string]interface{})
	A(externalepgAttributesMap, "op", externalepgAttributes.Ops)
	A(externalepgAttributesMap, "path", externalepgAttributes.Path)
	if externalepgAttributes.Value != nil {
		A(externalepgAttributesMap, "value", externalepgAttributes.Value)
	}

	return externalepgAttributesMap, nil
}
