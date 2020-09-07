package aci

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciSubnetCreate,
		Update: resourceAciSubnetUpdate,
		Read:   resourceAciSubnetRead,
		Delete: resourceAciSubnetDelete,

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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			},

			"scope": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"virtual": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

func setSubnetAttributes(fvSubnet *models.Subnet, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(fvSubnet.DistinguishedName)
	d.Set("description", fvSubnet.Description)
	if dn != fvSubnet.DistinguishedName {
		d.Set("parent_dn", "")
	}
	fvSubnetMap, _ := fvSubnet.ToMap()

	d.Set("ip", fvSubnetMap["ip"])

	d.Set("annotation", fvSubnetMap["annotation"])
	d.Set("ctrl", fvSubnetMap["ctrl"])
	d.Set("ip", fvSubnetMap["ip"])
	d.Set("name_alias", fvSubnetMap["nameAlias"])
	d.Set("preferred", fvSubnetMap["preferred"])
	d.Set("scope", fvSubnetMap["scope"])
	d.Set("virtual", fvSubnetMap["virtual"])
	return d
}

func resourceAciSubnetImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fvSubnet, err := getRemoteSubnet(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setSubnetAttributes(fvSubnet, d)

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

func resourceAciSubnetCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Subnet: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	ip := d.Get("ip").(string)

	BridgeDomainDn := d.Get("parent_dn").(string)

	err := checkForConflictingIP(aciClient, BridgeDomainDn, ip)
	if err != nil {
		return err
	}

	fvSubnetAttr := models.SubnetAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvSubnetAttr.Annotation = Annotation.(string)
	} else {
		fvSubnetAttr.Annotation = "{}"
	}
	if Ctrl, ok := d.GetOk("ctrl"); ok {
		fvSubnetAttr.Ctrl = Ctrl.(string)
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
	if Scope, ok := d.GetOk("scope"); ok {
		fvSubnetAttr.Scope = Scope.(string)
	}
	if Virtual, ok := d.GetOk("virtual"); ok {
		fvSubnetAttr.Virtual = Virtual.(string)
	}
	fvSubnet := models.NewSubnet(fmt.Sprintf("subnet-[%s]", ip), BridgeDomainDn, desc, fvSubnetAttr)

	err = aciClient.Save(fvSubnet)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("ip")

	d.Partial(false)

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
		return err
	}
	d.Partial(false)

	if relationTofvRsBDSubnetToOut, ok := d.GetOk("relation_fv_rs_bd_subnet_to_out"); ok {
		relationParamList := toStringList(relationTofvRsBDSubnetToOut.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsBDSubnetToOutFromSubnet(fvSubnet.DistinguishedName, relationParamName)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_bd_subnet_to_out")
			d.Partial(false)
		}
	}
	if relationTofvRsNdPfxPol, ok := d.GetOk("relation_fv_rs_nd_pfx_pol"); ok {
		relationParam := relationTofvRsNdPfxPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsNdPfxPolFromSubnet(fvSubnet.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_nd_pfx_pol")
		d.Partial(false)

	}
	if relationTofvRsBDSubnetToProfile, ok := d.GetOk("relation_fv_rs_bd_subnet_to_profile"); ok {
		relationParam := relationTofvRsBDSubnetToProfile.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsBDSubnetToProfileFromSubnet(fvSubnet.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_bd_subnet_to_profile")
		d.Partial(false)

	}

	d.SetId(fvSubnet.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciSubnetRead(d, m)
}

func resourceAciSubnetUpdate(d *schema.ResourceData, m interface{}) error {
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
	if Ctrl, ok := d.GetOk("ctrl"); ok {
		fvSubnetAttr.Ctrl = Ctrl.(string)
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
	if Scope, ok := d.GetOk("scope"); ok {
		fvSubnetAttr.Scope = Scope.(string)
	}
	if Virtual, ok := d.GetOk("virtual"); ok {
		fvSubnetAttr.Virtual = Virtual.(string)
	}
	fvSubnet := models.NewSubnet(fmt.Sprintf("subnet-[%s]", ip), BridgeDomainDn, desc, fvSubnetAttr)

	fvSubnet.Status = "modified"

	err := aciClient.Save(fvSubnet)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("ip")

	d.Partial(false)

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
		return err
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
				return err
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationfvRsBDSubnetToOutFromSubnet(fvSubnet.DistinguishedName, relDnName)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_bd_subnet_to_out")
			d.Partial(false)

		}

	}
	if d.HasChange("relation_fv_rs_nd_pfx_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_nd_pfx_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationfvRsNdPfxPolFromSubnet(fvSubnet.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationfvRsNdPfxPolFromSubnet(fvSubnet.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_nd_pfx_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_fv_rs_bd_subnet_to_profile") {
		_, newRelParam := d.GetChange("relation_fv_rs_bd_subnet_to_profile")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationfvRsBDSubnetToProfileFromSubnet(fvSubnet.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationfvRsBDSubnetToProfileFromSubnet(fvSubnet.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_bd_subnet_to_profile")
		d.Partial(false)

	}

	d.SetId(fvSubnet.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciSubnetRead(d, m)

}

func resourceAciSubnetRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fvSubnet, err := getRemoteSubnet(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setSubnetAttributes(fvSubnet, d)

	fvRsBDSubnetToOutData, err := aciClient.ReadRelationfvRsBDSubnetToOutFromSubnet(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsBDSubnetToOut %v", err)
		d.Set("relation_fv_rs_bd_subnet_to_out", make([]string, 0, 1))

	} else {
		if _, ok := d.GetOk("relation_fv_rs_bd_subnet_to_out"); ok {
			relationParamList := toStringList(d.Get("relation_fv_rs_bd_subnet_to_out").(*schema.Set).List())
			tfList := make([]string, 0, 1)
			for _, relationParam := range relationParamList {
				relationParamName := GetMOName(relationParam)
				tfList = append(tfList, relationParamName)
			}
			fvRsBDSubnetToOutDataList := toStringList(fvRsBDSubnetToOutData.(*schema.Set).List())
			sort.Strings(tfList)
			sort.Strings(fvRsBDSubnetToOutDataList)

			if !reflect.DeepEqual(tfList, fvRsBDSubnetToOutDataList) {
				d.Set("relation_fv_rs_bd_subnet_to_out", make([]string, 0, 1))
			}
		}
	}

	fvRsNdPfxPolData, err := aciClient.ReadRelationfvRsNdPfxPolFromSubnet(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsNdPfxPol %v", err)
		d.Set("relation_fv_rs_nd_pfx_pol", "")

	} else {
		if _, ok := d.GetOk("relation_fv_rs_nd_pfx_pol"); ok {
			tfName := GetMOName(d.Get("relation_fv_rs_nd_pfx_pol").(string))
			if tfName != fvRsNdPfxPolData {
				d.Set("relation_fv_rs_nd_pfx_pol", "")
			}
		}
	}

	fvRsBDSubnetToProfileData, err := aciClient.ReadRelationfvRsBDSubnetToProfileFromSubnet(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsBDSubnetToProfile %v", err)
		d.Set("relation_fv_rs_bd_subnet_to_profile", "")

	} else {
		if _, ok := d.GetOk("relation_fv_rs_bd_subnet_to_profile"); ok {
			tfName := GetMOName(d.Get("relation_fv_rs_bd_subnet_to_profile").(string))
			if tfName != fvRsBDSubnetToProfileData {
				d.Set("relation_fv_rs_bd_subnet_to_profile", "")
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciSubnetDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fvSubnet")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
