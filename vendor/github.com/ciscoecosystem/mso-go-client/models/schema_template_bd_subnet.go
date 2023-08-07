package models

type TemplateBDSubnet struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewTemplateBDSubnet(ops, path, ip, desc, scope string, shared, noDefaultGateway, querier, primary, virtual bool) *TemplateBDSubnet {
	var bdsubnetMap map[string]interface{}
	if ops != "remove" {
		bdsubnetMap = map[string]interface{}{
			"ip":               ip,
			"description":      desc,
			"scope":            scope,
			"shared":           shared,
			"noDefaultGateway": noDefaultGateway,
			"querier":          querier,
			"primary":          primary,
			"virtual":          virtual,
		}
	} else {
		bdsubnetMap = nil
	}

	return &TemplateBDSubnet{
		Ops:   ops,
		Path:  path,
		Value: bdsubnetMap,
	}

}

func (bdAttributes *TemplateBDSubnet) ToMap() (map[string]interface{}, error) {
	bdAttributesMap := make(map[string]interface{})
	A(bdAttributesMap, "op", bdAttributes.Ops)
	A(bdAttributesMap, "path", bdAttributes.Path)
	if bdAttributes.Value != nil {
		A(bdAttributesMap, "value", bdAttributes.Value)
	}

	return bdAttributesMap, nil
}
