package mso

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Test fixture for mso_schema_site_anp_epg_domain.
//
// Strategy:
//   - One shared schema/template/anp/epg (reused from the existing
//     testAccMSOSchemaSiteAnpEpgConfigCreate helper).
//   - One mso_schema_site_anp_epg_domain block per domain_type, all sharing
//     the parent EPG so they coexist as sibling entries under the EPG's
//     domainAssociations[].
//   - Synthetic NDO DNs (e.g. uni/phys-<rand>). NDO writes the DN string
//     into the schema document at PATCH time and does not validate the
//     domain reference until template deploy, which these tests do not
//     perform.
//   - Targets ND 4.1 / NDO 5.x (==4.2+ on the NDO side), so the full VMM
//     attribute matrix is always available.
//
// The lab must have the ansible_test and ansible_test_2 sites onboarded.

// Stable Terraform resource labels (one per domain_type covered).
const (
	msoSiteAnpEpgDomainLabelPhys    = "phys"
	msoSiteAnpEpgDomainLabelL2      = "l2"
	msoSiteAnpEpgDomainLabelL3      = "l3"
	msoSiteAnpEpgDomainLabelFc      = "fc"
	msoSiteAnpEpgDomainLabelVmmName = "vmm_by_name"
	msoSiteAnpEpgDomainLabelVmmDn   = "vmm_by_dn"
)

// Per-domain random names, kept distinct so each domain produces a unique
// DN under the shared EPG.
var (
	msoSiteAnpEpgDomainNamePhys    = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	msoSiteAnpEpgDomainNameL2      = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	msoSiteAnpEpgDomainNameL3      = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	msoSiteAnpEpgDomainNameFc      = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	msoSiteAnpEpgDomainNameVmmName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	msoSiteAnpEpgDomainNameVmmDn   = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
)

