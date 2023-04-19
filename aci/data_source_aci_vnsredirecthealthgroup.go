package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciL4L7RedirectHealthGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciL4L7RedirectHealthGroupRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		})),
	}
}

func dataSourceAciL4L7RedirectHealthGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)
	rn := fmt.Sprintf(models.RnvnsRedirectHealthGroup, name)
	dn := fmt.Sprintf("%s/%s", TenantDn, rn)

	vnsRedirectHealthGroup, err := getRemoteL4L7RedirectHealthGroup(aciClient, dn)
	if err != nil {
		return nil
	}

	d.SetId(dn)

	_, err = setL4L7RedirectHealthGroupAttributes(vnsRedirectHealthGroup, d)
	if err != nil {
		return nil
	}

	return nil
}
