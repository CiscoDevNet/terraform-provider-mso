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

func TestAccMSODHCPOptionPolicy_Basic(t *testing.T) {
	var optPolicy1 models.DHCPOptionPolicy
	var optPolicy2 models.DHCPOptionPolicy
	resourceName := "mso_dhcp_option_policy.test"
	tenant := makeTestVariable(acctest.RandString(5))
	tenantOther := makeTestVariable(acctest.RandString(5))
	name := makeTestVariable(acctest.RandString(5))
	nameOther := makeTestVariable(acctest.RandString(5))
	optionName := "acctest" + acctest.RandString(5)
	optionId := strconv.Itoa(acctest.RandIntRange(1, 1000))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSODHCPOptionPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      MSODHCPOptionPolicyWithoutRequired(tenant, name, "tenant_id"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSODHCPOptionPolicyWithoutRequired(tenant, name, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: MSODHCPOptionPolicyWithRequired(tenant, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPOptionPolicyExists(resourceName, &optPolicy1),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrSet(resourceName, "tenant_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "option.#", "0"),
				),
			},
			{
				Config:      MSODHCPOptionPolicyOptionWithoutRequired(tenant, name, optionName, optionId, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSODHCPOptionPolicyOptionWithoutRequired(tenant, name, optionName, optionId, "id"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: MSODHCPOptionPolicyWithOptioRequired(tenant, name, optionName, optionId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPOptionPolicyExists(resourceName, &optPolicy2),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrSet(resourceName, "tenant_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "option.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "option.0.name", optionName),
					resource.TestCheckResourceAttr(resourceName, "option.0.id", optionId),
					resource.TestCheckResourceAttr(resourceName, "option.0.data", ""),
					testAccCheckMSODHCPOptionPolicyIdEqual(&optPolicy1, &optPolicy2),
				),
			},
			{
				Config: MSODHCPOptionPolicyWithOptional(tenant, name, optionName, optionId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPOptionPolicyExists(resourceName, &optPolicy2),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttrSet(resourceName, "tenant_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "option.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "option.0.name", optionName),
					resource.TestCheckResourceAttr(resourceName, "option.0.id", optionId),
					resource.TestCheckResourceAttr(resourceName, "option.0.data", "test data"),
					testAccCheckMSODHCPOptionPolicyIdEqual(&optPolicy1, &optPolicy2),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      MSODHCPOptionPolicyWithOutRequiredParameters(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: MSODHCPOptionPolicyWithRequired(tenant, nameOther),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPOptionPolicyExists(resourceName, &optPolicy2),
					resource.TestCheckResourceAttrSet(resourceName, "tenant_id"),
					resource.TestCheckResourceAttr(resourceName, "name", nameOther),
					testAccCheckMSODHCPOptionPolicyIdNotEqual(&optPolicy1, &optPolicy2),
				),
			},
			{
				Config: MSODHCPOptionPolicyWithRequired(tenant, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPOptionPolicyExists(resourceName, &optPolicy1),
				),
			},
			{
				Config: MSODHCPOptionPolicyWithRequiredParam(tenant, tenantOther, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPOptionPolicyExists(resourceName, &optPolicy2),
					resource.TestCheckResourceAttrSet(resourceName, "tenant_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					testAccCheckMSODHCPOptionPolicyIdNotEqual(&optPolicy1, &optPolicy2),
				),
			},
		},
	})
}

func TestAccMSODHCPOptionPolicy_Update(t *testing.T) {
	var optPolicy1 models.DHCPOptionPolicy
	var optPolicy2 models.DHCPOptionPolicy
	resourceName := "mso_dhcp_option_policy.test"
	tenant := makeTestVariable(acctest.RandString(5))
	name := makeTestVariable(acctest.RandString(5))
	optionName := "acctest" + acctest.RandString(5)
	optionId := strconv.Itoa(acctest.RandIntRange(1, 1000))
	optionNameOther := "acctest" + acctest.RandString(5)
	optionIdOther := strconv.Itoa(acctest.RandIntRange(1, 1000))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSODHCPOptionPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: MSODHCPOptionPolicyWithMultipleOption(tenant, name, optionName, optionId, optionNameOther, optionIdOther),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPOptionPolicyExists(resourceName, &optPolicy1),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrSet(resourceName, "tenant_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "option.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "option.0.name", optionName),
					resource.TestCheckResourceAttr(resourceName, "option.0.id", optionId),
					resource.TestCheckResourceAttr(resourceName, "option.1.name", optionNameOther),
					resource.TestCheckResourceAttr(resourceName, "option.1.id", optionIdOther),
				),
			},
			{
				Config: MSODHCPOptionPolicyWithOptioRequired(tenant, name, optionName, optionId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPOptionPolicyExists(resourceName, &optPolicy2),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrSet(resourceName, "tenant_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "option.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "option.0.name", optionName),
					resource.TestCheckResourceAttr(resourceName, "option.0.id", optionId),
					testAccCheckMSODHCPOptionPolicyIdEqual(&optPolicy1, &optPolicy2),
				),
			},
		},
	})
}

