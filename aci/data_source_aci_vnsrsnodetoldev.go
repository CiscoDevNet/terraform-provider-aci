package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciRelationfromaAbsNodetoanLDev() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciRelationfromaAbsNodetoanLDevRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"function_node_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"t_dn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciRelationfromaAbsNodetoanLDevRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	rn := fmt.Sprintf("rsNodeToLDev")
	FunctionNodeDn := d.Get("function_node_dn").(string)

	dn := fmt.Sprintf("%s/%s", FunctionNodeDn, rn)

	vnsRsNodeToLDev, err := getRemoteRelationfromaAbsNodetoanLDev(aciClient, dn)

	if err != nil {
		return err
	}
	setRelationfromaAbsNodetoanLDevAttributes(vnsRsNodeToLDev, d)
	return nil
}
