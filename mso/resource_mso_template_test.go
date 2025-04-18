package mso

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const msoTemplateTenantName = "tf_test_mso_template_tenant"

var msoTemplateId string

func TestAccMSOTemplateResourceTenant(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Duplicate sites provided in Tenant Template configuration (error)") },
				Config:      testAccMSOTemplateResourceTenantErrorDuplicateSitesConfig(),
				ExpectError: regexp.MustCompile(`Duplication found in the sites list`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: No tenant provided in Tenant Template configuration (error)") },
				Config:      testAccMSOTemplateResourceTenanErrorNoTenantConfig(),
				ExpectError: regexp.MustCompile(`Tenant is required for template of type tenant.`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create Tenant Template without sites") },
				Config:    testAccMSOTemplateResourceTenantConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccMSOTemplateState(
						"mso_template.template_tenant",
						&TemplateTest{
							TemplateName: "test_template_tenant",
							TemplateType: "tenant",
							Tenant:       msoTemplateTenantName,
							Sites:        []string{},
						},
						true,
					),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import the Tenant Template with no sites configuration") },
				ResourceName:      "mso_template.template_tenant",
				ImportState:       true,
				ImportStateId:     msoTemplateId,
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Import the Tenant Template with no sites configuration with wrong ID (error)")
				},
				ResourceName:      "mso_template.template_tenant",
				ImportState:       true,
				ImportStateId:     "non_existing_template_id",
				ImportStateVerify: true,
				ExpectError:       regexp.MustCompile("Template ID non_existing_template_id invalid"),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update the Tenant Template with 1 site") },
				Config:    testAccMSOTemplateResourceTenanSiteAnsibleTestConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccMSOTemplateState(
						"mso_template.template_tenant",
						&TemplateTest{
							TemplateName: "test_template_tenant",
							TemplateType: "tenant",
							Tenant:       msoTemplateTenantName,
							Sites:        []string{msoTemplateSiteName1},
						},
						false,
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update the Tenant Template with 2 sites") },
				Config:    testAccMSOTemplateResourceTenanTwoSitesConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccMSOTemplateState(
						"mso_template.template_tenant",
						&TemplateTest{
							TemplateName: "test_template_tenant",
							TemplateType: "tenant",
							Tenant:       msoTemplateTenantName,
							Sites:        []string{msoTemplateSiteName1, msoTemplateSiteName2},
						},
						false,
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update the Tenant Template with reverse order of sites") },
				Config:    testAccMSOTemplateResourceTenanTwoSitesReversedConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccMSOTemplateState(
						"mso_template.template_tenant",
						&TemplateTest{
							TemplateName: "test_template_tenant",
							TemplateType: "tenant",
							Tenant:       msoTemplateTenantName,
							Sites:        []string{msoTemplateSiteName1, msoTemplateSiteName2},
						},
						false,
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update the Tenant Template with removal of site 2") },
				Config:    testAccMSOTemplateResourceTenanSiteAnsibleTest2Config(),
				Check: resource.ComposeTestCheckFunc(
					testAccMSOTemplateState(
						"mso_template.template_tenant",
						&TemplateTest{
							TemplateName: "test_template_tenant",
							TemplateType: "tenant",
							Tenant:       msoTemplateTenantName,
							Sites:        []string{msoTemplateSiteName2},
						},
						false,
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update the Tenant Template with change of site 2 to site 1") },
				Config:    testAccMSOTemplateResourceTenanSiteAnsibleTestConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccMSOTemplateState(
						"mso_template.template_tenant",
						&TemplateTest{
							TemplateName: "test_template_tenant",
							TemplateType: "tenant",
							Tenant:       msoTemplateTenantName,
							Sites:        []string{msoTemplateSiteName1},
						},
						false,
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update the Tenant Template with removal of sites configuration") },
				Config:    testAccMSOTemplateResourceTenantNoSitesConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccMSOTemplateState(
						"mso_template.template_tenant",
						&TemplateTest{
							TemplateName: "test_template_tenant",
							TemplateType: "tenant",
							Tenant:       msoTemplateTenantName,
							Sites:        []string{},
						},
						false,
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update the Tenant Template name") },
				Config:    testAccMSOTemplateResourceTenantNameChangeConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccMSOTemplateState(
						"mso_template.template_tenant",
						&TemplateTest{
							TemplateName: "test_template_tenant_changed",
							TemplateType: "tenant",
							Tenant:       msoTemplateTenantName,
							Sites:        []string{},
						},
						false,
					),
				),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Update the Tenant Template with duplicate sites (error)") },
				Config:      testAccMSOTemplateResourceTenantErrorDuplicateSitesConfig(),
				ExpectError: regexp.MustCompile(`Duplication found in the sites list`),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update the Tenant Template after manual removal from MSO")
					msoClient := testAccProvider.Meta().(*client.Client)
					err := msoClient.DeletebyId(fmt.Sprintf("api/v1/templates/%s", msoTemplateId))
					if err != nil {
						t.Fatalf("Failed to manually delete template '%s': %v", msoTemplateId, err)
					}
				},
				Config: testAccMSOTemplateResourceTenantNameChangeConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccMSOTemplateState(
						"mso_template.template_tenant",
						&TemplateTest{
							TemplateName: "test_template_tenant_changed",
							TemplateType: "tenant",
							Tenant:       msoTemplateTenantName,
							Sites:        []string{},
						},
						false,
					),
				),
			},
		},
	})
}

type TemplateTest struct {
	TemplateName string   `json:",omitempty"`
	TemplateType string   `json:",omitempty"`
	Tenant       string   `json:",omitempty"`
	Sites        []string `json:",omitempty"`
}

func testAccMSOTemplateState(resourceName string, stateTemplate *TemplateTest, setmsoTemplateId bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rootModule, err := s.RootModule().Resources[resourceName]
		if !err {
			return fmt.Errorf("%v", err)
		}

		if rootModule.Primary.ID == "" {
			return fmt.Errorf("No ID is set for the template")
		}

		// Set the ID for the template to global variable only when called from specific resource test
		// This is to avoid setting the ID issues when data source test is called first
		if setmsoTemplateId {
			msoTemplateId = rootModule.Primary.ID
		}

		if rootModule.Primary.Attributes["tenant_id"] == "" && stateTemplate.Tenant != "" {
			return fmt.Errorf("No tenant ID is set for the template")
		} else if stateTemplate.Tenant != "" {
			tenantState, err := s.RootModule().Resources[fmt.Sprintf("mso_tenant.%s", stateTemplate.Tenant)]
			if !err {
				return fmt.Errorf("Tenant %s not found in state", stateTemplate.Tenant)
			}
			if tenantState.Primary.Attributes["display_name"] != stateTemplate.Tenant {
				return fmt.Errorf("Tenant display name does not match, expected: %s, got: %s", stateTemplate.Tenant, tenantState.Primary.Attributes["display_name"])
			}
		}

		if rootModule.Primary.Attributes["template_name"] != stateTemplate.TemplateName {
			return fmt.Errorf("Template name does not match, expected: %s, got: %s", stateTemplate.TemplateName, rootModule.Primary.Attributes["template_name"])
		}

		if rootModule.Primary.Attributes["template_type"] != stateTemplate.TemplateType {
			return fmt.Errorf("Template type does not match, expected: %s, got: %s", stateTemplate.TemplateType, rootModule.Primary.Attributes["template_type"])
		}

		if sites, ok := rootModule.Primary.Attributes["sites.#"]; ok {
			if siteAmount, e := strconv.Atoi(sites); e != nil {
				return fmt.Errorf("Could not convert sites amount to integer")
			} else if siteAmount != len(stateTemplate.Sites) {
				return fmt.Errorf("Amount of sites do not match, expected: %d, got: %d", len(stateTemplate.Sites), len(rootModule.Primary.Attributes["sites.#"]))
			}

			for _, site := range stateTemplate.Sites {
				siteState, err := s.RootModule().Resources[fmt.Sprintf("data.mso_site.%s", site)]
				if !err {
					return fmt.Errorf("Site %s not found in state", site)
				}
				if siteState.Primary.Attributes["name"] != site {
					return fmt.Errorf("Site display name does not match, expected: %s, got: %s", site, siteState.Primary.Attributes["display_name"])
				}
			}
		} else {
			if len(stateTemplate.Sites) != 0 {
				return fmt.Errorf("Amount of sites do not match, expected: %d, got: 0", len(stateTemplate.Sites))
			}
		}

		return nil
	}
}

