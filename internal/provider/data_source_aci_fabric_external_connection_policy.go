// Code generated by "gen/generator.go"; DO NOT EDIT.
// In order to regenerate this file execute `go generate` from the repository root.
// More details can be found in the [README](https://github.com/CiscoDevNet/terraform-provider-aci/blob/master/README.md).

package provider

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &FvFabricExtConnPDataSource{}

func NewFvFabricExtConnPDataSource() datasource.DataSource {
	return &FvFabricExtConnPDataSource{}
}

// FvFabricExtConnPDataSource defines the data source implementation.
type FvFabricExtConnPDataSource struct {
	client *client.Client
}

func (d *FvFabricExtConnPDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of datasource: aci_fabric_external_connection_policy")
	resp.TypeName = req.ProviderTypeName + "_fabric_external_connection_policy"
	tflog.Debug(ctx, "End metadata of datasource: aci_fabric_external_connection_policy")
}

func (d *FvFabricExtConnPDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of datasource: aci_fabric_external_connection_policy")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "The fabric_external_connection_policy datasource for the 'fvFabricExtConnP' class",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The distinguished name (DN) of the Fabric External Connection Policy object.",
			},
			"parent_dn": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The distinguished name (DN) of the parent object.",
			},
			"annotation": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: `The annotation of the Fabric External Connection Policy object.`,
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: `The description of the Fabric External Connection Policy object.`,
			},
			"fabric_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: `A unique identifier of the fabric, associated with the Fabric External Connection Policy object.`,
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: `The name of the Fabric External Connection Policy object.`,
			},
			"name_alias": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: `The name alias of the Fabric External Connection Policy object.`,
			},
			"owner_key": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: `The key for enabling clients to own their data for entity correlation.`,
			},
			"owner_tag": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: `A tag for enabling clients to add their own data. For example, to indicate who created this object.`,
			},
			"community": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: `A global route target used to define communities for route leaking or redistribution in multi-pod or multi-site deployments to manage routing policies across fabrics.`,
			},
			"site_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: `A unique identifier for the site associated with the Fabric External Connection Policy object.`,
			},
			"peering_profile": schema.SingleNestedAttribute{
				MarkdownDescription: `Peering Profile`,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"annotation": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: `The annotation of the BGP EVPN Peering Profile object.`,
					},
					"description": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: `The description of the BGP EVPN Peering Profile object.`,
					},
					"name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: `The name of the BGP EVPN Peering Profile object.`,
					},
					"name_alias": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: `The name alias of the BGP EVPN Peering Profile object.`,
					},
					"owner_key": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: `The key for enabling clients to own their data for entity correlation.`,
					},
					"owner_tag": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: `A tag for enabling clients to add their own data. For example, to indicate who created this object.`,
					},
					"password": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: `The password used for establishing automatic BGP peering sessions.`,
					},
					"type": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: `The type of BGP EVPN Peering Profile object.`,
					},
					"annotations": schema.SetNestedAttribute{
						MarkdownDescription: ``,
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"key": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: `The key used to uniquely identify this configuration object.`,
								},
								"value": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: `The value of the property.`,
								},
							},
						},
					},
					"tags": schema.SetNestedAttribute{
						MarkdownDescription: ``,
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"key": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: `The key used to uniquely identify this configuration object.`,
								},
								"value": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: `The value of the property.`,
								},
							},
						},
					},
				},
			},
			"annotations": schema.SetNestedAttribute{
				MarkdownDescription: ``,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: `The key used to uniquely identify this configuration object.`,
						},
						"value": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: `The value of the property.`,
						},
					},
				},
			},
			"tags": schema.SetNestedAttribute{
				MarkdownDescription: ``,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: `The key used to uniquely identify this configuration object.`,
						},
						"value": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: `The value of the property.`,
						},
					},
				},
			},
		},
	}
	tflog.Debug(ctx, "End schema of datasource: aci_fabric_external_connection_policy")
}

func (d *FvFabricExtConnPDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of datasource: aci_fabric_external_connection_policy")
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
	tflog.Debug(ctx, "End configure of datasource: aci_fabric_external_connection_policy")
}

func (d *FvFabricExtConnPDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Start read of datasource: aci_fabric_external_connection_policy")
	var data *FvFabricExtConnPResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.ParentDn.IsNull() || data.ParentDn.IsUnknown() {
		data.ParentDn = basetypes.NewStringValue("uni/tn-infra")
	}

	setFvFabricExtConnPId(ctx, data)

	// Create a copy of the Id for when not found during getAndSetFvFabricExtConnPAttributes
	cachedId := data.Id.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Read of datasource aci_fabric_external_connection_policy with id '%s'", data.Id.ValueString()))

	getAndSetFvFabricExtConnPAttributes(ctx, &resp.Diagnostics, d.client, data)

	if data.Id.IsNull() {
		resp.Diagnostics.AddError(
			"Failed to read aci_fabric_external_connection_policy data source",
			fmt.Sprintf("The aci_fabric_external_connection_policy data source with id '%s' has not been found", cachedId),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End read of datasource aci_fabric_external_connection_policy with id '%s'", data.Id.ValueString()))
}
