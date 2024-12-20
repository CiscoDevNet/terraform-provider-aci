// Code generated by "gen/generator.go"; DO NOT EDIT.
// In order to regenerate this file execute `go generate` from the repository root.
// More details can be found in the [README](https://github.com/CiscoDevNet/terraform-provider-aci/blob/master/README.md).

package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &FvVmAttrResource{}
var _ resource.ResourceWithImportState = &FvVmAttrResource{}

func NewFvVmAttrResource() resource.Resource {
	return &FvVmAttrResource{}
}

// FvVmAttrResource defines the resource implementation.
type FvVmAttrResource struct {
	client *client.Client
}

// FvVmAttrResourceModel describes the resource data model.
type FvVmAttrResourceModel struct {
	Id            types.String `tfsdk:"id"`
	ParentDn      types.String `tfsdk:"parent_dn"`
	Annotation    types.String `tfsdk:"annotation"`
	Category      types.String `tfsdk:"category"`
	Descr         types.String `tfsdk:"description"`
	LabelName     types.String `tfsdk:"label_name"`
	Name          types.String `tfsdk:"name"`
	NameAlias     types.String `tfsdk:"name_alias"`
	Operator      types.String `tfsdk:"operator"`
	OwnerKey      types.String `tfsdk:"owner_key"`
	OwnerTag      types.String `tfsdk:"owner_tag"`
	Type          types.String `tfsdk:"type"`
	Value         types.String `tfsdk:"value"`
	TagAnnotation types.Set    `tfsdk:"annotations"`
	TagTag        types.Set    `tfsdk:"tags"`
}

func getEmptyFvVmAttrResourceModel() *FvVmAttrResourceModel {
	return &FvVmAttrResourceModel{
		Id:         basetypes.NewStringNull(),
		ParentDn:   basetypes.NewStringNull(),
		Annotation: basetypes.NewStringNull(),
		Category:   basetypes.NewStringNull(),
		Descr:      basetypes.NewStringNull(),
		LabelName:  basetypes.NewStringNull(),
		Name:       basetypes.NewStringNull(),
		NameAlias:  basetypes.NewStringNull(),
		Operator:   basetypes.NewStringNull(),
		OwnerKey:   basetypes.NewStringNull(),
		OwnerTag:   basetypes.NewStringNull(),
		Type:       basetypes.NewStringNull(),
		Value:      basetypes.NewStringNull(),
		TagAnnotation: types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"key":   types.StringType,
				"value": types.StringType,
			},
		}),
		TagTag: types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"key":   types.StringType,
				"value": types.StringType,
			},
		}),
	}
}

// TagAnnotationFvVmAttrResourceModel describes the resource data model for the children without relation ships.
type TagAnnotationFvVmAttrResourceModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

func getEmptyTagAnnotationFvVmAttrResourceModel() TagAnnotationFvVmAttrResourceModel {
	return TagAnnotationFvVmAttrResourceModel{
		Key:   basetypes.NewStringNull(),
		Value: basetypes.NewStringNull(),
	}
}

var TagAnnotationFvVmAttrType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	},
}

// TagTagFvVmAttrResourceModel describes the resource data model for the children without relation ships.
type TagTagFvVmAttrResourceModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

func getEmptyTagTagFvVmAttrResourceModel() TagTagFvVmAttrResourceModel {
	return TagTagFvVmAttrResourceModel{
		Key:   basetypes.NewStringNull(),
		Value: basetypes.NewStringNull(),
	}
}

var TagTagFvVmAttrType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	},
}

type FvVmAttrIdentifier struct {
	Name types.String
}

func (r *FvVmAttrResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if !req.Plan.Raw.IsNull() {
		var planData, stateData *FvVmAttrResourceModel
		resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
		resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

		if resp.Diagnostics.HasError() {
			return
		}

		if (planData.Id.IsUnknown() || planData.Id.IsNull()) && !planData.ParentDn.IsUnknown() && !planData.Name.IsUnknown() {
			setFvVmAttrId(ctx, planData)
		}

		if stateData == nil && !globalAllowExistingOnCreate && !planData.Id.IsUnknown() && !planData.Id.IsNull() {
			CheckDn(ctx, &resp.Diagnostics, r.client, "fvVmAttr", planData.Id.ValueString())
			if resp.Diagnostics.HasError() {
				return
			}
		}

		resp.Diagnostics.Append(resp.Plan.Set(ctx, &planData)...)
	}
}

