package tfplugin

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestRequireAttributesWhenEnabled(t *testing.T) {
	attributeTypes := map[string]attr.Type{"enabled": types.BoolType, "area_id": types.StringType}
	object := func(enabled types.Bool, areaID types.String) types.Object {
		return types.ObjectValueMust(attributeTypes, map[string]attr.Value{"enabled": enabled, "area_id": areaID})
	}
	tests := []struct {
		name             string
		config           types.Object
		wantDetailFormat string
	}{
		{name: "empty object", config: object(types.BoolNull(), types.StringNull())},
		{name: "explicitly disabled", config: object(types.BoolValue(false), types.StringNull())},
		{
			name:             "settings require enablement",
			config:           object(types.BoolNull(), types.StringValue("0.0.0.1")),
			wantDetailFormat: "%s.area_id requires %s.enabled = true.",
		},
		{
			name:             "explicitly disabled with settings",
			config:           object(types.BoolValue(false), types.StringValue("0.0.0.1")),
			wantDetailFormat: "%s.area_id requires %s.enabled = true.",
		},
		{
			name:             "enabled requires area",
			config:           object(types.BoolValue(true), types.StringNull()),
			wantDetailFormat: "%s.area_id must be configured when %s.enabled is true.",
		},
		{
			name:             "enabled rejects empty area",
			config:           object(types.BoolValue(true), types.StringValue("")),
			wantDetailFormat: "%s.area_id must be configured when %s.enabled is true.",
		},
		{name: "enabled with area", config: object(types.BoolValue(true), types.StringValue("0.0.0.1"))},
		{name: "unknown area defers validation", config: object(types.BoolValue(true), types.StringUnknown())},
	}
	objectPaths := []struct {
		name       string
		objectPath path.Path
		wantPrefix string
	}{
		{name: "OSPF object", objectPath: path.Root("ospf"), wantPrefix: "ospf"},
		{
			name:       "nested map object",
			objectPath: path.Root("groups").AtMapKey("edge").AtName("ospf"),
			wantPrefix: `groups["edge"].ospf`,
		},
		{
			name:       "nested list object",
			objectPath: path.Root("groups").AtListIndex(0).AtName("ospf"),
			wantPrefix: "groups[0].ospf",
		},
	}
	rule := RequireAttributesWhenEnabled("enabled", "area_id")
	for _, objectPath := range objectPaths {
		t.Run(objectPath.name, func(t *testing.T) {
			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					response := &validator.ObjectResponse{}
					rule.ValidateObject(
						context.Background(),
						validator.ObjectRequest{
							Path:        objectPath.objectPath,
							ConfigValue: test.config,
						},
						response,
					)
					if test.wantDetailFormat == "" {
						if response.Diagnostics.HasError() {
							t.Fatalf("unexpected diagnostics: %v", response.Diagnostics)
						}
						return
					}
					wantPath := fmt.Sprintf("%s.area_id", objectPath.wantPrefix)
					if len(response.Diagnostics.Errors()) != 1 {
						t.Fatalf("expected one diagnostic at %s, got %v", wantPath, response.Diagnostics)
					}
					diagnostic, ok := response.Diagnostics.Errors()[0].(diag.DiagnosticWithPath)
					if !ok || diagnostic.Path().String() != wantPath {
						t.Fatalf(
							"unexpected diagnostic path: got %#v, want %q",
							response.Diagnostics.Errors()[0],
							wantPath,
						)
					}
					wantDetail := fmt.Sprintf(test.wantDetailFormat, objectPath.wantPrefix, objectPath.wantPrefix)
					if diagnostic.Detail() != wantDetail {
						t.Fatalf("unexpected diagnostic detail: got %q, want %q", diagnostic.Detail(), wantDetail)
					}
				})
			}
		})
	}
}

