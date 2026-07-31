package mso

// Note: user_associations is only round-trip tested on platforms where the
// legacy tenant API preserves explicit associations. ND 4.2+ back-fills
// Tenant-domain users into userAssociations and rejects their removal, so that
// behavior is intentionally skipped here.
//
// Note: Cloud site association tests (AWS/Azure/GCP) are skipped because they
// require real cloud account credentials and vendor-specific configuration that
// is not available in the standard test environment.

import (
	"fmt"
	"os"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// msoTenantId is used to capture the tenant ID from the first test step for use in the manual delete/recreate step.
var msoTenantId string
var msoTenantNameUserAssociations = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

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

func TestAccMSOTenantResourceUserAssociations(t *testing.T) {
	resourceRef := "mso_tenant.tenant_user_associations"

	msoClient := testAccPreCheck(t)
	userID, err := testAccMSOTenantLookupUserID(msoClient, os.Getenv("MSO_USERNAME"))
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccVersionLessThanCheck(t, "5.2.0.0") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMsoTenantDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create Tenant with user association") },
				Config:    testAccMSOTenantConfigUserAssociations("Terraform test tenant user association", userID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", msoTenantNameUserAssociations),
					resource.TestCheckResourceAttr(resourceRef, "display_name", msoTenantNameUserAssociations),
					resource.TestCheckResourceAttr(resourceRef, "description", "Terraform test tenant user association"),
					CustomTestCheckTypeSetElemAttrs(resourceRef, "user_associations", map[string]string{
						"user_id": userID,
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Tenant with same user association") },
				Config:    testAccMSOTenantConfigUserAssociations("Terraform test tenant user association updated", userID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "description", "Terraform test tenant user association updated"),
					CustomTestCheckTypeSetElemAttrs(resourceRef, "user_associations", map[string]string{
						"user_id": userID,
					}),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Tenant with user association") },
				ResourceName:      resourceRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"orchestrator_only",
				},
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

func testAccMSOTenantConfigUserAssociations(description, userID string) string {
	return fmt.Sprintf(`
	resource "mso_tenant" "tenant_user_associations" {
		name         = "%s"
		display_name = "%s"
		description  = "%s"

		user_associations {
			user_id = "%s"
		}
	}`, msoTenantNameUserAssociations, msoTenantNameUserAssociations, description, userID)
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

func testAccMSOTenantLookupUserID(msoClient *client.Client, username string) (string, error) {
	paths := []string{"api/v1/users"}
	if msoClient.GetPlatform() == "nd" {
		paths = append(paths, "api/v2/users")
	}

	var lastErr error
	for _, path := range paths {
		con, err := msoClient.GetViaURL(path)
		if err != nil {
			lastErr = err
			continue
		}

		if userID, ok := testAccMSOTenantFindUserIDInContainer(con, username); ok {
			return userID, nil
		}
	}

	if lastErr != nil {
		return "", fmt.Errorf("failed to look up user %q for tenant user_associations test: %w", username, lastErr)
	}
	return "", fmt.Errorf("user %q not found for tenant user_associations test", username)
}

func testAccMSOTenantFindUserIDInContainer(con *container.Container, username string) (string, bool) {
	if users, ok := con.Data().([]interface{}); ok {
		for _, user := range users {
			userMap, ok := user.(map[string]interface{})
			if !ok {
				continue
			}
			if userID, ok := testAccMSOTenantMatchUserMap(userMap, username); ok {
				return userID, true
			}
		}
	}

	if users, ok := con.S("users").Data().([]interface{}); ok {
		for _, user := range users {
			userMap, ok := user.(map[string]interface{})
			if !ok {
				continue
			}
			if userID, ok := testAccMSOTenantMatchUserMap(userMap, username); ok {
				return userID, true
			}
		}
	}

	return "", false
}

func testAccMSOTenantMatchUserMap(userMap map[string]interface{}, username string) (string, bool) {
	for _, key := range []string{"loginID", "username"} {
		if fmt.Sprintf("%v", userMap[key]) != username {
			continue
		}

		for _, idKey := range []string{"userID", "id"} {
			if userID := fmt.Sprintf("%v", userMap[idKey]); userID != "" && userID != "<nil>" {
				return userID, true
			}
		}
	}

	return "", false
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
