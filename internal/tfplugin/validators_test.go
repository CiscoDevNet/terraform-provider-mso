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
		{name: "settings require enablement", config: object(types.BoolNull(), types.StringValue("0.0.0.1")), wantPath: "ospf.area_id"},
		{name: "enabled requires area", config: object(types.BoolValue(true), types.StringNull()), wantPath: "ospf.area_id"},
		{name: "enabled with area", config: object(types.BoolValue(true), types.StringValue("0.0.0.1"))},
		{name: "unknown area defers validation", config: object(types.BoolValue(true), types.StringUnknown())},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := &validator.ObjectResponse{}
			RequireAttributesWhenEnabled("enabled", "area_id").ValidateObject(context.Background(), validator.ObjectRequest{
				Path:        path.Root("ospf"),
				ConfigValue: test.config,
			}, response)
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
				t.Fatalf("unexpected diagnostic path: got %#v, want %q", response.Diagnostics.Errors()[0], test.wantPath)
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
	AllowAttributesOnlyWhenValue("authentication_enabled", []attr.Value{types.BoolValue(true)}, "key").ValidateObject(
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
