package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Note: NDO uppercases the `value` server-side for several useg_types (`vm-name`, `dns`,
// `hv`, `guest-os`, `vnic`), so the tests use uppercase values (e.g. "TEST-VM") to avoid
// a perpetual plan diff between the configured value and the value read back from the API.

// msoSchemaTemplateAnpEpgUsegAttrSchemaId is set during the first test step's Check to capture
// the dynamic schema ID for use in the manual deletion PreConfig step.
var msoSchemaTemplateAnpEpgUsegAttrSchemaId string

func TestAccMSOSchemaTemplateAnpEpgUsegAttrResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateAnpEpgUsegAttrDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create UsegAttr with required attributes only (ip type)") },
				Config:    testAccMSOSchemaTemplateAnpEpgUsegAttrConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName, "anp_name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName, "epg_name", msoSchemaTemplateAnpEpgName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName, "name", msoSchemaTemplateAnpEpgUsegAttrName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName, "useg_type", "ip"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName, "value", msoSchemaTemplateAnpEpgUsegAttrIp),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName, "operator", ""),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName, "category", ""),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName, "description", ""),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName, "useg_subnet", "false"),
					// Capture the dynamic schema ID for use in the manual deletion PreConfig step
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName]
						if !ok {
							return fmt.Errorf("UsegAttr resource not found in state")
						}
						msoSchemaTemplateAnpEpgUsegAttrSchemaId = rs.Primary.Attributes["schema_id"]
						return nil
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Operator round-trips when explicitly set on ip useg_type") },
				Config:    testAccMSOSchemaTemplateAnpEpgUsegAttrConfigOperatorOverride(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName, "useg_type", "ip"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName, "operator", "startsWith"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update UsegAttr all attributes (ip type)") },
				Config:    testAccMSOSchemaTemplateAnpEpgUsegAttrConfigUpdateAll(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName, "useg_type", "ip"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName, "value", msoSchemaTemplateAnpEpgUsegAttrIp),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName, "description", "test useg"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName, "category", "test_category"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName, "useg_subnet", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName, "operator", ""),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Switch UsegAttr to vm-name (operator preserved)") },
				Config:    testAccMSOSchemaTemplateAnpEpgUsegAttrConfigSwitchToVmName(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName, "useg_type", "vm-name"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName, "value", "TEST-VM"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName, "operator", "contains"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Reset UsegAttr optional attributes") },
				Config:    testAccMSOSchemaTemplateAnpEpgUsegAttrConfigReset(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName, "useg_type", "vm-name"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName, "value", "TEST-VM"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName, "operator", "equals"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName, "description", ""),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName, "category", ""),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: ForceNew UsegAttr by changing name") },
				Config:    testAccMSOSchemaTemplateAnpEpgUsegAttrConfigNewName(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName2, "name", msoSchemaTemplateAnpEpgUsegAttrName2),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName2, "useg_type", "vm-name"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName2, "value", "TEST-VM"),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import UsegAttr") },
				ResourceName: "mso_schema_template_anp_epg_useg_attr." + msoSchemaTemplateAnpEpgUsegAttrName2,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName2]
					if !ok {
						return "", fmt.Errorf("UsegAttr resource not found in state")
					}
					return fmt.Sprintf("%s/templates/%s/anps/%s/epgs/%s/uSegAttrs/%s",
						rs.Primary.Attributes["schema_id"],
						rs.Primary.Attributes["template_name"],
						rs.Primary.Attributes["anp_name"],
						rs.Primary.Attributes["epg_name"],
						rs.Primary.Attributes["name"]), nil
				},
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate UsegAttr after manual deletion from NDO")
					msoClient := testAccProvider.Meta().(*client.Client)
					path := fmt.Sprintf("/templates/%s/anps/%s/epgs/%s/uSegAttrs/%s",
						msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateAnpEpgUsegAttrName2)
					removePayload := models.GetRemovePatchPayload(path)
					_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", msoSchemaTemplateAnpEpgUsegAttrSchemaId), removePayload)
					if err != nil {
						t.Fatalf("Failed to manually delete useg_attr: %v", err)
					}
				},
				Config: testAccMSOSchemaTemplateAnpEpgUsegAttrConfigNewName(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName2, "name", msoSchemaTemplateAnpEpgUsegAttrName2),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName2, "useg_type", "vm-name"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName2, "value", "TEST-VM"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update parent EPG description with UsegAttr present") },
				Config:    testAccMSOSchemaTemplateAnpEpgUsegAttrConfigUpdateParentEpg(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName2, "name", msoSchemaTemplateAnpEpgUsegAttrName2),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr."+msoSchemaTemplateAnpEpgUsegAttrName2, "useg_type", "vm-name"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "description", "Updated EPG description with useg_attr"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create assorted UsegAttr type combinations") },
				Config:    testAccMSOSchemaTemplateAnpEpgUsegAttrConfigTypes(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// ip + description only (no operator stored on NDO 4.x)
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.ip_with_description", "useg_type", "ip"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.ip_with_description", "value", "10.20.20.1"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.ip_with_description", "description", "fooo"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.ip_with_description", "operator", ""),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.ip_with_description", "category", ""),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.ip_with_description", "useg_subnet", "false"),

					// ip + useg_subnet=true
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.ip_subnet", "useg_type", "ip"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.ip_subnet", "value", "0.0.0.0"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.ip_subnet", "useg_subnet", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.ip_subnet", "operator", ""),

					// mac (no operator stored on NDO 4.x)
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.mac_basic", "useg_type", "mac"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.mac_basic", "value", "AA:BB:CC:DD:EE:FF"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.mac_basic", "operator", ""),

					// dns (no operator stored on NDO 4.x; value uppercased server-side)
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.dns_basic", "useg_type", "dns"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.dns_basic", "value", "*.EXAMPLE.COM"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.dns_basic", "operator", ""),

					// vm-name + operator=equals (NDO uppercases value server-side)
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.vm_name_equals", "useg_type", "vm-name"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.vm_name_equals", "value", "TYPES-VM"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.vm_name_equals", "operator", "equals"),

					// rootContName (VM data center) + operator=equals
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.root_cont_name_equals", "useg_type", "rootContName"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.root_cont_name_equals", "value", "DC1"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.root_cont_name_equals", "operator", "equals"),

					// hv (Hypervisor) + operator=contains (value uppercased server-side)
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.hypervisor_contains", "useg_type", "hv"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.hypervisor_contains", "value", "HOST-1"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.hypervisor_contains", "operator", "contains"),

					// guest-os (Operating System) + operator=startsWith (value uppercased server-side)
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.guest_os_starts_with", "useg_type", "guest-os"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.guest_os_starts_with", "value", "UBUNTU"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.guest_os_starts_with", "operator", "startsWith"),

					// tag + operator=contains + category
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.tag_with_category", "useg_type", "tag"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.tag_with_category", "value", "SOMETHIGN"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.tag_with_category", "operator", "contains"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.tag_with_category", "category", "foo"),

					// vm (Identifier) + operator=equals
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.vm_id_equals", "useg_type", "vm"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.vm_id_equals", "value", "SPECIAL_VM"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.vm_id_equals", "operator", "equals"),

					// domain (VMM domain) + operator=contains
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.domain_contains", "useg_type", "domain"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.domain_contains", "value", "SPECIAL_VM"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.domain_contains", "operator", "contains"),

					// vnic (Vnic DN) + operator=endsWith (value uppercased server-side)
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.vnic_ends_with", "useg_type", "vnic"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.vnic_ends_with", "value", "NETWORK ADAPTER 1"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_useg_attr.vnic_ends_with", "operator", "endsWith"),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateAnpEpgUsegAttrConfigTypes() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp_epg_useg_attr" "ip_with_description" {
		schema_id     = mso_schema_template_anp_epg.%[2]s.schema_id
		template_name = "%[3]s"
		anp_name      = "%[4]s"
		epg_name      = mso_schema_template_anp_epg.%[2]s.name
		name          = "ip_with_description"
		useg_type     = "ip"
		value         = "10.20.20.1"
		description   = "fooo"
	}
	resource "mso_schema_template_anp_epg_useg_attr" "ip_subnet" {
		schema_id     = mso_schema_template_anp_epg.%[2]s.schema_id
		template_name = "%[3]s"
		anp_name      = "%[4]s"
		epg_name      = mso_schema_template_anp_epg.%[2]s.name
		name          = "ip_subnet"
		useg_type     = "ip"
		value         = "0.0.0.0"
		useg_subnet   = true
	}
	resource "mso_schema_template_anp_epg_useg_attr" "mac_basic" {
		schema_id     = mso_schema_template_anp_epg.%[2]s.schema_id
		template_name = "%[3]s"
		anp_name      = "%[4]s"
		epg_name      = mso_schema_template_anp_epg.%[2]s.name
		name          = "mac_basic"
		useg_type     = "mac"
		value         = "AA:BB:CC:DD:EE:FF"
	}
	resource "mso_schema_template_anp_epg_useg_attr" "dns_basic" {
		schema_id     = mso_schema_template_anp_epg.%[2]s.schema_id
		template_name = "%[3]s"
		anp_name      = "%[4]s"
		epg_name      = mso_schema_template_anp_epg.%[2]s.name
		name          = "dns_basic"
		useg_type     = "dns"
		value         = "*.EXAMPLE.COM"
	}
	resource "mso_schema_template_anp_epg_useg_attr" "vm_name_equals" {
		schema_id     = mso_schema_template_anp_epg.%[2]s.schema_id
		template_name = "%[3]s"
		anp_name      = "%[4]s"
		epg_name      = mso_schema_template_anp_epg.%[2]s.name
		name          = "vm_name_equals"
		useg_type     = "vm-name"
		operator      = "equals"
		value         = "TYPES-VM"
	}
	resource "mso_schema_template_anp_epg_useg_attr" "root_cont_name_equals" {
		schema_id     = mso_schema_template_anp_epg.%[2]s.schema_id
		template_name = "%[3]s"
		anp_name      = "%[4]s"
		epg_name      = mso_schema_template_anp_epg.%[2]s.name
		name          = "root_cont_name_equals"
		useg_type     = "rootContName"
		operator      = "equals"
		value         = "DC1"
	}
	resource "mso_schema_template_anp_epg_useg_attr" "hypervisor_contains" {
		schema_id     = mso_schema_template_anp_epg.%[2]s.schema_id
		template_name = "%[3]s"
		anp_name      = "%[4]s"
		epg_name      = mso_schema_template_anp_epg.%[2]s.name
		name          = "hypervisor_contains"
		useg_type     = "hv"
		operator      = "contains"
		value         = "HOST-1"
	}
	resource "mso_schema_template_anp_epg_useg_attr" "guest_os_starts_with" {
		schema_id     = mso_schema_template_anp_epg.%[2]s.schema_id
		template_name = "%[3]s"
		anp_name      = "%[4]s"
		epg_name      = mso_schema_template_anp_epg.%[2]s.name
		name          = "guest_os_starts_with"
		useg_type     = "guest-os"
		operator      = "startsWith"
		value         = "UBUNTU"
	}
	resource "mso_schema_template_anp_epg_useg_attr" "tag_with_category" {
		schema_id     = mso_schema_template_anp_epg.%[2]s.schema_id
		template_name = "%[3]s"
		anp_name      = "%[4]s"
		epg_name      = mso_schema_template_anp_epg.%[2]s.name
		name          = "tag_with_category"
		useg_type     = "tag"
		operator      = "contains"
		category      = "foo"
		value         = "SOMETHIGN"
	}
	resource "mso_schema_template_anp_epg_useg_attr" "vm_id_equals" {
		schema_id     = mso_schema_template_anp_epg.%[2]s.schema_id
		template_name = "%[3]s"
		anp_name      = "%[4]s"
		epg_name      = mso_schema_template_anp_epg.%[2]s.name
		name          = "vm_id_equals"
		useg_type     = "vm"
		operator      = "equals"
		value         = "SPECIAL_VM"
	}
	resource "mso_schema_template_anp_epg_useg_attr" "domain_contains" {
		schema_id     = mso_schema_template_anp_epg.%[2]s.schema_id
		template_name = "%[3]s"
		anp_name      = "%[4]s"
		epg_name      = mso_schema_template_anp_epg.%[2]s.name
		name          = "domain_contains"
		useg_type     = "domain"
		operator      = "contains"
		value         = "SPECIAL_VM"
	}
	resource "mso_schema_template_anp_epg_useg_attr" "vnic_ends_with" {
		schema_id     = mso_schema_template_anp_epg.%[2]s.schema_id
		template_name = "%[3]s"
		anp_name      = "%[4]s"
		epg_name      = mso_schema_template_anp_epg.%[2]s.name
		name          = "vnic_ends_with"
		useg_type     = "vnic"
		operator      = "endsWith"
		value         = "NETWORK ADAPTER 1"
	}`, testAccMSOSchemaTemplateAnpEpgUsegAttrConfigUpdateParentEpg(), msoSchemaTemplateAnpEpgName, msoSchemaTemplateName, msoSchemaTemplateAnpName)
}

func testAccMSOSchemaTemplateAnpEpgUsegAttrPrerequisiteConfig() string {
	return fmt.Sprintf(`%s%s%s%s%s%s`, testSiteConfigAnsibleTest(), testTenantConfig(), testSchemaConfig(), testSchemaTemplateVrfConfig(), testSchemaTemplateBdConfig(), testSchemaTemplateAnpConfig()) + fmt.Sprintf(`
