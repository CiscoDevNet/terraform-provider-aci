package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciTenanttoCloudAccountAssociation() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciTenanttoCloudAccountAssociationRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"annotation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cloud_account_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
		})),
	}
}

func dataSourceAciTenanttoCloudAccountAssociationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	TenantDn := d.Get("tenant_dn").(string)
	rn := fmt.Sprintf("rsCloudAccount")
	dn := fmt.Sprintf("%s/%s", TenantDn, rn)

	fvRsCloudAccount, err := getRemoteTenanttoCloudAccountAssociation(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setTenanttoCloudAccountAssociationAttributes(fvRsCloudAccount, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
