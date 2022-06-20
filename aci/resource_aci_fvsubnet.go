package aci

import (
	"context"
	"fmt"
	"log"
	"net"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciSubnet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciSubnetCreate,
		UpdateContext: resourceAciSubnetUpdate,
		ReadContext:   resourceAciSubnetRead,
		DeleteContext: resourceAciSubnetDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciSubnetImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"parent_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				StateFunc: func(val interface{}) string {
					splitVal := strings.Split(val.(string), "/")
					if len(splitVal) <= 1 {
						ip := net.ParseIP(val.(string))
						return ip.String()
					} else {
						ip := net.ParseIP(splitVal[0])
						return ip.String() + "/" + splitVal[1]
					}
				},
			},

			"ctrl": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"unspecified",
						"querier",
						"nd",
						"no-default-gateway",
					}, false),
				},
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"preferred": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"yes",
					"no",
				}, false),
			},

			"scope": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"public",
						"private",
						"shared",
					}, false),
				},
			},

			"virtual": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"yes",
					"no",
				}, false),
			},

			"relation_fv_rs_bd_subnet_to_out": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_fv_rs_nd_pfx_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_bd_subnet_to_profile": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			// EP Reachability
			"next_hop_addr": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"msnlb", "anycast_mac"},
				Computed:      true,
			},
			// MSNLB
			"msnlb": {
				Optional: true,
				Type:     schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:      true,
				ConflictsWith: []string{"next_hop_addr", "anycast_mac"},
			},
			// Anycast MAC
			"anycast_mac": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"msnlb", "next_hop_addr"},
				Computed:      true,
			},
		}),
	}
}

func getAndSetNexthopEpPReachability(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	NexthopEpPReachabilityList, err := client.ListNexthopEpPReachability(dn)
	if err == nil {
		d.Set("next_hop_addr", NexthopEpPReachabilityList[0].NhAddr)
	} else {
		d.Set("next_hop_addr", nil)
	}
	return d, nil
}

func getRemoteNlbEndpoint(client *client.Client, dn string) (*models.NlbEndpoint, error) {
	fvEpNlbCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	fvEpNlb := models.NlbEndpointFromContainer(fvEpNlbCont)
	if fvEpNlb.DistinguishedName == "" {
		return nil, fmt.Errorf("NlbEndpoint %s not found", dn)
	}
	return fvEpNlb, nil
}

func setNlbEndpointAttributes(fvEpNlb *models.NlbEndpoint, d *schema.ResourceData) (*schema.ResourceData, error) {
	fvEpNlbMap, err := fvEpNlb.ToMap()
	if err != nil {
		return d, err
	}
	newContent := make(map[string]interface{})
	newContent["mode"] = fvEpNlbMap["mode"]
	newContent["group"] = fvEpNlbMap["group"]
	newContent["mac"] = fvEpNlbMap["mac"]
	d.Set("msnlb", newContent)
	return d, nil
}

func getAndSetAnycastMac(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	AnycastMacList, err := client.ListAnycastEndpoint(dn)
	if err == nil {
		d.Set("anycast_mac", AnycastMacList[0].Mac)
	} else {
		d.Set("anycast_mac", nil)
	}
	return d, nil
}

func getRemoteSubnet(client *client.Client, dn string) (*models.Subnet, error) {
	fvSubnetCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fvSubnet := models.SubnetFromContainer(fvSubnetCont)

	if fvSubnet.DistinguishedName == "" {
		return nil, fmt.Errorf("Subnet %s not found", fvSubnet.DistinguishedName)
	}

	return fvSubnet, nil
}

