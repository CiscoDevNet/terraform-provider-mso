package tfplugin

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// RequireValueWhenBool requires a value at target when the validated boolean
// equals when. If values are provided, the target must equal one of them.
// The target may be relative to the validated attribute or rooted at the
// resource. Empty strings do not satisfy the requirement; unknowns are deferred.
func RequireValueWhenBool(when bool, target path.Expression, values ...attr.Value) validator.Bool {
	return requireValueWhenBool{
		when:   when,
		target: target,
		values: append([]attr.Value(nil), values...),
	}
}

type requireValueWhenBool struct {
	when   bool
	target path.Expression
	values []attr.Value
}

func (required requireValueWhenBool) Description(context.Context) string {
	return required.description("this attribute", required.target.String())
}

func (required requireValueWhenBool) description(attribute, targetName string) string {
	if len(required.values) == 0 {
		return fmt.Sprintf("%s must have a value when %s is %t", targetName, attribute, required.when)
	}
	values := make([]string, len(required.values))
	for index, value := range required.values {
		values[index] = value.String()
	}
	return fmt.Sprintf("%s must be one of %s when %s is %t", targetName, strings.Join(values, ", "), attribute, required.when)
}

func (required requireValueWhenBool) MarkdownDescription(ctx context.Context) string {
	return required.Description(ctx)
}

func (required requireValueWhenBool) ValidateBool(ctx context.Context, request validator.BoolRequest, response *validator.BoolResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() || request.ConfigValue.ValueBool() != required.when {
		return
	}

	target, ok := configTargetPath(ctx, request.Config, request.Path, required.target, "Invalid Validator Configuration", &response.Diagnostics)
	if !ok {
		return
	}
	var value attr.Value
	response.Diagnostics.Append(request.Config.GetAttribute(ctx, target, &value)...)
	if response.Diagnostics.HasError() || value.IsUnknown() {
		return
	}
	if value.IsNull() {
		response.Diagnostics.AddAttributeError(request.Path, "Invalid Required Attribute Value", required.description(request.Path.String(), target.String()))
		return
	}
	if stringValue, ok := value.(types.String); ok && stringValue.ValueString() == "" {
		response.Diagnostics.AddAttributeError(request.Path, "Invalid Required Attribute Value", required.description(request.Path.String(), target.String()))
		return
	}
	if len(required.values) > 0 {
		if slices.ContainsFunc(required.values, value.Equal) {
			return
		}
		response.Diagnostics.AddAttributeError(request.Path, "Invalid Required Attribute Value", required.description(request.Path.String(), target.String()))
	}
}

// RequireValuesWhenString requires values at targets when the validated string
// equals one of when. Targets may be relative or rooted. Empty strings do not
// satisfy the requirement; unknowns are deferred.
func RequireValuesWhenString(when []string, targets ...path.Expression) validator.String {
	return requireValuesWhenString{
		when:    append([]string(nil), when...),
		targets: append([]path.Expression(nil), targets...),
	}
}

type requireValuesWhenString struct {
	when    []string
	targets []path.Expression
}

func (required requireValuesWhenString) Description(context.Context) string {
	targets := make([]string, len(required.targets))
	for index, target := range required.targets {
		targets[index] = target.String()
	}
	return fmt.Sprintf("%s must have values when this attribute is one of %s", strings.Join(targets, ", "), strings.Join(required.when, ", "))
}

func (required requireValuesWhenString) description(attribute string, target path.Path) string {
	return fmt.Sprintf("%s must have a value when %s is one of %s", target, attribute, strings.Join(required.when, ", "))
}

func (required requireValuesWhenString) MarkdownDescription(ctx context.Context) string {
	return required.Description(ctx)
}

func (required requireValuesWhenString) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() || !slices.Contains(required.when, request.ConfigValue.ValueString()) {
		return
	}

	for _, targetExpression := range required.targets {
		target, ok := configTargetPath(ctx, request.Config, request.Path, targetExpression, "Invalid Validator Configuration", &response.Diagnostics)
		if !ok {
			continue
		}
		var value attr.Value
		diagnostics := request.Config.GetAttribute(ctx, target, &value)
		response.Diagnostics.Append(diagnostics...)
		if diagnostics.HasError() || value.IsUnknown() {
			continue
		}
		if value.IsNull() {
			response.Diagnostics.AddAttributeError(request.Path, "Invalid Required Attribute Value", required.description(request.Path.String(), target))
			continue
		}
		if stringValue, ok := value.(types.String); ok && stringValue.ValueString() == "" {
			response.Diagnostics.AddAttributeError(request.Path, "Invalid Required Attribute Value", required.description(request.Path.String(), target))
		}
	}
}

