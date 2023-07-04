package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciLogicalInterfaceProfile() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciLogicalInterfaceProfileRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"logical_node_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"prio": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"tag": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"relation_l3ext_rs_l_if_p_to_netflow_monitor_pol": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
			},

			"relation_l3ext_rs_egress_qos_dpp_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"relation_l3ext_rs_ingress_qos_dpp_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"relation_l3ext_rs_l_if_p_cust_qos_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"relation_l3ext_rs_arp_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"relation_l3ext_rs_nd_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"relation_l3ext_rs_pim_ip_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"relation_l3ext_rs_pim_ipv6_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"relation_l3ext_rs_igmp_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciLogicalInterfaceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf(models.Rnl3extlifp, name)
	LogicalNodeProfileDn := d.Get("logical_node_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", LogicalNodeProfileDn, rn)

	l3extLIfP, err := getRemoteLogicalInterfaceProfile(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setLogicalInterfaceProfileAttributes(l3extLIfP, d)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = getandSetL3extLIfPRelationshipAttributes(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] L3extLIfP Relationship Attributes - Read finished successfully")
	}

	return nil
}
