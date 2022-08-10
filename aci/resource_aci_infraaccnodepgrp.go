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

func getAndSetReadRelationinfraRsBfdIpv4InstPol(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsBfdIpv4InstPolData, err := client.ReadRelationinfraRsBfdIpv4InstPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsBfdIpv4InstPol %v", err)
		d.Set("relation_infra_rs_bfd_ipv4_inst_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_bfd_ipv4_inst_pol", infraRsBfdIpv4InstPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsBfdIpv6InstPol(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsBfdIpv6InstPolData, err := client.ReadRelationinfraRsBfdIpv6InstPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsBfdIpv6InstPol %v", err)
		d.Set("relation_infra_rs_bfd_ipv6_inst_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_bfd_ipv6_inst_pol", infraRsBfdIpv6InstPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsBfdMhIpv4InstPol(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsBfdMhIpv4InstPolData, err := client.ReadRelationinfraRsBfdMhIpv4InstPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsBfdMhIpv4InstPol %v", err)
		d.Set("relation_infra_rs_bfd_mh_ipv4_inst_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_bfd_mh_ipv4_inst_pol", infraRsBfdMhIpv4InstPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsBfdMhIpv6InstPol(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsBfdMhIpv6InstPolData, err := client.ReadRelationinfraRsBfdMhIpv6InstPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsBfdMhIpv6InstPol %v", err)
		d.Set("relation_infra_rs_bfd_mh_ipv6_inst_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_bfd_mh_ipv6_inst_pol", infraRsBfdMhIpv6InstPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsEquipmentFlashConfigPol(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsEquipmentFlashConfigPolData, err := client.ReadRelationinfraRsEquipmentFlashConfigPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsEquipmentFlashConfigPol %v", err)
		d.Set("relation_infra_rs_equipment_flash_config_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_equipment_flash_config_pol", infraRsEquipmentFlashConfigPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsFcFabricPol(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsFcFabricPolData, err := client.ReadRelationinfraRsFcFabricPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsFcFabricPol %v", err)
		d.Set("relation_infra_rs_fc_fabric_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_fc_fabric_pol", infraRsFcFabricPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsFcInstPol(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsFcInstPolData, err := client.ReadRelationinfraRsFcInstPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsFcInstPol %v", err)
		d.Set("relation_infra_rs_fc_inst_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_fc_inst_pol", infraRsFcInstPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsIaclLeafProfile(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsIaclLeafProfileData, err := client.ReadRelationinfraRsIaclLeafProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsIaclLeafProfile %v", err)
		d.Set("relation_infra_rs_iacl_leaf_profile", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_iacl_leaf_profile", infraRsIaclLeafProfileData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsL2NodeAuthPol(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsL2NodeAuthPolData, err := client.ReadRelationinfraRsL2NodeAuthPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsL2NodeAuthPol %v", err)
		d.Set("relation_infra_rs_l2_node_auth_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_l2_node_auth_pol", infraRsL2NodeAuthPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsLeafCoppProfile(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsLeafCoppProfileData, err := client.ReadRelationinfraRsLeafCoppProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsLeafCoppProfile %v", err)
		d.Set("relation_infra_rs_leaf_copp_profile", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_leaf_copp_profile", infraRsLeafCoppProfileData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsLeafPGrpToCdpIfPol(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsLeafPGrpToCdpIfPolData, err := client.ReadRelationinfraRsLeafPGrpToCdpIfPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsLeafPGrpToCdpIfPol %v", err)
		d.Set("relation_infra_rs_leaf_p_grp_to_cdp_if_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_leaf_p_grp_to_cdp_if_pol", infraRsLeafPGrpToCdpIfPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsLeafPGrpToLldpIfPol(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsLeafPGrpToLldpIfPolData, err := client.ReadRelationinfraRsLeafPGrpToLldpIfPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsLeafPGrpToLldpIfPol %v", err)
		d.Set("relation_infra_rs_leaf_p_grp_to_lldp_if_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_leaf_p_grp_to_lldp_if_pol", infraRsLeafPGrpToLldpIfPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsMonNodeInfraPol(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsMonNodeInfraPolData, err := client.ReadRelationinfraRsMonNodeInfraPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsMonNodeInfraPol %v", err)
		d.Set("relation_infra_rs_mon_node_infra_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_mon_node_infra_pol", infraRsMonNodeInfraPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsMstInstPol(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsMstInstPolData, err := client.ReadRelationinfraRsMstInstPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsMstInstPol %v", err)
		d.Set("relation_infra_rs_mst_inst_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_mst_inst_pol", infraRsMstInstPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsNetflowNodePol(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsNetflowNodePolData, err := client.ReadRelationinfraRsNetflowNodePol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsNetflowNodePol %v", err)
		d.Set("relation_infra_rs_netflow_node_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_netflow_node_pol", infraRsNetflowNodePolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsPoeInstPol(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsPoeInstPolData, err := client.ReadRelationinfraRsPoeInstPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsPoeInstPol %v", err)
		d.Set("relation_infra_rs_poe_inst_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_poe_inst_pol", infraRsPoeInstPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsTopoctrlFastLinkFailoverInstPol(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsTopoctrlFastLinkFailoverInstPolData, err := client.ReadRelationinfraRsTopoctrlFastLinkFailoverInstPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsTopoctrlFastLinkFailoverInstPol %v", err)
		d.Set("relation_infra_rs_topoctrl_fast_link_failover_inst_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_topoctrl_fast_link_failover_inst_pol", infraRsTopoctrlFastLinkFailoverInstPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsTopoctrlFwdScaleProfPol(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsTopoctrlFwdScaleProfPolData, err := client.ReadRelationinfraRsTopoctrlFwdScaleProfPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsTopoctrlFwdScaleProfPol %v", err)
		d.Set("relation_infra_rs_topoctrl_fwd_scale_prof_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_topoctrl_fwd_scale_prof_pol", infraRsTopoctrlFwdScaleProfPolData.(string))
	}
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

	// infraRsBfdIpv4InstPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsBfdIpv4InstPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsBfdIpv4InstPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsBfdIpv4InstPol - Import finished successfully", d.Get("relation_infra_rs_bfd_ipv4_inst_pol"))
	}
	// infraRsBfdIpv4InstPol - Import finished successfully

	// infraRsBfdIpv6InstPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsBfdIpv6InstPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsBfdIpv6InstPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsBfdIpv6InstPol - Import finished successfully", d.Get("relation_infra_rs_bfd_ipv6_inst_pol"))
	}
	// infraRsBfdIpv6InstPol - Import finished successfully

	// infraRsBfdMhIpv4InstPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsBfdMhIpv4InstPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsBfdMhIpv4InstPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsBfdMhIpv4InstPol - Import finished successfully", d.Get("relation_infra_rs_bfd_mh_ipv4_inst_pol"))
	}
	// infraRsBfdMhIpv4InstPol - Import finished successfully

	// infraRsBfdMhIpv6InstPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsBfdMhIpv6InstPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsBfdMhIpv6InstPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsBfdMhIpv6InstPol - Import finished successfully", d.Get("relation_infra_rs_bfd_mh_ipv6_inst_pol"))
	}
	// infraRsBfdMhIpv6InstPol - Import finished successfully

	// infraRsEquipmentFlashConfigPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsEquipmentFlashConfigPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsEquipmentFlashConfigPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsEquipmentFlashConfigPol - Import finished successfully", d.Get("relation_infra_rs_equipment_flash_config_pol"))
	}
	// infraRsEquipmentFlashConfigPol - Import finished successfully

	// infraRsFcFabricPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsFcFabricPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsFcFabricPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsFcFabricPol - Import finished successfully", d.Get("relation_infra_rs_fc_fabric_pol"))
	}
	// infraRsFcFabricPol - Import finished successfully

	// infraRsFcInstPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsFcInstPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsFcInstPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsFcInstPol - Import finished successfully", d.Get("relation_infra_rs_fc_inst_pol"))
	}
	// infraRsFcInstPol - Import finished successfully

	// infraRsIaclLeafProfile - Beginning Import
	log.Printf("[DEBUG] %s: infraRsIaclLeafProfile - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsIaclLeafProfile(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsIaclLeafProfile - Import finished successfully", d.Get("relation_infra_rs_iacl_leaf_profile"))
	}
	// infraRsIaclLeafProfile - Import finished successfully

	// infraRsL2NodeAuthPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsL2NodeAuthPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsL2NodeAuthPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsL2NodeAuthPol - Import finished successfully", d.Get("relation_infra_rs_l2_node_auth_pol"))
	}
	// infraRsL2NodeAuthPol - Import finished successfully

	// infraRsLeafCoppProfile - Beginning Import
	log.Printf("[DEBUG] %s: infraRsLeafCoppProfile - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsLeafCoppProfile(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsLeafCoppProfile - Import finished successfully", d.Get("relation_infra_rs_leaf_copp_profile"))
	}
	// infraRsLeafCoppProfile - Import finished successfully

	// infraRsLeafPGrpToCdpIfPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsLeafPGrpToCdpIfPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsLeafPGrpToCdpIfPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsLeafPGrpToCdpIfPol - Import finished successfully", d.Get("relation_infra_rs_leaf_p_grp_to_cdp_if_pol"))
	}
	// infraRsLeafPGrpToCdpIfPol - Import finished successfully

	// infraRsLeafPGrpToLldpIfPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsLeafPGrpToLldpIfPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsLeafPGrpToLldpIfPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsLeafPGrpToLldpIfPol - Import finished successfully", d.Get("relation_infra_rs_leaf_p_grp_to_lldp_if_pol"))
	}
	// infraRsLeafPGrpToLldpIfPol - Import finished successfully

	// infraRsMonNodeInfraPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsMonNodeInfraPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsMonNodeInfraPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsMonNodeInfraPol - Import finished successfully", d.Get("relation_infra_rs_mon_node_infra_pol"))
	}
	// infraRsMonNodeInfraPol - Import finished successfully

	// infraRsMstInstPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsMstInstPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsMstInstPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsMstInstPol - Import finished successfully", d.Get("relation_infra_rs_mst_inst_pol"))
	}
	// infraRsMstInstPol - Import finished successfully

	// infraRsNetflowNodePol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsNetflowNodePol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsNetflowNodePol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsNetflowNodePol - Import finished successfully", d.Get("relation_infra_rs_netflow_node_pol"))
	}
	// infraRsNetflowNodePol - Import finished successfully

	// infraRsPoeInstPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsPoeInstPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsPoeInstPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsPoeInstPol - Import finished successfully", d.Get("relation_infra_rs_poe_inst_pol"))
	}
	// infraRsPoeInstPol - Import finished successfully

	// infraRsTopoctrlFastLinkFailoverInstPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsTopoctrlFastLinkFailoverInstPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsTopoctrlFastLinkFailoverInstPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsTopoctrlFastLinkFailoverInstPol - Import finished successfully", d.Get("relation_infra_rs_topoctrl_fast_link_failover_inst_pol"))
	}
	// infraRsTopoctrlFastLinkFailoverInstPol - Import finished successfully

	// infraRsTopoctrlFwdScaleProfPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsTopoctrlFwdScaleProfPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsTopoctrlFwdScaleProfPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsTopoctrlFwdScaleProfPol - Import finished successfully", d.Get("relation_infra_rs_topoctrl_fwd_scale_prof_pol"))
	}
	// infraRsTopoctrlFwdScaleProfPol - Import finished successfully

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

	// infraRsBfdIpv4InstPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsBfdIpv4InstPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsBfdIpv4InstPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsBfdIpv4InstPol - Read finished successfully", d.Get("relation_infra_rs_bfd_ipv4_inst_pol"))
	}
	// infraRsBfdIpv4InstPol - Read finished successfully

	// infraRsBfdIpv6InstPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsBfdIpv6InstPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsBfdIpv6InstPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsBfdIpv6InstPol - Read finished successfully", d.Get("relation_infra_rs_bfd_ipv6_inst_pol"))
	}
	// infraRsBfdIpv6InstPol - Read finished successfully

	// infraRsBfdMhIpv4InstPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsBfdMhIpv4InstPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsBfdMhIpv4InstPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsBfdMhIpv4InstPol - Read finished successfully", d.Get("relation_infra_rs_bfd_mh_ipv4_inst_pol"))
	}
	// infraRsBfdMhIpv4InstPol - Read finished successfully

	// infraRsBfdMhIpv6InstPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsBfdMhIpv6InstPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsBfdMhIpv6InstPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsBfdMhIpv6InstPol - Read finished successfully", d.Get("relation_infra_rs_bfd_mh_ipv6_inst_pol"))
	}
	// infraRsBfdMhIpv6InstPol - Read finished successfully

	// infraRsEquipmentFlashConfigPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsEquipmentFlashConfigPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsEquipmentFlashConfigPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsEquipmentFlashConfigPol - Read finished successfully", d.Get("relation_infra_rs_equipment_flash_config_pol"))
	}
	// infraRsEquipmentFlashConfigPol - Read finished successfully

	// infraRsFcFabricPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsFcFabricPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsFcFabricPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsFcFabricPol - Read finished successfully", d.Get("relation_infra_rs_fc_fabric_pol"))
	}
	// infraRsFcFabricPol - Read finished successfully

	// infraRsFcInstPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsFcInstPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsFcInstPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsFcInstPol - Read finished successfully", d.Get("relation_infra_rs_fc_inst_pol"))
	}
	// infraRsFcInstPol - Read finished successfully

	// infraRsIaclLeafProfile - Beginning Read
	log.Printf("[DEBUG] %s: infraRsIaclLeafProfile - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsIaclLeafProfile(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsIaclLeafProfile - Read finished successfully", d.Get("relation_infra_rs_iacl_leaf_profile"))
	}
	// infraRsIaclLeafProfile - Read finished successfully

	// infraRsL2NodeAuthPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsL2NodeAuthPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsL2NodeAuthPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsL2NodeAuthPol - Read finished successfully", d.Get("relation_infra_rs_l2_node_auth_pol"))
	}
	// infraRsL2NodeAuthPol - Read finished successfully

	// infraRsLeafCoppProfile - Beginning Read
	log.Printf("[DEBUG] %s: infraRsLeafCoppProfile - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsLeafCoppProfile(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsLeafCoppProfile - Read finished successfully", d.Get("relation_infra_rs_leaf_copp_profile"))
	}
	// infraRsLeafCoppProfile - Read finished successfully

	// infraRsLeafPGrpToCdpIfPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsLeafPGrpToCdpIfPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsLeafPGrpToCdpIfPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsLeafPGrpToCdpIfPol - Read finished successfully", d.Get("relation_infra_rs_leaf_p_grp_to_cdp_if_pol"))
	}
	// infraRsLeafPGrpToCdpIfPol - Read finished successfully

	// infraRsLeafPGrpToLldpIfPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsLeafPGrpToLldpIfPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsLeafPGrpToLldpIfPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsLeafPGrpToLldpIfPol - Read finished successfully", d.Get("relation_infra_rs_leaf_p_grp_to_lldp_if_pol"))
	}
	// infraRsLeafPGrpToLldpIfPol - Read finished successfully

	// infraRsMonNodeInfraPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsMonNodeInfraPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsMonNodeInfraPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsMonNodeInfraPol - Read finished successfully", d.Get("relation_infra_rs_mon_node_infra_pol"))
	}
	// infraRsMonNodeInfraPol - Read finished successfully

	// infraRsMstInstPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsMstInstPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsMstInstPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsMstInstPol - Read finished successfully", d.Get("relation_infra_rs_mst_inst_pol"))
	}
	// infraRsMstInstPol - Read finished successfully

	// infraRsNetflowNodePol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsNetflowNodePol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsNetflowNodePol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsNetflowNodePol - Read finished successfully", d.Get("relation_infra_rs_netflow_node_pol"))
	}
	// infraRsNetflowNodePol - Read finished successfully

	// infraRsPoeInstPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsPoeInstPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsPoeInstPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsPoeInstPol - Read finished successfully", d.Get("relation_infra_rs_poe_inst_pol"))
	}
	// infraRsPoeInstPol - Read finished successfully

	// infraRsTopoctrlFastLinkFailoverInstPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsTopoctrlFastLinkFailoverInstPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsTopoctrlFastLinkFailoverInstPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsTopoctrlFastLinkFailoverInstPol - Read finished successfully", d.Get("relation_infra_rs_topoctrl_fast_link_failover_inst_pol"))
	}
	// infraRsTopoctrlFastLinkFailoverInstPol - Read finished successfully

	// infraRsTopoctrlFwdScaleProfPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsTopoctrlFwdScaleProfPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsTopoctrlFwdScaleProfPol(aciClient, dn, d)
	if err != nil {
		log.Printf("[DEBUG] %s: infraRsTopoctrlFwdScaleProfPol - Read finished successfully", d.Get("relation_infra_rs_topoctrl_fwd_scale_prof_pol"))
	}
	// infraRsTopoctrlFwdScaleProfPol - Read finished successfully

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
