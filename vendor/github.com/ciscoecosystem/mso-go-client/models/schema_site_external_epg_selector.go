package models

type SchemaSiteExternalEpgSelector struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewSchemaSiteExternalEpgSelector(ops, path string, selectorMap map[string]interface{}) *SchemaSiteExternalEpgSelector {
	var temp map[string]interface{}

	if ops != "remove" {
		temp = selectorMap
	} else {
		temp = nil
	}

	return &SchemaSiteExternalEpgSelector{
		Ops:   ops,
		Path:  path,
		Value: temp,
	}
}

func (schemaSiteExternalEpgSelector *SchemaSiteExternalEpgSelector) ToMap() (map[string]interface{}, error) {
	schemaSiteExternalEpgSelectorMap := make(map[string]interface{})

	A(schemaSiteExternalEpgSelectorMap, "op", schemaSiteExternalEpgSelector.Ops)
	A(schemaSiteExternalEpgSelectorMap, "path", schemaSiteExternalEpgSelector.Path)
	if schemaSiteExternalEpgSelector.Value != nil {
		A(schemaSiteExternalEpgSelectorMap, "value", schemaSiteExternalEpgSelector.Value)
	}

	return schemaSiteExternalEpgSelectorMap, nil
}
