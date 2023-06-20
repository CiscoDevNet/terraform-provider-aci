package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciPimInterfaceProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciInterfaceProfileRead,
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

			"relation_pim_rs_if_pol": {
				Type: schema.TypeString,

				Optional:    true,
				Description: "Query pim:IfPol relationship object",
			}})),
	}
}

func dataSourceAciPimInterfaceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	LogicalInterfaceProfileDn := d.Get("logical_interface_profile_dn").(string)
	rn := fmt.Sprintf(models.RnPimIfP)
	dn := fmt.Sprintf("%s/%s", LogicalInterfaceProfileDn, rn)

	pimIfP, err := getRemotePimInterfaceProfile(aciClient, dn)
	if err != nil {
		return nil
	}

	d.SetId(dn)

	_, err = setPimInterfaceProfileAttributes(pimIfP, d)
	if err != nil {
		return nil
	}

	// Get and Set Relational Attributes
	getAndSetPimInterfaceProfileRelationalAttributes(aciClient, dn, d)
	return nil
}
