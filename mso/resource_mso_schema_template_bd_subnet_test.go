package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// msoSchemaTemplateBdSubnetSchemaId is set during the first test step's Check to capture the dynamic schema ID for use in the manual deletion PreConfig step.
var msoSchemaTemplateBdSubnetSchemaId string

func TestAccMSOSchemaTemplateBdSubnetResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create BD Subnet with required attributes") },
				Config:    testAccMSOSchemaTemplateBdSubnetConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "bd_name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "ip", msoSchemaTemplateBdSubnetIp),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "scope", "private"),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "shared", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "no_default_gateway", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "querier", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "primary", "false"),
					// Capture the dynamic schema ID from state for use in the manual deletion PreConfig step
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet"]
						if !ok {
							return fmt.Errorf("BD Subnet resource not found in state")
						}
						msoSchemaTemplateBdSubnetSchemaId = rs.Primary.Attributes["schema_id"]
						return nil
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update BD Subnet scope to public") },
				Config:    testAccMSOSchemaTemplateBdSubnetConfigUpdateScope(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "ip", msoSchemaTemplateBdSubnetIp),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "scope", "public"),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "shared", "false"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update BD Subnet shared and querier") },
				Config:    testAccMSOSchemaTemplateBdSubnetConfigUpdateSharedAndQuerier(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "ip", msoSchemaTemplateBdSubnetIp),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "scope", "public"),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "shared", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "querier", "true"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update BD Subnet - set all attributes") },
				Config:    testAccMSOSchemaTemplateBdSubnetConfigUpdateAllAttributes(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "ip", msoSchemaTemplateBdSubnetIp),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "scope", "public"),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "shared", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "querier", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "no_default_gateway", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "primary", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "virtual", "true"),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "description", "Terraform test BD Subnet"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Reset BD Subnet optional attributes to defaults") },
				Config:    testAccMSOSchemaTemplateBdSubnetConfigResetAttributes(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "ip", msoSchemaTemplateBdSubnetIp),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "scope", "private"),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "shared", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "querier", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "no_default_gateway", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "primary", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "virtual", "false"),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "description", ""),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import BD Subnet") },
				ResourceName: "mso_schema_template_bd_subnet." + msoSchemaTemplateBdName + "_subnet",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet"]
					if !ok {
						return "", fmt.Errorf("BD Subnet resource not found in state")
					}
					return fmt.Sprintf("%s/templates/%s/bds/%s/ip/%s",
						rs.Primary.Attributes["schema_id"],
						rs.Primary.Attributes["template_name"],
						rs.Primary.Attributes["bd_name"],
						rs.Primary.Attributes["ip"],
					), nil
				},
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate BD Subnet after manual deletion from NDO")
					msoClient := testAccProvider.Meta().(*client.Client)
					cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", msoSchemaTemplateBdSubnetSchemaId))
					if err != nil {
						t.Fatalf("Failed to get schema: %v", err)
					}
					index, err := findBdSubnetIndex(cont, msoSchemaTemplateName, msoSchemaTemplateBdName, msoSchemaTemplateBdSubnetIp)
					if err != nil {
						t.Fatalf("Failed to fetch BD subnet index: %v", err)
					}
					if index == -1 {
						t.Fatalf("BD Subnet not found for manual deletion")
					}
					subnetRemovePatchPayload := models.GetRemovePatchPayload(fmt.Sprintf("/templates/%s/bds/%s/subnets/%d", msoSchemaTemplateName, msoSchemaTemplateBdName, index))
					_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", msoSchemaTemplateBdSubnetSchemaId), subnetRemovePatchPayload)
					if err != nil {
						t.Fatalf("Failed to manually delete BD Subnet: %v", err)
					}
				},
				Config: testAccMSOSchemaTemplateBdSubnetConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "ip", msoSchemaTemplateBdSubnetIp),
					resource.TestCheckResourceAttr("mso_schema_template_bd_subnet."+msoSchemaTemplateBdName+"_subnet", "scope", "private"),
				),
			},
		},
		CheckDestroy: testAccCheckMSOSchemaTemplateBdSubnetDestroy,
	})
}

func testAccMSOSchemaTemplateBdSubnetPrerequisiteConfig() string {
	return fmt.Sprintf(`%s%s%s%s%s`, testSiteConfigAnsibleTest(), testTenantConfig(), testSchemaConfig(), testSchemaTemplateVrfConfig(), testSchemaTemplateBdStretchedConfig())
}

func testAccMSOSchemaTemplateBdSubnetConfigCreate() string {
	return testAccMSOSchemaTemplateBdSubnetPrerequisiteConfig() + testSchemaTemplateBdSubnetConfig()
}

