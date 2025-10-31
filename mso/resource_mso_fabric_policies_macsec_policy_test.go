package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOMacsecPolicyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create MACsec Policy") },
				Config:    testAccMSOMacsecPolicyConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "name", "tf_test_macsec_policy"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "description", "Terraform test MACsec Policy"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "admin_state", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "interface_type", "access"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "cipher_suite", "256GcmAes"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "window_size", "128"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "security_policy", "shouldSecure"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "sak_expire_time", "60"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "confidentiality_offset", "offset30"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "key_server_priority", "8"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "macsec_keys.#", "1"),
					customTestCheckResourceTypeSetAttr("mso_fabric_policies_macsec_policy.macsec_policy", "macsec_keys",
						map[string]string{
							"key_name":   "abc123",
							"psk":        "AA111111111111111111111111111111111111111111111111111111111111aa",
							"start_time": "2027-09-23 00:00:00",
							"end_time":   "2030-09-23 00:00:00",
						},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update MACsec Policy adding extra MACsec Key") },
				Config:    testAccMSOMacsecPolicyConfigUpdateAddingExtraMacSecKey(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "name", "tf_test_macsec_policy"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "description", "Terraform test MACsec Policy adding extra MACsec Key"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "admin_state", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "interface_type", "access"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "cipher_suite", "256GcmAes"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "window_size", "128"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "security_policy", "shouldSecure"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "sak_expire_time", "60"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "confidentiality_offset", "offset30"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "key_server_priority", "8"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "macsec_keys.#", "2"),
					customTestCheckResourceTypeSetAttr("mso_fabric_policies_macsec_policy.macsec_policy", "macsec_keys",
						map[string]string{
							"key_name":   "abc123",
							"psk":        "AA111111111111111111111111111111111111111111111111111111111111aa",
							"start_time": "2027-09-23 00:00:00",
							"end_time":   "2030-09-23 00:00:00",
						},
					),
					customTestCheckResourceTypeSetAttr("mso_fabric_policies_macsec_policy.macsec_policy", "macsec_keys",
						map[string]string{
							"key_name":   "def456",
							"psk":        "AA11111111111111111111111111111111111111111111111111111111111aaa",
							"start_time": "2029-12-11 11:12:13",
							"end_time":   "2030-12-11 11:12:13",
						},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update MACsec Policy removing extra MACsec Key") },
				Config:    testAccMSOMacsecPolicyConfigUpdateRemovingExtraMacSecKey(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "admin_state", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "interface_type", "access"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "cipher_suite", "256GcmAes"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "window_size", "128"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "security_policy", "shouldSecure"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "sak_expire_time", "60"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "confidentiality_offset", "offset30"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "key_server_priority", "8"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "macsec_keys.#", "1"),
					customTestCheckResourceTypeSetAttr("mso_fabric_policies_macsec_policy.macsec_policy", "macsec_keys",
						map[string]string{
							"key_name":   "abc123",
							"psk":        "AA111111111111111111111111111111111111111111111111111111111111aa",
							"start_time": "2027-09-23 00:00:00",
							"end_time":   "2030-09-23 00:00:00",
						},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update MACsec Policy changing the interface type") },
				Config:    testAccMSOMacsecPolicyConfigUpdateChangingType(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "admin_state", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "interface_type", "fabric"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "cipher_suite", "256GcmAesXpn"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "window_size", "256"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "security_policy", "mustSecure"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "sak_expire_time", "120"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "confidentiality_offset", "offset30"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "key_server_priority", "8"),
					resource.TestCheckResourceAttr("mso_fabric_policies_macsec_policy.macsec_policy", "macsec_keys.#", "1"),
					customTestCheckResourceTypeSetAttr("mso_fabric_policies_macsec_policy.macsec_policy", "macsec_keys",
						map[string]string{
							"key_name":   "abc123",
							"psk":        "AA111111111111111111111111111111111111111111111111111111111111ab",
							"start_time": "2028-12-12 12:12:12",
							"end_time":   "infinite",
						},
					),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import MACsec Policy") },
				ResourceName:      "mso_fabric_policies_macsec_policy.macsec_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments("mso_fabric_policies_macsec_policy", "fabricPolicyTemplate", "template", "macsecPolicies"),
	})
}

