package models

type TemplateServiceGraph struct {
	Ops   string `json:",omitempty"`
	Path  string `json:",omitempty"`
	Value map[string]interface{}
}

type TemplateServiceGraphUpdate struct {
	Ops   string `json:",omitempty"`
	Path  string `json:",omitempty"`
	Value interface{}
}

func NewTemplateServiceGraphUpdate(ops, path string, graphRef interface{}) *TemplateServiceGraphUpdate {
	return &TemplateServiceGraphUpdate{
		Ops:   ops,
		Path:  path,
		Value: graphRef,
	}
}
func NewTemplateServiceGraph(ops, path string, graphRef map[string]interface{}) *TemplateServiceGraph {

	return &TemplateServiceGraph{
		Ops:   ops,
		Path:  path,
		Value: graphRef,
	}

}

func (graphAttributes *TemplateServiceGraphUpdate) ToMap() (map[string]interface{}, error) {
	graphAttributesMap := make(map[string]interface{})
	A(graphAttributesMap, "op", graphAttributes.Ops)
	A(graphAttributesMap, "path", graphAttributes.Path)
	if graphAttributes.Value != nil {
		A(graphAttributesMap, "value", graphAttributes.Value)
	}

	return graphAttributesMap, nil
}

func (graphAttributes *TemplateServiceGraph) ToMap() (map[string]interface{}, error) {
	graphAttributesMap := make(map[string]interface{})
	A(graphAttributesMap, "op", graphAttributes.Ops)
	A(graphAttributesMap, "path", graphAttributes.Path)
	if graphAttributes.Value != nil {
		A(graphAttributesMap, "value", graphAttributes.Value)
	}

	return graphAttributesMap, nil
}
