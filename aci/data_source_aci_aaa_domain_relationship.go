package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciDomainRelationship() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciDomainRelationshipRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"parent_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		})),
	}
}

func dataSourceAciDomainRelationshipRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	parent_dn := d.Get("parent_dn").(string)
	rn := fmt.Sprintf("domain-%s", name)
	dn := fmt.Sprintf("%s/%s", parent_dn, rn)

	aaaDomainRef, err := getRemoteAaaDomainRelationship(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setAaaDomainRelationshipAttributes(aaaDomainRef, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
