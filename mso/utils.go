package mso

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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

// extractServiceGraphNodesFromContainer extracts the nodes from the given container.
//
// Parameters:
// - cont: A pointer to the container.Container object.
//
// Returns:
// - nodes: A slice of interfaces that contains the extracted nodes.
func extractServiceGraphNodesFromContainer(cont *container.Container) []interface{} {
	nodes := make([]interface{}, 0, 1)
	for _, node := range cont.S("serviceNodes").Data().([]interface{}) {
		nodes = append(nodes, models.StripQuotes(node.(map[string]interface{})["name"].(string)))
	}
	return nodes
}

// getSchemaTemplateServiceGraphFromContainer retrieves the schema template service graph based on the provided parameters.
//
// Parameters:
// - cont: The container object.
// - templateName: The name of the template.
// - graphName: The name of the service graph.
//
// Returns:
// - cont: The template service graph container object.
// - int: The index of the service graph in the container.
// - error: An error indicating any issues encountered during the retrieval process.
func getSchemaTemplateServiceGraphFromContainer(cont *container.Container, templateName, graphName string) (*container.Container, int, error) {
	templateCount, err := cont.ArrayCount("templates")
	if err != nil {
		return nil, -1, fmt.Errorf("No Template found")
	}

	for i := 0; i < templateCount; i++ {
		templateCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return nil, -1, fmt.Errorf("Unable to get template element")
		}

		apiTemplate := models.StripQuotes(templateCont.S("name").String())

		if apiTemplate == templateName {
			log.Printf("[DEBUG] Template found")

			sgCount, err := templateCont.ArrayCount("serviceGraphs")

			if err != nil {
				return nil, -1, fmt.Errorf("No Service Graph found")
			}

			for j := 0; j < sgCount; j++ {
				sgCont, err := templateCont.ArrayElement(j, "serviceGraphs")

				if err != nil {
					return nil, -1, fmt.Errorf("Unable to get service graph element")
				}

				apiSgName := models.StripQuotes(sgCont.S("name").String())

				if apiSgName == graphName {
					return sgCont, j, nil
				}
			}

		}
	}
	return nil, -1, fmt.Errorf("unable to find service graph")
}

// Verifies, if the value (string) is in the list of strings
func valueInSliceofStrings(value string, list []string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

// convertInterfaceToString converts an interface object to a string.
//
// interfaceObject: The interface object to be converted.
// Returns: The converted string.
func convertInterfaceToString(interfaceObject interface{}) string {
	switch v := interfaceObject.(type) {
	case string:
		return v
	case nil:
		return ""
	default:
		return ""
	}
}

// convertInterfaceToInt converts an interface object to an integer.
//
// interfaceObject: The interface object to be converted.
// Returns: The converted integer value.
func convertInterfaceToInt(interfaceObject interface{}) int {
	switch v := interfaceObject.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case nil:
		return 0
	default:
		return 0
	}
}

func boolToYesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func yesNoToBool(s string) bool {
	return s == "yes"
}

func getMapKeys(inputMap map[string]string) []string {
	var keys []string
	for k := range inputMap {
		keys = append(keys, k)
	}
	return keys
}

func getKeyByValue(inputMap map[string]string, value string) string {
	for k, v := range inputMap {
		if v == value {
			return k
		}
	}
	return ""
}

// Removes the schema id from the id and returns the path needed in PATCH request
func getPathFromId(id string) string {
	return fmt.Sprintf("/%s", strings.Join(strings.Split(id, "/")[1:], "/"))
}

func addPatchPayloadToContainer(payloadContainer *container.Container, op, path string, value interface{}) error {

	payloadMap := map[string]interface{}{"op": op, "path": path, "value": value}

	payload, err := json.Marshal(payloadMap)
	if err != nil {
		return err
	}

	jsonContainer, err := container.ParseJSON([]byte(payload))
	if err != nil {
		return err
	}

	err = payloadContainer.ArrayAppend(jsonContainer.Data())
	if err != nil {
		return err
	}

	return nil
}

