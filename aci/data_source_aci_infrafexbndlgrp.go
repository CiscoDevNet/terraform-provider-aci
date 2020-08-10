package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciFexBundleGroup() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciFexBundleGroupRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"fex_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
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

func dataSourceAciFexBundleGroupRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("fexbundle-%s", name)
	FEXProfileDn := d.Get("fex_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", FEXProfileDn, rn)

	infraFexBndlGrp, err := getRemoteFexBundleGroup(aciClient, dn)

	if err != nil {
		return err
	}
	setFexBundleGroupAttributes(infraFexBndlGrp, d)
	return nil
}
