package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciLinkLevelPolicy() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "The datasource 'aci_fabric_if_pol' is deprecated, please refer to 'aci_link_level_interface_policy' instead. The datasource will be removed in the next major version of the provider.",

		ReadContext: dataSourceAciLinkLevelPolicyRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"auto_neg": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"fec_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"link_debounce": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"speed": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciLinkLevelPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("infra/hintfpol-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	fabricHIfPol, err := getRemoteLinkLevelPolicy(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setLinkLevelPolicyAttributes(fabricHIfPol, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