func doPatchRequest(msoClient *client.Client, path string, payloadCon *container.Container) error {

	req, err := msoClient.MakeRestRequest("PATCH", path, payloadCon, true)
	if err != nil {
		return err
	}

	cont, _, err := msoClient.Do(req)
	if err != nil {
		return err
	}

	err = client.CheckForErrors(cont, "PATCH")
	if err != nil {
		return err
	}

	return nil
}

func getSchemaIdFromName(msoClient *client.Client, name string) (string, error) {

	con, err := msoClient.GetViaURL("/api/v1/schemas/list-identity")

	if err != nil {
		return "", nil
	}

	schemas := con.S("schemas").Data().([]interface{})
	for _, schema := range schemas {
		if displayName, ok := schema.(map[string]interface{})["displayName"]; ok && displayName == name {
			if id, ok := schema.(map[string]interface{})["id"]; ok {
				return id.(string), nil
			}
		}
	}

	return "", fmt.Errorf("Schema of specified name not found")
}

func getStaticPortPathValues(pathValue string, re *regexp.Regexp) map[string]string {
	match := re.FindStringSubmatch(pathValue) //list of matched strings
	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	return result
}

func setValuesFromPortPath(staticPortMap map[string]interface{}, pathValue string) {

	portFexPath := regexp.MustCompile(`(topology\/(?P<podValue>.*)\/paths-(?P<leafValue>.*)\/extpaths-(?P<fexValue>.*)\/pathep-\[(?P<pathValue>.*)\])`)
	vpcFexPath := regexp.MustCompile(`(topology\/(?P<podValue>.*)\/protpaths-(?P<leafValue>.*)\/extprotpaths-(?P<fexValue>.*)\/pathep-\[(?P<pathValue>.*)\])`)
	vpcPath := regexp.MustCompile(`(topology\/(?P<podValue>.*)\/protpaths-(?P<leafValue>.*)\/pathep-\[(?P<pathValue>.*)\])`)
	// dpcPath also handles the port without FEX defined in the path
	dpcPath := regexp.MustCompile(`(topology\/(?P<podValue>.*)\/paths-(?P<leafValue>.*)\/pathep-\[(?P<pathValue>.*)\])`)

	matchedMap := make(map[string]string)

	if portFexPath.MatchString(pathValue) {
		matchedMap = getStaticPortPathValues(pathValue, portFexPath)
	} else if vpcFexPath.MatchString(pathValue) {
		matchedMap = getStaticPortPathValues(pathValue, vpcFexPath)
	} else if vpcPath.MatchString(pathValue) {
		matchedMap = getStaticPortPathValues(pathValue, vpcPath)
	} else if dpcPath.MatchString(pathValue) {
		matchedMap = getStaticPortPathValues(pathValue, dpcPath)
	}

	staticPortMap["pod"] = matchedMap["podValue"]
	staticPortMap["leaf"] = matchedMap["leafValue"]
	staticPortMap["path"] = matchedMap["pathValue"]
	if fexValue, ok := matchedMap["fexValue"]; ok {
		staticPortMap["fex"] = fexValue
	}

}

func createPortPath(path_type, static_port_pod, static_port_leaf, static_port_fex, static_port_path string) string {

	if path_type == "vpc" && static_port_fex != "" {
		return fmt.Sprintf("topology/%s/protpaths-%s/extprotpaths-%s/pathep-[%s]", static_port_pod, static_port_leaf, static_port_fex, static_port_path)
	} else if path_type == "vpc" {
		return fmt.Sprintf("topology/%s/protpaths-%s/pathep-[%s]", static_port_pod, static_port_leaf, static_port_path)
	} else if static_port_fex != "" {
		return fmt.Sprintf("topology/%s/paths-%s/extpaths-%s/pathep-[%s]", static_port_pod, static_port_leaf, static_port_fex, static_port_path)
	} else {
		return fmt.Sprintf("topology/%s/paths-%s/pathep-[%s]", static_port_pod, static_port_leaf, static_port_path)
	}
}

func getListOfStringsFromSchemaList(d *schema.ResourceData, key string) []string {
	if values, ok := d.GetOk(key); ok {
		return convertToListOfStrings(values.([]interface{}))
	}
	return nil
}

