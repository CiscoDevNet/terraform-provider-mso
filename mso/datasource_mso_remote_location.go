package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceMSORemoteLocation() *schema.Resource {
	return &schema.Resource{

		Read: datasourceMSORemoteRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hostname": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}
}

func getAndSetRemoteLocation(d *schema.ResourceData, remoteLocationsContainer *container.Container, name string) error {
	remoteLocations := remoteLocationsContainer.Search("remoteLocations").Data()
	if remoteLocations == nil || len(remoteLocations.([]interface{})) == 0 {
		return fmt.Errorf("no remote locations found")
	}
	for _, remoteLocation := range remoteLocations.([]interface{}) {
		remoteLocationDetails := remoteLocation.(map[string]interface{})
		if remoteLocationDetails["name"].(string) == name {
			setRemoteLocation(d, remoteLocationDetails)
			return nil
		}
	}
	d.SetId("")
	return fmt.Errorf("Unable to find remote location: %s", name)
}

func datasourceMSORemoteRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Get("name").(string))

	msoClient := m.(*client.Client)

	remoteLocations, err := msoClient.GetViaURL("api/v1/platform/remote-locations")
	if err != nil {
		return err
	}

	err = getAndSetRemoteLocation(d, remoteLocations, d.Get("name").(string))
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
