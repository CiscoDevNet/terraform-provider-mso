package models

import (
	"fmt"

	"github.com/CiscoDevNet/terraform-provider-mso/internal/ndoapi"
	"github.com/CiscoDevNet/terraform-provider-mso/internal/tfplugin"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var interfaceGroupBFDTypes = map[string]attr.Type{
	"enabled":                types.BoolType,
	"authentication_enabled": types.BoolType,
	"key_id":                 types.Int64Type,
	"key":                    types.StringType,
}

func interfaceGroupBFDResourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"enabled": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
			MarkdownDescription: "Whether BFD is enabled.",
		},
		"authentication_enabled": schema.BoolAttribute{
			Optional: true,
			Computed: true,
			Validators: []validator.Bool{
				tfplugin.RequireValueWhenBool(
					true,
					path.MatchRelative().AtParent().AtName("key"),
				),
				tfplugin.RequireValueWhenBool(
					true,
					path.MatchRelative().AtParent().AtName("key_id"),
				),
			},
			MarkdownDescription: "Whether BFD authentication is enabled.",
		},
		"key_id": schema.Int64Attribute{
			Optional:            true,
			MarkdownDescription: "BFD authentication key ID.",
		},
		"key": schema.StringAttribute{
			Optional:            true,
			Sensitive:           true,
			MarkdownDescription: "BFD authentication key.",
		},
	}
}

func interfaceGroupBFDDataSourceAttributes() map[string]datasourceschema.Attribute {
	return map[string]datasourceschema.Attribute{
		"enabled": datasourceschema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "Whether BFD is enabled.",
		},
		"authentication_enabled": datasourceschema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "Whether BFD authentication is enabled.",
		},
		"key_id": datasourceschema.Int64Attribute{
			Computed:            true,
			MarkdownDescription: "BFD authentication key ID.",
		},
	}
}

// BFD and BFD multi-hop have the same shape; field selects the NDO object.
func interfaceGroupBFDPayload(value types.Object) map[string]any {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	attributes := value.Attributes()
	enabled, _ := attributes["enabled"].(types.Bool)
	if !enabled.ValueBool() {
		return nil
	}
	payload := map[string]any{"enabled": true}
	auth, _ := attributes["authentication_enabled"].(types.Bool)
	if pointer := tfplugin.KnownBoolPointer(auth); pointer != nil {
		payload["authEnabled"] = *pointer
	}
	keyID, _ := attributes["key_id"].(types.Int64)
	if pointer := tfplugin.KnownInt64Pointer(keyID); pointer != nil {
		payload["keyID"] = *pointer
	}
	key, _ := attributes["key"].(types.String)
	if pointer := tfplugin.KnownStringPointer(key); pointer != nil {
		payload["key"] = map[string]any{"value": *pointer}
	}
	return payload
}

func interfaceGroupBFDFromNDO(object map[string]any, field string, previous types.Object) (types.Object, error) {
	protocol, exists, err := ndoapi.MapField(object, field, ndoapi.OptionalField)
	if err != nil {
		return types.ObjectNull(interfaceGroupBFDTypes), err
	}
	values := map[string]attr.Value{
		"enabled":                types.BoolValue(false),
		"authentication_enabled": types.BoolNull(),
		"key_id":                 types.Int64Null(),
		"key":                    types.StringNull(),
	}
	if exists {
		enabled, present, err := ndoapi.BoolField(protocol, "enabled", ndoapi.OptionalField)
		if err != nil {
			return types.ObjectNull(interfaceGroupBFDTypes), fmt.Errorf("NDO %s: %w", field, err)
		}
		if !present {
			enabled = true // Older responses omit enabled on configured objects.
		}
		values["enabled"] = types.BoolValue(enabled)
		if enabled {
			auth, present, err := ndoapi.BoolField(protocol, "authEnabled", ndoapi.OptionalField)
			if err != nil {
				return types.ObjectNull(interfaceGroupBFDTypes), fmt.Errorf("NDO %s: %w", field, err)
			}
			values["authentication_enabled"] = tfplugin.BoolOrNull(auth, present)
			keyID, present, err := ndoapi.Int64Field(protocol, "keyID", ndoapi.OptionalField)
			if err != nil {
				return types.ObjectNull(interfaceGroupBFDTypes), fmt.Errorf("NDO %s: %w", field, err)
			}
			if present && auth {
				values["key_id"] = types.Int64Value(keyID)
			}
			_, present, err = ndoapi.MapField(protocol, "key", ndoapi.OptionalField)
			if err != nil {
				return types.ObjectNull(interfaceGroupBFDTypes), fmt.Errorf("NDO %s: %w", field, err)
			}
			if auth && !present && !previous.IsNull() && !previous.IsUnknown() {
				if key, ok := previous.Attributes()["key"].(types.String); ok && !key.IsNull() && !key.IsUnknown() {
					return types.ObjectNull(interfaceGroupBFDTypes), fmt.Errorf("NDO %s.key reference is missing for configured authentication", field)
				}
			}
			if present && auth && !previous.IsNull() && !previous.IsUnknown() {
				if key, ok := previous.Attributes()["key"].(types.String); ok && !key.IsUnknown() {
					values["key"] = key
				}
			}
		}
	}
	result, diagnostics := types.ObjectValue(interfaceGroupBFDTypes, values)
	if diagnostics.HasError() {
		return types.ObjectNull(interfaceGroupBFDTypes), fmt.Errorf("unable to represent %s state: %s", field, diagnostics.Errors()[0].Detail())
	}
	return result, nil
}

func interfaceGroupBFDAddPatchOperation(operations *ndoapi.PatchOperations, existing map[string]any, field string, plan, configuration, state types.Object) error {
	if configuration.IsNull() {
		return nil
	}
	desired := interfaceGroupBFDPayload(plan)
	if desired == nil {
		operations.Remove(field)
		return nil
	}
	return setInterfaceGroupProtocolPatch(operations, existing, field, desired, configuration, state, func(protocol map[string]any) bool {
		authenticated, _ := protocol["authEnabled"].(bool)
		return authenticated
	})
}
