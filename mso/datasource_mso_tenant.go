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

func datasourceMSOTenant() *schema.Resource {
	return &schema.Resource{

		Read: datasourceMSOTenantRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_associations": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Computed: true,
			},
			"site_associations": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"site_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_domains": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"vendor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gcp_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gcp_access_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gcp_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gcp_key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gcp_private_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gcp_client_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gcp_email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"aws_account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_aws_account_trusted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"aws_access_key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"aws_secret_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"azure_subscription_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"azure_access_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"azure_shared_account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"azure_application_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"azure_client_secret": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"azure_active_directory_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Computed: true,
			},
		}),
	}
}

func datasourceMSOTenantRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	name := d.Get("name").(string)
	con, err := msoClient.GetViaURL("api/v1/tenants")
	if err != nil {
		return err
	}

	data := con.S("tenants").Data().([]interface{})
	var flag bool
	var count int
	for _, info := range data {
		val := info.(map[string]interface{})
		if val["name"].(string) == name {
			flag = true
			break
		}
		count = count + 1
	}

	if flag != true {
		return fmt.Errorf("Tenant of specified name not found")
	}

	dataCon := con.S("tenants").Index(count)

	d.SetId(models.StripQuotes(dataCon.S("id").String()))

	d.Set("name", models.StripQuotes(dataCon.S("name").String()))

	d.Set("display_name", models.StripQuotes(dataCon.S("displayName").String()))

	if dataCon.Exists("description") {
		d.Set("description", models.StripQuotes(dataCon.S("description").String()))
	}

	count1, _ := dataCon.ArrayCount("siteAssociations")
	site_associations := make([]interface{}, 0)
	for i := 0; i < count1; i++ {
		sitesCont, err := dataCon.ArrayElement(i, "siteAssociations")
		if err != nil {
			return fmt.Errorf("Unable to parse the site associations list")
		}

		mapSite := make(map[string]interface{})
		mapSite["site_id"] = models.StripQuotes(sitesCont.S("siteId").String())
		mapSite["security_domains"] = sitesCont.S("securityDomains").Data().([]interface{})

		readGcpAccountDataFromSchema(sitesCont, mapSite)
		readAwsAccountDataFromSchema(sitesCont, mapSite)
		readAzureAccountDataFromSchema(sitesCont, mapSite)

		if sitesCont.Exists("cloudAccount") && sitesCont.S("cloudAccount").String() != "{}" {
			setCloudAccountInfo(strings.Split(sitesCont.S("cloudAccount").String(), "/")[2], mapSite)
		}

		site_associations = append(site_associations, mapSite)
	}

	d.Set("site_associations", site_associations)

	count2, _ := dataCon.ArrayCount("userAssociations")
	if err != nil {
		d.Set("user_assocoations", make([]interface{}, 0))
	}

	user_associations := make([]interface{}, 0)
	for i := 0; i < count2; i++ {
		usersCont, err := dataCon.ArrayElement(i, "userAssociations")
		if err != nil {
			return fmt.Errorf("Unable to parse the user associations list")
		}

		mapUser := make(map[string]interface{})
		mapUser["user_id"] = models.StripQuotes(usersCont.S("userId").String())
		user_associations = append(user_associations, mapUser)
	}

	d.Set("user_associations", user_associations)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
