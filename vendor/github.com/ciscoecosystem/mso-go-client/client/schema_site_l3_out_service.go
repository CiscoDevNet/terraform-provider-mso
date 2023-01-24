package client

import (
	"fmt"

	"github.com/ciscoecosystem/mso-go-client/models"
)

func (client *Client) CreateIntersiteL3outs(obj *models.IntersiteL3outs) error {
	l3out := models.CreateIntersiteL3outsModel(obj)
	_, err := client.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", obj.SchemaID), l3out)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) DeleteIntersiteL3outs(obj *models.IntersiteL3outs) error {
	l3out := models.DeleteIntersiteL3outsModel(obj)
	_, err := client.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", obj.SchemaID), l3out)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) ReadIntersiteL3outs(obj *models.IntersiteL3outs) (*models.IntersiteL3outs, error) {
	schemaCont, err := client.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", obj.SchemaID))
	if err != nil {
		return nil, err
	}
	l3out, err := models.IntersiteL3outsFromContainer(schemaCont, obj)
	if err != nil {
		return nil, err
	}
	return l3out, nil
}