func TestAccMSOSchemaSiteAnpEpgDomainResource(t *testing.T) {
	physAddr := "mso_schema_site_anp_epg_domain." + msoSiteAnpEpgDomainLabelPhys
	l2Addr := "mso_schema_site_anp_epg_domain." + msoSiteAnpEpgDomainLabelL2
	l3Addr := "mso_schema_site_anp_epg_domain." + msoSiteAnpEpgDomainLabelL3
	fcAddr := "mso_schema_site_anp_epg_domain." + msoSiteAnpEpgDomainLabelFc
	vmmNameAddr := "mso_schema_site_anp_epg_domain." + msoSiteAnpEpgDomainLabelVmmName
	vmmDnAddr := "mso_schema_site_anp_epg_domain." + msoSiteAnpEpgDomainLabelVmmDn

	physDn := fmt.Sprintf("uni/phys-%s", msoSiteAnpEpgDomainNamePhys)
	l2Dn := fmt.Sprintf("uni/l2dom-%s", msoSiteAnpEpgDomainNameL2)
	l3Dn := fmt.Sprintf("uni/l3dom-%s", msoSiteAnpEpgDomainNameL3)
	fcDn := fmt.Sprintf("uni/fc-%s", msoSiteAnpEpgDomainNameFc)
	vmmNameDn := fmt.Sprintf("uni/vmmp-VMware/dom-%s", msoSiteAnpEpgDomainNameVmmName)
	vmmDnDn := fmt.Sprintf("uni/vmmp-VMware/dom-%s", msoSiteAnpEpgDomainNameVmmDn)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteAnpEpgDomainDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Create site ANP EPG domain with no DN inputs (expect error)") },
				Config:      testAccMSOSchemaSiteAnpEpgDomainConfigMissingDn(),
				ExpectError: regexp.MustCompile(`domain_dn or domain_name`),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Create site ANP EPG domain with domain_name but no domain_type (expect error)")
				},
				Config:      testAccMSOSchemaSiteAnpEpgDomainConfigNameWithoutType(),
				ExpectError: regexp.MustCompile(`domain_type is required when domain_name is provided`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Create vmmDomain without vmm_domain_type (expect error)") },
				Config:      testAccMSOSchemaSiteAnpEpgDomainConfigVmmWithoutVmmType(),
				ExpectError: regexp.MustCompile(`vmm_domain_type is required when domain_type is vmmDomain`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Conflicts between domain_dn and domain_name (expect error)") },
				Config:      testAccMSOSchemaSiteAnpEpgDomainConfigConflicts(),
				ExpectError: regexp.MustCompile(`"domain_dn": conflicts with`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create all six site ANP EPG domains") },
				Config:    testAccMSOSchemaSiteAnpEpgDomainConfigAll(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					// physicalDomain via domain_dn.
					// domain_type is intentionally NOT asserted here: the resource
					// Read deliberately skips d.Set("domain_type", ...) when the
					// user supplied domain_dn (the two attributes ConflictsWith
					// each other in the schema). The datasource Read does set it
					// unconditionally, so the datasource test asserts it instead.
					resource.TestCheckResourceAttr(physAddr, "domain_dn", physDn),
					resource.TestCheckResourceAttr(physAddr, "deploy_immediacy", "lazy"),
					resource.TestCheckResourceAttr(physAddr, "resolution_immediacy", "lazy"),
					resource.TestMatchResourceAttr(physAddr, "id", regexp.MustCompile(regexp.QuoteMeta("/domainAssociations/"+physDn)+"$")),

					// l2ExtDomain via domain_name + domain_type
					resource.TestCheckResourceAttr(l2Addr, "domain_name", msoSiteAnpEpgDomainNameL2),
					resource.TestCheckResourceAttr(l2Addr, "domain_type", "l2ExtDomain"),
					resource.TestCheckResourceAttr(l2Addr, "deploy_immediacy", "lazy"),
					resource.TestCheckResourceAttr(l2Addr, "resolution_immediacy", "lazy"),
					resource.TestMatchResourceAttr(l2Addr, "id", regexp.MustCompile(regexp.QuoteMeta("/domainAssociations/"+l2Dn)+"$")),

					// l3ExtDomain via domain_name + domain_type
					resource.TestCheckResourceAttr(l3Addr, "domain_name", msoSiteAnpEpgDomainNameL3),
					resource.TestCheckResourceAttr(l3Addr, "domain_type", "l3ExtDomain"),
					resource.TestMatchResourceAttr(l3Addr, "id", regexp.MustCompile(regexp.QuoteMeta("/domainAssociations/"+l3Dn)+"$")),

					// fibreChannelDomain via domain_name + domain_type
					resource.TestCheckResourceAttr(fcAddr, "domain_name", msoSiteAnpEpgDomainNameFc),
					resource.TestCheckResourceAttr(fcAddr, "domain_type", "fibreChannelDomain"),
					resource.TestMatchResourceAttr(fcAddr, "id", regexp.MustCompile(regexp.QuoteMeta("/domainAssociations/"+fcDn)+"$")),

					// vmmDomain via domain_name + domain_type + vmm_domain_type, full attribute matrix
					resource.TestCheckResourceAttr(vmmNameAddr, "domain_name", msoSiteAnpEpgDomainNameVmmName),
					resource.TestCheckResourceAttr(vmmNameAddr, "domain_type", "vmmDomain"),
					resource.TestCheckResourceAttr(vmmNameAddr, "vmm_domain_type", "VMware"),
					resource.TestCheckResourceAttr(vmmNameAddr, "deploy_immediacy", "lazy"),
					resource.TestCheckResourceAttr(vmmNameAddr, "resolution_immediacy", "immediate"),
					resource.TestCheckResourceAttr(vmmNameAddr, "vlan_encap_mode", "static"),
					resource.TestCheckResourceAttr(vmmNameAddr, "allow_micro_segmentation", "true"),
					resource.TestCheckResourceAttr(vmmNameAddr, "switching_mode", "native"),
					resource.TestCheckResourceAttr(vmmNameAddr, "switch_type", "default"),
					resource.TestCheckResourceAttr(vmmNameAddr, "micro_seg_vlan_type", "vlan"),
					resource.TestCheckResourceAttr(vmmNameAddr, "micro_seg_vlan", "46"),
					resource.TestCheckResourceAttr(vmmNameAddr, "port_encap_vlan_type", "vlan"),
					resource.TestCheckResourceAttr(vmmNameAddr, "port_encap_vlan", "45"),
					resource.TestCheckResourceAttr(vmmNameAddr, "enhanced_lag_policy_name", "Lacp"),
					resource.TestCheckResourceAttr(vmmNameAddr, "enhanced_lag_policy_dn", vmmNameDn+"/vswitchpolcont/enlacplagp-Lacp"),
					resource.TestCheckResourceAttr(vmmNameAddr, "delimiter", "|"),
					resource.TestCheckResourceAttr(vmmNameAddr, "binding_type", "static"),
					resource.TestCheckResourceAttr(vmmNameAddr, "port_allocation", "fixed"),
					resource.TestCheckResourceAttr(vmmNameAddr, "num_ports", "3"),
					resource.TestCheckResourceAttr(vmmNameAddr, "netflow", "disabled"),
					resource.TestCheckResourceAttr(vmmNameAddr, "allow_promiscuous", "accept"),
					resource.TestCheckResourceAttr(vmmNameAddr, "mac_changes", "reject"),
					resource.TestCheckResourceAttr(vmmNameAddr, "forged_transmits", "reject"),
					resource.TestCheckResourceAttr(vmmNameAddr, "custom_epg_name", "custom_epg_name_1"),
					resource.TestMatchResourceAttr(vmmNameAddr, "id", regexp.MustCompile(regexp.QuoteMeta("/domainAssociations/"+vmmNameDn)+"$")),

					// vmmDomain via domain_dn.
					// domain_type omitted for the same ConflictsWith reason as
					// the phys block above; the datasource test covers it.
					resource.TestCheckResourceAttr(vmmDnAddr, "domain_dn", vmmDnDn),
					resource.TestCheckResourceAttr(vmmDnAddr, "allow_micro_segmentation", "false"),
					resource.TestCheckResourceAttr(vmmDnAddr, "switching_mode", "native"),
					resource.TestCheckResourceAttr(vmmDnAddr, "switch_type", "default"),
					resource.TestCheckResourceAttr(vmmDnAddr, "vlan_encap_mode", "dynamic"),
					resource.TestMatchResourceAttr(vmmDnAddr, "id", regexp.MustCompile(regexp.QuoteMeta("/domainAssociations/"+vmmDnDn)+"$")),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update mutable attributes on phys and vmm domains") },
				Config:    testAccMSOSchemaSiteAnpEpgDomainConfigAll(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					// phys: only deploy_immediacy / resolution_immediacy are mutable
					resource.TestCheckResourceAttr(physAddr, "deploy_immediacy", "immediate"),
					resource.TestCheckResourceAttr(physAddr, "resolution_immediacy", "immediate"),

					// vmm_by_name: flipped VMM attribute block.
					// allow_micro_segmentation / micro_seg_vlan{,_type} are
					// intentionally NOT flipped here. They are Optional+Computed
					// in the schema, so dropping them from the updated config
					// preserves the prior state values, and the resource Update
					// would still send the previous micro_seg_vlan to MSO,
					// which then rejects the combo
					// (allow_micro_segmentation=false + microSegVlan set).
					resource.TestCheckResourceAttr(vmmNameAddr, "port_encap_vlan", "201"),
					resource.TestCheckResourceAttr(vmmNameAddr, "netflow", "enabled"),
					resource.TestCheckResourceAttr(vmmNameAddr, "allow_promiscuous", "reject"),
					resource.TestCheckResourceAttr(vmmNameAddr, "mac_changes", "accept"),
					resource.TestCheckResourceAttr(vmmNameAddr, "forged_transmits", "accept"),
					resource.TestCheckResourceAttr(vmmNameAddr, "num_ports", "5"),
				),
			},
			testAccMSOSchemaSiteAnpEpgDomainImportStep(physAddr),
			testAccMSOSchemaSiteAnpEpgDomainImportStep(l2Addr),
			testAccMSOSchemaSiteAnpEpgDomainImportStep(l3Addr),
			testAccMSOSchemaSiteAnpEpgDomainImportStep(fcAddr),
			testAccMSOSchemaSiteAnpEpgDomainImportStep(vmmNameAddr),
			testAccMSOSchemaSiteAnpEpgDomainImportStep(vmmDnAddr),
		},
	})
}