resource "mso_schema_template_anp_epg" "%[1]s" {
	name          = "%[1]s"
	display_name  = "%[1]s"
	anp_name      = mso_schema_template_anp.%[2]s.name
	schema_id     = mso_schema.%[3]s.id
	template_name = "%[4]s"
	bd_name       = mso_schema_template_bd.%[5]s.name
	useg_epg      = true
}
`, msoSchemaTemplateAnpEpgName, msoSchemaTemplateAnpName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateBdName)
}

func testAccMSOSchemaTemplateAnpEpgUsegAttrConfigCreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp_epg_useg_attr" "%[2]s" {
		schema_id     = mso_schema_template_anp_epg.%[7]s.schema_id
		template_name = "%[4]s"
		anp_name      = "%[5]s"
		epg_name      = mso_schema_template_anp_epg.%[7]s.name
		name          = "%[2]s"
		useg_type     = "ip"
		value         = "%[8]s"
	}`, testAccMSOSchemaTemplateAnpEpgUsegAttrPrerequisiteConfig(), msoSchemaTemplateAnpEpgUsegAttrName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateAnpEpgUsegAttrIp)
}

func testAccMSOSchemaTemplateAnpEpgUsegAttrConfigOperatorOverride() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp_epg_useg_attr" "%[2]s" {
		schema_id     = mso_schema_template_anp_epg.%[7]s.schema_id
		template_name = "%[4]s"
		anp_name      = "%[5]s"
		epg_name      = mso_schema_template_anp_epg.%[7]s.name
		name          = "%[2]s"
		useg_type     = "ip"
		value         = "%[8]s"
		operator      = "startsWith"
	}`, testAccMSOSchemaTemplateAnpEpgUsegAttrPrerequisiteConfig(), msoSchemaTemplateAnpEpgUsegAttrName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateAnpEpgUsegAttrIp)
}

func testAccMSOSchemaTemplateAnpEpgUsegAttrConfigUpdateAll() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp_epg_useg_attr" "%[2]s" {
		schema_id     = mso_schema_template_anp_epg.%[7]s.schema_id
		template_name = "%[4]s"
		anp_name      = "%[5]s"
		epg_name      = mso_schema_template_anp_epg.%[7]s.name
		name          = "%[2]s"
		useg_type     = "ip"
		value         = "%[8]s"
		description   = "test useg"
		category      = "test_category"
		useg_subnet   = true
	}`, testAccMSOSchemaTemplateAnpEpgUsegAttrPrerequisiteConfig(), msoSchemaTemplateAnpEpgUsegAttrName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateAnpEpgUsegAttrIp)
}

