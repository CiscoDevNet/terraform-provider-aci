package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciL3outPathAttachmentSecondaryIp() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciL3outPathAttachmentSecondaryIpRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"l3out_path_attachment_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"addr": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ipv6_dad": &schema.Schema{
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

func dataSourceAciL3outPathAttachmentSecondaryIpRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	addr := d.Get("addr").(string)

	rn := fmt.Sprintf("addr-[%s]", addr)
	LeafPortDn := d.Get("l3out_path_attachment_dn").(string)

	dn := fmt.Sprintf("%s/%s", LeafPortDn, rn)

	l3extIp, err := getRemoteL3outPathAttachmentSecondaryIp(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setL3outPathAttachmentSecondaryIpAttributes(l3extIp, d)
	return nil
}