func setSubnetAttributes(fvSubnet *models.Subnet, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(fvSubnet.DistinguishedName)
	d.Set("description", fvSubnet.Description)
	fvSubnetMap, err := fvSubnet.ToMap()
	if err != nil {
		return d, err
	}

	if dn != fvSubnet.DistinguishedName {
		d.Set("parent_dn", "")
	} else {
		d.Set("parent_dn", GetParentDn(dn, fmt.Sprintf("/subnet-[%s]", fvSubnetMap["ip"])))
	}

	d.Set("ip", fvSubnetMap["ip"])
	d.Set("annotation", fvSubnetMap["annotation"])
	d.Set("name_alias", fvSubnetMap["nameAlias"])
	d.Set("preferred", fvSubnetMap["preferred"])

	ctrlGet := make([]string, 0, 1)
	for _, val := range strings.Split(fvSubnetMap["ctrl"], ",") {
		if val == "" {
			ctrlGet = append(ctrlGet, "unspecified")
		} else {
			ctrlGet = append(ctrlGet, strings.Trim(val, " "))
		}
	}
	sort.Strings(ctrlGet)
	if ctrlInp, ok := d.GetOk("ctrl"); ok {
		ctrlAct := make([]string, 0, 1)
		for _, val := range ctrlInp.([]interface{}) {
			ctrlAct = append(ctrlAct, val.(string))
		}
		sort.Strings(ctrlAct)
		if reflect.DeepEqual(ctrlAct, ctrlGet) {
			d.Set("ctrl", d.Get("ctrl").([]interface{}))
		} else {
			d.Set("ctrl", ctrlGet)
		}
	} else {
		d.Set("ctrl", ctrlGet)
	}

	scopeGet := make([]string, 0, 1)
	for _, val := range strings.Split(fvSubnetMap["scope"], ",") {
		scopeGet = append(scopeGet, strings.Trim(val, " "))
	}
	sort.Strings(scopeGet)
	if scopeIntr, ok := d.GetOk("scope"); ok {
		scopeAct := make([]string, 0, 1)
		for _, val := range scopeIntr.([]interface{}) {
			scopeAct = append(scopeAct, val.(string))
		}
		sort.Strings(scopeAct)
		if reflect.DeepEqual(scopeAct, scopeGet) {
			d.Set("scope", d.Get("scope").([]interface{}))
		} else {
			d.Set("scope", scopeGet)
		}
	} else {
		d.Set("scope", scopeGet)
	}
	d.Set("virtual", fvSubnetMap["virtual"])
	return d, nil
}

func resourceAciSubnetImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fvSubnet, err := getRemoteSubnet(aciClient, dn)

	if err != nil {
		return nil, err
	}
	fvSubnetMap, err := fvSubnet.ToMap()
	if err != nil {
		return nil, err
	}
	ip := fvSubnetMap["ip"]
	pDN := GetParentDn(dn, fmt.Sprintf("/subnet-[%s]", ip))
	d.Set("parent_dn", pDN)
	schemaFilled, err := setSubnetAttributes(fvSubnet, d)
	if err != nil {
		return nil, err
	}

	// ipNexthopEpP - Beginning Import
	ipNexthopEpPParentDn := fvSubnet.DistinguishedName + "/epReach"
	log.Printf("[DEBUG] %s: ipNexthopEpP - Beginning Import with parent DN", ipNexthopEpPParentDn)
	_, err = getAndSetNexthopEpPReachability(aciClient, ipNexthopEpPParentDn, d)
	if err == nil {
		ipNexthopEpPDn := dn + "/epReach/" + fmt.Sprintf(models.RnipNexthopEpP, d.Get("next_hop_addr"))
		log.Printf("[DEBUG] %s: ipNexthopEpP - Import finished successfully", ipNexthopEpPDn)
	} else {
		log.Printf("[DEBUG] %s: ipNexthopEpP - Object not present in the parent", ipNexthopEpPParentDn)
	}
	// ipNexthopEpP - Import finished successfully

	// fvEpNlb - Beginning Import
	fvEpNlbDn := fvSubnet.DistinguishedName + "/" + models.RnfvEpNlb
	fvEpNlb, err := getRemoteNlbEndpoint(aciClient, fvEpNlbDn)
	if err == nil {
		log.Printf("[DEBUG] %s: fvEpNlb - Beginning Import", fvEpNlbDn)
		_, err = setNlbEndpointAttributes(fvEpNlb, d)
		if err != nil {
			return nil, err
		}
		log.Printf("[DEBUG] %s: fvEpNlb - Import finished successfully", fvEpNlbDn)
	}
	// fvEpNlb - Import finished successfully

	// fvEpAnycast - Beginning Import
	fvEpAnycastParentDn := fvSubnet.DistinguishedName
	log.Printf("[DEBUG] %s: fvEpAnycast - Beginning import with parent DN", fvEpAnycastParentDn)
	_, err = getAndSetAnycastMac(aciClient, fvEpAnycastParentDn, d)
	if err == nil {
		fvEpAnycastDn := dn + "/" + fmt.Sprintf(models.RnfvEpAnycast, d.Get("anycast_mac"))
		log.Printf("[DEBUG] %s: fvEpAnycast - Import finished successfully", fvEpAnycastDn)
	} else {
		log.Printf("[DEBUG] %s: fvEpAnycast - Import not present in the parent", dn)
	}
	// fvEpAnycast - Import finished successfully

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func checkForConflictingVRF(client *client.Client, tenantDN, bdName, vrfDn, ip string) bool {
	flag := false

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/BD-%s/fvRsCtx.json", baseurlStr, tenantDN, bdName)
	fvCtxCont, err := client.GetViaURL(dnUrl)
	if err != nil {
		return flag
	}

	fvCtxContList := models.ListFromContainer(fvCtxCont, "fvRsCtx")
	if len(fvCtxContList) > 0 {
		if vrfDn != models.G(fvCtxContList[0], "tDn") {
			return flag
		}

		dnUrl = fmt.Sprintf("%s/BD-%s/subnet-[%s]", tenantDN, bdName, ip)
		_, err = client.Get(dnUrl)
		if err == nil {
			flag = true
			return flag
		}
		return flag
	}

	return flag
}

