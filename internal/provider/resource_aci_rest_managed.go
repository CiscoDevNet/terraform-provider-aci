package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AciRestManagedResource{}
var _ resource.ResourceWithImportState = &AciRestManagedResource{}

func NewAciRestManagedResource() resource.Resource {
	return &AciRestManagedResource{}
}

// AciRestManagedResource defines the resource implementation.
type AciRestManagedResource struct {
	client *client.Client
}

// AciRestManagedResourceModel describes the resource data model.
type AciRestManagedResourceModel struct {
	Id         types.String `tfsdk:"id"`
	Dn         types.String `tfsdk:"dn"`
	ClassName  types.String `tfsdk:"class_name"`
	Content    types.Map    `tfsdk:"content"`
	Child      types.Set    `tfsdk:"child"`
	Annotation types.String `tfsdk:"annotation"`
}

// ChildAciRestManagedResourceModel describes the resource data model for the children without relationships.
type ChildAciRestManagedResourceModel struct {
	Rn        types.String `tfsdk:"rn"`
	ClassName types.String `tfsdk:"class_name"`
	Content   types.Map    `tfsdk:"content"`
}

type AciRestManagedChildIdentifier struct {
	Rn        types.String
	ClassName types.String
}

// List of classes where 'rsp-prop-include=config-only' does not return the desired objects/properties
var FullClasses = []string{"firmwareFwGrp", "maintMaintGrp", "maintMaintP", "firmwareFwP", "pkiExportEncryptionKey"}

// List of classes where an immediate GET following a POST might not reflect the created/updated object
var AllowEmptyReadClasses = []string{"firmwareFwGrp", "firmwareRsFwgrpp", "firmwareFwP", "fabricNodeBlk"}

// List of classes which do not support annotations
var NoAnnotationClasses = UnsupportedAnnotationClasses()

var UnableToDelete = "unable to delete"

func (r *AciRestManagedResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Debug(ctx, "Start schema of resource: aci_rest_managed")
	resp.TypeName = req.ProviderTypeName + "_rest_managed"
	tflog.Debug(ctx, "End schema of resource: aci_rest_managed")
}

func (r *AciRestManagedResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Manages ACI Model Objects via REST API calls. This resource can only manage a single API object and its direct children. It is able to read the state and therefore reconcile configuration drift.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The distinquised name (DN) of the object.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dn": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Distinguished name of object being managed including its relative name, e.g. uni/tn-EXAMPLE_TENANT.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"class_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Which class object is being created. (Make sure there is no colon in the classname)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content": schema.MapAttribute{
				MarkdownDescription: "Map of key-value pairs those needed to be passed to the Model object as parameters. Make sure the key name matches the name with the object parameter in ACI.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"annotation": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: `The annotation of the ACI object.`,
			},
		},
		Blocks: map[string]schema.Block{
			"child": schema.SetNestedBlock{
				//Optional:            true,
				MarkdownDescription: "List of children.",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"rn": schema.StringAttribute{
							MarkdownDescription: "The relative name of the child object.",
							Required:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"class_name": schema.StringAttribute{
							MarkdownDescription: "Class name of child object.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"content": schema.MapAttribute{
							MarkdownDescription: "Map of key-value pairs which represents the attributes for the child object.",
							Optional:            true,
							Computed:            true,
							ElementType:         types.StringType,
							PlanModifiers: []planmodifier.Map{
								mapplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
		},
	}
}

func (r *AciRestManagedResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of resource: aci_rest_managed")
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
	tflog.Debug(ctx, "End configure of resource: aci_rest_managed")
}

func (r *AciRestManagedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start create of resource: aci_rest_managed")

	var data *AciRestManagedResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	setAciRestManagedProperties(data)

	tflog.Debug(ctx, fmt.Sprintf("create of resource aci_rest_managed with id '%s'", data.Id.ValueString()))

	var childPlan, childState []ChildAciRestManagedResourceModel
	data.Child.ElementsAs(ctx, &childPlan, false)
	jsonPayload := getAciRestManagedCreateJsonPayload(ctx, &resp.Diagnostics, data, childPlan, childState)

	if resp.Diagnostics.HasError() {
		return
	}

	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("api/mo/%s.json", data.Id.ValueString()), "POST", jsonPayload)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, "End create of resource: aci_rest_managed")
}

func (r *AciRestManagedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start read of resource: aci_rest_managed")
	var data *AciRestManagedResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("read of resource aci_rest_managed with id '%s'", data.Id.ValueString()))

	setAciRestManagedAttributes(ctx, &resp.Diagnostics, r.client, data)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("AKINI DATA: '%s'", data))

	// Save updated data into Terraform state
	if data.Id.IsNull() {
		var emptyData *AciRestManagedResourceModel
		resp.Diagnostics.Append(resp.State.Set(ctx, &emptyData)...)
	} else {
		tflog.Debug(ctx, fmt.Sprintf("AKINI DATA NOT NULL: '%s'", data))
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	}

	tflog.Debug(ctx, "End read of resource: aci_rest_managed")
}

