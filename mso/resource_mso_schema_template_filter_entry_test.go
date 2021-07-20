package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccMSOSchemaTemplateFilterEntry_Basic(t *testing.T) {
	var ss TemplateFilterEntry
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateFilterEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateFilterEntryConfig_basic("unspecified"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateFilterEntryExists("mso_schema.schema1", "mso_schema_template_filter_entry.filter_entry", &ss),
					testAccCheckMSOSchemaTemplateFilterEntryAttributes("unspecified", &ss),
				),
			},
			{
				ResourceName:      "mso_schema_template_filter_entry.filter_entry",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMSOSchemaTemplateFilterEntry_Update(t *testing.T) {
	var ss TemplateFilterEntry

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateFilterEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOTemplateFilterEntryConfig_basic("unspecified"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateFilterEntryExists("mso_schema.schema1", "mso_schema_template_filter_entry.filter_entry", &ss),
					testAccCheckMSOSchemaTemplateFilterEntryAttributes("unspecified", &ss),
				),
			},
			{
				Config: testAccCheckMSOTemplateFilterEntryConfig_basic("trill"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOSchemaTemplateFilterEntryExists("mso_schema.schema1", "mso_schema_template_filter_entry.filter_entry", &ss),
					testAccCheckMSOSchemaTemplateFilterEntryAttributes("trill", &ss),
				),
			},
		},
	})
}

func testAccCheckMSOTemplateFilterEntryConfig_basic(ethertype string) string {
	return fmt.Sprintf(`

	resource "mso_schema" "schema1" {
		name          = "Schema2"
		template_name = "Template1"
		tenant_id     = "5fb5fed8520000452a9e8911"
	  
	  }

	resource "mso_schema_template_filter_entry" "filter_entry" {
		schema_id = mso_schema.schema1.id
		template_name = "Template1"
		name = "Any"
		display_name="Any"
		entry_name = "testAcc"
		entry_display_name="testAcc"
		entry_description="test1"
		destination_from="unspecified"
		destination_to="unspecified"
		source_from="unspecified"
		source_to="unspecified"
		arp_flag="unspecified"
		ip_protocol="unspecified"
		match_only_fragments=false
		tcp_session_rules=[
			"unspecified"
		]
		 
	}

	resource "mso_schema_template_filter_entry" "filter_entry2" {
		schema_id = mso_schema.schema1.id
		template_name = mso_schema_template_filter_entry.filter_entry.template_name
		name = "One"
		display_name="One"
		entry_name = "testAcc2"
		entry_display_name="testAcc2"
		entry_description="test2"
		destination_from="unspecified"
		destination_to="unspecified"
		source_from="unspecified"
		source_to="unspecified"
		arp_flag="unspecified"
		ip_protocol="unspecified"
		match_only_fragments=false
		ether_type="unspecified"
		tcp_session_rules=[
			"unspecified"
		]
		 
	}

	resource "mso_schema_template_filter_entry" "filter_entry3" {
		schema_id = mso_schema.schema1.id
		template_name = mso_schema_template_filter_entry.filter_entry2.template_name
		name = "One"
		display_name="One"
		entry_name = "testAcc2"
		entry_display_name="testAcc2"
		entry_description="test2"
		destination_from="unspecified"
		destination_to="unspecified"
		source_from="unspecified"
		source_to="unspecified"
		arp_flag="unspecified"
		ip_protocol="unspecified"
		match_only_fragments=false
		ether_type="%s"
		tcp_session_rules=[
			"unspecified"
		]
		 
	}
`, ethertype)
}

