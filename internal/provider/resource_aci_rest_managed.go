package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
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
	Child      types.Set    `tfsdk:"child"`
	Annotation types.String `tfsdk:"annotation"`
	EscapeHtml types.Bool   `tfsdk:"escape_html"`
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

type ImportJsonString struct {
	ParentDn string   `json:"parentDn"`
	ChildRns []string `json:"childRns"`
}

// Angle brackets are not allowed within ACI class object identifier fields.
// Because of that "<,>" string was used to merge and split the list elements.
// Only using "," is not enough when the object identifier contains "," for example: "annotationKey-[~!$([])_+-={};:|,.]"
const ListElementConcatenationDelimiter = "<,>"

// List of attributes to be not stored in state
var IgnoreAttr = []string{"dn", "configQual", "configSt", "virtualIp", "annotation"}

// List of classes where 'rsp-prop-include=config-only' does not return the desired objects/properties
// var FullClasses = []string{"firmwareFwGrp", "maintMaintGrp", "maintMaintP", "firmwareFwP", "pkiExportEncryptionKey"}
var ConfigOnlyDns = []string{"uni/fabric/fwgrp-", "uni/fabric/maintgrp-", "uni/fabric/maintpol-", "uni/fabric/fwpol-", "uni/exportcryptkey"}

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

func (r *AciRestManagedResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	tflog.Debug(ctx, "Start plan modification of resource: aci_rest_managed")
	if !req.Plan.Raw.IsNull() {
		var planData, stateData *AciRestManagedResourceModel
		resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
		resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

		// modify the plan when annotation is not a property of class
		if ContainsString(NoAnnotationClasses, planData.ClassName.ValueString()) {
			planData.Annotation = basetypes.NewStringNull()
		} else if !planData.Content.IsNull() && !planData.Content.IsUnknown() && planData.Content.Elements()["annotation"] != nil {
			if !planData.Content.Elements()["annotation"].IsNull() {
				resp.Diagnostics.AddError(
					"Annotation not supported in content",
					"Annotation is not supported in content, please remove annotation from content and specify at resource level",
				)
			}
		}

		if stateData == nil && !globalAllowExistingOnCreate && !planData.Dn.IsUnknown() {
			var createCheckData *AciRestManagedResourceModel
			resp.Diagnostics.Append(req.Plan.Get(ctx, &createCheckData)...)
			CheckDn(ctx, &resp.Diagnostics, r.client, createCheckData.ClassName.ValueString(), createCheckData.Dn.ValueString())
			if resp.Diagnostics.HasError() {
				return
			}
		}

		resp.Plan.Set(ctx, planData)
	}
	tflog.Debug(ctx, "End plan modification of resource: aci_rest_managed")
}

