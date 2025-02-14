// Code generated by "gen/generator.go"; DO NOT EDIT.
// In order to regenerate this file execute `go generate` from the repository root.
// More details can be found in the [README](https://github.com/CiscoDevNet/terraform-provider-aci/blob/master/README.md).

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &FvEpMacTagResource{}
var _ resource.ResourceWithImportState = &FvEpMacTagResource{}

func NewFvEpMacTagResource() resource.Resource {
	return &FvEpMacTagResource{}
}

// FvEpMacTagResource defines the resource implementation.
type FvEpMacTagResource struct {
	client *client.Client
}

// FvEpMacTagResourceModel describes the resource data model.
type FvEpMacTagResourceModel struct {
	Id            types.String `tfsdk:"id"`
	ParentDn      types.String `tfsdk:"parent_dn"`
	Annotation    types.String `tfsdk:"annotation"`
	BdName        types.String `tfsdk:"bd_name"`
	FvEpMacTagId  types.String `tfsdk:"id_attribute"`
	Mac           types.String `tfsdk:"mac"`
	Name          types.String `tfsdk:"name"`
	NameAlias     types.String `tfsdk:"name_alias"`
	TagAnnotation types.Set    `tfsdk:"annotations"`
	TagTag        types.Set    `tfsdk:"tags"`
}

func getEmptyFvEpMacTagResourceModel() *FvEpMacTagResourceModel {
	return &FvEpMacTagResourceModel{
		Id:           basetypes.NewStringNull(),
		ParentDn:     basetypes.NewStringNull(),
		Annotation:   basetypes.NewStringNull(),
		BdName:       basetypes.NewStringNull(),
		FvEpMacTagId: basetypes.NewStringNull(),
		Mac:          basetypes.NewStringNull(),
		Name:         basetypes.NewStringNull(),
		NameAlias:    basetypes.NewStringNull(),
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

// TagAnnotationFvEpMacTagResourceModel describes the resource data model for the children without relation ships.
type TagAnnotationFvEpMacTagResourceModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

func getEmptyTagAnnotationFvEpMacTagResourceModel() TagAnnotationFvEpMacTagResourceModel {
	return TagAnnotationFvEpMacTagResourceModel{
		Key:   basetypes.NewStringNull(),
		Value: basetypes.NewStringNull(),
	}
}

var TagAnnotationFvEpMacTagType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	},
}

// TagTagFvEpMacTagResourceModel describes the resource data model for the children without relation ships.
type TagTagFvEpMacTagResourceModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

func getEmptyTagTagFvEpMacTagResourceModel() TagTagFvEpMacTagResourceModel {
	return TagTagFvEpMacTagResourceModel{
		Key:   basetypes.NewStringNull(),
		Value: basetypes.NewStringNull(),
	}
}

var TagTagFvEpMacTagType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	},
}

type FvEpMacTagIdentifier struct {
	BdName types.String
	Mac    types.String
}

func (r *FvEpMacTagResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if !req.Plan.Raw.IsNull() {
		var planData, stateData *FvEpMacTagResourceModel
		resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
		resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

		if resp.Diagnostics.HasError() {
			return
		}

		if (planData.Id.IsUnknown() || planData.Id.IsNull()) && !planData.ParentDn.IsUnknown() && !planData.BdName.IsUnknown() && !planData.Mac.IsUnknown() {
			setFvEpMacTagId(ctx, planData)
		}

		if stateData == nil && !globalAllowExistingOnCreate && !planData.Id.IsUnknown() && !planData.Id.IsNull() {
			CheckDn(ctx, &resp.Diagnostics, r.client, "fvEpMacTag", planData.Id.ValueString())
			if resp.Diagnostics.HasError() {
				return
			}
		}

		resp.Diagnostics.Append(resp.Plan.Set(ctx, &planData)...)
	}
}

func (r *FvEpMacTagResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of resource: aci_endpoint_tag_mac")
	resp.TypeName = req.ProviderTypeName + "_endpoint_tag_mac"
	tflog.Debug(ctx, "End metadata of resource: aci_endpoint_tag_mac")
}

