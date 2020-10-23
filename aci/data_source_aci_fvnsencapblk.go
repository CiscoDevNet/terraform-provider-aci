package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciRanges() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciRangesRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"vlan_pool_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"_from": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"to": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"alloc_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"from": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"role": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciRangesRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	_from := d.Get("_from").(string)

	to := d.Get("to").(string)

	rn := fmt.Sprintf("from-[%s]-to-[%s]", _from, to)
	VLANPoolDn := d.Get("vlan_pool_dn").(string)

	dn := fmt.Sprintf("%s/%s", VLANPoolDn, rn)

	fvnsEncapBlk, err := getRemoteRanges(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setRangesAttributes(fvnsEncapBlk, d)
	return nil
}
