package models

type Schema struct {
	Id          string `json:",omitempty"`
	DisplayName string `json:",omitempty"`

	Templates []map[string]interface{} `json:",omitempty"`

	Sites []map[string]interface{} `json:",omitempty"`
}

func NewSchema(id, displayName, templateName, tenantId string) *Schema {
	templateMap := map[string]interface{}{
		"name":          templateName,
		"tenantId":      tenantId,
		"displayName":   templateName,
		"anps":          []interface{}{},
		"contracts":     []interface{}{},
		"vrfs":          []interface{}{},
		"bds":           []interface{}{},
		"filters":       []interface{}{},
		"externalEpgs":  []interface{}{},
		"serviceGraphs": []interface{}{},
	}
	templates := []map[string]interface{}{
		templateMap,
	}
	return &Schema{
		Id:          id,
		DisplayName: displayName,
		Templates:   templates,
		Sites:       []map[string]interface{}{},
	}
}

func (schema *Schema) ToMap() (map[string]interface{}, error) {
	schemaAttributeMap := make(map[string]interface{})
	A(schemaAttributeMap, "id", schema.Id)
	A(schemaAttributeMap, "displayName", schema.DisplayName)
	A(schemaAttributeMap, "templates", schema.Templates)
	A(schemaAttributeMap, "sites", schema.Sites)

	return schemaAttributeMap, nil
}
