package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciLogicalDeviceContext() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciLogicalDeviceContextCreate,
		UpdateContext: resourceAciLogicalDeviceContextUpdate,
		ReadContext:   resourceAciLogicalDeviceContextRead,
		DeleteContext: resourceAciLogicalDeviceContextDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciLogicalDeviceContextImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"tenant_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"ctrct_name_or_lbl": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"graph_name_or_lbl": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"node_name_or_lbl": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"context": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_vns_rs_l_dev_ctx_to_l_dev": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_vns_rs_l_dev_ctx_to_rtr_cfg": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteLogicalDeviceContext(client *client.Client, dn string) (*models.LogicalDeviceContext, error) {
	vnsLDevCtxCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsLDevCtx := models.LogicalDeviceContextFromContainer(vnsLDevCtxCont)

	if vnsLDevCtx.DistinguishedName == "" {
		return nil, fmt.Errorf("LogicalDeviceContext %s not found", vnsLDevCtx.DistinguishedName)
	}

	return vnsLDevCtx, nil
}

func setLogicalDeviceContextAttributes(vnsLDevCtx *models.LogicalDeviceContext, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()

	d.SetId(vnsLDevCtx.DistinguishedName)
	d.Set("description", vnsLDevCtx.Description)

	if dn != vnsLDevCtx.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	vnsLDevCtxMap, err := vnsLDevCtx.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("ctrct_name_or_lbl", vnsLDevCtxMap["ctrctNameOrLbl"])
	d.Set("graph_name_or_lbl", vnsLDevCtxMap["graphNameOrLbl"])
	d.Set("node_name_or_lbl", vnsLDevCtxMap["nodeNameOrLbl"])
	d.Set("annotation", vnsLDevCtxMap["annotation"])
	d.Set("context", vnsLDevCtxMap["context"])
	d.Set("name_alias", vnsLDevCtxMap["nameAlias"])

	return d, nil
}

func resourceAciLogicalDeviceContextImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	vnsLDevCtx, err := getRemoteLogicalDeviceContext(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setLogicalDeviceContextAttributes(vnsLDevCtx, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLogicalDeviceContextCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LogicalDeviceContext: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	ctrctNameOrLbl := d.Get("ctrct_name_or_lbl").(string)

	graphNameOrLbl := d.Get("graph_name_or_lbl").(string)

	nodeNameOrLbl := d.Get("node_name_or_lbl").(string)

	TenantDn := d.Get("tenant_dn").(string)

	vnsLDevCtxAttr := models.LogicalDeviceContextAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vnsLDevCtxAttr.Annotation = Annotation.(string)
	} else {
		vnsLDevCtxAttr.Annotation = "{}"
	}
	if Context, ok := d.GetOk("context"); ok {
		vnsLDevCtxAttr.Context = Context.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vnsLDevCtxAttr.NameAlias = NameAlias.(string)
	}

	vnsLDevCtx := models.NewLogicalDeviceContext(fmt.Sprintf("ldevCtx-c-%s-g-%s-n-%s", ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl), TenantDn, desc, vnsLDevCtxAttr)

	err := aciClient.Save(vnsLDevCtx)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTovnsRsLDevCtxToLDev, ok := d.GetOk("relation_vns_rs_l_dev_ctx_to_l_dev"); ok {
		relationParam := relationTovnsRsLDevCtxToLDev.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTovnsRsLDevCtxToRtrCfg, ok := d.GetOk("relation_vns_rs_l_dev_ctx_to_rtr_cfg"); ok {
		relationParam := relationTovnsRsLDevCtxToRtrCfg.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTovnsRsLDevCtxToLDev, ok := d.GetOk("relation_vns_rs_l_dev_ctx_to_l_dev"); ok {
		relationParam := relationTovnsRsLDevCtxToLDev.(string)
		err = aciClient.CreateRelationvnsRsLDevCtxToLDevFromLogicalDeviceContext(vnsLDevCtx.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTovnsRsLDevCtxToRtrCfg, ok := d.GetOk("relation_vns_rs_l_dev_ctx_to_rtr_cfg"); ok {
		relationParam := relationTovnsRsLDevCtxToRtrCfg.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationvnsRsLDevCtxToRtrCfgFromLogicalDeviceContext(vnsLDevCtx.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(vnsLDevCtx.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciLogicalDeviceContextRead(ctx, d, m)
}

func resourceAciLogicalDeviceContextUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LogicalDeviceContext: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	ctrctNameOrLbl := d.Get("ctrct_name_or_lbl").(string)

	graphNameOrLbl := d.Get("graph_name_or_lbl").(string)

	nodeNameOrLbl := d.Get("node_name_or_lbl").(string)

	TenantDn := d.Get("tenant_dn").(string)

	vnsLDevCtxAttr := models.LogicalDeviceContextAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vnsLDevCtxAttr.Annotation = Annotation.(string)
	} else {
		vnsLDevCtxAttr.Annotation = "{}"
	}
	if Context, ok := d.GetOk("context"); ok {
		vnsLDevCtxAttr.Context = Context.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vnsLDevCtxAttr.NameAlias = NameAlias.(string)
	}

	vnsLDevCtx := models.NewLogicalDeviceContext(fmt.Sprintf("ldevCtx-c-%s-g-%s-n-%s", ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl), TenantDn, desc, vnsLDevCtxAttr)

	vnsLDevCtx.Status = "modified"

	err := aciClient.Save(vnsLDevCtx)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_vns_rs_l_dev_ctx_to_l_dev") {
		_, newRelParam := d.GetChange("relation_vns_rs_l_dev_ctx_to_l_dev")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_vns_rs_l_dev_ctx_to_rtr_cfg") {
		_, newRelParam := d.GetChange("relation_vns_rs_l_dev_ctx_to_rtr_cfg")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_vns_rs_l_dev_ctx_to_l_dev") {
		_, newRelParam := d.GetChange("relation_vns_rs_l_dev_ctx_to_l_dev")
		err = aciClient.DeleteRelationvnsRsLDevCtxToLDevFromLogicalDeviceContext(vnsLDevCtx.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvnsRsLDevCtxToLDevFromLogicalDeviceContext(vnsLDevCtx.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_vns_rs_l_dev_ctx_to_rtr_cfg") {
		_, newRelParam := d.GetChange("relation_vns_rs_l_dev_ctx_to_rtr_cfg")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationvnsRsLDevCtxToRtrCfgFromLogicalDeviceContext(vnsLDevCtx.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvnsRsLDevCtxToRtrCfgFromLogicalDeviceContext(vnsLDevCtx.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(vnsLDevCtx.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciLogicalDeviceContextRead(ctx, d, m)

}

func resourceAciLogicalDeviceContextRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	vnsLDevCtx, err := getRemoteLogicalDeviceContext(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setLogicalDeviceContextAttributes(vnsLDevCtx, d)
	if err != nil {
		d.SetId("")
		return nil
	}
	vnsRsLDevCtxToLDevData, err := aciClient.ReadRelationvnsRsLDevCtxToLDevFromLogicalDeviceContext(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsLDevCtxToLDev %v", err)
		d.Set("relation_vns_rs_l_dev_ctx_to_l_dev", "")

	} else {
		d.Set("relation_vns_rs_l_dev_ctx_to_l_dev", vnsRsLDevCtxToLDevData.(string))
	}

	vnsRsLDevCtxToRtrCfgData, err := aciClient.ReadRelationvnsRsLDevCtxToRtrCfgFromLogicalDeviceContext(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsLDevCtxToRtrCfg %v", err)
		d.Set("relation_vns_rs_l_dev_ctx_to_rtr_cfg", "")

	} else {
		d.Set("relation_vns_rs_l_dev_ctx_to_rtr_cfg", vnsRsLDevCtxToRtrCfgData.(string))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciLogicalDeviceContextDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vnsLDevCtx")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
