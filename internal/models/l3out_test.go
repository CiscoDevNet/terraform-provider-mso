package models

import (
	"context"
	"strings"
	"testing"
)

func TestL3OutRoutingProtocolState(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		bgp       bool
		ospf      bool
		wantError bool
	}{
		{name: "omitted", value: ""},
		{name: "none", value: "none"},
		{name: "BGP", value: "bgp", bgp: true},
		{name: "OSPF", value: "ospf", ospf: true},
		{name: "BGP and OSPF", value: "bgpOspf", bgp: true, ospf: true},
		{name: "invalid reversed spelling", value: "ospfbgp", wantError: true},
		{name: "unknown", value: "unknown", wantError: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state, err := newL3OutRoutingProtocolState(test.value)
			if test.wantError {
				if err == nil {
					t.Fatalf("expected routing protocol %q to be rejected", test.value)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected routing protocol error: %v", err)
			}
			if state.bgp != test.bgp || state.ospf != test.ospf {
				t.Fatalf("unexpected routing protocol state: %+v", state)
			}
		})
	}
}

func TestL3OutSetFromNDOObjectRejectsRequiredFields(t *testing.T) {
	tests := []struct {
		name     string
		change   func(map[string]any)
		expected string
	}{
		{name: "missing UUID", change: func(object map[string]any) { delete(object, "uuid") }, expected: `"uuid" is required`},
		{name: "invalid name", change: func(object map[string]any) { object["name"] = 42 }, expected: `"name" has unexpected type int`},
		{name: "missing VRF", change: func(object map[string]any) { delete(object, "vrfRef") }, expected: `"vrfRef" is required`},
		{name: "invalid description", change: func(object map[string]any) { object["description"] = 42 }, expected: `"description" has unexpected type int`},
		{name: "invalid routing protocol type", change: func(object map[string]any) { object["routingProtocol"] = 42 }, expected: `"routingProtocol" has unexpected type int`},
		{name: "unknown routing protocol", change: func(object map[string]any) { object["routingProtocol"] = "unknown" }, expected: `unexpected NDO routingProtocol`},
		{name: "invalid default route", change: func(object map[string]any) { object["defaultRouteLeak"] = "invalid" }, expected: `"defaultRouteLeak" has unexpected type string`},
		{name: "invalid annotations", change: func(object map[string]any) { object["tagAnnotations"] = "invalid" }, expected: `"tagAnnotations" has unexpected type string`},
		{name: "invalid interface groups", change: func(object map[string]any) { object["interfaceGroups"] = "invalid" }, expected: `"interfaceGroups" has unexpected type string`},
		{name: "invalid interface group entry", change: func(object map[string]any) { object["interfaceGroups"] = []any{"invalid"} }, expected: `interfaceGroups[0] has unexpected type string`},
		{name: "duplicate interface group name", change: func(object map[string]any) {
			object["interfaceGroups"] = []any{map[string]any{"name": "edge"}, map[string]any{"name": "edge"}}
		}, expected: `duplicate name "edge"`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			object := map[string]any{
				"uuid":            "l3out-123",
				"name":            "edge",
				"vrfRef":          "vrf-123",
				"routingProtocol": "none",
			}
			test.change(object)
			var model L3OutModel
			if err := model.SetFromNDOObject(context.Background(), "template-123", object, L3OutModel{}); err == nil || !strings.Contains(err.Error(), test.expected) {
				t.Fatalf("expected error containing %q, got %v", test.expected, err)
			}
		})
	}
}