func (r *AciRestManagedResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Manages ACI Model Objects via REST API calls. This resource can only manage a single API object and its direct children. It is able to read the state and therefore reconcile configuration drift.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The distinguished name (DN) of the object.",
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
				MarkdownDescription: "Which class object is being created, eg. fvTenant. (Make sure there is no colon in the classname)",
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
				Default:             stringdefault.StaticString(globalAnnotation),
				MarkdownDescription: `The annotation of the ACI object.`,
			},
			"escape_html": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable escaping of HTML characters when encoding the JSON payload.",
			},
		},
		Blocks: map[string]schema.Block{
			"child": schema.SetNestedBlock{
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
							Required:            true,
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

	data.Id = data.Dn

	tflog.Debug(ctx, fmt.Sprintf("Create of resource aci_rest_managed with id '%s'", data.Id.ValueString()))

	var childPlan, childState []ChildAciRestManagedResourceModel
	data.Child.ElementsAs(ctx, &childPlan, false)
	jsonPayload := getAciRestManagedCreateJsonPayload(ctx, &resp.Diagnostics, true, data, childPlan, childState)

	if resp.Diagnostics.HasError() {
		return
	}

	DoRestRequestEscapeHtml(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("api/mo/%s.json", data.Id.ValueString()), "POST", jsonPayload, data.EscapeHtml.ValueBool())
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End create of resource aci_rest_managed with id '%s'", data.Id.ValueString()))
}

func (r *AciRestManagedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start read of resource: aci_rest_managed")
	var data *AciRestManagedResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Read of resource aci_rest_managed with id '%s'", data.Id.ValueString()))

	// Retrive the import rn values from private state when set during import operation
	value, diags := req.Private.GetKey(ctx, "import_rn")
	resp.Diagnostics.Append(diags...)
	rn_values := []string{}
	if value != nil {
		rnMap := make(map[string]string, 0)
		err := json.Unmarshal(value, &rnMap)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to parse import_rn",
				fmt.Sprintf("Err: %s. Please report this issue to the provider developers.", err),
			)
			return
		}

		// This conversion is required because the SetKey function, utilized by the Terraform private state.
		// The SetKey function accepts byte strings as the value for the key.
		if strings.Contains(rnMap["rn_values"], ListElementConcatenationDelimiter) {
			rn_values = strings.Split(rnMap["rn_values"], ListElementConcatenationDelimiter)
		} else {
			rn_values = strings.Split(rnMap["rn_values"], ",")
		}
	}

	setAciRestManagedAttributes(ctx, &resp.Diagnostics, r.client, data, rn_values)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	if data.Id.IsNull() {
		var emptyData *AciRestManagedResourceModel
		resp.Diagnostics.Append(resp.State.Set(ctx, &emptyData)...)
	} else {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	}

	tflog.Debug(ctx, fmt.Sprintf("End read of resource aci_rest_managed with id '%s'", data.Id.ValueString()))
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

	data.Id = types.StringValue(data.Dn.ValueString())

	var childPlan, childState []ChildAciRestManagedResourceModel
	data.Child.ElementsAs(ctx, &childPlan, false)
	stateData.Child.ElementsAs(ctx, &childState, false)
	jsonPayload := getAciRestManagedCreateJsonPayload(ctx, &resp.Diagnostics, false, data, childPlan, childState)

	if resp.Diagnostics.HasError() {
		return
	}

	DoRestRequestEscapeHtml(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("api/mo/%s.json", data.Id.ValueString()), "POST", jsonPayload, data.EscapeHtml.ValueBool())
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End update of resource aci_rest_managed with id '%s'", data.Id.ValueString()))
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

	DoRestRequestEscapeHtml(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("api/mo/%s.json", data.Id.ValueString()), "POST", jsonPayload, data.EscapeHtml.ValueBool())
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("End delete of resource aci_rest_managed with id '%s'", data.Id.ValueString()))
}

func parseImportId(ctx context.Context, importId string) (string, string) {
	var importJson ImportJsonString
	err := json.Unmarshal([]byte(importId), &importJson)
	if err == nil {
		// JSON parsing successful
		// This conversion is required because the SetKey function, utilized by the Terraform private state.
		// The SetKey function accepts byte strings as the value for the key.
		return importJson.ParentDn, strings.Join(importJson.ChildRns, ListElementConcatenationDelimiter)
	}

	// JSON parsing failed, fall back to legacy logic
	tflog.Debug(ctx, "JSON parsing of the import ID failed. Falling back to legacy parsing logic.")
	return legacyParseImportId(ctx, importId)
}

