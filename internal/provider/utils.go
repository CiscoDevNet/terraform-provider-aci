package provider

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func getChildClassesForGetRequest(childClasses []string) []string {
	classVersions := classVersions()
	subtreeClasses := []string{}
	for _, childClass := range childClasses {
		childVersion := classVersions[childClass]
		inside, err := CompareVersionsRange(apicVersion, childVersion, "inside")
		if err == nil && inside && !slices.Contains(subtreeClasses, childClass) {
			subtreeClasses = append(subtreeClasses, childClass)
		}
	}
	return subtreeClasses
}

func IsEmptySingleNestedAttribute(attributes map[string]attr.Value) bool {
	for _, value := range attributes {
		if !value.IsNull() {
			return false
		}
	}
	return true
}

func SingleNestedAttributeRequiredAttributesNotProvided(attributes map[string]attr.Value, requiredAttributes []string) bool {
	for _, requiredAttribute := range requiredAttributes {
		if attributes[requiredAttribute].IsNull() {
			return true
		}
	}
	return false
}

type AciObject struct {
	Attributes map[string]interface{}   `json:"attributes"`
	Children   []map[string]interface{} `json:"children"`
}

func NewAciObject() AciObject {
	return AciObject{
		Attributes: make(map[string]interface{}),
		Children:   []map[string]interface{}{},
	}
}

func ContainsString(strings []string, matchString string) bool {
	for _, stringValue := range strings {
		if stringValue == matchString {
			return true
		}
	}
	return false
}

func GetMOName(dn string) string {
	splittedDn := strings.Split(dn, "/")
	if len(splittedDn) > 1 {
		return strings.Join(strings.Split(splittedDn[len(splittedDn)-1], "-")[1:], "-")
	}
	return splittedDn[0]
}

func CheckDn(ctx context.Context, diags *diag.Diagnostics, client *client.Client, classname, dn string) {
	requestData := DoRestRequest(ctx, diags, client, fmt.Sprintf("api/mo/%s.json", dn), "GET", nil)
	if requestData.Search("imdata").Search(classname).Data() != nil {
		diags.AddError(
			"Object Already Exists",
			fmt.Sprintf("The %s object with DN '%s' already exists.", classname, dn),
		)
	}
}

func DoRestRequestEscapeHtml(ctx context.Context, diags *diag.Diagnostics, aciClient *client.Client, path, method string, payload *container.Container, escapeHtml bool) *container.Container {

	// Ensure path starts with a slash to assure signature is created correctly
	if !strings.HasPrefix("/", path) {
		path = fmt.Sprintf("/%s", path)
	}
	var restRequest *http.Request
	var err error
	if escapeHtml {
		restRequest, err = aciClient.MakeRestRequest(method, path, payload, true)
	} else {
		restRequest, err = aciClient.MakeRestRequestRaw(method, path, payload.EncodeJSON(), true)
	}
	if err != nil {
		diags.AddError(
			"Creation of rest request failed",
			fmt.Sprintf("err: %s. Please report this issue to the provider developers.", err),
		)
		return nil
	}

	cont, restResponse, err := aciClient.Do(restRequest)

	if err != nil {
		diags.AddError(
			fmt.Sprintf("The %s rest request failed", strings.ToLower(method)),
			fmt.Sprintf("err: %s. Please report this issue to the provider developers.", err),
		)
		return nil
	}

	if restResponse != nil && cont.Data() != nil && restResponse.StatusCode != 200 {
		errCode := models.StripQuotes(models.StripSquareBrackets(cont.Search("imdata", "error", "attributes", "code").String()))
		errText := models.StripQuotes(models.StripSquareBrackets(cont.Search("imdata", "error", "attributes", "text").String()))
		// Ignore errors of type "Cannot create object", "Cannot delete object", "Request in progress", error text containing "can not be deleted." when the error code is 120 and error text containing "cannot be deleted." when the error code is 1
		if errCode == "103" || errCode == "107" || errCode == "202" || (errCode == "120" && strings.HasSuffix(errText, "can not be deleted.")) || (errCode == "1" && strings.HasSuffix(errText, "cannot be deleted.")) {
			tflog.Debug(ctx, fmt.Sprintf("Exiting from error: Code: %s, Message: %s", errCode, errText))
			return nil
		} else if (errText == "" && errCode == "403") || errCode == "401" {
			diags.AddError(
				"Unable to authenticate. Please check your credentials",
				fmt.Sprintf("Response Status Code: %d, Error Code: %s, Error Message: %s.", restResponse.StatusCode, errCode, errText),
			)
			return nil
		} else {
			diags.AddError(
				fmt.Sprintf("The %s rest request failed", strings.ToLower(method)),
				fmt.Sprintf("Response Status Code: %d, Error Code: %s, Error Message: %s.", restResponse.StatusCode, errCode, errText),
			)
			return nil
		}
	}

	return cont
}

func DoRestRequest(ctx context.Context, diags *diag.Diagnostics, client *client.Client, path, method string, payload *container.Container) *container.Container {
	return DoRestRequestEscapeHtml(ctx, diags, client, path, method, payload, true)
}

func GetDeleteJsonPayload(ctx context.Context, diags *diag.Diagnostics, className, dn string) *container.Container {

	jsonString := fmt.Sprintf(`{"%s":{"attributes":{"dn": "%s","status": "deleted"}}}`, className, dn)
	jsonPayload, err := container.ParseJSON([]byte(jsonString))
	if err != nil {
		diags.AddError(
			"Construction of json payload failed",
			fmt.Sprintf("Err: %s. Please report this issue to the provider developers.", err),
		)
		return nil
	}
	return jsonPayload
}

type setToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate struct{}

func SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate() planmodifier.String {
	return setToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate{}
}