func testAccMSOSchemaTemplateAnpEpgUsegAttrConfigSwitchToVmName() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp_epg_useg_attr" "%[2]s" {
		schema_id     = mso_schema_template_anp_epg.%[7]s.schema_id
		template_name = "%[4]s"
		anp_name      = "%[5]s"
		epg_name      = mso_schema_template_anp_epg.%[7]s.name
		name          = "%[2]s"
		useg_type     = "vm-name"
		value         = "TEST-VM"
		operator      = "contains"
		description   = "test useg"
		category      = "test_category"
	}`, testAccMSOSchemaTemplateAnpEpgUsegAttrPrerequisiteConfig(), msoSchemaTemplateAnpEpgUsegAttrName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateAnpEpgName)
}

func testAccMSOSchemaTemplateAnpEpgUsegAttrConfigReset() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp_epg_useg_attr" "%[2]s" {
		schema_id     = mso_schema_template_anp_epg.%[7]s.schema_id
		template_name = "%[4]s"
		anp_name      = "%[5]s"
		epg_name      = mso_schema_template_anp_epg.%[7]s.name
		name          = "%[2]s"
		useg_type     = "vm-name"
		value         = "TEST-VM"
		operator      = "equals"
	}`, testAccMSOSchemaTemplateAnpEpgUsegAttrPrerequisiteConfig(), msoSchemaTemplateAnpEpgUsegAttrName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateAnpEpgName)
}

