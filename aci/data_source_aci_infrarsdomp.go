package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciInfraRsDomP() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciInfraRsDomPRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"attachable_access_entity_profile_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"annotation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"domain_dn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciInfraRsDomPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	domain_dn := d.Get("domain_dn").(string)
	attachableAccessEntityProfileDn := d.Get("attachable_access_entity_profile_dn").(string)
	rn := fmt.Sprintf(models.RninfraRsDomP, domain_dn)
	dn := fmt.Sprintf("%s/%s", attachableAccessEntityProfileDn, rn)

	infraRsDomP, err := getRemoteInfraRsDomP(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setInfraRsDomPAttributes(infraRsDomP, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
