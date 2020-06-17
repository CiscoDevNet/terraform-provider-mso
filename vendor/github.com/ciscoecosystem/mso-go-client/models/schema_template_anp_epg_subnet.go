package models

type SchemaTemplateAnpEpgSubnet struct {
	Ops         string `json:",omitempty"`
	Path         string `json:",omitempty"`
	Value       map[string]interface{} `json:",omitempty"`
	
}

func NewSchemaTemplateAnpEpgSubnet(ops, path, ip,scope string,shared bool) *SchemaTemplateAnpEpgSubnet{
	var SubnetMap map[string]interface{}
	
	if ops !="remove" {
		SubnetMap = map[string]interface{}{
			"ip":     ip,
			"scope":    scope,
			"shared":   shared,
		}
	}else{
		
		SubnetMap=nil
	}

	return &SchemaTemplateAnpEpgSubnet{
		Ops:   ops,
		Path:  path,
		Value: SubnetMap,
	}

}

func (schematemplateanpepgsubnetAttributes *SchemaTemplateAnpEpgSubnet) ToMap() (map[string]interface{}, error) {
	schematemplateanpepgsubnetAttributeMap := make(map[string]interface{})
	A(schematemplateanpepgsubnetAttributeMap, "op", schematemplateanpepgsubnetAttributes.Ops)
	A(schematemplateanpepgsubnetAttributeMap, "path", schematemplateanpepgsubnetAttributes.Path)
	if schematemplateanpepgsubnetAttributes.Value != nil {
		A(schematemplateanpepgsubnetAttributeMap, "value", schematemplateanpepgsubnetAttributes.Value)
	}

	return schematemplateanpepgsubnetAttributeMap, nil
}