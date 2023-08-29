package models

func NewSystemConfigBanner(alias, bannerState, bannerType, message string) *PatchPayloadList {

	bannerMap := make(map[string]interface{})

	if alias != "" {
		bannerMap["alias"] = alias
	}

	if bannerState != "" && bannerType != "" && message != "" {
		bannerMap["banner"] = map[string]interface{}{
			"bannerState": bannerState,
			"bannerType":  bannerType,
			"message":     message,
		}
	}

	return &PatchPayloadList{
		Ops:   "replace",
		Path:  "/bannerConfig",
		Value: []interface{}{bannerMap},
	}

}

func NewSystemConfigChangeControl(enable bool, numOfApprovers int) *PatchPayload {

	changeControlMap := map[string]interface{}{
		"enable": enable,
	}

	if numOfApprovers > 0 {
		changeControlMap["numOfApprovers"] = numOfApprovers
	}

	return &PatchPayload{
		Ops:   "replace",
		Path:  "/changeControl",
		Value: changeControlMap,
	}

}
