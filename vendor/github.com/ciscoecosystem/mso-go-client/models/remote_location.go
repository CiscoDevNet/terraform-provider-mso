package models

type RemoteLocation struct {
	Name        string                 `json:",omitempty"`
	Description string                 `json:",omitempty"`
	Id          string                 `json:",omitempty"`
	Credential  map[string]interface{} `json:",omitempty"`
}

func NewRemoteLocation(name, description, id string, credential map[string]interface{}) *RemoteLocation {
	return &RemoteLocation{Name: name, Description: description, Id: id, Credential: credential}
}

func (remoteLocation *RemoteLocation) ToMap() (map[string]interface{}, error) {
	remoteLocationMap := make(map[string]interface{})
	A(remoteLocationMap, "name", remoteLocation.Name)
	A(remoteLocationMap, "description", remoteLocation.Description)
	A(remoteLocationMap, "id", remoteLocation.Id)
	A(remoteLocationMap, "credential", remoteLocation.Credential)
	return remoteLocationMap, nil
}
