package client

import (
	"fmt"

	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
)

// GetTenantIDFromSchemaTemplate retrieves the Tenant ID from the schema template object.
func (client *Client) GetTenantIDFromSchemaTemplate(schemaID, templateName string) (string, error) {
	schemaObj, err := client.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return "", err
	}

	templatesCount, _ := schemaObj.ArrayCount("templates")
	if err != nil {
		return "", err
	}

	for i := 0; i < templatesCount; i++ {
		templateObj, err := schemaObj.ArrayElement(i, "templates")
		if err != nil {
			return "", err
		}

		apiTemplate := models.StripQuotes(templateObj.S("name").String())
		if templateName == apiTemplate {
			return models.StripQuotes(templateObj.S("tenantId").String()), nil
		}
	}
	return "", nil
}

// GetPoliciesByTenantID returns the policies container object based on the tenant id.
func (client *Client) GetPoliciesByTenantID(objectType, tenantID string) (*container.Container, error) {
	path := fmt.Sprintf("api/v1/templates/objects?type=%s&tenant-id=%s&include-common=true", objectType, tenantID)
	cont, err := client.GetViaURL(path)
	if err != nil {
		return nil, err
	}
	return cont, nil
}

// GetPolicyByTenantID retrieves a policy based on the given object type, object name, and tenant ID.
func (client *Client) GetPolicyByTenantID(objectType, objectName, tenantID string) (map[string]interface{}, error) {
	cont, _ := client.GetPoliciesByTenantID(objectType, tenantID)
	commonTenantPolicy := make(map[string]interface{})
	for _, policy := range cont.Data().([]interface{}) {
		if policyMap, ok := policy.(map[string]interface{}); ok {
			if objectName == policyMap["name"].(string) && tenantID == policyMap["tenantId"].(string) {
				return policyMap, nil
			} else if objectName == policyMap["name"].(string) && policyMap["tenantName"].(string) == "common" {
				commonTenantPolicy = policyMap
			}
		}
	}
	if len(commonTenantPolicy) != 0 {
		return commonTenantPolicy, nil
	}
	return nil, fmt.Errorf("%s policy with name: %s not found", objectType, objectName)
}

// GetObjectNameByUUID returns the name of an object given its UUID and boolean indicating whether the object was found or not.
func GetObjectNameByUUID(objectRef string, objectCont *container.Container) (string, bool) {
	for _, object := range objectCont.Data().([]interface{}) {
		if objectMap, ok := object.(map[string]interface{}); ok {
			if objectMap["uuid"].(string) == objectRef {
				return objectMap["name"].(string), true
			}
		}
	}
	return "", false
}

// GetObjectUUIDByName returns the UUID of an object given its name and boolean indicating whether the object was found or not.
func GetObjectUUIDByName(objectName string, objectCont *container.Container) (string, bool) {
	for _, object := range objectCont.Data().([]interface{}) {
		if objectMap, ok := object.(map[string]interface{}); ok {
			if objectMap["name"].(string) == objectName {
				return objectMap["uuid"].(string), true
			}
		}
	}
	return "", false
}

// GetDHCPPoliciesNameByUUID retrieves the DHCP policies' names by UUID.
// It takes in the tenant ID and a list of object references as parameters.
// The function returns a list of interface{} and an error.
func (client *Client) GetDHCPPoliciesNameByUUID(tenantID string, objectRefs []interface{}) ([]interface{}, error) {
	dhcpPoliciesList := make([]interface{}, 0)
	dhcpRelayCont, relayError := client.GetPoliciesByTenantID("dhcpRelay", tenantID)
	if relayError != nil {
		return nil, relayError
	}

	dhcpOptionCont, optionError := client.GetPoliciesByTenantID("dhcpOption", tenantID)
	if optionError != nil {
		return nil, optionError
	}

	for _, objectRef := range objectRefs {
		var relayObjectFound, optionObjectFound bool
		relayRef := objectRef.(map[string]interface{})["relayRef"].(string)
		optionRef := objectRef.(map[string]interface{})["optionRef"].(string)
		dhcpPolicyMap := make(map[string]interface{})
		dhcpPolicyMap["name"], relayObjectFound = GetObjectNameByUUID(relayRef, dhcpRelayCont)
		if !relayObjectFound {
			return nil, fmt.Errorf("DHCP Relay: %s policy reference not found", relayRef)
		}
		if optionRef != "{}" {
			dhcpPolicyMap["dhcp_option_policy_name"], optionObjectFound = GetObjectNameByUUID(optionRef, dhcpOptionCont)
			if !optionObjectFound {
				return nil, fmt.Errorf("DHCP Option: %s policy reference not found", optionRef)
			}
		} else {
			dhcpPolicyMap["dhcp_option_policy_name"] = ""
		}
		dhcpPoliciesList = append(dhcpPoliciesList, dhcpPolicyMap)
	}
	return dhcpPoliciesList, nil
}

// GetDHCPPoliciesUUIDByName retrieves the DHCP policies UUIDs by name for a given tenant ID.
//
// Parameters:
// - tenantID: The ID of the tenant.
// - objectNames: An array of objects containing the relay name and option name.
func (client *Client) GetDHCPPoliciesUUIDByName(tenantID string, objectNames []interface{}) ([]interface{}, error) {
	dhcpRelayCont, relayError := client.GetPoliciesByTenantID("dhcpRelay", tenantID)
	if relayError != nil {
		return nil, relayError
	}
	dhcpOptionCont, optionError := client.GetPoliciesByTenantID("dhcpOption", tenantID)
	if optionError != nil {
		return nil, optionError
	}
	dhcpPoliciesList := make([]interface{}, 0)
	for _, objectName := range objectNames {
		var relayObjectFound, optionObjectFound bool
		var relayUUID, optionUUID string

		relayName := objectName.(map[string]interface{})["relayName"].(string)
		optionName := objectName.(map[string]interface{})["optionName"].(string)

		relayUUID, relayObjectFound = GetObjectUUIDByName(relayName, dhcpRelayCont)
		if !relayObjectFound {
			return nil, fmt.Errorf("DHCP Relay: %s policy not name found", relayName)
		}

		if optionName != "" {
			optionUUID, optionObjectFound = GetObjectUUIDByName(optionName, dhcpOptionCont)
			if !optionObjectFound {
				return nil, fmt.Errorf("DHCP Option: %s policy not name found", optionName)
			}
		} else {
			optionObjectFound = true
		}

		dhcpPoliciesList = append(
			dhcpPoliciesList, map[string]interface{}{
				"ref": relayUUID,
				"dhcpOptionLabel": map[string]interface{}{
					"ref": optionUUID,
				},
			},
		)
	}
	return dhcpPoliciesList, nil
}
