package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// msoSchemaTemplateFilterEntrySchemaId is set during the first test step's Check to capture the dynamic schema ID for use in the manual deletion PreConfig step.
var msoSchemaTemplateFilterEntrySchemaId string

func TestAccMSOSchemaTemplateFilterEntryResource(t *testing.T) {
	filterEntryName := msoSchemaTemplateFilterName + "_entry"
	resourceName := "mso_schema_template_filter_entry." + msoSchemaTemplateFilterName
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateFilterEntryDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create filter entry with basic attributes") },
				Config:    testAccMSOSchemaTemplateFilterEntryConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttr(resourceName, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr(resourceName, "name", msoSchemaTemplateFilterName),
					resource.TestCheckResourceAttr(resourceName, "display_name", msoSchemaTemplateFilterName),
					resource.TestCheckResourceAttr(resourceName, "entry_name", filterEntryName),
					resource.TestCheckResourceAttr(resourceName, "entry_display_name", filterEntryName),
					resource.TestCheckResourceAttr(resourceName, "entry_description", ""),
					resource.TestCheckResourceAttr(resourceName, "ether_type", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "arp_flag", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "ip_protocol", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "match_only_fragments", "false"),
					resource.TestCheckResourceAttr(resourceName, "stateful", "false"),
					resource.TestCheckResourceAttr(resourceName, "source_from", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "source_to", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "destination_from", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "destination_to", "unspecified"),
					// Capture the dynamic schema ID from state for use in the manual deletion PreConfig step
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok {
							return fmt.Errorf("Filter entry resource not found in state")
						}
						msoSchemaTemplateFilterEntrySchemaId = rs.Primary.Attributes["schema_id"]
						return nil
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add filter entry description") },
				Config:    testAccMSOSchemaTemplateFilterEntryConfigAddDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttr(resourceName, "entry_name", filterEntryName),
					resource.TestCheckResourceAttr(resourceName, "entry_description", "Terraform test filter entry"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update filter entry description") },
				Config:    testAccMSOSchemaTemplateFilterEntryConfigUpdateDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttr(resourceName, "entry_name", filterEntryName),
					resource.TestCheckResourceAttr(resourceName, "entry_description", "Terraform test filter entry updated"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove filter entry description") },
				Config:    testAccMSOSchemaTemplateFilterEntryConfigRemoveDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttr(resourceName, "entry_name", filterEntryName),
					resource.TestCheckResourceAttr(resourceName, "entry_description", ""),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update ether_type to arp and arp_flag to req") },
				Config:    testAccMSOSchemaTemplateFilterEntryConfigUpdateArpFlag(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttr(resourceName, "entry_name", filterEntryName),
					resource.TestCheckResourceAttr(resourceName, "ether_type", "arp"),
					resource.TestCheckResourceAttr(resourceName, "arp_flag", "req"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update ether_type and ip_protocol") },
				Config:    testAccMSOSchemaTemplateFilterEntryConfigUpdateEtherIp(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttr(resourceName, "entry_name", filterEntryName),
					resource.TestCheckResourceAttr(resourceName, "ether_type", "ip"),
					resource.TestCheckResourceAttr(resourceName, "ip_protocol", "tcp"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update match_only_fragments to true") },
				Config:    testAccMSOSchemaTemplateFilterEntryConfigUpdateMatchOnlyFragments(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttr(resourceName, "entry_name", filterEntryName),
					resource.TestCheckResourceAttr(resourceName, "ether_type", "ip"),
					resource.TestCheckResourceAttr(resourceName, "ip_protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "match_only_fragments", "true"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update ports, match_only_fragments, and stateful") },
				Config:    testAccMSOSchemaTemplateFilterEntryConfigUpdatePorts(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttr(resourceName, "entry_name", filterEntryName),
					resource.TestCheckResourceAttr(resourceName, "ether_type", "ip"),
					resource.TestCheckResourceAttr(resourceName, "ip_protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "source_from", "http"),
					resource.TestCheckResourceAttr(resourceName, "source_to", "http"),
					resource.TestCheckResourceAttr(resourceName, "destination_from", "https"),
					resource.TestCheckResourceAttr(resourceName, "destination_to", "https"),
					resource.TestCheckResourceAttr(resourceName, "stateful", "true"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update tcp_session_rules") },
				Config:    testAccMSOSchemaTemplateFilterEntryConfigUpdateTcpSessionRules(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttr(resourceName, "entry_name", filterEntryName),
					resource.TestCheckResourceAttr(resourceName, "tcp_session_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tcp_session_rules.0", "established"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Add second filter entry to same filter") },
				Config:    testAccMSOSchemaTemplateFilterEntryConfigWithSecondEntry(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttr(resourceName, "entry_name", filterEntryName),
					resource.TestCheckResourceAttrSet("mso_schema_template_filter_entry."+msoSchemaTemplateFilterEntryName2, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_filter_entry."+msoSchemaTemplateFilterEntryName2, "name", msoSchemaTemplateFilterName),
					resource.TestCheckResourceAttr("mso_schema_template_filter_entry."+msoSchemaTemplateFilterEntryName2, "entry_name", msoSchemaTemplateFilterEntryName2),
					resource.TestCheckResourceAttr("mso_schema_template_filter_entry."+msoSchemaTemplateFilterEntryName2, "entry_display_name", msoSchemaTemplateFilterEntryName2),
					resource.TestCheckResourceAttr("mso_schema_template_filter_entry."+msoSchemaTemplateFilterEntryName2, "ether_type", "ip"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Remove first entry, keep second (filter should survive)") },
				Config:    testAccMSOSchemaTemplateFilterEntryConfigSecondEntryOnly(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_filter_entry."+msoSchemaTemplateFilterEntryName2, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_filter_entry."+msoSchemaTemplateFilterEntryName2, "name", msoSchemaTemplateFilterName),
					resource.TestCheckResourceAttr("mso_schema_template_filter_entry."+msoSchemaTemplateFilterEntryName2, "entry_name", msoSchemaTemplateFilterEntryName2),
					resource.TestCheckResourceAttr("mso_schema_template_filter_entry."+msoSchemaTemplateFilterEntryName2, "entry_display_name", msoSchemaTemplateFilterEntryName2),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import remaining second filter entry") },
				ResourceName: "mso_schema_template_filter_entry." + msoSchemaTemplateFilterEntryName2,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["mso_schema_template_filter_entry."+msoSchemaTemplateFilterEntryName2]
					if !ok {
						return "", fmt.Errorf("Filter entry resource not found in state")
					}
					return fmt.Sprintf("%s/templates/%s/filters/%s/entries/%s", rs.Primary.Attributes["schema_id"], rs.Primary.Attributes["template_name"], rs.Primary.Attributes["name"], rs.Primary.Attributes["entry_name"]), nil
				},
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate filter entry after manual deletion from NDO")
					msoClient := testAccProvider.Meta().(*client.Client)
					filterRemovePatchPayload := models.GetRemovePatchPayload(fmt.Sprintf("/templates/%s/filters/%s", msoSchemaTemplateName, msoSchemaTemplateFilterName))
					_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", msoSchemaTemplateFilterEntrySchemaId), filterRemovePatchPayload)
					if err != nil {
						t.Fatalf("Failed to manually delete filter: %v", err)
					}
				},
				Config: testAccMSOSchemaTemplateFilterEntryConfigSecondEntryOnly(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_filter_entry."+msoSchemaTemplateFilterEntryName2, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_filter_entry."+msoSchemaTemplateFilterEntryName2, "name", msoSchemaTemplateFilterName),
					resource.TestCheckResourceAttr("mso_schema_template_filter_entry."+msoSchemaTemplateFilterEntryName2, "entry_name", msoSchemaTemplateFilterEntryName2),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Delete last entry, filter should be deleted") },
				Config:    testAccMSOSchemaTemplateFilterEntryPrerequisiteConfig(),
				Check: func(s *terraform.State) error {
					client := testAccProvider.Meta().(*client.Client)
					cont, err := client.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", msoSchemaTemplateFilterEntrySchemaId))
					if err != nil {
						return err
					}
					count, err := cont.ArrayCount("templates")
					if err != nil {
						return fmt.Errorf("No Template found")
					}
					for i := 0; i < count; i++ {
						tempCont, err := cont.ArrayElement(i, "templates")
						if err != nil {
							return err
						}
						apiTemplate := models.StripQuotes(tempCont.S("name").String())
						if apiTemplate == msoSchemaTemplateName {
							filterCount, err := tempCont.ArrayCount("filters")
							if err != nil {
								return fmt.Errorf("Unable to get Filter list")
							}
							for j := 0; j < filterCount; j++ {
								filterCont, err := tempCont.ArrayElement(j, "filters")
								if err != nil {
									return err
								}
								if models.StripQuotes(filterCont.S("name").String()) == msoSchemaTemplateFilterName {
									return fmt.Errorf("Filter %s still exists after deleting last entry", msoSchemaTemplateFilterName)
								}
							}
						}
					}
					return nil
				},
			},
		},
	})
}

func testAccMSOSchemaTemplateFilterEntryPrerequisiteConfig() string {
	return fmt.Sprintf(`%s%s%s`, testSiteConfigAnsibleTest(), testTenantConfig(), testSchemaConfig())
}

func testAccMSOSchemaTemplateFilterEntryConfigCreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_filter_entry" "%[2]s" {
		schema_id          = mso_schema.%[3]s.id
		template_name      = "%[4]s"
		name               = "%[2]s"
		display_name       = "%[2]s"
		entry_name         = "%[2]s_entry"
		entry_display_name = "%[2]s_entry"
	}`, testAccMSOSchemaTemplateFilterEntryPrerequisiteConfig(), msoSchemaTemplateFilterName, msoSchemaName, msoSchemaTemplateName)
}

func testAccMSOSchemaTemplateFilterEntryConfigAddDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_filter_entry" "%[2]s" {
		schema_id          = mso_schema.%[3]s.id
		template_name      = "%[4]s"
		name               = "%[2]s"
		display_name       = "%[2]s"
		entry_name         = "%[2]s_entry"
		entry_display_name = "%[2]s_entry"
		entry_description  = "Terraform test filter entry"
	}`, testAccMSOSchemaTemplateFilterEntryPrerequisiteConfig(), msoSchemaTemplateFilterName, msoSchemaName, msoSchemaTemplateName)
}

func testAccMSOSchemaTemplateFilterEntryConfigUpdateDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_filter_entry" "%[2]s" {
		schema_id          = mso_schema.%[3]s.id
		template_name      = "%[4]s"
		name               = "%[2]s"
		display_name       = "%[2]s"
		entry_name         = "%[2]s_entry"
		entry_display_name = "%[2]s_entry"
		entry_description  = "Terraform test filter entry updated"
	}`, testAccMSOSchemaTemplateFilterEntryPrerequisiteConfig(), msoSchemaTemplateFilterName, msoSchemaName, msoSchemaTemplateName)
}

func testAccMSOSchemaTemplateFilterEntryConfigRemoveDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_filter_entry" "%[2]s" {
		schema_id          = mso_schema.%[3]s.id
		template_name      = "%[4]s"
		name               = "%[2]s"
		display_name       = "%[2]s"
		entry_name         = "%[2]s_entry"
		entry_display_name = "%[2]s_entry"
		entry_description  = ""
	}`, testAccMSOSchemaTemplateFilterEntryPrerequisiteConfig(), msoSchemaTemplateFilterName, msoSchemaName, msoSchemaTemplateName)
}

func testAccMSOSchemaTemplateFilterEntryConfigUpdateEtherIp() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_filter_entry" "%[2]s" {
		schema_id          = mso_schema.%[3]s.id
		template_name      = "%[4]s"
		name               = "%[2]s"
		display_name       = "%[2]s"
		entry_name         = "%[2]s_entry"
		entry_display_name = "%[2]s_entry"
		ether_type         = "ip"
		ip_protocol        = "tcp"
		arp_flag           = "unspecified"
	}`, testAccMSOSchemaTemplateFilterEntryPrerequisiteConfig(), msoSchemaTemplateFilterName, msoSchemaName, msoSchemaTemplateName)
}

func testAccMSOSchemaTemplateFilterEntryConfigUpdatePorts() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_filter_entry" "%[2]s" {
		schema_id            = mso_schema.%[3]s.id
		template_name        = "%[4]s"
		name                 = "%[2]s"
		display_name         = "%[2]s"
		entry_name           = "%[2]s_entry"
		entry_display_name   = "%[2]s_entry"
		ether_type           = "ip"
		ip_protocol          = "tcp"
		match_only_fragments = false
		source_from          = "http"
		source_to            = "http"
		destination_from     = "https"
		destination_to       = "https"
		stateful             = true
	}`, testAccMSOSchemaTemplateFilterEntryPrerequisiteConfig(), msoSchemaTemplateFilterName, msoSchemaName, msoSchemaTemplateName)
}

