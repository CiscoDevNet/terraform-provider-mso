package client

import (
	"fmt"
	"log"
	"sync"

	"github.com/ciscoecosystem/mso-go-client/models"
)

var regionHubNetworkMutex sync.Mutex

func (client *Client) CreateInterSchemaSiteVrfRegionHubNetwork(obj *models.InterSchemaSiteVrfRegionHubNetork) error {
	schemaCont, err := client.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", obj.SchemaID))
	if err != nil {
		return err
	}
	hubNetwork := models.CreateInterSchemaSiteVrfRegionNetworkModel(obj, schemaCont)
	regionHubNetworkMutex.Lock()
	log.Printf("hubNetwork: %v\n", hubNetwork)
	_, err = client.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", obj.SchemaID), hubNetwork)
	if err != nil {
		return err
	}
	regionHubNetworkMutex.Unlock()
	return nil
}

func (client *Client) DeleteInterSchemaSiteVrfRegionHubNetwork(obj *models.InterSchemaSiteVrfRegionHubNetork) error {
	schemaCont, err := client.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", obj.SchemaID))
	if err != nil {
		return err
	}
	hubNetwork := models.DeleteInterSchemaSiteVrfRegionNetworkModel(obj, schemaCont)
	regionHubNetworkMutex.Lock()
	_, err = client.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", obj.SchemaID), hubNetwork)
	if err != nil {
		return err
	}
	regionHubNetworkMutex.Unlock()
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
