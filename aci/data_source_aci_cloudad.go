package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciCloudActiveDirectory() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciCloudActiveDirectoryRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"active_directory_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		})),
	}
}

func dataSourceAciCloudActiveDirectoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	active_directory_id := d.Get("active_directory_id").(string)
	TenantDn := d.Get("tenant_dn").(string)
	rn := fmt.Sprintf(models.RncloudAD, active_directory_id)
	dn := fmt.Sprintf("%s/%s", TenantDn, rn)

	cloudAD, err := getRemoteCloudActiveDirectory(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setCloudActiveDirectoryAttributes(cloudAD, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
