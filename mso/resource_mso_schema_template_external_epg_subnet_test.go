package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// msoSchemaTemplateExtEpgSubnetSchemaId is set during the first test step's Check to capture the dynamic schema ID for use in the manual deletion PreConfig step.
var msoSchemaTemplateExtEpgSubnetSchemaId string

func TestAccMSOSchemaTemplateExternalEpgSubnetResource(t *testing.T) {
	resourceName := "mso_schema_template_external_epg_subnet." + msoSchemaTemplateExtEpgName + "_subnet"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateExtEpgSubnetDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create External EPG Subnet with required ip only") },
				Config:    testAccMSOSchemaTemplateExtEpgSubnetConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttr(resourceName, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr(resourceName, "external_epg_name", msoSchemaTemplateExtEpgName),
					resource.TestCheckResourceAttr(resourceName, "ip", msoSchemaTemplateExtEpgSubnetIp),
					resource.TestCheckResourceAttr(resourceName, "name", ""),
					// Verify defaults when scope and aggregate are not set in config:
					// scope is Computed so it reflects the server-side default (empty list);
					// aggregate is not Computed so it defaults to an empty list.
					resource.TestCheckResourceAttr(resourceName, "scope.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "aggregate.#", "0"),
					// Capture the dynamic schema ID from state for use in the manual deletion PreConfig step
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok {
							return fmt.Errorf("External EPG Subnet resource not found in state")
						}
						msoSchemaTemplateExtEpgSubnetSchemaId = rs.Primary.Attributes["schema_id"]
						return nil
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add External EPG Subnet name") },
				Config:    testAccMSOSchemaTemplateExtEpgSubnetConfigWithName(msoSchemaTemplateExtEpgSubnetName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ip", msoSchemaTemplateExtEpgSubnetIp),
					resource.TestCheckResourceAttr(resourceName, "name", msoSchemaTemplateExtEpgSubnetName),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update External EPG Subnet name") },
				Config:    testAccMSOSchemaTemplateExtEpgSubnetConfigWithName(msoSchemaTemplateExtEpgSubnetName2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ip", msoSchemaTemplateExtEpgSubnetIp),
					resource.TestCheckResourceAttr(resourceName, "name", msoSchemaTemplateExtEpgSubnetName2),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Set External EPG Subnet scope") },
				Config:    testAccMSOSchemaTemplateExtEpgSubnetConfigWithScope(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ip", msoSchemaTemplateExtEpgSubnetIp),
					resource.TestCheckResourceAttr(resourceName, "name", msoSchemaTemplateExtEpgSubnetName2),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "import-rtctrl"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "export-rtctrl"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Set External EPG Subnet aggregate (with required shared-rtctrl scope)") },
				Config:    testAccMSOSchemaTemplateExtEpgSubnetConfigWithAggregate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ip", msoSchemaTemplateExtEpgSubnetIp),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "shared-rtctrl"),
					resource.TestCheckResourceAttr(resourceName, "aggregate.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "aggregate.0", "shared-rtctrl"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update External EPG Subnet scope and aggregate") },
				Config:    testAccMSOSchemaTemplateExtEpgSubnetConfigWithScopeAndAggregateUpdated(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ip", msoSchemaTemplateExtEpgSubnetIp),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "shared-rtctrl"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "export-rtctrl"),
					resource.TestCheckResourceAttr(resourceName, "aggregate.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "aggregate.0", "shared-rtctrl"),
					resource.TestCheckResourceAttr(resourceName, "aggregate.1", "export-rtctrl"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Clear External EPG Subnet aggregate with empty list") },
				Config:    testAccMSOSchemaTemplateExtEpgSubnetConfigWithAggregateEmpty(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ip", msoSchemaTemplateExtEpgSubnetIp),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "shared-rtctrl"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "export-rtctrl"),
					resource.TestCheckResourceAttr(resourceName, "aggregate.#", "0"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Recreate External EPG Subnet on ip change (ForceNew)") },
				Config:    testAccMSOSchemaTemplateExtEpgSubnetConfigUpdateIp(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ip", msoSchemaTemplateExtEpgSubnetIp2),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import External EPG Subnet") },
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("External EPG Subnet resource not found in state")
					}
					return fmt.Sprintf("%s/templates/%s/externalEpgs/%s/ip/%s",
						rs.Primary.Attributes["schema_id"],
						rs.Primary.Attributes["template_name"],
						rs.Primary.Attributes["external_epg_name"],
						rs.Primary.Attributes["ip"]), nil
				},
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate External EPG Subnet after manual deletion from NDO")
					msoClient := testAccProvider.Meta().(*client.Client)
					subnetRemovePatchPayload := models.GetRemovePatchPayload(fmt.Sprintf("/templates/%s/externalEpgs/%s/subnets/0", msoSchemaTemplateName, msoSchemaTemplateExtEpgName))
					_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", msoSchemaTemplateExtEpgSubnetSchemaId), subnetRemovePatchPayload)
					if err != nil {
						t.Fatalf("Failed to manually delete External EPG Subnet: %v", err)
					}
				},
				Config: testAccMSOSchemaTemplateExtEpgSubnetConfigUpdateIp(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ip", msoSchemaTemplateExtEpgSubnetIp2),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateExtEpgSubnetPrerequisiteConfig() string {
	return fmt.Sprintf(`%s
	resource "mso_schema_template_external_epg" "%[2]s" {
		schema_id         = mso_schema.%[3]s.id
		template_name     = "%[4]s"
		external_epg_name = "%[2]s"
		display_name      = "%[2]s"
		vrf_name          = mso_schema_template_vrf.%[5]s.name
	}`, testAccMSOSchemaTemplateExtEpgPrerequisiteConfig(), msoSchemaTemplateExtEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfName)
}

func testAccMSOSchemaTemplateExtEpgSubnetConfigCreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_external_epg_subnet" "%[2]s_subnet" {
		schema_id         = mso_schema.%[3]s.id
		template_name     = "%[4]s"
		external_epg_name = mso_schema_template_external_epg.%[2]s.external_epg_name
		ip                = "%[5]s"
	}`, testAccMSOSchemaTemplateExtEpgSubnetPrerequisiteConfig(), msoSchemaTemplateExtEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateExtEpgSubnetIp)
}

func testAccMSOSchemaTemplateExtEpgSubnetConfigWithName(name string) string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_external_epg_subnet" "%[2]s_subnet" {
		schema_id         = mso_schema.%[3]s.id
		template_name     = "%[4]s"
		external_epg_name = mso_schema_template_external_epg.%[2]s.external_epg_name
		ip                = "%[5]s"
		name              = "%[6]s"
	}`, testAccMSOSchemaTemplateExtEpgSubnetPrerequisiteConfig(), msoSchemaTemplateExtEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateExtEpgSubnetIp, name)
}

func testAccMSOSchemaTemplateExtEpgSubnetConfigWithScope() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_external_epg_subnet" "%[2]s_subnet" {
		schema_id         = mso_schema.%[3]s.id
		template_name     = "%[4]s"
		external_epg_name = mso_schema_template_external_epg.%[2]s.external_epg_name
		ip                = "%[5]s"
		name              = "%[6]s"
		scope             = ["import-rtctrl", "export-rtctrl"]
	}`, testAccMSOSchemaTemplateExtEpgSubnetPrerequisiteConfig(), msoSchemaTemplateExtEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateExtEpgSubnetIp, msoSchemaTemplateExtEpgSubnetName2)
}

func testAccMSOSchemaTemplateExtEpgSubnetConfigWithAggregate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_external_epg_subnet" "%[2]s_subnet" {
		schema_id         = mso_schema.%[3]s.id
		template_name     = "%[4]s"
		external_epg_name = mso_schema_template_external_epg.%[2]s.external_epg_name
		ip                = "%[5]s"
		scope             = ["shared-rtctrl"]
		aggregate         = ["shared-rtctrl"]
	}`, testAccMSOSchemaTemplateExtEpgSubnetPrerequisiteConfig(), msoSchemaTemplateExtEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateExtEpgSubnetIp)
}

