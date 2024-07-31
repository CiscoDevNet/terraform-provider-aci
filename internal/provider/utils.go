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
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
		errText := models.StripQuotes(models.StripSquareBrackets(cont.Search("imdata", "error", "attributes", "text").String()))
		// Ignore errors of type "Cannot create object", "Cannot delete object", "Request in progress" and error text containing "can not be deleted." when the error code is 120
		if errCode == "103" || errCode == "107" || errCode == "202" || (errCode == "120" && strings.HasSuffix(errText, "can not be deleted.")) {
			tflog.Debug(ctx, errText, map[string]interface{}{
				"Error Code": errCode,
			})
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