func TestAllowAttributesOnlyWhenValueMissingTrigger(t *testing.T) {
	config := types.ObjectValueMust(
		map[string]attr.Type{"key": types.StringType},
		map[string]attr.Value{"key": types.StringValue("secret")},
	)
	response := &validator.ObjectResponse{}
	AllowAttributesOnlyWhenValue("authentication_enabled", []attr.Value{
		types.BoolValue(true),
	}, nil, "key").ValidateObject(
		context.Background(),
		validator.ObjectRequest{Path: path.Root("bfd"), ConfigValue: config},
		response,
	)
	if len(response.Diagnostics.Errors()) != 1 {
		t.Fatalf("expected one diagnostic, got %v", response.Diagnostics)
	}
	diagnostic, ok := response.Diagnostics.Errors()[0].(diag.DiagnosticWithPath)
	if !ok || diagnostic.Path().String() != "bfd.key" {
		t.Fatalf("unexpected diagnostic path: got %#v, want %q", response.Diagnostics.Errors()[0], "bfd.key")
	}
	wantDetail := "bfd.key may be set only when bfd.authentication_enabled is one of true"
	if diagnostic.Detail() != wantDetail {
		t.Fatalf("unexpected diagnostic detail: got %q, want %q", diagnostic.Detail(), wantDetail)
	}
}

func TestAllowAttributesOnlyWhenValueDiagnosticPaths(t *testing.T) {
	ctx := context.Background()
	rule := AllowAttributesOnlyWhenValue(
		"enabled",
		[]attr.Value{types.BoolValue(true)},
		[]attr.Value{types.StringValue("")},
		"route_dampening_ipv4_uuid",
		"route_dampening_ipv6_uuid",
	)
	config := types.ObjectValueMust(
		map[string]attr.Type{
			"enabled":                   types.BoolType,
			"route_dampening_ipv4_uuid": types.StringType,
			"route_dampening_ipv6_uuid": types.StringType,
		},
		map[string]attr.Value{
			"enabled":                   types.BoolValue(false),
			"route_dampening_ipv4_uuid": types.StringValue("ipv4-policy-uuid"),
			"route_dampening_ipv6_uuid": types.StringValue("ipv6-policy-uuid"),
		},
	)
	tests := []struct {
		name       string
		objectPath path.Path
		wantPrefix string
	}{
		{name: "BGP object", objectPath: path.Root("bgp"), wantPrefix: "bgp"},
		{
			name:       "nested map object",
			objectPath: path.Root("groups").AtMapKey("edge").AtName("bgp"),
			wantPrefix: `groups["edge"].bgp`,
		},
		{
			name:       "nested list object",
			objectPath: path.Root("groups").AtListIndex(0).AtName("bgp"),
			wantPrefix: "groups[0].bgp",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := &validator.ObjectResponse{}
			rule.ValidateObject(ctx, validator.ObjectRequest{Path: test.objectPath, ConfigValue: config}, response)
			if len(response.Diagnostics.Errors()) != 2 {
				t.Fatalf("expected two diagnostics, got %v", response.Diagnostics)
			}
			wantDetail := fmt.Sprintf(
				`%s.route_dampening_ipv4_uuid, %s.route_dampening_ipv6_uuid may be set only when %s.enabled is one of true; attribute values "" are exempt from this condition`,
				test.wantPrefix,
				test.wantPrefix,
				test.wantPrefix,
			)
			for index, name := range []string{"route_dampening_ipv4_uuid", "route_dampening_ipv6_uuid"} {
				diagnostic, ok := response.Diagnostics.Errors()[index].(diag.DiagnosticWithPath)
				wantPath := fmt.Sprintf("%s.%s", test.wantPrefix, name)
				if !ok || diagnostic.Path().String() != wantPath {
					t.Fatalf(
						"unexpected diagnostic path: got %#v, want %q",
						response.Diagnostics.Errors()[index],
						wantPath,
					)
				}
				if diagnostic.Detail() != wantDetail {
					t.Fatalf("unexpected diagnostic detail: got %q, want %q", diagnostic.Detail(), wantDetail)
				}
			}
		})
	}
}

