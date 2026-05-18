package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
)

// cleanupOrphanSchemaSiteTestResources deletes any leftover schema/tenant on
// MSO that match the package-global test names (msoSchemaName, msoTenantName).
//
// Why this exists: SDK v1's testing framework rolls back state when an apply
// errors -- including an apply that fails on a data source read. Error-flow
// acceptance tests intentionally trigger such failures, so the prereq
// tenant/schema/site association created during the failing apply are
// orphaned on MSO (present on the server, absent from Terraform state).
// Without explicit cleanup, the next step (or test run) fails on
// "Tenant/Schema: '...' already exists".
//
// Cleanup order: schema first (so the schema_site association inside it goes
// with it), then tenant. Lives here in test_cleanup_helpers.go so any
// acceptance test in the package that intentionally triggers apply-time
// errors can reuse it.
func cleanupOrphanSchemaSiteTestResources(t *testing.T) {
	msoClient := testAccPreCheck(t)
	deleteSchemaByDisplayName(t, msoClient, msoSchemaName)
	deleteTenantByName(t, msoClient, msoTenantName)
}

// deleteSchemaByDisplayName best-effort deletes a schema on MSO whose
// `displayName` matches the supplied value. Intended as a PreConfig cleanup
// helper for acceptance tests that intentionally trigger apply-time errors --
// SDK v1 rolls back state on any apply error, leaving the prereq schema
// orphaned on the server. Errors are logged via t.Logf and ignored so this is
// safe to call regardless of prior state.
func deleteSchemaByDisplayName(t *testing.T, msoClient *client.Client, displayName string) {
	con, err := msoClient.GetViaURL("api/v1/schemas")
	if err != nil {
		t.Logf("cleanup: list schemas failed (ignored): %v", err)
		return
	}
	raw, ok := con.S("schemas").Data().([]interface{})
	if !ok {
		return
	}
	for _, info := range raw {
		val, ok := info.(map[string]interface{})
		if !ok {
			continue
		}
		if val["displayName"] == displayName {
			id, _ := val["id"].(string)
			if id == "" {
				continue
			}
			if err := msoClient.DeletebyId("api/v1/schemas/" + id); err != nil {
				t.Logf("cleanup: delete schema %q (id=%s) failed (ignored): %v", displayName, id, err)
				return
			}
			t.Logf("cleanup: deleted orphan schema %q (id=%s)", displayName, id)
			return
		}
	}
}

// deleteTenantByName best-effort deletes a tenant on MSO whose `name` matches
// the supplied value. See deleteSchemaByDisplayName for the rationale; the
// orchestrator-only flag matches what resourceMSOTenantDelete uses so the
// cleanup mirrors a normal Terraform destroy.
func deleteTenantByName(t *testing.T, msoClient *client.Client, name string) {
	con, err := msoClient.GetViaURL("api/v1/tenants")
	if err != nil {
		t.Logf("cleanup: list tenants failed (ignored): %v", err)
		return
	}
	raw, ok := con.S("tenants").Data().([]interface{})
	if !ok {
		return
	}
	for _, info := range raw {
		val, ok := info.(map[string]interface{})
		if !ok {
			continue
		}
		if val["name"] == name {
			id, _ := val["id"].(string)
			if id == "" {
				continue
			}
			if err := msoClient.DeletebyId(fmt.Sprintf("api/v1/tenants/%s?msc-only=true", id)); err != nil {
				t.Logf("cleanup: delete tenant %q (id=%s) failed (ignored): %v", name, id, err)
				return
			}
			t.Logf("cleanup: deleted orphan tenant %q (id=%s)", name, id)
			return
		}
	}
}
