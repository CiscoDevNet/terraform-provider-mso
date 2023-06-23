package models

func NewSchemaSiteVrfRouteLeak(ops, path, tenantName, vrfRef string, includeAllSubnets bool, prefixSubnets []map[string]string, siteIds []string) *PatchPayload {

	siteVrfRouteLeakMap := map[string]interface{}{
		"tenantName":        tenantName,
		"vrfRef":            vrfRef,
		"includeAllSubnets": includeAllSubnets,
		"siteIds":           siteIds,
		"prefixsubnet":      prefixSubnets,
	}

	return &PatchPayload{
		Ops:   ops,
		Path:  path,
		Value: siteVrfRouteLeakMap,
	}

}
