package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciFexBundleGroup() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciFexBundleGroupRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"fex_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciFexBundleGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("fexbundle-%s", name)
	FEXProfileDn := d.Get("fex_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", FEXProfileDn, rn)

	infraFexBndlGrp, err := getRemoteFexBundleGroup(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setFexBundleGroupAttributes(infraFexBndlGrp, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
