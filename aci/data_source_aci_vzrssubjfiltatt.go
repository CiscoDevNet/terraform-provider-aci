package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
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
			"tn_vz_filter_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		}),
	}
}

func dataSourceAciSubjectFilterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	tnVzFilterName := d.Get("tn_vz_filter_name").(string)
	log.Printf("[TEST] datasrc name %v : ", tnVzFilterName)
	ContractSubjectDn := d.Get("contract_subject_dn").(string)
	log.Printf("[TEST] datasrc ContractSubjectDn %v : ", ContractSubjectDn)
	rn := fmt.Sprintf(models.RnvzRsSubjFiltAtt, tnVzFilterName)
	log.Printf("[TEST] datasrc rn %v : ", ContractSubjectDn)
	dn := fmt.Sprintf("%s/%s", ContractSubjectDn, rn)
	log.Printf("[TEST] datasrc dn %v : ", ContractSubjectDn)

	vzRsSubjFiltAtt, err := getRemoteSubjectFilter(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[TEST] vzRsSubjFiltAtt %v : ", vzRsSubjFiltAtt)

	d.SetId(dn)

	_, err = setSubjectFilterAttributes(vzRsSubjFiltAtt, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
