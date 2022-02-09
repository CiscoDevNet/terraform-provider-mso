package mso

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSODHCPOptionPolicy_Basic(t *testing.T) {
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
