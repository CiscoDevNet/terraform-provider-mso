package models

type TenantAttributes struct {
	Name        string        `json:",omitempty"`
	DisplayName string        `json:",omitempty"`
	Description string        `json:",omitempty"`
	Users       []interface{} `json:",omitempty"`
	Sites       []interface{} `json:",omitempty"`
}

func NewTenant(tenantAttr TenantAttributes) *TenantAttributes {

	TenantAttributes := tenantAttr
	return &TenantAttributes
}

func (tenantAttributes *TenantAttributes) ToMap() (map[string]interface{}, error) {
	tenantAttributeMap := make(map[string]interface{})
	A(tenantAttributeMap, "name", tenantAttributes.Name)
	A(tenantAttributeMap, "displayName", tenantAttributes.DisplayName)
	A(tenantAttributeMap, "description", tenantAttributes.Description)
	A(tenantAttributeMap, "userAssociations", tenantAttributes.Users)
	A(tenantAttributeMap, "siteAssociations", tenantAttributes.Sites)

	return tenantAttributeMap, nil
}
