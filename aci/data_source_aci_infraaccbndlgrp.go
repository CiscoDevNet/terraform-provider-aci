package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciPCVPCInterfacePolicyGroup() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciPCVPCInterfacePolicyGroupRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"lag_t": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"relation_infra_rs_span_v_src_grp": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_infra_rs_stormctrl_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_lldp_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_macsec_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_qos_dpp_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_h_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_netflow_monitor_pol": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_dn": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"flt_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"relation_infra_rs_l2_port_auth_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_mcp_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_l2_port_security_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_copp_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_span_v_dest_grp": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_infra_rs_lacp_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_cdp_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_qos_pfc_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_qos_sd_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_mon_if_infra_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_fc_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_qos_ingress_dpp_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_qos_egress_dpp_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_l2_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_stp_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_att_ent_p": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_infra_rs_l2_inst_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		}),
	}
}

func dataSourceAciPCVPCInterfacePolicyGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("infra/funcprof/accbundle-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)
	log.Printf("[DEBUG] %s: Data Source - Beginning Read", dn)
	infraAccBndlGrp, err := getRemotePCVPCInterfacePolicyGroup(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setPCVPCInterfacePolicyGroupAttributes(infraAccBndlGrp, d)
	if err != nil {
		return diag.FromErr(err)
	}

	// infraRsSpanVSrcGrp - Beginning Read
	log.Printf("[DEBUG] %s: infraRsSpanVSrcGrp - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsSpanVSrcGrpFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsSpanVSrcGrp - Read finished successfully", d.Get("relation_infra_rs_span_v_src_grp"))
	}
	// infraRsSpanVSrcGrp - Read finished successfully

	// infraRsStormctrlIfPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsStormctrlIfPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsStormctrlIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsStormctrlIfPol - Read finished successfully", d.Get("relation_infra_rs_stormctrl_if_pol"))
	}
	// infraRsStormctrlIfPol - Read finished successfully

	// infraRsLldpIfPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsLldpIfPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsLldpIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsLldpIfPol - Read finished successfully", d.Get("relation_infra_rs_lldp_if_pol"))
	}
	// infraRsLldpIfPol - Read finished successfully

	// infraRsMacsecIfPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsMacsecIfPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsMacsecIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsMacsecIfPol - Read finished successfully", d.Get("relation_infra_rs_macsec_if_pol"))
	}
	// infraRsMacsecIfPol - Read finished successfully

	// infraRsQosDppIfPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsQosDppIfPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsQosDppIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsQosDppIfPol - Read finished successfully", d.Get("relation_infra_rs_qos_dpp_if_pol"))
	}
	// infraRsQosDppIfPol - Read finished successfully

	// infraRsHIfPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsHIfPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsHIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsHIfPol - Read finished successfully", d.Get("relation_infra_rs_h_if_pol"))
	}
	// infraRsHIfPol - Read finished successfully

	// infraRsNetflowMonitorPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsNetflowMonitorPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsNetflowMonitorPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsNetflowMonitorPol - Read finished successfully", d.Get("relation_infra_rs_netflow_monitor_pol"))
	}
	// infraRsNetflowMonitorPol - Read finished successfully

	// infraRsL2PortAuthPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsL2PortAuthPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsL2PortAuthPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsL2PortAuthPol - Read finished successfully", d.Get("relation_infra_rs_l2_port_auth_pol"))
	}
	// infraRsL2PortAuthPol - Read finished successfully

	// infraRsMcpIfPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsMcpIfPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsMcpIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsMcpIfPol - Read finished successfully", d.Get("relation_infra_rs_mcp_if_pol"))
	}
	// infraRsMcpIfPol - Read finished successfully

	// infraRsL2PortSecurityPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsL2PortSecurityPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsL2PortSecurityPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsL2PortSecurityPol - Read finished successfully", d.Get("relation_infra_rs_l2_port_security_pol"))
	}
	// infraRsL2PortSecurityPol - Read finished successfully

	// infraRsCoppIfPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsCoppIfPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsCoppIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsCoppIfPol - Read finished successfully", d.Get("relation_infra_rs_copp_if_pol"))
	}
	// infraRsCoppIfPol - Read finished successfully

	// infraRsSpanVDestGrp - Beginning Read
	log.Printf("[DEBUG] %s: infraRsSpanVDestGrp - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsSpanVDestGrpFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsSpanVDestGrp - Read finished successfully", d.Get("relation_infra_rs_span_v_dest_grp"))
	}
	// infraRsSpanVDestGrp - Read finished successfully

	// infraRsLacpPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsLacpPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsLacpPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsLacpPol - Read finished successfully", d.Get("relation_infra_rs_lacp_pol"))
	}
	// infraRsLacpPol - Read finished successfully

	// infraRsCdpIfPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsCdpIfPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsCdpIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsCdpIfPol - Read finished successfully", d.Get("relation_infra_rs_cdp_if_pol"))
	}
	// infraRsCdpIfPol - Read finished successfully

	// infraRsQosPfcIfPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsQosPfcIfPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsQosPfcIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsQosPfcIfPol - Read finished successfully", d.Get("relation_infra_rs_qos_pfc_if_pol"))
	}
	// infraRsQosPfcIfPol - Read finished successfully

	// infraRsQosSdIfPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsQosSdIfPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsQosSdIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsQosSdIfPol - Read finished successfully", d.Get("relation_infra_rs_qos_sd_if_pol"))
	}
	// infraRsQosSdIfPol - Read finished successfully

	// infraRsMonIfInfraPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsMonIfInfraPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsMonIfInfraPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsMonIfInfraPol - Read finished successfully", d.Get("relation_infra_rs_mon_if_infra_pol"))
	}
	// infraRsMonIfInfraPol - Read finished successfully

	// infraRsFcIfPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsFcIfPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsFcIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsFcIfPol - Read finished successfully", d.Get("relation_infra_rs_fc_if_pol"))
	}
	// infraRsFcIfPol - Read finished successfully

	// infraRsQosIngressDppIfPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsQosIngressDppIfPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsQosIngressDppIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsQosIngressDppIfPol - Read finished successfully", d.Get("relation_infra_rs_qos_ingress_dpp_if_pol"))
	}
	// infraRsQosIngressDppIfPol - Read finished successfully

	// infraRsQosEgressDppIfPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsQosEgressDppIfPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsQosEgressDppIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsQosEgressDppIfPol - Read finished successfully", d.Get("relation_infra_rs_qos_egress_dpp_if_pol"))
	}
	// infraRsQosEgressDppIfPol - Read finished successfully

	// infraRsL2IfPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsL2IfPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsL2IfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsL2IfPol - Read finished successfully", d.Get("relation_infra_rs_l2_if_pol"))
	}
	// infraRsL2IfPol - Read finished successfully

	// infraRsStpIfPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsStpIfPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsStpIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsStpIfPol - Read finished successfully", d.Get("relation_infra_rs_stp_if_pol"))
	}
	// infraRsStpIfPol - Read finished successfully

	// infraRsAttEntP - Beginning Read
	log.Printf("[DEBUG] %s: infraRsAttEntP - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsAttEntPFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsAttEntP - Read finished successfully", d.Get("relation_infra_rs_att_ent_p"))
	}
	// infraRsAttEntP - Read finished successfully

	// infraRsL2InstPol - Beginning Read
	log.Printf("[DEBUG] %s: infraRsL2InstPol - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsL2InstPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsL2InstPol - Read finished successfully", d.Get("relation_infra_rs_l2_inst_pol"))
	}
	// infraRsL2InstPol - Read finished successfully

	log.Printf("[DEBUG] %s: Data Source - Read finished successfully", dn)
	return nil
}
