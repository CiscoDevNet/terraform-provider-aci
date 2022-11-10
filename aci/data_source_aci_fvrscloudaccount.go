package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciTenantToCloudAccountAssociation() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciTenantToCloudAccountAssociationRead,
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

func dataSourceAciTenantToCloudAccountAssociationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	TenantDn := d.Get("tenant_dn").(string)
	rn := fmt.Sprintf(models.RnfvRsCloudAccount)
	dn := fmt.Sprintf("%s/%s", TenantDn, rn)

	fvRsCloudAccount, err := getRemoteTenantToCloudAccountAssociation(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setTenantToCloudAccountAssociationAttributes(fvRsCloudAccount, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