func convertToListOfStrings(values []interface{}) []string {
	result := []string{}
	for _, item := range values {
		result = append(result, item.(string))
	}
	return result
}

func duplicatesInList(list []string) []string {
	duplicates := []string{}
	set := make(map[string]int)
	for index, item := range list {
		if _, ok := set[item]; ok {
			duplicates = append(duplicates, item)
		} else {
			set[item] = index
		}
	}
	return duplicates
}

func GetTemplateIdFromResourceId(input string) (string, error) {
	parts := strings.Split(input, "/")
	if parts[0] != "templateId" {
		return "", fmt.Errorf("Invalid resource id provided")
	}
	return parts[1], nil
}

func GetPolicyNameFromResourceId(input, policyType string) (string, error) {
	parts := strings.Split(input, "/")

	for i := 0; i < len(parts)-1; i++ {
		if parts[i] == policyType {
			if i+1 < len(parts) {
				return parts[i+1], nil
			}
			return "", fmt.Errorf("No value found after policyType")
		}
	}

	return "", fmt.Errorf("PolicyType not found in the id")
}

func GetPolicyIndexByKeyAndValue(cont *container.Container, policyIdentifier, policyIdentifierValue string, templateElements ...string) (int, error) {
	index := -1

	policyArray := cont.S(templateElements...)
	if policyArray.Data() == nil {
		return index, fmt.Errorf("Policy type %s is not a list or does not exist", templateElements[len(templateElements)-1])
	}

	policyCount, err := cont.ArrayCount(templateElements...)
	if err != nil {
		return index, err
	}

	for i := 0; i < policyCount; i++ {
		policy := policyArray.Index(i)
		identifierValue := policy.S(policyIdentifier).Data().(string)
		if identifierValue == policyIdentifierValue {
			index = i
			break
		}
	}

	if index == -1 {
		return index, fmt.Errorf("Policy %s %s not found in policy list", policyIdentifier, policyIdentifierValue)
	}

	return index, nil
}

func GetPolicyByName(cont *container.Container, policyName string, templateElements ...string) (*container.Container, error) {
	policyObject := cont.S(templateElements...)
	if policyObject.Data() != nil {
		policyCount, err := cont.ArrayCount(templateElements...)
		if err == nil {
			for i := 0; i < policyCount; i++ {
				policy := policyObject.Index(i)
				name, ok := policy.S("name").Data().(string)
				if ok && name == policyName {
					return policy, nil
				}
			}
		} else {
			name, ok := policyObject.S("name").Data().(string)
			if ok && name == policyName {
				return policyObject, nil
			}
		}
	}

	return nil, fmt.Errorf("Policy name %s not found", policyName)
}

func testAccVerifyKeyValue(resourceAttrsMap *map[string]string, stateKey, stateValue string) {
	for inputKey, inputValue := range *resourceAttrsMap {
		inputPattern := regexp.MustCompile(`\.\d+`).ReplaceAllString(inputKey, `\.\d+`)
		pattern := regexp.MustCompile("^" + inputPattern + "$")
		if pattern.MatchString(stateKey) && (stateValue == inputValue || (inputValue == "reference" && stateValue != "")) {
			delete(*resourceAttrsMap, inputKey)
			break
		}
	}
}

func customTestCheckResourceAttr(resourceName string, resourceAttrsMap map[string]string) resource.TestCheckFunc {
	return func(is *terraform.State) error {
		rootModule, err := is.RootModule().Resources[resourceName]
		if !err {
			return fmt.Errorf("%v", err)
		}
		if rootModule.Primary.ID == "" {
			return fmt.Errorf("No ID is set for the template")
		}
		for stateKey, stateValue := range rootModule.Primary.Attributes {
			testAccVerifyKeyValue(&resourceAttrsMap, stateKey, stateValue)
		}
		if len(resourceAttrsMap) > 0 {
			return fmt.Errorf("Assertion check failed,\nCurrent state file content: %v\nSimilar to unmatched values: %v", rootModule.Primary.Attributes, resourceAttrsMap)
		}
		return nil
	}
}
