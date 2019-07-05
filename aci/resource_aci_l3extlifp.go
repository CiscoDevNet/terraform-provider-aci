package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAciLogicalInterfaceProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciLogicalInterfaceProfileCreate,
		Update: resourceAciLogicalInterfaceProfileUpdate,
		Read:   resourceAciLogicalInterfaceProfileRead,
		Delete: resourceAciLogicalInterfaceProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciLogicalInterfaceProfileImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"logical_node_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"prio": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_l3ext_rs_l_if_p_to_netflow_monitor_pol": &schema.Schema{
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
			"relation_l3ext_rs_path_l3_out_att": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_l3ext_rs_egress_qos_dpp_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_l3ext_rs_ingress_qos_dpp_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_l3ext_rs_l_if_p_cust_qos_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_l3ext_rs_arp_if_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_l3ext_rs_nd_if_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteLogicalInterfaceProfile(client *client.Client, dn string) (*models.LogicalInterfaceProfile, error) {
	l3extLIfPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extLIfP := models.LogicalInterfaceProfileFromContainer(l3extLIfPCont)

	if l3extLIfP.DistinguishedName == "" {
		return nil, fmt.Errorf("LogicalInterfaceProfile %s not found", l3extLIfP.DistinguishedName)
	}

	return l3extLIfP, nil
}

func setLogicalInterfaceProfileAttributes(l3extLIfP *models.LogicalInterfaceProfile, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(l3extLIfP.DistinguishedName)
	d.Set("description", l3extLIfP.Description)
	d.Set("logical_node_profile_dn", GetParentDn(l3extLIfP.DistinguishedName))
	l3extLIfPMap, _ := l3extLIfP.ToMap()

	d.Set("name", l3extLIfPMap["name"])

	d.Set("annotation", l3extLIfPMap["annotation"])
	d.Set("name_alias", l3extLIfPMap["nameAlias"])
	d.Set("prio", l3extLIfPMap["prio"])
	d.Set("tag", l3extLIfPMap["tag"])
	return d
}

func resourceAciLogicalInterfaceProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	l3extLIfP, err := getRemoteLogicalInterfaceProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setLogicalInterfaceProfileAttributes(l3extLIfP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLogicalInterfaceProfileCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LogicalInterfaceProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	LogicalNodeProfileDn := d.Get("logical_node_profile_dn").(string)

	l3extLIfPAttr := models.LogicalInterfaceProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extLIfPAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extLIfPAttr.NameAlias = NameAlias.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		l3extLIfPAttr.Prio = Prio.(string)
	}
	if Tag, ok := d.GetOk("tag"); ok {
		l3extLIfPAttr.Tag = Tag.(string)
	}
	l3extLIfP := models.NewLogicalInterfaceProfile(fmt.Sprintf("lifp-%s", name), LogicalNodeProfileDn, desc, l3extLIfPAttr)

	err := aciClient.Save(l3extLIfP)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if relationTol3extRsLIfPToNetflowMonitorPol, ok := d.GetOk("relation_l3ext_rs_l_if_p_to_netflow_monitor_pol"); ok {

		relationParamList := relationTol3extRsLIfPToNetflowMonitorPol.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationl3extRsLIfPToNetflowMonitorPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, paramMap["tn_netflow_monitor_pol_name"].(string), paramMap["flt_type"].(string))
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_l3ext_rs_l_if_p_to_netflow_monitor_pol")
			d.Partial(false)
		}

	}
	if relationTol3extRsPathL3OutAtt, ok := d.GetOk("relation_l3ext_rs_path_l3_out_att"); ok {
		relationParamList := toStringList(relationTol3extRsPathL3OutAtt.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationl3extRsPathL3OutAttFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_l3ext_rs_path_l3_out_att")
			d.Partial(false)
		}
	}
	if relationTol3extRsEgressQosDppPol, ok := d.GetOk("relation_l3ext_rs_egress_qos_dpp_pol"); ok {
		relationParam := relationTol3extRsEgressQosDppPol.(string)
		err = aciClient.CreateRelationl3extRsEgressQosDppPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_l3ext_rs_egress_qos_dpp_pol")
		d.Partial(false)

	}
	if relationTol3extRsIngressQosDppPol, ok := d.GetOk("relation_l3ext_rs_ingress_qos_dpp_pol"); ok {
		relationParam := relationTol3extRsIngressQosDppPol.(string)
		err = aciClient.CreateRelationl3extRsIngressQosDppPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_l3ext_rs_ingress_qos_dpp_pol")
		d.Partial(false)

	}
	if relationTol3extRsLIfPCustQosPol, ok := d.GetOk("relation_l3ext_rs_l_if_p_cust_qos_pol"); ok {
		relationParam := relationTol3extRsLIfPCustQosPol.(string)
		err = aciClient.CreateRelationl3extRsLIfPCustQosPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_l3ext_rs_l_if_p_cust_qos_pol")
		d.Partial(false)

	}
	if relationTol3extRsArpIfPol, ok := d.GetOk("relation_l3ext_rs_arp_if_pol"); ok {
		relationParam := relationTol3extRsArpIfPol.(string)
		err = aciClient.CreateRelationl3extRsArpIfPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_l3ext_rs_arp_if_pol")
		d.Partial(false)

	}
	if relationTol3extRsNdIfPol, ok := d.GetOk("relation_l3ext_rs_nd_if_pol"); ok {
		relationParam := relationTol3extRsNdIfPol.(string)
		err = aciClient.CreateRelationl3extRsNdIfPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_l3ext_rs_nd_if_pol")
		d.Partial(false)

	}

	d.SetId(l3extLIfP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciLogicalInterfaceProfileRead(d, m)
}

func resourceAciLogicalInterfaceProfileUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LogicalInterfaceProfile: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	LogicalNodeProfileDn := d.Get("logical_node_profile_dn").(string)

	l3extLIfPAttr := models.LogicalInterfaceProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extLIfPAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extLIfPAttr.NameAlias = NameAlias.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		l3extLIfPAttr.Prio = Prio.(string)
	}
	if Tag, ok := d.GetOk("tag"); ok {
		l3extLIfPAttr.Tag = Tag.(string)
	}
	l3extLIfP := models.NewLogicalInterfaceProfile(fmt.Sprintf("lifp-%s", name), LogicalNodeProfileDn, desc, l3extLIfPAttr)

	l3extLIfP.Status = "modified"

	err := aciClient.Save(l3extLIfP)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if d.HasChange("relation_l3ext_rs_l_if_p_to_netflow_monitor_pol") {
		oldRel, newRel := d.GetChange("relation_l3ext_rs_l_if_p_to_netflow_monitor_pol")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.DeleteRelationl3extRsLIfPToNetflowMonitorPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, paramMap["tn_netflow_monitor_pol_name"].(string), paramMap["flt_type"].(string))
			if err != nil {
				return err
			}

		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationl3extRsLIfPToNetflowMonitorPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, paramMap["tn_netflow_monitor_pol_name"].(string), paramMap["flt_type"].(string))
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_l3ext_rs_l_if_p_to_netflow_monitor_pol")
			d.Partial(false)
		}

	}
	if d.HasChange("relation_l3ext_rs_path_l3_out_att") {
		oldRel, newRel := d.GetChange("relation_l3ext_rs_path_l3_out_att")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationl3extRsPathL3OutAttFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationl3extRsPathL3OutAttFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, relDn)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_l3ext_rs_path_l3_out_att")
			d.Partial(false)

		}

	}
	if d.HasChange("relation_l3ext_rs_egress_qos_dpp_pol") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_egress_qos_dpp_pol")
		err = aciClient.CreateRelationl3extRsEgressQosDppPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_l3ext_rs_egress_qos_dpp_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_l3ext_rs_ingress_qos_dpp_pol") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_ingress_qos_dpp_pol")
		err = aciClient.CreateRelationl3extRsIngressQosDppPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_l3ext_rs_ingress_qos_dpp_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_l3ext_rs_l_if_p_cust_qos_pol") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_l_if_p_cust_qos_pol")
		err = aciClient.CreateRelationl3extRsLIfPCustQosPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_l3ext_rs_l_if_p_cust_qos_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_l3ext_rs_arp_if_pol") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_arp_if_pol")
		err = aciClient.CreateRelationl3extRsArpIfPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_l3ext_rs_arp_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_l3ext_rs_nd_if_pol") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_nd_if_pol")
		err = aciClient.CreateRelationl3extRsNdIfPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_l3ext_rs_nd_if_pol")
		d.Partial(false)

	}

	d.SetId(l3extLIfP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciLogicalInterfaceProfileRead(d, m)

}

func resourceAciLogicalInterfaceProfileRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l3extLIfP, err := getRemoteLogicalInterfaceProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setLogicalInterfaceProfileAttributes(l3extLIfP, d)

	l3extRsLIfPToNetflowMonitorPolData, err := aciClient.ReadRelationl3extRsLIfPToNetflowMonitorPolFromLogicalInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsLIfPToNetflowMonitorPol %v", err)

	} else {
		d.Set("relation_l3ext_rs_l_if_p_to_netflow_monitor_pol", l3extRsLIfPToNetflowMonitorPolData)
	}

	l3extRsPathL3OutAttData, err := aciClient.ReadRelationl3extRsPathL3OutAttFromLogicalInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsPathL3OutAtt %v", err)

	} else {
		d.Set("relation_l3ext_rs_path_l3_out_att", l3extRsPathL3OutAttData)
	}

	l3extRsEgressQosDppPolData, err := aciClient.ReadRelationl3extRsEgressQosDppPolFromLogicalInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsEgressQosDppPol %v", err)

	} else {
		d.Set("relation_l3ext_rs_egress_qos_dpp_pol", l3extRsEgressQosDppPolData)
	}

	l3extRsIngressQosDppPolData, err := aciClient.ReadRelationl3extRsIngressQosDppPolFromLogicalInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsIngressQosDppPol %v", err)

	} else {
		d.Set("relation_l3ext_rs_ingress_qos_dpp_pol", l3extRsIngressQosDppPolData)
	}

	l3extRsLIfPCustQosPolData, err := aciClient.ReadRelationl3extRsLIfPCustQosPolFromLogicalInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsLIfPCustQosPol %v", err)

	} else {
		d.Set("relation_l3ext_rs_l_if_p_cust_qos_pol", l3extRsLIfPCustQosPolData)
	}

	l3extRsArpIfPolData, err := aciClient.ReadRelationl3extRsArpIfPolFromLogicalInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsArpIfPol %v", err)

	} else {
		d.Set("relation_l3ext_rs_arp_if_pol", l3extRsArpIfPolData)
	}

	l3extRsNdIfPolData, err := aciClient.ReadRelationl3extRsNdIfPolFromLogicalInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsNdIfPol %v", err)

	} else {
		d.Set("relation_l3ext_rs_nd_if_pol", l3extRsNdIfPolData)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciLogicalInterfaceProfileDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l3extLIfP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
