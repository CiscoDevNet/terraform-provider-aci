package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciConfigurationExportPolicy() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciConfigurationExportPolicyRead,

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

			"format": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"include_secure_fields": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"max_snapshot_count": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"snapshot": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"target_dn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciConfigurationExportPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("fabric/configexp-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	configExportP, err := getRemoteConfigurationExportPolicy(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setConfigurationExportPolicyAttributes(configExportP, d)

	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
