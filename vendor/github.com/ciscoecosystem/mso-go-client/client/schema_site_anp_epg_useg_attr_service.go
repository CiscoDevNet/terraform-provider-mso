package client

import (
	"fmt"

	"github.com/ciscoecosystem/mso-go-client/models"
)

func (client *Client) CreateAnpEpgUsegAttr(obj *models.SiteUsegAttr) error {
	useg := models.SiteAnpEpgUsegAttrForCreation(obj)
	_, err := client.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", obj.SchemaID), useg)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) DeleteAnpEpgUsegAttr(obj *models.SiteUsegAttr) error {
	_, useg_index, read_err := client.ReadAnpEpgUsegAttr(obj)
	if read_err != nil {
		return read_err
	}
	useg := models.SiteAnpEpgUsegAttrforDeletion(obj, useg_index)
	_, err := client.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", obj.SchemaID), useg)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) UpdateAnpEpgUsegAttr(obj *models.SiteUsegAttr) error {
	_, useg_index, read_err := client.ReadAnpEpgUsegAttr(obj)
	if read_err != nil {
		return read_err
	}
	useg := models.SiteAnpEpgUsegAttrforUpdate(obj, useg_index)
	_, err := client.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", obj.SchemaID), useg)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) ReadAnpEpgUsegAttr(obj *models.SiteUsegAttr) (*models.SiteUsegAttr, int, error) {
	schemaCont, err := client.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", obj.SchemaID))
	if err != nil {
		return nil, -1, err
	}
	useg, useg_index, err := models.SiteAnpEpgUsegAttrFromContainer(schemaCont, obj)
	if err != nil {
		return nil, -1, err
	}
	return useg, useg_index, nil
}
