package models

type Label struct {
	Id           string  `json:",omitempty"`
	DisplayName string `json:",omitempty"`
	Type      string `json:",omitempty"`
}




func NewLabel (id, labels, types string) *Label{
	
	return &Label{
		Id:  id,
		DisplayName: labels,
		Type:  types,
		

	}
}

func (label *Label) ToMap() (map[string]interface{}, error) {
	labelAttributeMap := make(map[string]interface{})
	A(labelAttributeMap, "id", label.Id)
	A(labelAttributeMap, "displayName", label.DisplayName)
	A(labelAttributeMap, "type", label.Type)



	return labelAttributeMap, nil
}