package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMSOTenantPoliciesNetflowExporterSiteDataSource(t *testing.T) {
	dataSourceName := "data.mso_tenant_policies_netflow_exporter_site.netflow_exporter_site"
	resourceName := "mso_tenant_policies_netflow_exporter_site.netflow_exporter_site"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccVersionCheck(t, "5.1") },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: NetFlow Exporter Site Data Source - Not Found") },
				Config:      testAccMSOTenantPoliciesNetflowExporterSiteDataSourceNotFound(),
				ExpectError: regexp.MustCompile(`NetFlow Exporter .* does not exist on site .*`),
			},
			{
				PreConfig: func() { fmt.Println("Test: NetFlow Exporter Site Data Source") },
				Config:    testAccMSOTenantPoliciesNetflowExporterSiteDataSource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "template_id", resourceName, "template_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "site_id", resourceName, "site_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "netflow_exporter_uuid", resourceName, "netflow_exporter_uuid"),
					resource.TestCheckResourceAttr(dataSourceName, "description", ""),
					resource.TestCheckResourceAttr(dataSourceName, "source_ip_address_type", "inbandMgmtIP"),
					resource.TestCheckResourceAttr(dataSourceName, "source_ip", ""),
					resource.TestCheckResourceAttr(dataSourceName, "destination_port", "10"),
					resource.TestCheckResourceAttr(dataSourceName, "destination_ip", "10.0.0.15"),
					resource.TestCheckResourceAttr(dataSourceName, "dscp", "cs1"),
					resource.TestCheckResourceAttr(dataSourceName, "application_epg.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "external_epg.#", "0"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: NetFlow Exporter Site Data Source with application EPG") },
				Config:    testAccMSOTenantPoliciesNetflowExporterSiteDataSourceApplicationEpg(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "source_ip_address_type", "inbandMgmtIP"),
					resource.TestCheckResourceAttr(dataSourceName, "application_epg.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "application_epg.0.tenant_name", msoTenantName),
					resource.TestCheckResourceAttr(dataSourceName, "application_epg.0.anp_name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr(dataSourceName, "application_epg.0.epg_name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr(dataSourceName, "application_epg.0.vrf_name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr(dataSourceName, "application_epg.0.vrf_tenant_name", msoTenantName),
					resource.TestCheckResourceAttr(dataSourceName, "external_epg.#", "0"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: NetFlow Exporter Site Data Source with external EPG") },
				Config:    testAccMSOTenantPoliciesNetflowExporterSiteDataSourceExternalEpg(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "source_ip_address_type", "oobMgmtIP"),
					resource.TestCheckResourceAttr(dataSourceName, "application_epg.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "external_epg.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "external_epg.0.tenant_name", msoTenantName),
					resource.TestCheckResourceAttr(dataSourceName, "external_epg.0.l3out_name", msoSchemaTemplateL3outName),
					resource.TestCheckResourceAttr(dataSourceName, "external_epg.0.external_epg_name", msoSchemaTemplateExtEpgName),
					resource.TestCheckResourceAttr(dataSourceName, "external_epg.0.vrf_name", msoSchemaTemplateVrfName),
					resource.TestCheckResourceAttr(dataSourceName, "external_epg.0.vrf_tenant_name", msoTenantName),
				),
			},
		},
	})
}

func testAccMSOTenantPoliciesNetflowExporterSiteDataSource() string {
	return testAccMSOTenantPoliciesNetflowExporterSiteDataSourceConfig(testAccMSOTenantPoliciesNetflowExporterSiteConfigCreate())
}

func testAccMSOTenantPoliciesNetflowExporterSiteDataSourceApplicationEpg() string {
	return testAccMSOTenantPoliciesNetflowExporterSiteDataSourceConfig(testAccMSOTenantPoliciesNetflowExporterSiteConfigApplicationEpg())
}

func testAccMSOTenantPoliciesNetflowExporterSiteDataSourceExternalEpg() string {
	return testAccMSOTenantPoliciesNetflowExporterSiteDataSourceConfig(testAccMSOTenantPoliciesNetflowExporterSiteConfigExternalEpgOobMgmtIp())
}

func testAccMSOTenantPoliciesNetflowExporterSiteDataSourceConfig(resourceConfig string) string {
	return fmt.Sprintf(`%[1]s
	data "mso_tenant_policies_netflow_exporter_site" "netflow_exporter_site" {
		template_id           = mso_tenant_policies_netflow_exporter_site.netflow_exporter_site.template_id
		site_id               = mso_tenant_policies_netflow_exporter_site.netflow_exporter_site.site_id
		netflow_exporter_uuid = mso_tenant_policies_netflow_exporter_site.netflow_exporter_site.netflow_exporter_uuid
	}
`, resourceConfig)
}

func testAccMSOTenantPoliciesNetflowExporterSiteDataSourceNotFound() string {
	return fmt.Sprintf(`%[1]s
	data "mso_tenant_policies_netflow_exporter_site" "netflow_exporter_site" {
		template_id           = mso_tenant_policies_netflow_exporter.netflow_exporter.template_id
		site_id               = data.mso_site.%[2]s.id
		netflow_exporter_uuid = "00000000-0000-0000-0000-000000000000"
	}
`, testAccMSOTenantPoliciesNetflowExporterSiteBaseConfig(), msoTemplateSiteName1)
}
