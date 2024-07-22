package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
		if errCode != "1" && errCode != "103" && errCode != "107" && errCode != "120" {
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

type setToSetNullWhenStateIsNullPlanIsUnknownDuringUpdate struct{}

func SetToSetNullWhenStateIsNullPlanIsUnknownDuringUpdate() planmodifier.Set {
	return setToSetNullWhenStateIsNullPlanIsUnknownDuringUpdate{}
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

	setName := req.Path.String()[0:strings.Index(req.Path.String(), "[")]
	attributeName := req.Path.String()[strings.Index(req.Path.String(), "]")+2:]

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Incorrect attribute value type",
		fmt.Sprintf("Inappropriate value for attribute \"%s\": attribute \"%s\" is required.", setName, attributeName),
	)
}

// StringNotNull returns an validator which ensures that the string attribute is
// configured. Most attributes should set Required: true instead, however in
// certain scenarios, such as a computed nested attribute, all underlying
// attributes must also be computed for planning to not show unexpected
// differences.
func MakeStringRequired() validator.String {
	return MakeStringRequiredValidator{}
}
