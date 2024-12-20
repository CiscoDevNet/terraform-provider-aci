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
var _ resource.Resource = &L3extRsRedistributePolResource{}
var _ resource.ResourceWithImportState = &L3extRsRedistributePolResource{}

func NewL3extRsRedistributePolResource() resource.Resource {
	return &L3extRsRedistributePolResource{}
}

// L3extRsRedistributePolResource defines the resource implementation.
type L3extRsRedistributePolResource struct {
	client *client.Client
}

// L3extRsRedistributePolResourceModel describes the resource data model.
type L3extRsRedistributePolResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	ParentDn            types.String `tfsdk:"parent_dn"`
	Annotation          types.String `tfsdk:"annotation"`
	Src                 types.String `tfsdk:"source"`
	TnRtctrlProfileName types.String `tfsdk:"route_control_profile_name"`
	TagAnnotation       types.Set    `tfsdk:"annotations"`
	TagTag              types.Set    `tfsdk:"tags"`
}

func getEmptyL3extRsRedistributePolResourceModel() *L3extRsRedistributePolResourceModel {
	return &L3extRsRedistributePolResourceModel{
		Id:                  basetypes.NewStringNull(),
		ParentDn:            basetypes.NewStringNull(),
		Annotation:          basetypes.NewStringNull(),
		Src:                 basetypes.NewStringNull(),
		TnRtctrlProfileName: basetypes.NewStringNull(),
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

// TagAnnotationL3extRsRedistributePolResourceModel describes the resource data model for the children without relation ships.
type TagAnnotationL3extRsRedistributePolResourceModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

func getEmptyTagAnnotationL3extRsRedistributePolResourceModel() TagAnnotationL3extRsRedistributePolResourceModel {
	return TagAnnotationL3extRsRedistributePolResourceModel{
		Key:   basetypes.NewStringNull(),
		Value: basetypes.NewStringNull(),
	}
}

var TagAnnotationL3extRsRedistributePolType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	},
}

// TagTagL3extRsRedistributePolResourceModel describes the resource data model for the children without relation ships.
type TagTagL3extRsRedistributePolResourceModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

func getEmptyTagTagL3extRsRedistributePolResourceModel() TagTagL3extRsRedistributePolResourceModel {
	return TagTagL3extRsRedistributePolResourceModel{
		Key:   basetypes.NewStringNull(),
		Value: basetypes.NewStringNull(),
	}
}

var TagTagL3extRsRedistributePolType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	},
}

type L3extRsRedistributePolIdentifier struct {
	Src                 types.String
	TnRtctrlProfileName types.String
}

func (r *L3extRsRedistributePolResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if !req.Plan.Raw.IsNull() {
		var planData, stateData *L3extRsRedistributePolResourceModel
		resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
		resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

		if resp.Diagnostics.HasError() {
			return
		}

		if (planData.Id.IsUnknown() || planData.Id.IsNull()) && !planData.ParentDn.IsUnknown() && !planData.Src.IsUnknown() && !planData.TnRtctrlProfileName.IsUnknown() {
			setL3extRsRedistributePolId(ctx, planData)
		}

		if stateData == nil && !globalAllowExistingOnCreate && !planData.Id.IsUnknown() && !planData.Id.IsNull() {
			CheckDn(ctx, &resp.Diagnostics, r.client, "l3extRsRedistributePol", planData.Id.ValueString())
			if resp.Diagnostics.HasError() {
				return
			}
		}

		resp.Diagnostics.Append(resp.Plan.Set(ctx, &planData)...)
	}
}

func (r *L3extRsRedistributePolResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of resource: aci_l3out_redistribute_policy")
	resp.TypeName = req.ProviderTypeName + "_l3out_redistribute_policy"
	tflog.Debug(ctx, "End metadata of resource: aci_l3out_redistribute_policy")
}

