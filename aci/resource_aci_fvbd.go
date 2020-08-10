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

func resourceAciBridgeDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciBridgeDomainCreate,
		Update: resourceAciBridgeDomainUpdate,
		Read:   resourceAciBridgeDomainRead,
		Delete: resourceAciBridgeDomainDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciBridgeDomainImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"tenant_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"optimize_wan_bandwidth": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"arp_flood": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ep_clear": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ep_move_detect_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"host_based_routing": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"intersite_bum_traffic_allow": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"intersite_l2_stretch": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ip_learning": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ipv6_mcast_allow": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"limit_ip_learn_to_subnets": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ll_addr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mac": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mcast_allow": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"multi_dst_pkt_act": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"bridge_domain_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"unicast_route": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"unk_mac_ucast_act": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"unk_mcast_act": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"v6unk_mcast_act": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"vmac": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_fv_rs_bd_to_profile": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_mldsn": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_abd_pol_mon_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_bd_to_nd_p": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_bd_flood_to": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_fv_rs_bd_to_fhs": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_bd_to_relay_p": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_ctx": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_bd_to_netflow_monitor_pol": &schema.Schema{
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
			"relation_fv_rs_igmpsn": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_bd_to_ep_ret": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_bd_to_out": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
		}),
	}
}
func getRemoteBridgeDomain(client *client.Client, dn string) (*models.BridgeDomain, error) {
	fvBDCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fvBD := models.BridgeDomainFromContainer(fvBDCont)

	if fvBD.DistinguishedName == "" {
		return nil, fmt.Errorf("BridgeDomain %s not found", fvBD.DistinguishedName)
	}

	return fvBD, nil
}

func setBridgeDomainAttributes(fvBD *models.BridgeDomain, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(fvBD.DistinguishedName)
	d.Set("description", fvBD.Description)
	// d.Set("tenant_dn", GetParentDn(fvBD.DistinguishedName))
	if dn != fvBD.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	fvBDMap, _ := fvBD.ToMap()

	d.Set("name", fvBDMap["name"])

	d.Set("optimize_wan_bandwidth", fvBDMap["OptimizeWanBandwidth"])
	d.Set("annotation", fvBDMap["annotation"])
	d.Set("arp_flood", fvBDMap["arpFlood"])
	d.Set("ep_clear", fvBDMap["epClear"])
	d.Set("ep_move_detect_mode", fvBDMap["epMoveDetectMode"])
	d.Set("host_based_routing", fvBDMap["hostBasedRouting"])
	d.Set("intersite_bum_traffic_allow", fvBDMap["intersiteBumTrafficAllow"])
	d.Set("intersite_l2_stretch", fvBDMap["intersiteL2Stretch"])
	d.Set("ip_learning", fvBDMap["ipLearning"])
	d.Set("ipv6_mcast_allow", fvBDMap["ipv6McastAllow"])
	d.Set("limit_ip_learn_to_subnets", fvBDMap["limitIpLearnToSubnets"])
	d.Set("ll_addr", fvBDMap["llAddr"])
	d.Set("mac", fvBDMap["mac"])
	d.Set("mcast_allow", fvBDMap["mcastAllow"])
	d.Set("multi_dst_pkt_act", fvBDMap["multiDstPktAct"])
	d.Set("name_alias", fvBDMap["nameAlias"])
	d.Set("bridge_domain_type", fvBDMap["type"])
	d.Set("unicast_route", fvBDMap["unicastRoute"])
	d.Set("unk_mac_ucast_act", fvBDMap["unkMacUcastAct"])
	d.Set("unk_mcast_act", fvBDMap["unkMcastAct"])
	d.Set("v6unk_mcast_act", fvBDMap["v6unkMcastAct"])
	d.Set("vmac", fvBDMap["vmac"])
	return d
}

func resourceAciBridgeDomainImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fvBD, err := getRemoteBridgeDomain(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setBridgeDomainAttributes(fvBD, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciBridgeDomainCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] BridgeDomain: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	fvBDAttr := models.BridgeDomainAttributes{}
	if OptimizeWanBandwidth, ok := d.GetOk("optimize_wan_bandwidth"); ok {
		fvBDAttr.OptimizeWanBandwidth = OptimizeWanBandwidth.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvBDAttr.Annotation = Annotation.(string)
	} else {
		fvBDAttr.Annotation = "{}"
	}
	if ArpFlood, ok := d.GetOk("arp_flood"); ok {
		fvBDAttr.ArpFlood = ArpFlood.(string)
	}
	if EpClear, ok := d.GetOk("ep_clear"); ok {
		fvBDAttr.EpClear = EpClear.(string)
	}
	if EpMoveDetectMode, ok := d.GetOk("ep_move_detect_mode"); ok {
		fvBDAttr.EpMoveDetectMode = EpMoveDetectMode.(string)
	}
	if HostBasedRouting, ok := d.GetOk("host_based_routing"); ok {
		fvBDAttr.HostBasedRouting = HostBasedRouting.(string)
	}
	if IntersiteBumTrafficAllow, ok := d.GetOk("intersite_bum_traffic_allow"); ok {
		fvBDAttr.IntersiteBumTrafficAllow = IntersiteBumTrafficAllow.(string)
	}
	if IntersiteL2Stretch, ok := d.GetOk("intersite_l2_stretch"); ok {
		fvBDAttr.IntersiteL2Stretch = IntersiteL2Stretch.(string)
	}
	if IpLearning, ok := d.GetOk("ip_learning"); ok {
		fvBDAttr.IpLearning = IpLearning.(string)
	}
	if Ipv6McastAllow, ok := d.GetOk("ipv6_mcast_allow"); ok {
		fvBDAttr.Ipv6McastAllow = Ipv6McastAllow.(string)
	}
	if LimitIpLearnToSubnets, ok := d.GetOk("limit_ip_learn_to_subnets"); ok {
		fvBDAttr.LimitIpLearnToSubnets = LimitIpLearnToSubnets.(string)
	}
	if LlAddr, ok := d.GetOk("ll_addr"); ok {
		fvBDAttr.LlAddr = LlAddr.(string)
	}
	if Mac, ok := d.GetOk("mac"); ok {
		fvBDAttr.Mac = Mac.(string)
	}
	if McastAllow, ok := d.GetOk("mcast_allow"); ok {
		fvBDAttr.McastAllow = McastAllow.(string)
	}
	if MultiDstPktAct, ok := d.GetOk("multi_dst_pkt_act"); ok {
		fvBDAttr.MultiDstPktAct = MultiDstPktAct.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvBDAttr.NameAlias = NameAlias.(string)
	}
	if BridgeDomain_type, ok := d.GetOk("bridge_domain_type"); ok {
		fvBDAttr.BridgeDomain_type = BridgeDomain_type.(string)
	}
	if UnicastRoute, ok := d.GetOk("unicast_route"); ok {
		fvBDAttr.UnicastRoute = UnicastRoute.(string)
	}
	if UnkMacUcastAct, ok := d.GetOk("unk_mac_ucast_act"); ok {
		fvBDAttr.UnkMacUcastAct = UnkMacUcastAct.(string)
	}
	if UnkMcastAct, ok := d.GetOk("unk_mcast_act"); ok {
		fvBDAttr.UnkMcastAct = UnkMcastAct.(string)
	}
	if V6unkMcastAct, ok := d.GetOk("v6unk_mcast_act"); ok {
		fvBDAttr.V6unkMcastAct = V6unkMcastAct.(string)
	}
	if Vmac, ok := d.GetOk("vmac"); ok {
		fvBDAttr.Vmac = Vmac.(string)
	}
	fvBD := models.NewBridgeDomain(fmt.Sprintf("BD-%s", name), TenantDn, desc, fvBDAttr)

	err := aciClient.Save(fvBD)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if relationTofvRsBDToProfile, ok := d.GetOk("relation_fv_rs_bd_to_profile"); ok {
		relationParam := relationTofvRsBDToProfile.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsBDToProfileFromBridgeDomain(fvBD.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_bd_to_profile")
		d.Partial(false)

	}
	if relationTofvRsMldsn, ok := d.GetOk("relation_fv_rs_mldsn"); ok {
		relationParam := relationTofvRsMldsn.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsMldsnFromBridgeDomain(fvBD.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_mldsn")
		d.Partial(false)

	}
	if relationTofvRsABDPolMonPol, ok := d.GetOk("relation_fv_rs_abd_pol_mon_pol"); ok {
		relationParam := relationTofvRsABDPolMonPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsABDPolMonPolFromBridgeDomain(fvBD.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_abd_pol_mon_pol")
		d.Partial(false)

	}
	if relationTofvRsBDToNdP, ok := d.GetOk("relation_fv_rs_bd_to_nd_p"); ok {
		relationParam := relationTofvRsBDToNdP.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsBDToNdPFromBridgeDomain(fvBD.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_bd_to_nd_p")
		d.Partial(false)

	}
	if relationTofvRsBdFloodTo, ok := d.GetOk("relation_fv_rs_bd_flood_to"); ok {
		relationParamList := toStringList(relationTofvRsBdFloodTo.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsBdFloodToFromBridgeDomain(fvBD.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_bd_flood_to")
			d.Partial(false)
		}
	}
	if relationTofvRsBDToFhs, ok := d.GetOk("relation_fv_rs_bd_to_fhs"); ok {
		relationParam := relationTofvRsBDToFhs.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsBDToFhsFromBridgeDomain(fvBD.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_bd_to_fhs")
		d.Partial(false)

	}
	if relationTofvRsBDToRelayP, ok := d.GetOk("relation_fv_rs_bd_to_relay_p"); ok {
		relationParam := relationTofvRsBDToRelayP.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsBDToRelayPFromBridgeDomain(fvBD.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_bd_to_relay_p")
		d.Partial(false)

	}
	if relationTofvRsCtx, ok := d.GetOk("relation_fv_rs_ctx"); ok {
		relationParam := relationTofvRsCtx.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsCtxFromBridgeDomain(fvBD.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_ctx")
		d.Partial(false)

	}
	if relationTofvRsBDToNetflowMonitorPol, ok := d.GetOk("relation_fv_rs_bd_to_netflow_monitor_pol"); ok {

		relationParamList := relationTofvRsBDToNetflowMonitorPol.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationfvRsBDToNetflowMonitorPolFromBridgeDomain(fvBD.DistinguishedName, paramMap["tn_netflow_monitor_pol_name"].(string), paramMap["flt_type"].(string))
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_bd_to_netflow_monitor_pol")
			d.Partial(false)
		}

	}
	if relationTofvRsIgmpsn, ok := d.GetOk("relation_fv_rs_igmpsn"); ok {
		relationParam := relationTofvRsIgmpsn.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsIgmpsnFromBridgeDomain(fvBD.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_igmpsn")
		d.Partial(false)

	}
	if relationTofvRsBdToEpRet, ok := d.GetOk("relation_fv_rs_bd_to_ep_ret"); ok {
		relationParam := relationTofvRsBdToEpRet.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsBdToEpRetFromBridgeDomain(fvBD.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_bd_to_ep_ret")
		d.Partial(false)

	}
	if relationTofvRsBDToOut, ok := d.GetOk("relation_fv_rs_bd_to_out"); ok {
		relationParamList := toStringList(relationTofvRsBDToOut.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsBDToOutFromBridgeDomain(fvBD.DistinguishedName, relationParamName)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_bd_to_out")
			d.Partial(false)
		}
	}

	d.SetId(fvBD.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciBridgeDomainRead(d, m)
}

func resourceAciBridgeDomainUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] BridgeDomain: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	fvBDAttr := models.BridgeDomainAttributes{}
	if OptimizeWanBandwidth, ok := d.GetOk("optimize_wan_bandwidth"); ok {
		fvBDAttr.OptimizeWanBandwidth = OptimizeWanBandwidth.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvBDAttr.Annotation = Annotation.(string)
	} else {
		fvBDAttr.Annotation = "{}"
	}
	if ArpFlood, ok := d.GetOk("arp_flood"); ok {
		fvBDAttr.ArpFlood = ArpFlood.(string)
	}
	if EpClear, ok := d.GetOk("ep_clear"); ok {
		fvBDAttr.EpClear = EpClear.(string)
	}
	if EpMoveDetectMode, ok := d.GetOk("ep_move_detect_mode"); ok {
		fvBDAttr.EpMoveDetectMode = EpMoveDetectMode.(string)
	}
	if HostBasedRouting, ok := d.GetOk("host_based_routing"); ok {
		fvBDAttr.HostBasedRouting = HostBasedRouting.(string)
	}
	if IntersiteBumTrafficAllow, ok := d.GetOk("intersite_bum_traffic_allow"); ok {
		fvBDAttr.IntersiteBumTrafficAllow = IntersiteBumTrafficAllow.(string)
	}
	if IntersiteL2Stretch, ok := d.GetOk("intersite_l2_stretch"); ok {
		fvBDAttr.IntersiteL2Stretch = IntersiteL2Stretch.(string)
	}
	if IpLearning, ok := d.GetOk("ip_learning"); ok {
		fvBDAttr.IpLearning = IpLearning.(string)
	}
	if Ipv6McastAllow, ok := d.GetOk("ipv6_mcast_allow"); ok {
		fvBDAttr.Ipv6McastAllow = Ipv6McastAllow.(string)
	}
	if LimitIpLearnToSubnets, ok := d.GetOk("limit_ip_learn_to_subnets"); ok {
		fvBDAttr.LimitIpLearnToSubnets = LimitIpLearnToSubnets.(string)
	}
	if LlAddr, ok := d.GetOk("ll_addr"); ok {
		fvBDAttr.LlAddr = LlAddr.(string)
	}
	if Mac, ok := d.GetOk("mac"); ok {
		fvBDAttr.Mac = Mac.(string)
	}
	if McastAllow, ok := d.GetOk("mcast_allow"); ok {
		fvBDAttr.McastAllow = McastAllow.(string)
	}
	if MultiDstPktAct, ok := d.GetOk("multi_dst_pkt_act"); ok {
		fvBDAttr.MultiDstPktAct = MultiDstPktAct.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvBDAttr.NameAlias = NameAlias.(string)
	}
	if BridgeDomain_type, ok := d.GetOk("bridge_domain_type"); ok {
		fvBDAttr.BridgeDomain_type = BridgeDomain_type.(string)
	}
	if UnicastRoute, ok := d.GetOk("unicast_route"); ok {
		fvBDAttr.UnicastRoute = UnicastRoute.(string)
	}
	if UnkMacUcastAct, ok := d.GetOk("unk_mac_ucast_act"); ok {
		fvBDAttr.UnkMacUcastAct = UnkMacUcastAct.(string)
	}
	if UnkMcastAct, ok := d.GetOk("unk_mcast_act"); ok {
		fvBDAttr.UnkMcastAct = UnkMcastAct.(string)
	}
	if V6unkMcastAct, ok := d.GetOk("v6unk_mcast_act"); ok {
		fvBDAttr.V6unkMcastAct = V6unkMcastAct.(string)
	}
	if Vmac, ok := d.GetOk("vmac"); ok {
		fvBDAttr.Vmac = Vmac.(string)
	}
	fvBD := models.NewBridgeDomain(fmt.Sprintf("BD-%s", name), TenantDn, desc, fvBDAttr)

	fvBD.Status = "modified"

	err := aciClient.Save(fvBD)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if d.HasChange("relation_fv_rs_bd_to_profile") {
		_, newRelParam := d.GetChange("relation_fv_rs_bd_to_profile")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationfvRsBDToProfileFromBridgeDomain(fvBD.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationfvRsBDToProfileFromBridgeDomain(fvBD.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_bd_to_profile")
		d.Partial(false)

	}
	if d.HasChange("relation_fv_rs_mldsn") {
		_, newRelParam := d.GetChange("relation_fv_rs_mldsn")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationfvRsMldsnFromBridgeDomain(fvBD.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_mldsn")
		d.Partial(false)

	}
	if d.HasChange("relation_fv_rs_abd_pol_mon_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_abd_pol_mon_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationfvRsABDPolMonPolFromBridgeDomain(fvBD.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationfvRsABDPolMonPolFromBridgeDomain(fvBD.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_abd_pol_mon_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_fv_rs_bd_to_nd_p") {
		_, newRelParam := d.GetChange("relation_fv_rs_bd_to_nd_p")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationfvRsBDToNdPFromBridgeDomain(fvBD.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_bd_to_nd_p")
		d.Partial(false)

	}
	if d.HasChange("relation_fv_rs_bd_flood_to") {
		oldRel, newRel := d.GetChange("relation_fv_rs_bd_flood_to")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsBdFloodToFromBridgeDomain(fvBD.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsBdFloodToFromBridgeDomain(fvBD.DistinguishedName, relDn)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_bd_flood_to")
			d.Partial(false)

		}

	}
	if d.HasChange("relation_fv_rs_bd_to_fhs") {
		_, newRelParam := d.GetChange("relation_fv_rs_bd_to_fhs")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationfvRsBDToFhsFromBridgeDomain(fvBD.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationfvRsBDToFhsFromBridgeDomain(fvBD.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_bd_to_fhs")
		d.Partial(false)

	}
	if d.HasChange("relation_fv_rs_bd_to_relay_p") {
		_, newRelParam := d.GetChange("relation_fv_rs_bd_to_relay_p")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationfvRsBDToRelayPFromBridgeDomain(fvBD.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationfvRsBDToRelayPFromBridgeDomain(fvBD.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_bd_to_relay_p")
		d.Partial(false)

	}
	if d.HasChange("relation_fv_rs_ctx") {
		_, newRelParam := d.GetChange("relation_fv_rs_ctx")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationfvRsCtxFromBridgeDomain(fvBD.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_ctx")
		d.Partial(false)

	}
	if d.HasChange("relation_fv_rs_bd_to_netflow_monitor_pol") {
		oldRel, newRel := d.GetChange("relation_fv_rs_bd_to_netflow_monitor_pol")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.DeleteRelationfvRsBDToNetflowMonitorPolFromBridgeDomain(fvBD.DistinguishedName, paramMap["tn_netflow_monitor_pol_name"].(string), paramMap["flt_type"].(string))
			if err != nil {
				return err
			}

		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationfvRsBDToNetflowMonitorPolFromBridgeDomain(fvBD.DistinguishedName, paramMap["tn_netflow_monitor_pol_name"].(string), paramMap["flt_type"].(string))
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_bd_to_netflow_monitor_pol")
			d.Partial(false)
		}

	}
	if d.HasChange("relation_fv_rs_igmpsn") {
		_, newRelParam := d.GetChange("relation_fv_rs_igmpsn")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationfvRsIgmpsnFromBridgeDomain(fvBD.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_igmpsn")
		d.Partial(false)

	}
	if d.HasChange("relation_fv_rs_bd_to_ep_ret") {
		_, newRelParam := d.GetChange("relation_fv_rs_bd_to_ep_ret")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationfvRsBdToEpRetFromBridgeDomain(fvBD.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_bd_to_ep_ret")
		d.Partial(false)

	}
	if d.HasChange("relation_fv_rs_bd_to_out") {
		oldRel, newRel := d.GetChange("relation_fv_rs_bd_to_out")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			relDnName := GetMOName(relDn)
			err = aciClient.DeleteRelationfvRsBDToOutFromBridgeDomain(fvBD.DistinguishedName, relDnName)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationfvRsBDToOutFromBridgeDomain(fvBD.DistinguishedName, relDnName)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_bd_to_out")
			d.Partial(false)

		}

	}

	d.SetId(fvBD.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciBridgeDomainRead(d, m)

}

func resourceAciBridgeDomainRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fvBD, err := getRemoteBridgeDomain(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setBridgeDomainAttributes(fvBD, d)

	fvRsBDToProfileData, err := aciClient.ReadRelationfvRsBDToProfileFromBridgeDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsBDToProfile %v", err)

	} else {
		if _, ok := d.GetOk("relation_fv_rs_bd_to_profile"); ok {
			tfName := GetMOName(d.Get("relation_fv_rs_bd_to_profile").(string))
			if tfName != fvRsBDToProfileData {
				d.Set("relation_fv_rs_bd_to_profile", "")
			}
		}
	}

	fvRsMldsnData, err := aciClient.ReadRelationfvRsMldsnFromBridgeDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsMldsn %v", err)

	} else {
		if _, ok := d.GetOk("relation_fv_rs_mldsn"); ok {
			tfName := GetMOName(d.Get("relation_fv_rs_mldsn").(string))
			if tfName != fvRsMldsnData {
				d.Set("relation_fv_rs_mldsn", "")
			}
		}
	}

	fvRsABDPolMonPolData, err := aciClient.ReadRelationfvRsABDPolMonPolFromBridgeDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsABDPolMonPol %v", err)

	} else {
		if _, ok := d.GetOk("relation_fv_rs_abd_pol_mon_pol"); ok {
			tfName := GetMOName(d.Get("relation_fv_rs_abd_pol_mon_pol").(string))
			if tfName != fvRsABDPolMonPolData {
				d.Set("relation_fv_rs_abd_pol_mon_pol", "")
			}
		}
	}

	fvRsBDToNdPData, err := aciClient.ReadRelationfvRsBDToNdPFromBridgeDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsBDToNdP %v", err)

	} else {
		if _, ok := d.GetOk("relation_fv_rs_bd_to_nd_p"); ok {
			tfName := GetMOName(d.Get("relation_fv_rs_bd_to_nd_p").(string))
			if tfName != fvRsBDToNdPData {
				d.Set("relation_fv_rs_bd_to_nd_p", "")
			}
		}
	}

	fvRsBdFloodToData, err := aciClient.ReadRelationfvRsBdFloodToFromBridgeDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsBdFloodTo %v", err)

	} else {
		d.Set("relation_fv_rs_bd_flood_to", fvRsBdFloodToData)
	}

	fvRsBDToFhsData, err := aciClient.ReadRelationfvRsBDToFhsFromBridgeDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsBDToFhs %v", err)

	} else {
		if _, ok := d.GetOk("relation_fv_rs_bd_to_fhs"); ok {
			tfName := GetMOName(d.Get("relation_fv_rs_bd_to_fhs").(string))
			if tfName != fvRsBDToFhsData {
				d.Set("relation_fv_rs_bd_to_fhs", "")
			}
		}
	}

	fvRsBDToRelayPData, err := aciClient.ReadRelationfvRsBDToRelayPFromBridgeDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsBDToRelayP %v", err)

	} else {
		if _, ok := d.GetOk("relation_fv_rs_bd_to_relay_p"); ok {
			tfName := GetMOName(d.Get("relation_fv_rs_bd_to_relay_p").(string))
			if tfName != fvRsBDToRelayPData {
				d.Set("relation_fv_rs_bd_to_relay_p", "")
			}
		}
	}

	fvRsCtxData, err := aciClient.ReadRelationfvRsCtxFromBridgeDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsCtx %v", err)

	} else {
		if _, ok := d.GetOk("relation_fv_rs_ctx"); ok {
			tfName := GetMOName(d.Get("relation_fv_rs_ctx").(string))
			if tfName != fvRsCtxData {
				d.Set("relation_fv_rs_ctx", "")
			}
		}
	}

	fvRsBDToNetflowMonitorPolData, err := aciClient.ReadRelationfvRsBDToNetflowMonitorPolFromBridgeDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsBDToNetflowMonitorPol %v", err)

	} else {
		d.Set("relation_fv_rs_bd_to_netflow_monitor_pol", fvRsBDToNetflowMonitorPolData)
	}

	fvRsIgmpsnData, err := aciClient.ReadRelationfvRsIgmpsnFromBridgeDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsIgmpsn %v", err)

	} else {
		if _, ok := d.GetOk("relation_fv_rs_igmpsn"); ok {
			tfName := GetMOName(d.Get("relation_fv_rs_igmpsn").(string))
			if tfName != fvRsIgmpsnData {
				d.Set("relation_fv_rs_igmpsn", "")
			}
		}
	}

	fvRsBdToEpRetData, err := aciClient.ReadRelationfvRsBdToEpRetFromBridgeDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsBdToEpRet %v", err)

	} else {
		if _, ok := d.GetOk("relation_fv_rs_bd_to_ep_ret"); ok {
			tfName := GetMOName(d.Get("relation_fv_rs_bd_to_ep_ret").(string))
			if tfName != fvRsBdToEpRetData {
				d.Set("relation_fv_rs_bd_to_ep_ret", "")
			}
		}
	}

	fvRsBDToOutData, err := aciClient.ReadRelationfvRsBDToOutFromBridgeDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsBDToOut %v", err)

	} else {
		if _, ok := d.GetOk("relation_fv_rs_bd_to_out"); ok {
			relationParamList := toStringList(d.Get("relation_fv_rs_bd_to_out").(*schema.Set).List())
			tfList := make([]string, 0, 1)
			for _, relationParam := range relationParamList {
				relationParamName := GetMOName(relationParam)
				tfList = append(tfList, relationParamName)
			}
			fvRsBDToOutDataList := toStringList(fvRsBDToOutData.(*schema.Set).List())
			sort.Strings(tfList)
			sort.Strings(fvRsBDToOutDataList)

			if !reflect.DeepEqual(tfList, fvRsBDToOutDataList) {
				d.Set("relation_fv_rs_bd_to_out", make([]string, 0, 1))
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciBridgeDomainDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fvBD")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
