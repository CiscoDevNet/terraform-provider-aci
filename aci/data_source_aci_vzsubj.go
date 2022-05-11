package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciContractSubject() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciContractSubjectRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"contract_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"cons_match_t": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"prio": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"prov_match_t": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"rev_flt_ports": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"target_dscp": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"apply_both_directions": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"consumer_to_provider": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"prio": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"target_dscp": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"relation_vz_rs_in_term_graph_att": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"provider_to_consumer": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"prio": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"target_dscp": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"relation_vz_rs_out_term_graph_att": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		}),
	}
}

func dataSourceAciContractSubjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("subj-%s", name)
	ContractDn := d.Get("contract_dn").(string)

	dn := fmt.Sprintf("%s/%s", ContractDn, rn)

	vzSubj, err := getRemoteContractSubject(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setContractSubjectAttributes(vzSubj, d)
	if err != nil {
		return diag.FromErr(err)
	}

	vzInTerm, err := getRemoteInTermSubject(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[TEST] DATASOURCE vzInTerm : %v ", vzInTerm)
	if vzInTerm != nil {
		abc, err := setInTermSubjectAttributes(vzInTerm, d)
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("[TEST] in DATA abc : %v ", abc)

		// vzRsInTermGraphAttData, err := aciClient.ReadRelationvzRsInTermGraphAtt(vzInTerm.DistinguishedName)
		// if err != nil {
		// 	log.Printf("[DEBUG] Error while reading relation vzRsInTermGraphAtt %v", err)
		// 	d.Set("relation_vz_rs_in_term_graph_att", "")
		// } else {
		// 	d.Set("relation_vz_rs_in_term_graph_att", vzRsInTermGraphAttData.(string))
		// }
	}

	vzOutTerm, err := getRemoteOutTermSubject(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[TEST] DATASOURCE vzOutTerm : %v ", vzOutTerm)
	if vzOutTerm != nil {
		bcd, err := setOutTermSubjectAttributes(vzOutTerm, d)
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("[TEST] in DATA bcd : %v ", bcd)
	}

	if vzInTerm == nil && vzOutTerm == nil {
		d.Set("apply_both_directions", "yes")
	} else {
		d.Set("apply_both_directions", "no")
	}

	return nil
}
