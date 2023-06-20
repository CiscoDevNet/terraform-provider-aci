package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciIGMPInterfaceProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciIGMPInterfaceProfileRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"logical_interface_profile_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"relation_igmp_rs_if_pol": {
				Type: schema.TypeString,

				Optional:    true,
				Description: "Query igmp:IfPol relationship object",
			}})),
	}
}

func dataSourceAciIGMPInterfaceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	LogicalInterfaceProfileDn := d.Get("logical_interface_profile_dn").(string)
	rn := fmt.Sprintf(models.RnIgmpIfP)
	dn := fmt.Sprintf("%s/%s", LogicalInterfaceProfileDn, rn)

	igmpIfP, err := getRemoteIGMPInterfaceProfile(aciClient, dn)
	if err != nil {
		return nil
	}

	d.SetId(dn)

	_, err = setIGMPInterfaceProfileAttributes(igmpIfP, d)
	if err != nil {
		return nil
	}

	// Get and Set Relational Attributes
	getAndSetIGMPInterfaceProfileRelationalAttributes(aciClient, dn, d)
	return nil
}