func testAccMSOSchemaTemplateFilterEntryConfigUpdateArpFlag() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_filter_entry" "%[2]s" {
		schema_id          = mso_schema.%[3]s.id
		template_name      = "%[4]s"
		name               = "%[2]s"
		display_name       = "%[2]s"
		entry_name         = "%[2]s_entry"
		entry_display_name = "%[2]s_entry"
		ether_type         = "arp"
		arp_flag           = "req"
	}`, testAccMSOSchemaTemplateFilterEntryPrerequisiteConfig(), msoSchemaTemplateFilterName, msoSchemaName, msoSchemaTemplateName)
}

func testAccMSOSchemaTemplateFilterEntryConfigUpdateMatchOnlyFragments() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_filter_entry" "%[2]s" {
		schema_id            = mso_schema.%[3]s.id
		template_name        = "%[4]s"
		name                 = "%[2]s"
		display_name         = "%[2]s"
		entry_name           = "%[2]s_entry"
		entry_display_name   = "%[2]s_entry"
		ether_type           = "ip"
		ip_protocol          = "tcp"
		match_only_fragments = true
	}`, testAccMSOSchemaTemplateFilterEntryPrerequisiteConfig(), msoSchemaTemplateFilterName, msoSchemaName, msoSchemaTemplateName)
}

