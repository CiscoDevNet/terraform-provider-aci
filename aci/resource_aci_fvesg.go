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

func resourceAciEndpointSecurityGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciEndpointSecurityGroupCreate,
		UpdateContext: resourceAciEndpointSecurityGroupUpdate,
		ReadContext:   resourceAciEndpointSecurityGroupRead,
		DeleteContext: resourceAciEndpointSecurityGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciEndpointSecurityGroupImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"application_profile_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"match_t": {
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pc_enf_pref": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"enforced",
					"unenforced",
				}, false),
			},
			"pref_gr_memb": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"exclude",
					"include",
				}, false),
			},
			"relation_fv_rs_cons": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Create relation to vzBrCP",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{"prio": {
						Optional: true,
						Type:     schema.TypeString,
						ValidateFunc: validation.StringInSlice([]string{
							"level1",
							"level2",
							"level3",
							"level4",
							"level5",
							"level6",
							"unspecified",
						}, false),
					},
						"target_dn": {
							Required: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"relation_fv_rs_cons_if": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Create relation to vzCPIf",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{"prio": {
						Optional: true,
						Type:     schema.TypeString,
						ValidateFunc: validation.StringInSlice([]string{
							"level1",
							"level2",
							"level3",
							"level4",
							"level5",
							"level6",
							"unspecified",
						}, false),
					},
						"target_dn": {
							Required: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"relation_fv_rs_cust_qos_pol": {
				Type:        schema.TypeString,
				Default:     "uni/tn-common/qoscustom-default",
				Optional:    true,
				Description: "Create relation to qos:CustomPol",
			},
			"relation_fv_rs_intra_epg": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Create relation to vz:BrCP",
				Set:         schema.HashString,
			},
			"relation_fv_rs_prot_by": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Create relation to vz:Taboo",
				Set:         schema.HashString,
			},
			"relation_fv_rs_prov": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Create relation to vzBrCP",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{"match_t": {
						Optional: true,
						Type:     schema.TypeString,
						ValidateFunc: validation.StringInSlice([]string{
							"All",
							"AtleastOne",
							"AtmostOne",
							"None",
						}, false),
					}, "prio": {
						Optional: true,
						Type:     schema.TypeString,
						ValidateFunc: validation.StringInSlice([]string{
							"level1",
							"level2",
							"level3",
							"level4",
							"level5",
							"level6",
							"unspecified",
						}, false),
					},
						"target_dn": {
							Required: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"relation_fv_rs_scope": {
				Type:        schema.TypeString,
				Default:     "uni/tn-common/qoscustom-default",
				Optional:    true,
				Description: "Create relation to fv:Ctx",
			},
			"relation_fv_rs_sec_inherited": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Create relation to fv:EPg",
				Set:         schema.HashString,
			}})),
	}
}

func getRemoteEndpointSecurityGroup(client *client.Client, dn string) (*models.EndpointSecurityGroup, error) {
	fvESgCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	fvESg := models.EndpointSecurityGroupFromContainer(fvESgCont)
	if fvESg.DistinguishedName == "" {
		return nil, fmt.Errorf("EndpointSecurityGroup %s not found", fvESg.DistinguishedName)
	}
	return fvESg, nil
}

func setEndpointSecurityGroupAttributes(fvESg *models.EndpointSecurityGroup, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(fvESg.DistinguishedName)
	d.Set("description", fvESg.Description)
	fvESgMap, err := fvESg.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("annotation", fvESgMap["annotation"])
	d.Set("flood_on_encap", fvESgMap["floodOnEncap"])
	d.Set("match_t", fvESgMap["matchT"])
	d.Set("name", fvESgMap["name"])
	d.Set("application_profile_dn", GetParentDn(fvESg.DistinguishedName, fmt.Sprintf("/esg-%s", fvESgMap["name"])))
	d.Set("pc_enf_pref", fvESgMap["pcEnfPref"])
	d.Set("pref_gr_memb", fvESgMap["prefGrMemb"])
	d.Set("prio", fvESgMap["prio"])
	d.Set("name_alias", fvESgMap["nameAlias"])
	return d, nil
}

func resourceAciEndpointSecurityGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	fvESg, err := getRemoteEndpointSecurityGroup(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setEndpointSecurityGroupAttributes(fvESg, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciEndpointSecurityGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] EndpointSecurityGroup: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	ApplicationProfileDn := d.Get("application_profile_dn").(string)

	fvESgAttr := models.EndpointSecurityGroupAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvESgAttr.Annotation = Annotation.(string)
	} else {
		fvESgAttr.Annotation = "{}"
	}

	if FloodOnEncap, ok := d.GetOk("flood_on_encap"); ok {
		fvESgAttr.FloodOnEncap = FloodOnEncap.(string)
	}

	if MatchT, ok := d.GetOk("match_t"); ok {
		fvESgAttr.MatchT = MatchT.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		fvESgAttr.Name = Name.(string)
	}

	if PcEnfPref, ok := d.GetOk("pc_enf_pref"); ok {
		fvESgAttr.PcEnfPref = PcEnfPref.(string)
	}

	if PrefGrMemb, ok := d.GetOk("pref_gr_memb"); ok {
		fvESgAttr.PrefGrMemb = PrefGrMemb.(string)
	}

	if Prio, ok := d.GetOk("prio"); ok {
		fvESgAttr.Prio = Prio.(string)
	}
	fvESg := models.NewEndpointSecurityGroup(fmt.Sprintf("esg-%s", name), ApplicationProfileDn, desc, nameAlias, fvESgAttr)

	err := aciClient.Save(fvESg)
	if err != nil {
		return diag.FromErr(err)
	}
	checkDns := make([]string, 0, 1)

	if relationTofvRsCons, ok := d.GetOk("relation_fv_rs_cons"); ok {
		relationParamList := toStringList(relationTofvRsCons.(*schema.Set).List())
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

	if relationTofvRsCustQosPol, ok := d.GetOk("relation_fv_rs_cust_qos_pol"); ok {
		relationParam := relationTofvRsCustQosPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationTofvRsIntraEpg, ok := d.GetOk("relation_fv_rs_intra_epg"); ok {
		relationParamList := toStringList(relationTofvRsIntraEpg.(*schema.Set).List())
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

	if relationTofvRsProv, ok := d.GetOk("relation_fv_rs_prov"); ok {
		relationParamList := toStringList(relationTofvRsProv.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTofvRsScope, ok := d.GetOk("relation_fv_rs_scope"); ok {
		relationParam := relationTofvRsScope.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationTofvRsSecInherited, ok := d.GetOk("relation_fv_rs_sec_inherited"); ok {
		relationParamList := toStringList(relationTofvRsSecInherited.(*schema.Set).List())
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

	if relationTofvRsCons, ok := d.GetOk("relation_fv_rs_cons"); ok {
		relationParamList := relationTofvRsCons.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationfvRsCons(fvESg.DistinguishedName, fvESgAttr.Annotation, paramMap["prio"].(string), GetMOName(paramMap["target_dn"].(string)))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if relationTofvRsConsIf, ok := d.GetOk("relation_fv_rs_cons_if"); ok {
		relationParamList := relationTofvRsConsIf.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})

			err = aciClient.CreateRelationfvRsConsIf(fvESg.DistinguishedName, fvESgAttr.Annotation, paramMap["prio"].(string), GetMOName(paramMap["target_dn"].(string)))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if relationTofvRsCustQosPol, ok := d.GetOk("relation_fv_rs_cust_qos_pol"); ok {
		relationParam := relationTofvRsCustQosPol.(string)
		err = aciClient.CreateRelationfvRsCustQosPol(fvESg.DistinguishedName, fvESgAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationTofvRsIntraEpg, ok := d.GetOk("relation_fv_rs_intra_epg"); ok {
		relationParamList := toStringList(relationTofvRsIntraEpg.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsIntraEpg(fvESg.DistinguishedName, fvESgAttr.Annotation, GetMOName(relationParam))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if relationTofvRsProtBy, ok := d.GetOk("relation_fv_rs_prot_by"); ok {
		relationParamList := toStringList(relationTofvRsProtBy.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsProtBy(fvESg.DistinguishedName, fvESgAttr.Annotation, GetMOName(relationParam))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if relationTofvRsProv, ok := d.GetOk("relation_fv_rs_prov"); ok {
		relationParamList := relationTofvRsProv.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})

			err = aciClient.CreateRelationfvRsProv(fvESg.DistinguishedName, fvESgAttr.Annotation, paramMap["match_t"].(string), paramMap["prio"].(string), GetMOName(paramMap["target_dn"].(string)))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if relationTofvRsScope, ok := d.GetOk("relation_fv_rs_scope"); ok {
		relationParam := relationTofvRsScope.(string)
		err = aciClient.CreateRelationfvRsScope(fvESg.DistinguishedName, fvESgAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationTofvRsSecInherited, ok := d.GetOk("relation_fv_rs_sec_inherited"); ok {
		relationParamList := toStringList(relationTofvRsSecInherited.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsSecInherited(fvESg.DistinguishedName, fvESgAttr.Annotation, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(fvESg.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciEndpointSecurityGroupRead(ctx, d, m)
}

func resourceAciEndpointSecurityGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] EndpointSecurityGroup: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	ApplicationProfileDn := d.Get("application_profile_dn").(string)
	fvESgAttr := models.EndpointSecurityGroupAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		fvESgAttr.Annotation = Annotation.(string)
	} else {
		fvESgAttr.Annotation = "{}"
	}

	if FloodOnEncap, ok := d.GetOk("flood_on_encap"); ok {
		fvESgAttr.FloodOnEncap = FloodOnEncap.(string)
	}

	if MatchT, ok := d.GetOk("match_t"); ok {
		fvESgAttr.MatchT = MatchT.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		fvESgAttr.Name = Name.(string)
	}

	if PcEnfPref, ok := d.GetOk("pc_enf_pref"); ok {
		fvESgAttr.PcEnfPref = PcEnfPref.(string)
	}

	if PrefGrMemb, ok := d.GetOk("pref_gr_memb"); ok {
		fvESgAttr.PrefGrMemb = PrefGrMemb.(string)
	}

	if Prio, ok := d.GetOk("prio"); ok {
		fvESgAttr.Prio = Prio.(string)
	}
	fvESg := models.NewEndpointSecurityGroup(fmt.Sprintf("esg-%s", name), ApplicationProfileDn, desc, nameAlias, fvESgAttr)

	fvESg.Status = "modified"
	err := aciClient.Save(fvESg)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_fv_rs_cons") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fv_rs_cons")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())
		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fv_rs_cons_if") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fv_rs_cons_if")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())
		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fv_rs_cust_qos_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_fv_rs_cust_qos_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_fv_rs_intra_epg") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fv_rs_intra_epg")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())
		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fv_rs_prot_by") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fv_rs_prot_by")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())
		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fv_rs_prov") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fv_rs_prov")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())
		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fv_rs_scope") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_fv_rs_scope")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_fv_rs_sec_inherited") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fv_rs_sec_inherited")
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

	if d.HasChange("relation_fv_rs_cons") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fv_rs_cons")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})

			err = aciClient.DeleteRelationfvRsCons(fvESg.DistinguishedName, GetMOName(paramMap["target_dn"].(string)))
			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})

			err = aciClient.CreateRelationfvRsCons(fvESg.DistinguishedName, fvESgAttr.Annotation, paramMap["prio"].(string), GetMOName(paramMap["target_dn"].(string)))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("relation_fv_rs_cons_if") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fv_rs_cons_if")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})

			err = aciClient.DeleteRelationfvRsConsIf(fvESg.DistinguishedName, GetMOName(paramMap["target_dn"].(string)))
			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})

			err = aciClient.CreateRelationfvRsConsIf(fvESg.DistinguishedName, fvESgAttr.Annotation, paramMap["prio"].(string), GetMOName(paramMap["target_dn"].(string)))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("relation_fv_rs_cust_qos_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_fv_rs_cust_qos_pol")
		err = aciClient.DeleteRelationfvRsCustQosPol(fvESg.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationfvRsCustQosPol(fvESg.DistinguishedName, fvESgAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_fv_rs_intra_epg") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fv_rs_intra_epg")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsIntraEpg(fvESg.DistinguishedName, GetMOName(relDn))
			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsIntraEpg(fvESg.DistinguishedName, fvESgAttr.Annotation, GetMOName(relDn))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("relation_fv_rs_prot_by") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fv_rs_prot_by")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsProtBy(fvESg.DistinguishedName, GetMOName(relDn))
			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsProtBy(fvESg.DistinguishedName, fvESgAttr.Annotation, GetMOName(relDn))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("relation_fv_rs_prov") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fv_rs_prov")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})

			err = aciClient.DeleteRelationfvRsProv(fvESg.DistinguishedName, GetMOName(paramMap["target_dn"].(string)))
			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})

			err = aciClient.CreateRelationfvRsProv(fvESg.DistinguishedName, fvESgAttr.Annotation, paramMap["match_t"].(string), paramMap["prio"].(string), GetMOName(paramMap["target_dn"].(string)))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("relation_fv_rs_scope") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_fv_rs_scope")
		err = aciClient.DeleteRelationfvRsScope(fvESg.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationfvRsScope(fvESg.DistinguishedName, fvESgAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_fv_rs_sec_inherited") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fv_rs_sec_inherited")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsSecInherited(fvESg.DistinguishedName, relDn)

			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsSecInherited(fvESg.DistinguishedName, fvESgAttr.Annotation, relDn)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(fvESg.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciEndpointSecurityGroupRead(ctx, d, m)
}

func resourceAciEndpointSecurityGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	fvESg, err := getRemoteEndpointSecurityGroup(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	_, err = setEndpointSecurityGroupAttributes(fvESg, d)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	fvRsConsData, err := aciClient.ReadRelationfvRsCons(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsCons %v", err)
	} else {
		listRelMap := make([]map[string]string, 0, 1)
		listfvRsConsData := fvRsConsData.([]map[string]string)
		for _, obj := range listfvRsConsData {
			listRelMap = append(listRelMap, map[string]string{
				"prio":      obj["prio"],
				"target_dn": obj["tDn"],
			})
		}
		d.Set("relation_fv_rs_cons", listRelMap)
	}

	fvRsConsIfData, err := aciClient.ReadRelationfvRsConsIf(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsConsIf %v", err)
	} else {
		listRelMap := make([]map[string]string, 0, 1)
		listfvRsConsIfData := fvRsConsIfData.([]map[string]string)
		for _, obj := range listfvRsConsIfData {
			listRelMap = append(listRelMap, map[string]string{
				"prio":      obj["prio"],
				"target_dn": obj["tDn"],
			})
		}
		d.Set("relation_fv_rs_cons_if", listRelMap)
	}

	fvRsCustQosPolData, err := aciClient.ReadRelationfvRsCustQosPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsCustQosPol %v", err)
		d.Set("fv_rs_cust_qos_pol", "")
	} else {
		d.Set("relation_fv_rs_cust_qos_pol", fvRsCustQosPolData.(string))
	}
	fvRsIntraEpgData, err := aciClient.ReadRelationfvRsIntraEpg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsIntraEpg %v", err)
		d.Set("relation_fv_rs_intra_epg", make([]string, 0, 1))
	} else {
		d.Set("relation_fv_rs_intra_epg", toStringList(fvRsIntraEpgData.(*schema.Set).List()))
	}
	fvRsProtByData, err := aciClient.ReadRelationfvRsProtBy(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsProtBy %v", err)
		d.Set("relation_fv_rs_prot_by", make([]string, 0, 1))
	} else {
		d.Set("relation_fv_rs_prot_by", toStringList(fvRsProtByData.(*schema.Set).List()))
	}

	fvRsProvData, err := aciClient.ReadRelationfvRsProv(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsProv %v", err)
	} else {
		listRelMap := make([]map[string]string, 0, 1)
		listfvRsProvData := fvRsProvData.([]map[string]string)
		for _, obj := range listfvRsProvData {
			listRelMap = append(listRelMap, map[string]string{
				"prio":      obj["prio"],
				"target_dn": obj["tDn"],
				"match_t":   obj["matchT"],
			})
		}
		d.Set("relation_fv_rs_prov", listRelMap)
	}

	fvRsScopeData, err := aciClient.ReadRelationfvRsScope(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsScope %v", err)
		d.Set("fv_rs_scope", "")
	} else {
		d.Set("relation_fv_rs_scope", fvRsScopeData.(string))
	}
	fvRsSecInheritedData, err := aciClient.ReadRelationfvRsSecInherited(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsSecInherited %v", err)
		d.Set("relation_fv_rs_sec_inherited", make([]string, 0, 1))
	} else {
		d.Set("relation_fv_rs_sec_inherited", toStringList(fvRsSecInheritedData.(*schema.Set).List()))
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciEndpointSecurityGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fvESg")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
