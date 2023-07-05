package aci

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
							Type:       schema.TypeString,
							Optional:   true,
							Deprecated: "Use tn_netflow_monitor_pol_dn instead",
						},
						"tn_netflow_monitor_pol_dn": {
							Type:     schema.TypeString,
							Optional: true,
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

			"relation_l3ext_rs_egress_qos_dpp_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_l3ext_rs_ingress_qos_dpp_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_l3ext_rs_l_if_p_cust_qos_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_l3ext_rs_arp_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_l3ext_rs_nd_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_l3ext_rs_pim_ip_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_l3ext_rs_pim_ipv6_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_l3ext_rs_igmp_if_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		}),
		CustomizeDiff: func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
			// Plan time validation.
			if relation, ok := diff.GetOk("relation_l3ext_rs_l_if_p_to_netflow_monitor_pol"); ok {
				for _, relationParam := range relation.(*schema.Set).List() {
					paramMap := relationParam.(map[string]interface{})
					if paramMap["tn_netflow_monitor_pol_dn"] != "" && paramMap["tn_netflow_monitor_pol_name"] != "" {
						return errors.New("Both tn_netflow_monitor_pol_dn and tn_netflow_monitor_pol_name cannot be used together. Use tn_netflow_monitor_pol_dn instead because tn_netflow_monitor_pol_name will be deprecated")
					}
					if paramMap["tn_netflow_monitor_pol_dn"] == "" && paramMap["tn_netflow_monitor_pol_name"] == "" {
						return errors.New("tn_netflow_monitor_pol_dn is required")
					}
				}
			}
			return nil
		},
	}
}

