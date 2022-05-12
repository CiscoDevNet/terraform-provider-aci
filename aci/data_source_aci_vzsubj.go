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

	if vzInTerm != nil {
		var vzRsInTermGraphAtt string
		vzRsInTermGraphAttData, err := aciClient.ReadRelationvzRsInTermGraphAtt(vzInTerm.DistinguishedName)
		if err != nil {
			log.Printf("[DEBUG] Error while reading relation vzRsInTermGraphAtt %v", err)
		} else {
			vzRsInTermGraphAtt = vzRsInTermGraphAttData.(string)
		}

		_, err = setInTermSubjectAttributes(vzInTerm, d, vzRsInTermGraphAtt)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	vzOutTerm, err := getRemoteOutTermSubject(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	if vzOutTerm != nil {
		var vzRsOutTermGraphAtt string
		vzRsOutTermGraphAttData, err := aciClient.ReadRelationvzRsOutTermGraphAtt(vzOutTerm.DistinguishedName)
		if err != nil {
			log.Printf("[DEBUG] Error while reading relation vzRsOutTermGraphAtt %v", err)
		} else {
			vzRsOutTermGraphAtt = vzRsOutTermGraphAttData.(string)
		}

		_, err = setOutTermSubjectAttributes(vzOutTerm, d, vzRsOutTermGraphAtt)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if vzInTerm == nil && vzOutTerm == nil {
		d.Set("apply_both_directions", "yes")
	} else {
		d.Set("apply_both_directions", "no")
	}

	return nil
}
