package models

func NewTemplateContractFilterRelationShip(ops, path, action, priority, desc string, filterRef map[string]interface{}, directives []interface{}) *PatchPayload {

	filterMap := map[string]interface{}{
		"filterRef": filterRef,
	}

	if len(directives) > 0 {
		filterMap["directives"] = directives
	}

	if action != "" {
		filterMap["action"] = action
	}

	if priority != "" {
		filterMap["priorityOverride"] = priority
	}

	return &PatchPayload{
		Ops:   ops,
		Path:  path,
		Value: filterMap,
	}

}