// testAccMSOSchemaSiteAnpEpgDomainImportStep builds an Import test step for
// the supplied resource address. The Id is reconstructed from the resource
// attributes already in state.
//
// Imported state diverges from applied state on the following attributes
// (covered by ImportStateVerifyIgnore):
//   - dn: deprecated alternative input, never set by Import.
//   - domain_name: the Import path extracts it from the DN's trailing
//     segment regardless of domain_type. For domains created via
//     `domain_dn` the apply state leaves domain_name unset, but Import
//     populates it - producing a benign diff.
//   - vmm_domain_type: same reason - Import always sets it from `vmmp-X/`
//     when present, even if apply was driven through `domain_dn`.
//   - domain_type: Import always sets it from the API response, but the
//     resource Read deliberately skips it when domain_dn was the input
//     (the two attributes ConflictsWith each other in the schema), so the
//     applied state lacks it for the phys / vmm_by_dn resources.
func testAccMSOSchemaSiteAnpEpgDomainImportStep(resourceAddr string) resource.TestStep {
	return resource.TestStep{
		PreConfig:    func() { fmt.Printf("Test: Import site ANP EPG domain %s\n", resourceAddr) },
		ResourceName: resourceAddr,
		ImportState:  true,
		ImportStateIdFunc: func(s *terraform.State) (string, error) {
			rs, ok := s.RootModule().Resources[resourceAddr]
			if !ok {
				return "", fmt.Errorf("resource not found in state: %s", resourceAddr)
			}
			return fmt.Sprintf("%s/sites/%s-%s/anps/%s/epgs/%s/domainAssociations/%s",
				rs.Primary.Attributes["schema_id"],
				rs.Primary.Attributes["site_id"],
				rs.Primary.Attributes["template_name"],
				rs.Primary.Attributes["anp_name"],
				rs.Primary.Attributes["epg_name"],
				rs.Primary.Attributes["domain_dn"],
			), nil
		},
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{"dn", "domain_name", "vmm_domain_type", "domain_type"},
	}
}

