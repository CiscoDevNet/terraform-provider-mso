package models

import (
	"context"
	"fmt"
	"slices"

	"github.com/CiscoDevNet/terraform-provider-mso/internal/ndoapi"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// InterfaceGroupPoliciesModel represents interfaceGroups by their unique name.
// A nil map is unmanaged; an empty non-nil map clears the collection.
type InterfaceGroupPoliciesModel map[string]InterfaceGroupPolicyModel

func interfaceGroupPolicyObjectType() attr.Type {
	return (schema.NestedAttributeObject{Attributes: interfaceGroupPolicyResourceAttributes()}).Type()
}

func interfaceGroupPolicyDataSourceObjectType() attr.Type {
	return (datasourceschema.NestedAttributeObject{Attributes: interfaceGroupPolicyDataSourceAttributes()}).Type()
}

func InterfaceGroupPoliciesFromTerraform(ctx context.Context, value types.Map, diagnostics *diag.Diagnostics) InterfaceGroupPoliciesModel {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	policies := make(InterfaceGroupPoliciesModel, len(value.Elements()))
	diagnostics.Append(value.ElementsAs(ctx, &policies, false)...)
	return policies
}

func (policies InterfaceGroupPoliciesModel) TerraformValue(ctx context.Context, diagnostics *diag.Diagnostics) types.Map {
	value, conversionDiagnostics := types.MapValueFrom(ctx, interfaceGroupPolicyObjectType(), map[string]InterfaceGroupPolicyModel(policies))
	diagnostics.Append(conversionDiagnostics...)
	return value
}

func (policies InterfaceGroupPoliciesModel) DataSourceValue(ctx context.Context, diagnostics *diag.Diagnostics) types.Map {
	result := make(map[string]InterfaceGroupPolicyModel, len(policies))
	for name, policy := range policies {
		result[name] = policy.DataSourceValue(ctx, diagnostics)
	}
	value, conversionDiagnostics := types.MapValueFrom(ctx, interfaceGroupPolicyDataSourceObjectType(), result)
	diagnostics.Append(conversionDiagnostics...)
	return value
}

func (policies *InterfaceGroupPoliciesModel) SetFromNDOObject(ctx context.Context, object map[string]any, previous InterfaceGroupPoliciesModel, diagnostics *diag.Diagnostics) error {
	items, _, err := ndoapi.ListField(object, "interfaceGroups", ndoapi.OptionalField)
	if err != nil {
		return err
	}
	result := make(InterfaceGroupPoliciesModel, len(items))
	for index, item := range items {
		group, ok := item.(map[string]any)
		if !ok {
			return fmt.Errorf("NDO interfaceGroups[%d] has unexpected type %T", index, item)
		}
		name, _, err := ndoapi.StringField(group, "name", ndoapi.RequiredNonEmptyField)
		if err != nil {
			return fmt.Errorf("NDO interfaceGroups[%d]: %w", index, err)
		}
		if _, exists := result[name]; exists {
			return fmt.Errorf("NDO interfaceGroups contains duplicate name %q", name)
		}
		var policy InterfaceGroupPolicyModel
		if err := policy.SetFromNDOObject(ctx, group, previous[name], diagnostics); err != nil {
			return fmt.Errorf("NDO interfaceGroups[%d]: %w", index, err)
		}
		if diagnostics.HasError() {
			return nil
		}
		result[name] = policy
	}
	*policies = result
	return nil
}

func (policies InterfaceGroupPoliciesModel) ToPayload(ctx context.Context, configuration InterfaceGroupPoliciesModel, diagnostics *diag.Diagnostics) []map[string]any {
	if policies == nil {
		return nil
	}
	names := make([]string, 0, len(policies))
	for name := range policies {
		names = append(names, name)
	}
	slices.Sort(names)
	payload := make([]map[string]any, 0, len(names))
	for _, name := range names {
		group := policies[name].ToPayload(ctx, configuration[name], diagnostics)
		group["name"] = name
		payload = append(payload, group)
	}
	return payload
}

func (policies InterfaceGroupPoliciesModel) AddPatchOperations(ctx context.Context, existing map[string]any, operations *ndoapi.PatchOperations, parentPath string, configuration, state InterfaceGroupPoliciesModel, diagnostics *diag.Diagnostics) error {
	if policies == nil {
		return nil
	}
	items, exists, err := ndoapi.ListField(existing, "interfaceGroups", ndoapi.OptionalField)
	if err != nil {
		return err
	}
	if !exists {
		if len(policies) > 0 {
			operations.Set("interfaceGroups", policies.ToPayload(ctx, configuration, diagnostics))
		}
		return nil
	}
	current := make(map[string]map[string]any, len(items))
	indexes := make(map[string]int, len(items))
	for index, item := range items {
		group, ok := item.(map[string]any)
		if !ok {
			return fmt.Errorf("NDO interfaceGroups[%d] has unexpected type %T", index, item)
		}
		name, _, err := ndoapi.StringField(group, "name", ndoapi.RequiredNonEmptyField)
		if err != nil {
			return fmt.Errorf("NDO interfaceGroups[%d]: %w", index, err)
		}
		if _, duplicate := current[name]; duplicate {
			return fmt.Errorf("NDO interfaceGroups contains duplicate name %q", name)
		}
		current[name] = group
		indexes[name] = index
	}

	// Update existing groups while their indexes still match the NDO response.
	names := make([]string, 0, len(policies))
	for name := range policies {
		names = append(names, name)
	}
	slices.Sort(names)
	for _, name := range names {
		group, found := current[name]
		if !found {
			continue
		}
		groupPath := fmt.Sprintf("%s/interfaceGroups/%d", parentPath, indexes[name])
		groupOperations := ndoapi.NewPatchOperations(group, groupPath)
		if err := policies[name].AddPatchOperations(ctx, group, groupOperations, configuration[name], state[name], diagnostics); err != nil {
			return fmt.Errorf("NDO interface group %q: %w", name, err)
		}
		operations.Append(groupOperations.Operations()...)
	}
	// Append new groups without changing existing indexes.
	for _, name := range names {
		if _, found := current[name]; found {
			continue
		}
		group := policies[name].ToPayload(ctx, configuration[name], diagnostics)
		group["name"] = name
		operations.Append(ndoapi.NewPatchPayload("add", parentPath+"/interfaceGroups/-", group))
	}
	// Remove existing groups last, from highest index to lowest.
	for index := len(items) - 1; index >= 0; index-- {
		group := items[index].(map[string]any)
		name, _, _ := ndoapi.StringField(group, "name", ndoapi.RequiredNonEmptyField)
		if _, configured := policies[name]; !configured {
			operations.Append(ndoapi.NewRemovePatchPayload(fmt.Sprintf("%s/interfaceGroups/%d", parentPath, index)))
		}
	}
	return nil
}

func interfaceGroupPoliciesResourceAttribute() schema.MapNestedAttribute {
	return schema.MapNestedAttribute{
		Optional: true,
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: interfaceGroupPolicyResourceAttributes(),
		},
		PlanModifiers: []planmodifier.Map{
			mapplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.Map{
			mapvalidator.KeysAre(stringvalidator.LengthAtLeast(1)),
			mapvalidator.NoNullValues(),
		},
		MarkdownDescription: "Interface group policies on the L3Out, keyed by name.",
	}
}

func interfaceGroupPoliciesDataSourceAttribute() datasourceschema.MapNestedAttribute {
	return datasourceschema.MapNestedAttribute{
		Computed: true,
		NestedObject: datasourceschema.NestedAttributeObject{
			Attributes: interfaceGroupPolicyDataSourceAttributes(),
		},
		MarkdownDescription: "Interface group policies keyed by name on the L3Out.",
	}
}
