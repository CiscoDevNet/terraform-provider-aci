package aci

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciBridgeDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciBridgeDomainCreate,
		UpdateContext: resourceAciBridgeDomainUpdate,
		ReadContext:   resourceAciBridgeDomainRead,
		DeleteContext: resourceAciBridgeDomainDelete,

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
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"arp_flood": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"ep_clear": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"ep_move_detect_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"garp",
					"disable",
				}, false),
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
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"intersite_l2_stretch": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"ip_learning": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
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
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
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
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"multi_dst_pkt_act": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"bd-flood",
					"encap-flood",
					"drop",
				}, false),
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
				ValidateFunc: validation.StringInSlice([]string{
					"regular",
					"fc",
				}, false),
			},

			"unicast_route": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"unk_mac_ucast_act": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"flood",
					"proxy",
				}, false),
			},

			"unk_mcast_act": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"flood",
					"opt-flood",
				}, false),
			},

			"v6unk_mcast_act": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"flood",
					"opt-flood",
				}, false),
			},

			"vmac": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"not-applicable",
				}, false),
			},

			"relation_fv_rs_bd_to_profile": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_mldsn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_fv_rs_abd_pol_mon_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_bd_to_nd_p": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
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
				Type:     schema.TypeString,
				Computed: true,
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
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_fv_rs_bd_to_ep_ret": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
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

func setBridgeDomainAttributes(fvBD *models.BridgeDomain, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(fvBD.DistinguishedName)
	d.Set("description", fvBD.Description)
	if dn != fvBD.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	fvBDMap, err := fvBD.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("tenant_dn", GetParentDn(dn, fmt.Sprintf("/BD-%s", fvBDMap["name"])))
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
	return d, err
}

func resourceAciBridgeDomainImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fvBD, err := getRemoteBridgeDomain(aciClient, dn)

	if err != nil {
		return nil, err
	}
	fvBDMap, err := fvBD.ToMap()
	if err != nil {
		return nil, err
	}
	name := fvBDMap["name"]
	pDN := GetParentDn(dn, fmt.Sprintf("/BD-%s", name))
	d.Set("tenant_dn", pDN)
	schemaFilled, err := setBridgeDomainAttributes(fvBD, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func checkForSubnetConflict(client *client.Client, bdDN, ctxRelation string) error {
	tokens := strings.Split(bdDN, "/")
	bdName := (strings.Split(tokens[2], "-"))[1]
	tenantDn := fmt.Sprintf("%s/%s", tokens[0], tokens[1])

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, tenantDn, "fvBD")

	domains, err := client.GetViaURL(dnUrl)
	if err != nil {
		return err
	}
	bdList := models.ListFromContainer(domains, "fvBD")

	dnUrl = fmt.Sprintf("%s/%s/%s.json", baseurlStr, bdDN, "fvSubnet")
	subnets, err := client.GetViaURL(dnUrl)
	if err != nil {
		return nil
	}
	subnetList := models.ListFromContainer(subnets, "fvSubnet")

	if len(bdList) > 1 {
		for i := 0; i < (len(bdList)); i++ {
			currName := models.G(bdList[i], "name")
			if currName != bdName {
				if len(subnetList) > 0 {
					for j := 0; j < len(subnetList); j++ {
						ip := models.G(subnetList[j], "ip")
						if checkForConflictingVRF(client, tenantDn, currName, ctxRelation, ip) {
							return fmt.Errorf("A subnet already exist with Bridge Domain %s and ip %s", currName, ip)
						}
					}
				}
			}
		}
	}
	return nil
}

func resourceAciBridgeDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		if EpMoveDetectMode == "disable" {
			fvBDAttr.EpMoveDetectMode = "{}"
		} else {
			fvBDAttr.EpMoveDetectMode = EpMoveDetectMode.(string)
		}
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
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTofvRsBDToProfile, ok := d.GetOk("relation_fv_rs_bd_to_profile"); ok {
		relationParam := relationTofvRsBDToProfile.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTofvRsMldsn, ok := d.GetOk("relation_fv_rs_mldsn"); ok {
		relationParam := relationTofvRsMldsn.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTofvRsABDPolMonPol, ok := d.GetOk("relation_fv_rs_abd_pol_mon_pol"); ok {
		relationParam := relationTofvRsABDPolMonPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTofvRsBDToNdP, ok := d.GetOk("relation_fv_rs_bd_to_nd_p"); ok {
		relationParam := relationTofvRsBDToNdP.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTofvRsBdFloodTo, ok := d.GetOk("relation_fv_rs_bd_flood_to"); ok {
		relationParamList := toStringList(relationTofvRsBdFloodTo.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTofvRsBDToFhs, ok := d.GetOk("relation_fv_rs_bd_to_fhs"); ok {
		relationParam := relationTofvRsBDToFhs.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTofvRsBDToRelayP, ok := d.GetOk("relation_fv_rs_bd_to_relay_p"); ok {
		relationParam := relationTofvRsBDToRelayP.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTofvRsCtx, ok := d.GetOk("relation_fv_rs_ctx"); ok {
		relationParam := relationTofvRsCtx.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTofvRsIgmpsn, ok := d.GetOk("relation_fv_rs_igmpsn"); ok {
		relationParam := relationTofvRsIgmpsn.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTofvRsBdToEpRet, ok := d.GetOk("relation_fv_rs_bd_to_ep_ret"); ok {
		relationParam := relationTofvRsBdToEpRet.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTofvRsBDToOut, ok := d.GetOk("relation_fv_rs_bd_to_out"); ok {
		relationParamList := toStringList(relationTofvRsBDToOut.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTofvRsBDToProfile, ok := d.GetOk("relation_fv_rs_bd_to_profile"); ok {
		relationParam := relationTofvRsBDToProfile.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsBDToProfileFromBridgeDomain(fvBD.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTofvRsMldsn, ok := d.GetOk("relation_fv_rs_mldsn"); ok {
		relationParam := relationTofvRsMldsn.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsMldsnFromBridgeDomain(fvBD.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTofvRsABDPolMonPol, ok := d.GetOk("relation_fv_rs_abd_pol_mon_pol"); ok {
		relationParam := relationTofvRsABDPolMonPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsABDPolMonPolFromBridgeDomain(fvBD.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTofvRsBDToNdP, ok := d.GetOk("relation_fv_rs_bd_to_nd_p"); ok {
		relationParam := relationTofvRsBDToNdP.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsBDToNdPFromBridgeDomain(fvBD.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTofvRsBdFloodTo, ok := d.GetOk("relation_fv_rs_bd_flood_to"); ok {
		relationParamList := toStringList(relationTofvRsBdFloodTo.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsBdFloodToFromBridgeDomain(fvBD.DistinguishedName, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if relationTofvRsBDToFhs, ok := d.GetOk("relation_fv_rs_bd_to_fhs"); ok {
		relationParam := relationTofvRsBDToFhs.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsBDToFhsFromBridgeDomain(fvBD.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTofvRsBDToRelayP, ok := d.GetOk("relation_fv_rs_bd_to_relay_p"); ok {
		relationParam := relationTofvRsBDToRelayP.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsBDToRelayPFromBridgeDomain(fvBD.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTofvRsCtx, ok := d.GetOk("relation_fv_rs_ctx"); ok {
		relationParam := relationTofvRsCtx.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsCtxFromBridgeDomain(fvBD.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTofvRsBDToNetflowMonitorPol, ok := d.GetOk("relation_fv_rs_bd_to_netflow_monitor_pol"); ok {

		relationParamList := relationTofvRsBDToNetflowMonitorPol.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationfvRsBDToNetflowMonitorPolFromBridgeDomain(fvBD.DistinguishedName, GetMOName(paramMap["tn_netflow_monitor_pol_name"].(string)), paramMap["flt_type"].(string))
			if err != nil {
				return diag.FromErr(err)
			}
		}

	}
	if relationTofvRsIgmpsn, ok := d.GetOk("relation_fv_rs_igmpsn"); ok {
		relationParam := relationTofvRsIgmpsn.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsIgmpsnFromBridgeDomain(fvBD.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTofvRsBdToEpRet, ok := d.GetOk("relation_fv_rs_bd_to_ep_ret"); ok {
		relationParam := relationTofvRsBdToEpRet.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsBdToEpRetFromBridgeDomain(fvBD.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTofvRsBDToOut, ok := d.GetOk("relation_fv_rs_bd_to_out"); ok {
		relationParamList := toStringList(relationTofvRsBDToOut.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsBDToOutFromBridgeDomain(fvBD.DistinguishedName, relationParamName)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(fvBD.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciBridgeDomainRead(ctx, d, m)
}

func resourceAciBridgeDomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		if EpMoveDetectMode == "disable" {
			fvBDAttr.EpMoveDetectMode = "{}"
		} else {
			fvBDAttr.EpMoveDetectMode = EpMoveDetectMode.(string)
		}
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
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_fv_rs_bd_to_profile") {
		_, newRelParam := d.GetChange("relation_fv_rs_bd_to_profile")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_fv_rs_mldsn") {
		_, newRelParam := d.GetChange("relation_fv_rs_mldsn")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_fv_rs_abd_pol_mon_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_abd_pol_mon_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_fv_rs_bd_to_nd_p") {
		_, newRelParam := d.GetChange("relation_fv_rs_bd_to_nd_p")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_fv_rs_bd_flood_to") {
		oldRel, newRel := d.GetChange("relation_fv_rs_bd_flood_to")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fv_rs_bd_to_fhs") {
		_, newRelParam := d.GetChange("relation_fv_rs_bd_to_fhs")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_fv_rs_bd_to_relay_p") {
		_, newRelParam := d.GetChange("relation_fv_rs_bd_to_relay_p")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_fv_rs_ctx") {
		_, newRelParam := d.GetChange("relation_fv_rs_ctx")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_fv_rs_igmpsn") {
		_, newRelParam := d.GetChange("relation_fv_rs_igmpsn")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_fv_rs_bd_to_ep_ret") {
		_, newRelParam := d.GetChange("relation_fv_rs_bd_to_ep_ret")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_fv_rs_bd_to_out") {
		oldRel, newRel := d.GetChange("relation_fv_rs_bd_to_out")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_fv_rs_bd_to_profile") {
		_, newRelParam := d.GetChange("relation_fv_rs_bd_to_profile")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationfvRsBDToProfileFromBridgeDomain(fvBD.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationfvRsBDToProfileFromBridgeDomain(fvBD.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_fv_rs_mldsn") {
		_, newRelParam := d.GetChange("relation_fv_rs_mldsn")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationfvRsMldsnFromBridgeDomain(fvBD.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_fv_rs_abd_pol_mon_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_abd_pol_mon_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationfvRsABDPolMonPolFromBridgeDomain(fvBD.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationfvRsABDPolMonPolFromBridgeDomain(fvBD.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_fv_rs_bd_to_nd_p") {
		_, newRelParam := d.GetChange("relation_fv_rs_bd_to_nd_p")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationfvRsBDToNdPFromBridgeDomain(fvBD.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

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
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsBdFloodToFromBridgeDomain(fvBD.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}
	if d.HasChange("relation_fv_rs_bd_to_fhs") {
		_, newRelParam := d.GetChange("relation_fv_rs_bd_to_fhs")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationfvRsBDToFhsFromBridgeDomain(fvBD.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationfvRsBDToFhsFromBridgeDomain(fvBD.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_fv_rs_bd_to_relay_p") {
		_, newRelParam := d.GetChange("relation_fv_rs_bd_to_relay_p")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationfvRsBDToRelayPFromBridgeDomain(fvBD.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationfvRsBDToRelayPFromBridgeDomain(fvBD.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_fv_rs_ctx") {
		_, newRelParam := d.GetChange("relation_fv_rs_ctx")
		err := checkForSubnetConflict(aciClient, d.Id(), newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationfvRsCtxFromBridgeDomain(fvBD.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_fv_rs_bd_to_netflow_monitor_pol") {
		oldRel, newRel := d.GetChange("relation_fv_rs_bd_to_netflow_monitor_pol")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.DeleteRelationfvRsBDToNetflowMonitorPolFromBridgeDomain(fvBD.DistinguishedName, GetMOName(paramMap["tn_netflow_monitor_pol_name"].(string)), paramMap["flt_type"].(string))
			if err != nil {
				return diag.FromErr(err)
			}

		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationfvRsBDToNetflowMonitorPolFromBridgeDomain(fvBD.DistinguishedName, GetMOName(paramMap["tn_netflow_monitor_pol_name"].(string)), paramMap["flt_type"].(string))
			if err != nil {
				return diag.FromErr(err)
			}
		}

	}
	if d.HasChange("relation_fv_rs_igmpsn") {
		_, newRelParam := d.GetChange("relation_fv_rs_igmpsn")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationfvRsIgmpsnFromBridgeDomain(fvBD.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_fv_rs_bd_to_ep_ret") {
		_, newRelParam := d.GetChange("relation_fv_rs_bd_to_ep_ret")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationfvRsBdToEpRetFromBridgeDomain(fvBD.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

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
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationfvRsBDToOutFromBridgeDomain(fvBD.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}

	d.SetId(fvBD.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciBridgeDomainRead(ctx, d, m)

}

func resourceAciBridgeDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fvBD, err := getRemoteBridgeDomain(aciClient, dn)

	if fvBD.EpMoveDetectMode == "" {
		fvBD.EpMoveDetectMode = "disable"
	}

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setBridgeDomainAttributes(fvBD, d)

	if err != nil {
		d.SetId("")
		return nil
	}

	fvRsBDToProfileData, err := aciClient.ReadRelationfvRsBDToProfileFromBridgeDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsBDToProfile %v", err)
		d.Set("relation_fv_rs_bd_to_profile", "")

	} else {
		d.Set("relation_fv_rs_bd_to_profile", fvRsBDToProfileData.(string))
	}

	fvRsMldsnData, err := aciClient.ReadRelationfvRsMldsnFromBridgeDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsMldsn %v", err)
		d.Set("relation_fv_rs_mldsn", "")

	} else {
		d.Set("relation_fv_rs_mldsn", fvRsMldsnData.(string))
	}

	fvRsABDPolMonPolData, err := aciClient.ReadRelationfvRsABDPolMonPolFromBridgeDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsABDPolMonPol %v", err)
		d.Set("relation_fv_rs_abd_pol_mon_pol", "")

	} else {
		d.Set("relation_fv_rs_abd_pol_mon_pol", fvRsABDPolMonPolData.(string))
	}

	fvRsBDToNdPData, err := aciClient.ReadRelationfvRsBDToNdPFromBridgeDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsBDToNdP %v", err)
		d.Set("relation_fv_rs_bd_to_nd_p", "")

	} else {
		d.Set("relation_fv_rs_bd_to_nd_p", fvRsBDToNdPData.(string))
	}

	fvRsBdFloodToData, err := aciClient.ReadRelationfvRsBdFloodToFromBridgeDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsBdFloodTo %v", err)
		d.Set("relation_fv_rs_bd_flood_to", make([]string, 0, 1))

	} else {
		d.Set("relation_fv_rs_bd_flood_to", fvRsBdFloodToData)
	}

	fvRsBDToFhsData, err := aciClient.ReadRelationfvRsBDToFhsFromBridgeDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsBDToFhs %v", err)
		d.Set("relation_fv_rs_bd_to_fhs", "")

	} else {
		d.Set("relation_fv_rs_bd_to_fhs", fvRsBDToFhsData.(string))
	}

	fvRsBDToRelayPData, err := aciClient.ReadRelationfvRsBDToRelayPFromBridgeDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsBDToRelayP %v", err)
		d.Set("relation_fv_rs_bd_to_relay_p", "")

	} else {
		d.Set("relation_fv_rs_bd_to_relay_p", fvRsBDToRelayPData.(string))
	}

	fvRsCtxData, err := aciClient.ReadRelationfvRsCtxFromBridgeDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsCtx %v", err)
		d.Set("relation_fv_rs_ctx", "")

	} else {
		d.Set("relation_fv_rs_ctx", fvRsCtxData.(string))
	}

	fvRsBDToNetflowMonitorPolData, err := aciClient.ReadRelationfvRsBDToNetflowMonitorPolFromBridgeDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsBDToNetflowMonitorPol %v", err)

	} else {
		listRelMap := make([]map[string]string, 0, 1)
		listfvRsBDToNetflowMonitorPolData := fvRsBDToNetflowMonitorPolData.([]map[string]string)
		for _, obj := range listfvRsBDToNetflowMonitorPolData {
			listRelMap = append(listRelMap, map[string]string{
				"tn_netflow_monitor_pol_name": obj["tnNetflowMonitorPolName"],
				"flt_type":                    obj["fltType"],
			})
		}
		d.Set("relation_fv_rs_bd_to_netflow_monitor_pol", listRelMap)
	}

	fvRsIgmpsnData, err := aciClient.ReadRelationfvRsIgmpsnFromBridgeDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsIgmpsn %v", err)
		d.Set("relation_fv_rs_igmpsn", "")

	} else {
		d.Set("relation_fv_rs_igmpsn", fvRsIgmpsnData.(string))
	}

	fvRsBdToEpRetData, err := aciClient.ReadRelationfvRsBdToEpRetFromBridgeDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsBdToEpRet %v", err)
		d.Set("relation_fv_rs_bd_to_ep_ret", "")

	} else {
		d.Set("relation_fv_rs_bd_to_ep_ret", fvRsBdToEpRetData.(string))
	}

	fvRsBDToOutData, err := aciClient.ReadRelationfvRsBDToOutFromBridgeDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsBDToOut %v", err)
		d.Set("relation_fv_rs_bd_to_out", make([]string, 0, 1))

	} else {
		d.Set("relation_fv_rs_bd_to_out", toStringList(fvRsBDToOutData.(*schema.Set).List()))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciBridgeDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fvBD")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
