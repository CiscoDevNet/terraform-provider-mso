package models

type SchemaSiteVrfRegionCidrSubnet struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewSchemaSiteVrfRegionCidrSubnet(ops, path, name, ip, zone, usage, subnetGroup string) *SchemaSiteVrfRegionCidrSubnet {
	var bdsubnetMap map[string]interface{}
	if ops != "remove" {
		bdsubnetMap = map[string]interface{}{
			"ip": ip,
		}
		if name != "" {
			bdsubnetMap["name"] = name
		}
		if zone != "" {
			bdsubnetMap["zone"] = zone
		}
		if usage != "" {
			bdsubnetMap["usage"] = usage
		}
		if subnetGroup != "" {
			bdsubnetMap["subnetGroup"] = subnetGroup
		}
	} else {
		bdsubnetMap = nil
	}

	return &SchemaSiteVrfRegionCidrSubnet{
		Ops:   ops,
		Path:  path,
		Value: bdsubnetMap,
	}

}

func (bdAttributes *SchemaSiteVrfRegionCidrSubnet) ToMap() (map[string]interface{}, error) {
	bdAttributesMap := make(map[string]interface{})
	A(bdAttributesMap, "op", bdAttributes.Ops)
	A(bdAttributesMap, "path", bdAttributes.Path)
	if bdAttributes.Value != nil {
		A(bdAttributesMap, "value", bdAttributes.Value)
	}

	return bdAttributesMap, nil
}
