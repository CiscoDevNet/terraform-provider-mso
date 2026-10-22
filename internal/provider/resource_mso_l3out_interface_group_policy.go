package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/CiscoDevNet/terraform-provider-mso/internal/models"
	"github.com/CiscoDevNet/terraform-provider-mso/internal/ndoapi"
	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &L3OutInterfaceGroupPolicyResource{}
var _ resource.ResourceWithIdentity = &L3OutInterfaceGroupPolicyResource{}
var _ resource.ResourceWithImportState = &L3OutInterfaceGroupPolicyResource{}

// Collection indexes must be resolved immediately before a mutation.
var l3OutInterfaceGroupPolicyMutationMu sync.Mutex

type L3OutInterfaceGroupPolicyResource struct {
	client *client.Client
}

func NewL3OutInterfaceGroupPolicyResource() resource.Resource {
	return &L3OutInterfaceGroupPolicyResource{}
}

func init() {
	registerResource("mso_l3out_interface_group_policy", NewL3OutInterfaceGroupPolicyResource)
}

func (r *L3OutInterfaceGroupPolicyResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_l3out_interface_group_policy"
}

func (r *L3OutInterfaceGroupPolicyResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, response *resource.IdentitySchemaResponse) {
	response.IdentitySchema = models.L3OutInterfaceGroupPolicyResourceIdentitySchema()
}

func (r *L3OutInterfaceGroupPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = models.L3OutInterfaceGroupPolicyResourceSchema()
}

func (r *L3OutInterfaceGroupPolicyResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	configuredClient, ok := request.ProviderData.(*client.Client)
	if !ok {
		response.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *client.Client, got %T", request.ProviderData))
		return
	}
	r.client = configuredClient
}

func (r *L3OutInterfaceGroupPolicyResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan, configuration models.L3OutInterfaceGroupPolicyModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.Config.Get(ctx, &configuration)...)
	if response.Diagnostics.HasError() {
		return
	}
	l3OutInterfaceGroupPolicyMutationMu.Lock()
	defer l3OutInterfaceGroupPolicyMutationMu.Unlock()
	template, err := ndoapi.GetTemplate(r.client, plan.TemplateID.ValueString())
	if err != nil {
		response.Diagnostics.AddError(l3OutDiagnosticTitle(err, ""), err.Error())
		return
	}
	appendPath, err := template.AppendPath(plan.Path())
	if err != nil {
		response.Diagnostics.AddError(l3OutDiagnosticTitle(err, "Failed to Resolve Interface Group Policy Path"), err.Error())
		return
	}
	_, found, err := template.Find(plan.Path())
	if err != nil {
		response.Diagnostics.AddError("Failed to Read L3Out Template", err.Error())
		return
	}
	if found {
		response.Diagnostics.AddError("L3Out Interface Group Policy Already Exists", fmt.Sprintf("Interface group policy %q already exists. Import it to manage the existing policy.", plan.Name.ValueString()))
		return
	}
	payload := plan.ToPayload(ctx, configuration, &response.Diagnostics)
	if response.Diagnostics.HasError() {
		return
	}
	if _, err := r.client.PatchbyID(ndoapi.TemplateEndpoint(plan.TemplateID.ValueString()), ndoapi.NewPatchPayload("add", appendPath, payload)); err != nil {
		response.Diagnostics.AddError("Failed to Create L3Out Interface Group Policy", err.Error())
		return
	}
	r.refreshState(ctx, plan, false, &response.State, response.Identity, &response.Diagnostics)
}

func (r *L3OutInterfaceGroupPolicyResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state models.L3OutInterfaceGroupPolicyModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	r.refreshState(ctx, state, true, &response.State, response.Identity, &response.Diagnostics)
}

func (r *L3OutInterfaceGroupPolicyResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan, state, configuration models.L3OutInterfaceGroupPolicyModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	response.Diagnostics.Append(request.Config.Get(ctx, &configuration)...)
	if response.Diagnostics.HasError() {
		return
	}
	l3OutInterfaceGroupPolicyMutationMu.Lock()
	defer l3OutInterfaceGroupPolicyMutationMu.Unlock()
	template, err := ndoapi.GetTemplate(r.client, state.TemplateID.ValueString())
	if err != nil {
		response.Diagnostics.AddError(l3OutDiagnosticTitle(err, ""), err.Error())
		return
	}
	resolved, err := template.FindRequired(state.Path())
	if err != nil {
		response.Diagnostics.AddError(l3OutDiagnosticTitle(err, "L3Out Interface Group Policy Not Found"), err.Error())
		return
	}
	operations := ndoapi.NewPatchOperations(resolved.Object, resolved.PatchPath())
	if err := plan.AddPatchOperations(ctx, resolved.Object, operations, configuration, state, &response.Diagnostics); err != nil {
		response.Diagnostics.AddError("Failed to Read L3Out Interface Group Policy", err.Error())
		return
	}
	if response.Diagnostics.HasError() {
		return
	}
	if err := operations.Apply(r.client, ndoapi.TemplateEndpoint(plan.TemplateID.ValueString())); err != nil {
		response.Diagnostics.AddError("Failed to Update L3Out Interface Group Policy", err.Error())
		return
	}
	r.refreshState(ctx, plan, false, &response.State, response.Identity, &response.Diagnostics)
}

func (r *L3OutInterfaceGroupPolicyResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state models.L3OutInterfaceGroupPolicyModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	l3OutInterfaceGroupPolicyMutationMu.Lock()
	defer l3OutInterfaceGroupPolicyMutationMu.Unlock()
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
		if _, err := r.client.PatchbyID(ndoapi.TemplateEndpoint(state.TemplateID.ValueString()), ndoapi.NewRemovePatchPayload(resolved.PatchPath())); err != nil {
			response.Diagnostics.AddError("Failed to Delete L3Out Interface Group Policy", err.Error())
			return
		}
	}
	response.State.RemoveResource(ctx)
}

func (r *L3OutInterfaceGroupPolicyResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var identity models.L3OutInterfaceGroupPolicyResourceIdentityModel
	if request.ID != "" {
		parts := strings.SplitN(request.ID, "/", 3)
		if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
			response.Diagnostics.AddError("Invalid L3Out Interface Group Policy Import ID", "The import ID must be in the form <template_id>/<l3out_uuid>/<name>.")
			return
		}
		identity = models.L3OutInterfaceGroupPolicyResourceIdentityModel{TemplateID: types.StringValue(parts[0]), L3OutUUID: types.StringValue(parts[1]), Name: types.StringValue(parts[2])}
	} else {
		response.Diagnostics.Append(request.Identity.Get(ctx, &identity)...)
		if response.Diagnostics.HasError() {
			return
		}
	}
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("template_id"), identity.TemplateID)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("l3out_uuid"), identity.L3OutUUID)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("name"), identity.Name)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("id"), identity.ID())...)
	response.Diagnostics.Append(response.Identity.Set(ctx, identity)...)
}

func (r *L3OutInterfaceGroupPolicyResource) refreshState(ctx context.Context, previous models.L3OutInterfaceGroupPolicyModel, removeIfMissing bool, state *tfsdk.State, identity *tfsdk.ResourceIdentity, diagnostics *diag.Diagnostics) {
	template, err := ndoapi.GetTemplate(r.client, previous.TemplateID.ValueString())
	if err != nil {
		if removeIfMissing && errors.Is(err, ndoapi.ErrNotFound) {
			state.RemoveResource(ctx)
		} else {
			diagnostics.AddError(l3OutDiagnosticTitle(err, ""), err.Error())
		}
		return
	}
	resolved, found, err := template.Find(previous.Path())
	if err != nil {
		diagnostics.AddError("Failed to Read L3Out Template", err.Error())
		return
	}
	if !found {
		if removeIfMissing {
			state.RemoveResource(ctx)
		} else {
			diagnostics.AddError("L3Out Interface Group Policy Not Found After Apply", fmt.Sprintf("Interface group policy %q was not found after the change was applied.", previous.Name.ValueString()))
		}
		return
	}
	var data models.L3OutInterfaceGroupPolicyModel
	if err := data.SetFromNDOObject(ctx, previous.TemplateID.ValueString(), previous.L3OutUUID.ValueString(), resolved.Object, previous, diagnostics); err != nil {
		diagnostics.AddError("Failed to Read L3Out Template", err.Error())
		return
	}
	if diagnostics.HasError() {
		return
	}
	diagnostics.Append(state.Set(ctx, &data)...)
	if identity != nil {
		diagnostics.Append(identity.Set(ctx, models.L3OutInterfaceGroupPolicyResourceIdentityModel{TemplateID: data.TemplateID, L3OutUUID: data.L3OutUUID, Name: data.Name})...)
	}
}
