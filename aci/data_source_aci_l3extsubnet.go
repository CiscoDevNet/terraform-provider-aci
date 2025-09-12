package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciL3ExtSubnet() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciL3ExtSubnetRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"external_network_instance_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"aggregate": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"scope": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"relation_l3ext_rs_subnet_to_profile": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tn_rtctrl_profile_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tn_rtctrl_profile_dn": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"direction": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"relation_l3ext_rs_subnet_to_rt_summ": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		})),
	}
}

func dataSourceAciL3ExtSubnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	ip := d.Get("ip").(string)

	rn := fmt.Sprintf("extsubnet-[%s]", ip)
	ExternalNetworkInstanceProfileDn := d.Get("external_network_instance_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", ExternalNetworkInstanceProfileDn, rn)

	l3extSubnet, err := getRemoteL3ExtSubnet(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setL3ExtSubnetAttributes(l3extSubnet, d)
	if err != nil {
		return diag.FromErr(err)
	}

	getAndSetl3extRsSubnetToProfileFromL3ExtSubnet(aciClient, dn, d)
	getAndSetl3extRsSubnetToRtSummFromL3ExtSubnet(aciClient, dn, d)

	return nil
}
