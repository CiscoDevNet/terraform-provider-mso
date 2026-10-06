package ndoapi

import (
	"reflect"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/models"
)

var _ models.Model = (*PatchPayload)(nil)

func TestPatchPayloadToMap(t *testing.T) {
	tests := []struct {
		name     string
		payload  *PatchPayload
		expected map[string]any
	}{
		{
			name:     "replace scalar",
			payload:  NewPatchPayload("replace", "/description", "updated"),
			expected: map[string]any{"op": "replace", "path": "/description", "value": "updated"},
		},
		{
			name:     "add object",
			payload:  NewPatchPayload("add", "/interfaces/-", map[string]any{"name": "ethernet1/1"}),
			expected: map[string]any{"op": "add", "path": "/interfaces/-", "value": map[string]any{"name": "ethernet1/1"}},
		},
		{
			name:     "replace explicit null",
			payload:  NewPatchPayload("replace", "/value", nil),
			expected: map[string]any{"op": "replace", "path": "/value", "value": nil},
		},
		{
			name:     "remove omits value",
			payload:  NewRemovePatchPayload("/description"),
			expected: map[string]any{"op": "remove", "path": "/description"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := test.payload.ToMap()
			if err != nil {
				t.Fatalf("converting patch payload: %v", err)
			}
			if !reflect.DeepEqual(actual, test.expected) {
				t.Fatalf("unexpected patch payload: got %#v, want %#v", actual, test.expected)
			}
		})
	}
}

func TestPatchOperations(t *testing.T) {
	operations := NewPatchOperations(map[string]any{
		"description": "existing",
		"pim":         false,
		"targetDscp":  "af11",
	}, "")
	operations.Set("description", "existing")
	operations.Set("l3domain", "external")
	operations.Set("pim", true)
	operations.Remove("targetDscp")
	operations.Remove("missing")

	actual := operations.Operations()
	expected := []*PatchPayload{
		NewPatchPayload("add", "/l3domain", "external"),
		NewPatchPayload("replace", "/pim", true),
		NewRemovePatchPayload("/targetDscp"),
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("unexpected patch operations: got %#v, want %#v", actual, expected)
	}
}

func TestPatchOperationsSetString(t *testing.T) {
	operations := NewPatchOperations(map[string]any{
		"description": "existing",
		"targetDscp":  "af11",
	}, "")
	operations.SetString("description", nil, nil, true)
	operations.SetString("l3domain", stringPointer("external"), nil, true)
	operations.SetString("targetDscp", stringPointer("expedited_forwarding"), func(value string) string {
		return "expeditedForwarding"
	}, true)
	operations.SetString("description", stringPointer(""), nil, true)

	actual := operations.Operations()
	expected := []*PatchPayload{
		NewPatchPayload("add", "/l3domain", "external"),
		NewPatchPayload("replace", "/targetDscp", "expeditedForwarding"),
		NewRemovePatchPayload("/description"),
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("unexpected patch operations: got %#v, want %#v", actual, expected)
	}
}

func stringPointer(value string) *string {
	return &value
}

func TestPatchOperationsPaths(t *testing.T) {
	operations := NewPatchOperations(map[string]any{
		"description": "existing",
		"targetDscp":  "af11",
	}, "/l3outTemplate/l3outs/2")
	operations.Set("description", "updated")
	operations.Remove("targetDscp")
	patchOperations := make([]map[string]any, 0, len(operations.Operations()))
	for index, operation := range operations.Operations() {
		payload, err := operation.ToMap()
		if err != nil {
			t.Fatalf("converting patch operation %d: %v", index, err)
		}
		patchOperations = append(patchOperations, payload)
	}

	expected := []map[string]any{
		{"op": "replace", "path": "/l3outTemplate/l3outs/2/description", "value": "updated"},
		{"op": "remove", "path": "/l3outTemplate/l3outs/2/targetDscp"},
	}
	if !reflect.DeepEqual(patchOperations, expected) {
		t.Fatalf("unexpected prefixed operations: got %#v, want %#v", patchOperations, expected)
	}
	if operations.Operations()[0].Path != "/l3outTemplate/l3outs/2/description" || operations.Operations()[1].Path != "/l3outTemplate/l3outs/2/targetDscp" {
		t.Fatal("unexpected operation paths")
	}
}

func TestPatchOperationsApplyNoOp(t *testing.T) {
	operations := NewPatchOperations(nil, "")

	if err := operations.Apply(nil, ""); err != nil {
		t.Fatalf("applying empty patch operations: %v", err)
	}
}
