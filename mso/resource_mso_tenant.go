package mso

import (
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOTenant() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOTenantCreate,
		Update: resourceMSOTenantUpdate,
		Read:   resourceMSOTenantRead,
		Delete: resourceMSOTenantDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOTenantImport,
		},

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"display_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"orchestrator_only": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"user_associations": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Optional: true,
				Computed: true,
			},

			"site_associations": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"site_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"security_domains": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
						"vendor": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"aws",
								"azure",
								"gcp",
							}, false),
						},
						"gcp_project_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"gcp_access_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"unmanaged",
								"managed",
							}, false),
						},
						"gcp_name": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"gcp_key_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"gcp_private_key": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"gcp_client_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"gcp_email": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"aws_account_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringLenValidator(12),
						},
						"is_aws_account_trusted": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"aws_access_key_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringLenValidator(20),
						},
						"aws_secret_key": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringLenValidator(40),
						},
						"azure_subscription_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"azure_access_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"credentials",
								"managed",
								"shared",
							}, false),
						},
						"azure_shared_account_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"azure_application_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"azure_client_secret": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"azure_active_directory_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
					},
				},
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func StringLenValidator(lengt int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		v, ok := i.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
			return warnings, errors
		}

		if len(v) != lengt {
			errors = append(errors, fmt.Errorf("expected length of %s to be %d , got %d", k, lengt, len(v)))
		}

		return warnings, errors
	}
}

func readAwsAccountDataFromSchema(sitesCont *container.Container, mapSite map[string]interface{}) {
	awsCont, err := sitesCont.ArrayElement(0, "awsAccount")
	mapSite["aws_account_id"] = ""
	mapSite["aws_access_key_id"] = ""
	mapSite["aws_secret_key"] = ""
	mapSite["is_aws_account_trusted"] = false
	if err == nil {
		mapSite["vendor"] = "aws"
		mapSite["aws_account_id"] = models.StripQuotes(awsCont.S("accountId").String())
		if awsCont.Exists("isTrusted") {
			mapSite["is_aws_account_trusted"] = awsCont.S("isTrusted").Data().(bool)
		}
		accessKey := models.StripQuotes(awsCont.S("accessKeyId").String())
		if accessKey != "{}" {
			mapSite["aws_access_key_id"] = accessKey
		}
		secretKey := models.StripQuotes(awsCont.S("secretKey").String())
		if secretKey != "{}" {
			mapSite["aws_secret_key"] = secretKey
		}
	}
}

func readAzureAccountDataFromSchema(sitesCont *container.Container, mapSite map[string]interface{}) {
	azureCont, err := sitesCont.ArrayElement(0, "azureAccount")
	mapSite["azure_access_type"] = ""
	mapSite["azure_client_secret"] = ""
	mapSite["azure_active_directory_id"] = ""
	mapSite["azure_subscription_id"] = ""
	mapSite["azure_application_id"] = ""
	if err == nil {
		mapSite["vendor"] = "azure"
		mapSite["azure_access_type"] = models.StripQuotes(azureCont.S("accessType").String())
		mapSite["azure_subscription_id"] = models.StripQuotes(azureCont.S("cloudSubscription", "cloudSubscriptionId").String())
		if mapSite["azure_access_type"] == "credentials" {
			mapSite["azure_application_id"] = models.StripQuotes(azureCont.S("cloudSubscription", "cloudApplicationId").String())
			appCont, err := azureCont.ArrayElement(0, "cloudApplication")
			if err == nil {
				mapSite["azure_client_secret"] = models.StripQuotes(appCont.S("secretKey").String())
				mapSite["azure_active_directory_id"] = models.StripQuotes(appCont.S("cloudActiveDirectoryId").String())
			}
		}
	}
}