func testAccMSOSchemaTemplateExtEpgSubnetConfigWithScopeAndAggregateUpdated() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_external_epg_subnet" "%[2]s_subnet" {
		schema_id         = mso_schema.%[3]s.id
		template_name     = "%[4]s"
		external_epg_name = mso_schema_template_external_epg.%[2]s.external_epg_name
		ip                = "%[5]s"
		scope             = ["shared-rtctrl", "export-rtctrl"]
		aggregate         = ["shared-rtctrl", "export-rtctrl"]
	}`, testAccMSOSchemaTemplateExtEpgSubnetPrerequisiteConfig(), msoSchemaTemplateExtEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateExtEpgSubnetIp)
}

func testAccMSOSchemaTemplateExtEpgSubnetConfigWithAggregateEmpty() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_external_epg_subnet" "%[2]s_subnet" {
		schema_id         = mso_schema.%[3]s.id
		template_name     = "%[4]s"
		external_epg_name = mso_schema_template_external_epg.%[2]s.external_epg_name
		ip                = "%[5]s"
		scope             = ["shared-rtctrl", "export-rtctrl"]
		aggregate         = []
	}`, testAccMSOSchemaTemplateExtEpgSubnetPrerequisiteConfig(), msoSchemaTemplateExtEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateExtEpgSubnetIp)
}

