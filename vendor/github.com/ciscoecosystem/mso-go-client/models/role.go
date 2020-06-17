package models

type RoleAttributes struct {
	Id          string `json:",omitempty"`
	Name        string `json:",omitempty`
	DisplayName string `json:",omitempty"`
	Description string `json:",omitempty"`

	ReadPermissions []interface{} `json:",omitempty"`

	WritePermissions []interface{} `json:",omitempty"`
}

func NewRole(roleAttr RoleAttributes) *RoleAttributes {

	RoleAttributes := roleAttr
	return &RoleAttributes
}

func (role *RoleAttributes) ToMap() (map[string]interface{}, error) {
	roleAttributeMap := make(map[string]interface{})
	A(roleAttributeMap, "id", role.Id)
	A(roleAttributeMap, "name", role.Name)
	A(roleAttributeMap, "displayName", role.DisplayName)
	A(roleAttributeMap, "description", role.Description)
	A(roleAttributeMap, "readPermissions", role.ReadPermissions)
	A(roleAttributeMap, "writePermissions", role.WritePermissions)

	return roleAttributeMap, nil
}
