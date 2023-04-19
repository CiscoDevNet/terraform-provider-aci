package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciExternalNetworkInstanceProfile() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciExternalNetworkInstanceProfileRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"l3_outside_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"exception_tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"flood_on_encap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

			"prio": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"target_dscp": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"relation_fv_rs_sec_inherited": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
				Computed: true,
			},
			"relation_fv_rs_prov": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
				Computed: true,
			},
			"relation_fv_rs_cons_if": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
				Computed: true,
			},
			"relation_l3ext_rs_inst_p_to_profile": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tn_rtctrl_profile_dn": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"direction": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"relation_fv_rs_cons": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
				Computed: true,
			},
			"relation_fv_rs_prot_by": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciExternalNetworkInstanceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf(models.Rnl3extinstp, name)
	L3OutsideDn := d.Get("l3_outside_dn").(string)

	dn := fmt.Sprintf("%s/%s", L3OutsideDn, rn)

	l3extInstP, err := getRemoteExternalNetworkInstanceProfile(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setExternalNetworkInstanceProfileAttributes(l3extInstP, d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Importing l3extRsInstPToProfile object
	getAndSetReadRelationl3extRsInstPToProfileFromExternalNetworkInstanceProfile(aciClient, dn, d)

	// Importing fvRsSecInherited object
	getAndSetReadRelationfvRsSecInheritedFromExternalNetworkInstanceProfile(aciClient, dn, d)

	// Importing fvRsProv object
	getAndSetReadRelationfvRsProvFromExternalNetworkInstanceProfile(aciClient, dn, d)

	// Importing fvRsConsIf object
	getAndSetReadRelationfvRsConsIfFromExternalNetworkInstanceProfile(aciClient, dn, d)

	// Importing fvRsCons object
	getAndSetReadRelationfvRsConsFromExternalNetworkInstanceProfile(aciClient, dn, d)

	// Importing fvRsProtBy object
	getAndSetReadRelationfvRsProtByFromExternalNetworkInstanceProfile(aciClient, dn, d)

	return nil
}
