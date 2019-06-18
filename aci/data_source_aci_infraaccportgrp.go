package aci

import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAciLeafAccessPortPolicyGroup() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciLeafAccessPortPolicyGroupRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
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

func dataSourceAciLeafAccessPortPolicyGroupRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("infra/funcprof/accportgrp-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	infraAccPortGrp, err := getRemoteLeafAccessPortPolicyGroup(aciClient, dn)

	if err != nil {
		return err
	}
	setLeafAccessPortPolicyGroupAttributes(infraAccPortGrp, d)
	return nil
}
