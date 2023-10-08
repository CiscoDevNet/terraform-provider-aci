package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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
	//Child      types.Set    `tfsdk:"child"`
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
var NoAnnotationClasses = []string{"tagTag"}

func (r *AciRestManagedResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Trace(ctx, "start schema of resource: aci_rest_managed")
	resp.TypeName = req.ProviderTypeName + "_rest_managed"
	tflog.Trace(ctx, "end schema of resource: aci_rest_managed")
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
				// TODO DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				// 	content := d.Get("content")
				// 	contentStrMap := toStrMap(content.(map[string]interface{}))
				// 	key := strings.Split(k, ".")[1]
				// 	if _, ok := contentStrMap[key]; ok {
				// 		return false
				// 	}
				// 	return true
				// },
			},
			"annotation": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Default:             stringdefault.StaticString(globalAnnotation),
				MarkdownDescription: `The annotation of the ACI object.`,
			},
		},
		// TODO Fix the support child blocks
		// Blocks: map[string]schema.Block{
		// 	"child": schema.SetNestedBlock{
		// 		//Optional:            true,
		// 		MarkdownDescription: "List of children.",
		// 		PlanModifiers: []planmodifier.Set{
		// 			setplanmodifier.UseStateForUnknown(),
		// 		},
		// 		NestedObject: schema.NestedBlockObject{
		// 			Attributes: map[string]schema.Attribute{
		// 				"rn": schema.StringAttribute{
		// 					MarkdownDescription: "The relative name of the child object.",
		// 					Required:            true,
		// 					PlanModifiers: []planmodifier.String{
		// 						stringplanmodifier.UseStateForUnknown(),
		// 					},
		// 				},
		// 				"class_name": schema.StringAttribute{
		// 					MarkdownDescription: "Class name of child object.",
		// 					Optional:            true,
		// 					Computed:            true,
		// 					PlanModifiers: []planmodifier.String{
		// 						stringplanmodifier.UseStateForUnknown(),
		// 					},
		// 				},
		// 				"content": schema.MapAttribute{
		// 					MarkdownDescription: "Map of key-value pairs which represents the attributes for the child object.",
		// 					Optional:            true,
		// 					Computed:            true,
		// 					ElementType:         types.StringType,
		// 					PlanModifiers: []planmodifier.Map{
		// 						mapplanmodifier.UseStateForUnknown(),
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// },
	}
}

func (r *AciRestManagedResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Trace(ctx, "start configure of resource: aci_rest_managed")
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
	tflog.Trace(ctx, "end configure of resource: aci_rest_managed")
}

func (r *AciRestManagedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Trace(ctx, "start create of resource: aci_rest_managed")
	// On create retrieve information on current state prior to making any changes in order to determine child delete operations
	var stateData *AciRestManagedResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &stateData)...)
	setAciRestManagedId(ctx, stateData)
	messageMap := setAciRestManagedAttributes(ctx, r.client, stateData)
	if messageMap != nil {
		resp.Diagnostics.AddError(messageMap.(map[string]string)["message"], messageMap.(map[string]string)["messageDetail"])
	}

	var data *AciRestManagedResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	setAciRestManagedId(ctx, data)

	tflog.Trace(ctx, fmt.Sprintf("create of resource aci_rest_managed with id '%s'", data.Id.ValueString()))

	var childPlan, childState []ChildAciRestManagedResourceModel
	//data.Child.ElementsAs(ctx, &childPlan, false)
	//stateData.Child.ElementsAs(ctx, &childState, false)
	jsonPayload, message, messageDetail := getAciRestManagedCreateJsonPayload(ctx, data, childPlan, childState)

	if jsonPayload == nil {
		resp.Diagnostics.AddError(message, messageDetail)
		return
	}

	requestData, message, messageDetail := doAciRestManagedRequest(ctx, r.client, fmt.Sprintf("api/mo/%s.json", data.Id.ValueString()), "POST", jsonPayload)
	if requestData == nil {
		resp.Diagnostics.AddError(message, messageDetail)
		return
	}

	messageMap = setAciRestManagedAttributes(ctx, r.client, data)
	if messageMap != nil {
		resp.Diagnostics.AddError(messageMap.(map[string]string)["message"], messageMap.(map[string]string)["messageDetail"])
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "end create of resource: aci_rest_managed")
}

