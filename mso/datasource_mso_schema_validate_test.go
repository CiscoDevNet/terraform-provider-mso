package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOSchemaValidate_DataSource(t *testing.T) {
	dataSourceName := "data.mso_schema_validate.test"
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      MSOSchemaValidateDataSourceWithoutSchameId(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaValidateDataSourceWithRequired(inValidScheamaId),
				ExpectError: regexp.MustCompile(`Bad Request`),
			},
			{
				Config:      MSOSchemaValidateDataSourceWithRequired(randomValue),
				ExpectError: regexp.MustCompile(`Resource Not Found`),
			},
			{
				Config:      MSOSchemaValidateDataSourceRandomAttr(validSchemaId, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)+is not expected here.`),
			},
			{
				Config: MSOSchemaValidateDataSourceWithRequired(validSchemaId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "schema_id", validSchemaId),
					resource.TestCheckResourceAttr(dataSourceName, "result", "true"),
				),
			},
		},
	})
}

func MSOSchemaValidateDataSourceWithoutSchameId() string {
	resource := fmt.Sprintln(`
	data "mso_schema_validate" "test" {}
	`)
	return resource
}

func MSOSchemaValidateDataSourceWithRequired(id string) string {
	resource := fmt.Sprintf(`
	data "mso_schema_validate" "test" {
		schema_id = "%s"
	}
	`, id)
	return resource
}

func MSOSchemaValidateDataSourceRandomAttr(id, key, value string) string {
	resource := fmt.Sprintf(`
	data "mso_schema_validate" "test" {
		schema_id = "%s"
		%s = "%s"
	}
	`, id, key, value)
	return resource
}
