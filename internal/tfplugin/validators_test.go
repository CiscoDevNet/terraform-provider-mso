package tfplugin

import (
	"context"
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
		name     string
		config   types.Object
		wantPath string
	}{
		{name: "empty object", config: object(types.BoolNull(), types.StringNull())},
		{name: "explicitly disabled", config: object(types.BoolValue(false), types.StringNull())},
		{
			name:     "settings require enablement",
			config:   object(types.BoolNull(), types.StringValue("0.0.0.1")),
			wantPath: "ospf.area_id",
		},
		{
			name:     "enabled requires area",
			config:   object(types.BoolValue(true), types.StringNull()),
			wantPath: "ospf.area_id",
		},
		{name: "enabled with area", config: object(types.BoolValue(true), types.StringValue("0.0.0.1"))},
		{name: "unknown area defers validation", config: object(types.BoolValue(true), types.StringUnknown())},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := &validator.ObjectResponse{}
			RequireAttributesWhenEnabled("enabled", "area_id").ValidateObject(
				context.Background(),
				validator.ObjectRequest{
					Path:        path.Root("ospf"),
					ConfigValue: test.config,
				},
				response,
			)
			if test.wantPath == "" {
				if response.Diagnostics.HasError() {
					t.Fatalf("unexpected diagnostics: %v", response.Diagnostics)
				}
				return
			}
			if len(response.Diagnostics.Errors()) != 1 {
				t.Fatalf("expected one diagnostic at %s, got %v", test.wantPath, response.Diagnostics)
			}
			diagnostic, ok := response.Diagnostics.Errors()[0].(diag.DiagnosticWithPath)
			if !ok || diagnostic.Path().String() != test.wantPath {
				t.Fatalf(
					"unexpected diagnostic path: got %#v, want %q",
					response.Diagnostics.Errors()[0],
					test.wantPath,
				)
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
