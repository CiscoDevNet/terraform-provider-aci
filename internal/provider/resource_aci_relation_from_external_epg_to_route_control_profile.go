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
var _ resource.Resource = &L3extRsInstPToProfileResource{}
var _ resource.ResourceWithImportState = &L3extRsInstPToProfileResource{}

func NewL3extRsInstPToProfileResource() resource.Resource {
	return &L3extRsInstPToProfileResource{}
}

// L3extRsInstPToProfileResource defines the resource implementation.
type L3extRsInstPToProfileResource struct {
	client *client.Client
}

// L3extRsInstPToProfileResourceModel describes the resource data model.
type L3extRsInstPToProfileResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	ParentDn            types.String `tfsdk:"parent_dn"`
	Annotation          types.String `tfsdk:"annotation"`
	Direction           types.String `tfsdk:"direction"`
	TnRtctrlProfileName types.String `tfsdk:"route_control_profile_name"`
	TagAnnotation       types.Set    `tfsdk:"annotations"`
	TagTag              types.Set    `tfsdk:"tags"`
}

func getEmptyL3extRsInstPToProfileResourceModel() *L3extRsInstPToProfileResourceModel {
	return &L3extRsInstPToProfileResourceModel{
		Id:                  basetypes.NewStringNull(),
		ParentDn:            basetypes.NewStringNull(),
		Annotation:          basetypes.NewStringNull(),
		Direction:           basetypes.NewStringNull(),
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

// TagAnnotationL3extRsInstPToProfileResourceModel describes the resource data model for the children without relation ships.
type TagAnnotationL3extRsInstPToProfileResourceModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

func getEmptyTagAnnotationL3extRsInstPToProfileResourceModel() TagAnnotationL3extRsInstPToProfileResourceModel {
	return TagAnnotationL3extRsInstPToProfileResourceModel{
		Key:   basetypes.NewStringNull(),
		Value: basetypes.NewStringNull(),
	}
}

var TagAnnotationL3extRsInstPToProfileType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	},
}

// TagTagL3extRsInstPToProfileResourceModel describes the resource data model for the children without relation ships.
type TagTagL3extRsInstPToProfileResourceModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

func getEmptyTagTagL3extRsInstPToProfileResourceModel() TagTagL3extRsInstPToProfileResourceModel {
	return TagTagL3extRsInstPToProfileResourceModel{
		Key:   basetypes.NewStringNull(),
		Value: basetypes.NewStringNull(),
	}
}

var TagTagL3extRsInstPToProfileType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	},
}

type L3extRsInstPToProfileIdentifier struct {
	Direction           types.String
	TnRtctrlProfileName types.String
}

func (r *L3extRsInstPToProfileResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if !req.Plan.Raw.IsNull() {
		var planData, stateData *L3extRsInstPToProfileResourceModel
		resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
		resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

		if resp.Diagnostics.HasError() {
			return
		}

		if (planData.Id.IsUnknown() || planData.Id.IsNull()) && !planData.ParentDn.IsUnknown() && !planData.Direction.IsUnknown() && !planData.TnRtctrlProfileName.IsUnknown() {
			setL3extRsInstPToProfileId(ctx, planData)
		}

		if stateData == nil && !globalAllowExistingOnCreate && !planData.Id.IsUnknown() && !planData.Id.IsNull() {
			CheckDn(ctx, &resp.Diagnostics, r.client, "l3extRsInstPToProfile", planData.Id.ValueString())
			if resp.Diagnostics.HasError() {
				return
			}
		}

		resp.Diagnostics.Append(resp.Plan.Set(ctx, &planData)...)
	}
}

func (r *L3extRsInstPToProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of resource: aci_relation_from_external_epg_to_route_control_profile")
	resp.TypeName = req.ProviderTypeName + "_relation_from_external_epg_to_route_control_profile"
	tflog.Debug(ctx, "End metadata of resource: aci_relation_from_external_epg_to_route_control_profile")
}

