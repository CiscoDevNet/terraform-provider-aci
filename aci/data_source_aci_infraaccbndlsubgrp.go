package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciOverridePCVPCPolicyGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciOverridePCVPCPolicyGroupRead,
		SchemaVersion: 1,
		Schema: AppendAttrSchemas(map[string]*schema.Schema{
			"leaf_access_bundle_policy_group_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"port_channel_member": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		}, GetBaseAttrSchema(), GetNameAliasAttrSchema()),
	}
}

func dataSourceAciOverridePCVPCPolicyGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	LeafAccessBundlePolicyGroupDn := d.Get("leaf_access_bundle_policy_group_dn").(string)
	rn := fmt.Sprintf(models.RninfraAccBndlSubgrp, name)
	dn := fmt.Sprintf("%s/%s", LeafAccessBundlePolicyGroupDn, rn)

	infraAccBndlSubgrp, err := getRemoteOverridePCVPCPolicyGroup(aciClient, dn)
	if err != nil {
		return nil
	}

	d.SetId(dn)

	_, err = setOverridePCVPCPolicyGroupAttributes(infraAccBndlSubgrp, d)
	if err != nil {
		return nil
	}

	infraRsLacpInterfacePolData := getInfraRsLacpInterfacePolData("Read", infraAccBndlSubgrp.DistinguishedName, aciClient)
	if infraRsLacpInterfacePolData != nil {
		d.Set("port_channel_member", fmt.Sprintf("uni/infra/lacpifp-%s", infraRsLacpInterfacePolData.(string)))
	}

	return nil
}
