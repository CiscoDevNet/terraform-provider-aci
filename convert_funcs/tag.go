package convert_funcs

import (
	"context"
	"encoding/json"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/provider"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func CreateTagTag(attributes map[string]interface{}, status string) map[string]interface{} {
	ctx := context.Background()
	var diags diag.Diagnostics
	data := &provider.TagTagResourceModel{}
	if v, ok := attributes["parent_dn"].(string); ok && v != "" {
		data.ParentDn = types.StringValue(v)
	}
	if v, ok := attributes["key"].(string); ok && v != "" {
		data.Key = types.StringValue(v)
	}
	if v, ok := attributes["value"].(string); ok && v != "" {
		data.Value = types.StringValue(v)
	}

	// Handle deletion logic
	if status == "deleted" {

		provider.SetTagTagId(ctx, data)

		deletePayload := provider.GetDeleteJsonPayload(ctx, &diags, "tagTag", data.Id.ValueString())
		if deletePayload != nil {
			jsonPayload := deletePayload.EncodeJSON(container.EncodeOptIndent("", "  "))
			var customData map[string]interface{}
			json.Unmarshal(jsonPayload, &customData)
			return customData
		}

	}

	newAciTagTag := provider.GetTagTagCreateJsonPayload(ctx, &diags, true, data)

	jsonPayload := newAciTagTag.EncodeJSON(container.EncodeOptIndent("", "  "))

	var customData map[string]interface{}
	json.Unmarshal(jsonPayload, &customData)

	payload := customData

	provider.SetTagTagId(ctx, data)
	attrs := payload["tagTag"].(map[string]interface{})["attributes"].(map[string]interface{})
	attrs["dn"] = data.Id.ValueString()

	return payload
}
