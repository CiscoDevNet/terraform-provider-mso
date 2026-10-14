package tfplugin

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestUseConfiguredDisabledObject(t *testing.T) {
	attributeTypes := map[string]attr.Type{"enabled": types.BoolType, "name": types.StringType}
	object := func(enabled types.Bool, name types.String) types.Object {
		return types.ObjectValueMust(attributeTypes, map[string]attr.Value{"enabled": enabled, "name": name})
	}
	previous := object(types.BoolValue(true), types.StringValue("previous"))
	disabled := object(types.BoolValue(false), types.StringNull())
	tests := []struct {
		name   string
		config types.Object
		want   types.Object
	}{
		{name: "empty object disables and clears prior settings", config: object(types.BoolNull(), types.StringNull()), want: disabled},
		{name: "explicit false clears prior settings", config: disabled, want: disabled},
		{name: "enabled object retains planned settings", config: object(types.BoolValue(true), types.StringNull()), want: previous},
		{name: "omitted object retains planned settings", config: types.ObjectNull(attributeTypes), want: previous},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := &planmodifier.ObjectResponse{PlanValue: previous}
			UseConfiguredDisabledObject("enabled").PlanModifyObject(context.Background(), planmodifier.ObjectRequest{
				ConfigValue: test.config,
				PlanValue:   previous,
			}, response)
			if response.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", response.Diagnostics)
			}
			if !response.PlanValue.Equal(test.want) {
				t.Fatalf("unexpected plan value: got %#v, want %#v", response.PlanValue, test.want)
			}
		})
	}
}

func TestUseNonNullStateForUnknownInt64(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		name                string
		config, plan, state types.Int64
		want                types.Int64
	}{
		{name: "new omitted integer remains unknown", config: types.Int64Null(), plan: types.Int64Unknown(), state: types.Int64Null(), want: types.Int64Unknown()},
		{name: "prior integer is retained", config: types.Int64Null(), plan: types.Int64Unknown(), state: types.Int64Value(10), want: types.Int64Value(10)},
		{name: "configured integer is retained", config: types.Int64Value(5), plan: types.Int64Value(5), state: types.Int64Null(), want: types.Int64Value(5)},
		{name: "unknown configuration remains unknown", config: types.Int64Unknown(), plan: types.Int64Unknown(), state: types.Int64Value(10), want: types.Int64Unknown()},
		{name: "known null plan is retained", config: types.Int64Null(), plan: types.Int64Null(), state: types.Int64Value(10), want: types.Int64Null()},
	} {
		t.Run(test.name, func(t *testing.T) {
			response := &planmodifier.Int64Response{PlanValue: test.plan}
			UseNonNullStateForUnknownInt64().PlanModifyInt64(ctx, planmodifier.Int64Request{
				ConfigValue: test.config,
				PlanValue:   test.plan,
				StateValue:  test.state,
			}, response)
			if !response.PlanValue.Equal(test.want) {
				t.Fatalf("got %v, want %v", response.PlanValue, test.want)
			}
		})
	}
}