func testAccTenantConfig() string {
	return fmt.Sprintf(`
	%s%s
	resource "mso_tenant" "%s" {
		name = "%s"
		display_name = "%s"
		site_associations { 
			site_id = data.mso_site.%s.id 
		}
		site_associations { 
			site_id = data.mso_site.%s.id 
		}
	}
	`, testSiteConfigAnsibleTest(), testSiteConfigAnsibleTest2(), msoTemplateTenantName, msoTemplateTenantName, msoTemplateTenantName, msoTemplateSiteName1, msoTemplateSiteName2)
}

func testAccMSOTemplateResourceTenantConfig() string {
	return fmt.Sprintf(`%s
	resource "mso_template" "template_tenant" {
		template_name = "test_template_tenant"
		template_type = "tenant"
		tenant_id = mso_tenant.%s.id
	}
	`, testAccTenantConfig(), msoTemplateTenantName)
}

func testAccMSOTemplateResourceTenantNameChangeConfig() string {
	return fmt.Sprintf(`%s
	resource "mso_template" "template_tenant" {
		template_name = "test_template_tenant_changed"
		template_type = "tenant"
		tenant_id = mso_tenant.%s.id
	}
	`, testAccTenantConfig(), msoTemplateTenantName)
}

func testAccMSOTemplateResourceTenantNoSitesConfig() string {
	return fmt.Sprintf(`%s
	resource "mso_template" "template_tenant" {
		template_name = "test_template_tenant"
		template_type = "tenant"
		tenant_id = mso_tenant.%s.id
		sites = []
	}
	`, testAccTenantConfig(), msoTemplateTenantName)
}

