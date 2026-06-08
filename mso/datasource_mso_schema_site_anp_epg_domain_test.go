package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test fixture for data.mso_schema_site_anp_epg_domain.
//
// Strategy:
//   - Reuse the resource-create config from the resource acceptance test
//     to seed one site ANP EPG domain per domain_type (phys, l2, l3, fc,
//     vmm-by-name, vmm-by-dn). All six are siblings under the shared
//     EPG's domainAssociations[] array.
//   - For each, declare a data.mso_schema_site_anp_epg_domain that
//     references the resource by either domain_dn or
//     domain_name+domain_type(+vmm_domain_type), then assert pair-equality
//     against the resource attributes.

func TestAccMSOSchemaSiteAnpEpgDomainDatasource(t *testing.T) {
	physAddr := "mso_schema_site_anp_epg_domain." + msoSiteAnpEpgDomainLabelPhys
	l2Addr := "mso_schema_site_anp_epg_domain." + msoSiteAnpEpgDomainLabelL2
	l3Addr := "mso_schema_site_anp_epg_domain." + msoSiteAnpEpgDomainLabelL3
	fcAddr := "mso_schema_site_anp_epg_domain." + msoSiteAnpEpgDomainLabelFc
	vmmNameAddr := "mso_schema_site_anp_epg_domain." + msoSiteAnpEpgDomainLabelVmmName
	vmmDnAddr := "mso_schema_site_anp_epg_domain." + msoSiteAnpEpgDomainLabelVmmDn

	physDs := "data.mso_schema_site_anp_epg_domain.ds_phys"
	l2Ds := "data.mso_schema_site_anp_epg_domain.ds_l2"
	l3Ds := "data.mso_schema_site_anp_epg_domain.ds_l3"
	fcDs := "data.mso_schema_site_anp_epg_domain.ds_fc"
	vmmNameDs := "data.mso_schema_site_anp_epg_domain.ds_vmm_name"
	vmmDnDs := "data.mso_schema_site_anp_epg_domain.ds_vmm_dn"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpEpgDomainDestroy,
		Steps: []resource.TestStep{
			// Each step's PreConfig calls cleanupOrphanSchemaSiteTestResources
			// (defined in datasource_mso_schema_site_test.go) to remove any
			// orphaned tenant/schema left on MSO from a previous step's
			// rolled-back apply. SDK v1 rolls back Terraform state on any
			// apply error (including data source read errors), but the
			// prereq resources have already been created server-side --
			// without cleanup the next step (or rerun) fails with
			// "Tenant/Schema: '...' already exists". The helper deletes by
			// the same package globals (msoSchemaName, msoTenantName) that
			// this test's prereq config uses, so it is directly reusable.
			{
				PreConfig: func() {
					fmt.Println("Test: Read site ANP EPG domain datasource with no DN inputs (expect error)")
					cleanupOrphanSchemaSiteTestResources(t)
				},
				Config:      testAccMSOSchemaSiteAnpEpgDomainDatasourceMissingDn(),
				ExpectError: regexp.MustCompile(`domain_dn or domain_name`),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Read site ANP EPG domain datasource for non-existing domain (expect error)")
					cleanupOrphanSchemaSiteTestResources(t)
				},
				Config:      testAccMSOSchemaSiteAnpEpgDomainDatasourceNotFound(),
				ExpectError: regexp.MustCompile(`Unable to find the Site ANP EPG Domain`),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Read all six site ANP EPG domain datasources")
					cleanupOrphanSchemaSiteTestResources(t)
				},
				Config: testAccMSOSchemaSiteAnpEpgDomainDatasourceAll(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// physicalDomain via domain_dn
					checkSiteAnpEpgDomainDatasourcePair(physDs, physAddr),
					resource.TestCheckResourceAttr(physDs, "domain_type", "physicalDomain"),

					// l2ExtDomain via domain_name + domain_type
					checkSiteAnpEpgDomainDatasourcePair(l2Ds, l2Addr),
					resource.TestCheckResourceAttr(l2Ds, "domain_type", "l2ExtDomain"),
					resource.TestCheckResourceAttrPair(l2Ds, "domain_dn", l2Addr, "domain_dn"),

					// l3ExtDomain via domain_name + domain_type
					checkSiteAnpEpgDomainDatasourcePair(l3Ds, l3Addr),
					resource.TestCheckResourceAttr(l3Ds, "domain_type", "l3ExtDomain"),
					resource.TestCheckResourceAttrPair(l3Ds, "domain_dn", l3Addr, "domain_dn"),

					// fibreChannelDomain via domain_name + domain_type
					checkSiteAnpEpgDomainDatasourcePair(fcDs, fcAddr),
					resource.TestCheckResourceAttr(fcDs, "domain_type", "fibreChannelDomain"),
					resource.TestCheckResourceAttrPair(fcDs, "domain_dn", fcAddr, "domain_dn"),

					// vmmDomain via domain_name + domain_type + vmm_domain_type, full attribute matrix
					checkSiteAnpEpgDomainDatasourcePair(vmmNameDs, vmmNameAddr),
					resource.TestCheckResourceAttr(vmmNameDs, "domain_type", "vmmDomain"),
					resource.TestCheckResourceAttrPair(vmmNameDs, "domain_dn", vmmNameAddr, "domain_dn"),
					resource.TestCheckResourceAttrPair(vmmNameDs, "vlan_encap_mode", vmmNameAddr, "vlan_encap_mode"),
					resource.TestCheckResourceAttrPair(vmmNameDs, "allow_micro_segmentation", vmmNameAddr, "allow_micro_segmentation"),
					resource.TestCheckResourceAttrPair(vmmNameDs, "switching_mode", vmmNameAddr, "switching_mode"),
					resource.TestCheckResourceAttrPair(vmmNameDs, "switch_type", vmmNameAddr, "switch_type"),
					resource.TestCheckResourceAttrPair(vmmNameDs, "micro_seg_vlan_type", vmmNameAddr, "micro_seg_vlan_type"),
					resource.TestCheckResourceAttrPair(vmmNameDs, "micro_seg_vlan", vmmNameAddr, "micro_seg_vlan"),
					resource.TestCheckResourceAttrPair(vmmNameDs, "port_encap_vlan_type", vmmNameAddr, "port_encap_vlan_type"),
					resource.TestCheckResourceAttrPair(vmmNameDs, "port_encap_vlan", vmmNameAddr, "port_encap_vlan"),
					resource.TestCheckResourceAttrPair(vmmNameDs, "enhanced_lag_policy_name", vmmNameAddr, "enhanced_lag_policy_name"),
					resource.TestCheckResourceAttrPair(vmmNameDs, "enhanced_lag_policy_dn", vmmNameAddr, "enhanced_lag_policy_dn"),
					resource.TestCheckResourceAttrPair(vmmNameDs, "delimiter", vmmNameAddr, "delimiter"),
					resource.TestCheckResourceAttrPair(vmmNameDs, "binding_type", vmmNameAddr, "binding_type"),
					resource.TestCheckResourceAttrPair(vmmNameDs, "port_allocation", vmmNameAddr, "port_allocation"),
					resource.TestCheckResourceAttrPair(vmmNameDs, "num_ports", vmmNameAddr, "num_ports"),
					resource.TestCheckResourceAttrPair(vmmNameDs, "netflow", vmmNameAddr, "netflow"),
					resource.TestCheckResourceAttrPair(vmmNameDs, "allow_promiscuous", vmmNameAddr, "allow_promiscuous"),
					resource.TestCheckResourceAttrPair(vmmNameDs, "mac_changes", vmmNameAddr, "mac_changes"),
					resource.TestCheckResourceAttrPair(vmmNameDs, "forged_transmits", vmmNameAddr, "forged_transmits"),
					resource.TestCheckResourceAttrPair(vmmNameDs, "custom_epg_name", vmmNameAddr, "custom_epg_name"),

					// vmmDomain via domain_dn (covers dn != "" branch in datasource Read)
					checkSiteAnpEpgDomainDatasourcePair(vmmDnDs, vmmDnAddr),
					resource.TestCheckResourceAttr(vmmDnDs, "domain_type", "vmmDomain"),
					resource.TestCheckResourceAttrPair(vmmDnDs, "domain_dn", vmmDnAddr, "domain_dn"),
				),
			},
		},
	})
}

