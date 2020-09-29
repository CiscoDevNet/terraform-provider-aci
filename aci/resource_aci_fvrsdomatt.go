package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAciDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciDomainCreate,
		Update: resourceAciDomainUpdate,
		Read:   resourceAciDomainRead,
		Delete: resourceAciDomainDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciDomainImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"application_epg_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"tdn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"binding_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"none",
					"staticBinding",
					"dynamicBinding",
					"ephemeral",
				}, false),
			},

			"allow_micro_seg": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"delimiter": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"encap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"encap_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"auto",
					"vlan",
					"vxlan",
				}, false),
			},

			"epg_cos": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Cos0",
					"Cos1",
					"Cos2",
					"Cos3",
					"Cos4",
					"Cos5",
					"Cos6",
					"Cos7",
				}, false),
			},

			"epg_cos_pref": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
				}, false),
			},

			"instr_imedcy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"immediate",
					"lazy",
				}, false),
			},

			"lag_policy_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"netflow_dir": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ingress",
					"egress",
					"both",
				}, false),
			},

			"netflow_pref": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
				}, false),
			},

			"num_ports": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"port_allocation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"primary_encap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"primary_encap_inner": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"res_imedcy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"immediate",
					"lazy",
					"pre-provision",
				}, false),
			},

			"secondary_encap_inner": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"switching_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"native",
					"AVE",
				}, false),
			},

			"vmm_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"vmm_allow_promiscuous": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"vmm_forged_transmits": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"vmm_mac_changes": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteDomain(client *client.Client, dn string) (*models.FVDomain, error) {
	fvRsDomAttCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fvRsDomAtt := models.FVDomainFromContainer(fvRsDomAttCont)

	if fvRsDomAtt.DistinguishedName == "" {
		return nil, fmt.Errorf("Domain %s not found", fvRsDomAtt.DistinguishedName)
	}

	return fvRsDomAtt, nil
}

func getRemoteVMMSecurityPolicy(client *client.Client, dn string) (*models.VMMSecurityPolicy, error) {
	vmmSecPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vmmSecP := models.VMMSecurityPolicyFromContainer(vmmSecPCont)

	if vmmSecP.DistinguishedName == "" {
		return nil, fmt.Errorf("VMMSecurityPolicy %s not found", vmmSecP.DistinguishedName)
	}

	return vmmSecP, nil
}

func setDomainAttributes(fvRsDomAtt *models.FVDomain, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(fvRsDomAtt.DistinguishedName)
	if dn != fvRsDomAtt.DistinguishedName {
		d.Set("application_epg_dn", "")
	}
	fvRsDomAttMap, _ := fvRsDomAtt.ToMap()

	d.Set("tdn", fvRsDomAttMap["tDn"])

	d.Set("annotation", fvRsDomAttMap["annotation"])
	d.Set("binding_type", fvRsDomAttMap["bindingType"])
	if fvRsDomAttMap["classPref"] == "useg" {
		d.Set("allow_micro_seg", true)
	} else {
		d.Set("allow_micro_seg", false)
	}
	d.Set("delimiter", fvRsDomAttMap["delimiter"])
	d.Set("encap", fvRsDomAttMap["encap"])
	d.Set("encap_mode", fvRsDomAttMap["encapMode"])
	d.Set("epg_cos", fvRsDomAttMap["epgCos"])
	d.Set("epg_cos_pref", fvRsDomAttMap["epgCosPref"])
	d.Set("instr_imedcy", fvRsDomAttMap["instrImedcy"])
	d.Set("lag_policy_name", fvRsDomAttMap["lagPolicyName"])
	d.Set("netflow_dir", fvRsDomAttMap["netflowDir"])
	d.Set("netflow_pref", fvRsDomAttMap["netflowPref"])
	d.Set("num_ports", fvRsDomAttMap["numPorts"])
	d.Set("port_allocation", fvRsDomAttMap["portAllocation"])
	d.Set("primary_encap", fvRsDomAttMap["primaryEncap"])
	d.Set("primary_encap_inner", fvRsDomAttMap["primaryEncapInner"])
	d.Set("res_imedcy", fvRsDomAttMap["resImedcy"])
	d.Set("secondary_encap_inner", fvRsDomAttMap["secondaryEncapInner"])
	d.Set("switching_mode", fvRsDomAttMap["switchingMode"])
	return d
}

