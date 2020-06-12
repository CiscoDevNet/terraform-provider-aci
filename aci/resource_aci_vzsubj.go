package aci

import (
	"fmt"
	"log"
	"reflect"
	"sort"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciContractSubject() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciContractSubjectCreate,
		Update: resourceAciContractSubjectUpdate,
		Read:   resourceAciContractSubjectRead,
		Delete: resourceAciContractSubjectDelete,

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

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

func setContractSubjectAttributes(vzSubj *models.ContractSubject, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(vzSubj.DistinguishedName)
	d.Set("description", vzSubj.Description)
	// d.Set("contract_dn", GetParentDn(vzSubj.DistinguishedName))
	if dn != vzSubj.DistinguishedName {
		d.Set("contract_dn", "")
	}
	vzSubjMap, _ := vzSubj.ToMap()

	d.Set("name", vzSubjMap["name"])

	d.Set("annotation", vzSubjMap["annotation"])
	d.Set("cons_match_t", vzSubjMap["consMatchT"])
	d.Set("name_alias", vzSubjMap["nameAlias"])
	d.Set("prio", vzSubjMap["prio"])
	d.Set("prov_match_t", vzSubjMap["provMatchT"])
	d.Set("rev_flt_ports", vzSubjMap["revFltPorts"])
	d.Set("target_dscp", vzSubjMap["targetDscp"])
	return d
}

func resourceAciContractSubjectImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	vzSubj, err := getRemoteContractSubject(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setContractSubjectAttributes(vzSubj, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciContractSubjectCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] ContractSubject: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	ContractDn := d.Get("contract_dn").(string)

	vzSubjAttr := models.ContractSubjectAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vzSubjAttr.Annotation = Annotation.(string)
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
	vzSubj := models.NewContractSubject(fmt.Sprintf("subj-%s", name), ContractDn, desc, vzSubjAttr)

	err := aciClient.Save(vzSubj)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if relationTovzRsSubjGraphAtt, ok := d.GetOk("relation_vz_rs_subj_graph_att"); ok {
		relationParam := relationTovzRsSubjGraphAtt.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationvzRsSubjGraphAttFromContractSubject(vzSubj.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vz_rs_subj_graph_att")
		d.Partial(false)

	}
	if relationTovzRsSdwanPol, ok := d.GetOk("relation_vz_rs_sdwan_pol"); ok {
		relationParam := relationTovzRsSdwanPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationvzRsSdwanPolFromContractSubject(vzSubj.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vz_rs_sdwan_pol")
		d.Partial(false)

	}
	if relationTovzRsSubjFiltAtt, ok := d.GetOk("relation_vz_rs_subj_filt_att"); ok {
		relationParamList := toStringList(relationTovzRsSubjFiltAtt.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationvzRsSubjFiltAttFromContractSubject(vzSubj.DistinguishedName, relationParamName)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_vz_rs_subj_filt_att")
			d.Partial(false)
		}
	}

	d.SetId(vzSubj.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciContractSubjectRead(d, m)
}

func resourceAciContractSubjectUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] ContractSubject: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	ContractDn := d.Get("contract_dn").(string)

	vzSubjAttr := models.ContractSubjectAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vzSubjAttr.Annotation = Annotation.(string)
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
	vzSubj := models.NewContractSubject(fmt.Sprintf("subj-%s", name), ContractDn, desc, vzSubjAttr)

	vzSubj.Status = "modified"

	err := aciClient.Save(vzSubj)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if d.HasChange("relation_vz_rs_subj_graph_att") {
		_, newRelParam := d.GetChange("relation_vz_rs_subj_graph_att")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationvzRsSubjGraphAttFromContractSubject(vzSubj.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationvzRsSubjGraphAttFromContractSubject(vzSubj.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vz_rs_subj_graph_att")
		d.Partial(false)

	}
	if d.HasChange("relation_vz_rs_sdwan_pol") {
		_, newRelParam := d.GetChange("relation_vz_rs_sdwan_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationvzRsSdwanPolFromContractSubject(vzSubj.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationvzRsSdwanPolFromContractSubject(vzSubj.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vz_rs_sdwan_pol")
		d.Partial(false)

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
				return err
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationvzRsSubjFiltAttFromContractSubject(vzSubj.DistinguishedName, relDnName)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_vz_rs_subj_filt_att")
			d.Partial(false)

		}

	}

	d.SetId(vzSubj.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciContractSubjectRead(d, m)

}

func resourceAciContractSubjectRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	vzSubj, err := getRemoteContractSubject(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setContractSubjectAttributes(vzSubj, d)

	vzRsSubjGraphAttData, err := aciClient.ReadRelationvzRsSubjGraphAttFromContractSubject(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vzRsSubjGraphAtt %v", err)

	} else {
		if _, ok := d.GetOk("relation_vz_rs_subj_graph_att"); ok {
			tfName := GetMOName(d.Get("relation_vz_rs_subj_graph_att").(string))
			if tfName != vzRsSubjGraphAttData {
				d.Set("relation_vz_rs_subj_graph_att", "")
			}
		}
	}

	vzRsSdwanPolData, err := aciClient.ReadRelationvzRsSdwanPolFromContractSubject(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vzRsSdwanPol %v", err)

	} else {
		if _, ok := d.GetOk("relation_vz_rs_sdwan_pol"); ok {
			tfName := GetMOName(d.Get("relation_vz_rs_sdwan_pol").(string))
			if tfName != vzRsSdwanPolData {
				d.Set("relation_vz_rs_sdwan_pol", "")
			}
		}
	}

	vzRsSubjFiltAttData, err := aciClient.ReadRelationvzRsSubjFiltAttFromContractSubject(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vzRsSubjFiltAtt %v", err)

	} else {
		if _, ok := d.GetOk("relation_vz_rs_subj_filt_att"); ok {
			relationParamList := toStringList(d.Get("relation_vz_rs_subj_filt_att").(*schema.Set).List())
			tfList := make([]string, 0, 1)
			for _, relationParam := range relationParamList {
				relationParamName := GetMOName(relationParam)
				tfList = append(tfList, relationParamName)
			}
			vzRsSubjFiltAttDataList := toStringList(vzRsSubjFiltAttData.(*schema.Set).List())
			sort.Strings(tfList)
			sort.Strings(vzRsSubjFiltAttDataList)

			if !reflect.DeepEqual(tfList, vzRsSubjFiltAttDataList) {
				d.Set("relation_vz_rs_subj_filt_att", make([]string, 0, 1))
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciContractSubjectDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vzSubj")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