// checkSiteAnpEpgDomainDatasourcePair asserts pair-equality between a
// datasource and its source resource for the identifying attributes that
// are common to every domain_type.
func checkSiteAnpEpgDomainDatasourcePair(ds, res string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrPair(ds, "id", res, "id"),
		resource.TestCheckResourceAttrPair(ds, "schema_id", res, "schema_id"),
		resource.TestCheckResourceAttrPair(ds, "site_id", res, "site_id"),
		resource.TestCheckResourceAttrPair(ds, "template_name", res, "template_name"),
		resource.TestCheckResourceAttrPair(ds, "anp_name", res, "anp_name"),
		resource.TestCheckResourceAttrPair(ds, "epg_name", res, "epg_name"),
		resource.TestCheckResourceAttrPair(ds, "deploy_immediacy", res, "deploy_immediacy"),
		resource.TestCheckResourceAttrPair(ds, "resolution_immediacy", res, "resolution_immediacy"),
	)
}

// testAccMSOSchemaSiteAnpEpgDomainDatasourceMissingDn emits the
// resource-create config plus a datasource block with no DN inputs.
func testAccMSOSchemaSiteAnpEpgDomainDatasourceMissingDn() string {
	return fmt.Sprintf(`%[1]s
data "mso_schema_site_anp_epg_domain" "bad" {
	schema_id     = mso_schema_site_anp_epg_domain.%[2]s.schema_id
	site_id       = mso_schema_site_anp_epg_domain.%[2]s.site_id
	template_name = mso_schema_site_anp_epg_domain.%[2]s.template_name
	anp_name      = mso_schema_site_anp_epg_domain.%[2]s.anp_name
	epg_name      = mso_schema_site_anp_epg_domain.%[2]s.epg_name
}
`, testAccMSOSchemaSiteAnpEpgDomainConfigAll(false), msoSiteAnpEpgDomainLabelPhys)
}

