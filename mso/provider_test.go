package mso

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"mso": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

var (
	msoClientTest     *client.Client
	msoClientTestOnce sync.Once
)

func testAccPreCheck(t *testing.T) *client.Client {
	msoClientTestOnce.Do(func() {
		var mso_url, mso_username, mso_password, mso_platform string
		if v := os.Getenv("MSO_USERNAME"); v == "" {
			t.Fatal("MSO_USERNAME must be set for acceptance tests")
		} else {
			mso_username = v
		}
		if v := os.Getenv("MSO_PASSWORD"); v == "" {
			t.Fatal("MSO_PASSWORD must be set for acceptance tests")
		} else {
			mso_password = v
		}
		if v := os.Getenv("MSO_URL"); v == "" {
			t.Fatal("MSO_URL must be set for acceptance tests")
		} else {
			mso_url = v
		}
		if v := os.Getenv("MSO_PLATFORM"); v == "" {
			mso_platform = "mso"
		} else {
			mso_platform = v
		}

		msoClientTest = client.GetClient(mso_url, mso_username, client.Password(mso_password), client.Insecure(true), client.Platform(mso_platform))
	})
	return msoClientTest

}

func testCheckResourceDestroyPolicyWithArguments(resource, policyType string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		return testCheckResourceDestroyPolicy(s, resource, policyType)
	}
}

func testCheckResourceDestroyPolicy(s *terraform.State, resource, policyType string) error {
	msoClient := testAccPreCheck(nil)
	for name, rs := range s.RootModule().Resources {
		if rs.Type == resource {
			response, err := msoClient.GetViaURL((fmt.Sprintf("api/v1/templates/objects?type=%s&uuid=%s", policyType, rs.Primary.Attributes["uuid"])))
			if err != nil {
				if response.S("code").Data().(float64) == 404 {
					continue
				} else {
					return fmt.Errorf("error checking if resource '%s' with ID '%s' still exists: %s", name, rs.Primary.ID, err)
				}
			}
			return fmt.Errorf("terraform destroy was unsuccessful. The resource '%s' with ID '%s' still exists", name, rs.Primary.ID)
		}
	}
	return nil
}

func testCheckResourceDestroyPolicyWithPathAttributesAndArguments(resource string, objectPath ...string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		return testCheckResourceDestroyPolicyWithPathAttributes(s, resource, objectPath...)
	}
}

func testCheckResourceDestroyPolicyWithPathAttributes(s *terraform.State, resource string, objectPath ...string) error {
	msoClient := testAccPreCheck(nil)
	for name, rs := range s.RootModule().Resources {
		if rs.Type == resource {
			response, err := msoClient.GetViaURL((fmt.Sprintf("api/v1/templates/%s", rs.Primary.Attributes["template_id"])))
			if err != nil {
				continue
			}
			policyObjects := response.S(objectPath...)
			if policyObjects.Data() != nil {
				policyCount, err := response.ArrayCount(objectPath...)
				if err == nil {
					for i := range policyCount {
						policy := policyObjects.Index(i)
						uuid, ok := policy.S("uuid").Data().(string)
						if ok && uuid == rs.Primary.Attributes["uuid"] {
							return fmt.Errorf("terraform destroy was unsuccessful. The resource '%s' with ID '%s' still exists", name, rs.Primary.ID)
						}
					}
				} else {
					uuid, ok := policyObjects.S("uuid").Data().(string)
					if ok && uuid == rs.Primary.Attributes["uuid"] {
						return fmt.Errorf("terraform destroy was unsuccessful. The resource '%s' with ID '%s' still exists", name, rs.Primary.ID)
					}
				}
			}
		}
	}
	return nil
}

func IsReference(s string) bool {
	return strings.HasPrefix(s, "mso_") || strings.HasPrefix(s, "data.mso_")
}

func testAccVerifyKeyValue(resourceAttrsMap *map[string]string, resourceAttrRootkey, stateKey, stateValue string) {
	stateKeySplit := strings.Split(stateKey, ".")
	for inputKey, inputValue := range *resourceAttrsMap {
		if strings.Contains(stateKey, resourceAttrRootkey) && stateKeySplit[len(stateKeySplit)-1] == inputKey && (stateValue == inputValue || (IsReference(inputValue) && stateValue != "")) {
			delete(*resourceAttrsMap, inputKey)
			break
		}
	}
}

func customTestCheckResourceTypeSetAttr(resourceName, resourceAttrRootkey string, resourceAttrsMap map[string]string) resource.TestCheckFunc {
	return func(is *terraform.State) error {
		rootModule, err := is.RootModule().Resources[resourceName]
		if !err {
			return fmt.Errorf("%v", err)
		}
		if rootModule.Primary.ID == "" {
			return fmt.Errorf("No ID is set for the template")
		}
		for stateKey, stateValue := range rootModule.Primary.Attributes {
			testAccVerifyKeyValue(&resourceAttrsMap, resourceAttrRootkey, stateKey, stateValue)
		}
		if len(resourceAttrsMap) > 0 {
			return fmt.Errorf("Assertion check failed,\nCurrent state file content: %v\nComparable to unmatched values: %v", rootModule.Primary.Attributes, resourceAttrsMap)
		}
		return nil
	}
}

func customTestCheckLogs(logFilePath string, patterns []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		file, err := os.Open(logFilePath)
		if err != nil {
			return fmt.Errorf("failed to open log file %s: %w", logFilePath, err)
		}
		defer file.Close()

		// To check a sequence, we can read the whole file line-by-line into a builder
		var logBuilder strings.Builder
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			logBuilder.WriteString(scanner.Text() + "\n")
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error while scanning log file: %w", err)
		}

		logOutput := logBuilder.String()

		fullPattern := "(?s)" + strings.Join(patterns, ".*")

		matched, err := regexp.MatchString(fullPattern, logOutput)
		if err != nil {
			return fmt.Errorf("error compiling regex pattern: %w", err)
		}

		if !matched {
			expectedSequence := strings.Join(patterns, "\n...\n")
			return fmt.Errorf(
				"expected log sequence not found.\n--- Expected Sequence (regex) ---\n%s\n\n--- Full Log Output ---\n%s",
				expectedSequence,
				logOutput,
			)
		}

		if err := os.Truncate(logFilePath, 0); err != nil {
			return fmt.Errorf("failed to truncate log file: %w", err)
		}

		return nil
	}
}
