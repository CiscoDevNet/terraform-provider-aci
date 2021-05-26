package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciVlanEncapsulationforVxlanTraffic() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciVlanEncapsulationforVxlanTrafficRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"attachable_access_entity_profile_dn": &schema.Schema{
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

func dataSourceAciVlanEncapsulationforVxlanTrafficRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	rn := fmt.Sprintf("provacc")
	AttachableAccessEntityProfileDn := d.Get("attachable_access_entity_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", AttachableAccessEntityProfileDn, rn)

	infraProvAcc, err := getRemoteVlanEncapsulationforVxlanTraffic(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	
	d.SetId(dn)
	_, err = setVlanEncapsulationforVxlanTrafficAttributes(infraProvAcc, d)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
