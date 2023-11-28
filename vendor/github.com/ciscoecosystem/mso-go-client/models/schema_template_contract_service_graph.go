package models

type TemplateContractServiceGraph struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

type SiteContractServiceGraph struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewTemplateContractServiceGraph(ops, path string, serviceGraph map[string]interface{}, nodeRelation []interface{}) *TemplateExternalepg {
	var serviceGraphMap map[string]interface{}
	if ops != "remove" {
		serviceGraphMap = map[string]interface{}{
			"serviceGraphRef":          serviceGraph,
			"serviceNodesRelationship": nodeRelation,
		}
	}

	return &TemplateExternalepg{
		Ops:   ops,
		Path:  path,
		Value: serviceGraphMap,
	}
}

func NewSiteContractServiceGraph(ops, path string, serviceGraph map[string]interface{}, nodeRelation []interface{}) *TemplateExternalepg {
	var serviceGraphMap map[string]interface{}
	if ops != "remove" {
		serviceGraphMap = map[string]interface{}{
			"serviceGraphRef":          serviceGraph,
			"serviceNodesRelationship": nodeRelation,
		}
	}

	return &TemplateExternalepg{
		Ops:   ops,
		Path:  path,
		Value: serviceGraphMap,
	}
}

func (graphAttr *TemplateContractServiceGraph) ToMap() (map[string]interface{}, error) {
	graphAttrMap := make(map[string]interface{})
	A(graphAttrMap, "op", graphAttr.Ops)
	A(graphAttrMap, "path", graphAttr.Path)
	if graphAttr.Value != nil {
		A(graphAttrMap, "value", graphAttr.Value)
	}

	return graphAttrMap, nil
}

func (graphAttr *SiteContractServiceGraph) ToMap() (map[string]interface{}, error) {
	graphAttrMap := make(map[string]interface{})
	A(graphAttrMap, "op", graphAttr.Ops)
	A(graphAttrMap, "path", graphAttr.Path)
	if graphAttr.Value != nil {
		A(graphAttrMap, "value", graphAttr.Value)
	}

	return graphAttrMap, nil
}

func NewSiteContractServiceGraphListener(ops, path, name, protocol, securityPolicy string, port int, certificates []interface{}, rules []interface{}, frontendIpDnMap map[string]string) *PatchPayload {

	listenerMap := map[string]interface{}{
		"name":         name,
		"port":         port,
		"protocol":     protocol,
		"secPolicy":    securityPolicy,
		"certificates": certificates,
		"rules":        rules,
		"nlbDevIp":     frontendIpDnMap,
	}

	return &PatchPayload{
		Ops:   ops,
		Path:  path,
		Value: listenerMap,
	}
}
