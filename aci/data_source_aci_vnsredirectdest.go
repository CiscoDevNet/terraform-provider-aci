package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciDestinationofredirectedtraffic() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciDestinationofredirectedtrafficRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"service_redirect_policy_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"dest_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ip2": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mac": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pod_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciDestinationofredirectedtrafficRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	ip := d.Get("ip").(string)

	rn := fmt.Sprintf("RedirectDest_ip-[%s]", ip)
	ServiceRedirectPolicyDn := d.Get("service_redirect_policy_dn").(string)

	dn := fmt.Sprintf("%s/%s", ServiceRedirectPolicyDn, rn)

	vnsRedirectDest, err := getRemoteDestinationofredirectedtraffic(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setDestinationofredirectedtrafficAttributes(vnsRedirectDest, d)
	return nil
}
