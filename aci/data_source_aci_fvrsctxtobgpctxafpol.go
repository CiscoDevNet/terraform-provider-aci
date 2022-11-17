package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAciBGPAddressFamilyContextPolicyRelationship() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciBGPAddressFamilyContextPolicyRelationshipRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"vrf_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"address_family": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ipv4-ucast",
					"ipv6-ucast",
				}, false),
			},
			"annotation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bgp_address_family_context_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
		})),
	}
}

func dataSourceAciBGPAddressFamilyContextPolicyRelationshipRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	tnBgpCtxAfPolName := GetMOName(d.Get("bgp_address_family_context_dn").(string))
	address_family := d.Get("address_family").(string)
	VRFDn := d.Get("vrf_dn").(string)
	rn := fmt.Sprintf(models.RnfvRsCtxToBgpCtxAfPol, tnBgpCtxAfPolName, address_family)
	dn := fmt.Sprintf("%s/%s", VRFDn, rn)

	fvRsCtxToBgpCtxAfPol, err := getRemoteBGPAddressFamilyContextPolicyRelationship(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setBGPAddressFamilyContextPolicyRelationshipAttributes(fvRsCtxToBgpCtxAfPol, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
