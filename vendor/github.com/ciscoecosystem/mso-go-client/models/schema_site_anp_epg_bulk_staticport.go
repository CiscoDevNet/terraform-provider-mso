package models

func NewSchemaSiteAnpEpgBulkStaticPort(ops, path string, staticPortsList []interface{}) *PatchPayloadList {

	return &PatchPayloadList{
		Ops:   ops,
		Path:  path,
		Value: staticPortsList,
	}

}