func testAccMSOSchemaTemplateFilterEntryConfigUpdateTcpSessionRules() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_filter_entry" "%[2]s" {
		schema_id          = mso_schema.%[3]s.id
		template_name      = "%[4]s"
		name               = "%[2]s"
		display_name       = "%[2]s"
		entry_name         = "%[2]s_entry"
		entry_display_name = "%[2]s_entry"
		ether_type         = "ip"
		ip_protocol        = "tcp"
		source_from        = "http"
		source_to          = "http"
		destination_from   = "https"
		destination_to     = "https"
		stateful           = true
		tcp_session_rules  = ["established"]
	}`, testAccMSOSchemaTemplateFilterEntryPrerequisiteConfig(), msoSchemaTemplateFilterName, msoSchemaName, msoSchemaTemplateName)
}

func testAccMSOSchemaTemplateFilterEntryConfigWithSecondEntry() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_filter_entry" "%[2]s" {
		schema_id          = mso_schema.%[3]s.id
		template_name      = "%[4]s"
		name               = "%[2]s"
		display_name       = "%[2]s"
		entry_name         = "%[2]s_entry"
		entry_display_name = "%[2]s_entry"
		ether_type         = "ip"
		ip_protocol        = "tcp"
		source_from        = "http"
		source_to          = "http"
		destination_from   = "https"
		destination_to     = "https"
		stateful           = true
	}
	resource "mso_schema_template_filter_entry" "%[5]s" {
		schema_id          = mso_schema.%[3]s.id
		template_name      = "%[4]s"
		name               = "%[2]s"
		display_name       = "%[2]s"
		entry_name         = "%[5]s"
		entry_display_name = "%[5]s"
		ether_type         = "ip"
		depends_on         = [mso_schema_template_filter_entry.%[2]s]
	}`, testAccMSOSchemaTemplateFilterEntryPrerequisiteConfig(), msoSchemaTemplateFilterName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateFilterEntryName2)
}

