package aci

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciL3Outside() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciL3OutsideCreate,
		UpdateContext: resourceAciL3OutsideUpdate,
		ReadContext:   resourceAciL3OutsideRead,
		DeleteContext: resourceAciL3OutsideDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciL3OutsideImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
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

			"enforce_rtctrl": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"export",
						"import",
					}, false),
					Default: "export",
				},
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
			"mpls_enabled": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			// Relation to Route Control for Dampening
			"relation_l3ext_rs_dampening_pol": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tn_rtctrl_profile_name": {
							Type:       schema.TypeString,
							Optional:   true,
							Deprecated: "Use tn_rtctrl_profile_dn instead of tn_rtctrl_profile_name",
						},
						"tn_rtctrl_profile_dn": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"af": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ipv4-ucast",
								"ipv6-ucast",
							}, false),
							Default: "ipv4-ucast",
						},
					},
				},
			},
			// Target VRF object should belong to the parent tenant or be a shared object.
			"relation_l3ext_rs_ectx": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			// Relation to Route Profile for Interleak - L3 Out Context Interleak Policy object should belong to the parent tenant or be a shared object.
			"relation_l3ext_rs_interleak_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// Relation to L3 Domain
			"relation_l3ext_rs_l3_dom_att": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// Relation to Route Profile for Redistribution
			"relation_l3extrs_redistribute_pol": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Create relation to rtctrlProfile",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Optional: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"attached-host",
								"direct",
								"static",
							}, false),
						},
						"target_dn": {
							Required: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		})),
	}
}
func getRemoteL3Outside(client *client.Client, dn string) (*models.L3Outside, error) {
	l3extOutCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extOut := models.L3OutsideFromContainer(l3extOutCont)

	if l3extOut.DistinguishedName == "" {
		return nil, fmt.Errorf("L3 Outside %s not found", dn)
	}

	return l3extOut, nil
}

func setL3OutsideAttributes(l3extOut *models.L3Outside, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(l3extOut.DistinguishedName)
	d.Set("description", l3extOut.Description)

	if dn != l3extOut.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	l3extOutMap, err := l3extOut.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("tenant_dn", GetParentDn(dn, fmt.Sprintf("/%s", fmt.Sprintf(models.Rnl3extOut, l3extOutMap["name"]))))

	d.Set("name", l3extOutMap["name"])

	d.Set("annotation", l3extOutMap["annotation"])

	enforceRtctrlGet := make([]string, 0, 1)
	for _, val := range strings.Split(l3extOutMap["enforceRtctrl"], ",") {
		enforceRtctrlGet = append(enforceRtctrlGet, strings.Trim(val, " "))
	}
	sort.Strings(enforceRtctrlGet)
	if enforceRtctrlIntr, ok := d.GetOk("enforce_rtctrl"); ok {
		enforceRtctrlAct := make([]string, 0, 1)
		for _, val := range enforceRtctrlIntr.([]interface{}) {
			enforceRtctrlAct = append(enforceRtctrlAct, val.(string))
		}
		sort.Strings(enforceRtctrlAct)
		if reflect.DeepEqual(enforceRtctrlAct, enforceRtctrlGet) {
			d.Set("enforce_rtctrl", d.Get("enforce_rtctrl").([]interface{}))
		} else {
			d.Set("enforce_rtctrl", enforceRtctrlGet)
		}
	} else {
		d.Set("enforce_rtctrl", enforceRtctrlGet)
	}

	d.Set("name_alias", l3extOutMap["nameAlias"])
	d.Set("target_dscp", l3extOutMap["targetDscp"])
	d.Set("mpls_enabled", l3extOutMap["mplsEnabled"])
	return d, nil
}

func getAndSetReadRelationl3extRsDampeningPolFromL3Outside(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	l3extRsDampeningPolData, err := client.ReadRelationl3extRsDampeningPolFromL3Outside(dn)
	l3extRsDampeningPolList := make([]map[string]string, 0)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsDampeningPol %v", err)
		d.Set("relation_l3ext_rs_dampening_pol", l3extRsDampeningPolList)
		return nil, err
	} else {
		for _, obj := range l3extRsDampeningPolData.([]map[string]string) {
			l3extRsDampeningPolList = append(l3extRsDampeningPolList, map[string]string{
				"tn_rtctrl_profile_dn": obj["tDn"],
				"af":                   obj["af"],
			})
		}
		d.Set("relation_l3ext_rs_dampening_pol", l3extRsDampeningPolList)
		log.Printf("[DEBUG]: l3extRsDampeningPol: %s reading finished successfully", l3extRsDampeningPolList)
	}
	return d, nil
}