func testAccMSOMacsecPolicyConfigCreate() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_macsec_policy" "macsec_policy" {
		template_id            = mso_template.template_fabric_policy.id
		name                   = "tf_test_macsec_policy"
		description            = "Terraform test MACsec Policy"
		admin_state            = "enabled"
		interface_type         = "access"
		cipher_suite           = "256GcmAes"
		window_size            = 128
		security_policy        = "shouldSecure"
		sak_expire_time        = 60
		confidentiality_offset = "offset30"
		key_server_priority    = 8
		macsec_keys {
			key_name           = "abc123"
			psk                = "AA111111111111111111111111111111111111111111111111111111111111aa"
			start_time         = "2027-09-23 00:00:00"
			end_time           = "2030-09-23 00:00:00"
		}
	}`, testAccMSOTemplateResourceFabricPolicyConfig())
}

func testAccMSOMacsecPolicyConfigUpdateAddingExtraMacSecKey() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_macsec_policy" "macsec_policy" {
		template_id            = mso_template.template_fabric_policy.id
		name                   = "tf_test_macsec_policy"
		description            = "Terraform test MACsec Policy adding extra MACsec Key"
		admin_state            = "enabled"
		interface_type         = "access"
		cipher_suite           = "256GcmAes"
		window_size            = 128
		security_policy        = "shouldSecure"
		sak_expire_time        = 60
		confidentiality_offset = "offset30"
		key_server_priority    = 8
		macsec_keys {
			key_name           = "abc123"
			psk                = "AA111111111111111111111111111111111111111111111111111111111111aa"
			start_time         = "2027-09-23 00:00:00"
			end_time           = "2030-09-23 00:00:00"
		}
		macsec_keys {
			key_name           = "def456"
			psk                = "AA11111111111111111111111111111111111111111111111111111111111aaa"
			start_time         = "2029-12-11 11:12:13"
			end_time           = "2030-12-11 11:12:13"
		}
	}`, testAccMSOTemplateResourceFabricPolicyConfig())
}

func testAccMSOMacsecPolicyConfigUpdateRemovingExtraMacSecKey() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_macsec_policy" "macsec_policy" {
		template_id            = mso_template.template_fabric_policy.id
		name                   = "tf_test_macsec_policy"
		description            = "Terraform test MACsec Policy removing extra MACsec Key"
		admin_state            = "enabled"
		interface_type         = "access"
		cipher_suite           = "256GcmAes"
		window_size            = 128
		security_policy        = "shouldSecure"
		sak_expire_time        = 60
		confidentiality_offset = "offset30"
		key_server_priority    = 8
		macsec_keys {
			key_name           = "abc123"
			psk                = "AA111111111111111111111111111111111111111111111111111111111111aa"
			start_time         = "2027-09-23 00:00:00"
			end_time           = "2030-09-23 00:00:00"
		}
	}`, testAccMSOTemplateResourceFabricPolicyConfig())
}

func testAccMSOMacsecPolicyConfigUpdateChangingType() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_macsec_policy" "macsec_policy" {
		template_id            = mso_template.template_fabric_policy.id
		name                   = "tf_test_macsec_policy"
		description            = "Terraform test MACsec Policy changing interface type"
		admin_state            = "disabled"
		interface_type         = "fabric"
		cipher_suite           = "256GcmAesXpn"
		window_size            = 256
		security_policy        = "mustSecure"
		sak_expire_time        = 120
		macsec_keys {
			key_name           = "abc123"
			psk                = "AA111111111111111111111111111111111111111111111111111111111111ab"
			start_time         = "2028-12-12 12:12:12"
			end_time           = "infinite"
		}
	}`, testAccMSOTemplateResourceFabricPolicyConfig())
}
