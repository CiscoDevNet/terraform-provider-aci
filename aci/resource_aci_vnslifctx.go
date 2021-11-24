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

func resourceAciLogicalInterfaceContext() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciLogicalInterfaceContextCreate,
		UpdateContext: resourceAciLogicalInterfaceContextUpdate,
		ReadContext:   resourceAciLogicalInterfaceContextRead,
		DeleteContext: resourceAciLogicalInterfaceContextDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciLogicalInterfaceContextImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"logical_device_context_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"conn_name_or_lbl": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"l3_dest": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"yes",
					"no",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"permit_log": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"yes",
					"no",
				}, false),
			},

			"relation_vns_rs_l_if_ctx_to_cust_qos_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_vns_rs_l_if_ctx_to_svc_e_pg_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_vns_rs_l_if_ctx_to_svc_redirect_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_vns_rs_l_if_ctx_to_l_if": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_vns_rs_l_if_ctx_to_out_def": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_vns_rs_l_if_ctx_to_inst_p": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_vns_rs_l_if_ctx_to_bd": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_vns_rs_l_if_ctx_to_out": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteLogicalInterfaceContext(client *client.Client, dn string) (*models.LogicalInterfaceContext, error) {
	vnsLIfCtxCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsLIfCtx := models.LogicalInterfaceContextFromContainer(vnsLIfCtxCont)

	if vnsLIfCtx.DistinguishedName == "" {
		return nil, fmt.Errorf("LogicalInterfaceContext %s not found", vnsLIfCtx.DistinguishedName)
	}

	return vnsLIfCtx, nil
}

func setLogicalInterfaceContextAttributes(vnsLIfCtx *models.LogicalInterfaceContext, d *schema.ResourceData) (*schema.ResourceData, error) {

	dn := d.Id()
	d.SetId(vnsLIfCtx.DistinguishedName)
	d.Set("description", vnsLIfCtx.Description)

	if dn != vnsLIfCtx.DistinguishedName {
		d.Set("logical_device_context_dn", "")
	}

	vnsLIfCtxMap, err := vnsLIfCtx.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("conn_name_or_lbl", vnsLIfCtxMap["connNameOrLbl"])

	d.Set("annotation", vnsLIfCtxMap["annotation"])
	d.Set("l3_dest", vnsLIfCtxMap["l3Dest"])
	d.Set("name_alias", vnsLIfCtxMap["nameAlias"])
	d.Set("permit_log", vnsLIfCtxMap["permitLog"])
	return d, nil
}

func resourceAciLogicalInterfaceContextImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	vnsLIfCtx, err := getRemoteLogicalInterfaceContext(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setLogicalInterfaceContextAttributes(vnsLIfCtx, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLogicalInterfaceContextCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LogicalInterfaceContext: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	connNameOrLbl := d.Get("conn_name_or_lbl").(string)

	LogicalDeviceContextDn := d.Get("logical_device_context_dn").(string)

	vnsLIfCtxAttr := models.LogicalInterfaceContextAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vnsLIfCtxAttr.Annotation = Annotation.(string)
	} else {
		vnsLIfCtxAttr.Annotation = "{}"
	}
	if ConnNameOrLbl, ok := d.GetOk("conn_name_or_lbl"); ok {
		vnsLIfCtxAttr.ConnNameOrLbl = ConnNameOrLbl.(string)
	}
	if L3Dest, ok := d.GetOk("l3_dest"); ok {
		vnsLIfCtxAttr.L3Dest = L3Dest.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vnsLIfCtxAttr.NameAlias = NameAlias.(string)
	}
	if PermitLog, ok := d.GetOk("permit_log"); ok {
		vnsLIfCtxAttr.PermitLog = PermitLog.(string)
	}
	vnsLIfCtx := models.NewLogicalInterfaceContext(fmt.Sprintf("lIfCtx-c-%s", connNameOrLbl), LogicalDeviceContextDn, desc, vnsLIfCtxAttr)

	err := aciClient.Save(vnsLIfCtx)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTovnsRsLIfCtxToCustQosPol, ok := d.GetOk("relation_vns_rs_l_if_ctx_to_cust_qos_pol"); ok {
		relationParam := relationTovnsRsLIfCtxToCustQosPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTovnsRsLIfCtxToSvcEPgPol, ok := d.GetOk("relation_vns_rs_l_if_ctx_to_svc_e_pg_pol"); ok {
		relationParam := relationTovnsRsLIfCtxToSvcEPgPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTovnsRsLIfCtxToSvcRedirectPol, ok := d.GetOk("relation_vns_rs_l_if_ctx_to_svc_redirect_pol"); ok {
		relationParam := relationTovnsRsLIfCtxToSvcRedirectPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTovnsRsLIfCtxToLIf, ok := d.GetOk("relation_vns_rs_l_if_ctx_to_l_if"); ok {
		relationParam := relationTovnsRsLIfCtxToLIf.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTovnsRsLIfCtxToOutDef, ok := d.GetOk("relation_vns_rs_l_if_ctx_to_out_def"); ok {
		relationParam := relationTovnsRsLIfCtxToOutDef.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTovnsRsLIfCtxToInstP, ok := d.GetOk("relation_vns_rs_l_if_ctx_to_inst_p"); ok {
		relationParam := relationTovnsRsLIfCtxToInstP.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTovnsRsLIfCtxToBD, ok := d.GetOk("relation_vns_rs_l_if_ctx_to_bd"); ok {
		relationParam := relationTovnsRsLIfCtxToBD.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTovnsRsLIfCtxToOut, ok := d.GetOk("relation_vns_rs_l_if_ctx_to_out"); ok {
		relationParam := relationTovnsRsLIfCtxToOut.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTovnsRsLIfCtxToCustQosPol, ok := d.GetOk("relation_vns_rs_l_if_ctx_to_cust_qos_pol"); ok {
		relationParam := relationTovnsRsLIfCtxToCustQosPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationvnsRsLIfCtxToCustQosPolFromLogicalInterfaceContext(vnsLIfCtx.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationTovnsRsLIfCtxToSvcEPgPol, ok := d.GetOk("relation_vns_rs_l_if_ctx_to_svc_e_pg_pol"); ok {
		relationParam := relationTovnsRsLIfCtxToSvcEPgPol.(string)
		err = aciClient.CreateRelationvnsRsLIfCtxToSvcEPgPolFromLogicalInterfaceContext(vnsLIfCtx.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTovnsRsLIfCtxToSvcRedirectPol, ok := d.GetOk("relation_vns_rs_l_if_ctx_to_svc_redirect_pol"); ok {
		relationParam := relationTovnsRsLIfCtxToSvcRedirectPol.(string)
		err = aciClient.CreateRelationvnsRsLIfCtxToSvcRedirectPolFromLogicalInterfaceContext(vnsLIfCtx.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTovnsRsLIfCtxToLIf, ok := d.GetOk("relation_vns_rs_l_if_ctx_to_l_if"); ok {
		relationParam := relationTovnsRsLIfCtxToLIf.(string)
		err = aciClient.CreateRelationvnsRsLIfCtxToLIfFromLogicalInterfaceContext(vnsLIfCtx.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTovnsRsLIfCtxToOutDef, ok := d.GetOk("relation_vns_rs_l_if_ctx_to_out_def"); ok {
		relationParam := relationTovnsRsLIfCtxToOutDef.(string)
		err = aciClient.CreateRelationvnsRsLIfCtxToOutDefFromLogicalInterfaceContext(vnsLIfCtx.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationTovnsRsLIfCtxToInstP, ok := d.GetOk("relation_vns_rs_l_if_ctx_to_inst_p"); ok {
		relationParam := relationTovnsRsLIfCtxToInstP.(string)
		err = aciClient.CreateRelationvnsRsLIfCtxToInstPFromLogicalInterfaceContext(vnsLIfCtx.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationTovnsRsLIfCtxToBD, ok := d.GetOk("relation_vns_rs_l_if_ctx_to_bd"); ok {
		relationParam := relationTovnsRsLIfCtxToBD.(string)
		err = aciClient.CreateRelationvnsRsLIfCtxToBDFromLogicalInterfaceContext(vnsLIfCtx.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if relationTovnsRsLIfCtxToOut, ok := d.GetOk("relation_vns_rs_l_if_ctx_to_out"); ok {
		relationParam := relationTovnsRsLIfCtxToOut.(string)
		err = aciClient.CreateRelationvnsRsLIfCtxToOutFromLogicalInterfaceContext(vnsLIfCtx.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(vnsLIfCtx.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciLogicalInterfaceContextRead(ctx, d, m)
}

func resourceAciLogicalInterfaceContextUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LogicalInterfaceContext: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	connNameOrLbl := d.Get("conn_name_or_lbl").(string)

	LogicalDeviceContextDn := d.Get("logical_device_context_dn").(string)

	vnsLIfCtxAttr := models.LogicalInterfaceContextAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vnsLIfCtxAttr.Annotation = Annotation.(string)
	} else {
		vnsLIfCtxAttr.Annotation = "{}"
	}
	if ConnNameOrLbl, ok := d.GetOk("conn_name_or_lbl"); ok {
		vnsLIfCtxAttr.ConnNameOrLbl = ConnNameOrLbl.(string)
	}
	if L3Dest, ok := d.GetOk("l3_dest"); ok {
		vnsLIfCtxAttr.L3Dest = L3Dest.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vnsLIfCtxAttr.NameAlias = NameAlias.(string)
	}
	if PermitLog, ok := d.GetOk("permit_log"); ok {
		vnsLIfCtxAttr.PermitLog = PermitLog.(string)
	}
	vnsLIfCtx := models.NewLogicalInterfaceContext(fmt.Sprintf("lIfCtx-c-%s", connNameOrLbl), LogicalDeviceContextDn, desc, vnsLIfCtxAttr)

	vnsLIfCtx.Status = "modified"

	err := aciClient.Save(vnsLIfCtx)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_vns_rs_l_if_ctx_to_cust_qos_pol") {
		_, newRelParam := d.GetChange("relation_vns_rs_l_if_ctx_to_cust_qos_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_vns_rs_l_if_ctx_to_svc_e_pg_pol") {
		_, newRelParam := d.GetChange("relation_vns_rs_l_if_ctx_to_svc_e_pg_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_vns_rs_l_if_ctx_to_svc_redirect_pol") {
		_, newRelParam := d.GetChange("relation_vns_rs_l_if_ctx_to_svc_redirect_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_vns_rs_l_if_ctx_to_l_if") {
		_, newRelParam := d.GetChange("relation_vns_rs_l_if_ctx_to_l_if")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_vns_rs_l_if_ctx_to_out_def") {
		_, newRelParam := d.GetChange("relation_vns_rs_l_if_ctx_to_out_def")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_vns_rs_l_if_ctx_to_inst_p") {
		_, newRelParam := d.GetChange("relation_vns_rs_l_if_ctx_to_inst_p")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_vns_rs_l_if_ctx_to_bd") {
		_, newRelParam := d.GetChange("relation_vns_rs_l_if_ctx_to_bd")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_vns_rs_l_if_ctx_to_out") {
		_, newRelParam := d.GetChange("relation_vns_rs_l_if_ctx_to_out")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_vns_rs_l_if_ctx_to_cust_qos_pol") {
		_, newRelParam := d.GetChange("relation_vns_rs_l_if_ctx_to_cust_qos_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationvnsRsLIfCtxToCustQosPolFromLogicalInterfaceContext(vnsLIfCtx.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_vns_rs_l_if_ctx_to_svc_e_pg_pol") {
		_, newRelParam := d.GetChange("relation_vns_rs_l_if_ctx_to_svc_e_pg_pol")
		err = aciClient.DeleteRelationvnsRsLIfCtxToSvcEPgPolFromLogicalInterfaceContext(vnsLIfCtx.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvnsRsLIfCtxToSvcEPgPolFromLogicalInterfaceContext(vnsLIfCtx.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_vns_rs_l_if_ctx_to_svc_redirect_pol") {
		_, newRelParam := d.GetChange("relation_vns_rs_l_if_ctx_to_svc_redirect_pol")
		err = aciClient.DeleteRelationvnsRsLIfCtxToSvcRedirectPolFromLogicalInterfaceContext(vnsLIfCtx.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvnsRsLIfCtxToSvcRedirectPolFromLogicalInterfaceContext(vnsLIfCtx.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_vns_rs_l_if_ctx_to_l_if") {
		_, newRelParam := d.GetChange("relation_vns_rs_l_if_ctx_to_l_if")
		err = aciClient.DeleteRelationvnsRsLIfCtxToLIfFromLogicalInterfaceContext(vnsLIfCtx.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvnsRsLIfCtxToLIfFromLogicalInterfaceContext(vnsLIfCtx.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_vns_rs_l_if_ctx_to_out_def") {
		_, newRelParam := d.GetChange("relation_vns_rs_l_if_ctx_to_out_def")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationvnsRsLIfCtxToOutDefFromLogicalInterfaceContext(vnsLIfCtx.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_vns_rs_l_if_ctx_to_inst_p") {
		_, newRelParam := d.GetChange("relation_vns_rs_l_if_ctx_to_inst_p")
		err = aciClient.DeleteRelationvnsRsLIfCtxToInstPFromLogicalInterfaceContext(vnsLIfCtx.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvnsRsLIfCtxToInstPFromLogicalInterfaceContext(vnsLIfCtx.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_vns_rs_l_if_ctx_to_bd") {
		_, newRelParam := d.GetChange("relation_vns_rs_l_if_ctx_to_bd")
		err = aciClient.DeleteRelationvnsRsLIfCtxToBDFromLogicalInterfaceContext(vnsLIfCtx.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvnsRsLIfCtxToBDFromLogicalInterfaceContext(vnsLIfCtx.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_vns_rs_l_if_ctx_to_out") {
		_, newRelParam := d.GetChange("relation_vns_rs_l_if_ctx_to_out")
		err = aciClient.DeleteRelationvnsRsLIfCtxToOutFromLogicalInterfaceContext(vnsLIfCtx.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvnsRsLIfCtxToOutFromLogicalInterfaceContext(vnsLIfCtx.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(vnsLIfCtx.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciLogicalInterfaceContextRead(ctx, d, m)

}

func resourceAciLogicalInterfaceContextRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	vnsLIfCtx, err := getRemoteLogicalInterfaceContext(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setLogicalInterfaceContextAttributes(vnsLIfCtx, d)

	if err != nil {
		d.SetId("")
		return nil
	}

	vnsRsLIfCtxToCustQosPolData, err := aciClient.ReadRelationvnsRsLIfCtxToCustQosPolFromLogicalInterfaceContext(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsLIfCtxToCustQosPol %v", err)
		d.Set("relation_vns_rs_l_if_ctx_to_cust_qos_pol", "")

	} else {
		d.Set("relation_vns_rs_l_if_ctx_to_cust_qos_pol", vnsRsLIfCtxToCustQosPolData.(string))
	}

	vnsRsLIfCtxToSvcEPgPolData, err := aciClient.ReadRelationvnsRsLIfCtxToSvcEPgPolFromLogicalInterfaceContext(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsLIfCtxToSvcEPgPol %v", err)
		d.Set("relation_vns_rs_l_if_ctx_to_svc_e_pg_pol", "")

	} else {
		d.Set("relation_vns_rs_l_if_ctx_to_svc_e_pg_pol", vnsRsLIfCtxToSvcEPgPolData.(string))
	}

	vnsRsLIfCtxToSvcRedirectPolData, err := aciClient.ReadRelationvnsRsLIfCtxToSvcRedirectPolFromLogicalInterfaceContext(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsLIfCtxToSvcRedirectPol %v", err)
		d.Set("relation_vns_rs_l_if_ctx_to_svc_redirect_pol", "")

	} else {
		d.Set("relation_vns_rs_l_if_ctx_to_svc_redirect_pol", vnsRsLIfCtxToSvcRedirectPolData.(string))
	}

	vnsRsLIfCtxToLIfData, err := aciClient.ReadRelationvnsRsLIfCtxToLIfFromLogicalInterfaceContext(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsLIfCtxToLIf %v", err)
		d.Set("relation_vns_rs_l_if_ctx_to_l_if", "")

	} else {
		d.Set("relation_vns_rs_l_if_ctx_to_l_if", vnsRsLIfCtxToLIfData.(string))
	}

	vnsRsLIfCtxToOutDefData, err := aciClient.ReadRelationvnsRsLIfCtxToOutDefFromLogicalInterfaceContext(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsLIfCtxToOutDef %v", err)
		d.Set("relation_vns_rs_l_if_ctx_to_out_def", "")

	} else {
		d.Set("relation_vns_rs_l_if_ctx_to_out_def", vnsRsLIfCtxToOutDefData.(string))
	}

	vnsRsLIfCtxToInstPData, err := aciClient.ReadRelationvnsRsLIfCtxToInstPFromLogicalInterfaceContext(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsLIfCtxToInstP %v", err)
		d.Set("relation_vns_rs_l_if_ctx_to_inst_p", "")

	} else {
		d.Set("relation_vns_rs_l_if_ctx_to_inst_p", vnsRsLIfCtxToInstPData.(string))
	}

	vnsRsLIfCtxToBDData, err := aciClient.ReadRelationvnsRsLIfCtxToBDFromLogicalInterfaceContext(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsLIfCtxToBD %v", err)
		d.Set("relation_vns_rs_l_if_ctx_to_bd", "")
	} else {
		d.Set("relation_vns_rs_l_if_ctx_to_bd", vnsRsLIfCtxToBDData.(string))
	}

	vnsRsLIfCtxToOutData, err := aciClient.ReadRelationvnsRsLIfCtxToOutFromLogicalInterfaceContext(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsLIfCtxToOut %v", err)
		d.Set("relation_vns_rs_l_if_ctx_to_out", "")

	} else {
		d.Set("relation_vns_rs_l_if_ctx_to_out", vnsRsLIfCtxToOutData.(string))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciLogicalInterfaceContextDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vnsLIfCtx")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