func testAccMSOSchemaTemplateAnpEpgUsegAttrConfigNewName() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp_epg_useg_attr" "%[2]s" {
		schema_id     = mso_schema_template_anp_epg.%[7]s.schema_id
		template_name = "%[4]s"
		anp_name      = "%[5]s"
		epg_name      = mso_schema_template_anp_epg.%[7]s.name
		name          = "%[2]s"
		useg_type     = "vm-name"
		value         = "TEST-VM"
		operator      = "equals"
	}`, testAccMSOSchemaTemplateAnpEpgUsegAttrPrerequisiteConfig(), msoSchemaTemplateAnpEpgUsegAttrName2, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateAnpEpgName)
}

func testAccMSOSchemaTemplateAnpEpgUsegAttrConfigUpdateParentEpg() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp_epg" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		template_name = "%[4]s"
		anp_name      = mso_schema_template_anp.%[5]s.name
		name          = "%[2]s"
		display_name  = "%[2]s"
		description   = "Updated EPG description with useg_attr"
		bd_name       = mso_schema_template_bd.%[6]s.name
		useg_epg      = true
	}
	resource "mso_schema_template_anp_epg_useg_attr" "%[7]s" {
		schema_id     = mso_schema_template_anp_epg.%[2]s.schema_id
		template_name = "%[4]s"
		anp_name      = "%[5]s"
		epg_name      = mso_schema_template_anp_epg.%[2]s.name
		name          = "%[7]s"
		useg_type     = "vm-name"
		value         = "TEST-VM"
		operator      = "equals"
	}`, testSiteConfigAnsibleTest()+testTenantConfig()+testSchemaConfig()+testSchemaTemplateVrfConfig()+testSchemaTemplateBdConfig()+testSchemaTemplateAnpConfig(), msoSchemaTemplateAnpEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateBdName, msoSchemaTemplateAnpEpgUsegAttrName2)
}