// testAccMSOSchemaSiteAnpEpgDomainBadBlock wraps a domain resource block
// around the shared site-ANP-EPG prerequisite config.
func testAccMSOSchemaSiteAnpEpgDomainBadBlock(badBlock string) string {
	return fmt.Sprintf(`%s
%s`, testAccMSOSchemaSiteAnpEpgConfigCreate(), badBlock)
}

func testAccMSOSchemaSiteAnpEpgDomainConfigMissingDn() string {
	return testAccMSOSchemaSiteAnpEpgDomainBadBlock(fmt.Sprintf(`
resource "mso_schema_site_anp_epg_domain" "bad" {
	schema_id            = mso_schema_site_anp_epg.%[1]s.schema_id
	site_id              = mso_schema_site_anp_epg.%[1]s.site_id
	template_name        = mso_schema_site_anp_epg.%[1]s.template_name
	anp_name             = mso_schema_site_anp_epg.%[1]s.anp_name
	epg_name             = mso_schema_site_anp_epg.%[1]s.epg_name
	deploy_immediacy     = "lazy"
	resolution_immediacy = "lazy"
}`, msoSchemaTemplateAnpEpgName))
}

func testAccMSOSchemaSiteAnpEpgDomainConfigNameWithoutType() string {
	return testAccMSOSchemaSiteAnpEpgDomainBadBlock(fmt.Sprintf(`
resource "mso_schema_site_anp_epg_domain" "bad" {
	schema_id            = mso_schema_site_anp_epg.%[1]s.schema_id
	site_id              = mso_schema_site_anp_epg.%[1]s.site_id
	template_name        = mso_schema_site_anp_epg.%[1]s.template_name
	anp_name             = mso_schema_site_anp_epg.%[1]s.anp_name
	epg_name             = mso_schema_site_anp_epg.%[1]s.epg_name
	domain_name          = "no_type"
	deploy_immediacy     = "lazy"
	resolution_immediacy = "lazy"
}`, msoSchemaTemplateAnpEpgName))
}