func (r *L3extRsInstPToProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of resource: aci_relation_from_external_epg_to_route_control_profile")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "The relation_from_external_epg_to_route_control_profile resource for the 'l3extRsInstPToProfile' class",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The distinguished name (DN) of the Relation From External EPG To Route Control Profile object.",
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
				MarkdownDescription: `The annotation of the Relation From External EPG To Route Control Profile object.`,
			},
			"direction": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("export", "import"),
				},
				MarkdownDescription: `The direction on which to apply the Route Map associated with the Route Control Profile.`,
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
	tflog.Debug(ctx, "End schema of resource: aci_relation_from_external_epg_to_route_control_profile")
}

func (r *L3extRsInstPToProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of resource: aci_relation_from_external_epg_to_route_control_profile")
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
	tflog.Debug(ctx, "End configure of resource: aci_relation_from_external_epg_to_route_control_profile")
}

func (r *L3extRsInstPToProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start create of resource: aci_relation_from_external_epg_to_route_control_profile")
	// On create retrieve information on current state prior to making any changes in order to determine child delete operations
	var stateData *L3extRsInstPToProfileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &stateData)...)
	if stateData.Id.IsUnknown() || stateData.Id.IsNull() {
		setL3extRsInstPToProfileId(ctx, stateData)
	}
	getAndSetL3extRsInstPToProfileAttributes(ctx, &resp.Diagnostics, r.client, stateData)
	if !globalAllowExistingOnCreate && !stateData.Id.IsNull() {
		resp.Diagnostics.AddError(
			"Object Already Exists",
			fmt.Sprintf("The l3extRsInstPToProfile object with DN '%s' already exists.", stateData.Id.ValueString()),
		)
		return
	}

	var data *L3extRsInstPToProfileResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Id.IsUnknown() || data.Id.IsNull() {
		setL3extRsInstPToProfileId(ctx, data)
	}

	tflog.Debug(ctx, fmt.Sprintf("Create of resource aci_relation_from_external_epg_to_route_control_profile with id '%s'", data.Id.ValueString()))

	var tagAnnotationPlan, tagAnnotationState []TagAnnotationL3extRsInstPToProfileResourceModel
	data.TagAnnotation.ElementsAs(ctx, &tagAnnotationPlan, false)
	stateData.TagAnnotation.ElementsAs(ctx, &tagAnnotationState, false)
	var tagTagPlan, tagTagState []TagTagL3extRsInstPToProfileResourceModel
	data.TagTag.ElementsAs(ctx, &tagTagPlan, false)
	stateData.TagTag.ElementsAs(ctx, &tagTagState, false)
	jsonPayload := getL3extRsInstPToProfileCreateJsonPayload(ctx, &resp.Diagnostics, true, data, tagAnnotationPlan, tagAnnotationState, tagTagPlan, tagTagState)

	if resp.Diagnostics.HasError() {
		return
	}

	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("api/mo/%s.json", data.Id.ValueString()), "POST", jsonPayload)

	if resp.Diagnostics.HasError() {
		return
	}

	getAndSetL3extRsInstPToProfileAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End create of resource aci_relation_from_external_epg_to_route_control_profile with id '%s'", data.Id.ValueString()))
}

func (r *L3extRsInstPToProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start read of resource: aci_relation_from_external_epg_to_route_control_profile")
	var data *L3extRsInstPToProfileResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Read of resource aci_relation_from_external_epg_to_route_control_profile with id '%s'", data.Id.ValueString()))

	getAndSetL3extRsInstPToProfileAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	if data.Id.IsNull() {
		var emptyData *L3extRsInstPToProfileResourceModel
		resp.Diagnostics.Append(resp.State.Set(ctx, &emptyData)...)
	} else {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}

	tflog.Debug(ctx, fmt.Sprintf("End read of resource aci_relation_from_external_epg_to_route_control_profile with id '%s'", data.Id.ValueString()))
}