func (r *FvVmAttrResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of resource: aci_epg_useg_vm_attribute")
	resp.TypeName = req.ProviderTypeName + "_epg_useg_vm_attribute"
	tflog.Debug(ctx, "End metadata of resource: aci_epg_useg_vm_attribute")
}

func (r *FvVmAttrResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of resource: aci_epg_useg_vm_attribute")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "The epg_useg_vm_attribute resource for the 'fvVmAttr' class",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The distinguished name (DN) of the EPG uSeg VM Attribute object.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"parent_dn": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The distinguished name (DN) of the parent object.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"annotation": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
				Default:             stringdefault.StaticString(globalAnnotation),
				MarkdownDescription: `The annotation of the EPG uSeg VM Attribute object.`,
			},
			"category": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
				MarkdownDescription: `The category of the EPG uSeg VM Attribute object.`,
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
				MarkdownDescription: `The description of the EPG uSeg VM Attribute object.`,
			},
			"label_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
				MarkdownDescription: `The label name of the EPG uSeg VM Attribute object.`,
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: `The name of the EPG uSeg VM Attribute object.`,
			},
			"name_alias": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
				MarkdownDescription: `The name alias of the EPG uSeg VM Attribute object.`,
			},
			"operator": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("contains", "endsWith", "equals", "notEquals", "startsWith"),
				},
				MarkdownDescription: `The operator of the EPG uSeg VM Attribute object.`,
			},
			"owner_key": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
				MarkdownDescription: `The key for enabling clients to own their data for entity correlation.`,
			},
			"owner_tag": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
				MarkdownDescription: `A tag for enabling clients to add their own data. For example, to indicate who created this object.`,
			},
			"type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("custom-label", "domain", "guest-os", "hv", "rootContName", "tag", "vm", "vm-folder", "vm-name", "vmfolder-path", "vnic"),
				},
				MarkdownDescription: `The type of the EPG uSeg VM Attribute object.`,
			},
			"value": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
				MarkdownDescription: `The value of the EPG uSeg VM Attribute object.`,
			},
			"annotations": schema.SetNestedAttribute{
				MarkdownDescription: ``,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
							MarkdownDescription: `The key used to uniquely identify this configuration object.`,
						},
						"value": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
							MarkdownDescription: `The value of the property.`,
						},
					},
				},
			},
			"tags": schema.SetNestedAttribute{
				MarkdownDescription: ``,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
							MarkdownDescription: `The key used to uniquely identify this configuration object.`,
						},
						"value": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
							MarkdownDescription: `The value of the property.`,
						},
					},
				},
			},
		},
	}
	tflog.Debug(ctx, "End schema of resource: aci_epg_useg_vm_attribute")
}

func (r *FvVmAttrResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of resource: aci_epg_useg_vm_attribute")
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
	tflog.Debug(ctx, "End configure of resource: aci_epg_useg_vm_attribute")
}

func (r *FvVmAttrResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start create of resource: aci_epg_useg_vm_attribute")
	// On create retrieve information on current state prior to making any changes in order to determine child delete operations
	var stateData *FvVmAttrResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &stateData)...)
	if stateData.Id.IsUnknown() || stateData.Id.IsNull() {
		setFvVmAttrId(ctx, stateData)
	}
	getAndSetFvVmAttrAttributes(ctx, &resp.Diagnostics, r.client, stateData)
	if !globalAllowExistingOnCreate && !stateData.Id.IsNull() {
		resp.Diagnostics.AddError(
			"Object Already Exists",
			fmt.Sprintf("The fvVmAttr object with DN '%s' already exists.", stateData.Id.ValueString()),
		)
		return
	}

	var data *FvVmAttrResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Id.IsUnknown() || data.Id.IsNull() {
		setFvVmAttrId(ctx, data)
	}

	tflog.Debug(ctx, fmt.Sprintf("Create of resource aci_epg_useg_vm_attribute with id '%s'", data.Id.ValueString()))

	var tagAnnotationPlan, tagAnnotationState []TagAnnotationFvVmAttrResourceModel
	data.TagAnnotation.ElementsAs(ctx, &tagAnnotationPlan, false)
	stateData.TagAnnotation.ElementsAs(ctx, &tagAnnotationState, false)
	var tagTagPlan, tagTagState []TagTagFvVmAttrResourceModel
	data.TagTag.ElementsAs(ctx, &tagTagPlan, false)
	stateData.TagTag.ElementsAs(ctx, &tagTagState, false)
	jsonPayload := getFvVmAttrCreateJsonPayload(ctx, &resp.Diagnostics, true, data, tagAnnotationPlan, tagAnnotationState, tagTagPlan, tagTagState)

	if resp.Diagnostics.HasError() {
		return
	}

	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("api/mo/%s.json", data.Id.ValueString()), "POST", jsonPayload)

	if resp.Diagnostics.HasError() {
		return
	}

	getAndSetFvVmAttrAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End create of resource aci_epg_useg_vm_attribute with id '%s'", data.Id.ValueString()))
}

