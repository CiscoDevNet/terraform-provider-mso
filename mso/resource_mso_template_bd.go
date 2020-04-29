package mso

import (
	"fmt"
	"log"
	"regexp"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceMSOTemplateBD() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOTemplateBDCreate,
		Read:   resourceMSOTemplateBDRead,
		Update: resourceMSOTemplateBDUpdate,
		Delete: resourceMSOTemplateBDDelete,

		// Importer: &schema.ResourceImporter{
		//     State: resourceMSOSchemaSiteImport,
		// },

		SchemaVersion: 1,

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

			"intersite_bum_traffic": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"optimize_wan_bandwidth": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"layer2_stretch": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"layer3_multicast": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"layer2_unknown_unicast": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vrf_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vrf_schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vrf_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dhcp_policy": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"version": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"dhcp_option_policy_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"dhcp_option_policy_version": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		}),
	}
}

func resourceMSOTemplateBDCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template BD: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	name := d.Get("name").(string)
	displayName := d.Get("display_name").(string)
	templateName := d.Get("template_name").(string)
	vrfName := d.Get("vrf_name").(string)

	var intersite_bum_traffic, optimize_wan_bandwidth, layer2_stretch, layer3_multicast bool
	var layer2_unknown_unicast, vrf_schema_id, vrf_template_name string

	if tempVar, ok := d.GetOk("intersite_bum_traffic"); ok {
		intersite_bum_traffic = tempVar.(bool)
	}
	if tempVar, ok := d.GetOk("optimize_wan_bandwidth"); ok {
		optimize_wan_bandwidth = tempVar.(bool)
	}
	if tempVar, ok := d.GetOk("layer2_stretch"); ok {
		layer2_stretch = tempVar.(bool)
	}
	if tempVar, ok := d.GetOk("layer3_multicast"); ok {
		layer3_multicast = tempVar.(bool)
	}
	if tempVar, ok := d.GetOk("layer2_unknown_unicast"); ok {
		layer2_unknown_unicast = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("vrf_schema_id"); ok {
		vrf_schema_id = tempVar.(string)
	} else {
		vrf_schema_id = schemaID
	}
	if tempVar, ok := d.GetOk("vrf_template_name"); ok {
		vrf_template_name = tempVar.(string)
	} else {
		vrf_template_name = templateName
	}
	var dhcpPolMap map[string]interface{}
	if tempVar, ok := d.GetOk("dhcp_policy"); ok {
		dhcp_policy := tempVar.(map[string]interface{})
		dhcpPolMap["name"] = dhcp_policy["name"]
		dhcpPolMap["version"] = dhcp_policy["version"]

		optionName := dhcp_policy["dhcp_option_policy_name"]
		optionVersion := dhcp_policy["dhcp_option_policy_version"]
		if optionName != nil {
			if optionVersion != nil {
				dhcpOptionMap := make(map[string]interface{})
				dhcpOptionMap["name"] = optionName
				dhcpOptionMap["version"] = optionVersion
				dhcpPolMap["dhcpOptionLabel"] = dhcpOptionMap
			} else {
				return fmt.Errorf("dhcp_option_policy_version is required with dhcp_option_policy_name")
			}
		}
	} else {
		dhcpPolMap = nil
	}

	vrfRefMap := make(map[string]interface{})
	vrfRefMap["schemaId"] = vrf_schema_id
	vrfRefMap["templateName"] = vrf_template_name
	vrfRefMap["vrfName"] = vrfName
	path := fmt.Sprintf("/templates/%s/bds/-", templateName)
	bdStruct := models.NewTemplateBD("add", path, name, displayName, layer2_unknown_unicast, intersite_bum_traffic, optimize_wan_bandwidth, layer2_stretch, layer3_multicast, vrfRefMap, dhcpPolMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), bdStruct)

	if err != nil {
		return err
	}
	return resourceMSOTemplateBDRead(d, m)
}