func testAccMSOSchemaSiteAnpEpgDomainConfigVmmWithoutVmmType() string {
	return testAccMSOSchemaSiteAnpEpgDomainBadBlock(fmt.Sprintf(`
resource "mso_schema_site_anp_epg_domain" "bad" {
	schema_id            = mso_schema_site_anp_epg.%[1]s.schema_id
	site_id              = mso_schema_site_anp_epg.%[1]s.site_id
	template_name        = mso_schema_site_anp_epg.%[1]s.template_name
	anp_name             = mso_schema_site_anp_epg.%[1]s.anp_name
	epg_name             = mso_schema_site_anp_epg.%[1]s.epg_name
	domain_name          = "no_vmm_type"
	domain_type          = "vmmDomain"
	deploy_immediacy     = "lazy"
	resolution_immediacy = "lazy"
}`, msoSchemaTemplateAnpEpgName))
}

func testAccMSOSchemaSiteAnpEpgDomainConfigConflicts() string {
	return testAccMSOSchemaSiteAnpEpgDomainBadBlock(fmt.Sprintf(`
resource "mso_schema_site_anp_epg_domain" "bad" {
	schema_id            = mso_schema_site_anp_epg.%[1]s.schema_id
	site_id              = mso_schema_site_anp_epg.%[1]s.site_id
	template_name        = mso_schema_site_anp_epg.%[1]s.template_name
	anp_name             = mso_schema_site_anp_epg.%[1]s.anp_name
	epg_name             = mso_schema_site_anp_epg.%[1]s.epg_name
	domain_name          = "conflicting"
	domain_type          = "physicalDomain"
	domain_dn            = "uni/phys-conflicting"
	deploy_immediacy     = "lazy"
	resolution_immediacy = "lazy"
}`, msoSchemaTemplateAnpEpgName))
}

