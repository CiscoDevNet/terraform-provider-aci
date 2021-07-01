package aci

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciL3Outside() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciL3OutsideCreate,
		Update: resourceAciL3OutsideUpdate,
		Read:   resourceAciL3OutsideRead,
		Delete: resourceAciL3OutsideDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciL3OutsideImport,
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

			"enforce_rtctrl": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"export",
						"import",
					}, false),
				},
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"target_dscp": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"CS0",
					"CS1",
					"AF11",
					"AF12",
					"AF13",
					"CS2",
					"AF21",
					"AF22",
					"AF23",
					"CS3",
					"CS4",
					"CS5",
					"CS6",
					"CS7",
					"AF31",
					"AF32",
					"AF33",
					"AF41",
					"AF42",
					"AF43",
					"VA",
					"EF",
					"unspecified",
				}, false),
			},

			"relation_l3ext_rs_dampening_pol": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tn_rtctrl_profile_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"af": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"relation_l3ext_rs_ectx": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_l3ext_rs_out_to_bd_public_subnet_holder": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_l3ext_rs_interleak_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_l3ext_rs_l3_dom_att": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteL3Outside(client *client.Client, dn string) (*models.L3Outside, error) {
	l3extOutCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extOut := models.L3OutsideFromContainer(l3extOutCont)

	if l3extOut.DistinguishedName == "" {
		return nil, fmt.Errorf("L3Outside %s not found", l3extOut.DistinguishedName)
	}

	return l3extOut, nil
}

func setL3OutsideAttributes(l3extOut *models.L3Outside, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(l3extOut.DistinguishedName)
	d.Set("description", l3extOut.Description)
	// d.Set("tenant_dn", GetParentDn(l3extOut.DistinguishedName))
	if dn != l3extOut.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	l3extOutMap, _ := l3extOut.ToMap()

	d.Set("name", l3extOutMap["name"])

	d.Set("annotation", l3extOutMap["annotation"])
	enforceRtctrlGet := make([]string, 0, 1)
	for _, val := range strings.Split(l3extOutMap["enforceRtctrl"], ",") {
		enforceRtctrlGet = append(enforceRtctrlGet, strings.Trim(val, " "))
	}
	sort.Strings(enforceRtctrlGet)
	if enforceRtctrlIntr, ok := d.GetOk("enforce_rtctrl"); ok {
		enforceRtctrlAct := make([]string, 0, 1)
		for _, val := range enforceRtctrlIntr.([]interface{}) {
			enforceRtctrlAct = append(enforceRtctrlAct, val.(string))
		}
		sort.Strings(enforceRtctrlAct)
		if reflect.DeepEqual(enforceRtctrlAct, enforceRtctrlGet) {
			d.Set("enforce_rtctrl", d.Get("enforce_rtctrl").([]interface{}))
		} else {
			d.Set("enforce_rtctrl", enforceRtctrlGet)
		}
	} else {
		d.Set("enforce_rtctrl", enforceRtctrlGet)
	}
	d.Set("name_alias", l3extOutMap["nameAlias"])
	d.Set("target_dscp", l3extOutMap["targetDscp"])
	return d
}