func (r *L3extRsRedistributePolResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of resource: aci_l3out_redistribute_policy")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "The l3out_redistribute_policy resource for the 'l3extRsRedistributePol' class",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The distinguished name (DN) of the L3Out Redistribute Policy object.",
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
				MarkdownDescription: `The annotation of the L3Out Redistribute Policy object.`,
			},
			"source": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("attached-host", "direct", "static"),
				},
				MarkdownDescription: `The source of the L3Out Redistribute Policy object.`,
			},
			"route_control_profile_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: `The name of the Route Control Profile object.`,
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
	tflog.Debug(ctx, "End schema of resource: aci_l3out_redistribute_policy")
}

func (r *L3extRsRedistributePolResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of resource: aci_l3out_redistribute_policy")
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
	tflog.Debug(ctx, "End configure of resource: aci_l3out_redistribute_policy")
}

func (r *L3extRsRedistributePolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start create of resource: aci_l3out_redistribute_policy")
	// On create retrieve information on current state prior to making any changes in order to determine child delete operations
	var stateData *L3extRsRedistributePolResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &stateData)...)
	if stateData.Id.IsUnknown() || stateData.Id.IsNull() {
		setL3extRsRedistributePolId(ctx, stateData)
	}
	getAndSetL3extRsRedistributePolAttributes(ctx, &resp.Diagnostics, r.client, stateData)
	if !globalAllowExistingOnCreate && !stateData.Id.IsNull() {
		resp.Diagnostics.AddError(
			"Object Already Exists",
			fmt.Sprintf("The l3extRsRedistributePol object with DN '%s' already exists.", stateData.Id.ValueString()),
		)
		return
	}

	var data *L3extRsRedistributePolResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Id.IsUnknown() || data.Id.IsNull() {
		setL3extRsRedistributePolId(ctx, data)
	}

	tflog.Debug(ctx, fmt.Sprintf("Create of resource aci_l3out_redistribute_policy with id '%s'", data.Id.ValueString()))

	var tagAnnotationPlan, tagAnnotationState []TagAnnotationL3extRsRedistributePolResourceModel
	data.TagAnnotation.ElementsAs(ctx, &tagAnnotationPlan, false)
	stateData.TagAnnotation.ElementsAs(ctx, &tagAnnotationState, false)
	var tagTagPlan, tagTagState []TagTagL3extRsRedistributePolResourceModel
	data.TagTag.ElementsAs(ctx, &tagTagPlan, false)
	stateData.TagTag.ElementsAs(ctx, &tagTagState, false)
	jsonPayload := getL3extRsRedistributePolCreateJsonPayload(ctx, &resp.Diagnostics, true, data, tagAnnotationPlan, tagAnnotationState, tagTagPlan, tagTagState)

	if resp.Diagnostics.HasError() {
		return
	}

	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("api/mo/%s.json", data.Id.ValueString()), "POST", jsonPayload)

	if resp.Diagnostics.HasError() {
		return
	}

	getAndSetL3extRsRedistributePolAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End create of resource aci_l3out_redistribute_policy with id '%s'", data.Id.ValueString()))
}

func (r *L3extRsRedistributePolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start read of resource: aci_l3out_redistribute_policy")
	var data *L3extRsRedistributePolResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Read of resource aci_l3out_redistribute_policy with id '%s'", data.Id.ValueString()))

	getAndSetL3extRsRedistributePolAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	if data.Id.IsNull() {
		var emptyData *L3extRsRedistributePolResourceModel
		resp.Diagnostics.Append(resp.State.Set(ctx, &emptyData)...)
	} else {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}

	tflog.Debug(ctx, fmt.Sprintf("End read of resource aci_l3out_redistribute_policy with id '%s'", data.Id.ValueString()))
}

