package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciFilterRelationship() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciFilterRelationshipRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"contract_subject_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"directives": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"priority_override": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tn_vz_filter_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filter_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
		})),
	}
}

func dataSourceAciFilterRelationshipRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	tnVzFilterName := GetMOName(d.Get("filter_dn").(string))
	ContractSubjectDn := d.Get("contract_subject_dn").(string)
	rn := fmt.Sprintf(models.RnvzRsFiltAtt, tnVzFilterName)
	dn := fmt.Sprintf("%s/%s", ContractSubjectDn, rn)

	vzRsFiltAtt, err := getRemoteFilter(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setFilterAttributes(vzRsFiltAtt, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
