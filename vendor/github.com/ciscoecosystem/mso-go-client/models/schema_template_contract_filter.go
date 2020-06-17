package models

type TemplateContractFilter struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value string `json:",omitempty"`
}

type TemplateContractFilterRelationShip struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewTemplateContractFilter(ops, path, filterType string) *TemplateContractFilter {
	
	return &TemplateContractFilter{
		Ops:   ops,
		Path:  path,
		Value: filterType,
	}

}

func NewTemplateContractFilterRelationShip(ops, path string, filterRef map[string]interface{}, directives []interface{}) *TemplateContractFilterRelationShip {
	var contractMap map[string]interface{}
	if ops != "remove" {
		contractMap = map[string]interface{}{
			"filterRef":                              filterRef,
			"directives":                             directives,
			
		}
	} else {
		contractMap = nil
	}
	return &TemplateContractFilterRelationShip{
		Ops:   ops,
		Path:  path,
		Value: contractMap,
	}

}

func (FilterAttributes *TemplateContractFilter) ToMap() (map[string]interface{}, error) {
	FilterAttributesMap := make(map[string]interface{})
	A(FilterAttributesMap, "op", FilterAttributes.Ops)
	A(FilterAttributesMap, "path", FilterAttributes.Path)
	if FilterAttributes.Value != "" {
		A(FilterAttributesMap, "value", FilterAttributes.Value)
	}

	return FilterAttributesMap, nil
}

func (FilterAttributes *TemplateContractFilterRelationShip) ToMap() (map[string]interface{}, error) {
	FilterAttributesMap := make(map[string]interface{})
	A(FilterAttributesMap, "op", FilterAttributes.Ops)
	A(FilterAttributesMap, "path", FilterAttributes.Path)
	if FilterAttributes.Value != nil {
		A(FilterAttributesMap, "value", FilterAttributes.Value)
	}

	return FilterAttributesMap, nil
}

