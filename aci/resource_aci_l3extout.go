package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"enforce_rtctrl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
	d.SetId(l3extOut.DistinguishedName)
	d.Set("description", l3extOut.Description)
	d.Set("tenant_dn", GetParentDn(l3extOut.DistinguishedName))
	l3extOutMap, _ := l3extOut.ToMap()

	d.Set("name", l3extOutMap["name"])

	d.Set("annotation", l3extOutMap["annotation"])
	d.Set("enforce_rtctrl", l3extOutMap["enforceRtctrl"])
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
	}
	if EnforceRtctrl, ok := d.GetOk("enforce_rtctrl"); ok {
		l3extOutAttr.EnforceRtctrl = EnforceRtctrl.(string)
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

	d.SetPartial("name")

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
			d.SetPartial("relation_l3ext_rs_dampening_pol")
			d.Partial(false)
		}

	}
	if relationTol3extRsEctx, ok := d.GetOk("relation_l3ext_rs_ectx"); ok {
		relationParam := relationTol3extRsEctx.(string)
		err = aciClient.CreateRelationl3extRsEctxFromL3Outside(l3extOut.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_l3ext_rs_ectx")
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
			d.SetPartial("relation_l3ext_rs_out_to_bd_public_subnet_holder")
			d.Partial(false)
		}
	}
	if relationTol3extRsInterleakPol, ok := d.GetOk("relation_l3ext_rs_interleak_pol"); ok {
		relationParam := relationTol3extRsInterleakPol.(string)
		err = aciClient.CreateRelationl3extRsInterleakPolFromL3Outside(l3extOut.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_l3ext_rs_interleak_pol")
		d.Partial(false)

	}
	if relationTol3extRsL3DomAtt, ok := d.GetOk("relation_l3ext_rs_l3_dom_att"); ok {
		relationParam := relationTol3extRsL3DomAtt.(string)
		err = aciClient.CreateRelationl3extRsL3DomAttFromL3Outside(l3extOut.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_l3ext_rs_l3_dom_att")
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
	}
	if EnforceRtctrl, ok := d.GetOk("enforce_rtctrl"); ok {
		l3extOutAttr.EnforceRtctrl = EnforceRtctrl.(string)
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

	d.SetPartial("name")

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
			d.SetPartial("relation_l3ext_rs_dampening_pol")
			d.Partial(false)
		}

	}
	if d.HasChange("relation_l3ext_rs_ectx") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_ectx")
		err = aciClient.CreateRelationl3extRsEctxFromL3Outside(l3extOut.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_l3ext_rs_ectx")
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
			d.SetPartial("relation_l3ext_rs_out_to_bd_public_subnet_holder")
			d.Partial(false)

		}

	}
	if d.HasChange("relation_l3ext_rs_interleak_pol") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_interleak_pol")
		err = aciClient.DeleteRelationl3extRsInterleakPolFromL3Outside(l3extOut.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationl3extRsInterleakPolFromL3Outside(l3extOut.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_l3ext_rs_interleak_pol")
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
		d.SetPartial("relation_l3ext_rs_l3_dom_att")
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

	} else {
		d.Set("relation_l3ext_rs_ectx", l3extRsEctxData)
	}

	l3extRsOutToBDPublicSubnetHolderData, err := aciClient.ReadRelationl3extRsOutToBDPublicSubnetHolderFromL3Outside(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsOutToBDPublicSubnetHolder %v", err)

	} else {
		d.Set("relation_l3ext_rs_out_to_bd_public_subnet_holder", l3extRsOutToBDPublicSubnetHolderData)
	}

	l3extRsInterleakPolData, err := aciClient.ReadRelationl3extRsInterleakPolFromL3Outside(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsInterleakPol %v", err)

	} else {
		d.Set("relation_l3ext_rs_interleak_pol", l3extRsInterleakPolData)
	}

	l3extRsL3DomAttData, err := aciClient.ReadRelationl3extRsL3DomAttFromL3Outside(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsL3DomAtt %v", err)

	} else {
		d.Set("relation_l3ext_rs_l3_dom_att", l3extRsL3DomAttData)
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
