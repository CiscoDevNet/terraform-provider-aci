package aci

import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciCloudCIDRPool() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciCloudCIDRPoolRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"cloud_context_profile_dn": &schema.Schema{
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

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"primary": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciCloudCIDRPoolRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	addr := d.Get("addr").(string)

	rn := fmt.Sprintf("cidr-[%s]", addr)
	CloudContextProfileDn := d.Get("cloud_context_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", CloudContextProfileDn, rn)

	cloudCidr, err := getRemoteCloudCIDRPool(aciClient, dn)

	if err != nil {
		return err
	}
	setCloudCIDRPoolAttributes(cloudCidr, d)
	return nil
}