func testAccMSOSchemaTemplateBdSubnetConfigUpdateScope() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_template_bd_subnet" "%[2]s_subnet" {
	schema_id     = mso_schema.%[3]s.id
	template_name = "%[4]s"
	bd_name       = mso_schema_template_bd.%[2]s.name
	ip            = "%[5]s"
	scope         = "public"
	shared        = false
}`, testAccMSOSchemaTemplateBdSubnetPrerequisiteConfig(), msoSchemaTemplateBdName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateBdSubnetIp)
}

func testAccMSOSchemaTemplateBdSubnetConfigUpdateSharedAndQuerier() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_template_bd_subnet" "%[2]s_subnet" {
	schema_id     = mso_schema.%[3]s.id
	template_name = "%[4]s"
	bd_name       = mso_schema_template_bd.%[2]s.name
	ip            = "%[5]s"
	scope         = "public"
	shared        = true
	querier       = true
}`, testAccMSOSchemaTemplateBdSubnetPrerequisiteConfig(), msoSchemaTemplateBdName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateBdSubnetIp)
}

func testAccMSOSchemaTemplateBdSubnetConfigUpdateAllAttributes() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_template_bd_subnet" "%[2]s_subnet" {
	schema_id          = mso_schema.%[3]s.id
	template_name      = "%[4]s"
	bd_name            = mso_schema_template_bd.%[2]s.name
	ip                 = "%[5]s"
	scope              = "public"
	shared             = true
	querier            = true
	no_default_gateway = true
	primary            = true
	virtual            = true
	description        = "Terraform test BD Subnet"
}`, testAccMSOSchemaTemplateBdSubnetPrerequisiteConfig(), msoSchemaTemplateBdName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateBdSubnetIp)
}

func testAccMSOSchemaTemplateBdSubnetConfigResetAttributes() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_template_bd_subnet" "%[2]s_subnet" {
	schema_id          = mso_schema.%[3]s.id
	template_name      = "%[4]s"
	bd_name            = mso_schema_template_bd.%[2]s.name
	ip                 = "%[5]s"
	scope              = "private"
	shared             = false
	querier            = false
	no_default_gateway = false
	primary            = false
	virtual            = false
	description        = ""
}`, testAccMSOSchemaTemplateBdSubnetPrerequisiteConfig(), msoSchemaTemplateBdName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateBdSubnetIp)
}

// findBdSubnetIndex returns the index of the subnet with the given IP in the BD, or -1 if not found.
func findBdSubnetIndex(cont *container.Container, templateName, bdName, ip string) (int, error) {
	found := false
	index := -1
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return index, fmt.Errorf("No Template found")
	}
	for i := 0; i < count && !found; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return index, err
		}
		currentTemplateName := models.StripQuotes(tempCont.S("name").String())
		if currentTemplateName == templateName {
			bdCount, err := tempCont.ArrayCount("bds")
			if err != nil {
				return index, fmt.Errorf("Unable to get BD list")
			}
			for j := 0; j < bdCount && !found; j++ {
				bdCont, err := tempCont.ArrayElement(j, "bds")
				if err != nil {
					return index, err
				}
				currentBDName := models.StripQuotes(bdCont.S("name").String())
				if currentBDName == bdName {
					subnetCount, err := bdCont.ArrayCount("subnets")
					if err != nil {
						return index, fmt.Errorf("Unable to get Subnet list")
					}
					for k := 0; k < subnetCount; k++ {
						subnetCont, err := bdCont.ArrayElement(k, "subnets")
						if err != nil {
							return index, err
						}
						currentIP := models.StripQuotes(subnetCont.S("ip").String())
						if currentIP == ip {
							index = k
							found = true
							break
						}
					}
				}
			}
		}
	}
	return index, nil
}

func testAccCheckMSOSchemaTemplateBdSubnetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "mso_schema_template_bd_subnet" {
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
				apiTemplate := models.StripQuotes(tempCont.S("name").String())
				if apiTemplate == rs.Primary.Attributes["template_name"] {
					bdCount, err := tempCont.ArrayCount("bds")
					if err != nil {
						return fmt.Errorf("Unable to get BD list")
					}
					for j := 0; j < bdCount; j++ {
						bdCont, err := tempCont.ArrayElement(j, "bds")
						if err != nil {
							return err
						}
						apiBD := models.StripQuotes(bdCont.S("name").String())
						if apiBD == rs.Primary.Attributes["bd_name"] {
							subnetCount, err := bdCont.ArrayCount("subnets")
							if err != nil {
								return fmt.Errorf("Unable to get BD subnet list")
							}
							for k := 0; k < subnetCount; k++ {
								subnetCont, err := bdCont.ArrayElement(k, "subnets")
								if err != nil {
									return err
								}
								apiIP := models.StripQuotes(subnetCont.S("ip").String())
								if apiIP == rs.Primary.Attributes["ip"] {
									return fmt.Errorf("Schema Template BD Subnet still exists")
								}
							}
						}
					}
				}
			}
		}
	}
	return nil
}