func (r *FvVmAttrResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start read of resource: aci_epg_useg_vm_attribute")
	var data *FvVmAttrResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Read of resource aci_epg_useg_vm_attribute with id '%s'", data.Id.ValueString()))

	getAndSetFvVmAttrAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	if data.Id.IsNull() {
		var emptyData *FvVmAttrResourceModel
		resp.Diagnostics.Append(resp.State.Set(ctx, &emptyData)...)
	} else {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}

	tflog.Debug(ctx, fmt.Sprintf("End read of resource aci_epg_useg_vm_attribute with id '%s'", data.Id.ValueString()))
}

func (r *FvVmAttrResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start update of resource: aci_epg_useg_vm_attribute")
	var data *FvVmAttrResourceModel
	var stateData *FvVmAttrResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Update of resource aci_epg_useg_vm_attribute with id '%s'", data.Id.ValueString()))

	var tagAnnotationPlan, tagAnnotationState []TagAnnotationFvVmAttrResourceModel
	data.TagAnnotation.ElementsAs(ctx, &tagAnnotationPlan, false)
	stateData.TagAnnotation.ElementsAs(ctx, &tagAnnotationState, false)
	var tagTagPlan, tagTagState []TagTagFvVmAttrResourceModel
	data.TagTag.ElementsAs(ctx, &tagTagPlan, false)
	stateData.TagTag.ElementsAs(ctx, &tagTagState, false)
	jsonPayload := getFvVmAttrCreateJsonPayload(ctx, &resp.Diagnostics, false, data, tagAnnotationPlan, tagAnnotationState, tagTagPlan, tagTagState)

	if resp.Diagnostics.HasError() {
		return
	}

	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("api/mo/%s.json", data.Id.ValueString()), "POST", jsonPayload)

	if resp.Diagnostics.HasError() {
		return
	}

	getAndSetFvVmAttrAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End update of resource aci_epg_useg_vm_attribute with id '%s'", data.Id.ValueString()))
}

func (r *FvVmAttrResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start delete of resource: aci_epg_useg_vm_attribute")
	var data *FvVmAttrResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Delete of resource aci_epg_useg_vm_attribute with id '%s'", data.Id.ValueString()))
	jsonPayload := GetDeleteJsonPayload(ctx, &resp.Diagnostics, "fvVmAttr", data.Id.ValueString())
	if resp.Diagnostics.HasError() {
		return
	}
	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("api/mo/%s.json", data.Id.ValueString()), "POST", jsonPayload)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("End delete of resource aci_epg_useg_vm_attribute with id '%s'", data.Id.ValueString()))
}

func (r *FvVmAttrResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start import state of resource: aci_epg_useg_vm_attribute")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	var stateData *FvVmAttrResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &stateData)...)
	tflog.Debug(ctx, fmt.Sprintf("Import state of resource aci_epg_useg_vm_attribute with id '%s'", stateData.Id.ValueString()))

	tflog.Debug(ctx, "End import of state resource: aci_epg_useg_vm_attribute")
}

