package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciSPANSourcedestinationGroupMatchLabel() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciSPANSourcedestinationGroupMatchLabelRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"span_source_group_dn": &schema.Schema{
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

			"tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciSPANSourcedestinationGroupMatchLabelRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("spanlbl-%s", name)
	SPANSourceGroupDn := d.Get("span_source_group_dn").(string)

	dn := fmt.Sprintf("%s/%s", SPANSourceGroupDn, rn)

	spanSpanLbl, err := getRemoteSPANSourcedestinationGroupMatchLabel(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setSPANSourcedestinationGroupMatchLabelAttributes(spanSpanLbl, d)
	return nil
}
