package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciEPgDef() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciEPgDefRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"legacy_virtual_node_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"encap": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fabric_encap": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"delete_pbr_scenario": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"logical_context_dn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"member_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"router_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func setEPgDefAttributes(vnsEPgDef *models.EPgDef, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(vnsEPgDef.DistinguishedName)
	d.Set("description", vnsEPgDef.Description)
	vnsEPgDefMap, err := vnsEPgDef.ToMap()
	if err != nil {
		return d, err
	}
	dn := d.Id()
	if dn != vnsEPgDef.DistinguishedName {
		d.Set("legacy_virtual_node_dn", "")
	} else {
		d.Set("legacy_virtual_node_dn", GetParentDn(vnsEPgDef.DistinguishedName, fmt.Sprintf("/"+models.RnvnsEPgDef, vnsEPgDefMap["name"])))
	}
	d.Set("name", vnsEPgDefMap["name"])
	d.Set("encap", vnsEPgDefMap["encap"])
	d.Set("fabric_encap", vnsEPgDefMap["fabEncap"])
	d.Set("delete_pbr_scenario", vnsEPgDefMap["isDelPbr"])
	d.Set("member_type", vnsEPgDefMap["membType"])
	d.Set("logical_context_dn", vnsEPgDefMap["lIfCtxDn"])
	d.Set("router_id", vnsEPgDefMap["rtrId"])
	return d, nil
}

func getRemoteEPgDef(client *client.Client, dn string) (*models.EPgDef, error) {
	vnsEPgDefCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	vnsEPgDef := models.EPgDefFromContainer(vnsEPgDefCont)
	if vnsEPgDef.DistinguishedName == "" {
		return nil, fmt.Errorf("EPG Def %s not found", dn)
	}
	return vnsEPgDef, nil
}

func dataSourceAciEPgDefRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	LegVNodeDn := d.Get("leg_v_node_dn").(string)
	rn := fmt.Sprintf(models.RnvnsEPgDef, name)
	dn := fmt.Sprintf("%s/%s", LegVNodeDn, rn)

	vnsEPgDef, err := getRemoteEPgDef(aciClient, dn)
	if err != nil {
		return nil
	}

	d.SetId(dn)

	_, err = setEPgDefAttributes(vnsEPgDef, d)
	if err != nil {
		return nil
	}

	return nil
}
