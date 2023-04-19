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

func dataSourceAciLogicalInterfaceProfile() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciLogicalInterfaceProfileRead,

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
				Computed: true,
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
		}),
	}
}

func dataSourceAciLogicalInterfaceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf(models.Rnl3extlifp, name)
	LogicalNodeProfileDn := d.Get("logical_node_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", LogicalNodeProfileDn, rn)

	l3extLIfP, err := getRemoteLogicalInterfaceProfile(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setLogicalInterfaceProfileAttributes(l3extLIfP, d)
	if err != nil {
		return diag.FromErr(err)
	}

	l3extRsLIfPToNetflowMonitorPolData, err := aciClient.ReadRelationl3extRsLIfPToNetflowMonitorPolFromLogicalInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsLIfPToNetflowMonitorPol %v", err)

	} else {
		relParamList := make([]map[string]string, 0, 1)
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
		d.Set("relation_l3ext_rs_l_if_p_to_netflow_monitor_pol", relParamList)
	}

	l3extRsEgressQosDppPolData, err := aciClient.ReadRelationl3extRsEgressQosDppPolFromLogicalInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsEgressQosDppPol %v", err)
		d.Set("relation_l3ext_rs_egress_qos_dpp_pol", "")

	} else {
		setRelationAttribute(d, "relation_l3ext_rs_egress_qos_dpp_pol", l3extRsEgressQosDppPolData.(string))
	}

	l3extRsIngressQosDppPolData, err := aciClient.ReadRelationl3extRsIngressQosDppPolFromLogicalInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsIngressQosDppPol %v", err)
		d.Set("relation_l3ext_rs_ingress_qos_dpp_pol", "")

	} else {
		setRelationAttribute(d, "relation_l3ext_rs_ingress_qos_dpp_pol", l3extRsIngressQosDppPolData.(string))
	}

	l3extRsLIfPCustQosPolData, err := aciClient.ReadRelationl3extRsLIfPCustQosPolFromLogicalInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsLIfPCustQosPol %v", err)
		d.Set("relation_l3ext_rs_l_if_p_cust_qos_pol", "")

	} else {
		setRelationAttribute(d, "relation_l3ext_rs_l_if_p_cust_qos_pol", l3extRsLIfPCustQosPolData.(string))
	}

	l3extRsArpIfPolData, err := aciClient.ReadRelationl3extRsArpIfPolFromLogicalInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsArpIfPol %v", err)
		d.Set("relation_l3ext_rs_arp_if_pol", "")

	} else {
		setRelationAttribute(d, "relation_l3ext_rs_arp_if_pol", l3extRsArpIfPolData.(string))
	}

	l3extRsNdIfPolData, err := aciClient.ReadRelationl3extRsNdIfPolFromLogicalInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsNdIfPol %v", err)
		d.Set("relation_l3ext_rs_nd_if_pol", "")

	} else {
		setRelationAttribute(d, "relation_l3ext_rs_nd_if_pol", l3extRsNdIfPolData.(string))
	}

	return nil
}
