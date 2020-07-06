package models

type TemplateVRFContract struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewTemplateVRFContract(ops, path string, contractRef map[string]interface{}) *TemplateVRFContract {

	return &TemplateVRFContract{
		Ops:   ops,
		Path:  path,
		Value: contractRef,
	}

}

func (vrfConAttributes *TemplateVRFContract) ToMap() (map[string]interface{}, error) {
	vrfConAttributesMap := make(map[string]interface{})
	A(vrfConAttributesMap, "op", vrfConAttributes.Ops)
	A(vrfConAttributesMap, "path", vrfConAttributes.Path)
	if vrfConAttributes.Value != nil {
		A(vrfConAttributesMap, "value", vrfConAttributes.Value)
	}

	return vrfConAttributesMap, nil
}
