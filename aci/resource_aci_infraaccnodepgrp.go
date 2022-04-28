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

func resourceAciAccessSwitchPolicyGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciAccessSwitchPolicyGroupCreate,
		UpdateContext: resourceAciAccessSwitchPolicyGroupUpdate,
		ReadContext:   resourceAciAccessSwitchPolicyGroupRead,
		DeleteContext: resourceAciAccessSwitchPolicyGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciAccessSwitchPolicyGroupImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"relation_infra_rs_bfd_ipv4_inst_pol": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to bfd:Ipv4InstPol",
			},
			"relation_infra_rs_bfd_ipv6_inst_pol": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to bfd:Ipv6InstPol",
			},
			"relation_infra_rs_bfd_mh_ipv4_inst_pol": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to bfd:MhIpv4InstPol",
			},
			"relation_infra_rs_bfd_mh_ipv6_inst_pol": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to bfd:MhIpv6InstPol",
			},
			"relation_infra_rs_equipment_flash_config_pol": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to equipment:FlashConfigPol",
			},
			"relation_infra_rs_fc_fabric_pol": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to fc:FabricPol",
			},
			"relation_infra_rs_fc_inst_pol": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to fc:InstPol",
			},
			"relation_infra_rs_iacl_leaf_profile": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to iacl:LeafProfile",
			},
			"relation_infra_rs_l2_node_auth_pol": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to l2:NodeAuthPol",
			},
			"relation_infra_rs_leaf_copp_profile": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to copp:LeafProfile",
			},
			"relation_infra_rs_leaf_p_grp_to_cdp_if_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to cdp:IfPol",
			},
			"relation_infra_rs_leaf_p_grp_to_lldp_if_pol": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to lldp:IfPol",
			},
			"relation_infra_rs_mon_node_infra_pol": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to mon:InfraPol",
			},
			"relation_infra_rs_mst_inst_pol": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to stp:InstPol",
			},
			"relation_infra_rs_netflow_node_pol": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to netflow:NodePol",
			},
			"relation_infra_rs_poe_inst_pol": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to poe:InstPol",
			},
			"relation_infra_rs_topoctrl_fast_link_failover_inst_pol": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to topoctrl:FastLinkFailoverInstPol",
			},
			"relation_infra_rs_topoctrl_fwd_scale_prof_pol": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to topoctrl:FwdScaleProfilePol",
			}})),
	}
}

func getRemoteAccessSwitchPolicyGroup(client *client.Client, dn string) (*models.AccessSwitchPolicyGroup, error) {
	infraAccNodePGrpCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	infraAccNodePGrp := models.AccessSwitchPolicyGroupFromContainer(infraAccNodePGrpCont)
	if infraAccNodePGrp.DistinguishedName == "" {
		return nil, fmt.Errorf("AccessSwitchPolicyGroup %s not found", infraAccNodePGrp.DistinguishedName)
	}
	return infraAccNodePGrp, nil
}

func setAccessSwitchPolicyGroupAttributes(infraAccNodePGrp *models.AccessSwitchPolicyGroup, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(infraAccNodePGrp.DistinguishedName)
	d.Set("description", infraAccNodePGrp.Description)
	infraAccNodePGrpMap, err := infraAccNodePGrp.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", infraAccNodePGrpMap["annotation"])
	d.Set("name", infraAccNodePGrpMap["name"])
	d.Set("name_alias", infraAccNodePGrpMap["nameAlias"])
	return d, nil
}

func resourceAciAccessSwitchPolicyGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	infraAccNodePGrp, err := getRemoteAccessSwitchPolicyGroup(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setAccessSwitchPolicyGroupAttributes(infraAccNodePGrp, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciAccessSwitchPolicyGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] AccessSwitchPolicyGroup: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	infraAccNodePGrpAttr := models.AccessSwitchPolicyGroupAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraAccNodePGrpAttr.Annotation = Annotation.(string)
	} else {
		infraAccNodePGrpAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		infraAccNodePGrpAttr.Name = Name.(string)
	}
	infraAccNodePGrp := models.NewAccessSwitchPolicyGroup(fmt.Sprintf("infra/funcprof/accnodepgrp-%s", name), "uni", desc, nameAlias, infraAccNodePGrpAttr)
	err := aciClient.Save(infraAccNodePGrp)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationToinfraRsBfdIpv4InstPol, ok := d.GetOk("relation_infra_rs_bfd_ipv4_inst_pol"); ok {
		relationParam := relationToinfraRsBfdIpv4InstPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationToinfraRsBfdIpv6InstPol, ok := d.GetOk("relation_infra_rs_bfd_ipv6_inst_pol"); ok {
		relationParam := relationToinfraRsBfdIpv6InstPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationToinfraRsBfdMhIpv4InstPol, ok := d.GetOk("relation_infra_rs_bfd_mh_ipv4_inst_pol"); ok {
		relationParam := relationToinfraRsBfdMhIpv4InstPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationToinfraRsBfdMhIpv6InstPol, ok := d.GetOk("relation_infra_rs_bfd_mh_ipv6_inst_pol"); ok {
		relationParam := relationToinfraRsBfdMhIpv6InstPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationToinfraRsEquipmentFlashConfigPol, ok := d.GetOk("relation_infra_rs_equipment_flash_config_pol"); ok {
		relationParam := relationToinfraRsEquipmentFlashConfigPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationToinfraRsFcFabricPol, ok := d.GetOk("relation_infra_rs_fc_fabric_pol"); ok {
		relationParam := relationToinfraRsFcFabricPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationToinfraRsFcInstPol, ok := d.GetOk("relation_infra_rs_fc_inst_pol"); ok {
		relationParam := relationToinfraRsFcInstPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationToinfraRsIaclLeafProfile, ok := d.GetOk("relation_infra_rs_iacl_leaf_profile"); ok {
		relationParam := relationToinfraRsIaclLeafProfile.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationToinfraRsL2NodeAuthPol, ok := d.GetOk("relation_infra_rs_l2_node_auth_pol"); ok {
		relationParam := relationToinfraRsL2NodeAuthPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationToinfraRsLeafCoppProfile, ok := d.GetOk("relation_infra_rs_leaf_copp_profile"); ok {
		relationParam := relationToinfraRsLeafCoppProfile.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationToinfraRsLeafPGrpToCdpIfPol, ok := d.GetOk("relation_infra_rs_leaf_p_grp_to_cdp_if_pol"); ok {
		relationParam := relationToinfraRsLeafPGrpToCdpIfPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationToinfraRsLeafPGrpToLldpIfPol, ok := d.GetOk("relation_infra_rs_leaf_p_grp_to_lldp_if_pol"); ok {
		relationParam := relationToinfraRsLeafPGrpToLldpIfPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationToinfraRsMonNodeInfraPol, ok := d.GetOk("relation_infra_rs_mon_node_infra_pol"); ok {
		relationParam := relationToinfraRsMonNodeInfraPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationToinfraRsMstInstPol, ok := d.GetOk("relation_infra_rs_mst_inst_pol"); ok {
		relationParam := relationToinfraRsMstInstPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationToinfraRsNetflowNodePol, ok := d.GetOk("relation_infra_rs_netflow_node_pol"); ok {
		relationParam := relationToinfraRsNetflowNodePol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationToinfraRsPoeInstPol, ok := d.GetOk("relation_infra_rs_poe_inst_pol"); ok {
		relationParam := relationToinfraRsPoeInstPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationToinfraRsTopoctrlFastLinkFailoverInstPol, ok := d.GetOk("relation_infra_rs_topoctrl_fast_link_failover_inst_pol"); ok {
		relationParam := relationToinfraRsTopoctrlFastLinkFailoverInstPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationToinfraRsTopoctrlFwdScaleProfPol, ok := d.GetOk("relation_infra_rs_topoctrl_fwd_scale_prof_pol"); ok {
		relationParam := relationToinfraRsTopoctrlFwdScaleProfPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}

	if relationToinfraRsBfdIpv4InstPol, ok := d.GetOk("relation_infra_rs_bfd_ipv4_inst_pol"); ok {
		relationParam := relationToinfraRsBfdIpv4InstPol.(string)
		err = aciClient.CreateRelationinfraRsBfdIpv4InstPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationToinfraRsBfdIpv6InstPol, ok := d.GetOk("relation_infra_rs_bfd_ipv6_inst_pol"); ok {
		relationParam := relationToinfraRsBfdIpv6InstPol.(string)
		err = aciClient.CreateRelationinfraRsBfdIpv6InstPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationToinfraRsBfdMhIpv4InstPol, ok := d.GetOk("relation_infra_rs_bfd_mh_ipv4_inst_pol"); ok {
		relationParam := relationToinfraRsBfdMhIpv4InstPol.(string)
		err = aciClient.CreateRelationinfraRsBfdMhIpv4InstPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationToinfraRsBfdMhIpv6InstPol, ok := d.GetOk("relation_infra_rs_bfd_mh_ipv6_inst_pol"); ok {
		relationParam := relationToinfraRsBfdMhIpv6InstPol.(string)
		err = aciClient.CreateRelationinfraRsBfdMhIpv6InstPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationToinfraRsEquipmentFlashConfigPol, ok := d.GetOk("relation_infra_rs_equipment_flash_config_pol"); ok {
		relationParam := relationToinfraRsEquipmentFlashConfigPol.(string)
		err = aciClient.CreateRelationinfraRsEquipmentFlashConfigPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationToinfraRsFcFabricPol, ok := d.GetOk("relation_infra_rs_fc_fabric_pol"); ok {
		relationParam := relationToinfraRsFcFabricPol.(string)
		err = aciClient.CreateRelationinfraRsFcFabricPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationToinfraRsFcInstPol, ok := d.GetOk("relation_infra_rs_fc_inst_pol"); ok {
		relationParam := relationToinfraRsFcInstPol.(string)
		err = aciClient.CreateRelationinfraRsFcInstPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationToinfraRsIaclLeafProfile, ok := d.GetOk("relation_infra_rs_iacl_leaf_profile"); ok {
		relationParam := relationToinfraRsIaclLeafProfile.(string)
		err = aciClient.CreateRelationinfraRsIaclLeafProfile(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationToinfraRsL2NodeAuthPol, ok := d.GetOk("relation_infra_rs_l2_node_auth_pol"); ok {
		relationParam := relationToinfraRsL2NodeAuthPol.(string)
		err = aciClient.CreateRelationinfraRsL2NodeAuthPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationToinfraRsLeafCoppProfile, ok := d.GetOk("relation_infra_rs_leaf_copp_profile"); ok {
		relationParam := relationToinfraRsLeafCoppProfile.(string)
		err = aciClient.CreateRelationinfraRsLeafCoppProfile(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationToinfraRsLeafPGrpToCdpIfPol, ok := d.GetOk("relation_infra_rs_leaf_p_grp_to_cdp_if_pol"); ok {
		relationParam := relationToinfraRsLeafPGrpToCdpIfPol.(string)
		err = aciClient.CreateRelationinfraRsLeafPGrpToCdpIfPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationToinfraRsLeafPGrpToLldpIfPol, ok := d.GetOk("relation_infra_rs_leaf_p_grp_to_lldp_if_pol"); ok {
		relationParam := relationToinfraRsLeafPGrpToLldpIfPol.(string)
		err = aciClient.CreateRelationinfraRsLeafPGrpToLldpIfPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationToinfraRsMonNodeInfraPol, ok := d.GetOk("relation_infra_rs_mon_node_infra_pol"); ok {
		relationParam := relationToinfraRsMonNodeInfraPol.(string)
		err = aciClient.CreateRelationinfraRsMonNodeInfraPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationToinfraRsMstInstPol, ok := d.GetOk("relation_infra_rs_mst_inst_pol"); ok {
		relationParam := relationToinfraRsMstInstPol.(string)
		err = aciClient.CreateRelationinfraRsMstInstPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationToinfraRsNetflowNodePol, ok := d.GetOk("relation_infra_rs_netflow_node_pol"); ok {
		relationParam := relationToinfraRsNetflowNodePol.(string)
		err = aciClient.CreateRelationinfraRsNetflowNodePol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationToinfraRsPoeInstPol, ok := d.GetOk("relation_infra_rs_poe_inst_pol"); ok {
		relationParam := relationToinfraRsPoeInstPol.(string)
		err = aciClient.CreateRelationinfraRsPoeInstPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationToinfraRsTopoctrlFastLinkFailoverInstPol, ok := d.GetOk("relation_infra_rs_topoctrl_fast_link_failover_inst_pol"); ok {
		relationParam := relationToinfraRsTopoctrlFastLinkFailoverInstPol.(string)
		err = aciClient.CreateRelationinfraRsTopoctrlFastLinkFailoverInstPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationToinfraRsTopoctrlFwdScaleProfPol, ok := d.GetOk("relation_infra_rs_topoctrl_fwd_scale_prof_pol"); ok {
		relationParam := relationToinfraRsTopoctrlFwdScaleProfPol.(string)
		err = aciClient.CreateRelationinfraRsTopoctrlFwdScaleProfPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(infraAccNodePGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciAccessSwitchPolicyGroupRead(ctx, d, m)
}

func resourceAciAccessSwitchPolicyGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] AccessSwitchPolicyGroup: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	infraAccNodePGrpAttr := models.AccessSwitchPolicyGroupAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		infraAccNodePGrpAttr.Annotation = Annotation.(string)
	} else {
		infraAccNodePGrpAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		infraAccNodePGrpAttr.Name = Name.(string)
	}
	infraAccNodePGrp := models.NewAccessSwitchPolicyGroup(fmt.Sprintf("infra/funcprof/accnodepgrp-%s", name), "uni", desc, nameAlias, infraAccNodePGrpAttr)
	infraAccNodePGrp.Status = "modified"
	err := aciClient.Save(infraAccNodePGrp)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_infra_rs_bfd_ipv4_inst_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_bfd_ipv4_inst_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_infra_rs_bfd_ipv6_inst_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_bfd_ipv6_inst_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_infra_rs_bfd_mh_ipv4_inst_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_bfd_mh_ipv4_inst_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_infra_rs_bfd_mh_ipv6_inst_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_bfd_mh_ipv6_inst_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_infra_rs_equipment_flash_config_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_equipment_flash_config_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_infra_rs_fc_fabric_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_fc_fabric_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_infra_rs_fc_inst_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_fc_inst_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_infra_rs_iacl_leaf_profile") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_iacl_leaf_profile")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_infra_rs_l2_node_auth_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_l2_node_auth_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_infra_rs_leaf_copp_profile") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_leaf_copp_profile")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_infra_rs_leaf_p_grp_to_cdp_if_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_leaf_p_grp_to_cdp_if_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_infra_rs_leaf_p_grp_to_lldp_if_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_leaf_p_grp_to_lldp_if_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_infra_rs_mon_node_infra_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_mon_node_infra_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_infra_rs_mst_inst_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_mst_inst_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_infra_rs_netflow_node_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_netflow_node_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_infra_rs_poe_inst_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_poe_inst_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_infra_rs_topoctrl_fast_link_failover_inst_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_topoctrl_fast_link_failover_inst_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_infra_rs_topoctrl_fwd_scale_prof_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_topoctrl_fwd_scale_prof_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("relation_infra_rs_bfd_ipv4_inst_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_bfd_ipv4_inst_pol")
		err = aciClient.DeleteRelationinfraRsBfdIpv4InstPol(infraAccNodePGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsBfdIpv4InstPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_bfd_ipv6_inst_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_bfd_ipv6_inst_pol")
		err = aciClient.DeleteRelationinfraRsBfdIpv6InstPol(infraAccNodePGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsBfdIpv6InstPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_bfd_mh_ipv4_inst_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_bfd_mh_ipv4_inst_pol")
		err = aciClient.DeleteRelationinfraRsBfdMhIpv4InstPol(infraAccNodePGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsBfdMhIpv4InstPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_bfd_mh_ipv6_inst_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_bfd_mh_ipv6_inst_pol")
		err = aciClient.DeleteRelationinfraRsBfdMhIpv6InstPol(infraAccNodePGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsBfdMhIpv6InstPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_equipment_flash_config_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_equipment_flash_config_pol")
		err = aciClient.DeleteRelationinfraRsEquipmentFlashConfigPol(infraAccNodePGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsEquipmentFlashConfigPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_fc_fabric_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_fc_fabric_pol")
		err = aciClient.DeleteRelationinfraRsFcFabricPol(infraAccNodePGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsFcFabricPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_fc_inst_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_fc_inst_pol")
		err = aciClient.DeleteRelationinfraRsFcInstPol(infraAccNodePGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsFcInstPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_iacl_leaf_profile") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_iacl_leaf_profile")
		err = aciClient.DeleteRelationinfraRsIaclLeafProfile(infraAccNodePGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsIaclLeafProfile(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_l2_node_auth_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_l2_node_auth_pol")
		err = aciClient.DeleteRelationinfraRsL2NodeAuthPol(infraAccNodePGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsL2NodeAuthPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_leaf_copp_profile") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_leaf_copp_profile")
		err = aciClient.DeleteRelationinfraRsLeafCoppProfile(infraAccNodePGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsLeafCoppProfile(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_leaf_p_grp_to_cdp_if_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_leaf_p_grp_to_cdp_if_pol")
		err = aciClient.DeleteRelationinfraRsLeafPGrpToCdpIfPol(infraAccNodePGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsLeafPGrpToCdpIfPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_leaf_p_grp_to_lldp_if_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_leaf_p_grp_to_lldp_if_pol")
		err = aciClient.DeleteRelationinfraRsLeafPGrpToLldpIfPol(infraAccNodePGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsLeafPGrpToLldpIfPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_mon_node_infra_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_mon_node_infra_pol")
		err = aciClient.DeleteRelationinfraRsMonNodeInfraPol(infraAccNodePGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsMonNodeInfraPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_mst_inst_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_mst_inst_pol")
		err = aciClient.DeleteRelationinfraRsMstInstPol(infraAccNodePGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsMstInstPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_netflow_node_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_netflow_node_pol")
		err = aciClient.DeleteRelationinfraRsNetflowNodePol(infraAccNodePGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsNetflowNodePol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_poe_inst_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_poe_inst_pol")
		err = aciClient.DeleteRelationinfraRsPoeInstPol(infraAccNodePGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsPoeInstPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_topoctrl_fast_link_failover_inst_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_topoctrl_fast_link_failover_inst_pol")
		err = aciClient.DeleteRelationinfraRsTopoctrlFastLinkFailoverInstPol(infraAccNodePGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsTopoctrlFastLinkFailoverInstPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_topoctrl_fwd_scale_prof_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_topoctrl_fwd_scale_prof_pol")
		err = aciClient.DeleteRelationinfraRsTopoctrlFwdScaleProfPol(infraAccNodePGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsTopoctrlFwdScaleProfPol(infraAccNodePGrp.DistinguishedName, infraAccNodePGrpAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(infraAccNodePGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciAccessSwitchPolicyGroupRead(ctx, d, m)
}

func resourceAciAccessSwitchPolicyGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	infraAccNodePGrp, err := getRemoteAccessSwitchPolicyGroup(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}
	setAccessSwitchPolicyGroupAttributes(infraAccNodePGrp, d)

	infraRsBfdIpv4InstPolData, err := aciClient.ReadRelationinfraRsBfdIpv4InstPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsBfdIpv4InstPol %v", err)
		d.Set("relation_infra_rs_bfd_ipv4_inst_pol", "")
	} else {
		if _, ok := d.GetOk("relation_infra_rs_bfd_ipv4_inst_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_bfd_ipv4_inst_pol").(string))
			if tfName != infraRsBfdIpv4InstPolData {
				d.Set("relation_infra_rs_bfd_ipv4_inst_pol", "")
			}
		}
	}

	infraRsBfdIpv6InstPolData, err := aciClient.ReadRelationinfraRsBfdIpv6InstPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsBfdIpv6InstPol %v", err)
		d.Set("relation_infra_rs_bfd_ipv6_inst_pol", "")
	} else {
		if _, ok := d.GetOk("relation_infra_rs_bfd_ipv6_inst_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_bfd_ipv6_inst_pol").(string))
			if tfName != infraRsBfdIpv6InstPolData {
				d.Set("relation_infra_rs_bfd_ipv6_inst_pol", "")
			}
		}
	}

	infraRsBfdMhIpv4InstPolData, err := aciClient.ReadRelationinfraRsBfdMhIpv4InstPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsBfdMhIpv4InstPol %v", err)
		d.Set("relation_infra_rs_bfd_mh_ipv4_inst_pol", "")
	} else {
		if _, ok := d.GetOk("relation_infra_rs_bfd_mh_ipv4_inst_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_bfd_mh_ipv4_inst_pol").(string))
			if tfName != infraRsBfdMhIpv4InstPolData {
				d.Set("relation_infra_rs_bfd_mh_ipv4_inst_pol", "")
			}
		}
	}

	infraRsBfdMhIpv6InstPolData, err := aciClient.ReadRelationinfraRsBfdMhIpv6InstPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsBfdMhIpv6InstPol %v", err)
		d.Set("relation_infra_rs_bfd_mh_ipv6_inst_pol", "")
	} else {
		if _, ok := d.GetOk("relation_infra_rs_bfd_mh_ipv6_inst_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_bfd_mh_ipv6_inst_pol").(string))
			if tfName != infraRsBfdMhIpv6InstPolData {
				d.Set("relation_infra_rs_bfd_mh_ipv6_inst_pol", "")
			}
		}
	}

	infraRsEquipmentFlashConfigPolData, err := aciClient.ReadRelationinfraRsEquipmentFlashConfigPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsEquipmentFlashConfigPol %v", err)
		d.Set("relation_infra_rs_equipment_flash_config_pol", "")
	} else {
		if _, ok := d.GetOk("relation_infra_rs_equipment_flash_config_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_equipment_flash_config_pol").(string))
			if tfName != infraRsEquipmentFlashConfigPolData {
				d.Set("relation_infra_rs_equipment_flash_config_pol", "")
			}
		}
	}

	infraRsFcFabricPolData, err := aciClient.ReadRelationinfraRsFcFabricPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsFcFabricPol %v", err)
		d.Set("relation_infra_rs_fc_fabric_pol", "")
	} else {
		if _, ok := d.GetOk("relation_infra_rs_fc_fabric_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_fc_fabric_pol").(string))
			if tfName != infraRsFcFabricPolData {
				d.Set("relation_infra_rs_fc_fabric_pol", "")
			}
		}
	}

	infraRsFcInstPolData, err := aciClient.ReadRelationinfraRsFcInstPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsFcInstPol %v", err)
		d.Set("relation_infra_rs_fc_inst_pol", "")
	} else {
		if _, ok := d.GetOk("relation_infra_rs_fc_inst_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_fc_inst_pol").(string))
			if tfName != infraRsFcInstPolData {
				d.Set("relation_infra_rs_fc_inst_pol", "")
			}
		}
	}

	infraRsIaclLeafProfileData, err := aciClient.ReadRelationinfraRsIaclLeafProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsIaclLeafProfile %v", err)
		d.Set("relation_infra_rs_iacl_leaf_profile", "")
	} else {
		if _, ok := d.GetOk("relation_infra_rs_iacl_leaf_profile"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_iacl_leaf_profile").(string))
			if tfName != infraRsIaclLeafProfileData {
				d.Set("relation_infra_rs_iacl_leaf_profile", "")
			}
		}
	}

	infraRsL2NodeAuthPolData, err := aciClient.ReadRelationinfraRsL2NodeAuthPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsL2NodeAuthPol %v", err)
		d.Set("relation_infra_rs_l2_node_auth_pol", "")
	} else {
		if _, ok := d.GetOk("relation_infra_rs_l2_node_auth_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_l2_node_auth_pol").(string))
			if tfName != infraRsL2NodeAuthPolData {
				d.Set("relation_infra_rs_l2_node_auth_pol", "")
			}
		}
	}

	infraRsLeafCoppProfileData, err := aciClient.ReadRelationinfraRsLeafCoppProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsLeafCoppProfile %v", err)
		d.Set("relation_infra_rs_leaf_copp_profile", "")
	} else {
		if _, ok := d.GetOk("relation_infra_rs_leaf_copp_profile"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_leaf_copp_profile").(string))
			if tfName != infraRsLeafCoppProfileData {
				d.Set("relation_infra_rs_leaf_copp_profile", "")
			}
		}
	}

	infraRsLeafPGrpToCdpIfPolData, err := aciClient.ReadRelationinfraRsLeafPGrpToCdpIfPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsLeafPGrpToCdpIfPol %v", err)
		d.Set("relation_infra_rs_leaf_p_grp_to_cdp_if_pol", "")
	} else {
		if _, ok := d.GetOk("relation_infra_rs_leaf_p_grp_to_cdp_if_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_leaf_p_grp_to_cdp_if_pol").(string))
			if tfName != infraRsLeafPGrpToCdpIfPolData {
				d.Set("relation_infra_rs_leaf_p_grp_to_cdp_if_pol", "")
			}
		}
	}

	infraRsLeafPGrpToLldpIfPolData, err := aciClient.ReadRelationinfraRsLeafPGrpToLldpIfPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsLeafPGrpToLldpIfPol %v", err)
		d.Set("relation_infra_rs_leaf_p_grp_to_lldp_if_pol", "")
	} else {
		if _, ok := d.GetOk("relation_infra_rs_leaf_p_grp_to_lldp_if_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_leaf_p_grp_to_lldp_if_pol").(string))
			if tfName != infraRsLeafPGrpToLldpIfPolData {
				d.Set("relation_infra_rs_leaf_p_grp_to_lldp_if_pol", "")
			}
		}
	}

	infraRsMonNodeInfraPolData, err := aciClient.ReadRelationinfraRsMonNodeInfraPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsMonNodeInfraPol %v", err)
		d.Set("relation_infra_rs_mon_node_infra_pol", "")
	} else {
		if _, ok := d.GetOk("relation_infra_rs_mon_node_infra_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_mon_node_infra_pol").(string))
			if tfName != infraRsMonNodeInfraPolData {
				d.Set("relation_infra_rs_mon_node_infra_pol", "")
			}
		}
	}

	infraRsMstInstPolData, err := aciClient.ReadRelationinfraRsMstInstPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsMstInstPol %v", err)
		d.Set("relation_infra_rs_mst_inst_pol", "")
	} else {
		if _, ok := d.GetOk("relation_infra_rs_mst_inst_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_mst_inst_pol").(string))
			if tfName != infraRsMstInstPolData {
				d.Set("relation_infra_rs_mst_inst_pol", "")
			}
		}
	}

	infraRsNetflowNodePolData, err := aciClient.ReadRelationinfraRsNetflowNodePol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsNetflowNodePol %v", err)
		d.Set("relation_infra_rs_netflow_node_pol", "")
	} else {
		if _, ok := d.GetOk("relation_infra_rs_netflow_node_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_netflow_node_pol").(string))
			if tfName != infraRsNetflowNodePolData {
				d.Set("relation_infra_rs_netflow_node_pol", "")
			}
		}
	}

	infraRsPoeInstPolData, err := aciClient.ReadRelationinfraRsPoeInstPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsPoeInstPol %v", err)
		d.Set("relation_infra_rs_poe_inst_pol", "")
	} else {
		if _, ok := d.GetOk("relation_infra_rs_poe_inst_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_poe_inst_pol").(string))
			if tfName != infraRsPoeInstPolData {
				d.Set("relation_infra_rs_poe_inst_pol", "")
			}
		}
	}

	infraRsTopoctrlFastLinkFailoverInstPolData, err := aciClient.ReadRelationinfraRsTopoctrlFastLinkFailoverInstPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsTopoctrlFastLinkFailoverInstPol %v", err)
		d.Set("relation_infra_rs_topoctrl_fast_link_failover_inst_pol", "")
	} else {
		if _, ok := d.GetOk("relation_infra_rs_topoctrl_fast_link_failover_inst_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_topoctrl_fast_link_failover_inst_pol").(string))
			if tfName != infraRsTopoctrlFastLinkFailoverInstPolData {
				d.Set("relation_infra_rs_topoctrl_fast_link_failover_inst_pol", "")
			}
		}
	}

	infraRsTopoctrlFwdScaleProfPolData, err := aciClient.ReadRelationinfraRsTopoctrlFwdScaleProfPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsTopoctrlFwdScaleProfPol %v", err)
		d.Set("relation_infra_rs_topoctrl_fwd_scale_prof_pol", "")
	} else {
		if _, ok := d.GetOk("relation_infra_rs_topoctrl_fwd_scale_prof_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_topoctrl_fwd_scale_prof_pol").(string))
			if tfName != infraRsTopoctrlFwdScaleProfPolData {
				d.Set("relation_infra_rs_topoctrl_fwd_scale_prof_pol", "")
			}
		}
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciAccessSwitchPolicyGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraAccNodePGrp")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