func (r *AciRestManagedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start update of resource: aci_rest_managed")
	var data *AciRestManagedResourceModel
	var stateData *AciRestManagedResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Update of resource aci_rest_managed with id '%s'", data.Id.ValueString()))

	setAciRestManagedProperties(data)

	var childPlan, childState []ChildAciRestManagedResourceModel
	data.Child.ElementsAs(ctx, &childPlan, false)
	stateData.Child.ElementsAs(ctx, &childState, false)
	jsonPayload := getAciRestManagedCreateJsonPayload(ctx, &resp.Diagnostics, data, childPlan, childState)

	if resp.Diagnostics.HasError() {
		return
	}

	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("api/mo/%s.json", data.Id.ValueString()), "POST", jsonPayload)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, "End update of resource: aci_rest_managed")
}

func (r *AciRestManagedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start delete of resource: aci_rest_managed")
	var data *AciRestManagedResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Delete of resource aci_rest_managed with id '%s'", data.Id.ValueString()))
	jsonPayload := GetDeleteJsonPayload(ctx, &resp.Diagnostics, data.ClassName.ValueString(), data.Id.ValueString())
	if resp.Diagnostics.HasError() {
		return
	}

	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("api/mo/%s.json", data.Id.ValueString()), "POST", jsonPayload)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "End delete of resource: aci_rest_managed")
}

func (r *AciRestManagedResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start import state of resource: aci_rest_managed")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	var stateData *AciRestManagedResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &stateData)...)
	tflog.Debug(ctx, fmt.Sprintf("Import state of resource aci_annotation with id '%s'", stateData.Id.ValueString()))

	tflog.Debug(ctx, "End import of state resource: aci_rest_managed")
}

func setAciRestManagedAttributes(ctx context.Context, diags *diag.Diagnostics, client *client.Client, data *AciRestManagedResourceModel) {
	requestData := DoRestRequest(ctx, diags, client, fmt.Sprintf("api/mo/%s.json?rsp-subtree=children", data.Id.ValueString()), "GET", nil)

	if diags.HasError() {
		return
	}

	if requestData.Search("imdata").Index(0).Data() == nil {
		data.Id = basetypes.NewStringNull()
		return
	}

	for className := range requestData.Search("imdata").Index(0).Data().(map[string]interface{}) {
		tflog.Debug(ctx, fmt.Sprintf("Setting ClassName to %s", className))
		data.ClassName = basetypes.NewStringValue(className)
		break
	}

	classData := requestData.Search("imdata").Search(data.ClassName.ValueString()).Data().([]interface{})

	var classAttributes map[string]interface{}
	var classChildren []interface{}
	if len(classData) == 1 {
		classAttributes = classData[0].(map[string]interface{})["attributes"].(map[string]interface{})
		if classData[0].(map[string]interface{})["children"] != nil {
			classChildren = classData[0].(map[string]interface{})["children"].([]interface{})
		}
	} else {
		diags.AddError(
			"Too many results in response",
			fmt.Sprintf("%v matches returned for class '%s'. Please report this issue to the provider developers.", len(classData), data.ClassName),
		)
		return
	}

	content := make(map[string]attr.Value, 0)
	for contentKey := range data.Content.Elements() {
		if val, ok := classAttributes[contentKey]; ok && !data.Content.Elements()[contentKey].IsNull() {
			content[contentKey] = basetypes.NewStringValue(val.(string))
		} else {
			content[contentKey] = basetypes.NewStringNull()
		}
	}
	data.Content, _ = types.MapValue(types.StringType, content)

	var children, foundChildren []ChildAciRestManagedResourceModel
	data.Child.ElementsAs(ctx, &children, false)

	for _, child := range children {
		aciRestManagedChild := ChildAciRestManagedResourceModel{}
		aciRestManagedChild.Rn = child.Rn

		var childClassDetails map[string]interface{}

		for _, childClass := range classChildren {
			childClassDetails = findChildClass(childClass.(map[string]interface{}), &aciRestManagedChild)
			if childClassDetails != nil {
				break
			}
		}

		childContent := make(map[string]attr.Value, 0)
		for contentKey := range child.Content.Elements() {
			if val, ok := childClassDetails["attributes"].(map[string]interface{})[contentKey]; ok && !child.Content.Elements()[contentKey].IsNull() {
				childContent[contentKey] = basetypes.NewStringValue(val.(string))
			} else {
				childContent[contentKey] = basetypes.NewStringNull()
			}
		}
		aciRestManagedChild.Content, _ = types.MapValue(types.StringType, childContent)
		foundChildren = append(foundChildren, aciRestManagedChild)
	}

	if len(foundChildren) > 0 {
		data.Child, _ = types.SetValueFrom(ctx, data.Child.ElementType(ctx), foundChildren)
	} else {
		data.Child = types.SetNull(data.Child.ElementType(ctx))
	}

}

