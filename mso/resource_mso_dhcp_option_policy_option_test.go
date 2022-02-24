package mso

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccMSODHCPOptionPolicyOption_Basic(t *testing.T) {
	var opt1 models.DHCPOptionPolicyOption
	var opt2 models.DHCPOptionPolicyOption
	resourceName := "mso_dhcp_option_policy_option.test"
	tenant := tenantNames[0]
	optionPolicyName := makeTestVariable(acctest.RandString(5))
	name := acctest.RandString(5)
	nameOther := acctest.RandString(5)
	id := strconv.Itoa(acctest.RandIntRange(1, 1000))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSODHCPOptionPolicyOptionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      MSODHCPOptionPolicyOptionWithoutRequiredAttr(tenant, optionPolicyName, name, "option_policy_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSODHCPOptionPolicyOptionWithoutRequiredAttr(tenant, optionPolicyName, name, "option_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: MSODHCPOptionPolicyOptionWithRequired(tenant, optionPolicyName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPOptionPolicyOptionExists(resourceName, &opt1),
					resource.TestCheckResourceAttr(resourceName, "option_data", ""),
					resource.TestCheckResourceAttr(resourceName, "option_policy_name", optionPolicyName),
					resource.TestCheckResourceAttr(resourceName, "option_name", name),
					resource.TestCheckResourceAttr(resourceName, "option_id", ""),
				),
			},
			{
				Config: MSODHCPOptionPolicyOptionWithOptional(tenant, optionPolicyName, name, id),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPOptionPolicyOptionExists(resourceName, &opt2),
					resource.TestCheckResourceAttr(resourceName, "option_data", "test_data"),
					resource.TestCheckResourceAttr(resourceName, "option_policy_name", optionPolicyName),
					resource.TestCheckResourceAttr(resourceName, "option_name", name),
					resource.TestCheckResourceAttr(resourceName, "option_id", id),
					testAccCheckMSODHCPOptionPolicyOptionIdEqual(&opt1, &opt2),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      MSODHCPOptionPolicyOptionWithOutRequiredParameters(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: MSODHCPOptionPolicyOptionWithRequired(tenant, optionPolicyName, nameOther),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPOptionPolicyOptionExists(resourceName, &opt2),
					resource.TestCheckResourceAttr(resourceName, "option_policy_name", optionPolicyName),
					resource.TestCheckResourceAttr(resourceName, "option_name", nameOther),
					testAccCheckMSODHCPOptionPolicyOptionIdNotEqual(&opt1, &opt2),
				),
			},
			{
				Config: MSODHCPOptionPolicyOptionWithRequired(tenant, nameOther, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPOptionPolicyOptionExists(resourceName, &opt2),
					resource.TestCheckResourceAttr(resourceName, "option_policy_name", nameOther),
					resource.TestCheckResourceAttr(resourceName, "option_name", name),
					testAccCheckMSODHCPOptionPolicyOptionIdNotEqual(&opt1, &opt2),
				),
			},
		},
	})
}