func testAccMSOSchemaTemplateFilterEntryConfigSecondEntryOnly() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_schema_template_filter_entry" "%[5]s" {
		schema_id          = mso_schema.%[3]s.id
		template_name      = "%[4]s"
		name               = "%[2]s"
		display_name       = "%[2]s"
		entry_name         = "%[5]s"
		entry_display_name = "%[5]s"
		ether_type         = "ip"
	}`, testAccMSOSchemaTemplateFilterEntryPrerequisiteConfig(), msoSchemaTemplateFilterName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateFilterEntryName2)
}

func testAccCheckMSOSchemaTemplateFilterEntryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "mso_schema_template_filter_entry" {
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
				filterCount, err := tempCont.ArrayCount("filters")
				if err != nil {
					return fmt.Errorf("Unable to get Filter list")
				}
				for j := 0; j < filterCount; j++ {
					filterCont, err := tempCont.ArrayElement(j, "filters")
					if err != nil {
						return err
					}
					apiFilterName := models.StripQuotes(filterCont.S("name").String())
					if apiFilterName == rs.Primary.Attributes["name"] {
						entryCount, err := filterCont.ArrayCount("entries")
						if err != nil {
							return fmt.Errorf("Unable to get Entry list")
						}
						for k := 0; k < entryCount; k++ {
							entryCont, err := filterCont.ArrayElement(k, "entries")
							if err != nil {
								return err
							}
							apiEntryName := models.StripQuotes(entryCont.S("name").String())
							if apiEntryName == rs.Primary.Attributes["entry_name"] {
								return fmt.Errorf("Schema Template Filter Entry record still exists")
							}
						}
					}
				}
			}
		}
	}
	return nil
}
