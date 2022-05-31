package aci

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciContractSubject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciContractSubjectCreate,
		UpdateContext: resourceAciContractSubjectUpdate,
		ReadContext:   resourceAciContractSubjectRead,
		DeleteContext: resourceAciContractSubjectDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciContractSubjectImport,
		},

		SchemaVersion: 1,

		CustomizeDiff: func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
			// ApplyBothDirections := diff.Get("apply_both_directions")
			if diff.Get("apply_both_directions") == "yes" {
				if len(diff.Get("consumer_to_provider").(*schema.Set).List()) != 0 || len(diff.Get("provider_to_consumer").(*schema.Set).List()) != 0 {
					// diff.SetNew("consumer_to_provider", nil)
					// diff.SetNew("provider_to_consumer", nil)
					return errors.New("you cannot set consumer_to_provider and provider_to_consumer when apply_both_directions is set to yes")
				}
			}

			result := false
			if diff.HasChange("consumer_to_provider") {
				old, new := diff.GetChange("consumer_to_provider")
				oldInTerm := old.(*schema.Set).List()
				newInTerm := new.(*schema.Set).List()
				if len(newInTerm) == len(oldInTerm) {
					oldInTermParam := oldInTerm[0].(map[string]interface{})
					newInTermParam := newInTerm[0].(map[string]interface{})

					if oldInTermParam["prio"] != newInTermParam["prio"] || oldInTermParam["target_dscp"] != newInTermParam["target_dscp"] || oldInTermParam["relation_vz_rs_in_term_graph_att"] != newInTermParam["relation_vz_rs_in_term_graph_att"] {
						result = true
					} else if oldInTermParam["relation_vz_rs_filt_att"] != newInTermParam["relation_vz_rs_filt_att"] {

						oldInTermFilterList := oldInTermParam["relation_vz_rs_filt_att"].(*schema.Set).List()
						newInTermFilterList := newInTermParam["relation_vz_rs_filt_att"].(*schema.Set).List()

						if len(oldInTermFilterList) != len(newInTermFilterList) {
							result = true
						} else {
							for _, oldInTermFilter := range oldInTermFilterList {
								oldVzRsFiltAttMap := oldInTermFilter.(map[string]interface{})
								oldInTermFilterDirectives := oldVzRsFiltAttMap["directives"].(*schema.Set).List()
								found_same_filter_dn := false
								for _, newInTermFilter := range newInTermFilterList {
									newVzRsFiltAttMap := newInTermFilter.(map[string]interface{})
									newInTermFilterDirectives := newVzRsFiltAttMap["directives"].(*schema.Set).List()

									if oldVzRsFiltAttMap["filter_dn"] == newVzRsFiltAttMap["filter_dn"] {
										found_same_filter_dn = true
										if oldVzRsFiltAttMap["action"] != newVzRsFiltAttMap["action"] || oldVzRsFiltAttMap["priority_override"] != newVzRsFiltAttMap["priority_override"] || !reflect.DeepEqual(oldInTermFilterDirectives, newInTermFilterDirectives) {
											if !(oldVzRsFiltAttMap["priority_override"] == "default" && newVzRsFiltAttMap["priority_override"] == "") {
												result = true
												break
											}
										}

									}
								}
								if !found_same_filter_dn {
									result = true
									break
								}
							}
						}
					}
				} else {
					result = true
				}
			}

			if diff.HasChange("provider_to_consumer") {
				oldParam, newParam := diff.GetChange("provider_to_consumer")
				oldOutTerm := oldParam.(*schema.Set).List()
				newOutTerm := newParam.(*schema.Set).List()
				if len(oldOutTerm) == len(newOutTerm) {
					oldOutTermParam := oldOutTerm[0].(map[string]interface{})
					newOutTermParam := newOutTerm[0].(map[string]interface{})
					if oldOutTermParam["prio"] != newOutTermParam["prio"] || oldOutTermParam["target_dscp"] != newOutTermParam["target_dscp"] || oldOutTermParam["relation_vz_rs_out_term_graph_att"] != newOutTermParam["relation_vz_rs_out_term_graph_att"] {
						result = true
					} else if oldOutTermParam["relation_vz_rs_filt_att"] != newOutTermParam["relation_vz_rs_filt_att"] {

						oldOutTermFilterList := oldOutTermParam["relation_vz_rs_filt_att"].(*schema.Set).List()
						newOutTermFilterList := newOutTermParam["relation_vz_rs_filt_att"].(*schema.Set).List()

						if len(oldOutTermFilterList) != len(newOutTermFilterList) {
							result = true
						} else {
							for _, oldOutTermFilter := range oldOutTermFilterList {
								oldVzRsFiltAttMap := oldOutTermFilter.(map[string]interface{})
								oldOutTermFilterDirectives := oldVzRsFiltAttMap["directives"].(*schema.Set).List()
								found_same_filter_dn := false
								for _, newOutTermFilter := range newOutTermFilterList {
									newVzRsFiltAttMap := newOutTermFilter.(map[string]interface{})
									newOutTermFilterDirectives := newVzRsFiltAttMap["directives"].(*schema.Set).List()

									if oldVzRsFiltAttMap["filter_dn"] == newVzRsFiltAttMap["filter_dn"] {
										found_same_filter_dn = true
										if oldVzRsFiltAttMap["action"] != newVzRsFiltAttMap["action"] || oldVzRsFiltAttMap["priority_override"] != newVzRsFiltAttMap["priority_override"] || !reflect.DeepEqual(oldOutTermFilterDirectives, newOutTermFilterDirectives) {
											if !(oldVzRsFiltAttMap["priority_override"] == "default" && newVzRsFiltAttMap["priority_override"] == "") {
												result = true
												break
											}
										}

									}
								}
								if !found_same_filter_dn {
									result = true
									break
								}
							}
						}
					}
				} else {
					result = true
				}
			}

			if result {
				return nil
			} else {
				diff.Clear("consumer_to_provider")
				diff.Clear("provider_to_consumer")
			}
			return nil
		},

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"contract_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cons_match_t": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"All",
					"AtleastOne",
					"AtmostOne",
					"None",
				}, false),
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
				ValidateFunc: validation.StringInSlice([]string{
					"unspecified",
					"level3",
					"level2",
					"level1",
					"level4",
					"level5",
					"level6",
				}, false),
			},

			"prov_match_t": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"All",
					"AtleastOne",
					"AtmostOne",
					"None",
				}, false),
			},

			"rev_flt_ports": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"apply_both_directions": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "yes",
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"target_dscp": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"CS0",
					"CS1",
					"AF11",
					"AF12",
					"AF13",
					"CS2",
					"AF21",
					"AF22",
					"AF23",
					"CS3",
					"CS4",
					"CS5",
					"CS6",
					"CS7",
					"AF31",
					"AF32",
					"AF33",
					"AF41",
					"AF42",
					"AF43",
					"VA",
					"EF",
					"unspecified",
				}, false),
			},

			"relation_vz_rs_subj_graph_att": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_vz_rs_sdwan_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_vz_rs_subj_filt_att": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"consumer_to_provider": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Set InTerm attributes",
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
							ValidateFunc: validation.StringInSlice([]string{
								"unspecified",
								"level3",
								"level2",
								"level1",
								"level4",
								"level5",
								"level6",
							}, false),
						},
						"target_dscp": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"CS0",
								"CS1",
								"AF11",
								"AF12",
								"AF13",
								"CS2",
								"AF21",
								"AF22",
								"AF23",
								"CS3",
								"CS4",
								"CS5",
								"CS6",
								"CS7",
								"AF31",
								"AF32",
								"AF33",
								"AF41",
								"AF42",
								"AF43",
								"VA",
								"EF",
								"unspecified",
							}, false),
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
							Description: "Create relation to vzFilter",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											"deny",
											"permit",
										}, false),
									},
									"directives": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												"log",
												"no_stats",
												"none",
											}, false),
										},
									},
									"priority_override": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											"default",
											"level1",
											"level2",
											"level3",
										}, false),
										// Used when action is deny
									},
									"filter_dn": {
										Required: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
					},
				},
			},
			"provider_to_consumer": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Set OutTerm attributes",
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
							ValidateFunc: validation.StringInSlice([]string{
								"unspecified",
								"level3",
								"level2",
								"level1",
								"level4",
								"level5",
								"level6",
							}, false),
						},
						"target_dscp": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"CS0",
								"CS1",
								"AF11",
								"AF12",
								"AF13",
								"CS2",
								"AF21",
								"AF22",
								"AF23",
								"CS3",
								"CS4",
								"CS5",
								"CS6",
								"CS7",
								"AF31",
								"AF32",
								"AF33",
								"AF41",
								"AF42",
								"AF43",
								"VA",
								"EF",
								"unspecified",
							}, false),
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
							Description: "Create relation to vzFilter",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											"deny",
											"permit",
										}, false),
									},
									"directives": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												"log",
												"no_stats",
												"none",
											}, false),
										},
									},
									"priority_override": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											"default",
											"level1",
											"level2",
											"level3",
										}, false),
										// Used when action is deny
									},
									"filter_dn": {
										Required: true,
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
func getRemoteContractSubject(client *client.Client, dn string) (*models.ContractSubject, error) {
	vzSubjCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vzSubj := models.ContractSubjectFromContainer(vzSubjCont)

	if vzSubj.DistinguishedName == "" {
		return nil, fmt.Errorf("ContractSubject %s not found", vzSubj.DistinguishedName)
	}

	return vzSubj, nil
}

func setContractSubjectAttributes(vzSubj *models.ContractSubject, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(vzSubj.DistinguishedName)
	d.Set("description", vzSubj.Description)
	if dn != vzSubj.DistinguishedName {
		d.Set("contract_dn", "")
	}
	vzSubjMap, err := vzSubj.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("contract_dn", GetParentDn(dn, fmt.Sprintf("/subj-%s", vzSubjMap["name"])))

	d.Set("name", vzSubjMap["name"])

	d.Set("annotation", vzSubjMap["annotation"])
	d.Set("cons_match_t", vzSubjMap["consMatchT"])
	d.Set("name_alias", vzSubjMap["nameAlias"])
	d.Set("prio", vzSubjMap["prio"])
	d.Set("prov_match_t", vzSubjMap["provMatchT"])
	d.Set("rev_flt_ports", vzSubjMap["revFltPorts"])
	d.Set("apply_both_directions", vzSubjMap["applyBothDirections"])
	d.Set("target_dscp", vzSubjMap["targetDscp"])
	return d, nil
}

func getRemoteInTermSubject(client *client.Client, dn string) (*models.InTermSubject, error) {

	vzInTermCont, err := client.Get(dn + "/" + models.RnvzInTerm)
	if err != nil {
		if fmt.Sprint(err) == "Error retrieving Object: Object may not exists" {
			return nil, nil
		} else {
			return nil, err
		}
	}
	vzInTerm := models.InTermSubjectFromContainer(vzInTermCont)
	if vzInTerm.DistinguishedName == "" {
		return nil, fmt.Errorf("InTermSubject child of %s not found", dn)
	}
	return vzInTerm, nil
}

func setInTermSubjectAttributes(vzInTerm *models.InTermSubject, d *schema.ResourceData) (map[string]interface{}, error) {

	vzInTermMap, err := vzInTerm.ToMap()
	if err != nil {
		return nil, err
	}

	cTopMap := make(map[string]interface{})
	cTopMap["id"] = vzInTerm.DistinguishedName
	cTopMap["prio"] = vzInTermMap["prio"]
	cTopMap["target_dscp"] = vzInTermMap["targetDscp"]

	return cTopMap, nil
}

func getRemoteOutTermSubject(client *client.Client, dn string) (*models.OutTermSubject, error) {
	vzOutTermCont, err := client.Get(dn + "/" + models.RnvzOutTerm)
	if err != nil {
		if fmt.Sprint(err) == "Error retrieving Object: Object may not exists" {
			return nil, nil
		} else {
			return nil, err
		}
	}
	vzOutTerm := models.OutTermSubjectFromContainer(vzOutTermCont)
	if vzOutTerm.DistinguishedName == "" {
		return nil, fmt.Errorf("OutTermSubject child of %s not found", dn)
	}
	return vzOutTerm, nil
}

func setOutTermSubjectAttributes(vzOutTerm *models.OutTermSubject, d *schema.ResourceData) (map[string]interface{}, error) {

	vzOutTermMap, err := vzOutTerm.ToMap()
	if err != nil {
		return nil, err
	}

	pTocMap := make(map[string]interface{})

	pTocMap["id"] = vzOutTerm.DistinguishedName
	pTocMap["prio"] = vzOutTermMap["prio"]
	pTocMap["target_dscp"] = vzOutTermMap["targetDscp"]

	return pTocMap, nil
}

func resourceAciContractSubjectImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	vzSubj, err := getRemoteContractSubject(aciClient, dn)
	if err != nil {
		return nil, err
	}

	vzSubjMap, err := vzSubj.ToMap()
	if err != nil {
		return nil, err
	}

	name := vzSubjMap["name"]
	pDN := GetParentDn(dn, fmt.Sprintf("/subj-%s", name))
	d.Set("contract_dn", pDN)

	schemaFilled, err := setContractSubjectAttributes(vzSubj, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciContractSubjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ContractSubject: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	ContractDn := d.Get("contract_dn").(string)

	vzSubjAttr := models.ContractSubjectAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vzSubjAttr.Annotation = Annotation.(string)
	} else {
		vzSubjAttr.Annotation = "{}"
	}
	if ConsMatchT, ok := d.GetOk("cons_match_t"); ok {
		vzSubjAttr.ConsMatchT = ConsMatchT.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vzSubjAttr.NameAlias = NameAlias.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		vzSubjAttr.Prio = Prio.(string)
	}
	if ProvMatchT, ok := d.GetOk("prov_match_t"); ok {
		vzSubjAttr.ProvMatchT = ProvMatchT.(string)
	}
	if RevFltPorts, ok := d.GetOk("rev_flt_ports"); ok {
		vzSubjAttr.RevFltPorts = RevFltPorts.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		vzSubjAttr.TargetDscp = TargetDscp.(string)
	}
	ApplyBothDirections := d.Get("apply_both_directions")

	vzSubj := models.NewContractSubject(fmt.Sprintf("subj-%s", name), ContractDn, desc, vzSubjAttr)

	err := aciClient.Save(vzSubj)
	if err != nil {
		return diag.FromErr(err)
	}

	if ApplyBothDirections == "no" {
		if ConsumerToProvider, ok := d.GetOk("consumer_to_provider"); ok {
			consumerToProviderSet := ConsumerToProvider.(*schema.Set).List()
			for _, val := range consumerToProviderSet {
				var inTermGraphAttribute string
				var inTermFilterAtt interface{}
				ConsumerToProviderMap := val.(map[string]interface{})
				vzInTermAttr := models.InTermSubjectAttributes{}

				vzInTermAttr.Prio = ConsumerToProviderMap["prio"].(string)
				vzInTermAttr.TargetDscp = ConsumerToProviderMap["target_dscp"].(string)
				inTermGraphAttribute = ConsumerToProviderMap["relation_vz_rs_in_term_graph_att"].(string)
				inTermFilterAtt = ConsumerToProviderMap["relation_vz_rs_filt_att"]

				vzInTerm := models.NewInTermSubject(fmt.Sprintf(models.RnvzInTerm), vzSubj.DistinguishedName, desc, "", vzInTermAttr)

				err := aciClient.Save(vzInTerm)
				if err != nil {
					return diag.FromErr(err)
				}

				if inTermGraphAttribute != "" {
					d.Partial(true)
					err = checkTDn(aciClient, []string{inTermGraphAttribute})
					if err != nil {
						return diag.FromErr(err)
					}
					d.Partial(false)

					err = aciClient.CreateRelationvzRsInTermGraphAtt(vzInTerm.DistinguishedName, vzInTermAttr.Annotation, GetMOName(inTermGraphAttribute))
					if err != nil {
						return diag.FromErr(err)
					}
				}

				inTermFilterList := inTermFilterAtt.(*schema.Set).List()
				if len(inTermFilterList) > 0 {
					checkDns := make([]string, 0, 1)

					for _, inTermFilterAttribute := range inTermFilterList {
						inTermFilterDn := inTermFilterAttribute.(map[string]interface{})["filter_dn"].(string)
						checkDns = append(checkDns, inTermFilterDn)
					}

					d.Partial(true)
					err = checkTDn(aciClient, checkDns)
					if err != nil {
						return diag.FromErr(err)
					}
					d.Partial(false)

					for _, inTermFilterAttribute := range inTermFilterList {
						ConsumerToProviderFilterMap := inTermFilterAttribute.(map[string]interface{})
						vzInTermFilterAttr := models.FilterRelationshipAttributes{}
						tnVzFilterName := GetMOName(ConsumerToProviderFilterMap["filter_dn"].(string))
						vzInTermFilterAttr.Action = ConsumerToProviderFilterMap["action"].(string)
						directivesCheck := ConsumerToProviderFilterMap["directives"]
						directivesSetList := directivesCheck.(*schema.Set).List()

						directivesList := make([]string, 0, 1)
						for _, val := range directivesSetList {
							directivesList = append(directivesList, val.(string))
						}
						Directives := strings.Join(directivesList, ",")
						vzInTermFilterAttr.Directives = Directives
						vzInTermFilterAttr.PriorityOverride = ConsumerToProviderFilterMap["priority_override"].(string)
						vzInTermFilterAttr.TnVzFilterName = tnVzFilterName

						_, err = aciClient.CreateFilterRelationship(tnVzFilterName, vzInTerm.DistinguishedName, vzInTermFilterAttr)
						if err != nil {
							return diag.FromErr(err)
						}
					}

				}

			}
		} else {

			vzInTermAttr := models.InTermSubjectAttributes{}
			vzInTerm := models.NewInTermSubject(fmt.Sprintf(models.RnvzInTerm), vzSubj.DistinguishedName, desc, "", vzInTermAttr)

			err := aciClient.Save(vzInTerm)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		if ProviderToConsumer, ok := d.GetOk("provider_to_consumer"); ok {
			providerToConsumerSet := ProviderToConsumer.(*schema.Set).List()
			for _, val := range providerToConsumerSet {
				var outTermGraphAttribute string
				var outTermFilterAtt interface{}
				ProviderToConsumerMap := val.(map[string]interface{})
				vzOutTermAttr := models.OutTermSubjectAttributes{}

				vzOutTermAttr.Prio = ProviderToConsumerMap["prio"].(string)
				vzOutTermAttr.TargetDscp = ProviderToConsumerMap["target_dscp"].(string)
				outTermGraphAttribute = ProviderToConsumerMap["relation_vz_rs_out_term_graph_att"].(string)
				outTermFilterAtt = ProviderToConsumerMap["relation_vz_rs_filt_att"]

				vzOutTerm := models.NewOutTermSubject(fmt.Sprintf(models.RnvzOutTerm), vzSubj.DistinguishedName, desc, "", vzOutTermAttr)

				err := aciClient.Save(vzOutTerm)
				if err != nil {
					return diag.FromErr(err)
				}

				if outTermGraphAttribute != "" {
					d.Partial(true)
					err = checkTDn(aciClient, []string{outTermGraphAttribute})
					if err != nil {
						return diag.FromErr(err)
					}
					d.Partial(false)

					err = aciClient.CreateRelationvzRsOutTermGraphAtt(vzOutTerm.DistinguishedName, vzOutTermAttr.Annotation, GetMOName(outTermGraphAttribute))
					if err != nil {
						return diag.FromErr(err)
					}
				}

				outTermFilterList := outTermFilterAtt.(*schema.Set).List()
				if len(outTermFilterList) > 0 {
					checkDns := make([]string, 0, 1)

					for _, outTermFilterAttribute := range outTermFilterList {
						outTermFilterDn := outTermFilterAttribute.(map[string]interface{})["filter_dn"].(string)
						checkDns = append(checkDns, outTermFilterDn)
					}

					d.Partial(true)
					err = checkTDn(aciClient, checkDns)
					if err != nil {
						return diag.FromErr(err)
					}
					d.Partial(false)

					for _, outTermFilterAttribute := range outTermFilterList {
						ProviderToConsumerFilterMap := outTermFilterAttribute.(map[string]interface{})
						vzOutTermFilterAttr := models.FilterRelationshipAttributes{}
						tnVzFilterName := GetMOName(ProviderToConsumerFilterMap["filter_dn"].(string))
						vzOutTermFilterAttr.Action = ProviderToConsumerFilterMap["action"].(string)
						directivesCheck := ProviderToConsumerFilterMap["directives"]
						directivesSetList := directivesCheck.(*schema.Set).List()

						directivesList := make([]string, 0, 1)
						for _, val := range directivesSetList {
							directivesList = append(directivesList, val.(string))
						}
						Directives := strings.Join(directivesList, ",")
						vzOutTermFilterAttr.Directives = Directives
						vzOutTermFilterAttr.PriorityOverride = ProviderToConsumerFilterMap["priority_override"].(string)
						vzOutTermFilterAttr.TnVzFilterName = tnVzFilterName

						_, err = aciClient.CreateFilterRelationship(tnVzFilterName, vzOutTerm.DistinguishedName, vzOutTermFilterAttr)
						if err != nil {
							return diag.FromErr(err)
						}
					}

				}

			}
		} else {
			vzOutTermAttr := models.OutTermSubjectAttributes{}
			vzOutTerm := models.NewOutTermSubject(fmt.Sprintf(models.RnvzOutTerm), vzSubj.DistinguishedName, desc, "", vzOutTermAttr)

			err := aciClient.Save(vzOutTerm)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	checkDns := make([]string, 0, 1)

	if relationTovzRsSubjGraphAtt, ok := d.GetOk("relation_vz_rs_subj_graph_att"); ok {
		relationParam := relationTovzRsSubjGraphAtt.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTovzRsSdwanPol, ok := d.GetOk("relation_vz_rs_sdwan_pol"); ok {
		relationParam := relationTovzRsSdwanPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTovzRsSubjFiltAtt, ok := d.GetOk("relation_vz_rs_subj_filt_att"); ok {
		relationParamList := toStringList(relationTovzRsSubjFiltAtt.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTovzRsSubjGraphAtt, ok := d.GetOk("relation_vz_rs_subj_graph_att"); ok {
		relationParam := relationTovzRsSubjGraphAtt.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationvzRsSubjGraphAttFromContractSubject(vzSubj.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationTovzRsSdwanPol, ok := d.GetOk("relation_vz_rs_sdwan_pol"); ok {
		relationParam := relationTovzRsSdwanPol.(string)
		err = aciClient.CreateRelationvzRsSdwanPolFromContractSubject(vzSubj.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationTovzRsSubjFiltAtt, ok := d.GetOk("relation_vz_rs_subj_filt_att"); ok {
		relationParamList := toStringList(relationTovzRsSubjFiltAtt.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationvzRsSubjFiltAttFromContractSubject(vzSubj.DistinguishedName, relationParamName)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(vzSubj.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciContractSubjectRead(ctx, d, m)
}

func resourceAciContractSubjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ContractSubject: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	ContractDn := d.Get("contract_dn").(string)

	vzSubjAttr := models.ContractSubjectAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vzSubjAttr.Annotation = Annotation.(string)
	} else {
		vzSubjAttr.Annotation = "{}"
	}
	if ConsMatchT, ok := d.GetOk("cons_match_t"); ok {
		vzSubjAttr.ConsMatchT = ConsMatchT.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vzSubjAttr.NameAlias = NameAlias.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		vzSubjAttr.Prio = Prio.(string)
	}
	if ProvMatchT, ok := d.GetOk("prov_match_t"); ok {
		vzSubjAttr.ProvMatchT = ProvMatchT.(string)
	}
	if RevFltPorts, ok := d.GetOk("rev_flt_ports"); ok {
		vzSubjAttr.RevFltPorts = RevFltPorts.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		vzSubjAttr.TargetDscp = TargetDscp.(string)
	}
	ApplyBothDirections := d.Get("apply_both_directions")

	// if ApplyBothDirections == "yes" {
	// 	if (d.HasChange("consumer_to_provider") || d.HasChange("provider_to_consumer")) && (d.Get("consumer_to_provider") != nil || d.Get("provider_to_consumer") != nil) {
	// 		d.Set("consumer_to_provider", nil)
	// 		d.Set("provider_to_consumer", nil)
	// 		return diag.FromErr(fmt.Errorf("you cannot set consumer_to_provider and provider_to_consumer when apply_both_directions is set to yes"))
	// 	}
	// }

	vzSubj := models.NewContractSubject(fmt.Sprintf("subj-%s", name), ContractDn, desc, vzSubjAttr)

	vzSubj.Status = "modified"

	err := aciClient.Save(vzSubj)
	if err != nil {
		return diag.FromErr(err)
	}

	if ApplyBothDirections == "no" && (d.HasChange("consumer_to_provider") || d.HasChange("provider_to_consumer")) {
		if d.HasChange("consumer_to_provider") {
			if ConsumerToProvider, ok := d.GetOk("consumer_to_provider"); ok {
				consumerToProviderSet := ConsumerToProvider.(*schema.Set).List()
				for _, val := range consumerToProviderSet {
					var inTermGraphAttribute string
					var inTermFilterAtt interface{}
					ConsumerToProviderMap := val.(map[string]interface{})
					vzInTermAttr := models.InTermSubjectAttributes{}

					vzInTermAttr.Prio = ConsumerToProviderMap["prio"].(string)
					vzInTermAttr.TargetDscp = ConsumerToProviderMap["target_dscp"].(string)
					inTermGraphAttribute = ConsumerToProviderMap["relation_vz_rs_in_term_graph_att"].(string)
					inTermFilterAtt = ConsumerToProviderMap["relation_vz_rs_filt_att"]

					vzInTerm := models.NewInTermSubject(fmt.Sprintf(models.RnvzInTerm), vzSubj.DistinguishedName, desc, "", vzInTermAttr)
					vzInTerm.Status = "modified"
					err := aciClient.Save(vzInTerm)
					if err != nil {
						return diag.FromErr(err)
					}

					if inTermGraphAttribute != "" {
						d.Partial(true)
						err = checkTDn(aciClient, []string{inTermGraphAttribute})
						if err != nil {
							return diag.FromErr(err)
						}
						d.Partial(false)

						err = aciClient.DeleteRelationvzRsInTermGraphAtt(vzInTerm.DistinguishedName)
						if err != nil {
							return diag.FromErr(err)
						}
						err = aciClient.CreateRelationvzRsInTermGraphAtt(vzInTerm.DistinguishedName, vzInTermAttr.Annotation, GetMOName(inTermGraphAttribute))
						if err != nil {
							return diag.FromErr(err)
						}
					}

					old, _ := d.GetChange("consumer_to_provider")
					oldInTermParam := old.(*schema.Set).List()[0].(map[string]interface{})
					oldInTermFilterList := oldInTermParam["relation_vz_rs_filt_att"].(*schema.Set).List()

					inTermFilterList := inTermFilterAtt.(*schema.Set).List()
					if len(inTermFilterList) > 0 {
						checkDns := make([]string, 0, 1)

						for _, inTermFilterAttribute := range inTermFilterList {
							inTermFilterDn := inTermFilterAttribute.(map[string]interface{})["filter_dn"].(string)
							checkDns = append(checkDns, inTermFilterDn)
						}

						d.Partial(true)
						err = checkTDn(aciClient, checkDns)
						if err != nil {
							return diag.FromErr(err)
						}
						d.Partial(false)

						for _, oldInTermFilterAttribute := range oldInTermFilterList {
							oldInTermFilterMap := oldInTermFilterAttribute.(map[string]interface{})
							tnVzFilterName := GetMOName(oldInTermFilterMap["filter_dn"].(string))
							err = aciClient.DeleteFilterRelationship(tnVzFilterName, vzInTerm.DistinguishedName)
							if err != nil {
								return diag.FromErr(err)
							}
						}
						for _, inTermFilterAttribute := range inTermFilterList {
							ConsumerToProviderFilterMap := inTermFilterAttribute.(map[string]interface{})
							vzInTermFilterAttr := models.FilterRelationshipAttributes{}
							tnVzFilterName := GetMOName(ConsumerToProviderFilterMap["filter_dn"].(string))
							vzInTermFilterAttr.Action = ConsumerToProviderFilterMap["action"].(string)

							directivesCheck := ConsumerToProviderFilterMap["directives"]
							directivesSetList := directivesCheck.(*schema.Set).List()
							directivesList := make([]string, 0, 1)
							for _, val := range directivesSetList {
								directivesList = append(directivesList, val.(string))
							}
							Directives := strings.Join(directivesList, ",")

							vzInTermFilterAttr.Directives = Directives
							vzInTermFilterAttr.PriorityOverride = ConsumerToProviderFilterMap["priority_override"].(string)
							vzInTermFilterAttr.TnVzFilterName = tnVzFilterName

							_, err = aciClient.CreateFilterRelationship(tnVzFilterName, vzInTerm.DistinguishedName, vzInTermFilterAttr)
							if err != nil {
								return diag.FromErr(err)
							}
						}
					}

				}
			}
		}
		if d.HasChange("provider_to_consumer") {
			if ProviderToConsumer, ok := d.GetOk("provider_to_consumer"); ok {
				providerToConsumerSet := ProviderToConsumer.(*schema.Set).List()
				for _, val := range providerToConsumerSet {
					var outTermGraphAttribute string
					var outTermFilterAtt interface{}
					ProviderToConsumerMap := val.(map[string]interface{})
					vzOutTermAttr := models.OutTermSubjectAttributes{}

					vzOutTermAttr.Prio = ProviderToConsumerMap["prio"].(string)
					vzOutTermAttr.TargetDscp = ProviderToConsumerMap["target_dscp"].(string)
					outTermGraphAttribute = ProviderToConsumerMap["relation_vz_rs_out_term_graph_att"].(string)
					outTermFilterAtt = ProviderToConsumerMap["relation_vz_rs_filt_att"]

					vzOutTerm := models.NewOutTermSubject(fmt.Sprintf(models.RnvzOutTerm), vzSubj.DistinguishedName, desc, "", vzOutTermAttr)
					vzOutTerm.Status = "modified"
					error := aciClient.Save(vzOutTerm)
					if error != nil {
						return diag.FromErr(error)
					}

					if outTermGraphAttribute != "" {
						d.Partial(true)
						err = checkTDn(aciClient, []string{outTermGraphAttribute})
						if err != nil {
							return diag.FromErr(err)
						}
						d.Partial(false)

						err = aciClient.DeleteRelationvzRsOutTermGraphAtt(vzOutTerm.DistinguishedName)
						if err != nil {
							return diag.FromErr(err)
						}
						err = aciClient.CreateRelationvzRsOutTermGraphAtt(vzOutTerm.DistinguishedName, vzOutTermAttr.Annotation, GetMOName(outTermGraphAttribute))
						if err != nil {
							return diag.FromErr(err)
						}
					}

					old, _ := d.GetChange("provider_to_consumer")
					oldOutTermParam := old.(*schema.Set).List()[0].(map[string]interface{})
					oldOutTermFilterList := oldOutTermParam["relation_vz_rs_filt_att"].(*schema.Set).List()

					outTermFilterList := outTermFilterAtt.(*schema.Set).List()
					if len(outTermFilterList) > 0 {
						checkDns := make([]string, 0, 1)

						for _, outTermFilterAttribute := range outTermFilterList {
							outTermFilterDn := outTermFilterAttribute.(map[string]interface{})["filter_dn"].(string)
							checkDns = append(checkDns, outTermFilterDn)
						}

						d.Partial(true)
						err = checkTDn(aciClient, checkDns)
						if err != nil {
							return diag.FromErr(err)
						}
						d.Partial(false)

						for _, oldOutTermFilterAttribute := range oldOutTermFilterList {
							oldOutTermFilterMap := oldOutTermFilterAttribute.(map[string]interface{})
							tnVzFilterName := GetMOName(oldOutTermFilterMap["filter_dn"].(string))
							err = aciClient.DeleteFilterRelationship(tnVzFilterName, vzOutTerm.DistinguishedName)
							if err != nil {
								return diag.FromErr(err)
							}
						}

						for _, outTermFilterAttribute := range outTermFilterList {
							ProviderToConsumerFilterMap := outTermFilterAttribute.(map[string]interface{})
							vzOutTermFilterAttr := models.FilterRelationshipAttributes{}
							tnVzFilterName := GetMOName(ProviderToConsumerFilterMap["filter_dn"].(string))
							vzOutTermFilterAttr.Action = ProviderToConsumerFilterMap["action"].(string)

							directivesCheck := ProviderToConsumerFilterMap["directives"]
							directivesSetList := directivesCheck.(*schema.Set).List()
							directivesList := make([]string, 0, 1)
							for _, val := range directivesSetList {
								directivesList = append(directivesList, val.(string))
							}
							Directives := strings.Join(directivesList, ",")

							vzOutTermFilterAttr.Directives = Directives
							vzOutTermFilterAttr.PriorityOverride = ProviderToConsumerFilterMap["priority_override"].(string)
							vzOutTermFilterAttr.TnVzFilterName = tnVzFilterName

							_, err = aciClient.CreateFilterRelationship(tnVzFilterName, vzOutTerm.DistinguishedName, vzOutTermFilterAttr)
							if err != nil {
								return diag.FromErr(err)
							}

						}

					}
				}
			}
		}
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_vz_rs_subj_graph_att") {
		_, newRelParam := d.GetChange("relation_vz_rs_subj_graph_att")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_vz_rs_sdwan_pol") {
		_, newRelParam := d.GetChange("relation_vz_rs_sdwan_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_vz_rs_subj_filt_att") {
		oldRel, newRel := d.GetChange("relation_vz_rs_subj_filt_att")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_vz_rs_subj_graph_att") {
		_, newRelParam := d.GetChange("relation_vz_rs_subj_graph_att")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationvzRsSubjGraphAttFromContractSubject(vzSubj.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvzRsSubjGraphAttFromContractSubject(vzSubj.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_vz_rs_sdwan_pol") {
		_, newRelParam := d.GetChange("relation_vz_rs_sdwan_pol")
		err = aciClient.DeleteRelationvzRsSdwanPolFromContractSubject(vzSubj.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvzRsSdwanPolFromContractSubject(vzSubj.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_vz_rs_subj_filt_att") {
		oldRel, newRel := d.GetChange("relation_vz_rs_subj_filt_att")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			relDnName := GetMOName(relDn)
			err = aciClient.DeleteRelationvzRsSubjFiltAttFromContractSubject(vzSubj.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationvzRsSubjFiltAttFromContractSubject(vzSubj.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}
		}

	}

	d.SetId(vzSubj.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciContractSubjectRead(ctx, d, m)

}

func resourceAciContractSubjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	vzSubj, err := getRemoteContractSubject(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setContractSubjectAttributes(vzSubj, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	vzInTerm, err := getRemoteInTermSubject(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}

	if vzInTerm != nil {
		vzInTermMap, err := setInTermSubjectAttributes(vzInTerm, d)
		if err != nil {
			d.SetId("")
			return nil
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
					if val != "" {
						directivesGet = append(directivesGet, strings.Trim(val, " "))
					}
				}
				inTermFilterMap["directives"] = directivesGet
				vzInTermFilterSet = append(vzInTermFilterSet, inTermFilterMap)
			}
			vzInTermMap["relation_vz_rs_filt_att"] = vzInTermFilterSet

		}

		vzInTermSet := make([]interface{}, 0, 1)
		vzInTermSet = append(vzInTermSet, vzInTermMap)
		d.Set("consumer_to_provider", vzInTermSet)

	} else {
		d.Set("consumer_to_provider", nil)
	}

	vzOutTerm, err := getRemoteOutTermSubject(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}
	if vzOutTerm != nil {

		vzOutTermMap, err := setOutTermSubjectAttributes(vzOutTerm, d)
		if err != nil {
			d.SetId("")
			return nil
		}

		vzRsOutTermGraphAttData, err := aciClient.ReadRelationvzRsOutTermGraphAtt(vzOutTerm.DistinguishedName)
		if err != nil {
			log.Printf("[DEBUG] Error while reading relation vzRsOutTermGraphAtt %v", err)
		} else {
			vzOutTermMap["relation_vz_rs_out_term_graph_att"] = vzRsOutTermGraphAttData.(string)
		}

		vzRsOutTermFilterAttData, err := aciClient.ListFilterRelationship(vzOutTerm.DistinguishedName)
		if err != nil {
			log.Printf("[DEBUG] Error while reading relation vzRsOutTermFilterAtt %v", err)
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
					if val != "" {
						directivesGet = append(directivesGet, strings.Trim(val, " "))
					}
				}
				outTermFilterMap["directives"] = directivesGet
				vzOutTermFilterSet = append(vzOutTermFilterSet, outTermFilterMap)
			}
			vzOutTermMap["relation_vz_rs_filt_att"] = vzOutTermFilterSet

		}

		vzOutTermSet := make([]interface{}, 0, 1)
		vzOutTermSet = append(vzOutTermSet, vzOutTermMap)
		d.Set("provider_to_consumer", vzOutTermSet)

	} else {
		d.Set("provider_to_consumer", nil)
	}

	if vzInTerm == nil && vzOutTerm == nil {
		d.Set("apply_both_directions", "yes")
	} else {
		d.Set("apply_both_directions", "no")
	}

	vzRsSubjGraphAttData, err := aciClient.ReadRelationvzRsSubjGraphAttFromContractSubject(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vzRsSubjGraphAtt %v", err)
		d.Set("relation_vz_rs_subj_graph_att", "")

	} else {
		setRelationAttribute(d, "relation_vz_rs_subj_graph_att", vzRsSubjGraphAttData.(string))
	}

	vzRsSdwanPolData, err := aciClient.ReadRelationvzRsSdwanPolFromContractSubject(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vzRsSdwanPol %v", err)
		d.Set("relation_vz_rs_sdwan_pol", "")

	} else {
		setRelationAttribute(d, "relation_vz_rs_sdwan_pol", vzRsSdwanPolData.(string))
	}

	vzRsSubjFiltAttData, err := aciClient.ReadRelationvzRsSubjFiltAttFromContractSubject(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vzRsSubjFiltAtt %v", err)
		setRelationAttribute(d, "relation_vz_rs_subj_filt_att", make([]interface{}, 0, 1))
	} else {
		setRelationAttribute(d, "relation_vz_rs_subj_filt_att", toStringList(vzRsSubjFiltAttData.(*schema.Set).List()))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciContractSubjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vzSubj")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