func readGcpAccountDataFromSchema(sitesCont *container.Container, mapSite map[string]interface{}) {
	gcpCont, err := sitesCont.ArrayElement(0, "gcpAccount")
	mapSite["gcp_project_id"] = ""
	mapSite["gcp_access_type"] = ""
	mapSite["gcp_client_id"] = ""
	mapSite["gcp_email"] = ""
	mapSite["gcp_name"] = ""
	mapSite["gcp_key_id"] = ""
	mapSite["gcp_private_key"] = ""
	if err == nil {
		mapSite["vendor"] = "gcp"
		mapSite["gcp_project_id"] = models.StripQuotes(gcpCont.S("projectID").String())
		mapSite["gcp_access_type"] = models.StripQuotes(gcpCont.S("accessType").String())
		if models.StripQuotes(gcpCont.S("accessType").String()) == "unmanaged" {
			mapSite["gcp_client_id"] = models.StripQuotes(gcpCont.S("cloudCredentials", "clientId").String())
			mapSite["gcp_email"] = models.StripQuotes(gcpCont.S("cloudCredentials", "email").String())
			mapSite["gcp_name"] = models.StripQuotes(gcpCont.S("cloudCredentials", "name").String())
			mapSite["gcp_key_id"] = models.StripQuotes(gcpCont.S("cloudCredentials", "keyId").String())
			mapSite["gcp_private_key"] = models.StripQuotes(gcpCont.S("cloudCredentials", "rsaPrivateKey").String())
		}
	}
}

func setGcpAccountDetails(mapSite, accountDetails map[string]interface{}, tenant string, new bool) {
	gcpAccountMap := make(map[string]interface{})
	gcpAccountMap["vendor"] = "gcp"
	gcpAccountMap["isNew"] = new
	gcpAccountMap["projectID"] = accountDetails["gcp_project_id"]
	gcpAccountMap["securityDomains"] = make([]interface{}, 0)
	gcpAccountMap["accessType"] = accountDetails["gcp_access_type"]
	if accountDetails["gcp_access_type"] == "unmanaged" {
		cloudCredentials := make(map[string]interface{})
		cloudCredentials["clientId"] = accountDetails["gcp_client_id"]
		cloudCredentials["email"] = accountDetails["gcp_email"]
		cloudCredentials["keyId"] = accountDetails["gcp_key_id"]
		cloudCredentials["name"] = accountDetails["gcp_name"]
		cloudCredentials["rsaPrivateKey"] = accountDetails["gcp_private_key"]
		gcpAccountMap["cloudCredentials"] = cloudCredentials
	}
	mapSite["cloudAccount"] = fmt.Sprintf("uni/tn-%s/act-[%s]-vendor-gcp", tenant, accountDetails["gcp_project_id"])
	mapSite["gcpAccount"] = [...]interface{}{gcpAccountMap}
}

func resourceMSOTenantImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] Tenant: Beginning Import")
	msoClient := m.(*client.Client)
	con, err := msoClient.GetViaURL("api/v1/tenants/" + d.Id())
	if err != nil {
		return nil, err
	}
	d.SetId(models.StripQuotes(con.S("id").String()))
	d.Set("name", models.StripQuotes(con.S("name").String()))
	d.Set("display_name", models.StripQuotes(con.S("displayName").String()))
	d.Set("description", models.StripQuotes(con.S("description").String()))
	count1, _ := con.ArrayCount("siteAssociations")
	site_associations := make([]interface{}, 0)
	for i := 0; i < count1; i++ {
		sitesCont, err := con.ArrayElement(i, "siteAssociations")
		if err != nil {
			return nil, fmt.Errorf("Unable to parse the site associations list")
		}
		mapSite := make(map[string]interface{})
		mapSite["site_id"] = models.StripQuotes(sitesCont.S("siteId").String())
		mapSite["security_domains"] = sitesCont.S("securityDomains").Data().([]interface{})

		readGcpAccountDataFromSchema(sitesCont, mapSite)
		readAwsAccountDataFromSchema(sitesCont, mapSite)
		readAzureAccountDataFromSchema(sitesCont, mapSite)

		if sitesCont.Exists("cloudAccount") && sitesCont.S("cloudAccount").String() != "{}" {
			splitStr := strings.Split(sitesCont.S("cloudAccount").String(), "/")
			if len(splitStr) > 2 {
				setCloudAccountInfo(splitStr[2], mapSite)
			}
		}

		site_associations = append(site_associations, mapSite)
	}
	d.Set("site_associations", site_associations)
	count2, _ := con.ArrayCount("userAssociations")
	if err != nil {
		d.Set("user_assocoations", make([]interface{}, 0))
	}
	user_associations := make([]interface{}, 0)
	for i := 0; i < count2; i++ {
		usersCont, err := con.ArrayElement(i, "userAssociations")
		if err != nil {
			return nil, fmt.Errorf("Unable to parse the user associations list")
		}
		mapUser := make(map[string]interface{})
		mapUser["user_id"] = models.StripQuotes(usersCont.S("userId").String())
		user_associations = append(user_associations, mapUser)
	}
	d.Set("user_associations", user_associations)
	log.Printf("[DEBUG] %s: Tenant Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOTenantCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Tenant: Beginning Creation")
	msoClient := m.(*client.Client)
	tenantAttr := models.TenantAttributes{}

	if name, ok := d.GetOk("name"); ok {
		tenantAttr.Name = name.(string)
	}

	if display_name, ok := d.GetOk("display_name"); ok {
		tenantAttr.DisplayName = display_name.(string)
	}

	if description, ok := d.GetOk("description"); ok {
		tenantAttr.Description = description.(string)
	}

	site_associations := make([]interface{}, 0, 1)
	if val, ok := d.GetOk("site_associations"); ok {
		siteList := val.([]interface{})
		for _, val := range siteList {

			mapSite := make(map[string]interface{})
			inner := val.(map[string]interface{})
			if inner["site_id"] != "" {
				mapSite["siteId"] = fmt.Sprintf("%v", inner["site_id"])
			}
			if inner["vendor"] != "" {
				if inner["vendor"] == "gcp" {

					if inner["gcp_project_id"] == "" {
						return fmt.Errorf("gcp_project_id is required with vendor = gcp")
					}

					setGcpAccountDetails(mapSite, inner, tenantAttr.Name, true)

				} else if inner["vendor"] == "aws" {

					awsAccountMap := make(map[string]interface{})

					if inner["aws_account_id"] != "" {
						awsAccountMap["accountId"] = inner["aws_account_id"]
						mapSite["cloudAccount"] = inner["aws_account_id"]
					} else {
						return fmt.Errorf("aws_account_id is required with vendor = aws")
					}

					trusted := inner["is_aws_account_trusted"]
					awsAccountMap["isTrusted"] = trusted.(bool)
					awsAccountMap["vendor"] = "aws"

					if !(trusted.(bool)) {
						if inner["aws_access_key_id"] != "" {
							awsAccountMap["accessKeyId"] = inner["aws_access_key_id"]
						} else {
							return fmt.Errorf("aws_access_key_id is required if the AWS account is not trusted.")
						}

						if inner["aws_secret_key"] != "" {
							awsAccountMap["secretKey"] = inner["aws_secret_key"]
						} else {
							return fmt.Errorf("aws_secret_key is required if the AWS account is not trusted.")
						}
					}

					mapSite["awsAccount"] = [...]interface{}{awsAccountMap}

				} else if inner["vendor"] == "azure" {
					azureAccountMap := make(map[string]interface{})
					var subscriptionId string

					azureAccessType := inner["azure_access_type"].(string)
					applicationId := inner["azure_application_id"].(string)
					clientSecret := inner["azure_client_secret"].(string)
					activeDirectoryId := inner["azure_active_directory_id"].(string)
					sharedAccID := inner["azure_shared_account_id"].(string)

					if azureAccessType == "" {
						azureAccessType = "managed"
					}
					cloudSubMap := make(map[string]interface{})
					if azureAccessType == "managed" {
						if inner["azure_subscription_id"] != "" {
							subscriptionId = inner["azure_subscription_id"].(string)
							mapSite["cloudAccount"] = fmt.Sprintf("uni/tn-%s/act-[%s]-vendor-azure", tenantAttr.Name, subscriptionId)
						} else {
							return fmt.Errorf("azure_subscription_id is required when vendor = azure and azure_access_type = managed or credentials")
						}

						cloudSubMap["cloudSubscriptionId"] = subscriptionId
						azureAccountMap["cloudSubscription"] = cloudSubMap

						azureAccountMap["securityDomains"] = make([]interface{}, 0)
						azureAccountMap["vendor"] = "azure"
						azureAccountMap["accessType"] = azureAccessType

						mapSite["azureAccount"] = [...]interface{}{azureAccountMap}

					} else if azureAccessType == "credentials" {
						if applicationId == "" || clientSecret == "" || activeDirectoryId == "" {
							return fmt.Errorf("azure_application_id, azure_client_secret and azure_active_directory_id are required with azure_access_type = credentials")
						}

						if inner["azure_subscription_id"] != "" {
							subscriptionId = inner["azure_subscription_id"].(string)
							mapSite["cloudAccount"] = fmt.Sprintf("uni/tn-%s/act-[%s]-vendor-azure", tenantAttr.Name, subscriptionId)
						} else {
							return fmt.Errorf("azure_subscription_id is required when vendor = azure and azure_access_type = managed or credentials")
						}

						cloudSubMap["cloudSubscriptionId"] = subscriptionId
						cloudSubMap["cloudApplicationId"] = applicationId
						azureAccountMap["cloudSubscription"] = cloudSubMap
						cloudApplicationMap := make(map[string]interface{})

						cloudApplicationMap["cloudApplicationId"] = applicationId
						cloudApplicationMap["secretKey"] = clientSecret
						cloudApplicationMap["cloudActiveDirectoryId"] = activeDirectoryId
						cloudApplicationMap["cloudCredentialName"] = "cApicApp"

						azureAccountMap["cloudApplication"] = [...]interface{}{cloudApplicationMap}

						activeDirectoryMap := make(map[string]interface{})
						activeDirectoryMap["cloudActiveDirectoryId"] = activeDirectoryId
						activeDirectoryMap["cloudActiveDirectoryName"] = "CiscoINSBUAd"

						azureAccountMap["cloudActiveDirectory"] = [...]interface{}{activeDirectoryMap}

						azureAccountMap["securityDomains"] = make([]interface{}, 0)
						azureAccountMap["vendor"] = "azure"
						azureAccountMap["accessType"] = azureAccessType

						mapSite["azureAccount"] = [...]interface{}{azureAccountMap}

					} else if azureAccessType == "shared" {
						if sharedAccID == "" || inner["site_id"] == "" {
							return fmt.Errorf("azure_shared_account_id and site_id are required with azure_access_type = shared")
						}
						durl := fmt.Sprintf("api/v1/sites/%s/aci/cloud-accounts", inner["site_id"].(string))
						cont, err := msoClient.GetViaURL(durl)
						if err != nil {
							return err
						}

						count, err := cont.ArrayCount("cloudAccounts")
						if err != nil {
							return err
						}

						for i := 0; i < count; i++ {
							dn := models.StripQuotes(cont.S("cloudAccounts").Index(i).S("id").String())
							if dn == sharedAccID {
								mapSite["cloudAccount"] = models.StripQuotes(cont.S("cloudAccounts").Index(i).S("dn").String())
								break
							}
						}
					}
				}
			}

			mapSite["securityDomains"] = make([]interface{}, 0)
			site_associations = append(site_associations, mapSite)
		}
	}
	log.Println("check .... : ", site_associations)
	tenantAttr.Sites = site_associations

	user_associations := make([]interface{}, 0, 1)
	if val, ok := d.GetOk("user_associations"); ok {
		userList := val.(*schema.Set).List()
		for _, val := range userList {

			mapUser := make(map[string]interface{})
			inner := val.(map[string]interface{})
			if inner["user_id"] != "" {
				mapUser["userId"] = fmt.Sprintf("%v", inner["user_id"])
			}
			user_associations = append(user_associations, mapUser)
		}
	}
	tenantAttr.Users = user_associations

	tenantApp := models.NewTenant(tenantAttr)

	cont, err := msoClient.Save("api/v1/tenants", tenantApp)
	if err != nil {
		log.Println(err)
		return err
	}

	id := models.StripQuotes(cont.S("id").String())
	d.SetId(fmt.Sprintf("%v", id))
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceMSOTenantRead(d, m)
}

func resourceMSOTenantUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Tenant: Beginning Update")

	msoClient := m.(*client.Client)

	tenantAttr := models.TenantAttributes{}

	if name, ok := d.GetOk("name"); ok {
		tenantAttr.Name = name.(string)
	}

	if display_name, ok := d.GetOk("display_name"); ok {
		tenantAttr.DisplayName = display_name.(string)
	}

	if description, ok := d.GetOk("description"); ok {
		tenantAttr.Description = description.(string)
	}

	site_associations := make([]interface{}, 0, 1)
	if val, ok := d.GetOk("site_associations"); ok {
		siteList := val.([]interface{})
		for _, val := range siteList {

			mapSite := make(map[string]interface{})
			inner := val.(map[string]interface{})
			if inner["site_id"] != "" {
				mapSite["siteId"] = fmt.Sprintf("%v", inner["site_id"])
			}
			if inner["vendor"] != "" {
				if inner["vendor"] == "gcp" {

					if inner["gcp_project_id"] == "" {
						return fmt.Errorf("gcp_project_id is required with vendor = gcp")
					}

					setGcpAccountDetails(mapSite, inner, tenantAttr.Name, false)

				} else if inner["vendor"] == "aws" {

					awsAccountMap := make(map[string]interface{})

					if inner["aws_account_id"] != "" {
						awsAccountMap["accountId"] = inner["aws_account_id"]
						mapSite["cloudAccount"] = inner["aws_account_id"]
					} else {
						return fmt.Errorf("aws_account_id is required with vendor = aws")
					}

					trusted := inner["is_aws_account_trusted"]
					awsAccountMap["isTrusted"] = trusted.(bool)
					awsAccountMap["vendor"] = "aws"

					if !(trusted.(bool)) {
						if inner["aws_access_key_id"] != "" {
							awsAccountMap["accessKeyId"] = inner["aws_access_key_id"]
						} else {
							return fmt.Errorf("aws_access_key_id is required if the AWS account is not trusted.")
						}

						if inner["aws_secret_key"] != "" {
							awsAccountMap["secretKey"] = inner["aws_secret_key"]
						} else {
							return fmt.Errorf("aws_secret_key is required if the AWS account is not trusted.")
						}
					}

					mapSite["awsAccount"] = [...]interface{}{awsAccountMap}

				} else if inner["vendor"] == "azure" {
					azureAccountMap := make(map[string]interface{})
					var subscriptionId string

					azureAccessType := inner["azure_access_type"].(string)
					applicationId := inner["azure_application_id"].(string)
					clientSecret := inner["azure_client_secret"].(string)
					activeDirectoryId := inner["azure_active_directory_id"].(string)
					sharedAccID := inner["azure_shared_account_id"].(string)

					if azureAccessType == "" {
						azureAccessType = "managed"
					}
					cloudSubMap := make(map[string]interface{})
					if azureAccessType == "managed" {
						if inner["azure_subscription_id"] != "" {
							subscriptionId = inner["azure_subscription_id"].(string)
							mapSite["cloudAccount"] = fmt.Sprintf("uni/tn-%s/act-[%s]-vendor-azure", tenantAttr.Name, subscriptionId)
						} else {
							return fmt.Errorf("azure_subscription_id is required when vendor = azure and azure_access_type = managed or credentials")
						}

						cloudSubMap["cloudSubscriptionId"] = subscriptionId
						azureAccountMap["cloudSubscription"] = cloudSubMap

						azureAccountMap["securityDomains"] = make([]interface{}, 0)
						azureAccountMap["vendor"] = "azure"
						azureAccountMap["accessType"] = azureAccessType

						mapSite["azureAccount"] = [...]interface{}{azureAccountMap}

					} else if azureAccessType == "credentials" {
						if applicationId == "" || clientSecret == "" || activeDirectoryId == "" {
							return fmt.Errorf("azure_application_id, azure_client_secret and azure_active_directory_id are required with azure_access_type = credentials")
						}

						if inner["azure_subscription_id"] != "" {
							subscriptionId = inner["azure_subscription_id"].(string)
							mapSite["cloudAccount"] = fmt.Sprintf("uni/tn-%s/act-[%s]-vendor-azure", tenantAttr.Name, subscriptionId)
						} else {
							return fmt.Errorf("azure_subscription_id is required when vendor = azure and azure_access_type = managed or credentials")
						}

						cloudSubMap["cloudSubscriptionId"] = subscriptionId
						cloudSubMap["cloudApplicationId"] = applicationId
						azureAccountMap["cloudSubscription"] = cloudSubMap
						cloudApplicationMap := make(map[string]interface{})

						cloudApplicationMap["cloudApplicationId"] = applicationId
						cloudApplicationMap["secretKey"] = clientSecret
						cloudApplicationMap["cloudActiveDirectoryId"] = activeDirectoryId
						cloudApplicationMap["cloudCredentialName"] = "cApicApp"

						azureAccountMap["cloudApplication"] = [...]interface{}{cloudApplicationMap}

						activeDirectoryMap := make(map[string]interface{})
						activeDirectoryMap["cloudActiveDirectoryId"] = activeDirectoryId
						activeDirectoryMap["cloudActiveDirectoryName"] = "CiscoINSBUAd"

						azureAccountMap["cloudActiveDirectory"] = [...]interface{}{activeDirectoryMap}

						azureAccountMap["securityDomains"] = make([]interface{}, 0)
						azureAccountMap["vendor"] = "azure"
						azureAccountMap["accessType"] = azureAccessType

						mapSite["azureAccount"] = [...]interface{}{azureAccountMap}

					} else if azureAccessType == "shared" {
						if sharedAccID == "" || inner["site_id"] == "" {
							return fmt.Errorf("azure_shared_account_id and site_id are required with azure_access_type = shared")
						}
						durl := fmt.Sprintf("api/v1/sites/%s/aci/cloud-accounts", inner["site_id"].(string))
						cont, err := msoClient.GetViaURL(durl)
						if err != nil {
							return err
						}

						count, err := cont.ArrayCount("cloudAccounts")
						if err != nil {
							return err
						}

						for i := 0; i < count; i++ {
							dn := models.StripQuotes(cont.S("cloudAccounts").Index(i).S("id").String())
							if dn == sharedAccID {
								mapSite["cloudAccount"] = models.StripQuotes(cont.S("cloudAccounts").Index(i).S("dn").String())
								break
							}
						}
					}
				}
			}
			mapSite["securityDomains"] = make([]interface{}, 0)
			site_associations = append(site_associations, mapSite)
		}
	}
	tenantAttr.Sites = site_associations

	user_associations := make([]interface{}, 0, 1)
	if val, ok := d.GetOk("user_associations"); ok {
		userList := val.(*schema.Set).List()
		for _, val := range userList {

			mapUser := make(map[string]interface{})
			inner := val.(map[string]interface{})
			if inner["user_id"] != "" {
				mapUser["userId"] = fmt.Sprintf("%v", inner["user_id"])
			}
			user_associations = append(user_associations, mapUser)
		}
	}
	tenantAttr.Users = user_associations

	tenantApp := models.NewTenant(tenantAttr)
	cont, err := msoClient.Put(fmt.Sprintf("api/v1/tenants/%s", d.Id()), tenantApp)
	if err != nil {
		return err
	}

	id := models.StripQuotes(cont.S("id").String())
	d.SetId(fmt.Sprintf("%v", id))
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceMSOTenantRead(d, m)
}

