package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciAccessPortSelector() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciAccessPortSelectorRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"leaf_interface_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"access_port_selector_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"relation_infra_rs_acc_base_grp": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		}),
	}
}

func dataSourceAciAccessPortSelectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	access_port_selector_type := d.Get("access_port_selector_type").(string)

	rn := fmt.Sprintf("hports-%s-typ-%s", name, access_port_selector_type)
	LeafInterfaceProfileDn := d.Get("leaf_interface_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", LeafInterfaceProfileDn, rn)
	log.Printf("[DEBUG] %s: Data Source - Beginning Read", dn)

	infraHPortS, err := getRemoteAccessPortSelector(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setAccessPortSelectorAttributes(infraHPortS, d)
	if err != nil {
		return diag.FromErr(err)
	}

	// infraRsAccBaseGrp - Beginning Read
	log.Printf("[DEBUG] %s: infraRsAccBaseGrp - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsAccBaseGrpFromAccessPortSelector(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsAccBaseGrp - Read finished successfully", d.Get("relation_infra_rs_acc_base_grp"))
	}
	// infraRsAccBaseGrp - Read finished successfully

	log.Printf("[DEBUG] %s: Data Source - Read finished successfully", dn)
	return nil
}
