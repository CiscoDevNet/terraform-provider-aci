package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciISISDomainPolicy() *schema.Resource {
	return &schema.Resource{
		Read:          dataSourceAciISISDomainPolicyRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mtu": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"redistrib_metric": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciISISDomainPolicyRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)

	rn := fmt.Sprintf("fabric/isisDomP-%s", name)
	dn := fmt.Sprintf("uni/%s", rn)
	isisDomPol, err := getRemoteISISDomainPolicy(aciClient, dn)
	if err != nil {
		return err
	}
	d.SetId(dn)
	setISISDomainPolicyAttributes(isisDomPol, d)
	return nil
}