func resourceMSOTemplateBDRead(d *schema.ResourceData, m interface{}) error {
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
	stateBD := d.Get("name")
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {
			bdCount, err := tempCont.ArrayCount("bds")
			if err != nil {
				return fmt.Errorf("Unable to get BD list")
			}
			for j := 0; j < bdCount; j++ {
				bdCont, err := tempCont.ArrayElement(j, "bds")
				if err != nil {
					return err
				}
				apiBD := models.StripQuotes(bdCont.S("name").String())
				if apiBD == stateBD {
					d.SetId(apiBD)
					d.Set("name", apiBD)
					d.Set("schema_id", schemaId)
					d.Set("template_name", apiTemplate)
					d.Set("display_name", models.StripQuotes(bdCont.S("displayName").String()))
					d.Set("layer2_unknown_unicast", models.StripQuotes(bdCont.S("l2UnknownUnicast").String()))
					d.Set("intersite_bum_traffic", bdCont.S("intersiteBumTrafficAllow").Data().(bool))
					d.Set("optimize_wan_bandwidth", bdCont.S("optimizeWanBandwidth").Data().(bool))
					d.Set("layer3_multicast", bdCont.S("l3MCast").Data().(bool))
					d.Set("layer2_stretch", bdCont.S("l2Stretch").Data().(bool))

					vrfRef := models.StripQuotes(bdCont.S("vrfRef").String())
					re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/vrfs/(.*)")
					match := re.FindStringSubmatch(vrfRef)
					d.Set("vrf_name", match[3])
					d.Set("vrf_schema_id", match[1])
					d.Set("vrf_template_name", match[2])

					if bdCont.Exists("dhcpLabel") {
						dhcpPolMap := make(map[string]interface{})
						dhcpPolMap["name"] = models.StripQuotes(bdCont.S("dhcpLabel", "name").String())
						dhcpPolMap["version"] = models.StripQuotes(bdCont.S("dhcpLabel", "version").String())
						if bdCont.Exists("dhcpLabel", "dhcpOptionLabel") {
							dhcpPolMap["dhcp_option_policy_name"] = models.StripQuotes(bdCont.S("dhcpLabel", "dhcpOptionLabel", "name").String())
							dhcpPolMap["dhcp_option_policy_version"] = models.StripQuotes(bdCont.S("dhcpLabel", "dhcpOptionLabel", "version").String())
						}
						d.Set("dhcp_policy", dhcpPolMap)
					} else {
						d.Set("dhcp_policy", make(map[string]interface{}))
					}
					found = true
					break
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

func resourceMSOTemplateBDUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template BD: Beginning Update")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	name := d.Get("name").(string)
	displayName := d.Get("display_name").(string)
	templateName := d.Get("template_name").(string)
	vrfName := d.Get("vrf_name").(string)

	var intersite_bum_traffic, optimize_wan_bandwidth, layer2_stretch, layer3_multicast bool
	var layer2_unknown_unicast, vrf_schema_id, vrf_template_name string

	if tempVar, ok := d.GetOk("intersite_bum_traffic"); ok {
		intersite_bum_traffic = tempVar.(bool)
	}
	if tempVar, ok := d.GetOk("optimize_wan_bandwidth"); ok {
		optimize_wan_bandwidth = tempVar.(bool)
	}
	if tempVar, ok := d.GetOk("layer2_stretch"); ok {
		layer2_stretch = tempVar.(bool)
	}
	if tempVar, ok := d.GetOk("layer3_multicast"); ok {
		layer3_multicast = tempVar.(bool)
	}
	if tempVar, ok := d.GetOk("layer2_unknown_unicast"); ok {
		layer2_unknown_unicast = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("vrf_schema_id"); ok {
		vrf_schema_id = tempVar.(string)
	} else {
		vrf_schema_id = schemaID
	}
	if tempVar, ok := d.GetOk("vrf_template_name"); ok {
		vrf_template_name = tempVar.(string)
	} else {
		vrf_template_name = templateName
	}
	var dhcpPolMap map[string]interface{}
	if tempVar, ok := d.GetOk("dhcp_policy"); ok {
		dhcp_policy := tempVar.(map[string]interface{})
		dhcpPolMap["name"] = dhcp_policy["name"]
		dhcpPolMap["version"] = dhcp_policy["version"]

		optionName := dhcp_policy["dhcp_option_policy_name"]
		optionVersion := dhcp_policy["dhcp_option_policy_version"]
		if optionName != nil {
			if optionVersion != nil {
				dhcpOptionMap := make(map[string]interface{})
				dhcpOptionMap["name"] = optionName
				dhcpOptionMap["version"] = optionVersion
				dhcpPolMap["dhcpOptionLabel"] = dhcpOptionMap
			} else {
				return fmt.Errorf("dhcp_option_policy_version is required with dhcp_option_policy_name")
			}
		}
	} else {
		dhcpPolMap = nil
	}

	vrfRefMap := make(map[string]interface{})
	vrfRefMap["schemaId"] = vrf_schema_id
	vrfRefMap["templateName"] = vrf_template_name
	vrfRefMap["vrfName"] = vrfName
	path := fmt.Sprintf("/templates/%s/bds/%s", templateName, name)
	bdStruct := models.NewTemplateBD("replace", path, name, displayName, layer2_unknown_unicast, intersite_bum_traffic, optimize_wan_bandwidth, layer2_stretch, layer3_multicast, vrfRefMap, dhcpPolMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), bdStruct)

	if err != nil {
		return err
	}
	return resourceMSOTemplateBDRead(d, m)
}

func resourceMSOTemplateBDDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template BD: Beginning Update")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	name := d.Get("name").(string)
	displayName := d.Get("display_name").(string)
	templateName := d.Get("template_name").(string)
	vrfName := d.Get("vrf_name").(string)

	var intersite_bum_traffic, optimize_wan_bandwidth, layer2_stretch, layer3_multicast bool
	var layer2_unknown_unicast, vrf_schema_id, vrf_template_name string

	if tempVar, ok := d.GetOk("intersite_bum_traffic"); ok {
		intersite_bum_traffic = tempVar.(bool)
	}
	if tempVar, ok := d.GetOk("optimize_wan_bandwidth"); ok {
		optimize_wan_bandwidth = tempVar.(bool)
	}
	if tempVar, ok := d.GetOk("layer3_multicast"); ok {
		layer3_multicast = tempVar.(bool)
	}
	if tempVar, ok := d.GetOk("layer2_unknown_unicast"); ok {
		layer2_unknown_unicast = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("vrf_schema_id"); ok {
		vrf_schema_id = tempVar.(string)
	} else {
		vrf_schema_id = schemaID
	}
	if tempVar, ok := d.GetOk("vrf_template_name"); ok {
		vrf_template_name = tempVar.(string)
	} else {
		vrf_template_name = templateName
	}
	var dhcpPolMap map[string]interface{}
	if tempVar, ok := d.GetOk("dhcp_policy"); ok {
		dhcp_policy := tempVar.(map[string]interface{})
		dhcpPolMap["name"] = dhcp_policy["name"]
		dhcpPolMap["version"] = dhcp_policy["version"]

		optionName := dhcp_policy["dhcp_option_policy_name"]
		optionVersion := dhcp_policy["dhcp_option_policy_version"]
		if optionName != nil {
			if optionVersion != nil {
				dhcpOptionMap := make(map[string]interface{})
				dhcpOptionMap["name"] = optionName
				dhcpOptionMap["version"] = optionVersion
				dhcpPolMap["dhcpOptionLabel"] = dhcpOptionMap
			} else {
				return fmt.Errorf("dhcp_option_policy_version is required with dhcp_option_policy_name")
			}
		}
	} else {
		dhcpPolMap = nil
	}

	vrfRefMap := make(map[string]interface{})
	vrfRefMap["schemaId"] = vrf_schema_id
	vrfRefMap["templateName"] = vrf_template_name
	vrfRefMap["vrfName"] = vrfName
	path := fmt.Sprintf("/templates/%s/bds/%s", templateName, name)
	bdStruct := models.NewTemplateBD("remove", path, name, displayName, layer2_unknown_unicast, intersite_bum_traffic, optimize_wan_bandwidth, layer2_stretch, layer3_multicast, vrfRefMap, dhcpPolMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), bdStruct)

	if err != nil {
		return err
	}
	d.SetId("")
	return resourceMSOTemplateBDRead(d, m)
	return nil
}
