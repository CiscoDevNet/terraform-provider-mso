package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccMSODHCPRelayPolicy_Basic(t *testing.T) {
	var relayPolicy1 models.DHCPRelayPolicy
	var relayPolicy2 models.DHCPRelayPolicy
	resourceName := "mso_dhcp_relay_policy.test"
	tenant := tenantNames[0]
	epg := epg
	name := makeTestVariable(acctest.RandString(5))
	nameOther := makeTestVariable(acctest.RandString(5))
	displayName := makeTestVariable(acctest.RandString(5))
	schemaName := makeTestVariable(acctest.RandString(5))
	templateName := makeTestVariable(acctest.RandString(5))
	dhcpServerAddress, _ := acctest.RandIpAddress("1.2.0.0/16")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSODHCPRelayPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      MSODHCPRelayPolicyWithoutRequired(tenant, name, "tenant_id"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSODHCPRelayPolicyWithoutRequired(tenant, name, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: MSODHCPRelayPolicyWithRequired(tenant, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPRelayPolicyExists(resourceName, &relayPolicy1),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrSet(resourceName, "tenant_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "dhcp_relay_policy_provider.#", "0"),
				),
			},
			{
				Config:      MSODHCPRelayPolicyProviderWithoutEPG(tenant, name, dhcpServerAddress, epg, "epg"),
				ExpectError: regexp.MustCompile(`expected any one of the epg or external_epg.`),
			},
			{
				Config:      MSODHCPRelayPolicyProviderWithoutEPG(tenant, name, dhcpServerAddress, epg, "dhcp_server_address"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSODHCPRelayPolicyProviderWithoutExternalEPG(tenant, name, dhcpServerAddress, schemaName, templateName, displayName),
				ExpectError: regexp.MustCompile(`expected any one of the epg or external_epg.`),
			},
			{
				Config: MSODHCPRelayPolicyWithEPGRequired(tenant, name, dhcpServerAddress, epg),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPRelayPolicyExists(resourceName, &relayPolicy2),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrSet(resourceName, "tenant_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "dhcp_relay_policy_provider.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dhcp_relay_policy_provider.0.epg", epg),
					resource.TestCheckResourceAttr(resourceName, "dhcp_relay_policy_provider.0.dhcp_server_address", dhcpServerAddress),
					testAccCheckMSODHCPRelayPolicyIdEqual(&relayPolicy1, &relayPolicy2),
				),
			},
			{
				Config: MSODHCPRelayPolicyWithExternalEPGRequired(tenant, name, schemaName, templateName, dhcpServerAddress, displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPRelayPolicyExists(resourceName, &relayPolicy2),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrSet(resourceName, "tenant_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "dhcp_relay_policy_provider.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "dhcp_relay_policy_provider.0.external_epg"),
					resource.TestCheckResourceAttr(resourceName, "dhcp_relay_policy_provider.0.dhcp_server_address", dhcpServerAddress),
					testAccCheckMSODHCPRelayPolicyIdEqual(&relayPolicy1, &relayPolicy2),
				),
			},
			{
				Config: MSODHCPRelayPolicyWithOptional(tenant, name, epg, dhcpServerAddress),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPRelayPolicyExists(resourceName, &relayPolicy2),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttrSet(resourceName, "tenant_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "dhcp_relay_policy_provider.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "dhcp_relay_policy_provider.0.epg"),
					resource.TestCheckResourceAttr(resourceName, "dhcp_relay_policy_provider.0.dhcp_server_address", dhcpServerAddress),
					testAccCheckMSODHCPRelayPolicyIdEqual(&relayPolicy1, &relayPolicy2),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      MSODHCPRelayPolicyWithOutRequiredParameters(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: MSODHCPRelayPolicyWithRequired(tenant, nameOther),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPRelayPolicyExists(resourceName, &relayPolicy2),
					resource.TestCheckResourceAttrSet(resourceName, "tenant_id"),
					resource.TestCheckResourceAttr(resourceName, "name", nameOther),
					testAccCheckMSODHCPRelayPolicyIdNotEqual(&relayPolicy1, &relayPolicy2),
				),
			},
		},
	})
}

