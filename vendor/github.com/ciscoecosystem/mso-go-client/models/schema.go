package models

type Schema struct {
	Id          string                   `json:",omitempty"`
	DisplayName string                   `json:",omitempty"`
	Description string                   `json:",omitempty"`
	Templates   []map[string]interface{} `json:",omitempty"`

	Sites []map[string]interface{} `json:",omitempty"`
}

func NewSchema(id, displayName, description, templateName, tenantId string, template []interface{}) *Schema {
	result := []map[string]interface{}{}
	if templateName != "" {
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
		result = []map[string]interface{}{
			templateMap,
		}
	} else {
		for _, map_values := range template {
			map_template_values := map_values.(map[string]interface{})
			templateMap := map[string]interface{}{
				"name":            map_template_values["name"],
				"tenantId":        map_template_values["tenantId"],
				"displayName":     map_template_values["displayName"],
				"description":     map_template_values["description"],
				"templateType":    map_template_values["templateType"],
				"templateSubType": map_template_values["templateSubType"],
			}
			result = append(result, templateMap)
		}
	}

	return &Schema{
		Id:          id,
		Description: description,
		DisplayName: displayName,
		Templates:   result,
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
