package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciSubnetPoolforIpSecTunnels() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciSubnetPoolforIpSecTunnelsRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subnet_pool": {
				Type:     schema.TypeString,
				Required: true,
			},
		})),
	}
}

func dataSourceAciSubnetPoolforIpSecTunnelsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	subnetpool := d.Get("subnet_pool").(string)
	rn := fmt.Sprintf(models.RncloudtemplateIpSecTunnelSubnetPool, subnetpool)
	dn := fmt.Sprintf("%s/%s", models.CloudInfraNetworkDefaultTemplateDn, rn)

	cloudtemplateIpSecTunnelSubnetPool, err := getRemoteSubnetPoolforIpSecTunnels(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setSubnetPoolforIpSecTunnelsAttributes(cloudtemplateIpSecTunnelSubnetPool, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
