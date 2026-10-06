package models

import (
	"fmt"

	"github.com/CiscoDevNet/terraform-provider-mso/internal/ndoapi"
	"github.com/CiscoDevNet/terraform-provider-mso/internal/tfplugin"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var interfaceGroupOSPFTypes = map[string]attr.Type{
	"enabled":             types.BoolType,
	"authentication_type": types.StringType,
	"key_id":              types.Int64Type,
	"key":                 types.StringType,
}

func interfaceGroupOSPFResourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"enabled": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
			MarkdownDescription: "Whether OSPF authentication settings are enabled.",
		},
		"authentication_type": schema.StringAttribute{
			Optional: true,
			Computed: true,
			Validators: []validator.String{
				stringvalidator.OneOf("none", "simple", "md5"),
				tfplugin.RequireValuesWhenString(
					[]string{"simple", "md5"},
					path.MatchRelative().AtParent().AtName("key"),
					path.MatchRelative().AtParent().AtName("key_id"),
				),
			},
			MarkdownDescription: "OSPF authentication mode.",
		},
		"key_id": schema.Int64Attribute{
			Optional:            true,
			MarkdownDescription: "OSPF authentication key ID.",
		},
		"key": schema.StringAttribute{
			Optional:            true,
			Sensitive:           true,
			MarkdownDescription: "OSPF authentication key.",
		},
	}
}

func interfaceGroupOSPFDataSourceAttributes() map[string]datasourceschema.Attribute {
	return map[string]datasourceschema.Attribute{
		"enabled": datasourceschema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "Whether OSPF authentication settings are enabled.",
		},
		"authentication_type": datasourceschema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "OSPF authentication mode.",
		},
		"key_id": datasourceschema.Int64Attribute{
			Computed:            true,
			MarkdownDescription: "OSPF authentication key ID.",
		},
	}
}

func interfaceGroupOSPFPayload(value types.Object) map[string]any {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	attributes := value.Attributes()
	enabled, _ := attributes["enabled"].(types.Bool)
	if !enabled.ValueBool() {
		return nil
	}
	payload := map[string]any{"enabled": true}
	mode, _ := attributes["authentication_type"].(types.String)
	if pointer := tfplugin.KnownStringPointer(mode); pointer != nil {
		payload["authType"] = *pointer
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

func interfaceGroupOSPFFromNDO(object map[string]any, previous types.Object) (types.Object, error) {
	protocol, exists, err := ndoapi.MapField(object, "ospf", ndoapi.OptionalField)
	if err != nil {
		return types.ObjectNull(interfaceGroupOSPFTypes), err
	}
	values := map[string]attr.Value{
		"enabled":             types.BoolValue(false),
		"authentication_type": types.StringNull(),
		"key_id":              types.Int64Null(),
		"key":                 types.StringNull(),
	}
	if exists {
		enabled, present, err := ndoapi.BoolField(protocol, "enabled", ndoapi.OptionalField)
		if err != nil {
			return types.ObjectNull(interfaceGroupOSPFTypes), fmt.Errorf("NDO ospf: %w", err)
		}
		if !present {
			enabled = true // Older responses omit enabled on configured objects.
		}
		values["enabled"] = types.BoolValue(enabled)
		if enabled {
			mode, present, err := ndoapi.StringField(protocol, "authType", ndoapi.OptionalField)
			if err != nil {
				return types.ObjectNull(interfaceGroupOSPFTypes), fmt.Errorf("NDO ospf: %w", err)
			}
			values["authentication_type"] = tfplugin.StringOrNull(mode, present, nil)
			keyID, present, err := ndoapi.Int64Field(protocol, "keyID", ndoapi.OptionalField)
			if err != nil {
				return types.ObjectNull(interfaceGroupOSPFTypes), fmt.Errorf("NDO ospf: %w", err)
			}
			authenticated := mode == "md5" || mode == "simple"
			if present && authenticated {
				values["key_id"] = types.Int64Value(keyID)
			}
			_, present, err = ndoapi.MapField(protocol, "key", ndoapi.OptionalField)
			if err != nil {
				return types.ObjectNull(interfaceGroupOSPFTypes), fmt.Errorf("NDO ospf: %w", err)
			}
			if authenticated && !present && !previous.IsNull() && !previous.IsUnknown() {
				if key, ok := previous.Attributes()["key"].(types.String); ok && !key.IsNull() && !key.IsUnknown() {
					return types.ObjectNull(interfaceGroupOSPFTypes), fmt.Errorf("NDO ospf.key reference is missing for configured authentication")
				}
			}
			if present && authenticated && !previous.IsNull() && !previous.IsUnknown() {
				if key, ok := previous.Attributes()["key"].(types.String); ok && !key.IsUnknown() {
					values["key"] = key
				}
			}
		}
	}
	result, diagnostics := types.ObjectValue(interfaceGroupOSPFTypes, values)
	if diagnostics.HasError() {
		return types.ObjectNull(interfaceGroupOSPFTypes), fmt.Errorf("unable to represent ospf state: %s", diagnostics.Errors()[0].Detail())
	}
	return result, nil
}

func interfaceGroupOSPFAddPatchOperation(operations *ndoapi.PatchOperations, existing map[string]any, plan, configuration, state types.Object) error {
	if configuration.IsNull() {
		return nil
	}
	desired := interfaceGroupOSPFPayload(plan)
	if desired == nil {
		operations.Remove("ospf")
		return nil
	}
	return setInterfaceGroupProtocolPatch(operations, existing, "ospf", desired, configuration, state, func(protocol map[string]any) bool {
		return protocol["authType"] == "md5" || protocol["authType"] == "simple"
	})
}