func testAccMSOTemplateResourceTenanSiteAnsibleTestConfig() string {
	return fmt.Sprintf(`%s
	resource "mso_template" "template_tenant" {
		template_name = "test_template_tenant"
		template_type = "tenant"
		tenant_id = mso_tenant.%s.id
		sites = [data.mso_site.%s.id]
	}
	`, testAccTenantConfig(), msoTemplateTenantName, msoTemplateSiteName1)
}

func testAccMSOTemplateResourceTenanSiteAnsibleTest2Config() string {
	return fmt.Sprintf(`%s
	resource "mso_template" "template_tenant" {
		template_name = "test_template_tenant"
		template_type = "tenant"
		tenant_id = mso_tenant.%s.id
		sites = [data.mso_site.%s.id]
	}
	`, testAccTenantConfig(), msoTemplateTenantName, msoTemplateSiteName2)
}

func testAccMSOTemplateResourceTenanTwoSitesConfig() string {
	return fmt.Sprintf(`%s
	resource "mso_template" "template_tenant" {
		template_name = "test_template_tenant"
		template_type = "tenant"
		tenant_id = mso_tenant.%s.id
		sites = [data.mso_site.%s.id, data.mso_site.%s.id]
	}
	`, testAccTenantConfig(), msoTemplateTenantName, msoTemplateSiteName1, msoTemplateSiteName2)
}

func testAccMSOTemplateResourceTenanTwoSitesReversedConfig() string {
	return fmt.Sprintf(`%s
	resource "mso_template" "template_tenant" {
		template_name = "test_template_tenant"
		template_type = "tenant"
		tenant_id = mso_tenant.%s.id
		sites = [data.mso_site.%s.id, data.mso_site.%s.id]
	}
	`, testAccTenantConfig(), msoTemplateTenantName, msoTemplateSiteName2, msoTemplateSiteName1)
}

func testAccMSOTemplateResourceTenantErrorDuplicateSitesConfig() string {
	return fmt.Sprintf(`%s
	resource "mso_template" "template_tenant" {
		template_name = "test_template_tenant"
		template_type = "tenant"
		tenant_id = mso_tenant.%s.id
		sites = [data.mso_site.%s.id, data.mso_site.%s.id]
	}
	`, testAccTenantConfig(), msoTemplateTenantName, msoTemplateSiteName1, msoTemplateSiteName1)
}

func testAccMSOTemplateResourceTenanErrorNoTenantConfig() string {
	return fmt.Sprintf(`%s%s
	resource "mso_template" "template_tenant" {
		template_name = "test_template_tenant"
		template_type = "tenant"
		sites = [data.mso_site.%s.id, data.mso_site.%s.id]
	}
	`, testSiteConfigAnsibleTest(), testSiteConfigAnsibleTest2(), msoTemplateSiteName1, msoTemplateSiteName2)
}

