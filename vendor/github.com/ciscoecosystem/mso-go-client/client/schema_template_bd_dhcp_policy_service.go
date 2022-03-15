package client

import (
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
)

func (client *Client) CreateTemplateBDDHCPPolicy(obj *models.TemplateBDDHCPPolicy) (*container.Container, error) {
	path := "api/v1/schemas/" + obj.SchemaID
	cont, err := client.PatchbyID(path, models.TemplateBDDHCPPolicyModelForCreation(obj))
	if err != nil {
		return nil, err
	}
	return cont, nil
}

func (client *Client) ReadTemplateBDDHCPPolicy(schemaID string) (*container.Container, error) {
	path := "api/v1/schemas/" + schemaID
	cont, err := client.GetViaURL(path)
	if err != nil {
		return nil, err
	}
	return cont, nil
}

func (client *Client) UpdateTemplateBDDHCPPolicy(obj *models.TemplateBDDHCPPolicy) (*container.Container, error) {
	path := "api/v1/schemas/" + obj.SchemaID
	cont, err := client.PatchbyID(path, models.TemplateBDDHCPPolicyModelForUpdate(obj))
	if err != nil {
		return nil, err
	}
	return cont, nil
}

func (client *Client) DeleteTemplateBDDHCPPolicy(obj *models.TemplateBDDHCPPolicy) (*container.Container, error) {
	path := "api/v1/schemas/" + obj.SchemaID
	cont, err := client.PatchbyID(path, models.TemplateBDDHCPPolicyModelForDeletion(obj))
	if err != nil {
		return nil, CheckForErrors(cont, "PATCH")
	}
	return cont, nil
}