func testAccMSOSchemaTemplateExtEpgSubnetConfigUpdateIp() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_external_epg_subnet" "%[2]s_subnet" {
		schema_id         = mso_schema.%[3]s.id
		template_name     = "%[4]s"
		external_epg_name = mso_schema_template_external_epg.%[2]s.external_epg_name
		ip                = "%[5]s"
	}`, testAccMSOSchemaTemplateExtEpgSubnetPrerequisiteConfig(), msoSchemaTemplateExtEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateExtEpgSubnetIp2)
}

func testAccCheckMSOSchemaTemplateExtEpgSubnetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "mso_schema_template_external_epg_subnet" {
			schemaID := rs.Primary.Attributes["schema_id"]
			con, err := client.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
			if err != nil {
				return nil
			}
			count, err := con.ArrayCount("templates")
			if err != nil {
				return fmt.Errorf("No Template found")
			}
			for i := 0; i < count; i++ {
				tempCont, err := con.ArrayElement(i, "templates")
				if err != nil {
					return fmt.Errorf("No template exists")
				}
				if models.StripQuotes(tempCont.S("name").String()) != rs.Primary.Attributes["template_name"] {
					continue
				}
				externalEpgCount, err := tempCont.ArrayCount("externalEpgs")
				if err != nil {
					return fmt.Errorf("Unable to get External EPG list")
				}
				for j := 0; j < externalEpgCount; j++ {
					externalEpgCont, err := tempCont.ArrayElement(j, "externalEpgs")
					if err != nil {
						return err
					}
					if models.StripQuotes(externalEpgCont.S("name").String()) != rs.Primary.Attributes["external_epg_name"] {
						continue
					}
					subnetCount, err := externalEpgCont.ArrayCount("subnets")
					if err != nil {
						return nil
					}
					for k := 0; k < subnetCount; k++ {
						subnetCont, err := externalEpgCont.ArrayElement(k, "subnets")
						if err != nil {
							return err
						}
						ip := models.StripQuotes(subnetCont.S("ip").String())
						if rs.Primary.Attributes["ip"] == ip {
							return fmt.Errorf("Schema Template External EPG Subnet record still exists")
						}
					}
				}
			}
		}
	}
	return nil
}
