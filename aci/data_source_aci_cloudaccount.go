package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciCloudAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciCloudAccountRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vendor": {
				Type:     schema.TypeString,
				Required: true,
			},
		})),
	}
}

func dataSourceAciCloudAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	account_id := d.Get("account_id").(string)
	vendor := d.Get("vendor").(string)
	TenantDn := d.Get("tenant_dn").(string)
	rn := fmt.Sprintf(models.RncloudAccount, account_id, vendor)
	dn := fmt.Sprintf("%s/%s", TenantDn, rn)

	cloudAccount, err := getRemoteCloudAccount(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setCloudAccountAttributes(cloudAccount, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
