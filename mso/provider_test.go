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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"mso": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
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

func testCheckTypeSetStringElemAttr(resourceName, setKey, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		prefix := setKey + "."
		for k, v := range rs.Primary.Attributes {
			if strings.HasPrefix(k, prefix) && !strings.HasSuffix(k, ".#") && v == value {
				return nil
			}
		}
		return fmt.Errorf("no element with value %q found in set %q", value, setKey)
	}
}

// Deprecated: This check has a bug because it matches key value pairs across different set elements instead of verifying all attributes
// belong to the same set element, leading to false positives. Use CustomTestCheckTypeSetElemAttrs instead.
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

// resolveStateReference resolves a value that looks like a Terraform resource reference
// (e.g. "resource_type.resource_name.attribute") by looking up the actual value in state.
// If the value does not match this pattern or the referenced resource/attribute is not found,
// the original value is returned unchanged.
func resolveStateReference(s *terraform.State, value string) string {
	parts := strings.SplitN(value, ".", 3)
	if len(parts) != 3 {
		return value
	}
	resourceKey := parts[0] + "." + parts[1]
	attrName := parts[2]
	if rs, ok := s.RootModule().Resources[resourceKey]; ok {
		if attrVal, ok := rs.Primary.Attributes[attrName]; ok {
			return attrVal
		}
	}
	return value
}

func CustomTestCheckTypeSetElemAttrs(resourceName, setName string, attrsToCheck map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		resolvedAttrs := make(map[string]string, len(attrsToCheck))
		for k, v := range attrsToCheck {
			resolvedAttrs[k] = resolveStateReference(s, v)
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

		// Numeric path segments in expected keys (e.g. "pbr_destination.0.ip")
		// are treated as wildcards so the same expected key can match a nested
		// TypeSet whose real index is a hash. Keys without numeric segments
		// keep the original exact-match lookup.
		keyMatchers := make(map[string]*regexp.Regexp, len(resolvedAttrs))
		for expectedKey := range resolvedAttrs {
			parts := strings.Split(expectedKey, ".")
			patternParts := make([]string, len(parts))
			hasNumericSegment := false
			for i, p := range parts {
				if _, err := strconv.Atoi(p); err == nil {
					patternParts[i] = `\d+`
					hasNumericSegment = true
				} else {
					patternParts[i] = regexp.QuoteMeta(p)
				}
			}
			if hasNumericSegment {
				keyMatchers[expectedKey] = regexp.MustCompile("^" + strings.Join(patternParts, `\.`) + "$")
			}
		}

		for _, elemAttrs := range groupedAttrs {
			match := true
			for expectedKey, expectedVal := range resolvedAttrs {
				if val, ok := elemAttrs[expectedKey]; ok {
					if fmt.Sprintf("%v", val) != expectedVal {
						match = false
						break
					}
					continue
				} else if expectedVal != "" {
					// SDKv2 omits zero-value Optional fields (empty string, false) from
					// TypeSet element flat state. Treat an absent key as matching only
					// when the expected value is also the zero value ("").
					match = false
					break
				}
				if matcher, ok := keyMatchers[expectedKey]; ok {
					found := false
					for k, v := range elemAttrs {
						if matcher.MatchString(k) && fmt.Sprintf("%v", v) == expectedVal {
							found = true
							break
						}
					}
					if !found {
						match = false
						break
					}
					continue
				}
				match = false
				break
			}

			if match {
				return nil
			}
		}
		return fmt.Errorf("No element in set '%s' found with the following attributes: %v\nResolved to: %v\nState attributes for resource: %v", setName, attrsToCheck, resolvedAttrs, rs.Primary.Attributes)
	}
}

// testAccVersionCheck skips the test if the NDO version is older than minVersion.
// Must be called after testAccPreCheck in the same PreCheck function.
// minVersion should be a version string like "5.1" or "4.0.0.0".
func testAccVersionCheck(t *testing.T, minVersion string) {
	t.Helper()
	result, err := msoClientTest.CompareVersion(minVersion)
	if err != nil {
		t.Skipf("Skipping: could not determine NDO version: %s", err)
	}
	if result > 0 {
		t.Skipf("Skipping: requires NDO >= %s", minVersion)
	}
}

// CustomTestCheckCollectionElemAttrsByKeys locates the TypeSet element whose
// attributes match every key/value pair in matchAttrs and compares every key
// in attrsToCheck against the value stored in state, producing a
// per-attribute diff. Unlike CustomTestCheckTypeSetElemAttrs, which only
// reports "no element matched the whole map", this helper pinpoints exactly
// which attribute(s) differ — useful when an upstream server silently coerces
// a subset of fields and the test just needs to see which ones. Both
// matchAttrs values and attrsToCheck values are resolved through
// resolveStateReference so they can refer to other resources in state.
func CustomTestCheckCollectionElemAttrsByKeys(resourceName, setName string, matchAttrs, attrsToCheck map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}
		if len(matchAttrs) == 0 {
			return fmt.Errorf("matchAttrs must contain at least one key/value pair")
		}

		resolvedMatch := make(map[string]string, len(matchAttrs))
		for k, v := range matchAttrs {
			resolvedMatch[k] = resolveStateReference(s, v)
		}
		resolvedAttrs := make(map[string]string, len(attrsToCheck))
		for k, v := range attrsToCheck {
			resolvedAttrs[k] = resolveStateReference(s, v)
		}

		groupedAttrs := make(map[string]map[string]string)
		re := regexp.MustCompile(fmt.Sprintf(`^%s\.(\d+)\.(.*)$`, regexp.QuoteMeta(setName)))
		for key, val := range rs.Primary.Attributes {
			if m := re.FindStringSubmatch(key); len(m) == 3 {
				hash := m[1]
				if _, ok := groupedAttrs[hash]; !ok {
					groupedAttrs[hash] = make(map[string]string)
				}
				groupedAttrs[hash][m[2]] = val
			}
		}

		var matchingHashes []string
		for hash, elemAttrs := range groupedAttrs {
			allMatch := true
			for mk, mv := range resolvedMatch {
				if av, ok := elemAttrs[mk]; !ok || av != mv {
					allMatch = false
					break
				}
			}
			if allMatch {
				matchingHashes = append(matchingHashes, hash)
			}
		}

		if len(matchingHashes) == 0 {
			return fmt.Errorf("%s element matching %v not found in state", setName, resolvedMatch)
		}
		if len(matchingHashes) > 1 {
			return fmt.Errorf("%s element match %v is ambiguous: %d elements matched, refine matchAttrs", setName, resolvedMatch, len(matchingHashes))
		}

		elemAttrs := groupedAttrs[matchingHashes[0]]
		var diffs []string
		for ek, ev := range resolvedAttrs {
			av, present := elemAttrs[ek]
			switch {
			case !present:
				diffs = append(diffs, fmt.Sprintf("%s: expected %q, missing in state", ek, ev))
			case av != ev:
				diffs = append(diffs, fmt.Sprintf("%s: expected %q, got %q", ek, ev, av))
			}
		}
		if len(diffs) > 0 {
			return fmt.Errorf("%s element %v mismatches:\n  %s", setName, resolvedMatch, strings.Join(diffs, "\n  "))
		}
		return nil
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