// testAccMSOSchemaSiteAnpEpgDomainConfigAll renders the prerequisite config
// plus one mso_schema_site_anp_epg_domain block per domain_type covered.
// When updated is false, the VMM domain is created with the "pre-update"
// matrix; when true, the mutable attributes are flipped to exercise the
// resource Update path.
func testAccMSOSchemaSiteAnpEpgDomainConfigAll(updated bool) string {
	physDeploy := "lazy"
	physResolution := "lazy"
	// allow_micro_segmentation / micro_seg_vlan{,_type} are NOT flipped
	// between the pre- and post-update matrices. They are Optional+Computed
	// in the resource schema, so removing them from the updated HCL would
	// leave the prior state values in place and the Update path would still
	// send the previous microSegVlan to MSO, which rejects that combo
	// whenever allow_micro_segmentation is false
	// ("Microsegmentation not specified, microSegVlan must not be specified").
	vmmPortEncapVlan := "45"
	vmmNetflow := "disabled"
	vmmAllowPromiscuous := "accept"
	vmmMacChanges := "reject"
	vmmForgedTransmits := "reject"
	vmmNumPorts := "3"
	if updated {
		physDeploy = "immediate"
		physResolution = "immediate"
		vmmPortEncapVlan = "201"
		vmmNetflow = "enabled"
		vmmAllowPromiscuous = "reject"
		vmmMacChanges = "accept"
		vmmForgedTransmits = "accept"
		vmmNumPorts = "5"
	}

	return fmt.Sprintf(`%[1]s
resource "mso_schema_site_anp_epg_domain" "%[2]s" {
	schema_id            = mso_schema_site_anp_epg.%[3]s.schema_id
	site_id              = mso_schema_site_anp_epg.%[3]s.site_id
	template_name        = mso_schema_site_anp_epg.%[3]s.template_name
	anp_name             = mso_schema_site_anp_epg.%[3]s.anp_name
	epg_name             = mso_schema_site_anp_epg.%[3]s.epg_name
	domain_dn            = "uni/phys-%[4]s"
	deploy_immediacy     = "%[5]s"
	resolution_immediacy = "%[6]s"
}

resource "mso_schema_site_anp_epg_domain" "%[7]s" {
	schema_id            = mso_schema_site_anp_epg.%[3]s.schema_id
	site_id              = mso_schema_site_anp_epg.%[3]s.site_id
	template_name        = mso_schema_site_anp_epg.%[3]s.template_name
	anp_name             = mso_schema_site_anp_epg.%[3]s.anp_name
	epg_name             = mso_schema_site_anp_epg.%[3]s.epg_name
	domain_name          = "%[8]s"
	domain_type          = "l2ExtDomain"
	deploy_immediacy     = "lazy"
	resolution_immediacy = "lazy"

	# Serialize against the previous mso_schema_site_anp_epg_domain so the
	# concurrent writers don't race on domainAssociations[] indices.
	depends_on = [mso_schema_site_anp_epg_domain.%[2]s]
}

resource "mso_schema_site_anp_epg_domain" "%[9]s" {
	schema_id            = mso_schema_site_anp_epg.%[3]s.schema_id
	site_id              = mso_schema_site_anp_epg.%[3]s.site_id
	template_name        = mso_schema_site_anp_epg.%[3]s.template_name
	anp_name             = mso_schema_site_anp_epg.%[3]s.anp_name
	epg_name             = mso_schema_site_anp_epg.%[3]s.epg_name
	domain_name          = "%[10]s"
	domain_type          = "l3ExtDomain"
	deploy_immediacy     = "lazy"
	resolution_immediacy = "lazy"

	depends_on = [mso_schema_site_anp_epg_domain.%[7]s]
}

resource "mso_schema_site_anp_epg_domain" "%[11]s" {
	schema_id            = mso_schema_site_anp_epg.%[3]s.schema_id
	site_id              = mso_schema_site_anp_epg.%[3]s.site_id
	template_name        = mso_schema_site_anp_epg.%[3]s.template_name
	anp_name             = mso_schema_site_anp_epg.%[3]s.anp_name
	epg_name             = mso_schema_site_anp_epg.%[3]s.epg_name
	domain_name          = "%[12]s"
	domain_type          = "fibreChannelDomain"
	deploy_immediacy     = "lazy"
	resolution_immediacy = "lazy"

	depends_on = [mso_schema_site_anp_epg_domain.%[9]s]
}

resource "mso_schema_site_anp_epg_domain" "%[13]s" {
	schema_id                = mso_schema_site_anp_epg.%[3]s.schema_id
	site_id                  = mso_schema_site_anp_epg.%[3]s.site_id
	template_name            = mso_schema_site_anp_epg.%[3]s.template_name
	anp_name                 = mso_schema_site_anp_epg.%[3]s.anp_name
	epg_name                 = mso_schema_site_anp_epg.%[3]s.epg_name
	domain_name              = "%[14]s"
	domain_type              = "vmmDomain"
	vmm_domain_type          = "VMware"
	deploy_immediacy         = "lazy"
	resolution_immediacy     = "immediate"
	vlan_encap_mode          = "static"
	allow_micro_segmentation = true
	switching_mode           = "native"
	switch_type              = "default"
	micro_seg_vlan_type      = "vlan"
	micro_seg_vlan           = 46
	port_encap_vlan_type     = "vlan"
	port_encap_vlan          = %[15]s
	enhanced_lag_policy_name = "Lacp"
	enhanced_lag_policy_dn   = "uni/vmmp-VMware/dom-%[14]s/vswitchpolcont/enlacplagp-Lacp"
	delimiter                = "|"
	binding_type             = "static"
	port_allocation          = "fixed"
	num_ports                = %[16]s
	netflow                  = "%[17]s"
	allow_promiscuous        = "%[18]s"
	mac_changes              = "%[19]s"
	forged_transmits         = "%[20]s"
	custom_epg_name          = "custom_epg_name_1"

	depends_on = [mso_schema_site_anp_epg_domain.%[11]s]
}

resource "mso_schema_site_anp_epg_domain" "%[21]s" {
	schema_id            = mso_schema_site_anp_epg.%[3]s.schema_id
	site_id              = mso_schema_site_anp_epg.%[3]s.site_id
	template_name        = mso_schema_site_anp_epg.%[3]s.template_name
	anp_name             = mso_schema_site_anp_epg.%[3]s.anp_name
	epg_name             = mso_schema_site_anp_epg.%[3]s.epg_name
	domain_dn            = "uni/vmmp-VMware/dom-%[22]s"
	binding_type         = "static"
	port_allocation      = "fixed"
	netflow              = "disabled"
	allow_promiscuous    = "accept"
	mac_changes          = "reject"
	forged_transmits     = "reject"
	deploy_immediacy     = "lazy"
	resolution_immediacy = "lazy"

	depends_on = [mso_schema_site_anp_epg_domain.%[13]s]
}
`,
		testAccMSOSchemaSiteAnpEpgConfigCreate(), // 1
		msoSiteAnpEpgDomainLabelPhys,             // 2
		msoSchemaTemplateAnpEpgName,              // 3
		msoSiteAnpEpgDomainNamePhys,              // 4
		physDeploy,                               // 5
		physResolution,                           // 6
		msoSiteAnpEpgDomainLabelL2,               // 7
		msoSiteAnpEpgDomainNameL2,                // 8
		msoSiteAnpEpgDomainLabelL3,               // 9
		msoSiteAnpEpgDomainNameL3,                // 10
		msoSiteAnpEpgDomainLabelFc,               // 11
		msoSiteAnpEpgDomainNameFc,                // 12
		msoSiteAnpEpgDomainLabelVmmName,          // 13
		msoSiteAnpEpgDomainNameVmmName,           // 14
		vmmPortEncapVlan,                         // 15
		vmmNumPorts,                              // 16
		vmmNetflow,                               // 17
		vmmAllowPromiscuous,                      // 18
		vmmMacChanges,                            // 19
		vmmForgedTransmits,                       // 20
		msoSiteAnpEpgDomainLabelVmmDn,            // 21
		msoSiteAnpEpgDomainNameVmmDn,             // 22
	)
}