func (r *AciRestManagedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Trace(ctx, "start read of resource: aci_rest_managed")
	var data *AciRestManagedResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("read of resource aci_rest_managed with id '%s'", data.Id.ValueString()))

	messageMap := setAciRestManagedAttributes(ctx, r.client, data)
	if messageMap != nil {
		resp.Diagnostics.AddError(messageMap.(map[string]string)["message"], messageMap.(map[string]string)["messageDetail"])
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "end read of resource: aci_rest_managed")
}

func (r *AciRestManagedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Trace(ctx, "start update of resource: aci_rest_managed")
	var data *AciRestManagedResourceModel
	var stateData *AciRestManagedResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("update of resource aci_rest_managed with id '%s'", data.Id.ValueString()))

	var childPlan, childState []ChildAciRestManagedResourceModel
	//data.Child.ElementsAs(ctx, &childPlan, false)
	//stateData.Child.ElementsAs(ctx, &childState, false)
	jsonPayload, message, messageDetail := getAciRestManagedCreateJsonPayload(ctx, data, childPlan, childState)

	if jsonPayload == nil {
		resp.Diagnostics.AddError(message, messageDetail)
		return
	}

	requestData, message, messageDetail := doAciRestManagedRequest(ctx, r.client, fmt.Sprintf("api/mo/%s.json", data.Id.ValueString()), "POST", jsonPayload)
	if requestData == nil {
		resp.Diagnostics.AddError(message, messageDetail)
		return
	}

	messageMap := setAciRestManagedAttributes(ctx, r.client, data)
	if messageMap != nil {
		resp.Diagnostics.AddError(messageMap.(map[string]string)["message"], messageMap.(map[string]string)["messageDetail"])
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "end update of resource: aci_rest_managed")
}

func (r *AciRestManagedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "start delete of resource: aci_rest_managed")
	var data *AciRestManagedResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("delete of resource aci_rest_managed with id '%s'", data.Id.ValueString()))
	jsonPayload, message, messageDetail := getAciRestManagedDeleteJsonPayload(ctx, data)
	if jsonPayload == nil {
		resp.Diagnostics.AddError(message, messageDetail)
		return
	}
	requestData, message, messageDetail := doAciRestManagedRequest(ctx, r.client, fmt.Sprintf("api/mo/%s.json", data.Id.ValueString()), "POST", jsonPayload)
	if requestData == nil {
		resp.Diagnostics.AddError(message, messageDetail)
		return
	}
	tflog.Trace(ctx, "end delete of resource: aci_rest_managed")
}

