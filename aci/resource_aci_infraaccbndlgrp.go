package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciPCVPCInterfacePolicyGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciPCVPCInterfacePolicyGroupCreate,
		UpdateContext: resourceAciPCVPCInterfacePolicyGroupUpdate,
		ReadContext:   resourceAciPCVPCInterfacePolicyGroupRead,
		DeleteContext: resourceAciPCVPCInterfacePolicyGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciPCVPCInterfacePolicyGroupImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"lag_t": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"not-aggregated",
					"link",
					"node",
				}, false),
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
				Computed: true,
				Optional: true,
			},
			"relation_infra_rs_lldp_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_infra_rs_macsec_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_infra_rs_qos_dpp_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_infra_rs_h_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_infra_rs_netflow_monitor_pol": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_dn": {
							Type:     schema.TypeString,
							Required: true,
						},
						"flt_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ipv4",
								"ipv6",
								"ce",
							}, false),
						},
					},
				},
			},
			"relation_infra_rs_l2_port_auth_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_infra_rs_mcp_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_infra_rs_l2_port_security_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_infra_rs_copp_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
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
				Computed: true,
				Optional: true,
			},
			"relation_infra_rs_cdp_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_infra_rs_qos_pfc_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_infra_rs_qos_sd_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_infra_rs_mon_if_infra_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_infra_rs_fc_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_infra_rs_qos_ingress_dpp_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_infra_rs_qos_egress_dpp_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_infra_rs_l2_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_infra_rs_stp_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_infra_rs_att_ent_p": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_infra_rs_l2_inst_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		}),
	}
}

func getRemotePCVPCInterfacePolicyGroup(client *client.Client, dn string) (*models.PCVPCInterfacePolicyGroup, error) {
	infraAccBndlGrpCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraAccBndlGrp := models.PCVPCInterfacePolicyGroupFromContainer(infraAccBndlGrpCont)

	if infraAccBndlGrp.DistinguishedName == "" {
		return nil, fmt.Errorf("PC/VPC Interface Policy Group %s not found", dn)
	}

	return infraAccBndlGrp, nil
}

func setPCVPCInterfacePolicyGroupAttributes(infraAccBndlGrp *models.PCVPCInterfacePolicyGroup, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(infraAccBndlGrp.DistinguishedName)
	d.Set("description", infraAccBndlGrp.Description)
	infraAccBndlGrpMap, err := infraAccBndlGrp.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("name", infraAccBndlGrpMap["name"])

	d.Set("annotation", infraAccBndlGrpMap["annotation"])
	d.Set("lag_t", infraAccBndlGrpMap["lagT"])
	d.Set("name_alias", infraAccBndlGrpMap["nameAlias"])
	return d, nil
}

