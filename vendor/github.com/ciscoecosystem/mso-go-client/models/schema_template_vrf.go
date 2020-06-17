package models

type SchemaTemplateVrf struct {
	Ops         string `json:",omitempty"`
	Path         string `json:",omitempty"`
	Value       map[string]interface{} `json:",omitempty"`
	
}

func NewSchemaTemplateVrf(ops, path, Name,displayName string,l3m bool) *SchemaSite {
	var VrfMap map[string]interface{}
	
	if ops !="remove" {
		VrfMap = map[string]interface{}{
			"displayName":     displayName,
			"name":            Name,
			"l3MCast":         l3m,
		}
	}else{
		
		VrfMap=nil
	}

	return &SchemaSite{
		Ops:   ops,
		Path:  path,
		Value: VrfMap,
	}

}

func (schematemplatevrfAttributes *SchemaTemplateVrf) ToMap() (map[string]interface{}, error) {
	schematemplatevrfAttributeMap := make(map[string]interface{})
	A(schematemplatevrfAttributeMap, "op", schematemplatevrfAttributes.Ops)
	A(schematemplatevrfAttributeMap, "path", schematemplatevrfAttributes.Path)
	if schematemplatevrfAttributes.Value != nil {
		A(schematemplatevrfAttributeMap, "value", schematemplatevrfAttributes.Value)
	}

	return schematemplatevrfAttributeMap, nil
}