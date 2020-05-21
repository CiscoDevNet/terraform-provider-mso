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

func TestAccMSOLabel_Basic(t *testing.T) {
	var s LabelTest
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOLabelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMSOLabelConfig_basic("site"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMSOLabelExists("mso_label.label1", &s),
					testAccCheckMSOLabelAttributes("site", &s),
				),
			},
		},
	})
}

func testAccCheckMSOLabelConfig_basic(types string) string {
	return fmt.Sprintf(`
	resource "mso_label" "label1" {
	 label = "hello4"
	 type  = "%s"
		 }
	`, types)
}

func testAccCheckMSOLabelExists(userName string, st *LabelTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs1, err1 := s.RootModule().Resources[userName]

		if !err1 {
			return fmt.Errorf("Label %s not found", userName)
		}

		if rs1.Primary.ID == "" {
			return fmt.Errorf("No Label id was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.GetViaURL("api/v1/labels/" + rs1.Primary.ID)

		if err != nil {
			return err
		}

		sts, _ := labelFromcontainer(cont)

		*st = *sts
		return nil
	}
}

func labelFromcontainer(con *container.Container) (*LabelTest, error) {

	s := LabelTest{}

	s.DisplayName = models.StripQuotes(con.S("displayName").String())
	s.Type = models.StripQuotes(con.S("type").String())

	return &s, nil
}

func testAccCheckMSOLabelDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "mso_label" {
			_, err := client.GetViaURL("api/v1/labels/" + rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Label still exists")
			}
		} else {
			continue
		}

	}
	return nil
}

func testAccCheckMSOLabelAttributes(name string, st *LabelTest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if name != st.Type {
			return fmt.Errorf("Bad Type %s", st.Type)
		}
		return nil
	}
}

type LabelTest struct {
	Id          string `json:",omitempty"`
	DisplayName string `json:",omitempty"`
	Type        string `json:",omitempty"`
}