func TestAllowAttributesOnlyWhenValueExceptions(t *testing.T) {
	ctx := context.Background()
	emptyString := []attr.Value{types.StringValue("")}
	tests := []struct {
		name        string
		trigger     types.Bool
		value       attr.Value
		exceptions  []attr.Value
		omitTrigger bool
		wantError   bool
	}{
		{
			name:      "strict rejects empty string",
			trigger:   types.BoolValue(false),
			value:     types.StringValue(""),
			wantError: true,
		},
		{
			name:      "strict omitted trigger rejects empty string",
			trigger:   types.BoolNull(),
			value:     types.StringValue(""),
			wantError: true,
		},
		{
			name:       "empty clear while disabled",
			trigger:    types.BoolValue(false),
			value:      types.StringValue(""),
			exceptions: emptyString,
		},
		{
			name:       "empty clear with omitted trigger",
			trigger:    types.BoolNull(),
			value:      types.StringValue(""),
			exceptions: emptyString,
		},
		{
			name:        "empty clear with missing trigger",
			value:       types.StringValue(""),
			exceptions:  emptyString,
			omitTrigger: true,
		},
		{
			name:       "nonempty value while disabled",
			trigger:    types.BoolValue(false),
			value:      types.StringValue("policy-uuid"),
			exceptions: emptyString,
			wantError:  true,
		},
		{
			name:       "nonempty value with omitted trigger",
			trigger:    types.BoolNull(),
			value:      types.StringValue("policy-uuid"),
			exceptions: emptyString,
			wantError:  true,
		},
		{
			name:        "nonempty value with missing trigger",
			value:       types.StringValue("policy-uuid"),
			exceptions:  emptyString,
			omitTrigger: true,
			wantError:   true,
		},
		{
			name:    "matching trigger allows value",
			trigger: types.BoolValue(true),
			value:   types.StringValue("policy-uuid"),
		},
		{
			name:       "unknown trigger defers validation",
			trigger:    types.BoolUnknown(),
			value:      types.StringValue("policy-uuid"),
			exceptions: emptyString,
		},
		{
			name:       "unknown value defers validation",
			trigger:    types.BoolValue(false),
			value:      types.StringUnknown(),
			exceptions: emptyString,
		},
		{
			name:       "null value is allowed",
			trigger:    types.BoolValue(false),
			value:      types.StringNull(),
			exceptions: emptyString,
		},
		{
			name:       "string exception does not allow false",
			trigger:    types.BoolValue(false),
			value:      types.BoolValue(false),
			exceptions: emptyString,
			wantError:  true,
		},
		{
			name:       "string exception does not allow zero",
			trigger:    types.BoolValue(false),
			value:      types.Int64Value(0),
			exceptions: emptyString,
			wantError:  true,
		},
		{
			name:       "explicit false exception",
			trigger:    types.BoolValue(false),
			value:      types.BoolValue(false),
			exceptions: []attr.Value{types.BoolValue(false)},
		},
		{
			name:       "explicit zero exception",
			trigger:    types.BoolValue(false),
			value:      types.Int64Value(0),
			exceptions: []attr.Value{types.Int64Value(0)},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			attributeTypes := map[string]attr.Type{"reference": test.value.Type(ctx)}
			attributes := map[string]attr.Value{"reference": test.value}
			if !test.omitTrigger {
				attributeTypes["enabled"] = types.BoolType
				attributes["enabled"] = test.trigger
			}
			response := &validator.ObjectResponse{}
			AllowAttributesOnlyWhenValue("enabled", []attr.Value{
				types.BoolValue(true),
			}, test.exceptions, "reference").ValidateObject(
				ctx,
				validator.ObjectRequest{
					Path:        path.Root("settings"),
					ConfigValue: types.ObjectValueMust(attributeTypes, attributes),
				},
				response,
			)
			if !test.wantError {
				if response.Diagnostics.HasError() {
					t.Fatalf("unexpected diagnostics: %v", response.Diagnostics)
				}
				return
			}
			if len(response.Diagnostics.Errors()) != 1 {
				t.Fatalf("expected one diagnostic, got %v", response.Diagnostics)
			}
			diagnostic, ok := response.Diagnostics.Errors()[0].(diag.DiagnosticWithPath)
			if !ok || diagnostic.Path().String() != "settings.reference" {
				t.Fatalf(
					"unexpected diagnostic path: got %#v, want settings.reference",
					response.Diagnostics.Errors()[0],
				)
			}
		})
	}
}
