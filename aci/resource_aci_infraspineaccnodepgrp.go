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

func resourceAciSpineSwitchPolicyGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciSpineSwitchPolicyGroupCreate,
		UpdateContext: resourceAciSpineSwitchPolicyGroupUpdate,
		ReadContext:   resourceAciSpineSwitchPolicyGroupRead,
		DeleteContext: resourceAciSpineSwitchPolicyGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciSpineSwitchPolicyGroupImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"relation_infra_rs_iacl_spine_profile": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to iacl:SpineProfile",
			},
			"relation_infra_rs_spine_bfd_ipv4_inst_pol": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to bfd:Ipv4InstPol",
			},
			"relation_infra_rs_spine_bfd_ipv6_inst_pol": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to bfd:Ipv6InstPol",
			},
			"relation_infra_rs_spine_copp_profile": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to copp:SpineProfile",
			},
			"relation_infra_rs_spine_p_grp_to_cdp_if_pol": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to cdp:IfPol",
			},
			"relation_infra_rs_spine_p_grp_to_lldp_if_pol": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to lldp:IfPol",
			}})),
	}
}

func getRemoteSpineSwitchPolicyGroup(client *client.Client, dn string) (*models.SpineSwitchPolicyGroup, error) {
	infraSpineAccNodePGrpCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	infraSpineAccNodePGrp := models.SpineSwitchPolicyGroupFromContainer(infraSpineAccNodePGrpCont)
	if infraSpineAccNodePGrp.DistinguishedName == "" {
		return nil, fmt.Errorf("SpineSwitchPolicyGroup %s not found", infraSpineAccNodePGrp.DistinguishedName)
	}
	return infraSpineAccNodePGrp, nil
}

func setSpineSwitchPolicyGroupAttributes(infraSpineAccNodePGrp *models.SpineSwitchPolicyGroup, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(infraSpineAccNodePGrp.DistinguishedName)
	d.Set("description", infraSpineAccNodePGrp.Description)
	infraSpineAccNodePGrpMap, err := infraSpineAccNodePGrp.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", infraSpineAccNodePGrpMap["annotation"])
	d.Set("name", infraSpineAccNodePGrpMap["name"])
	d.Set("name_alias", infraSpineAccNodePGrpMap["nameAlias"])
	return d, nil
}

