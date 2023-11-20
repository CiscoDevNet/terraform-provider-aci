package provider

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AciRestManagedDataSource{}

func NewAciRestManagedDataSource() datasource.DataSource {
	return &AciRestManagedDataSource{}
}

// AciRestManagedDataSource defines the data source implementation.
type AciRestManagedDataSource struct {
	client *client.Client
}

func (d *AciRestManagedDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	tflog.Debug(ctx, "Start schema of datasource: aci_rest_managed")
	resp.TypeName = req.ProviderTypeName + "_rest_managed"
	tflog.Debug(ctx, "End schema of datasource: aci_rest_managed")
}

func (d *AciRestManagedDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "The _rest_managed datasource for the 'AciRestManaged' class",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The distinquised name (DN) of the object.",
			},
			"dn": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The distinquised name (DN) of the parent object. e.g. uni/tn-EXAMPLE_TENANT",
			},
			"class_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Which class object is being created. (Make sure there is no colon in the classname)",
			},
			"content": schema.MapAttribute{
				MarkdownDescription: "Map of key-value pairs those needed to be passed to the Model object as parameters. Make sure the key name matches the name with the object parameter in ACI.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"annotation": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: `The annotation of the ACI object.`,
			},
		},
		Blocks: map[string]schema.Block{
			"child": schema.SetNestedBlock{
				//Optional:            true,
				MarkdownDescription: "List of children.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"rn": schema.StringAttribute{
							MarkdownDescription: "The relative name of the child object.",
							Required:            true,
						},
						"class_name": schema.StringAttribute{
							MarkdownDescription: "Class name of child object.",
							Optional:            true,
							Computed:            true,
						},
						"content": schema.MapAttribute{
							MarkdownDescription: "Map of key-value pairs which represents the attributes for the child object.",
							Optional:            true,
							Computed:            true,
							ElementType:         types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *AciRestManagedDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of datasource: aci_rest_managed")
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
	tflog.Debug(ctx, "End configure of datasource: aci_rest_managed")
}

func (d *AciRestManagedDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Start read of datasource: aci_rest_managed")
	var data *AciRestManagedResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	setAciRestManagedProperties(data)

	// Create a copy of the Id for when not found during setAciRestManagedAttributes
	cachedId := data.Id.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("read of datasource aci_rest_managed with id '%s'", data.Id.ValueString()))

	setAciRestManagedAttributes(ctx, &resp.Diagnostics, d.client, data)

	if data.Id.IsNull() {
		resp.Diagnostics.AddError(
			"Failed to read aci_rest_managed data source",
			fmt.Sprintf("The aci_rest_managed data source with id '%s' has not been found", cachedId),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, "End read of datasource: aci_rest_managed")
}
