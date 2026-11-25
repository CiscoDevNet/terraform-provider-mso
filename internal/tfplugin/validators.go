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
	return fmt.Sprintf(
		"%s must be one of %s when %s is %t",
		targetName,
		strings.Join(values, ", "),
		attribute,
		required.when,
	)
}

func (required requireValueWhenBool) MarkdownDescription(ctx context.Context) string {
	return required.Description(ctx)
}

func (required requireValueWhenBool) ValidateBool(ctx context.Context, request validator.BoolRequest, response *validator.BoolResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() ||
		request.ConfigValue.ValueBool() != required.when {
		return
	}

	target, ok := configTargetPath(
		ctx,
		request.Config,
		request.Path,
		required.target,
		"Invalid Validator Configuration",
		&response.Diagnostics,
	)
	if !ok {
		return
	}
	var value attr.Value
	response.Diagnostics.Append(request.Config.GetAttribute(ctx, target, &value)...)
	if response.Diagnostics.HasError() || value.IsUnknown() {
		return
	}
	if value.IsNull() {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid Required Attribute Value",
			required.description(request.Path.String(), target.String()),
		)
		return
	}
	if stringValue, ok := value.(types.String); ok && stringValue.ValueString() == "" {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid Required Attribute Value",
			required.description(request.Path.String(), target.String()),
		)
		return
	}
	if len(required.values) > 0 {
		if slices.ContainsFunc(required.values, value.Equal) {
			return
		}
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid Required Attribute Value",
			required.description(request.Path.String(), target.String()),
		)
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
	return fmt.Sprintf(
		"%s must have values when this attribute is one of %s",
		strings.Join(targets, ", "),
		strings.Join(required.when, ", "),
	)
}

func (required requireValuesWhenString) description(attribute string, target path.Path) string {
	return fmt.Sprintf(
		"%s must have a value when %s is one of %s",
		target,
		attribute,
		strings.Join(required.when, ", "),
	)
}

func (required requireValuesWhenString) MarkdownDescription(ctx context.Context) string {
	return required.Description(ctx)
}

func (required requireValuesWhenString) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() ||
		!slices.Contains(required.when, request.ConfigValue.ValueString()) {
		return
	}

	for _, targetExpression := range required.targets {
		target, ok := configTargetPath(
			ctx,
			request.Config,
			request.Path,
			targetExpression,
			"Invalid Validator Configuration",
			&response.Diagnostics,
		)
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
			response.Diagnostics.AddAttributeError(
				request.Path,
				"Invalid Required Attribute Value",
				required.description(request.Path.String(), target),
			)
			continue
		}
		if stringValue, ok := value.(types.String); ok && stringValue.ValueString() == "" {
			response.Diagnostics.AddAttributeError(
				request.Path,
				"Invalid Required Attribute Value",
				required.description(request.Path.String(), target),
			)
		}
	}
}

// AllowAttributesOnlyWhenValue rejects named object settings unless the
// triggerAttribute equals one of allowedTriggerValues. Values of attributeNames
// listed in attributeValueExceptions are permitted without a trigger match.
// An unknown trigger defers validation. Pass nil to allow no exceptions.
func AllowAttributesOnlyWhenValue(triggerAttribute string, allowedTriggerValues []attr.Value, attributeValueExceptions []attr.Value, attributeNames ...string) validator.Object {
	return allowAttributesOnlyWhenValue{
		triggerAttribute:         triggerAttribute,
		allowedTriggerValues:     append([]attr.Value(nil), allowedTriggerValues...),
		attributeValueExceptions: append([]attr.Value(nil), attributeValueExceptions...),
		attributeNames:           append([]string(nil), attributeNames...),
	}
}

type allowAttributesOnlyWhenValue struct {
	triggerAttribute         string
	allowedTriggerValues     []attr.Value
	attributeValueExceptions []attr.Value
	attributeNames           []string
}

func (rule allowAttributesOnlyWhenValue) Description(context.Context) string {
	return "Validates conditional relationships between object attributes."
}

func (rule allowAttributesOnlyWhenValue) MarkdownDescription(ctx context.Context) string {
	return rule.Description(ctx)
}

func (rule allowAttributesOnlyWhenValue) ValidateObject(_ context.Context, request validator.ObjectRequest, response *validator.ObjectResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	attributes := request.ConfigValue.Attributes()
	trigger, exists := attributes[rule.triggerAttribute]
	if exists && (trigger.IsUnknown() || slices.ContainsFunc(rule.allowedTriggerValues, trigger.Equal)) {
		return
	}

	names := make([]string, len(rule.attributeNames))
	for index, name := range rule.attributeNames {
		names[index] = request.Path.AtName(name).String()
	}

	values := make([]string, len(rule.allowedTriggerValues))
	for index, value := range rule.allowedTriggerValues {
		values[index] = value.String()
	}
	message := fmt.Sprintf(
		"%s may be set only when %s is one of %s",
		strings.Join(names, ", "),
		request.Path.AtName(rule.triggerAttribute).String(),
		strings.Join(values, ", "),
	)

	if len(rule.attributeValueExceptions) > 0 {
		exceptions := make([]string, len(rule.attributeValueExceptions))
		for index, value := range rule.attributeValueExceptions {
			exceptions[index] = value.String()
		}
		message = fmt.Sprintf(
			"%s; attribute values %s are exempt from this condition",
			message,
			strings.Join(exceptions, ", "),
		)
	}

	for _, name := range rule.attributeNames {
		value, exists := attributes[name]
		if !exists || value.IsNull() || value.IsUnknown() ||
			slices.ContainsFunc(rule.attributeValueExceptions, value.Equal) {
			continue
		}
		response.Diagnostics.AddAttributeError(
			request.Path.AtName(name),
			"Invalid Conditional Attribute Value",
			message,
		)
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
	return fmt.Sprintf(
		"Requires %s when %s is true; settings require it to be true.",
		strings.Join(required.names, ", "),
		required.enabledAttribute,
	)
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
	enabledPath := request.Path.AtName(required.enabledAttribute)
	if !enabled.ValueBool() {
		for name, value := range attributes {
			if name != required.enabledAttribute && !value.IsNull() && !value.IsUnknown() {
				attributePath := request.Path.AtName(name)
				response.Diagnostics.AddAttributeError(
					attributePath,
					"Invalid Disabled Object Attribute",
					fmt.Sprintf("%s requires %s = true.", attributePath, enabledPath),
				)
			}
		}
		return
	}
	for _, name := range required.names {
		attributePath := request.Path.AtName(name)
		value, exists := attributes[name]
		if !exists || value.IsNull() {
			response.Diagnostics.AddAttributeError(
				attributePath,
				"Missing required object attribute",
				fmt.Sprintf("%s must be configured when %s is true.", attributePath, enabledPath),
			)
			continue
		}
		if value.IsUnknown() {
			continue
		}
		if stringValue, ok := value.(types.String); ok && stringValue.ValueString() == "" {
			response.Diagnostics.AddAttributeError(
				attributePath,
				"Missing required object attribute",
				fmt.Sprintf("%s must be configured when %s is true.", attributePath, enabledPath),
			)
		}
	}
}
