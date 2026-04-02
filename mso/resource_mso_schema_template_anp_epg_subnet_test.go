package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Note: The `querier` and `primary` attributes are not tested because they are BD (Bridge Domain) specific attributes.
// - querier: API returns "EPG: <name> in Schema: <schema>, Template: <template>, 'Querier' is only supported for Bridge Domain subnets"
// - primary: API returns "EPG: <name> in Schema: <schema>, Template: <template>, EPG Subnet <ip> cannot be marked as primary"
// See https://www.cisco.com/c/en/us/td/docs/dcn/ndo/4x/articles-441/nexus-dashboard-orchestrator-aci-schemas-and-application-templates-441.html#_configuring_bridge_domains for more details
// msoSchemaTemplateAnpEpgSubnetSchemaId is set during the first test step's Check to capture the dynamic schema ID for use in the manual deletion PreConfig step.
var msoSchemaTemplateAnpEpgSubnetSchemaId string

func TestAccMSOSchemaTemplateAnpEpgSubnetResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateAnpEpgSubnetDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create EPG Subnet with required attributes only") },
				Config:    testAccMSOSchemaTemplateAnpEpgSubnetConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "template", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "anp_name", msoSchemaTemplateAnpName),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "epg_name", msoSchemaTemplateAnpEpgName),
					// Capture the dynamic schema ID from state for use in the manual deletion PreConfig step
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet"]
						if !ok {
							return fmt.Errorf("EPG Subnet resource not found in state")
						}
						msoSchemaTemplateAnpEpgSubnetSchemaId = rs.Primary.Attributes["schema_id"]
						return nil
					},
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "ip", msoSchemaTemplateAnpEpgSubnetIp),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "scope", "private"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "shared", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "querier", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "no_default_gateway", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "primary", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "description", ""),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update EPG Subnet scope to public") },
				Config:    testAccMSOSchemaTemplateAnpEpgSubnetConfigUpdateScope(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "ip", msoSchemaTemplateAnpEpgSubnetIp),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "scope", "public"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "shared", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "querier", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "no_default_gateway", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "primary", "false"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update EPG Subnet shared and description") },
				Config:    testAccMSOSchemaTemplateAnpEpgSubnetConfigUpdateSharedAndDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "ip", msoSchemaTemplateAnpEpgSubnetIp),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "scope", "public"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "shared", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "description", "test subnet"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update EPG Subnet no_default_gateway") },
				Config:    testAccMSOSchemaTemplateAnpEpgSubnetConfigUpdateAllAttributes(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "ip", msoSchemaTemplateAnpEpgSubnetIp),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "scope", "public"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "shared", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "description", "test subnet"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "no_default_gateway", "true"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Reset EPG Subnet optional attributes") },
				Config:    testAccMSOSchemaTemplateAnpEpgSubnetConfigResetAttributes(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "ip", msoSchemaTemplateAnpEpgSubnetIp),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "scope", "private"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "shared", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "description", ""),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "querier", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "no_default_gateway", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "primary", "false"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update EPG Subnet IP") },
				Config:    testAccMSOSchemaTemplateAnpEpgSubnetConfigUpdateIp(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "ip", msoSchemaTemplateAnpEpgSubnetIp2),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "scope", "private"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "shared", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "querier", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "no_default_gateway", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "primary", "false"),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import EPG Subnet") },
				ResourceName: "mso_schema_template_anp_epg_subnet." + msoSchemaTemplateAnpEpgName + "_subnet",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet"]
					if !ok {
						return "", fmt.Errorf("EPG Subnet resource not found in state")
					}
					return fmt.Sprintf("%s/templates/%s/anps/%s/epgs/%s/ip/%s", rs.Primary.Attributes["schema_id"], rs.Primary.Attributes["template"], rs.Primary.Attributes["anp_name"], rs.Primary.Attributes["epg_name"], rs.Primary.Attributes["ip"]), nil
				},
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate EPG Subnet after manual deletion from NDO")
					msoClient := testAccProvider.Meta().(*client.Client)
					cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", msoSchemaTemplateAnpEpgSubnetSchemaId))
					if err != nil {
						t.Fatalf("Failed to get schema: %v", err)
					}
					index, err := fetchIndex(cont, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateAnpEpgSubnetIp2)
					if err != nil {
						t.Fatalf("Failed to fetch subnet index: %v", err)
					}
					if index == -1 {
						t.Fatalf("Subnet not found for manual deletion")
					}
					subnetRemovePatchPayload := models.GetRemovePatchPayload(fmt.Sprintf("/templates/%s/anps/%s/epgs/%s/subnets/%d", msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName, index))
					_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", msoSchemaTemplateAnpEpgSubnetSchemaId), subnetRemovePatchPayload)
					if err != nil {
						t.Fatalf("Failed to manually delete subnet: %v", err)
					}
				},
				Config: testAccMSOSchemaTemplateAnpEpgSubnetConfigUpdateIp(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "ip", msoSchemaTemplateAnpEpgSubnetIp2),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "scope", "private"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "shared", "false"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update parent EPG description with subnet present") },
				Config:    testAccMSOSchemaTemplateAnpEpgSubnetConfigUpdateParentEpg(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "ip", msoSchemaTemplateAnpEpgSubnetIp2),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "scope", "private"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg_subnet."+msoSchemaTemplateAnpEpgName+"_subnet", "shared", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_anp_epg."+msoSchemaTemplateAnpEpgName, "description", "Updated EPG description with subnet"),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateAnpEpgSubnetPrerequisiteConfig() string {
	return fmt.Sprintf(`%s%s%s%s%s%s`, testSiteConfigAnsibleTest(), testTenantConfig(), testSchemaConfig(), testSchemaTemplateVrfConfig(), testSchemaTemplateBdConfig(), testSchemaTemplateAnpConfig()) + fmt.Sprintf(`
resource "mso_schema_template_anp_epg" "%[1]s" {
	name          = "%[1]s"
	display_name  = "%[1]s"
	anp_name      = mso_schema_template_anp.%[2]s.name
	schema_id     = mso_schema.%[3]s.id
	template_name = "%[4]s"
	bd_name       = mso_schema_template_bd.%[5]s.name
}
`, msoSchemaTemplateAnpEpgName, msoSchemaTemplateAnpName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateBdName)
}

