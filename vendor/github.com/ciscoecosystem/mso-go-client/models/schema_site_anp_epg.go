package models

type SchemaSiteAnpEpg struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewSchemaSiteAnpEpg(ops, path string, privateLinkLabel, epgRef map[string]interface{}) *SchemaSiteAnpEpg {
	var siteAnpEpgMap map[string]interface{}
	siteAnpEpgMap = map[string]interface{}{
		"epgRef":             epgRef,
		"domainAssociations": []interface{}{},
		"staticPorts":        []interface{}{},
		"contracts":          []interface{}{},
		"staticLeafs":        []interface{}{},
		"uSegAttrs":          []interface{}{},
		"subnets":            []interface{}{},
		"selectors":          []interface{}{},
		"privateLinkLabel":   privateLinkLabel,
	}

	return &SchemaSiteAnpEpg{
		Ops:   ops,
		Path:  path,
		Value: siteAnpEpgMap,
	}

}

func (siteAnpEpgAttributes *SchemaSiteAnpEpg) ToMap() (map[string]interface{}, error) {
	siteAnpEpgAttributesMap := make(map[string]interface{})
	A(siteAnpEpgAttributesMap, "op", siteAnpEpgAttributes.Ops)
	A(siteAnpEpgAttributesMap, "path", siteAnpEpgAttributes.Path)
	if siteAnpEpgAttributes.Value != nil {
		A(siteAnpEpgAttributesMap, "value", siteAnpEpgAttributes.Value)
	}

	return siteAnpEpgAttributesMap, nil
}
