package models

type SiteBd struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewSchemaSiteBd(ops, path, mac string, bdRef map[string]interface{}, host bool) *SiteBd {
	siteBdMap := map[string]interface{}{
		"bdRef":            bdRef,
		"hostBasedRouting": host,
	}

	if mac != "" {
		siteBdMap["mac"] = mac
	}

	return &SiteBd{
		Ops:   ops,
		Path:  path,
		Value: siteBdMap,
	}

}

func (externalepgAttributes *SiteBd) ToMap() (map[string]interface{}, error) {
	externalepgAttributesMap := make(map[string]interface{})
	A(externalepgAttributesMap, "op", externalepgAttributes.Ops)
	A(externalepgAttributesMap, "path", externalepgAttributes.Path)
	if externalepgAttributes.Value != nil {
		A(externalepgAttributesMap, "value", externalepgAttributes.Value)
	}

	return externalepgAttributesMap, nil
}
