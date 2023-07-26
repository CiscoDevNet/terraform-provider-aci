package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciApplicationEPG() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciApplicationEPGCreate,
		UpdateContext: resourceAciApplicationEPGUpdate,
		ReadContext:   resourceAciApplicationEPGRead,
		DeleteContext: resourceAciApplicationEPGDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciApplicationEPGImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"application_profile_dn": &schema.Schema{
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

			"fwd_ctrl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"none",
					"proxy-arp",
				}, false),
			},

			"has_mcast_source": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"is_attr_based_epg": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
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

			"pc_enf_pref": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"enforced",
					"unenforced",
				}, false),
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

			"shutdown": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_fv_rs_bd": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_fv_rs_cust_qos_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_fv_rs_fc_path_att": &schema.Schema{
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
			"relation_fv_rs_node_att": &schema.Schema{ // also have to change read function
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"node_dn": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"encap": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"deployment_immediacy": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"immediate",
								"lazy",
							}, false),
							Default: "lazy",
						},
						"mode": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"regular",
								"native",
								"untagged",
							}, false),
							Default: "regular",
						},
						"primary_encap": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  "unknown",
						},
					},
				},
			},
			"relation_fv_rs_dpp_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_cons": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
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
			"relation_fv_rs_prot_by": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_fv_rs_path_att": &schema.Schema{
				Type:       schema.TypeSet,
				Elem:       &schema.Schema{Type: schema.TypeString},
				Optional:   true,
				Set:        schema.HashString,
				Deprecated: "use resource aci_epg_to_static_path instead",
			},
			"relation_fv_rs_aepg_mon_pol": &schema.Schema{
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

	if fvAEPg.DistinguishedName == "" || fvAEPg.DistinguishedName == "{}" {
		return nil, fmt.Errorf("Application EPG %s not found", dn)
	}

	return fvAEPg, nil
}

func setApplicationEPGAttributes(fvAEPg *models.ApplicationEPG, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(fvAEPg.DistinguishedName)
	d.Set("description", fvAEPg.Description)
	if dn != fvAEPg.DistinguishedName {
		d.Set("application_profile_dn", "")
	}
	fvAEPgMap, err := fvAEPg.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("application_profile_dn", GetParentDn(dn, fmt.Sprintf("/epg-%s", fvAEPgMap["name"])))

	d.Set("name", fvAEPgMap["name"])

	d.Set("annotation", fvAEPgMap["annotation"])
	d.Set("exception_tag", fvAEPgMap["exceptionTag"])
	d.Set("flood_on_encap", fvAEPgMap["floodOnEncap"])
	if fvAEPgMap["fwdCtrl"] == "" {
		d.Set("fwd_ctrl", "none")
	} else {
		d.Set("fwd_ctrl", fvAEPgMap["fwdCtrl"])
	}
	d.Set("has_mcast_source", fvAEPgMap["hasMcastSource"])
	d.Set("is_attr_based_epg", fvAEPgMap["isAttrBasedEPg"])
	d.Set("match_t", fvAEPgMap["matchT"])
	d.Set("name_alias", fvAEPgMap["nameAlias"])
	d.Set("pc_enf_pref", fvAEPgMap["pcEnfPref"])
	d.Set("pref_gr_memb", fvAEPgMap["prefGrMemb"])
	d.Set("prio", fvAEPgMap["prio"])
	d.Set("shutdown", fvAEPgMap["shutdown"])
	return d, nil
}

func resourceAciApplicationEPGImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fvAEPg, err := getRemoteApplicationEPG(aciClient, dn)

	if err != nil {
		return nil, err
	}
	fvAEPgMap, err := fvAEPg.ToMap()
	if err != nil {
		return nil, err
	}
	name := fvAEPgMap["name"]
	pDN := GetParentDn(dn, fmt.Sprintf("/epg-%s", name))
	d.Set("application_profile_dn", pDN)
	schemaFilled, err := setApplicationEPGAttributes(fvAEPg, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciApplicationEPGCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ApplicationEPG: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	ApplicationProfileDn := d.Get("application_profile_dn").(string)

	fvAEPgAttr := models.ApplicationEPGAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvAEPgAttr.Annotation = Annotation.(string)
	} else {
		fvAEPgAttr.Annotation = "{}"
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
	if IsAttrBasedEPg, ok := d.GetOk("is_attr_based_epg"); ok {
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
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTofvRsBd, ok := d.GetOk("relation_fv_rs_bd"); ok {
		relationParam := relationTofvRsBd.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTofvRsCustQosPol, ok := d.GetOk("relation_fv_rs_cust_qos_pol"); ok {
		relationParam := relationTofvRsCustQosPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTofvRsFcPathAtt, ok := d.GetOk("relation_fv_rs_fc_path_att"); ok {
		relationParamList := toStringList(relationTofvRsFcPathAtt.(*schema.Set).List())
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

	if relationTofvRsConsIf, ok := d.GetOk("relation_fv_rs_cons_if"); ok {
		relationParamList := toStringList(relationTofvRsConsIf.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTofvRsSecInherited, ok := d.GetOk("relation_fv_rs_sec_inherited"); ok {
		relationParamList := toStringList(relationTofvRsSecInherited.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTofvRsNodeAtt, ok := d.GetOk("relation_fv_rs_node_att"); ok {
		relationParamList := toStringList(relationTofvRsNodeAtt.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTofvRsDppPol, ok := d.GetOk("relation_fv_rs_dpp_pol"); ok {
		relationParam := relationTofvRsDppPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTofvRsCons, ok := d.GetOk("relation_fv_rs_cons"); ok {
		relationParamList := toStringList(relationTofvRsCons.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTofvRsProvDef, ok := d.GetOk("relation_fv_rs_prov_def"); ok {
		relationParamList := toStringList(relationTofvRsProvDef.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTofvRsTrustCtrl, ok := d.GetOk("relation_fv_rs_trust_ctrl"); ok {
		relationParam := relationTofvRsTrustCtrl.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTofvRsPathAtt, ok := d.GetOk("relation_fv_rs_path_att"); ok {
		relationParamList := toStringList(relationTofvRsPathAtt.(*schema.Set).List())
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

	if relationTofvRsAEPgMonPol, ok := d.GetOk("relation_fv_rs_aepg_mon_pol"); ok {
		relationParam := relationTofvRsAEPgMonPol.(string)
		checkDns = append(checkDns, relationParam)
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

	if relationTofvRsBd, ok := d.GetOk("relation_fv_rs_bd"); ok {
		relationParam := relationTofvRsBd.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsBdFromApplicationEPG(fvAEPg.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationTofvRsCustQosPol, ok := d.GetOk("relation_fv_rs_cust_qos_pol"); ok {
		relationParam := relationTofvRsCustQosPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsCustQosPolFromApplicationEPG(fvAEPg.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTofvRsFcPathAtt, ok := d.GetOk("relation_fv_rs_fc_path_att"); ok {
		relationParamList := toStringList(relationTofvRsFcPathAtt.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsFcPathAttFromApplicationEPG(fvAEPg.DistinguishedName, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if relationTofvRsProv, ok := d.GetOk("relation_fv_rs_prov"); ok {
		relationParamList := toStringList(relationTofvRsProv.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsProvFromApplicationEPG(fvAEPg.DistinguishedName, relationParamName)

			if err != nil {
				return diag.FromErr(err)
			}

		}
	}
	if relationTofvRsConsIf, ok := d.GetOk("relation_fv_rs_cons_if"); ok {
		relationParamList := toStringList(relationTofvRsConsIf.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsConsIfFromApplicationEPG(fvAEPg.DistinguishedName, relationParamName)

			if err != nil {
				return diag.FromErr(err)
			}

		}
	}
	if relationTofvRsSecInherited, ok := d.GetOk("relation_fv_rs_sec_inherited"); ok {
		relationParamList := toStringList(relationTofvRsSecInherited.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsSecInheritedFromApplicationEPG(fvAEPg.DistinguishedName, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}

		}
	}
	if relationTofvRsNodeAtt, ok := d.GetOk("relation_fv_rs_node_att"); ok {
		relationParamList := relationTofvRsNodeAtt.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationfvRsNodeAtt(fvAEPg.DistinguishedName, paramMap["encap"].(string), paramMap["mode"].(string), paramMap["description"].(string), paramMap["deployment_immediacy"].(string), paramMap["node_dn"].(string))

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if relationTofvRsDppPol, ok := d.GetOk("relation_fv_rs_dpp_pol"); ok {
		relationParam := relationTofvRsDppPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsDppPolFromApplicationEPG(fvAEPg.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTofvRsCons, ok := d.GetOk("relation_fv_rs_cons"); ok {
		relationParamList := toStringList(relationTofvRsCons.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsConsFromApplicationEPG(fvAEPg.DistinguishedName, relationParamName)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if relationTofvRsProvDef, ok := d.GetOk("relation_fv_rs_prov_def"); ok {
		relationParamList := toStringList(relationTofvRsProvDef.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsProvDefFromApplicationEPG(fvAEPg.DistinguishedName, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if relationTofvRsTrustCtrl, ok := d.GetOk("relation_fv_rs_trust_ctrl"); ok {
		relationParam := relationTofvRsTrustCtrl.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsTrustCtrlFromApplicationEPG(fvAEPg.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTofvRsPathAtt, ok := d.GetOk("relation_fv_rs_path_att"); ok {
		relationParamList := toStringList(relationTofvRsPathAtt.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsPathAttFromApplicationEPG(fvAEPg.DistinguishedName, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}

		}
	}
	if relationTofvRsProtBy, ok := d.GetOk("relation_fv_rs_prot_by"); ok {
		relationParamList := toStringList(relationTofvRsProtBy.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsProtByFromApplicationEPG(fvAEPg.DistinguishedName, relationParamName)

			if err != nil {
				return diag.FromErr(err)
			}

		}
	}
	if relationTofvRsAEPgMonPol, ok := d.GetOk("relation_fv_rs_aepg_mon_pol"); ok {
		relationParam := relationTofvRsAEPgMonPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsAEPgMonPolFromApplicationEPG(fvAEPg.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTofvRsIntraEpg, ok := d.GetOk("relation_fv_rs_intra_epg"); ok {
		relationParamList := toStringList(relationTofvRsIntraEpg.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsIntraEpgFromApplicationEPG(fvAEPg.DistinguishedName, relationParamName)

			if err != nil {
				return diag.FromErr(err)
			}

		}
	}

	d.SetId(fvAEPg.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciApplicationEPGRead(ctx, d, m)
}

func resourceAciApplicationEPGUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ApplicationEPG: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	ApplicationProfileDn := d.Get("application_profile_dn").(string)

	fvAEPgAttr := models.ApplicationEPGAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvAEPgAttr.Annotation = Annotation.(string)
	} else {
		fvAEPgAttr.Annotation = "{}"
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
	if IsAttrBasedEPg, ok := d.GetOk("is_attr_based_epg"); ok {
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
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_fv_rs_bd") {
		_, newRelParam := d.GetChange("relation_fv_rs_bd")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_fv_rs_cust_qos_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_cust_qos_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_fv_rs_fc_path_att") {
		oldRel, newRel := d.GetChange("relation_fv_rs_fc_path_att")
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

	if d.HasChange("relation_fv_rs_cons_if") {
		oldRel, newRel := d.GetChange("relation_fv_rs_cons_if")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fv_rs_sec_inherited") {
		oldRel, newRel := d.GetChange("relation_fv_rs_sec_inherited")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fv_rs_node_att") {
		oldRel, newRel := d.GetChange("relation_fv_rs_node_att")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fv_rs_dpp_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_dpp_pol")
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

	if d.HasChange("relation_fv_rs_prov_def") {
		oldRel, newRel := d.GetChange("relation_fv_rs_prov_def")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fv_rs_trust_ctrl") {
		_, newRelParam := d.GetChange("relation_fv_rs_trust_ctrl")
		checkDns = append(checkDns, newRelParam.(string))

	}
	if d.HasChange("relation_fv_rs_path_att") {
		oldRel, newRel := d.GetChange("relation_fv_rs_path_att")
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

	if d.HasChange("relation_fv_rs_aepg_mon_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_aepg_mon_pol")
		checkDns = append(checkDns, newRelParam.(string))
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

	if d.HasChange("relation_fv_rs_bd") {
		_, newRelParam := d.GetChange("relation_fv_rs_bd")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationfvRsBdFromApplicationEPG(fvAEPg.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_fv_rs_cust_qos_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_cust_qos_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationfvRsCustQosPolFromApplicationEPG(fvAEPg.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
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
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsFcPathAttFromApplicationEPG(fvAEPg.DistinguishedName, relDn)
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
			err = aciClient.DeleteRelationfvRsProvFromApplicationEPG(fvAEPg.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationfvRsProvFromApplicationEPG(fvAEPg.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
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
			relDnName := GetMOName(relDn)
			err = aciClient.DeleteRelationfvRsConsIfFromApplicationEPG(fvAEPg.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationfvRsConsIfFromApplicationEPG(fvAEPg.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
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
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsSecInheritedFromApplicationEPG(fvAEPg.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}
	if d.HasChange("relation_fv_rs_node_att") {
		oldRel, newRel := d.GetChange("relation_fv_rs_node_att")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.DeleteRelationfvRsNodeAttFromApplicationEPG(fvAEPg.DistinguishedName, paramMap["node_dn"].(string))

			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationfvRsNodeAtt(fvAEPg.DistinguishedName, paramMap["encap"].(string), paramMap["mode"].(string), paramMap["description"].(string), paramMap["deployment_immediacy"].(string), paramMap["node_dn"].(string))

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("relation_fv_rs_dpp_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_dpp_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationfvRsDppPolFromApplicationEPG(fvAEPg.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationfvRsDppPolFromApplicationEPG(fvAEPg.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
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
			err = aciClient.DeleteRelationfvRsConsFromApplicationEPG(fvAEPg.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationfvRsConsFromApplicationEPG(fvAEPg.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

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
				return diag.FromErr(err)
			}

		}

	}
	if d.HasChange("relation_fv_rs_trust_ctrl") {
		_, newRelParam := d.GetChange("relation_fv_rs_trust_ctrl")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationfvRsTrustCtrlFromApplicationEPG(fvAEPg.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationfvRsTrustCtrlFromApplicationEPG(fvAEPg.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
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
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsPathAttFromApplicationEPG(fvAEPg.DistinguishedName, relDn)
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
			err = aciClient.DeleteRelationfvRsProtByFromApplicationEPG(fvAEPg.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationfvRsProtByFromApplicationEPG(fvAEPg.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}
	if d.HasChange("relation_fv_rs_aepg_mon_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_aepg_mon_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationfvRsAEPgMonPolFromApplicationEPG(fvAEPg.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationfvRsAEPgMonPolFromApplicationEPG(fvAEPg.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
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
			err = aciClient.DeleteRelationfvRsIntraEpgFromApplicationEPG(fvAEPg.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationfvRsIntraEpgFromApplicationEPG(fvAEPg.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}

	d.SetId(fvAEPg.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciApplicationEPGRead(ctx, d, m)

}

func resourceAciApplicationEPGRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fvAEPg, err := getRemoteApplicationEPG(aciClient, dn)

	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	_, err = setApplicationEPGAttributes(fvAEPg, d)
	if err != nil {
		d.SetId("")
		return nil
	}
	fvRsBdData, err := aciClient.ReadRelationfvRsBdFromApplicationEPG(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsBd %v", err)
		d.Set("relation_fv_rs_bd", "")

	} else {
		setRelationAttribute(d, "relation_fv_rs_bd", fvRsBdData.(string))
	}

	fvRsCustQosPolData, err := aciClient.ReadRelationfvRsCustQosPolFromApplicationEPG(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsCustQosPol %v", err)
		d.Set("relation_fv_rs_cust_qos_pol", "")

	} else {
		setRelationAttribute(d, "relation_fv_rs_cust_qos_pol", fvRsCustQosPolData.(string))
	}

	fvRsFcPathAttData, err := aciClient.ReadRelationfvRsFcPathAttFromApplicationEPG(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsFcPathAtt %v", err)
		setRelationAttribute(d, "relation_fv_rs_fc_path_att", make([]interface{}, 0, 1))
	} else {
		setRelationAttribute(d, "relation_fv_rs_fc_path_att", toStringList(fvRsFcPathAttData.(*schema.Set).List()))
	}

	fvRsProvData, err := aciClient.ReadRelationfvRsProvFromApplicationEPG(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsProv %v", err)
		setRelationAttribute(d, "relation_fv_rs_prov", make([]interface{}, 0, 1))
	} else {
		setRelationAttribute(d, "relation_fv_rs_prov", toStringList(fvRsProvData.(*schema.Set).List()))
	}

	fvRsConsIfData, err := aciClient.ReadRelationfvRsConsIfFromApplicationEPG(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsConsIf %v", err)
		setRelationAttribute(d, "relation_fv_rs_cons_if", make([]interface{}, 0, 1))
	} else {
		setRelationAttribute(d, "relation_fv_rs_cons_if", toStringList(fvRsConsIfData.(*schema.Set).List()))
	}

	fvRsSecInheritedData, err := aciClient.ReadRelationfvRsSecInheritedFromApplicationEPG(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsSecInherited %v", err)
		setRelationAttribute(d, "relation_fv_rs_sec_inherited", make([]interface{}, 0, 1))
	} else {
		setRelationAttribute(d, "relation_fv_rs_sec_inherited", toStringList(fvRsSecInheritedData.(*schema.Set).List()))
	}

	fvRsNodeAttData, err := aciClient.ReadRelationfvRsNodeAtt(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsNodeAtt %v", err)
	} else {
		relParams := make([]map[string]string, 0, 1)
		relParamsList := fvRsNodeAttData.([]map[string]string)
		for _, obj := range relParamsList {
			relParams = append(relParams, map[string]string{
				"encap":                obj["encap"],
				"description":          obj["descr"],
				"deployment_immediacy": obj["instrImedcy"],
				"mode":                 obj["mode"],
			})
		}
		d.Set("relation_fv_rs_node_att", relParams)
	}

	fvRsDppPolData, err := aciClient.ReadRelationfvRsDppPolFromApplicationEPG(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsDppPol %v", err)
		d.Set("relation_fv_rs_dpp_pol", "")

	} else {
		setRelationAttribute(d, "relation_fv_rs_dpp_pol", fvRsDppPolData.(string))
	}

	fvRsConsData, err := aciClient.ReadRelationfvRsConsFromApplicationEPG(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsCons %v", err)
		setRelationAttribute(d, "relation_fv_rs_cons", make([]interface{}, 0, 1))
	} else {
		setRelationAttribute(d, "relation_fv_rs_cons", toStringList(fvRsConsData.(*schema.Set).List()))
	}

	fvRsProvDefData, err := aciClient.ReadRelationfvRsProvDefFromApplicationEPG(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsProvDef %v", err)
		setRelationAttribute(d, "relation_fv_rs_prov_def", make([]interface{}, 0, 1))
	} else {
		setRelationAttribute(d, "relation_fv_rs_prov_def", toStringList(fvRsProvDefData.(*schema.Set).List()))
	}

	fvRsTrustCtrlData, err := aciClient.ReadRelationfvRsTrustCtrlFromApplicationEPG(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsTrustCtrl %v", err)
		d.Set("relation_fv_rs_trust_ctrl", "")

	} else {
		setRelationAttribute(d, "relation_fv_rs_trust_ctrl", fvRsTrustCtrlData.(string))
	}

	fvRsPathAttData, err := aciClient.ReadRelationfvRsPathAttFromApplicationEPG(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsPathAtt %v", err)
		setRelationAttribute(d, "relation_fv_rs_path_att", make([]interface{}, 0, 1))
	} else {
		setRelationAttribute(d, "relation_fv_rs_path_att", toStringList(fvRsPathAttData.(*schema.Set).List()))
	}

	fvRsProtByData, err := aciClient.ReadRelationfvRsProtByFromApplicationEPG(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsProtBy %v", err)
		setRelationAttribute(d, "relation_fv_rs_prot_by", make([]interface{}, 0, 1))
	} else {
		setRelationAttribute(d, "relation_fv_rs_prot_by", toStringList(fvRsProtByData.(*schema.Set).List()))
	}

	fvRsAEPgMonPolData, err := aciClient.ReadRelationfvRsAEPgMonPolFromApplicationEPG(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsAEPgMonPol %v", err)
		d.Set("relation_fv_rs_aepg_mon_pol", "")

	} else {
		setRelationAttribute(d, "relation_fv_rs_aepg_mon_pol", fvRsAEPgMonPolData.(string))
	}

	fvRsIntraEpgData, err := aciClient.ReadRelationfvRsIntraEpgFromApplicationEPG(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsIntraEpg %v", err)
		setRelationAttribute(d, "relation_fv_rs_intra_epg", make([]interface{}, 0, 1))
	} else {
		setRelationAttribute(d, "relation_fv_rs_intra_epg", toStringList(fvRsIntraEpgData.(*schema.Set).List()))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciApplicationEPGDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fvAEPg")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
