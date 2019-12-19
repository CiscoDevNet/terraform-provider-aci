package aci

import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciVlanEncapsulationforVxlanTraffic() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciVlanEncapsulationforVxlanTrafficRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"attachable_access_entity_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciVlanEncapsulationforVxlanTrafficRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	rn := fmt.Sprintf("provacc")
	AttachableAccessEntityProfileDn := d.Get("attachable_access_entity_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", AttachableAccessEntityProfileDn, rn)

	infraProvAcc, err := getRemoteVlanEncapsulationforVxlanTraffic(aciClient, dn)

	if err != nil {
		return err
	}
	setVlanEncapsulationforVxlanTrafficAttributes(infraProvAcc, d)
	return nil
}
