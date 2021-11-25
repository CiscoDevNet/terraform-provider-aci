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

func resourceAciSpineAccessPortPolicyGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciSpineAccessPortPolicyGroupCreate,
		UpdateContext: resourceAciSpineAccessPortPolicyGroupUpdate,
		ReadContext:   resourceAciSpineAccessPortPolicyGroupRead,
		DeleteContext: resourceAciSpineAccessPortPolicyGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciSpineAccessPortPolicyGroupImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_infra_rs_h_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Default:  "uni/infra/hintfpol-default",
				Optional: true,
			},
			"relation_infra_rs_cdp_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Default:  "uni/infra/cdpIfP-default",
				Optional: true,
			},
			"relation_infra_rs_copp_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Default:  "uni/infra/coppifpol-default",
				Optional: true,
			},
			"relation_infra_rs_att_ent_p": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_macsec_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Default:  "uni/infra/macsecifp-default",
				Optional: true,
			},
		}),
	}
}
func getRemoteSpineAccessPortPolicyGroup(client *client.Client, dn string) (*models.SpineAccessPortPolicyGroup, error) {
	infraSpAccPortGrpCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraSpAccPortGrp := models.SpineAccessPortPolicyGroupFromContainer(infraSpAccPortGrpCont)

	if infraSpAccPortGrp.DistinguishedName == "" {
		return nil, fmt.Errorf("SpineAccessPortPolicyGroup %s not found", infraSpAccPortGrp.DistinguishedName)
	}

	return infraSpAccPortGrp, nil
}

func setSpineAccessPortPolicyGroupAttributes(infraSpAccPortGrp *models.SpineAccessPortPolicyGroup, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(infraSpAccPortGrp.DistinguishedName)
	d.Set("description", infraSpAccPortGrp.Description)
	infraSpAccPortGrpMap, err := infraSpAccPortGrp.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("name", infraSpAccPortGrpMap["name"])

	d.Set("annotation", infraSpAccPortGrpMap["annotation"])
	d.Set("name_alias", infraSpAccPortGrpMap["nameAlias"])
	return d, nil
}

func resourceAciSpineAccessPortPolicyGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraSpAccPortGrp, err := getRemoteSpineAccessPortPolicyGroup(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setSpineAccessPortPolicyGroupAttributes(infraSpAccPortGrp, d)

	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciSpineAccessPortPolicyGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SpineAccessPortPolicyGroup: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraSpAccPortGrpAttr := models.SpineAccessPortPolicyGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraSpAccPortGrpAttr.Annotation = Annotation.(string)
	} else {
		infraSpAccPortGrpAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraSpAccPortGrpAttr.NameAlias = NameAlias.(string)
	}
	infraSpAccPortGrp := models.NewSpineAccessPortPolicyGroup(fmt.Sprintf("infra/funcprof/spaccportgrp-%s", name), "uni", desc, infraSpAccPortGrpAttr)

	err := aciClient.Save(infraSpAccPortGrp)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationToinfraRsHIfPol, ok := d.GetOk("relation_infra_rs_h_if_pol"); ok {
		relationParam := relationToinfraRsHIfPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsCdpIfPol, ok := d.GetOk("relation_infra_rs_cdp_if_pol"); ok {
		relationParam := relationToinfraRsCdpIfPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsCoppIfPol, ok := d.GetOk("relation_infra_rs_copp_if_pol"); ok {
		relationParam := relationToinfraRsCoppIfPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsAttEntP, ok := d.GetOk("relation_infra_rs_att_ent_p"); ok {
		relationParam := relationToinfraRsAttEntP.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsMacsecIfPol, ok := d.GetOk("relation_infra_rs_macsec_if_pol"); ok {
		relationParam := relationToinfraRsMacsecIfPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationToinfraRsHIfPol, ok := d.GetOk("relation_infra_rs_h_if_pol"); ok {
		relationParam := relationToinfraRsHIfPol.(string)
		relationParamName := models.GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsHIfPolFromSpineAccessPortPolicyGroup(infraSpAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationToinfraRsCdpIfPol, ok := d.GetOk("relation_infra_rs_cdp_if_pol"); ok {
		relationParam := relationToinfraRsCdpIfPol.(string)
		relationParamName := models.GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsCdpIfPolFromSpineAccessPortPolicyGroup(infraSpAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationToinfraRsCoppIfPol, ok := d.GetOk("relation_infra_rs_copp_if_pol"); ok {
		relationParam := relationToinfraRsCoppIfPol.(string)
		relationParamName := models.GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsCoppIfPolFromSpineAccessPortPolicyGroup(infraSpAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationToinfraRsAttEntP, ok := d.GetOk("relation_infra_rs_att_ent_p"); ok {
		relationParam := relationToinfraRsAttEntP.(string)
		err = aciClient.CreateRelationinfraRsAttEntPFromSpineAccessPortPolicyGroup(infraSpAccPortGrp.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationToinfraRsMacsecIfPol, ok := d.GetOk("relation_infra_rs_macsec_if_pol"); ok {
		relationParam := relationToinfraRsMacsecIfPol.(string)
		relationParamName := models.GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsMacsecIfPolFromSpineAccessPortPolicyGroup(infraSpAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(infraSpAccPortGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciSpineAccessPortPolicyGroupRead(ctx, d, m)
}

func resourceAciSpineAccessPortPolicyGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SpineAccessPortPolicyGroup: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraSpAccPortGrpAttr := models.SpineAccessPortPolicyGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraSpAccPortGrpAttr.Annotation = Annotation.(string)
	} else {
		infraSpAccPortGrpAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraSpAccPortGrpAttr.NameAlias = NameAlias.(string)
	}
	infraSpAccPortGrp := models.NewSpineAccessPortPolicyGroup(fmt.Sprintf("infra/funcprof/spaccportgrp-%s", name), "uni", desc, infraSpAccPortGrpAttr)

	infraSpAccPortGrp.Status = "modified"

	err := aciClient.Save(infraSpAccPortGrp)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_infra_rs_h_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_h_if_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_cdp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_cdp_if_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_copp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_copp_if_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_att_ent_p") {
		_, newRelParam := d.GetChange("relation_infra_rs_att_ent_p")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_macsec_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_macsec_if_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_infra_rs_h_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_h_if_pol")
		newRelParamName := models.GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsHIfPolFromSpineAccessPortPolicyGroup(infraSpAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_infra_rs_cdp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_cdp_if_pol")
		newRelParamName := models.GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsCdpIfPolFromSpineAccessPortPolicyGroup(infraSpAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_infra_rs_copp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_copp_if_pol")
		newRelParamName := models.GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsCoppIfPolFromSpineAccessPortPolicyGroup(infraSpAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_infra_rs_att_ent_p") {
		_, newRelParam := d.GetChange("relation_infra_rs_att_ent_p")
		err = aciClient.DeleteRelationinfraRsAttEntPFromSpineAccessPortPolicyGroup(infraSpAccPortGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsAttEntPFromSpineAccessPortPolicyGroup(infraSpAccPortGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_infra_rs_macsec_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_macsec_if_pol")
		newRelParamName := models.GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsMacsecIfPolFromSpineAccessPortPolicyGroup(infraSpAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(infraSpAccPortGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciSpineAccessPortPolicyGroupRead(ctx, d, m)

}

func resourceAciSpineAccessPortPolicyGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraSpAccPortGrp, err := getRemoteSpineAccessPortPolicyGroup(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setSpineAccessPortPolicyGroupAttributes(infraSpAccPortGrp, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	infraRsHIfPolData, err := aciClient.ReadRelationinfraRsHIfPolFromSpineAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsHIfPol %v", err)
		d.Set("relation_infra_rs_h_if_pol", "")

	} else {
		d.Set("relation_infra_rs_h_if_pol", infraRsHIfPolData.(string))
	}

	infraRsCdpIfPolData, err := aciClient.ReadRelationinfraRsCdpIfPolFromSpineAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsCdpIfPol %v", err)
		d.Set("relation_infra_rs_cdp_if_pol", "")

	} else {
		d.Set("relation_infra_rs_cdp_if_pol", infraRsCdpIfPolData.(string))
	}

	infraRsCoppIfPolData, err := aciClient.ReadRelationinfraRsCoppIfPolFromSpineAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsCoppIfPol %v", err)
		d.Set("relation_infra_rs_copp_if_pol", "")

	} else {
		d.Set("relation_infra_rs_copp_if_pol", infraRsCoppIfPolData.(string))
	}

	infraRsAttEntPData, err := aciClient.ReadRelationinfraRsAttEntPFromSpineAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsAttEntP %v", err)
		d.Set("relation_infra_rs_att_ent_p", "")

	} else {
		d.Set("relation_infra_rs_att_ent_p", infraRsAttEntPData.(string))
	}

	infraRsMacsecIfPolData, err := aciClient.ReadRelationinfraRsMacsecIfPolFromSpineAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsMacsecIfPol %v", err)
		d.Set("relation_infra_rs_macsec_if_pol", "")

	} else {
		d.Set("relation_infra_rs_macsec_if_pol", infraRsMacsecIfPolData.(string))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciSpineAccessPortPolicyGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraSpAccPortGrp")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
