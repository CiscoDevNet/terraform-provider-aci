package aci

import (
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciContract() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciContractCreate,
		Update: resourceAciContractUpdate,
		Read:   resourceAciContractRead,
		Delete: resourceAciContractDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciContractImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"tenant_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"annotation": &schema.Schema{
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

			"scope": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"target_dscp": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_vz_rs_graph_att": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},

			"filter": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"filter_name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"description": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"annotation": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"name_alias": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"filter_entry": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},

									"filter_entry_name": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},

									"entry_description": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},

									"entry_annotation": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},

									"apply_to_frag": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},

									"arp_opc": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},

									"d_from_port": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											constantPortMapping := map[string]string{
												"smtp":        "25",
												"dns":         "53",
												"http":        "80",
												"https":       "443",
												"pop3":        "110",
												"rtsp":        "554",
												"ftpData":     "20",
												"ssh":         "22",
												"unspecified": "0",
											}
											if old != "" {
												if constantPortMapping[new] == old {
													return true
												}
											}
											return false
										},
									},

									"d_to_port": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},

									"ether_t": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},

									"icmpv4_t": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},

									"icmpv6_t": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},

									"match_dscp": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},

									"entry_name_alias": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},

									"prot": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},

									"s_from_port": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},

									"s_to_port": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},

									"stateful": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},

									"tcp_rules": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"filter_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},

			"filter_entry_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		}),
	}
}
func getRemoteContract(client *client.Client, dn string) (*models.Contract, error) {
	vzBrCPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vzBrCP := models.ContractFromContainer(vzBrCPCont)

	if vzBrCP.DistinguishedName == "" {
		return nil, fmt.Errorf("Contract %s not found", vzBrCP.DistinguishedName)
	}

	return vzBrCP, nil
}

func getRemoteFilterFromContract(client *client.Client, dn string) (*models.Filter, error) {
	vzFilterCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vzFilter := models.FilterFromContainer(vzFilterCont)

	if vzFilter.DistinguishedName == "" {
		return nil, fmt.Errorf("Filter %s not found", vzFilter.DistinguishedName)
	}

	return vzFilter, nil
}

func getRemoteFilterEntryFromContract(client *client.Client, dn string) (*models.FilterEntry, error) {
	vzEntryCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vzEntry := models.FilterEntryFromContainer(vzEntryCont)

	if vzEntry.DistinguishedName == "" {
		return nil, fmt.Errorf("FilterEntry %s not found", vzEntry.DistinguishedName)
	}

	return vzEntry, nil
}

func setContractAttributes(vzBrCP *models.Contract, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(vzBrCP.DistinguishedName)
	d.Set("description", vzBrCP.Description)
	// d.Set("tenant_dn", GetParentDn(vzBrCP.DistinguishedName))
	if dn != vzBrCP.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	vzBrCPMap, _ := vzBrCP.ToMap()

	d.Set("name", vzBrCPMap["name"])

	d.Set("annotation", vzBrCPMap["annotation"])
	d.Set("name_alias", vzBrCPMap["nameAlias"])
	d.Set("prio", vzBrCPMap["prio"])
	d.Set("scope", vzBrCPMap["scope"])
	d.Set("target_dscp", vzBrCPMap["targetDscp"])
	return d
}

func setFilterAttributesFromContract(vzfilters []*models.Filter, vzEntries []*models.FilterEntry, d *schema.ResourceData) *schema.ResourceData {
	log.Println("Check .... :", vzfilters)
	log.Println("Check ... Filter :", vzEntries)
	filterSet := make([]interface{}, 0, 1)
	for _, filter := range vzfilters {
		fMap := make(map[string]interface{})
		fMap["description"] = filter.Description
		fMap["id"] = filter.DistinguishedName

		vzFilterMap, _ := filter.ToMap()
		fMap["filter_name"] = vzFilterMap["name"]
		fMap["annotation"] = vzFilterMap["annotation"]
		fMap["name_alias"] = vzFilterMap["nameAlias"]

		entrySet := make([]interface{}, 0, 1)
		for _, entry := range vzEntries {
			if strings.Contains(entry.DistinguishedName, filter.DistinguishedName) {
				entryMap := setFilterEntryAttributesFromContract(entry, d)
				entrySet = append(entrySet, entryMap)
			}
		}
		fMap["filter_entry"] = entrySet
		filterSet = append(filterSet, fMap)
	}
	log.Println("Check ...:", filterSet)
	d.Set("filter", filterSet)
	return d
}

