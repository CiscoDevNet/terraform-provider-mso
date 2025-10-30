package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOMacsecPolicyDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: MACsec Policy Data Source") },
				Config:    testAccMSOMacsecPolicyDataSource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_fabric_policies_macsec_policy.macsec_policy", "name", "tf_test_macsec_policy"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_macsec_policy.macsec_policy", "description", "Terraform test MACsec Policy"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_macsec_policy.macsec_policy", "admin_state", "enabled"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_macsec_policy.macsec_policy", "interface_type", "access"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_macsec_policy.macsec_policy", "cipher_suite", "256GcmAes"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_macsec_policy.macsec_policy", "window_size", "128"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_macsec_policy.macsec_policy", "security_policy", "shouldSecure"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_macsec_policy.macsec_policy", "sak_expire_time", "60"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_macsec_policy.macsec_policy", "confidentiality_offset", "offset30"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_macsec_policy.macsec_policy", "key_server_priority", "8"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_macsec_policy.macsec_policy", "macsec_key.#", "1"),
					customTestCheckResourceTypeSetAttr("data.mso_fabric_policies_macsec_policy.macsec_policy", "macsec_key",
						map[string]string{
							"key_name":   "abc123",
							"psk":        "AA111111111111111111111111111111111111111111111111111111111111aa",
							"start_time": "2027-09-23 00:00:00",
							"end_time":   "2030-09-23 00:00:00",
						},
					),
				),
			},
		},
	})
}

func testAccMSOMacsecPolicyDataSource() string {
	return fmt.Sprintf(`%s
	data "mso_fabric_policies_macsec_policy" "macsec_policy" {
	    template_id        = mso_fabric_policies_macsec_policy.macsec_policy.template_id
	    name               = "tf_test_macsec_policy"
    }`, testAccMSOMacsecPolicyConfigCreate())
}
