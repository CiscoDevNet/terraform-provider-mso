package models

type ServiceNodeTypeAttributes struct {
	Name        string `json:",omitempty"`
	DisplayName string `json:",omitempty"`
}

func NewServiceNodeType(typeAttr ServiceNodeTypeAttributes) *ServiceNodeTypeAttributes {

	TypeAttributes := typeAttr
	return &TypeAttributes
}

func (typeAttributes *ServiceNodeTypeAttributes) ToMap() (map[string]interface{}, error) {
	typeAttributeMap := make(map[string]interface{})
	A(typeAttributeMap, "name", typeAttributes.Name)
	A(typeAttributeMap, "displayName", typeAttributes.DisplayName)

	return typeAttributeMap, nil
}
