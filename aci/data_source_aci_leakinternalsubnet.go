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
			// True -> public, False -> private, Default -> false(private)
			"allow_l3out_advertisement": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// Tenant and VRF Destinations
			"leak_to": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vrf_dn": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						// True -> public, False -> private, Default -> "inherit"
						"allow_l3out_advertisement": &schema.Schema{
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
	log.Printf("[DEBUG] LeakInternalSubnet: Beginning of data source read")
	aciClient := m.(*client.Client)
	ip := d.Get("ip").(string)
	VRFDn := d.Get("vrf_dn").(string)
	rn := fmt.Sprintf(models.RnleakInternalSubnet, ip)
	dn := fmt.Sprintf("%s/%s/%s", VRFDn, models.RnleakRoutes, rn)

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
	log.Printf("[DEBUG] LeakInternalSubnet: %s data source read finished successfully", d.Id())
	return nil
}
