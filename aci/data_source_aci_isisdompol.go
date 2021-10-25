package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciISISDomainPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciISISDomainPolicyRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"mtu": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"redistrib_metric": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"lsp_fast_flood": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"lsp_gen_init_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"lsp_gen_max_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"lsp_gen_sec_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"isis_level_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"spf_comp_init_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"spf_comp_max_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"spf_comp_sec_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
		)),
	}
}

func dataSourceAciISISDomainPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := "default"

	rn := fmt.Sprintf("fabric/isisDomP-%s", name)
	dn := fmt.Sprintf("uni/%s", rn)
	isisDomPol, err := getRemoteISISDomainPolicy(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	_, err = setISISDomainPolicyAttributes(isisDomPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	lvlCompUrl := fmt.Sprintf("/api/node/mo/%s.json?query-target=children", dn)

	lvlCompCont, err := aciClient.GetViaURL(lvlCompUrl)
	if err != nil {
		return diag.FromErr(err)
	}
	childCont := lvlCompCont.S("imdata")
	isisCompDn := ""

	for i := 0; i < len(childCont.Data().([]interface{})); i++ {
		isisCompId := G(childCont.Index(i).S("isisLvlComp", "attributes"), "dn")
		if isisCompId != "{}" {
			isisCompDn = isisCompId
		}
	}

	if isisCompDn != "" {
		isisLvlComp, err := getRemoteISISLevel(aciClient, isisCompDn)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = setISISLevelAttributes(isisLvlComp, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}
