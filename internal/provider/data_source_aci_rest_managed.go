package provider

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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

// AciRestManagedDataSourceModel describes the data source model.
type AciRestManagedDataSourceModel struct {
	Id         types.String `tfsdk:"id"`
	Dn         types.String `tfsdk:"dn"`
	ClassName  types.String `tfsdk:"class_name"`
	Content    types.Map    `tfsdk:"content"`
	Child      types.Set    `tfsdk:"child"`
	Annotation types.String `tfsdk:"annotation"`
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
				MarkdownDescription: "The distinguished name (DN) of the object.",
			},
			"dn": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The distinguished name (DN) of the parent object. e.g. uni/tn-EXAMPLE_TENANT",
			},
			"class_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Which class object is being created. (Make sure there is no colon in the classname)",
			},
			"content": schema.MapAttribute{
				MarkdownDescription: "Map of key-value pairs those needed to be passed to the Model object as parameters. Make sure the key name matches the name with the object parameter in ACI.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"annotation": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: `The annotation of the ACI object.`,
			},
			"escape_html": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Enable escaping of HTML characters when encoding the JSON payload.",
			},
		},
		Blocks: map[string]schema.Block{
			"child": schema.SetNestedBlock{
				MarkdownDescription: "List of children.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"rn": schema.StringAttribute{
							MarkdownDescription: "The relative name of the child object.",
							Computed:            true,
						},
						"class_name": schema.StringAttribute{
							MarkdownDescription: "Class name of child object.",
							Computed:            true,
						},
						"content": schema.MapAttribute{
							MarkdownDescription: "Map of key-value pairs which represents the attributes for the child object.",
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
	var data *AciRestManagedDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	requestData := DoRestRequest(ctx, &resp.Diagnostics, d.client, fmt.Sprintf("api/mo/%s.json?rsp-subtree=children", data.Dn.ValueString()), "GET", nil)

	if resp.Diagnostics.HasError() {
		return
	}

	if requestData.Search("imdata").Index(0).Data() == nil {
		data.Id = basetypes.NewStringNull()
		dataSourceRestManagedNotFoundError(&resp.Diagnostics, data)
		return
	}

	for className := range requestData.Search("imdata").Index(0).Data().(map[string]interface{}) {
		tflog.Debug(ctx, fmt.Sprintf("Setting ClassName to %s", className))
		data.ClassName = basetypes.NewStringValue(className)
		break
	}

	if requestData.Search("imdata").Search(data.ClassName.ValueString()).Data() != nil {
		classReadInfo := requestData.Search("imdata").Search(data.ClassName.ValueString()).Data().([]interface{})
		if len(classReadInfo) == 1 {
			content := map[string]attr.Value{}
			for attributeName, attributeValue := range classReadInfo[0].(map[string]interface{})["attributes"].(map[string]interface{}) {
				if attributeName == "dn" {
					dn := attributeValue.(string)
					data.Id = basetypes.NewStringValue(dn)
					data.Dn = basetypes.NewStringValue(dn)
				} else if attributeName == "annotation" {
					data.Annotation = basetypes.NewStringValue(attributeValue.(string))
				} else {
					content[attributeName] = basetypes.NewStringValue(attributeValue.(string))
				}
			}
			data.Content, _ = types.MapValue(types.StringType, content)

			childList := make([]ChildAciRestManagedResourceModel, 0)
			if _, ok := classReadInfo[0].(map[string]interface{})["children"]; ok {
				for _, child := range classReadInfo[0].(map[string]interface{})["children"].([]interface{}) {
					for childClassName, childClassDetails := range child.(map[string]interface{}) {
						childAttributes := childClassDetails.(map[string]interface{})["attributes"].(map[string]interface{})
						childContents := map[string]attr.Value{}
						ChildAciRestManaged := ChildAciRestManagedResourceModel{}
						ChildAciRestManaged.ClassName = basetypes.NewStringValue(childClassName)

						if val, ok := childAttributes["rn"]; ok {
							ChildAciRestManaged.Rn = basetypes.NewStringValue(val.(string))
						}

						for childAttributeName, childAttributeValue := range childAttributes {
							if childAttributeName != "rn" {
								childContents[childAttributeName] = basetypes.NewStringValue(childAttributeValue.(string))
							}
						}
						ChildAciRestManaged.Content, _ = types.MapValue(types.StringType, childContents)
						childList = append(childList, ChildAciRestManaged)
					}
				}
			}

			if len(childList) > 0 {
				data.Child, _ = types.SetValueFrom(ctx, data.Child.ElementType(ctx), childList)
			}

		} else {
			resp.Diagnostics.AddError(
				"Too many results in response",
				fmt.Sprintf("%v matches returned for class '%s'. Please report this issue to the provider developers.", len(classReadInfo), data.ClassName),
			)
		}
	} else {
		dataSourceRestManagedNotFoundError(&resp.Diagnostics, data)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End read of datasource aci_rest_managed with id '%s'", data.Id.ValueString()))

}

func dataSourceRestManagedNotFoundError(diags *diag.Diagnostics, data *AciRestManagedDataSourceModel) {
	diags.AddError(
		"Failed to read aci_rest_managed data source",
		fmt.Sprintf("The aci_rest_managed data source with dn '%s' has not been found", data.Dn),
	)
}
