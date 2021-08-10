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

func resourceAciLeafAccessPortPolicyGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciLeafAccessPortPolicyGroupCreate,
		UpdateContext: resourceAciLeafAccessPortPolicyGroupUpdate,
		ReadContext:   resourceAciLeafAccessPortPolicyGroupRead,
		DeleteContext: resourceAciLeafAccessPortPolicyGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciLeafAccessPortPolicyGroupImport,
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
			"relation_infra_rs_poe_if_pol": &schema.Schema{
				Type: schema.TypeString,

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
			"relation_infra_rs_dwdm_if_pol": &schema.Schema{
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
			"relation_infra_rs_cdp_if_pol": &schema.Schema{
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
			"relation_infra_rs_qos_egress_dpp_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_infra_rs_att_ent_p": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_l2_inst_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteLeafAccessPortPolicyGroup(client *client.Client, dn string) (*models.LeafAccessPortPolicyGroup, error) {
	infraAccPortGrpCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraAccPortGrp := models.LeafAccessPortPolicyGroupFromContainer(infraAccPortGrpCont)

	if infraAccPortGrp.DistinguishedName == "" {
		return nil, fmt.Errorf("LeafAccessPortPolicyGroup %s not found", infraAccPortGrp.DistinguishedName)
	}

	return infraAccPortGrp, nil
}

func setLeafAccessPortPolicyGroupAttributes(infraAccPortGrp *models.LeafAccessPortPolicyGroup, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(infraAccPortGrp.DistinguishedName)
	d.Set("description", infraAccPortGrp.Description)
	infraAccPortGrpMap, err := infraAccPortGrp.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("name", infraAccPortGrpMap["name"])

	d.Set("annotation", infraAccPortGrpMap["annotation"])
	d.Set("name_alias", infraAccPortGrpMap["nameAlias"])
	return d, nil
}

func resourceAciLeafAccessPortPolicyGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraAccPortGrp, err := getRemoteLeafAccessPortPolicyGroup(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setLeafAccessPortPolicyGroupAttributes(infraAccPortGrp, d)

	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLeafAccessPortPolicyGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LeafAccessPortPolicyGroup: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraAccPortGrpAttr := models.LeafAccessPortPolicyGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraAccPortGrpAttr.Annotation = Annotation.(string)
	} else {
		infraAccPortGrpAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraAccPortGrpAttr.NameAlias = NameAlias.(string)
	}
	infraAccPortGrp := models.NewLeafAccessPortPolicyGroup(fmt.Sprintf("infra/funcprof/accportgrp-%s", name), "uni", desc, infraAccPortGrpAttr)

	err := aciClient.Save(infraAccPortGrp)
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

	if relationToinfraRsPoeIfPol, ok := d.GetOk("relation_infra_rs_poe_if_pol"); ok {
		relationParam := relationToinfraRsPoeIfPol.(string)
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

	if relationToinfraRsDwdmIfPol, ok := d.GetOk("relation_infra_rs_dwdm_if_pol"); ok {
		relationParam := relationToinfraRsDwdmIfPol.(string)
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

	if relationToinfraRsCdpIfPol, ok := d.GetOk("relation_infra_rs_cdp_if_pol"); ok {
		relationParam := relationToinfraRsCdpIfPol.(string)
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

	if relationToinfraRsQosEgressDppIfPol, ok := d.GetOk("relation_infra_rs_qos_egress_dpp_if_pol"); ok {
		relationParam := relationToinfraRsQosEgressDppIfPol.(string)
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
			err = aciClient.CreateRelationinfraRsSpanVSrcGrpFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if relationToinfraRsStormctrlIfPol, ok := d.GetOk("relation_infra_rs_stormctrl_if_pol"); ok {
		relationParam := relationToinfraRsStormctrlIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsStormctrlIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsPoeIfPol, ok := d.GetOk("relation_infra_rs_poe_if_pol"); ok {
		relationParam := relationToinfraRsPoeIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsPoeIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationToinfraRsLldpIfPol, ok := d.GetOk("relation_infra_rs_lldp_if_pol"); ok {
		relationParam := relationToinfraRsLldpIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsLldpIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsMacsecIfPol, ok := d.GetOk("relation_infra_rs_macsec_if_pol"); ok {
		relationParam := relationToinfraRsMacsecIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsMacsecIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsQosDppIfPol, ok := d.GetOk("relation_infra_rs_qos_dpp_if_pol"); ok {
		relationParam := relationToinfraRsQosDppIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsQosDppIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationToinfraRsHIfPol, ok := d.GetOk("relation_infra_rs_h_if_pol"); ok {
		relationParam := relationToinfraRsHIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsHIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsNetflowMonitorPol, ok := d.GetOk("relation_infra_rs_netflow_monitor_pol"); ok {

		relationParamList := relationToinfraRsNetflowMonitorPol.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationinfraRsNetflowMonitorPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, GetMOName(paramMap["target_dn"].(string)), paramMap["flt_type"].(string))
			if err != nil {
				return diag.FromErr(err)
			}
		}

	}
	if relationToinfraRsL2PortAuthPol, ok := d.GetOk("relation_infra_rs_l2_port_auth_pol"); ok {
		relationParam := relationToinfraRsL2PortAuthPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsL2PortAuthPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsMcpIfPol, ok := d.GetOk("relation_infra_rs_mcp_if_pol"); ok {
		relationParam := relationToinfraRsMcpIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsMcpIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationToinfraRsL2PortSecurityPol, ok := d.GetOk("relation_infra_rs_l2_port_security_pol"); ok {
		relationParam := relationToinfraRsL2PortSecurityPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsL2PortSecurityPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsCoppIfPol, ok := d.GetOk("relation_infra_rs_copp_if_pol"); ok {
		relationParam := relationToinfraRsCoppIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsCoppIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsSpanVDestGrp, ok := d.GetOk("relation_infra_rs_span_v_dest_grp"); ok {
		relationParamList := toStringList(relationToinfraRsSpanVDestGrp.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationinfraRsSpanVDestGrpFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if relationToinfraRsDwdmIfPol, ok := d.GetOk("relation_infra_rs_dwdm_if_pol"); ok {
		relationParam := relationToinfraRsDwdmIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsDwdmIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsQosPfcIfPol, ok := d.GetOk("relation_infra_rs_qos_pfc_if_pol"); ok {
		relationParam := relationToinfraRsQosPfcIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsQosPfcIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsQosSdIfPol, ok := d.GetOk("relation_infra_rs_qos_sd_if_pol"); ok {
		relationParam := relationToinfraRsQosSdIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsQosSdIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsMonIfInfraPol, ok := d.GetOk("relation_infra_rs_mon_if_infra_pol"); ok {
		relationParam := relationToinfraRsMonIfInfraPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsMonIfInfraPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsFcIfPol, ok := d.GetOk("relation_infra_rs_fc_if_pol"); ok {
		relationParam := relationToinfraRsFcIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsFcIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsQosIngressDppIfPol, ok := d.GetOk("relation_infra_rs_qos_ingress_dpp_if_pol"); ok {
		relationParam := relationToinfraRsQosIngressDppIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsQosIngressDppIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsCdpIfPol, ok := d.GetOk("relation_infra_rs_cdp_if_pol"); ok {
		relationParam := relationToinfraRsCdpIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsCdpIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsL2IfPol, ok := d.GetOk("relation_infra_rs_l2_if_pol"); ok {
		relationParam := relationToinfraRsL2IfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsL2IfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsStpIfPol, ok := d.GetOk("relation_infra_rs_stp_if_pol"); ok {
		relationParam := relationToinfraRsStpIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsStpIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsQosEgressDppIfPol, ok := d.GetOk("relation_infra_rs_qos_egress_dpp_if_pol"); ok {
		relationParam := relationToinfraRsQosEgressDppIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsQosEgressDppIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsAttEntP, ok := d.GetOk("relation_infra_rs_att_ent_p"); ok {
		relationParam := relationToinfraRsAttEntP.(string)
		err = aciClient.CreateRelationinfraRsAttEntPFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsL2InstPol, ok := d.GetOk("relation_infra_rs_l2_inst_pol"); ok {
		relationParam := relationToinfraRsL2InstPol.(string)
		err = aciClient.CreateRelationinfraRsL2InstPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(infraAccPortGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciLeafAccessPortPolicyGroupRead(ctx, d, m)
}

func resourceAciLeafAccessPortPolicyGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LeafAccessPortPolicyGroup: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraAccPortGrpAttr := models.LeafAccessPortPolicyGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraAccPortGrpAttr.Annotation = Annotation.(string)
	} else {
		infraAccPortGrpAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraAccPortGrpAttr.NameAlias = NameAlias.(string)
	}
	infraAccPortGrp := models.NewLeafAccessPortPolicyGroup(fmt.Sprintf("infra/funcprof/accportgrp-%s", name), "uni", desc, infraAccPortGrpAttr)

	infraAccPortGrp.Status = "modified"

	err := aciClient.Save(infraAccPortGrp)

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

	if d.HasChange("relation_infra_rs_poe_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_poe_if_pol")
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

	if d.HasChange("relation_infra_rs_dwdm_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_dwdm_if_pol")
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

	if d.HasChange("relation_infra_rs_cdp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_cdp_if_pol")
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

	if d.HasChange("relation_infra_rs_qos_egress_dpp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_qos_egress_dpp_if_pol")
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
			err = aciClient.DeleteRelationinfraRsSpanVSrcGrpFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationinfraRsSpanVSrcGrpFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}
	if d.HasChange("relation_infra_rs_stormctrl_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_stormctrl_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsStormctrlIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_poe_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_poe_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationinfraRsPoeIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsPoeIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_lldp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_lldp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsLldpIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_macsec_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_macsec_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsMacsecIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_qos_dpp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_qos_dpp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsQosDppIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_h_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_h_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsHIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
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
			err = aciClient.DeleteRelationinfraRsNetflowMonitorPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, GetMOName(paramMap["target_dn"].(string)), paramMap["flt_type"].(string))
			if err != nil {
				return diag.FromErr(err)
			}

		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationinfraRsNetflowMonitorPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, GetMOName(paramMap["target_dn"].(string)), paramMap["flt_type"].(string))
			if err != nil {
				return diag.FromErr(err)
			}
		}

	}
	if d.HasChange("relation_infra_rs_l2_port_auth_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_l2_port_auth_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsL2PortAuthPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_mcp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_mcp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsMcpIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_l2_port_security_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_l2_port_security_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsL2PortSecurityPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_copp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_copp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsCoppIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
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
			err = aciClient.DeleteRelationinfraRsSpanVDestGrpFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationinfraRsSpanVDestGrpFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}
	if d.HasChange("relation_infra_rs_dwdm_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_dwdm_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsDwdmIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_qos_pfc_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_qos_pfc_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsQosPfcIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_qos_sd_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_qos_sd_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsQosSdIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_mon_if_infra_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_mon_if_infra_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsMonIfInfraPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_fc_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_fc_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsFcIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_qos_ingress_dpp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_qos_ingress_dpp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsQosIngressDppIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_cdp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_cdp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsCdpIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_l2_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_l2_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsL2IfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_infra_rs_stp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_stp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsStpIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_qos_egress_dpp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_qos_egress_dpp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsQosEgressDppIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_att_ent_p") {
		_, newRelParam := d.GetChange("relation_infra_rs_att_ent_p")
		err = aciClient.DeleteRelationinfraRsAttEntPFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsAttEntPFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_l2_inst_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_l2_inst_pol")
		err = aciClient.DeleteRelationinfraRsL2InstPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsL2InstPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(infraAccPortGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciLeafAccessPortPolicyGroupRead(ctx, d, m)

}

func resourceAciLeafAccessPortPolicyGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraAccPortGrp, err := getRemoteLeafAccessPortPolicyGroup(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setLeafAccessPortPolicyGroupAttributes(infraAccPortGrp, d)
	if err != nil {
		d.SetId("")
		return nil
	}
	infraRsSpanVSrcGrpData, err := aciClient.ReadRelationinfraRsSpanVSrcGrpFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsSpanVSrcGrp %v", err)
		d.Set("relation_infra_rs_span_v_src_grp", make([]string, 0, 1))

	} else {
		d.Set("relation_infra_rs_span_v_src_grp", toStringList(infraRsSpanVSrcGrpData.(*schema.Set).List()))
	}

	infraRsStormctrlIfPolData, err := aciClient.ReadRelationinfraRsStormctrlIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsStormctrlIfPol %v", err)
		d.Set("relation_infra_rs_stormctrl_if_pol", "")

	} else {
		d.Set("relation_infra_rs_stormctrl_if_pol", infraRsStormctrlIfPolData.(string))
	}

	infraRsPoeIfPolData, err := aciClient.ReadRelationinfraRsPoeIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsPoeIfPol %v", err)
		d.Set("relation_infra_rs_poe_if_pol", "")

	} else {
		d.Set("relation_infra_rs_poe_if_pol", infraRsPoeIfPolData.(string))
	}

	infraRsLldpIfPolData, err := aciClient.ReadRelationinfraRsLldpIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsLldpIfPol %v", err)
		d.Set("relation_infra_rs_lldp_if_pol", "")

	} else {
		d.Set("relation_infra_rs_lldp_if_pol", infraRsLldpIfPolData.(string))
	}

	infraRsMacsecIfPolData, err := aciClient.ReadRelationinfraRsMacsecIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsMacsecIfPol %v", err)
		d.Set("relation_infra_rs_macsec_if_pol", "")

	} else {
		d.Set("relation_infra_rs_macsec_if_pol", infraRsMacsecIfPolData.(string))
	}

	infraRsQosDppIfPolData, err := aciClient.ReadRelationinfraRsQosDppIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsQosDppIfPol %v", err)
		d.Set("relation_infra_rs_qos_dpp_if_pol", "")

	} else {
		d.Set("relation_infra_rs_qos_dpp_if_pol", infraRsQosDppIfPolData.(string))
	}

	infraRsHIfPolData, err := aciClient.ReadRelationinfraRsHIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsHIfPol %v", err)
		d.Set("relation_infra_rs_h_if_pol", "")

	} else {
		d.Set("relation_infra_rs_h_if_pol", infraRsHIfPolData.(string))
	}

	infraRsNetflowMonitorPolData, err := aciClient.ReadRelationinfraRsNetflowMonitorPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsNetflowMonitorPol %v", err)
	} else {
		relParamList := make([]map[string]string, 0, 1)
		relParams := infraRsNetflowMonitorPolData.([]map[string]string)
		for _, obj := range relParams {
			relParamList = append(relParamList, map[string]string{
				"target_dn": obj["tnNetflowMonitorPolName"],
				"flt_type":  obj["fltType"],
			})
		}
		d.Set("relation_infra_rs_netflow_monitor_pol", relParamList)
	}

	infraRsL2PortAuthPolData, err := aciClient.ReadRelationinfraRsL2PortAuthPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsL2PortAuthPol %v", err)
		d.Set("relation_infra_rs_l2_port_auth_pol", "")

	} else {
		d.Set("relation_infra_rs_l2_port_auth_pol", infraRsL2PortAuthPolData.(string))
	}

	infraRsMcpIfPolData, err := aciClient.ReadRelationinfraRsMcpIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsMcpIfPol %v", err)
		d.Set("relation_infra_rs_mcp_if_pol", "")

	} else {
		d.Set("relation_infra_rs_mcp_if_pol", infraRsMcpIfPolData.(string))
	}

	infraRsL2PortSecurityPolData, err := aciClient.ReadRelationinfraRsL2PortSecurityPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsL2PortSecurityPol %v", err)
		d.Set("relation_infra_rs_l2_port_security_pol", "")

	} else {
		d.Set("relation_infra_rs_l2_port_security_pol", infraRsL2PortSecurityPolData.(string))
	}

	infraRsCoppIfPolData, err := aciClient.ReadRelationinfraRsCoppIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsCoppIfPol %v", err)
		d.Set("relation_infra_rs_copp_if_pol", "")

	} else {
		d.Set("relation_infra_rs_copp_if_pol", infraRsCoppIfPolData.(string))
	}

	infraRsSpanVDestGrpData, err := aciClient.ReadRelationinfraRsSpanVDestGrpFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsSpanVDestGrp %v", err)
		d.Set("relation_infra_rs_span_v_dest_grp", make([]string, 0, 1))

	} else {
		d.Set("relation_infra_rs_span_v_dest_grp", toStringList(infraRsSpanVDestGrpData.(*schema.Set).List()))
	}

	infraRsDwdmIfPolData, err := aciClient.ReadRelationinfraRsDwdmIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsDwdmIfPol %v", err)
		d.Set("relation_infra_rs_dwdm_if_pol", "")

	} else {
		d.Set("relation_infra_rs_dwdm_if_pol", infraRsDwdmIfPolData.(string))
	}

	infraRsQosPfcIfPolData, err := aciClient.ReadRelationinfraRsQosPfcIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsQosPfcIfPol %v", err)
		d.Set("relation_infra_rs_qos_pfc_if_pol", "")

	} else {
		d.Set("relation_infra_rs_qos_pfc_if_pol", infraRsQosPfcIfPolData.(string))
	}

	infraRsQosSdIfPolData, err := aciClient.ReadRelationinfraRsQosSdIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsQosSdIfPol %v", err)
		d.Set("relation_infra_rs_qos_sd_if_pol", "")

	} else {
		d.Set("relation_infra_rs_qos_sd_if_pol", infraRsQosSdIfPolData.(string))
	}

	infraRsMonIfInfraPolData, err := aciClient.ReadRelationinfraRsMonIfInfraPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsMonIfInfraPol %v", err)
		d.Set("relation_infra_rs_mon_if_infra_pol", "")

	} else {
		d.Set("relation_infra_rs_mon_if_infra_pol", infraRsMonIfInfraPolData.(string))
	}

	infraRsFcIfPolData, err := aciClient.ReadRelationinfraRsFcIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsFcIfPol %v", err)
		d.Set("relation_infra_rs_fc_if_pol", "")

	} else {
		d.Set("relation_infra_rs_fc_if_pol", infraRsFcIfPolData.(string))
	}

	infraRsQosIngressDppIfPolData, err := aciClient.ReadRelationinfraRsQosIngressDppIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsQosIngressDppIfPol %v", err)
		d.Set("relation_infra_rs_qos_ingress_dpp_if_pol", "")

	} else {
		d.Set("relation_infra_rs_qos_ingress_dpp_if_pol", infraRsQosIngressDppIfPolData.(string))
	}

	infraRsCdpIfPolData, err := aciClient.ReadRelationinfraRsCdpIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsCdpIfPol %v", err)
		d.Set("relation_infra_rs_cdp_if_pol", "")

	} else {
		d.Set("relation_infra_rs_cdp_if_pol", infraRsCdpIfPolData.(string))
	}

	infraRsL2IfPolData, err := aciClient.ReadRelationinfraRsL2IfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsL2IfPol %v", err)
		d.Set("relation_infra_rs_l2_if_pol", "")

	} else {
		d.Set("relation_infra_rs_l2_if_pol", infraRsL2IfPolData.(string))
	}

	infraRsStpIfPolData, err := aciClient.ReadRelationinfraRsStpIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsStpIfPol %v", err)
		d.Set("relation_infra_rs_stp_if_pol", "")

	} else {
		d.Set("relation_infra_rs_stp_if_pol", infraRsStpIfPolData.(string))
	}

	infraRsQosEgressDppIfPolData, err := aciClient.ReadRelationinfraRsQosEgressDppIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsQosEgressDppIfPol %v", err)
		d.Set("relation_infra_rs_qos_egress_dpp_if_pol", "")

	} else {
		d.Set("relation_infra_rs_qos_egress_dpp_if_pol", infraRsQosEgressDppIfPolData.(string))
	}

	infraRsAttEntPData, err := aciClient.ReadRelationinfraRsAttEntPFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsAttEntP %v", err)
		d.Set("relation_infra_rs_att_ent_p", "")

	} else {
		d.Set("relation_infra_rs_att_ent_p", infraRsAttEntPData.(string))
	}

	infraRsL2InstPolData, err := aciClient.ReadRelationinfraRsL2InstPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsL2InstPol %v", err)
		d.Set("relation_infra_rs_l2_inst_pol", "")

	} else {
		d.Set("relation_infra_rs_l2_inst_pol", infraRsL2InstPolData.(string))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciLeafAccessPortPolicyGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraAccPortGrp")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
