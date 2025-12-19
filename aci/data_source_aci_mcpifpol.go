package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciMiscablingProtocolInterfacePolicy() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "The datasource 'aci_miscabling_protocol_interface_policy' is deprecated, please refer to 'aci_mcp_interface_policy' instead. The datasource will be removed in the next major version of the provider.",

		ReadContext: dataSourceAciMiscablingProtocolInterfacePolicyRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"admin_st": &schema.Schema{
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

func dataSourceAciMiscablingProtocolInterfacePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("infra/mcpIfP-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	mcpIfPol, err := getRemoteMiscablingProtocolInterfacePolicy(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setMiscablingProtocolInterfacePolicyAttributes(mcpIfPol, d)

	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
