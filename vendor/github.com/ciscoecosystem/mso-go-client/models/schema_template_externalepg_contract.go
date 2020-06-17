package models

type ExternalEpgContract struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewTemplateExternalEpgContract(ops, path, relationshipType string, contractRefMap map[string]interface{}) *ExternalEpgContract {
	var contractMap map[string]interface{}
	if(ops != "remove"){
	contractMap = map[string]interface{}{
		"relationshipType":                      relationshipType,
		"contractRef":                           contractRefMap,
	}
	} else{
		contractMap = nil;
	}


	return &ExternalEpgContract{
		Ops:   ops,
		Path:  path,
		Value: contractMap,
	}

}

func (anpAttributes *ExternalEpgContract) ToMap() (map[string]interface{}, error) {
	anpAttributesMap := make(map[string]interface{})
	A(anpAttributesMap, "op", anpAttributes.Ops)
	A(anpAttributesMap, "path", anpAttributes.Path)
	if anpAttributes.Value != nil {
		A(anpAttributesMap, "value", anpAttributes.Value)
	}

	return anpAttributesMap, nil
}