// testAccCheckMSOSchemaSiteAnpEpgDomainDestroy walks state for any
// mso_schema_site_anp_epg_domain resources, fetches the schema, and asserts
// no domainAssociations[] entry with a matching dn remains under the parent
// site/anp/epg. A missing schema or sites array is treated as a successful
// destroy.
func testAccCheckMSOSchemaSiteAnpEpgDomainDestroy(s *terraform.State) error {
	msoClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mso_schema_site_anp_epg_domain" {
			continue
		}
		schemaId := rs.Primary.Attributes["schema_id"]
		stateSiteId := rs.Primary.Attributes["site_id"]
		stateTemplate := rs.Primary.Attributes["template_name"]
		stateAnp := rs.Primary.Attributes["anp_name"]
		stateEpg := rs.Primary.Attributes["epg_name"]
		stateDn := rs.Primary.Attributes["domain_dn"]

		cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
		if err != nil {
			return nil
		}
		count, err := cont.ArrayCount("sites")
		if err != nil {
			return nil
		}
		for i := 0; i < count; i++ {
			siteCont, err := cont.ArrayElement(i, "sites")
			if err != nil {
				return err
			}
			apiSiteId := models.StripQuotes(siteCont.S("siteId").String())
			apiTemplate := models.StripQuotes(siteCont.S("templateName").String())
			if apiSiteId != stateSiteId || apiTemplate != stateTemplate {
				continue
			}
			anpCount, err := siteCont.ArrayCount("anps")
			if err != nil {
				continue
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := siteCont.ArrayElement(j, "anps")
				if err != nil {
					return err
				}
				anpRef := models.StripQuotes(anpCont.S("anpRef").String())
				anpSplit := strings.Split(anpRef, "/")
				if len(anpSplit) < 7 || anpSplit[6] != stateAnp {
					continue
				}
				epgCount, err := anpCont.ArrayCount("epgs")
				if err != nil {
					continue
				}
				for k := 0; k < epgCount; k++ {
					epgCont, err := anpCont.ArrayElement(k, "epgs")
					if err != nil {
						return err
					}
					epgRef := models.StripQuotes(epgCont.S("epgRef").String())
					epgSplit := strings.Split(epgRef, "/")
					if len(epgSplit) < 9 || epgSplit[8] != stateEpg {
						continue
					}
					domainCount, err := epgCont.ArrayCount("domainAssociations")
					if err != nil {
						continue
					}
					for l := 0; l < domainCount; l++ {
						domainCont, err := epgCont.ArrayElement(l, "domainAssociations")
						if err != nil {
							return err
						}
						apiDomain := models.StripQuotes(domainCont.S("dn").String())
						if apiDomain == stateDn {
							return fmt.Errorf("mso_schema_site_anp_epg_domain (site=%s, template=%s, anp=%s, epg=%s, dn=%s) still exists on schema %s", stateSiteId, stateTemplate, stateAnp, stateEpg, stateDn, schemaId)
						}
					}
				}
			}
		}
	}
	return nil
}
