package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSite_Basic(t *testing.T) {
	var s SiteTest
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMsoSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMsoSiteConfig_basic("mso123"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMsoSiteExists("mso_site.site1", &s),
					testAccCheckMsoSiteAttributes("mso123", &s),
				),
			},
		},
	})
}

func TestAccMsoSite_Update(t *testing.T) {
	var s SiteTest

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMsoSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMsoSiteConfig_basic("mso123"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMsoSiteExists("mso_site.site1", &s),
					testAccCheckMsoSiteAttributes("mso123", &s),
				),
			},
			{
				Config: testAccCheckMsoSiteConfig_basic("mso234"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMsoSiteExists("mso_site.site1", &s),
					testAccCheckMsoSiteAttributes("mso234", &s),
				),
			},
		},
	})
}

func testAccCheckMsoSiteConfig_basic(name string) string {
	return fmt.Sprintf(`
	resource "mso_site" "site1" {
		name = "%v"
		username = "admin"
		password = "noir0!234"
		apic_site_id = "18"
		urls = [ "https://3.208.123.222" ]
	}
	`, name)
}

func testAccCheckMsoSiteExists(siteName string, st *SiteTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs1, err1 := s.RootModule().Resources[siteName]

		if !err1 {
			return fmt.Errorf("Site record %s not found", siteName)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Site record id was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		resp, err := client.GetViaURL("api/v1/sites/" + rs1.Primary.ID)

		if err != nil {
			return err
		}

		tp, _ := sitefromcontainer(resp)

		*st = *tp
		return nil
	}
}

func sitefromcontainer(con *container.Container) (*SiteTest, error) {

	s := SiteTest{}
	s.Name = models.StripQuotes(con.S("name").String())
	s.ApicUsername = models.StripQuotes(con.S("username").String())

	s.ApicSiteId = models.StripQuotes(con.S("apic_site_id").String())
	s.Labels = con.S("labels").Data().([]interface{})
	s.Url = con.S("urls").Data().([]interface{})
	s.CloudProviders = con.S("cloudProviders").Data().([]interface{})

	return &s, nil

}

func testAccCheckMsoSiteDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_site" {
			_, err := client.GetViaURL("api/v1/sites/" + rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Site still exists")
			}
		} else {
			continue
		}

	}
	return nil
}
func testAccCheckMsoSiteAttributes(name string, st *SiteTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if name != st.Name {
			return fmt.Errorf("Bad Site Name %s", st.Name)
		}
		return nil
	}
}

type Location struct {
	Lat  float64 `json:"lat,omitempty"`
	Long float64 `json:"long,omitempty"`
}

type SiteTest struct {
	Id             string        `json:",omitempty"`
	Name           string        `json:",omitempty"`
	ApicUsername   string        `json:",omitempty"`
	ApicPassword   string        `json:",omitempty"`
	ApicSiteId     string        `json:",omitempty"`
	Labels         []interface{} `json:",omitempty"`
	Location       []interface{} `json:",omitempty"`
	Url            []interface{} `json:",omitempty"`
	CloudProviders []interface{} `json:",omitempty"`
}
