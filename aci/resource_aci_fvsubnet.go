package aci

import (
	"context"
	"fmt"
	"log"
	"reflect"
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
		}),
	}
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
	if dn != fvSubnet.DistinguishedName {
		d.Set("parent_dn", "")
	}
	fvSubnetMap, err := fvSubnet.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("parent_dn", GetParentDn(dn, fmt.Sprintf("/subnet-[%s]", fvSubnetMap["ip"])))
	d.Set("ip", fvSubnetMap["ip"])

	d.Set("annotation", fvSubnetMap["annotation"])
	d.Set("ip", fvSubnetMap["ip"])
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

	err := checkForConflictingIP(aciClient, BridgeDomainDn, ip)
	if err != nil {
		return diag.FromErr(err)
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

	err = aciClient.Save(fvSubnet)
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
		d.Set("relation_fv_rs_bd_subnet_to_out", make([]string, 0, 1))

	} else {
		d.Set("relation_fv_rs_bd_subnet_to_out", toStringList(fvRsBDSubnetToOutData.(*schema.Set).List()))
	}

	fvRsNdPfxPolData, err := aciClient.ReadRelationfvRsNdPfxPolFromSubnet(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsNdPfxPol %v", err)
		d.Set("relation_fv_rs_nd_pfx_pol", "")

	} else {
		d.Set("relation_fv_rs_nd_pfx_pol", fvRsNdPfxPolData.(string))
	}

	fvRsBDSubnetToProfileData, err := aciClient.ReadRelationfvRsBDSubnetToProfileFromSubnet(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsBDSubnetToProfile %v", err)
		d.Set("relation_fv_rs_bd_subnet_to_profile", "")

	} else {
		d.Set("relation_fv_rs_bd_subnet_to_profile", fvRsBDSubnetToProfileData.(string))
	}

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
