package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciDestinationofredirectedtraffic() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciDestinationofredirectedtrafficRead,

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

func dataSourceAciDestinationofredirectedtrafficRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	ip := d.Get("ip").(string)

	rn := fmt.Sprintf("RedirectDest_ip-[%s]", ip)
	ServiceRedirectPolicyDn := d.Get("service_redirect_policy_dn").(string)

	dn := fmt.Sprintf("%s/%s", ServiceRedirectPolicyDn, rn)

	vnsRedirectDest, err := getRemoteDestinationofredirectedtraffic(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setDestinationofredirectedtrafficAttributes(vnsRedirectDest, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