// testAccMSOSchemaSiteAnpEpgDomainDatasourceNotFound queries a synthetic
// DN that does not exist on the EPG, exercising the !found error path in
// dataSourceMSOSchemaSiteAnpEpgDomainRead.
func testAccMSOSchemaSiteAnpEpgDomainDatasourceNotFound() string {
	return fmt.Sprintf(`%[1]s
data "mso_schema_site_anp_epg_domain" "missing" {
	schema_id     = mso_schema_site_anp_epg_domain.%[2]s.schema_id
	site_id       = mso_schema_site_anp_epg_domain.%[2]s.site_id
	template_name = mso_schema_site_anp_epg_domain.%[2]s.template_name
	anp_name      = mso_schema_site_anp_epg_domain.%[2]s.anp_name
	epg_name      = mso_schema_site_anp_epg_domain.%[2]s.epg_name
	domain_dn     = "uni/phys-does_not_exist"
}
`, testAccMSOSchemaSiteAnpEpgDomainConfigAll(false), msoSiteAnpEpgDomainLabelPhys)
}

// testAccMSOSchemaSiteAnpEpgDomainDatasourceAll emits the resource-create
// config plus one datasource block per domain_type covered.
func testAccMSOSchemaSiteAnpEpgDomainDatasourceAll() string {
	return fmt.Sprintf(`%[1]s
data "mso_schema_site_anp_epg_domain" "ds_phys" {
	schema_id     = mso_schema_site_anp_epg_domain.%[2]s.schema_id
	site_id       = mso_schema_site_anp_epg_domain.%[2]s.site_id
	template_name = mso_schema_site_anp_epg_domain.%[2]s.template_name
	anp_name      = mso_schema_site_anp_epg_domain.%[2]s.anp_name
	epg_name      = mso_schema_site_anp_epg_domain.%[2]s.epg_name
	domain_dn     = mso_schema_site_anp_epg_domain.%[2]s.domain_dn
}

data "mso_schema_site_anp_epg_domain" "ds_l2" {
	schema_id     = mso_schema_site_anp_epg_domain.%[3]s.schema_id
	site_id       = mso_schema_site_anp_epg_domain.%[3]s.site_id
	template_name = mso_schema_site_anp_epg_domain.%[3]s.template_name
	anp_name      = mso_schema_site_anp_epg_domain.%[3]s.anp_name
	epg_name      = mso_schema_site_anp_epg_domain.%[3]s.epg_name
	domain_name   = mso_schema_site_anp_epg_domain.%[3]s.domain_name
	domain_type   = mso_schema_site_anp_epg_domain.%[3]s.domain_type
}

data "mso_schema_site_anp_epg_domain" "ds_l3" {
	schema_id     = mso_schema_site_anp_epg_domain.%[4]s.schema_id
	site_id       = mso_schema_site_anp_epg_domain.%[4]s.site_id
	template_name = mso_schema_site_anp_epg_domain.%[4]s.template_name
	anp_name      = mso_schema_site_anp_epg_domain.%[4]s.anp_name
	epg_name      = mso_schema_site_anp_epg_domain.%[4]s.epg_name
	domain_name   = mso_schema_site_anp_epg_domain.%[4]s.domain_name
	domain_type   = mso_schema_site_anp_epg_domain.%[4]s.domain_type
}

data "mso_schema_site_anp_epg_domain" "ds_fc" {
	schema_id     = mso_schema_site_anp_epg_domain.%[5]s.schema_id
	site_id       = mso_schema_site_anp_epg_domain.%[5]s.site_id
	template_name = mso_schema_site_anp_epg_domain.%[5]s.template_name
	anp_name      = mso_schema_site_anp_epg_domain.%[5]s.anp_name
	epg_name      = mso_schema_site_anp_epg_domain.%[5]s.epg_name
	domain_name   = mso_schema_site_anp_epg_domain.%[5]s.domain_name
	domain_type   = mso_schema_site_anp_epg_domain.%[5]s.domain_type
}

data "mso_schema_site_anp_epg_domain" "ds_vmm_name" {
	schema_id       = mso_schema_site_anp_epg_domain.%[6]s.schema_id
	site_id         = mso_schema_site_anp_epg_domain.%[6]s.site_id
	template_name   = mso_schema_site_anp_epg_domain.%[6]s.template_name
	anp_name        = mso_schema_site_anp_epg_domain.%[6]s.anp_name
	epg_name        = mso_schema_site_anp_epg_domain.%[6]s.epg_name
	domain_name     = mso_schema_site_anp_epg_domain.%[6]s.domain_name
	domain_type     = mso_schema_site_anp_epg_domain.%[6]s.domain_type
	vmm_domain_type = mso_schema_site_anp_epg_domain.%[6]s.vmm_domain_type
}

data "mso_schema_site_anp_epg_domain" "ds_vmm_dn" {
	schema_id     = mso_schema_site_anp_epg_domain.%[7]s.schema_id
	site_id       = mso_schema_site_anp_epg_domain.%[7]s.site_id
	template_name = mso_schema_site_anp_epg_domain.%[7]s.template_name
	anp_name      = mso_schema_site_anp_epg_domain.%[7]s.anp_name
	epg_name      = mso_schema_site_anp_epg_domain.%[7]s.epg_name
	domain_dn     = mso_schema_site_anp_epg_domain.%[7]s.domain_dn
}
`,
		testAccMSOSchemaSiteAnpEpgDomainConfigAll(false),
		msoSiteAnpEpgDomainLabelPhys,
		msoSiteAnpEpgDomainLabelL2,
		msoSiteAnpEpgDomainLabelL3,
		msoSiteAnpEpgDomainLabelFc,
		msoSiteAnpEpgDomainLabelVmmName,
		msoSiteAnpEpgDomainLabelVmmDn,
	)
}
