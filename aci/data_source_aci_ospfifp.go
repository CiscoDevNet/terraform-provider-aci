package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciOSPFInterfaceProfile() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciOSPFInterfaceProfileRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"logical_interface_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"auth_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"auth_key_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"auth_type": &schema.Schema{
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

func dataSourceAciOSPFInterfaceProfileRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	rn := fmt.Sprintf("ospfIfP")
	LogicalInterfaceProfileDn := d.Get("logical_interface_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", LogicalInterfaceProfileDn, rn)

	ospfIfP, err := getRemoteOSPFInterfaceProfile(aciClient, dn)

	if err != nil {
		return err
	}

	d.SetId(dn)
	setOSPFInterfaceProfileAttributes(ospfIfP, d)
	return nil
}
