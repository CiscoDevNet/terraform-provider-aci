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

func resourceAciLogicalInterfaceProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciLogicalInterfaceProfileCreate,
		UpdateContext: resourceAciLogicalInterfaceProfileUpdate,
		ReadContext:   resourceAciLogicalInterfaceProfileRead,
		DeleteContext: resourceAciLogicalInterfaceProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciLogicalInterfaceProfileImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"logical_node_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

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

func setLogicalInterfaceProfileAttributes(l3extLIfP *models.LogicalInterfaceProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(l3extLIfP.DistinguishedName)
	d.Set("description", l3extLIfP.Description)

	if dn != l3extLIfP.DistinguishedName {
		d.Set("logical_node_profile_dn", "")
	}
	l3extLIfPMap, err := l3extLIfP.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("logical_node_profile_dn", GetParentDn(dn, fmt.Sprintf("/lifp-%s", l3extLIfPMap["name"])))
	d.Set("name", l3extLIfPMap["name"])

	d.Set("annotation", l3extLIfPMap["annotation"])
	d.Set("name_alias", l3extLIfPMap["nameAlias"])
	d.Set("prio", l3extLIfPMap["prio"])
	d.Set("tag", l3extLIfPMap["tag"])
	return d, nil
}

func resourceAciLogicalInterfaceProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	l3extLIfP, err := getRemoteLogicalInterfaceProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	l3extLIfPMap, err := l3extLIfP.ToMap()
	if err != nil {
		return nil, err
	}
	name := l3extLIfPMap["name"]
	pDN := GetParentDn(dn, fmt.Sprintf("/lifp-%s", name))
	d.Set("logical_node_profile_dn", pDN)
	schemaFilled, err := setLogicalInterfaceProfileAttributes(l3extLIfP, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLogicalInterfaceProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LogicalInterfaceProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	LogicalNodeProfileDn := d.Get("logical_node_profile_dn").(string)

	l3extLIfPAttr := models.LogicalInterfaceProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extLIfPAttr.Annotation = Annotation.(string)
	} else {
		l3extLIfPAttr.Annotation = "{}"
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
		return diag.FromErr(err)

	}

	checkDns := make([]string, 0, 1)

	if relationTol3extRsEgressQosDppPol, ok := d.GetOk("relation_l3ext_rs_egress_qos_dpp_pol"); ok {
		relationParam := relationTol3extRsEgressQosDppPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTol3extRsIngressQosDppPol, ok := d.GetOk("relation_l3ext_rs_ingress_qos_dpp_pol"); ok {
		relationParam := relationTol3extRsIngressQosDppPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTol3extRsLIfPCustQosPol, ok := d.GetOk("relation_l3ext_rs_l_if_p_cust_qos_pol"); ok {
		relationParam := relationTol3extRsLIfPCustQosPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTol3extRsArpIfPol, ok := d.GetOk("relation_l3ext_rs_arp_if_pol"); ok {
		relationParam := relationTol3extRsArpIfPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTol3extRsNdIfPol, ok := d.GetOk("relation_l3ext_rs_nd_if_pol"); ok {
		relationParam := relationTol3extRsNdIfPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTol3extRsLIfPToNetflowMonitorPol, ok := d.GetOk("relation_l3ext_rs_l_if_p_to_netflow_monitor_pol"); ok {

		relationParamList := relationTol3extRsLIfPToNetflowMonitorPol.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationl3extRsLIfPToNetflowMonitorPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, paramMap["tn_netflow_monitor_pol_name"].(string), paramMap["flt_type"].(string))
			if err != nil {
				return diag.FromErr(err)

			}
		}

	}
	if relationTol3extRsEgressQosDppPol, ok := d.GetOk("relation_l3ext_rs_egress_qos_dpp_pol"); ok {
		relationParam := relationTol3extRsEgressQosDppPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationl3extRsEgressQosDppPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)

		}

	}
	if relationTol3extRsIngressQosDppPol, ok := d.GetOk("relation_l3ext_rs_ingress_qos_dpp_pol"); ok {
		relationParam := relationTol3extRsIngressQosDppPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationl3extRsIngressQosDppPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)

		}

	}
	if relationTol3extRsLIfPCustQosPol, ok := d.GetOk("relation_l3ext_rs_l_if_p_cust_qos_pol"); ok {
		relationParam := relationTol3extRsLIfPCustQosPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationl3extRsLIfPCustQosPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)

		}

	}
	if relationTol3extRsArpIfPol, ok := d.GetOk("relation_l3ext_rs_arp_if_pol"); ok {
		relationParam := relationTol3extRsArpIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationl3extRsArpIfPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)

		}

	}
	if relationTol3extRsNdIfPol, ok := d.GetOk("relation_l3ext_rs_nd_if_pol"); ok {
		relationParam := relationTol3extRsNdIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationl3extRsNdIfPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)

		}

	}

	d.SetId(l3extLIfP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciLogicalInterfaceProfileRead(ctx, d, m)
}

func resourceAciLogicalInterfaceProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LogicalInterfaceProfile: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	LogicalNodeProfileDn := d.Get("logical_node_profile_dn").(string)

	l3extLIfPAttr := models.LogicalInterfaceProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extLIfPAttr.Annotation = Annotation.(string)
	} else {
		l3extLIfPAttr.Annotation = "{}"
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
		return diag.FromErr(err)

	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_l3ext_rs_egress_qos_dpp_pol") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_egress_qos_dpp_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_l3ext_rs_ingress_qos_dpp_pol") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_ingress_qos_dpp_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_l3ext_rs_l_if_p_cust_qos_pol") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_l_if_p_cust_qos_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_l3ext_rs_arp_if_pol") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_arp_if_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_l3ext_rs_nd_if_pol") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_nd_if_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)

	}
	d.Partial(false)

	if d.HasChange("relation_l3ext_rs_l_if_p_to_netflow_monitor_pol") {
		oldRel, newRel := d.GetChange("relation_l3ext_rs_l_if_p_to_netflow_monitor_pol")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.DeleteRelationl3extRsLIfPToNetflowMonitorPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, paramMap["tn_netflow_monitor_pol_name"].(string), paramMap["flt_type"].(string))
			if err != nil {
				return diag.FromErr(err)

			}

		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationl3extRsLIfPToNetflowMonitorPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, paramMap["tn_netflow_monitor_pol_name"].(string), paramMap["flt_type"].(string))
			if err != nil {
				return diag.FromErr(err)

			}
		}

	}

	if d.HasChange("relation_l3ext_rs_egress_qos_dpp_pol") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_egress_qos_dpp_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationl3extRsEgressQosDppPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)

		}

	}
	if d.HasChange("relation_l3ext_rs_ingress_qos_dpp_pol") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_ingress_qos_dpp_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationl3extRsIngressQosDppPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)

		}
	}
	if d.HasChange("relation_l3ext_rs_l_if_p_cust_qos_pol") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_l_if_p_cust_qos_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationl3extRsLIfPCustQosPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)

		}
	}
	if d.HasChange("relation_l3ext_rs_arp_if_pol") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_arp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationl3extRsArpIfPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)

		}
	}
	if d.HasChange("relation_l3ext_rs_nd_if_pol") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_nd_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationl3extRsNdIfPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)

		}
	}

	d.SetId(l3extLIfP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciLogicalInterfaceProfileRead(ctx, d, m)

}

func resourceAciLogicalInterfaceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l3extLIfP, err := getRemoteLogicalInterfaceProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setLogicalInterfaceProfileAttributes(l3extLIfP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	l3extRsLIfPToNetflowMonitorPolData, err := aciClient.ReadRelationl3extRsLIfPToNetflowMonitorPolFromLogicalInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsLIfPToNetflowMonitorPol %v", err)

	} else {
		d.Set("relation_l3ext_rs_l_if_p_to_netflow_monitor_pol", l3extRsLIfPToNetflowMonitorPolData)
	}

	l3extRsEgressQosDppPolData, err := aciClient.ReadRelationl3extRsEgressQosDppPolFromLogicalInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsEgressQosDppPol %v", err)
		d.Set("relation_l3ext_rs_egress_qos_dpp_pol", "")

	} else {
		if _, ok := d.GetOk("relation_l3ext_rs_egress_qos_dpp_pol"); ok {
			tfName := GetMOName(d.Get("relation_l3ext_rs_egress_qos_dpp_pol").(string))
			if tfName != l3extRsEgressQosDppPolData {
				d.Set("relation_l3ext_rs_egress_qos_dpp_pol", "")
			}
		}
	}

	l3extRsIngressQosDppPolData, err := aciClient.ReadRelationl3extRsIngressQosDppPolFromLogicalInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsIngressQosDppPol %v", err)
		d.Set("relation_l3ext_rs_ingress_qos_dpp_pol", "")

	} else {
		if _, ok := d.GetOk("relation_l3ext_rs_ingress_qos_dpp_pol"); ok {
			tfName := GetMOName(d.Get("relation_l3ext_rs_ingress_qos_dpp_pol").(string))
			if tfName != l3extRsIngressQosDppPolData {
				d.Set("relation_l3ext_rs_ingress_qos_dpp_pol", "")
			}
		}
	}

	l3extRsLIfPCustQosPolData, err := aciClient.ReadRelationl3extRsLIfPCustQosPolFromLogicalInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsLIfPCustQosPol %v", err)
		d.Set("relation_l3ext_rs_l_if_p_cust_qos_pol", "")

	} else {
		if _, ok := d.GetOk("relation_l3ext_rs_l_if_p_cust_qos_pol"); ok {
			tfName := GetMOName(d.Get("relation_l3ext_rs_l_if_p_cust_qos_pol").(string))
			if tfName != l3extRsLIfPCustQosPolData {
				d.Set("relation_l3ext_rs_l_if_p_cust_qos_pol", "")
			}
		}
	}

	l3extRsArpIfPolData, err := aciClient.ReadRelationl3extRsArpIfPolFromLogicalInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsArpIfPol %v", err)
		d.Set("relation_l3ext_rs_arp_if_pol", "")

	} else {
		if _, ok := d.GetOk("relation_l3ext_rs_arp_if_pol"); ok {
			tfName := GetMOName(d.Get("relation_l3ext_rs_arp_if_pol").(string))
			if tfName != l3extRsArpIfPolData {
				d.Set("relation_l3ext_rs_arp_if_pol", "")
			}
		}
	}

	l3extRsNdIfPolData, err := aciClient.ReadRelationl3extRsNdIfPolFromLogicalInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsNdIfPol %v", err)
		d.Set("relation_l3ext_rs_nd_if_pol", "")

	} else {
		if _, ok := d.GetOk("relation_l3ext_rs_nd_if_pol"); ok {
			tfName := GetMOName(d.Get("relation_l3ext_rs_nd_if_pol").(string))
			if tfName != l3extRsNdIfPolData {
				d.Set("relation_l3ext_rs_nd_if_pol", "")
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciLogicalInterfaceProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l3extLIfP")
	if err != nil {
		return diag.FromErr(err)

	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)

}
