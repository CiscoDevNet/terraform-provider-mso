package models

type SchemaSiteAnpEpgDomain struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewSchemaSiteAnpEpgDomain(ops, path, domainType, dn, deploymentImmediacy, resolutionImmediacy string, vmmDomainProperties map[string]interface{}) *SchemaSiteAnpEpgDomain {
	siteAnpEpgDomainMap := map[string]interface{}{
		"domainType":          domainType,
		"dn":                  dn,
		"deploymentImmediacy": deploymentImmediacy, // keeping for backworths compatibility
		"deployImmediacy":     deploymentImmediacy, // rename of deploymentImmediacy
		"resolutionImmediacy": resolutionImmediacy,
		"vmmDomainProperties": vmmDomainProperties,
	}

	if len(vmmDomainProperties) > 0 {
		injectVmmDomainProperties(siteAnpEpgDomainMap, vmmDomainProperties)
	}

	return &SchemaSiteAnpEpgDomain{
		Ops:   ops,
		Path:  path,
		Value: siteAnpEpgDomainMap,
	}

}

func injectVmmDomainProperties(siteAnpEpgDomainMap, vmmDomainProperties map[string]interface{}) {

	properties := []string{
		"allowMicroSegmentation",
		"epgLagPol",
		"switchType",
		"switchingMode",
		"vlanEncapMode",
		"portEncapVlan",
		"microSegVlan",
		"delimiter",
		"bindingType",
		"numPorts",
		"portAllocation",
		"netflowPref",
		"allowPromiscuous",
		"forgedTransmits",
		"macChanges",
		"customEpgName",
	}
	for _, property := range properties {
		value, exists := vmmDomainProperties[property]
		if exists {
			siteAnpEpgDomainMap[property] = value
		}
	}
}

func (siteAnpEpgDomainAttributes *SchemaSiteAnpEpgDomain) ToMap() (map[string]interface{}, error) {
	siteAnpEpgDomainAttributesMap := make(map[string]interface{})
	A(siteAnpEpgDomainAttributesMap, "op", siteAnpEpgDomainAttributes.Ops)
	A(siteAnpEpgDomainAttributesMap, "path", siteAnpEpgDomainAttributes.Path)
	if siteAnpEpgDomainAttributes.Value != nil {
		A(siteAnpEpgDomainAttributesMap, "value", siteAnpEpgDomainAttributes.Value)
	}

	return siteAnpEpgDomainAttributesMap, nil
}