func setCloudAccountInfo(accInfo string, mapSite map[string]interface{}) {
	if strings.Contains(accInfo, "azure") {
		mapSite["vendor"] = "azure"
		mapSite["azure_access_type"] = "shared"
		mapSite["azure_shared_account_id"] = accInfo[strings.Index(accInfo, "[")+1 : strings.Index(accInfo, "]")]
	}
}

func resourceMSOTenantRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	msoClient := m.(*client.Client)

	dn := d.Id()

	con, err := msoClient.GetViaURL("api/v1/tenants/" + dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, con, d)
	}

	d.SetId(models.StripQuotes(con.S("id").String()))
	d.Set("name", models.StripQuotes(con.S("name").String()))
	d.Set("display_name", models.StripQuotes(con.S("displayName").String()))
	d.Set("description", models.StripQuotes(con.S("description").String()))

	count1, _ := con.ArrayCount("siteAssociations")
	site_associations := make([]interface{}, 0)
	for i := 0; i < count1; i++ {
		sitesCont, err := con.ArrayElement(i, "siteAssociations")
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
			splitStr := strings.Split(sitesCont.S("cloudAccount").String(), "/")
			if len(splitStr) > 2 {
				setCloudAccountInfo(splitStr[2], mapSite)
			}
		}

		site_associations = append(site_associations, mapSite)
	}

	d.Set("site_associations", site_associations)

	count2, _ := con.ArrayCount("userAssociations")
	if err != nil {
		d.Set("user_assocoations", make([]interface{}, 0))
	}

	user_associations := make([]interface{}, 0)
	for i := 0; i < count2; i++ {
		usersCont, err := con.ArrayElement(i, "userAssociations")
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

func resourceMSOTenantDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	msoClient := m.(*client.Client)
	dn := d.Id()
	orchestratorOnly := d.Get("orchestrator_only").(bool)
	err := msoClient.DeletebyId(fmt.Sprintf("api/v1/tenants/%v?msc-only=%v", dn, orchestratorOnly))
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return nil
}
