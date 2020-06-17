package models

type SchemaSiteVrfRegionCidr struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewSchemaSiteVrfRegionCidr(ops, path, ip string, primary bool) *SchemaSiteVrfRegionCidr {
	var siteVrfRegionCidrMap map[string]interface{}
	siteVrfRegionCidrMap = map[string]interface{}{
		"ip":      ip,
		"primary": primary,
	}

	return &SchemaSiteVrfRegionCidr{
		Ops:   ops,
		Path:  path,
		Value: siteVrfRegionCidrMap,
	}

}

func (siteVrfRegionCidrAttributes *SchemaSiteVrfRegionCidr) ToMap() (map[string]interface{}, error) {
	siteVrfRegionCidrAttributesMap := make(map[string]interface{})
	A(siteVrfRegionCidrAttributesMap, "op", siteVrfRegionCidrAttributes.Ops)
	A(siteVrfRegionCidrAttributesMap, "path", siteVrfRegionCidrAttributes.Path)
	if siteVrfRegionCidrAttributes.Value != nil {
		A(siteVrfRegionCidrAttributesMap, "value", siteVrfRegionCidrAttributes.Value)
	}

	return siteVrfRegionCidrAttributesMap, nil
}
