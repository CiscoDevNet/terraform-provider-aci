package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciSpineSwitchPolicyGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciSpineSwitchPolicyGroupRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"relation_infra_rs_iacl_spine_profile": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_spine_bfd_ipv4_inst_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_spine_bfd_ipv6_inst_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_spine_copp_profile": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_spine_p_grp_to_cdp_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_spine_p_grp_to_lldp_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		})),
	}
}

func dataSourceAciSpineSwitchPolicyGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)

	rn := fmt.Sprintf("infra/funcprof/spaccnodepgrp-%s", name)
	dn := fmt.Sprintf("uni/%s", rn)
	log.Printf("[DEBUG] %s: Data Source - Beginning Read", dn)
	infraSpineAccNodePGrp, err := getRemoteSpineSwitchPolicyGroup(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setSpineSwitchPolicyGroupAttributes(infraSpineAccNodePGrp, d)
	if err != nil {
		return diag.FromErr(err)
	}

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

	log.Printf("[DEBUG] %s: Data Source - Read finished successfully", dn)
	return nil
}