func testAccCheckMSOSchemaTemplateFilterEntryExists(schemaName string, filterName string, ss *TemplateFilterEntry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		rs1, err1 := s.RootModule().Resources[schemaName]
		rs2, err2 := s.RootModule().Resources[filterName]
		if !err1 {
			return fmt.Errorf("Schema %s not found", schemaName)
		}
		if !err2 {
			return fmt.Errorf("Entry %s not found", filterName)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Schema id was set")
		}
		if rs2.Primary.ID == "" {
			return fmt.Errorf("No Filter was set")
		}

		cont, err := client.GetViaURL("api/v1/schemas/" + rs1.Primary.ID)
		if err != nil {
			return err
		}

		count, err := cont.ArrayCount("templates")
		if err != nil {
			return fmt.Errorf("No Template found")
		}
		tp := TemplateFilterEntry{}
		found := false
		for i := 0; i < count; i++ {
			tempCont, err := cont.ArrayElement(i, "templates")
			if err != nil {
				return err
			}

			apiTemplateName := models.StripQuotes(tempCont.S("name").String())
			if apiTemplateName == "Template1" {
				bdCount, err := tempCont.ArrayCount("filters")
				if err != nil {
					return fmt.Errorf("Unable to get Filter list")
				}
				for j := 0; j < bdCount; j++ {
					bdCont, err := tempCont.ArrayElement(j, "filters")
					if err != nil {
						return err
					}
					apiFilter := models.StripQuotes(bdCont.S("name").String())
					if apiFilter == "Any" {
						entryCount, err := bdCont.ArrayCount("entries")
						if err != nil {
							return fmt.Errorf("Unable to get Entry list")
						}
						for k := 0; k < entryCount; k++ {
							entryCont, err := bdCont.ArrayElement(k, "entries")
							if err != nil {
								return err
							}
							apiFilterEntry := models.StripQuotes(entryCont.S("name").String())
							if apiFilterEntry == "testAcc" {
								tp.entry_display_name = models.StripQuotes(entryCont.S("displayName").String())
								tp.arp_flag = models.StripQuotes(entryCont.S("arpFlag").String())
								tp.ip_protocol = models.StripQuotes(entryCont.S("ipProtocol").String())
								tp.ether_type = models.StripQuotes(entryCont.S("etherType").String())
								found = true
								break
							}
						}

					}
				}
			}
		}

		if !found {
			return fmt.Errorf("Entry not found from API")
		}

		tp1 := &tp

		*ss = *tp1
		return nil
	}
}

func testAccCheckMSOSchemaTemplateFilterEntryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	rs1, err1 := s.RootModule().Resources["mso_schema.schema1"]
	if !err1 {
		return fmt.Errorf("Schema %s not found", "mso_schema.schema1")
	}
	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_schema_template_filter_entry" {
			cont, err := client.GetViaURL("api/v1/schemas/" + rs1.Primary.ID)
			if err != nil {
				return nil
			} else {
				count, err := cont.ArrayCount("templates")
				if err != nil {
					return fmt.Errorf("No Template found")
				}
				for i := 0; i < count; i++ {
					tempCont, err := cont.ArrayElement(i, "templates")
					if err != nil {
						return fmt.Errorf("No Template exists")
					}
					apiTemplateName := models.StripQuotes(tempCont.S("name").String())
					if apiTemplateName == "Template1" {
						bdCount, err := tempCont.ArrayCount("filters")
						if err != nil {
							return fmt.Errorf("Unable to get Filter list")
						}
						for j := 0; j < bdCount; j++ {
							bdCont, err := tempCont.ArrayElement(j, "filters")
							if err != nil {
								return err
							}
							apiFilter := models.StripQuotes(bdCont.S("name").String())
							if apiFilter == "Any" {
								entryCount, err := bdCont.ArrayCount("entries")
								if err != nil {
									return fmt.Errorf("Unable to get Entry list")
								}
								for k := 0; k < entryCount; k++ {
									entryCont, err := bdCont.ArrayElement(k, "entries")
									if err != nil {
										return err
									}
									apiFilterEntry := models.StripQuotes(entryCont.S("name").String())
									if apiFilterEntry == "testAcc" {
										return fmt.Errorf("Template Filter Entry still exists.")
									}
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

func testAccCheckMSOSchemaTemplateFilterEntryAttributes(ethertype string, ss *TemplateFilterEntry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// if ethertype != ss.ether_type {
		// 	return fmt.Errorf("Bad Template Filter Entry Ether Type %s", ss.ether_type)
		// }

		if "unspecified" != ss.ip_protocol {
			return fmt.Errorf("Bad Template Filter Entry Ip Protocol %s", ss.ip_protocol)
		}

		if "unspecified" != ss.arp_flag {
			return fmt.Errorf("Bad Template Filter Entry ARP Flag %s", ss.arp_flag)
		}
		return nil
	}
}

type TemplateFilterEntry struct {
	entry_display_name string
	ether_type         string
	arp_flag           string
	ip_protocol        string
}
