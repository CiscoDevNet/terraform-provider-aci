package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciPhysicalDomain() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciPhysicalDomainRead,

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

func dataSourceAciPhysicalDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("phys-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	physDomP, err := getRemotePhysicalDomain(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setPhysicalDomainAttributes(physDomP, d)

	if err != nil {
		return diag.FromErr(err)
	}

	// infraRsVlanNs - Beginning Read
	log.Printf("[DEBUG] %s: infraRsVlanNs - Beginning Read with parent DN", dn)
	_, err = getAndSetRelationinfraRsVlanNsFromPhysicalDomain(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsVlanNs - Read finished successfully", d.Get("relation_infra_rs_vlan_ns"))
	}
	// infraRsVlanNs - Read finished successfully

	// infraRsVlanNsDef - Beginning Read
	log.Printf("[DEBUG] %s: infraRsVlanNsDef - Beginning Read with parent DN", dn)
	_, err = getAndSetRelationinfraRsVlanNsDefFromPhysicalDomain(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsVlanNsDef - Read finished successfully", d.Get("relation_infra_rs_vlan_ns_def"))
	}
	// infraRsVlanNsDef - Read finished successfully

	// infraRsVipAddrNs - Beginning Read
	log.Printf("[DEBUG] %s: infraRsVipAddrNs - Beginning Read with parent DN", dn)
	_, err = getAndSetRelationinfraRsVipAddrNsFromPhysicalDomain(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsVipAddrNs - Read finished successfully", d.Get("relation_infra_rs_vip_addr_ns"))
	}
	// infraRsVipAddrNs - Read finished successfully

	// infraRsDomVxlanNsDef - Beginning Read
	log.Printf("[DEBUG] %s: infraRsDomVxlanNsDef - Beginning Read with parent DN", dn)
	_, err = getAndSetRelationinfraRsDomVxlanNsDefFromPhysicalDomain(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsDomVxlanNsDef - Read finished successfully", d.Get("relation_infra_rs_dom_vxlan_ns_def"))
	}
	// infraRsDomVxlanNsDef - Read finished successfully

	return nil
}
