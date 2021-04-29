package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciSpineAccessPortPolicyGroup() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciSpineAccessPortPolicyGroupRead,

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
		}),
	}
}

func dataSourceAciSpineAccessPortPolicyGroupRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("infra/funcprof/spaccportgrp-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	infraSpAccPortGrp, err := getRemoteSpineAccessPortPolicyGroup(aciClient, dn)

	if err != nil {
		return err
	}
	setSpineAccessPortPolicyGroupAttributes(infraSpAccPortGrp, d)
	return nil
}
