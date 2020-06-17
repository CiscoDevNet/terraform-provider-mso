package models

type SchemaSiteAnpEpgSubnet struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewSchemaSiteAnpEpgSubnet(ops, path, ip, desc, scope string, shared, noDefaultGateway, querier bool) *SchemaSiteAnpEpgSubnet {
	var bdsubnetMap map[string]interface{}
	if ops != "remove" {
		bdsubnetMap = map[string]interface{}{
			"ip":               ip,
			"description":      desc,
			"scope":            scope,
			"shared":           shared,
			"noDefaultGateway": noDefaultGateway,
			"querier":          querier,
		}
	} else {
		bdsubnetMap = nil
	}

	return &SchemaSiteAnpEpgSubnet{
		Ops:   ops,
		Path:  path,
		Value: bdsubnetMap,
	}

}

func (bdAttributes *SchemaSiteAnpEpgSubnet) ToMap() (map[string]interface{}, error) {
	bdAttributesMap := make(map[string]interface{})
	A(bdAttributesMap, "op", bdAttributes.Ops)
	A(bdAttributesMap, "path", bdAttributes.Path)
	if bdAttributes.Value != nil {
		A(bdAttributesMap, "value", bdAttributes.Value)
	}

	return bdAttributesMap, nil
}
