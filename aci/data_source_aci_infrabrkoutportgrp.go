package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciLeafBreakoutPortGroup() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciLeafBreakoutPortGroupRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"brkout_map": &schema.Schema{
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

func dataSourceAciLeafBreakoutPortGroupRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("infra/funcprof/brkoutportgrp-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	infraBrkoutPortGrp, err := getRemoteLeafBreakoutPortGroup(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setLeafBreakoutPortGroupAttributes(infraBrkoutPortGrp, d)
	return nil
}
