package models

type SchemaSiteAnpEpgSelector struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewSchemaSiteAnpEpgSelector(ops, path string, SiteAnpEpgSelectorMap map[string]interface{}) *SchemaSiteAnpEpgSelector {
	var temp map[string]interface{}

	if ops != "remove" {
		temp = SiteAnpEpgSelectorMap
	} else {
		temp = nil
	}

	return &SchemaSiteAnpEpgSelector{
		Ops:   ops,
		Path:  path,
		Value: temp,
	}
}

func (schemasiteanpepgselectorattr *SchemaSiteAnpEpgSelector) ToMap() (map[string]interface{}, error) {
	schemasiteanpepgselectorMap := make(map[string]interface{})

	A(schemasiteanpepgselectorMap, "op", schemasiteanpepgselectorattr.Ops)
	A(schemasiteanpepgselectorMap, "path", schemasiteanpepgselectorattr.Path)
	if schemasiteanpepgselectorattr.Value != nil {
		A(schemasiteanpepgselectorMap, "value", schemasiteanpepgselectorattr.Value)
	}

	return schemasiteanpepgselectorMap, nil
}