func TestAccMSODHCPRelayPolicy_Update(t *testing.T) {
	var relayPolicy1 models.DHCPRelayPolicy
	var relayPolicy2 models.DHCPRelayPolicy
	resourceName := "mso_dhcp_relay_policy.test"
	tenant := tenantNames[0]
	name := makeTestVariable(acctest.RandString(5))
	dhcpServerAddress, _ := acctest.RandIpAddress("1.2.0.0/16")
	otherDhcpServerAddress, _ := acctest.RandIpAddress("1.2.0.0/16")
	displayName := makeTestVariable(acctest.RandString(5))
	otherDisplayName := makeTestVariable(acctest.RandString(5))
	schemaName := makeTestVariable(acctest.RandString(5))
	templateName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSODHCPRelayPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: MSODHCPRelayPolicyWithMultipleProvider(tenant, name, dhcpServerAddress, otherDhcpServerAddress, schemaName, templateName, displayName, otherDisplayName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPRelayPolicyExists(resourceName, &relayPolicy1),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrSet(resourceName, "tenant_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "dhcp_relay_policy_provider.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "dhcp_relay_policy_provider.0.external_epg"),
					resource.TestCheckResourceAttrSet(resourceName, "dhcp_relay_policy_provider.0.dhcp_server_address"),
					resource.TestCheckResourceAttrSet(resourceName, "dhcp_relay_policy_provider.1.external_epg"),
					resource.TestCheckResourceAttrSet(resourceName, "dhcp_relay_policy_provider.1.dhcp_server_address"),
				),
			},
			{
				Config: MSODHCPRelayPolicyWithEPGRequired(tenant, name, dhcpServerAddress, epg),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSODHCPRelayPolicyExists(resourceName, &relayPolicy2),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrSet(resourceName, "tenant_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "dhcp_relay_policy_provider.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "dhcp_relay_policy_provider.0.epg"),
					resource.TestCheckResourceAttrSet(resourceName, "dhcp_relay_policy_provider.0.dhcp_server_address"),
					testAccCheckMSODHCPRelayPolicyIdEqual(&relayPolicy1, &relayPolicy2),
				),
			},
		},
	})
}

