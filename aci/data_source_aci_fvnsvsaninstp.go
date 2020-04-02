package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciVSANPool() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciVSANPoolRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"alloc_mode": &schema.Schema{
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

func dataSourceAciVSANPoolRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	allocMode := d.Get("alloc_mode").(string)

	rn := fmt.Sprintf("infra/vsanns-[%s]-%s", name, allocMode)

	dn := fmt.Sprintf("uni/%s", rn)

	fvnsVsanInstP, err := getRemoteVSANPool(aciClient, dn)

	if err != nil {
		return err
	}
	setVSANPoolAttributes(fvnsVsanInstP, d)
	return nil
}
