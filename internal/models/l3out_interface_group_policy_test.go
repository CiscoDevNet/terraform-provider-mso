package models

import (
	"context"
	"reflect"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-mso/internal/ndoapi"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestInterfaceGroupPoliciesToPayloadQoSPriority(t *testing.T) {
	tests := []struct {
		name          string
		priority      types.String
		configuration types.String
		wantPriority  string
		wantPresent   bool
	}{
		{
			name:          "unknown computed priority is omitted",
			priority:      types.StringUnknown(),
			configuration: types.StringNull(),
		},
		{
			name:          "null priority is omitted",
			priority:      types.StringNull(),
			configuration: types.StringNull(),
		},
		{
			name:          "configured priority is included",
			priority:      types.StringValue("level6"),
			configuration: types.StringValue("level6"),
			wantPriority:  "level6",
			wantPresent:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			policy := InterfaceGroupPolicyModel{
				Description:                types.StringNull(),
				InterfaceRoutingPolicyUUID: types.StringNull(),
				CustomQoSPolicyUUID:        types.StringNull(),
				QoSPriority:                test.priority,
				NetFlowMonitorUUIDs:        types.MapNull(types.StringType),
				BFD:                        types.ObjectNull(interfaceGroupBFDTypes),
				BFDMultiHop:                types.ObjectNull(interfaceGroupBFDTypes),
				OSPF:                       types.ObjectNull(interfaceGroupOSPFTypes),
			}
			configuration := policy
			configuration.QoSPriority = test.configuration
			var diagnostics diag.Diagnostics
			payload := (InterfaceGroupPoliciesModel{"edge": policy}).ToPayload(
				context.Background(),
				InterfaceGroupPoliciesModel{"edge": configuration},
				&diagnostics,
			)
			if diagnostics.HasError() {
				t.Fatalf("unexpected payload diagnostics: %v", diagnostics)
			}
			if len(payload) != 1 || payload[0]["name"] != "edge" {
				t.Fatalf("unexpected interface group payload: %#v", payload)
			}
			priority, present := payload[0]["qosPriority"]
			if present != test.wantPresent || (present && priority != test.wantPriority) {
				t.Fatalf("unexpected QoS priority in payload: value=%#v, present=%t", priority, present)
			}
		})
	}
}

func TestInterfaceGroupPoliciesAddPatchOperationsRemovesSeveralGroups(t *testing.T) {
	existing := map[string]any{
		"interfaceGroups": []any{
			map[string]any{"name": "remove-first"},
			map[string]any{"name": "keep", "description": "before"},
			map[string]any{"name": "remove-second"},
			map[string]any{"name": "remove-third"},
		},
	}
	policy := InterfaceGroupPolicyModel{
		Description:                types.StringValue("after"),
		InterfaceRoutingPolicyUUID: types.StringNull(),
		CustomQoSPolicyUUID:        types.StringNull(),
		QoSPriority:                types.StringNull(),
		NetFlowMonitorUUIDs:        types.MapNull(types.StringType),
		BFD:                        types.ObjectNull(interfaceGroupBFDTypes),
		BFDMultiHop:                types.ObjectNull(interfaceGroupBFDTypes),
		OSPF:                       types.ObjectNull(interfaceGroupOSPFTypes),
	}
	policies := InterfaceGroupPoliciesModel{"keep": policy}
	operations := ndoapi.NewPatchOperations(existing, "/l3outs/0")
	var diagnostics diag.Diagnostics
	if err := policies.AddPatchOperations(
		context.Background(),
		existing,
		operations,
		"/l3outs/0",
		policies,
		nil,
		&diagnostics,
	); err != nil {
		t.Fatalf("unexpected patch error: %v", err)
	}
	if diagnostics.HasError() {
		t.Fatalf("unexpected patch diagnostics: %v", diagnostics)
	}
	want := []*ndoapi.PatchPayload{
		ndoapi.NewPatchPayload("replace", "/l3outs/0/interfaceGroups/1/description", "after"),
		ndoapi.NewRemovePatchPayload("/l3outs/0/interfaceGroups/3"),
		ndoapi.NewRemovePatchPayload("/l3outs/0/interfaceGroups/2"),
		ndoapi.NewRemovePatchPayload("/l3outs/0/interfaceGroups/0"),
	}
	if got := operations.Operations(); !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected patch operations:\ngot  %#v\nwant %#v", got, want)
	}
}

func TestInterfaceGroupPoliciesAddPatchOperationsAppendsBeforeRemovals(t *testing.T) {
	existing := map[string]any{
		"interfaceGroups": []any{
			map[string]any{"name": "remove-first"},
			map[string]any{"name": "keep", "description": "before"},
			map[string]any{"name": "remove-last"},
		},
	}
	policy := InterfaceGroupPolicyModel{
		Description:                types.StringValue("after"),
		InterfaceRoutingPolicyUUID: types.StringNull(),
		CustomQoSPolicyUUID:        types.StringNull(),
		QoSPriority:                types.StringNull(),
		NetFlowMonitorUUIDs:        types.MapNull(types.StringType),
		BFD:                        types.ObjectNull(interfaceGroupBFDTypes),
		BFDMultiHop:                types.ObjectNull(interfaceGroupBFDTypes),
		OSPF:                       types.ObjectNull(interfaceGroupOSPFTypes),
	}
	newAPolicy := policy
	newAPolicy.Description = types.StringValue("new a")
	newZPolicy := policy
	newZPolicy.Description = types.StringValue("new z")
	policies := InterfaceGroupPoliciesModel{
		"new-z": newZPolicy,
		"keep":  policy,
		"new-a": newAPolicy,
	}
	operations := ndoapi.NewPatchOperations(existing, "/l3outs/0")
	var diagnostics diag.Diagnostics
	if err := policies.AddPatchOperations(
		context.Background(),
		existing,
		operations,
		"/l3outs/0",
		policies,
		nil,
		&diagnostics,
	); err != nil {
		t.Fatalf("unexpected patch error: %v", err)
	}
	if diagnostics.HasError() {
		t.Fatalf("unexpected patch diagnostics: %v", diagnostics)
	}
	want := []*ndoapi.PatchPayload{
		ndoapi.NewPatchPayload("replace", "/l3outs/0/interfaceGroups/1/description", "after"),
		ndoapi.NewPatchPayload("add", "/l3outs/0/interfaceGroups/-", map[string]any{"name": "new-a", "description": "new a"}),
		ndoapi.NewPatchPayload("add", "/l3outs/0/interfaceGroups/-", map[string]any{"name": "new-z", "description": "new z"}),
		ndoapi.NewRemovePatchPayload("/l3outs/0/interfaceGroups/2"),
		ndoapi.NewRemovePatchPayload("/l3outs/0/interfaceGroups/0"),
	}
	if got := operations.Operations(); !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected patch operations:\ngot  %#v\nwant %#v", got, want)
	}
}
