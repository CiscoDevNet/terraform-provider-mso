package mso

import (
	"fmt"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var HTTP_METHODS = []string{"GET", "PUT", "PATCH", "POST", "DELETE"}

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
			"retrigger": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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

	if !contains(HTTP_METHODS, method) {
		return fmt.Errorf("Invalid method %s passed", method)
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
	d.Set("retrigger", false)
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
	if !contains(HTTP_METHODS, method) {
		return fmt.Errorf("Invalid method %s passed", method)
	}
	msoClient := m.(*client.Client)
	_, err := MakeRestRequest(msoClient, path, method, payload)

	if err != nil {
		return err
	}
	d.SetId(path)
	return resourceMSORestRead(d, m)
}

// resourceMSORestDelete has surprising semantics that are preserved for
// backwards compatibility:
//
//   - A DELETE request against `path` is issued ONLY when the `method`
//     attribute is unset in state (`d.GetOk("method")` returns false). In
//     that case the function uses the hard-coded method "DELETE" and ignores
//     whatever the user may have configured.
//   - When `method` is set in state to ANY value -- including "DELETE" --
//     the destroy short-circuits, performs no HTTP call, and simply clears
//     the resource id. The underlying object on MSO is left untouched.
//
// Consequences for callers:
//
//   - To get a real DELETE on `terraform destroy`, leave `method` unset on
//     the resource for its entire lifecycle. Once `method` has ever been
//     set, it stays in state (the attribute is Optional+Computed) and the
//     destroy short-circuits.
//   - If you need to remove the remote object while `method` is already
//     set, issue an Update with `method = "DELETE"` (the Update path
//     respects the configured method). That is what the
//     `Update with method=DELETE` step in resource_mso_rest_test.go
//     exercises.
//
// This shape is non-obvious -- a user might reasonably expect
// `method = "DELETE"` to fire a DELETE on destroy too -- but changing it
// would break existing configurations that rely on the short-circuit to
// keep a remote object alive across `terraform destroy` runs.
//
// History: this behaviour was introduced in commit 76f1ca6b3ffb70012973bf80657a4ca97073d305
// (https://github.com/CiscoDevNet/terraform-provider-mso/commit/76f1ca6b3ffb70012973bf80657a4ca97073d305),
// which added the `method` short-circuit to resourceMSORestDelete.
func resourceMSORestDelete(d *schema.ResourceData, m interface{}) error {
	var method, path, payload string
	path = d.Get("path").(string)
	payload = d.Get("payload").(string)

	if _, ok := d.GetOk("method"); !ok {
		method = "DELETE"
		msoClient := m.(*client.Client)
		_, err := MakeRestRequest(msoClient, path, method, payload)

		if err != nil {
			return err
		}
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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
