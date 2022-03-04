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
							}, false),
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
		awsCount, err := sitesCont.ArrayCount("awsAccount")
		if err == nil {
			if awsCount > 0 {
				awsCont, err := sitesCont.ArrayElement(0, "awsAccount")
				if err == nil {
					mapSite["aws_account_id"] = models.StripQuotes(awsCont.S("accountId").String())
					if awsCont.Exists("isTrusted") && awsCont.S("isTrusted").Data() != nil {
						mapSite["is_aws_account_trusted"] = awsCont.S("isTrusted").Data().(bool)
					}
					mapSite["vendor"] = "aws"
					accessKey := models.StripQuotes(awsCont.S("accessKeyId").String())
					secretKey := models.StripQuotes(awsCont.S("secretKey").String())

					if accessKey != "{}" {
						mapSite["aws_access_key_id"] = accessKey
					}

					if secretKey != "{}" {
						mapSite["aws_secret_key"] = secretKey
					}

				} else {
					log.Printf("Unable to load AWS credentials")
				}
			} else {
				mapSite["aws_account_id"] = ""
				mapSite["aws_access_key_id"] = ""
				mapSite["aws_secret_key"] = ""
				mapSite["is_aws_account_trusted"] = false
			}
		} else {
			log.Printf("Error occurred while loading AWS creds")
			mapSite["aws_account_id"] = ""
			mapSite["aws_access_key_id"] = ""
			mapSite["aws_secret_key"] = ""
			mapSite["is_aws_account_trusted"] = false

		}

		azureCount, err := sitesCont.ArrayCount("azureAccount")
		if err == nil {
			if azureCount > 0 {
				azureCont, err := sitesCont.ArrayElement(0, "azureAccount")
				if err == nil {
					mapSite["vendor"] = "azure"
					mapSite["azure_access_type"] = models.StripQuotes(azureCont.S("accessType").String())
					mapSite["azure_subscription_id"] = models.StripQuotes(azureCont.S("cloudSubscription", "cloudSubscriptionId").String())
					if mapSite["azure_access_type"] == "credentials" {
						mapSite["azure_application_id"] = models.StripQuotes(azureCont.S("cloudSubscription", "cloudApplicationId").String())

						applicationCount, err := azureCont.ArrayCount("cloudApplication")
						if err == nil {
							if applicationCount > 0 {
								appCont, err := azureCont.ArrayElement(0, "cloudApplication")
								if err == nil {
									mapSite["azure_client_secret"] = models.StripQuotes(appCont.S("secretKey").String())
									mapSite["azure_active_directory_id"] = models.StripQuotes(appCont.S("cloudActiveDirectoryId").String())
								} else {
									mapSite["azure_client_secret"] = ""
									mapSite["azure_active_directory_id"] = ""
								}
							} else {
								// Set to empty string
								mapSite["azure_client_secret"] = ""
								mapSite["azure_active_directory_id"] = ""
							}
						} else {
							// Set to empty string
							mapSite["azure_client_secret"] = ""
							mapSite["azure_active_directory_id"] = ""
						}
					}

				} else {
					if sitesCont.Exists("cloudAccount") && sitesCont.S("cloudAccount").String() != "{}" {
						mapSite["azure_access_type"] = "shared"
						cldAcc := strings.Split(models.StripQuotes(sitesCont.S("cloudAccount").String()), "/")
						accInfo := strings.Split(cldAcc[2], "-")

						mapSite["vendor"] = accInfo[3]
						mapSite["azure_shared_account_id"] = (accInfo[1])[1 : len(accInfo[1])-1]

					} else {
						mapSite["azure_access_type"] = ""

					}
					mapSite["azure_client_secret"] = ""
					mapSite["azure_active_directory_id"] = ""
					mapSite["azure_subscription_id"] = ""
					mapSite["azure_application_id"] = ""
				}

			} else {
				log.Printf("Error occurred while loading count for azureAccount.")
				mapSite["azure_client_secret"] = ""
				mapSite["azure_active_directory_id"] = ""
				mapSite["azure_access_type"] = ""
				mapSite["azure_subscription_id"] = ""
				mapSite["azure_application_id"] = ""
				mapSite["azure_shared_account_id"] = ""
			}
		} else {
			log.Printf("Error ocurred while loading azure credentials")
			mapSite["azure_client_secret"] = ""
			mapSite["azure_active_directory_id"] = ""
			mapSite["azure_access_type"] = ""
			mapSite["azure_subscription_id"] = ""
			mapSite["azure_application_id"] = ""
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
				if inner["vendor"] == "aws" {

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
				if inner["vendor"] == "aws" {

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

func resourceMSOTenantRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	msoClient := m.(*client.Client)

	dn := d.Id()

	con, err := msoClient.GetViaURL("api/v1/tenants/" + dn)
	if err != nil {
		return err
	}

	d.SetId(models.StripQuotes(con.S("id").String()))
	d.Set("name", models.StripQuotes(con.S("name").String()))
	d.Set("display_name", models.StripQuotes(con.S("displayName").String()))
	if con.Exists("description") {
		d.Set("description", models.StripQuotes(con.S("description").String()))
	} else {
		d.Set("description", "")
	}

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
		awsCount, err := sitesCont.ArrayCount("awsAccount")
		if err == nil {
			if awsCount > 0 {
				awsCont, err := sitesCont.ArrayElement(0, "awsAccount")
				if err == nil {
					mapSite["aws_account_id"] = models.StripQuotes(awsCont.S("accountId").String())
					if awsCont.Exists("isTrusted") && awsCont.S("isTrusted").Data() != nil {
						mapSite["is_aws_account_trusted"] = awsCont.S("isTrusted").Data().(bool)
					}
					mapSite["vendor"] = "aws"
					accessKey := models.StripQuotes(awsCont.S("accessKeyId").String())
					secretKey := models.StripQuotes(awsCont.S("secretKey").String())

					if accessKey != "{}" {
						mapSite["aws_access_key_id"] = accessKey
					}

					if secretKey != "{}" {
						mapSite["aws_secret_key"] = secretKey
					}

				} else {
					log.Printf("Unable to load AWS credentials")
				}
			} else {
				mapSite["aws_account_id"] = ""
				mapSite["aws_access_key_id"] = ""
				mapSite["aws_secret_key"] = ""
				mapSite["is_aws_account_trusted"] = false
			}
		} else {
			log.Printf("Error occurred while loading AWS creds")
			mapSite["aws_account_id"] = ""
			mapSite["aws_access_key_id"] = ""
			mapSite["aws_secret_key"] = ""
			mapSite["is_aws_account_trusted"] = false

		}

		azureCount, err := sitesCont.ArrayCount("azureAccount")
		if err == nil {
			if azureCount > 0 {
				azureCont, err := sitesCont.ArrayElement(0, "azureAccount")
				if err == nil {
					mapSite["vendor"] = "azure"
					mapSite["azure_access_type"] = models.StripQuotes(azureCont.S("accessType").String())
					mapSite["azure_subscription_id"] = models.StripQuotes(azureCont.S("cloudSubscription", "cloudSubscriptionId").String())
					if mapSite["azure_access_type"] == "credentials" {
						mapSite["azure_application_id"] = models.StripQuotes(azureCont.S("cloudSubscription", "cloudApplicationId").String())

						applicationCount, err := azureCont.ArrayCount("cloudApplication")
						if err == nil {
							if applicationCount > 0 {
								appCont, err := azureCont.ArrayElement(0, "cloudApplication")
								if err == nil {
									mapSite["azure_client_secret"] = models.StripQuotes(appCont.S("secretKey").String())
									mapSite["azure_active_directory_id"] = models.StripQuotes(appCont.S("cloudActiveDirectoryId").String())
								} else {
									mapSite["azure_client_secret"] = ""
									mapSite["azure_active_directory_id"] = ""
								}
							} else {
								// Set to empty string
								mapSite["azure_client_secret"] = ""
								mapSite["azure_active_directory_id"] = ""
							}
						} else {
							// Set to empty string
							mapSite["azure_client_secret"] = ""
							mapSite["azure_active_directory_id"] = ""
						}
					}

				} else {
					if sitesCont.Exists("cloudAccount") && sitesCont.S("cloudAccount").String() != "{}" {
						mapSite["azure_access_type"] = "shared"
						cldAcc := strings.Split(models.StripQuotes(sitesCont.S("cloudAccount").String()), "/")
						accInfo := strings.Split(cldAcc[2], "-")

						mapSite["vendor"] = accInfo[3]
						mapSite["azure_shared_account_id"] = (accInfo[1])[1 : len(accInfo[1])-1]

					} else {
						mapSite["azure_access_type"] = ""

					}
					mapSite["azure_client_secret"] = ""
					mapSite["azure_active_directory_id"] = ""
					mapSite["azure_subscription_id"] = ""
					mapSite["azure_application_id"] = ""
				}

			} else {
				log.Printf("Error occurred while loading count for azureAccount.")
				mapSite["azure_client_secret"] = ""
				mapSite["azure_active_directory_id"] = ""
				mapSite["azure_access_type"] = ""
				mapSite["azure_subscription_id"] = ""
				mapSite["azure_application_id"] = ""
				mapSite["azure_shared_account_id"] = ""
			}
		} else {
			log.Printf("Error ocurred while loading azure credentials")
			mapSite["azure_client_secret"] = ""
			mapSite["azure_active_directory_id"] = ""
			mapSite["azure_access_type"] = ""
			mapSite["azure_subscription_id"] = ""
			mapSite["azure_application_id"] = ""
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
	err := msoClient.DeletebyId(fmt.Sprintf("api/v1/tenants/%v", dn))
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return nil
}