func getAndSetReadRelationl3extRsEctxFromL3Outside(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	l3extRsEctxData, err := client.ReadRelationl3extRsEctxFromL3Outside(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsEctx %v", err)
		d.Set("relation_l3ext_rs_ectx", nil)
		return nil, err
	} else {
		d.Set("relation_l3ext_rs_ectx", l3extRsEctxData.(string))
		log.Printf("[DEBUG]: l3extRsEctx: %s reading finished successfully", l3extRsEctxData.(string))
	}
	return d, nil
}

func getAndSetReadRelationl3extRsInterleakPolFromL3Outside(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	l3extRsInterleakPolData, err := client.ReadRelationl3extRsInterleakPolFromL3Outside(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsInterleakPol %v", err)
		d.Set("relation_l3ext_rs_interleak_pol", nil)
		return nil, err
	} else {
		d.Set("relation_l3ext_rs_interleak_pol", l3extRsInterleakPolData.(string))
		log.Printf("[DEBUG]: l3extRsInterleakPol: %s reading finished successfully", l3extRsInterleakPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationl3extRsL3DomAttFromL3Outside(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	l3extRsL3DomAttData, err := client.ReadRelationl3extRsL3DomAttFromL3Outside(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsL3DomAtt %v", err)
		d.Set("relation_l3ext_rs_l3_dom_att", nil)
		return nil, err
	} else {
		d.Set("relation_l3ext_rs_l3_dom_att", l3extRsL3DomAttData.(string))
		log.Printf("[DEBUG]: l3extRsL3DomAtt: %s reading finished successfully", l3extRsL3DomAttData.(string))

	}
	return d, nil
}

func getAndSetReadRelationl3extRsRedistributePol(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	l3extRsRedistributePolData, err := client.ReadRelationl3extRsRedistributePol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsRedistributePol %v", err)
		d.Set("relation_l3extrs_redistribute_pol", make([]interface{}, 0, 1))
	} else {
		l3extRsRedistributePolResultData := make([]map[string]string, 0, 1)
		for _, obj := range l3extRsRedistributePolData.([]map[string]string) {
			l3extRsRedistributePolResultData = append(l3extRsRedistributePolResultData, map[string]string{
				"source":    obj["src"],
				"target_dn": obj["tDn"],
			})
		}
		d.Set("relation_l3extrs_redistribute_pol", l3extRsRedistributePolResultData)
		log.Printf("[DEBUG]: l3extRsRedistributePol: Reading finished successfully")
	}
	return d, nil
}

func resourceAciL3OutsideImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	l3extOut, err := getRemoteL3Outside(aciClient, dn)

	if err != nil {
		return nil, err
	}
	l3extOutMap, err := l3extOut.ToMap()
	if err != nil {
		return nil, err
	}
	name := l3extOutMap["name"]
	pDN := GetParentDn(dn, fmt.Sprintf("/%s", fmt.Sprintf(models.Rnl3extOut, name)))
	d.Set("tenant_dn", pDN)
	schemaFilled, err := setL3OutsideAttributes(l3extOut, d)
	if err != nil {
		return nil, err
	}

	// Importing l3extRsDampeningPol object
	getAndSetReadRelationl3extRsDampeningPolFromL3Outside(aciClient, dn, d)

	// Importing l3extRsEctx object
	getAndSetReadRelationl3extRsEctxFromL3Outside(aciClient, dn, d)

	// Importing l3extRsInterleakPol object
	getAndSetReadRelationl3extRsInterleakPolFromL3Outside(aciClient, dn, d)

	// Importing l3extRsL3DomAtt object
	getAndSetReadRelationl3extRsL3DomAttFromL3Outside(aciClient, dn, d)

	// Importing l3extRsRedistributePol object
	getAndSetReadRelationl3extRsRedistributePol(aciClient, dn, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL3OutsideCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L3Outside: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	l3extOutAttr := models.L3OutsideAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extOutAttr.Annotation = Annotation.(string)
	} else {
		l3extOutAttr.Annotation = "{}"
	}

	if EnforceRtctrl, ok := d.GetOk("enforce_rtctrl"); ok {
		exportFlag := false
		enforceRtctrlList := make([]string, 0, 1)
		for _, val := range EnforceRtctrl.([]interface{}) {
			enforceRtctrlList = append(enforceRtctrlList, val.(string))
			if val.(string) == "export" {
				exportFlag = true
			}
		}
		if !exportFlag {
			enforceRtctrlList = append(enforceRtctrlList, "export")
		}
		EnforceRtctrl := strings.Join(enforceRtctrlList, ",")
		l3extOutAttr.EnforceRtctrl = EnforceRtctrl
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extOutAttr.NameAlias = NameAlias.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		l3extOutAttr.TargetDscp = TargetDscp.(string)
	}

	if MplsEnabled, ok := d.GetOk("mpls_enabled"); ok {
		l3extOutAttr.MplsEnabled = MplsEnabled.(string)
	}

	l3extOut := models.NewL3Outside(fmt.Sprintf(models.Rnl3extOut, name), TenantDn, desc, l3extOutAttr)

	err := aciClient.Save(l3extOut)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTol3extRsEctx, ok := d.GetOk("relation_l3ext_rs_ectx"); ok {
		relationParam := relationTol3extRsEctx.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTol3extRsInterleakPol, ok := d.GetOk("relation_l3ext_rs_interleak_pol"); ok {
		relationParam := relationTol3extRsInterleakPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTol3extRsL3DomAtt, ok := d.GetOk("relation_l3ext_rs_l3_dom_att"); ok {
		relationParam := relationTol3extRsL3DomAtt.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTol3extRsRedistributePol, ok := d.GetOk("relation_l3extrs_redistribute_pol"); ok {
		relationParamList := toStringList(relationTol3extRsRedistributePol.(*schema.Set).List())
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

	if relationTol3extRsDampeningPol, ok := d.GetOk("relation_l3ext_rs_dampening_pol"); ok {

		relationParamList := relationTol3extRsDampeningPol.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			l3extRsDampeningPolName := getTargetObjectName(paramMap, "tn_rtctrl_profile_dn", "tn_rtctrl_profile_name")
			err = aciClient.CreateRelationl3extRsDampeningPolFromL3Outside(l3extOut.DistinguishedName, l3extRsDampeningPolName, paramMap["af"].(string))
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}
	if relationTol3extRsEctx, ok := d.GetOk("relation_l3ext_rs_ectx"); ok {
		relationParam := relationTol3extRsEctx.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationl3extRsEctxFromL3Outside(l3extOut.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTol3extRsInterleakPol, ok := d.GetOk("relation_l3ext_rs_interleak_pol"); ok {
		relationParam := relationTol3extRsInterleakPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationl3extRsInterleakPolFromL3Outside(l3extOut.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTol3extRsL3DomAtt, ok := d.GetOk("relation_l3ext_rs_l3_dom_att"); ok {
		relationParam := relationTol3extRsL3DomAtt.(string)
		err = aciClient.CreateRelationl3extRsL3DomAttFromL3Outside(l3extOut.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationTol3extRsRedistributePol, ok := d.GetOk("relation_l3extrs_redistribute_pol"); ok {
		relationParamList := relationTol3extRsRedistributePol.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})

			err = aciClient.CreateRelationl3extRsRedistributePol(l3extOut.DistinguishedName, l3extOutAttr.Annotation, paramMap["source"].(string), GetMOName(paramMap["target_dn"].(string)))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(l3extOut.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciL3OutsideRead(ctx, d, m)
}

func resourceAciL3OutsideUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L3Outside: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	l3extOutAttr := models.L3OutsideAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extOutAttr.Annotation = Annotation.(string)
	} else {
		l3extOutAttr.Annotation = "{}"
	}

	if EnforceRtctrl, ok := d.GetOk("enforce_rtctrl"); ok {
		exportFlag := false
		enforceRtctrlList := make([]string, 0, 1)
		for _, val := range EnforceRtctrl.([]interface{}) {
			enforceRtctrlList = append(enforceRtctrlList, val.(string))
			if val.(string) == "export" {
				exportFlag = true
			}
		}
		if !exportFlag {
			enforceRtctrlList = append(enforceRtctrlList, "export")
		}
		EnforceRtctrl := strings.Join(enforceRtctrlList, ",")
		l3extOutAttr.EnforceRtctrl = EnforceRtctrl
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extOutAttr.NameAlias = NameAlias.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		l3extOutAttr.TargetDscp = TargetDscp.(string)
	}

	if MplsEnabled, ok := d.GetOk("mpls_enabled"); ok {
		l3extOutAttr.MplsEnabled = MplsEnabled.(string)
	}

	l3extOut := models.NewL3Outside(fmt.Sprintf(models.Rnl3extOut, name), TenantDn, desc, l3extOutAttr)

	l3extOut.Status = "modified"

	err := aciClient.Save(l3extOut)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_l3ext_rs_ectx") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_ectx")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_l3ext_rs_interleak_pol") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_interleak_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_l3ext_rs_l3_dom_att") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_l3_dom_att")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_l3extrs_redistribute_pol") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_l3extrs_redistribute_pol")
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

	if d.HasChange("relation_l3ext_rs_dampening_pol") {
		oldRel, newRel := d.GetChange("relation_l3ext_rs_dampening_pol")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})
			l3extRsDampeningPolName := getTargetObjectName(paramMap, "tn_rtctrl_profile_dn", "tn_rtctrl_profile_name")
			err = aciClient.DeleteRelationl3extRsDampeningPolFromL3Outside(l3extOut.DistinguishedName, l3extRsDampeningPolName, paramMap["af"].(string))
			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			l3extRsDampeningPolName := getTargetObjectName(paramMap, "tn_rtctrl_profile_dn", "tn_rtctrl_profile_name")
			err = aciClient.CreateRelationl3extRsDampeningPolFromL3Outside(l3extOut.DistinguishedName, l3extRsDampeningPolName, paramMap["af"].(string))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("relation_l3ext_rs_ectx") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_ectx")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationl3extRsEctxFromL3Outside(l3extOut.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_l3ext_rs_interleak_pol") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_interleak_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationl3extRsInterleakPolFromL3Outside(l3extOut.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationl3extRsInterleakPolFromL3Outside(l3extOut.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_l3ext_rs_l3_dom_att") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_l3_dom_att")
		err = aciClient.DeleteRelationl3extRsL3DomAttFromL3Outside(l3extOut.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationl3extRsL3DomAttFromL3Outside(l3extOut.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("relation_l3extrs_redistribute_pol") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_l3extrs_redistribute_pol")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})

			err = aciClient.DeleteRelationl3extRsRedistributePol(l3extOut.DistinguishedName, GetMOName(paramMap["target_dn"].(string)), paramMap["source"].(string))
			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})

			err = aciClient.CreateRelationl3extRsRedistributePol(l3extOut.DistinguishedName, l3extOutAttr.Annotation, paramMap["source"].(string), GetMOName(paramMap["target_dn"].(string)))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(l3extOut.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciL3OutsideRead(ctx, d, m)

}

func resourceAciL3OutsideRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l3extOut, err := getRemoteL3Outside(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setL3OutsideAttributes(l3extOut, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	// Importing l3extRsDampeningPol object
	getAndSetReadRelationl3extRsDampeningPolFromL3Outside(aciClient, dn, d)

	// Importing l3extRsEctx object
	getAndSetReadRelationl3extRsEctxFromL3Outside(aciClient, dn, d)

	// Importing l3extRsInterleakPol object
	getAndSetReadRelationl3extRsInterleakPolFromL3Outside(aciClient, dn, d)

	// Importing l3extRsL3DomAtt object
	getAndSetReadRelationl3extRsL3DomAttFromL3Outside(aciClient, dn, d)

	// Importing l3extRsRedistributePol object
	getAndSetReadRelationl3extRsRedistributePol(aciClient, dn, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciL3OutsideDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l3extOut")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
