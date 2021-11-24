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

func resourceAciVSwitchPolicyGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciVSwitchPolicyGroupCreate,
		UpdateContext: resourceAciVSwitchPolicyGroupUpdate,
		ReadContext:   resourceAciVSwitchPolicyGroupRead,
		DeleteContext: resourceAciVSwitchPolicyGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciVSwitchPolicyGroupImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"vmm_domain_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"relation_vmm_rs_vswitch_exporter_pol": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Create relation to netflowVmmExporterPol",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"active_flow_time_out": {
							Optional: true,
							Type:     schema.TypeString,
						},
						"idle_flow_time_out": {
							Optional: true,
							Type:     schema.TypeString,
						},
						"sampling_rate": {
							Optional: true,
							Type:     schema.TypeString,
						},
						"target_dn": {
							Required: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"relation_vmm_rs_vswitch_override_cdp_if_pol": {
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to cdp:IfPol",
			},
			"relation_vmm_rs_vswitch_override_fw_pol": {
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to nws:FwPol",
			},
			"relation_vmm_rs_vswitch_override_lacp_pol": {
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to lacp:LagPol",
			},
			"relation_vmm_rs_vswitch_override_lldp_if_pol": {
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to lldp:IfPol",
			},
			"relation_vmm_rs_vswitch_override_mcp_if_pol": {
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to mcp:IfPol",
			},
			"relation_vmm_rs_vswitch_override_mtu_pol": {
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to l2:InstPol",
			},
			"relation_vmm_rs_vswitch_override_stp_pol": {
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to stp:IfPol",
			}})),
	}
}

func getRemoteVSwitchPolicyGroup(client *client.Client, dn string) (*models.VSwitchPolicyGroup, error) {
	vmmVSwitchPolicyContCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	vmmVSwitchPolicyCont := models.VSwitchPolicyGroupFromContainer(vmmVSwitchPolicyContCont)
	if vmmVSwitchPolicyCont.DistinguishedName == "" {
		return nil, fmt.Errorf("VSwitchPolicyGroup %s not found", vmmVSwitchPolicyCont.DistinguishedName)
	}
	return vmmVSwitchPolicyCont, nil
}

func setVSwitchPolicyGroupAttributes(vmmVSwitchPolicyCont *models.VSwitchPolicyGroup, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(vmmVSwitchPolicyCont.DistinguishedName)
	d.Set("description", vmmVSwitchPolicyCont.Description)
	vmmVSwitchPolicyContMap, err := vmmVSwitchPolicyCont.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("vmm_domain_dn", GetParentDn(vmmVSwitchPolicyCont.DistinguishedName, fmt.Sprintf("/vswitchpolcont")))
	d.Set("annotation", vmmVSwitchPolicyContMap["annotation"])
	d.Set("name_alias", vmmVSwitchPolicyContMap["nameAlias"])
	return d, nil
}

func resourceAciVSwitchPolicyGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	vmmVSwitchPolicyCont, err := getRemoteVSwitchPolicyGroup(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setVSwitchPolicyGroupAttributes(vmmVSwitchPolicyCont, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciVSwitchPolicyGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] VSwitchPolicyGroup: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	VMMDomainDn := d.Get("vmm_domain_dn").(string)

	vmmVSwitchPolicyContAttr := models.VSwitchPolicyGroupAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vmmVSwitchPolicyContAttr.Annotation = Annotation.(string)
	} else {
		vmmVSwitchPolicyContAttr.Annotation = "{}"
	}

	vmmVSwitchPolicyCont := models.NewVSwitchPolicyGroup(fmt.Sprintf("vswitchpolcont"), VMMDomainDn, desc, nameAlias, vmmVSwitchPolicyContAttr)

	err := aciClient.Save(vmmVSwitchPolicyCont)
	if err != nil {
		return diag.FromErr(err)
	}
	checkDns := make([]string, 0, 1)

	if relationTovmmRsVswitchExporterPol, ok := d.GetOk("relation_vmm_rs_vswitch_exporter_pol"); ok {
		relationParamList := toStringList(relationTovmmRsVswitchExporterPol.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTovmmRsVswitchOverrideCdpIfPol, ok := d.GetOk("relation_vmm_rs_vswitch_override_cdp_if_pol"); ok {
		relationParam := relationTovmmRsVswitchOverrideCdpIfPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationTovmmRsVswitchOverrideFwPol, ok := d.GetOk("relation_vmm_rs_vswitch_override_fw_pol"); ok {
		relationParam := relationTovmmRsVswitchOverrideFwPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationTovmmRsVswitchOverrideLacpPol, ok := d.GetOk("relation_vmm_rs_vswitch_override_lacp_pol"); ok {
		relationParam := relationTovmmRsVswitchOverrideLacpPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationTovmmRsVswitchOverrideLldpIfPol, ok := d.GetOk("relation_vmm_rs_vswitch_override_lldp_if_pol"); ok {
		relationParam := relationTovmmRsVswitchOverrideLldpIfPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationTovmmRsVswitchOverrideMcpIfPol, ok := d.GetOk("relation_vmm_rs_vswitch_override_mcp_if_pol"); ok {
		relationParam := relationTovmmRsVswitchOverrideMcpIfPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationTovmmRsVswitchOverrideMtuPol, ok := d.GetOk("relation_vmm_rs_vswitch_override_mtu_pol"); ok {
		relationParam := relationTovmmRsVswitchOverrideMtuPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationTovmmRsVswitchOverrideStpPol, ok := d.GetOk("relation_vmm_rs_vswitch_override_stp_pol"); ok {
		relationParam := relationTovmmRsVswitchOverrideStpPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTovmmRsVswitchExporterPol, ok := d.GetOk("relation_vmm_rs_vswitch_exporter_pol"); ok {
		relationParamList := relationTovmmRsVswitchExporterPol.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationvmmRsVswitchExporterPol(vmmVSwitchPolicyCont.DistinguishedName, vmmVSwitchPolicyContAttr.Annotation, paramMap["active_flow_time_out"].(string), paramMap["idle_flow_time_out"].(string), paramMap["sampling_rate"].(string), paramMap["target_dn"].(string))

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if relationTovmmRsVswitchOverrideCdpIfPol, ok := d.GetOk("relation_vmm_rs_vswitch_override_cdp_if_pol"); ok {
		relationParam := relationTovmmRsVswitchOverrideCdpIfPol.(string)
		err = aciClient.CreateRelationvmmRsVswitchOverrideCdpIfPol(vmmVSwitchPolicyCont.DistinguishedName, vmmVSwitchPolicyContAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationTovmmRsVswitchOverrideFwPol, ok := d.GetOk("relation_vmm_rs_vswitch_override_fw_pol"); ok {
		relationParam := relationTovmmRsVswitchOverrideFwPol.(string)
		err = aciClient.CreateRelationvmmRsVswitchOverrideFwPol(vmmVSwitchPolicyCont.DistinguishedName, vmmVSwitchPolicyContAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationTovmmRsVswitchOverrideLacpPol, ok := d.GetOk("relation_vmm_rs_vswitch_override_lacp_pol"); ok {
		relationParam := relationTovmmRsVswitchOverrideLacpPol.(string)
		err = aciClient.CreateRelationvmmRsVswitchOverrideLacpPol(vmmVSwitchPolicyCont.DistinguishedName, vmmVSwitchPolicyContAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationTovmmRsVswitchOverrideLldpIfPol, ok := d.GetOk("relation_vmm_rs_vswitch_override_lldp_if_pol"); ok {
		relationParam := relationTovmmRsVswitchOverrideLldpIfPol.(string)
		err = aciClient.CreateRelationvmmRsVswitchOverrideLldpIfPol(vmmVSwitchPolicyCont.DistinguishedName, vmmVSwitchPolicyContAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	if relationTovmmRsVswitchOverrideMcpIfPol, ok := d.GetOk("relation_vmm_rs_vswitch_override_mcp_if_pol"); ok {
		relationParam := relationTovmmRsVswitchOverrideMcpIfPol.(string)
		err = aciClient.CreateRelationvmmRsVswitchOverrideMcpIfPol(vmmVSwitchPolicyCont.DistinguishedName, vmmVSwitchPolicyContAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationTovmmRsVswitchOverrideMtuPol, ok := d.GetOk("relation_vmm_rs_vswitch_override_mtu_pol"); ok {
		relationParam := relationTovmmRsVswitchOverrideMtuPol.(string)
		err = aciClient.CreateRelationvmmRsVswitchOverrideMtuPol(vmmVSwitchPolicyCont.DistinguishedName, vmmVSwitchPolicyContAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationTovmmRsVswitchOverrideStpPol, ok := d.GetOk("relation_vmm_rs_vswitch_override_stp_pol"); ok {
		relationParam := relationTovmmRsVswitchOverrideStpPol.(string)
		err = aciClient.CreateRelationvmmRsVswitchOverrideStpPol(vmmVSwitchPolicyCont.DistinguishedName, vmmVSwitchPolicyContAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(vmmVSwitchPolicyCont.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciVSwitchPolicyGroupRead(ctx, d, m)
}

func resourceAciVSwitchPolicyGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] VSwitchPolicyGroup: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	VMMDomainDn := d.Get("vmm_domain_dn").(string)
	vmmVSwitchPolicyContAttr := models.VSwitchPolicyGroupAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		vmmVSwitchPolicyContAttr.Annotation = Annotation.(string)
	} else {
		vmmVSwitchPolicyContAttr.Annotation = "{}"
	}

	vmmVSwitchPolicyCont := models.NewVSwitchPolicyGroup(fmt.Sprintf("vswitchpolcont"), VMMDomainDn, desc, nameAlias, vmmVSwitchPolicyContAttr)

	vmmVSwitchPolicyCont.Status = "modified"
	err := aciClient.Save(vmmVSwitchPolicyCont)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_vmm_rs_vswitch_exporter_pol") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_vmm_rs_vswitch_exporter_pol")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())
		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_vmm_rs_vswitch_override_cdp_if_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vmm_rs_vswitch_override_cdp_if_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_vmm_rs_vswitch_override_fw_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vmm_rs_vswitch_override_fw_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_vmm_rs_vswitch_override_lacp_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vmm_rs_vswitch_override_lacp_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_vmm_rs_vswitch_override_lldp_if_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vmm_rs_vswitch_override_lldp_if_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_vmm_rs_vswitch_override_mcp_if_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vmm_rs_vswitch_override_mcp_if_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_vmm_rs_vswitch_override_mtu_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vmm_rs_vswitch_override_mtu_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_vmm_rs_vswitch_override_stp_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vmm_rs_vswitch_override_stp_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_vmm_rs_vswitch_exporter_pol") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_vmm_rs_vswitch_exporter_pol")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.DeleteRelationvmmRsVswitchExporterPol(vmmVSwitchPolicyCont.DistinguishedName, paramMap["target_dn"].(string))

			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationvmmRsVswitchExporterPol(vmmVSwitchPolicyCont.DistinguishedName, vmmVSwitchPolicyContAttr.Annotation, paramMap["active_flow_time_out"].(string), paramMap["idle_flow_time_out"].(string), paramMap["sampling_rate"].(string), paramMap["target_dn"].(string))

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("relation_vmm_rs_vswitch_override_cdp_if_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vmm_rs_vswitch_override_cdp_if_pol")
		err = aciClient.DeleteRelationvmmRsVswitchOverrideCdpIfPol(vmmVSwitchPolicyCont.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvmmRsVswitchOverrideCdpIfPol(vmmVSwitchPolicyCont.DistinguishedName, vmmVSwitchPolicyContAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_vmm_rs_vswitch_override_fw_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vmm_rs_vswitch_override_fw_pol")
		err = aciClient.DeleteRelationvmmRsVswitchOverrideFwPol(vmmVSwitchPolicyCont.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvmmRsVswitchOverrideFwPol(vmmVSwitchPolicyCont.DistinguishedName, vmmVSwitchPolicyContAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_vmm_rs_vswitch_override_lacp_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vmm_rs_vswitch_override_lacp_pol")
		err = aciClient.DeleteRelationvmmRsVswitchOverrideLacpPol(vmmVSwitchPolicyCont.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvmmRsVswitchOverrideLacpPol(vmmVSwitchPolicyCont.DistinguishedName, vmmVSwitchPolicyContAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_vmm_rs_vswitch_override_lldp_if_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vmm_rs_vswitch_override_lldp_if_pol")
		err = aciClient.DeleteRelationvmmRsVswitchOverrideLldpIfPol(vmmVSwitchPolicyCont.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvmmRsVswitchOverrideLldpIfPol(vmmVSwitchPolicyCont.DistinguishedName, vmmVSwitchPolicyContAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_vmm_rs_vswitch_override_mcp_if_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vmm_rs_vswitch_override_mcp_if_pol")
		err = aciClient.DeleteRelationvmmRsVswitchOverrideMcpIfPol(vmmVSwitchPolicyCont.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvmmRsVswitchOverrideMcpIfPol(vmmVSwitchPolicyCont.DistinguishedName, vmmVSwitchPolicyContAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_vmm_rs_vswitch_override_mtu_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vmm_rs_vswitch_override_mtu_pol")
		err = aciClient.DeleteRelationvmmRsVswitchOverrideMtuPol(vmmVSwitchPolicyCont.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvmmRsVswitchOverrideMtuPol(vmmVSwitchPolicyCont.DistinguishedName, vmmVSwitchPolicyContAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_vmm_rs_vswitch_override_stp_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vmm_rs_vswitch_override_stp_pol")
		err = aciClient.DeleteRelationvmmRsVswitchOverrideStpPol(vmmVSwitchPolicyCont.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvmmRsVswitchOverrideStpPol(vmmVSwitchPolicyCont.DistinguishedName, vmmVSwitchPolicyContAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(vmmVSwitchPolicyCont.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciVSwitchPolicyGroupRead(ctx, d, m)
}

func resourceAciVSwitchPolicyGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	vmmVSwitchPolicyCont, err := getRemoteVSwitchPolicyGroup(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	_, err = setVSwitchPolicyGroupAttributes(vmmVSwitchPolicyCont, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	vmmRsVswitchExporterPolData, err := aciClient.ReadRelationvmmRsVswitchExporterPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsVswitchExporterPol %v", err)
	} else {
		relParams := make([]map[string]string, 0, 1)
		relParamsList := vmmRsVswitchExporterPolData.([]map[string]string)
		for _, obj := range relParamsList {
			relParams = append(relParams, map[string]string{
				"active_flow_time_out": obj["activeFlowTimeOut"],
				"idle_flow_time_out":   obj["idleFlowTimeOut"],
				"sampling_rate":        obj["samplingRate"],
				"target_dn":            obj["tDn"],
			})
		}
		d.Set("relation_vmm_rs_vswitch_exporter_pol", relParams)
	}

	vmmRsVswitchOverrideCdpIfPolData, err := aciClient.ReadRelationvmmRsVswitchOverrideCdpIfPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsVswitchOverrideCdpIfPol %v", err)
		d.Set("relation_vmm_rs_vswitch_override_cdp_if_pol", "")
	} else {
		d.Set("relation_vmm_rs_vswitch_override_cdp_if_pol", vmmRsVswitchOverrideCdpIfPolData.(string))
	}

	vmmRsVswitchOverrideFwPolData, err := aciClient.ReadRelationvmmRsVswitchOverrideFwPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsVswitchOverrideFwPol %v", err)
		d.Set("relation_vmm_rs_vswitch_override_fw_pol", "")
	} else {
		d.Set("relation_vmm_rs_vswitch_override_fw_pol", vmmRsVswitchOverrideFwPolData.(string))
	}

	vmmRsVswitchOverrideLacpPolData, err := aciClient.ReadRelationvmmRsVswitchOverrideLacpPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsVswitchOverrideLacpPol %v", err)
		d.Set("relation_vmm_rs_vswitch_override_lacp_pol", "")
	} else {
		d.Set("relation_vmm_rs_vswitch_override_lacp_pol", vmmRsVswitchOverrideLacpPolData.(string))
	}

	vmmRsVswitchOverrideLldpIfPolData, err := aciClient.ReadRelationvmmRsVswitchOverrideLldpIfPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsVswitchOverrideLldpIfPol %v", err)
		d.Set("relation_vmm_rs_vswitch_override_lldp_if_pol", "")
	} else {
		d.Set("relation_vmm_rs_vswitch_override_lldp_if_pol", vmmRsVswitchOverrideLldpIfPolData.(string))
	}

	vmmRsVswitchOverrideMcpIfPolData, err := aciClient.ReadRelationvmmRsVswitchOverrideMcpIfPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsVswitchOverrideMcpIfPol %v", err)
		d.Set("relation_vmm_rs_vswitch_override_mcp_if_pol", "")
	} else {
		d.Set("relation_vmm_rs_vswitch_override_mcp_if_pol", vmmRsVswitchOverrideMcpIfPolData.(string))
	}

	vmmRsVswitchOverrideMtuPolData, err := aciClient.ReadRelationvmmRsVswitchOverrideMtuPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsVswitchOverrideMtuPol %v", err)
		d.Set("relation_vmm_rs_vswitch_override_mtu_pol", "")
	} else {
		d.Set("relation_vmm_rs_vswitch_override_mtu_pol", vmmRsVswitchOverrideMtuPolData.(string))
	}

	vmmRsVswitchOverrideStpPolData, err := aciClient.ReadRelationvmmRsVswitchOverrideStpPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsVswitchOverrideStpPol %v", err)
		d.Set("relation_vmm_rs_vswitch_override_stp_pol", "")
	} else {
		d.Set("relation_vmm_rs_vswitch_override_stp_pol", vmmRsVswitchOverrideStpPolData.(string))
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciVSwitchPolicyGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vmmVSwitchPolicyCont")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
