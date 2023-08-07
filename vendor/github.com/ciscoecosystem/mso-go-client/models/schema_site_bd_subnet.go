package models

type SchemaSiteBdSubnet struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewSchemaSiteBdSubnet(ops, path, ip, desc, scope string, shared, noDefaultGateway, querier, primary, virtual bool) *SchemaSiteBdSubnet {
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

	return &SchemaSiteBdSubnet{
		Ops:   ops,
		Path:  path,
		Value: bdsubnetMap,
	}

}

func (bdAttributes *SchemaSiteBdSubnet) ToMap() (map[string]interface{}, error) {
	bdAttributesMap := make(map[string]interface{})
	A(bdAttributesMap, "op", bdAttributes.Ops)
	A(bdAttributesMap, "path", bdAttributes.Path)
	if bdAttributes.Value != nil {
		A(bdAttributesMap, "value", bdAttributes.Value)
	}

	return bdAttributesMap, nil
}
