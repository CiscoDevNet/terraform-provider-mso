package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOTenantPoliciesIPSLATrackListDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: IPSLA Track List Data Source") },
				Config:    testAccMSOTenantPoliciesIPSLATrackListDataSource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName+"_data", "name", msoTenantPolicyTemplateIPSLATrackListName),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName+"_data", "description", "Terraform test IPSLA Track List"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName+"_data", "threshold_down", "11"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName+"_data", "threshold_up", "12"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName+"_data", "type", "percentage"),
					resource.TestCheckResourceAttr("data.mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName+"_data", "members.#", "1"),
					CustomTestCheckTypeSetElemAttrs("data.mso_tenant_policies_ipsla_track_list."+msoTenantPolicyTemplateIPSLATrackListName+"_data", "members",
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
		},
	})
}

func testAccMSOTenantPoliciesIPSLATrackListDataSource() string {
	return fmt.Sprintf(`%[1]s
	data "mso_tenant_policies_ipsla_track_list" "%[2]s_data" {
	    template_id	= mso_tenant_policies_ipsla_track_list.%[2]s.template_id
	    name		= "%[2]s"
    }`, testAccMSOTenantPoliciesIPSLATrackListConfigUpdateRemovingExtraEntry(), msoTenantPolicyTemplateIPSLATrackListName)
}
