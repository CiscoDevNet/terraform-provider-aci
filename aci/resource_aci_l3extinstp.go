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

func resourceAciExternalNetworkInstanceProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciExternalNetworkInstanceProfileCreate,
		Update: resourceAciExternalNetworkInstanceProfileUpdate,
		Read:   resourceAciExternalNetworkInstanceProfileRead,
		Delete: resourceAciExternalNetworkInstanceProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciExternalNetworkInstanceProfileImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"l3_outside_dn": &schema.Schema{
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

			"target_dscp": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_fv_rs_sec_inherited": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_fv_rs_prov": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_l3ext_rs_l3_inst_p_to_dom_p": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_l3ext_rs_inst_p_to_nat_mapping_e_pg": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_cons_if": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_fv_rs_cust_qos_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_l3ext_rs_inst_p_to_profile": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tn_rtctrl_profile_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"direction": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"relation_fv_rs_cons": &schema.Schema{
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
			"relation_fv_rs_intra_epg": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
		}),
	}
}
func getRemoteExternalNetworkInstanceProfile(client *client.Client, dn string) (*models.ExternalNetworkInstanceProfile, error) {
	l3extInstPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extInstP := models.ExternalNetworkInstanceProfileFromContainer(l3extInstPCont)

	if l3extInstP.DistinguishedName == "" {
		return nil, fmt.Errorf("ExternalNetworkInstanceProfile %s not found", l3extInstP.DistinguishedName)
	}

	return l3extInstP, nil
}

func setExternalNetworkInstanceProfileAttributes(l3extInstP *models.ExternalNetworkInstanceProfile, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(l3extInstP.DistinguishedName)
	d.Set("description", l3extInstP.Description)
	// d.Set("l3_outside_dn", GetParentDn(l3extInstP.DistinguishedName))
	if dn != l3extInstP.DistinguishedName {
		d.Set("l3_outside_dn", "")
	}
	l3extInstPMap, _ := l3extInstP.ToMap()

	d.Set("name", l3extInstPMap["name"])

	d.Set("annotation", l3extInstPMap["annotation"])
	d.Set("exception_tag", l3extInstPMap["exceptionTag"])
	d.Set("flood_on_encap", l3extInstPMap["floodOnEncap"])
	d.Set("match_t", l3extInstPMap["matchT"])
	d.Set("name_alias", l3extInstPMap["nameAlias"])
	d.Set("pref_gr_memb", l3extInstPMap["prefGrMemb"])
	d.Set("prio", l3extInstPMap["prio"])
	d.Set("target_dscp", l3extInstPMap["targetDscp"])
	return d
}

func resourceAciExternalNetworkInstanceProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	l3extInstP, err := getRemoteExternalNetworkInstanceProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setExternalNetworkInstanceProfileAttributes(l3extInstP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciExternalNetworkInstanceProfileCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] ExternalNetworkInstanceProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	L3OutsideDn := d.Get("l3_outside_dn").(string)

	l3extInstPAttr := models.ExternalNetworkInstanceProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extInstPAttr.Annotation = Annotation.(string)
	}
	if ExceptionTag, ok := d.GetOk("exception_tag"); ok {
		l3extInstPAttr.ExceptionTag = ExceptionTag.(string)
	}
	if FloodOnEncap, ok := d.GetOk("flood_on_encap"); ok {
		l3extInstPAttr.FloodOnEncap = FloodOnEncap.(string)
	}
	if MatchT, ok := d.GetOk("match_t"); ok {
		l3extInstPAttr.MatchT = MatchT.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extInstPAttr.NameAlias = NameAlias.(string)
	}
	if PrefGrMemb, ok := d.GetOk("pref_gr_memb"); ok {
		l3extInstPAttr.PrefGrMemb = PrefGrMemb.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		l3extInstPAttr.Prio = Prio.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		l3extInstPAttr.TargetDscp = TargetDscp.(string)
	}
	l3extInstP := models.NewExternalNetworkInstanceProfile(fmt.Sprintf("instP-%s", name), L3OutsideDn, desc, l3extInstPAttr)

	err := aciClient.Save(l3extInstP)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if relationTofvRsSecInherited, ok := d.GetOk("relation_fv_rs_sec_inherited"); ok {
		relationParamList := toStringList(relationTofvRsSecInherited.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsSecInheritedFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_sec_inherited")
			d.Partial(false)
		}
	}
	if relationTofvRsProv, ok := d.GetOk("relation_fv_rs_prov"); ok {
		relationParamList := toStringList(relationTofvRsProv.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsProvFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relationParamName)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_prov")
			d.Partial(false)
		}
	}
	if relationTol3extRsL3InstPToDomP, ok := d.GetOk("relation_l3ext_rs_l3_inst_p_to_dom_p"); ok {
		relationParam := relationTol3extRsL3InstPToDomP.(string)
		err = aciClient.CreateRelationl3extRsL3InstPToDomPFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_l3ext_rs_l3_inst_p_to_dom_p")
		d.Partial(false)

	}
	if relationTol3extRsInstPToNatMappingEPg, ok := d.GetOk("relation_l3ext_rs_inst_p_to_nat_mapping_e_pg"); ok {
		relationParam := relationTol3extRsInstPToNatMappingEPg.(string)
		err = aciClient.CreateRelationl3extRsInstPToNatMappingEPgFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_l3ext_rs_inst_p_to_nat_mapping_e_pg")
		d.Partial(false)

	}
	if relationTofvRsConsIf, ok := d.GetOk("relation_fv_rs_cons_if"); ok {
		relationParamList := toStringList(relationTofvRsConsIf.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsConsIfFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relationParamName)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_cons_if")
			d.Partial(false)
		}
	}
	if relationTofvRsCustQosPol, ok := d.GetOk("relation_fv_rs_cust_qos_pol"); ok {
		relationParam := relationTofvRsCustQosPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsCustQosPolFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_cust_qos_pol")
		d.Partial(false)

	}
	if relationTol3extRsInstPToProfile, ok := d.GetOk("relation_l3ext_rs_inst_p_to_profile"); ok {

		relationParamList := relationTol3extRsInstPToProfile.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationl3extRsInstPToProfileFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, paramMap["tn_rtctrl_profile_name"].(string), paramMap["direction"].(string))
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_l3ext_rs_inst_p_to_profile")
			d.Partial(false)
		}

	}
	if relationTofvRsCons, ok := d.GetOk("relation_fv_rs_cons"); ok {
		relationParamList := toStringList(relationTofvRsCons.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsConsFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relationParamName)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_cons")
			d.Partial(false)
		}
	}
	if relationTofvRsProtBy, ok := d.GetOk("relation_fv_rs_prot_by"); ok {
		relationParamList := toStringList(relationTofvRsProtBy.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsProtByFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relationParamName)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_prot_by")
			d.Partial(false)
		}
	}
	if relationTofvRsIntraEpg, ok := d.GetOk("relation_fv_rs_intra_epg"); ok {
		relationParamList := toStringList(relationTofvRsIntraEpg.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsIntraEpgFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relationParamName)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_intra_epg")
			d.Partial(false)
		}
	}

	d.SetId(l3extInstP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciExternalNetworkInstanceProfileRead(d, m)
}

func resourceAciExternalNetworkInstanceProfileUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] ExternalNetworkInstanceProfile: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	L3OutsideDn := d.Get("l3_outside_dn").(string)

	l3extInstPAttr := models.ExternalNetworkInstanceProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extInstPAttr.Annotation = Annotation.(string)
	}
	if ExceptionTag, ok := d.GetOk("exception_tag"); ok {
		l3extInstPAttr.ExceptionTag = ExceptionTag.(string)
	}
	if FloodOnEncap, ok := d.GetOk("flood_on_encap"); ok {
		l3extInstPAttr.FloodOnEncap = FloodOnEncap.(string)
	}
	if MatchT, ok := d.GetOk("match_t"); ok {
		l3extInstPAttr.MatchT = MatchT.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extInstPAttr.NameAlias = NameAlias.(string)
	}
	if PrefGrMemb, ok := d.GetOk("pref_gr_memb"); ok {
		l3extInstPAttr.PrefGrMemb = PrefGrMemb.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		l3extInstPAttr.Prio = Prio.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		l3extInstPAttr.TargetDscp = TargetDscp.(string)
	}
	l3extInstP := models.NewExternalNetworkInstanceProfile(fmt.Sprintf("instP-%s", name), L3OutsideDn, desc, l3extInstPAttr)

	l3extInstP.Status = "modified"

	err := aciClient.Save(l3extInstP)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if d.HasChange("relation_fv_rs_sec_inherited") {
		oldRel, newRel := d.GetChange("relation_fv_rs_sec_inherited")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsSecInheritedFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsSecInheritedFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDn)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_sec_inherited")
			d.Partial(false)

		}

	}
	if d.HasChange("relation_fv_rs_prov") {
		oldRel, newRel := d.GetChange("relation_fv_rs_prov")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			relDnName := GetMOName(relDn)
			err = aciClient.DeleteRelationfvRsProvFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDnName)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationfvRsProvFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDnName)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_prov")
			d.Partial(false)

		}

	}
	if d.HasChange("relation_l3ext_rs_l3_inst_p_to_dom_p") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_l3_inst_p_to_dom_p")
		err = aciClient.CreateRelationl3extRsL3InstPToDomPFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_l3ext_rs_l3_inst_p_to_dom_p")
		d.Partial(false)

	}
	if d.HasChange("relation_l3ext_rs_inst_p_to_nat_mapping_e_pg") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_inst_p_to_nat_mapping_e_pg")
		err = aciClient.DeleteRelationl3extRsInstPToNatMappingEPgFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationl3extRsInstPToNatMappingEPgFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_l3ext_rs_inst_p_to_nat_mapping_e_pg")
		d.Partial(false)

	}
	if d.HasChange("relation_fv_rs_cons_if") {
		oldRel, newRel := d.GetChange("relation_fv_rs_cons_if")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			relDnName := GetMOName(relDn)
			err = aciClient.DeleteRelationfvRsConsIfFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDnName)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationfvRsConsIfFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDnName)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_cons_if")
			d.Partial(false)

		}

	}
	if d.HasChange("relation_fv_rs_cust_qos_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_cust_qos_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationfvRsCustQosPolFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_cust_qos_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_l3ext_rs_inst_p_to_profile") {
		oldRel, newRel := d.GetChange("relation_l3ext_rs_inst_p_to_profile")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.DeleteRelationl3extRsInstPToProfileFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, paramMap["tn_rtctrl_profile_name"].(string), paramMap["direction"].(string))
			if err != nil {
				return err
			}

		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationl3extRsInstPToProfileFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, paramMap["tn_rtctrl_profile_name"].(string), paramMap["direction"].(string))
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_l3ext_rs_inst_p_to_profile")
			d.Partial(false)
		}

	}
	if d.HasChange("relation_fv_rs_cons") {
		oldRel, newRel := d.GetChange("relation_fv_rs_cons")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			relDnName := GetMOName(relDn)
			err = aciClient.DeleteRelationfvRsConsFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDnName)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationfvRsConsFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDnName)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_cons")
			d.Partial(false)

		}

	}
	if d.HasChange("relation_fv_rs_prot_by") {
		oldRel, newRel := d.GetChange("relation_fv_rs_prot_by")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			relDnName := GetMOName(relDn)
			err = aciClient.DeleteRelationfvRsProtByFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDnName)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationfvRsProtByFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDnName)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_prot_by")
			d.Partial(false)

		}

	}
	if d.HasChange("relation_fv_rs_intra_epg") {
		oldRel, newRel := d.GetChange("relation_fv_rs_intra_epg")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			relDnName := GetMOName(relDn)
			err = aciClient.DeleteRelationfvRsIntraEpgFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDnName)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationfvRsIntraEpgFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDnName)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_intra_epg")
			d.Partial(false)

		}

	}

	d.SetId(l3extInstP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciExternalNetworkInstanceProfileRead(d, m)

}

func resourceAciExternalNetworkInstanceProfileRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l3extInstP, err := getRemoteExternalNetworkInstanceProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setExternalNetworkInstanceProfileAttributes(l3extInstP, d)

	fvRsSecInheritedData, err := aciClient.ReadRelationfvRsSecInheritedFromExternalNetworkInstanceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsSecInherited %v", err)

	} else {
		d.Set("relation_fv_rs_sec_inherited", fvRsSecInheritedData)
	}

	fvRsProvData, err := aciClient.ReadRelationfvRsProvFromExternalNetworkInstanceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsProv %v", err)

	} else {
		if _, ok := d.GetOk("relation_fv_rs_prov"); ok {
			relationParamList := toStringList(d.Get("relation_fv_rs_prov").(*schema.Set).List())
			tfList := make([]string, 0, 1)
			for _, relationParam := range relationParamList {
				relationParamName := GetMOName(relationParam)
				tfList = append(tfList, relationParamName)
			}
			fvRsProvDataList := toStringList(fvRsProvData.(*schema.Set).List())
			sort.Strings(tfList)
			sort.Strings(fvRsProvDataList)

			if !reflect.DeepEqual(tfList, fvRsProvDataList) {
				d.Set("relation_fv_rs_prov", make([]string, 0, 1))
			}
		}
	}

	l3extRsL3InstPToDomPData, err := aciClient.ReadRelationl3extRsL3InstPToDomPFromExternalNetworkInstanceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsL3InstPToDomP %v", err)

	} else {
		if _, ok := d.GetOk("relation_l3ext_rs_l3_inst_p_to_dom_p"); ok {
			tfName := d.Get("relation_l3ext_rs_l3_inst_p_to_dom_p").(string)
			if tfName != l3extRsL3InstPToDomPData {
				d.Set("relation_fv_rs_nd_pfx_pol", "")
			}
		}
	}

	l3extRsInstPToNatMappingEPgData, err := aciClient.ReadRelationl3extRsInstPToNatMappingEPgFromExternalNetworkInstanceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsInstPToNatMappingEPg %v", err)

	} else {
		if _, ok := d.GetOk("relation_l3ext_rs_inst_p_to_nat_mapping_e_pg"); ok {
			tfName := d.Get("relation_l3ext_rs_inst_p_to_nat_mapping_e_pg").(string)
			if tfName != l3extRsInstPToNatMappingEPgData {
				d.Set("relation_l3ext_rs_inst_p_to_nat_mapping_e_pg", "")
			}
		}
	}

	fvRsConsIfData, err := aciClient.ReadRelationfvRsConsIfFromExternalNetworkInstanceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsConsIf %v", err)

	} else {
		if _, ok := d.GetOk("relation_fv_rs_cons_if"); ok {
			relationParamList := toStringList(d.Get("relation_fv_rs_cons_if").(*schema.Set).List())
			tfList := make([]string, 0, 1)
			for _, relationParam := range relationParamList {
				relationParamName := GetMOName(relationParam)
				tfList = append(tfList, relationParamName)
			}
			fvRsConsIfDataList := toStringList(fvRsConsIfData.(*schema.Set).List())
			sort.Strings(tfList)
			sort.Strings(fvRsConsIfDataList)

			if !reflect.DeepEqual(tfList, fvRsConsIfDataList) {
				d.Set("relation_fv_rs_cons_if", make([]string, 0, 1))
			}
		}
	}

	fvRsCustQosPolData, err := aciClient.ReadRelationfvRsCustQosPolFromExternalNetworkInstanceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsCustQosPol %v", err)

	} else {
		if _, ok := d.GetOk("relation_fv_rs_cust_qos_pol"); ok {
			tfName := GetMOName(d.Get("relation_fv_rs_cust_qos_pol").(string))
			if tfName != fvRsCustQosPolData {
				d.Set("relation_fv_rs_cust_qos_pol", "")
			}
		}
	}

	l3extRsInstPToProfileData, err := aciClient.ReadRelationl3extRsInstPToProfileFromExternalNetworkInstanceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsInstPToProfile %v", err)

	} else {
		d.Set("relation_l3ext_rs_inst_p_to_profile", l3extRsInstPToProfileData)
	}

	fvRsConsData, err := aciClient.ReadRelationfvRsConsFromExternalNetworkInstanceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsCons %v", err)

	} else {
		if _, ok := d.GetOk("relation_fv_rs_cons"); ok {
			relationParamList := toStringList(d.Get("relation_fv_rs_cons").(*schema.Set).List())
			tfList := make([]string, 0, 1)
			for _, relationParam := range relationParamList {
				relationParamName := GetMOName(relationParam)
				tfList = append(tfList, relationParamName)
			}
			fvRsConsDataList := toStringList(fvRsConsData.(*schema.Set).List())
			sort.Strings(tfList)
			sort.Strings(fvRsConsDataList)

			if !reflect.DeepEqual(tfList, fvRsConsDataList) {
				d.Set("relation_fv_rs_cons", make([]string, 0, 1))
			}
		}
	}

	fvRsProtByData, err := aciClient.ReadRelationfvRsProtByFromExternalNetworkInstanceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsProtBy %v", err)

	} else {
		if _, ok := d.GetOk("relation_fv_rs_prot_by"); ok {
			relationParamList := toStringList(d.Get("relation_fv_rs_prot_by").(*schema.Set).List())
			tfList := make([]string, 0, 1)
			for _, relationParam := range relationParamList {
				relationParamName := GetMOName(relationParam)
				tfList = append(tfList, relationParamName)
			}
			fvRsProtByDataList := toStringList(fvRsProtByData.(*schema.Set).List())
			sort.Strings(tfList)
			sort.Strings(fvRsProtByDataList)

			if !reflect.DeepEqual(tfList, fvRsProtByDataList) {
				d.Set("relation_fv_rs_prot_by", make([]string, 0, 1))
			}
		}
	}

	fvRsIntraEpgData, err := aciClient.ReadRelationfvRsIntraEpgFromExternalNetworkInstanceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsIntraEpg %v", err)

	} else {
		if _, ok := d.GetOk("relation_fv_rs_intra_epg"); ok {
			relationParamList := toStringList(d.Get("relation_fv_rs_intra_epg").(*schema.Set).List())
			tfList := make([]string, 0, 1)
			for _, relationParam := range relationParamList {
				relationParamName := GetMOName(relationParam)
				tfList = append(tfList, relationParamName)
			}
			fvRsIntraEpgDataList := toStringList(fvRsIntraEpgData.(*schema.Set).List())
			sort.Strings(tfList)
			sort.Strings(fvRsIntraEpgDataList)

			if !reflect.DeepEqual(tfList, fvRsIntraEpgDataList) {
				d.Set("relation_fv_rs_intra_epg", make([]string, 0, 1))
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciExternalNetworkInstanceProfileDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l3extInstP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
