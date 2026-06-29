package mso

// NOTE: This file uses a bespoke net/http APIC client rather than
// github.com/ciscoecosystem/aci-go-client because aci-go-client/v2 requires
// terraform-plugin-sdk/v2, which conflicts with this provider's current use of
// SDKv1. Once the provider is migrated to SDKv2, this client
// should be replaced with aci-go-client/v2 using client.NewClient +
// client.PostObjectConfig, which removes the boilerplate below.

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"testing"
	"time"
)

const (
	envAPICURL      = "APIC_URL"
	envAPICUsername = "APIC_USERNAME"
	envAPICPassword = "APIC_PASSWORD"
)

// testAPICPreCheck skips t if the required APIC env vars are absent, otherwise
// runs apicSetup to ensure all APIC objects exist. Called in each test's PreCheck
// so objects are re-created after each test's TF destroy wipes the tenant.
func testAPICPreCheck(t *testing.T) {
	t.Helper()
	for _, env := range []string{envAPICURL, envAPICUsername, envAPICPassword} {
		if os.Getenv(env) == "" {
			t.Skipf("Skipping: %s must be set for tests that require APIC configuration", env)
		}
	}
	if err := apicSetup(); err != nil {
		t.Fatalf("APIC setup failed: %v", err)
	}
}

// apicTestClient is a minimal HTTP client for making APIC REST calls during
// acceptance tests. It authenticates with username/password and stores the
// session token in a cookie jar for subsequent requests.
type apicTestClient struct {
	httpClient *http.Client
	baseURL    string
}

// newAPICTestClient authenticates to APIC and returns a ready-to-use client.
func newAPICTestClient() (*apicTestClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("creating cookie jar: %w", err)
	}

	c := &apicTestClient{
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // #nosec G402 — lab env
			},
			Jar: jar,
		},
		baseURL: os.Getenv(envAPICURL),
	}

	if err := c.post("/api/aaaLogin.json", map[string]interface{}{
		"aaaUser": map[string]interface{}{
			"attributes": map[string]string{
				"name": os.Getenv(envAPICUsername),
				"pwd":  os.Getenv(envAPICPassword),
			},
		},
	}); err != nil {
		return nil, fmt.Errorf("APIC login failed: %w", err)
	}
	return c, nil
}

// post sends a POST request to the given APIC path with a JSON body.
func (c *apicTestClient) post(path string, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshalling request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+path, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request to %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("APIC returned HTTP %d for %s: %s", resp.StatusCode, path, string(body))
	}
	return nil
}

// get sends a GET request to the given APIC path and returns the response body.
func (c *apicTestClient) get(path string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request to %s: %w", path, err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("APIC returned HTTP %d for %s: %s", resp.StatusCode, path, string(body))
	}
	return body, nil
}

// waitForAPICMOs polls APIC every 2 seconds until every DN in dns is present
// (totalCount > 0) on APIC, or the timeout elapses. NDO's deploy task reports
// "Complete" when it dispatches a template to APIC, NOT when APIC has finished
// applying it, so acceptance steps that reference deployed-by-name objects
// (e.g. a physical_domain or PC/VPC policy group referenced from a service
// device cluster site) can race APIC and fail with "X does not exist on the
// fabric Y". Use this in a post-deploy Check to gate the next apply step on
// APIC convergence. A single APIC login is reused for all polls.
//
// Each DN is polled to completion in order; the overall timeout applies to
// the whole batch. The loop sleeps 2s between attempts on the currently
// outstanding DN.
func waitForAPICMOs(timeout time.Duration, dns ...string) error {
	if len(dns) == 0 {
		return fmt.Errorf("waitForAPICMOs: no DNs supplied")
	}
	c, err := newAPICTestClient()
	if err != nil {
		return fmt.Errorf("APIC client: %w", err)
	}
	deadline := time.Now().Add(timeout)
	for _, dn := range dns {
		path := fmt.Sprintf("/api/node/mo/%s.json?query-target=self", dn)
		var lastErr error
		for {
			body, err := c.get(path)
			if err == nil {
				var parsed struct {
					TotalCount string `json:"totalCount"`
				}
				if jerr := json.Unmarshal(body, &parsed); jerr == nil && parsed.TotalCount != "" && parsed.TotalCount != "0" {
					log.Printf("[APIC] MO %s present (totalCount=%s)", dn, parsed.TotalCount)
					break
				}
				lastErr = fmt.Errorf("MO %s not yet present (totalCount=%q)", dn, parsed.TotalCount)
			} else {
				lastErr = err
			}
			if !time.Now().Before(deadline) {
				return fmt.Errorf("timeout waiting for APIC MO %s: %w", dn, lastErr)
			}
			time.Sleep(2 * time.Second)
		}
	}
	return nil
}