func (r *FvEpMacTagResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of resource: aci_endpoint_tag_mac")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "The endpoint_tag_mac resource for the 'fvEpMacTag' class",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The distinguished name (DN) of the Endpoint Tag MAC object.",
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
				MarkdownDescription: `The annotation of the Endpoint Tag MAC object.`,
			},
			"bd_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: `The bridge domain name of the Endpoint Tag MAC object.`,
			},
			"id_attribute": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
				MarkdownDescription: `The identifier of the Endpoint Tag MAC object.`,
			},
			"mac": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: `The MAC address of the Endpoint Tag MAC object.`,
			},
			"name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
				MarkdownDescription: `The name of the Endpoint Tag MAC object.`,
			},
			"name_alias": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
				MarkdownDescription: `The name alias of the Endpoint Tag MAC object.`,
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
	tflog.Debug(ctx, "End schema of resource: aci_endpoint_tag_mac")
}

func (r *FvEpMacTagResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of resource: aci_endpoint_tag_mac")
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
	tflog.Debug(ctx, "End configure of resource: aci_endpoint_tag_mac")
}

func (r *FvEpMacTagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start create of resource: aci_endpoint_tag_mac")
	// On create retrieve information on current state prior to making any changes in order to determine child delete operations
	var stateData *FvEpMacTagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &stateData)...)
	if stateData.Id.IsUnknown() || stateData.Id.IsNull() {
		setFvEpMacTagId(ctx, stateData)
	}
	getAndSetFvEpMacTagAttributes(ctx, &resp.Diagnostics, r.client, stateData)
	if !globalAllowExistingOnCreate && !stateData.Id.IsNull() {
		resp.Diagnostics.AddError(
			"Object Already Exists",
			fmt.Sprintf("The fvEpMacTag object with DN '%s' already exists.", stateData.Id.ValueString()),
		)
		return
	}

	var data *FvEpMacTagResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Id.IsUnknown() || data.Id.IsNull() {
		setFvEpMacTagId(ctx, data)
	}

	tflog.Debug(ctx, fmt.Sprintf("Create of resource aci_endpoint_tag_mac with id '%s'", data.Id.ValueString()))

	var tagAnnotationPlan, tagAnnotationState []TagAnnotationFvEpMacTagResourceModel
	data.TagAnnotation.ElementsAs(ctx, &tagAnnotationPlan, false)
	stateData.TagAnnotation.ElementsAs(ctx, &tagAnnotationState, false)
	var tagTagPlan, tagTagState []TagTagFvEpMacTagResourceModel
	data.TagTag.ElementsAs(ctx, &tagTagPlan, false)
	stateData.TagTag.ElementsAs(ctx, &tagTagState, false)
	jsonPayload := getFvEpMacTagCreateJsonPayload(ctx, &resp.Diagnostics, true, data, tagAnnotationPlan, tagAnnotationState, tagTagPlan, tagTagState)

	if resp.Diagnostics.HasError() {
		return
	}

	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("api/mo/%s.json", data.Id.ValueString()), "POST", jsonPayload)

	if resp.Diagnostics.HasError() {
		return
	}

	getAndSetFvEpMacTagAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End create of resource aci_endpoint_tag_mac with id '%s'", data.Id.ValueString()))
}

func (r *FvEpMacTagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start read of resource: aci_endpoint_tag_mac")
	var data *FvEpMacTagResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Read of resource aci_endpoint_tag_mac with id '%s'", data.Id.ValueString()))

	getAndSetFvEpMacTagAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	if data.Id.IsNull() {
		var emptyData *FvEpMacTagResourceModel
		resp.Diagnostics.Append(resp.State.Set(ctx, &emptyData)...)
	} else {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}

	tflog.Debug(ctx, fmt.Sprintf("End read of resource aci_endpoint_tag_mac with id '%s'", data.Id.ValueString()))
}

