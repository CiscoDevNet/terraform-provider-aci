package aci

import (
	"fmt"
	"log"
	"reflect"
	"sort"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciLeafAccessPortPolicyGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciLeafAccessPortPolicyGroupCreate,
		Update: resourceAciLeafAccessPortPolicyGroupUpdate,
		Read:   resourceAciLeafAccessPortPolicyGroupRead,
		Delete: resourceAciLeafAccessPortPolicyGroupDelete,

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
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_poe_if_pol": &schema.Schema{
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
			"relation_infra_rs_dwdm_if_pol": &schema.Schema{
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
			"relation_infra_rs_cdp_if_pol": &schema.Schema{
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
			"relation_infra_rs_qos_egress_dpp_if_pol": &schema.Schema{
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

func setLeafAccessPortPolicyGroupAttributes(infraAccPortGrp *models.LeafAccessPortPolicyGroup, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(infraAccPortGrp.DistinguishedName)
	d.Set("description", infraAccPortGrp.Description)
	infraAccPortGrpMap, _ := infraAccPortGrp.ToMap()

	d.Set("name", infraAccPortGrpMap["name"])

	d.Set("annotation", infraAccPortGrpMap["annotation"])
	d.Set("name_alias", infraAccPortGrpMap["nameAlias"])
	return d
}

func resourceAciLeafAccessPortPolicyGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraAccPortGrp, err := getRemoteLeafAccessPortPolicyGroup(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setLeafAccessPortPolicyGroupAttributes(infraAccPortGrp, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLeafAccessPortPolicyGroupCreate(d *schema.ResourceData, m interface{}) error {
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
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if relationToinfraRsSpanVSrcGrp, ok := d.GetOk("relation_infra_rs_span_v_src_grp"); ok {
		relationParamList := toStringList(relationToinfraRsSpanVSrcGrp.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationinfraRsSpanVSrcGrpFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_infra_rs_span_v_src_grp")
			d.Partial(false)
		}
	}
	if relationToinfraRsStormctrlIfPol, ok := d.GetOk("relation_infra_rs_stormctrl_if_pol"); ok {
		relationParam := relationToinfraRsStormctrlIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsStormctrlIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_stormctrl_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsPoeIfPol, ok := d.GetOk("relation_infra_rs_poe_if_pol"); ok {
		relationParam := relationToinfraRsPoeIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsPoeIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_poe_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsLldpIfPol, ok := d.GetOk("relation_infra_rs_lldp_if_pol"); ok {
		relationParam := relationToinfraRsLldpIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsLldpIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_lldp_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsMacsecIfPol, ok := d.GetOk("relation_infra_rs_macsec_if_pol"); ok {
		relationParam := relationToinfraRsMacsecIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsMacsecIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_macsec_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsQosDppIfPol, ok := d.GetOk("relation_infra_rs_qos_dpp_if_pol"); ok {
		relationParam := relationToinfraRsQosDppIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsQosDppIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_qos_dpp_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsHIfPol, ok := d.GetOk("relation_infra_rs_h_if_pol"); ok {
		relationParam := relationToinfraRsHIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsHIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
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
			err = aciClient.CreateRelationinfraRsNetflowMonitorPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, paramMap["tn_netflow_monitor_pol_name"].(string), paramMap["flt_type"].(string))
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
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsL2PortAuthPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_l2_port_auth_pol")
		d.Partial(false)

	}
	if relationToinfraRsMcpIfPol, ok := d.GetOk("relation_infra_rs_mcp_if_pol"); ok {
		relationParam := relationToinfraRsMcpIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsMcpIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_mcp_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsL2PortSecurityPol, ok := d.GetOk("relation_infra_rs_l2_port_security_pol"); ok {
		relationParam := relationToinfraRsL2PortSecurityPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsL2PortSecurityPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_l2_port_security_pol")
		d.Partial(false)

	}
	if relationToinfraRsCoppIfPol, ok := d.GetOk("relation_infra_rs_copp_if_pol"); ok {
		relationParam := relationToinfraRsCoppIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsCoppIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
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
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationinfraRsSpanVDestGrpFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_infra_rs_span_v_dest_grp")
			d.Partial(false)
		}
	}
	if relationToinfraRsDwdmIfPol, ok := d.GetOk("relation_infra_rs_dwdm_if_pol"); ok {
		relationParam := relationToinfraRsDwdmIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsDwdmIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_dwdm_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsQosPfcIfPol, ok := d.GetOk("relation_infra_rs_qos_pfc_if_pol"); ok {
		relationParam := relationToinfraRsQosPfcIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsQosPfcIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_qos_pfc_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsQosSdIfPol, ok := d.GetOk("relation_infra_rs_qos_sd_if_pol"); ok {
		relationParam := relationToinfraRsQosSdIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsQosSdIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_qos_sd_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsMonIfInfraPol, ok := d.GetOk("relation_infra_rs_mon_if_infra_pol"); ok {
		relationParam := relationToinfraRsMonIfInfraPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsMonIfInfraPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_mon_if_infra_pol")
		d.Partial(false)

	}
	if relationToinfraRsFcIfPol, ok := d.GetOk("relation_infra_rs_fc_if_pol"); ok {
		relationParam := relationToinfraRsFcIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsFcIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_fc_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsQosIngressDppIfPol, ok := d.GetOk("relation_infra_rs_qos_ingress_dpp_if_pol"); ok {
		relationParam := relationToinfraRsQosIngressDppIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsQosIngressDppIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_qos_ingress_dpp_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsCdpIfPol, ok := d.GetOk("relation_infra_rs_cdp_if_pol"); ok {
		relationParam := relationToinfraRsCdpIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsCdpIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_cdp_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsL2IfPol, ok := d.GetOk("relation_infra_rs_l2_if_pol"); ok {
		relationParam := relationToinfraRsL2IfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsL2IfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_l2_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsStpIfPol, ok := d.GetOk("relation_infra_rs_stp_if_pol"); ok {
		relationParam := relationToinfraRsStpIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsStpIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_stp_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsQosEgressDppIfPol, ok := d.GetOk("relation_infra_rs_qos_egress_dpp_if_pol"); ok {
		relationParam := relationToinfraRsQosEgressDppIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsQosEgressDppIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_qos_egress_dpp_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsAttEntP, ok := d.GetOk("relation_infra_rs_att_ent_p"); ok {
		relationParam := relationToinfraRsAttEntP.(string)
		err = aciClient.CreateRelationinfraRsAttEntPFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_att_ent_p")
		d.Partial(false)

	}
	if relationToinfraRsL2InstPol, ok := d.GetOk("relation_infra_rs_l2_inst_pol"); ok {
		relationParam := relationToinfraRsL2InstPol.(string)
		err = aciClient.CreateRelationinfraRsL2InstPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_l2_inst_pol")
		d.Partial(false)

	}

	d.SetId(infraAccPortGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciLeafAccessPortPolicyGroupRead(d, m)
}

func resourceAciLeafAccessPortPolicyGroupUpdate(d *schema.ResourceData, m interface{}) error {
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
			relDnName := GetMOName(relDn)
			err = aciClient.DeleteRelationinfraRsSpanVSrcGrpFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relDnName)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationinfraRsSpanVSrcGrpFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relDnName)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_infra_rs_span_v_src_grp")
			d.Partial(false)

		}

	}
	if d.HasChange("relation_infra_rs_stormctrl_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_stormctrl_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsStormctrlIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_stormctrl_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_poe_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_poe_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationinfraRsPoeIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationinfraRsPoeIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_poe_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_lldp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_lldp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsLldpIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_lldp_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_macsec_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_macsec_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsMacsecIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_macsec_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_qos_dpp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_qos_dpp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsQosDppIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_qos_dpp_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_h_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_h_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsHIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
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
			err = aciClient.DeleteRelationinfraRsNetflowMonitorPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, paramMap["tn_netflow_monitor_pol_name"].(string), paramMap["flt_type"].(string))
			if err != nil {
				return err
			}

		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationinfraRsNetflowMonitorPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, paramMap["tn_netflow_monitor_pol_name"].(string), paramMap["flt_type"].(string))
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
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsL2PortAuthPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_l2_port_auth_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_mcp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_mcp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsMcpIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_mcp_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_l2_port_security_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_l2_port_security_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsL2PortSecurityPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_l2_port_security_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_copp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_copp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsCoppIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
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
			relDnName := GetMOName(relDn)
			err = aciClient.DeleteRelationinfraRsSpanVDestGrpFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relDnName)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationinfraRsSpanVDestGrpFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, relDnName)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_infra_rs_span_v_dest_grp")
			d.Partial(false)

		}

	}
	if d.HasChange("relation_infra_rs_dwdm_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_dwdm_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsDwdmIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_dwdm_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_qos_pfc_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_qos_pfc_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsQosPfcIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_qos_pfc_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_qos_sd_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_qos_sd_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsQosSdIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_qos_sd_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_mon_if_infra_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_mon_if_infra_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsMonIfInfraPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_mon_if_infra_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_fc_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_fc_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsFcIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_fc_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_qos_ingress_dpp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_qos_ingress_dpp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsQosIngressDppIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_qos_ingress_dpp_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_cdp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_cdp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsCdpIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_cdp_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_l2_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_l2_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsL2IfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_l2_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_stp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_stp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsStpIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_stp_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_qos_egress_dpp_if_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_qos_egress_dpp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsQosEgressDppIfPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_qos_egress_dpp_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_att_ent_p") {
		_, newRelParam := d.GetChange("relation_infra_rs_att_ent_p")
		err = aciClient.DeleteRelationinfraRsAttEntPFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationinfraRsAttEntPFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_att_ent_p")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_l2_inst_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_l2_inst_pol")
		err = aciClient.DeleteRelationinfraRsL2InstPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationinfraRsL2InstPolFromLeafAccessPortPolicyGroup(infraAccPortGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_l2_inst_pol")
		d.Partial(false)

	}

	d.SetId(infraAccPortGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciLeafAccessPortPolicyGroupRead(d, m)

}

func resourceAciLeafAccessPortPolicyGroupRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraAccPortGrp, err := getRemoteLeafAccessPortPolicyGroup(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setLeafAccessPortPolicyGroupAttributes(infraAccPortGrp, d)

	infraRsSpanVSrcGrpData, err := aciClient.ReadRelationinfraRsSpanVSrcGrpFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsSpanVSrcGrp %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_span_v_src_grp"); ok {
			relationParamList := toStringList(d.Get("relation_infra_rs_span_v_src_grp").(*schema.Set).List())
			tfList := make([]string, 0, 1)
			for _, relationParam := range relationParamList {
				relationParamName := GetMOName(relationParam)
				tfList = append(tfList, relationParamName)
			}
			infraRsSpanVSrcGrpDataList := toStringList(infraRsSpanVSrcGrpData.(*schema.Set).List())
			sort.Strings(tfList)
			sort.Strings(infraRsSpanVSrcGrpDataList)

			if !reflect.DeepEqual(tfList, infraRsSpanVSrcGrpDataList) {
				d.Set("relation_infra_rs_span_v_src_grp", make([]string, 0, 1))
			}
		}
	}

	infraRsStormctrlIfPolData, err := aciClient.ReadRelationinfraRsStormctrlIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsStormctrlIfPol %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_stormctrl_if_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_stormctrl_if_pol").(string))
			if tfName != infraRsStormctrlIfPolData {
				d.Set("relation_infra_rs_stormctrl_if_pol", "")
			}
		}
	}

	infraRsPoeIfPolData, err := aciClient.ReadRelationinfraRsPoeIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsPoeIfPol %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_poe_if_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_poe_if_pol").(string))
			if tfName != infraRsPoeIfPolData {
				d.Set("relation_infra_rs_poe_if_pol", "")
			}
		}
	}

	infraRsLldpIfPolData, err := aciClient.ReadRelationinfraRsLldpIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsLldpIfPol %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_lldp_if_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_lldp_if_pol").(string))
			if tfName != infraRsLldpIfPolData {
				d.Set("relation_infra_rs_lldp_if_pol", "")
			}
		}
	}

	infraRsMacsecIfPolData, err := aciClient.ReadRelationinfraRsMacsecIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsMacsecIfPol %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_macsec_if_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_macsec_if_pol").(string))
			if tfName != infraRsMacsecIfPolData {
				d.Set("relation_infra_rs_macsec_if_pol", "")
			}
		}
	}

	infraRsQosDppIfPolData, err := aciClient.ReadRelationinfraRsQosDppIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsQosDppIfPol %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_qos_dpp_if_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_qos_dpp_if_pol").(string))
			if tfName != infraRsQosDppIfPolData {
				d.Set("relation_infra_rs_qos_dpp_if_pol", "")
			}
		}
	}

	infraRsHIfPolData, err := aciClient.ReadRelationinfraRsHIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsHIfPol %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_h_if_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_h_if_pol").(string))
			if tfName != infraRsHIfPolData {
				d.Set("relation_infra_rs_h_if_pol", "")
			}
		}
	}

	infraRsNetflowMonitorPolData, err := aciClient.ReadRelationinfraRsNetflowMonitorPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsNetflowMonitorPol %v", err)

	} else {
		d.Set("relation_infra_rs_netflow_monitor_pol", infraRsNetflowMonitorPolData)
	}

	infraRsL2PortAuthPolData, err := aciClient.ReadRelationinfraRsL2PortAuthPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsL2PortAuthPol %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_l2_port_auth_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_l2_port_auth_pol").(string))
			if tfName != infraRsL2PortAuthPolData {
				d.Set("relation_infra_rs_l2_port_auth_pol", "")
			}
		}
	}

	infraRsMcpIfPolData, err := aciClient.ReadRelationinfraRsMcpIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsMcpIfPol %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_mcp_if_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_mcp_if_pol").(string))
			if tfName != infraRsMcpIfPolData {
				d.Set("relation_infra_rs_mcp_if_pol", "")
			}
		}
	}

	infraRsL2PortSecurityPolData, err := aciClient.ReadRelationinfraRsL2PortSecurityPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsL2PortSecurityPol %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_l2_port_security_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_l2_port_security_pol").(string))
			if tfName != infraRsL2PortSecurityPolData {
				d.Set("relation_infra_rs_l2_port_security_pol", "")
			}
		}
	}

	infraRsCoppIfPolData, err := aciClient.ReadRelationinfraRsCoppIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsCoppIfPol %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_copp_if_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_copp_if_pol").(string))
			if tfName != infraRsCoppIfPolData {
				d.Set("relation_infra_rs_copp_if_pol", "")
			}
		}
	}

	infraRsSpanVDestGrpData, err := aciClient.ReadRelationinfraRsSpanVDestGrpFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsSpanVDestGrp %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_span_v_dest_grp"); ok {
			relationParamList := toStringList(d.Get("relation_infra_rs_span_v_dest_grp").(*schema.Set).List())
			tfList := make([]string, 0, 1)
			for _, relationParam := range relationParamList {
				relationParamName := GetMOName(relationParam)
				tfList = append(tfList, relationParamName)
			}
			infraRsSpanVDestGrpDataList := toStringList(infraRsSpanVDestGrpData.(*schema.Set).List())
			sort.Strings(tfList)
			sort.Strings(infraRsSpanVDestGrpDataList)

			if !reflect.DeepEqual(tfList, infraRsSpanVDestGrpDataList) {
				d.Set("relation_infra_rs_span_v_dest_grp", make([]string, 0, 1))
			}
		}
	}

	infraRsDwdmIfPolData, err := aciClient.ReadRelationinfraRsDwdmIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsDwdmIfPol %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_dwdm_if_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_dwdm_if_pol").(string))
			if tfName != infraRsDwdmIfPolData {
				d.Set("relation_infra_rs_dwdm_if_pol", "")
			}
		}
	}

	infraRsQosPfcIfPolData, err := aciClient.ReadRelationinfraRsQosPfcIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsQosPfcIfPol %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_qos_pfc_if_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_qos_pfc_if_pol").(string))
			if tfName != infraRsQosPfcIfPolData {
				d.Set("relation_infra_rs_qos_pfc_if_pol", "")
			}
		}
	}

	infraRsQosSdIfPolData, err := aciClient.ReadRelationinfraRsQosSdIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsQosSdIfPol %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_qos_sd_if_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_qos_sd_if_pol").(string))
			if tfName != infraRsQosSdIfPolData {
				d.Set("relation_infra_rs_qos_sd_if_pol", "")
			}
		}
	}

	infraRsMonIfInfraPolData, err := aciClient.ReadRelationinfraRsMonIfInfraPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsMonIfInfraPol %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_mon_if_infra_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_mon_if_infra_pol").(string))
			if tfName != infraRsMonIfInfraPolData {
				d.Set("relation_infra_rs_mon_if_infra_pol", "")
			}
		}
	}

	infraRsFcIfPolData, err := aciClient.ReadRelationinfraRsFcIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsFcIfPol %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_fc_if_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_fc_if_pol").(string))
			if tfName != infraRsFcIfPolData {
				d.Set("relation_infra_rs_fc_if_pol", "")
			}
		}
	}

	infraRsQosIngressDppIfPolData, err := aciClient.ReadRelationinfraRsQosIngressDppIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsQosIngressDppIfPol %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_qos_ingress_dpp_if_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_qos_ingress_dpp_if_pol").(string))
			if tfName != infraRsQosIngressDppIfPolData {
				d.Set("relation_infra_rs_qos_ingress_dpp_if_pol", "")
			}
		}
	}

	infraRsCdpIfPolData, err := aciClient.ReadRelationinfraRsCdpIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsCdpIfPol %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_cdp_if_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_cdp_if_pol").(string))
			if tfName != infraRsCdpIfPolData {
				d.Set("relation_infra_rs_cdp_if_pol", "")
			}
		}
	}

	infraRsL2IfPolData, err := aciClient.ReadRelationinfraRsL2IfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsL2IfPol %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_l2_if_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_l2_if_pol").(string))
			if tfName != infraRsL2IfPolData {
				d.Set("relation_infra_rs_l2_if_pol", "")
			}
		}
	}

	infraRsStpIfPolData, err := aciClient.ReadRelationinfraRsStpIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsStpIfPol %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_stp_if_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_stp_if_pol").(string))
			if tfName != infraRsStpIfPolData {
				d.Set("relation_infra_rs_stp_if_pol", "")
			}
		}
	}

	infraRsQosEgressDppIfPolData, err := aciClient.ReadRelationinfraRsQosEgressDppIfPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsQosEgressDppIfPol %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_qos_egress_dpp_if_pol"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_qos_egress_dpp_if_pol").(string))
			if tfName != infraRsQosEgressDppIfPolData {
				d.Set("relation_infra_rs_qos_egress_dpp_if_pol", "")
			}
		}
	}

	infraRsAttEntPData, err := aciClient.ReadRelationinfraRsAttEntPFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsAttEntP %v", err)

	} else {
		d.Set("relation_infra_rs_att_ent_p", infraRsAttEntPData)
	}

	infraRsL2InstPolData, err := aciClient.ReadRelationinfraRsL2InstPolFromLeafAccessPortPolicyGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsL2InstPol %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_l2_inst_pol"); ok {
			tfName := d.Get("relation_infra_rs_l2_inst_pol").(string)
			if tfName != infraRsL2InstPolData {
				d.Set("relation_infra_rs_l2_inst_pol", "")
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciLeafAccessPortPolicyGroupDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraAccPortGrp")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
