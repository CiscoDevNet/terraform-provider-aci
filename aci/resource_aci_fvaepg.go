package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciApplicationEPG() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciApplicationEPGCreate,
		Update: resourceAciApplicationEPGUpdate,
		Read:   resourceAciApplicationEPGRead,
		Delete: resourceAciApplicationEPGDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciApplicationEPGImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"application_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"exception_tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"flood_on_encap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"fwd_ctrl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"has_mcast_source": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"is_attr_based_e_pg": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"match_t": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pc_enf_pref": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pref_gr_memb": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"prio": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"shutdown": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_fv_rs_bd": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_cust_qos_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_dom_att": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_fv_rs_fc_path_att": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_fv_rs_graph_def": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_fv_rs_cons_if": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_fv_rs_sec_inherited": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_fv_rs_node_att": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_fv_rs_dpp_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},

			"relation_fv_rs_prov_def": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_fv_rs_trust_ctrl": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_path_att": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_fv_rs_prot_by": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_fv_rs_ae_pg_mon_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_intra_epg": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
		}),
	}
}
func getRemoteApplicationEPG(client *client.Client, dn string) (*models.ApplicationEPG, error) {
	fvAEPgCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fvAEPg := models.ApplicationEPGFromContainer(fvAEPgCont)

	if fvAEPg.DistinguishedName == "" {
		return nil, fmt.Errorf("ApplicationEPG %s not found", fvAEPg.DistinguishedName)
	}

	return fvAEPg, nil
}

func setApplicationEPGAttributes(fvAEPg *models.ApplicationEPG, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(fvAEPg.DistinguishedName)
	d.Set("description", fvAEPg.Description)
	// d.Set("application_profile_dn", GetParentDn(fvAEPg.DistinguishedName))
	if dn != fvAEPg.DistinguishedName {
		d.Set("application_profile_dn", "")
	}
	fvAEPgMap, _ := fvAEPg.ToMap()

	d.Set("name", fvAEPgMap["name"])

	d.Set("annotation", fvAEPgMap["annotation"])
	d.Set("exception_tag", fvAEPgMap["exceptionTag"])
	d.Set("flood_on_encap", fvAEPgMap["floodOnEncap"])
	d.Set("fwd_ctrl", fvAEPgMap["fwdCtrl"])
	d.Set("has_mcast_source", fvAEPgMap["hasMcastSource"])
	d.Set("is_attr_based_e_pg", fvAEPgMap["isAttrBasedEPg"])
	d.Set("match_t", fvAEPgMap["matchT"])
	d.Set("name_alias", fvAEPgMap["nameAlias"])
	d.Set("pc_enf_pref", fvAEPgMap["pcEnfPref"])
	d.Set("pref_gr_memb", fvAEPgMap["prefGrMemb"])
	d.Set("prio", fvAEPgMap["prio"])
	d.Set("shutdown", fvAEPgMap["shutdown"])
	return d
}

func resourceAciApplicationEPGImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fvAEPg, err := getRemoteApplicationEPG(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setApplicationEPGAttributes(fvAEPg, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciApplicationEPGCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] ApplicationEPG: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	ApplicationProfileDn := d.Get("application_profile_dn").(string)

	fvAEPgAttr := models.ApplicationEPGAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvAEPgAttr.Annotation = Annotation.(string)
	}
	if ExceptionTag, ok := d.GetOk("exception_tag"); ok {
		fvAEPgAttr.ExceptionTag = ExceptionTag.(string)
	}
	if FloodOnEncap, ok := d.GetOk("flood_on_encap"); ok {
		fvAEPgAttr.FloodOnEncap = FloodOnEncap.(string)
	}
	if FwdCtrl, ok := d.GetOk("fwd_ctrl"); ok {
		fvAEPgAttr.FwdCtrl = FwdCtrl.(string)
	}
	if HasMcastSource, ok := d.GetOk("has_mcast_source"); ok {
		fvAEPgAttr.HasMcastSource = HasMcastSource.(string)
	}
	if IsAttrBasedEPg, ok := d.GetOk("is_attr_based_e_pg"); ok {
		fvAEPgAttr.IsAttrBasedEPg = IsAttrBasedEPg.(string)
	}
	if MatchT, ok := d.GetOk("match_t"); ok {
		fvAEPgAttr.MatchT = MatchT.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvAEPgAttr.NameAlias = NameAlias.(string)
	}
	if PcEnfPref, ok := d.GetOk("pc_enf_pref"); ok {
		fvAEPgAttr.PcEnfPref = PcEnfPref.(string)
	}
	if PrefGrMemb, ok := d.GetOk("pref_gr_memb"); ok {
		fvAEPgAttr.PrefGrMemb = PrefGrMemb.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		fvAEPgAttr.Prio = Prio.(string)
	}
	if Shutdown, ok := d.GetOk("shutdown"); ok {
		fvAEPgAttr.Shutdown = Shutdown.(string)
	}
	fvAEPg := models.NewApplicationEPG(fmt.Sprintf("epg-%s", name), ApplicationProfileDn, desc, fvAEPgAttr)

	err := aciClient.Save(fvAEPg)
	if err != nil {
		return err
	}

	if relationTofvRsBd, ok := d.GetOk("relation_fv_rs_bd"); ok {
		relationParam := relationTofvRsBd.(string)
		err = aciClient.CreateRelationfvRsBdFromApplicationEPG(fvAEPg.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}
	if relationTofvRsCustQosPol, ok := d.GetOk("relation_fv_rs_cust_qos_pol"); ok {
		relationParam := relationTofvRsCustQosPol.(string)
		err = aciClient.CreateRelationfvRsCustQosPolFromApplicationEPG(fvAEPg.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}
	if relationTofvRsDomAtt, ok := d.GetOk("relation_fv_rs_dom_att"); ok {
		relationParamList := toStringList(relationTofvRsDomAtt.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsDomAttFromApplicationEPG(fvAEPg.DistinguishedName, relationParam)

			if err != nil {
				return err
			}

		}
	}
	if relationTofvRsFcPathAtt, ok := d.GetOk("relation_fv_rs_fc_path_att"); ok {
		relationParamList := toStringList(relationTofvRsFcPathAtt.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsFcPathAttFromApplicationEPG(fvAEPg.DistinguishedName, relationParam)

			if err != nil {
				return err
			}

		}
	}
	if relationTofvRsGraphDef, ok := d.GetOk("relation_fv_rs_graph_def"); ok {
		relationParamList := toStringList(relationTofvRsGraphDef.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsGraphDefFromApplicationEPG(fvAEPg.DistinguishedName, relationParam)

			if err != nil {
				return err
			}

		}
	}
	if relationTofvRsConsIf, ok := d.GetOk("relation_fv_rs_cons_if"); ok {
		relationParamList := toStringList(relationTofvRsConsIf.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsConsIfFromApplicationEPG(fvAEPg.DistinguishedName, relationParam)

			if err != nil {
				return err
			}

		}
	}
	if relationTofvRsSecInherited, ok := d.GetOk("relation_fv_rs_sec_inherited"); ok {
		relationParamList := toStringList(relationTofvRsSecInherited.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsSecInheritedFromApplicationEPG(fvAEPg.DistinguishedName, relationParam)

			if err != nil {
				return err
			}

		}
	}
	if relationTofvRsNodeAtt, ok := d.GetOk("relation_fv_rs_node_att"); ok {
		relationParamList := toStringList(relationTofvRsNodeAtt.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsNodeAttFromApplicationEPG(fvAEPg.DistinguishedName, relationParam)

			if err != nil {
				return err
			}

		}
	}
	if relationTofvRsDppPol, ok := d.GetOk("relation_fv_rs_dpp_pol"); ok {
		relationParam := relationTofvRsDppPol.(string)
		err = aciClient.CreateRelationfvRsDppPolFromApplicationEPG(fvAEPg.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}
	if relationTofvRsProvDef, ok := d.GetOk("relation_fv_rs_prov_def"); ok {
		relationParamList := toStringList(relationTofvRsProvDef.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsProvDefFromApplicationEPG(fvAEPg.DistinguishedName, relationParam)

			if err != nil {
				return err
			}

		}
	}
	if relationTofvRsTrustCtrl, ok := d.GetOk("relation_fv_rs_trust_ctrl"); ok {
		relationParam := relationTofvRsTrustCtrl.(string)
		err = aciClient.CreateRelationfvRsTrustCtrlFromApplicationEPG(fvAEPg.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}
	if relationTofvRsPathAtt, ok := d.GetOk("relation_fv_rs_path_att"); ok {
		relationParamList := toStringList(relationTofvRsPathAtt.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsPathAttFromApplicationEPG(fvAEPg.DistinguishedName, relationParam)

			if err != nil {
				return err
			}

		}
	}
	if relationTofvRsProtBy, ok := d.GetOk("relation_fv_rs_prot_by"); ok {
		relationParamList := toStringList(relationTofvRsProtBy.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsProtByFromApplicationEPG(fvAEPg.DistinguishedName, relationParam)

			if err != nil {
				return err
			}

		}
	}
	if relationTofvRsAEPgMonPol, ok := d.GetOk("relation_fv_rs_ae_pg_mon_pol"); ok {
		relationParam := relationTofvRsAEPgMonPol.(string)
		err = aciClient.CreateRelationfvRsAEPgMonPolFromApplicationEPG(fvAEPg.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}
	if relationTofvRsIntraEpg, ok := d.GetOk("relation_fv_rs_intra_epg"); ok {
		relationParamList := toStringList(relationTofvRsIntraEpg.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsIntraEpgFromApplicationEPG(fvAEPg.DistinguishedName, relationParam)

			if err != nil {
				return err
			}

		}
	}

	d.SetId(fvAEPg.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciApplicationEPGRead(d, m)
}

func resourceAciApplicationEPGUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] ApplicationEPG: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	ApplicationProfileDn := d.Get("application_profile_dn").(string)

	fvAEPgAttr := models.ApplicationEPGAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvAEPgAttr.Annotation = Annotation.(string)
	}
	if ExceptionTag, ok := d.GetOk("exception_tag"); ok {
		fvAEPgAttr.ExceptionTag = ExceptionTag.(string)
	}
	if FloodOnEncap, ok := d.GetOk("flood_on_encap"); ok {
		fvAEPgAttr.FloodOnEncap = FloodOnEncap.(string)
	}
	if FwdCtrl, ok := d.GetOk("fwd_ctrl"); ok {
		fvAEPgAttr.FwdCtrl = FwdCtrl.(string)
	}
	if HasMcastSource, ok := d.GetOk("has_mcast_source"); ok {
		fvAEPgAttr.HasMcastSource = HasMcastSource.(string)
	}
	if IsAttrBasedEPg, ok := d.GetOk("is_attr_based_e_pg"); ok {
		fvAEPgAttr.IsAttrBasedEPg = IsAttrBasedEPg.(string)
	}
	if MatchT, ok := d.GetOk("match_t"); ok {
		fvAEPgAttr.MatchT = MatchT.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvAEPgAttr.NameAlias = NameAlias.(string)
	}
	if PcEnfPref, ok := d.GetOk("pc_enf_pref"); ok {
		fvAEPgAttr.PcEnfPref = PcEnfPref.(string)
	}
	if PrefGrMemb, ok := d.GetOk("pref_gr_memb"); ok {
		fvAEPgAttr.PrefGrMemb = PrefGrMemb.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		fvAEPgAttr.Prio = Prio.(string)
	}
	if Shutdown, ok := d.GetOk("shutdown"); ok {
		fvAEPgAttr.Shutdown = Shutdown.(string)
	}
	fvAEPg := models.NewApplicationEPG(fmt.Sprintf("epg-%s", name), ApplicationProfileDn, desc, fvAEPgAttr)

	fvAEPg.Status = "modified"

	err := aciClient.Save(fvAEPg)

	if err != nil {
		return err
	}

	if d.HasChange("relation_fv_rs_bd") {
		_, newRelParam := d.GetChange("relation_fv_rs_bd")
		err = aciClient.CreateRelationfvRsBdFromApplicationEPG(fvAEPg.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}
	if d.HasChange("relation_fv_rs_cust_qos_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_cust_qos_pol")
		err = aciClient.CreateRelationfvRsCustQosPolFromApplicationEPG(fvAEPg.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}
	if d.HasChange("relation_fv_rs_dom_att") {
		oldRel, newRel := d.GetChange("relation_fv_rs_dom_att")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsDomAttFromApplicationEPG(fvAEPg.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsDomAttFromApplicationEPG(fvAEPg.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

	}
	if d.HasChange("relation_fv_rs_fc_path_att") {
		oldRel, newRel := d.GetChange("relation_fv_rs_fc_path_att")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsFcPathAttFromApplicationEPG(fvAEPg.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsFcPathAttFromApplicationEPG(fvAEPg.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

	}
	if d.HasChange("relation_fv_rs_graph_def") {
		oldRel, newRel := d.GetChange("relation_fv_rs_graph_def")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsGraphDefFromApplicationEPG(fvAEPg.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

	}
	if d.HasChange("relation_fv_rs_cons_if") {
		oldRel, newRel := d.GetChange("relation_fv_rs_cons_if")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsConsIfFromApplicationEPG(fvAEPg.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsConsIfFromApplicationEPG(fvAEPg.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

	}
	if d.HasChange("relation_fv_rs_sec_inherited") {
		oldRel, newRel := d.GetChange("relation_fv_rs_sec_inherited")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsSecInheritedFromApplicationEPG(fvAEPg.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsSecInheritedFromApplicationEPG(fvAEPg.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

	}
	if d.HasChange("relation_fv_rs_node_att") {
		oldRel, newRel := d.GetChange("relation_fv_rs_node_att")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsNodeAttFromApplicationEPG(fvAEPg.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsNodeAttFromApplicationEPG(fvAEPg.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

	}
	if d.HasChange("relation_fv_rs_dpp_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_dpp_pol")
		err = aciClient.DeleteRelationfvRsDppPolFromApplicationEPG(fvAEPg.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationfvRsDppPolFromApplicationEPG(fvAEPg.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}
	if d.HasChange("relation_fv_rs_prov_def") {
		oldRel, newRel := d.GetChange("relation_fv_rs_prov_def")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsProvDefFromApplicationEPG(fvAEPg.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

	}
	if d.HasChange("relation_fv_rs_trust_ctrl") {
		_, newRelParam := d.GetChange("relation_fv_rs_trust_ctrl")
		err = aciClient.DeleteRelationfvRsTrustCtrlFromApplicationEPG(fvAEPg.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationfvRsTrustCtrlFromApplicationEPG(fvAEPg.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}
	if d.HasChange("relation_fv_rs_path_att") {
		oldRel, newRel := d.GetChange("relation_fv_rs_path_att")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsPathAttFromApplicationEPG(fvAEPg.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsPathAttFromApplicationEPG(fvAEPg.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

	}
	if d.HasChange("relation_fv_rs_prot_by") {
		oldRel, newRel := d.GetChange("relation_fv_rs_prot_by")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsProtByFromApplicationEPG(fvAEPg.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsProtByFromApplicationEPG(fvAEPg.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

	}
	if d.HasChange("relation_fv_rs_ae_pg_mon_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_ae_pg_mon_pol")
		err = aciClient.DeleteRelationfvRsAEPgMonPolFromApplicationEPG(fvAEPg.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationfvRsAEPgMonPolFromApplicationEPG(fvAEPg.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}
	if d.HasChange("relation_fv_rs_intra_epg") {
		oldRel, newRel := d.GetChange("relation_fv_rs_intra_epg")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsIntraEpgFromApplicationEPG(fvAEPg.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsIntraEpgFromApplicationEPG(fvAEPg.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

	}

	d.SetId(fvAEPg.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciApplicationEPGRead(d, m)

}

func resourceAciApplicationEPGRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	log.Printf("[DEBUG] %s: Start Get ", d.Id())
	fvAEPg, err := getRemoteApplicationEPG(aciClient, dn)
	log.Printf("[DEBUG] %s: Finished Get ", d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}
	log.Printf("[DEBUG] %s: Start set attributes ", d.Id())
	setApplicationEPGAttributes(fvAEPg, d)
	log.Printf("[DEBUG] %s: Finish set attributes ", d.Id())

	fvRsBdData, err := aciClient.ReadRelationfvRsBdFromApplicationEPG(dn)
	log.Printf("[DEBUG] %s: Finished reading relation fvRsBd", d.Id())
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsBd %v", err)

	} else {
		d.Set("relation_fv_rs_bd", fvRsBdData)
	}

	fvRsCustQosPolData, err := aciClient.ReadRelationfvRsCustQosPolFromApplicationEPG(dn)
	log.Printf("[DEBUG] %s: Finished reading relation fvRsCustQosPol", d.Id())

	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsCustQosPol %v", err)

	} else {
		d.Set("relation_fv_rs_cust_qos_pol", fvRsCustQosPolData)
	}

	fvRsDomAttData, err := aciClient.ReadRelationfvRsDomAttFromApplicationEPG(dn)
	log.Printf("[DEBUG] %s: Finished reading relation fvRsCustQosPol", d.Id())

	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsDomAtt %v", err)

	} else {
		d.Set("relation_fv_rs_dom_att", fvRsDomAttData)
	}

	fvRsFcPathAttData, err := aciClient.ReadRelationfvRsFcPathAttFromApplicationEPG(dn)
	log.Printf("[DEBUG] %s: Finished reading relation fvRsCustQosPol", d.Id())

	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsFcPathAtt %v", err)

	} else {
		d.Set("relation_fv_rs_fc_path_att", fvRsFcPathAttData)
	}

	fvRsGraphDefData, err := aciClient.ReadRelationfvRsGraphDefFromApplicationEPG(dn)
	log.Printf("[DEBUG] %s: Finished reading relation fvRsGraphDef", d.Id())

	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsGraphDef %v", err)

	} else {
		d.Set("relation_fv_rs_graph_def", fvRsGraphDefData)
	}

	fvRsConsIfData, err := aciClient.ReadRelationfvRsConsIfFromApplicationEPG(dn)
	log.Printf("[DEBUG] %s: Finished reading relation fvRsConsIf", d.Id())

	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsConsIf %v", err)

	} else {
		d.Set("relation_fv_rs_cons_if", fvRsConsIfData)
	}

	fvRsSecInheritedData, err := aciClient.ReadRelationfvRsSecInheritedFromApplicationEPG(dn)
	log.Printf("[DEBUG] %s: Finished reading relation fvRsSecInherited", d.Id())

	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsSecInherited %v", err)

	} else {
		d.Set("relation_fv_rs_sec_inherited", fvRsSecInheritedData)
	}

	fvRsNodeAttData, err := aciClient.ReadRelationfvRsNodeAttFromApplicationEPG(dn)
	log.Printf("[DEBUG] %s: Finished reading relation fvRsNodeAtt", d.Id())

	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsNodeAtt %v", err)

	} else {
		d.Set("relation_fv_rs_node_att", fvRsNodeAttData)
	}

	fvRsDppPolData, err := aciClient.ReadRelationfvRsDppPolFromApplicationEPG(dn)
	log.Printf("[DEBUG] %s: Finished reading relation fvRsDppPol", d.Id())

	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsDppPol %v", err)

	} else {
		d.Set("relation_fv_rs_dpp_pol", fvRsDppPolData)
	}

	fvRsProvDefData, err := aciClient.ReadRelationfvRsProvDefFromApplicationEPG(dn)
	log.Printf("[DEBUG] %s: Finished reading relation fvRsProvDef", d.Id())

	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsProvDef %v", err)

	} else {
		d.Set("relation_fv_rs_prov_def", fvRsProvDefData)
	}

	fvRsTrustCtrlData, err := aciClient.ReadRelationfvRsTrustCtrlFromApplicationEPG(dn)
	log.Printf("[DEBUG] %s: Finished reading relation fvRsTrustCtrl", d.Id())

	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsTrustCtrl %v", err)

	} else {
		d.Set("relation_fv_rs_trust_ctrl", fvRsTrustCtrlData)
	}

	fvRsPathAttData, err := aciClient.ReadRelationfvRsPathAttFromApplicationEPG(dn)
	log.Printf("[DEBUG] %s: Finished reading relation fvRsPathAtt", d.Id())

	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsPathAtt %v", err)

	} else {
		d.Set("relation_fv_rs_path_att", fvRsPathAttData)
	}

	fvRsProtByData, err := aciClient.ReadRelationfvRsProtByFromApplicationEPG(dn)
	log.Printf("[DEBUG] %s: Finished reading relation fvRsProtBy", d.Id())

	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsProtBy %v", err)

	} else {
		d.Set("relation_fv_rs_prot_by", fvRsProtByData)
	}

	fvRsAEPgMonPolData, err := aciClient.ReadRelationfvRsAEPgMonPolFromApplicationEPG(dn)
	log.Printf("[DEBUG] %s: Finished reading relation fvRsAEPgMonPol", d.Id())

	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsAEPgMonPol %v", err)

	} else {
		d.Set("relation_fv_rs_ae_pg_mon_pol", fvRsAEPgMonPolData)
	}

	fvRsIntraEpgData, err := aciClient.ReadRelationfvRsIntraEpgFromApplicationEPG(dn)
	log.Printf("[DEBUG] %s: Finished reading relation fvRsIntraEpg", d.Id())

	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsIntraEpg %v", err)

	} else {
		d.Set("relation_fv_rs_intra_epg", fvRsIntraEpgData)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciApplicationEPGDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fvAEPg")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
