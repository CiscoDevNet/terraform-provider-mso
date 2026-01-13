package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOTenantPoliciesIPSLATrackListResource(t *testing.T) {
	print(testAccMSOTenantPoliciesIPSLATrackListConfigCreate())
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create IPSLA Track List") },
				Config:    testAccMSOTenantPoliciesIPSLATrackListConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName, "name", msoTenantPolicyTemplateIPSLATrackListName),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName, "description", "Terraform test IPSLA Track List"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName, "threshold_down", "11"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName, "threshold_up", "12"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName, "type", "weight"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName, "members.#", "1"),
					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName, "members",
						map[string]string{
							"destination_ip":               "1.1.1.1",
							"ipsla_monitoring_policy_uuid": fmt.Sprintf("mso_tenant_policies_ipsla_monitoring_policy.%s.uuid", msoTenantPolicyTemplateIPSLAMonitoringPolicyName),
							"scope_type":                   "bd",
							"scope_uuid":                   fmt.Sprintf("mso_schema_template_bd.%s.uuid", msoSchemaTemplateBdName),
							"weight":                       "10",
						},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update IPSLA Track List adding extra entry") },
				Config:    testAccMSOTenantPoliciesIPSLATrackListConfigUpdateAddingExtraEntry(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName, "name", msoTenantPolicyTemplateIPSLATrackListName),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName, "description", "Terraform test IPSLA Track List"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName, "threshold_down", "21"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName, "threshold_up", "22"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName, "type", "percentage"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName, "members.#", "2"),
					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName, "members",
						map[string]string{
							"destination_ip":               "1.1.1.3",
							"ipsla_monitoring_policy_uuid": fmt.Sprintf("mso_tenant_policies_ipsla_monitoring_policy.%s.uuid", msoTenantPolicyTemplateIPSLAMonitoringPolicyName),
							"scope_type":                   "bd",
							"scope_uuid":                   fmt.Sprintf("mso_schema_template_bd.%s.uuid", msoSchemaTemplateBdName),
							"weight":                       "10",
						},
					),
					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName, "members",
						map[string]string{
							"destination_ip":               "1.1.1.2",
							"ipsla_monitoring_policy_uuid": fmt.Sprintf("mso_tenant_policies_ipsla_monitoring_policy.%s.uuid", msoTenantPolicyTemplateIPSLAMonitoringPolicyName),
							"scope_type":                   "bd",
							"scope_uuid":                   fmt.Sprintf("mso_schema_template_bd.%s.uuid", msoSchemaTemplateBdName),
							"weight":                       "10",
						},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update IPSLA Track List removing extra entry") },
				Config:    testAccMSOTenantPoliciesIPSLATrackListConfigUpdateRemovingExtraEntry(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName, "name", msoTenantPolicyTemplateIPSLATrackListName),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName, "description", "Terraform test IPSLA Track List"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName, "threshold_down", "11"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName, "threshold_up", "12"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName, "type", "percentage"),
					resource.TestCheckResourceAttr("mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName, "members.#", "1"),
					CustomTestCheckTypeSetElemAttrs("mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName, "members",
						map[string]string{
							"destination_ip":               "1.1.1.2",
							"ipsla_monitoring_policy_uuid": fmt.Sprintf("mso_tenant_policies_ipsla_monitoring_policy.%s.uuid", msoTenantPolicyTemplateIPSLAMonitoringPolicyName),
							"scope_type":                   "bd",
							"scope_uuid":                   fmt.Sprintf("mso_schema_template_bd.%s.uuid", msoSchemaTemplateBdName),
							"weight":                       "10",
						},
					),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import IPSLA Track List") },
				ResourceName:      "mso_tenant_policies_ipsla_track_list." + msoTenantPolicyTemplateIPSLATrackListName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithArguments("mso_tenant_policies_ipsla_track_list", "ipslaTrackList"),
	})
}

var IPSLATrackListPreConfig = testSiteConfigAnsibleTest() + testTenantConfig() + testSchemaConfig() + testSchemaTemplateVrfConfig() + testSchemaTemplateBdConfig() + testTenantPolicyTemplateConfig() + testTenantPolicyTemplateIPSLAMonitoringPolicyConfig()

func testAccMSOTenantPoliciesIPSLATrackListConfigCreate() string {
	return fmt.Sprintf(`%[5]s
	resource "mso_tenant_policies_ipsla_track_list" "%[1]s" {
		template_id    = mso_template.%[2]s.id
		name           = "%[1]s"
		description    = "Terraform test IPSLA Track List"
		threshold_down = 11
		threshold_up   = 12
		type           = "weight"
		members {
			destination_ip               = "1.1.1.1"
			ipsla_monitoring_policy_uuid = mso_tenant_policies_ipsla_monitoring_policy.%[3]s.uuid
			scope_type                   = "bd"
			scope_uuid                   = mso_schema_template_bd.%[4]s.uuid
		}
	}`, msoTenantPolicyTemplateIPSLATrackListName, msoTenantPolicyTemplateName, msoTenantPolicyTemplateIPSLAMonitoringPolicyName, msoSchemaTemplateBdName, IPSLATrackListPreConfig)
}

func testAccMSOTenantPoliciesIPSLATrackListConfigUpdateAddingExtraEntry() string {
	return fmt.Sprintf(`%[5]s
	resource "mso_tenant_policies_ipsla_track_list" "%[1]s" {
		template_id    = mso_template.%[2]s.id
		name           = "%[1]s"
		description    = "Terraform test IPSLA Track List"
		threshold_down = 21
		threshold_up   = 22
		type           = "percentage"
		members {
			destination_ip               = "1.1.1.3"
			ipsla_monitoring_policy_uuid = mso_tenant_policies_ipsla_monitoring_policy.%[3]s.uuid
			scope_type                   = "bd"
			scope_uuid                   = mso_schema_template_bd.%[4]s.uuid
		}
		members {
			destination_ip               = "1.1.1.2"
			ipsla_monitoring_policy_uuid = mso_tenant_policies_ipsla_monitoring_policy.%[3]s.uuid
			scope_type                   = "bd"
			scope_uuid                   = mso_schema_template_bd.%[4]s.uuid
		}
	}`, msoTenantPolicyTemplateIPSLATrackListName, msoTenantPolicyTemplateName, msoTenantPolicyTemplateIPSLAMonitoringPolicyName, msoSchemaTemplateBdName, IPSLATrackListPreConfig)
}

func testAccMSOTenantPoliciesIPSLATrackListConfigUpdateRemovingExtraEntry() string {
	return fmt.Sprintf(`%[5]s
	resource "mso_tenant_policies_ipsla_track_list" "%[1]s" {
		template_id    = mso_template.%[2]s.id
		name           = "%[1]s"
		description    = "Terraform test IPSLA Track List"
		threshold_down = 11
		threshold_up   = 12
		type           = "percentage"
		members {
			destination_ip               = "1.1.1.2"
			ipsla_monitoring_policy_uuid = mso_tenant_policies_ipsla_monitoring_policy.%[3]s.uuid
			scope_type                   = "bd"
			scope_uuid                   = mso_schema_template_bd.%[4]s.uuid
		}
	}`, msoTenantPolicyTemplateIPSLATrackListName, msoTenantPolicyTemplateName, msoTenantPolicyTemplateIPSLAMonitoringPolicyName, msoSchemaTemplateBdName, IPSLATrackListPreConfig)
}
