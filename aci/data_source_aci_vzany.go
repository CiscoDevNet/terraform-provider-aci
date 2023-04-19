package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciAny() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciAnyRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"vrf_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"match_t": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pref_gr_memb": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"relation_vz_rs_any_to_cons": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"relation_vz_rs_any_to_cons_if": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"relation_vz_rs_any_to_prov": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		})),
	}
}

func dataSourceAciAnyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	VRFDn := d.Get("vrf_dn").(string)

	dn := fmt.Sprintf("%s/any", VRFDn)

	vzAny, err := getRemoteAny(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setAnyAttributes(vzAny, d)
	if err != nil {
		return diag.FromErr(err)
	}

	getAndSetVzRsAnyToConsFromAny(aciClient, dn, d)
	getAndSetVzRsAnyToConsIfFromAny(aciClient, dn, d)
	getAndSetVzRsAnyToProvFromAny(aciClient, dn, d)

	return nil
}