func getAndSetReadRelationinfraRsSpanVSrcGrpFromPCVPCInterfacePolicyGroup(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsSpanVSrcGrpData, err := client.ReadRelationinfraRsSpanVSrcGrpFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsSpanVSrcGrp %v", err)
		d.Set("relation_infra_rs_span_v_src_grp", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_span_v_src_grp", toStringList(infraRsSpanVSrcGrpData.(*schema.Set).List()))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsStormctrlIfPolFromPCVPCInterfacePolicyGroup(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsStormctrlIfPolData, err := client.ReadRelationinfraRsStormctrlIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsStormctrlIfPol %v", err)
		d.Set("relation_infra_rs_stormctrl_if_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_stormctrl_if_pol", infraRsStormctrlIfPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsLldpIfPolFromPCVPCInterfacePolicyGroup(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsLldpIfPolData, err := client.ReadRelationinfraRsLldpIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsLldpIfPol %v", err)
		d.Set("relation_infra_rs_lldp_if_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_lldp_if_pol", infraRsLldpIfPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsMacsecIfPolFromPCVPCInterfacePolicyGroup(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsMacsecIfPolData, err := client.ReadRelationinfraRsMacsecIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsMacsecIfPol %v", err)
		d.Set("relation_infra_rs_macsec_if_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_macsec_if_pol", infraRsMacsecIfPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsQosDppIfPolFromPCVPCInterfacePolicyGroup(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsQosDppIfPolData, err := client.ReadRelationinfraRsQosDppIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsQosDppIfPol %v", err)
		d.Set("relation_infra_rs_qos_dpp_if_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_qos_dpp_if_pol", infraRsQosDppIfPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsHIfPolFromPCVPCInterfacePolicyGroup(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsHIfPolData, err := client.ReadRelationinfraRsHIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsHIfPol %v", err)
		d.Set("relation_infra_rs_h_if_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_h_if_pol", infraRsHIfPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsNetflowMonitorPolFromPCVPCInterfacePolicyGroup(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsNetflowMonitorPolData, err := client.ReadRelationinfraRsNetflowMonitorPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsNetflowMonitorPol %v", err)
		d.Set("relation_infra_rs_netflow_monitor_pol", nil)
		return d, err
	} else {
		relParamMapList := make([]map[string]string, 0, 1)
		relParams := infraRsNetflowMonitorPolData.([]map[string]string)
		for _, obj := range relParams {
			relParamMapList = append(relParamMapList, map[string]string{
				"target_dn": obj["tDn"],
				"flt_type":  obj["fltType"],
			})
		}
		d.Set("relation_infra_rs_netflow_monitor_pol", relParamMapList)
	}
	return d, nil
}

func getAndSetReadRelationinfraRsL2PortAuthPolFromPCVPCInterfacePolicyGroup(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsL2PortAuthPolData, err := client.ReadRelationinfraRsL2PortAuthPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsL2PortAuthPol %v", err)
		d.Set("relation_infra_rs_l2_port_auth_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_l2_port_auth_pol", infraRsL2PortAuthPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsMcpIfPolFromPCVPCInterfacePolicyGroup(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsMcpIfPolData, err := client.ReadRelationinfraRsMcpIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsMcpIfPol %v", err)
		d.Set("relation_infra_rs_mcp_if_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_mcp_if_pol", infraRsMcpIfPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsL2PortSecurityPolFromPCVPCInterfacePolicyGroup(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsL2PortSecurityPolData, err := client.ReadRelationinfraRsL2PortSecurityPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsL2PortSecurityPol %v", err)
		d.Set("relation_infra_rs_l2_port_security_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_l2_port_security_pol", infraRsL2PortSecurityPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsCoppIfPolFromPCVPCInterfacePolicyGroup(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsCoppIfPolData, err := client.ReadRelationinfraRsCoppIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsCoppIfPol %v", err)
		d.Set("relation_infra_rs_copp_if_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_copp_if_pol", infraRsCoppIfPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsSpanVDestGrpFromPCVPCInterfacePolicyGroup(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsSpanVDestGrpData, err := client.ReadRelationinfraRsSpanVDestGrpFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsSpanVDestGrp %v", err)
		d.Set("relation_infra_rs_span_v_dest_grp", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_span_v_dest_grp", toStringList(infraRsSpanVDestGrpData.(*schema.Set).List()))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsLacpPolFromPCVPCInterfacePolicyGroup(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsLacpPolData, err := client.ReadRelationinfraRsLacpPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsLacpPol %v", err)
		d.Set("relation_infra_rs_lacp_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_lacp_pol", infraRsLacpPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsCdpIfPolFromPCVPCInterfacePolicyGroup(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsCdpIfPolData, err := client.ReadRelationinfraRsCdpIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsCdpIfPol %v", err)
		d.Set("relation_infra_rs_cdp_if_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_cdp_if_pol", infraRsCdpIfPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsQosPfcIfPolFromPCVPCInterfacePolicyGroup(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsQosPfcIfPolData, err := client.ReadRelationinfraRsQosPfcIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsQosPfcIfPol %v", err)
		d.Set("relation_infra_rs_qos_pfc_if_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_qos_pfc_if_pol", infraRsQosPfcIfPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsQosSdIfPolFromPCVPCInterfacePolicyGroup(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsQosSdIfPolData, err := client.ReadRelationinfraRsQosSdIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsQosSdIfPol %v", err)
		d.Set("relation_infra_rs_qos_sd_if_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_qos_sd_if_pol", infraRsQosSdIfPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsMonIfInfraPolFromPCVPCInterfacePolicyGroup(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsMonIfInfraPolData, err := client.ReadRelationinfraRsMonIfInfraPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsMonIfInfraPol %v", err)
		d.Set("relation_infra_rs_mon_if_infra_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_mon_if_infra_pol", infraRsMonIfInfraPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsFcIfPolFromPCVPCInterfacePolicyGroup(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsFcIfPolData, err := client.ReadRelationinfraRsFcIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsFcIfPol %v", err)
		d.Set("relation_infra_rs_fc_if_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_fc_if_pol", infraRsFcIfPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsQosIngressDppIfPolFromPCVPCInterfacePolicyGroup(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsQosIngressDppIfPolData, err := client.ReadRelationinfraRsQosIngressDppIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsQosIngressDppIfPol %v", err)
		d.Set("relation_infra_rs_qos_ingress_dpp_if_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_qos_ingress_dpp_if_pol", infraRsQosIngressDppIfPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsQosEgressDppIfPolFromPCVPCInterfacePolicyGroup(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsQosEgressDppIfPolData, err := client.ReadRelationinfraRsQosEgressDppIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsQosEgressDppIfPol %v", err)
		d.Set("relation_infra_rs_qos_egress_dpp_if_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_qos_egress_dpp_if_pol", infraRsQosEgressDppIfPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsL2IfPolFromPCVPCInterfacePolicyGroup(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsL2IfPolData, err := client.ReadRelationinfraRsL2IfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsL2IfPol %v", err)
		d.Set("relation_infra_rs_l2_if_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_l2_if_pol", infraRsL2IfPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsStpIfPolFromPCVPCInterfacePolicyGroup(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsStpIfPolData, err := client.ReadRelationinfraRsStpIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsStpIfPol %v", err)
		d.Set("relation_infra_rs_stp_if_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_stp_if_pol", infraRsStpIfPolData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsAttEntPFromPCVPCInterfacePolicyGroup(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsAttEntPData, err := client.ReadRelationinfraRsAttEntPFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsAttEntP %v", err)
		d.Set("relation_infra_rs_att_ent_p", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_att_ent_p", infraRsAttEntPData.(string))
	}
	return d, nil
}

func getAndSetReadRelationinfraRsL2InstPolFromPCVPCInterfacePolicyGroup(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsL2InstPolData, err := client.ReadRelationinfraRsL2InstPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsL2InstPol %v", err)
		d.Set("relation_infra_rs_l2_inst_pol", nil)
		return d, err
	} else {
		d.Set("relation_infra_rs_l2_inst_pol", infraRsL2InstPolData.(string))
	}
	return d, nil
}

func resourceAciPCVPCInterfacePolicyGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraAccBndlGrp, err := getRemotePCVPCInterfacePolicyGroup(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setPCVPCInterfacePolicyGroupAttributes(infraAccBndlGrp, d)
	if err != nil {
		return nil, err
	}

	// infraRsSpanVSrcGrp - Beginning Import
	log.Printf("[DEBUG] %s: infraRsSpanVSrcGrp - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsSpanVSrcGrpFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsSpanVSrcGrp - Import finished successfully", d.Get("relation_infra_rs_span_v_src_grp"))
	}
	// infraRsSpanVSrcGrp - Import finished successfully

	// infraRsStormctrlIfPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsStormctrlIfPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsStormctrlIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsStormctrlIfPol - Import finished successfully", d.Get("relation_infra_rs_stormctrl_if_pol"))
	}
	// infraRsStormctrlIfPol - Import finished successfully

	// infraRsLldpIfPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsLldpIfPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsLldpIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsLldpIfPol - Import finished successfully", d.Get("relation_infra_rs_lldp_if_pol"))
	}
	// infraRsLldpIfPol - Import finished successfully

	// infraRsMacsecIfPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsMacsecIfPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsMacsecIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsMacsecIfPol - Import finished successfully", d.Get("relation_infra_rs_macsec_if_pol"))
	}
	// infraRsMacsecIfPol - Import finished successfully

	// infraRsQosDppIfPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsQosDppIfPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsQosDppIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsQosDppIfPol - Import finished successfully", d.Get("relation_infra_rs_qos_dpp_if_pol"))
	}
	// infraRsQosDppIfPol - Import finished successfully

	// infraRsHIfPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsHIfPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsHIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsHIfPol - Import finished successfully", d.Get("relation_infra_rs_h_if_pol"))
	}
	// infraRsHIfPol - Import finished successfully

	// infraRsNetflowMonitorPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsNetflowMonitorPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsNetflowMonitorPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsNetflowMonitorPol - Import finished successfully", d.Get("relation_infra_rs_netflow_monitor_pol"))
	}
	// infraRsNetflowMonitorPol - Import finished successfully

	// infraRsL2PortAuthPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsL2PortAuthPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsL2PortAuthPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsL2PortAuthPol - Import finished successfully", d.Get("relation_infra_rs_l2_port_auth_pol"))
	}
	// infraRsL2PortAuthPol - Import finished successfully

	// infraRsMcpIfPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsMcpIfPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsMcpIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsMcpIfPol - Import finished successfully", d.Get("relation_infra_rs_mcp_if_pol"))
	}
	// infraRsMcpIfPol - Import finished successfully

	// infraRsL2PortSecurityPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsL2PortSecurityPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsL2PortSecurityPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsL2PortSecurityPol - Import finished successfully", d.Get("relation_infra_rs_l2_port_security_pol"))
	}
	// infraRsL2PortSecurityPol - Import finished successfully

	// infraRsCoppIfPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsCoppIfPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsCoppIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsCoppIfPol - Import finished successfully", d.Get("relation_infra_rs_copp_if_pol"))
	}
	// infraRsCoppIfPol - Import finished successfully

	// infraRsSpanVDestGrp - Beginning Import
	log.Printf("[DEBUG] %s: infraRsSpanVDestGrp - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsSpanVDestGrpFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsSpanVDestGrp - Import finished successfully", d.Get("relation_infra_rs_span_v_dest_grp"))
	}
	// infraRsSpanVDestGrp - Import finished successfully

	// infraRsLacpPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsLacpPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsLacpPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsLacpPol - Import finished successfully", d.Get("relation_infra_rs_lacp_pol"))
	}
	// infraRsLacpPol - Import finished successfully

	// infraRsCdpIfPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsCdpIfPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsCdpIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsCdpIfPol - Import finished successfully", d.Get("relation_infra_rs_cdp_if_pol"))
	}
	// infraRsCdpIfPol - Import finished successfully

	// infraRsQosPfcIfPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsQosPfcIfPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsQosPfcIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsQosPfcIfPol - Import finished successfully", d.Get("relation_infra_rs_qos_pfc_if_pol"))
	}
	// infraRsQosPfcIfPol - Import finished successfully

	// infraRsQosSdIfPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsQosSdIfPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsQosSdIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsQosSdIfPol - Import finished successfully", d.Get("relation_infra_rs_qos_sd_if_pol"))
	}
	// infraRsQosSdIfPol - Import finished successfully

	// infraRsMonIfInfraPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsMonIfInfraPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsMonIfInfraPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsMonIfInfraPol - Import finished successfully", d.Get("relation_infra_rs_mon_if_infra_pol"))
	}
	// infraRsMonIfInfraPol - Import finished successfully

	// infraRsFcIfPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsFcIfPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsFcIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsFcIfPol - Import finished successfully", d.Get("relation_infra_rs_fc_if_pol"))
	}
	// infraRsFcIfPol - Import finished successfully

	// infraRsQosIngressDppIfPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsQosIngressDppIfPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsQosIngressDppIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsQosIngressDppIfPol - Import finished successfully", d.Get("relation_infra_rs_qos_ingress_dpp_if_pol"))
	}
	// infraRsQosIngressDppIfPol - Import finished successfully

	// infraRsQosEgressDppIfPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsQosEgressDppIfPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsQosEgressDppIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsQosEgressDppIfPol - Import finished successfully", d.Get("relation_infra_rs_qos_egress_dpp_if_pol"))
	}
	// infraRsQosEgressDppIfPol - Import finished successfully

	// infraRsL2IfPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsL2IfPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsL2IfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsL2IfPol - Import finished successfully", d.Get("relation_infra_rs_l2_if_pol"))
	}
	// infraRsL2IfPol - Import finished successfully

	// infraRsStpIfPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsStpIfPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsStpIfPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsStpIfPol - Import finished successfully", d.Get("relation_infra_rs_stp_if_pol"))
	}
	// infraRsStpIfPol - Import finished successfully

	// infraRsAttEntP - Beginning Import
	log.Printf("[DEBUG] %s: infraRsAttEntP - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsAttEntPFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsAttEntP - Import finished successfully", d.Get("relation_infra_rs_att_ent_p"))
	}
	// infraRsAttEntP - Import finished successfully

	// infraRsL2InstPol - Beginning Import
	log.Printf("[DEBUG] %s: infraRsL2InstPol - Beginning Import with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsL2InstPolFromPCVPCInterfacePolicyGroup(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsL2InstPol - Import finished successfully", d.Get("relation_infra_rs_l2_inst_pol"))
	}
	// infraRsL2InstPol - Import finished successfully

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciPCVPCInterfacePolicyGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] PCVPCInterfacePolicyGroup: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraAccBndlGrpAttr := models.PCVPCInterfacePolicyGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraAccBndlGrpAttr.Annotation = Annotation.(string)
	} else {
		infraAccBndlGrpAttr.Annotation = "{}"
	}
	if LagT, ok := d.GetOk("lag_t"); ok {
		infraAccBndlGrpAttr.LagT = LagT.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraAccBndlGrpAttr.NameAlias = NameAlias.(string)
	}
	infraAccBndlGrp := models.NewPCVPCInterfacePolicyGroup(fmt.Sprintf("infra/funcprof/accbundle-%s", name), "uni", desc, infraAccBndlGrpAttr)

	err := aciClient.Save(infraAccBndlGrp)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationToinfraRsSpanVSrcGrp, ok := d.GetOk("relation_infra_rs_span_v_src_grp"); ok {
		relationParamList := toStringList(relationToinfraRsSpanVSrcGrp.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationToinfraRsStormctrlIfPol, ok := d.GetOk("relation_infra_rs_stormctrl_if_pol"); ok {
		relationParam := relationToinfraRsStormctrlIfPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsLldpIfPol, ok := d.GetOk("relation_infra_rs_lldp_if_pol"); ok {
		relationParam := relationToinfraRsLldpIfPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsMacsecIfPol, ok := d.GetOk("relation_infra_rs_macsec_if_pol"); ok {
		relationParam := relationToinfraRsMacsecIfPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsQosDppIfPol, ok := d.GetOk("relation_infra_rs_qos_dpp_if_pol"); ok {
		relationParam := relationToinfraRsQosDppIfPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsHIfPol, ok := d.GetOk("relation_infra_rs_h_if_pol"); ok {
		relationParam := relationToinfraRsHIfPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsL2PortAuthPol, ok := d.GetOk("relation_infra_rs_l2_port_auth_pol"); ok {
		relationParam := relationToinfraRsL2PortAuthPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsMcpIfPol, ok := d.GetOk("relation_infra_rs_mcp_if_pol"); ok {
		relationParam := relationToinfraRsMcpIfPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsL2PortSecurityPol, ok := d.GetOk("relation_infra_rs_l2_port_security_pol"); ok {
		relationParam := relationToinfraRsL2PortSecurityPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsCoppIfPol, ok := d.GetOk("relation_infra_rs_copp_if_pol"); ok {
		relationParam := relationToinfraRsCoppIfPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsSpanVDestGrp, ok := d.GetOk("relation_infra_rs_span_v_dest_grp"); ok {
		relationParamList := toStringList(relationToinfraRsSpanVDestGrp.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationToinfraRsLacpPol, ok := d.GetOk("relation_infra_rs_lacp_pol"); ok {
		relationParam := relationToinfraRsLacpPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsCdpIfPol, ok := d.GetOk("relation_infra_rs_cdp_if_pol"); ok {
		relationParam := relationToinfraRsCdpIfPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsQosPfcIfPol, ok := d.GetOk("relation_infra_rs_qos_pfc_if_pol"); ok {
		relationParam := relationToinfraRsQosPfcIfPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsQosSdIfPol, ok := d.GetOk("relation_infra_rs_qos_sd_if_pol"); ok {
		relationParam := relationToinfraRsQosSdIfPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsMonIfInfraPol, ok := d.GetOk("relation_infra_rs_mon_if_infra_pol"); ok {
		relationParam := relationToinfraRsMonIfInfraPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsFcIfPol, ok := d.GetOk("relation_infra_rs_fc_if_pol"); ok {
		relationParam := relationToinfraRsFcIfPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsQosIngressDppIfPol, ok := d.GetOk("relation_infra_rs_qos_ingress_dpp_if_pol"); ok {
		relationParam := relationToinfraRsQosIngressDppIfPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsQosEgressDppIfPol, ok := d.GetOk("relation_infra_rs_qos_egress_dpp_if_pol"); ok {
		relationParam := relationToinfraRsQosEgressDppIfPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsL2IfPol, ok := d.GetOk("relation_infra_rs_l2_if_pol"); ok {
		relationParam := relationToinfraRsL2IfPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsStpIfPol, ok := d.GetOk("relation_infra_rs_stp_if_pol"); ok {
		relationParam := relationToinfraRsStpIfPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsAttEntP, ok := d.GetOk("relation_infra_rs_att_ent_p"); ok {
		relationParam := relationToinfraRsAttEntP.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsL2InstPol, ok := d.GetOk("relation_infra_rs_l2_inst_pol"); ok {
		relationParam := relationToinfraRsL2InstPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationToinfraRsSpanVSrcGrp, ok := d.GetOk("relation_infra_rs_span_v_src_grp"); ok {
		relationParamList := toStringList(relationToinfraRsSpanVSrcGrp.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationinfraRsSpanVSrcGrpFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParamName)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if relationToinfraRsStormctrlIfPol, ok := d.GetOk("relation_infra_rs_stormctrl_if_pol"); ok {
		relationParam := relationToinfraRsStormctrlIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsStormctrlIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsLldpIfPol, ok := d.GetOk("relation_infra_rs_lldp_if_pol"); ok {
		relationParam := relationToinfraRsLldpIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsLldpIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsMacsecIfPol, ok := d.GetOk("relation_infra_rs_macsec_if_pol"); ok {
		relationParam := relationToinfraRsMacsecIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsMacsecIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsQosDppIfPol, ok := d.GetOk("relation_infra_rs_qos_dpp_if_pol"); ok {
		relationParam := relationToinfraRsQosDppIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsQosDppIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsHIfPol, ok := d.GetOk("relation_infra_rs_h_if_pol"); ok {
		relationParam := relationToinfraRsHIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsHIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsNetflowMonitorPol, ok := d.GetOk("relation_infra_rs_netflow_monitor_pol"); ok {

		relationParamList := relationToinfraRsNetflowMonitorPol.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationinfraRsNetflowMonitorPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, GetMOName(paramMap["target_dn"].(string)), paramMap["flt_type"].(string))
			if err != nil {
				return diag.FromErr(err)
			}
		}

	}
	if relationToinfraRsL2PortAuthPol, ok := d.GetOk("relation_infra_rs_l2_port_auth_pol"); ok {
		relationParam := relationToinfraRsL2PortAuthPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsL2PortAuthPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsMcpIfPol, ok := d.GetOk("relation_infra_rs_mcp_if_pol"); ok {
		relationParam := relationToinfraRsMcpIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsMcpIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsL2PortSecurityPol, ok := d.GetOk("relation_infra_rs_l2_port_security_pol"); ok {
		relationParam := relationToinfraRsL2PortSecurityPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsL2PortSecurityPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationToinfraRsCoppIfPol, ok := d.GetOk("relation_infra_rs_copp_if_pol"); ok {
		relationParam := relationToinfraRsCoppIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsCoppIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsSpanVDestGrp, ok := d.GetOk("relation_infra_rs_span_v_dest_grp"); ok {
		relationParamList := toStringList(relationToinfraRsSpanVDestGrp.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationinfraRsSpanVDestGrpFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParamName)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if relationToinfraRsLacpPol, ok := d.GetOk("relation_infra_rs_lacp_pol"); ok {
		relationParam := relationToinfraRsLacpPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsLacpPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationToinfraRsCdpIfPol, ok := d.GetOk("relation_infra_rs_cdp_if_pol"); ok {
		relationParam := relationToinfraRsCdpIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsCdpIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsQosPfcIfPol, ok := d.GetOk("relation_infra_rs_qos_pfc_if_pol"); ok {
		relationParam := relationToinfraRsQosPfcIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsQosPfcIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsQosSdIfPol, ok := d.GetOk("relation_infra_rs_qos_sd_if_pol"); ok {
		relationParam := relationToinfraRsQosSdIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsQosSdIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsMonIfInfraPol, ok := d.GetOk("relation_infra_rs_mon_if_infra_pol"); ok {
		relationParam := relationToinfraRsMonIfInfraPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsMonIfInfraPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsFcIfPol, ok := d.GetOk("relation_infra_rs_fc_if_pol"); ok {
		relationParam := relationToinfraRsFcIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsFcIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsQosIngressDppIfPol, ok := d.GetOk("relation_infra_rs_qos_ingress_dpp_if_pol"); ok {
		relationParam := relationToinfraRsQosIngressDppIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsQosIngressDppIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsQosEgressDppIfPol, ok := d.GetOk("relation_infra_rs_qos_egress_dpp_if_pol"); ok {
		relationParam := relationToinfraRsQosEgressDppIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsQosEgressDppIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsL2IfPol, ok := d.GetOk("relation_infra_rs_l2_if_pol"); ok {
		relationParam := relationToinfraRsL2IfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsL2IfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsStpIfPol, ok := d.GetOk("relation_infra_rs_stp_if_pol"); ok {
		relationParam := relationToinfraRsStpIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsStpIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsAttEntP, ok := d.GetOk("relation_infra_rs_att_ent_p"); ok {
		relationParam := relationToinfraRsAttEntP.(string)
		err = aciClient.CreateRelationinfraRsAttEntPFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsL2InstPol, ok := d.GetOk("relation_infra_rs_l2_inst_pol"); ok {
		relationParam := relationToinfraRsL2InstPol.(string)
		err = aciClient.CreateRelationinfraRsL2InstPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(infraAccBndlGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciPCVPCInterfacePolicyGroupRead(ctx, d, m)
}

func resourceAciPCVPCInterfacePolicyGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] PCVPCInterfacePolicyGroup: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraAccBndlGrpAttr := models.PCVPCInterfacePolicyGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraAccBndlGrpAttr.Annotation = Annotation.(string)
	} else {
		infraAccBndlGrpAttr.Annotation = "{}"
	}
	if LagT, ok := d.GetOk("lag_t"); ok {
		infraAccBndlGrpAttr.LagT = LagT.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraAccBndlGrpAttr.NameAlias = NameAlias.(string)
	}
	infraAccBndlGrp := models.NewPCVPCInterfacePolicyGroup(fmt.Sprintf("infra/funcprof/accbundle-%s", name), "uni", desc, infraAccBndlGrpAttr)

	infraAccBndlGrp.Status = "modified"

	err := aciClient.Save(infraAccBndlGrp)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_infra_rs_span_v_src_grp") {
		oldRel, newRel := d.GetChange("relation_infra_rs_span_v_src_grp")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_infra_rs_stormctrl_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_stormctrl_if_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_lldp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_lldp_if_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_macsec_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_macsec_if_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_qos_dpp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_qos_dpp_if_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_h_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_h_if_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_l2_port_auth_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_l2_port_auth_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_mcp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_mcp_if_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_l2_port_security_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_l2_port_security_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_copp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_copp_if_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_span_v_dest_grp") {
		oldRel, newRel := d.GetChange("relation_infra_rs_span_v_dest_grp")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_infra_rs_lacp_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_lacp_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_cdp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_cdp_if_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_qos_pfc_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_qos_pfc_if_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_qos_sd_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_qos_sd_if_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_mon_if_infra_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_mon_if_infra_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_fc_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_fc_if_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_qos_ingress_dpp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_qos_ingress_dpp_if_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_qos_egress_dpp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_qos_egress_dpp_if_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_l2_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_l2_if_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_stp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_stp_if_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_att_ent_p") {
		_, newRelParam := d.GetChange("relation_infra_rs_att_ent_p")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_l2_inst_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_l2_inst_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_infra_rs_span_v_src_grp") {
		oldRel, newRel := d.GetChange("relation_infra_rs_span_v_src_grp")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			relDnName := GetMOName(relDn)
			err = aciClient.DeleteRelationinfraRsSpanVSrcGrpFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationinfraRsSpanVSrcGrpFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}
	if d.HasChange("relation_infra_rs_stormctrl_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_stormctrl_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsStormctrlIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_lldp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_lldp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsLldpIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_macsec_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_macsec_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsMacsecIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_qos_dpp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_qos_dpp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsQosDppIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_h_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_h_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsHIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_netflow_monitor_pol") {
		oldRel, newRel := d.GetChange("relation_infra_rs_netflow_monitor_pol")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.DeleteRelationinfraRsNetflowMonitorPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, GetMOName(paramMap["target_dn"].(string)), paramMap["flt_type"].(string))
			if err != nil {
				return diag.FromErr(err)
			}

		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationinfraRsNetflowMonitorPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, GetMOName(paramMap["target_dn"].(string)), paramMap["flt_type"].(string))
			if err != nil {
				return diag.FromErr(err)
			}
		}

	}
	if d.HasChange("relation_infra_rs_l2_port_auth_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_l2_port_auth_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsL2PortAuthPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_mcp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_mcp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsMcpIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_l2_port_security_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_l2_port_security_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsL2PortSecurityPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_copp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_copp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsCoppIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_span_v_dest_grp") {
		oldRel, newRel := d.GetChange("relation_infra_rs_span_v_dest_grp")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			relDnName := GetMOName(relDn)
			err = aciClient.DeleteRelationinfraRsSpanVDestGrpFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationinfraRsSpanVDestGrpFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}
	if d.HasChange("relation_infra_rs_lacp_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_lacp_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsLacpPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_cdp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_cdp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsCdpIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_qos_pfc_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_qos_pfc_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsQosPfcIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_qos_sd_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_qos_sd_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsQosSdIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_mon_if_infra_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_mon_if_infra_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsMonIfInfraPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_fc_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_fc_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsFcIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_qos_ingress_dpp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_qos_ingress_dpp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsQosIngressDppIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_qos_egress_dpp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_qos_egress_dpp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsQosEgressDppIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_l2_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_l2_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsL2IfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_stp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_stp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsStpIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_att_ent_p") {
		_, newRelParam := d.GetChange("relation_infra_rs_att_ent_p")
		err = aciClient.DeleteRelationinfraRsAttEntPFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsAttEntPFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_l2_inst_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_l2_inst_pol")
		err = aciClient.DeleteRelationinfraRsL2InstPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsL2InstPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(infraAccBndlGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciPCVPCInterfacePolicyGroupRead(ctx, d, m)

}

func resourceAciPCVPCInterfacePolicyGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraAccBndlGrp, err := getRemotePCVPCInterfacePolicyGroup(aciClient, dn)

	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	_, err = setPCVPCInterfacePolicyGroupAttributes(infraAccBndlGrp, d)
	if err != nil {
		d.SetId("")
		return nil
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

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciPCVPCInterfacePolicyGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraAccBndlGrp")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
