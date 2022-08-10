package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciL3DomainProfile() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciL3DomainProfileRead,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"relation_infra_rs_vlan_ns": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_vlan_ns_def": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_vip_addr_ns": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_extnw_rs_out": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_infra_rs_dom_vxlan_ns_def": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAciL3DomainProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("l3dom-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)
	log.Printf("[DEBUG] %s: Data Source - Beginning Read", dn)

	l3extDomP, err := getRemoteL3DomainProfile(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setL3DomainProfileAttributes(l3extDomP, d)
	if err != nil {
		return diag.FromErr(err)
	}

	// infraRsVlanNs - Beginning Read
	log.Printf("[DEBUG] %s: infraRsVlanNs - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsVlanNsFromL3DomainProfile(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsVlanNs - Read finished successfully", d.Get("relation_infra_rs_vlan_ns"))
	}
	// infraRsVlanNs - Read finished successfully

	// infraRsVlanNsDef - Beginning Read
	log.Printf("[DEBUG] %s: infraRsVlanNsDef - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsVlanNsDefFromL3DomainProfile(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsVlanNsDef - Read finished successfully", d.Get("relation_infra_rs_vlan_ns_def"))
	}
	// infraRsVlanNsDef - Read finished successfully

	// infraRsVipAddrNs - Beginning Read
	log.Printf("[DEBUG] %s: infraRsVipAddrNs - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsVipAddrNsFromL3DomainProfile(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsVipAddrNs - Read finished successfully", d.Get("relation_infra_rs_vip_addr_ns"))
	}
	// infraRsVipAddrNs - Read finished successfully

	// extnwRsOut - Beginning Read
	log.Printf("[DEBUG] %s: extnwRsOut - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationextnwRsOutFromL3DomainProfile(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: extnwRsOut - Read finished successfully", d.Get("relation_extnw_rs_out"))
	}
	// extnwRsOut - Read finished successfully

	// infraRsDomVxlanNsDef - Beginning Read
	log.Printf("[DEBUG] %s: infraRsDomVxlanNsDef - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsDomVxlanNsDefFromL3DomainProfile(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsDomVxlanNsDef - Read finished successfully", d.Get("relation_infra_rs_dom_vxlan_ns_def"))
	}
	// infraRsDomVxlanNsDef - Read finished successfully

	log.Printf("[DEBUG] %s: Data Source - Read finished successfully", dn)
	return nil
}
