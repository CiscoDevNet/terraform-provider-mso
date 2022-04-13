package models

type Schema struct {
	Id          string `json:",omitempty"`
	DisplayName string `json:",omitempty"`

	Templates []map[string]interface{} `json:",omitempty"`

	Sites []map[string]interface{} `json:",omitempty"`
}

type SchemaReplace struct {
	Ops   string                   `json:",omitempty"`
	Path  string                   `json:",omitempty"`
	Value []map[string]interface{} `json:",omitempty"`
}

func NewSchema(id, displayName, templateName, tenantId string, template []interface{}) (*Schema, *SchemaReplace) {
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
				"name":          map_template_values["name"],
				"tenantId":      map_template_values["tenantId"],
				"displayName":   map_template_values["displayName"],
				"anps":          []interface{}{},
				"contracts":     []interface{}{},
				"vrfs":          []interface{}{},
				"bds":           []interface{}{},
				"filters":       []interface{}{},
				"externalEpgs":  []interface{}{},
				"serviceGraphs": []interface{}{},
			}
			result = append(result, templateMap)
		}
	}
	if id == "" {
		return &Schema{
			Id:          id,
			DisplayName: displayName,
			Templates:   result,
			Sites:       []map[string]interface{}{},
		}, nil
	} else {
		return nil, &SchemaReplace{
			Ops:   "replace",
			Path:  "/templates",
			Value: result,
		}
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

func (schemaReplace *SchemaReplace) ToMap() (map[string]interface{}, error) {
	schemaReplaceMap := make(map[string]interface{})
	A(schemaReplaceMap, "op", schemaReplace.Ops)
	A(schemaReplaceMap, "path", schemaReplace.Path)
	if schemaReplace.Value != nil {
		A(schemaReplaceMap, "value", schemaReplace.Value)
	}

	return schemaReplaceMap, nil
}