func (m setToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate) Description(_ context.Context) string {
	return "During the update phase, set the value of this attribute to StringNull when the state value is null and the plan value is unknown."
}

func (m setToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate) MarkdownDescription(_ context.Context) string {
	return "During the update phase, set the value of this attribute to StringNull when the state value is null and the plan value is unknown."
}

// Custom plan modifier to set the plan value to null under certain conditions
func (m setToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Set the plan value to StringType null when state value is null and plan value is unknown during an Update
	if !req.State.Raw.IsNull() && req.StateValue.IsNull() && req.PlanValue.IsUnknown() {
		resp.PlanValue = types.StringNull()
	}
}

type setToSetNullWhenStateIsNullPlanIsUnknownDuringUpdate struct {
	resourceFunction func(ctx context.Context, planValue types.Set, stateValue types.Set) types.Set
}

func SetToSetNullWhenStateIsNullPlanIsUnknownDuringUpdate(resourceFunction func(ctx context.Context, planValue types.Set, stateValue types.Set) types.Set) planmodifier.Set {
	return setToSetNullWhenStateIsNullPlanIsUnknownDuringUpdate{
		resourceFunction: resourceFunction,
	}
}

func (m setToSetNullWhenStateIsNullPlanIsUnknownDuringUpdate) Description(_ context.Context) string {
	return "During the update phase, set the value of this attribute to StringNull when the state value is null and the plan value is unknown."
}

func (m setToSetNullWhenStateIsNullPlanIsUnknownDuringUpdate) MarkdownDescription(_ context.Context) string {
	return "During the update phase, set the value of this attribute to StringNull when the state value is null and the plan value is unknown."
}

// Custom plan modifier to set the plan value to null under certain conditions
func (m setToSetNullWhenStateIsNullPlanIsUnknownDuringUpdate) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	// Set the plan value to SetType null when state value is null and plan value is unknown during an Update
	if !req.State.Raw.IsNull() && req.StateValue.IsNull() && req.PlanValue.IsUnknown() {
		resp.PlanValue = types.SetNull(req.StateValue.ElementType(ctx))
	} else if !req.State.Raw.IsNull() && !req.StateValue.IsNull() && !req.PlanValue.IsUnknown() && m.resourceFunction != nil {
		resp.PlanValue = m.resourceFunction(ctx, req.PlanValue, req.StateValue)
	}
}

// MakeStringRequiredValidator validates that an attribute is not null, as a workaround for when a resource has read-only attributes and nested sets with required attributes
// https://github.com/hashicorp/terraform-plugin-framework/issues/898

var _ validator.String = MakeStringRequiredValidator{}

// MakeStringRequiredValidator validates that an attribute is not null. Most
// attributes should set Required: true instead, however in certain scenarios,
// such as a computed nested attribute, all underlying attributes must also be
// computed for planning to not show unexpected differences.
type MakeStringRequiredValidator struct{}

// Description describes the validation in plain text formatting.
func (v MakeStringRequiredValidator) Description(_ context.Context) string {
	return "is required"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v MakeStringRequiredValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v MakeStringRequiredValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if !req.ConfigValue.IsNull() {
		return
	}

	path := req.Path.String()
	if strings.Contains(path, "[") {
		setName := path[0:strings.Index(path, "[")]
		attributeName := path[strings.Index(path, "]")+2:]
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Incorrect attribute value type",
			fmt.Sprintf("Inappropriate value for attribute \"%s\": attribute \"%s\" is required.", setName, attributeName),
		)
	} else {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Incorrect attribute value type",
			fmt.Sprintf("Attribute is required for path: %s", path),
		)
	}

}

// StringNotNull returns an validator which ensures that the string attribute is
// configured. Most attributes should set Required: true instead, however in
// certain scenarios, such as a computed nested attribute, all underlying
// attributes must also be computed for planning to not show unexpected
// differences.
func MakeStringRequired() validator.String {
	return MakeStringRequiredValidator{}
}

// SingleNestedAttributeRequiredAttributesNotProvidedValidator validates that all required attributes are provided when it is not {}
// {} logic is needed in order to remove the child object from APIC

var _ validator.Object = SingleNestedAttributeRequiredAttributesNotProvidedValidator{}

type SingleNestedAttributeRequiredAttributesNotProvidedValidator struct {
	attributeName      string
	requiredAttributes []string
}

func (v SingleNestedAttributeRequiredAttributesNotProvidedValidator) Description(_ context.Context) string {
	return "is required"
}

func (v SingleNestedAttributeRequiredAttributesNotProvidedValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v SingleNestedAttributeRequiredAttributesNotProvidedValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() {
		return
	}

	if !IsEmptySingleNestedAttribute(req.ConfigValue.Attributes()) && SingleNestedAttributeRequiredAttributesNotProvided(req.ConfigValue.Attributes(), v.requiredAttributes) {
		errMessage := fmt.Sprintf("Inappropriate value for attribute \"%s\": attribute \"%s\" is required.", v.attributeName, strings.Join(v.requiredAttributes, ", "))
		if len(v.requiredAttributes) > 1 {
			errMessage = fmt.Sprintf("Inappropriate values for attribute \"%s\": attributes \"%s\" are required.", v.attributeName, strings.Join(v.requiredAttributes, "\", \""))
		}
		resp.Diagnostics.AddError(
			"Incorrect attribute value type",
			errMessage,
		)
	}
}

func MakeSingleNestedAttributeRequiredAttributesNotProvidedValidator(atributeName string, requiredAttributes []string) validator.Object {
	return SingleNestedAttributeRequiredAttributesNotProvidedValidator{
		attributeName:      atributeName,
		requiredAttributes: requiredAttributes,
	}
}
