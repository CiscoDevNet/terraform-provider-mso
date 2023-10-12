package mso

import (
	"log"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const version = 1

func toStringList(configured interface{}) []string {
	vs := make([]string, 0, 1)
	val, ok := configured.(string)
	if ok && val != "" {
		vs = append(vs, val)
	}
	return vs
}

func errorForObjectNotFound(err error, dn string, con *container.Container, d *schema.ResourceData) error {
	if err != nil {
		if con.S("code").String() == "404" || strings.HasSuffix(err.Error(), "not found") || strings.HasSuffix(models.StripQuotes(con.S("error").String()), "no documents in result") {
			log.Printf("[WARN] %s, removing from state: %s", err, dn)
			d.SetId("")
			return nil
		} else {
			return err
		}
	}
	return nil
}

// CHANGE TO UTILS IF YOU ADD MORE FUNCTIONS

// checkNodeAttr checks the attribute of a node in a list of objects.
//
// Parameters:
// - object: the list of objects.
// - attrName: the name of the attribute to check.
// - index: the index of the object in the list.
//
// Returns true if the attribute is not empty, false otherwise.
func checkNodeAttr(object interface{}, attrName string, index int) bool {
	objList := object.([]interface{})

	instance := objList[index].(map[string]interface{})

	if instance[attrName] != "" {
		return true
	}
	return false
}

// extractNodes extracts the nodes from the given container.
//
// Parameters:
// - cont: A pointer to the container.Container object.
//
// Returns:
// - nodes: A slice of interfaces that contains the extracted nodes.
// - error: An error object if there is any error encountered during the extraction process.
func extractNodes(cont *container.Container) ([]interface{}, error) {
	nodes := make([]interface{}, 0, 1)
	count, err := cont.ArrayCount("serviceNodes")
	if err != nil {
		return nodes, err
	}

	for i := 0; i < count; i++ {
		node, err := cont.ArrayElement(i, "serviceNodes", "name")
		if err != nil {
			return nodes, err
		}

		nodes = append(nodes, models.StripQuotes(node.String()))
	}
	return nodes, nil
}