func (r *L3extRsInstPToProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start update of resource: aci_relation_from_external_epg_to_route_control_profile")
	var data *L3extRsInstPToProfileResourceModel
	var stateData *L3extRsInstPToProfileResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Update of resource aci_relation_from_external_epg_to_route_control_profile with id '%s'", data.Id.ValueString()))

	var tagAnnotationPlan, tagAnnotationState []TagAnnotationL3extRsInstPToProfileResourceModel
	data.TagAnnotation.ElementsAs(ctx, &tagAnnotationPlan, false)
	stateData.TagAnnotation.ElementsAs(ctx, &tagAnnotationState, false)
	var tagTagPlan, tagTagState []TagTagL3extRsInstPToProfileResourceModel
	data.TagTag.ElementsAs(ctx, &tagTagPlan, false)
	stateData.TagTag.ElementsAs(ctx, &tagTagState, false)
	jsonPayload := getL3extRsInstPToProfileCreateJsonPayload(ctx, &resp.Diagnostics, false, data, tagAnnotationPlan, tagAnnotationState, tagTagPlan, tagTagState)

	if resp.Diagnostics.HasError() {
		return
	}

	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("api/mo/%s.json", data.Id.ValueString()), "POST", jsonPayload)

	if resp.Diagnostics.HasError() {
		return
	}

	getAndSetL3extRsInstPToProfileAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End update of resource aci_relation_from_external_epg_to_route_control_profile with id '%s'", data.Id.ValueString()))
}

func (r *L3extRsInstPToProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start delete of resource: aci_relation_from_external_epg_to_route_control_profile")
	var data *L3extRsInstPToProfileResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Delete of resource aci_relation_from_external_epg_to_route_control_profile with id '%s'", data.Id.ValueString()))
	jsonPayload := GetDeleteJsonPayload(ctx, &resp.Diagnostics, "l3extRsInstPToProfile", data.Id.ValueString())
	if resp.Diagnostics.HasError() {
		return
	}
	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("api/mo/%s.json", data.Id.ValueString()), "POST", jsonPayload)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("End delete of resource aci_relation_from_external_epg_to_route_control_profile with id '%s'", data.Id.ValueString()))
}

func (r *L3extRsInstPToProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start import state of resource: aci_relation_from_external_epg_to_route_control_profile")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	var stateData *L3extRsInstPToProfileResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &stateData)...)
	tflog.Debug(ctx, fmt.Sprintf("Import state of resource aci_relation_from_external_epg_to_route_control_profile with id '%s'", stateData.Id.ValueString()))

	tflog.Debug(ctx, "End import of state resource: aci_relation_from_external_epg_to_route_control_profile")
}

