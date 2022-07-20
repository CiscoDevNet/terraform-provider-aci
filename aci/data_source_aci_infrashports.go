package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciSpineAccessPortSelector() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciSpineAccessPortSelectorRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"spine_interface_profile_dn": {
				Type:     schema.TypeString,
				Required: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"spine_access_port_selector_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_sp_acc_grp": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to infra:SpAccGrp",
			},
		}),
	}
}

func dataSourceAciSpineAccessPortSelectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	spine_access_port_selector_type := d.Get("spine_access_port_selector_type").(string)

	rn := fmt.Sprintf("shports-%s-typ-%s", name, spine_access_port_selector_type)
	SpineInterfaceProfileDn := d.Get("spine_interface_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", SpineInterfaceProfileDn, rn)

	infraSHPortS, err := getRemoteSpineAccessPortSelector(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setSpineAccessPortSelectorAttributes(infraSHPortS, d)
	if err != nil {
		return diag.FromErr(err)
	}

	// infraRsSpAccGrp - Beginning Read
	log.Printf("[DEBUG] %s: infraRsSpAccGrp - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsSpAccGrp(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsSpAccGrp - Read finished successfully", d.Get("relation_infra_rs_sp_acc_grp"))
	}
	// infraRsSpAccGrp - Read finished successfully

	return nil
}
