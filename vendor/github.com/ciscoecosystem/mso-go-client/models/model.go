package models

type Model interface {
	ToMap() (map[string]interface{}, error)
}

type PatchPayload struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

type PatchPayloadList struct {
	Ops   string        `json:",omitempty"`
	Path  string        `json:",omitempty"`
	Value []interface{} `json:",omitempty"`
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

func (patchPayloadListAttributes *PatchPayloadList) ToMap() (map[string]interface{}, error) {
	patchPayloadListMap := make(map[string]interface{})
	A(patchPayloadListMap, "op", patchPayloadListAttributes.Ops)
	A(patchPayloadListMap, "path", patchPayloadListAttributes.Path)
	if patchPayloadListAttributes.Value != nil {
		A(patchPayloadListMap, "value", patchPayloadListAttributes.Value)
	}

	return patchPayloadListMap, nil
}

func GetRemovePatchPayload(path string) *PatchPayload {
	return &PatchPayload{
		Ops:  "remove",
		Path: path,
	}
}
