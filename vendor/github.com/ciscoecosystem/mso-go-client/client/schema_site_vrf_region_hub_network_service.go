package client

import (
	"fmt"

	"github.com/ciscoecosystem/mso-go-client/models"
)

func (client *Client) CreateInterSchemaSiteVrfRegionHubNetwork(obj *models.InterSchemaSiteVrfRegionHubNetork) error {
	schemaCont, err := client.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", obj.SchemaID))
	if err != nil {
		return err
	}
	hubNetwork, err := models.CreateInterSchemaSiteVrfRegionNetworkModel(obj, schemaCont)
	if err != nil {
		return err
	}
	_, err = client.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", obj.SchemaID), hubNetwork)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) DeleteInterSchemaSiteVrfRegionHubNetwork(obj *models.InterSchemaSiteVrfRegionHubNetork) error {
	schemaCont, err := client.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", obj.SchemaID))
	if err != nil {
		return err
	}
	hubNetwork, err := models.DeleteInterSchemaSiteVrfRegionNetworkModel(obj, schemaCont)
	if err != nil {
		return err
	}
	_, err = client.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", obj.SchemaID), hubNetwork)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) ReadInterSchemaSiteVrfRegionHubNetwork(obj *models.InterSchemaSiteVrfRegionHubNetork) (*models.InterSchemaSiteVrfRegionHubNetork, error) {
	schemaCont, err := client.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", obj.SchemaID))
	if err != nil {
		return nil, err
	}
	hubNetwork, err := models.InterSchemaSiteVrfRegionHubNetworkFromContainer(schemaCont, obj)
	if err != nil {
		return nil, err
	}
	return hubNetwork, nil
}
