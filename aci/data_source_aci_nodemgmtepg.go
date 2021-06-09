package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAciNodeManagementEPg() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciNodeManagementEPgRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"in_band",
					"out_of_band",
				}, false),
			},

			"management_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"encap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

			"name_alias": &schema.Schema{
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
		}),
	}
}

func dataSourceAciNodeManagementEPgRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.Get("type").(string) == "in_band" {
		return dataSourceInBandManagementEPgRead(ctx, d, m)
	}
	return dataSourceOutOfBandManagementEPgRead(ctx, d, m)
}

func dataSourceInBandManagementEPgRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("inb-%s", name)
	ManagementProfileDn := d.Get("management_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", ManagementProfileDn, rn)

	mgmtInB, err := getRemoteInBandManagementEPg(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setInBandManagementEPgAttributes(mgmtInB, d)

	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func dataSourceOutOfBandManagementEPgRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("oob-%s", name)
	ManagementProfileDn := d.Get("management_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", ManagementProfileDn, rn)

	mgmtOoB, err := getRemoteOutOfBandManagementEPg(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setOutOfBandManagementEPgAttributes(mgmtOoB, d)

	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