func getAndSetFvVmAttrAttributes(ctx context.Context, diags *diag.Diagnostics, client *client.Client, data *FvVmAttrResourceModel) {
	requestData := DoRestRequest(ctx, diags, client, fmt.Sprintf("api/mo/%s.json?rsp-subtree=full&rsp-subtree-class=%s", data.Id.ValueString(), "fvVmAttr,tagAnnotation,tagTag"), "GET", nil)

	readData := getEmptyFvVmAttrResourceModel()

	if diags.HasError() {
		return
	}
	if requestData.Search("imdata").Search("fvVmAttr").Data() != nil {
		classReadInfo := requestData.Search("imdata").Search("fvVmAttr").Data().([]interface{})
		if len(classReadInfo) == 1 {
			attributes := classReadInfo[0].(map[string]interface{})["attributes"].(map[string]interface{})
			for attributeName, attributeValue := range attributes {
				if attributeName == "dn" {
					readData.Id = basetypes.NewStringValue(attributeValue.(string))
					setFvVmAttrParentDn(ctx, attributeValue.(string), readData)
				}
				if attributeName == "annotation" {
					readData.Annotation = basetypes.NewStringValue(attributeValue.(string))
				}
				if attributeName == "category" {
					readData.Category = basetypes.NewStringValue(attributeValue.(string))
				}
				if attributeName == "descr" {
					readData.Descr = basetypes.NewStringValue(attributeValue.(string))
				}
				if attributeName == "labelName" {
					readData.LabelName = basetypes.NewStringValue(attributeValue.(string))
				}
				if attributeName == "name" {
					readData.Name = basetypes.NewStringValue(attributeValue.(string))
				}
				if attributeName == "nameAlias" {
					readData.NameAlias = basetypes.NewStringValue(attributeValue.(string))
				}
				if attributeName == "operator" {
					readData.Operator = basetypes.NewStringValue(attributeValue.(string))
				}
				if attributeName == "ownerKey" {
					readData.OwnerKey = basetypes.NewStringValue(attributeValue.(string))
				}
				if attributeName == "ownerTag" {
					readData.OwnerTag = basetypes.NewStringValue(attributeValue.(string))
				}
				if attributeName == "type" {
					readData.Type = basetypes.NewStringValue(attributeValue.(string))
				}
				if attributeName == "value" {
					readData.Value = basetypes.NewStringValue(attributeValue.(string))
				}
			}
			TagAnnotationFvVmAttrList := make([]TagAnnotationFvVmAttrResourceModel, 0)
			TagTagFvVmAttrList := make([]TagTagFvVmAttrResourceModel, 0)
			_, ok := classReadInfo[0].(map[string]interface{})["children"]
			if ok {
				children := classReadInfo[0].(map[string]interface{})["children"].([]interface{})
				for _, child := range children {
					for childClassName, childClassDetails := range child.(map[string]interface{}) {
						childAttributes := childClassDetails.(map[string]interface{})["attributes"].(map[string]interface{})
						if childClassName == "tagAnnotation" {
							TagAnnotationFvVmAttr := getEmptyTagAnnotationFvVmAttrResourceModel()
							for childAttributeName, childAttributeValue := range childAttributes {
								if childAttributeName == "key" {
									TagAnnotationFvVmAttr.Key = basetypes.NewStringValue(childAttributeValue.(string))
								}
								if childAttributeName == "value" {
									TagAnnotationFvVmAttr.Value = basetypes.NewStringValue(childAttributeValue.(string))
								}

							}
							TagAnnotationFvVmAttrList = append(TagAnnotationFvVmAttrList, TagAnnotationFvVmAttr)
						}
						if childClassName == "tagTag" {
							TagTagFvVmAttr := getEmptyTagTagFvVmAttrResourceModel()
							for childAttributeName, childAttributeValue := range childAttributes {
								if childAttributeName == "key" {
									TagTagFvVmAttr.Key = basetypes.NewStringValue(childAttributeValue.(string))
								}
								if childAttributeName == "value" {
									TagTagFvVmAttr.Value = basetypes.NewStringValue(childAttributeValue.(string))
								}

							}
							TagTagFvVmAttrList = append(TagTagFvVmAttrList, TagTagFvVmAttr)
						}
					}
				}
			}
			tagAnnotationSet, _ := types.SetValueFrom(ctx, readData.TagAnnotation.ElementType(ctx), TagAnnotationFvVmAttrList)
			readData.TagAnnotation = tagAnnotationSet
			tagTagSet, _ := types.SetValueFrom(ctx, readData.TagTag.ElementType(ctx), TagTagFvVmAttrList)
			readData.TagTag = tagTagSet
		} else {
			diags.AddError(
				"too many results in response",
				fmt.Sprintf("%v matches returned for class 'fvVmAttr'. Please report this issue to the provider developers.", len(classReadInfo)),
			)
		}
	} else {
		readData.Id = basetypes.NewStringNull()
	}
	*data = *readData
}

func getFvVmAttrRn(ctx context.Context, data *FvVmAttrResourceModel) string {
	return fmt.Sprintf("vmattr-%s", data.Name.ValueString())
}