func setVMMSecurityPolicyAttributes(vmmSecP *models.VMMSecurityPolicy, d *schema.ResourceData) *schema.ResourceData {
	vmmSecPMap, _ := vmmSecP.ToMap()
	d.Set("vmm_allow_promiscuous", vmmSecPMap["allowPromiscuous"])
	d.Set("vmm_forged_transmits", vmmSecPMap["forgedTransmits"])
	d.Set("vmm_mac_changes", vmmSecPMap["macChanges"])
	d.Set("vmm_id", vmmSecP.DistinguishedName)
	return d
}

func resourceAciDomainImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fvRsDomAtt, err := getRemoteDomain(aciClient, dn)

	if err != nil {
		return nil, err
	}
	fvRsDomAttMap, _ := fvRsDomAtt.ToMap()
	tDn := fvRsDomAttMap["tDn"]
	pDN := GetParentDn(dn, fmt.Sprintf("/rsdomAtt-[%s]", tDn))
	d.Set("application_epg_dn", pDN)
	schemaFilled := setDomainAttributes(fvRsDomAtt, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciDomainCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Domain: Beginning Creation")
	aciClient := m.(*client.Client)

	tDn := d.Get("tdn").(string)

	ApplicationEPGDn := d.Get("application_epg_dn").(string)

	fvRsDomAttAttr := models.FVDomainAttributes{}
	fvRsDomAttAttr.TDn = tDn
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvRsDomAttAttr.Annotation = Annotation.(string)
	} else {
		fvRsDomAttAttr.Annotation = "{}"
	}
	if BindingType, ok := d.GetOk("binding_type"); ok {
		fvRsDomAttAttr.BindingType = BindingType.(string)
	}
	if flag, ok := d.GetOk("allow_micro_seg"); ok {
		if flag.(bool) == true {
			fvRsDomAttAttr.ClassPref = "useg"
		} else {
			fvRsDomAttAttr.ClassPref = "encap"
		}
	} else {
		fvRsDomAttAttr.ClassPref = "encap"
	}
	if Delimiter, ok := d.GetOk("delimiter"); ok {
		fvRsDomAttAttr.Delimiter = Delimiter.(string)
	}
	if Encap, ok := d.GetOk("encap"); ok {
		fvRsDomAttAttr.Encap = Encap.(string)
	}
	if EncapMode, ok := d.GetOk("encap_mode"); ok {
		fvRsDomAttAttr.EncapMode = EncapMode.(string)
	}
	if EpgCos, ok := d.GetOk("epg_cos"); ok {
		fvRsDomAttAttr.EpgCos = EpgCos.(string)
	}
	if EpgCosPref, ok := d.GetOk("epg_cos_pref"); ok {
		fvRsDomAttAttr.EpgCosPref = EpgCosPref.(string)
	}
	if InstrImedcy, ok := d.GetOk("instr_imedcy"); ok {
		fvRsDomAttAttr.InstrImedcy = InstrImedcy.(string)
	}
	if LagPolicyName, ok := d.GetOk("lag_policy_name"); ok {
		fvRsDomAttAttr.LagPolicyName = LagPolicyName.(string)
	}
	if NetflowDir, ok := d.GetOk("netflow_dir"); ok {
		fvRsDomAttAttr.NetflowDir = NetflowDir.(string)
	}
	if NetflowPref, ok := d.GetOk("netflow_pref"); ok {
		fvRsDomAttAttr.NetflowPref = NetflowPref.(string)
	}
	if NumPorts, ok := d.GetOk("num_ports"); ok {
		fvRsDomAttAttr.NumPorts = NumPorts.(string)
	}
	if PortAllocation, ok := d.GetOk("port_allocation"); ok {
		fvRsDomAttAttr.PortAllocation = PortAllocation.(string)
	}
	if PrimaryEncap, ok := d.GetOk("primary_encap"); ok {
		fvRsDomAttAttr.PrimaryEncap = PrimaryEncap.(string)
	}
	if PrimaryEncapInner, ok := d.GetOk("primary_encap_inner"); ok {
		fvRsDomAttAttr.PrimaryEncapInner = PrimaryEncapInner.(string)
	}
	if ResImedcy, ok := d.GetOk("res_imedcy"); ok {
		fvRsDomAttAttr.ResImedcy = ResImedcy.(string)
	}
	if SecondaryEncapInner, ok := d.GetOk("secondary_encap_inner"); ok {
		fvRsDomAttAttr.SecondaryEncapInner = SecondaryEncapInner.(string)
	}
	if SwitchingMode, ok := d.GetOk("switching_mode"); ok {
		fvRsDomAttAttr.SwitchingMode = SwitchingMode.(string)
	}
	fvRsDomAtt := models.NewFVDomain(fmt.Sprintf("rsdomAtt-[%s]", tDn), ApplicationEPGDn, "", fvRsDomAttAttr)

	err := aciClient.Save(fvRsDomAtt)
	if err != nil {
		return err
	}

	vmmSecPAttr := models.VMMSecurityPolicyAttributes{}
	if AllowPromiscuous, ok := d.GetOk("vmm_allow_promiscuous"); ok {
		vmmSecPAttr.AllowPromiscuous = AllowPromiscuous.(string)
	}
	if ForgedTransmits, ok := d.GetOk("vmm_forged_transmits"); ok {
		vmmSecPAttr.ForgedTransmits = ForgedTransmits.(string)
	}
	if MacChanges, ok := d.GetOk("vmm_mac_changes"); ok {
		vmmSecPAttr.MacChanges = MacChanges.(string)
	}
	vmmSecP := models.NewVMMSecurityPolicy(fmt.Sprintf("sec"), fvRsDomAtt.DistinguishedName, "", vmmSecPAttr)

	err = aciClient.Save(vmmSecP)
	if err != nil {
		return err
	}
	d.Set("vmm_id", vmmSecP.DistinguishedName)

	d.Partial(true)

	d.SetPartial("tdn")

	d.Partial(false)

	d.SetId(fvRsDomAtt.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciDomainRead(d, m)
}

func resourceAciDomainUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Domain: Beginning Update")

	aciClient := m.(*client.Client)

	tDn := d.Get("tdn").(string)

	ApplicationEPGDn := d.Get("application_epg_dn").(string)

	fvRsDomAttAttr := models.FVDomainAttributes{}
	fvRsDomAttAttr.TDn = tDn
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvRsDomAttAttr.Annotation = Annotation.(string)
	} else {
		fvRsDomAttAttr.Annotation = "{}"
	}
	if BindingType, ok := d.GetOk("binding_type"); ok {
		fvRsDomAttAttr.BindingType = BindingType.(string)
	}
	if flag, ok := d.GetOk("allow_micro_seg"); ok {
		if flag.(bool) == true {
			fvRsDomAttAttr.ClassPref = "useg"
		} else {
			fvRsDomAttAttr.ClassPref = "encap"
		}
	} else {
		fvRsDomAttAttr.ClassPref = "encap"
	}
	if Delimiter, ok := d.GetOk("delimiter"); ok {
		fvRsDomAttAttr.Delimiter = Delimiter.(string)
	}
	if Encap, ok := d.GetOk("encap"); ok {
		fvRsDomAttAttr.Encap = Encap.(string)
	}
	if EncapMode, ok := d.GetOk("encap_mode"); ok {
		fvRsDomAttAttr.EncapMode = EncapMode.(string)
	}
	if EpgCos, ok := d.GetOk("epg_cos"); ok {
		fvRsDomAttAttr.EpgCos = EpgCos.(string)
	}
	if EpgCosPref, ok := d.GetOk("epg_cos_pref"); ok {
		fvRsDomAttAttr.EpgCosPref = EpgCosPref.(string)
	}
	if InstrImedcy, ok := d.GetOk("instr_imedcy"); ok {
		fvRsDomAttAttr.InstrImedcy = InstrImedcy.(string)
	}
	if LagPolicyName, ok := d.GetOk("lag_policy_name"); ok {
		fvRsDomAttAttr.LagPolicyName = LagPolicyName.(string)
	}
	if NetflowDir, ok := d.GetOk("netflow_dir"); ok {
		fvRsDomAttAttr.NetflowDir = NetflowDir.(string)
	}
	if NetflowPref, ok := d.GetOk("netflow_pref"); ok {
		fvRsDomAttAttr.NetflowPref = NetflowPref.(string)
	}
	if NumPorts, ok := d.GetOk("num_ports"); ok {
		fvRsDomAttAttr.NumPorts = NumPorts.(string)
	}
	if PortAllocation, ok := d.GetOk("port_allocation"); ok {
		fvRsDomAttAttr.PortAllocation = PortAllocation.(string)
	}
	if PrimaryEncap, ok := d.GetOk("primary_encap"); ok {
		fvRsDomAttAttr.PrimaryEncap = PrimaryEncap.(string)
	}
	if PrimaryEncapInner, ok := d.GetOk("primary_encap_inner"); ok {
		fvRsDomAttAttr.PrimaryEncapInner = PrimaryEncapInner.(string)
	}
	if ResImedcy, ok := d.GetOk("res_imedcy"); ok {
		fvRsDomAttAttr.ResImedcy = ResImedcy.(string)
	}
	if SecondaryEncapInner, ok := d.GetOk("secondary_encap_inner"); ok {
		fvRsDomAttAttr.SecondaryEncapInner = SecondaryEncapInner.(string)
	}
	if SwitchingMode, ok := d.GetOk("switching_mode"); ok {
		fvRsDomAttAttr.SwitchingMode = SwitchingMode.(string)
	}
	fvRsDomAtt := models.NewFVDomain(fmt.Sprintf("rsdomAtt-[%s]", tDn), ApplicationEPGDn, "", fvRsDomAttAttr)

	fvRsDomAtt.Status = "modified"

	err := aciClient.Save(fvRsDomAtt)
	if err != nil {
		return err
	}

	vmmSecPAttr := models.VMMSecurityPolicyAttributes{}
	if AllowPromiscuous, ok := d.GetOk("vmm_allow_promiscuous"); ok {
		vmmSecPAttr.AllowPromiscuous = AllowPromiscuous.(string)
	}
	if ForgedTransmits, ok := d.GetOk("vmm_forged_transmits"); ok {
		vmmSecPAttr.ForgedTransmits = ForgedTransmits.(string)
	}
	if MacChanges, ok := d.GetOk("vmm_mac_changes"); ok {
		vmmSecPAttr.MacChanges = MacChanges.(string)
	}
	vmmSecP := models.NewVMMSecurityPolicy(fmt.Sprintf("sec"), fvRsDomAtt.DistinguishedName, "", vmmSecPAttr)

	vmmSecP.Status = "modified"

	err = aciClient.Save(vmmSecP)
	if err != nil {
		return err
	}
	d.Set("vmm_id", vmmSecP.DistinguishedName)

	d.Partial(true)

	d.SetPartial("tdn")

	d.Partial(false)

	d.SetId(fvRsDomAtt.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciDomainRead(d, m)

}

func resourceAciDomainRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fvRsDomAtt, err := getRemoteDomain(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setDomainAttributes(fvRsDomAtt, d)

	if d.Get("vmm_id") != nil {
		vmmDn := d.Get("vmm_id").(string)
		vmmSecP, err := getRemoteVMMSecurityPolicy(aciClient, vmmDn)
		if err == nil {
			setVMMSecurityPolicyAttributes(vmmSecP, d)
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciDomainDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fvRsDomAtt")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
