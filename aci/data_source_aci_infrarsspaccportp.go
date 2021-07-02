package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciInterfaceProfile() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciInterfaceProfileRead,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"spine_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				// Default:  "orchestrator:terraform",
				Computed: true,
				DefaultFunc: func() (interface{}, error) {
					return "orchestrator:terraform", nil
				},
			},

			"tdn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceAciInterfaceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	tDn := d.Get("tdn").(string)

	rn := fmt.Sprintf("rsspAccPortP-[%s]", tDn)
	SpineProfileDn := d.Get("spine_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", SpineProfileDn, rn)

	infraRsSpAccPortP, err := getRemoteInterfaceProfile(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setInterfaceProfileAttributes(infraRsSpAccPortP, d)

	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
