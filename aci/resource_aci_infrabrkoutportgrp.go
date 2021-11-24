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

func resourceAciLeafBreakoutPortGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciLeafBreakoutPortGroupCreate,
		UpdateContext: resourceAciLeafBreakoutPortGroupUpdate,
		ReadContext:   resourceAciLeafBreakoutPortGroupRead,
		DeleteContext: resourceAciLeafBreakoutPortGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciLeafBreakoutPortGroupImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"brkout_map": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"100g-2x",
					"100g-4x",
					"10g-4x",
					"25g-4x",
					"50g-8x",
					"none",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_infra_rs_mon_brkout_infra_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		}),
	}
}
func getRemoteLeafBreakoutPortGroup(client *client.Client, dn string) (*models.LeafBreakoutPortGroup, error) {
	infraBrkoutPortGrpCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraBrkoutPortGrp := models.LeafBreakoutPortGroupFromContainer(infraBrkoutPortGrpCont)

	if infraBrkoutPortGrp.DistinguishedName == "" {
		return nil, fmt.Errorf("LeafBreakoutPortGroup %s not found", infraBrkoutPortGrp.DistinguishedName)
	}

	return infraBrkoutPortGrp, nil
}

func setLeafBreakoutPortGroupAttributes(infraBrkoutPortGrp *models.LeafBreakoutPortGroup, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(infraBrkoutPortGrp.DistinguishedName)
	d.Set("description", infraBrkoutPortGrp.Description)
	infraBrkoutPortGrpMap, err := infraBrkoutPortGrp.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("name", infraBrkoutPortGrpMap["name"])

	d.Set("annotation", infraBrkoutPortGrpMap["annotation"])
	d.Set("brkout_map", infraBrkoutPortGrpMap["brkoutMap"])
	d.Set("name_alias", infraBrkoutPortGrpMap["nameAlias"])
	return d, nil
}

func resourceAciLeafBreakoutPortGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraBrkoutPortGrp, err := getRemoteLeafBreakoutPortGroup(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setLeafBreakoutPortGroupAttributes(infraBrkoutPortGrp, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLeafBreakoutPortGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LeafBreakoutPortGroup: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraBrkoutPortGrpAttr := models.LeafBreakoutPortGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraBrkoutPortGrpAttr.Annotation = Annotation.(string)
	} else {
		infraBrkoutPortGrpAttr.Annotation = "{}"
	}
	if BrkoutMap, ok := d.GetOk("brkout_map"); ok {
		infraBrkoutPortGrpAttr.BrkoutMap = BrkoutMap.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraBrkoutPortGrpAttr.NameAlias = NameAlias.(string)
	}
	infraBrkoutPortGrp := models.NewLeafBreakoutPortGroup(fmt.Sprintf("infra/funcprof/brkoutportgrp-%s", name), "uni", desc, infraBrkoutPortGrpAttr)

	err := aciClient.Save(infraBrkoutPortGrp)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationToinfraRsMonBrkoutInfraPol, ok := d.GetOk("relation_infra_rs_mon_brkout_infra_pol"); ok {
		relationParam := relationToinfraRsMonBrkoutInfraPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationToinfraRsMonBrkoutInfraPol, ok := d.GetOk("relation_infra_rs_mon_brkout_infra_pol"); ok {
		relationParam := relationToinfraRsMonBrkoutInfraPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsMonBrkoutInfraPolFromLeafBreakoutPortGroup(infraBrkoutPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(infraBrkoutPortGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciLeafBreakoutPortGroupRead(ctx, d, m)
}

func resourceAciLeafBreakoutPortGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LeafBreakoutPortGroup: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraBrkoutPortGrpAttr := models.LeafBreakoutPortGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraBrkoutPortGrpAttr.Annotation = Annotation.(string)
	} else {
		infraBrkoutPortGrpAttr.Annotation = "{}"
	}
	if BrkoutMap, ok := d.GetOk("brkout_map"); ok {
		infraBrkoutPortGrpAttr.BrkoutMap = BrkoutMap.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraBrkoutPortGrpAttr.NameAlias = NameAlias.(string)
	}
	infraBrkoutPortGrp := models.NewLeafBreakoutPortGroup(fmt.Sprintf("infra/funcprof/brkoutportgrp-%s", name), "uni", desc, infraBrkoutPortGrpAttr)

	infraBrkoutPortGrp.Status = "modified"

	err := aciClient.Save(infraBrkoutPortGrp)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_infra_rs_mon_brkout_infra_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_mon_brkout_infra_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_infra_rs_mon_brkout_infra_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_mon_brkout_infra_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsMonBrkoutInfraPolFromLeafBreakoutPortGroup(infraBrkoutPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(infraBrkoutPortGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciLeafBreakoutPortGroupRead(ctx, d, m)

}

func resourceAciLeafBreakoutPortGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraBrkoutPortGrp, err := getRemoteLeafBreakoutPortGroup(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setLeafBreakoutPortGroupAttributes(infraBrkoutPortGrp, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	infraRsMonBrkoutInfraPolData, err := aciClient.ReadRelationinfraRsMonBrkoutInfraPolFromLeafBreakoutPortGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsMonBrkoutInfraPol %v", err)
		d.Set("relation_infra_rs_mon_brkout_infra_pol", "")

	} else {
		d.Set("relation_infra_rs_mon_brkout_infra_pol", infraRsMonBrkoutInfraPolData.(string))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciLeafBreakoutPortGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraBrkoutPortGrp")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
