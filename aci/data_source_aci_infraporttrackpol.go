package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciPortTracking() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciPortTrackingReadContext,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"admin_st": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"delay": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"include_apic_ports": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"minlinks": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciPortTrackingReadContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := "default"

	rn := fmt.Sprintf("infra/trackEqptFabP-%s", name)
	dn := fmt.Sprintf("uni/%s", rn)
	infraPortTrackPol, err := GetRemotePortTracking(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)

	_, err = setPortTrackingAttributes(infraPortTrackPol, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