func getRemoteLogicalInterfaceProfile(client *client.Client, dn string) (*models.LogicalInterfaceProfile, error) {
	l3extLIfPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extLIfP := models.LogicalInterfaceProfileFromContainer(l3extLIfPCont)

	if l3extLIfP.DistinguishedName == "" {
		return nil, fmt.Errorf("Logical Interface Profile %s not found", dn)
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

func getTnNetflowMonitorPolName(paramMap map[string]interface{}) string {
	if paramMap["tn_netflow_monitor_pol_dn"] != "" {
		return GetMOName(paramMap["tn_netflow_monitor_pol_dn"].(string))
	} else {
		return GetMOName(paramMap["tn_netflow_monitor_pol_name"].(string))
	}
}

func getandSetL3extLIfPRelationshipAttributes(aciClient *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	l3extRsLIfPToNetflowMonitorPolData, err := aciClient.ReadRelationl3extRsLIfPToNetflowMonitorPolFromLogicalInterfaceProfile(dn)
	relParamList := make([]map[string]string, 0, 1)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsLIfPToNetflowMonitorPol %v", err)
	} else {
		relParams := l3extRsLIfPToNetflowMonitorPolData.([]map[string]string)
		for _, obj := range relParams {
			relParamList = append(relParamList, map[string]string{
				"tn_netflow_monitor_pol_dn": obj["tDn"],
				"flt_type":                  obj["fltType"],
			})
		}
		if relationTol3extRsLIfPToNetflowMonitorPol, ok := d.GetOk("relation_l3ext_rs_l_if_p_to_netflow_monitor_pol"); ok {
			relationParamListUser := relationTol3extRsLIfPToNetflowMonitorPol.(*schema.Set).List()
			for _, relationParamUser := range relationParamListUser {
				paramMapUser := relationParamUser.(map[string]interface{})
				if paramMapUser["tn_netflow_monitor_pol_dn"] == "" {
					for _, sub_attributes_map_apic := range relParamList {
						if sub_attribute_apic_key, ok := sub_attributes_map_apic["tn_netflow_monitor_pol_dn"]; ok {
							sub_attributes_map_apic["tn_netflow_monitor_pol_name"] = sub_attribute_apic_key
							delete(sub_attributes_map_apic, "tn_netflow_monitor_pol_dn")
						}
					}
				}
			}
		}
	}
	d.Set("relation_l3ext_rs_l_if_p_to_netflow_monitor_pol", relParamList)

	l3extRsEgressQosDppPolData, err := aciClient.ReadRelationl3extRsEgressQosDppPolFromLogicalInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsEgressQosDppPol %v", err)
		d.Set("relation_l3ext_rs_egress_qos_dpp_pol", "")

	} else {
		d.Set("relation_l3ext_rs_egress_qos_dpp_pol", l3extRsEgressQosDppPolData.(string))
	}

	l3extRsIngressQosDppPolData, err := aciClient.ReadRelationl3extRsIngressQosDppPolFromLogicalInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsIngressQosDppPol %v", err)
		d.Set("relation_l3ext_rs_ingress_qos_dpp_pol", "")

	} else {
		d.Set("relation_l3ext_rs_ingress_qos_dpp_pol", l3extRsIngressQosDppPolData.(string))
	}

	l3extRsLIfPCustQosPolData, err := aciClient.ReadRelationl3extRsLIfPCustQosPolFromLogicalInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsLIfPCustQosPol %v", err)
		d.Set("relation_l3ext_rs_l_if_p_cust_qos_pol", "")

	} else {
		d.Set("relation_l3ext_rs_l_if_p_cust_qos_pol", l3extRsLIfPCustQosPolData.(string))
	}

	l3extRsArpIfPolData, err := aciClient.ReadRelationl3extRsArpIfPolFromLogicalInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsArpIfPol %v", err)
		d.Set("relation_l3ext_rs_arp_if_pol", "")

	} else {
		d.Set("relation_l3ext_rs_arp_if_pol", l3extRsArpIfPolData.(string))
	}

	l3extRsNdIfPolData, err := aciClient.ReadRelationl3extRsNdIfPolFromLogicalInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsNdIfPol %v", err)
		d.Set("relation_l3ext_rs_nd_if_pol", "")

	} else {
		d.Set("relation_l3ext_rs_nd_if_pol", l3extRsNdIfPolData.(string))
	}

	pimRsIfPolData, err := aciClient.ReadRelationPIMRsIfPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation pimRsIfPol %v", err)
		d.Set("relation_l3ext_rs_pim_ip_if_pol", "")
	} else {
		d.Set("relation_l3ext_rs_pim_ip_if_pol", pimRsIfPolData.(string))
	}

	pimRsV6IfPolData, err := aciClient.ReadRelationPIMIPv6RsIfPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation pimRsV6IfPol %v", err)
		d.Set("relation_l3ext_rs_pim_ipv6_if_pol", "")
	} else {
		d.Set("relation_l3ext_rs_pim_ipv6_if_pol", pimRsV6IfPolData.(string))
	}

	igmpRsIfPolData, err := aciClient.ReadRelationIGMPRsIfPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation igmpRsIfPol %v", err)
		d.Set("relation_l3ext_rs_igmp_if_pol", "")
	} else {
		d.Set("relation_l3ext_rs_igmp_if_pol", igmpRsIfPolData.(string))
	}
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

	_, err = getandSetL3extLIfPRelationshipAttributes(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] L3extLIfP Relationship Attributes - Read finished successfully")
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

	if relationTopimRsIfPol, ok := d.GetOk("relation_l3ext_rs_pim_ip_if_pol"); ok {
		relationParam := relationTopimRsIfPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationTopimRsV6IfPol, ok := d.GetOk("relation_l3ext_rs_pim_ipv6_if_pol"); ok {
		relationParam := relationTopimRsV6IfPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationToigmpRsIfPol, ok := d.GetOk("relation_l3ext_rs_igmp_if_pol"); ok {
		relationParam := relationToigmpRsIfPol.(string)
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
			relationParamName := getTnNetflowMonitorPolName(paramMap)
			err = aciClient.CreateRelationl3extRsLIfPToNetflowMonitorPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, relationParamName, paramMap["flt_type"].(string))
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

	if relationTopimRsIfPol, ok := d.GetOk("relation_l3ext_rs_pim_ip_if_pol"); ok {

		pimIfP := models.NewPIMInterfaceProfile(l3extLIfP.DistinguishedName, "", models.PIMInterfaceProfileAttributes{})
		err := aciClient.Save(pimIfP)
		if err != nil {
			return diag.FromErr(err)
		}
		relationParam := relationTopimRsIfPol.(string)
		err = aciClient.CreateRelationPIMRsIfPolFromLogicalInterfaceProfile(pimIfP.DistinguishedName, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationTopimRsV6IfPol, ok := d.GetOk("relation_l3ext_rs_pim_ipv6_if_pol"); ok {
		pimIPV6IfP := models.NewPIMIPv6InterfaceProfile(l3extLIfP.DistinguishedName, "", models.PIMIPv6InterfaceProfileAttributes{})
		err := aciClient.Save(pimIPV6IfP)
		if err != nil {
			return diag.FromErr(err)
		}
		relationParam := relationTopimRsV6IfPol.(string)
		err = aciClient.CreateRelationPIMIPv6RsIfPolFromLogicalInterfaceProfile(pimIPV6IfP.DistinguishedName, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationToigmpRsIfPol, ok := d.GetOk("relation_l3ext_rs_igmp_if_pol"); ok {
		igmpIfP := models.NewIGMPInterfaceProfile(l3extLIfP.DistinguishedName, "", models.IGMPInterfaceProfileAttributes{})
		err := aciClient.Save(igmpIfP)
		if err != nil {
			return diag.FromErr(err)
		}
		relationParam := relationToigmpRsIfPol.(string)
		err = aciClient.CreateRelationIGMPRsIfPolFromLogicalInterfaceProfile(igmpIfP.DistinguishedName, relationParam)

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

	d.SetId(l3extLIfP.DistinguishedName)

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
			relationParamName := getTnNetflowMonitorPolName(paramMap)
			err = aciClient.DeleteRelationl3extRsLIfPToNetflowMonitorPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, relationParamName, paramMap["flt_type"].(string))
			if err != nil {
				return diag.FromErr(err)

			}

		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			relationParamName := getTnNetflowMonitorPolName(paramMap)
			err = aciClient.CreateRelationl3extRsLIfPToNetflowMonitorPolFromLogicalInterfaceProfile(l3extLIfP.DistinguishedName, relationParamName, paramMap["flt_type"].(string))
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

	if d.HasChange("relation_l3ext_rs_pim_ip_if_pol") {
		pimIfP := models.NewPIMInterfaceProfile(l3extLIfP.DistinguishedName, "", models.PIMInterfaceProfileAttributes{})
		pimIfP.Status = "modified"
		err_ifp := aciClient.Save(pimIfP)
		if err_ifp != nil {
			return diag.FromErr(err)
		}
		_, newRelParam := d.GetChange("relation_l3ext_rs_pim_ip_if_pol")
		err = aciClient.CreateRelationPIMRsIfPolFromLogicalInterfaceProfile(pimIfP.DistinguishedName, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("relation_l3ext_rs_pim_ipv6_if_pol") {
		pimIPV6IfP := models.NewPIMIPv6InterfaceProfile(l3extLIfP.DistinguishedName, "", models.PIMIPv6InterfaceProfileAttributes{})
		pimIPV6IfP.Status = "modified"
		err_ifp := aciClient.Save(pimIPV6IfP)
		if err_ifp != nil {
			return diag.FromErr(err)
		}
		_, newRelParam := d.GetChange("relation_l3ext_rs_pim_ipv6_if_pol")
		err = aciClient.CreateRelationPIMIPv6RsIfPolFromLogicalInterfaceProfile(pimIPV6IfP.DistinguishedName, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("relation_l3ext_rs_igmp_if_pol") {
		igmpIfP := models.NewIGMPInterfaceProfile(l3extLIfP.DistinguishedName, "", models.IGMPInterfaceProfileAttributes{})
		igmpIfP.Status = "modified"
		err_ifp := aciClient.Save(igmpIfP)
		if err_ifp != nil {
			return diag.FromErr(err)
		}
		_, newRelParam := d.GetChange("relation_l3ext_rs_igmp_if_pol")
		err = aciClient.CreateRelationIGMPRsIfPolFromLogicalInterfaceProfile(igmpIfP.DistinguishedName, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciLogicalInterfaceProfileRead(ctx, d, m)

}

func resourceAciLogicalInterfaceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l3extLIfP, err := getRemoteLogicalInterfaceProfile(aciClient, dn)

	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	_, err = setLogicalInterfaceProfileAttributes(l3extLIfP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	_, err = getandSetL3extLIfPRelationshipAttributes(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] L3extLIfP Relationship Attributes - Read finished successfully")
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
