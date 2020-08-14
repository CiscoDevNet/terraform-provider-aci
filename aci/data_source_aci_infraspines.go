package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciSwitchSpineAssociation() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciSwitchSpineAssociationRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"spine_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"spine_switch_association_type": &schema.Schema{
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

func dataSourceAciSwitchSpineAssociationRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	switchAssociationType := d.Get("spine_switch_association_type").(string)

	rn := fmt.Sprintf("spines-%s-typ-%s", name, switchAssociationType)
	SpineProfileDn := d.Get("spine_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", SpineProfileDn, rn)

	infraSpineS, err := getRemoteSwitchSpineAssociation(aciClient, dn)

	if err != nil {
		return err
	}
	setSwitchSpineAssociationAttributes(infraSpineS, d)
	return nil
}