func getAndSetL3extRsInstPToProfileAttributes(ctx context.Context, diags *diag.Diagnostics, client *client.Client, data *L3extRsInstPToProfileResourceModel) {
	requestData := DoRestRequest(ctx, diags, client, fmt.Sprintf("api/mo/%s.json?rsp-subtree=full&rsp-subtree-class=%s", data.Id.ValueString(), "l3extRsInstPToProfile,tagAnnotation,tagTag"), "GET", nil)

	readData := getEmptyL3extRsInstPToProfileResourceModel()

	if diags.HasError() {
		return
	}
	if requestData.Search("imdata").Search("l3extRsInstPToProfile").Data() != nil {
		classReadInfo := requestData.Search("imdata").Search("l3extRsInstPToProfile").Data().([]interface{})
		if len(classReadInfo) == 1 {
			attributes := classReadInfo[0].(map[string]interface{})["attributes"].(map[string]interface{})
			for attributeName, attributeValue := range attributes {
				if attributeName == "dn" {
					readData.Id = basetypes.NewStringValue(attributeValue.(string))
					setL3extRsInstPToProfileParentDn(ctx, attributeValue.(string), readData)
				}
				if attributeName == "annotation" {
					readData.Annotation = basetypes.NewStringValue(attributeValue.(string))
				}
				if attributeName == "direction" {
					readData.Direction = basetypes.NewStringValue(attributeValue.(string))
				}
				if attributeName == "tnRtctrlProfileName" {
					readData.TnRtctrlProfileName = basetypes.NewStringValue(attributeValue.(string))
				}
			}
			TagAnnotationL3extRsInstPToProfileList := make([]TagAnnotationL3extRsInstPToProfileResourceModel, 0)
			TagTagL3extRsInstPToProfileList := make([]TagTagL3extRsInstPToProfileResourceModel, 0)
			_, ok := classReadInfo[0].(map[string]interface{})["children"]
			if ok {
				children := classReadInfo[0].(map[string]interface{})["children"].([]interface{})
				for _, child := range children {
					for childClassName, childClassDetails := range child.(map[string]interface{}) {
						childAttributes := childClassDetails.(map[string]interface{})["attributes"].(map[string]interface{})
						if childClassName == "tagAnnotation" {
							TagAnnotationL3extRsInstPToProfile := getEmptyTagAnnotationL3extRsInstPToProfileResourceModel()
							for childAttributeName, childAttributeValue := range childAttributes {
								if childAttributeName == "key" {
									TagAnnotationL3extRsInstPToProfile.Key = basetypes.NewStringValue(childAttributeValue.(string))
								}
								if childAttributeName == "value" {
									TagAnnotationL3extRsInstPToProfile.Value = basetypes.NewStringValue(childAttributeValue.(string))
								}

							}
							TagAnnotationL3extRsInstPToProfileList = append(TagAnnotationL3extRsInstPToProfileList, TagAnnotationL3extRsInstPToProfile)
						}
						if childClassName == "tagTag" {
							TagTagL3extRsInstPToProfile := getEmptyTagTagL3extRsInstPToProfileResourceModel()
							for childAttributeName, childAttributeValue := range childAttributes {
								if childAttributeName == "key" {
									TagTagL3extRsInstPToProfile.Key = basetypes.NewStringValue(childAttributeValue.(string))
								}
								if childAttributeName == "value" {
									TagTagL3extRsInstPToProfile.Value = basetypes.NewStringValue(childAttributeValue.(string))
								}

							}
							TagTagL3extRsInstPToProfileList = append(TagTagL3extRsInstPToProfileList, TagTagL3extRsInstPToProfile)
						}
					}
				}
			}
			tagAnnotationSet, _ := types.SetValueFrom(ctx, readData.TagAnnotation.ElementType(ctx), TagAnnotationL3extRsInstPToProfileList)
			readData.TagAnnotation = tagAnnotationSet
			tagTagSet, _ := types.SetValueFrom(ctx, readData.TagTag.ElementType(ctx), TagTagL3extRsInstPToProfileList)
			readData.TagTag = tagTagSet
		} else {
			diags.AddError(
				"too many results in response",
				fmt.Sprintf("%v matches returned for class 'l3extRsInstPToProfile'. Please report this issue to the provider developers.", len(classReadInfo)),
			)
		}
	} else {
		readData.Id = basetypes.NewStringNull()
	}
	*data = *readData
}

func getL3extRsInstPToProfileRn(ctx context.Context, data *L3extRsInstPToProfileResourceModel) string {
	return fmt.Sprintf("rsinstPToProfile-[%s]-%s", data.TnRtctrlProfileName.ValueString(), data.Direction.ValueString())
}

