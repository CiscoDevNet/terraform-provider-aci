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
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"ctrl": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		})),
	}
}

func dataSourceAciSpanningTreeInterfacePolicyRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	rn := fmt.Sprintf("ifPol-%s", name)
	dn := fmt.Sprintf("uni/%s", rn)
	stpIfPol, err := getRemoteSpanningTreeInterfacePolicy(aciClient, dn)
	if err != nil {
		return err
	}
	setSpanningTreeInterfacePolicyAttributes(stpIfPol, d)
	return nil
}
