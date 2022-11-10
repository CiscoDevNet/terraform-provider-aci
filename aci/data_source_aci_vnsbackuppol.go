package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciPBRBackupPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciPBRBackupPolicyRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		})),
	}
}

func dataSourceAciPBRBackupPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)
	rn := fmt.Sprintf(models.RnvnsBackupPol, name)
	dn := fmt.Sprintf("%s/%s/%s", TenantDn, models.RnvnsSvcCont, rn)

	vnsBackupPol, err := getRemotePBRBackupPolicy(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setPBRBackupPolicyAttributes(vnsBackupPol, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