// AllowAttributesOnlyWhenValue rejects named object settings unless the
// trigger attribute equals one of allowed. An absent trigger does not allow
// the settings, while an unknown trigger defers validation.
func AllowAttributesOnlyWhenValue(trigger string, allowed []attr.Value, names ...string) validator.Object {
	return allowAttributesOnlyWhenValue{
		trigger: trigger,
		allowed: append([]attr.Value(nil), allowed...),
		names:   append([]string(nil), names...),
	}
}

type allowAttributesOnlyWhenValue struct {
	trigger string
	allowed []attr.Value
	names   []string
}

func (allowed allowAttributesOnlyWhenValue) Description(context.Context) string {
	values := make([]string, len(allowed.allowed))
	for index, value := range allowed.allowed {
		values[index] = value.String()
	}
	return fmt.Sprintf("%s may be set only when %s is one of %s", strings.Join(allowed.names, ", "), allowed.trigger, strings.Join(values, ", "))
}

func (allowed allowAttributesOnlyWhenValue) MarkdownDescription(ctx context.Context) string {
	return allowed.Description(ctx)
}

func (allowed allowAttributesOnlyWhenValue) ValidateObject(ctx context.Context, request validator.ObjectRequest, response *validator.ObjectResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	attributes := request.ConfigValue.Attributes()
	trigger, exists := attributes[allowed.trigger]
	if exists && (trigger.IsUnknown() || slices.ContainsFunc(allowed.allowed, trigger.Equal)) {
		return
	}
	for _, name := range allowed.names {
		value, exists := attributes[name]
		if exists && !value.IsNull() && !value.IsUnknown() {
			response.Diagnostics.AddAttributeError(request.Path.AtName(name), "Invalid Conditional Attribute Value", allowed.Description(ctx))
		}
	}
}

// RequireAttributesWhenEnabled requires named settings when enabled is true
// and rejects settings when enabled is false or omitted. An omitted enabled
// value is treated as false, matching a default-disabled object schema.
func RequireAttributesWhenEnabled(enabledAttribute string, names ...string) validator.Object {
	return requireAttributesWhenEnabled{enabledAttribute: enabledAttribute, names: append([]string(nil), names...)}
}

type requireAttributesWhenEnabled struct {
	enabledAttribute string
	names            []string
}

func (required requireAttributesWhenEnabled) Description(context.Context) string {
	return "Requires " + strings.Join(required.names, ", ") + " when " + required.enabledAttribute + " is true; settings require it to be true."
}

func (required requireAttributesWhenEnabled) MarkdownDescription(ctx context.Context) string {
	return required.Description(ctx)
}

func (required requireAttributesWhenEnabled) ValidateObject(_ context.Context, request validator.ObjectRequest, response *validator.ObjectResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	attributes := request.ConfigValue.Attributes()
	enabled, ok := attributes[required.enabledAttribute].(types.Bool)
	if !ok || enabled.IsUnknown() {
		return
	}
	if !enabled.ValueBool() {
		for name, value := range attributes {
			if name != required.enabledAttribute && !value.IsNull() && !value.IsUnknown() {
				response.Diagnostics.AddAttributeError(
					request.Path.AtName(name),
					"Invalid Disabled Object Attribute",
					name+" requires "+required.enabledAttribute+" = true.",
				)
			}
		}
		return
	}
	for _, name := range required.names {
		value, exists := attributes[name]
		if !exists || value.IsNull() {
			response.Diagnostics.AddAttributeError(request.Path.AtName(name), "Missing required object attribute", name+" must be configured when "+required.enabledAttribute+" is true.")
			continue
		}
		if value.IsUnknown() {
			continue
		}
		if stringValue, ok := value.(types.String); ok && stringValue.ValueString() == "" {
			response.Diagnostics.AddAttributeError(request.Path.AtName(name), "Missing required object attribute", name+" must be configured when "+required.enabledAttribute+" is true.")
		}
	}
}
