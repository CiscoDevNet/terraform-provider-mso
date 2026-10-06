package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/CiscoDevNet/terraform-provider-mso/internal/models"
	ndoapi "github.com/CiscoDevNet/terraform-provider-mso/internal/ndoapi"
	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &L3OutResource{}
var _ resource.ResourceWithIdentity = &L3OutResource{}
var _ resource.ResourceWithImportState = &L3OutResource{}

func NewL3OutResource() resource.Resource {
	return &L3OutResource{}
}

func init() {
	registerResource("mso_l3out", NewL3OutResource)
}

// L3OutResource manages the aggregate L3Out configuration from
// ndo_l3out_template. Omitted optional attributes and blocks are treated as
// unmanaged so child resources can own those parts independently.
type L3OutResource struct {
	client *client.Client
}

func (r *L3OutResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_l3out"
}

func (r *L3OutResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, response *resource.IdentitySchemaResponse) {
	response.IdentitySchema = models.L3OutResourceIdentitySchema()
}

func (r *L3OutResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = models.L3OutResourceSchema()
}

func (r *L3OutResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	tflog.Debug(ctx, "MSO L3Out Resource: Beginning Configure")
	if request.ProviderData == nil {
		return
	}

	configuredClient, ok := request.ProviderData.(*client.Client)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)
		return
	}

	r.client = configuredClient
	tflog.Debug(ctx, "MSO L3Out Resource: Configure Completed")
}

func (r *L3OutResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	tflog.Debug(ctx, "MSO L3Out Resource: Beginning Create")
	var plan models.L3OutModel
	var configuration models.L3OutModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.Config.Get(ctx, &configuration)...)
	if response.Diagnostics.HasError() {
		return
	}
	appendPath, found, err := plan.Path().AppendPath()
	if err != nil {
		response.Diagnostics.AddError("Failed to Create L3Out", err.Error())
		return
	}
	if !found {
		response.Diagnostics.AddError("Failed to Create L3Out", "Unable to resolve the L3Out collection path.")
		return
	}

	payload := plan.ToPayload(ctx, configuration, &response.Diagnostics)
	if response.Diagnostics.HasError() {
		return
	}
	_, err = r.client.PatchbyID(ndoapi.TemplateEndpoint(plan.TemplateID.ValueString()), ndoapi.NewPatchPayload(
		"add",
		appendPath,
		payload,
	))
	if err != nil {
		response.Diagnostics.AddError("Failed to Create L3Out", err.Error())
		return
	}

	r.refreshL3OutState(ctx, plan.TemplateID.ValueString(), plan.UUID.ValueString(), plan.Name.ValueString(), plan, false, &response.State, response.Identity, &response.Diagnostics)
	if !response.Diagnostics.HasError() {
		tflog.Debug(ctx, "MSO L3Out Resource: Create Completed", map[string]interface{}{"name": plan.Name.ValueString()})
	}
}

func (r *L3OutResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	tflog.Debug(ctx, "MSO L3Out Resource: Beginning Read")
	var state models.L3OutModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	r.refreshL3OutState(ctx, state.TemplateID.ValueString(), state.UUID.ValueString(), state.Name.ValueString(), state, true, &response.State, response.Identity, &response.Diagnostics)
	if !response.Diagnostics.HasError() {
		tflog.Debug(ctx, "MSO L3Out Resource: Read Completed", map[string]interface{}{"name": state.Name.ValueString()})
	}
}

func (r *L3OutResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	tflog.Debug(ctx, "MSO L3Out Resource: Beginning Update")
	var plan models.L3OutModel
	var state models.L3OutModel
	var configuration models.L3OutModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	response.Diagnostics.Append(request.Config.Get(ctx, &configuration)...)
	if response.Diagnostics.HasError() {
		return
	}
	_, resolved, found, err := getL3OutResourceModel(
		ctx,
		r.client,
		plan.TemplateID.ValueString(),
		state.Path(),
		state,
	)
	if err != nil {
		// Terraform normally refreshes the resource before Update. This is a
		// defensive error path for a template read or decode failure during the update.
		response.Diagnostics.AddError(l3OutDiagnosticTitle(err, ""), err.Error())
		return
	}
	if !found {
		// Only triggers when the L3Out is deleted after Terraform's read and before this update.
		// A subsequent Terraform plan should detect the missing object and schedule its recreation.
		response.Diagnostics.AddError("L3Out Not Found", fmt.Sprintf("L3Out %q was not found in template %q.", state.Name.ValueString(), plan.TemplateID.ValueString()))
		return
	}

	operations := models.BuildPatchOperations(ctx, resolved, plan, configuration, state, &response.Diagnostics)
	if response.Diagnostics.HasError() {
		return
	}
	if err = operations.Apply(r.client, ndoapi.TemplateEndpoint(plan.TemplateID.ValueString())); err != nil {
		response.Diagnostics.AddError("Failed to Update L3Out", err.Error())
		return
	}

	r.refreshL3OutState(ctx, plan.TemplateID.ValueString(), state.UUID.ValueString(), plan.Name.ValueString(), plan, false, &response.State, response.Identity, &response.Diagnostics)
	if !response.Diagnostics.HasError() {
		tflog.Debug(ctx, "MSO L3Out Resource: Update Completed", map[string]interface{}{"name": plan.Name.ValueString()})
	}
}