func findChildClass(childClass map[string]interface{}, child *ChildAciRestManagedResourceModel) map[string]interface{} {

	for className, classDetails := range childClass {
		classRn := classDetails.(map[string]interface{})["attributes"].(map[string]interface{})["rn"].(string)
		if child.Rn.ValueString() == classRn {
			child.ClassName = basetypes.NewStringValue(className)
			return classDetails.(map[string]interface{})
		}
	}
	return nil
}

func setAciRestManagedProperties(data *AciRestManagedResourceModel) {
	// Set Id
	data.Id = types.StringValue(data.Dn.ValueString())
	// Remove annotation when unsupported
	if ContainsString(NoAnnotationClasses, data.ClassName.ValueString()) {
		data.Annotation = types.StringNull()
		// Add default annotation is not set
	} else if data.Annotation.IsNull() || data.Annotation.IsUnknown() {
		data.Annotation = types.StringValue(globalAnnotation)
	}
}

func getAciRestManagedChildPayloads(ctx context.Context, diags *diag.Diagnostics, data *AciRestManagedResourceModel, childPlan, childState []ChildAciRestManagedResourceModel) []interface{} {
	childPayloads := []interface{}{}
	if !data.Child.IsUnknown() && !data.Child.IsNull() {
		childIdentifiers := []AciRestManagedChildIdentifier{}
		for _, child := range childPlan {
			childMap := map[string]map[string]string{"attributes": {}}
			if !child.Rn.IsUnknown() {
				childMap["attributes"]["rn"] = child.Rn.ValueString()
			}
			if !child.Content.IsNull() && !child.Content.IsUnknown() {
				for k, v := range child.Content.Elements() {
					if !v.(basetypes.StringValue).IsNull() && !v.(basetypes.StringValue).IsUnknown() {
						childMap["attributes"][k] = v.(basetypes.StringValue).ValueString()
					}
				}
			}
			childPayloads = append(childPayloads, map[string]interface{}{child.ClassName.ValueString(): childMap})
			childIdentifier := AciRestManagedChildIdentifier{}
			childIdentifier.Rn = child.Rn
			childIdentifier.ClassName = child.ClassName
			childIdentifiers = append(childIdentifiers, childIdentifier)
		}
		for _, child := range childState {
			delete := true
			for _, childIdentifier := range childIdentifiers {
				if childIdentifier.Rn == child.Rn && childIdentifier.ClassName == child.ClassName {
					delete = false
					break
				}
			}
			if delete {
				// TODO check if this works
				// Changing the RN or ClassName should delete previous objects
				childMap := map[string]map[string]interface{}{"attributes": {}}
				childMap["attributes"]["status"] = "deleted"
				childMap["attributes"]["rn"] = child.Rn.ValueString()
				childPayloads = append(childPayloads, map[string]interface{}{child.ClassName.ValueString(): childMap})
			}
		}
	} else {
		data.Child = types.SetNull(data.Child.ElementType(ctx))
	}

	return childPayloads
}

func getAciRestManagedCreateJsonPayload(ctx context.Context, diags *diag.Diagnostics, data *AciRestManagedResourceModel, childPlan, childState []ChildAciRestManagedResourceModel) *container.Container {
	payloadMap := map[string]interface{}{}
	payloadMap["attributes"] = map[string]string{}

	childPayloads := getAciRestManagedChildPayloads(ctx, diags, data, childPlan, childState)
	if diags.HasError() {
		return nil
	}

	payloadMap["children"] = childPayloads

	if !data.Annotation.IsNull() && !data.Annotation.IsUnknown() {
		payloadMap["attributes"].(map[string]string)["annotation"] = data.Annotation.ValueString()
	}
	if !data.Content.IsNull() && !data.Content.IsUnknown() {
		for k, v := range data.Content.Elements() {
			if !v.(basetypes.StringValue).IsNull() && !v.(basetypes.StringValue).IsUnknown() {
				payloadMap["attributes"].(map[string]string)[k] = v.(basetypes.StringValue).ValueString()
			}
		}
	} else {
		data.Content = types.MapNull(data.Content.ElementType(ctx))
	}

	payload, err := json.Marshal(map[string]interface{}{data.ClassName.ValueString(): payloadMap})

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
