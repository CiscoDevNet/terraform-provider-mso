package models

type TemplateAnpEpgContract struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewTemplateAnpEpgContract(ops, path string, contractRef map[string]interface{}, relationshipType string) *TemplateAnpEpgContract {
	var epgcontractMap map[string]interface{}
	if ops != "remove" {
		epgcontractMap = map[string]interface{}{
			"contractRef":      contractRef,
			"relationshipType": relationshipType,
		}
	} else {
		epgcontractMap = nil
	}

	return &TemplateAnpEpgContract{
		Ops:   ops,
		Path:  path,
		Value: epgcontractMap,
	}

}

func (bdAttributes *TemplateAnpEpgContract) ToMap() (map[string]interface{}, error) {
	bdAttributesMap := make(map[string]interface{})
	A(bdAttributesMap, "op", bdAttributes.Ops)
	A(bdAttributesMap, "path", bdAttributes.Path)
	if bdAttributes.Value != nil {
		A(bdAttributesMap, "value", bdAttributes.Value)
	}

	return bdAttributesMap, nil
}
