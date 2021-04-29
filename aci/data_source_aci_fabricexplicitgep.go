package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciVPCExplicitProtectionGroup() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciVPCExplicitProtectionGroupRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"vpc_explicit_protection_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciVPCExplicitProtectionGroupRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("fabric/protpol/expgep-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	fabricExplicitGEp, err := getRemoteVPCExplicitProtectionGroup(aciClient, dn)

	if err != nil {
		return err
	}
	setVPCExplicitProtectionGroupAttributes(fabricExplicitGEp, d)
	return nil
}