func TestAccMSOTemplateResourceL3out(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: No tenant provided in L3out Template configuration (error)") },
				Config:      testAccMSOTemplateResourceL3outErrorNoTenantConfig(),
				ExpectError: regexp.MustCompile(`Tenant is required for template of type l3out.`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: No sites provided in L3out Template configuration (error)") },
				Config:      testAccMSOTemplateResourceL3outErrorNoSitesConfig(),
				ExpectError: regexp.MustCompile(`Site is required for template of type l3out.`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Two sites provided in L3out Template configuration (error)") },
				Config:      testAccMSOTemplateResourceL3outErrorTwositesConfig(),
				ExpectError: regexp.MustCompile(`Only one site is allowed for template of type l3out.`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create L3out Template with 1 site") },
				Config:    testAccMSOTemplateResourceL3outConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccMSOTemplateState(
						"mso_template.template_l3out",
						&TemplateTest{
							TemplateName: "test_template_l3out",
							TemplateType: "l3out",
							Tenant:       msoTemplateTenantName,
							Sites:        []string{msoTemplateSiteName1},
						},
						false,
					),
				),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Update the L3out Template with change of site 1 to site 2 (error)") },
				Config:      testAccMSOTemplateResourceL3outErrorChangeSiteConfig(),
				ExpectError: regexp.MustCompile(`Cannot change site for template of type l3out.`),
			},
		},
	})
}

func testAccMSOTemplateResourceL3outConfig() string {
	return fmt.Sprintf(`%s
	resource "mso_template" "template_l3out" {
		template_name = "test_template_l3out"
		template_type = "l3out"
		tenant_id = mso_tenant.%s.id
		sites = [data.mso_site.%s.id]
	}
	`, testAccTenantConfig(), msoTemplateTenantName, msoTemplateSiteName1)
}

func testAccMSOTemplateResourceL3outErrorNoTenantConfig() string {
	return fmt.Sprintf(`%s
	resource "mso_template" "template_l3out" {
		template_name = "test_template_l3out"
		template_type = "l3out"
		sites = [data.mso_site.%s.id]
	}
	`, testSiteConfigAnsibleTest(), msoTemplateSiteName1)
}

func testAccMSOTemplateResourceL3outErrorNoSitesConfig() string {
	return fmt.Sprintf(`%s
	resource "mso_template" "template_l3out" {
		template_name = "test_template_l3out"
		template_type = "l3out"
		tenant_id = mso_tenant.%s.id
	}
	`, testAccTenantConfig(), msoTemplateTenantName)
}

func testAccMSOTemplateResourceL3outErrorTwositesConfig() string {
	return fmt.Sprintf(`%s
	resource "mso_template" "template_l3out" {
		template_name = "test_template_l3out"
		template_type = "l3out"
		tenant_id = mso_tenant.%s.id
		sites = [data.mso_site.%s.id, data.mso_site.%s.id]
	}
	`, testAccTenantConfig(), msoTemplateTenantName, msoTemplateSiteName1, msoTemplateSiteName2)
}

func testAccMSOTemplateResourceL3outErrorChangeSiteConfig() string {
	return fmt.Sprintf(`%s
	resource "mso_template" "template_l3out" {
		template_name = "test_template_l3out"
		template_type = "l3out"
		tenant_id = mso_tenant.%s.id
		sites = [data.mso_site.%s.id]
	}
	`, testAccTenantConfig(), msoTemplateTenantName, msoTemplateSiteName2)
}

func TestAccMSOTemplateResourceFabricPolicy(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: No tenant provided in Fabric Policy Template configuration (error)") },
				Config:      testAccMSOTemplateResourceFabricPolicyErrorTenantConfig(),
				ExpectError: regexp.MustCompile(`Tenant cannot be attached to template of type fabric_policy.`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create Fabric Policy Template without sites") },
				Config:    testAccMSOTemplateResourceFabricPolicyConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccMSOTemplateState(
						"mso_template.template_fabric_policy",
						&TemplateTest{
							TemplateName: "test_template_fabric_policy",
							TemplateType: "fabric_policy",
							Sites:        []string{},
						},
						false,
					),
				),
			},
		},
	})
}

func testAccMSOTemplateResourceFabricPolicyConfig() string {
	return fmt.Sprintf(`
	resource "mso_template" "template_fabric_policy" {
		template_name = "test_template_fabric_policy"
		template_type = "fabric_policy"
	}
	`)
}

func testAccMSOTemplateResourceFabricPolicyErrorTenantConfig() string {
	return fmt.Sprintf(`%s
	resource "mso_template" "template_fabric_policy" {
		template_name = "test_template_fabric_policy"
		template_type = "fabric_policy"
		tenant_id = mso_tenant.%s.id
	}
	`, testAccTenantConfig(), msoTemplateTenantName)
}

func TestAccMSOTemplateResourceFabricResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Tenant provided in Fabric Resource Template configuration (error)") },
				Config:      testAccMSOTemplateResourceFabricResourceErrorTenantConfig(),
				ExpectError: regexp.MustCompile(`Tenant cannot be attached to template of type fabric_resource.`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create Fabric Resource Template without sites") },
				Config:    testAccMSOTemplateResourceFabricResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccMSOTemplateState(
						"mso_template.template_fabric_resource",
						&TemplateTest{
							TemplateName: "test_template_fabric_resource",
							TemplateType: "fabric_resource",
							Sites:        []string{},
						},
						false,
					),
				),
			},
		},
	})
}

func testAccMSOTemplateResourceFabricResourceConfig() string {
	return fmt.Sprintf(`
	resource "mso_template" "template_fabric_resource" {
		template_name = "test_template_fabric_resource"
		template_type = "fabric_resource"
	}
	`)
}

func testAccMSOTemplateResourceFabricResourceErrorTenantConfig() string {
	return fmt.Sprintf(`%s
	resource "mso_template" "template_fabric_resource" {
		template_name = "test_template_fabric_resource"
		template_type = "fabric_resource"
		tenant_id = mso_tenant.%s.id
	}
	`, testAccTenantConfig(), msoTemplateTenantName)
}

func TestAccMSOTemplateResourceMonitoringTenant(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: No tenant provided in Monitoring Tenant Template configuration (error)") },
				Config:      testAccMSOTemplateResourceMonitoringTenantErrorNoTenantConfig(),
				ExpectError: regexp.MustCompile(`Tenant is required for template of type monitoring_tenant.`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: No site provided in Monitoring Tenant Template configuration (error)") },
				Config:      testAccMSOTemplateResourceMonitoringTenantErrorNoSiteConfig(),
				ExpectError: regexp.MustCompile(`Site is required for template of type monitoring_tenant.`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Two sites provided in Monitoring Tenant Template configuration (error)") },
				Config:      testAccMSOTemplateResourceMonitoringTenantErrorTwoSitesConfig(),
				ExpectError: regexp.MustCompile(`Only one site is allowed for template of type monitoring_tenant.`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create Monitoring Tenant Template with 1 site") },
				Config:    testAccMSOTemplateResourceMonitoringTenantConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccMSOTemplateState(
						"mso_template.template_monitoring_tenant",
						&TemplateTest{
							TemplateName: "test_template_monitoring_tenant",
							TemplateType: "monitoring_tenant",
							Sites:        []string{msoTemplateSiteName1},
						},
						false,
					),
				),
			},
		},
	})
}

func testAccMSOTemplateResourceMonitoringTenantErrorNoTenantConfig() string {
	return fmt.Sprintf(`%s
	resource "mso_template" "template_monitoring_tenant" {
		template_name = "test_template_monitoring_tenant"
		template_type = "monitoring_tenant"
		sites = [data.mso_site.%s.id]
	}
	`, testAccTenantConfig(), msoTemplateSiteName1)
}

func testAccMSOTemplateResourceMonitoringTenantErrorNoSiteConfig() string {
	return fmt.Sprintf(`%s
	resource "mso_template" "template_monitoring_tenant" {
		template_name = "test_template_monitoring_tenant"
		template_type = "monitoring_tenant"
		tenant_id = mso_tenant.%s.id
	}
	`, testAccTenantConfig(), msoTemplateTenantName)
}

func testAccMSOTemplateResourceMonitoringTenantErrorTwoSitesConfig() string {
	return fmt.Sprintf(`%s
	resource "mso_template" "template_monitoring_tenant" {
		template_name = "test_template_monitoring_tenant"
		template_type = "monitoring_tenant"
		tenant_id = mso_tenant.%s.id
		sites = [data.mso_site.%s.id, data.mso_site.%s.id]
	}
	`, testAccTenantConfig(), msoTemplateTenantName, msoTemplateSiteName1, msoTemplateSiteName2)
}

func testAccMSOTemplateResourceMonitoringTenantConfig() string {
	return fmt.Sprintf(`%s
	resource "mso_template" "template_monitoring_tenant" {
		template_name = "test_template_monitoring_tenant"
		template_type = "monitoring_tenant"
		tenant_id = mso_tenant.%s.id
		sites = [data.mso_site.%s.id]
	}
	`, testAccTenantConfig(), msoTemplateTenantName, msoTemplateSiteName1)
}

