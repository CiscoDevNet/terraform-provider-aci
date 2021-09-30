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
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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
	name := d.Get("name").(string)

	rn := fmt.Sprintf("fabric/isisDomP-%s", name)
	dn := fmt.Sprintf("uni/%s", rn)
	isisDomPol, err := getRemoteISISDomainPolicy(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setISISDomainPolicyAttributes(isisDomPol, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