func testAccCheckMSOSchemaTemplateAnpEpgUsegAttrDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "mso_schema_template_anp_epg_useg_attr" {
			schemaID := rs.Primary.Attributes["schema_id"]
			cont, err := client.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
			if err != nil {
				return nil
			}
			count, err := cont.ArrayCount("templates")
			if err != nil {
				return fmt.Errorf("No Template found")
			}
			for i := 0; i < count; i++ {
				tempCont, err := cont.ArrayElement(i, "templates")
				if err != nil {
					return fmt.Errorf("No template exists")
				}
				anpCount, err := tempCont.ArrayCount("anps")
				if err != nil {
					return fmt.Errorf("No Anp found")
				}
				for j := 0; j < anpCount; j++ {
					anpCont, err := tempCont.ArrayElement(j, "anps")
					if err != nil {
						return err
					}
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return fmt.Errorf("Unable to get EPG list")
					}
					for k := 0; k < epgCount; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return err
						}
						usegCount, err := epgCont.ArrayCount("uSegAttrs")
						if err != nil {
							continue
						}
						for l := 0; l < usegCount; l++ {
							usegCont, err := epgCont.ArrayElement(l, "uSegAttrs")
							if err != nil {
								return err
							}
							name := models.StripQuotes(usegCont.S("name").String())
							if rs.Primary.Attributes["name"] == name {
								return fmt.Errorf("Schema Template Anp Epg UsegAttr record still exists")
							}
						}
					}
				}
			}
		}
	}
	return nil
}
