package mso

import (
	"errors"
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSORemoteLocation() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSORemoteLocationCreate,
		Update: resourceMSORemoteLocationUpdate,
		Read:   resourceMSORemoteLocationRead,
		Delete: resourceMSORemoteLocationDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSORemoteLocationImport,
		},

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"scp",
					"sftp",
				}, false),
			},
			"hostname": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"path": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  22,
			},
			"username": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"password": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"ssh_key": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"passphrase": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		}),
		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			_, password_ok := diff.GetOk("password")
			_, ssh_key_ok := diff.GetOk("ssh_key")
			if password_ok && ssh_key_ok {
				return errors.New(`"password" and "ssh_key" cannot be provided for the same MSO remote location.`)
			}
			if !password_ok && !ssh_key_ok {
				return errors.New(`"password" or "ssh_key" is required to manage the MSO remote location.`)
			}
			return nil
		},
	}
}

func setAuthenticationInState(d *schema.ResourceData) {

	password := d.Get("password").(string)
	if password != "" {
		d.Set("password", password)
	}
	ssh_key := d.Get("ssh_key").(string)
	if ssh_key != "" {
		d.Set("ssh_key", ssh_key)
	}
	passphrase := d.Get("passphrase").(string)
	if passphrase != "" {
		d.Set("passphrase", passphrase)
	}
}

func setRemoteLocation(d *schema.ResourceData, remote map[string]interface{}) {

	d.SetId(remote["id"].(string))
	d.Set("name", remote["name"].(string))
	description, ok := remote["description"].(string)
	if ok {
		d.Set("description", description)
	}
	credential := remote["credential"].(map[string]interface{})
	d.Set("protocol", credential["protocolType"].(string))
	d.Set("hostname", credential["hostname"].(string))
	d.Set("path", credential["remotePath"].(string))
	d.Set("port", credential["port"].(float64))
	d.Set("username", credential["username"].(string))

}

func getCredentialMap(d *schema.ResourceData) map[string]interface{} {

	credentialMap := map[string]interface{}{
		"hostname":     d.Get("hostname").(string),
		"port":         d.Get("port").(int),
		"protocolType": d.Get("protocol").(string),
		"remotePath":   d.Get("path").(string),
		"username":     d.Get("username").(string),
	}

	if password, ok := d.GetOk("password"); ok {
		credentialMap["authType"] = "password"
		credentialMap["password"] = password.(string)
	}

	if ssh_key, ok := d.GetOk("ssh_key"); ok {
		credentialMap["authType"] = "sshKey"
		credentialMap["sshKey"] = ssh_key.(string)
		passphrase := d.Get("passphrase").(string)
		if passphrase != "" {
			credentialMap["passPhrase"] = passphrase
		}
	}

	return credentialMap
}

func resourceMSORemoteLocationImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())

	err := resourceMSORemoteLocationRead(d, m)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil

}

func resourceMSORemoteLocationCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Remote Location: Beginning Create")

	msoClient := m.(*client.Client)

	remoteLocation := models.NewRemoteLocation(d.Get("name").(string), d.Get("description").(string), "", getCredentialMap(d))
	cont, err := msoClient.Save("api/v1/platform/remote-locations", remoteLocation)
	if err != nil {
		return err
	}

	d.SetId(models.StripQuotes(cont.S("id").String()))
	setAuthenticationInState(d)

	return resourceMSORemoteLocationRead(d, m)
}

func resourceMSORemoteLocationUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Remote Location: Beginning Update")

	msoClient := m.(*client.Client)

	remoteLocation := models.NewRemoteLocation(d.Get("name").(string), d.Get("description").(string), d.Id(), getCredentialMap(d))
	_, err := msoClient.Put(fmt.Sprintf("api/v1/platform/remote-locations/%s", d.Id()), remoteLocation)
	if err != nil {
		return err
	}

	setAuthenticationInState(d)

	return resourceMSORemoteLocationRead(d, m)
}

func resourceMSORemoteLocationRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	remoteLocation, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/platform/remote-locations/%s", d.Id()))
	if err != nil {
		return errorForObjectNotFound(err, d.Id(), remoteLocation, d)
	}
	setRemoteLocation(d, remoteLocation.Data().(map[string]interface{}))

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSORemoteLocationDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Remote Location: Beginning Delete")

	msoClient := m.(*client.Client)

	if d.Id() != "" {
		err := msoClient.DeletebyId(fmt.Sprintf("api/v1/platform/remote-locations/%s", d.Id()))
		if err != nil {
			return err
		}
		d.SetId("")
	}

	log.Printf("[DEBUG] Delete finished successfully")
	return nil
}