func (r *AciRestManagedResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func setAciRestManagedAttributes(ctx context.Context, client *client.Client, data *AciRestManagedResourceModel) interface{} {
	requestData, message, messageDetail := doAciRestManagedRequest(ctx, client, fmt.Sprintf("api/mo/%s.json?rsp-subtree=children", data.Id.ValueString()), "GET", nil)

	if requestData == nil {
		return map[string]string{"message": message, "messageDetail": messageDetail}
	}

	if requestData.Search("imdata").Index(0).Data() == nil {
		return nil
	}

	classData := requestData.Search("imdata").Index(0).Data().(map[string]interface{})
	for className := range classData {
		tflog.Trace(ctx, fmt.Sprintf("Setting ClassName to %s", className))
		data.ClassName = basetypes.NewStringValue(className)
		break
	}

	// Only attributes set in the content should be saved into state
	contentKeys := make([]string, 0, 5)
	for k := range data.Content.Elements(){
		contentKeys = append(contentKeys, k)
	}

	if requestData.Search("imdata").Search(data.ClassName.ValueString()).Data() != nil {
		classReadInfo := requestData.Search("imdata").Search(data.ClassName.ValueString()).Data().([]interface{})
		if len(classReadInfo) == 1 {
			attributes := classReadInfo[0].(map[string]interface{})["attributes"].(map[string]interface{})
			contents := map[string]attr.Value{}
			for attributeName, attributeValue := range attributes {
				if attributeName == "dn" {
					dn := attributeValue.(string)
					data.Id = basetypes.NewStringValue(dn)
					data.Dn = basetypes.NewStringValue(dn)
				} else if attributeName == "annotation" {
					data.Annotation = basetypes.NewStringValue(attributeValue.(string))
				} else if containsString(contentKeys, attributeName) {
					contents[attributeName] = basetypes.NewStringValue(attributeValue.(string))
				}
			}
			data.Content, _ = types.MapValue(types.StringType, contents)
			// ChildAciRestManagedResourceList := make([]ChildAciRestManagedResourceModel, 0)
			// _, ok := classReadInfo[0].(map[string]interface{})["children"]
			// if ok {
			// 	children := classReadInfo[0].(map[string]interface{})["children"].([]interface{})
			// 	for _, child := range children {
			// 		for childClassName, childClassDetails := range child.(map[string]interface{}) {
			// 			childAttributes := childClassDetails.(map[string]interface{})["attributes"].(map[string]interface{})
			// 			ChildAciRestManaged := ChildAciRestManagedResourceModel{}
			// 			ChildAciRestManaged.ClassName = basetypes.NewStringValue(childClassName)
			// 			for childAttributeName, childAttributeValue := range childAttributes {
			// 				if childAttributeName == "rn" {
			// 					ChildAciRestManaged.Rn = basetypes.NewStringValue(childAttributeValue.(string))
			// 				}
			// 				if childAttributeName == "content" {
			// 					// TODO set child content map
			// 					//ChildAciRestManaged.Content = basetypes.NewStringValue(childAttributeValue.(string))
			// 				}
			// 			}
			// 			ChildAciRestManagedResourceList = append(ChildAciRestManagedResourceList, ChildAciRestManaged)
			// 		}
			// 	}
			// }
			// if len(ChildAciRestManagedResourceList) > 0 {
			// 	tflog.Trace(ctx, "Setting Child Set Data")
			// 	childSet, _ := types.SetValueFrom(ctx, data.Child.ElementType(ctx), ChildAciRestManagedResourceList)
			// 	data.Child = childSet
			// }
		} else {
			return map[string]string{
				"message":       "too many results in response",
				"messageDetail": fmt.Sprintf("%v matches returned for class 'l3extConsLbl'. Please report this issue to the provider developers.", len(classReadInfo)),
			}
		}
	}
	return nil
}

func setAciRestManagedId(ctx context.Context, data *AciRestManagedResourceModel) {
	data.Id = types.StringValue(data.Dn.ValueString())
}

// TODO This needs more attention
// Child payloads may not be in the correct format
// func getAciRestManagedChildPayloads(ctx context.Context, data *AciRestManagedResourceModel, childPlan, childState []ChildAciRestManagedResourceModel) ([]map[string]interface{}, string, string) {
// 	childPayloads := []map[string]interface{}{}
// 	if !data.Child.IsUnknown() {
// 		childIdentifiers := []AciRestManagedChildIdentifier{}
// 		for _, child := range childPlan {
// 			childMap := map[string]map[string]interface{}{"attributes": {}}
// 			if !child.Rn.IsUnknown() {
// 				childMap["attributes"]["rn"] = child.Rn.ValueString()
// 			}
// 			if !child.ClassName.IsUnknown() {
// 				childMap["attributes"]["class_name"] = child.ClassName.ValueString()
// 			}
// 			if !child.Content.IsUnknown() {
// 				// TODO This will need to be fixed
// 				// Because child contents should be the actual attributes.
// 				childMap["attributes"]["content"] = child.Content.Elements()
// 			}
// 			childPayloads = append(childPayloads, map[string]interface{}{"children": childMap})
// 			childIdentifier := AciRestManagedChildIdentifier{}
// 			childIdentifier.Rn = child.Rn
// 			childIdentifier.ClassName = child.ClassName
// 			childIdentifiers = append(childIdentifiers, childIdentifier)
// 		}
// 		for _, child := range childState {
// 			delete := true
// 			for _, childIdentifier := range childIdentifiers {
// 				if childIdentifier.Rn == child.Rn && childIdentifier.ClassName == child.ClassName {
// 					delete = false
// 					break
// 				}
// 			}
// 			if delete {
// 				// TODO check if this works
// 				// Changing the RN or ClassName should delete previous objects
// 				childMap := map[string]map[string]interface{}{"attributes": {}}
// 				childMap["attributes"]["status"] = "deleted"
// 				childMap["attributes"]["rn"] = child.Rn.ValueString()
// 				childPayloads = append(childPayloads, map[string]interface{}{"children": childMap})
// 			}
// 		}
// 	} else {
// 		data.Child = types.SetNull(data.Child.ElementType(ctx))
// 	}

// 	return childPayloads, "", ""
// }

func getAciRestManagedCreateJsonPayload(ctx context.Context, data *AciRestManagedResourceModel, childPlan, childState []ChildAciRestManagedResourceModel) (*container.Container, string, string) {
	payloadMap := map[string]interface{}{}
	payloadMap["attributes"] = map[string]string{}

	// childPayloads := []map[string]interface{}{}
	// ChildPayloads, errMessage, errMessageDetail := getAciRestManagedChildPayloads(ctx, data, childPlan, childState)
	// if ChildPayloads == nil {
	// 	return nil, errMessage, errMessageDetail
	// }
	// childPayloads = append(childPayloads, ChildPayloads...)
	// payloadMap["children"] = childPayloads

	if !data.Annotation.IsNull() && !data.Annotation.IsUnknown() {
		payloadMap["attributes"].(map[string]string)["annotation"] = data.Annotation.ValueString()
	}
	if !data.Content.IsNull() && !data.Content.IsUnknown() {
		for k, v := range data.Content.Elements() {
			payloadMap["attributes"].(map[string]string)[k] = v.(basetypes.StringValue).ValueString()
		}
	}

	payload, err := json.Marshal(map[string]interface{}{data.ClassName.ValueString(): payloadMap})
	if err != nil {
		errMessage := "marshalling of json payload failed"
		errMessageDetail := fmt.Sprintf("err: %s. Please report this issue to the provider developers.", err)
		return nil, errMessage, errMessageDetail
	}

	jsonPayload, err := container.ParseJSON(payload)

	if err != nil {
		errMessage := "construction of json payload failed"
		errMessageDetail := fmt.Sprintf("err: %s. Please report this issue to the provider developers.", err)
		return nil, errMessage, errMessageDetail
	}
	return jsonPayload, "", ""
}

func getAciRestManagedDeleteJsonPayload(ctx context.Context, data *AciRestManagedResourceModel) (*container.Container, string, string) {

	jsonString := fmt.Sprintf(`{"%s":{"attributes":{"dn": "%s","status": "deleted"}}}`, data.ClassName.ValueString(), data.Id.ValueString())
	jsonPayload, err := container.ParseJSON([]byte(jsonString))
	if err != nil {
		errMessage := "construction of json payload failed"
		errMessageDetail := fmt.Sprintf("err: %s. Please report this issue to the provider developers.", err)
		return nil, errMessage, errMessageDetail
	}
	return jsonPayload, "", ""
}

// TODO make this a generic function in the generator.
func doAciRestManagedRequest(ctx context.Context, client *client.Client, path, method string, payload *container.Container) (*container.Container, string, string) {

	restRequest, err := client.MakeRestRequest(method, path, payload, true)
	if err != nil {
		message := fmt.Sprintf("creation of %s rest request failed", strings.ToLower(method))
		messageDetail := fmt.Sprintf("Err: %s. Please report this issue to the provider developers.", err)
		return nil, message, messageDetail
	}

	cont, restResponse, err := client.Do(restRequest)

	if restResponse != nil && restResponse.StatusCode != 200 {
		message := fmt.Sprintf("%s rest request failed", strings.ToLower(method))
		messageDetail := fmt.Sprintf("Response: %s, err: %s. Please report this issue to the provider developers.", cont.Data().(map[string]interface{})["imdata"], err)
		return nil, message, messageDetail
	} else if err != nil {
		message := fmt.Sprintf("%s rest request failed", strings.ToLower(method))
		messageDetail := fmt.Sprintf("Err: %s. Please report this issue to the provider developers.", err)
		return nil, message, messageDetail
	}

	return cont, "", ""
}

func containsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
