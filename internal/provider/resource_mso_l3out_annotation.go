package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/CiscoDevNet/terraform-provider-mso/internal/models"
	ndoapi "github.com/CiscoDevNet/terraform-provider-mso/internal/ndoapi"
	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &L3OutAnnotationResource{}
var _ resource.ResourceWithIdentity = &L3OutAnnotationResource{}
var _ resource.ResourceWithImportState = &L3OutAnnotationResource{}

func NewL3OutAnnotationResource() resource.Resource {
	return &L3OutAnnotationResource{}
}

func init() {
	registerResource("mso_l3out_annotation", NewL3OutAnnotationResource)
}

// L3OutAnnotationResource manages one tagAnnotations entry by tagKey.
type L3OutAnnotationResource struct {
	client *client.Client
}

func (r *L3OutAnnotationResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = fmt.Sprintf("%s_l3out_annotation", request.ProviderTypeName)
}

func (r *L3OutAnnotationResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, response *resource.IdentitySchemaResponse) {
	response.IdentitySchema = models.L3OutAnnotationResourceIdentitySchema()
}

func (r *L3OutAnnotationResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = models.L3OutAnnotationResourceSchema()
}

func (r *L3OutAnnotationResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	configuredClient, ok := request.ProviderData.(*client.Client)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected *client.Client, got: %T. Please report this issue to the provider developers.",
				request.ProviderData,
			),
		)
		return
	}
	r.client = configuredClient
	tflog.Debug(ctx, "MSO L3Out Annotation Resource: Configure Completed")
}

func (r *L3OutAnnotationResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan models.L3OutAnnotationModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	templateID := plan.TemplateID.ValueString()
	template, err := ndoapi.GetTemplate(r.client, templateID)
	if err != nil {
		response.Diagnostics.AddError(l3OutDiagnosticTitle(err, ""), err.Error())
		return
	}
	appendPath, err := template.AppendPath(plan.Path())
	if err != nil {
		response.Diagnostics.AddError(l3OutDiagnosticTitle(err, "Failed to Resolve L3Out Annotation Path"), err.Error())
		return
	}
	_, found, err := template.Find(plan.Path())
	if err != nil {
		response.Diagnostics.AddError("Failed to Read L3Out Template", err.Error())
		return
	}
	if found {
		response.Diagnostics.AddError(
			"L3Out Annotation Already Exists",
			fmt.Sprintf(
				"Annotation key %q already exists on L3Out %q. Import it to manage the existing annotation.",
				plan.Key.ValueString(),
				plan.L3OutUUID.ValueString(),
			),
		)
		return
	}

	patch := ndoapi.NewPatchPayload("add", appendPath, plan.AnnotationModel.ToPayload())
	patchResponse, err := r.client.PatchbyID(ndoapi.TemplateEndpoint(templateID), patch)
	if err != nil {
		response.Diagnostics.AddError("Failed to Create L3Out Annotation", err.Error())
		return
	}
	r.refreshState(ctx, plan, patchResponse, true, &response.State, response.Identity, &response.Diagnostics)
}

func (r *L3OutAnnotationResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state models.L3OutAnnotationModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	r.refreshState(ctx, state, nil, false, &response.State, response.Identity, &response.Diagnostics)
}

func (r *L3OutAnnotationResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan models.L3OutAnnotationModel
	var state models.L3OutAnnotationModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	template, err := ndoapi.GetTemplate(r.client, state.TemplateID.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Failed to Read L3Out Template", err.Error())
		return
	}
	resolved, err := template.FindRequired(state.Path())
	if err != nil {
		response.Diagnostics.AddError(l3OutDiagnosticTitle(err, "L3Out Annotation Not Found"), err.Error())
		return
	}
	operations := ndoapi.NewPatchOperations(resolved.Object, resolved.PatchPath())
	plan.AnnotationModel.AddPatchOperations(operations)
	// A concurrent change can make NDO already match the new plan. Save that
	// object without sending an empty PATCH or issuing another GET.
	if len(operations.Operations()) == 0 {
		var data models.L3OutAnnotationModel
		decodeErr := data.SetFromNDOObject(
			plan.TemplateID.ValueString(), plan.L3OutUUID.ValueString(), resolved.Object,
		)
		setL3OutAnnotationResourceState(ctx, data, decodeErr, &response.State, response.Identity, &response.Diagnostics)
		if decodeErr != nil {
			response.Diagnostics.AddError("Failed to Read L3Out Annotation After Apply", decodeErr.Error())
		}
		return
	}
	patchResponse, err := operations.Apply(r.client, ndoapi.TemplateEndpoint(plan.TemplateID.ValueString()))
	if err != nil {
		response.Diagnostics.AddError("Failed to Update L3Out Annotation", err.Error())
		return
	}
	r.refreshState(ctx, plan, patchResponse, true, &response.State, response.Identity, &response.Diagnostics)
}

func (r *L3OutAnnotationResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state models.L3OutAnnotationModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	template, err := ndoapi.GetTemplate(r.client, state.TemplateID.ValueString())
	if err != nil {
		if errors.Is(err, ndoapi.ErrNotFound) {
			response.State.RemoveResource(ctx)
		} else {
			response.Diagnostics.AddError("Failed to Read L3Out Template", err.Error())
		}
		return
	}
	resolved, found, err := template.Find(state.Path())
	if err != nil {
		response.Diagnostics.AddError("Failed to Read L3Out Template", err.Error())
		return
	}
	if found {
		if _, err := r.client.PatchbyID(
			ndoapi.TemplateEndpoint(state.TemplateID.ValueString()),
			ndoapi.NewRemovePatchPayload(resolved.PatchPath()),
		); err != nil {
			response.Diagnostics.AddError("Failed to Delete L3Out Annotation", err.Error())
			return
		}
	}
	response.State.RemoveResource(ctx)
}

