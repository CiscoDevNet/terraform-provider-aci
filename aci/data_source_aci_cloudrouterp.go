package aci

import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciCloudVpnGateway() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciCloudVpnGatewayRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"cloud_context_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
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

			"num_instances": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"cloud_router_profile_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciCloudVpnGatewayRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("routerp-%s", name)
	CloudContextProfileDn := d.Get("cloud_context_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", CloudContextProfileDn, rn)

	cloudRouterP, err := getRemoteCloudVpnGateway(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setCloudVpnGatewayAttributes(cloudRouterP, d)
	return nil
}
