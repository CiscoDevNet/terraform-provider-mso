package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOSchemaSiteServiceGraphNode_DataSource(t *testing.T) {
	resourceName := "mso_schema_site_service_graph_node.test"
	datasourceName := "data.mso_schema_site_service_graph_node.test"
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	name := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMsoSchemaSiteServiceGraphNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config:      MSOSchemaSiteServiceGraphNodeDataSourceWithoutRequired(name, "schema_id"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteServiceGraphNodeDataSourceWithoutRequired(name, "template_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteServiceGraphNodeDataSourceWithoutRequired(name, "service_graph_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteServiceGraphNodeDataSourceWithoutRequired(name, "service_node_type"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteServiceGraphNodeDataSourceWithoutRequired(name, "service_node_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      MSOSchemaSiteServiceGraphNodeDataSourceForRandomAttrName(name, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)+is not expected here.`),
			},
			{
				Config:      MSOSchemaSiteServiceGraphNodeDataSourceWithInvalidParentReference(name),
				ExpectError: regexp.MustCompile(`Resource Not Found`),
			},
			{
				Config: MSOSchemaSiteServiceGraphNodeDataSourceRequired(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "schema_id", datasourceName, "schema_id"),
					resource.TestCheckResourceAttrPair(resourceName, "template_name", datasourceName, "template_name"),
					resource.TestCheckResourceAttrPair(resourceName, "service_graph_name", datasourceName, "service_graph_name"),
					resource.TestCheckResourceAttrPair(resourceName, "service_node_type", datasourceName, "service_node_type"),
					resource.TestCheckResourceAttr(datasourceName, "service_node_name", "tfnode2"),
				),
			},
		},
	})
}

func MSOSchemaSiteServiceGraphNodeDataSourceWithInvalidParentReference(name string) string {
	resource := CreateServiceNodeResource(name)
	resource += fmt.Sprintln(`
	data "mso_schema_site_service_graph_node" "test" {
		schema_id          = "${mso_schema_site_service_graph_node.test.schema_id}_invalid"
		template_name      = mso_schema_site_service_graph_node.test.template_name
		service_graph_name = mso_schema_site_service_graph_node.test.service_graph_name
		service_node_type  = mso_schema_site_service_graph_node.test.service_node_type
		service_node_name  = "tfnode2"
	  }
	`)
	return resource
}

func MSOSchemaSiteServiceGraphNodeDataSourceForRandomAttrName(name, key, value string) string {
	resource := CreateServiceNodeResource(name)
	resource += fmt.Sprintf(`
	data "mso_schema_site_service_graph_node" "test" {
		schema_id          = mso_schema_site_service_graph_node.test.schema_id
		template_name      = mso_schema_site_service_graph_node.test.template_name
		service_graph_name = mso_schema_site_service_graph_node.test.service_graph_name
		service_node_type  = mso_schema_site_service_graph_node.test.service_node_type
		service_node_name  = "tfnode2"
		%s				   = "%s"
	  }
	`, key, value)
	return resource
}

func MSOSchemaSiteServiceGraphNodeDataSourceRequired(name string) string {
	resource := CreateServiceNodeResource(name)
	resource += fmt.Sprintln(`
	data "mso_schema_site_service_graph_node" "test" {
		schema_id          = mso_schema_site_service_graph_node.test.schema_id
		template_name      = mso_schema_site_service_graph_node.test.template_name
		service_graph_name = mso_schema_site_service_graph_node.test.service_graph_name
		service_node_type  = mso_schema_site_service_graph_node.test.service_node_type
		service_node_name  = "tfnode2"
	  }
	`)
	return resource
}

func MSOSchemaSiteServiceGraphNodeDataSourceWithoutRequired(name, attr string) string {
	resource := CreateServiceNodeResource(name)
	switch attr {
	case "schema_id":
		resource += `
		data "mso_schema_site_service_graph_node" "test" {
		#	schema_id          = mso_schema_site_service_graph_node.test.schema_id
			template_name      = mso_schema_site_service_graph_node.test.template_name
			service_graph_name = mso_schema_site_service_graph_node.test.service_graph_name
			service_node_type  = mso_schema_site_service_graph_node.test.service_node_type
			service_node_name  = "tfnode2"
		  }
		`
	case "template_name":
		resource += `
		data "mso_schema_site_service_graph_node" "test" {
			schema_id          = mso_schema_site_service_graph_node.test.schema_id
		#	template_name      = mso_schema_site_service_graph_node.test.template_name
			service_graph_name = mso_schema_site_service_graph_node.test.service_graph_name
			service_node_type  = mso_schema_site_service_graph_node.test.service_node_type
			service_node_name  = "tfnode2"
		}	
		`
	case "service_graph_name":
		resource += `
		data "mso_schema_site_service_graph_node" "test" {
			schema_id          = mso_schema_site_service_graph_node.test.schema_id
			template_name      = mso_schema_site_service_graph_node.test.template_name
		#	service_graph_name = mso_schema_site_service_graph_node.test.service_graph_name
			service_node_type  = mso_schema_site_service_graph_node.test.service_node_type
			service_node_name  = "tfnode2"
		}	
		`
	case "service_node_type":
		resource += `
		data "mso_schema_site_service_graph_node" "test" {
			schema_id          = mso_schema_site_service_graph_node.test.schema_id
			template_name      = mso_schema_site_service_graph_node.test.template_name
			service_graph_name = mso_schema_site_service_graph_node.test.service_graph_name
		#	service_node_type  = mso_schema_site_service_graph_node.test.service_node_type
			service_node_name  = "tfnode2"
		}
		`
	case "service_node_name":
		resource += `
		data "mso_schema_site_service_graph_node" "test" {
			schema_id          = mso_schema_site_service_graph_node.test.schema_id
			template_name      = mso_schema_site_service_graph_node.test.template_name
			service_graph_name = mso_schema_site_service_graph_node.test.service_graph_name
			service_node_type  = mso_schema_site_service_graph_node.test.service_node_type
		#	service_node_name  = "tfnode2"
		}
		`
	}
	return fmt.Sprintln(resource)
}
