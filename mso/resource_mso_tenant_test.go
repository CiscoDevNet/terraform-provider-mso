package mso

// Note: User association tests are skipped because the mso_user data source uses
// api/v1/users (or api/v2/users on ND), which no longer works on ND 4.1+ where
// the endpoint changed to /api/v1/infra/aaa/localUsers.
//
// Note: Cloud site association tests (AWS/Azure/GCP) are skipped because they
// require real cloud account credentials and vendor-specific configuration that
// is not available in the standard test environment.

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// msoTenantId is used to capture the tenant ID from the first test step for use in the manual delete/recreate step.
var msoTenantId string

func TestAccMSOTenantResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMsoTenantDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create Tenant") },
				Config:    testAccMSOTenantConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant.tenant", "name", msoTenantName),
					resource.TestCheckResourceAttr("mso_tenant.tenant", "display_name", msoTenantName),
					resource.TestCheckResourceAttr("mso_tenant.tenant", "description", "Terraform test tenant"),
					// Capture the tenant ID for the manual delete/recreate step.
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["mso_tenant.tenant"]
						if !ok {
							return fmt.Errorf("mso_tenant.tenant not found in state")
						}
						msoTenantId = rs.Primary.ID
						return nil
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Tenant basic fields") },
				Config:    testAccMSOTenantConfigUpdateBasicFields(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant.tenant", "name", msoTenantName),
					resource.TestCheckResourceAttr("mso_tenant.tenant", "display_name", msoTenantName),
					resource.TestCheckResourceAttr("mso_tenant.tenant", "description", "Terraform test tenant updated"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove Tenant description") },
				Config:    testAccMSOTenantConfigRemoveDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant.tenant", "name", msoTenantName),
					resource.TestCheckResourceAttr("mso_tenant.tenant", "display_name", msoTenantName),
					resource.TestCheckResourceAttr("mso_tenant.tenant", "description", ""),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Tenant") },
				ResourceName:      "mso_tenant.tenant",
				ImportState:       true,
				ImportStateVerify: true,
				// orchestrator_only is client-side only (controls delete behavior) and is not returned by the API.
				ImportStateVerifyIgnore: []string{"orchestrator_only"},
			},
			{
				PreConfig: func() { fmt.Println("Test: Add site association") },
				Config:    testAccMSOTenantConfigAddSiteAssociation(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant.tenant", "name", msoTenantName),
					resource.TestCheckResourceAttr("mso_tenant.tenant", "site_associations.#", "1"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add extra site association") },
				Config:    testAccMSOTenantConfigAddExtraSiteAssociation(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant.tenant", "name", msoTenantName),
					resource.TestCheckResourceAttr("mso_tenant.tenant", "site_associations.#", "2"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove extra site association") },
				Config:    testAccMSOTenantConfigRemoveSiteAssociation(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant.tenant", "name", msoTenantName),
					resource.TestCheckResourceAttr("mso_tenant.tenant", "site_associations.#", "1"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Tenant with site associations") },
				ResourceName:      "mso_tenant.tenant",
				ImportState:       true,
				ImportStateVerify: true,
				// orchestrator_only is client-side only (controls delete behavior) and is not returned by the API.
				ImportStateVerifyIgnore: []string{"orchestrator_only"},
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Verify recreation of manually deleted Tenant")
					client := testAccProvider.Meta().(*client.Client)
					err := client.DeletebyId("api/v1/tenants/" + msoTenantId)
					if err != nil {
						panic(fmt.Sprintf("Failed to delete tenant %s via API: %s", msoTenantId, err))
					}
				},
				Config: testAccMSOTenantConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant.tenant", "name", msoTenantName),
					resource.TestCheckResourceAttr("mso_tenant.tenant", "display_name", msoTenantName),
					resource.TestCheckResourceAttr("mso_tenant.tenant", "description", "Terraform test tenant"),
				),
			},
		},
	})
}

func testAccMSOTenantConfigCreate() string {
	return fmt.Sprintf(`
	resource "mso_tenant" "tenant" {
		name         = "%s"
		display_name = "%s"
		description  = "Terraform test tenant"
	}`, msoTenantName, msoTenantName)
}

func testAccMSOTenantConfigUpdateBasicFields() string {
	return fmt.Sprintf(`
	resource "mso_tenant" "tenant" {
		name         = "%s"
		display_name = "%s"
		description  = "Terraform test tenant updated"
	}`, msoTenantName, msoTenantName)
}

func testAccMSOTenantConfigRemoveDescription() string {
	return fmt.Sprintf(`
	resource "mso_tenant" "tenant" {
		name         = "%s"
		display_name = "%s"
	}`, msoTenantName, msoTenantName)
}

func testAccMSOTenantConfigAddSiteAssociation() string {
	return fmt.Sprintf(`%s
	resource "mso_tenant" "tenant" {
		name         = "%s"
		display_name = "%s"
		site_associations {
			site_id = data.mso_site.%s.id
		}
	}`, testSiteConfigAnsibleTest(), msoTenantName, msoTenantName, msoTemplateSiteName1)
}

func testAccMSOTenantConfigAddExtraSiteAssociation() string {
	return fmt.Sprintf(`%s%s
	resource "mso_tenant" "tenant" {
		name         = "%s"
		display_name = "%s"
		site_associations {
			site_id = data.mso_site.%s.id
		}
		site_associations {
			site_id = data.mso_site.%s.id
		}
	}`, testSiteConfigAnsibleTest(), testSiteConfigAnsibleTest2(), msoTenantName, msoTenantName, msoTemplateSiteName1, msoTemplateSiteName2)
}

func testAccMSOTenantConfigRemoveSiteAssociation() string {
	return fmt.Sprintf(`%s
	resource "mso_tenant" "tenant" {
		name         = "%s"
		display_name = "%s"
		site_associations {
			site_id = data.mso_site.%s.id
		}
	}`, testSiteConfigAnsibleTest(), msoTenantName, msoTenantName, msoTemplateSiteName1)
}

// testAccCheckMsoTenantDestroy verifies that the tenant is deleted from MSO.
// The generic testCheckResourceDestroyPolicy helpers cannot be used here because
// they query the policy/template API (api/v1/templates/objects), whereas tenants
// use a separate API endpoint (api/v1/tenants).
func testAccCheckMsoTenantDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_tenant" {
			_, err := client.GetViaURL("api/v1/tenants/" + rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Tenant still exists")
			}
		} else {
			continue
		}
	}
	return nil
}

