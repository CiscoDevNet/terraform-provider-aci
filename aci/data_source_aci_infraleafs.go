package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciSwitchAssociation() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciSwitchAssociationRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"leaf_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"switch_association_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciSwitchAssociationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	switch_association_type := d.Get("switch_association_type").(string)

	rn := fmt.Sprintf("leaves-%s-typ-%s", name, switch_association_type)
	LeafProfileDn := d.Get("leaf_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", LeafProfileDn, rn)

	infraLeafS, err := getRemoteSwitchAssociation(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	setSwitchAssociationAttributes(infraLeafS, d)
	return nil
}