func getAndSetReadRelationinfraRsIaclSpineProfile(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsIaclSpineProfileData, err := client.ReadRelationinfraRsIaclSpineProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsIaclSpineProfile %v", err)
		d.Set("relation_infra_rs_iacl_spine_profile", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_iacl_spine_profile", infraRsIaclSpineProfileData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsSpineBfdIpv4InstPol(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsSpineBfdIpv4InstPolData, err := client.ReadRelationinfraRsSpineBfdIpv4InstPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsSpineBfdIpv4InstPol %v", err)
		d.Set("relation_infra_rs_spine_bfd_ipv4_inst_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_spine_bfd_ipv4_inst_pol", infraRsSpineBfdIpv4InstPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsSpineBfdIpv6InstPol(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsSpineBfdIpv6InstPolData, err := client.ReadRelationinfraRsSpineBfdIpv6InstPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsSpineBfdIpv6InstPol %v", err)
		d.Set("relation_infra_rs_spine_bfd_ipv6_inst_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_spine_bfd_ipv6_inst_pol", infraRsSpineBfdIpv6InstPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsSpineCoppProfile(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsSpineCoppProfileData, err := client.ReadRelationinfraRsSpineCoppProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsSpineCoppProfile %v", err)
		d.Set("relation_infra_rs_spine_copp_profile", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_spine_copp_profile", infraRsSpineCoppProfileData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsSpinePGrpToCdpIfPol(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsSpinePGrpToCdpIfPolData, err := client.ReadRelationinfraRsSpinePGrpToCdpIfPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsSpinePGrpToCdpIfPol %v", err)
		d.Set("relation_infra_rs_spine_p_grp_to_cdp_if_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_spine_p_grp_to_cdp_if_pol", infraRsSpinePGrpToCdpIfPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsSpinePGrpToLldpIfPol(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsSpinePGrpToLldpIfPolData, err := client.ReadRelationinfraRsSpinePGrpToLldpIfPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsSpinePGrpToLldpIfPol %v", err)
		d.Set("relation_infra_rs_spine_p_grp_to_lldp_if_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_spine_p_grp_to_lldp_if_pol", infraRsSpinePGrpToLldpIfPolData.(string))
	}
	return d, nil
}

func resourceAciSpineSwitchPolicyGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	infraSpineAccNodePGrp, err := getRemoteSpineSwitchPolicyGroup(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setSpineSwitchPolicyGroupAttributes(infraSpineAccNodePGrp, d)
	if err != nil {
		return nil, err
	}

	// infraRsIaclSpineProfile - Beginning Import
	log.Printf("[DEBUG] %s: infraRsIaclSpineProfile - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsIaclSpineProfile(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsIaclSpineProfile - Import finished successfully", d.Get("relation_infra_rs_iacl_spine_profile"))
	}
	// infraRsIaclSpineProfile - Import finished successfully

	// infraRsSpineBfdIpv4InstPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsSpineBfdIpv4InstPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsSpineBfdIpv4InstPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsSpineBfdIpv4InstPol - Import finished successfully", d.Get("relation_infra_rs_spine_bfd_ipv4_inst_pol"))
	}
	// infraRsSpineBfdIpv4InstPol - Import finished successfully

	// infraRsSpineBfdIpv6InstPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsSpineBfdIpv6InstPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsSpineBfdIpv6InstPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsSpineBfdIpv6InstPol - Import finished successfully", d.Get("relation_infra_rs_spine_bfd_ipv6_inst_pol"))
	}
	// infraRsSpineBfdIpv6InstPol - Import finished successfully

	// infraRsSpineCoppProfile - Beginning Import
	log.Printf("[DEBUG] %s: infraRsSpineCoppProfile - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsSpineCoppProfile(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsSpineCoppProfile - Import finished successfully", d.Get("relation_infra_rs_spine_copp_profile"))
	}
	// infraRsSpineCoppProfile - Import finished successfully

	// infraRsSpinePGrpToCdpIfPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsSpinePGrpToCdpIfPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsSpinePGrpToCdpIfPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsSpinePGrpToCdpIfPol - Import finished successfully", d.Get("relation_infra_rs_spine_p_grp_to_cdp_if_pol"))
	}
	// infraRsSpinePGrpToCdpIfPol - Import finished successfully

	// infraRsSpinePGrpToLldpIfPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsSpinePGrpToLldpIfPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsSpinePGrpToLldpIfPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsSpinePGrpToLldpIfPol - Import finished successfully", d.Get("relation_infra_rs_spine_p_grp_to_lldp_if_pol"))
	}
	// infraRsSpinePGrpToLldpIfPol - Import finished successfully

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciSpineSwitchPolicyGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SpineSwitchPolicyGroup: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	infraSpineAccNodePGrpAttr := models.SpineSwitchPolicyGroupAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraSpineAccNodePGrpAttr.Annotation = Annotation.(string)
	} else {
		infraSpineAccNodePGrpAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		infraSpineAccNodePGrpAttr.Name = Name.(string)
	}
	infraSpineAccNodePGrp := models.NewSpineSwitchPolicyGroup(fmt.Sprintf("infra/funcprof/spaccnodepgrp-%s", name), "uni", desc, nameAlias, infraSpineAccNodePGrpAttr)
	err := aciClient.Save(infraSpineAccNodePGrp)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationToinfraRsIaclSpineProfile, ok := d.GetOk("relation_infra_rs_iacl_spine_profile"); ok {
		relationParam := relationToinfraRsIaclSpineProfile.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationToinfraRsSpineBfdIpv4InstPol, ok := d.GetOk("relation_infra_rs_spine_bfd_ipv4_inst_pol"); ok {
		relationParam := relationToinfraRsSpineBfdIpv4InstPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationToinfraRsSpineBfdIpv6InstPol, ok := d.GetOk("relation_infra_rs_spine_bfd_ipv6_inst_pol"); ok {
		relationParam := relationToinfraRsSpineBfdIpv6InstPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationToinfraRsSpineCoppProfile, ok := d.GetOk("relation_infra_rs_spine_copp_profile"); ok {
		relationParam := relationToinfraRsSpineCoppProfile.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationToinfraRsSpinePGrpToCdpIfPol, ok := d.GetOk("relation_infra_rs_spine_p_grp_to_cdp_if_pol"); ok {
		relationParam := relationToinfraRsSpinePGrpToCdpIfPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationToinfraRsSpinePGrpToLldpIfPol, ok := d.GetOk("relation_infra_rs_spine_p_grp_to_lldp_if_pol"); ok {
		relationParam := relationToinfraRsSpinePGrpToLldpIfPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}

	if relationToinfraRsIaclSpineProfile, ok := d.GetOk("relation_infra_rs_iacl_spine_profile"); ok {
		relationParam := relationToinfraRsIaclSpineProfile.(string)
		err = aciClient.CreateRelationinfraRsIaclSpineProfile(infraSpineAccNodePGrp.DistinguishedName, infraSpineAccNodePGrpAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationToinfraRsSpineBfdIpv4InstPol, ok := d.GetOk("relation_infra_rs_spine_bfd_ipv4_inst_pol"); ok {
		relationParam := relationToinfraRsSpineBfdIpv4InstPol.(string)
		err = aciClient.CreateRelationinfraRsSpineBfdIpv4InstPol(infraSpineAccNodePGrp.DistinguishedName, infraSpineAccNodePGrpAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationToinfraRsSpineBfdIpv6InstPol, ok := d.GetOk("relation_infra_rs_spine_bfd_ipv6_inst_pol"); ok {
		relationParam := relationToinfraRsSpineBfdIpv6InstPol.(string)
		err = aciClient.CreateRelationinfraRsSpineBfdIpv6InstPol(infraSpineAccNodePGrp.DistinguishedName, infraSpineAccNodePGrpAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationToinfraRsSpineCoppProfile, ok := d.GetOk("relation_infra_rs_spine_copp_profile"); ok {
		relationParam := relationToinfraRsSpineCoppProfile.(string)
		err = aciClient.CreateRelationinfraRsSpineCoppProfile(infraSpineAccNodePGrp.DistinguishedName, infraSpineAccNodePGrpAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationToinfraRsSpinePGrpToCdpIfPol, ok := d.GetOk("relation_infra_rs_spine_p_grp_to_cdp_if_pol"); ok {
		relationParam := relationToinfraRsSpinePGrpToCdpIfPol.(string)
		err = aciClient.CreateRelationinfraRsSpinePGrpToCdpIfPol(infraSpineAccNodePGrp.DistinguishedName, infraSpineAccNodePGrpAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationToinfraRsSpinePGrpToLldpIfPol, ok := d.GetOk("relation_infra_rs_spine_p_grp_to_lldp_if_pol"); ok {
		relationParam := relationToinfraRsSpinePGrpToLldpIfPol.(string)
		err = aciClient.CreateRelationinfraRsSpinePGrpToLldpIfPol(infraSpineAccNodePGrp.DistinguishedName, infraSpineAccNodePGrpAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(infraSpineAccNodePGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciSpineSwitchPolicyGroupRead(ctx, d, m)
}

func resourceAciSpineSwitchPolicyGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SpineSwitchPolicyGroup: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	infraSpineAccNodePGrpAttr := models.SpineSwitchPolicyGroupAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		infraSpineAccNodePGrpAttr.Annotation = Annotation.(string)
	} else {
		infraSpineAccNodePGrpAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		infraSpineAccNodePGrpAttr.Name = Name.(string)
	}
	infraSpineAccNodePGrp := models.NewSpineSwitchPolicyGroup(fmt.Sprintf("infra/funcprof/spaccnodepgrp-%s", name), "uni", desc, nameAlias, infraSpineAccNodePGrpAttr)
	infraSpineAccNodePGrp.Status = "modified"
	err := aciClient.Save(infraSpineAccNodePGrp)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_infra_rs_iacl_spine_profile") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_iacl_spine_profile")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_infra_rs_spine_bfd_ipv4_inst_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_spine_bfd_ipv4_inst_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_infra_rs_spine_bfd_ipv6_inst_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_spine_bfd_ipv6_inst_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_infra_rs_spine_copp_profile") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_spine_copp_profile")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_infra_rs_spine_p_grp_to_cdp_if_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_spine_p_grp_to_cdp_if_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_infra_rs_spine_p_grp_to_lldp_if_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_spine_p_grp_to_lldp_if_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("relation_infra_rs_iacl_spine_profile") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_iacl_spine_profile")
		err = aciClient.DeleteRelationinfraRsIaclSpineProfile(infraSpineAccNodePGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsIaclSpineProfile(infraSpineAccNodePGrp.DistinguishedName, infraSpineAccNodePGrpAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_spine_bfd_ipv4_inst_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_spine_bfd_ipv4_inst_pol")
		err = aciClient.DeleteRelationinfraRsSpineBfdIpv4InstPol(infraSpineAccNodePGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsSpineBfdIpv4InstPol(infraSpineAccNodePGrp.DistinguishedName, infraSpineAccNodePGrpAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_spine_bfd_ipv6_inst_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_spine_bfd_ipv6_inst_pol")
		err = aciClient.DeleteRelationinfraRsSpineBfdIpv6InstPol(infraSpineAccNodePGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsSpineBfdIpv6InstPol(infraSpineAccNodePGrp.DistinguishedName, infraSpineAccNodePGrpAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_spine_copp_profile") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_spine_copp_profile")
		err = aciClient.DeleteRelationinfraRsSpineCoppProfile(infraSpineAccNodePGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsSpineCoppProfile(infraSpineAccNodePGrp.DistinguishedName, infraSpineAccNodePGrpAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_spine_p_grp_to_cdp_if_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_spine_p_grp_to_cdp_if_pol")
		err = aciClient.DeleteRelationinfraRsSpinePGrpToCdpIfPol(infraSpineAccNodePGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsSpinePGrpToCdpIfPol(infraSpineAccNodePGrp.DistinguishedName, infraSpineAccNodePGrpAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_spine_p_grp_to_lldp_if_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_spine_p_grp_to_lldp_if_pol")
		err = aciClient.DeleteRelationinfraRsSpinePGrpToLldpIfPol(infraSpineAccNodePGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsSpinePGrpToLldpIfPol(infraSpineAccNodePGrp.DistinguishedName, infraSpineAccNodePGrpAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(infraSpineAccNodePGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciSpineSwitchPolicyGroupRead(ctx, d, m)
}

func resourceAciSpineSwitchPolicyGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	infraSpineAccNodePGrp, err := getRemoteSpineSwitchPolicyGroup(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}
	setSpineSwitchPolicyGroupAttributes(infraSpineAccNodePGrp, d)

	// infraRsIaclSpineProfileData, err := aciClient.ReadRelationinfraRsIaclSpineProfile(dn)
	// if err != nil {
	// 	log.Printf("[DEBUG] Error while reading relation infraRsIaclSpineProfile %v", err)
	// 	d.Set("infra_rs_iacl_spine_profile", "")
	// } else {
	// 	if _, ok := d.GetOk("relation_infra_rs_iacl_spine_profile"); ok {
	// 		tfName := GetMOName(d.Get("relation_infra_rs_iacl_spine_profile").(string))
	// 		if tfName != infraRsIaclSpineProfileData {
	// 			d.Set("relation_infra_rs_iacl_spine_profile", "")
	// 		}
	// 	}
	// }

	// infraRsSpineBfdIpv4InstPolData, err := aciClient.ReadRelationinfraRsSpineBfdIpv4InstPol(dn)
	// if err != nil {
	// 	log.Printf("[DEBUG] Error while reading relation infraRsSpineBfdIpv4InstPol %v", err)
	// 	d.Set("infra_rs_spine_bfd_ipv4_inst_pol", "")
	// } else {
	// 	if _, ok := d.GetOk("relation_infra_rs_spine_bfd_ipv4_inst_pol"); ok {
	// 		tfName := GetMOName(d.Get("relation_infra_rs_spine_bfd_ipv4_inst_pol").(string))
	// 		if tfName != infraRsSpineBfdIpv4InstPolData {
	// 			d.Set("relation_infra_rs_spine_bfd_ipv4_inst_pol", "")
	// 		}
	// 	}
	// }

	// infraRsSpineBfdIpv6InstPolData, err := aciClient.ReadRelationinfraRsSpineBfdIpv6InstPol(dn)
	// if err != nil {
	// 	log.Printf("[DEBUG] Error while reading relation infraRsSpineBfdIpv6InstPol %v", err)
	// 	d.Set("infra_rs_spine_bfd_ipv6_inst_pol", "")
	// } else {
	// 	if _, ok := d.GetOk("relation_infra_rs_spine_bfd_ipv6_inst_pol"); ok {
	// 		tfName := GetMOName(d.Get("relation_infra_rs_spine_bfd_ipv6_inst_pol").(string))
	// 		if tfName != infraRsSpineBfdIpv6InstPolData {
	// 			d.Set("relation_infra_rs_spine_bfd_ipv6_inst_pol", "")
	// 		}
	// 	}
	// }

	// infraRsSpineCoppProfileData, err := aciClient.ReadRelationinfraRsSpineCoppProfile(dn)
	// if err != nil {
	// 	log.Printf("[DEBUG] Error while reading relation infraRsSpineCoppProfile %v", err)
	// 	d.Set("infra_rs_spine_copp_profile", "")
	// } else {
	// 	if _, ok := d.GetOk("relation_infra_rs_spine_copp_profile"); ok {
	// 		tfName := GetMOName(d.Get("relation_infra_rs_spine_copp_profile").(string))
	// 		if tfName != infraRsSpineCoppProfileData {
	// 			d.Set("relation_infra_rs_spine_copp_profile", "")
	// 		}
	// 	}
	// }

	// infraRsSpinePGrpToCdpIfPolData, err := aciClient.ReadRelationinfraRsSpinePGrpToCdpIfPol(dn)
	// if err != nil {
	// 	log.Printf("[DEBUG] Error while reading relation infraRsSpinePGrpToCdpIfPol %v", err)
	// 	d.Set("infra_rs_spine_p_grp_to_cdp_if_pol", "")
	// } else {
	// 	if _, ok := d.GetOk("relation_infra_rs_spine_p_grp_to_cdp_if_pol"); ok {
	// 		tfName := GetMOName(d.Get("relation_infra_rs_spine_p_grp_to_cdp_if_pol").(string))
	// 		if tfName != infraRsSpinePGrpToCdpIfPolData {
	// 			d.Set("relation_infra_rs_spine_p_grp_to_cdp_if_pol", "")
	// 		}
	// 	}
	// }

	// infraRsSpinePGrpToLldpIfPolData, err := aciClient.ReadRelationinfraRsSpinePGrpToLldpIfPol(dn)
	// if err != nil {
	// 	log.Printf("[DEBUG] Error while reading relation infraRsSpinePGrpToLldpIfPol %v", err)
	// 	d.Set("infra_rs_spine_p_grp_to_lldp_if_pol", "")
	// } else {
	// 	if _, ok := d.GetOk("relation_infra_rs_spine_p_grp_to_lldp_if_pol"); ok {
	// 		tfName := GetMOName(d.Get("relation_infra_rs_spine_p_grp_to_lldp_if_pol").(string))
	// 		if tfName != infraRsSpinePGrpToLldpIfPolData {
	// 			d.Set("relation_infra_rs_spine_p_grp_to_lldp_if_pol", "")
	// 		}
	// 	}
	// }

	// infraRsIaclSpineProfile - Beginning Read
	log.Printf("[DEBUG] %s: infraRsIaclSpineProfile - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsIaclSpineProfile(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsIaclSpineProfile - Read finished successfully", d.Get("relation_infra_rs_iacl_spine_profile"))
	}
	// infraRsIaclSpineProfile - Read finished successfully

	// infraRsSpineBfdIpv4InstPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsSpineBfdIpv4InstPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsSpineBfdIpv4InstPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsSpineBfdIpv4InstPol - Read finished successfully", d.Get("relation_infra_rs_spine_bfd_ipv4_inst_pol"))
	}
	// infraRsSpineBfdIpv4InstPol - Read finished successfully

	// infraRsSpineBfdIpv6InstPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsSpineBfdIpv6InstPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsSpineBfdIpv6InstPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsSpineBfdIpv6InstPol - Read finished successfully", d.Get("relation_infra_rs_spine_bfd_ipv6_inst_pol"))
	}
	// infraRsSpineBfdIpv6InstPol - Read finished successfully

	// infraRsSpineCoppProfile - Beginning Read
	log.Printf("[DEBUG] %s: infraRsSpineCoppProfile - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsSpineCoppProfile(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsSpineCoppProfile - Read finished successfully", d.Get("relation_infra_rs_spine_copp_profile"))
	}
	// infraRsSpineCoppProfile - Read finished successfully

	// infraRsSpinePGrpToCdpIfPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsSpinePGrpToCdpIfPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsSpinePGrpToCdpIfPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsSpinePGrpToCdpIfPol - Read finished successfully", d.Get("relation_infra_rs_spine_p_grp_to_cdp_if_pol"))
	}
	// infraRsSpinePGrpToCdpIfPol - Read finished successfully

	// infraRsSpinePGrpToLldpIfPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsSpinePGrpToLldpIfPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsSpinePGrpToLldpIfPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsSpinePGrpToLldpIfPol - Read finished successfully", d.Get("relation_infra_rs_spine_p_grp_to_lldp_if_pol"))
	}
	// infraRsSpinePGrpToLldpIfPol - Read finished successfully

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciSpineSwitchPolicyGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraSpineAccNodePGrp")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
