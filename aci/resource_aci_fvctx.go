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

func resourceAciVRF() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciVRFCreate,
		UpdateContext: resourceAciVRFUpdate,
		ReadContext:   resourceAciVRFRead,
		DeleteContext: resourceAciVRFDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciVRFImport,
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

			"bd_enforced_enable": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"ip_data_plane_learning": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"knw_mcast_act": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"deny",
					"permit",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pc_enf_dir": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ingress",
					"egress",
				}, false),
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

			"relation_fv_rs_ospf_ctx_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_vrf_validation_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_ctx_mcast_to": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_fv_rs_ctx_to_eigrp_ctx_af_pol": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tn_eigrp_ctx_af_pol_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"af": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"relation_fv_rs_ctx_to_ospf_ctx_pol": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tn_ospf_ctx_pol_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"af": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"relation_fv_rs_ctx_to_ep_ret": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_bgp_ctx_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_ctx_mon_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_ctx_to_ext_route_tag_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_ctx_to_bgp_ctx_af_pol": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tn_bgp_ctx_af_pol_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"af": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		}),
	}
}
func getRemoteVRF(client *client.Client, dn string) (*models.VRF, error) {
	fvCtxCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fvCtx := models.VRFFromContainer(fvCtxCont)

	if fvCtx.DistinguishedName == "" {
		return nil, fmt.Errorf("VRF %s not found", fvCtx.DistinguishedName)
	}

	return fvCtx, nil
}

func setVRFAttributes(fvCtx *models.VRF, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(fvCtx.DistinguishedName)
	d.Set("description", fvCtx.Description)

	if dn != fvCtx.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	fvCtxMap, err := fvCtx.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("tenant_dn", GetParentDn(dn, fmt.Sprintf("/ctx-%s", fvCtxMap["name"])))
	d.Set("name", fvCtxMap["name"])

	d.Set("annotation", fvCtxMap["annotation"])
	d.Set("bd_enforced_enable", fvCtxMap["bdEnforcedEnable"])
	d.Set("ip_data_plane_learning", fvCtxMap["ipDataPlaneLearning"])
	d.Set("knw_mcast_act", fvCtxMap["knwMcastAct"])
	d.Set("name_alias", fvCtxMap["nameAlias"])
	d.Set("pc_enf_dir", fvCtxMap["pcEnfDir"])
	d.Set("pc_enf_pref", fvCtxMap["pcEnfPref"])
	return d, nil
}

func resourceAciVRFImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fvCtx, err := getRemoteVRF(aciClient, dn)
	if err != nil {
		return nil, err
	}
	fvCtxMap, err := fvCtx.ToMap()
	if err != nil {
		return nil, err
	}
	name := fvCtxMap["name"]
	pDN := GetParentDn(dn, fmt.Sprintf("/ctx-%s", name))
	d.Set("tenant_dn", pDN)
	schemaFilled, err := setVRFAttributes(fvCtx, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciVRFCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] VRF: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	fvCtxAttr := models.VRFAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvCtxAttr.Annotation = Annotation.(string)
	} else {
		fvCtxAttr.Annotation = "{}"
	}
	if BdEnforcedEnable, ok := d.GetOk("bd_enforced_enable"); ok {
		fvCtxAttr.BdEnforcedEnable = BdEnforcedEnable.(string)
	}
	if IpDataPlaneLearning, ok := d.GetOk("ip_data_plane_learning"); ok {
		fvCtxAttr.IpDataPlaneLearning = IpDataPlaneLearning.(string)
	}
	if KnwMcastAct, ok := d.GetOk("knw_mcast_act"); ok {
		fvCtxAttr.KnwMcastAct = KnwMcastAct.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvCtxAttr.NameAlias = NameAlias.(string)
	}
	if PcEnfDir, ok := d.GetOk("pc_enf_dir"); ok {
		fvCtxAttr.PcEnfDir = PcEnfDir.(string)
	}
	if PcEnfPref, ok := d.GetOk("pc_enf_pref"); ok {
		fvCtxAttr.PcEnfPref = PcEnfPref.(string)
	}
	fvCtx := models.NewVRF(fmt.Sprintf("ctx-%s", name), TenantDn, desc, fvCtxAttr)

	err := aciClient.Save(fvCtx)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTofvRsOspfCtxPol, ok := d.GetOk("relation_fv_rs_ospf_ctx_pol"); ok {
		relationParam := relationTofvRsOspfCtxPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTofvRsVrfValidationPol, ok := d.GetOk("relation_fv_rs_vrf_validation_pol"); ok {
		relationParam := relationTofvRsVrfValidationPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTofvRsCtxMcastTo, ok := d.GetOk("relation_fv_rs_ctx_mcast_to"); ok {
		relationParamList := toStringList(relationTofvRsCtxMcastTo.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTofvRsCtxToEpRet, ok := d.GetOk("relation_fv_rs_ctx_to_ep_ret"); ok {
		relationParam := relationTofvRsCtxToEpRet.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTofvRsBgpCtxPol, ok := d.GetOk("relation_fv_rs_bgp_ctx_pol"); ok {
		relationParam := relationTofvRsBgpCtxPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTofvRsCtxMonPol, ok := d.GetOk("relation_fv_rs_ctx_mon_pol"); ok {
		relationParam := relationTofvRsCtxMonPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTofvRsCtxToExtRouteTagPol, ok := d.GetOk("relation_fv_rs_ctx_to_ext_route_tag_pol"); ok {
		relationParam := relationTofvRsCtxToExtRouteTagPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTofvRsOspfCtxPol, ok := d.GetOk("relation_fv_rs_ospf_ctx_pol"); ok {
		relationParam := relationTofvRsOspfCtxPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsOspfCtxPolFromVRF(fvCtx.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTofvRsVrfValidationPol, ok := d.GetOk("relation_fv_rs_vrf_validation_pol"); ok {
		relationParam := relationTofvRsVrfValidationPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsVrfValidationPolFromVRF(fvCtx.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTofvRsCtxMcastTo, ok := d.GetOk("relation_fv_rs_ctx_mcast_to"); ok {
		relationParamList := toStringList(relationTofvRsCtxMcastTo.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsCtxMcastToFromVRF(fvCtx.DistinguishedName, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}

		}
	}
	if relationTofvRsCtxToEigrpCtxAfPol, ok := d.GetOk("relation_fv_rs_ctx_to_eigrp_ctx_af_pol"); ok {

		relationParamList := relationTofvRsCtxToEigrpCtxAfPol.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationfvRsCtxToEigrpCtxAfPolFromVRF(fvCtx.DistinguishedName, paramMap["tn_eigrp_ctx_af_pol_name"].(string), paramMap["af"].(string))
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}
	if relationTofvRsCtxToOspfCtxPol, ok := d.GetOk("relation_fv_rs_ctx_to_ospf_ctx_pol"); ok {

		relationParamList := relationTofvRsCtxToOspfCtxPol.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationfvRsCtxToOspfCtxPolFromVRF(fvCtx.DistinguishedName, paramMap["tn_ospf_ctx_pol_name"].(string), paramMap["af"].(string))
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}
	if relationTofvRsCtxToEpRet, ok := d.GetOk("relation_fv_rs_ctx_to_ep_ret"); ok {
		relationParam := relationTofvRsCtxToEpRet.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsCtxToEpRetFromVRF(fvCtx.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTofvRsBgpCtxPol, ok := d.GetOk("relation_fv_rs_bgp_ctx_pol"); ok {
		relationParam := relationTofvRsBgpCtxPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsBgpCtxPolFromVRF(fvCtx.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTofvRsCtxMonPol, ok := d.GetOk("relation_fv_rs_ctx_mon_pol"); ok {
		relationParam := relationTofvRsCtxMonPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsCtxMonPolFromVRF(fvCtx.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTofvRsCtxToExtRouteTagPol, ok := d.GetOk("relation_fv_rs_ctx_to_ext_route_tag_pol"); ok {
		relationParam := relationTofvRsCtxToExtRouteTagPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsCtxToExtRouteTagPolFromVRF(fvCtx.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTofvRsCtxToBgpCtxAfPol, ok := d.GetOk("relation_fv_rs_ctx_to_bgp_ctx_af_pol"); ok {

		relationParamList := relationTofvRsCtxToBgpCtxAfPol.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationfvRsCtxToBgpCtxAfPolFromVRF(fvCtx.DistinguishedName, paramMap["tn_bgp_ctx_af_pol_name"].(string), paramMap["af"].(string))
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}

	d.SetId(fvCtx.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciVRFRead(ctx, d, m)
}

func resourceAciVRFUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] VRF: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	fvCtxAttr := models.VRFAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvCtxAttr.Annotation = Annotation.(string)
	} else {
		fvCtxAttr.Annotation = "{}"
	}
	if BdEnforcedEnable, ok := d.GetOk("bd_enforced_enable"); ok {
		fvCtxAttr.BdEnforcedEnable = BdEnforcedEnable.(string)
	}
	if IpDataPlaneLearning, ok := d.GetOk("ip_data_plane_learning"); ok {
		fvCtxAttr.IpDataPlaneLearning = IpDataPlaneLearning.(string)
	}
	if KnwMcastAct, ok := d.GetOk("knw_mcast_act"); ok {
		fvCtxAttr.KnwMcastAct = KnwMcastAct.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvCtxAttr.NameAlias = NameAlias.(string)
	}
	if PcEnfDir, ok := d.GetOk("pc_enf_dir"); ok {
		fvCtxAttr.PcEnfDir = PcEnfDir.(string)
	}
	if PcEnfPref, ok := d.GetOk("pc_enf_pref"); ok {
		fvCtxAttr.PcEnfPref = PcEnfPref.(string)
	}
	fvCtx := models.NewVRF(fmt.Sprintf("ctx-%s", name), TenantDn, desc, fvCtxAttr)

	fvCtx.Status = "modified"

	err := aciClient.Save(fvCtx)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_fv_rs_ospf_ctx_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_ospf_ctx_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_fv_rs_vrf_validation_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_vrf_validation_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_fv_rs_ctx_mcast_to") {
		oldRel, newRel := d.GetChange("relation_fv_rs_ctx_mcast_to")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fv_rs_ctx_to_ep_ret") {
		_, newRelParam := d.GetChange("relation_fv_rs_ctx_to_ep_ret")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_fv_rs_bgp_ctx_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_bgp_ctx_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_fv_rs_ctx_mon_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_ctx_mon_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_fv_rs_ctx_to_ext_route_tag_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_ctx_to_ext_route_tag_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_fv_rs_ospf_ctx_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_ospf_ctx_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationfvRsOspfCtxPolFromVRF(fvCtx.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_fv_rs_vrf_validation_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_vrf_validation_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationfvRsVrfValidationPolFromVRF(fvCtx.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_fv_rs_ctx_mcast_to") {
		oldRel, newRel := d.GetChange("relation_fv_rs_ctx_mcast_to")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsCtxMcastToFromVRF(fvCtx.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}
	if d.HasChange("relation_fv_rs_ctx_to_eigrp_ctx_af_pol") {
		oldRel, newRel := d.GetChange("relation_fv_rs_ctx_to_eigrp_ctx_af_pol")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.DeleteRelationfvRsCtxToEigrpCtxAfPolFromVRF(fvCtx.DistinguishedName, paramMap["tn_eigrp_ctx_af_pol_name"].(string), paramMap["af"].(string))
			if err != nil {
				return diag.FromErr(err)
			}

		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationfvRsCtxToEigrpCtxAfPolFromVRF(fvCtx.DistinguishedName, paramMap["tn_eigrp_ctx_af_pol_name"].(string), paramMap["af"].(string))
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}
	if d.HasChange("relation_fv_rs_ctx_to_ospf_ctx_pol") {
		oldRel, newRel := d.GetChange("relation_fv_rs_ctx_to_ospf_ctx_pol")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.DeleteRelationfvRsCtxToOspfCtxPolFromVRF(fvCtx.DistinguishedName, paramMap["tn_ospf_ctx_pol_name"].(string), paramMap["af"].(string))
			if err != nil {
				return diag.FromErr(err)
			}

		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationfvRsCtxToOspfCtxPolFromVRF(fvCtx.DistinguishedName, paramMap["tn_ospf_ctx_pol_name"].(string), paramMap["af"].(string))
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}
	if d.HasChange("relation_fv_rs_ctx_to_ep_ret") {
		_, newRelParam := d.GetChange("relation_fv_rs_ctx_to_ep_ret")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationfvRsCtxToEpRetFromVRF(fvCtx.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_fv_rs_bgp_ctx_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_bgp_ctx_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationfvRsBgpCtxPolFromVRF(fvCtx.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_fv_rs_ctx_mon_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_ctx_mon_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationfvRsCtxMonPolFromVRF(fvCtx.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationfvRsCtxMonPolFromVRF(fvCtx.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_fv_rs_ctx_to_ext_route_tag_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_ctx_to_ext_route_tag_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationfvRsCtxToExtRouteTagPolFromVRF(fvCtx.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_fv_rs_ctx_to_bgp_ctx_af_pol") {
		oldRel, newRel := d.GetChange("relation_fv_rs_ctx_to_bgp_ctx_af_pol")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.DeleteRelationfvRsCtxToBgpCtxAfPolFromVRF(fvCtx.DistinguishedName, paramMap["tn_bgp_ctx_af_pol_name"].(string), paramMap["af"].(string))
			if err != nil {
				return diag.FromErr(err)
			}

		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationfvRsCtxToBgpCtxAfPolFromVRF(fvCtx.DistinguishedName, paramMap["tn_bgp_ctx_af_pol_name"].(string), paramMap["af"].(string))
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}

	d.SetId(fvCtx.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciVRFRead(ctx, d, m)

}

func resourceAciVRFRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fvCtx, err := getRemoteVRF(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setVRFAttributes(fvCtx, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	fvRsOspfCtxPolData, err := aciClient.ReadRelationfvRsOspfCtxPolFromVRF(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsOspfCtxPol %v", err)
		d.Set("relation_fv_rs_ospf_ctx_pol", "")

	} else {
		if _, ok := d.GetOk("relation_fv_rs_ospf_ctx_pol"); ok {
			tfName := GetMOName(d.Get("relation_fv_rs_ospf_ctx_pol").(string))
			if tfName != fvRsOspfCtxPolData {
				d.Set("relation_fv_rs_ospf_ctx_pol", "")
			}
		}
	}

	fvRsVrfValidationPolData, err := aciClient.ReadRelationfvRsVrfValidationPolFromVRF(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsVrfValidationPol %v", err)
		d.Set("relation_fv_rs_vrf_validation_pol", "")

	} else {
		if _, ok := d.GetOk("relation_fv_rs_vrf_validation_pol"); ok {
			tfName := GetMOName(d.Get("relation_fv_rs_vrf_validation_pol").(string))
			if tfName != fvRsVrfValidationPolData {
				d.Set("relation_fv_rs_vrf_validation_pol", "")
			}
		}
	}

	fvRsCtxMcastToData, err := aciClient.ReadRelationfvRsCtxMcastToFromVRF(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsCtxMcastTo %v", err)
		d.Set("relation_fv_rs_ctx_mcast_to", fvRsCtxMcastToData)

	} else {
		d.Set("relation_fv_rs_ctx_mcast_to", fvRsCtxMcastToData)
	}

	fvRsCtxToEigrpCtxAfPolData, err := aciClient.ReadRelationfvRsCtxToEigrpCtxAfPolFromVRF(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsCtxToEigrpCtxAfPol %v", err)

	} else {
		d.Set("relation_fv_rs_ctx_to_eigrp_ctx_af_pol", fvRsCtxToEigrpCtxAfPolData)
	}

	fvRsCtxToOspfCtxPolData, err := aciClient.ReadRelationfvRsCtxToOspfCtxPolFromVRF(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsCtxToOspfCtxPol %v", err)

	} else {
		d.Set("relation_fv_rs_ctx_to_ospf_ctx_pol", fvRsCtxToOspfCtxPolData)
	}

	fvRsCtxToEpRetData, err := aciClient.ReadRelationfvRsCtxToEpRetFromVRF(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsCtxToEpRet %v", err)
		d.Set("relation_fv_rs_ctx_to_ep_ret", "")

	} else {
		if _, ok := d.GetOk("relation_fv_rs_ctx_to_ep_ret"); ok {
			tfName := GetMOName(d.Get("relation_fv_rs_ctx_to_ep_ret").(string))
			if tfName != fvRsCtxToEpRetData {
				d.Set("relation_fv_rs_ctx_to_ep_ret", "")
			}
		}
	}

	fvRsBgpCtxPolData, err := aciClient.ReadRelationfvRsBgpCtxPolFromVRF(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsBgpCtxPol %v", err)
		d.Set("relation_fv_rs_bgp_ctx_pol", "")

	} else {
		if _, ok := d.GetOk("relation_fv_rs_bgp_ctx_pol"); ok {
			tfName := GetMOName(d.Get("relation_fv_rs_bgp_ctx_pol").(string))
			if tfName != fvRsBgpCtxPolData {
				d.Set("relation_fv_rs_bgp_ctx_pol", "")
			}
		}
	}

	fvRsCtxMonPolData, err := aciClient.ReadRelationfvRsCtxMonPolFromVRF(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsCtxMonPol %v", err)
		d.Set("relation_fv_rs_ctx_mon_pol", "")

	} else {
		if _, ok := d.GetOk("relation_fv_rs_ctx_mon_pol"); ok {
			tfName := GetMOName(d.Get("relation_fv_rs_ctx_mon_pol").(string))
			if tfName != fvRsCtxMonPolData {
				d.Set("relation_fv_rs_ctx_mon_pol", "")
			}
		}
	}

	fvRsCtxToExtRouteTagPolData, err := aciClient.ReadRelationfvRsCtxToExtRouteTagPolFromVRF(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsCtxToExtRouteTagPol %v", err)
		d.Set("relation_fv_rs_ctx_to_ext_route_tag_pol", "")

	} else {
		if _, ok := d.GetOk("relation_fv_rs_ctx_to_ext_route_tag_pol"); ok {
			tfName := GetMOName(d.Get("relation_fv_rs_ctx_to_ext_route_tag_pol").(string))
			if tfName != fvRsCtxToExtRouteTagPolData {
				d.Set("relation_fv_rs_ctx_to_ext_route_tag_pol", "")
			}
		}
	}

	fvRsCtxToBgpCtxAfPolData, err := aciClient.ReadRelationfvRsCtxToBgpCtxAfPolFromVRF(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsCtxToBgpCtxAfPol %v", err)

	} else {
		d.Set("relation_fv_rs_ctx_to_bgp_ctx_af_pol", fvRsCtxToBgpCtxAfPolData)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciVRFDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fvCtx")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
