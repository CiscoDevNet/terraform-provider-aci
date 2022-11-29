package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciL3Outside() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciL3OutsideRead,

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

			"enforce_rtctrl": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"target_dscp": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"relation_l3ext_rs_dampening_pol": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tn_rtctrl_profile_name": {
							Type:       schema.TypeString,
							Optional:   true,
							Computed:   true,
							Deprecated: "Use tn_rtctrl_profile_dn instead of tn_rtctrl_profile_name",
						},
						"tn_rtctrl_profile_dn": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"af": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"relation_l3ext_rs_ectx": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_l3ext_rs_out_to_bd_public_subnet_holder": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_l3ext_rs_interleak_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_l3ext_rs_l3_dom_att": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		})),
	}
}

func dataSourceAciL3OutsideRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf(models.Rnl3extOut, name)
	TenantDn := d.Get("tenant_dn").(string)

	dn := fmt.Sprintf("%s/%s", TenantDn, rn)

	l3extOut, err := getRemoteL3Outside(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setL3OutsideAttributes(l3extOut, d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Importing l3extRsDampeningPol object
	getAndSetReadRelationl3extRsDampeningPolFromL3Outside(aciClient, dn, d)

	// Importing l3extRsEctx object
	getAndSetReadRelationl3extRsEctxFromL3Outside(aciClient, dn, d)

	// Importing l3extRsOutToBDPublicSubnetHolder object
	getAndSetReadRelationl3extRsOutToBDPublicSubnetHolderFromL3Outside(aciClient, dn, d)

	// Importing l3extRsInterleakPol object
	getAndSetReadRelationl3extRsInterleakPolFromL3Outside(aciClient, dn, d)

	// Importing l3extRsL3DomAtt object
	getAndSetReadRelationl3extRsL3DomAttFromL3Outside(aciClient, dn, d)

	return nil
}