func (r *L3OutResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	tflog.Debug(ctx, "MSO L3Out Resource: Beginning Delete")
	var state models.L3OutModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	l3OutTemplate, err := ndoapi.GetTemplate(r.client, state.TemplateID.ValueString())
	if err != nil {
		if errors.Is(err, ndoapi.ErrNotFound) {
			response.State.RemoveResource(ctx)
		} else {
			response.Diagnostics.AddError(l3OutDiagnosticTitle(err, ""), err.Error())
		}
		return
	}

	resolved, found, err := l3OutTemplate.Find(state.Path())
	if err != nil {
		response.Diagnostics.AddError("Failed to Read L3Out Template", err.Error())
		return
	}
	if found {
		if _, err = r.client.PatchbyID(ndoapi.TemplateEndpoint(state.TemplateID.ValueString()), ndoapi.NewRemovePatchPayload(resolved.PatchPath())); err != nil {
			response.Diagnostics.AddError("Failed to Delete L3Out", err.Error())
			return
		}
	}

	response.State.RemoveResource(ctx)
	tflog.Debug(ctx, "MSO L3Out Resource: Delete Completed", map[string]interface{}{"name": state.Name.ValueString()})
}

func (r *L3OutResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	tflog.Debug(ctx, "MSO L3Out Resource: Beginning Import")
	var templateID, uuid string
	if request.ID != "" {
		parts := strings.SplitN(request.ID, "/", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			response.Diagnostics.AddError(
				"Invalid L3Out Import ID",
				"The import ID must be in the form <template_id>/<l3out_uuid>.",
			)
			return
		}
		templateID = parts[0]
		uuid = parts[1]
	} else {
		var identity models.L3OutResourceIdentityModel
		response.Diagnostics.Append(request.Identity.Get(ctx, &identity)...)
		if response.Diagnostics.HasError() {
			return
		}
		templateID = identity.TemplateID.ValueString()
		uuid = identity.UUID.ValueString()
	}

	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("template_id"), templateID)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("uuid"), uuid)...)
	identity := models.L3OutResourceIdentityModel{
		TemplateID: types.StringValue(templateID),
		UUID:       types.StringValue(uuid),
	}
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("id"), identity.ID())...)
	response.Diagnostics.Append(response.Identity.Set(ctx, identity)...)
	if !response.Diagnostics.HasError() {
		tflog.Debug(ctx, "MSO L3Out Resource: Import Completed", map[string]interface{}{"template_id": templateID, "uuid": uuid})
	}
}

// l3OutDiagnosticTitle maps template and L3Out path errors to a diagnostic
// title. Callers decide whether a missing object should be removed from state.
func l3OutDiagnosticTitle(err error, missingChildTitle string) string {
	if errors.Is(err, ndoapi.ErrNotFound) {
		return "L3Out Template Not Found"
	}
	var pathErr *ndoapi.TemplatePathNotFoundError
	if errors.As(err, &pathErr) {
		if pathErr.Field == "l3outs" {
			return "L3Out Not Found"
		}
		if missingChildTitle != "" {
			return missingChildTitle
		}
	}
	return "Failed to Read L3Out Template"
}

// getL3OutResourceModel also returns the resolved raw object and patch path
// needed to update the L3Out without another template lookup.
func getL3OutResourceModel(ctx context.Context, msoClient *client.Client, templateID string, path ndoapi.Path, previousL3OutModel models.L3OutModel) (models.L3OutModel, ndoapi.ResolvedObject, bool, error) {
	l3OutTemplate, err := ndoapi.GetTemplate(msoClient, templateID)
	if err != nil {
		return models.L3OutModel{}, ndoapi.ResolvedObject{}, false, err
	}

	resolved, found, err := l3OutTemplate.Find(path)
	if err != nil {
		return models.L3OutModel{}, ndoapi.ResolvedObject{}, false, err
	}
	if !found {
		return models.L3OutModel{}, ndoapi.ResolvedObject{}, false, nil
	}

	var data models.L3OutModel
	if err := data.SetFromNDOObject(ctx, templateID, resolved.Object, previousL3OutModel); err != nil {
		return models.L3OutModel{}, resolved, true, err
	}

	return data, resolved, true, nil
}

func (r *L3OutResource) refreshL3OutState(ctx context.Context, templateID, uuid, name string, previousL3OutModel models.L3OutModel, removeIfMissing bool, state *tfsdk.State, identity *tfsdk.ResourceIdentity, diagnostics *diag.Diagnostics) {
	data, _, found, err := getL3OutResourceModel(
		ctx,
		r.client,
		templateID,
		models.NewL3OutPath(uuid, name),
		previousL3OutModel,
	)
	if err != nil {
		if removeIfMissing && errors.Is(err, ndoapi.ErrNotFound) {
			state.RemoveResource(ctx)
			return
		}
		diagnostics.AddError(l3OutDiagnosticTitle(err, ""), err.Error())
		return
	}
	if !found {
		if removeIfMissing {
			state.RemoveResource(ctx)
		} else {
			diagnostics.AddError("L3Out Not Found After Apply", fmt.Sprintf("L3Out %q was not found in template %q after the change was applied.", name, templateID))
		}
		return
	}
	diagnostics.Append(state.Set(ctx, &data)...)
	if identity != nil {
		diagnostics.Append(identity.Set(ctx, models.L3OutResourceIdentityModel{
			TemplateID: data.TemplateID,
			UUID:       data.UUID,
		})...)
	}
}
