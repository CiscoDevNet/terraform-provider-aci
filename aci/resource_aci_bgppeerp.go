package aci

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

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
			"logical_node_profile_dn": &schema.Schema{
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "use parent_dn instead",
			},

			"parent_dn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"addr": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"addr_t_ctrl": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"af-ucast",
						"af-mcast",
					}, false),
				},
			},

			"admin_state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
				}, false),
			},

			"allowed_self_as_cnt": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ctrl": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"allow-self-as",
						"as-override",
						"dis-peer-as-check",
						"nh-self",
						"send-com",
						"send-ext-com",
					}, false),
				},
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
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"bfd",
						"dis-conn-check",
					}, false),
				},
			},

			"private_a_sctrl": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"remove-all",
						"remove-exclusive",
						"replace-as",
					}, false),
				},
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
				Type:     schema.TypeString,
				Default:  "uni/tn-common/bgpPfxP-default",
				Optional: true,
			},
			"relation_bgp_rs_peer_to_profile": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Create relation to rtctrlProfile",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"direction": {
							Required: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"export",
								"import",
							}, false),
						},

						"target_dn": {
							Required: true,
							Type:     schema.TypeString,
						},
					},
				},
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
	d.Set("parent_dn", GetParentDn(dn, fmt.Sprintf("/peerP-[%s]", bgpPeerPMap["addr"])))

	d.Set("admin_state", bgpPeerPMap["adminSt"])
	d.Set("addr", bgpPeerPMap["addr"])
	addrTCtrlGet := make([]string, 0, 1)
	for _, val := range strings.Split(bgpPeerPMap["addrTCtrl"], ",") {
		addrTCtrlGet = append(addrTCtrlGet, strings.Trim(val, " "))
	}
	sort.Strings(addrTCtrlGet)
	if len(addrTCtrlGet) == 1 && addrTCtrlGet[0] == "" {
		d.Set("addr_t_ctrl", make([]string, 0, 1))
	} else {
		d.Set("addr_t_ctrl", addrTCtrlGet)
	}
	d.Set("allowed_self_as_cnt", bgpPeerPMap["allowedSelfAsCnt"])
	d.Set("annotation", bgpPeerPMap["annotation"])
	ctrlGet := make([]string, 0, 1)
	for _, val := range strings.Split(bgpPeerPMap["ctrl"], ",") {
		ctrlGet = append(ctrlGet, strings.Trim(val, " "))
	}
	sort.Strings(ctrlGet)
	if len(ctrlGet) == 1 && ctrlGet[0] == "" {
		d.Set("ctrl", make([]string, 0, 1))
	} else {
		d.Set("ctrl", ctrlGet)
	}
	d.Set("name_alias", bgpPeerPMap["nameAlias"])
	peerCtrlGet := make([]string, 0, 1)
	for _, val := range strings.Split(bgpPeerPMap["peerCtrl"], ",") {
		peerCtrlGet = append(peerCtrlGet, strings.Trim(val, " "))
	}
	sort.Strings(peerCtrlGet)
	if len(peerCtrlGet) == 1 && peerCtrlGet[0] == "" {
		d.Set("peer_ctrl", make([]string, 0, 1))
	} else {
		d.Set("peer_ctrl", peerCtrlGet)
	}
	privateASctrlGet := make([]string, 0, 1)
	for _, val := range strings.Split(bgpPeerPMap["privateASctrl"], ",") {
		privateASctrlGet = append(privateASctrlGet, strings.Trim(val, " "))
	}
	sort.Strings(privateASctrlGet)
	if len(privateASctrlGet) == 1 && privateASctrlGet[0] == "" {
		d.Set("private_a_sctrl", make([]string, 0, 1))
	} else {
		d.Set("private_a_sctrl", privateASctrlGet)
	}
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
	var parentDn string

	if d.Get("logical_node_profile_dn").(string) != "" && d.Get("parent_dn").(string) != "" {
		return diag.FromErr(fmt.Errorf("Usage of both parent_dn and logical_node_profile_dn parameters is not supported. logical_node_profile_dn parameter will be deprecated use parent_dn instead."))
	} else if d.Get("parent_dn").(string) != "" {
		parentDn = d.Get("parent_dn").(string)
	} else if d.Get("logical_node_profile_dn").(string) != "" {
		parentDn = d.Get("logical_node_profile_dn").(string)
	} else {
		return diag.FromErr(fmt.Errorf("parent_dn is required to create a BGP Peer Connectivity Profile"))
	}

	bgpPeerPAttr := models.BgpPeerConnectivityProfileAttributes{}
	if Addr, ok := d.GetOk("addr"); ok {
		bgpPeerPAttr.Addr = Addr.(string)
	}
	if AddrTCtrl, ok := d.GetOk("addr_t_ctrl"); ok {
		addrTCtrlList := make([]string, 0, 1)
		for _, val := range AddrTCtrl.([]interface{}) {
			addrTCtrlList = append(addrTCtrlList, val.(string))
		}
		AddrTCtrl := strings.Join(addrTCtrlList, ",")
		bgpPeerPAttr.AddrTCtrl = AddrTCtrl
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
		ctrlList := make([]string, 0, 1)
		for _, val := range Ctrl.([]interface{}) {
			ctrlList = append(ctrlList, val.(string))
		}
		Ctrl := strings.Join(ctrlList, ",")
		bgpPeerPAttr.Ctrl = Ctrl
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		bgpPeerPAttr.NameAlias = NameAlias.(string)
	}
	if Password, ok := d.GetOk("password"); ok {
		bgpPeerPAttr.Password = Password.(string)
	}
	if PeerCtrl, ok := d.GetOk("peer_ctrl"); ok {
		peerCtrlList := make([]string, 0, 1)
		for _, val := range PeerCtrl.([]interface{}) {
			peerCtrlList = append(peerCtrlList, val.(string))
		}
		PeerCtrl := strings.Join(peerCtrlList, ",")
		bgpPeerPAttr.PeerCtrl = PeerCtrl
	}
	if PrivateASctrl, ok := d.GetOk("private_a_sctrl"); ok {
		privateASctrlList := make([]string, 0, 1)
		for _, val := range PrivateASctrl.([]interface{}) {
			privateASctrlList = append(privateASctrlList, val.(string))
		}
		PrivateASctrl := strings.Join(privateASctrlList, ",")
		bgpPeerPAttr.PrivateASctrl = PrivateASctrl
	}
	if Ttl, ok := d.GetOk("ttl"); ok {
		bgpPeerPAttr.Ttl = Ttl.(string)
	}
	if Weight, ok := d.GetOk("weight"); ok {
		bgpPeerPAttr.Weight = Weight.(string)
	}
	if AdminSt, ok := d.GetOk("admin_state"); ok {
		bgpPeerPAttr.AdminSt = AdminSt.(string)
	}
	bgpPeerP := models.NewBgpPeerConnectivityProfile(fmt.Sprintf("peerP-[%s]", addr), parentDn, desc, bgpPeerPAttr)

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

	if relationTobgpRsPeerToProfile, ok := d.GetOk("relation_bgp_rs_peer_to_profile"); ok {
		relationParamList := toStringList(relationTobgpRsPeerToProfile.(*schema.Set).List())
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

	if relationTobgpRsPeerPfxPol, ok := d.GetOk("relation_bgp_rs_peer_pfx_pol"); ok {
		relationParam := GetMOName(relationTobgpRsPeerPfxPol.(string))
		err = aciClient.CreateRelationbgpRsPeerPfxPolFromBgpPeerConnectivityProfile(bgpPeerP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationTobgpRsPeerToProfile, ok := d.GetOk("relation_bgp_rs_peer_to_profile"); ok {
		relationParamList := relationTobgpRsPeerToProfile.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationbgpRsPeerToProfile(bgpPeerP.DistinguishedName, bgpPeerPAttr.Annotation, paramMap["direction"].(string), paramMap["target_dn"].(string))

			if err != nil {
				return diag.FromErr(err)
			}
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
	var parentDn string

	if d.Get("logical_node_profile_dn").(string) != "" && d.Get("parent_dn").(string) != "" {
		return diag.FromErr(fmt.Errorf("Usage of both parent_dn and logical_node_profile_dn parameters is not supported. logical_node_profile_dn parameter will be deprecated use parent_dn instead."))
	} else if d.Get("parent_dn").(string) != "" {
		parentDn = d.Get("parent_dn").(string)
	} else if d.Get("logical_node_profile_dn").(string) != "" {
		parentDn = d.Get("logical_node_profile_dn").(string)
	} else {
		return diag.FromErr(fmt.Errorf("parent_dn is required to update a BGP Peer Connectivity Profile"))
	}

	bgpPeerPAttr := models.BgpPeerConnectivityProfileAttributes{}
	if Addr, ok := d.GetOk("addr"); ok {
		bgpPeerPAttr.Addr = Addr.(string)
	}
	if AddrTCtrl, ok := d.GetOk("addr_t_ctrl"); ok {
		addrTCtrlList := make([]string, 0, 1)
		for _, val := range AddrTCtrl.([]interface{}) {
			addrTCtrlList = append(addrTCtrlList, val.(string))
		}
		AddrTCtrl := strings.Join(addrTCtrlList, ",")
		bgpPeerPAttr.AddrTCtrl = AddrTCtrl
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
		ctrlList := make([]string, 0, 1)
		for _, val := range Ctrl.([]interface{}) {
			ctrlList = append(ctrlList, val.(string))
		}
		Ctrl := strings.Join(ctrlList, ",")
		bgpPeerPAttr.Ctrl = Ctrl
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		bgpPeerPAttr.NameAlias = NameAlias.(string)
	}
	if PeerCtrl, ok := d.GetOk("peer_ctrl"); ok {
		peerCtrlList := make([]string, 0, 1)
		for _, val := range PeerCtrl.([]interface{}) {
			peerCtrlList = append(peerCtrlList, val.(string))
		}
		PeerCtrl := strings.Join(peerCtrlList, ",")
		bgpPeerPAttr.PeerCtrl = PeerCtrl
	}
	if PrivateASctrl, ok := d.GetOk("private_a_sctrl"); ok {
		privateASctrlList := make([]string, 0, 1)
		for _, val := range PrivateASctrl.([]interface{}) {
			privateASctrlList = append(privateASctrlList, val.(string))
		}
		PrivateASctrl := strings.Join(privateASctrlList, ",")
		bgpPeerPAttr.PrivateASctrl = PrivateASctrl
	}
	if Ttl, ok := d.GetOk("ttl"); ok {
		bgpPeerPAttr.Ttl = Ttl.(string)
	}
	if Weight, ok := d.GetOk("weight"); ok {
		bgpPeerPAttr.Weight = Weight.(string)
	}
	if AdminSt, ok := d.GetOk("admin_state"); ok {
		bgpPeerPAttr.AdminSt = AdminSt.(string)
	}
	bgpPeerP := models.NewBgpPeerConnectivityProfile(fmt.Sprintf("peerP-[%s]", addr), parentDn, desc, bgpPeerPAttr)

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
		bgpLocalAsnP.Status = "modified"

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
		bgpAsP.Status = "modified"

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

	if d.HasChange("relation_bgp_rs_peer_to_profile") {
		oldRel, newRel := d.GetChange("relation_bgp_rs_peer_to_profile")
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

	if d.HasChange("relation_bgp_rs_peer_pfx_pol") {
		_, newRelParam := d.GetChange("relation_bgp_rs_peer_pfx_pol")
		err = aciClient.CreateRelationbgpRsPeerPfxPolFromBgpPeerConnectivityProfile(bgpPeerP.DistinguishedName, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if d.HasChange("relation_bgp_rs_peer_to_profile") {
		oldRel, newRel := d.GetChange("relation_bgp_rs_peer_to_profile")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.DeleteRelationbgpRsPeerToProfile(bgpPeerP.DistinguishedName, paramMap["target_dn"].(string), paramMap["direction"].(string))

			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationbgpRsPeerToProfile(bgpPeerP.DistinguishedName, bgpPeerPAttr.Annotation, paramMap["direction"].(string), paramMap["target_dn"].(string))

			if err != nil {
				return diag.FromErr(err)
			}
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
		d.Set("relation_bgp_rs_peer_pfx_pol", bgpRsPeerPfxPolData.(string))
	}

	bgpRsPeerToProfileData, err := aciClient.ReadRelationbgpRsPeerToProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation bgpRsPeerToProfile %v", err)
	} else {
		d.Set("relation_bgp_rs_peer_to_profile", bgpRsPeerToProfileData)
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