func legacyParseImportId(ctx context.Context, importId string) (string, string) {
	var dn, rnValues string
	var idParts []string

	// Check for simple colon-separated format without brackets
	colonCount := strings.Count(importId, ":")
	if !strings.Contains(importId, "[") && colonCount > 0 && colonCount < 3 {
		tflog.Warn(ctx, "The use of the colon-separated format to import children for the resource is deprecated and will be removed in the next release.")
		tflog.Warn(ctx, "Please use the JSON format string to import children, instead of using a colon-separated import statement.")
		idParts = strings.Split(importId, ":")
	} else if strings.Contains(importId, "[") {
		// Custom splitting logic that respects brackets
		var openBrackets int
		var currentPartStart int

		for i, r := range importId {
			switch r {
			case '[':
				openBrackets++
			case ']':
				openBrackets--
			case ':':
				if openBrackets == 0 {
					idParts = append(idParts, importId[currentPartStart:i])
					currentPartStart = i + 1
				}
			}
		}
		// Append the last part after the loop
		if currentPartStart < len(importId) {
			idParts = append(idParts, importId[currentPartStart:])
		}
	} else {
		// When the import ID is a simple colon-delimited string (no brackets)
		return importId, ""
	}

	if len(idParts) == 1 {
		dn = idParts[0]
	} else if len(idParts) == 2 {
		if strings.HasPrefix(idParts[0], "uni/") {
			dn = idParts[0]
			rnValues = idParts[1]
		} else {
			dn = idParts[1]
		}
	} else if len(idParts) == 3 {
		dn = idParts[1]
		rnValues = idParts[2]
	}

	return dn, rnValues
}

func (r *AciRestManagedResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start import state of resource: aci_rest_managed")

	dn, rnValues := parseImportId(ctx, req.ID)
	if dn == "" && rnValues == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Please use the JSON format string to import children, instead of using a colon-separated import statement,... Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), dn)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("dn"), dn)...)

	// Set the import rn values in private state to be used in read
	if rnValues != "" {
		value := []byte(fmt.Sprintf(`{"rn_values": "%s"}`, rnValues))
		diags := resp.Private.SetKey(ctx, "import_rn", value)
		resp.Diagnostics.Append(diags...)
	}

	tflog.Debug(ctx, fmt.Sprintf("Import state of resource aci_annotation with id '%s'", dn))

	tflog.Debug(ctx, "End import of state resource: aci_rest_managed")
}

