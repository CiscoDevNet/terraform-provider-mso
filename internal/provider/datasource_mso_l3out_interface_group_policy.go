package provider

import (
	"context"
	"fmt"

	"github.com/CiscoDevNet/terraform-provider-mso/internal/models"
	"github.com/CiscoDevNet/terraform-provider-mso/internal/ndoapi"
	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &L3OutInterfaceGroupPolicyDataSource{}
var _ datasource.DataSourceWithConfigure = &L3OutInterfaceGroupPolicyDataSource{}

type L3OutInterfaceGroupPolicyDataSource struct {
	client *client.Client
}

func NewL3OutInterfaceGroupPolicyDataSource() datasource.DataSource {
	return &L3OutInterfaceGroupPolicyDataSource{}
}

func init() {
	registerDataSource("mso_l3out_interface_group_policy", NewL3OutInterfaceGroupPolicyDataSource)
}

func (d *L3OutInterfaceGroupPolicyDataSource) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_l3out_interface_group_policy"
}

func (d *L3OutInterfaceGroupPolicyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = models.L3OutInterfaceGroupPolicyDataSourceSchema()
}

func (d *L3OutInterfaceGroupPolicyDataSource) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	configuredClient, ok := request.ProviderData.(*client.Client)
	if !ok {
		response.Diagnostics.AddError("Unexpected Data Source Configure Type", fmt.Sprintf("Expected *client.Client, got %T", request.ProviderData))
		return
	}
	d.client = configuredClient
}

func (d *L3OutInterfaceGroupPolicyDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var config models.L3OutInterfaceGroupPolicyModel
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	template, err := ndoapi.GetTemplate(d.client, config.TemplateID.ValueString())
	if err != nil {
		response.Diagnostics.AddError(l3OutDiagnosticTitle(err, ""), err.Error())
		return
	}
	resolved, err := template.FindRequired(config.Path())
	if err != nil {
		response.Diagnostics.AddError(l3OutDiagnosticTitle(err, "L3Out Interface Group Policy Not Found"), err.Error())
		return
	}
	var data models.L3OutInterfaceGroupPolicyModel
	if err := data.SetFromNDOObject(ctx, config.TemplateID.ValueString(), config.L3OutUUID.ValueString(), resolved.Object, models.L3OutInterfaceGroupPolicyModel{}, &response.Diagnostics); err != nil {
		response.Diagnostics.AddError("Failed to Read L3Out Template", err.Error())
		return
	}
	if response.Diagnostics.HasError() {
		return
	}
	data = data.DataSourceValue(ctx, &response.Diagnostics)
	if response.Diagnostics.HasError() {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
