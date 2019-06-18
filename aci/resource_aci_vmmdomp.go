package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAciVMMDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciVMMDomainCreate,
		Update: resourceAciVMMDomainUpdate,
		Read:   resourceAciVMMDomainRead,
		Delete: resourceAciVMMDomainDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciVMMDomainImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"provider_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"access_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"arp_learning": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			},

			"ctrl_knob": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			},

			"enable_tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"encap_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"enf_pref": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ep_inventory_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
		}),
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

func setVMMDomainAttributes(vmmDomP *models.VMMDomain, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(vmmDomP.DistinguishedName)
	d.Set("description", vmmDomP.Description)
	d.Set("provider_profile_dn", GetParentDn(vmmDomP.DistinguishedName))
	vmmDomPMap, _ := vmmDomP.ToMap()

	d.Set("name", vmmDomPMap["name"])

	d.Set("access_mode", vmmDomPMap["accessMode"])
	d.Set("annotation", vmmDomPMap["annotation"])
	d.Set("arp_learning", vmmDomPMap["arpLearning"])
	d.Set("ave_time_out", vmmDomPMap["aveTimeOut"])
	d.Set("config_infra_pg", vmmDomPMap["configInfraPg"])
	d.Set("ctrl_knob", vmmDomPMap["ctrlKnob"])
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
	return d
}

func resourceAciVMMDomainImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	vmmDomP, err := getRemoteVMMDomain(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setVMMDomainAttributes(vmmDomP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciVMMDomainCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] VMMDomain: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	ProviderProfileDn := d.Get("provider_profile_dn").(string)

	vmmDomPAttr := models.VMMDomainAttributes{}
	if AccessMode, ok := d.GetOk("access_mode"); ok {
		vmmDomPAttr.AccessMode = AccessMode.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vmmDomPAttr.Annotation = Annotation.(string)
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
	vmmDomP := models.NewVMMDomain(fmt.Sprintf("dom-%s", name), ProviderProfileDn, desc, vmmDomPAttr)

	err := aciClient.Save(vmmDomP)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if relationTovmmRsPrefEnhancedLagPol, ok := d.GetOk("relation_vmm_rs_pref_enhanced_lag_pol"); ok {
		relationParam := relationTovmmRsPrefEnhancedLagPol.(string)
		err = aciClient.CreateRelationvmmRsPrefEnhancedLagPolFromVMMDomain(vmmDomP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vmm_rs_pref_enhanced_lag_pol")
		d.Partial(false)

	}
	if relationToinfraRsVlanNs, ok := d.GetOk("relation_infra_rs_vlan_ns"); ok {
		relationParam := relationToinfraRsVlanNs.(string)
		err = aciClient.CreateRelationinfraRsVlanNsFromVMMDomain(vmmDomP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_vlan_ns")
		d.Partial(false)

	}
	if relationTovmmRsDomMcastAddrNs, ok := d.GetOk("relation_vmm_rs_dom_mcast_addr_ns"); ok {
		relationParam := relationTovmmRsDomMcastAddrNs.(string)
		err = aciClient.CreateRelationvmmRsDomMcastAddrNsFromVMMDomain(vmmDomP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vmm_rs_dom_mcast_addr_ns")
		d.Partial(false)

	}
	if relationTovmmRsDefaultCdpIfPol, ok := d.GetOk("relation_vmm_rs_default_cdp_if_pol"); ok {
		relationParam := relationTovmmRsDefaultCdpIfPol.(string)
		err = aciClient.CreateRelationvmmRsDefaultCdpIfPolFromVMMDomain(vmmDomP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vmm_rs_default_cdp_if_pol")
		d.Partial(false)

	}
	if relationTovmmRsDefaultLacpLagPol, ok := d.GetOk("relation_vmm_rs_default_lacp_lag_pol"); ok {
		relationParam := relationTovmmRsDefaultLacpLagPol.(string)
		err = aciClient.CreateRelationvmmRsDefaultLacpLagPolFromVMMDomain(vmmDomP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vmm_rs_default_lacp_lag_pol")
		d.Partial(false)

	}
	if relationToinfraRsVlanNsDef, ok := d.GetOk("relation_infra_rs_vlan_ns_def"); ok {
		relationParam := relationToinfraRsVlanNsDef.(string)
		err = aciClient.CreateRelationinfraRsVlanNsDefFromVMMDomain(vmmDomP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_vlan_ns_def")
		d.Partial(false)

	}
	if relationToinfraRsVipAddrNs, ok := d.GetOk("relation_infra_rs_vip_addr_ns"); ok {
		relationParam := relationToinfraRsVipAddrNs.(string)
		err = aciClient.CreateRelationinfraRsVipAddrNsFromVMMDomain(vmmDomP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_vip_addr_ns")
		d.Partial(false)

	}
	if relationTovmmRsDefaultLldpIfPol, ok := d.GetOk("relation_vmm_rs_default_lldp_if_pol"); ok {
		relationParam := relationTovmmRsDefaultLldpIfPol.(string)
		err = aciClient.CreateRelationvmmRsDefaultLldpIfPolFromVMMDomain(vmmDomP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vmm_rs_default_lldp_if_pol")
		d.Partial(false)

	}
	if relationTovmmRsDefaultStpIfPol, ok := d.GetOk("relation_vmm_rs_default_stp_if_pol"); ok {
		relationParam := relationTovmmRsDefaultStpIfPol.(string)
		err = aciClient.CreateRelationvmmRsDefaultStpIfPolFromVMMDomain(vmmDomP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vmm_rs_default_stp_if_pol")
		d.Partial(false)

	}
	if relationToinfraRsDomVxlanNsDef, ok := d.GetOk("relation_infra_rs_dom_vxlan_ns_def"); ok {
		relationParam := relationToinfraRsDomVxlanNsDef.(string)
		err = aciClient.CreateRelationinfraRsDomVxlanNsDefFromVMMDomain(vmmDomP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_dom_vxlan_ns_def")
		d.Partial(false)

	}
	if relationTovmmRsDefaultFwPol, ok := d.GetOk("relation_vmm_rs_default_fw_pol"); ok {
		relationParam := relationTovmmRsDefaultFwPol.(string)
		err = aciClient.CreateRelationvmmRsDefaultFwPolFromVMMDomain(vmmDomP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vmm_rs_default_fw_pol")
		d.Partial(false)

	}
	if relationTovmmRsDefaultL2InstPol, ok := d.GetOk("relation_vmm_rs_default_l2_inst_pol"); ok {
		relationParam := relationTovmmRsDefaultL2InstPol.(string)
		err = aciClient.CreateRelationvmmRsDefaultL2InstPolFromVMMDomain(vmmDomP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vmm_rs_default_l2_inst_pol")
		d.Partial(false)

	}

	d.SetId(vmmDomP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciVMMDomainRead(d, m)
}

func resourceAciVMMDomainUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] VMMDomain: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	ProviderProfileDn := d.Get("provider_profile_dn").(string)

	vmmDomPAttr := models.VMMDomainAttributes{}
	if AccessMode, ok := d.GetOk("access_mode"); ok {
		vmmDomPAttr.AccessMode = AccessMode.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vmmDomPAttr.Annotation = Annotation.(string)
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
	vmmDomP := models.NewVMMDomain(fmt.Sprintf("dom-%s", name), ProviderProfileDn, desc, vmmDomPAttr)

	vmmDomP.Status = "modified"

	err := aciClient.Save(vmmDomP)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if d.HasChange("relation_vmm_rs_pref_enhanced_lag_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_pref_enhanced_lag_pol")
		err = aciClient.DeleteRelationvmmRsPrefEnhancedLagPolFromVMMDomain(vmmDomP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationvmmRsPrefEnhancedLagPolFromVMMDomain(vmmDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vmm_rs_pref_enhanced_lag_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_vlan_ns") {
		_, newRelParam := d.GetChange("relation_infra_rs_vlan_ns")
		err = aciClient.DeleteRelationinfraRsVlanNsFromVMMDomain(vmmDomP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationinfraRsVlanNsFromVMMDomain(vmmDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_vlan_ns")
		d.Partial(false)

	}
	if d.HasChange("relation_vmm_rs_dom_mcast_addr_ns") {
		_, newRelParam := d.GetChange("relation_vmm_rs_dom_mcast_addr_ns")
		err = aciClient.DeleteRelationvmmRsDomMcastAddrNsFromVMMDomain(vmmDomP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationvmmRsDomMcastAddrNsFromVMMDomain(vmmDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vmm_rs_dom_mcast_addr_ns")
		d.Partial(false)

	}
	if d.HasChange("relation_vmm_rs_default_cdp_if_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_default_cdp_if_pol")
		err = aciClient.CreateRelationvmmRsDefaultCdpIfPolFromVMMDomain(vmmDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vmm_rs_default_cdp_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_vmm_rs_default_lacp_lag_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_default_lacp_lag_pol")
		err = aciClient.CreateRelationvmmRsDefaultLacpLagPolFromVMMDomain(vmmDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vmm_rs_default_lacp_lag_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_vlan_ns_def") {
		_, newRelParam := d.GetChange("relation_infra_rs_vlan_ns_def")
		err = aciClient.CreateRelationinfraRsVlanNsDefFromVMMDomain(vmmDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_vlan_ns_def")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_vip_addr_ns") {
		_, newRelParam := d.GetChange("relation_infra_rs_vip_addr_ns")
		err = aciClient.DeleteRelationinfraRsVipAddrNsFromVMMDomain(vmmDomP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationinfraRsVipAddrNsFromVMMDomain(vmmDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_vip_addr_ns")
		d.Partial(false)

	}
	if d.HasChange("relation_vmm_rs_default_lldp_if_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_default_lldp_if_pol")
		err = aciClient.CreateRelationvmmRsDefaultLldpIfPolFromVMMDomain(vmmDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vmm_rs_default_lldp_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_vmm_rs_default_stp_if_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_default_stp_if_pol")
		err = aciClient.CreateRelationvmmRsDefaultStpIfPolFromVMMDomain(vmmDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vmm_rs_default_stp_if_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_dom_vxlan_ns_def") {
		_, newRelParam := d.GetChange("relation_infra_rs_dom_vxlan_ns_def")
		err = aciClient.CreateRelationinfraRsDomVxlanNsDefFromVMMDomain(vmmDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_dom_vxlan_ns_def")
		d.Partial(false)

	}
	if d.HasChange("relation_vmm_rs_default_fw_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_default_fw_pol")
		err = aciClient.CreateRelationvmmRsDefaultFwPolFromVMMDomain(vmmDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vmm_rs_default_fw_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_vmm_rs_default_l2_inst_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_default_l2_inst_pol")
		err = aciClient.CreateRelationvmmRsDefaultL2InstPolFromVMMDomain(vmmDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vmm_rs_default_l2_inst_pol")
		d.Partial(false)

	}

	d.SetId(vmmDomP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciVMMDomainRead(d, m)

}

func resourceAciVMMDomainRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	vmmDomP, err := getRemoteVMMDomain(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setVMMDomainAttributes(vmmDomP, d)

	vmmRsPrefEnhancedLagPolData, err := aciClient.ReadRelationvmmRsPrefEnhancedLagPolFromVMMDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsPrefEnhancedLagPol %v", err)

	} else {
		d.Set("relation_vmm_rs_pref_enhanced_lag_pol", vmmRsPrefEnhancedLagPolData)
	}

	infraRsVlanNsData, err := aciClient.ReadRelationinfraRsVlanNsFromVMMDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsVlanNs %v", err)

	} else {
		d.Set("relation_infra_rs_vlan_ns", infraRsVlanNsData)
	}

	vmmRsDomMcastAddrNsData, err := aciClient.ReadRelationvmmRsDomMcastAddrNsFromVMMDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsDomMcastAddrNs %v", err)

	} else {
		d.Set("relation_vmm_rs_dom_mcast_addr_ns", vmmRsDomMcastAddrNsData)
	}

	vmmRsDefaultCdpIfPolData, err := aciClient.ReadRelationvmmRsDefaultCdpIfPolFromVMMDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsDefaultCdpIfPol %v", err)

	} else {
		d.Set("relation_vmm_rs_default_cdp_if_pol", vmmRsDefaultCdpIfPolData)
	}

	vmmRsDefaultLacpLagPolData, err := aciClient.ReadRelationvmmRsDefaultLacpLagPolFromVMMDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsDefaultLacpLagPol %v", err)

	} else {
		d.Set("relation_vmm_rs_default_lacp_lag_pol", vmmRsDefaultLacpLagPolData)
	}

	infraRsVlanNsDefData, err := aciClient.ReadRelationinfraRsVlanNsDefFromVMMDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsVlanNsDef %v", err)

	} else {
		d.Set("relation_infra_rs_vlan_ns_def", infraRsVlanNsDefData)
	}

	infraRsVipAddrNsData, err := aciClient.ReadRelationinfraRsVipAddrNsFromVMMDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsVipAddrNs %v", err)

	} else {
		d.Set("relation_infra_rs_vip_addr_ns", infraRsVipAddrNsData)
	}

	vmmRsDefaultLldpIfPolData, err := aciClient.ReadRelationvmmRsDefaultLldpIfPolFromVMMDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsDefaultLldpIfPol %v", err)

	} else {
		d.Set("relation_vmm_rs_default_lldp_if_pol", vmmRsDefaultLldpIfPolData)
	}

	vmmRsDefaultStpIfPolData, err := aciClient.ReadRelationvmmRsDefaultStpIfPolFromVMMDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsDefaultStpIfPol %v", err)

	} else {
		d.Set("relation_vmm_rs_default_stp_if_pol", vmmRsDefaultStpIfPolData)
	}

	infraRsDomVxlanNsDefData, err := aciClient.ReadRelationinfraRsDomVxlanNsDefFromVMMDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsDomVxlanNsDef %v", err)

	} else {
		d.Set("relation_infra_rs_dom_vxlan_ns_def", infraRsDomVxlanNsDefData)
	}

	vmmRsDefaultFwPolData, err := aciClient.ReadRelationvmmRsDefaultFwPolFromVMMDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsDefaultFwPol %v", err)

	} else {
		d.Set("relation_vmm_rs_default_fw_pol", vmmRsDefaultFwPolData)
	}

	vmmRsDefaultL2InstPolData, err := aciClient.ReadRelationvmmRsDefaultL2InstPolFromVMMDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsDefaultL2InstPol %v", err)

	} else {
		d.Set("relation_vmm_rs_default_l2_inst_pol", vmmRsDefaultL2InstPolData)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciVMMDomainDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vmmDomP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