func TestAccMSODHCPRelayPolicy_Negative(t *testing.T) {
	tenant := tenantNames[0]
	name := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	dhcpServerAddress, _ := acctest.RandIpAddress("1.2.0.0/16")
	displayName := makeTestVariable(acctest.RandString(5))
	schemaName := makeTestVariable(acctest.RandString(5))
	templateName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSODHCPRelayPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      MSODHCPRelayPolicyWithInvalidTenantId(randomValue, name),
				ExpectError: regexp.MustCompile(`tenant with id (.)+ is not found`),
			},
			{
				Config:      MSODHCPRelayPolicyWithRequired(tenant, acctest.RandString(1001)),
				ExpectError: regexp.MustCompile(`1 - 1000`),
			},
			{
				Config:      MSODHCPRelayPolicyWithRequiredUpdatedAttr(tenant, name, "description", acctest.RandString(1001)),
				ExpectError: regexp.MustCompile(`1 - 1000`),
			},
			{
				Config:      MSODHCPRelayPolicyWithRequiredUpdatedAttr(tenant, name, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)+is not expected here.`),
			},
			{
				Config:      MSODHCPRelayPolicyWithMultipleProvider(tenant, name, dhcpServerAddress, dhcpServerAddress, schemaName, templateName, displayName, displayName),
				ExpectError: regexp.MustCompile(`duplicate`),
			},
			{
				Config:      MSODHCPRelayPolicyWithExternalEPGRequired(tenant, name, schemaName, templateName, randomValue, displayName),
				ExpectError: regexp.MustCompile(`to contain a valid IP`),
			},
			{
				Config:      MSODHCPRelayPolicyWithExtraAttr(tenant, name, dhcpServerAddress, schemaName, templateName, displayName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`Unsupported argument`),
			},
			{
				Config: MSODHCPRelayPolicyWithRequired(tenant, name),
			},
		},
	})
}

func TestAccMSODHCPRelayPolicy_MultipleCreateDelete(t *testing.T) {
	tenant := tenantNames[0]
	name := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSODHCPRelayPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: MSODHCPRelayPolicyMultiple(tenant, name),
			},
		},
	})
}

func testAccCheckMSODHCPRelayPolicyExists(name string, m *models.DHCPRelayPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs1, err1 := s.RootModule().Resources[name]

		if !err1 {
			return fmt.Errorf("DHCP Relay Policy %s not found", name)
		}

		if rs1.Primary.ID == "" {
			return fmt.Errorf("DHCP Relay Policy was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.ReadDHCPRelayPolicy(rs1.Primary.ID)

		if err != nil {
			return err
		}

		sts, _ := models.DHCPRelayPolicyFromContainer(cont)

		*m = *sts
		return nil
	}
}

func testAccCheckMSODHCPRelayPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_dhcp_relay_policy" {
			_, err := client.ReadDHCPRelayPolicy(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Label still exists")
			}
		} else {
			continue
		}

	}
	return nil
}

func testAccCheckMSODHCPRelayPolicyIdEqual(m1, m2 *models.DHCPRelayPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.ID != m2.ID {
			return fmt.Errorf("DHCP Relay Policies are not equal")
		}
		return nil
	}
}

func testAccCheckMSODHCPRelayPolicyIdNotEqual(m1, m2 *models.DHCPRelayPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.ID == m2.ID {
			return fmt.Errorf("DHCP relay Policies are equal")
		}
		return nil
	}
}

func MSODHCPRelayPolicyWithRequiredUpdatedAttr(tenant, name, key, val string) string {
	resource := fmt.Sprintf(`
	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_relay_policy" "test" {
		tenant_id = data.mso_tenant.test.id
		name = "%s"		
		%s = "%s"
	}
	`, tenant, tenant, name, key, val)
	return resource
}

func MSODHCPRelayPolicyWithInvalidTenantId(tenant, name string) string {
	resource := fmt.Sprintf(`
	resource "mso_dhcp_relay_policy" "test" {
		tenant_id = "%s"
		name = "%s"		
	}
	`, tenant, name)
	return resource
}

func MSODHCPRelayPolicyMultiple(tenant, name string) string {
	resource := fmt.Sprintf(`
	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_relay_policy" "test" {
		tenant_id = data.mso_tenant.test.id
		name = "%s${count.index}"
		count = 5		
	}
	`, tenant, tenant, name)
	return resource
}

func MSODHCPRelayPolicyWithRequired(tenant, name string) string {
	resource := fmt.Sprintf(`
	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_relay_policy" "test" {
		tenant_id = data.mso_tenant.test.id
		name = "%s"		
	}
	`, tenant, tenant, name)
	return resource
}

func MSODHCPRelayPolicyWithMultipleProvider(tenant, name, dhcpServerAddress, otherDhcpServerAddress, schemaName, templateName, displayName, otherDisplayName string) string {
	resource := fmt.Sprintf(`
	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource mso_schema "test"{
		name = "%s"
		template_name = "%s"
		tenant_id = data.mso_tenant.test.id
	}
	resource mso_schema_template_vrf "test" {
		schema_id = mso_schema.test.id
		template= mso_schema.test.template_name
		name= "%s"
		display_name= "%s"
	}
	resource "mso_schema_template_external_epg" "test" {
		schema_id = mso_schema.test.id
		template_name = mso_schema.test.template_name
		external_epg_name = "%s"
		display_name = "%s"
		vrf_name = mso_schema_template_vrf.test.name
	}
	resource "mso_schema_template_external_epg" "test1" {
		schema_id = mso_schema.test.id
		template_name = mso_schema.test.template_name
		external_epg_name = "%s"
		display_name = "%s"
		vrf_name = mso_schema_template_vrf.test.name
	}
	resource "mso_dhcp_relay_policy" "test" {
		tenant_id = data.mso_tenant.test.id
		name = "%s"
		dhcp_relay_policy_provider {
			dhcp_server_address = "%s"
			external_epg = mso_schema_template_external_epg.test.id
		}
		dhcp_relay_policy_provider {
			dhcp_server_address = "%s"
			external_epg = mso_schema_template_external_epg.test1.id
		}
	}
	`, tenant, tenant, schemaName, templateName, displayName, displayName, displayName, displayName, otherDisplayName, otherDisplayName, name, dhcpServerAddress, otherDhcpServerAddress)
	return resource
}

func MSODHCPRelayPolicyWithExtraAttr(tenant, name, dhcpServerAddress, schemaName, templateName, displayName, key, val string) string {
	resource := fmt.Sprintf(`
	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource mso_schema "test"{
		name = "%s"
		template_name = "%s"
		tenant_id = data.mso_tenant.test.id
	}
	resource mso_schema_template_vrf "test" {
		schema_id = mso_schema.test.id
		template= mso_schema.test.template_name
		name= "%s"
		display_name= "%s"
	}
	resource "mso_schema_template_external_epg" "test" {
		schema_id = mso_schema.test.id
		template_name = mso_schema.test.template_name
		external_epg_name = "%s"
		display_name = "%s"
		vrf_name = mso_schema_template_vrf.test.name
	}
	resource "mso_dhcp_relay_policy" "test" {
		tenant_id = data.mso_tenant.test.id
		name = "%s"
		dhcp_relay_policy_provider {
			dhcp_server_address = "%s"
			external_epg = mso_schema_template_external_epg.test.id
			%s = "%s"
		}
	}
	`, tenant, tenant, schemaName, templateName, displayName, displayName, displayName, displayName, name, dhcpServerAddress, key, val)
	return resource
}

func MSODHCPRelayPolicyWithEPGRequired(tenant, name, dhcpServerAddress, epg string) string {
	resource := fmt.Sprintf(`
	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_relay_policy" "test" {
		tenant_id = data.mso_tenant.test.id
		name = "%s"
		dhcp_relay_policy_provider {
			epg = "%s"
			dhcp_server_address = "%s"
		}
	}
	`, tenant, tenant, name, epg, dhcpServerAddress)
	return resource
}

func MSODHCPRelayPolicyWithExternalEPGRequired(tenant, name, schemaName, templateName, dhcpServerAddress, displayName string) string {
	resource := fmt.Sprintf(`
	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource mso_schema "test"{
		name = "%s"
		template_name = "%s"
		tenant_id = data.mso_tenant.test.id
	}
	resource mso_schema_template_vrf "test" {
		schema_id = mso_schema.test.id
		template= mso_schema.test.template_name
		name= "%s"
		display_name= "%s"
	}
	resource "mso_schema_template_external_epg" "test" {
		schema_id = mso_schema.test.id
		template_name = mso_schema.test.template_name
		external_epg_name = "%s"
		display_name = "%s"
		vrf_name = mso_schema_template_vrf.test.name
	}
	resource "mso_dhcp_relay_policy" "test" {
		tenant_id = data.mso_tenant.test.id
		name = "%s"
		dhcp_relay_policy_provider {
			external_epg = mso_schema_template_external_epg.test.id
			dhcp_server_address = "%s"
		}
	}
	`, tenant, tenant, schemaName, templateName, displayName, displayName, displayName, displayName, name, dhcpServerAddress)
	return resource
}

func MSODHCPRelayPolicyWithOutRequiredParameters() string {
	resource := fmt.Sprintln(`
	resource "mso_dhcp_relay_policy" "test" {
		description = "test description updated"
	}
	`)
	return resource
}

func MSODHCPRelayPolicyWithOptional(tenant, name, epg, dhcpServerAddress string) string {
	resource := fmt.Sprintf(`
	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource "mso_dhcp_relay_policy" "test" {
		tenant_id = data.mso_tenant.test.id
		name = "%s"
		description = "test description"
		dhcp_relay_policy_provider {
			epg = "%s"
	        dhcp_server_address = "%s"
		}
	}
	`, tenant, tenant, name, epg, dhcpServerAddress)
	return resource
}

func MSODHCPRelayPolicyWithoutRequired(tenant, name, attr string) string {
	rBlock := `
	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	`
	switch attr {
	case "tenant_id":
		rBlock += `
		resource "mso_dhcp_relay_policy" "test" {
		#	tenant_id = data.mso_tenant.test.id
			name = "%s"
		}`
	case "name":
		rBlock += `
		resource "mso_dhcp_relay_policy" "test" {
			tenant_id = data.mso_tenant.test.id
		#	name = "%s"	
		}`
	}
	return fmt.Sprintf(rBlock, tenant, tenant, name)
}

func MSODHCPRelayPolicyProviderWithoutEPG(tenant, name, dhcpServerAddress, epg, attr string) string {
	rBlock := `
	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	`
	switch attr {
	case "epg":
		rBlock += `
		resource "mso_dhcp_relay_policy" "test" {
			tenant_id = data.mso_tenant.test.id
			name = "%s"
			dhcp_relay_policy_provider {
		#		epg = "%s"
				dhcp_server_address = "%s"
			}
		}`
	case "dhcp_server_address":
		rBlock += `
		resource "mso_dhcp_relay_policy" "test" {
			tenant_id = data.mso_tenant.test.id
			name = "%s"
			dhcp_relay_policy_provider {
				epg = "%s"
		#       dhcp_server_address = "%s"
			}	
		}`
	}
	return fmt.Sprintf(rBlock, tenant, tenant, name, epg, dhcpServerAddress)
}

func MSODHCPRelayPolicyProviderWithoutExternalEPG(tenant, name, dhcpServerAddress, schemaName, templateName, displayName string) string {
	rBlock := `
	data "mso_tenant" "test" {
		name = "%s"
		display_name = "%s"
	}
	resource mso_schema "test"{
		name = "%s"
		template_name = "%s"
		tenant_id = data.mso_tenant.test.id
	}
	resource mso_schema_template_vrf "test" {
		schema_id = mso_schema.test.id
		template= mso_schema.test.template_name
		name= "%s"
		display_name= "%s"
	}
	resource "mso_schema_template_external_epg" "test" {
		schema_id = mso_schema.test.id
		template_name = mso_schema.test.template_name
		external_epg_name = "%s"
		display_name = "%s"
		vrf_name = mso_schema_template_vrf.test.name
	}
	resource "mso_dhcp_relay_policy" "test" {
		tenant_id = data.mso_tenant.test.id
		name = "%s"
		dhcp_relay_policy_provider {
			dhcp_server_address = "%s"
	#		external_epg = mso_schema_template_external_epg.test.id
		}
	}`
	return fmt.Sprintf(rBlock, tenant, tenant, schemaName, templateName, displayName, displayName, displayName, displayName, name, dhcpServerAddress)
}
