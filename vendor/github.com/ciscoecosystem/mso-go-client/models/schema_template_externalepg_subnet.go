package models

type ExternalEpgSubnet struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewTemplateExternalEpgSubnet(ops, path, ip, name string, scope, aggregate []interface{}) *ExternalEpgSubnet {
	var bdsubnetMap map[string]interface{}
	bdsubnetMap = map[string]interface{}{
		"ip":        ip,
		"name":      name,
		"scope":     scope,
		"aggregate": aggregate,
	}

	return &ExternalEpgSubnet{
		Ops:   ops,
		Path:  path,
		Value: bdsubnetMap,
	}

}

func (subnetAttribute *ExternalEpgSubnet) ToMap() (map[string]interface{}, error) {
	subnetAttributeMap := make(map[string]interface{})
	A(subnetAttributeMap, "op", subnetAttribute.Ops)
	A(subnetAttributeMap, "path", subnetAttribute.Path)
	if subnetAttribute.Value != nil {
		A(subnetAttributeMap, "value", subnetAttribute.Value)
	}

	return subnetAttributeMap, nil
}
