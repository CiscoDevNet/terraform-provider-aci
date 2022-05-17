package aci

import (
	"context"
	"fmt"
	"log"
	"sort"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciLogicalInterface() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciLogicalInterfaceRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"l4_l7_devices_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"encap": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"lag_policy_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"relation_vns_rs_c_if_att_n": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Create relation to vns:CIf",
				Set:         schema.HashString,
			}})),
	}
}

func dataSourceAciLogicalInterfaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	L4_L7DevicesDn := d.Get("l4_l7_devices_dn").(string)
	rn := fmt.Sprintf(models.RnvnsLIf, name)
	dn := fmt.Sprintf("%s/%s", L4_L7DevicesDn, rn)

	vnsLIf, err := getRemoteLogicalInterface(aciClient, dn)
	if err != nil {
		return nil
	}

	d.SetId(dn)

	_, err = setLogicalInterfaceAttributes(vnsLIf, d)
	if err != nil {
		return nil
	}

	vnsRsCIfAttNData, err := aciClient.ReadRelationvnsRsCIfAttN(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsCIfAttN %v", err)
		d.Set("relation_vns_rs_c_if_att_n", make([]string, 0, 1))
	} else {
		vnsRsCIfAttNDataList := toStringList(vnsRsCIfAttNData.(*schema.Set).List())
		sort.Strings(vnsRsCIfAttNDataList)
		d.Set("relation_vns_rs_c_if_att_n", vnsRsCIfAttNDataList)

	}

	return nil
}
