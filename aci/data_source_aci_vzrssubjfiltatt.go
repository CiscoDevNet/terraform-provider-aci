package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciSubjectFilter() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciSubjectFilterRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
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
			"filter_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
		}),
	}
}

func dataSourceAciSubjectFilterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	tnVzFilterName := GetMOName(d.Get("filter_dn").(string))
	ContractSubjectDn := d.Get("contract_subject_dn").(string)
	rn := fmt.Sprintf(models.RnvzRsSubjFiltAtt, tnVzFilterName)
	dn := fmt.Sprintf("%s/%s", ContractSubjectDn, rn)

	vzRsSubjFiltAtt, err := getRemoteSubjectFilter(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setSubjectFilterAttributes(vzRsSubjFiltAtt, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