// apicSetup creates all APIC prerequisites for service graph acceptance tests:
//   - msoTenantName with three L4-L7 firewall devices:
//     two for mso_schema_site_service_graph (no cluster interfaces needed),
//     one for mso_schema_site_contract_service_graph (with named cluster interfaces).
//   - A redirect policy in msoTenantName (provider).
//   - msoTenantName2 with a redirect policy (consumer), exercising the
//     cross-tenant redirect policy path.
func apicSetup() error {
	c, err := newAPICTestClient()
	if err != nil {
		return err
	}

	if err := apicCreateTenant(c, msoTenantName); err != nil {
		return fmt.Errorf("creating APIC tenant %q: %w", msoTenantName, err)
	}
	log.Printf("[APIC] Created tenant %s", msoTenantName)

	for _, deviceName := range []string{
		msoSchemaSiteServiceGraphDeviceName,
		msoSchemaSiteServiceGraphDeviceName2,
	} {
		if err := apicCreateL4L7FirewallDevice(c, msoTenantName, deviceName, "cn_"+deviceName, nil); err != nil {
			return fmt.Errorf("creating site_service_graph device %q: %w", deviceName, err)
		}
		log.Printf("[APIC] Created L4-L7 device uni/tn-%s/lDevVip-%s", msoTenantName, deviceName)
	}

	clusterIfs := []string{
		msoSchemaSiteContractServiceGraphProviderClusterInterface,
		msoSchemaSiteContractServiceGraphConsumerClusterInterface,
	}
	if err := apicCreateL4L7FirewallDevice(c, msoTenantName, msoSchemaSiteContractServiceGraphDeviceName, "cn_"+msoSchemaSiteContractServiceGraphDeviceName, clusterIfs); err != nil {
		return fmt.Errorf("creating contract_service_graph device: %w", err)
	}
	log.Printf("[APIC] Created L4-L7 device uni/tn-%s/lDevVip-%s with interfaces %v",
		msoTenantName, msoSchemaSiteContractServiceGraphDeviceName, clusterIfs)

	if err := apicCreateRedirectPolicy(c, msoTenantName, msoSchemaSiteContractServiceGraphProviderRedirectPolicy); err != nil {
		return fmt.Errorf("creating provider redirect policy: %w", err)
	}
	log.Printf("[APIC] Created redirect policy %s/%s", msoTenantName, msoSchemaSiteContractServiceGraphProviderRedirectPolicy)

	if err := apicCreateTenant(c, msoTenantName2); err != nil {
		return fmt.Errorf("creating APIC tenant %q: %w", msoTenantName2, err)
	}
	log.Printf("[APIC] Created tenant %s", msoTenantName2)

	if err := apicCreateRedirectPolicy(c, msoTenantName2, msoSchemaSiteContractServiceGraphConsumerRedirectPolicy); err != nil {
		return fmt.Errorf("creating consumer redirect policy: %w", err)
	}
	log.Printf("[APIC] Created redirect policy %s/%s", msoTenantName2, msoSchemaSiteContractServiceGraphConsumerRedirectPolicy)

	return nil
}

// apicCreateL4L7FirewallDevice creates a vnsLDevVip of type FW with one
// concrete device and optional logical interfaces (cluster interfaces).
// When clusterIfNames is nil or empty, no logical interfaces are created —
// sufficient for mso_schema_site_service_graph tests that only reference the DN.
// When provided, each name gets a vnsLIf child — required for
// mso_schema_site_contract_service_graph tests that reference cluster interface names.
//
// The concrete interface path topology/pod-1/paths-101/pathep-[eth1/1] is a
// typical lab topology path for a leaf switch interface.
func apicCreateL4L7FirewallDevice(c *apicTestClient, tenant, deviceName, concreteIfName string, clusterIfNames []string) error {
	cIfTDn := fmt.Sprintf("uni/tn-%s/lDevVip-%s/cDev-%s/cIf-[%s]", tenant, deviceName, deviceName, concreteIfName)

	lIfChildren := make([]interface{}, 0, len(clusterIfNames))
	for _, ifName := range clusterIfNames {
		lIfChildren = append(lIfChildren, map[string]interface{}{
			"vnsLIf": map[string]interface{}{
				"attributes": map[string]string{"name": ifName},
				"children": []interface{}{
					map[string]interface{}{
						"vnsRsCIfAttN": map[string]interface{}{
							"attributes": map[string]string{"tDn": cIfTDn},
						},
					},
				},
			},
		})
	}

	cDevChild := map[string]interface{}{
		"vnsCDev": map[string]interface{}{
			"attributes": map[string]string{"name": deviceName},
			"children": []interface{}{
				map[string]interface{}{
					"vnsCIf": map[string]interface{}{
						"attributes": map[string]string{"name": concreteIfName},
						"children": []interface{}{
							map[string]interface{}{
								"vnsRsCIfPathAtt": map[string]interface{}{
									"attributes": map[string]string{
										"tDn": "topology/pod-1/paths-101/pathep-[eth1/1]",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return c.post(fmt.Sprintf("/api/node/mo/uni/tn-%s.json", tenant), map[string]interface{}{
		"vnsLDevVip": map[string]interface{}{
			"attributes": map[string]string{
				"svcType": "FW",
				"managed": "false",
				"name":    deviceName,
			},
			"children": append([]interface{}{cDevChild}, lIfChildren...),
		},
	})
}

// apicCreateRedirectPolicy creates a vnsSvcRedirectPol in the given tenant.
func apicCreateRedirectPolicy(c *apicTestClient, tenant, policyName string) error {
	return c.post(fmt.Sprintf("/api/node/mo/uni/tn-%s/svcCont.json", tenant), map[string]interface{}{
		"vnsSvcRedirectPol": map[string]interface{}{
			"attributes": map[string]string{"name": policyName},
		},
	})
}

// apicCreateTenant creates an fvTenant in APIC. The call is idempotent — if the
// tenant already exists (e.g. because the other setup function ran first) APIC
// returns 200 and the call succeeds.
func apicCreateTenant(c *apicTestClient, tenantName string) error {
	return c.post(fmt.Sprintf("/api/node/mo/uni/tn-%s.json", tenantName), map[string]interface{}{
		"fvTenant": map[string]interface{}{
			"attributes": map[string]string{"name": tenantName},
		},
	})
}
