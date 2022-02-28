package client

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
)

func (c *Client) GetViaURL(endpoint string) (*container.Container, error) {

	req, err := c.MakeRestRequest("GET", endpoint, nil, true)

	if err != nil {
		return nil, err
	}

	obj, _, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if obj == nil {
		return nil, errors.New("Empty response body")
	}
	return obj, CheckForErrors(obj, "GET")

}

func (c *Client) GetPlatform() string {
	return c.platform
}

func (c *Client) Put(endpoint string, obj models.Model) (*container.Container, error) {
	jsonPayload, err := c.PrepareModel(obj)

	if err != nil {
		return nil, err
	}
	req, err := c.MakeRestRequest("PUT", endpoint, jsonPayload, true)
	if err != nil {
		return nil, err
	}

	cont, _, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return cont, CheckForErrors(cont, "PUT")
}

func (c *Client) Save(endpoint string, obj models.Model) (*container.Container, error) {

	jsonPayload, err := c.PrepareModel(obj)

	if err != nil {
		return nil, err
	}
	req, err := c.MakeRestRequest("POST", endpoint, jsonPayload, true)
	if err != nil {
		return nil, err
	}

	cont, _, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return cont, CheckForErrors(cont, "POST")
}

// CheckForErrors parses the response and checks of there is an error attribute in the response
func CheckForErrors(cont *container.Container, method string) error {

	if cont.Exists("code") && cont.Exists("message") {

		return errors.New(fmt.Sprintf("%s%s", cont.S("message"), cont.S("info")))
	} else {
		return nil
	}

	return nil
}

func (c *Client) DeletebyId(url string) error {

	req, err := c.MakeRestRequest("DELETE", url, nil, true)
	if err != nil {
		return err
	}

	_, resp, err1 := c.Do(req)
	if err1 != nil {
		return err1
	}
	if resp != nil {
		if resp.StatusCode == 204 || resp.StatusCode == 200 {
			return nil
		} else {
			return fmt.Errorf("Unable to delete the object")
		}
	}

	return nil
}

func (c *Client) PatchbyID(endpoint string, objList ...models.Model) (*container.Container, error) {

	contJs := container.New()
	contJs.Array()
	for _, obj := range objList {
		jsonPayload, err := c.PrepareModel(obj)
		if err != nil {
			return nil, err
		}
		contJs.ArrayAppend(jsonPayload.Data())

	}

	// URL encoding
	baseUrl, _ := url.Parse(endpoint)
	qs := url.Values{}
	qs.Add("validate", "false")
	baseUrl.RawQuery = qs.Encode()

	req, err := c.MakeRestRequest("PATCH", baseUrl.String(), contJs, true)
	if err != nil {
		return nil, err
	}

	cont, _, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return cont, CheckForErrors(cont, "PATCH")
}

func (c *Client) PrepareModel(obj models.Model) (*container.Container, error) {
	con, err := obj.ToMap()
	if err != nil {
		return nil, err
	}

	payload := &container.Container{}
	if err != nil {
		return nil, err
	}

	for key, value := range con {
		payload.Set(value, key)
	}
	return payload, nil
}