func setFvVmAttrParentDn(ctx context.Context, dn string, data *FvVmAttrResourceModel) {
	bracketIndex := 0
	rnIndex := 0
	for i := len(dn) - 1; i >= 0; i-- {
		if string(dn[i]) == "]" {
			bracketIndex = bracketIndex + 1
		} else if string(dn[i]) == "[" {
			bracketIndex = bracketIndex - 1
		} else if string(dn[i]) == "/" && bracketIndex == 0 {
			rnIndex = i
			break
		}
	}
	data.ParentDn = basetypes.NewStringValue(dn[:rnIndex])
}

func setFvVmAttrId(ctx context.Context, data *FvVmAttrResourceModel) {
	rn := getFvVmAttrRn(ctx, data)
	data.Id = types.StringValue(fmt.Sprintf("%s/%s", data.ParentDn.ValueString(), rn))
}

func getFvVmAttrTagAnnotationChildPayloads(ctx context.Context, diags *diag.Diagnostics, data *FvVmAttrResourceModel, tagAnnotationFvVmAttrPlan, tagAnnotationFvVmAttrState []TagAnnotationFvVmAttrResourceModel) []map[string]interface{} {
	childPayloads := []map[string]interface{}{}
	if !data.TagAnnotation.IsNull() && !data.TagAnnotation.IsUnknown() {
		tagAnnotationIdentifiers := []TagAnnotationIdentifier{}
		for _, tagAnnotationFvVmAttr := range tagAnnotationFvVmAttrPlan {
			childMap := NewAciObject()
			if !tagAnnotationFvVmAttr.Key.IsNull() && !tagAnnotationFvVmAttr.Key.IsUnknown() {
				childMap.Attributes["key"] = tagAnnotationFvVmAttr.Key.ValueString()
			}
			if !tagAnnotationFvVmAttr.Value.IsNull() && !tagAnnotationFvVmAttr.Value.IsUnknown() {
				childMap.Attributes["value"] = tagAnnotationFvVmAttr.Value.ValueString()
			}
			childPayloads = append(childPayloads, map[string]interface{}{"tagAnnotation": childMap})
			tagAnnotationIdentifier := TagAnnotationIdentifier{}
			tagAnnotationIdentifier.Key = tagAnnotationFvVmAttr.Key
			tagAnnotationIdentifiers = append(tagAnnotationIdentifiers, tagAnnotationIdentifier)
		}
		for _, tagAnnotation := range tagAnnotationFvVmAttrState {
			delete := true
			for _, tagAnnotationIdentifier := range tagAnnotationIdentifiers {
				if tagAnnotationIdentifier.Key == tagAnnotation.Key {
					delete = false
					break
				}
			}
			if delete {
				tagAnnotationChildMapForDelete := NewAciObject()
				tagAnnotationChildMapForDelete.Attributes["status"] = "deleted"
				tagAnnotationChildMapForDelete.Attributes["key"] = tagAnnotation.Key.ValueString()
				childPayloads = append(childPayloads, map[string]interface{}{"tagAnnotation": tagAnnotationChildMapForDelete})
			}
		}
	} else {
		data.TagAnnotation = types.SetNull(data.TagAnnotation.ElementType(ctx))
	}

	return childPayloads
}

func getFvVmAttrTagTagChildPayloads(ctx context.Context, diags *diag.Diagnostics, data *FvVmAttrResourceModel, tagTagFvVmAttrPlan, tagTagFvVmAttrState []TagTagFvVmAttrResourceModel) []map[string]interface{} {
	childPayloads := []map[string]interface{}{}
	if !data.TagTag.IsNull() && !data.TagTag.IsUnknown() {
		tagTagIdentifiers := []TagTagIdentifier{}
		for _, tagTagFvVmAttr := range tagTagFvVmAttrPlan {
			childMap := NewAciObject()
			if !tagTagFvVmAttr.Key.IsNull() && !tagTagFvVmAttr.Key.IsUnknown() {
				childMap.Attributes["key"] = tagTagFvVmAttr.Key.ValueString()
			}
			if !tagTagFvVmAttr.Value.IsNull() && !tagTagFvVmAttr.Value.IsUnknown() {
				childMap.Attributes["value"] = tagTagFvVmAttr.Value.ValueString()
			}
			childPayloads = append(childPayloads, map[string]interface{}{"tagTag": childMap})
			tagTagIdentifier := TagTagIdentifier{}
			tagTagIdentifier.Key = tagTagFvVmAttr.Key
			tagTagIdentifiers = append(tagTagIdentifiers, tagTagIdentifier)
		}
		for _, tagTag := range tagTagFvVmAttrState {
			delete := true
			for _, tagTagIdentifier := range tagTagIdentifiers {
				if tagTagIdentifier.Key == tagTag.Key {
					delete = false
					break
				}
			}
			if delete {
				tagTagChildMapForDelete := NewAciObject()
				tagTagChildMapForDelete.Attributes["status"] = "deleted"
				tagTagChildMapForDelete.Attributes["key"] = tagTag.Key.ValueString()
				childPayloads = append(childPayloads, map[string]interface{}{"tagTag": tagTagChildMapForDelete})
			}
		}
	} else {
		data.TagTag = types.SetNull(data.TagTag.ElementType(ctx))
	}

	return childPayloads
}

