package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var (
	testNetflowExporterTemplateID string
	testNetflowExporterUUID       string
)

func TestAccMSOTenantPoliciesNetflowExporterResource(t *testing.T) {
	resourceName := "mso_tenant_policies_netflow_exporter.netflow_exporter"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccVersionCheck(t, "5.1") },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create NetFlow Exporter") },
				Config:    testAccMSOTenantPoliciesNetflowExporterConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_netflow_exporter"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok {
							return fmt.Errorf("resource %s not found in state", resourceName)
						}
						testNetflowExporterTemplateID = rs.Primary.Attributes["template_id"]
						testNetflowExporterUUID = rs.Primary.Attributes["uuid"]
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate NetFlow Exporter after out-of-band deletion")
					if err := manuallyDeleteNetflowExporterFromTemplate(testNetflowExporterTemplateID, testNetflowExporterUUID); err != nil {
						t.Fatalf("Failed to manually delete NetFlow Exporter: %v", err)
					}
				},
				Config: testAccMSOTenantPoliciesNetflowExporterConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_netflow_exporter"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update NetFlow Exporter Name") },
				Config:    testAccMSOTenantPoliciesNetflowExporterConfigUpdateName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_netflow_exporter_updated"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import NetFlow Exporter") },
				ResourceName:      "mso_tenant_policies_netflow_exporter.netflow_exporter",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments("mso_tenant_policies_netflow_exporter", "tenantPolicyTemplate", "template", "netFlowExporters"),
	})
}

func manuallyDeleteNetflowExporterFromTemplate(templateID, exporterUUID string) error {
	msoClient := testAccProvider.Meta().(*client.Client)
	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
	if err != nil {
		return fmt.Errorf("manual delete: get template %s: %w", templateID, err)
	}
	exporterIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", exporterUUID, "tenantPolicyTemplate", "template", "netFlowExporters")
	if err != nil {
		return fmt.Errorf("manual delete: locate NetFlow Exporter %q: %w", exporterUUID, err)
	}
	path := fmt.Sprintf("/tenantPolicyTemplate/template/netFlowExporters/%d", exporterIndex)
	if _, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateID), models.GetRemovePatchPayload(path)); err != nil {
		return fmt.Errorf("manual delete: remove NetFlow Exporter %q: %w", exporterUUID, err)
	}
	return nil
}

func testAccMSOTenantPoliciesNetflowExporterConfigCreate() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_netflow_exporter" "netflow_exporter" {
        template_id = mso_template.template_tenant.id
        name        = "test_netflow_exporter"
    }`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesNetflowExporterConfigUpdateName() string {
	return fmt.Sprintf(`%s
    resource "mso_tenant_policies_netflow_exporter" "netflow_exporter" {
        template_id = mso_template.template_tenant.id
        name        = "test_netflow_exporter_updated"
    }`, testAccMSOTemplateResourceTenantConfig())
}
