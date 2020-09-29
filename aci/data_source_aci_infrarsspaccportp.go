package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciInterfaceProfile() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciInterfaceProfileRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"spine_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"tdn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		}),
	}
}

func dataSourceAciInterfaceProfileRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	tDn := d.Get("tdn").(string)

	rn := fmt.Sprintf("rsspAccPortP-[%s]", tDn)
	SpineProfileDn := d.Get("spine_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", SpineProfileDn, rn)

	infraRsSpAccPortP, err := getRemoteInterfaceProfile(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setInterfaceProfileAttributes(infraRsSpAccPortP, d)
	return nil
}
