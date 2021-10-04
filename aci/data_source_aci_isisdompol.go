package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
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
			"isis_level": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(
						map[string]*schema.Schema{
							"id": &schema.Schema{
								Type:     schema.TypeString,
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
							"name": &schema.Schema{
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
							"isis_level_type": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
								Computed: true,
							},
						},
					)),
				},
			},
		})),
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
	lvlCompIds, err := getISISDomPolChildernIds(aciClient, isisDomPol.DistinguishedName)
	if err != nil {
		return diag.FromErr(err)
	}
	comps := lvlCompIds
	lvlComps := make([]*models.ISISLevel, 0, 1)

	for _, val := range comps {
		lvlCompDN := val
		lvlComp, err := getRemoteISISLevel(aciClient, lvlCompDN)
		if err != nil {
			return diag.FromErr(err)
		}
		lvlComps = append(lvlComps, lvlComp)
	}

	_, err = setISISLevelAttributes(lvlComps, d)

	if err != nil {
		d.SetId("")
		return nil
	}
	return nil
}

func getISISDomPolChildernIds(aciClient *client.Client, isisDn string) ([]string, error) {
	lvlCompUrl := fmt.Sprintf("/api/node/mo/%s.json?query-target=children", isisDn)
	lvlCompCont, err := aciClient.GetViaURL(lvlCompUrl)
	if err != nil {
		return nil, err
	}
	childCont := lvlCompCont.S("imdata")
	isisCompIds := make([]string, 0, 1)
	for i := 0; i < len(childCont.Data().([]interface{})); i++ {
		isisCompId := G(childCont.Index(i).S("isisLvlComp", "attributes"), "dn")
		if isisCompId != "{}" {
			isisCompIds = append(isisCompIds, isisCompId)
		}
	}
	return isisCompIds, nil
}