func testAccMSOSchemaTemplateAnpEpgSubnetConfigCreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp_epg_subnet" "%[2]s_subnet" {
		schema_id  = mso_schema_template_anp_epg.%[6]s.schema_id
		template   = "%[4]s"
		anp_name   = "%[5]s"
		epg_name   = mso_schema_template_anp_epg.%[6]s.name
		ip         = "%[7]s"
	}`, testAccMSOSchemaTemplateAnpEpgSubnetPrerequisiteConfig(), msoSchemaTemplateAnpEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateAnpEpgSubnetIp)
}

func testAccMSOSchemaTemplateAnpEpgSubnetConfigUpdateScope() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp_epg_subnet" "%[2]s_subnet" {
		schema_id  = mso_schema_template_anp_epg.%[6]s.schema_id
		template   = "%[4]s"
		anp_name   = "%[5]s"
		epg_name   = mso_schema_template_anp_epg.%[6]s.name
		ip         = "%[7]s"
		scope      = "public"
	}`, testAccMSOSchemaTemplateAnpEpgSubnetPrerequisiteConfig(), msoSchemaTemplateAnpEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateAnpEpgSubnetIp)
}

func testAccMSOSchemaTemplateAnpEpgSubnetConfigUpdateSharedAndDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp_epg_subnet" "%[2]s_subnet" {
		schema_id   = mso_schema_template_anp_epg.%[6]s.schema_id
		template    = "%[4]s"
		anp_name    = "%[5]s"
		epg_name    = mso_schema_template_anp_epg.%[6]s.name
		ip          = "%[7]s"
		scope       = "public"
		shared      = true
		description = "test subnet"
	}`, testAccMSOSchemaTemplateAnpEpgSubnetPrerequisiteConfig(), msoSchemaTemplateAnpEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateAnpEpgSubnetIp)
}

func testAccMSOSchemaTemplateAnpEpgSubnetConfigUpdateAllAttributes() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp_epg_subnet" "%[2]s_subnet" {
		schema_id          = mso_schema_template_anp_epg.%[6]s.schema_id
		template           = "%[4]s"
		anp_name           = "%[5]s"
		epg_name           = mso_schema_template_anp_epg.%[6]s.name
		ip                 = "%[7]s"
		scope              = "public"
		shared             = true
		description        = "test subnet"
		no_default_gateway = true
	}`, testAccMSOSchemaTemplateAnpEpgSubnetPrerequisiteConfig(), msoSchemaTemplateAnpEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateAnpEpgSubnetIp)
}