func resourceAciL3OutsideImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	l3extOut, err := getRemoteL3Outside(aciClient, dn)

	if err != nil {
		return nil, err
	}
	l3extOutMap, _ := l3extOut.ToMap()
	name := l3extOutMap["name"]
	pDN := GetParentDn(dn, fmt.Sprintf("/out-%s", name))
	d.Set("tenant_dn", pDN)
	schemaFilled := setL3OutsideAttributes(l3extOut, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL3OutsideCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] L3Outside: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	l3extOutAttr := models.L3OutsideAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extOutAttr.Annotation = Annotation.(string)
	} else {
		l3extOutAttr.Annotation = "{}"
	}
	if EnforceRtctrl, ok := d.GetOk("enforce_rtctrl"); ok {
		enforceRtctrlList := make([]string, 0, 1)
		for _, val := range EnforceRtctrl.([]interface{}) {
			enforceRtctrlList = append(enforceRtctrlList, val.(string))
		}
		EnforceRtctrl := strings.Join(enforceRtctrlList, ",")
		l3extOutAttr.EnforceRtctrl = EnforceRtctrl
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extOutAttr.NameAlias = NameAlias.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		l3extOutAttr.TargetDscp = TargetDscp.(string)
	}
	l3extOut := models.NewL3Outside(fmt.Sprintf("out-%s", name), TenantDn, desc, l3extOutAttr)

	err := aciClient.Save(l3extOut)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.Partial(false)

	checkDns := make([]string, 0, 1)

	if relationTol3extRsEctx, ok := d.GetOk("relation_l3ext_rs_ectx"); ok {
		relationParam := relationTol3extRsEctx.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTol3extRsOutToBDPublicSubnetHolder, ok := d.GetOk("relation_l3ext_rs_out_to_bd_public_subnet_holder"); ok {
		relationParamList := toStringList(relationTol3extRsOutToBDPublicSubnetHolder.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTol3extRsInterleakPol, ok := d.GetOk("relation_l3ext_rs_interleak_pol"); ok {
		relationParam := relationTol3extRsInterleakPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTol3extRsL3DomAtt, ok := d.GetOk("relation_l3ext_rs_l3_dom_att"); ok {
		relationParam := relationTol3extRsL3DomAtt.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return err
	}
	d.Partial(false)

	if relationTol3extRsDampeningPol, ok := d.GetOk("relation_l3ext_rs_dampening_pol"); ok {

		relationParamList := relationTol3extRsDampeningPol.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationl3extRsDampeningPolFromL3Outside(l3extOut.DistinguishedName, paramMap["tn_rtctrl_profile_name"].(string), paramMap["af"].(string))
			if err != nil {
				return err
			}
			d.Partial(true)
			d.Partial(false)
		}

	}
	if relationTol3extRsEctx, ok := d.GetOk("relation_l3ext_rs_ectx"); ok {
		relationParam := relationTol3extRsEctx.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationl3extRsEctxFromL3Outside(l3extOut.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.Partial(false)

	}
	if relationTol3extRsOutToBDPublicSubnetHolder, ok := d.GetOk("relation_l3ext_rs_out_to_bd_public_subnet_holder"); ok {
		relationParamList := toStringList(relationTol3extRsOutToBDPublicSubnetHolder.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationl3extRsOutToBDPublicSubnetHolderFromL3Outside(l3extOut.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.Partial(false)
		}
	}
	if relationTol3extRsInterleakPol, ok := d.GetOk("relation_l3ext_rs_interleak_pol"); ok {
		relationParam := relationTol3extRsInterleakPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationl3extRsInterleakPolFromL3Outside(l3extOut.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.Partial(false)

	}
	if relationTol3extRsL3DomAtt, ok := d.GetOk("relation_l3ext_rs_l3_dom_att"); ok {
		relationParam := relationTol3extRsL3DomAtt.(string)
		err = aciClient.CreateRelationl3extRsL3DomAttFromL3Outside(l3extOut.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.Partial(false)

	}

	d.SetId(l3extOut.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciL3OutsideRead(d, m)
}

func resourceAciL3OutsideUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] L3Outside: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	l3extOutAttr := models.L3OutsideAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extOutAttr.Annotation = Annotation.(string)
	} else {
		l3extOutAttr.Annotation = "{}"
	}
	if EnforceRtctrl, ok := d.GetOk("enforce_rtctrl"); ok {
		enforceRtctrlList := make([]string, 0, 1)
		for _, val := range EnforceRtctrl.([]interface{}) {
			enforceRtctrlList = append(enforceRtctrlList, val.(string))
		}
		EnforceRtctrl := strings.Join(enforceRtctrlList, ",")
		l3extOutAttr.EnforceRtctrl = EnforceRtctrl
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extOutAttr.NameAlias = NameAlias.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		l3extOutAttr.TargetDscp = TargetDscp.(string)
	}
	l3extOut := models.NewL3Outside(fmt.Sprintf("out-%s", name), TenantDn, desc, l3extOutAttr)

	l3extOut.Status = "modified"

	err := aciClient.Save(l3extOut)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.Partial(false)

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_l3ext_rs_ectx") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_ectx")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_l3ext_rs_out_to_bd_public_subnet_holder") {
		oldRel, newRel := d.GetChange("relation_l3ext_rs_out_to_bd_public_subnet_holder")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_l3ext_rs_interleak_pol") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_interleak_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_l3ext_rs_l3_dom_att") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_l3_dom_att")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return err
	}
	d.Partial(false)

	if d.HasChange("relation_l3ext_rs_dampening_pol") {
		oldRel, newRel := d.GetChange("relation_l3ext_rs_dampening_pol")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.DeleteRelationl3extRsDampeningPolFromL3Outside(l3extOut.DistinguishedName, paramMap["tn_rtctrl_profile_name"].(string), paramMap["af"].(string))
			if err != nil {
				return err
			}

		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationl3extRsDampeningPolFromL3Outside(l3extOut.DistinguishedName, paramMap["tn_rtctrl_profile_name"].(string), paramMap["af"].(string))
			if err != nil {
				return err
			}
			d.Partial(true)
			d.Partial(false)
		}

	}
	if d.HasChange("relation_l3ext_rs_ectx") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_ectx")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationl3extRsEctxFromL3Outside(l3extOut.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.Partial(false)

	}
	if d.HasChange("relation_l3ext_rs_out_to_bd_public_subnet_holder") {
		oldRel, newRel := d.GetChange("relation_l3ext_rs_out_to_bd_public_subnet_holder")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationl3extRsOutToBDPublicSubnetHolderFromL3Outside(l3extOut.DistinguishedName, relDn)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.Partial(false)

		}

	}
	if d.HasChange("relation_l3ext_rs_interleak_pol") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_interleak_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationl3extRsInterleakPolFromL3Outside(l3extOut.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationl3extRsInterleakPolFromL3Outside(l3extOut.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.Partial(false)

	}
	if d.HasChange("relation_l3ext_rs_l3_dom_att") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_l3_dom_att")
		err = aciClient.DeleteRelationl3extRsL3DomAttFromL3Outside(l3extOut.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationl3extRsL3DomAttFromL3Outside(l3extOut.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.Partial(false)

	}

	d.SetId(l3extOut.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciL3OutsideRead(d, m)

}

func resourceAciL3OutsideRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l3extOut, err := getRemoteL3Outside(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setL3OutsideAttributes(l3extOut, d)

	l3extRsDampeningPolData, err := aciClient.ReadRelationl3extRsDampeningPolFromL3Outside(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsDampeningPol %v", err)

	} else {
		d.Set("relation_l3ext_rs_dampening_pol", l3extRsDampeningPolData)
	}

	l3extRsEctxData, err := aciClient.ReadRelationl3extRsEctxFromL3Outside(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsEctx %v", err)
		d.Set("relation_l3ext_rs_ectx", "")

	} else {
		if _, ok := d.GetOk("relation_l3ext_rs_ectx"); ok {
			tfName := GetMOName(d.Get("relation_l3ext_rs_ectx").(string))
			if tfName != l3extRsEctxData {
				d.Set("relation_l3ext_rs_ectx", "")
			}
		}
	}

	l3extRsOutToBDPublicSubnetHolderData, err := aciClient.ReadRelationl3extRsOutToBDPublicSubnetHolderFromL3Outside(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsOutToBDPublicSubnetHolder %v", err)
		d.Set("relation_l3ext_rs_out_to_bd_public_subnet_holder", make([]string, 0, 1))

	} else {
		if _, ok := d.GetOk("relation_l3ext_rs_out_to_bd_public_subnet_holder"); ok {
			relationParamList := toStringList(d.Get("relation_l3ext_rs_out_to_bd_public_subnet_holder").(*schema.Set).List())
			l3extRsOutToBDPublicSubnetHolderDataList := toStringList(l3extRsOutToBDPublicSubnetHolderData.(*schema.Set).List())
			sort.Strings(relationParamList)
			sort.Strings(l3extRsOutToBDPublicSubnetHolderDataList)
			if !reflect.DeepEqual(relationParamList, l3extRsOutToBDPublicSubnetHolderDataList) {
				d.Set("relation_l3ext_rs_out_to_bd_public_subnet_holder", make([]string, 0, 1))
			}
		}
	}

	l3extRsInterleakPolData, err := aciClient.ReadRelationl3extRsInterleakPolFromL3Outside(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsInterleakPol %v", err)
		d.Set("relation_l3ext_rs_interleak_pol", "")

	} else {
		if _, ok := d.GetOk("relation_l3ext_rs_interleak_pol"); ok {
			tfName := GetMOName(d.Get("relation_l3ext_rs_interleak_pol").(string))
			if tfName != l3extRsInterleakPolData {
				d.Set("relation_l3ext_rs_interleak_pol", "")
			}
		}
	}

	l3extRsL3DomAttData, err := aciClient.ReadRelationl3extRsL3DomAttFromL3Outside(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsL3DomAtt %v", err)
		d.Set("relation_l3ext_rs_l3_dom_att", "")

	} else {
		if _, ok := d.GetOk("relation_l3ext_rs_l3_dom_att"); ok {
			tfName := d.Get("relation_l3ext_rs_l3_dom_att").(string)
			if tfName != l3extRsL3DomAttData {
				d.Set("relation_l3ext_rs_l3_dom_att", "")
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciL3OutsideDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l3extOut")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
