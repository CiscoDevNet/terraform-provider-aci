package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func GetMOName(dn string) string {
	splittedDn := strings.Split(dn, "/")
	if len(splittedDn) > 1 {
		return strings.Join(strings.Split(splittedDn[len(splittedDn)-1], "-")[1:], "-")
	}
	return splittedDn[0]
}

func CheckDn(ctx context.Context, client *client.Client, dn string, diags *diag.Diagnostics) {
	tflog.Debug(ctx, fmt.Sprintf("validate relation dn: %s", dn))
	_, err := client.Get(dn)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("failed validate relation dn: %s", dn))
		diags.AddError(
			"Relation target dn validation failed",
			fmt.Sprintf("The relation target dn is not found: %s", dn),
		)
	}
}
