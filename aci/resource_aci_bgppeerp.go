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

func resourceAciBgpPeerConnectivityProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciBgpPeerConnectivityProfileCreate,
		UpdateContext: resourceAciBgpPeerConnectivityProfileUpdate,
		ReadContext:   resourceAciBgpPeerConnectivityProfileRead,
		DeleteContext: resourceAciBgpPeerConnectivityProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciBgpPeerConnectivityProfileImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"parent_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"addr": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"addr_t_ctrl": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressBitMaskDiffFunc(),
				ValidateFunc: schema.SchemaValidateFunc(validateCommaSeparatedStringInSlice([]string{
					"af-mcast",
					"af-ucast",
					"",
				}, false, "")),
			},

			"allowed_self_as_cnt": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ctrl": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressBitMaskDiffFunc(),
				ValidateFunc: schema.SchemaValidateFunc(validateCommaSeparatedStringInSlice([]string{
					"allow-self-as",
					"as-override",
					"dis-peer-as-check",
					"nh-self",
					"send-com",
					"send-ext-com",
					"",
				}, false, "")),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"password": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},

			"peer_ctrl": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressBitMaskDiffFunc(),
				ValidateFunc: schema.SchemaValidateFunc(validateCommaSeparatedStringInSlice([]string{
					"bfd",
					"dis-conn-check",
					"",
				}, false, "")),
			},

			"private_a_sctrl": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressBitMaskDiffFunc(),
				ValidateFunc: schema.SchemaValidateFunc(validateCommaSeparatedStringInSlice([]string{
					"remove-all",
					"remove-exclusive",
					"replace-as",
					"",
				}, false, "")),
			},

			"ttl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"weight": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"as_number": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"local_asn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"local_asn_propagate": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"local_asn"},

				ValidateFunc: validation.StringInSlice([]string{
					"dual-as",
					"no-prepend",
					"none",
					"replace-as",
				}, false),
			},

			"relation_bgp_rs_peer_pfx_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}

func getRemoteBgpPeerConnectivityProfile(client *client.Client, dn string) (*models.BgpPeerConnectivityProfile, error) {
	bgpPeerPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	bgpPeerP := models.BgpPeerConnectivityProfileFromContainer(bgpPeerPCont)

	if bgpPeerP.DistinguishedName == "" {
		return nil, fmt.Errorf("BgpPeerConnectivityProfile %s not found", bgpPeerP.DistinguishedName)
	}

	return bgpPeerP, nil
}

func getRemoteLocalAutonomousSystemProfileFromBgpPeerConnectivityProfile(client *client.Client, dn string) (*models.LocalAutonomousSystemProfile, error) {
	bgpLocalAsnPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	bgpLocalAsnP := models.LocalAutonomousSystemProfileFromContainer(bgpLocalAsnPCont)

	if bgpLocalAsnP.DistinguishedName == "" {
		return nil, fmt.Errorf("LocalAutonomousSystemProfile %s not found", bgpLocalAsnP.DistinguishedName)
	}

	return bgpLocalAsnP, nil
}

func getRemoteBgpAutonomousSystemProfileFromBgpPeerConnectivityProfile(client *client.Client, dn string) (*models.BgpAutonomousSystemProfile, error) {
	bgpAsPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	bgpAsP := models.BgpAutonomousSystemProfileFromContainer(bgpAsPCont)

	if bgpAsP.DistinguishedName == "" {
		return nil, fmt.Errorf("BgpAutonomousSystemProfile %s not found", bgpAsP.DistinguishedName)
	}

	return bgpAsP, nil
}

