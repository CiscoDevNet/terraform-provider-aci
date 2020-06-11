package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciL3DomainProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciL3DomainProfileCreate,
		Update: resourceAciL3DomainProfileUpdate,
		Read:   resourceAciL3DomainProfileRead,
		Delete: resourceAciL3DomainProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciL3DomainProfileImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_infra_rs_vlan_ns": &schema.Schema{
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
			"relation_extnw_rs_out": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_infra_rs_dom_vxlan_ns_def": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteL3DomainProfile(client *client.Client, dn string) (*models.L3DomainProfile, error) {
	l3extDomPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extDomP := models.L3DomainProfileFromContainer(l3extDomPCont)

	if l3extDomP.DistinguishedName == "" {
		return nil, fmt.Errorf("L3DomainProfile %s not found", l3extDomP.DistinguishedName)
	}

	return l3extDomP, nil
}

func setL3DomainProfileAttributes(l3extDomP *models.L3DomainProfile, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(l3extDomP.DistinguishedName)
	d.Set("description", l3extDomP.Description)
	l3extDomPMap, _ := l3extDomP.ToMap()

	d.Set("name", l3extDomPMap["name"])

	d.Set("annotation", l3extDomPMap["annotation"])
	d.Set("name_alias", l3extDomPMap["nameAlias"])
	return d
}

func resourceAciL3DomainProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	l3extDomP, err := getRemoteL3DomainProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setL3DomainProfileAttributes(l3extDomP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL3DomainProfileCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] L3DomainProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	l3extDomPAttr := models.L3DomainProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extDomPAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extDomPAttr.NameAlias = NameAlias.(string)
	}
	l3extDomP := models.NewL3DomainProfile(fmt.Sprintf("l3dom-%s", name), "uni", desc, l3extDomPAttr)

	err := aciClient.Save(l3extDomP)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if relationToinfraRsVlanNs, ok := d.GetOk("relation_infra_rs_vlan_ns"); ok {
		relationParam := relationToinfraRsVlanNs.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsVlanNsFromL3DomainProfile(l3extDomP.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_vlan_ns")
		d.Partial(false)

	}
	if relationToinfraRsVlanNsDef, ok := d.GetOk("relation_infra_rs_vlan_ns_def"); ok {
		relationParam := relationToinfraRsVlanNsDef.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsVlanNsDefFromL3DomainProfile(l3extDomP.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_vlan_ns_def")
		d.Partial(false)

	}
	if relationToinfraRsVipAddrNs, ok := d.GetOk("relation_infra_rs_vip_addr_ns"); ok {
		relationParam := relationToinfraRsVipAddrNs.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsVipAddrNsFromL3DomainProfile(l3extDomP.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_vip_addr_ns")
		d.Partial(false)

	}
	if relationToextnwRsOut, ok := d.GetOk("relation_extnw_rs_out"); ok {
		relationParamList := toStringList(relationToextnwRsOut.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationextnwRsOutFromL3DomainProfile(l3extDomP.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_extnw_rs_out")
			d.Partial(false)
		}
	}
	if relationToinfraRsDomVxlanNsDef, ok := d.GetOk("relation_infra_rs_dom_vxlan_ns_def"); ok {
		relationParam := relationToinfraRsDomVxlanNsDef.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsDomVxlanNsDefFromL3DomainProfile(l3extDomP.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_dom_vxlan_ns_def")
		d.Partial(false)

	}

	d.SetId(l3extDomP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciL3DomainProfileRead(d, m)
}

func resourceAciL3DomainProfileUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] L3DomainProfile: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	l3extDomPAttr := models.L3DomainProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extDomPAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extDomPAttr.NameAlias = NameAlias.(string)
	}
	l3extDomP := models.NewL3DomainProfile(fmt.Sprintf("l3dom-%s", name), "uni", desc, l3extDomPAttr)

	l3extDomP.Status = "modified"

	err := aciClient.Save(l3extDomP)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if d.HasChange("relation_infra_rs_vlan_ns") {
		_, newRelParam := d.GetChange("relation_infra_rs_vlan_ns")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationinfraRsVlanNsFromL3DomainProfile(l3extDomP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationinfraRsVlanNsFromL3DomainProfile(l3extDomP.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_vlan_ns")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_vlan_ns_def") {
		_, newRelParam := d.GetChange("relation_infra_rs_vlan_ns_def")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsVlanNsDefFromL3DomainProfile(l3extDomP.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_vlan_ns_def")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_vip_addr_ns") {
		_, newRelParam := d.GetChange("relation_infra_rs_vip_addr_ns")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationinfraRsVipAddrNsFromL3DomainProfile(l3extDomP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationinfraRsVipAddrNsFromL3DomainProfile(l3extDomP.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_vip_addr_ns")
		d.Partial(false)

	}
	if d.HasChange("relation_extnw_rs_out") {
		oldRel, newRel := d.GetChange("relation_extnw_rs_out")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationextnwRsOutFromL3DomainProfile(l3extDomP.DistinguishedName, relDn)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_extnw_rs_out")
			d.Partial(false)

		}

	}
	if d.HasChange("relation_infra_rs_dom_vxlan_ns_def") {
		_, newRelParam := d.GetChange("relation_infra_rs_dom_vxlan_ns_def")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsDomVxlanNsDefFromL3DomainProfile(l3extDomP.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_dom_vxlan_ns_def")
		d.Partial(false)

	}

	d.SetId(l3extDomP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciL3DomainProfileRead(d, m)

}

func resourceAciL3DomainProfileRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l3extDomP, err := getRemoteL3DomainProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setL3DomainProfileAttributes(l3extDomP, d)

	infraRsVlanNsData, err := aciClient.ReadRelationinfraRsVlanNsFromL3DomainProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsVlanNs %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_vlan_ns"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_vlan_ns").(string))
			if tfName != infraRsVlanNsData {
				d.Set("relation_infra_rs_vlan_ns", "")
			}
		}
	}

	infraRsVlanNsDefData, err := aciClient.ReadRelationinfraRsVlanNsDefFromL3DomainProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsVlanNsDef %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_vlan_ns_def"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_vlan_ns_def").(string))
			if tfName != infraRsVlanNsDefData {
				d.Set("relation_infra_rs_vlan_ns_def", "")
			}
		}
	}

	infraRsVipAddrNsData, err := aciClient.ReadRelationinfraRsVipAddrNsFromL3DomainProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsVipAddrNs %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_vip_addr_ns"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_vip_addr_ns").(string))
			if tfName != infraRsVipAddrNsData {
				d.Set("relation_infra_rs_vip_addr_ns", "")
			}
		}
	}

	extnwRsOutData, err := aciClient.ReadRelationextnwRsOutFromL3DomainProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation extnwRsOut %v", err)

	} else {
		d.Set("relation_extnw_rs_out", extnwRsOutData)
	}

	infraRsDomVxlanNsDefData, err := aciClient.ReadRelationinfraRsDomVxlanNsDefFromL3DomainProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsDomVxlanNsDef %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_dom_vxlan_ns_def"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_dom_vxlan_ns_def").(string))
			if tfName != infraRsDomVxlanNsDefData {
				d.Set("relation_infra_rs_dom_vxlan_ns_def", "")
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciL3DomainProfileDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l3extDomP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
