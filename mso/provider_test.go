package mso

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
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
		var msoUrl, msoUsername, msoPassword, msoPlatform string
		var msoRetries int
		if username := os.Getenv("MSO_USERNAME"); username == "" {
			t.Fatal("MSO_USERNAME must be set for acceptance tests")
		} else {
			msoUsername = username
		}
		if password := os.Getenv("MSO_PASSWORD"); password == "" {
			t.Fatal("MSO_PASSWORD must be set for acceptance tests")
		} else {
			msoPassword = password
		}
		if url := os.Getenv("MSO_URL"); url == "" {
			t.Fatal("MSO_URL must be set for acceptance tests")
		} else {
			msoUrl = url
		}
		if platform := os.Getenv("MSO_PLATFORM"); platform == "" {
			msoPlatform = "mso"
		} else {
			msoPlatform = platform
		}
		if retries := os.Getenv("MSO_RETRIES"); retries == "" {
			msoRetries = 2
		} else {
			retriesInt, err := strconv.Atoi(retries)
			if err != nil {
				t.Log("Warning: MSO_RETRIES is not a valid integer, using default value 2. Error:", err.Error())
				msoRetries = 2
			} else {
				msoRetries = retriesInt
			}
		}

		msoClientTest = client.GetClient(msoUrl, msoUsername, client.Password(msoPassword), client.Insecure(true), client.Platform(msoPlatform), client.MaxRetries(msoRetries))
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

func CustomTestCheckTypeSetElemAttrs(resourceName, setName string, attrsToCheck map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}
		groupedAttrs := make(map[string]map[string]string)
		re := regexp.MustCompile(fmt.Sprintf(`^%s\.(\d+)\.(.*)$`, setName))

		for key, val := range rs.Primary.Attributes {
			matches := re.FindStringSubmatch(key)
			if len(matches) == 3 {
				hash := matches[1]
				attrName := matches[2]
				if _, ok := groupedAttrs[hash]; !ok {
					groupedAttrs[hash] = make(map[string]string)
				}
				groupedAttrs[hash][attrName] = val
			}
		}

		for _, elemAttrs := range groupedAttrs {
			match := true
			for expectedKey, expectedVal := range attrsToCheck {
				if val, ok := elemAttrs[expectedKey]; ok {
					if fmt.Sprintf("%v", val) != expectedVal {
						match = false
						break
					}
				} else {
					match = false
					break
				}
			}

			if match {
				return nil
			}
		}
		return fmt.Errorf("No element in set '%s' found with the following attributes: %v", setName, attrsToCheck)
	}
}

func setupTestLogCapture(t *testing.T, logLevel string) string {
	logFile, err := os.CreateTemp("", "tf-acc-test-*.log")
	if err != nil {
		t.Fatalf("Failed to create temp log file: %v", err)
	}

	logFileName := logFile.Name()

	t.Cleanup(func() {
		logFile.Close()
		os.Remove(logFileName)
	})

	t.Setenv("TF_LOG", logLevel)
	t.Setenv("TF_LOG_PATH", logFileName)

	return logFileName
}

func customTestCheckLogs(logFilePath string, patterns []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		file, err := os.Open(logFilePath)
		if err != nil {
			return fmt.Errorf("failed to open log file %s: %w", logFilePath, err)
		}
		defer file.Close()

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
