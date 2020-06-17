package models

type SchemaSiteBdL3out struct {
	Ops   string `json:",omitempty"`
	Path  string `json:",omitempty"`
	Value string `json:",omitempty"`
}

func NewSchemaSiteBdL3out(ops, path, l3out string) *SchemaSiteBdL3out {

	return &SchemaSiteBdL3out{
		Ops:   ops,
		Path:  path,
		Value: l3out,
	}

}

func (siteBdL3outAttributes *SchemaSiteBdL3out) ToMap() (map[string]interface{}, error) {
	siteBdL3outAttributesMap := make(map[string]interface{})
	A(siteBdL3outAttributesMap, "op", siteBdL3outAttributes.Ops)
	A(siteBdL3outAttributesMap, "path", siteBdL3outAttributes.Path)
	A(siteBdL3outAttributesMap, "value", siteBdL3outAttributes.Value)

	return siteBdL3outAttributesMap, nil
}
