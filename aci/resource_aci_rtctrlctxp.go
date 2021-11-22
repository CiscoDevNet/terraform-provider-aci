package aci

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciRouteControlContext() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciRouteControlContextCreate,
		UpdateContext: resourceAciRouteControlContextUpdate,
		ReadContext:   resourceAciRouteControlContextRead,
		DeleteContext: resourceAciRouteControlContextDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciRouteControlContextImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"route_control_profile_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"deny",
					"permit",
				}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"order": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_rtctrl_rs_ctx_p_to_subj_p": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Create relation to rtctrl:SubjP",
				Set:         schema.HashString,
			}})),
	}
}

func getRemoteRouteControlContext(client *client.Client, dn string) (*models.RouteControlContext, error) {
	rtctrlCtxPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	rtctrlCtxP := models.RouteControlContextFromContainer(rtctrlCtxPCont)
	if rtctrlCtxP.DistinguishedName == "" {
		return nil, fmt.Errorf("RouteControlContext %s not found", rtctrlCtxP.DistinguishedName)
	}
	return rtctrlCtxP, nil
}

func setRouteControlContextAttributes(rtctrlCtxP *models.RouteControlContext, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(rtctrlCtxP.DistinguishedName)
	d.Set("description", rtctrlCtxP.Description)
	rtctrlCtxPMap, err := rtctrlCtxP.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("action", rtctrlCtxPMap["action"])
	d.Set("annotation", rtctrlCtxPMap["annotation"])
	d.Set("name", rtctrlCtxPMap["name"])
	d.Set("order", rtctrlCtxPMap["order"])
	d.Set("name_alias", rtctrlCtxPMap["nameAlias"])
	return d, nil
}

func resourceAciRouteControlContextImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	rtctrlCtxP, err := getRemoteRouteControlContext(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setRouteControlContextAttributes(rtctrlCtxP, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciRouteControlContextCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] RouteControlContext: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	RouteControlProfileDn := d.Get("route_control_profile_dn").(string)

	rtctrlCtxPAttr := models.RouteControlContextAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		rtctrlCtxPAttr.Annotation = Annotation.(string)
	} else {
		rtctrlCtxPAttr.Annotation = "{}"
	}

	if Action, ok := d.GetOk("action"); ok {
		rtctrlCtxPAttr.Action = Action.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		rtctrlCtxPAttr.Name = Name.(string)
	}

	if Order, ok := d.GetOk("order"); ok {
		rtctrlCtxPAttr.Order = Order.(string)
	}
	rtctrlCtxP := models.NewRouteControlContext(fmt.Sprintf("ctx-%s", name), RouteControlProfileDn, desc, nameAlias, rtctrlCtxPAttr)

	err := aciClient.Save(rtctrlCtxP)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTortctrlRsCtxPToSubjP, ok := d.GetOk("relation_rtctrl_rs_ctx_p_to_subj_p"); ok {
		relationParamList := toStringList(relationTortctrlRsCtxPToSubjP.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}

	if relationTortctrlRsCtxPToSubjP, ok := d.GetOk("relation_rtctrl_rs_ctx_p_to_subj_p"); ok {
		relationParamList := toStringList(relationTortctrlRsCtxPToSubjP.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationrtctrlRsCtxPToSubjP(rtctrlCtxP.DistinguishedName, rtctrlCtxPAttr.Annotation, GetMOName(relationParam))
			if err != nil {
				return diag.FromErr(err)
			}

		}
	}

	d.SetId(rtctrlCtxP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciRouteControlContextRead(ctx, d, m)
}

func resourceAciRouteControlContextUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] RouteControlContext: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	RouteControlProfileDn := d.Get("route_control_profile_dn").(string)
	rtctrlCtxPAttr := models.RouteControlContextAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		rtctrlCtxPAttr.Annotation = Annotation.(string)
	} else {
		rtctrlCtxPAttr.Annotation = "{}"
	}

	if Action, ok := d.GetOk("action"); ok {
		rtctrlCtxPAttr.Action = Action.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		rtctrlCtxPAttr.Name = Name.(string)
	}

	if Order, ok := d.GetOk("order"); ok {
		rtctrlCtxPAttr.Order = Order.(string)
	}
	rtctrlCtxP := models.NewRouteControlContext(fmt.Sprintf("ctx-%s", name), RouteControlProfileDn, desc, nameAlias, rtctrlCtxPAttr)

	rtctrlCtxP.Status = "modified"
	err := aciClient.Save(rtctrlCtxP)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_rtctrl_rs_ctx_p_to_subj_p") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_rtctrl_rs_ctx_p_to_subj_p")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())
		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("relation_rtctrl_rs_ctx_p_to_subj_p") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_rtctrl_rs_ctx_p_to_subj_p")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationrtctrlRsCtxPToSubjP(rtctrlCtxP.DistinguishedName, GetMOName(relDn))
			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationrtctrlRsCtxPToSubjP(rtctrlCtxP.DistinguishedName, rtctrlCtxPAttr.Annotation, GetMOName(relDn))
			if err != nil {
				return diag.FromErr(err)
			}

		}
	}

	d.SetId(rtctrlCtxP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciRouteControlContextRead(ctx, d, m)
}

func resourceAciRouteControlContextRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	rtctrlCtxP, err := getRemoteRouteControlContext(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	_, err = setRouteControlContextAttributes(rtctrlCtxP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	rtctrlRsCtxPToSubjPData, err := aciClient.ReadRelationrtctrlRsCtxPToSubjP(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation rtctrlRsCtxPToSubjP %v", err)
		d.Set("relation_rtctrl_rs_ctx_p_to_subj_p", make([]string, 0, 1))
	} else {
		if _, ok := d.GetOk("relation_rtctrl_rs_ctx_p_to_subj_p"); ok {
			relationParamList := toStringList(d.Get("relation_rtctrl_rs_ctx_p_to_subj_p").(*schema.Set).List())
			tfList := make([]string, 0, 1)
			for _, relationParam := range relationParamList {
				relationParamName := GetMOName(relationParam)
				tfList = append(tfList, relationParamName)
			}
			rtctrlRsCtxPToSubjPDataList := toStringList(rtctrlRsCtxPToSubjPData.(*schema.Set).List())
			sort.Strings(tfList)
			sort.Strings(rtctrlRsCtxPToSubjPDataList)
			if !reflect.DeepEqual(tfList, rtctrlRsCtxPToSubjPDataList) {
				d.Set("relation_rtctrl_rs_ctx_p_to_subj_p", make([]string, 0, 1))
			}
		}
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciRouteControlContextDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "rtctrlCtxP")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
