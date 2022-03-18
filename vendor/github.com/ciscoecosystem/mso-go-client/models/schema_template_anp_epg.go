package models

type TemplateAnpEpg struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewTemplateAnpEpg(ops, path, name, displayName, intraEpg, epgType string, uSegEpg, intersiteMulticasteSource, preferredGroup, proxyArp bool, vrfRef, bdRef, cloudServiceEpgConfig map[string]interface{}) *TemplateAnpEpg {
	var anpepgMap map[string]interface{}
	anpepgMap = map[string]interface{}{
		"name":                  name,
		"displayName":           displayName,
		"subnets":               []interface{}{},
		"uSegEpg":               uSegEpg,
		"intraEpg":              intraEpg,
		"epgType":               epgType,
		"mCastSource":           intersiteMulticasteSource,
		"proxyArp":              proxyArp,
		"preferredGroup":        preferredGroup,
		"cloudServiceEpgConfig": cloudServiceEpgConfig,
	}

	if _, ok := vrfRef["vrfName"]; ok {
		anpepgMap["vrfRef"] = vrfRef
	}

	if _, ok := bdRef["bdName"]; ok {
		anpepgMap["bdRef"] = bdRef
	}

	if anpepgMap["intraEpg"] == "" {
		anpepgMap["intraEpg"] = "unenforced"
	}

	return &TemplateAnpEpg{
		Ops:   ops,
		Path:  path,
		Value: anpepgMap,
	}

}

func (anpAttributes *TemplateAnpEpg) ToMap() (map[string]interface{}, error) {
	anpAttributesMap := make(map[string]interface{})
	A(anpAttributesMap, "op", anpAttributes.Ops)
	A(anpAttributesMap, "path", anpAttributes.Path)
	if anpAttributes.Value != nil {
		A(anpAttributesMap, "value", anpAttributes.Value)
	}

	return anpAttributesMap, nil
}
