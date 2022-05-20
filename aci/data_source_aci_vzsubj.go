package aci

import (
	"context"
	"fmt"
	"log"
	"strings"

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
						"relation_vz_rs_filt_att": {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Description: "Represent a relation to vzFilter",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeString,
									},
									"directives": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"priority_override": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeString,
									},
									"filter_dn": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeString,
									},
								},
							},
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
						"relation_vz_rs_filt_att": {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Description: "Represent a relation to vzFilter",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeString,
									},
									"directives": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"priority_override": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeString,
									},
									"filter_dn": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeString,
									},
								},
							},
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
		vzInTermMap, err := setInTermSubjectAttributes(vzInTerm, d)
		if err != nil {
			return diag.FromErr(err)
		}

		vzRsInTermGraphAttData, err := aciClient.ReadRelationvzRsInTermGraphAtt(vzInTerm.DistinguishedName)
		if err != nil {
			log.Printf("[DEBUG] Error while reading relation vzRsInTermGraphAtt %v", err)
		} else {
			vzInTermMap["relation_vz_rs_in_term_graph_att"] = vzRsInTermGraphAttData.(string)
		}

		vzRsInTermFilterAttData, err := aciClient.ListFilterRelationship(vzInTerm.DistinguishedName)
		if err != nil {
			log.Printf("[DEBUG] Error while reading relation vzRsInTermFilterAtt %v", err)
		} else {
			vzInTermFilterSet := make([]interface{}, 0, 1)
			for _, vzRsInTermFilter := range vzRsInTermFilterAttData {
				vzRsFiltAttMap, err := vzRsInTermFilter.ToMap()
				if err != nil {
					log.Printf("[DEBUG] Error while creating map for relation vzRsInTermFilterAtt %v", err)
				}

				inTermFilterMap := make(map[string]interface{})
				inTermFilterMap["action"] = vzRsFiltAttMap["action"]
				inTermFilterMap["priority_override"] = vzRsFiltAttMap["priorityOverride"]
				inTermFilterMap["filter_dn"] = vzRsFiltAttMap["tDn"]
				directivesGet := make([]string, 0, 1)
				for _, val := range strings.Split(vzRsFiltAttMap["directives"], ",") {
					directivesGet = append(directivesGet, strings.Trim(val, " "))
				}
				inTermFilterMap["directives"] = directivesGet
				vzInTermFilterSet = append(vzInTermFilterSet, inTermFilterMap)
			}
			vzInTermMap["relation_vz_rs_filt_att"] = vzInTermFilterSet

		}

		vzInTermSet := make([]interface{}, 0, 1)
		vzInTermSet = append(vzInTermSet, vzInTermMap)
		d.Set("consumer_to_provider", vzInTermSet)
	}

	vzOutTerm, err := getRemoteOutTermSubject(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	if vzOutTerm != nil {
		vzOutTermMap, err := setOutTermSubjectAttributes(vzOutTerm, d)
		if err != nil {
			return diag.FromErr(err)
		}

		vzRsOutTermGraphAttData, err := aciClient.ReadRelationvzRsOutTermGraphAtt(vzOutTerm.DistinguishedName)
		if err != nil {
			log.Printf("[DEBUG] Error while reading relation vzRsOutTermGraphAtt %v", err)
		} else {
			vzOutTermMap["relation_vz_rs_out_term_graph_att"] = vzRsOutTermGraphAttData.(string)
		}

		vzRsOutTermFilterAttData, err := aciClient.ListFilterRelationship(vzOutTerm.DistinguishedName)
		if err != nil {
			log.Printf("[DEBUG] Error while reading relation vzRsInTermFilterAtt %v", err)
		} else {
			vzOutTermFilterSet := make([]interface{}, 0, 1)
			for _, vzRsOutTermFilter := range vzRsOutTermFilterAttData {
				vzRsFiltAttMap, err := vzRsOutTermFilter.ToMap()
				if err != nil {
					log.Printf("[DEBUG] Error while creating map for relation vzRsOutTermFilterAtt %v", err)
				}

				outTermFilterMap := make(map[string]interface{})
				outTermFilterMap["action"] = vzRsFiltAttMap["action"]
				outTermFilterMap["priority_override"] = vzRsFiltAttMap["priorityOverride"]
				outTermFilterMap["filter_dn"] = vzRsFiltAttMap["tDn"]
				directivesGet := make([]string, 0, 1)
				for _, val := range strings.Split(vzRsFiltAttMap["directives"], ",") {
					directivesGet = append(directivesGet, strings.Trim(val, " "))
				}
				outTermFilterMap["directives"] = directivesGet
				vzOutTermFilterSet = append(vzOutTermFilterSet, outTermFilterMap)
			}
			vzOutTermMap["relation_vz_rs_filt_att"] = vzOutTermFilterSet
		}
		vzOutTermSet := make([]interface{}, 0, 1)
		vzOutTermSet = append(vzOutTermSet, vzOutTermMap)
		d.Set("provider_to_consumer", vzOutTermSet)
	}

	if vzInTerm == nil && vzOutTerm == nil {
		d.Set("apply_both_directions", "yes")
	} else {
		d.Set("apply_both_directions", "no")
	}

	return nil
}
