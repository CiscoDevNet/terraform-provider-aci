package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciSubnet() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciSubnetRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"parent_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"ctrl": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"preferred": &schema.Schema{
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

			"virtual": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// EP Reachability
			"next_hop_addr": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"msnlb", "anycast_mac"},
			},
			// MSNLB
			"msnlb": {
				Type:          schema.TypeSet,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"next_hop_addr", "anycast_mac"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mac": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"group": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"mode": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			// Anycast MAC
			"anycast_mac": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"msnlb", "next_hop_addr"},
			},
		}),
	}
}

func dataSourceAciSubnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	ip := d.Get("ip").(string)

	rn := fmt.Sprintf("subnet-[%s]", ip)
	BridgeDomainDn := d.Get("parent_dn").(string)

	dn := fmt.Sprintf("%s/%s", BridgeDomainDn, rn)

	fvSubnet, err := getRemoteSubnet(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setSubnetAttributes(fvSubnet, d)
	if err != nil {
		return diag.FromErr(err)
	}

	// ipNexthopEpP - Beginning Read
	ipNexthopEpPParentDn := dn + "/epReach"
	log.Printf("[DEBUG] %s: ipNexthopEpP - Beginning Read with parent DN", ipNexthopEpPParentDn)
	_, err = getAndSetNexthopEpPReachability(aciClient, ipNexthopEpPParentDn, d)
	if err == nil {
		ipNexthopEpPDn := dn + "/epReach/" + fmt.Sprintf(models.RnipNexthopEpP, d.Get("next_hop_addr"))
		log.Printf("[DEBUG] %s: ipNexthopEpP - Read finished successfully", ipNexthopEpPDn)
	} else {
		log.Printf("[DEBUG] %s: ipNexthopEpP - Object not present in the parent", ipNexthopEpPParentDn)
	}
	// ipNexthopEpP - Read finished successfully

	// fvEpNlb - Beginning Read
	fvEpNlbDn := dn + fmt.Sprintf("/"+models.RnfvEpNlb)
	fvEpNlb, err := getRemoteNlbEndpoint(aciClient, fvEpNlbDn)
	if err == nil {
		log.Printf("[DEBUG] %s: fvEpNlb - Beginning Read", fvEpNlbDn)
		_, err = setNlbEndpointAttributes(fvEpNlb, d)
		if err != nil {
			return nil
		}
		log.Printf("[DEBUG] %s: fvEpNlb - Read finished successfully", fvEpNlbDn)
	} else {
		d.Set("msnlb", nil)
	}
	// fvEpNlb - Read finished successfully

	// fvEpAnycast - Beginning of Read
	log.Printf("[DEBUG] %s: fvEpAnycast - Beginning Read with parent DN", dn)
	_, err = getAndSetAnycastMac(aciClient, dn, d)
	if err == nil {
		fvEpAnycastDn := dn + "/" + fmt.Sprintf(models.RnfvEpAnycast, d.Get("anycast_mac"))
		log.Printf("[DEBUG] %s: fvEpAnycast - Read finished successfully", fvEpAnycastDn)
	} else {
		log.Printf("[DEBUG] %s: fvEpAnycast - Object not present in the parent", dn)
	}
	// fvEpAnycast - Read finished successfully

	return nil
}
