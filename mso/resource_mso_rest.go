package mso

import (
	"fmt"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceMSORest() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSORestCreate,
		Read:   resourceMSORestRead,
		Delete: resourceMSORestDelete,
		Update: resourceMSORestUpdate,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"path": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"method": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"payload": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		}),
	}
}

func resourceMSORestCreate(d *schema.ResourceData, m interface{}) error {
	var method, path, payload string
	path = d.Get("path").(string)
	payload = d.Get("payload").(string)

	if tempVar, ok := d.GetOk("method"); ok {
		method = tempVar.(string)
	} else {
		method = "POST"
	}
	msoClient := m.(*client.Client)
	_, err := MakeRestRequest(msoClient, path, method, payload)

	if err != nil {
		return err
	}
	d.SetId(path)
	return resourceMSORestRead(d, m)
}

func resourceMSORestRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceMSORestUpdate(d *schema.ResourceData, m interface{}) error {
	var method, path, payload string
	path = d.Get("path").(string)
	payload = d.Get("payload").(string)

	if tempVar, ok := d.GetOk("method"); ok {
		method = tempVar.(string)
	} else {
		method = "PATCH"
	}
	msoClient := m.(*client.Client)
	_, err := MakeRestRequest(msoClient, path, method, payload)

	if err != nil {
		return err
	}
	d.SetId(path)
	return resourceMSORestRead(d, m)
}

func resourceMSORestDelete(d *schema.ResourceData, m interface{}) error {
	var method, path, payload string
	path = d.Get("path").(string)
	payload = d.Get("payload").(string)

	if tempVar, ok := d.GetOk("method"); ok {
		method = tempVar.(string)
	} else {
		method = "DELETE"
	}
	msoClient := m.(*client.Client)
	_, err := MakeRestRequest(msoClient, path, method, payload)

	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func MakeRestRequest(cli *client.Client, path, method, payload string) (*container.Container, error) {

	jsonPayload, err := container.ParseJSON([]byte(payload))

	if err != nil {
		return nil, fmt.Errorf("Unable to parse the payload to JSON. Please check your payload")
	}

	if len(payload) == 0 {
		jsonPayload = nil
	}

	req, err := cli.MakeRestRequest(method, path, jsonPayload, true)

	if err != nil {
		return nil, err
	}

	respCont, _, err := cli.Do(req)

	return respCont, client.CheckForErrors(respCont, method)
}