func setL3extRsInstPToProfileParentDn(ctx context.Context, dn string, data *L3extRsInstPToProfileResourceModel) {
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

func setL3extRsInstPToProfileId(ctx context.Context, data *L3extRsInstPToProfileResourceModel) {
	rn := getL3extRsInstPToProfileRn(ctx, data)
	data.Id = types.StringValue(fmt.Sprintf("%s/%s", data.ParentDn.ValueString(), rn))
}

func getL3extRsInstPToProfileTagAnnotationChildPayloads(ctx context.Context, diags *diag.Diagnostics, data *L3extRsInstPToProfileResourceModel, tagAnnotationL3extRsInstPToProfilePlan, tagAnnotationL3extRsInstPToProfileState []TagAnnotationL3extRsInstPToProfileResourceModel) []map[string]interface{} {
	childPayloads := []map[string]interface{}{}
	if !data.TagAnnotation.IsNull() && !data.TagAnnotation.IsUnknown() {
		tagAnnotationIdentifiers := []TagAnnotationIdentifier{}
		for _, tagAnnotationL3extRsInstPToProfile := range tagAnnotationL3extRsInstPToProfilePlan {
			childMap := NewAciObject()
			if !tagAnnotationL3extRsInstPToProfile.Key.IsNull() && !tagAnnotationL3extRsInstPToProfile.Key.IsUnknown() {
				childMap.Attributes["key"] = tagAnnotationL3extRsInstPToProfile.Key.ValueString()
			}
			if !tagAnnotationL3extRsInstPToProfile.Value.IsNull() && !tagAnnotationL3extRsInstPToProfile.Value.IsUnknown() {
				childMap.Attributes["value"] = tagAnnotationL3extRsInstPToProfile.Value.ValueString()
			}
			childPayloads = append(childPayloads, map[string]interface{}{"tagAnnotation": childMap})
			tagAnnotationIdentifier := TagAnnotationIdentifier{}
			tagAnnotationIdentifier.Key = tagAnnotationL3extRsInstPToProfile.Key
			tagAnnotationIdentifiers = append(tagAnnotationIdentifiers, tagAnnotationIdentifier)
		}
		for _, tagAnnotation := range tagAnnotationL3extRsInstPToProfileState {
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

func getL3extRsInstPToProfileTagTagChildPayloads(ctx context.Context, diags *diag.Diagnostics, data *L3extRsInstPToProfileResourceModel, tagTagL3extRsInstPToProfilePlan, tagTagL3extRsInstPToProfileState []TagTagL3extRsInstPToProfileResourceModel) []map[string]interface{} {
	childPayloads := []map[string]interface{}{}
	if !data.TagTag.IsNull() && !data.TagTag.IsUnknown() {
		tagTagIdentifiers := []TagTagIdentifier{}
		for _, tagTagL3extRsInstPToProfile := range tagTagL3extRsInstPToProfilePlan {
			childMap := NewAciObject()
			if !tagTagL3extRsInstPToProfile.Key.IsNull() && !tagTagL3extRsInstPToProfile.Key.IsUnknown() {
				childMap.Attributes["key"] = tagTagL3extRsInstPToProfile.Key.ValueString()
			}
			if !tagTagL3extRsInstPToProfile.Value.IsNull() && !tagTagL3extRsInstPToProfile.Value.IsUnknown() {
				childMap.Attributes["value"] = tagTagL3extRsInstPToProfile.Value.ValueString()
			}
			childPayloads = append(childPayloads, map[string]interface{}{"tagTag": childMap})
			tagTagIdentifier := TagTagIdentifier{}
			tagTagIdentifier.Key = tagTagL3extRsInstPToProfile.Key
			tagTagIdentifiers = append(tagTagIdentifiers, tagTagIdentifier)
		}
		for _, tagTag := range tagTagL3extRsInstPToProfileState {
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

func getL3extRsInstPToProfileCreateJsonPayload(ctx context.Context, diags *diag.Diagnostics, createType bool, data *L3extRsInstPToProfileResourceModel, tagAnnotationPlan, tagAnnotationState []TagAnnotationL3extRsInstPToProfileResourceModel, tagTagPlan, tagTagState []TagTagL3extRsInstPToProfileResourceModel) *container.Container {
	payloadMap := map[string]interface{}{}
	payloadMap["attributes"] = map[string]string{}

	if createType && !globalAllowExistingOnCreate {
		payloadMap["attributes"].(map[string]string)["status"] = "created"
	}
	childPayloads := []map[string]interface{}{}

	TagAnnotationchildPayloads := getL3extRsInstPToProfileTagAnnotationChildPayloads(ctx, diags, data, tagAnnotationPlan, tagAnnotationState)
	if TagAnnotationchildPayloads == nil {
		return nil
	}
	childPayloads = append(childPayloads, TagAnnotationchildPayloads...)

	TagTagchildPayloads := getL3extRsInstPToProfileTagTagChildPayloads(ctx, diags, data, tagTagPlan, tagTagState)
	if TagTagchildPayloads == nil {
		return nil
	}
	childPayloads = append(childPayloads, TagTagchildPayloads...)

	payloadMap["children"] = childPayloads
	if !data.Annotation.IsNull() && !data.Annotation.IsUnknown() {
		payloadMap["attributes"].(map[string]string)["annotation"] = data.Annotation.ValueString()
	}
	if !data.Direction.IsNull() && !data.Direction.IsUnknown() {
		payloadMap["attributes"].(map[string]string)["direction"] = data.Direction.ValueString()
	}
	if !data.TnRtctrlProfileName.IsNull() && !data.TnRtctrlProfileName.IsUnknown() {
		payloadMap["attributes"].(map[string]string)["tnRtctrlProfileName"] = data.TnRtctrlProfileName.ValueString()
	}
	payload, err := json.Marshal(map[string]interface{}{"l3extRsInstPToProfile": payloadMap})
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
