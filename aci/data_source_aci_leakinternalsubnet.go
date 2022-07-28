package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciLeakInternalSubnet() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciLeakInternalSubnetRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"vrf_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vrf_scope": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// Tenant and VRF Destinations
			"leak_to": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_vrf_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"destination_vrf_scope": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"destination_tenant_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		})),
	}
}

func dataSourceAciLeakInternalSubnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	ip := d.Get("ip").(string)
	VRFDn := d.Get("vrf_dn").(string)
	rn := fmt.Sprintf("leakintsubnet-[%s]", ip)
	dn := fmt.Sprintf("%s/%s", VRFDn, rn)

	leakInternalSubnet, err := getRemoteLeakInternalSubnet(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setLeakInternalSubnetAttributes(leakInternalSubnet, d)
	if err != nil {
		return diag.FromErr(err)
	}

	// leakTo - Beginning Read
	_, err = getAndSetLeakToObjects(aciClient, dn, d)
	if err != nil {
		return diag.FromErr(err)
	}
	// leakTo - Read finished successfully

	return nil
}
