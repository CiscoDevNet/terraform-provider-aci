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

func resourceAciExternalNetworkInstanceProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciExternalNetworkInstanceProfileCreate,
		UpdateContext: resourceAciExternalNetworkInstanceProfileUpdate,
		ReadContext:   resourceAciExternalNetworkInstanceProfileRead,
		DeleteContext: resourceAciExternalNetworkInstanceProfileDelete,

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

			"exception_tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"flood_on_encap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
				}, false),
			},

			"match_t": &schema.Schema{
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

			"pref_gr_memb": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"include",
					"exclude",
				}, false),
			},

			"prio": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"unspecified",
					"level6",
					"level5",
					"level4",
					"level3",
					"level2",
					"level1",
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
			"relation_l3ext_rs_inst_p_to_nat_mapping_epg": &schema.Schema{
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
				Type:     schema.TypeString,
				Computed: true,
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

func setExternalNetworkInstanceProfileAttributes(l3extInstP *models.ExternalNetworkInstanceProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(l3extInstP.DistinguishedName)
	d.Set("description", l3extInstP.Description)

	if dn != l3extInstP.DistinguishedName {
		d.Set("l3_outside_dn", "")
	}
	l3extInstPMap, err := l3extInstP.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("l3_outside_dn", GetParentDn(dn, fmt.Sprintf("/instP-%s", l3extInstPMap["name"])))
	d.Set("name", l3extInstPMap["name"])

	d.Set("annotation", l3extInstPMap["annotation"])
	d.Set("exception_tag", l3extInstPMap["exceptionTag"])
	d.Set("flood_on_encap", l3extInstPMap["floodOnEncap"])
	d.Set("match_t", l3extInstPMap["matchT"])
	d.Set("name_alias", l3extInstPMap["nameAlias"])
	d.Set("pref_gr_memb", l3extInstPMap["prefGrMemb"])
	d.Set("prio", l3extInstPMap["prio"])
	d.Set("target_dscp", l3extInstPMap["targetDscp"])
	return d, nil
}

func resourceAciExternalNetworkInstanceProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	l3extInstP, err := getRemoteExternalNetworkInstanceProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	l3extInstPMap, err := l3extInstP.ToMap()
	if err != nil {
		return nil, err
	}
	name := l3extInstPMap["name"]
	pDN := GetParentDn(dn, fmt.Sprintf("/instP-%s", name))
	d.Set("l3_outside_dn", pDN)
	schemaFilled, err := setExternalNetworkInstanceProfileAttributes(l3extInstP, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciExternalNetworkInstanceProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ExternalNetworkInstanceProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	L3OutsideDn := d.Get("l3_outside_dn").(string)

	l3extInstPAttr := models.ExternalNetworkInstanceProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extInstPAttr.Annotation = Annotation.(string)
	} else {
		l3extInstPAttr.Annotation = "{}"
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
		return diag.FromErr(err)
	}
	checkDns := make([]string, 0, 1)

	if relationTofvRsSecInherited, ok := d.GetOk("relation_fv_rs_sec_inherited"); ok {
		relationParamList := toStringList(relationTofvRsSecInherited.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTofvRsProv, ok := d.GetOk("relation_fv_rs_prov"); ok {
		relationParamList := toStringList(relationTofvRsProv.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTol3extRsL3InstPToDomP, ok := d.GetOk("relation_l3ext_rs_l3_inst_p_to_dom_p"); ok {
		relationParam := relationTol3extRsL3InstPToDomP.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTol3extRsInstPToNatMappingEPg, ok := d.GetOk("relation_l3ext_rs_inst_p_to_nat_mapping_epg"); ok {
		relationParam := relationTol3extRsInstPToNatMappingEPg.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTofvRsConsIf, ok := d.GetOk("relation_fv_rs_cons_if"); ok {
		relationParamList := toStringList(relationTofvRsConsIf.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTofvRsCustQosPol, ok := d.GetOk("relation_fv_rs_cust_qos_pol"); ok {
		relationParam := relationTofvRsCustQosPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTofvRsCons, ok := d.GetOk("relation_fv_rs_cons"); ok {
		relationParamList := toStringList(relationTofvRsCons.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTofvRsProtBy, ok := d.GetOk("relation_fv_rs_prot_by"); ok {
		relationParamList := toStringList(relationTofvRsProtBy.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTofvRsIntraEpg, ok := d.GetOk("relation_fv_rs_intra_epg"); ok {
		relationParamList := toStringList(relationTofvRsIntraEpg.(*schema.Set).List())
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

	if relationTofvRsSecInherited, ok := d.GetOk("relation_fv_rs_sec_inherited"); ok {
		relationParamList := toStringList(relationTofvRsSecInherited.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsSecInheritedFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if relationTofvRsProv, ok := d.GetOk("relation_fv_rs_prov"); ok {
		relationParamList := toStringList(relationTofvRsProv.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsProvFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relationParamName)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if relationTol3extRsL3InstPToDomP, ok := d.GetOk("relation_l3ext_rs_l3_inst_p_to_dom_p"); ok {
		relationParam := relationTol3extRsL3InstPToDomP.(string)
		err = aciClient.CreateRelationl3extRsL3InstPToDomPFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTol3extRsInstPToNatMappingEPg, ok := d.GetOk("relation_l3ext_rs_inst_p_to_nat_mapping_epg"); ok {
		relationParam := relationTol3extRsInstPToNatMappingEPg.(string)
		err = aciClient.CreateRelationl3extRsInstPToNatMappingEPgFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTofvRsConsIf, ok := d.GetOk("relation_fv_rs_cons_if"); ok {
		relationParamList := toStringList(relationTofvRsConsIf.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsConsIfFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relationParamName)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if relationTofvRsCustQosPol, ok := d.GetOk("relation_fv_rs_cust_qos_pol"); ok {
		relationParam := relationTofvRsCustQosPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsCustQosPolFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTol3extRsInstPToProfile, ok := d.GetOk("relation_l3ext_rs_inst_p_to_profile"); ok {

		relationParamList := relationTol3extRsInstPToProfile.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationl3extRsInstPToProfileFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, GetMOName(paramMap["tn_rtctrl_profile_name"].(string)), paramMap["direction"].(string))
			if err != nil {
				return diag.FromErr(err)
			}
		}

	}
	if relationTofvRsCons, ok := d.GetOk("relation_fv_rs_cons"); ok {
		relationParamList := toStringList(relationTofvRsCons.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsConsFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relationParamName)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if relationTofvRsProtBy, ok := d.GetOk("relation_fv_rs_prot_by"); ok {
		relationParamList := toStringList(relationTofvRsProtBy.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsProtByFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relationParamName)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if relationTofvRsIntraEpg, ok := d.GetOk("relation_fv_rs_intra_epg"); ok {
		relationParamList := toStringList(relationTofvRsIntraEpg.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsIntraEpgFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relationParamName)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(l3extInstP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciExternalNetworkInstanceProfileRead(ctx, d, m)
}

func resourceAciExternalNetworkInstanceProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ExternalNetworkInstanceProfile: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	L3OutsideDn := d.Get("l3_outside_dn").(string)

	l3extInstPAttr := models.ExternalNetworkInstanceProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extInstPAttr.Annotation = Annotation.(string)
	} else {
		l3extInstPAttr.Annotation = "{}"
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
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_fv_rs_sec_inherited") {
		oldRel, newRel := d.GetChange("relation_fv_rs_sec_inherited")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fv_rs_prov") {
		oldRel, newRel := d.GetChange("relation_fv_rs_prov")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_l3ext_rs_l3_inst_p_to_dom_p") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_l3_inst_p_to_dom_p")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_l3ext_rs_inst_p_to_nat_mapping_epg") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_inst_p_to_nat_mapping_epg")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_fv_rs_cons_if") {
		oldRel, newRel := d.GetChange("relation_fv_rs_cons_if")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fv_rs_cust_qos_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_cust_qos_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_fv_rs_cons") {
		oldRel, newRel := d.GetChange("relation_fv_rs_cons")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fv_rs_prot_by") {
		oldRel, newRel := d.GetChange("relation_fv_rs_prot_by")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fv_rs_intra_epg") {
		oldRel, newRel := d.GetChange("relation_fv_rs_intra_epg")
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

	if d.HasChange("relation_fv_rs_sec_inherited") {
		oldRel, newRel := d.GetChange("relation_fv_rs_sec_inherited")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsSecInheritedFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsSecInheritedFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}
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
			if relDnName == "" {
				relDnName = relDn
			}
			err = aciClient.DeleteRelationfvRsProvFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			if relDnName == "" {
				relDnName = relDn
			}
			err = aciClient.CreateRelationfvRsProvFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}
		}

	}
	if d.HasChange("relation_l3ext_rs_l3_inst_p_to_dom_p") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_l3_inst_p_to_dom_p")
		err = aciClient.CreateRelationl3extRsL3InstPToDomPFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_l3ext_rs_inst_p_to_nat_mapping_epg") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_inst_p_to_nat_mapping_epg")
		err = aciClient.DeleteRelationl3extRsInstPToNatMappingEPgFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationl3extRsInstPToNatMappingEPgFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}
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
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationfvRsConsIfFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}
		}

	}
	if d.HasChange("relation_fv_rs_cust_qos_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_cust_qos_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationfvRsCustQosPolFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_l3ext_rs_inst_p_to_profile") {
		oldRel, newRel := d.GetChange("relation_l3ext_rs_inst_p_to_profile")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.DeleteRelationl3extRsInstPToProfileFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, GetMOName(paramMap["tn_rtctrl_profile_name"].(string)), paramMap["direction"].(string))
			if err != nil {
				return diag.FromErr(err)
			}

		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationl3extRsInstPToProfileFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, GetMOName(paramMap["tn_rtctrl_profile_name"].(string)), paramMap["direction"].(string))
			if err != nil {
				return diag.FromErr(err)
			}
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
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationfvRsConsFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
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
			relDnName := GetMOName(relDn)
			err = aciClient.DeleteRelationfvRsProtByFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationfvRsProtByFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

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
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationfvRsIntraEpgFromExternalNetworkInstanceProfile(l3extInstP.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}

	d.SetId(l3extInstP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciExternalNetworkInstanceProfileRead(ctx, d, m)

}

func resourceAciExternalNetworkInstanceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l3extInstP, err := getRemoteExternalNetworkInstanceProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setExternalNetworkInstanceProfileAttributes(l3extInstP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	fvRsSecInheritedData, err := aciClient.ReadRelationfvRsSecInheritedFromExternalNetworkInstanceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsSecInherited %v", err)
		d.Set("relation_fv_rs_sec_inherited", make([]string, 0, 1))

	} else {
		d.Set("relation_fv_rs_sec_inherited", toStringList(fvRsSecInheritedData.(*schema.Set).List()))
	}

	fvRsProvData, err := aciClient.ReadRelationfvRsProvFromExternalNetworkInstanceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsProv %v", err)
		d.Set("relation_fv_rs_prov", make([]string, 0, 1))
	} else {
		d.Set("relation_fv_rs_prov", toStringList(fvRsProvData.(*schema.Set).List()))
	}

	l3extRsL3InstPToDomPData, err := aciClient.ReadRelationl3extRsL3InstPToDomPFromExternalNetworkInstanceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsL3InstPToDomP %v", err)
		d.Set("relation_l3ext_rs_l3_inst_p_to_dom_p", "")

	} else {
		d.Set("relation_l3ext_rs_l3_inst_p_to_dom_p", l3extRsL3InstPToDomPData.(string))
	}

	l3extRsInstPToNatMappingEPgData, err := aciClient.ReadRelationl3extRsInstPToNatMappingEPgFromExternalNetworkInstanceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsInstPToNatMappingEPg %v", err)
		d.Set("relation_l3ext_rs_inst_p_to_nat_mapping_epg", "")

	} else {
		d.Set("relation_l3ext_rs_inst_p_to_nat_mapping_epg", l3extRsInstPToNatMappingEPgData.(string))
	}

	fvRsConsIfData, err := aciClient.ReadRelationfvRsConsIfFromExternalNetworkInstanceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsConsIf %v", err)
		d.Set("relation_fv_rs_cons_if", make([]string, 0, 1))

	} else {
		d.Set("relation_fv_rs_cons_if", toStringList(fvRsConsIfData.(*schema.Set).List()))
	}

	fvRsCustQosPolData, err := aciClient.ReadRelationfvRsCustQosPolFromExternalNetworkInstanceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsCustQosPol %v", err)
		d.Set("relation_fv_rs_cust_qos_pol", "")

	} else {
		d.Set("relation_fv_rs_cust_qos_pol", fvRsCustQosPolData.(string))
	}

	l3extRsInstPToProfileData, err := aciClient.ReadRelationl3extRsInstPToProfileFromExternalNetworkInstanceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsInstPToProfile %v", err)

	} else {
		relMapList := make([]map[string]string, 0, 1)
		relMaps := l3extRsInstPToProfileData.([]map[string]string)
		for _, obj := range relMaps {
			relMapList = append(relMapList, map[string]string{
				"tn_rtctrl_profile_name": obj["tnRtctrlProfileName"],
				"direction":              obj["direction"],
			})
		}
		d.Set("relation_l3ext_rs_inst_p_to_profile", relMapList)
	}

	fvRsConsData, err := aciClient.ReadRelationfvRsConsFromExternalNetworkInstanceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsCons %v", err)
		d.Set("relation_fv_rs_cons", make([]string, 0, 1))

	} else {
		d.Set("relation_fv_rs_cons", toStringList(fvRsConsData.(*schema.Set).List()))
	}

	fvRsProtByData, err := aciClient.ReadRelationfvRsProtByFromExternalNetworkInstanceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsProtBy %v", err)
		d.Set("relation_fv_rs_prot_by", make([]string, 0, 1))

	} else {
		d.Set("relation_fv_rs_prot_by", toStringList(fvRsProtByData.(*schema.Set).List()))
	}

	fvRsIntraEpgData, err := aciClient.ReadRelationfvRsIntraEpgFromExternalNetworkInstanceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsIntraEpg %v", err)
		d.Set("relation_fv_rs_intra_epg", make([]string, 0, 1))

	} else {
		d.Set("relation_fv_rs_intra_epg", toStringList(fvRsIntraEpgData.(*schema.Set).List()))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciExternalNetworkInstanceProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l3extInstP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
