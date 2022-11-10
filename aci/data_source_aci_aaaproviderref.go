package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciProviderGroupMember() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciProviderGroupMemberRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"parent_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"order": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciProviderGroupMemberRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	ParentDn := d.Get("parent_dn").(string)
	rn := fmt.Sprintf("providerref-%s", name)
	dn := fmt.Sprintf("%s/%s", ParentDn, rn)
	aaaProviderRef, err := getRemoteProviderGroupMember(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setProviderGroupMemberAttributes(aaaProviderRef, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
