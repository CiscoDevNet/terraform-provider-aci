package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciPIMInterfacePolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciPIMInterfacePolicyRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"auth_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auth_t": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ctrl": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dr_delay": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dr_prio": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hello_itvl": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"jp_interval": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"secure_auth_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciPIMInterfacePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)
	rn := fmt.Sprintf(models.RnPimIfPol, name)
	dn := fmt.Sprintf("%s/%s", TenantDn, rn)

	pimIfPol, err := getRemotePIMInterfacePolicy(aciClient, dn)
	if err != nil {
		return nil
	}

	d.SetId(dn)

	_, err = setPIMInterfacePolicyAttributes(pimIfPol, d)
	if err != nil {
		return nil
	}

	return nil
}
