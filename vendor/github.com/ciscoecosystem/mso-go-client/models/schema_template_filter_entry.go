package models

type TemplateFilterEntry struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewTemplateFilterEntry(ops, path, entryName, entryDisplayName, entryDescription, etherType, arpFlag, ipProtocol, sourceFrom, sourceTo, destinationFrom, destinationTo string, matchOnlyFragments, stateful bool, tcpSessionRules []interface{}) *TemplateFilterEntry {
	var anpepgMap map[string]interface{}

	anpepgMap = map[string]interface{}{

		"name":               entryName,
		"displayName":        entryDisplayName,
		"description":        entryDescription,
		"etherType":          etherType,
		"arpFlag":            arpFlag,
		"ipProtocol":         ipProtocol,
		"matchOnlyFragments": matchOnlyFragments,
		"stateful":           stateful,
		"sourceFrom":         sourceFrom,
		"sourceTo":           sourceTo,
		"destinationFrom":    destinationFrom,
		"destinationTo":      destinationTo,
		"tcpSessionRules":    tcpSessionRules,
	}

	if anpepgMap["etherType"] == "" {
		anpepgMap["etherType"] = "unspecified"
	}
	if anpepgMap["arpFlag"] == "" {
		anpepgMap["arpFlag"] = "unspecified"
	}
	if anpepgMap["ipProtocol"] == "" {
		anpepgMap["ipProtocol"] = "unspecified"
	}

	if anpepgMap["sourceFrom"] == "" {
		anpepgMap["sourceFrom"] = "unspecified"
	}
	if anpepgMap["sourceTo"] == "" {
		anpepgMap["sourceTo"] = "unspecified"
	}
	if anpepgMap["destinationTo"] == "" {
		anpepgMap["destinationTo"] = "unspecified"
	}
	if anpepgMap["destinationFrom"] == "" {
		anpepgMap["destinationFrom"] = "unspecified"
	}

	return &TemplateFilterEntry{
		Ops:   ops,
		Path:  path,
		Value: anpepgMap,
	}

}

func NewTemplateFilter(ops, path, filterName, filterDisplayName string, entries []interface{}) *TemplateFilterEntry {
	var anpepgMap map[string]interface{}
	anpepgMap = map[string]interface{}{

		"name":        filterName,
		"displayName": filterDisplayName,
		"entries":     entries,
	}

	return &TemplateFilterEntry{
		Ops:   ops,
		Path:  path,
		Value: anpepgMap,
	}
}

func (anpAttributes *TemplateFilterEntry) ToMap() (map[string]interface{}, error) {
	anpAttributesMap := make(map[string]interface{})
	A(anpAttributesMap, "op", anpAttributes.Ops)
	A(anpAttributesMap, "path", anpAttributes.Path)
	if anpAttributes.Value != nil {
		A(anpAttributesMap, "value", anpAttributes.Value)
	}

	return anpAttributesMap, nil
}
