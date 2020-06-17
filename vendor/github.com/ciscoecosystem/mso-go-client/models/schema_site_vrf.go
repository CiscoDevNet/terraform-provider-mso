package models

type SiteVrf struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewSchemaSiteVrf(ops, path string, vrfRef map[string]interface{}) *SiteVrf {
	var externalepgMap map[string]interface{}
	externalepgMap = map[string]interface{}{
		"vrfRef":      vrfRef,
		"regions":     []interface{}{},
	}

	return &SiteVrf{
		Ops:   ops,
		Path:  path,
		Value: externalepgMap,
	}

}

func (externalepgAttributes *SiteVrf) ToMap() (map[string]interface{}, error) {
	externalepgAttributesMap := make(map[string]interface{})
	A(externalepgAttributesMap, "op", externalepgAttributes.Ops)
	A(externalepgAttributesMap, "path", externalepgAttributes.Path)
	if externalepgAttributes.Value != nil {
		A(externalepgAttributesMap, "value", externalepgAttributes.Value)
	}

	return externalepgAttributesMap, nil
}