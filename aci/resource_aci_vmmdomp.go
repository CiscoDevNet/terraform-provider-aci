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

func resourceAciVMMDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciVMMDomainCreate,
		UpdateContext: resourceAciVMMDomainUpdate,
		ReadContext:   resourceAciVMMDomainRead,
		DeleteContext: resourceAciVMMDomainDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciVMMDomainImport,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"provider_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "orchestrator:terraform",
			},

			"access_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"read-write",
					"read-only",
				}, false),
			},

			"arp_learning": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"enabled",
					"disabled",
				}, false),
			},

			"ave_time_out": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"config_infra_pg": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"ctrl_knob": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"none",
					"epDpVerify",
				}, false),
			},

			"delimiter": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"enable_ave": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"enable_tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"encap_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"unknown",
					"vlan",
					"vxlan",
				}, false),
			},

			"enf_pref": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"sw",
					"hw",
					"unknown",
				}, false),
			},

			"ep_inventory_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"none",
					"on-link",
				}, false),
			},

			"ep_ret_time": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"hv_avail_monitor": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"mcast_addr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"default",
					"n1kv",
					"unknown",
					"ovs",
					"k8s",
					"rhev",
					"cf",
					"openshift",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pref_encap_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"unspecified",
					"vlan",
					"vxlan",
				}, false),
			},

			"relation_vmm_rs_pref_enhanced_lag_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_vlan_ns": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_vmm_rs_dom_mcast_addr_ns": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_vmm_rs_default_cdp_if_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_vmm_rs_default_lacp_lag_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_vlan_ns_def": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_vip_addr_ns": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_vmm_rs_default_lldp_if_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_vmm_rs_default_stp_if_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_dom_vxlan_ns_def": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_vmm_rs_default_fw_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_vmm_rs_default_l2_inst_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		},
	}
}
func getRemoteVMMDomain(client *client.Client, dn string) (*models.VMMDomain, error) {
	vmmDomPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vmmDomP := models.VMMDomainFromContainer(vmmDomPCont)

	if vmmDomP.DistinguishedName == "" {
		return nil, fmt.Errorf("VMMDomain %s not found", vmmDomP.DistinguishedName)
	}

	return vmmDomP, nil
}

func setVMMDomainAttributes(vmmDomP *models.VMMDomain, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(vmmDomP.DistinguishedName)

	if dn != vmmDomP.DistinguishedName {
		d.Set("provider_profile_dn", "")
	}
	vmmDomPMap, err := vmmDomP.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("provider_profile_dn", GetParentDn(dn, fmt.Sprintf("/dom-%s", vmmDomPMap["name"])))
	d.Set("name", vmmDomPMap["name"])

	d.Set("access_mode", vmmDomPMap["accessMode"])
	d.Set("annotation", vmmDomPMap["annotation"])
	if vmmDomPMap["arpLearning"] == "" {
		d.Set("arp_learning", "disabled")
	} else {
		d.Set("arp_learning", vmmDomPMap["arpLearning"])
	}
	if vmmDomPMap["ctrlKnob"] == "" {
		d.Set("ctrl_knob", "none")
	} else {
		d.Set("ctrl_knob", vmmDomPMap["ctrlKnob"])
	}
	d.Set("ave_time_out", vmmDomPMap["aveTimeOut"])
	d.Set("config_infra_pg", vmmDomPMap["configInfraPg"])
	d.Set("delimiter", vmmDomPMap["delimiter"])
	d.Set("enable_ave", vmmDomPMap["enableAVE"])
	d.Set("enable_tag", vmmDomPMap["enableTag"])
	d.Set("encap_mode", vmmDomPMap["encapMode"])
	d.Set("enf_pref", vmmDomPMap["enfPref"])
	d.Set("ep_inventory_type", vmmDomPMap["epInventoryType"])
	d.Set("ep_ret_time", vmmDomPMap["epRetTime"])
	d.Set("hv_avail_monitor", vmmDomPMap["hvAvailMonitor"])
	d.Set("mcast_addr", vmmDomPMap["mcastAddr"])
	d.Set("mode", vmmDomPMap["mode"])
	d.Set("name_alias", vmmDomPMap["nameAlias"])
	d.Set("pref_encap_mode", vmmDomPMap["prefEncapMode"])
	return d, nil
}

func resourceAciVMMDomainImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	vmmDomP, err := getRemoteVMMDomain(aciClient, dn)

	if err != nil {
		return nil, err
	}
	vmmDomPMap, err := vmmDomP.ToMap()
	if err != nil {
		return nil, err
	}
	name := vmmDomPMap["name"]
	pDN := GetParentDn(dn, fmt.Sprintf("/dom-%s", name))
	d.Set("provider_profile_dn", pDN)
	schemaFilled, err := setVMMDomainAttributes(vmmDomP, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciVMMDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] VMMDomain: Beginning Creation")
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	ProviderProfileDn := d.Get("provider_profile_dn").(string)

	vmmDomPAttr := models.VMMDomainAttributes{}
	if AccessMode, ok := d.GetOk("access_mode"); ok {
		vmmDomPAttr.AccessMode = AccessMode.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vmmDomPAttr.Annotation = Annotation.(string)
	} else {
		vmmDomPAttr.Annotation = "{}"
	}
	if ArpLearning, ok := d.GetOk("arp_learning"); ok {
		vmmDomPAttr.ArpLearning = ArpLearning.(string)
	}
	if AveTimeOut, ok := d.GetOk("ave_time_out"); ok {
		vmmDomPAttr.AveTimeOut = AveTimeOut.(string)
	}
	if ConfigInfraPg, ok := d.GetOk("config_infra_pg"); ok {
		vmmDomPAttr.ConfigInfraPg = ConfigInfraPg.(string)
	}
	if CtrlKnob, ok := d.GetOk("ctrl_knob"); ok {
		vmmDomPAttr.CtrlKnob = CtrlKnob.(string)
	}
	if Delimiter, ok := d.GetOk("delimiter"); ok {
		vmmDomPAttr.Delimiter = Delimiter.(string)
	}
	if EnableAVE, ok := d.GetOk("enable_ave"); ok {
		vmmDomPAttr.EnableAVE = EnableAVE.(string)
	}
	if EnableTag, ok := d.GetOk("enable_tag"); ok {
		vmmDomPAttr.EnableTag = EnableTag.(string)
	}
	if EncapMode, ok := d.GetOk("encap_mode"); ok {
		vmmDomPAttr.EncapMode = EncapMode.(string)
	}
	if EnfPref, ok := d.GetOk("enf_pref"); ok {
		vmmDomPAttr.EnfPref = EnfPref.(string)
	}
	if EpInventoryType, ok := d.GetOk("ep_inventory_type"); ok {
		vmmDomPAttr.EpInventoryType = EpInventoryType.(string)
	}
	if EpRetTime, ok := d.GetOk("ep_ret_time"); ok {
		vmmDomPAttr.EpRetTime = EpRetTime.(string)
	}
	if HvAvailMonitor, ok := d.GetOk("hv_avail_monitor"); ok {
		vmmDomPAttr.HvAvailMonitor = HvAvailMonitor.(string)
	}
	if McastAddr, ok := d.GetOk("mcast_addr"); ok {
		vmmDomPAttr.McastAddr = McastAddr.(string)
	}
	if Mode, ok := d.GetOk("mode"); ok {
		vmmDomPAttr.Mode = Mode.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vmmDomPAttr.NameAlias = NameAlias.(string)
	}
	if PrefEncapMode, ok := d.GetOk("pref_encap_mode"); ok {
		vmmDomPAttr.PrefEncapMode = PrefEncapMode.(string)
	}
	vmmDomP := models.NewVMMDomain(fmt.Sprintf("dom-%s", name), ProviderProfileDn, vmmDomPAttr)

	err := aciClient.Save(vmmDomP)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTovmmRsPrefEnhancedLagPol, ok := d.GetOk("relation_vmm_rs_pref_enhanced_lag_pol"); ok {
		relationParam := relationTovmmRsPrefEnhancedLagPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsVlanNs, ok := d.GetOk("relation_infra_rs_vlan_ns"); ok {
		relationParam := relationToinfraRsVlanNs.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTovmmRsDomMcastAddrNs, ok := d.GetOk("relation_vmm_rs_dom_mcast_addr_ns"); ok {
		relationParam := relationTovmmRsDomMcastAddrNs.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTovmmRsDefaultCdpIfPol, ok := d.GetOk("relation_vmm_rs_default_cdp_if_pol"); ok {
		relationParam := relationTovmmRsDefaultCdpIfPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTovmmRsDefaultLacpLagPol, ok := d.GetOk("relation_vmm_rs_default_lacp_lag_pol"); ok {
		relationParam := relationTovmmRsDefaultLacpLagPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsVlanNsDef, ok := d.GetOk("relation_infra_rs_vlan_ns_def"); ok {
		relationParam := relationToinfraRsVlanNsDef.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsVipAddrNs, ok := d.GetOk("relation_infra_rs_vip_addr_ns"); ok {
		relationParam := relationToinfraRsVipAddrNs.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTovmmRsDefaultLldpIfPol, ok := d.GetOk("relation_vmm_rs_default_lldp_if_pol"); ok {
		relationParam := relationTovmmRsDefaultLldpIfPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTovmmRsDefaultStpIfPol, ok := d.GetOk("relation_vmm_rs_default_stp_if_pol"); ok {
		relationParam := relationTovmmRsDefaultStpIfPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsDomVxlanNsDef, ok := d.GetOk("relation_infra_rs_dom_vxlan_ns_def"); ok {
		relationParam := relationToinfraRsDomVxlanNsDef.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTovmmRsDefaultFwPol, ok := d.GetOk("relation_vmm_rs_default_fw_pol"); ok {
		relationParam := relationTovmmRsDefaultFwPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTovmmRsDefaultL2InstPol, ok := d.GetOk("relation_vmm_rs_default_l2_inst_pol"); ok {
		relationParam := relationTovmmRsDefaultL2InstPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTovmmRsPrefEnhancedLagPol, ok := d.GetOk("relation_vmm_rs_pref_enhanced_lag_pol"); ok {
		relationParam := relationTovmmRsPrefEnhancedLagPol.(string)
		err = aciClient.CreateRelationvmmRsPrefEnhancedLagPolFromVMMDomain(vmmDomP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationToinfraRsVlanNs, ok := d.GetOk("relation_infra_rs_vlan_ns"); ok {
		relationParam := relationToinfraRsVlanNs.(string)
		err = aciClient.CreateRelationinfraRsVlanNsFromVMMDomain(vmmDomP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationTovmmRsDomMcastAddrNs, ok := d.GetOk("relation_vmm_rs_dom_mcast_addr_ns"); ok {
		relationParam := relationTovmmRsDomMcastAddrNs.(string)
		err = aciClient.CreateRelationvmmRsDomMcastAddrNsFromVMMDomain(vmmDomP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationTovmmRsDefaultCdpIfPol, ok := d.GetOk("relation_vmm_rs_default_cdp_if_pol"); ok {
		relationParam := relationTovmmRsDefaultCdpIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationvmmRsDefaultCdpIfPolFromVMMDomain(vmmDomP.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationTovmmRsDefaultLacpLagPol, ok := d.GetOk("relation_vmm_rs_default_lacp_lag_pol"); ok {
		relationParam := relationTovmmRsDefaultLacpLagPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationvmmRsDefaultLacpLagPolFromVMMDomain(vmmDomP.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationToinfraRsVlanNsDef, ok := d.GetOk("relation_infra_rs_vlan_ns_def"); ok {
		relationParam := relationToinfraRsVlanNsDef.(string)
		err = aciClient.CreateRelationinfraRsVlanNsDefFromVMMDomain(vmmDomP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationToinfraRsVipAddrNs, ok := d.GetOk("relation_infra_rs_vip_addr_ns"); ok {
		relationParam := relationToinfraRsVipAddrNs.(string)
		err = aciClient.CreateRelationinfraRsVipAddrNsFromVMMDomain(vmmDomP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationTovmmRsDefaultLldpIfPol, ok := d.GetOk("relation_vmm_rs_default_lldp_if_pol"); ok {
		relationParam := relationTovmmRsDefaultLldpIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationvmmRsDefaultLldpIfPolFromVMMDomain(vmmDomP.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationTovmmRsDefaultStpIfPol, ok := d.GetOk("relation_vmm_rs_default_stp_if_pol"); ok {
		relationParam := relationTovmmRsDefaultStpIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationvmmRsDefaultStpIfPolFromVMMDomain(vmmDomP.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationToinfraRsDomVxlanNsDef, ok := d.GetOk("relation_infra_rs_dom_vxlan_ns_def"); ok {
		relationParam := relationToinfraRsDomVxlanNsDef.(string)
		err = aciClient.CreateRelationinfraRsDomVxlanNsDefFromVMMDomain(vmmDomP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationTovmmRsDefaultFwPol, ok := d.GetOk("relation_vmm_rs_default_fw_pol"); ok {
		relationParam := relationTovmmRsDefaultFwPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationvmmRsDefaultFwPolFromVMMDomain(vmmDomP.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationTovmmRsDefaultL2InstPol, ok := d.GetOk("relation_vmm_rs_default_l2_inst_pol"); ok {
		relationParam := relationTovmmRsDefaultL2InstPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationvmmRsDefaultL2InstPolFromVMMDomain(vmmDomP.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(vmmDomP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciVMMDomainRead(ctx, d, m)
}

func resourceAciVMMDomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] VMMDomain: Beginning Update")

	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	ProviderProfileDn := d.Get("provider_profile_dn").(string)

	vmmDomPAttr := models.VMMDomainAttributes{}
	if AccessMode, ok := d.GetOk("access_mode"); ok {
		vmmDomPAttr.AccessMode = AccessMode.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vmmDomPAttr.Annotation = Annotation.(string)
	} else {
		vmmDomPAttr.Annotation = "{}"
	}
	if ArpLearning, ok := d.GetOk("arp_learning"); ok {
		vmmDomPAttr.ArpLearning = ArpLearning.(string)
	}
	if AveTimeOut, ok := d.GetOk("ave_time_out"); ok {
		vmmDomPAttr.AveTimeOut = AveTimeOut.(string)
	}
	if ConfigInfraPg, ok := d.GetOk("config_infra_pg"); ok {
		vmmDomPAttr.ConfigInfraPg = ConfigInfraPg.(string)
	}
	if CtrlKnob, ok := d.GetOk("ctrl_knob"); ok {
		vmmDomPAttr.CtrlKnob = CtrlKnob.(string)
	}
	if Delimiter, ok := d.GetOk("delimiter"); ok {
		vmmDomPAttr.Delimiter = Delimiter.(string)
	}
	if EnableAVE, ok := d.GetOk("enable_ave"); ok {
		vmmDomPAttr.EnableAVE = EnableAVE.(string)
	}
	if EnableTag, ok := d.GetOk("enable_tag"); ok {
		vmmDomPAttr.EnableTag = EnableTag.(string)
	}
	if EncapMode, ok := d.GetOk("encap_mode"); ok {
		vmmDomPAttr.EncapMode = EncapMode.(string)
	}
	if EnfPref, ok := d.GetOk("enf_pref"); ok {
		vmmDomPAttr.EnfPref = EnfPref.(string)
	}
	if EpInventoryType, ok := d.GetOk("ep_inventory_type"); ok {
		vmmDomPAttr.EpInventoryType = EpInventoryType.(string)
	}
	if EpRetTime, ok := d.GetOk("ep_ret_time"); ok {
		vmmDomPAttr.EpRetTime = EpRetTime.(string)
	}
	if HvAvailMonitor, ok := d.GetOk("hv_avail_monitor"); ok {
		vmmDomPAttr.HvAvailMonitor = HvAvailMonitor.(string)
	}
	if McastAddr, ok := d.GetOk("mcast_addr"); ok {
		vmmDomPAttr.McastAddr = McastAddr.(string)
	}
	if Mode, ok := d.GetOk("mode"); ok {
		vmmDomPAttr.Mode = Mode.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vmmDomPAttr.NameAlias = NameAlias.(string)
	}
	if PrefEncapMode, ok := d.GetOk("pref_encap_mode"); ok {
		vmmDomPAttr.PrefEncapMode = PrefEncapMode.(string)
	}
	vmmDomP := models.NewVMMDomain(fmt.Sprintf("dom-%s", name), ProviderProfileDn, vmmDomPAttr)

	vmmDomP.Status = "modified"

	err := aciClient.Save(vmmDomP)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_vmm_rs_pref_enhanced_lag_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_pref_enhanced_lag_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_vlan_ns") {
		_, newRelParam := d.GetChange("relation_infra_rs_vlan_ns")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_vmm_rs_dom_mcast_addr_ns") {
		_, newRelParam := d.GetChange("relation_vmm_rs_dom_mcast_addr_ns")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_vmm_rs_default_cdp_if_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_default_cdp_if_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_vmm_rs_default_lacp_lag_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_default_lacp_lag_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_vlan_ns_def") {
		_, newRelParam := d.GetChange("relation_infra_rs_vlan_ns_def")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_vip_addr_ns") {
		_, newRelParam := d.GetChange("relation_infra_rs_vip_addr_ns")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_vmm_rs_default_lldp_if_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_default_lldp_if_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_vmm_rs_default_stp_if_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_default_stp_if_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_dom_vxlan_ns_def") {
		_, newRelParam := d.GetChange("relation_infra_rs_dom_vxlan_ns_def")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_vmm_rs_default_fw_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_default_fw_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_vmm_rs_default_l2_inst_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_default_l2_inst_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_vmm_rs_pref_enhanced_lag_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_pref_enhanced_lag_pol")
		err = aciClient.DeleteRelationvmmRsPrefEnhancedLagPolFromVMMDomain(vmmDomP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvmmRsPrefEnhancedLagPolFromVMMDomain(vmmDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_infra_rs_vlan_ns") {
		_, newRelParam := d.GetChange("relation_infra_rs_vlan_ns")
		err = aciClient.DeleteRelationinfraRsVlanNsFromVMMDomain(vmmDomP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsVlanNsFromVMMDomain(vmmDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_vmm_rs_dom_mcast_addr_ns") {
		_, newRelParam := d.GetChange("relation_vmm_rs_dom_mcast_addr_ns")
		err = aciClient.DeleteRelationvmmRsDomMcastAddrNsFromVMMDomain(vmmDomP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvmmRsDomMcastAddrNsFromVMMDomain(vmmDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_vmm_rs_default_cdp_if_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_default_cdp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationvmmRsDefaultCdpIfPolFromVMMDomain(vmmDomP.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_vmm_rs_default_lacp_lag_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_default_lacp_lag_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationvmmRsDefaultLacpLagPolFromVMMDomain(vmmDomP.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_infra_rs_vlan_ns_def") {
		_, newRelParam := d.GetChange("relation_infra_rs_vlan_ns_def")
		err = aciClient.CreateRelationinfraRsVlanNsDefFromVMMDomain(vmmDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_infra_rs_vip_addr_ns") {
		_, newRelParam := d.GetChange("relation_infra_rs_vip_addr_ns")
		err = aciClient.DeleteRelationinfraRsVipAddrNsFromVMMDomain(vmmDomP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsVipAddrNsFromVMMDomain(vmmDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_vmm_rs_default_lldp_if_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_default_lldp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationvmmRsDefaultLldpIfPolFromVMMDomain(vmmDomP.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_vmm_rs_default_stp_if_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_default_stp_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationvmmRsDefaultStpIfPolFromVMMDomain(vmmDomP.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_infra_rs_dom_vxlan_ns_def") {
		_, newRelParam := d.GetChange("relation_infra_rs_dom_vxlan_ns_def")
		err = aciClient.CreateRelationinfraRsDomVxlanNsDefFromVMMDomain(vmmDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_vmm_rs_default_fw_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_default_fw_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationvmmRsDefaultFwPolFromVMMDomain(vmmDomP.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_vmm_rs_default_l2_inst_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_default_l2_inst_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationvmmRsDefaultL2InstPolFromVMMDomain(vmmDomP.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(vmmDomP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciVMMDomainRead(ctx, d, m)

}

func resourceAciVMMDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	vmmDomP, err := getRemoteVMMDomain(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setVMMDomainAttributes(vmmDomP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	vmmRsPrefEnhancedLagPolData, err := aciClient.ReadRelationvmmRsPrefEnhancedLagPolFromVMMDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsPrefEnhancedLagPol %v", err)
		d.Set("relation_vmm_rs_pref_enhanced_lag_pol", "")

	} else {
		if _, ok := d.GetOk("relation_vmm_rs_pref_enhanced_lag_pol"); ok {
			tfName := d.Get("relation_vmm_rs_pref_enhanced_lag_pol").(string)
			if tfName != vmmRsPrefEnhancedLagPolData {
				d.Set("relation_vmm_rs_pref_enhanced_lag_pol", "")
			}
		}
	}

	infraRsVlanNsData, err := aciClient.ReadRelationinfraRsVlanNsFromVMMDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsVlanNs %v", err)
		d.Set("relation_infra_rs_vlan_ns", "")

	} else {
		if _, ok := d.GetOk("relation_infra_rs_vlan_ns"); ok {
			tfName := d.Get("relation_infra_rs_vlan_ns").(string)
			if tfName != infraRsVlanNsData {
				d.Set("relation_infra_rs_vlan_ns", "")
			}
		}
	}

	vmmRsDomMcastAddrNsData, err := aciClient.ReadRelationvmmRsDomMcastAddrNsFromVMMDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsDomMcastAddrNs %v", err)
		d.Set("relation_vmm_rs_dom_mcast_addr_ns", "")

	} else {
		if _, ok := d.GetOk("relation_vmm_rs_dom_mcast_addr_ns"); ok {
			tfName := d.Get("relation_vmm_rs_dom_mcast_addr_ns").(string)
			if tfName != vmmRsDomMcastAddrNsData {
				d.Set("relation_vmm_rs_dom_mcast_addr_ns", "")
			}
		}
	}

	vmmRsDefaultCdpIfPolData, err := aciClient.ReadRelationvmmRsDefaultCdpIfPolFromVMMDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsDefaultCdpIfPol %v", err)
		d.Set("relation_vmm_rs_default_cdp_if_pol", "")

	} else {
		if _, ok := d.GetOk("relation_vmm_rs_default_cdp_if_pol"); ok {
			tfName := GetMOName(d.Get("relation_vmm_rs_default_cdp_if_pol").(string))
			if tfName != vmmRsDefaultCdpIfPolData {
				d.Set("relation_vmm_rs_default_cdp_if_pol", "")
			}
		}
	}

	vmmRsDefaultLacpLagPolData, err := aciClient.ReadRelationvmmRsDefaultLacpLagPolFromVMMDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsDefaultLacpLagPol %v", err)
		d.Set("relation_vmm_rs_default_lacp_lag_pol", "")

	} else {
		if _, ok := d.GetOk("relation_vmm_rs_default_lacp_lag_pol"); ok {
			tfName := GetMOName(d.Get("relation_vmm_rs_default_lacp_lag_pol").(string))
			if tfName != vmmRsDefaultLacpLagPolData {
				d.Set("relation_vmm_rs_default_lacp_lag_pol", "")
			}
		}
	}

	infraRsVlanNsDefData, err := aciClient.ReadRelationinfraRsVlanNsDefFromVMMDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsVlanNsDef %v", err)
		d.Set("relation_infra_rs_vlan_ns_def", "")

	} else {
		if _, ok := d.GetOk("relation_infra_rs_vlan_ns_def"); ok {
			tfName := d.Get("relation_infra_rs_vlan_ns_def").(string)
			if tfName != infraRsVlanNsDefData {
				d.Set("relation_infra_rs_vlan_ns_def", "")
			}
		}
	}

	infraRsVipAddrNsData, err := aciClient.ReadRelationinfraRsVipAddrNsFromVMMDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsVipAddrNs %v", err)
		d.Set("relation_infra_rs_vip_addr_ns", "")

	} else {
		if _, ok := d.GetOk("relation_infra_rs_vip_addr_ns"); ok {
			tfName := d.Get("relation_infra_rs_vip_addr_ns").(string)
			if tfName != infraRsVipAddrNsData {
				d.Set("relation_infra_rs_vip_addr_ns", "")
			}
		}
	}

	vmmRsDefaultLldpIfPolData, err := aciClient.ReadRelationvmmRsDefaultLldpIfPolFromVMMDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsDefaultLldpIfPol %v", err)
		d.Set("relation_vmm_rs_default_lldp_if_pol", "")

	} else {
		if _, ok := d.GetOk("relation_vmm_rs_default_lldp_if_pol"); ok {
			tfName := GetMOName(d.Get("relation_vmm_rs_default_lldp_if_pol").(string))
			if tfName != vmmRsDefaultLldpIfPolData {
				d.Set("relation_vmm_rs_default_lldp_if_pol", "")
			}
		}
	}

	vmmRsDefaultStpIfPolData, err := aciClient.ReadRelationvmmRsDefaultStpIfPolFromVMMDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsDefaultStpIfPol %v", err)
		d.Set("relation_vmm_rs_default_stp_if_pol", "")

	} else {
		if _, ok := d.GetOk("relation_vmm_rs_default_stp_if_pol"); ok {
			tfName := GetMOName(d.Get("relation_vmm_rs_default_stp_if_pol").(string))
			if tfName != vmmRsDefaultStpIfPolData {
				d.Set("relation_vmm_rs_default_stp_if_pol", "")
			}
		}
	}

	infraRsDomVxlanNsDefData, err := aciClient.ReadRelationinfraRsDomVxlanNsDefFromVMMDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsDomVxlanNsDef %v", err)
		d.Set("relation_infra_rs_dom_vxlan_ns_def", "")

	} else {
		if _, ok := d.GetOk("relation_infra_rs_dom_vxlan_ns_def"); ok {
			tfName := d.Get("relation_infra_rs_dom_vxlan_ns_def").(string)
			if tfName != infraRsDomVxlanNsDefData {
				d.Set("relation_infra_rs_dom_vxlan_ns_def", "")
			}
		}
	}

	vmmRsDefaultFwPolData, err := aciClient.ReadRelationvmmRsDefaultFwPolFromVMMDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsDefaultFwPol %v", err)
		d.Set("relation_vmm_rs_default_fw_pol", "")

	} else {
		if _, ok := d.GetOk("relation_vmm_rs_default_fw_pol"); ok {
			tfName := GetMOName(d.Get("relation_vmm_rs_default_fw_pol").(string))
			if tfName != vmmRsDefaultFwPolData {
				d.Set("relation_vmm_rs_default_fw_pol", "")
			}
		}
	}

	vmmRsDefaultL2InstPolData, err := aciClient.ReadRelationvmmRsDefaultL2InstPolFromVMMDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsDefaultL2InstPol %v", err)
		d.Set("relation_vmm_rs_default_l2_inst_pol", "")

	} else {
		if _, ok := d.GetOk("relation_vmm_rs_default_l2_inst_pol"); ok {
			tfName := GetMOName(d.Get("relation_vmm_rs_default_l2_inst_pol").(string))
			if tfName != vmmRsDefaultL2InstPolData {
				d.Set("relation_vmm_rs_default_l2_inst_pol", "")
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciVMMDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vmmDomP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