func (r *FvEpMacTagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start update of resource: aci_endpoint_tag_mac")
	var data *FvEpMacTagResourceModel
	var stateData *FvEpMacTagResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Update of resource aci_endpoint_tag_mac with id '%s'", data.Id.ValueString()))

	var tagAnnotationPlan, tagAnnotationState []TagAnnotationFvEpMacTagResourceModel
	data.TagAnnotation.ElementsAs(ctx, &tagAnnotationPlan, false)
	stateData.TagAnnotation.ElementsAs(ctx, &tagAnnotationState, false)
	var tagTagPlan, tagTagState []TagTagFvEpMacTagResourceModel
	data.TagTag.ElementsAs(ctx, &tagTagPlan, false)
	stateData.TagTag.ElementsAs(ctx, &tagTagState, false)
	jsonPayload := getFvEpMacTagCreateJsonPayload(ctx, &resp.Diagnostics, false, data, tagAnnotationPlan, tagAnnotationState, tagTagPlan, tagTagState)

	if resp.Diagnostics.HasError() {
		return
	}

	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("api/mo/%s.json", data.Id.ValueString()), "POST", jsonPayload)

	if resp.Diagnostics.HasError() {
		return
	}

	getAndSetFvEpMacTagAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End update of resource aci_endpoint_tag_mac with id '%s'", data.Id.ValueString()))
}

func (r *FvEpMacTagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start delete of resource: aci_endpoint_tag_mac")
	var data *FvEpMacTagResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Delete of resource aci_endpoint_tag_mac with id '%s'", data.Id.ValueString()))
	jsonPayload := GetDeleteJsonPayload(ctx, &resp.Diagnostics, "fvEpMacTag", data.Id.ValueString())
	if resp.Diagnostics.HasError() {
		return
	}
	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("api/mo/%s.json", data.Id.ValueString()), "POST", jsonPayload)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("End delete of resource aci_endpoint_tag_mac with id '%s'", data.Id.ValueString()))
}

func (r *FvEpMacTagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start import state of resource: aci_endpoint_tag_mac")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	var stateData *FvEpMacTagResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &stateData)...)
	tflog.Debug(ctx, fmt.Sprintf("Import state of resource aci_endpoint_tag_mac with id '%s'", stateData.Id.ValueString()))

	tflog.Debug(ctx, "End import of state resource: aci_endpoint_tag_mac")
}

func getAndSetFvEpMacTagAttributes(ctx context.Context, diags *diag.Diagnostics, client *client.Client, data *FvEpMacTagResourceModel) {
	requestData := DoRestRequest(ctx, diags, client, fmt.Sprintf("api/mo/%s.json?rsp-subtree=full&rsp-subtree-class=%s", data.Id.ValueString(), "fvEpMacTag,tagAnnotation,tagTag"), "GET", nil)

	readData := getEmptyFvEpMacTagResourceModel()

	if diags.HasError() {
		return
	}
	if requestData.Search("imdata").Search("fvEpMacTag").Data() != nil {
		classReadInfo := requestData.Search("imdata").Search("fvEpMacTag").Data().([]interface{})
		if len(classReadInfo) == 1 {
			attributes := classReadInfo[0].(map[string]interface{})["attributes"].(map[string]interface{})
			for attributeName, attributeValue := range attributes {
				if attributeName == "dn" {
					readData.Id = basetypes.NewStringValue(attributeValue.(string))
					setFvEpMacTagParentDn(ctx, attributeValue.(string), readData)
				}
				if attributeName == "annotation" {
					readData.Annotation = basetypes.NewStringValue(attributeValue.(string))
				}
				if attributeName == "bdName" {
					readData.BdName = basetypes.NewStringValue(attributeValue.(string))
				}
				if attributeName == "id" {
					readData.FvEpMacTagId = basetypes.NewStringValue(attributeValue.(string))
				}
				if attributeName == "mac" {
					readData.Mac = basetypes.NewStringValue(attributeValue.(string))
				}
				if attributeName == "name" {
					readData.Name = basetypes.NewStringValue(attributeValue.(string))
				}
				if attributeName == "nameAlias" {
					readData.NameAlias = basetypes.NewStringValue(attributeValue.(string))
				}
			}
			TagAnnotationFvEpMacTagList := make([]TagAnnotationFvEpMacTagResourceModel, 0)
			TagTagFvEpMacTagList := make([]TagTagFvEpMacTagResourceModel, 0)
			_, ok := classReadInfo[0].(map[string]interface{})["children"]
			if ok {
				children := classReadInfo[0].(map[string]interface{})["children"].([]interface{})
				for _, child := range children {
					for childClassName, childClassDetails := range child.(map[string]interface{}) {
						childAttributes := childClassDetails.(map[string]interface{})["attributes"].(map[string]interface{})
						if childClassName == "tagAnnotation" {
							TagAnnotationFvEpMacTag := getEmptyTagAnnotationFvEpMacTagResourceModel()
							for childAttributeName, childAttributeValue := range childAttributes {
								if childAttributeName == "key" {
									TagAnnotationFvEpMacTag.Key = basetypes.NewStringValue(childAttributeValue.(string))
								}
								if childAttributeName == "value" {
									TagAnnotationFvEpMacTag.Value = basetypes.NewStringValue(childAttributeValue.(string))
								}

							}
							TagAnnotationFvEpMacTagList = append(TagAnnotationFvEpMacTagList, TagAnnotationFvEpMacTag)
						}
						if childClassName == "tagTag" {
							TagTagFvEpMacTag := getEmptyTagTagFvEpMacTagResourceModel()
							for childAttributeName, childAttributeValue := range childAttributes {
								if childAttributeName == "key" {
									TagTagFvEpMacTag.Key = basetypes.NewStringValue(childAttributeValue.(string))
								}
								if childAttributeName == "value" {
									TagTagFvEpMacTag.Value = basetypes.NewStringValue(childAttributeValue.(string))
								}

							}
							TagTagFvEpMacTagList = append(TagTagFvEpMacTagList, TagTagFvEpMacTag)
						}
					}
				}
			}
			tagAnnotationSet, _ := types.SetValueFrom(ctx, readData.TagAnnotation.ElementType(ctx), TagAnnotationFvEpMacTagList)
			readData.TagAnnotation = tagAnnotationSet
			tagTagSet, _ := types.SetValueFrom(ctx, readData.TagTag.ElementType(ctx), TagTagFvEpMacTagList)
			readData.TagTag = tagTagSet
		} else {
			diags.AddError(
				"too many results in response",
				fmt.Sprintf("%v matches returned for class 'fvEpMacTag'. Please report this issue to the provider developers.", len(classReadInfo)),
			)
		}
	} else {
		readData.Id = basetypes.NewStringNull()
	}
	*data = *readData
}

