package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceMSOLabel() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOLabelCreate,
		Read:   resourceMSOLabelRead,
		Delete: resourceMSOLabelDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOLabelImport,
		},

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"label": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		}),
	}
}

func resourceMSOLabelImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] Label: Beginning Import")

	msoClient := m.(*client.Client)
	con, err := msoClient.GetViaURL("api/v1/labels" + d.Id())
	if err != nil {
		return nil, err
	}

	d.SetId(models.StripQuotes(con.S("id").String()))
	if con.Exists("displayName") {
		d.Set("label", models.StripQuotes(con.S("displayName").String()))
	}
	if con.Exists("type") {
		d.Set("type", models.StripQuotes(con.S("type").String()))
	}

	log.Printf("[DEBUG] %s: Label Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOLabelCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Label: Beginning Creation")
	msoClient := m.(*client.Client)

	var label string
	if labels, ok := d.GetOk("label"); ok {
		label = labels.(string)
	}

	var types string
	if typed, ok := d.GetOk("type"); ok {
		types = typed.(string)
	}

	labelApp := models.NewLabel("", label, types)

	cont, err := msoClient.Save("api/v1/labels", labelApp)
	if err != nil {
		return err
	}

	id := models.StripQuotes(cont.S("id").String())
	d.SetId(fmt.Sprintf("%v", id))
	log.Printf("[DEBUG] %s: Label Creation finished successfully", d.Id())

	return resourceMSOLabelRead(d, m)
}

func resourceMSOLabelRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	dn := d.Id()
	con, err := msoClient.GetViaURL("api/v1/labels/" + dn)

	if err != nil {
		return errorForObjectNotFound(err, dn, con, d)
	}

	d.SetId(models.StripQuotes(con.S("id").String()))

	if con.Exists("displayName") {
		d.Set("label", models.StripQuotes(con.S("displayName").String()))
	}

	if con.Exists("type") {

		d.Set("type", models.StripQuotes(con.S("type").String()))
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSOLabelDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	msoClient := m.(*client.Client)
	dn := d.Id()
	err := msoClient.DeletebyId("api/v1/labels/" + dn)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
