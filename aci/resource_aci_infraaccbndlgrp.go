package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciPCVPCInterfacePolicyGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciPCVPCInterfacePolicyGroupCreate,
		Update: resourceAciPCVPCInterfacePolicyGroupUpdate,
		Read:   resourceAciPCVPCInterfacePolicyGroupRead,
		Delete: resourceAciPCVPCInterfacePolicyGroupDelete,

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

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"relation_infra_rs_acc_bndl_grp_to_aggr_if": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_infra_rs_stormctrl_if_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_lldp_if_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_macsec_if_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_qos_dpp_if_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_h_if_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_netflow_monitor_pol": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tn_netflow_monitor_pol_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"flt_type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"relation_infra_rs_l2_port_auth_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_mcp_if_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_l2_port_security_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_copp_if_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_span_v_dest_grp": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_infra_rs_lacp_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_cdp_if_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_qos_pfc_if_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_qos_sd_if_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_mon_if_infra_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_fc_if_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_qos_ingress_dpp_if_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_qos_egress_dpp_if_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_l2_if_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_stp_if_pol": &schema.Schema{
				Type: schema.TypeString,

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
func getRemotePCVPCInterfacePolicyGroup(client *client.Client, dn string) (*models.PCVPCInterfacePolicyGroup, error) {
	infraAccBndlGrpCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraAccBndlGrp := models.PCVPCInterfacePolicyGroupFromContainer(infraAccBndlGrpCont)

	if infraAccBndlGrp.DistinguishedName == "" {
		return nil, fmt.Errorf("PCVPCInterfacePolicyGroup %s not found", infraAccBndlGrp.DistinguishedName)
	}

	return infraAccBndlGrp, nil
}

func setPCVPCInterfacePolicyGroupAttributes(infraAccBndlGrp *models.PCVPCInterfacePolicyGroup, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(infraAccBndlGrp.DistinguishedName)
	d.Set("description", infraAccBndlGrp.Description)
	infraAccBndlGrpMap, _ := infraAccBndlGrp.ToMap()

	d.Set("name", infraAccBndlGrpMap["name"])

	d.Set("annotation", infraAccBndlGrpMap["annotation"])
	d.Set("lag_t", infraAccBndlGrpMap["lagT"])
	d.Set("name_alias", infraAccBndlGrpMap["nameAlias"])
	return d
}

func resourceAciPCVPCInterfacePolicyGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraAccBndlGrp, err := getRemotePCVPCInterfacePolicyGroup(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setPCVPCInterfacePolicyGroupAttributes(infraAccBndlGrp, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciPCVPCInterfacePolicyGroupCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] PCVPCInterfacePolicyGroup: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraAccBndlGrpAttr := models.PCVPCInterfacePolicyGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraAccBndlGrpAttr.Annotation = Annotation.(string)
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
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if relationToinfraRsSpanVSrcGrp, ok := d.GetOk("relation_infra_rs_span_v_src_grp"); ok {
		relationParamList := toStringList(relationToinfraRsSpanVSrcGrp.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationinfraRsSpanVSrcGrpFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_infra_rs_span_v_src_grp")
			d.Partial(false)
		}
	}
	if relationToinfraRsAccBndlGrpToAggrIf, ok := d.GetOk("relation_infra_rs_acc_bndl_grp_to_aggr_if"); ok {
		relationParamList := toStringList(relationToinfraRsAccBndlGrpToAggrIf.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationinfraRsAccBndlGrpToAggrIfFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_infra_rs_acc_bndl_grp_to_aggr_if")
			d.Partial(false)
		}
	}
	if relationToinfraRsStormctrlIfPol, ok := d.GetOk("relation_infra_rs_stormctrl_if_pol"); ok {
		relationParam := relationToinfraRsStormctrlIfPol.(string)
		err = aciClient.CreateRelationinfraRsStormctrlIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_stormctrl_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsLldpIfPol, ok := d.GetOk("relation_infra_rs_lldp_if_pol"); ok {
		relationParam := relationToinfraRsLldpIfPol.(string)
		err = aciClient.CreateRelationinfraRsLldpIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_lldp_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsMacsecIfPol, ok := d.GetOk("relation_infra_rs_macsec_if_pol"); ok {
		relationParam := relationToinfraRsMacsecIfPol.(string)
		err = aciClient.CreateRelationinfraRsMacsecIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_macsec_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsQosDppIfPol, ok := d.GetOk("relation_infra_rs_qos_dpp_if_pol"); ok {
		relationParam := relationToinfraRsQosDppIfPol.(string)
		err = aciClient.CreateRelationinfraRsQosDppIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_qos_dpp_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsHIfPol, ok := d.GetOk("relation_infra_rs_h_if_pol"); ok {
		relationParam := relationToinfraRsHIfPol.(string)
		err = aciClient.CreateRelationinfraRsHIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_h_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsNetflowMonitorPol, ok := d.GetOk("relation_infra_rs_netflow_monitor_pol"); ok {

		relationParamList := relationToinfraRsNetflowMonitorPol.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationinfraRsNetflowMonitorPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, paramMap["tn_netflow_monitor_pol_name"].(string), paramMap["flt_type"].(string))
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_infra_rs_netflow_monitor_pol")
			d.Partial(false)
		}

	}
	if relationToinfraRsL2PortAuthPol, ok := d.GetOk("relation_infra_rs_l2_port_auth_pol"); ok {
		relationParam := relationToinfraRsL2PortAuthPol.(string)
		err = aciClient.CreateRelationinfraRsL2PortAuthPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_l2_port_auth_pol")
		d.Partial(false)

	}
	if relationToinfraRsMcpIfPol, ok := d.GetOk("relation_infra_rs_mcp_if_pol"); ok {
		relationParam := relationToinfraRsMcpIfPol.(string)
		err = aciClient.CreateRelationinfraRsMcpIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_mcp_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsL2PortSecurityPol, ok := d.GetOk("relation_infra_rs_l2_port_security_pol"); ok {
		relationParam := relationToinfraRsL2PortSecurityPol.(string)
		err = aciClient.CreateRelationinfraRsL2PortSecurityPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_l2_port_security_pol")
		d.Partial(false)

	}
	if relationToinfraRsCoppIfPol, ok := d.GetOk("relation_infra_rs_copp_if_pol"); ok {
		relationParam := relationToinfraRsCoppIfPol.(string)
		err = aciClient.CreateRelationinfraRsCoppIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_copp_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsSpanVDestGrp, ok := d.GetOk("relation_infra_rs_span_v_dest_grp"); ok {
		relationParamList := toStringList(relationToinfraRsSpanVDestGrp.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationinfraRsSpanVDestGrpFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_infra_rs_span_v_dest_grp")
			d.Partial(false)
		}
	}
	if relationToinfraRsLacpPol, ok := d.GetOk("relation_infra_rs_lacp_pol"); ok {
		relationParam := relationToinfraRsLacpPol.(string)
		err = aciClient.CreateRelationinfraRsLacpPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_lacp_pol")
		d.Partial(false)

	}
	if relationToinfraRsCdpIfPol, ok := d.GetOk("relation_infra_rs_cdp_if_pol"); ok {
		relationParam := relationToinfraRsCdpIfPol.(string)
		err = aciClient.CreateRelationinfraRsCdpIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_cdp_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsQosPfcIfPol, ok := d.GetOk("relation_infra_rs_qos_pfc_if_pol"); ok {
		relationParam := relationToinfraRsQosPfcIfPol.(string)
		err = aciClient.CreateRelationinfraRsQosPfcIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_qos_pfc_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsQosSdIfPol, ok := d.GetOk("relation_infra_rs_qos_sd_if_pol"); ok {
		relationParam := relationToinfraRsQosSdIfPol.(string)
		err = aciClient.CreateRelationinfraRsQosSdIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_qos_sd_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsMonIfInfraPol, ok := d.GetOk("relation_infra_rs_mon_if_infra_pol"); ok {
		relationParam := relationToinfraRsMonIfInfraPol.(string)
		err = aciClient.CreateRelationinfraRsMonIfInfraPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_mon_if_infra_pol")
		d.Partial(false)

	}
	if relationToinfraRsFcIfPol, ok := d.GetOk("relation_infra_rs_fc_if_pol"); ok {
		relationParam := relationToinfraRsFcIfPol.(string)
		err = aciClient.CreateRelationinfraRsFcIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_fc_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsQosIngressDppIfPol, ok := d.GetOk("relation_infra_rs_qos_ingress_dpp_if_pol"); ok {
		relationParam := relationToinfraRsQosIngressDppIfPol.(string)
		err = aciClient.CreateRelationinfraRsQosIngressDppIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_qos_ingress_dpp_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsQosEgressDppIfPol, ok := d.GetOk("relation_infra_rs_qos_egress_dpp_if_pol"); ok {
		relationParam := relationToinfraRsQosEgressDppIfPol.(string)
		err = aciClient.CreateRelationinfraRsQosEgressDppIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_qos_egress_dpp_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsL2IfPol, ok := d.GetOk("relation_infra_rs_l2_if_pol"); ok {
		relationParam := relationToinfraRsL2IfPol.(string)
		err = aciClient.CreateRelationinfraRsL2IfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_l2_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsStpIfPol, ok := d.GetOk("relation_infra_rs_stp_if_pol"); ok {
		relationParam := relationToinfraRsStpIfPol.(string)
		err = aciClient.CreateRelationinfraRsStpIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_stp_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsAttEntP, ok := d.GetOk("relation_infra_rs_att_ent_p"); ok {
		relationParam := relationToinfraRsAttEntP.(string)
		err = aciClient.CreateRelationinfraRsAttEntPFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_att_ent_p")
		d.Partial(false)

	}
	if relationToinfraRsL2InstPol, ok := d.GetOk("relation_infra_rs_l2_inst_pol"); ok {
		relationParam := relationToinfraRsL2InstPol.(string)
		err = aciClient.CreateRelationinfraRsL2InstPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_l2_inst_pol")
		d.Partial(false)

	}

	d.SetId(infraAccBndlGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciPCVPCInterfacePolicyGroupRead(d, m)
}

func resourceAciPCVPCInterfacePolicyGroupUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] PCVPCInterfacePolicyGroup: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraAccBndlGrpAttr := models.PCVPCInterfacePolicyGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraAccBndlGrpAttr.Annotation = Annotation.(string)
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
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if d.HasChange("relation_infra_rs_span_v_src_grp") {
		oldRel, newRel := d.GetChange("relation_infra_rs_span_v_src_grp")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationinfraRsSpanVSrcGrpFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationinfraRsSpanVSrcGrpFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relDn)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_infra_rs_span_v_src_grp")
			d.Partial(false)

		}

	}
	if d.HasChange("relation_infra_rs_acc_bndl_grp_to_aggr_if") {
		oldRel, newRel := d.GetChange("relation_infra_rs_acc_bndl_grp_to_aggr_if")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationinfraRsAccBndlGrpToAggrIfFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relDn)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_infra_rs_acc_bndl_grp_to_aggr_if")
			d.Partial(false)

		}

	}
	if d.HasChange("relation_infra_rs_stormctrl_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_stormctrl_if_pol")
		err = aciClient.CreateRelationinfraRsStormctrlIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_stormctrl_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_lldp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_lldp_if_pol")
		err = aciClient.CreateRelationinfraRsLldpIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_lldp_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_macsec_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_macsec_if_pol")
		err = aciClient.CreateRelationinfraRsMacsecIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_macsec_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_qos_dpp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_qos_dpp_if_pol")
		err = aciClient.CreateRelationinfraRsQosDppIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_qos_dpp_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_h_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_h_if_pol")
		err = aciClient.CreateRelationinfraRsHIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_h_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_netflow_monitor_pol") {
		oldRel, newRel := d.GetChange("relation_infra_rs_netflow_monitor_pol")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.DeleteRelationinfraRsNetflowMonitorPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, paramMap["tn_netflow_monitor_pol_name"].(string), paramMap["flt_type"].(string))
			if err != nil {
				return err
			}

		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationinfraRsNetflowMonitorPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, paramMap["tn_netflow_monitor_pol_name"].(string), paramMap["flt_type"].(string))
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_infra_rs_netflow_monitor_pol")
			d.Partial(false)
		}

	}
	if d.HasChange("relation_infra_rs_l2_port_auth_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_l2_port_auth_pol")
		err = aciClient.CreateRelationinfraRsL2PortAuthPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_l2_port_auth_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_mcp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_mcp_if_pol")
		err = aciClient.CreateRelationinfraRsMcpIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_mcp_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_l2_port_security_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_l2_port_security_pol")
		err = aciClient.CreateRelationinfraRsL2PortSecurityPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_l2_port_security_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_copp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_copp_if_pol")
		err = aciClient.CreateRelationinfraRsCoppIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_copp_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_span_v_dest_grp") {
		oldRel, newRel := d.GetChange("relation_infra_rs_span_v_dest_grp")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationinfraRsSpanVDestGrpFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationinfraRsSpanVDestGrpFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, relDn)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_infra_rs_span_v_dest_grp")
			d.Partial(false)

		}

	}
	if d.HasChange("relation_infra_rs_lacp_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_lacp_pol")
		err = aciClient.CreateRelationinfraRsLacpPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_lacp_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_cdp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_cdp_if_pol")
		err = aciClient.CreateRelationinfraRsCdpIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_cdp_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_qos_pfc_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_qos_pfc_if_pol")
		err = aciClient.CreateRelationinfraRsQosPfcIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_qos_pfc_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_qos_sd_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_qos_sd_if_pol")
		err = aciClient.CreateRelationinfraRsQosSdIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_qos_sd_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_mon_if_infra_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_mon_if_infra_pol")
		err = aciClient.CreateRelationinfraRsMonIfInfraPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_mon_if_infra_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_fc_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_fc_if_pol")
		err = aciClient.CreateRelationinfraRsFcIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_fc_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_qos_ingress_dpp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_qos_ingress_dpp_if_pol")
		err = aciClient.CreateRelationinfraRsQosIngressDppIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_qos_ingress_dpp_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_qos_egress_dpp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_qos_egress_dpp_if_pol")
		err = aciClient.CreateRelationinfraRsQosEgressDppIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_qos_egress_dpp_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_l2_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_l2_if_pol")
		err = aciClient.CreateRelationinfraRsL2IfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_l2_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_stp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_stp_if_pol")
		err = aciClient.CreateRelationinfraRsStpIfPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_stp_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_att_ent_p") {
		_, newRelParam := d.GetChange("relation_infra_rs_att_ent_p")
		err = aciClient.DeleteRelationinfraRsAttEntPFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationinfraRsAttEntPFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_att_ent_p")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_l2_inst_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_l2_inst_pol")
		err = aciClient.DeleteRelationinfraRsL2InstPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationinfraRsL2InstPolFromPCVPCInterfacePolicyGroup(infraAccBndlGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_l2_inst_pol")
		d.Partial(false)

	}

	d.SetId(infraAccBndlGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciPCVPCInterfacePolicyGroupRead(d, m)

}

func resourceAciPCVPCInterfacePolicyGroupRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraAccBndlGrp, err := getRemotePCVPCInterfacePolicyGroup(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setPCVPCInterfacePolicyGroupAttributes(infraAccBndlGrp, d)

	infraRsSpanVSrcGrpData, err := aciClient.ReadRelationinfraRsSpanVSrcGrpFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsSpanVSrcGrp %v", err)

	} else {
		d.Set("relation_infra_rs_span_v_src_grp", infraRsSpanVSrcGrpData)
	}

	infraRsAccBndlGrpToAggrIfData, err := aciClient.ReadRelationinfraRsAccBndlGrpToAggrIfFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsAccBndlGrpToAggrIf %v", err)

	} else {
		d.Set("relation_infra_rs_acc_bndl_grp_to_aggr_if", infraRsAccBndlGrpToAggrIfData)
	}

	infraRsStormctrlIfPolData, err := aciClient.ReadRelationinfraRsStormctrlIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsStormctrlIfPol %v", err)

	} else {
		d.Set("relation_infra_rs_stormctrl_if_pol", infraRsStormctrlIfPolData)
	}

	infraRsLldpIfPolData, err := aciClient.ReadRelationinfraRsLldpIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsLldpIfPol %v", err)

	} else {
		d.Set("relation_infra_rs_lldp_if_pol", infraRsLldpIfPolData)
	}

	infraRsMacsecIfPolData, err := aciClient.ReadRelationinfraRsMacsecIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsMacsecIfPol %v", err)

	} else {
		d.Set("relation_infra_rs_macsec_if_pol", infraRsMacsecIfPolData)
	}

	infraRsQosDppIfPolData, err := aciClient.ReadRelationinfraRsQosDppIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsQosDppIfPol %v", err)

	} else {
		d.Set("relation_infra_rs_qos_dpp_if_pol", infraRsQosDppIfPolData)
	}

	infraRsHIfPolData, err := aciClient.ReadRelationinfraRsHIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsHIfPol %v", err)

	} else {
		d.Set("relation_infra_rs_h_if_pol", infraRsHIfPolData)
	}

	infraRsNetflowMonitorPolData, err := aciClient.ReadRelationinfraRsNetflowMonitorPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsNetflowMonitorPol %v", err)

	} else {
		d.Set("relation_infra_rs_netflow_monitor_pol", infraRsNetflowMonitorPolData)
	}

	infraRsL2PortAuthPolData, err := aciClient.ReadRelationinfraRsL2PortAuthPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsL2PortAuthPol %v", err)

	} else {
		d.Set("relation_infra_rs_l2_port_auth_pol", infraRsL2PortAuthPolData)
	}

	infraRsMcpIfPolData, err := aciClient.ReadRelationinfraRsMcpIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsMcpIfPol %v", err)

	} else {
		d.Set("relation_infra_rs_mcp_if_pol", infraRsMcpIfPolData)
	}

	infraRsL2PortSecurityPolData, err := aciClient.ReadRelationinfraRsL2PortSecurityPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsL2PortSecurityPol %v", err)

	} else {
		d.Set("relation_infra_rs_l2_port_security_pol", infraRsL2PortSecurityPolData)
	}

	infraRsCoppIfPolData, err := aciClient.ReadRelationinfraRsCoppIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsCoppIfPol %v", err)

	} else {
		d.Set("relation_infra_rs_copp_if_pol", infraRsCoppIfPolData)
	}

	infraRsSpanVDestGrpData, err := aciClient.ReadRelationinfraRsSpanVDestGrpFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsSpanVDestGrp %v", err)

	} else {
		d.Set("relation_infra_rs_span_v_dest_grp", infraRsSpanVDestGrpData)
	}

	infraRsLacpPolData, err := aciClient.ReadRelationinfraRsLacpPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsLacpPol %v", err)

	} else {
		d.Set("relation_infra_rs_lacp_pol", infraRsLacpPolData)
	}

	infraRsCdpIfPolData, err := aciClient.ReadRelationinfraRsCdpIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsCdpIfPol %v", err)

	} else {
		d.Set("relation_infra_rs_cdp_if_pol", infraRsCdpIfPolData)
	}

	infraRsQosPfcIfPolData, err := aciClient.ReadRelationinfraRsQosPfcIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsQosPfcIfPol %v", err)

	} else {
		d.Set("relation_infra_rs_qos_pfc_if_pol", infraRsQosPfcIfPolData)
	}

	infraRsQosSdIfPolData, err := aciClient.ReadRelationinfraRsQosSdIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsQosSdIfPol %v", err)

	} else {
		d.Set("relation_infra_rs_qos_sd_if_pol", infraRsQosSdIfPolData)
	}

	infraRsMonIfInfraPolData, err := aciClient.ReadRelationinfraRsMonIfInfraPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsMonIfInfraPol %v", err)

	} else {
		d.Set("relation_infra_rs_mon_if_infra_pol", infraRsMonIfInfraPolData)
	}

	infraRsFcIfPolData, err := aciClient.ReadRelationinfraRsFcIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsFcIfPol %v", err)

	} else {
		d.Set("relation_infra_rs_fc_if_pol", infraRsFcIfPolData)
	}

	infraRsQosIngressDppIfPolData, err := aciClient.ReadRelationinfraRsQosIngressDppIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsQosIngressDppIfPol %v", err)

	} else {
		d.Set("relation_infra_rs_qos_ingress_dpp_if_pol", infraRsQosIngressDppIfPolData)
	}

	infraRsQosEgressDppIfPolData, err := aciClient.ReadRelationinfraRsQosEgressDppIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsQosEgressDppIfPol %v", err)

	} else {
		d.Set("relation_infra_rs_qos_egress_dpp_if_pol", infraRsQosEgressDppIfPolData)
	}

	infraRsL2IfPolData, err := aciClient.ReadRelationinfraRsL2IfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsL2IfPol %v", err)

	} else {
		d.Set("relation_infra_rs_l2_if_pol", infraRsL2IfPolData)
	}

	infraRsStpIfPolData, err := aciClient.ReadRelationinfraRsStpIfPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsStpIfPol %v", err)

	} else {
		d.Set("relation_infra_rs_stp_if_pol", infraRsStpIfPolData)
	}

	infraRsAttEntPData, err := aciClient.ReadRelationinfraRsAttEntPFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsAttEntP %v", err)

	} else {
		d.Set("relation_infra_rs_att_ent_p", infraRsAttEntPData)
	}

	infraRsL2InstPolData, err := aciClient.ReadRelationinfraRsL2InstPolFromPCVPCInterfacePolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsL2InstPol %v", err)

	} else {
		d.Set("relation_infra_rs_l2_inst_pol", infraRsL2InstPolData)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciPCVPCInterfacePolicyGroupDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraAccBndlGrp")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