func (r *L3OutAnnotationResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var identity models.L3OutAnnotationResourceIdentityModel
	if request.ID != "" {
		parts := strings.SplitN(request.ID, "/", 3)
		if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
			response.Diagnostics.AddError(
				"Invalid L3Out Annotation Import ID",
				"The import ID must be in the form <template_id>/<l3out_uuid>/<key>.",
			)
			return
		}
		identity = models.L3OutAnnotationResourceIdentityModel{
			TemplateID: types.StringValue(parts[0]),
			L3OutUUID:  types.StringValue(parts[1]),
			Key:        types.StringValue(parts[2]),
		}
	} else {
		response.Diagnostics.Append(request.Identity.Get(ctx, &identity)...)
		if response.Diagnostics.HasError() {
			return
		}
	}
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("template_id"), identity.TemplateID)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("l3out_uuid"), identity.L3OutUUID)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("key"), identity.Key)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("id"), identity.ID())...)
	response.Diagnostics.Append(response.Identity.Set(ctx, identity)...)
}

// refreshState uses a successful PATCH response after Create or Update, with
// one GET fallback if the response is missing or cannot be decoded. A missing
// object after PATCH is an error; Read removes missing objects from state.
func (r *L3OutAnnotationResource) refreshState(ctx context.Context, selector models.L3OutAnnotationModel, patchResponse *container.Container, afterPatch bool, state *tfsdk.State, identity *tfsdk.ResourceIdentity, diagnostics *diag.Diagnostics) {
	templateID := selector.TemplateID.ValueString()
	var patchError error
	if afterPatch {
		template, err := ndoapi.TemplateFromResponse(templateID, patchResponse)
		if err == nil {
			err = setL3OutAnnotationStateFromTemplate(ctx, template, selector, state, identity, diagnostics)
		}
		if diagnostics.HasError() || err == nil {
			return
		}
		patchError = err
	}

	template, err := ndoapi.GetTemplate(r.client, templateID)
	if err == nil {
		err = setL3OutAnnotationStateFromTemplate(ctx, template, selector, state, identity, diagnostics)
	}
	if diagnostics.HasError() {
		return
	}
	if err == nil {
		if afterPatch {
			tflog.Debug(
				ctx,
				"MSO L3Out Annotation Resource: Recovered PATCH response with template GET",
				map[string]interface{}{"reason": patchError.Error()},
			)
		}
		return
	}
	if !afterPatch {
		var missingPath *ndoapi.TemplatePathNotFoundError
		if errors.Is(err, ndoapi.ErrNotFound) || errors.As(err, &missingPath) {
			state.RemoveResource(ctx)
			return
		}
		diagnostics.AddError(l3OutDiagnosticTitle(err, ""), err.Error())
		return
	}
	diagnostics.AddError(
		"Failed to Refresh L3Out Annotation After Apply",
		fmt.Sprintf(
			"The annotation PATCH succeeded, but its state could not be refreshed. PATCH response: %s. GET fallback: %s.",
			patchError,
			err,
		),
	)
}

func setL3OutAnnotationStateFromTemplate(ctx context.Context, template *ndoapi.Template, selector models.L3OutAnnotationModel, state *tfsdk.State, identity *tfsdk.ResourceIdentity, diagnostics *diag.Diagnostics) error {
	resolved, err := template.FindRequired(selector.Path())
	if err != nil {
		return err
	}
	var data models.L3OutAnnotationModel
	decodeErr := data.SetFromNDOObject(
		selector.TemplateID.ValueString(), selector.L3OutUUID.ValueString(), resolved.Object,
	)
	setL3OutAnnotationResourceState(ctx, data, decodeErr, state, identity, diagnostics)
	return decodeErr
}

// setL3OutAnnotationResourceState records confirmed identity before full state,
// retaining the identity if decoding the annotation value fails.
func setL3OutAnnotationResourceState(ctx context.Context, data models.L3OutAnnotationModel, decodeErr error, state *tfsdk.State, identity *tfsdk.ResourceIdentity, diagnostics *diag.Diagnostics) {
	if data.Key.IsNull() || data.Key.IsUnknown() || data.Key.ValueString() == "" {
		return
	}
	diagnostics.Append(state.SetAttribute(ctx, path.Root("template_id"), data.TemplateID)...)
	diagnostics.Append(state.SetAttribute(ctx, path.Root("l3out_uuid"), data.L3OutUUID)...)
	diagnostics.Append(state.SetAttribute(ctx, path.Root("key"), data.Key)...)
	diagnostics.Append(state.SetAttribute(ctx, path.Root("id"), data.ID)...)
	if diagnostics.HasError() {
		return
	}
	if identity != nil {
		diagnostics.Append(identity.Set(ctx, models.L3OutAnnotationResourceIdentityModel{
			TemplateID: data.TemplateID,
			L3OutUUID:  data.L3OutUUID,
			Key:        data.Key,
		})...)
	}
	if diagnostics.HasError() || decodeErr != nil {
		return
	}
	diagnostics.Append(state.Set(ctx, &data)...)
}