func setBgpPeerConnectivityProfileAttributes(bgpPeerP *models.BgpPeerConnectivityProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(bgpPeerP.DistinguishedName)
	d.Set("description", bgpPeerP.Description)
	dn := d.Id()
	if dn != bgpPeerP.DistinguishedName {
		d.Set("parent_dn", "")
	}
	bgpPeerPMap, err := bgpPeerP.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("addr", bgpPeerPMap["addr"])

	d.Set("addr", bgpPeerPMap["addr"])
	d.Set("addr_t_ctrl", bgpPeerPMap["addrTCtrl"])
	d.Set("allowed_self_as_cnt", bgpPeerPMap["allowedSelfAsCnt"])
	d.Set("annotation", bgpPeerPMap["annotation"])
	d.Set("ctrl", bgpPeerPMap["ctrl"])
	d.Set("name_alias", bgpPeerPMap["nameAlias"])
	d.Set("peer_ctrl", bgpPeerPMap["peerCtrl"])
	d.Set("private_a_sctrl", bgpPeerPMap["privateASctrl"])
	d.Set("ttl", bgpPeerPMap["ttl"])
	d.Set("weight", bgpPeerPMap["weight"])
	return d, nil
}

func setBgpAutonomousSystemProfileAttributesFromBgpPeerConnectivityProfile(bgpAsP *models.BgpAutonomousSystemProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	bgpAsPMap, err := bgpAsP.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("as_number", bgpAsPMap["asn"])
	return d, nil
}

func setLocalAutonomousSystemProfileAttributesFromBgpPeerConnectivityProfile(bgpLocalAsnP *models.LocalAutonomousSystemProfile, d *schema.ResourceData) (*schema.ResourceData, error) {

	bgpLocalAsnPMap, err := bgpLocalAsnP.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("local_asn_propagate", bgpLocalAsnPMap["asnPropagate"])
	d.Set("local_asn", bgpLocalAsnPMap["localAsn"])
	return d, nil
}

func resourceAciBgpPeerConnectivityProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	bgpPeerP, err := getRemoteBgpPeerConnectivityProfile(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setBgpPeerConnectivityProfileAttributes(bgpPeerP, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciBgpPeerConnectivityProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] BgpPeerConnectivityProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	addr := d.Get("addr").(string)

	ParentDn := d.Get("parent_dn").(string)

	bgpPeerPAttr := models.BgpPeerConnectivityProfileAttributes{}
	if Addr, ok := d.GetOk("addr"); ok {
		bgpPeerPAttr.Addr = Addr.(string)
	}
	if AddrTCtrl, ok := d.GetOk("addr_t_ctrl"); ok {
		bgpPeerPAttr.AddrTCtrl = AddrTCtrl.(string)
	}
	if AllowedSelfAsCnt, ok := d.GetOk("allowed_self_as_cnt"); ok {
		bgpPeerPAttr.AllowedSelfAsCnt = AllowedSelfAsCnt.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		bgpPeerPAttr.Annotation = Annotation.(string)
	} else {
		bgpPeerPAttr.Annotation = "{}"
	}
	if Ctrl, ok := d.GetOk("ctrl"); ok {
		bgpPeerPAttr.Ctrl = Ctrl.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		bgpPeerPAttr.NameAlias = NameAlias.(string)
	}
	if Password, ok := d.GetOk("password"); ok {
		bgpPeerPAttr.Password = Password.(string)
	}
	if PeerCtrl, ok := d.GetOk("peer_ctrl"); ok {
		bgpPeerPAttr.PeerCtrl = PeerCtrl.(string)
	}
	if PrivateASctrl, ok := d.GetOk("private_a_sctrl"); ok {
		bgpPeerPAttr.PrivateASctrl = PrivateASctrl.(string)
	}
	if Ttl, ok := d.GetOk("ttl"); ok {
		bgpPeerPAttr.Ttl = Ttl.(string)
	}
	if Weight, ok := d.GetOk("weight"); ok {
		bgpPeerPAttr.Weight = Weight.(string)
	}
	bgpPeerP := models.NewBgpPeerConnectivityProfile(fmt.Sprintf("peerP-[%s]", addr), ParentDn, desc, bgpPeerPAttr)

	err := aciClient.Save(bgpPeerP)
	if err != nil {
		return diag.FromErr(err)
	}
	PeerConnectivityProfileDn := bgpPeerP.DistinguishedName

	if _, ok := d.GetOk("local_asn"); ok {
		bgpLocalAsnPAttr := models.LocalAutonomousSystemProfileAttributes{}

		if AsnPropagate, ok := d.GetOk("local_asn_propagate"); ok {
			bgpLocalAsnPAttr.AsnPropagate = AsnPropagate.(string)
		}
		if LocalAsn, ok := d.GetOk("local_asn"); ok {
			bgpLocalAsnPAttr.LocalAsn = LocalAsn.(string)
		}

		bgpLocalAsnP := models.NewLocalAutonomousSystemProfile(fmt.Sprintf("localasn"), PeerConnectivityProfileDn, desc, bgpLocalAsnPAttr)

		err = aciClient.Save(bgpLocalAsnP)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if _, ok := d.GetOk("as_number"); ok {
		bgpAsPAttr := models.BgpAutonomousSystemProfileAttributes{}

		if Asn, ok := d.GetOk("as_number"); ok {
			bgpAsPAttr.Asn = Asn.(string)
		}

		bgpAsP := models.NewBgpAutonomousSystemProfile(fmt.Sprintf("as"), PeerConnectivityProfileDn, desc, bgpAsPAttr)
		err = aciClient.Save(bgpAsP)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	checkDns := make([]string, 0, 1)

	if relationTobgpRsPeerPfxPol, ok := d.GetOk("relation_bgp_rs_peer_pfx_pol"); ok {
		relationParam := relationTobgpRsPeerPfxPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTobgpRsPeerPfxPol, ok := d.GetOk("relation_bgp_rs_peer_pfx_pol"); ok {
		relationParam := GetMOName(relationTobgpRsPeerPfxPol.(string))
		err = aciClient.CreateRelationbgpRsPeerPfxPolFromBgpPeerConnectivityProfile(bgpPeerP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(bgpPeerP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciBgpPeerConnectivityProfileRead(ctx, d, m)
}

func resourceAciBgpPeerConnectivityProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] BgpPeerConnectivityProfile: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	addr := d.Get("addr").(string)

	ParentDn := d.Get("parent_dn").(string)

	bgpPeerPAttr := models.BgpPeerConnectivityProfileAttributes{}
	if Addr, ok := d.GetOk("addr"); ok {
		bgpPeerPAttr.Addr = Addr.(string)
	}
	if AddrTCtrl, ok := d.GetOk("addr_t_ctrl"); ok {
		bgpPeerPAttr.AddrTCtrl = AddrTCtrl.(string)
	}
	if AllowedSelfAsCnt, ok := d.GetOk("allowed_self_as_cnt"); ok {
		bgpPeerPAttr.AllowedSelfAsCnt = AllowedSelfAsCnt.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		bgpPeerPAttr.Annotation = Annotation.(string)
	} else {
		bgpPeerPAttr.Annotation = "{}"
	}
	if Ctrl, ok := d.GetOk("ctrl"); ok {
		bgpPeerPAttr.Ctrl = Ctrl.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		bgpPeerPAttr.NameAlias = NameAlias.(string)
	}
	if Password, ok := d.GetOk("password"); ok {
		bgpPeerPAttr.Password = Password.(string)
	}
	if PeerCtrl, ok := d.GetOk("peer_ctrl"); ok {
		bgpPeerPAttr.PeerCtrl = PeerCtrl.(string)
	}
	if PrivateASctrl, ok := d.GetOk("private_a_sctrl"); ok {
		bgpPeerPAttr.PrivateASctrl = PrivateASctrl.(string)
	}
	if Ttl, ok := d.GetOk("ttl"); ok {
		bgpPeerPAttr.Ttl = Ttl.(string)
	}
	if Weight, ok := d.GetOk("weight"); ok {
		bgpPeerPAttr.Weight = Weight.(string)
	}
	bgpPeerP := models.NewBgpPeerConnectivityProfile(fmt.Sprintf("peerP-[%s]", addr), ParentDn, desc, bgpPeerPAttr)

	bgpPeerP.Status = "modified"

	err := aciClient.Save(bgpPeerP)

	if err != nil {
		return diag.FromErr(err)
	}

	PeerConnectivityProfileDn := bgpPeerP.DistinguishedName
	if _, ok := d.GetOk("local_asn"); ok {
		bgpLocalAsnPAttr := models.LocalAutonomousSystemProfileAttributes{}

		if AsnPropagate, ok := d.GetOk("local_asn_propagate"); ok {
			bgpLocalAsnPAttr.AsnPropagate = AsnPropagate.(string)
		}
		if LocalAsn, ok := d.GetOk("local_asn"); ok {
			bgpLocalAsnPAttr.LocalAsn = LocalAsn.(string)
		}

		bgpLocalAsnP := models.NewLocalAutonomousSystemProfile(fmt.Sprintf("localasn"), PeerConnectivityProfileDn, desc, bgpLocalAsnPAttr)

		err = aciClient.Save(bgpLocalAsnP)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if _, ok := d.GetOk("as_number"); ok {
		bgpAsPAttr := models.BgpAutonomousSystemProfileAttributes{}

		if Asn, ok := d.GetOk("as_number"); ok {
			bgpAsPAttr.Asn = Asn.(string)
		}

		bgpAsP := models.NewBgpAutonomousSystemProfile(fmt.Sprintf("as"), PeerConnectivityProfileDn, desc, bgpAsPAttr)
		err = aciClient.Save(bgpAsP)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_bgp_rs_peer_pfx_pol") {
		_, newRelParam := d.GetChange("relation_bgp_rs_peer_pfx_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_bgp_rs_peer_pfx_pol") {
		_, newRelParam := d.GetChange("relation_bgp_rs_peer_pfx_pol")
		err = aciClient.CreateRelationbgpRsPeerPfxPolFromBgpPeerConnectivityProfile(bgpPeerP.DistinguishedName, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(bgpPeerP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciBgpPeerConnectivityProfileRead(ctx, d, m)

}

func resourceAciBgpPeerConnectivityProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()

	bgpPeerP, err := getRemoteBgpPeerConnectivityProfile(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setBgpPeerConnectivityProfileAttributes(bgpPeerP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	if _, ok := d.GetOk("as_number"); ok {
		bgpAsP, err := getRemoteBgpAutonomousSystemProfileFromBgpPeerConnectivityProfile(aciClient, fmt.Sprintf("%s/as", dn))
		if err != nil {
			d.SetId("")
			return nil
		}
		setBgpAutonomousSystemProfileAttributesFromBgpPeerConnectivityProfile(bgpAsP, d)
	}

	if _, ok := d.GetOk("local_asn"); ok {
		bgpLocalAsnP, err := getRemoteLocalAutonomousSystemProfileFromBgpPeerConnectivityProfile(aciClient, fmt.Sprintf("%s/localasn", dn))
		if err != nil {
			d.SetId("")
			return nil
		}
		setLocalAutonomousSystemProfileAttributesFromBgpPeerConnectivityProfile(bgpLocalAsnP, d)
	}

	bgpRsPeerPfxPolData, err := aciClient.ReadRelationbgpRsPeerPfxPolFromBgpPeerConnectivityProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation bgpRsPeerPfxPol %v", err)
		d.Set("relation_bgp_rs_peer_pfx_pol", "")

	} else {
		if _, ok := d.GetOk("relation_bgp_rs_peer_pfx_pol"); ok {
			tfName := GetMOName(d.Get("relation_bgp_rs_peer_pfx_pol").(string))
			if tfName != bgpRsPeerPfxPolData {
				d.Set("relation_bgp_rs_peer_pfx_pol", "")
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciBgpPeerConnectivityProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "bgpPeerP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