// TestAccMSOTenantResourceDisplayNameDefault validates that display_name is
// optional and defaults to the value of name on create, and that omitting it
// on a subsequent update keeps display_name equal to name (the Computed
// behavior retains the value across updates).
func TestAccMSOTenantResourceDisplayNameDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMsoTenantDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: Create Tenant without display_name (defaults to name)")
				},
				Config: testAccMSOTenantConfigNoDisplayName(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant.tenant_default", "name", msoTenantName2),
					resource.TestCheckResourceAttr("mso_tenant.tenant_default", "display_name", msoTenantName2),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update Tenant without display_name (Computed keeps prior value)")
				},
				Config: testAccMSOTenantConfigNoDisplayNameWithDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant.tenant_default", "name", msoTenantName2),
					resource.TestCheckResourceAttr("mso_tenant.tenant_default", "display_name", msoTenantName2),
					resource.TestCheckResourceAttr("mso_tenant.tenant_default", "description", "display_name default test"),
				),
			},
		},
	})
}

func testAccMSOTenantConfigNoDisplayName() string {
	return fmt.Sprintf(`
	resource "mso_tenant" "tenant_default" {
		name = "%s"
	}`, msoTenantName2)
}

func testAccMSOTenantConfigNoDisplayNameWithDescription() string {
	return fmt.Sprintf(`
	resource "mso_tenant" "tenant_default" {
		name        = "%s"
		description = "display_name default test"
	}`, msoTenantName2)
}
