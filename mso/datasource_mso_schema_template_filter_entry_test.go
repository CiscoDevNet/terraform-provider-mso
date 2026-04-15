package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOSchemaTemplateFilterEntryDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateFilterEntryDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Read filter entry datasource not found error") },
				Config:      testAccMSOSchemaTemplateFilterEntryDatasourceNotFound(),
				ExpectError: regexp.MustCompile("Unable to find the Filter"),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read filter entry datasource") },
				Config:    testAccMSOSchemaTemplateFilterEntryDatasource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.mso_schema_template_filter_entry.filter_entry", "schema_id"),
					resource.TestCheckResourceAttr("data.mso_schema_template_filter_entry.filter_entry", "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("data.mso_schema_template_filter_entry.filter_entry", "name", msoSchemaTemplateFilterName),
					resource.TestCheckResourceAttr("data.mso_schema_template_filter_entry.filter_entry", "display_name", msoSchemaTemplateFilterName),
					resource.TestCheckResourceAttr("data.mso_schema_template_filter_entry.filter_entry", "entry_name", msoSchemaTemplateFilterName+"_entry"),
					resource.TestCheckResourceAttr("data.mso_schema_template_filter_entry.filter_entry", "entry_display_name", msoSchemaTemplateFilterName+"_entry"),
					resource.TestCheckResourceAttr("data.mso_schema_template_filter_entry.filter_entry", "ether_type", "unspecified"),
					resource.TestCheckResourceAttr("data.mso_schema_template_filter_entry.filter_entry", "ip_protocol", "unspecified"),
					resource.TestCheckResourceAttr("data.mso_schema_template_filter_entry.filter_entry", "arp_flag", "unspecified"),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateFilterEntryDatasource() string {
	return fmt.Sprintf(`%s
	data "mso_schema_template_filter_entry" "filter_entry" {
		schema_id     = mso_schema.%[2]s.id
		template_name = "%[3]s"
		name          = mso_schema_template_filter_entry.%[4]s.name
		entry_name    = mso_schema_template_filter_entry.%[4]s.entry_name
	}`, testAccMSOSchemaTemplateFilterEntryConfigCreate(), msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateFilterName)
}

func testAccMSOSchemaTemplateFilterEntryDatasourceNotFound() string {
	return fmt.Sprintf(`%s
	data "mso_schema_template_filter_entry" "filter_entry" {
		schema_id     = mso_schema.%[2]s.id
		template_name = "%[3]s"
		name          = mso_schema_template_filter_entry.%[4]s.name
		entry_name    = "non_existing_entry_name"
	}`, testAccMSOSchemaTemplateFilterEntryConfigCreate(), msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateFilterName)
}
