package models

type SchemaSiteExternalEpg struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewSchemaSiteExternalEpg(ops, path string, siteEpgMap map[string]interface{}) *SchemaSiteExternalEpg {
	var externalepgMap map[string]interface{}
	externalepgMap = map[string]interface{}{
		"externalEpgRef": siteEpgMap["externalEpgRef"],
		"l3outDn":        siteEpgMap["l3outDn"],
		"l3outRef":       siteEpgMap["l3outRef"],
	}

	return &SchemaSiteExternalEpg{
		Ops:   ops,
		Path:  path,
		Value: externalepgMap,
	}
}

func (schemaSiteExternalEpgAttributes *SchemaSiteExternalEpg) ToMap() (map[string]interface{}, error) {
	schemaSiteExternalEpgAttributesMap := make(map[string]interface{})

	A(schemaSiteExternalEpgAttributesMap, "op", schemaSiteExternalEpgAttributes.Ops)
	A(schemaSiteExternalEpgAttributesMap, "path", schemaSiteExternalEpgAttributes.Path)
	if schemaSiteExternalEpgAttributes.Value != nil {
		A(schemaSiteExternalEpgAttributesMap, "value", schemaSiteExternalEpgAttributes.Value)
	}

	return schemaSiteExternalEpgAttributesMap, nil
}