func TestAccMSOTemplateResourceMonitoringAccess(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Tenant provided in Monitoring Access Template configuration (error)") },
				Config:      testAccMSOTemplateResourceMonitoringAccessErrorTenantConfig(),
				ExpectError: regexp.MustCompile(`Tenant cannot be attached to template of type monitoring_access.`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: No site provided in Monitoring Access Template configuration (error)") },
				Config:      testAccMSOTemplateResourceMonitoringAccessErrorNoSiteConfig(),
				ExpectError: regexp.MustCompile(`Site is required for template of type monitoring_access.`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Two sites provided in Monitoring Access Template configuration (error)") },
				Config:      testAccMSOTemplateResourceMonitoringAccessErrorTwoSitesConfig(),
				ExpectError: regexp.MustCompile(`Only one site is allowed for template of type monitoring_access.`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create Monitoring Access Template with 1 site") },
				Config:    testAccMSOTemplateResourceMonitoringAccessConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccMSOTemplateState(
						"mso_template.template_monitoring_access",
						&TemplateTest{
							TemplateName: "test_template_monitoring_access",
							TemplateType: "monitoring_access",
							Sites:        []string{msoTemplateSiteName1},
						},
						false,
					),
				),
			},
		},
	})
}

func testAccMSOTemplateResourceMonitoringAccessErrorTenantConfig() string {
	return fmt.Sprintf(`%s
	resource "mso_template" "template_monitoring_access" {
		template_name = "test_template_monitoring_access"
		template_type = "monitoring_access"
		tenant_id = mso_tenant.%s.id
		sites = [data.mso_site.%s.id]
	}
	`, testAccTenantConfig(), msoTemplateTenantName, msoTemplateSiteName1)
}

func testAccMSOTemplateResourceMonitoringAccessErrorNoSiteConfig() string {
	return fmt.Sprintf(`%s
	resource "mso_template" "template_monitoring_access" {
		template_name = "test_template_monitoring_access"
		template_type = "monitoring_access"
	}
	`, testAccTenantConfig())
}

func testAccMSOTemplateResourceMonitoringAccessErrorTwoSitesConfig() string {
	return fmt.Sprintf(`%s
	resource "mso_template" "template_monitoring_access" {
		template_name = "test_template_monitoring_access"
		template_type = "monitoring_access"
		sites = [data.mso_site.%s.id, data.mso_site.%s.id]
	}
	`, testAccTenantConfig(), msoTemplateSiteName1, msoTemplateSiteName2)
}

func testAccMSOTemplateResourceMonitoringAccessConfig() string {
	return fmt.Sprintf(`%s
	resource "mso_template" "template_monitoring_access" {
		template_name = "test_template_monitoring_access"
		template_type = "monitoring_access"
		sites = [data.mso_site.%s.id]
	}
	`, testAccTenantConfig(), msoTemplateSiteName1)
}

func TestAccMSOTemplateResourceServiceDevice(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: No tenant provided in Service Device Template configuration (error)") },
				Config:      testAccMSOTemplateResourceServiceDeviceErrorNoTenantConfig(),
				ExpectError: regexp.MustCompile(`Tenant is required for template of type service_device.`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create Service Device Template with 2 sites") },
				Config:    testAccMSOTemplateResourceServiceDeviceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccMSOTemplateState(
						"mso_template.template_service_device",
						&TemplateTest{
							TemplateName: "test_template_service_device",
							TemplateType: "service_device",
							Sites:        []string{msoTemplateSiteName1, msoTemplateSiteName2},
						},
						false,
					),
				),
			},
		},
	})
}

func testAccMSOTemplateResourceServiceDeviceErrorNoTenantConfig() string {
	return fmt.Sprintf(`%s
	resource "mso_template" "template_service_device" {
		template_name = "test_template_service_device"
		template_type = "service_device"
	}
	`, testAccTenantConfig())
}

func testAccMSOTemplateResourceServiceDeviceConfig() string {
	return fmt.Sprintf(`%s
	resource "mso_template" "template_service_device" {
		template_name = "test_template_service_device"
		template_type = "service_device"
		tenant_id = mso_tenant.%s.id
		sites = [data.mso_site.%s.id, data.mso_site.%s.id]
	}
	`, testAccTenantConfig(), msoTemplateTenantName, msoTemplateSiteName1, msoTemplateSiteName2)
}