func (r *L3extRsRedistributePolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start update of resource: aci_l3out_redistribute_policy")
	var data *L3extRsRedistributePolResourceModel
	var stateData *L3extRsRedistributePolResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Update of resource aci_l3out_redistribute_policy with id '%s'", data.Id.ValueString()))

	var tagAnnotationPlan, tagAnnotationState []TagAnnotationL3extRsRedistributePolResourceModel
	data.TagAnnotation.ElementsAs(ctx, &tagAnnotationPlan, false)
	stateData.TagAnnotation.ElementsAs(ctx, &tagAnnotationState, false)
	var tagTagPlan, tagTagState []TagTagL3extRsRedistributePolResourceModel
	data.TagTag.ElementsAs(ctx, &tagTagPlan, false)
	stateData.TagTag.ElementsAs(ctx, &tagTagState, false)
	jsonPayload := getL3extRsRedistributePolCreateJsonPayload(ctx, &resp.Diagnostics, false, data, tagAnnotationPlan, tagAnnotationState, tagTagPlan, tagTagState)

	if resp.Diagnostics.HasError() {
		return
	}

	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("api/mo/%s.json", data.Id.ValueString()), "POST", jsonPayload)

	if resp.Diagnostics.HasError() {
		return
	}

	getAndSetL3extRsRedistributePolAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End update of resource aci_l3out_redistribute_policy with id '%s'", data.Id.ValueString()))
}

func (r *L3extRsRedistributePolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start delete of resource: aci_l3out_redistribute_policy")
	var data *L3extRsRedistributePolResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Delete of resource aci_l3out_redistribute_policy with id '%s'", data.Id.ValueString()))
	jsonPayload := GetDeleteJsonPayload(ctx, &resp.Diagnostics, "l3extRsRedistributePol", data.Id.ValueString())
	if resp.Diagnostics.HasError() {
		return
	}
	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("api/mo/%s.json", data.Id.ValueString()), "POST", jsonPayload)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("End delete of resource aci_l3out_redistribute_policy with id '%s'", data.Id.ValueString()))
}

func (r *L3extRsRedistributePolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start import state of resource: aci_l3out_redistribute_policy")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	var stateData *L3extRsRedistributePolResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &stateData)...)
	tflog.Debug(ctx, fmt.Sprintf("Import state of resource aci_l3out_redistribute_policy with id '%s'", stateData.Id.ValueString()))

	tflog.Debug(ctx, "End import of state resource: aci_l3out_redistribute_policy")
}

func getAndSetL3extRsRedistributePolAttributes(ctx context.Context, diags *diag.Diagnostics, client *client.Client, data *L3extRsRedistributePolResourceModel) {
	requestData := DoRestRequest(ctx, diags, client, fmt.Sprintf("api/mo/%s.json?rsp-subtree=full&rsp-subtree-class=%s", data.Id.ValueString(), "l3extRsRedistributePol,tagAnnotation,tagTag"), "GET", nil)

	readData := getEmptyL3extRsRedistributePolResourceModel()

	if diags.HasError() {
		return
	}
	if requestData.Search("imdata").Search("l3extRsRedistributePol").Data() != nil {
		classReadInfo := requestData.Search("imdata").Search("l3extRsRedistributePol").Data().([]interface{})
		if len(classReadInfo) == 1 {
			attributes := classReadInfo[0].(map[string]interface{})["attributes"].(map[string]interface{})
			for attributeName, attributeValue := range attributes {
				if attributeName == "dn" {
					readData.Id = basetypes.NewStringValue(attributeValue.(string))
					setL3extRsRedistributePolParentDn(ctx, attributeValue.(string), readData)
				}
				if attributeName == "annotation" {
					readData.Annotation = basetypes.NewStringValue(attributeValue.(string))
				}
				if attributeName == "src" {
					readData.Src = basetypes.NewStringValue(attributeValue.(string))
				}
				if attributeName == "tnRtctrlProfileName" {
					readData.TnRtctrlProfileName = basetypes.NewStringValue(attributeValue.(string))
				}
			}
			TagAnnotationL3extRsRedistributePolList := make([]TagAnnotationL3extRsRedistributePolResourceModel, 0)
			TagTagL3extRsRedistributePolList := make([]TagTagL3extRsRedistributePolResourceModel, 0)
			_, ok := classReadInfo[0].(map[string]interface{})["children"]
			if ok {
				children := classReadInfo[0].(map[string]interface{})["children"].([]interface{})
				for _, child := range children {
					for childClassName, childClassDetails := range child.(map[string]interface{}) {
						childAttributes := childClassDetails.(map[string]interface{})["attributes"].(map[string]interface{})
						if childClassName == "tagAnnotation" {
							TagAnnotationL3extRsRedistributePol := getEmptyTagAnnotationL3extRsRedistributePolResourceModel()
							for childAttributeName, childAttributeValue := range childAttributes {
								if childAttributeName == "key" {
									TagAnnotationL3extRsRedistributePol.Key = basetypes.NewStringValue(childAttributeValue.(string))
								}
								if childAttributeName == "value" {
									TagAnnotationL3extRsRedistributePol.Value = basetypes.NewStringValue(childAttributeValue.(string))
								}

							}
							TagAnnotationL3extRsRedistributePolList = append(TagAnnotationL3extRsRedistributePolList, TagAnnotationL3extRsRedistributePol)
						}
						if childClassName == "tagTag" {
							TagTagL3extRsRedistributePol := getEmptyTagTagL3extRsRedistributePolResourceModel()
							for childAttributeName, childAttributeValue := range childAttributes {
								if childAttributeName == "key" {
									TagTagL3extRsRedistributePol.Key = basetypes.NewStringValue(childAttributeValue.(string))
								}
								if childAttributeName == "value" {
									TagTagL3extRsRedistributePol.Value = basetypes.NewStringValue(childAttributeValue.(string))
								}

							}
							TagTagL3extRsRedistributePolList = append(TagTagL3extRsRedistributePolList, TagTagL3extRsRedistributePol)
						}
					}
				}
			}
			tagAnnotationSet, _ := types.SetValueFrom(ctx, readData.TagAnnotation.ElementType(ctx), TagAnnotationL3extRsRedistributePolList)
			readData.TagAnnotation = tagAnnotationSet
			tagTagSet, _ := types.SetValueFrom(ctx, readData.TagTag.ElementType(ctx), TagTagL3extRsRedistributePolList)
			readData.TagTag = tagTagSet
		} else {
			diags.AddError(
				"too many results in response",
				fmt.Sprintf("%v matches returned for class 'l3extRsRedistributePol'. Please report this issue to the provider developers.", len(classReadInfo)),
			)
		}
	} else {
		readData.Id = basetypes.NewStringNull()
	}
	*data = *readData
}