func checkForConflictingIP(client *client.Client, parentDN string, ip string) error {
	tokens := strings.Split(parentDN, "/")
	bdName := (strings.Split(tokens[2], "-"))[1]
	tenantDn := fmt.Sprintf("%s/%s", tokens[0], tokens[1])

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, tenantDn, "fvBD")

	domains, err := client.GetViaURL(dnUrl)
	if err != nil {
		return err
	}
	bdList := models.ListFromContainer(domains, "fvBD")

	dnUrl = fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDN, "fvRsCtx")
	fvCtxCont, err := client.GetViaURL(dnUrl)
	if err != nil {
		return nil
	}
	fvCtxContList := models.ListFromContainer(fvCtxCont, "fvRsCtx")
	var ctxDN string
	if len(fvCtxContList) > 0 {
		ctxDN = models.G(fvCtxContList[0], "tDn")
	} else {
		return nil
	}

	if len(bdList) > 1 {
		for i := 0; i < (len(bdList)); i++ {
			currName := models.G(bdList[i], "name")
			if currName != bdName {
				if checkForConflictingVRF(client, tenantDn, currName, ctxDN, ip) {
					return fmt.Errorf("A subnet already exist with Bridge Domain %s and ip %s", currName, ip)
				}
			}
		}
	}
	return nil
}

func resourceAciSubnetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Subnet: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	ip := d.Get("ip").(string)

	BridgeDomainDn := d.Get("parent_dn").(string)

	ap_epg_subnet_pattern, _ := regexp.Compile("^uni/tn-(.+)/ap-(.+)/epg-(.+)/subnet-(.+)")

	if ap_epg_subnet_pattern.Match([]byte(BridgeDomainDn)) {
		err := checkForConflictingIP(aciClient, BridgeDomainDn, ip)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	fvSubnetAttr := models.SubnetAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvSubnetAttr.Annotation = Annotation.(string)
	} else {
		fvSubnetAttr.Annotation = "{}"
	}
	if ctrlInp, ok := d.GetOk("ctrl"); ok {
		ctrlList := make([]string, 0, 1)
		for _, val := range ctrlInp.([]interface{}) {
			ctrlList = append(ctrlList, val.(string))
		}
		ctrl := strings.Join(ctrlList, ",")
		fvSubnetAttr.Ctrl = ctrl
	}
	if Ip, ok := d.GetOk("ip"); ok {
		fvSubnetAttr.Ip = Ip.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvSubnetAttr.NameAlias = NameAlias.(string)
	}
	if Preferred, ok := d.GetOk("preferred"); ok {
		fvSubnetAttr.Preferred = Preferred.(string)
	}
	if scIntr, ok := d.GetOk("scope"); ok {
		scopeList := make([]string, 0, 1)
		for _, val := range scIntr.([]interface{}) {
			scopeList = append(scopeList, val.(string))
		}
		Scope := strings.Join(scopeList, ",")
		fvSubnetAttr.Scope = Scope
	}
	if Virtual, ok := d.GetOk("virtual"); ok {
		fvSubnetAttr.Virtual = Virtual.(string)
	}
	fvSubnet := models.NewSubnet(fmt.Sprintf("subnet-[%s]", ip), BridgeDomainDn, desc, fvSubnetAttr)

	err := aciClient.Save(fvSubnet)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTofvRsBDSubnetToOut, ok := d.GetOk("relation_fv_rs_bd_subnet_to_out"); ok {
		relationParamList := toStringList(relationTofvRsBDSubnetToOut.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTofvRsNdPfxPol, ok := d.GetOk("relation_fv_rs_nd_pfx_pol"); ok {
		relationParam := relationTofvRsNdPfxPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTofvRsBDSubnetToProfile, ok := d.GetOk("relation_fv_rs_bd_subnet_to_profile"); ok {
		relationParam := relationTofvRsBDSubnetToProfile.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTofvRsBDSubnetToOut, ok := d.GetOk("relation_fv_rs_bd_subnet_to_out"); ok {
		relationParamList := toStringList(relationTofvRsBDSubnetToOut.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsBDSubnetToOutFromSubnet(fvSubnet.DistinguishedName, relationParamName)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if relationTofvRsNdPfxPol, ok := d.GetOk("relation_fv_rs_nd_pfx_pol"); ok {
		relationParam := relationTofvRsNdPfxPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsNdPfxPolFromSubnet(fvSubnet.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationTofvRsBDSubnetToProfile, ok := d.GetOk("relation_fv_rs_bd_subnet_to_profile"); ok {
		relationParam := relationTofvRsBDSubnetToProfile.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsBDSubnetToProfileFromSubnet(fvSubnet.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// fvEpReachability - Create
	if next_hop_addr, ok := d.GetOk("next_hop_addr"); ok {
		log.Printf("[DEBUG] fvEpReachability: Beginning Creation")

		fvEpReachabilityAttr := models.EpReachabilityAttributes{}
		fvEpReachability := models.NewEpReachability(fmt.Sprintf("epReach"), fvSubnet.DistinguishedName, fvEpReachabilityAttr)

		ep_reachability_create_err := aciClient.Save(fvEpReachability)
		if ep_reachability_create_err != nil {
			return diag.FromErr(ep_reachability_create_err)
		}

		log.Printf("[DEBUG] %s: Creation finished successfully", fvEpReachability.DistinguishedName)

		log.Printf("[DEBUG] ipNexthopEpP: Beginning Creation")

		ipNexthopEpPAttr := models.NexthopEpPReachabilityAttributes{}

		ipNexthopEpPAttr.NhAddr = next_hop_addr.(string)

		ipNexthopEpP := models.NewNexthopEpPReachability(fmt.Sprintf(models.RnipNexthopEpP, next_hop_addr), fvEpReachability.DistinguishedName, "", "", ipNexthopEpPAttr)

		next_hop_addr_create_err := aciClient.Save(ipNexthopEpP)
		if next_hop_addr_create_err != nil {
			return diag.FromErr(next_hop_addr_create_err)
		}

		log.Printf("[DEBUG] %s: Creation finished successfully", ipNexthopEpP.DistinguishedName)

	}

	// fvEpNlb - Create
	if msnlb, ok := d.GetOk("msnlb"); ok {
		log.Printf("[DEBUG] fvEpNlb: Beginning Creation")

		fvEpNlbAttr := models.NlbEndpointAttributes{}

		msnlbMap := toStrMap(msnlb.(map[string]interface{}))

		fvEpNlbAttr.Mode = msnlbMap["mode"]

		if msnlbMap["mode"] == "mode-mcast--static" || msnlbMap["mode"] == "mode-uc" {
			fvEpNlbAttr.Mac = msnlbMap["mac"]
			_, flag := msnlbMap["group"]
			if flag {
				if msnlbMap["group"] != "0.0.0.0" && msnlbMap["group"] != "" {
					return diag.FromErr(fmt.Errorf("Invalid configuration, \"group\" must be \"0.0.0.0\" or empty(\"\") string when the mode is other than \"mode-mcast-igmp\""))
				}
			} else {
				return diag.FromErr(fmt.Errorf("Invalid configuration, \"group\" must be \"0.0.0.0\" or empty(\"\") string when the mode is other than \"mode-mcast-igmp\""))
			}
		}

		if msnlbMap["mode"] == "mode-mcast-igmp" {
			fvEpNlbAttr.Group = msnlbMap["group"]
			_, flag := msnlbMap["mac"]
			if flag {
				if msnlbMap["mac"] != "00:00:00:00:00:00" && msnlbMap["mac"] != "" {
					return diag.FromErr(fmt.Errorf("Invalid configuration, \"mac\" must be \"00:00:00:00:00:00\" or empty(\"\") string when the mode is \"mode-mcast-igmp\""))
				}
			} else {
				return diag.FromErr(fmt.Errorf("Invalid configuration, \"mac\" must be \"00:00:00:00:00:00\" or empty(\"\") string when the mode is \"mode-mcast-igmp\""))
			}
		}

		fvEpNlb := models.NewNlbEndpoint(fmt.Sprintf(models.RnfvEpNlb), fvSubnet.DistinguishedName, "", "", fvEpNlbAttr)

		msnlb_create_err := aciClient.Save(fvEpNlb)
		if msnlb_create_err != nil {
			return diag.FromErr(msnlb_create_err)
		}

		log.Printf("[DEBUG] %s: Creation finished successfully", fvEpNlb.DistinguishedName)
	}

	// fvEpAnycast - Create
	if anycastMac, ok := d.GetOk("anycast_mac"); ok {
		log.Printf("[DEBUG] fvEpAnycast: Beginning Creation")

		fvEpAnycastAttr := models.AnycastEndpointAttributes{}
		fvEpAnycastAttr.Mac = anycastMac.(string)
		fvEpAnycast := models.NewAnycastEndpoint(fmt.Sprintf(models.RnfvEpAnycast, anycastMac), fvSubnet.DistinguishedName, "", "", fvEpAnycastAttr)

		anycast_create_err := aciClient.Save(fvEpAnycast)
		if anycast_create_err != nil {
			return diag.FromErr(anycast_create_err)
		}

		log.Printf("[DEBUG] %s: Creation finished successfully", fvEpAnycast.DistinguishedName)
	}

	d.SetId(fvSubnet.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciSubnetRead(ctx, d, m)
}

func resourceAciSubnetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Subnet: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	ip := d.Get("ip").(string)

	BridgeDomainDn := d.Get("parent_dn").(string)

	fvSubnetAttr := models.SubnetAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvSubnetAttr.Annotation = Annotation.(string)
	} else {
		fvSubnetAttr.Annotation = "{}"
	}
	if ctrlInp, ok := d.GetOk("ctrl"); ok {
		ctrlList := make([]string, 0, 1)
		for _, val := range ctrlInp.([]interface{}) {
			ctrlList = append(ctrlList, val.(string))
		}
		ctrl := strings.Join(ctrlList, ",")
		fvSubnetAttr.Ctrl = ctrl
	}
	if Ip, ok := d.GetOk("ip"); ok {
		fvSubnetAttr.Ip = Ip.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvSubnetAttr.NameAlias = NameAlias.(string)
	}
	if Preferred, ok := d.GetOk("preferred"); ok {
		fvSubnetAttr.Preferred = Preferred.(string)
	}
	if scIntr, ok := d.GetOk("scope"); ok {
		scopeList := make([]string, 0, 1)
		for _, val := range scIntr.([]interface{}) {
			scopeList = append(scopeList, val.(string))
		}
		Scope := strings.Join(scopeList, ",")
		fvSubnetAttr.Scope = Scope
	}
	if Virtual, ok := d.GetOk("virtual"); ok {
		fvSubnetAttr.Virtual = Virtual.(string)
	}
	fvSubnet := models.NewSubnet(fmt.Sprintf("subnet-[%s]", ip), BridgeDomainDn, desc, fvSubnetAttr)

	fvSubnet.Status = "modified"

	err := aciClient.Save(fvSubnet)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_fv_rs_bd_subnet_to_out") {
		oldRel, newRel := d.GetChange("relation_fv_rs_bd_subnet_to_out")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fv_rs_nd_pfx_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_nd_pfx_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_fv_rs_bd_subnet_to_profile") {
		_, newRelParam := d.GetChange("relation_fv_rs_bd_subnet_to_profile")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_fv_rs_bd_subnet_to_out") {
		oldRel, newRel := d.GetChange("relation_fv_rs_bd_subnet_to_out")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			relDnName := GetMOName(relDn)
			err = aciClient.DeleteRelationfvRsBDSubnetToOutFromSubnet(fvSubnet.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationfvRsBDSubnetToOutFromSubnet(fvSubnet.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("relation_fv_rs_nd_pfx_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_nd_pfx_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationfvRsNdPfxPolFromSubnet(fvSubnet.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationfvRsNdPfxPolFromSubnet(fvSubnet.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_fv_rs_bd_subnet_to_profile") {
		_, newRelParam := d.GetChange("relation_fv_rs_bd_subnet_to_profile")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationfvRsBDSubnetToProfileFromSubnet(fvSubnet.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationfvRsBDSubnetToProfileFromSubnet(fvSubnet.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// fvEpReachability - Update
	if d.HasChange("next_hop_addr") {
		if next_hop_addr, ok := d.GetOk("next_hop_addr"); ok {
			log.Printf("[DEBUG] fvEpReachability - Beginning Update")

			fvEpReachabilityAttr := models.EpReachabilityAttributes{}

			nextHopAddrDn := fvSubnet.DistinguishedName + fmt.Sprintf("/"+models.RnfvEpReachability)

			ep_reachability_delete_error := aciClient.DeleteByDn(nextHopAddrDn, "fvEpReachability")
			if ep_reachability_delete_error != nil {
				return diag.FromErr(ep_reachability_delete_error)
			}

			fvEpReachability := models.NewEpReachability(fmt.Sprintf("epReach"), fvSubnet.DistinguishedName, fvEpReachabilityAttr)

			ep_reachability_update_err := aciClient.Save(fvEpReachability)
			if ep_reachability_update_err != nil {
				return diag.FromErr(ep_reachability_update_err)
			}

			log.Printf("[DEBUG] ipNexthopEpP: Beginning Update")

			ipNexthopEpPAttr := models.NexthopEpPReachabilityAttributes{}

			ipNexthopEpPAttr.NhAddr = next_hop_addr.(string)

			ipNexthopEpP := models.NewNexthopEpPReachability(fmt.Sprintf(models.RnipNexthopEpP, next_hop_addr), fvEpReachability.DistinguishedName, "", "", ipNexthopEpPAttr)

			next_hop_addr_update_err := aciClient.Save(ipNexthopEpP)
			if next_hop_addr_update_err != nil {
				return diag.FromErr(next_hop_addr_update_err)
			}

			log.Printf("[DEBUG] %s: ipNexthopEpP - Update finished successfully", ipNexthopEpP.DistinguishedName)

			log.Printf("[DEBUG] %s: fvEpReachability - Update finished successfully", fvEpReachability.DistinguishedName)

		} else {
			nextHopAddrDn := fvSubnet.DistinguishedName + fmt.Sprintf("/"+models.RnfvEpReachability)
			log.Printf("[DEBUG] %s: fvEpReachability - Beginning Destroy", nextHopAddrDn)
			ep_reachability_delete_error := aciClient.DeleteByDn(nextHopAddrDn, "fvEpReachability")
			if ep_reachability_delete_error != nil {
				return diag.FromErr(ep_reachability_delete_error)
			}
			log.Printf("[DEBUG] %s: fvEpReachability - Destroy finished successfully", nextHopAddrDn)
		}
	}

	// fvEpNlb - Update
	if d.HasChange("msnlb") {
		if msnlb, ok := d.GetOk("msnlb"); ok {
			log.Printf("[DEBUG] fvEpNlb - Beginning Update")

			fvEpNlbAttr := models.NlbEndpointAttributes{}
			msnlbMap := toStrMap(msnlb.(map[string]interface{}))
			fvEpNlbAttr.Mode = msnlbMap["mode"]

			if msnlbMap["mode"] == "mode-mcast--static" || msnlbMap["mode"] == "mode-uc" {
				fvEpNlbAttr.Mac = msnlbMap["mac"]
				_, flag := msnlbMap["group"]
				if flag {
					if msnlbMap["group"] != "0.0.0.0" && msnlbMap["group"] != "" {
						return diag.FromErr(fmt.Errorf("Invalid configuration, \"group\" must be \"0.0.0.0\" or empty(\"\") string when the mode is other than \"mode-mcast-igmp\""))
					}
				} else {
					return diag.FromErr(fmt.Errorf("Invalid configuration, \"group\" must be \"0.0.0.0\" or empty(\"\") string when the mode is other than \"mode-mcast-igmp\""))
				}
			}

			if msnlbMap["mode"] == "mode-mcast-igmp" {
				fvEpNlbAttr.Group = msnlbMap["group"]
				_, flag := msnlbMap["mac"]
				if flag {
					if msnlbMap["mac"] != "00:00:00:00:00:00" && msnlbMap["mac"] != "" {
						return diag.FromErr(fmt.Errorf("Invalid configuration, \"mac\" must be \"00:00:00:00:00:00\" or empty(\"\") string when the mode is \"mode-mcast-igmp\""))
					}
				} else {
					return diag.FromErr(fmt.Errorf("Invalid configuration, \"mac\" must be \"00:00:00:00:00:00\" or empty(\"\") string when the mode is \"mode-mcast-igmp\""))
				}
			}

			msnlbDn := fvSubnet.DistinguishedName + fmt.Sprintf("/"+models.RnfvEpNlb)
			deletion_err := aciClient.DeleteByDn(msnlbDn, "fvEpNlb")
			if deletion_err != nil {
				return diag.FromErr(err)
			}

			fvEpNlb := models.NewNlbEndpoint(fmt.Sprintf(models.RnfvEpNlb), fvSubnet.DistinguishedName, "", "", fvEpNlbAttr)
			msnlb_update_err := aciClient.Save(fvEpNlb)
			if msnlb_update_err != nil {
				return diag.FromErr(msnlb_update_err)
			}

			log.Printf("[DEBUG] %s: fvEpNlb - Update finished successfully", fvEpNlb.DistinguishedName)
		} else {
			msnlbDn := fvSubnet.DistinguishedName + fmt.Sprintf("/"+models.RnfvEpNlb)
			log.Printf("[DEBUG] %s: fvEpNlb - Beginning Destroy", msnlbDn)
			err := aciClient.DeleteByDn(msnlbDn, "fvEpNlb")
			if err != nil {
				return diag.FromErr(err)
			}
			log.Printf("[DEBUG] %s: fvEpNlb - Destroy finished successfully", msnlbDn)
		}
	}

	// fvEpAnycast - Update
	if d.HasChange("anycast_mac") {
		if anycastMac, ok := d.GetOk("anycast_mac"); ok {
			log.Printf("[DEBUG] fvEpAnycast - Beginning Update")

			fvEpAnycastAttr := models.AnycastEndpointAttributes{}
			fvEpAnycastAttr.Mac = anycastMac.(string)

			AnycastMacList, err := aciClient.ListAnycastEndpoint(fvSubnet.DistinguishedName)

			if err != nil {
				return diag.FromErr(err)
			}

			anycastMacDn := fvSubnet.DistinguishedName + "/" + fmt.Sprintf(models.RnfvEpAnycast, AnycastMacList[0].Mac)

			anycastMac_delete_err := aciClient.DeleteByDn(anycastMacDn, "fvEpAnycast")
			if anycastMac_delete_err != nil {
				return diag.FromErr(err)
			}

			fvEpAnycast := models.NewAnycastEndpoint(fmt.Sprintf(models.RnfvEpAnycast, anycastMac), fvSubnet.DistinguishedName, "", "", fvEpAnycastAttr)

			anycast_update_err := aciClient.Save(fvEpAnycast)
			if anycast_update_err != nil {
				return diag.FromErr(anycast_update_err)
			}

			log.Printf("[DEBUG] %s: fvEpAnycast - Update finished successfully", fvEpAnycast.DistinguishedName)
		} else {
			AnycastMacList, err := aciClient.ListAnycastEndpoint(fvSubnet.DistinguishedName)
			if err != nil {
				return diag.FromErr(err)
			}
			anycastMacDn := fvSubnet.DistinguishedName + "/" + fmt.Sprintf(models.RnfvEpAnycast, AnycastMacList[0].Mac)
			log.Printf("[DEBUG] %s: fvEpAnycast - Beginning Destroy", anycastMacDn)
			anycastMac_delete_err := aciClient.DeleteByDn(anycastMacDn, "fvEpAnycast")
			if anycastMac_delete_err != nil {
				return diag.FromErr(anycastMac_delete_err)
			}
			log.Printf("[DEBUG] %s: fvEpAnycast - Destroy finished successfully", anycastMacDn)
		}
	}

	d.SetId(fvSubnet.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciSubnetRead(ctx, d, m)

}

func resourceAciSubnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fvSubnet, err := getRemoteSubnet(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setSubnetAttributes(fvSubnet, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	fvRsBDSubnetToOutData, err := aciClient.ReadRelationfvRsBDSubnetToOutFromSubnet(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsBDSubnetToOut %v", err)
		setRelationAttribute(d, "relation_fv_rs_bd_subnet_to_out", make([]interface{}, 0, 1))
	} else {
		setRelationAttribute(d, "relation_fv_rs_bd_subnet_to_out", toStringList(fvRsBDSubnetToOutData.(*schema.Set).List()))
	}

	fvRsNdPfxPolData, err := aciClient.ReadRelationfvRsNdPfxPolFromSubnet(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsNdPfxPol %v", err)
		d.Set("relation_fv_rs_nd_pfx_pol", "")

	} else {
		setRelationAttribute(d, "relation_fv_rs_nd_pfx_pol", fvRsNdPfxPolData.(string))
	}

	fvRsBDSubnetToProfileData, err := aciClient.ReadRelationfvRsBDSubnetToProfileFromSubnet(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsBDSubnetToProfile %v", err)
		d.Set("relation_fv_rs_bd_subnet_to_profile", "")

	} else {
		setRelationAttribute(d, "relation_fv_rs_bd_subnet_to_profile", fvRsBDSubnetToProfileData.(string))
	}

	// ipNexthopEpP - Beginning Read
	ipNexthopEpPParentDn := dn + "/epReach"
	log.Printf("[DEBUG] %s: ipNexthopEpP - Beginning Read with parent DN", ipNexthopEpPParentDn)
	_, err = getAndSetNexthopEpPReachability(aciClient, ipNexthopEpPParentDn, d)
	if err == nil {
		ipNexthopEpPDn := dn + "/epReach/" + fmt.Sprintf(models.RnipNexthopEpP, d.Get("next_hop_addr"))
		log.Printf("[DEBUG] %s: ipNexthopEpP - Read finished successfully", ipNexthopEpPDn)
	} else {
		log.Printf("[DEBUG] %s: ipNexthopEpP - Object not present in the parent", ipNexthopEpPParentDn)
	}
	// ipNexthopEpP - Read finished successfully

	// fvEpNlb - Beginning Read
	fvEpNlbDn := dn + "/" + models.RnfvEpNlb
	fvEpNlb, err := getRemoteNlbEndpoint(aciClient, fvEpNlbDn)
	if err == nil {
		log.Printf("[DEBUG] %s: fvEpNlb - Beginning Read", fvEpNlbDn)
		_, err = setNlbEndpointAttributes(fvEpNlb, d)
		if err != nil {
			return nil
		}
		log.Printf("[DEBUG] %s: fvEpNlb - Read finished successfully", fvEpNlbDn)
	} else {
		d.Set("msnlb", nil)
	}
	// fvEpNlb - Read finished successfully

	// fvEpAnycast - Beginning Read
	log.Printf("[DEBUG] %s: fvEpAnycast - Beginning Read with parent DN", dn)
	_, err = getAndSetAnycastMac(aciClient, dn, d)
	if err == nil {
		fvEpAnycastDn := dn + "/" + fmt.Sprintf(models.RnfvEpAnycast, d.Get("anycast_mac"))
		log.Printf("[DEBUG] %s: fvEpAnycast - Read finished successfully", fvEpAnycastDn)
	} else {
		log.Printf("[DEBUG] %s: fvEpAnycast - Object not present in the parent", dn)
	}
	// fvEpAnycast - Read finished successfully

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciSubnetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fvSubnet")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
