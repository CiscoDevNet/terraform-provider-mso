package models

type Model interface {
	ToMap() (map[string]interface{}, error)
}

type PatchPayload struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func (patchPayloadAttributes *PatchPayload) ToMap() (map[string]interface{}, error) {
	patchPayloadAttributesMap := make(map[string]interface{})
	A(patchPayloadAttributesMap, "op", patchPayloadAttributes.Ops)
	A(patchPayloadAttributesMap, "path", patchPayloadAttributes.Path)
	if patchPayloadAttributes.Value != nil {
		A(patchPayloadAttributesMap, "value", patchPayloadAttributes.Value)
	}
	return patchPayloadAttributesMap, nil
}

func GetRemovePatchPayload(path string) *PatchPayload {
	return &PatchPayload{
		Ops:  "remove",
		Path: path,
	}
}
