package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciFirmwareGroup() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciFirmwareGroupRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"firmware_group_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciFirmwareGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("fabric/fwgrp-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	firmwareFwGrp, err := getRemoteFirmwareGroup(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	setFirmwareGroupAttributes(firmwareFwGrp, d)
	return nil
}
