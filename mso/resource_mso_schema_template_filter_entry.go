package mso

import (
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOSchemaTemplateFilterEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaTemplateFilterEntryCreate,
		Read:   resourceMSOSchemaTemplateFilterEntryRead,
		Update: resourceMSOSchemaTemplateFilterEntryUpdate,
		Delete: resourceMSOSchemaTemplateFilterEntryDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaTemplateFilterEntryImport,
		},

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"display_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"entry_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"entry_display_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"entry_description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ether_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"arp_flag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip_protocol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"match_only_fragments": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"stateful": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"source_from": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_to": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"destination_from": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"destination_to": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tcp_session_rules": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		}),
	}
}

func resourceMSOSchemaTemplateFilterEntryImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())

	msoClient := m.(*client.Client)
	get_attribute := strings.Split(d.Id(), "/")
	schemaId := get_attribute[0]
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return nil, err
	}
	d.Set("schema_id", schemaId)
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return nil, fmt.Errorf("No Template found")
	}
	stateTemplate := get_attribute[2]
	found := false
	stateFilter := get_attribute[4]
	stateFilterEntry := get_attribute[6]

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return nil, err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {
			d.Set("template_name", apiTemplate)
			anpCount, err := tempCont.ArrayCount("filters")
			if err != nil {
				return nil, fmt.Errorf("Unable to get Filter list")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "filters")
				if err != nil {
					return nil, err
				}
				apiFilter := models.StripQuotes(anpCont.S("name").String())
				if apiFilter == stateFilter {
					d.Set("name", apiFilter)
					d.Set("display_name", models.StripQuotes(anpCont.S("displayName").String()))
					entriesCount, err := anpCont.ArrayCount("entries")
					if err != nil {
						return nil, fmt.Errorf("Unable to get Entry list")
					}
					for k := 0; k < entriesCount; k++ {
						entriesCont, err := anpCont.ArrayElement(k, "entries")
						if err != nil {
							return nil, err
						}
						apiFilterEntry := models.StripQuotes(entriesCont.S("name").String())
						if apiFilterEntry == stateFilterEntry {
							d.SetId(apiFilterEntry)
							d.Set("entry_name", apiFilterEntry)
							d.Set("entry_display_name", models.StripQuotes(entriesCont.S("displayName").String()))
							if entriesCont.Exists("descirption") {
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
		return nil, fmt.Errorf("Unable to find the filter entry")
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaTemplateFilterEntryCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Filter Entry Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	stateTemplate := d.Get("template_name").(string)
	filterName := d.Get("name").(string)
	displayFilterName := d.Get("display_name").(string)
	entryName := d.Get("entry_name").(string)
	entryDisplayName := d.Get("entry_display_name").(string)

	var entryDescription, etherType, arpFlag, ipProtocol, sourceFrom, sourceTo, destinationFrom, destinationTo string
	var matchOnlyFragments, stateful bool
	var tcpSessionRules []interface{}
	entries := make([]interface{}, 0, 1)
	entryMap := make(map[string]interface{})
	entryMap["name"] = entryName
	entryMap["displayName"] = entryDisplayName

	if tempVar, ok := d.GetOk("entry_description"); ok {
		entryDescription = tempVar.(string)
		entryMap["description"] = entryDescription
	}
	if tempVar, ok := d.GetOk("ether_type"); ok {
		etherType = tempVar.(string)
		entryMap["etherType"] = etherType
	}
	if tempVar, ok := d.GetOk("arp_flag"); ok {
		arpFlag = tempVar.(string)
		entryMap["arpFlag"] = arpFlag
	}
	if tempVar, ok := d.GetOk("ip_protocol"); ok {
		ipProtocol = tempVar.(string)
		entryMap["ipProtocol"] = ipProtocol
	}
	if tempVar, ok := d.GetOk("source_from"); ok {
		sourceFrom = tempVar.(string)
		entryMap["sourceFrom"] = sourceFrom

	}
	if tempVar, ok := d.GetOk("source_to"); ok {
		sourceTo = tempVar.(string)
		entryMap["sourceTo"] = sourceTo
	}
	if tempVar, ok := d.GetOk("destination_from"); ok {
		destinationFrom = tempVar.(string)
		entryMap["destinationFrom"] = destinationFrom

	}
	if tempVar, ok := d.GetOk("destination_to"); ok {
		destinationTo = tempVar.(string)
		entryMap["destinationTo"] = destinationTo
	}
	if tempVar, ok := d.GetOk("match_only_fragments"); ok {
		matchOnlyFragments = tempVar.(bool)
		entryMap["matchOnlyFragments"] = matchOnlyFragments
	}
	if tempVar, ok := d.GetOk("stateful"); ok {
		stateful = tempVar.(bool)
		entryMap["stateful"] = stateful
	}
	if tempVar, ok := d.GetOk("tcp_session_rules"); ok {
		tcpSessionRules = tempVar.([]interface{})
		entryMap["tcpSessionRules"] = tcpSessionRules
	}
	entries = append(entries, entryMap)
	foundEntry := false
	foundFilter := false
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("templates")

	if err != nil {
		return fmt.Errorf("No Template found")
	}
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}

		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {

			filterCount, err := tempCont.ArrayCount("filters")
			if err != nil {
				return fmt.Errorf("Unable to get Filter list")
			}

			for j := 0; j < filterCount; j++ {
				filterCont, err := tempCont.ArrayElement(j, "filters")
				if err != nil {
					return err
				}
				apiFilterName := models.StripQuotes(filterCont.S("name").String())

				if apiFilterName == filterName {

					foundFilter = true
					entriesCount, err := filterCont.ArrayCount("entries")
					if err == nil {
						for k := 0; k < entriesCount; k++ {
							entriesCont, err1 := filterCont.ArrayElement(k, "entries")
							if err1 != nil {
								return err1
							}

							apiFilterEntry := models.StripQuotes(entriesCont.S("name").String())
							if apiFilterEntry == entryName {
								foundEntry = true
								err2 := resourceMSOSchemaTemplateFilterEntryUpdate(d, m)
								if err2 != nil {
									return err2
								}
							}
						}
					}
					if !foundEntry {
						pathf := fmt.Sprintf("/templates/%s/filters/%s/entries/-", stateTemplate, filterName)
						filterStruct := models.NewTemplateFilterEntry("add", pathf, entryName, entryDisplayName, entryDescription, etherType, arpFlag, ipProtocol, sourceFrom, sourceTo, destinationFrom, destinationTo, matchOnlyFragments, stateful, tcpSessionRules)
						_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), filterStruct)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}
	if !foundFilter {
		pathf := fmt.Sprintf("/templates/%s/filters/-", stateTemplate)
		filterStruct := models.NewTemplateFilter("add", pathf, filterName, displayFilterName, entries)
		_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), filterStruct)
		if err != nil {
			return err
		}
	}
	return resourceMSOSchemaTemplateFilterEntryRead(d, m)
}

func resourceMSOSchemaTemplateFilterEntryRead(d *schema.ResourceData, m interface{}) error {
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
							if entriesCont.Exists("descirption") {
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
							} else {
								d.Set("tcp_session_rules", make([]interface{}, 0))
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
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}

func resourceMSOSchemaTemplateFilterEntryUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Filter Entry: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	stateTemplate := d.Get("template_name").(string)
	filterName := d.Get("name").(string)
	entryName := d.Get("entry_name").(string)
	entryDisplayName := d.Get("entry_display_name").(string)

	var entryDescription, etherType, arpFlag, ipProtocol, sourceFrom, sourceTo, destinationFrom, destinationTo string
	var matchOnlyFragments, stateful bool
	var tcpSessionRules []interface{}

	if tempVar, ok := d.GetOk("entry_description"); ok {
		entryDescription = tempVar.(string)

	}
	if tempVar, ok := d.GetOk("ether_type"); ok {
		etherType = tempVar.(string)

	}
	if tempVar, ok := d.GetOk("arp_flag"); ok {
		arpFlag = tempVar.(string)

	}
	if tempVar, ok := d.GetOk("ip_protocol"); ok {
		ipProtocol = tempVar.(string)

	}
	if tempVar, ok := d.GetOk("source_from"); ok {
		sourceFrom = tempVar.(string)

	}
	if tempVar, ok := d.GetOk("source_to"); ok {
		sourceTo = tempVar.(string)

	}
	if tempVar, ok := d.GetOk("destination_from"); ok {
		destinationFrom = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("destination_to"); ok {

		destinationTo = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("match_only_fragments"); ok {
		matchOnlyFragments = tempVar.(bool)

	}
	if tempVar, ok := d.GetOk("stateful"); ok {
		stateful = tempVar.(bool)

	}
	if tempVar, ok := d.GetOk("tcp_session_rules"); ok {

		tcpSessionRules = tempVar.([]interface{})

	}

	pathf := fmt.Sprintf("/templates/%s/filters/%s/entries/%s", stateTemplate, filterName, entryName)
	filterStruct := models.NewTemplateFilterEntry("replace", pathf, entryName, entryDisplayName, entryDescription, etherType, arpFlag, ipProtocol, sourceFrom, sourceTo, destinationFrom, destinationTo, matchOnlyFragments, stateful, tcpSessionRules)
	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), filterStruct)
	if err != nil {
		return err
	}
	return resourceMSOSchemaTemplateFilterEntryRead(d, m)
}

func resourceMSOSchemaTemplateFilterEntryDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Filter Entry: Beginning Deletion")
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	stateTemplate := d.Get("template_name").(string)
	filterName := d.Get("name").(string)
	displayName := d.Get("display_name").(string)
	entryName := d.Get("entry_name").(string)
	entryDisplayName := d.Get("entry_display_name").(string)

	var entryDescription, etherType, arpFlag, ipProtocol, sourceFrom, sourceTo, destinationFrom, destinationTo string
	var matchOnlyFragments, stateful bool
	var tcpSessionRules []interface{}
	entries := make([]interface{}, 0, 1)
	entryMap := make(map[string]interface{})
	entryMap["name"] = entryName
	entryMap["displayName"] = entryDisplayName

	if tempVar, ok := d.GetOk("entry_description"); ok {
		entryDescription = tempVar.(string)
		entryMap["description"] = entryDescription
	}
	if tempVar, ok := d.GetOk("ether_type"); ok {
		etherType = tempVar.(string)
		entryMap["etherType"] = etherType
	}
	if tempVar, ok := d.GetOk("arp_flag"); ok {
		arpFlag = tempVar.(string)
		entryMap["arpFlag"] = arpFlag
	}
	if tempVar, ok := d.GetOk("ip_protocol"); ok {
		ipProtocol = tempVar.(string)
		entryMap["ipProtocol"] = ipProtocol
	}
	if tempVar, ok := d.GetOk("source_from"); ok {
		sourceFrom = tempVar.(string)
		entryMap["sourceFrom"] = sourceFrom

	}
	if tempVar, ok := d.GetOk("source_to"); ok {
		sourceTo = tempVar.(string)
		entryMap["sourceTo"] = sourceTo
	}
	if tempVar, ok := d.GetOk("destination_from"); ok {
		destinationFrom = tempVar.(string)
		entryMap["destinationFrom"] = destinationFrom

	}
	if tempVar, ok := d.GetOk("destination_to"); ok {
		destinationTo = tempVar.(string)
		entryMap["destinationTo"] = destinationTo
	}
	if tempVar, ok := d.GetOk("match_only_fragments"); ok {
		matchOnlyFragments = tempVar.(bool)
		entryMap["matchOnlyFragments"] = matchOnlyFragments
	}
	if tempVar, ok := d.GetOk("stateful"); ok {
		stateful = tempVar.(bool)
		entryMap["stateful"] = stateful
	}
	if tempVar, ok := d.GetOk("tcp_session_rules"); ok {

		tcpSessionRules = tempVar.([]interface{})
		entryMap["tcpSessionRules"] = tcpSessionRules
	}
	entries = append(entries, entryMap)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No Template found")
	}
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}

		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {
			filterCount, err := tempCont.ArrayCount("filters")
			if err != nil {
				return fmt.Errorf("Unable to get Filter list")
			}

			for j := 0; j < filterCount; j++ {
				filterCont, err := tempCont.ArrayElement(j, "filters")
				if err != nil {
					return err
				}

				apiFilterName := models.StripQuotes(filterCont.S("name").String())
				if apiFilterName == "" {
					return fmt.Errorf("There was no filter with name %s to delete", apiFilterName)
				}
				if apiFilterName == filterName {
					entriesCount, err := filterCont.ArrayCount("entries")
					if err != nil {
						return fmt.Errorf("Unable to get Entry list")
					}
					if entriesCount == 1 {
						path := fmt.Sprintf("/templates/%s/filters/%s", apiTemplate, apiFilterName)
						filterStruct := models.NewTemplateFilter("remove", path, apiFilterName, displayName, entries)
						response, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), filterStruct)

						// Ignoring Error with code 141: Resource Not Found when deleting
						if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
							return err
						}
					} else {
						pathf := fmt.Sprintf("/templates/%s/filters/%s/entries/%s", stateTemplate, filterName, entryName)
						filterStruct := models.NewTemplateFilterEntry("remove", pathf, entryName, entryDisplayName, entryDescription, etherType, arpFlag, ipProtocol, sourceFrom, sourceTo, destinationFrom, destinationTo, matchOnlyFragments, stateful, tcpSessionRules)
						response, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), filterStruct)

						// Ignoring Error with code 141: Resource Not Found when deleting
						if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
							return err
						}
					}
				}
			}
		}
	}

	d.SetId("")
	return resourceMSOSchemaTemplateFilterEntryRead(d, m)

}