func setFilterEntryAttributesFromContract(vzentry *models.FilterEntry, d *schema.ResourceData) map[string]interface{} {
	eMap := make(map[string]interface{})
	constantPortMapping := map[string]string{
		"smtp":        "25",
		"dns":         "53",
		"http":        "80",
		"https":       "443",
		"pop3":        "110",
		"rtsp":        "554",
		"ftpData":     "20",
		"ssh":         "22",
		"unspecified": "0",
	}
	eMap["id"] = vzentry.DistinguishedName
	eMap["entry_description"] = vzentry.Description

	vzEntryMap, _ := vzentry.ToMap()
	eMap["filter_entry_name"] = vzEntryMap["name"]
	eMap["entry_annotation"] = vzEntryMap["annotation"]
	eMap["apply_to_frag"] = vzEntryMap["applyToFrag"]
	eMap["arp_opc"] = vzEntryMap["arpOpc"]
	eMap["ether_t"] = vzEntryMap["etherT"]
	eMap["icmpv4_t"] = vzEntryMap["icmpv4T"]
	eMap["icmpv6_t"] = vzEntryMap["icmpv6T"]
	eMap["match_dscp"] = vzEntryMap["matchDscp"]
	eMap["entry_name_alias"] = vzEntryMap["nameAlias"]
	eMap["prot"] = vzEntryMap["prot"]
	eMap["s_from_port"] = vzEntryMap["sFromPort"]
	eMap["s_to_port"] = vzEntryMap["sToPort"]
	eMap["stateful"] = vzEntryMap["stateful"]
	eMap["tcp_rules"] = vzEntryMap["tcpRules"]
	eMap["d_from_port"] = constantPortMapping[vzEntryMap["dFromPort"]]
	eMap["d_to_port"] = constantPortMapping[vzEntryMap["dToPort"]]
	return eMap
}

func resourceAciContractImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	vzBrCP, err := getRemoteContract(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setContractAttributes(vzBrCP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciContractCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Contract: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	vzBrCPAttr := models.ContractAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vzBrCPAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vzBrCPAttr.NameAlias = NameAlias.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		vzBrCPAttr.Prio = Prio.(string)
	}
	if Scope, ok := d.GetOk("scope"); ok {
		vzBrCPAttr.Scope = Scope.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		vzBrCPAttr.TargetDscp = TargetDscp.(string)
	}
	vzBrCP := models.NewContract(fmt.Sprintf("brc-%s", name), TenantDn, desc, vzBrCPAttr)

	err := aciClient.Save(vzBrCP)
	if err != nil {
		return err
	}

	filterIDS := make([]string, 0, 1)
	filterentryIDS := make([]string, 0, 1)
	if filters, ok := d.GetOk("filter"); ok {
		// filterSet := make([]interface{}, 0, 1)
		vzfilters := filters.([]interface{})
		for _, val := range vzfilters {
			vzFilterAttr := models.FilterAttributes{}
			filter := val.(map[string]interface{})

			name := filter["filter_name"].(string)

			desc := filter["description"].(string)

			if filter["annotation"] != nil {
				vzFilterAttr.Annotation = filter["annotation"].(string)
			}

			if filter["name_alias"] != nil {
				vzFilterAttr.NameAlias = filter["name_alias"].(string)
			}

			vzFilter := models.NewFilter(fmt.Sprintf("flt-%s", name), TenantDn, desc, vzFilterAttr)

			err := aciClient.Save(vzFilter)
			if err != nil {
				return err
			}

			if filter["filter_entry"] != nil {
				vzfilterentries := filter["filter_entry"].([]interface{})
				log.Println("Filter entries ... :", vzfilterentries)
				for _, entry := range vzfilterentries {
					vzEntryAttr := models.FilterEntryAttributes{}
					vzEntry := entry.(map[string]interface{})

					log.Println("Entries ......... :", vzEntry)
					entryDesc := vzEntry["entry_description"].(string)

					entryName := vzEntry["filter_entry_name"].(string)

					filterDn := vzFilter.DistinguishedName

					if vzEntry["entry_annotation"] != nil {
						vzEntryAttr.Annotation = vzEntry["entry_annotation"].(string)
					}
					if vzEntry["apply_to_frag"] != nil {
						vzEntryAttr.ApplyToFrag = vzEntry["apply_to_frag"].(string)
					}
					if vzEntry["arp_opc"] != nil {
						vzEntryAttr.ArpOpc = vzEntry["arp_opc"].(string)
					}
					if vzEntry["d_from_port"] != nil {
						vzEntryAttr.DFromPort = vzEntry["d_from_port"].(string)
					}
					if vzEntry["d_to_port"] != nil {
						vzEntryAttr.DToPort = vzEntry["d_to_port"].(string)
					}
					if vzEntry["ether_t"] != nil {
						vzEntryAttr.EtherT = vzEntry["ether_t"].(string)
					}
					if vzEntry["icmpv4_t"] != nil {
						vzEntryAttr.Icmpv4T = vzEntry["icmpv4_t"].(string)
					}
					if vzEntry["icmpv6_t"] != nil {
						vzEntryAttr.Icmpv6T = vzEntry["icmpv6_t"].(string)
					}
					if vzEntry["match_dscp"] != nil {
						vzEntryAttr.MatchDscp = vzEntry["match_dscp"].(string)
					}
					if vzEntry["entry_name_alias"] != nil {
						vzEntryAttr.NameAlias = vzEntry["entry_name_alias"].(string)
					}
					if vzEntry["prot"] != nil {
						vzEntryAttr.Prot = vzEntry["prot"].(string)
					}
					if vzEntry["s_from_port"] != nil {
						vzEntryAttr.SFromPort = vzEntry["s_from_port"].(string)
					}
					if vzEntry["s_to_port"] != nil {
						vzEntryAttr.SToPort = vzEntry["s_to_port"].(string)
					}
					if vzEntry["stateful"] != nil {
						vzEntryAttr.Stateful = vzEntry["stateful"].(string)
					}
					if vzEntry["tcp_rules"] != nil {
						vzEntryAttr.TcpRules = vzEntry["tcp_rules"].(string)
					}

					vzFilterEntry := models.NewFilterEntry(fmt.Sprintf("e-%s", entryName), filterDn, entryDesc, vzEntryAttr)
					err := aciClient.Save(vzFilterEntry)
					if err != nil {
						return err
					}

					filterentryIDS = append(filterentryIDS, vzFilterEntry.DistinguishedName)
				}

			}

			// fMap := make(map[string]interface{})
			// fMap["id"] = vzFilter.DistinguishedName
			filterIDS = append(filterIDS, vzFilter.DistinguishedName)
		}
		log.Println("Check ... :", filterIDS)
		d.Set("filter_ids", filterIDS)
		d.Set("filter_entry_ids", filterentryIDS)
	} else {
		d.Set("filter_ids", filterIDS)
		d.Set("filter_entry_ids", filterentryIDS)
	}

	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if relationTovzRsGraphAtt, ok := d.GetOk("relation_vz_rs_graph_att"); ok {
		relationParam := relationTovzRsGraphAtt.(string)
		err = aciClient.CreateRelationvzRsGraphAttFromContract(vzBrCP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vz_rs_graph_att")
		d.Partial(false)

	}

	d.SetId(vzBrCP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciContractRead(d, m)
}

func resourceAciContractUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Contract: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	vzBrCPAttr := models.ContractAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vzBrCPAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vzBrCPAttr.NameAlias = NameAlias.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		vzBrCPAttr.Prio = Prio.(string)
	}
	if Scope, ok := d.GetOk("scope"); ok {
		vzBrCPAttr.Scope = Scope.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		vzBrCPAttr.TargetDscp = TargetDscp.(string)
	}
	vzBrCP := models.NewContract(fmt.Sprintf("brc-%s", name), TenantDn, desc, vzBrCPAttr)

	vzBrCP.Status = "modified"

	err := aciClient.Save(vzBrCP)

	if err != nil {
		return err
	}

	if d.HasChange("filter") {
		filter := d.Get("filter_ids").([]interface{})
		for _, val := range filter {
			filterDN := val.(string)
			err := aciClient.DeleteByDn(filterDN, "vzFilter")
			if err != nil {
				return err
			}
		}

		filters := d.Get("filter")
		filterIDS := make([]string, 0, 1)
		filterentryIDS := make([]string, 0, 1)
		vzfilters := filters.([]interface{})
		for _, val := range vzfilters {
			vzFilterAttr := models.FilterAttributes{}
			filter := val.(map[string]interface{})

			name := filter["filter_name"].(string)

			desc := filter["description"].(string)

			if filter["annotation"] != nil {
				vzFilterAttr.Annotation = filter["annotation"].(string)
			}

			if filter["name_alias"] != nil {
				vzFilterAttr.NameAlias = filter["name_alias"].(string)
			}

			vzFilter := models.NewFilter(fmt.Sprintf("flt-%s", name), TenantDn, desc, vzFilterAttr)

			// vzFilter.Status = "modified"
			err := aciClient.Save(vzFilter)
			if err != nil {
				return err
			}

			if filter["filter_entry"] != nil {
				vzfilterentries := filter["filter_entry"].([]interface{})
				log.Println("Filter entries ... :", vzfilterentries)
				for _, entry := range vzfilterentries {
					vzEntryAttr := models.FilterEntryAttributes{}
					vzEntry := entry.(map[string]interface{})

					log.Println("Entries ......... :", vzEntry)
					entryDesc := vzEntry["entry_description"].(string)

					entryName := vzEntry["filter_entry_name"].(string)

					filterDn := vzFilter.DistinguishedName

					if vzEntry["entry_annotation"] != nil {
						vzEntryAttr.Annotation = vzEntry["entry_annotation"].(string)
					}
					if vzEntry["apply_to_frag"] != nil {
						vzEntryAttr.ApplyToFrag = vzEntry["apply_to_frag"].(string)
					}
					if vzEntry["arp_opc"] != nil {
						vzEntryAttr.ArpOpc = vzEntry["arp_opc"].(string)
					}
					if vzEntry["d_from_port"] != nil {
						vzEntryAttr.DFromPort = vzEntry["d_from_port"].(string)
					}
					if vzEntry["d_to_port"] != nil {
						vzEntryAttr.DToPort = vzEntry["d_to_port"].(string)
					}
					if vzEntry["ether_t"] != nil {
						vzEntryAttr.EtherT = vzEntry["ether_t"].(string)
					}
					if vzEntry["icmpv4_t"] != nil {
						vzEntryAttr.Icmpv4T = vzEntry["icmpv4_t"].(string)
					}
					if vzEntry["icmpv6_t"] != nil {
						vzEntryAttr.Icmpv6T = vzEntry["icmpv6_t"].(string)
					}
					if vzEntry["match_dscp"] != nil {
						vzEntryAttr.MatchDscp = vzEntry["match_dscp"].(string)
					}
					if vzEntry["entry_name_alias"] != nil {
						vzEntryAttr.NameAlias = vzEntry["entry_name_alias"].(string)
					}
					if vzEntry["prot"] != nil {
						vzEntryAttr.Prot = vzEntry["prot"].(string)
					}
					if vzEntry["s_from_port"] != nil {
						vzEntryAttr.SFromPort = vzEntry["s_from_port"].(string)
					}
					if vzEntry["s_to_port"] != nil {
						vzEntryAttr.SToPort = vzEntry["s_to_port"].(string)
					}
					if vzEntry["stateful"] != nil {
						vzEntryAttr.Stateful = vzEntry["stateful"].(string)
					}
					if vzEntry["tcp_rules"] != nil {
						vzEntryAttr.TcpRules = vzEntry["tcp_rules"].(string)
					}

					vzFilterEntry := models.NewFilterEntry(fmt.Sprintf("e-%s", entryName), filterDn, entryDesc, vzEntryAttr)
					err := aciClient.Save(vzFilterEntry)
					if err != nil {
						return err
					}

					filterentryIDS = append(filterentryIDS, vzFilterEntry.DistinguishedName)
				}
			}

			filterIDS = append(filterIDS, vzFilter.DistinguishedName)
		}

		d.Set("filter_ids", filterIDS)
		d.Set("filter_entry_ids", filterentryIDS)
	}

	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if d.HasChange("relation_vz_rs_graph_att") {
		_, newRelParam := d.GetChange("relation_vz_rs_graph_att")
		err = aciClient.DeleteRelationvzRsGraphAttFromContract(vzBrCP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationvzRsGraphAttFromContract(vzBrCP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vz_rs_graph_att")
		d.Partial(false)

	}

	d.SetId(vzBrCP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciContractRead(d, m)

}

func resourceAciContractRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	vzBrCP, err := getRemoteContract(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setContractAttributes(vzBrCP, d)

	filters := d.Get("filter_ids").([]interface{})
	log.Println("Check ... :", filters)
	vzFilters := make([]*models.Filter, 0, 1)
	vzEntries := make([]*models.FilterEntry, 0, 1)
	for _, val := range filters {
		filterDN := val.(string)
		vzfilter, err := getRemoteFilterFromContract(aciClient, filterDN)
		if err == nil {
			for _, entry := range d.Get("filter_entry_ids").([]interface{}) {
				if strings.Contains(entry.(string), filterDN) {
					vzentry, err := getRemoteFilterEntryFromContract(aciClient, entry.(string))
					if err == nil {
						vzEntries = append(vzEntries, vzentry)
					}
				}
			}
			vzFilters = append(vzFilters, vzfilter)
		}
	}
	setFilterAttributesFromContract(vzFilters, vzEntries, d)

	vzRsGraphAttData, err := aciClient.ReadRelationvzRsGraphAttFromContract(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vzRsGraphAtt %v", err)

	} else {
		d.Set("relation_vz_rs_graph_att", vzRsGraphAttData)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciContractDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vzBrCP")
	if err != nil {
		return err
	}

	filters := d.Get("filter_ids").([]interface{})
	for _, val := range filters {
		filterDN := val.(string)
		err := aciClient.DeleteByDn(filterDN, "vzFilter")
		if err != nil {
			return err
		}
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