func setAciRestManagedAttributes(ctx context.Context, diags *diag.Diagnostics, client *client.Client, data *AciRestManagedResourceModel, rnValues []string) {

	// 'rsp-prop-include=config-only' does not return the desired objects/properties for certain rn
	// in that case ?rsp-prop-include=config-only should not be set
	var paramString string
	var match bool
	for _, configOnlyDn := range ConfigOnlyDns {
		match, _ = regexp.MatchString(fmt.Sprintf("%s[a-zA-Z0-9_.:-]*$[^/]*$", configOnlyDn), data.Dn.ValueString())
		if match {
			break
		}
	}

	if !data.Child.IsNull() && !data.Child.IsUnknown() || len(rnValues) > 0 {
		paramString = "?rsp-subtree=children"
	} else if !match {
		paramString = "?rsp-prop-include=config-only"
	}

	requestData := DoRestRequest(ctx, diags, client, fmt.Sprintf("api/mo/%s.json%s", data.Dn.ValueString(), paramString), "GET", nil)

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
		if val, ok := requestData.Search("imdata").Index(0).Data().(map[string]interface{})[className].(map[string]interface{})["attributes"]; ok {
			if dn, ok := val.(map[string]interface{})["dn"]; ok {
				data.Dn = basetypes.NewStringValue(dn.(string))
				data.Id = basetypes.NewStringValue(dn.(string))
			}
			if annotation, ok := val.(map[string]interface{})["annotation"]; ok {
				data.Annotation = basetypes.NewStringValue(annotation.(string))
			}
		}
		break
	}

	classData := requestData.Search("imdata").Search(data.ClassName.ValueString()).Data()

	var classAttributes map[string]interface{}
	var classChildren []interface{}
	if classData == nil {
		// plugin framework will error automatically when import is not found
		// Error: Cannot import non-existent remote object
		return
	} else if len(classData.([]interface{})) == 1 {
		if paramString == "?rsp-subtree=children" {
			requestDataConfigOnly := DoRestRequest(ctx, diags, client, fmt.Sprintf("api/mo/%s.json?rsp-prop-include=config-only", data.Dn.ValueString()), "GET", nil)
			if diags.HasError() {
				return
			}
			if requestDataConfigOnly.Search("imdata").Index(0).Data() == nil {
				data.Id = basetypes.NewStringNull()
				return
			}
			classDataConfigOnly := requestDataConfigOnly.Search("imdata").Search(data.ClassName.ValueString()).Data()
			if classDataConfigOnly == nil {
				// plugin framework will error automatically when import is not found
				// Error: Cannot import non-existent remote object
				return
			} else if len(classDataConfigOnly.([]interface{})) == 1 {
				classAttributes = classDataConfigOnly.([]interface{})[0].(map[string]interface{})["attributes"].(map[string]interface{})
				if classData.([]interface{})[0].(map[string]interface{})["children"] != nil {
					classChildren = classData.([]interface{})[0].(map[string]interface{})["children"].([]interface{})
				}
			} else {
				diags.AddError(
					"Too many results in response",
					fmt.Sprintf("%v matches returned for class '%s'. Please report this issue to the provider developers.", len(classDataConfigOnly.([]interface{})), data.ClassName),
				)
				return
			}
		} else {
			classAttributes = classData.([]interface{})[0].(map[string]interface{})["attributes"].(map[string]interface{})
			if classData.([]interface{})[0].(map[string]interface{})["children"] != nil {
				classChildren = classData.([]interface{})[0].(map[string]interface{})["children"].([]interface{})
			}
		}
	} else {
		diags.AddError(
			"Too many results in response",
			fmt.Sprintf("%v matches returned for class '%s'. Please report this issue to the provider developers.", len(classData.([]interface{})), data.ClassName),
		)
		return
	}

	content := make(map[string]attr.Value, 0)
	if !data.Content.IsNull() {
		for contentKey := range data.Content.Elements() {
			if val, ok := classAttributes[contentKey]; ok && !data.Content.Elements()[contentKey].IsNull() && !ContainsString(IgnoreAttr, contentKey) {
				content[contentKey] = basetypes.NewStringValue(val.(string))
			} else {
				content[contentKey] = basetypes.NewStringNull()
			}
		}
	} else {
		for attributeName, attributeValue := range classAttributes {
			if !ContainsString(IgnoreAttr, attributeName) {
				content[attributeName] = basetypes.NewStringValue(attributeValue.(string))
			}
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

		// continue if child is not found so it can be recreated in the update
		if childClassDetails == nil {
			continue
		}

		childRequestData := DoRestRequest(ctx, diags, client, fmt.Sprintf("api/mo/%s/%s.json?rsp-prop-include=config-only", data.Dn.ValueString(), child.Rn.ValueString()), "GET", nil)
		childClassData := childRequestData.Search("imdata").Search(aciRestManagedChild.ClassName.ValueString()).Data()
		if childClassData == nil {
			continue
		} else if len(childClassData.([]interface{})) != 1 {
			diags.AddError(
				"Too many results in response",
				fmt.Sprintf("%v matches returned for rn '%s'. Please report this issue to the provider developers.", len(childClassData.([]interface{})), child.Rn.ValueString()),
			)
			return
		}

		childContent := make(map[string]attr.Value, 0)
		for contentKey := range child.Content.Elements() {
			if val, ok := childClassData.([]interface{})[0].(map[string]interface{})["attributes"].(map[string]interface{})[contentKey]; ok && !child.Content.Elements()[contentKey].IsNull() {
				childContent[contentKey] = basetypes.NewStringValue(val.(string))
			} else {
				childContent[contentKey] = basetypes.NewStringNull()
			}
		}

		aciRestManagedChild.Content, _ = types.MapValue(types.StringType, childContent)
		foundChildren = append(foundChildren, aciRestManagedChild)
	}

	if len(rnValues) > 0 && len(foundChildren) == 0 {
		for _, rn := range rnValues {
			aciRestManagedChild := ChildAciRestManagedResourceModel{}
			aciRestManagedChild.Rn = basetypes.NewStringValue(rn)

			var childClassDetails map[string]interface{}

			for _, childClass := range classChildren {
				childClassDetails = findChildClass(childClass.(map[string]interface{}), &aciRestManagedChild)
				if childClassDetails != nil {
					break
				}
			}

			// err if child is not found because imported rn should be present to import
			if childClassDetails == nil {
				diags.AddError(
					"Import Failed",
					fmt.Sprintf("Unable to find specified child '%s'", rn),
				)
				return
			}

			childRequestData := DoRestRequest(ctx, diags, client, fmt.Sprintf("api/mo/%s/%s.json?rsp-prop-include=config-only", data.Dn.ValueString(), rn), "GET", nil)
			childClassData := childRequestData.Search("imdata").Search(aciRestManagedChild.ClassName.ValueString()).Data()
			if childClassData != nil {
				childContent := make(map[string]attr.Value, 0)
				if childClassDetails == nil {
					// this is not tested because it would involve importing a child corner case
					// deleting the child between the GET of parent with children and the config-only GET for child
					diags.AddError(
						"Import Failed",
						fmt.Sprintf("Unable to find specified child '%s'", rn),
					)
					return
				} else if len(childClassData.([]interface{})) != 1 {
					diags.AddError(
						"Too many results in response",
						fmt.Sprintf("%v matches returned for rn '%s'. Please report this issue to the provider developers.", len(childClassData.([]interface{})), rn))
					return
				}

				for attributeName, attributeValue := range childClassData.([]interface{})[0].(map[string]interface{})["attributes"].(map[string]interface{}) {
					if !ContainsString(IgnoreAttr, attributeName) || attributeName == "annotation" {
						childContent[attributeName] = basetypes.NewStringValue(attributeValue.(string))
					}
				}
				aciRestManagedChild.Content, _ = types.MapValue(types.StringType, childContent)
				foundChildren = append(foundChildren, aciRestManagedChild)
			}
		}
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

func getAciRestManagedChildPayloads(ctx context.Context, diags *diag.Diagnostics, data *AciRestManagedResourceModel, childPlan, childState []ChildAciRestManagedResourceModel) []interface{} {
	childPayloads := []interface{}{}
	if !data.Child.IsUnknown() && !data.Child.IsNull() {
		childIdentifiers := []AciRestManagedChildIdentifier{}
		for _, child := range childPlan {
			childMap := map[string]map[string]string{"attributes": {}}
			if !child.Rn.IsUnknown() {
				childMap["attributes"]["rn"] = child.Rn.ValueString()
			}
			if !data.Annotation.IsNull() && !data.Annotation.IsUnknown() && !ContainsString(NoAnnotationClasses, child.ClassName.ValueString()) {
				childMap["attributes"]["annotation"] = data.Annotation.ValueString()
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
				childPayloads = append(childPayloads, getAciRestManagedRemoveChildPayload(child.ClassName.ValueString(), child.Rn.ValueString()))
			}
		}
	} else if data.Child.IsNull() {
		for _, child := range childState {
			childPayloads = append(childPayloads, getAciRestManagedRemoveChildPayload(child.ClassName.ValueString(), child.Rn.ValueString()))
		}
	} else {
		tflog.Debug(ctx, fmt.Sprintf("Child with null state set to '%v', and unknown state set to '%v'", data.Child.IsNull(), data.Child.IsUnknown()))
		data.Child = types.SetNull(data.Child.ElementType(ctx))
	}

	return childPayloads
}

func getAciRestManagedRemoveChildPayload(className, rn string) map[string]map[string]map[string]interface{} {
	return map[string]map[string]map[string]interface{}{className: {"attributes": {"status": "deleted", "rn": rn}}}
}

func getAciRestManagedCreateJsonPayload(ctx context.Context, diags *diag.Diagnostics, createType bool, data *AciRestManagedResourceModel, childPlan, childState []ChildAciRestManagedResourceModel) *container.Container {
	payloadMap := map[string]interface{}{}
	payloadMap["attributes"] = map[string]string{}
	if createType && !globalAllowExistingOnCreate {
		payloadMap["attributes"].(map[string]string)["status"] = "created"
	}

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
				if k == "annotation" {
					data.Annotation = v.(basetypes.StringValue)
				}
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
