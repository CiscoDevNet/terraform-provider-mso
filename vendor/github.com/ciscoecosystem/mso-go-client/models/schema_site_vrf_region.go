package models

type SchemaSiteVrfRegion struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewSchemaSiteVrfRegion(ops, path, name string) *SchemaSiteVrfRegion {
	var siteVrfRegionMap map[string]interface{}
	siteVrfRegionMap = map[string]interface{}{
		"name": name,
	}

	return &SchemaSiteVrfRegion{
		Ops:   ops,
		Path:  path,
		Value: siteVrfRegionMap,
	}

}

func (siteVrfRegionAttributes *SchemaSiteVrfRegion) ToMap() (map[string]interface{}, error) {
	siteVrfRegionAttributesMap := make(map[string]interface{})
	A(siteVrfRegionAttributesMap, "op", siteVrfRegionAttributes.Ops)
	A(siteVrfRegionAttributesMap, "path", siteVrfRegionAttributes.Path)
	if siteVrfRegionAttributes.Value != nil {
		A(siteVrfRegionAttributesMap, "value", siteVrfRegionAttributes.Value)
	}

	return siteVrfRegionAttributesMap, nil
}