func getFvVmAttrCreateJsonPayload(ctx context.Context, diags *diag.Diagnostics, createType bool, data *FvVmAttrResourceModel, tagAnnotationPlan, tagAnnotationState []TagAnnotationFvVmAttrResourceModel, tagTagPlan, tagTagState []TagTagFvVmAttrResourceModel) *container.Container {
	payloadMap := map[string]interface{}{}
	payloadMap["attributes"] = map[string]string{}

	if createType && !globalAllowExistingOnCreate {
		payloadMap["attributes"].(map[string]string)["status"] = "created"
	}
	childPayloads := []map[string]interface{}{}

	TagAnnotationchildPayloads := getFvVmAttrTagAnnotationChildPayloads(ctx, diags, data, tagAnnotationPlan, tagAnnotationState)
	if TagAnnotationchildPayloads == nil {
		return nil
	}
	childPayloads = append(childPayloads, TagAnnotationchildPayloads...)

	TagTagchildPayloads := getFvVmAttrTagTagChildPayloads(ctx, diags, data, tagTagPlan, tagTagState)
	if TagTagchildPayloads == nil {
		return nil
	}
	childPayloads = append(childPayloads, TagTagchildPayloads...)

	payloadMap["children"] = childPayloads
	if !data.Annotation.IsNull() && !data.Annotation.IsUnknown() {
		payloadMap["attributes"].(map[string]string)["annotation"] = data.Annotation.ValueString()
	}
	if !data.Category.IsNull() && !data.Category.IsUnknown() {
		payloadMap["attributes"].(map[string]string)["category"] = data.Category.ValueString()
	}
	if !data.Descr.IsNull() && !data.Descr.IsUnknown() {
		payloadMap["attributes"].(map[string]string)["descr"] = data.Descr.ValueString()
	}
	if !data.LabelName.IsNull() && !data.LabelName.IsUnknown() {
		payloadMap["attributes"].(map[string]string)["labelName"] = data.LabelName.ValueString()
	}
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		payloadMap["attributes"].(map[string]string)["name"] = data.Name.ValueString()
	}
	if !data.NameAlias.IsNull() && !data.NameAlias.IsUnknown() {
		payloadMap["attributes"].(map[string]string)["nameAlias"] = data.NameAlias.ValueString()
	}
	if !data.Operator.IsNull() && !data.Operator.IsUnknown() {
		payloadMap["attributes"].(map[string]string)["operator"] = data.Operator.ValueString()
	}
	if !data.OwnerKey.IsNull() && !data.OwnerKey.IsUnknown() {
		payloadMap["attributes"].(map[string]string)["ownerKey"] = data.OwnerKey.ValueString()
	}
	if !data.OwnerTag.IsNull() && !data.OwnerTag.IsUnknown() {
		payloadMap["attributes"].(map[string]string)["ownerTag"] = data.OwnerTag.ValueString()
	}
	if !data.Type.IsNull() && !data.Type.IsUnknown() {
		payloadMap["attributes"].(map[string]string)["type"] = data.Type.ValueString()
	}
	if !data.Value.IsNull() && !data.Value.IsUnknown() {
		payloadMap["attributes"].(map[string]string)["value"] = data.Value.ValueString()
	}
	payload, err := json.Marshal(map[string]interface{}{"fvVmAttr": payloadMap})
	if err != nil {
		diags.AddError(
			"Marshalling of json payload failed",
			fmt.Sprintf("Err: %s. Please report this issue to the provider developers.", err),
		)
		return nil
	}

	jsonPayload, err := container.ParseJSON(payload)

	if err != nil {
		diags.AddError(
			"Construction of json payload failed",
			fmt.Sprintf("Err: %s. Please report this issue to the provider developers.", err),
		)
		return nil
	}
	return jsonPayload
}