func TestAccMSODHCPOptionPolicyOption_Negative(t *testing.T) {
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	tenant := tenantNames[0]
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSODHCPOptionPolicyOptionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      MSODHCPOptionPolicyOptionWithRequired(tenant, randomValue, makeTestVariable(acctest.RandString(5))),
				ExpectError: regexp.MustCompile(`value should be alphanumeric`),
			},
			{
				Config:      MSODHCPOptionPolicyOptionWithUpdatedAttr(tenant, randomValue, randomValue, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)+is not expected here.`),
			},
			{
				Config: MSODHCPOptionPolicyOptionWithRequired(tenant, randomValue, randomValue),
			},
		},
	})
}

func TestAccMSODHCPOptionPolicyOption_MultipleCreateDelete(t *testing.T) {
	tenant := tenantNames[0]
	name := acctest.RandString(5)
	optionPolicyName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSODHCPOptionPolicyOptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: MSODHCPOptionPolicyOptionMultiple(tenant, optionPolicyName, name),
			},
		},
	})
}

func testAccCheckMSODHCPOptionPolicyOptionExists(name string, m *models.DHCPOptionPolicyOption) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs1, err1 := s.RootModule().Resources[name]

		if !err1 {
			return fmt.Errorf("DHCP Option Policy Option %s not found", name)
		}

		if rs1.Primary.ID == "" {
			return fmt.Errorf("DHCP Option Policy Option was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		sts, _ := client.ReadDHCPOptionPolicyOption(rs1.Primary.ID)

		*m = *sts
		return nil
	}
}

func testAccCheckMSODHCPOptionPolicyOptionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_dhcp_option_policy_option" {
			_, err := client.ReadDHCPOptionPolicyOption(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Option still exists")
			}
		} else {
			continue
		}

	}
	return nil
}

func testAccCheckMSODHCPOptionPolicyOptionIdEqual(m1, m2 *models.DHCPOptionPolicyOption) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.Name != m2.Name && m1.PolicyName != m2.PolicyName {
			return fmt.Errorf("DHCP Option Policy Options are not equal")
		}
		return nil
	}
}

func testAccCheckMSODHCPOptionPolicyOptionIdNotEqual(m1, m2 *models.DHCPOptionPolicyOption) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.Name == m2.Name && m1.PolicyName == m2.PolicyName {
			return fmt.Errorf("DHCP Option Policy Options are equal")
		}
		return nil
	}
}

func MSODHCPOptionPolicyOptionWithUpdatedAttr(tenant, optionName, name, key, val string) string {
	resource := fmt.Sprintf(`
	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_option_policy" "test" {
		tenant_id = data.mso_tenant.test.id
		name = "%s"		
	}

	resource "mso_dhcp_option_policy_option" "test"{
		option_policy_name = mso_dhcp_option_policy.test.name
		option_name = "%s"
		%s = "%s"
	}
	`, tenant, tenant, optionName, name, key, val)
	return resource
}

func MSODHCPOptionPolicyOptionMultiple(tenant, policyName, name string) string {
	resource := fmt.Sprintf(`
	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_option_policy" "test" {
		tenant_id = data.mso_tenant.test.id
		name = "%s"
	}

	resource "mso_dhcp_option_policy_option" "test"{
		option_policy_name = mso_dhcp_option_policy.test.name
		option_name = "%s${count.index}"
		count = 5		
	}
	`, tenant, tenant, policyName, name)
	return resource
}

func MSODHCPOptionPolicyOptionWithRequired(tenant, policyName, name string) string {
	resource := fmt.Sprintf(`
	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_option_policy" "test" {
		tenant_id = data.mso_tenant.test.id
		name = "%s"		
	}
	resource "mso_dhcp_option_policy_option" "test"{
		option_policy_name = mso_dhcp_option_policy.test.name
		option_name = "%s"
	}
	`, tenant, tenant, policyName, name)
	return resource
}

func MSODHCPOptionPolicyOptionWithOutRequiredParameters() string {
	resource := fmt.Sprintln(`
	resource "mso_dhcp_option_policy_option" "test" {
		option_data = "test_data"
	}
	`)
	return resource
}

func MSODHCPOptionPolicyOptionWithOptional(tenant, name, optionname, optionid string) string {
	resource := fmt.Sprintf(`
	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_option_policy" "test" {
		tenant_id = data.mso_tenant.test.id
		name = "%s"
	}
	resource "mso_dhcp_option_policy_option" "test" {
		option_data = "test_data"
		option_policy_name = mso_dhcp_option_policy.test.name
		option_name = "%s"
		option_id = "%s"
	}
	`, tenant, tenant, name, optionname, optionid)
	return resource
}

func MSODHCPOptionPolicyOptionWithoutRequiredAttr(tenant, policyName, name, attr string) string {
	rBlock := `
	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_option_policy" "test" {
		tenant_id = data.mso_tenant.test.id
		name = "%s"
	}	
	`
	switch attr {
	case "option_policy_name":
		rBlock += `
		resource "mso_dhcp_option_policy_option" "test"{
		#	option_policy_name = mso_dhcp_option_policy.test.name
			option_name = "%s"
		}
		`
	case "option_name":
		rBlock += `
		resource "mso_dhcp_option_policy_option" "test"{
			option_policy_name = mso_dhcp_option_policy.test.name
		#	option_name = "%s"
		}
		`
	}
	return fmt.Sprintf(rBlock, tenant, tenant, policyName, name)
}
