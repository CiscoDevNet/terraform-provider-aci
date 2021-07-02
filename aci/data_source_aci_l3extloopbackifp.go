package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciLoopBackInterfaceProfile() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciLoopBackInterfaceProfileRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"fabric_node_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"addr": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciLoopBackInterfaceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	addr := d.Get("addr").(string)

	rn := fmt.Sprintf("lbp-[%s]", addr)
	FabricNodeDn := d.Get("fabric_node_dn").(string)

	dn := fmt.Sprintf("%s/%s", FabricNodeDn, rn)

	l3extLoopBackIfP, err := getRemoteLoopBackInterfaceProfile(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)
	_, err = setLoopBackInterfaceProfileAttributes(l3extLoopBackIfP, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
