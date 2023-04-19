package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciAccessSwitchPolicyGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciAccessSwitchPolicyGroupRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"relation_infra_rs_bfd_ipv4_inst_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_bfd_ipv6_inst_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_bfd_mh_ipv4_inst_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_bfd_mh_ipv6_inst_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_equipment_flash_config_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_fc_fabric_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_fc_inst_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_iacl_leaf_profile": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_l2_node_auth_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_leaf_copp_profile": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_leaf_p_grp_to_cdp_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_leaf_p_grp_to_lldp_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_mon_node_infra_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_mst_inst_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_netflow_node_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_poe_inst_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_topoctrl_fast_link_failover_inst_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_topoctrl_fwd_scale_prof_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		})),
	}
}

func dataSourceAciAccessSwitchPolicyGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)

	rn := fmt.Sprintf("infra/funcprof/accnodepgrp-%s", name)
	dn := fmt.Sprintf("uni/%s", rn)
	log.Printf("[DEBUG] %s: Data Source - Beginning Read", dn)
	infraAccNodePGrp, err := getRemoteAccessSwitchPolicyGroup(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setAccessSwitchPolicyGroupAttributes(infraAccNodePGrp, d)
	if err != nil {
		return diag.FromErr(err)
	}

	// infraRsBfdIpv4InstPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsBfdIpv4InstPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsBfdIpv4InstPol(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsBfdIpv4InstPol - Read finished successfully", d.Get("relation_infra_rs_bfd_ipv4_inst_pol"))
	}
	// infraRsBfdIpv4InstPol - Read finished successfully

	// infraRsBfdIpv6InstPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsBfdIpv6InstPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsBfdIpv6InstPol(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsBfdIpv6InstPol - Read finished successfully", d.Get("relation_infra_rs_bfd_ipv6_inst_pol"))
	}
	// infraRsBfdIpv6InstPol - Read finished successfully

	// infraRsBfdMhIpv4InstPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsBfdMhIpv4InstPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsBfdMhIpv4InstPol(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsBfdMhIpv4InstPol - Read finished successfully", d.Get("relation_infra_rs_bfd_mh_ipv4_inst_pol"))
	}
	// infraRsBfdMhIpv4InstPol - Read finished successfully

	// infraRsBfdMhIpv6InstPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsBfdMhIpv6InstPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsBfdMhIpv6InstPol(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsBfdMhIpv6InstPol - Read finished successfully", d.Get("relation_infra_rs_bfd_mh_ipv6_inst_pol"))
	}
	// infraRsBfdMhIpv6InstPol - Read finished successfully

	// infraRsEquipmentFlashConfigPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsEquipmentFlashConfigPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsEquipmentFlashConfigPol(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsEquipmentFlashConfigPol - Read finished successfully", d.Get("relation_infra_rs_equipment_flash_config_pol"))
	}
	// infraRsEquipmentFlashConfigPol - Read finished successfully

	// infraRsFcFabricPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsFcFabricPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsFcFabricPol(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsFcFabricPol - Read finished successfully", d.Get("relation_infra_rs_fc_fabric_pol"))
	}
	// infraRsFcFabricPol - Read finished successfully

	// infraRsFcInstPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsFcInstPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsFcInstPol(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsFcInstPol - Read finished successfully", d.Get("relation_infra_rs_fc_inst_pol"))
	}
	// infraRsFcInstPol - Read finished successfully

	// infraRsIaclLeafProfile - Beginning Read
	log.Printf("[DEBUG] %s: infraRsIaclLeafProfile - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsIaclLeafProfile(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsIaclLeafProfile - Read finished successfully", d.Get("relation_infra_rs_iacl_leaf_profile"))
	}
	// infraRsIaclLeafProfile - Read finished successfully

	// infraRsL2NodeAuthPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsL2NodeAuthPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsL2NodeAuthPol(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsL2NodeAuthPol - Read finished successfully", d.Get("relation_infra_rs_l2_node_auth_pol"))
	}
	// infraRsL2NodeAuthPol - Read finished successfully

	// infraRsLeafCoppProfile - Beginning Read
	log.Printf("[DEBUG] %s: infraRsLeafCoppProfile - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsLeafCoppProfile(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsLeafCoppProfile - Read finished successfully", d.Get("relation_infra_rs_leaf_copp_profile"))
	}
	// infraRsLeafCoppProfile - Read finished successfully

	// infraRsLeafPGrpToCdpIfPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsLeafPGrpToCdpIfPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsLeafPGrpToCdpIfPol(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsLeafPGrpToCdpIfPol - Read finished successfully", d.Get("relation_infra_rs_leaf_p_grp_to_cdp_if_pol"))
	}
	// infraRsLeafPGrpToCdpIfPol - Read finished successfully

	// infraRsLeafPGrpToLldpIfPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsLeafPGrpToLldpIfPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsLeafPGrpToLldpIfPol(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsLeafPGrpToLldpIfPol - Read finished successfully", d.Get("relation_infra_rs_leaf_p_grp_to_lldp_if_pol"))
	}
	// infraRsLeafPGrpToLldpIfPol - Read finished successfully

	// infraRsMonNodeInfraPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsMonNodeInfraPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsMonNodeInfraPol(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsMonNodeInfraPol - Read finished successfully", d.Get("relation_infra_rs_mon_node_infra_pol"))
	}
	// infraRsMonNodeInfraPol - Read finished successfully

	// infraRsMstInstPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsMstInstPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsMstInstPol(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsMstInstPol - Read finished successfully", d.Get("relation_infra_rs_mst_inst_pol"))
	}
	// infraRsMstInstPol - Read finished successfully

	// infraRsNetflowNodePol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsNetflowNodePol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsNetflowNodePol(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsNetflowNodePol - Read finished successfully", d.Get("relation_infra_rs_netflow_node_pol"))
	}
	// infraRsNetflowNodePol - Read finished successfully

	// infraRsPoeInstPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsPoeInstPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsPoeInstPol(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsPoeInstPol - Read finished successfully", d.Get("relation_infra_rs_poe_inst_pol"))
	}
	// infraRsPoeInstPol - Read finished successfully

	// infraRsTopoctrlFastLinkFailoverInstPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsTopoctrlFastLinkFailoverInstPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsTopoctrlFastLinkFailoverInstPol(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsTopoctrlFastLinkFailoverInstPol - Read finished successfully", d.Get("relation_infra_rs_topoctrl_fast_link_failover_inst_pol"))
	}
	// infraRsTopoctrlFastLinkFailoverInstPol - Read finished successfully

	// infraRsTopoctrlFwdScaleProfPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsTopoctrlFwdScaleProfPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsTopoctrlFwdScaleProfPol(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsTopoctrlFwdScaleProfPol - Read finished successfully", d.Get("relation_infra_rs_topoctrl_fwd_scale_prof_pol"))
	}
	// infraRsTopoctrlFwdScaleProfPol - Read finished successfully

	log.Printf("[DEBUG] %s: Data Source - Read finished successfully", dn)
	return nil
}
