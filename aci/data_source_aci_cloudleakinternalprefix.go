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

func dataSourceAciLeakInternalPrefix() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciLeakInternalPrefixRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"vrf_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"leak_to": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vrf_dn": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		})),
	}
}

func dataSourceAciLeakInternalPrefixRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LeakInternalPrefix: Beginning of data source read")
	aciClient := m.(*client.Client)
	VRFDn := d.Get("vrf_dn").(string)
	rn := fmt.Sprintf(models.RnleakInternalPrefix, leakInternalPrefixIp)
	dn := fmt.Sprintf("%s/%s/%s", VRFDn, models.RnleakRoutes, rn)

	LeakInternalPrefix, err := getRemoteLeakInternalPrefix(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setLeakInternalPrefixAttributes(LeakInternalPrefix, d)
	if err != nil {
		return diag.FromErr(err)
	}

	// leakTo - Beginning Read
	_, err = getAndSetCloudLeakToObjects(aciClient, dn, d)
	if err != nil {
		return diag.FromErr(err)
	}
	// leakTo - Read finished successfully
	log.Printf("[DEBUG] LeakInternalPrefix: %s data source read finished successfully", d.Id())
	return nil
}
