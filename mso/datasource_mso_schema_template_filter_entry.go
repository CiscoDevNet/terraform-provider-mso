package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceMSOSchemaTemplateFilterEntry() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceMSOSchemaTemplateFilterEntryRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"entry_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"entry_display_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"entry_description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ether_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"arp_flag": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_protocol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"match_only_fragments": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"stateful": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"source_from": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_to": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"destination_from": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"destination_to": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"tcp_session_rules": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		}),
	}
}

func dataSourceMSOSchemaTemplateFilterEntryRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No Template found")
	}
	stateTemplate := d.Get("template_name").(string)
	found := false
	stateFilter := d.Get("name").(string)
	stateFilterEntry := d.Get("entry_name").(string)

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {
			d.Set("template_name", apiTemplate)
			anpCount, err := tempCont.ArrayCount("filters")
			if err != nil {
				return fmt.Errorf("Unable to get Filter list")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "filters")
				if err != nil {
					return err
				}
				apiFilter := models.StripQuotes(anpCont.S("name").String())
				if apiFilter == stateFilter {
					d.Set("name", apiFilter)
					d.Set("display_name", models.StripQuotes(anpCont.S("displayName").String()))
					entriesCount, err := anpCont.ArrayCount("entries")
					if err != nil {
						return fmt.Errorf("Unable to get Entry list")
					}
					for k := 0; k < entriesCount; k++ {
						entriesCont, err := anpCont.ArrayElement(k, "entries")
						if err != nil {
							return err
						}
						apiFilterEntry := models.StripQuotes(entriesCont.S("name").String())
						if apiFilterEntry == stateFilterEntry {
							d.SetId(apiFilterEntry)
							d.Set("entry_name", apiFilterEntry)
							d.Set("entry_display_name", models.StripQuotes(entriesCont.S("displayName").String()))
							if entriesCont.Exists("description") {
								d.Set("entry_description", models.StripQuotes(entriesCont.S("description").String()))
							}
							if entriesCont.Exists("etherType") {
								d.Set("ether_type", models.StripQuotes(entriesCont.S("etherType").String()))
							}
							if entriesCont.Exists("arpFlag") {
								d.Set("arp_flag", models.StripQuotes(entriesCont.S("arpFlag").String()))
							}
							if entriesCont.Exists("ipProtocol") {
								d.Set("ip_protocol", models.StripQuotes(entriesCont.S("ipProtocol").String()))
							}
							if entriesCont.Exists("matchOnlyFragments") {
								d.Set("match_only_fragments", entriesCont.S("matchOnlyFragments").Data().(bool))
							}
							if entriesCont.Exists("stateful") {
								d.Set("stateful", entriesCont.S("stateful").Data().(bool))
							}
							if entriesCont.Exists("sourceFrom") {
								d.Set("source_from", models.StripQuotes(entriesCont.S("sourceFrom").String()))
							}
							if entriesCont.Exists("sourceTo") {
								d.Set("source_to", models.StripQuotes(entriesCont.S("sourceTo").String()))
							}
							if entriesCont.Exists("destinationFrom") {
								d.Set("destination_from", models.StripQuotes(entriesCont.S("destinationFrom").String()))
							}
							if entriesCont.Exists("destinationTo") {
								d.Set("destination_to", models.StripQuotes(entriesCont.S("destinationTo").String()))
							}
							if entriesCont.Exists("tcpSessionRules") {
								d.Set("tcp_session_rules", entriesCont.S("tcpSessionRules").Data().([]interface{}))
							}
							found = true
							break
						}
					}
				}
			}
		}
	}
	if !found {
		d.SetId("")
		return fmt.Errorf("Unable to find the filter %s with entry %s", stateFilter, stateFilterEntry)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