func getFvEpMacTagRn(ctx context.Context, data *FvEpMacTagResourceModel) string {
	return fmt.Sprintf("eptags/epmactag-%s-[%s]", data.Mac.ValueString(), data.BdName.ValueString())
}

func setFvEpMacTagParentDn(ctx context.Context, dn string, data *FvEpMacTagResourceModel) {
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
	parentDn := dn[:rnIndex]
	parentDn = parentDn[:strings.Index(parentDn, "/eptags")]
	data.ParentDn = basetypes.NewStringValue(parentDn)
}

func setFvEpMacTagId(ctx context.Context, data *FvEpMacTagResourceModel) {
	rn := getFvEpMacTagRn(ctx, data)
	data.Id = types.StringValue(fmt.Sprintf("%s/%s", data.ParentDn.ValueString(), rn))
}

func getFvEpMacTagTagAnnotationChildPayloads(ctx context.Context, diags *diag.Diagnostics, data *FvEpMacTagResourceModel, tagAnnotationFvEpMacTagPlan, tagAnnotationFvEpMacTagState []TagAnnotationFvEpMacTagResourceModel) []map[string]interface{} {
	childPayloads := []map[string]interface{}{}
	if !data.TagAnnotation.IsNull() && !data.TagAnnotation.IsUnknown() {
		tagAnnotationIdentifiers := []TagAnnotationIdentifier{}
		for _, tagAnnotationFvEpMacTag := range tagAnnotationFvEpMacTagPlan {
			childMap := NewAciObject()
			if !tagAnnotationFvEpMacTag.Key.IsNull() && !tagAnnotationFvEpMacTag.Key.IsUnknown() {
				childMap.Attributes["key"] = tagAnnotationFvEpMacTag.Key.ValueString()
			}
			if !tagAnnotationFvEpMacTag.Value.IsNull() && !tagAnnotationFvEpMacTag.Value.IsUnknown() {
				childMap.Attributes["value"] = tagAnnotationFvEpMacTag.Value.ValueString()
			}
			childPayloads = append(childPayloads, map[string]interface{}{"tagAnnotation": childMap})
			tagAnnotationIdentifier := TagAnnotationIdentifier{}
			tagAnnotationIdentifier.Key = tagAnnotationFvEpMacTag.Key
			tagAnnotationIdentifiers = append(tagAnnotationIdentifiers, tagAnnotationIdentifier)
		}
		for _, tagAnnotation := range tagAnnotationFvEpMacTagState {
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

func getFvEpMacTagTagTagChildPayloads(ctx context.Context, diags *diag.Diagnostics, data *FvEpMacTagResourceModel, tagTagFvEpMacTagPlan, tagTagFvEpMacTagState []TagTagFvEpMacTagResourceModel) []map[string]interface{} {
	childPayloads := []map[string]interface{}{}
	if !data.TagTag.IsNull() && !data.TagTag.IsUnknown() {
		tagTagIdentifiers := []TagTagIdentifier{}
		for _, tagTagFvEpMacTag := range tagTagFvEpMacTagPlan {
			childMap := NewAciObject()
			if !tagTagFvEpMacTag.Key.IsNull() && !tagTagFvEpMacTag.Key.IsUnknown() {
				childMap.Attributes["key"] = tagTagFvEpMacTag.Key.ValueString()
			}
			if !tagTagFvEpMacTag.Value.IsNull() && !tagTagFvEpMacTag.Value.IsUnknown() {
				childMap.Attributes["value"] = tagTagFvEpMacTag.Value.ValueString()
			}
			childPayloads = append(childPayloads, map[string]interface{}{"tagTag": childMap})
			tagTagIdentifier := TagTagIdentifier{}
			tagTagIdentifier.Key = tagTagFvEpMacTag.Key
			tagTagIdentifiers = append(tagTagIdentifiers, tagTagIdentifier)
		}
		for _, tagTag := range tagTagFvEpMacTagState {
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

func getFvEpMacTagCreateJsonPayload(ctx context.Context, diags *diag.Diagnostics, createType bool, data *FvEpMacTagResourceModel, tagAnnotationPlan, tagAnnotationState []TagAnnotationFvEpMacTagResourceModel, tagTagPlan, tagTagState []TagTagFvEpMacTagResourceModel) *container.Container {
	payloadMap := map[string]interface{}{}
	payloadMap["attributes"] = map[string]string{}

	if createType && !globalAllowExistingOnCreate {
		payloadMap["attributes"].(map[string]string)["status"] = "created"
	}
	childPayloads := []map[string]interface{}{}

	TagAnnotationchildPayloads := getFvEpMacTagTagAnnotationChildPayloads(ctx, diags, data, tagAnnotationPlan, tagAnnotationState)
	if TagAnnotationchildPayloads == nil {
		return nil
	}
	childPayloads = append(childPayloads, TagAnnotationchildPayloads...)

	TagTagchildPayloads := getFvEpMacTagTagTagChildPayloads(ctx, diags, data, tagTagPlan, tagTagState)
	if TagTagchildPayloads == nil {
		return nil
	}
	childPayloads = append(childPayloads, TagTagchildPayloads...)

	payloadMap["children"] = childPayloads
	if !data.Annotation.IsNull() && !data.Annotation.IsUnknown() {
		payloadMap["attributes"].(map[string]string)["annotation"] = data.Annotation.ValueString()
	}
	if !data.BdName.IsNull() && !data.BdName.IsUnknown() {
		payloadMap["attributes"].(map[string]string)["bdName"] = data.BdName.ValueString()
	}
	if !data.FvEpMacTagId.IsNull() && !data.FvEpMacTagId.IsUnknown() {
		payloadMap["attributes"].(map[string]string)["id"] = data.FvEpMacTagId.ValueString()
	}
	if !data.Mac.IsNull() && !data.Mac.IsUnknown() {
		payloadMap["attributes"].(map[string]string)["mac"] = data.Mac.ValueString()
	}
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		payloadMap["attributes"].(map[string]string)["name"] = data.Name.ValueString()
	}
	if !data.NameAlias.IsNull() && !data.NameAlias.IsUnknown() {
		payloadMap["attributes"].(map[string]string)["nameAlias"] = data.NameAlias.ValueString()
	}
	payload, err := json.Marshal(map[string]interface{}{"fvEpMacTag": payloadMap})
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
