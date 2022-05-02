package aci

import (
	"context"
	"fmt"
	"log"

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
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_vz_rs_sdwan_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_vz_rs_subj_filt_att": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"consumer_to_provider": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if (k == "consumer_to_provider.target_dscp" || k == "consumer_to_provider.prio" || k == "consumer_to_provider.relation_vz_rs_in_term_graph_att") && new != old {
						return false
					}
					return true
				},
				Description: "Set InTerm attributes",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"provider_to_consumer": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if (k == "provider_to_consumer.target_dscp" || k == "provider_to_consumer.prio" || k == "provider_to_consumer.relation_vz_rs_out_term_graph_att") && new != old {
						return false
					}
					return true
				},
				Description: "Set OutTerm attributes",
				Elem: &schema.Schema{
					Type: schema.TypeString,
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

	vzInTermCont, err := client.Get(dn + "/intmnl")
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

func setInTermSubjectAttributes(vzInTerm *models.InTermSubject, d map[string]string) (map[string]string, error) {

	vzInTermMap, err := vzInTerm.ToMap()
	if err != nil {
		return d, err
	}
	d = map[string]string{
		"id":          vzInTerm.DistinguishedName,
		"prio":        vzInTermMap["prio"],
		"target_dscp": vzInTermMap["targetDscp"],
	}
	return d, nil
}

func getRemoteOutTermSubject(client *client.Client, dn string) (*models.OutTermSubject, error) {
	vzOutTermCont, err := client.Get(dn + "/outtmnl")
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

func setOutTermSubjectAttributes(vzOutTerm *models.OutTermSubject, d map[string]string) (map[string]string, error) {

	vzOutTermMap, err := vzOutTerm.ToMap()
	if err != nil {
		return d, err
	}
	d = map[string]string{
		"id":          vzOutTerm.DistinguishedName,
		"prio":        vzOutTermMap["prio"],
		"target_dscp": vzOutTermMap["targetDscp"],
	}
	return d, nil
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

	if ApplyBothDirections == "yes" {
		if len(d.Get("consumer_to_provider").(map[string]interface{})) != 0 || len(d.Get("provider_to_consumer").(map[string]interface{})) != 0 {
			d.Set("consumer_to_provider", nil)
			d.Set("provider_to_consumer", nil)
			return diag.FromErr(fmt.Errorf("you cannot set consumer_to_provider and provider_to_consumer when apply_both_directions is set to yes"))
		}
	}

	vzSubj := models.NewContractSubject(fmt.Sprintf("subj-%s", name), ContractDn, desc, vzSubjAttr)

	err := aciClient.Save(vzSubj)
	if err != nil {
		return diag.FromErr(err)
	}

	if ApplyBothDirections == "no" {
		var inTermGraphAttribute string
		vzInTermAttr := models.InTermSubjectAttributes{}
		if ConsumerToProvider, ok := d.GetOk("consumer_to_provider"); ok {
			ConsumerToProviderMap := ConsumerToProvider.(map[string]interface{})
			vzInTermAttr.Prio = ConsumerToProviderMap["prio"].(string)
			vzInTermAttr.TargetDscp = ConsumerToProviderMap["target_dscp"].(string)
			inTermGraphAttribute = ConsumerToProviderMap["relation_vz_rs_in_term_graph_att"].(string)

		}
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

		vzOutTermAttr := models.OutTermSubjectAttributes{}
		var outTermGraphAttribute string
		if ProviderToConsumer, ok := d.GetOk("provider_to_consumer"); ok {
			ProviderToConsumerMap := ProviderToConsumer.(map[string]interface{})
			vzOutTermAttr.Prio = ProviderToConsumerMap["prio"].(string)
			vzOutTermAttr.TargetDscp = ProviderToConsumerMap["target_dscp"].(string)
			outTermGraphAttribute = ProviderToConsumerMap["relation_vz_rs_out_term_graph_att"].(string)
		}
		vzOutTerm := models.NewOutTermSubject(fmt.Sprintf(models.RnvzOutTerm), vzSubj.DistinguishedName, desc, "", vzOutTermAttr)

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

			err = aciClient.CreateRelationvzRsOutTermGraphAtt(vzOutTerm.DistinguishedName, vzOutTermAttr.Annotation, GetMOName(outTermGraphAttribute))
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

	if ApplyBothDirections == "yes" {
		if (d.HasChange("consumer_to_provider") || d.HasChange("provider_to_consumer")) && (d.Get("consumer_to_provider") != nil || d.Get("provider_to_consumer") != nil) {
			d.Set("consumer_to_provider", nil)
			d.Set("provider_to_consumer", nil)
			return diag.FromErr(fmt.Errorf("you cannot set consumer_to_provider and provider_to_consumer when apply_both_directions is set to yes"))
		}
	}

	vzSubj := models.NewContractSubject(fmt.Sprintf("subj-%s", name), ContractDn, desc, vzSubjAttr)

	vzSubj.Status = "modified"

	err := aciClient.Save(vzSubj)
	if err != nil {
		return diag.FromErr(err)
	}

	if ApplyBothDirections == "no" && (d.HasChange("consumer_to_provider") || d.HasChange("provider_to_consumer")) {
		if d.HasChange("consumer_to_provider") {
			vzInTermAttr := models.InTermSubjectAttributes{}
			var inTermGraphAttribute string
			if ConsumerToProvider, ok := d.GetOk("consumer_to_provider"); ok {
				ConsumerToProviderMap := ConsumerToProvider.(map[string]interface{})
				vzInTermAttr.Prio = ConsumerToProviderMap["prio"].(string)
				vzInTermAttr.TargetDscp = ConsumerToProviderMap["target_dscp"].(string)
				inTermGraphAttribute = ConsumerToProviderMap["relation_vz_rs_in_term_graph_att"].(string)
			}
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
		}
		if d.HasChange("provider_to_consumer") {
			vzOutTermAttr := models.OutTermSubjectAttributes{}
			var outTermGraphAttribute string
			if ProviderToConsumer, ok := d.GetOk("provider_to_consumer"); ok {
				ProviderToConsumerMap := ProviderToConsumer.(map[string]interface{})
				vzOutTermAttr.Prio = ProviderToConsumerMap["prio"].(string)
				vzOutTermAttr.TargetDscp = ProviderToConsumerMap["target_dscp"].(string)
				outTermGraphAttribute = ProviderToConsumerMap["relation_vz_rs_out_term_graph_att"].(string)
			}
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
		vzInTermFactor, err := setInTermSubjectAttributes(vzInTerm, make(map[string]string))
		if err != nil {
			d.SetId("")
			return nil
		}

		vzRsInTermGraphAttData, err := aciClient.ReadRelationvzRsInTermGraphAtt(vzInTerm.DistinguishedName)
		if err != nil {
			log.Printf("[DEBUG] Error while reading relation vzRsInTermGraphAtt %v", err)
			d.Set("relation_vz_rs_in_term_graph_att", "")
		} else {
			vzInTermFactor["relation_vz_rs_in_term_graph_att"] = vzRsInTermGraphAttData.(string)
		}

		d.Set("consumer_to_provider", vzInTermFactor)

	} else {
		d.Set("consumer_to_provider", nil)
	}

	vzOutTerm, err := getRemoteOutTermSubject(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}
	if vzOutTerm != nil {
		vzOutTermFactor, err := setOutTermSubjectAttributes(vzOutTerm, make(map[string]string))
		if err != nil {
			d.SetId("")
			return nil
		}

		vzRsOutTermGraphAttData, err := aciClient.ReadRelationvzRsOutTermGraphAtt(vzOutTerm.DistinguishedName)
		if err != nil {
			log.Printf("[DEBUG] Error while reading relation vzRsOutTermGraphAtt %v", err)
			d.Set("relation_vz_rs_out_term_graph_att", "")
		} else {
			vzOutTermFactor["relation_vz_rs_out_term_graph_att"] = vzRsOutTermGraphAttData.(string)
		}

		d.Set("provider_to_consumer", vzOutTermFactor)

	} else {
		d.Set("consumer_to_provider", nil)
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