func getL3extRsRedistributePolRn(ctx context.Context, data *L3extRsRedistributePolResourceModel) string {
	return fmt.Sprintf("rsredistributePol-[%s]-%s", data.TnRtctrlProfileName.ValueString(), data.Src.ValueString())
}

func setL3extRsRedistributePolParentDn(ctx context.Context, dn string, data *L3extRsRedistributePolResourceModel) {
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

func setL3extRsRedistributePolId(ctx context.Context, data *L3extRsRedistributePolResourceModel) {
	rn := getL3extRsRedistributePolRn(ctx, data)
	data.Id = types.StringValue(fmt.Sprintf("%s/%s", data.ParentDn.ValueString(), rn))
}

func getL3extRsRedistributePolTagAnnotationChildPayloads(ctx context.Context, diags *diag.Diagnostics, data *L3extRsRedistributePolResourceModel, tagAnnotationL3extRsRedistributePolPlan, tagAnnotationL3extRsRedistributePolState []TagAnnotationL3extRsRedistributePolResourceModel) []map[string]interface{} {
	childPayloads := []map[string]interface{}{}
	if !data.TagAnnotation.IsNull() && !data.TagAnnotation.IsUnknown() {
		tagAnnotationIdentifiers := []TagAnnotationIdentifier{}
		for _, tagAnnotationL3extRsRedistributePol := range tagAnnotationL3extRsRedistributePolPlan {
			childMap := NewAciObject()
			if !tagAnnotationL3extRsRedistributePol.Key.IsNull() && !tagAnnotationL3extRsRedistributePol.Key.IsUnknown() {
				childMap.Attributes["key"] = tagAnnotationL3extRsRedistributePol.Key.ValueString()
			}
			if !tagAnnotationL3extRsRedistributePol.Value.IsNull() && !tagAnnotationL3extRsRedistributePol.Value.IsUnknown() {
				childMap.Attributes["value"] = tagAnnotationL3extRsRedistributePol.Value.ValueString()
			}
			childPayloads = append(childPayloads, map[string]interface{}{"tagAnnotation": childMap})
			tagAnnotationIdentifier := TagAnnotationIdentifier{}
			tagAnnotationIdentifier.Key = tagAnnotationL3extRsRedistributePol.Key
			tagAnnotationIdentifiers = append(tagAnnotationIdentifiers, tagAnnotationIdentifier)
		}
		for _, tagAnnotation := range tagAnnotationL3extRsRedistributePolState {
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

func getL3extRsRedistributePolTagTagChildPayloads(ctx context.Context, diags *diag.Diagnostics, data *L3extRsRedistributePolResourceModel, tagTagL3extRsRedistributePolPlan, tagTagL3extRsRedistributePolState []TagTagL3extRsRedistributePolResourceModel) []map[string]interface{} {
	childPayloads := []map[string]interface{}{}
	if !data.TagTag.IsNull() && !data.TagTag.IsUnknown() {
		tagTagIdentifiers := []TagTagIdentifier{}
		for _, tagTagL3extRsRedistributePol := range tagTagL3extRsRedistributePolPlan {
			childMap := NewAciObject()
			if !tagTagL3extRsRedistributePol.Key.IsNull() && !tagTagL3extRsRedistributePol.Key.IsUnknown() {
				childMap.Attributes["key"] = tagTagL3extRsRedistributePol.Key.ValueString()
			}
			if !tagTagL3extRsRedistributePol.Value.IsNull() && !tagTagL3extRsRedistributePol.Value.IsUnknown() {
				childMap.Attributes["value"] = tagTagL3extRsRedistributePol.Value.ValueString()
			}
			childPayloads = append(childPayloads, map[string]interface{}{"tagTag": childMap})
			tagTagIdentifier := TagTagIdentifier{}
			tagTagIdentifier.Key = tagTagL3extRsRedistributePol.Key
			tagTagIdentifiers = append(tagTagIdentifiers, tagTagIdentifier)
		}
		for _, tagTag := range tagTagL3extRsRedistributePolState {
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

func getL3extRsRedistributePolCreateJsonPayload(ctx context.Context, diags *diag.Diagnostics, createType bool, data *L3extRsRedistributePolResourceModel, tagAnnotationPlan, tagAnnotationState []TagAnnotationL3extRsRedistributePolResourceModel, tagTagPlan, tagTagState []TagTagL3extRsRedistributePolResourceModel) *container.Container {
	payloadMap := map[string]interface{}{}
	payloadMap["attributes"] = map[string]string{}

	if createType && !globalAllowExistingOnCreate {
		payloadMap["attributes"].(map[string]string)["status"] = "created"
	}
	childPayloads := []map[string]interface{}{}

	TagAnnotationchildPayloads := getL3extRsRedistributePolTagAnnotationChildPayloads(ctx, diags, data, tagAnnotationPlan, tagAnnotationState)
	if TagAnnotationchildPayloads == nil {
		return nil
	}
	childPayloads = append(childPayloads, TagAnnotationchildPayloads...)

	TagTagchildPayloads := getL3extRsRedistributePolTagTagChildPayloads(ctx, diags, data, tagTagPlan, tagTagState)
	if TagTagchildPayloads == nil {
		return nil
	}
	childPayloads = append(childPayloads, TagTagchildPayloads...)

	payloadMap["children"] = childPayloads
	if !data.Annotation.IsNull() && !data.Annotation.IsUnknown() {
		payloadMap["attributes"].(map[string]string)["annotation"] = data.Annotation.ValueString()
	}
	if !data.Src.IsNull() && !data.Src.IsUnknown() {
		payloadMap["attributes"].(map[string]string)["src"] = data.Src.ValueString()
	}
	if !data.TnRtctrlProfileName.IsNull() && !data.TnRtctrlProfileName.IsUnknown() {
		payloadMap["attributes"].(map[string]string)["tnRtctrlProfileName"] = data.TnRtctrlProfileName.ValueString()
	}
	payload, err := json.Marshal(map[string]interface{}{"l3extRsRedistributePol": payloadMap})
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