func testAccMSOSchemaTemplateAnpEpgSubnetConfigResetAttributes() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp_epg_subnet" "%[2]s_subnet" {
		schema_id          = mso_schema_template_anp_epg.%[6]s.schema_id
		template           = "%[4]s"
		anp_name           = "%[5]s"
		epg_name           = mso_schema_template_anp_epg.%[6]s.name
		ip                 = "%[7]s"
		scope              = "private"
		shared             = false
		description        = ""
		querier            = false
		no_default_gateway = false
		primary            = false
	}`, testAccMSOSchemaTemplateAnpEpgSubnetPrerequisiteConfig(), msoSchemaTemplateAnpEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateAnpEpgSubnetIp)
}

func testAccMSOSchemaTemplateAnpEpgSubnetConfigUpdateIp() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp_epg_subnet" "%[2]s_subnet" {
		schema_id          = mso_schema_template_anp_epg.%[6]s.schema_id
		template           = "%[4]s"
		anp_name           = "%[5]s"
		epg_name           = mso_schema_template_anp_epg.%[6]s.name
		ip                 = "%[7]s"
		scope              = "private"
		shared             = false
		querier            = false
		no_default_gateway = false
		primary            = false
	}`, testAccMSOSchemaTemplateAnpEpgSubnetPrerequisiteConfig(), msoSchemaTemplateAnpEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateAnpEpgSubnetIp2)
}

func testAccMSOSchemaTemplateAnpEpgSubnetConfigUpdateParentEpg() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_anp_epg" "%[2]s" {
		schema_id     = mso_schema.%[3]s.id
		template_name = "%[4]s"
		anp_name      = mso_schema_template_anp.%[5]s.name
		name          = "%[2]s"
		display_name  = "%[2]s"
		description   = "Updated EPG description with subnet"
		bd_name       = mso_schema_template_bd.%[6]s.name
	}
	resource "mso_schema_template_anp_epg_subnet" "%[2]s_subnet" {
		schema_id          = mso_schema_template_anp_epg.%[2]s.schema_id
		template           = "%[4]s"
		anp_name           = "%[5]s"
		epg_name           = mso_schema_template_anp_epg.%[2]s.name
		ip                 = "%[7]s"
		scope              = "private"
		shared             = false
		querier            = false
		no_default_gateway = false
		primary            = false
	}`, testSiteConfigAnsibleTest()+testTenantConfig()+testSchemaConfig()+testSchemaTemplateVrfConfig()+testSchemaTemplateBdConfig()+testSchemaTemplateAnpConfig(), msoSchemaTemplateAnpEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateBdName, msoSchemaTemplateAnpEpgSubnetIp2)
}

func testAccCheckMSOSchemaTemplateAnpEpgSubnetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "mso_schema_template_anp_epg_subnet" {
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
						subnetCount, err := epgCont.ArrayCount("subnets")
						if err != nil {
							return fmt.Errorf("Unable to get Subnet list")
						}
						for l := 0; l < subnetCount; l++ {
							subnetCont, err := epgCont.ArrayElement(l, "subnets")
							if err != nil {
								return err
							}
							ip := models.StripQuotes(subnetCont.S("ip").String())
							if rs.Primary.Attributes["ip"] == ip {
								return fmt.Errorf("Schema Template Anp Epg Subnet record still exists")
							}
						}
					}
				}
			}
		}
	}
	return nil
}
