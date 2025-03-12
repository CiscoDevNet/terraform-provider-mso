package models

type SchemaTemplateVrf struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewSchemaTemplateVrf(ops, path, Name, displayName, ipDataPlaneLearning, desc string, l3m, vzany, preferredGroup, siteAwarePolicyEnforcementMode bool, rendezvousPoints []interface{}) *SchemaSite {
	var VrfMap map[string]interface{}

	if ops != "remove" {
		VrfMap = map[string]interface{}{
			"displayName":                    displayName,
			"description":                    desc,
			"name":                           Name,
			"l3MCast":                        l3m,
			"vzAnyEnabled":                   vzany,
			"preferredGroup":                 preferredGroup,
			"siteAwarePolicyEnforcementMode": siteAwarePolicyEnforcementMode,
		}
		if ipDataPlaneLearning != "" {
			VrfMap["ipDataPlaneLearning"] = ipDataPlaneLearning
		}
		if rendezvousPoints != nil {
			VrfMap["rpConfigs"] = rendezvousPoints
		}
	} else {
		VrfMap = nil
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