func TestAccMSODHCPOptionPolicy_Negative(t *testing.T) {
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSODHCPOptionPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      MSODHCPOptionPolicyWithInvalidTenantId(randomValue, randomValue),
				ExpectError: regexp.MustCompile(`tenant with id (.)+ is not found`),
			},
			{
				Config:      MSODHCPOptionPolicyWithRequired(randomValue, acctest.RandString(1001)),
				ExpectError: regexp.MustCompile(`1 - 1000`),
			},
			{
				Config:      MSODHCPOptionPolicyWithRequiredUpdatedAttr(randomValue, randomValue, "description", acctest.RandString(1001)),
				ExpectError: regexp.MustCompile(`1 - 1000`),
			},
			{
				Config:      MSODHCPOptionPolicyWithRequiredUpdatedAttr(randomValue, randomValue, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)+is not expected here.`),
			},
			{
				Config:      MSODHCPOptionPolicyWithOptioRequired(randomValue, randomValue, makeTestVariable(acctest.RandString(5)), strconv.Itoa(acctest.RandIntRange(1, 1000))),
				ExpectError: regexp.MustCompile(`value should be alphanumeric`),
			},
			{
				Config:      MSODHCPOptionPolicyWithOptioRequired(randomValue, randomValue, randomValue, randomValue),
				ExpectError: regexp.MustCompile(`value should be numeric`),
			},
			{
				Config:      MSODHCPOptionPolicyWithOptioAttr(randomValue, randomValue, randomValue, strconv.Itoa(acctest.RandIntRange(1, 1000)), randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`Unsupported argument`),
			},
			{
				Config: MSODHCPOptionPolicyWithRequired(randomValue, randomValue),
			},
		},
	})
}

func TestAccMSODHCPOptionPolicy_MultipleCreateDelete(t *testing.T) {
	tenant := makeTestVariable(acctest.RandString(5))
	name := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSODHCPOptionPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: MSODHCPOptionPolicyMultiple(tenant, name),
			},
		},
	})
}

func testAccCheckMSODHCPOptionPolicyExists(name string, m *models.DHCPOptionPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs1, err1 := s.RootModule().Resources[name]

		if !err1 {
			return fmt.Errorf("DHCP Option Policy %s not found", name)
		}

		if rs1.Primary.ID == "" {
			return fmt.Errorf("DHCP Option Policy was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.ReadDHCPOptionPolicy(rs1.Primary.ID)

		if err != nil {
			return err
		}

		sts, _ := models.DHCPOptionPolicyFromContainer(cont)

		*m = *sts
		return nil
	}
}

func testAccCheckMSODHCPOptionPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_dhcp_option_policy" {
			_, err := client.ReadDHCPOptionPolicy(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("DHCP Option Policy still exists")
			}
		} else {
			continue
		}

	}
	return nil
}

func testAccCheckMSODHCPOptionPolicyIdEqual(m1, m2 *models.DHCPOptionPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.ID != m2.ID {
			return fmt.Errorf("DHCP Option Policies are not equal")
		}
		return nil
	}
}

func testAccCheckMSODHCPOptionPolicyIdNotEqual(m1, m2 *models.DHCPOptionPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.ID == m2.ID {
			return fmt.Errorf("DHCP Option Policies are equal")
		}
		return nil
	}
}

func MSODHCPOptionPolicyWithRequiredUpdatedAttr(tenant, name, key, val string) string {
	resource := fmt.Sprintf(`
	resource "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_option_policy" "test" {
		tenant_id = mso_tenant.test.id
		name = "%s"		
		%s = "%s"
	}
	`, tenant, tenant, name, key, val)
	return resource
}

func MSODHCPOptionPolicyWithInvalidTenantId(tenant, name string) string {
	resource := fmt.Sprintf(`
	resource "mso_dhcp_option_policy" "test" {
		tenant_id = "%s"
		name = "%s"		
	}
	`, tenant, name)
	return resource
}

func MSODHCPOptionPolicyMultiple(tenant, name string) string {
	resource := fmt.Sprintf(`
	resource "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_option_policy" "test" {
		tenant_id = mso_tenant.test.id
		name = "%s${count.index}"
		count = 5		
	}
	`, tenant, tenant, name)
	return resource
}

func MSODHCPOptionPolicyWithRequired(tenant, name string) string {
	resource := fmt.Sprintf(`
	resource "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_option_policy" "test" {
		tenant_id = mso_tenant.test.id
		name = "%s"		
	}
	`, tenant, tenant, name)
	return resource
}

func MSODHCPOptionPolicyWithRequiredParam(tenant, tenantother, name string) string {
	resource := fmt.Sprintf(`
	resource "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_tenant" "test1" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_option_policy" "test" {
		tenant_id = mso_tenant.test1.id
		name = "%s"		
	}
	`, tenant, tenant, tenantother, tenantother, name)
	return resource
}

func MSODHCPOptionPolicyWithMultipleOption(tenant, name, opname, opid, opname1, opid1 string) string {
	resource := fmt.Sprintf(`
	resource "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_option_policy" "test" {
		tenant_id = mso_tenant.test.id
		name = "%s"
		option {
			name = "%s"
			id = "%s"
		}
		option {
			name = "%s"
			id = "%s"
		}
	}
	`, tenant, tenant, name, opname, opid, opname1, opid1)
	return resource
}

func MSODHCPOptionPolicyWithOptioAttr(tenant, name, optionname, optionid, key, val string) string {
	resource := fmt.Sprintf(`
	resource "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_option_policy" "test" {
		tenant_id = mso_tenant.test.id
		name = "%s"
		option {
			name = "%s"
			id = "%s"
			%s = "%s"
		}
	}
	`, tenant, tenant, name, optionname, optionid, key, val)
	return resource
}

func MSODHCPOptionPolicyWithOptioRequired(tenant, name, optionname, optionid string) string {
	resource := fmt.Sprintf(`
	resource "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_option_policy" "test" {
		tenant_id = mso_tenant.test.id
		name = "%s"
		option {
			name = "%s"
			id = "%s"
		}
	}
	`, tenant, tenant, name, optionname, optionid)
	return resource
}

func MSODHCPOptionPolicyWithOutRequiredParameters() string {
	resource := fmt.Sprintln(`
	resource "mso_dhcp_option_policy" "test" {
		description = "test description updated"
	}
	`)
	return resource
}

func MSODHCPOptionPolicyWithOptional(tenant, name, optionname, optionid string) string {
	resource := fmt.Sprintf(`
	resource "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_option_policy" "test" {
		tenant_id = mso_tenant.test.id
		name = "%s"
		description = "test description"
		option {
			name = "%s"
			id = "%s"
			data = "test data"
		}
	}
	`, tenant, tenant, name, optionname, optionid)
	return resource
}

func MSODHCPOptionPolicyWithoutRequired(tenant, name, attr string) string {
	rBlock := `
	resource "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	`
	switch attr {
	case "tenant_id":
		rBlock += `
		resource "mso_dhcp_option_policy" "test" {
		#	tenant_id = mso_tenant.test.id
			name = "%s"
			
		}`
	case "name":
		rBlock += `
		resource "mso_dhcp_option_policy" "test" {
			tenant_id = mso_tenant.test.id
		#	name = "%s"	
		}`
	}
	return fmt.Sprintf(rBlock, tenant, tenant, name)
}

func MSODHCPOptionPolicyOptionWithoutRequired(tenant, name, optionname, optionid, attr string) string {
	rBlock := `
	resource "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	`
	switch attr {
	case "name":
		rBlock += `
		resource "mso_dhcp_option_policy" "test" {
			tenant_id = mso_tenant.test.id
			name = "%s"
			option {
		#		name = "%s"
				id = "%s"
			}
		}`
	case "id":
		rBlock += `
		resource "mso_dhcp_option_policy" "test" {
			tenant_id = mso_tenant.test.id
			name = "%s"
			option {
				name = "%s"
		#		id = "%s"
			}	
		}`
	}
	return fmt.Sprintf(rBlock, tenant, tenant, name, optionname, optionid)
}
