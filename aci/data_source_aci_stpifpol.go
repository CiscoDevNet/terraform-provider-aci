package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciSpanningTreeInterfacePolicy() *schema.Resource {
	return &schema.Resource{
		Read:          dataSourceAciSpanningTreeInterfacePolicyRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"annotation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ctrl": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		})),
	}
}

func dataSourceAciSpanningTreeInterfacePolicyRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)

	rn := fmt.Sprintf("infra/ifPol-%s", name)
	dn := fmt.Sprintf("uni/%s", rn)
	stpIfPol, err := getRemoteSpanningTreeInterfacePolicy(aciClient, dn)
	if err != nil {
		return err
	}
	setSpanningTreeInterfacePolicyAttributes(stpIfPol, d)
	return nil
}
