package tfplugin

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UseConfiguredDisabledObject prevents computed settings from being carried
// forward when an object is explicitly disabled, including with {}.
func UseConfiguredDisabledObject(enabledAttribute string) planmodifier.Object {
	return useConfiguredDisabledObject{enabledAttribute: enabledAttribute}
}

type useConfiguredDisabledObject struct {
	enabledAttribute string
}

func (disabled useConfiguredDisabledObject) Description(context.Context) string {
	return "Clear prior computed settings when this object is explicitly disabled."
}

func (disabled useConfiguredDisabledObject) MarkdownDescription(ctx context.Context) string {
	return disabled.Description(ctx)
}

func (disabled useConfiguredDisabledObject) PlanModifyObject(ctx context.Context, request planmodifier.ObjectRequest, response *planmodifier.ObjectResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	attributes := request.ConfigValue.Attributes()
	enabled, ok := attributes[disabled.enabledAttribute].(types.Bool)
	if !ok || enabled.IsUnknown() || enabled.ValueBool() {
		return
	}
	for name, value := range attributes {
		if name != disabled.enabledAttribute && !value.IsNull() {
			return
		}
	}
	configured := make(map[string]attr.Value, len(attributes))
	for name, value := range attributes {
		configured[name] = value
	}
	configured[disabled.enabledAttribute] = types.BoolValue(false)
	response.PlanValue, response.Diagnostics = types.ObjectValue(request.ConfigValue.AttributeTypes(ctx), configured)
}

// SetBoolWhenStringEquals sets an omitted boolean to value when a configured
// string at target equals match. The target may be relative or rooted. Other
// planned values are preserved.
func SetBoolWhenStringEquals(target path.Expression, match string, value bool) planmodifier.Bool {
	return setBoolWhenStringEquals{target: target, match: match, value: value}
}

type setBoolWhenStringEquals struct {
	target path.Expression
	match  string
	value  bool
}

func (modifier setBoolWhenStringEquals) Description(context.Context) string {
	return fmt.Sprintf("Set this boolean to %t when %s is %q.", modifier.value, modifier.target, modifier.match)
}

func (modifier setBoolWhenStringEquals) MarkdownDescription(ctx context.Context) string {
	return modifier.Description(ctx)
}

func (modifier setBoolWhenStringEquals) PlanModifyBool(ctx context.Context, request planmodifier.BoolRequest, response *planmodifier.BoolResponse) {
	if !request.ConfigValue.IsNull() || (!request.PlanValue.IsNull() && !request.PlanValue.IsUnknown()) {
		return
	}
	target, ok := configTargetPath(ctx, request.Config, request.Path, modifier.target, "Invalid Plan Modifier Configuration", &response.Diagnostics)
	if !ok {
		return
	}
	var configured types.String
	response.Diagnostics.Append(request.Config.GetAttribute(ctx, target, &configured)...)
	if response.Diagnostics.HasError() || configured.IsNull() || configured.IsUnknown() {
		return
	}
	if configured.ValueString() == modifier.match {
		response.PlanValue = types.BoolValue(modifier.value)
	}
}

// NullBoolWhenParentDisabledOrAbsent clears an omitted boolean when its parent
// object is configured without enabled=true, or when both the parent config
// and the prior boolean are null. The parent path may be relative or rooted.
func NullBoolWhenParentDisabledOrAbsent(parent path.Expression, enabledAttribute string) planmodifier.Bool {
	return nullBoolWhenParentDisabledOrAbsent{parent: parent, enabledAttribute: enabledAttribute}
}

type nullBoolWhenParentDisabledOrAbsent struct {
	parent           path.Expression
	enabledAttribute string
}

func (modifier nullBoolWhenParentDisabledOrAbsent) Description(context.Context) string {
	return fmt.Sprintf("Clear this boolean when %s.%s is disabled or the parent has no prior value.", modifier.parent, modifier.enabledAttribute)
}

func (modifier nullBoolWhenParentDisabledOrAbsent) MarkdownDescription(ctx context.Context) string {
	return modifier.Description(ctx)
}

func (modifier nullBoolWhenParentDisabledOrAbsent) PlanModifyBool(ctx context.Context, request planmodifier.BoolRequest, response *planmodifier.BoolResponse) {
	if !request.ConfigValue.IsNull() {
		return
	}
	parentPath, ok := configTargetPath(ctx, request.Config, request.Path, modifier.parent, "Invalid Plan Modifier Configuration", &response.Diagnostics)
	if !ok {
		return
	}
	var parent types.Object
	response.Diagnostics.Append(request.Config.GetAttribute(ctx, parentPath, &parent)...)
	if response.Diagnostics.HasError() || parent.IsUnknown() {
		return
	}
	if parent.IsNull() {
		if request.StateValue.IsNull() {
			response.PlanValue = types.BoolNull()
		}
		return
	}
	enabled, ok := parent.Attributes()[modifier.enabledAttribute].(types.Bool)
	if !ok {
		response.Diagnostics.AddError(
			"Invalid Plan Modifier Configuration",
			fmt.Sprintf("%s.%s must be a boolean attribute.", parentPath, modifier.enabledAttribute),
		)
		return
	}
	if enabled.IsNull() || (!enabled.IsUnknown() && !enabled.ValueBool()) {
		response.PlanValue = types.BoolNull()
	}
}

// UseNonNullStateForUnknownInt64 retains a known prior integer while leaving
// a new omitted value unknown for the read after apply.
func UseNonNullStateForUnknownInt64() planmodifier.Int64 {
	return nonNullStateForUnknownInt64{}
}

type nonNullStateForUnknownInt64 struct{}

func (modifier nonNullStateForUnknownInt64) Description(context.Context) string {
	return "Retain a non-null prior integer when the planned value is unknown."
}

func (modifier nonNullStateForUnknownInt64) MarkdownDescription(ctx context.Context) string {
	return modifier.Description(ctx)
}

func (modifier nonNullStateForUnknownInt64) PlanModifyInt64(_ context.Context, request planmodifier.Int64Request, response *planmodifier.Int64Response) {
	if request.StateValue.IsNull() || request.StateValue.IsUnknown() || !request.PlanValue.IsUnknown() || request.ConfigValue.IsUnknown() {
		return
	}
	response.PlanValue = request.StateValue
}
