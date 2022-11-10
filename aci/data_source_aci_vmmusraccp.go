package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciVMMCredential() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciVMMCredentialRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"vmm_domain_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"annotation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pwd": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"usr": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciVMMCredentialRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	VMMDomainDn := d.Get("vmm_domain_dn").(string)
	rn := fmt.Sprintf("usracc-%s", name)
	dn := fmt.Sprintf("%s/%s", VMMDomainDn, rn)
	vmmUsrAccP, err := getRemoteVMMCredential(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setVMMCredentialAttributes(vmmUsrAccP, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
