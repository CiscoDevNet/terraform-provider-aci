package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciL1L2RedirectDestTraffic() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciL1L2RedirectDestTrafficRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"policy_based_redirect_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"dest_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"mac": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pod_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"relation_vns_rs_l1_l2_redirect_health_group": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Create relation to vns:RedirectHealthGroup",
			},
			"relation_vns_rs_to_c_if": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Create relation to vns:CIf",
			},
		})),
	}
}

func dataSourceAciL1L2RedirectDestTrafficRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	destName := d.Get("dest_name").(string)
	policyBasedRedirectDn := d.Get("policy_based_redirect_dn").(string)
	rn := fmt.Sprintf(models.RnvnsL1L2RedirectDest, destName)
	dn := fmt.Sprintf("%s/%s", policyBasedRedirectDn, rn)
	vnsL1L2RedirectDest, err := getRemoteL1L2RedirectDestTraffic(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setL1L2RedirectDestTrafficAttributes(vnsL1L2RedirectDest, d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Importing vnsRsL1L2RedirectHealthGroup object
	getAndSetRemoteReadRelationvnsRsL1L2RedirectHealthGroup(aciClient, dn, d)

	// Importing vnsRsToCIf object
	getAndSetRemoteReadRelationvnsRsToCIf(aciClient, dn, d)

	return nil
}
