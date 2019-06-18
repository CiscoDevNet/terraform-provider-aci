package aci

import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAciAutonomousSystemProfile() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciAutonomousSystemProfileRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"asn": &schema.Schema{
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

func dataSourceAciAutonomousSystemProfileRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	rn := fmt.Sprintf("clouddomp/as")

	dn := fmt.Sprintf("uni/%s", rn)

	cloudBgpAsP, err := getRemoteAutonomousSystemProfile(aciClient, dn)

	if err != nil {
		return err
	}
	setAutonomousSystemProfileAttributes(cloudBgpAsP, d)
	return nil
}
