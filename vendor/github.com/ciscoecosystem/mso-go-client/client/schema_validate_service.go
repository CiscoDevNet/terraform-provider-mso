package client

import (
	"errors"
	"fmt"

	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
)

func (client *Client) ReadSchemaValidate(obj *models.SchemValidate) (*models.SchemValidate, error) {
	cont, err := client.GetSchemaValidate(fmt.Sprintf("api/v1/schemas/%s/validate", obj.SchmaId))
	if err != nil {
		return nil, err
	}
	remoteSchemaValidate := models.SchemValidate{
		SchmaId: obj.SchmaId,
		Result:  models.G(cont, "result"),
	}
	return &remoteSchemaValidate, nil
}

func (c *Client) GetSchemaValidate(endpoint string) (*container.Container, error) {

	req, err := c.MakeRestRequest("GET", endpoint, nil, true)

	if err != nil {
		return nil, err
	}
	req.Header.Del("Content-Type")

	obj, _, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if obj == nil {
		return nil, errors.New("empty response body")
	}
	return obj, CheckForErrors(obj, "GET")

}
